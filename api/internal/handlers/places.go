package handlers

import (
	"net/http"

	"github.com/dw-souza/pcity/api/internal/config"
	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/dw-souza/pcity/api/internal/middleware"
	"github.com/dw-souza/pcity/api/internal/services"
	"github.com/go-chi/chi/v5"
)

type PlacesHandler struct {
	places *services.PlacesService
}

func NewPlacesHandler(places *services.PlacesService) *PlacesHandler {
	return &PlacesHandler{places: places}
}

func (h *PlacesHandler) Get(w http.ResponseWriter, r *http.Request) {
	placeID := chi.URLParam(r, "placeId")
	q := r.URL.Query()

	var lat, lng *float64
	if v, ok := config.ParseFloatQuery(q.Get("lat")); ok {
		lat = &v
	}
	if v, ok := config.ParseFloatQuery(q.Get("lng")); ok {
		lng = &v
	}

	userID, _ := middleware.UserIDFromContext(r.Context())

	place, err := h.places.GetPlace(r.Context(), placeID, lat, lng, userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteData(w, http.StatusOK, place)
}
