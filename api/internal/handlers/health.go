package handlers

import (
	"net/http"

	"github.com/dw-souza/pcity/api/internal/httputil"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	httputil.WriteData(w, http.StatusOK, map[string]string{"status": "ok"})
}
