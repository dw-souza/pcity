package handlers

import (
	"errors"
	"net/http"

	"github.com/dw-souza/pcity/api/internal/apperrors"
	"github.com/dw-souza/pcity/api/internal/config"
	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/dw-souza/pcity/api/internal/middleware"
	"github.com/dw-souza/pcity/api/internal/models"
	"github.com/dw-souza/pcity/api/internal/services"
	"github.com/go-chi/chi/v5"
)

type CitiesHandler struct {
	places *services.PlacesService
}

func NewCitiesHandler(places *services.PlacesService) *CitiesHandler {
	return &CitiesHandler{places: places}
}

func (h *CitiesHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	city, err := h.places.GetCity(r.Context(), slug)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteData(w, http.StatusOK, city)
}

func (h *CitiesHandler) ListPlaces(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	q := r.URL.Query()

	limit := config.ParseIntQuery(q.Get("limit"), 50)
	if limit > 100 {
		limit = 100
	}
	offset := config.ParseIntQuery(q.Get("offset"), 0)

	var lat, lng *float64
	if v, ok := config.ParseFloatQuery(q.Get("lat")); ok {
		lat = &v
	}
	if v, ok := config.ParseFloatQuery(q.Get("lng")); ok {
		lng = &v
	}

	var category *models.PlaceCategory
	if c := q.Get("category"); c != "" {
		cat := models.PlaceCategory(c)
		category = &cat
	}

	userID, _ := middleware.UserIDFromContext(r.Context())

	resp, err := h.places.ListCityPlaces(r.Context(), slug, models.ListPlacesParams{
		Lat:      lat,
		Lng:      lng,
		Category: category,
		Search:   q.Get("q"),
		Limit:    limit,
		Offset:   offset,
		UserID:   userID,
	})
	if errors.Is(err, apperrors.ErrIndexingInProgress) {
		writeIndexingInProgress(w, resp.City)
		return
	}
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, resp)
}
