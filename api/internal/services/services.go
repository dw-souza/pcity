package services

import (
	"context"

	"github.com/dw-souza/pcity/api/internal/apperrors"
	"github.com/dw-souza/pcity/api/internal/models"
	"github.com/dw-souza/pcity/api/internal/repository"
)

type PlacesService struct {
	repo    *repository.Repository
	indexer *Indexer
}

func NewPlacesService(repo *repository.Repository, indexer *Indexer) *PlacesService {
	return &PlacesService{repo: repo, indexer: indexer}
}

func (s *PlacesService) GetCity(ctx context.Context, slug string) (models.City, error) {
	return s.repo.GetCityBySlug(ctx, slug)
}

func (s *PlacesService) ListCityPlaces(ctx context.Context, slug string, p models.ListPlacesParams) (models.PlaceListResponse, error) {
	city, err := s.indexer.EnsureIndexed(ctx, slug)
	if err != nil {
		if err == apperrors.ErrIndexingInProgress {
			return models.PlaceListResponse{City: city}, err
		}
		return models.PlaceListResponse{}, err
	}

	p.CityID = city.ID
	total, err := s.repo.CountPlaces(ctx, city.ID, p.Category, p.Search)
	if err != nil {
		return models.PlaceListResponse{}, err
	}

	places, err := s.repo.ListPlaces(ctx, p)
	if err != nil {
		return models.PlaceListResponse{}, err
	}
	if places == nil {
		places = []models.PlaceSummary{}
	}

	return models.PlaceListResponse{
		Data: places,
		Meta: models.PaginationMeta{Limit: p.Limit, Offset: p.Offset, Total: total},
		City: city,
	}, nil
}

func (s *PlacesService) GetPlace(ctx context.Context, placeID string, lat, lng *float64, userID string) (models.PlaceDetail, error) {
	return s.repo.GetPlace(ctx, placeID, lat, lng, userID)
}

type AmenitiesService struct {
	repo *repository.Repository
}

func NewAmenitiesService(repo *repository.Repository) *AmenitiesService {
	return &AmenitiesService{repo: repo}
}

func (s *AmenitiesService) ListTypes(ctx context.Context) ([]models.AmenityType, error) {
	return s.repo.ListAmenityTypes(ctx)
}

func (s *AmenitiesService) Report(ctx context.Context, placeID, amenityTypeID, userID string) (models.PlaceAmenity, error) {
	exists, err := s.repo.PlaceExists(ctx, placeID)
	if err != nil {
		return models.PlaceAmenity{}, err
	}
	if !exists {
		return models.PlaceAmenity{}, apperrors.ErrNotFound
	}

	ok, err := s.repo.AmenityTypeExists(ctx, amenityTypeID)
	if err != nil {
		return models.PlaceAmenity{}, err
	}
	if !ok {
		return models.PlaceAmenity{}, apperrors.ErrInvalidInput
	}

	return s.repo.CreatePlaceAmenity(ctx, placeID, amenityTypeID, userID)
}

func (s *AmenitiesService) Confirm(ctx context.Context, placeID, amenityTypeID, userID string) (models.PlaceAmenity, error) {
	exists, err := s.repo.PlaceExists(ctx, placeID)
	if err != nil {
		return models.PlaceAmenity{}, err
	}
	if !exists {
		return models.PlaceAmenity{}, apperrors.ErrNotFound
	}

	return s.repo.ConfirmPlaceAmenity(ctx, placeID, amenityTypeID, userID)
}

type ProfileService struct {
	repo *repository.Repository
}

func NewProfileService(repo *repository.Repository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) Get(ctx context.Context, userID string) (models.Profile, error) {
	return s.repo.GetProfile(ctx, userID)
}

func (s *ProfileService) Update(ctx context.Context, userID string, displayName *string) (models.Profile, error) {
	if displayName != nil && *displayName == "" {
		return models.Profile{}, apperrors.ErrInvalidInput
	}
	return s.repo.UpdateProfile(ctx, userID, displayName)
}
