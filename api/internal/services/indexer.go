package services

import (
	"context"
	"math"

	"github.com/dw-souza/pcity/api/internal/apperrors"
	"github.com/dw-souza/pcity/api/internal/googleplaces"
	"github.com/dw-souza/pcity/api/internal/models"
	"github.com/dw-souza/pcity/api/internal/repository"
)

type Indexer struct {
	repo   *repository.Repository
	places *googleplaces.Client
}

func NewIndexer(repo *repository.Repository, places *googleplaces.Client) *Indexer {
	return &Indexer{repo: repo, places: places}
}

const (
	gridStepDeg    = 0.04    // ~4.4 km at Franca's latitude
	cellRadiusM    = 5000.0  // circle per grid cell
	cityMaxRadiusM = 15000.0 // municipality envelope from city center
)

// Plain terms only — city name in textQuery overrides locationBias on Places API (New).
var indexQueries = []string{
	"bares",
	"restaurantes",
	"lanchonetes",
	"cafés",
	"pizzarias",
	"hamburguerias",
}

func (s *Indexer) IndexCity(ctx context.Context, slug string) (models.CityIndexResult, error) {
	city, claimed, err := s.repo.TryClaimIndexing(ctx, slug)
	if err != nil {
		return models.CityIndexResult{}, err
	}
	if !claimed {
		if city.IndexingStatus == models.IndexingInProgress {
			return models.CityIndexResult{}, apperrors.ErrIndexingInProgress
		}
		return models.CityIndexResult{
			City:           city,
			PlacesImported: 0,
			IndexingStatus: city.IndexingStatus,
		}, nil
	}

	imported, indexErr := s.fetchAndStore(ctx, city)
	status := models.IndexingCompleted
	if indexErr != nil {
		status = models.IndexingFailed
	}

	updated, err := s.repo.FinishIndexing(ctx, city.ID, status)
	if err != nil {
		return models.CityIndexResult{}, err
	}
	if indexErr != nil {
		return models.CityIndexResult{}, indexErr
	}

	return models.CityIndexResult{
		City:           updated,
		PlacesImported: imported,
		IndexingStatus: updated.IndexingStatus,
	}, nil
}

func (s *Indexer) fetchAndStore(ctx context.Context, city models.City) (int, error) {
	seen := make(map[string]struct{})
	var batch []models.ImportedPlace

	for _, cell := range searchGrid(city.Location) {
		area := &googleplaces.AreaSearch{
			Lat:     cell.Lat,
			Lng:     cell.Lng,
			RadiusM: cellRadiusM,
		}
		for _, query := range indexQueries {
			places, err := s.places.SearchPlaces(ctx, query, area)
			if err != nil {
				return 0, err
			}
			for _, p := range places {
				if !withinCityRadius(p, city.Location, cityMaxRadiusM) {
					continue
				}
				if _, ok := seen[p.GooglePlaceID]; ok {
					continue
				}
				seen[p.GooglePlaceID] = struct{}{}
				batch = append(batch, p)
			}
		}
	}

	return s.repo.UpsertImportedPlaces(ctx, city.ID, batch)
}

func searchGrid(center models.GeoPoint) []models.GeoPoint {
	step := gridStepDeg
	cells := make([]models.GeoPoint, 0, 9)
	for dLat := -step; dLat <= step+1e-9; dLat += step {
		for dLng := -step; dLng <= step+1e-9; dLng += step {
			cells = append(cells, models.GeoPoint{
				Lat: center.Lat + dLat,
				Lng: center.Lng + dLng,
			})
		}
	}
	return cells
}

func withinCityRadius(p models.ImportedPlace, center models.GeoPoint, maxM float64) bool {
	return haversineM(center.Lat, center.Lng, p.Lat, p.Lng) <= maxM
}

func haversineM(lat1, lng1, lat2, lng2 float64) float64 {
	const earthR = 6371000.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return earthR * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (s *Indexer) EnsureIndexed(ctx context.Context, slug string) (models.City, error) {
	city, err := s.repo.GetCityBySlug(ctx, slug)
	if err != nil {
		return city, err
	}
	switch city.IndexingStatus {
	case models.IndexingCompleted:
		return city, nil
	case models.IndexingInProgress:
		return city, apperrors.ErrIndexingInProgress
	case models.IndexingPending, models.IndexingFailed:
		result, err := s.IndexCity(ctx, slug)
		if err != nil {
			return city, err
		}
		return result.City, nil
	default:
		return city, nil
	}
}
