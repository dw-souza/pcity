package googleplaces

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dw-souza/pcity/api/internal/models"
)

const textSearchURL = "https://maps.googleapis.com/maps/api/place/textsearch/json"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type textSearchResponse struct {
	Results []struct {
		PlaceID          string `json:"place_id"`
		Name             string `json:"name"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
		Types []string `json:"types"`
	} `json:"results"`
	NextPageToken string `json:"next_page_token"`
	Status        string `json:"status"`
}

func (c *Client) SearchPlaces(ctx context.Context, query string) ([]models.ImportedPlace, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("google places api key not configured")
	}

	var all []models.ImportedPlace
	pageToken := ""
	for page := 0; page < 3; page++ {
		results, next, err := c.searchPage(ctx, query, pageToken)
		if err != nil {
			return all, err
		}
		all = append(all, results...)
		if next == "" {
			break
		}
		pageToken = next
		time.Sleep(2 * time.Second) // Google requires delay before next_page_token
	}
	return all, nil
}

func (c *Client) searchPage(ctx context.Context, query, pageToken string) ([]models.ImportedPlace, string, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("key", c.apiKey)
	params.Set("language", "pt-BR")
	if pageToken != "" {
		params.Set("pagetoken", pageToken)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, textSearchURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var body textSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, "", err
	}
	if body.Status != "OK" && body.Status != "ZERO_RESULTS" {
		return nil, "", fmt.Errorf("google places api: %s", body.Status)
	}

	out := make([]models.ImportedPlace, 0, len(body.Results))
	for _, r := range body.Results {
		out = append(out, models.ImportedPlace{
			GooglePlaceID: r.PlaceID,
			Name:          r.Name,
			Address:       r.FormattedAddress,
			Lat:           r.Geometry.Location.Lat,
			Lng:           r.Geometry.Location.Lng,
			Category:      mapCategory(r.Types),
		})
	}
	return out, body.NextPageToken, nil
}

func mapCategory(types []string) models.PlaceCategory {
	joined := strings.ToLower(strings.Join(types, " "))
	switch {
	case strings.Contains(joined, "bar"):
		return models.CategoryBar
	case strings.Contains(joined, "cafe") || strings.Contains(joined, "bakery"):
		return models.CategoryCafe
	case strings.Contains(joined, "restaurant") || strings.Contains(joined, "food"):
		return models.CategoryRestaurant
	default:
		return models.CategoryOther
	}
}
