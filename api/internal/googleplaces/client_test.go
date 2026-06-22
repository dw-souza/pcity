package googleplaces

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dw-souza/pcity/api/internal/models"
)

func TestMapCategory(t *testing.T) {
	tests := []struct {
		types []string
		want  models.PlaceCategory
	}{
		{[]string{"bar", "point_of_interest"}, models.CategoryBar},
		{[]string{"restaurant", "food"}, models.CategoryRestaurant},
		{[]string{"cafe", "store"}, models.CategoryCafe},
		{[]string{"gym"}, models.CategoryOther},
	}

	for _, tt := range tests {
		if got := mapCategory(tt.types); got != tt.want {
			t.Errorf("mapCategory(%v) = %q, want %q", tt.types, got, tt.want)
		}
	}
}

func TestSearchPlaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.Header.Get("X-Goog-Api-Key") != "test-key" {
			t.Errorf("missing or wrong API key header")
		}
		if !strings.Contains(r.Header.Get("X-Goog-FieldMask"), "places.id") {
			t.Errorf("missing field mask")
		}

		var req searchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.TextQuery != "bares em Franca, SP, Brasil" {
			t.Errorf("textQuery = %q", req.TextQuery)
		}
		if req.LanguageCode != "pt-BR" {
			t.Errorf("languageCode = %q", req.LanguageCode)
		}
		if req.PageSize != pageSize {
			t.Errorf("pageSize = %d, want %d", req.PageSize, pageSize)
		}

		resp := searchResponse{
			Places: []placeResult{
				{
					ID: "ChIJ123",
					DisplayName: struct {
						Text string `json:"text"`
					}{Text: "Bar do Zé"},
					FormattedAddress: "Rua A, Franca - SP",
					Location: struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					}{Latitude: -20.53, Longitude: -47.40},
					Types: []string{"bar", "point_of_interest"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client := NewClient("test-key")
	client.baseURL = srv.URL

	places, err := client.SearchPlaces(context.Background(), "bares em Franca, SP, Brasil", nil)
	if err != nil {
		t.Fatalf("SearchPlaces: %v", err)
	}
	if len(places) != 1 {
		t.Fatalf("len(places) = %d, want 1", len(places))
	}
	if places[0].GooglePlaceID != "ChIJ123" {
		t.Errorf("GooglePlaceID = %q", places[0].GooglePlaceID)
	}
	if places[0].Name != "Bar do Zé" {
		t.Errorf("Name = %q", places[0].Name)
	}
	if places[0].Category != models.CategoryBar {
		t.Errorf("Category = %q", places[0].Category)
	}
}

func TestSearchPlacesAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(apiErrorResponse{
			Error: struct {
				Message string `json:"message"`
				Status  string `json:"status"`
			}{
				Status:  "PERMISSION_DENIED",
				Message: "API key not valid",
			},
		})
	}))
	defer srv.Close()

	client := NewClient("bad-key")
	client.baseURL = srv.URL

	_, err := client.SearchPlaces(context.Background(), "bares em Franca, SP, Brasil", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "PERMISSION_DENIED") {
		t.Errorf("error = %v", err)
	}
}

func TestSearchPlacesWithArea(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req searchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatal(err)
		}
		if req.LocationBias == nil {
			t.Fatal("expected locationBias")
		}
		if req.LocationBias.Circle.Radius != 5000 {
			t.Errorf("radius = %v", req.LocationBias.Circle.Radius)
		}
		json.NewEncoder(w).Encode(searchResponse{})
	}))
	defer srv.Close()

	client := NewClient("test-key")
	client.baseURL = srv.URL

	_, err := client.SearchPlaces(context.Background(), "bares em Franca, SP", &AreaSearch{
		Lat: -20.53, Lng: -47.40, RadiusM: 5000,
	})
	if err != nil {
		t.Fatal(err)
	}
}
