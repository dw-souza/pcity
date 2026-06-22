package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/dw-souza/pcity/api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type DevAuthHandler struct {
	repo      *repository.Repository
	jwtSecret string
}

func NewDevAuthHandler(repo *repository.Repository, jwtSecret string) *DevAuthHandler {
	return &DevAuthHandler{repo: repo, jwtSecret: jwtSecret}
}

type devAuthRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type devAuthResponse struct {
	AccessToken string `json:"access_token"`
	User        struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	} `json:"user"`
}

func (h *DevAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if h.jwtSecret == "" {
		httputil.WriteError(w, http.StatusServiceUnavailable, "dev_auth_disabled", "JWT secret not configured")
		return
	}

	var body devAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "bad_request", "Invalid request body")
		return
	}

	email := strings.TrimSpace(strings.ToLower(body.Email))
	displayName := strings.TrimSpace(body.DisplayName)
	if email == "" || displayName == "" {
		httputil.WriteError(w, http.StatusBadRequest, "bad_request", "email and display_name are required")
		return
	}

	userID, err := h.repo.UpsertDevUser(r.Context(), email, displayName)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"role": "authenticated",
		"exp":  time.Now().Add(30 * 24 * time.Hour).Unix(),
	})
	accessToken, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
		return
	}

	var resp devAuthResponse
	resp.AccessToken = accessToken
	resp.User.ID = userID
	resp.User.Email = email
	resp.User.DisplayName = displayName
	httputil.WriteData(w, http.StatusOK, resp)
}
