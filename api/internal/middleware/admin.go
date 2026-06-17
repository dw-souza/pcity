package middleware

import (
	"net/http"

	"github.com/dw-souza/pcity/api/internal/httputil"
)

func RequireAdmin(adminKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if adminKey == "" || r.Header.Get("X-Admin-Key") != adminKey {
				httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "Invalid admin key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
