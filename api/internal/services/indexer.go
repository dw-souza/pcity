package services

import (
	"context"
	"fmt"

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

var indexQueries = []string{
	"bares em %s, %s, Brasil",
	"restaurantes em %s, %s, Brasil",
	"lanchonetes em %s, %s, Brasil",
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

	for _, tmpl := range indexQueries {
		query := fmt.Sprintf(tmpl, city.Name, city.State)
		places, err := s.places.SearchPlaces(ctx, query)
		if err != nil {
			return 0, err
		}
		for _, p := range places {
			if _, ok := seen[p.GooglePlaceID]; ok {
				continue
			}
			seen[p.GooglePlaceID] = struct{}{}
			batch = append(batch, p)
		}
	}

	return s.repo.UpsertImportedPlaces(ctx, city.ID, batch)
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
