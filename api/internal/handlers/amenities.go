package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dw-souza/pcity/api/internal/apperrors"
	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/dw-souza/pcity/api/internal/middleware"
	"github.com/dw-souza/pcity/api/internal/models"
	"github.com/dw-souza/pcity/api/internal/services"
	"github.com/go-chi/chi/v5"
)

type AmenitiesHandler struct {
	amenities *services.AmenitiesService
}

func NewAmenitiesHandler(amenities *services.AmenitiesService) *AmenitiesHandler {
	return &AmenitiesHandler{amenities: amenities}
}

func (h *AmenitiesHandler) ListTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.amenities.ListTypes(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if types == nil {
		types = []models.AmenityType{}
	}
	httputil.WriteData(w, http.StatusOK, types)
}

type reportAmenityRequest struct {
	AmenityTypeID string `json:"amenity_type_id"`
}

func (h *AmenitiesHandler) Report(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "placeId")
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "Invalid or missing token")
		return
	}

	var body reportAmenityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.AmenityTypeID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "bad_request", "Invalid request body")
		return
	}

	pa, err := h.amenities.Report(r.Context(), placeID, body.AmenityTypeID, userID)
	if errors.Is(err, apperrors.ErrAmenityAlreadyExists) {
		writeAmenityExists(w, pa)
		return
	}
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteData(w, http.StatusCreated, pa)
}

func (h *AmenitiesHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "placeId")
	amenityTypeID := chi.URLParam(r, "amenityTypeId")
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "Invalid or missing token")
		return
	}

	pa, err := h.amenities.Confirm(r.Context(), placeID, amenityTypeID, userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteData(w, http.StatusOK, pa)
}
