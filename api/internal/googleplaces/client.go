package googleplaces

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dw-souza/pcity/api/internal/models"
)

const (
	textSearchURL = "https://places.googleapis.com/v1/places:searchText"
	fieldMask     = "places.id,places.displayName,places.formattedAddress,places.location,places.types,nextPageToken"
	pageSize      = 20
	maxPages      = 3
)

type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string // override in tests
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: textSearchURL,
	}
}

type searchRequest struct {
	TextQuery    string         `json:"textQuery"`
	LanguageCode string         `json:"languageCode"`
	PageSize     int            `json:"pageSize"`
	PageToken    string         `json:"pageToken,omitempty"`
	LocationBias *locationBias  `json:"locationBias,omitempty"`
}

type locationBias struct {
	Circle circleRestriction `json:"circle"`
}

type circleRestriction struct {
	Center struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"center"`
	Radius float64 `json:"radius"`
}

// AreaSearch biases results toward a circular region (radius in meters).
type AreaSearch struct {
	Lat     float64
	Lng     float64
	RadiusM float64
}

type searchResponse struct {
	Places        []placeResult `json:"places"`
	NextPageToken string        `json:"nextPageToken"`
}

type placeResult struct {
	ID               string `json:"id"`
	DisplayName      struct {
		Text string `json:"text"`
	} `json:"displayName"`
	FormattedAddress string `json:"formattedAddress"`
	Location         struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Types []string `json:"types"`
}

type apiErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (c *Client) SearchPlaces(ctx context.Context, query string, area *AreaSearch) ([]models.ImportedPlace, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("google places api key not configured")
	}

	var all []models.ImportedPlace
	pageToken := ""
	for page := 0; page < maxPages; page++ {
		results, next, err := c.searchPage(ctx, query, pageToken, area)
		if err != nil {
			return all, err
		}
		all = append(all, results...)
		if next == "" {
			break
		}
		pageToken = next
	}
	return all, nil
}

func (c *Client) searchPage(ctx context.Context, query, pageToken string, area *AreaSearch) ([]models.ImportedPlace, string, error) {
	reqBody := searchRequest{
		TextQuery:    query,
		LanguageCode: "pt-BR",
		PageSize:     pageSize,
		PageToken:    pageToken,
	}
	if area != nil {
		bias := &locationBias{}
		bias.Circle.Center.Latitude = area.Lat
		bias.Circle.Center.Longitude = area.Lng
		bias.Circle.Radius = area.RadiusM
		reqBody.LocationBias = bias
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", fieldMask)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", parseAPIError(resp.StatusCode, raw)
	}

	var parsed searchResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, "", err
	}

	out := make([]models.ImportedPlace, 0, len(parsed.Places))
	for _, p := range parsed.Places {
		out = append(out, models.ImportedPlace{
			GooglePlaceID: p.ID,
			Name:          p.DisplayName.Text,
			Address:       p.FormattedAddress,
			Lat:           p.Location.Latitude,
			Lng:           p.Location.Longitude,
			Category:      mapCategory(p.Types),
		})
	}
	return out, parsed.NextPageToken, nil
}

func parseAPIError(statusCode int, raw []byte) error {
	var apiErr apiErrorResponse
	if json.Unmarshal(raw, &apiErr) == nil && apiErr.Error.Message != "" {
		if apiErr.Error.Status != "" {
			return fmt.Errorf("google places api: %s: %s", apiErr.Error.Status, apiErr.Error.Message)
		}
		return fmt.Errorf("google places api: %s", apiErr.Error.Message)
	}
	return fmt.Errorf("google places api: HTTP %d", statusCode)
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
