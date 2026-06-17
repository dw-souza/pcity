package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/dw-souza/pcity/api/internal/middleware"
	"github.com/dw-souza/pcity/api/internal/services"
)

type ProfileHandler struct {
	profiles *services.ProfileService
}

func NewProfileHandler(profiles *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profiles: profiles}
}

func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "Invalid or missing token")
		return
	}

	profile, err := h.profiles.Get(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteData(w, http.StatusOK, profile)
}

type updateProfileRequest struct {
	DisplayName *string `json:"display_name"`
}

func (h *ProfileHandler) Patch(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "Invalid or missing token")
		return
	}

	var body updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "bad_request", "Invalid request body")
		return
	}

	profile, err := h.profiles.Update(r.Context(), userID, body.DisplayName)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteData(w, http.StatusOK, profile)
}
