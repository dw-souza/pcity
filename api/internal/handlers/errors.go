package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dw-souza/pcity/api/internal/apperrors"
	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/dw-souza/pcity/api/internal/models"
)

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		httputil.WriteError(w, http.StatusNotFound, "not_found", "Resource not found")
	case errors.Is(err, apperrors.ErrInvalidInput):
		httputil.WriteError(w, http.StatusBadRequest, "bad_request", "Invalid request")
	case errors.Is(err, apperrors.ErrIndexingInProgress):
		httputil.WriteError(w, http.StatusConflict, "indexing_in_progress", "City indexing is in progress")
	case errors.Is(err, apperrors.ErrAlreadyConfirmed):
		httputil.WriteError(w, http.StatusConflict, "already_confirmed", "You have already confirmed this amenity")
	case errors.Is(err, apperrors.ErrCannotConfirmOwnReport):
		httputil.WriteError(w, http.StatusConflict, "cannot_confirm_own_report", "You cannot confirm your own report")
	default:
		httputil.WriteError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
	}
}

func writeIndexingInProgress(w http.ResponseWriter, city models.City) {
	cityJSON, _ := json.Marshal(city)
	httputil.WriteErrorWithFields(w, http.StatusConflict, httputil.ErrorBody{
		Code:    "indexing_in_progress",
		Message: "City indexing is in progress",
		City:    cityJSON,
	})
}

func writeAmenityExists(w http.ResponseWriter, existing models.PlaceAmenity) {
	existingJSON, _ := json.Marshal(existing)
	httputil.WriteErrorWithFields(w, http.StatusConflict, httputil.ErrorBody{
		Code:     "amenity_already_exists",
		Message:  "Amenity already exists for this place",
		Existing: existingJSON,
	})
}
