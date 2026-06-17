package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dw-souza/pcity/api/internal/httputil"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func OptionalAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if userID, ok := parseBearer(r, jwtSecret); ok {
				r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userID))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := parseBearer(r, jwtSecret)
			if !ok {
				httputil.WriteError(w, http.StatusUnauthorized, "unauthorized", "Invalid or missing token")
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userID))
			next.ServeHTTP(w, r)
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(UserIDKey).(string)
	return v, ok
}

func parseBearer(r *http.Request, jwtSecret string) (string, bool) {
	if jwtSecret == "" {
		return "", false
	}
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return "", false
	}
	tokenStr := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return "", false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	sub, ok := claims["sub"].(string)
	return sub, ok && sub != ""
}
