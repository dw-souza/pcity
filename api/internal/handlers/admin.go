package handlers

import (
	"errors"
	"net/http"

	"github.com/dw-souza/pcity/api/internal/apperrors"
	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/dw-souza/pcity/api/internal/services"
	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	indexer *services.Indexer
}

func NewAdminHandler(indexer *services.Indexer) *AdminHandler {
	return &AdminHandler{indexer: indexer}
}

func (h *AdminHandler) IndexCity(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	result, err := h.indexer.IndexCity(r.Context(), slug)
	if errors.Is(err, apperrors.ErrNotFound) {
		writeServiceError(w, err)
		return
	}
	if errors.Is(err, apperrors.ErrIndexingInProgress) {
		writeIndexingInProgress(w, result.City)
		return
	}
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httputil.WriteData(w, http.StatusOK, result)
}
