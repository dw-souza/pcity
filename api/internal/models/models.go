package models

import "time"

type IndexingStatus string

const (
	IndexingPending    IndexingStatus = "pending"
	IndexingInProgress IndexingStatus = "in_progress"
	IndexingCompleted  IndexingStatus = "completed"
	IndexingFailed     IndexingStatus = "failed"
)

type PlaceSource string

const (
	PlaceSourceGoogle PlaceSource = "google"
	PlaceSourceUser   PlaceSource = "user"
)

type PlaceCategory string

const (
	CategoryBar        PlaceCategory = "bar"
	CategoryRestaurant PlaceCategory = "restaurant"
	CategoryCafe       PlaceCategory = "cafe"
	CategoryOther      PlaceCategory = "other"
)

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type City struct {
	ID              string         `json:"id"`
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	State           string         `json:"state"`
	Country         string         `json:"country"`
	Location        GeoPoint       `json:"location"`
	IndexingStatus  IndexingStatus `json:"indexing_status"`
	IndexedAt       *time.Time     `json:"indexed_at,omitempty"`
}

type AmenityType struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}

type PlaceAmenity struct {
	ID                string       `json:"id"`
	AmenityTypeID     string       `json:"amenity_type_id"`
	AmenityType       *AmenityType `json:"amenity_type,omitempty"`
	ConfirmationCount int          `json:"confirmation_count"`
	IsVerified        bool         `json:"is_verified"`
	FirstReportedAt   time.Time    `json:"first_reported_at"`
	VerifiedAt        *time.Time   `json:"verified_at,omitempty"`
	UserHasConfirmed  *bool        `json:"user_has_confirmed,omitempty"`
}

type PlaceSummary struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Address   string         `json:"address"`
	Location  GeoPoint       `json:"location"`
	Phone     *string        `json:"phone,omitempty"`
	Category  PlaceCategory  `json:"category"`
	Source    PlaceSource    `json:"source"`
	DistanceM *float64       `json:"distance_m,omitempty"`
	Amenities []PlaceAmenity `json:"amenities,omitempty"`
}

type PlaceDetail struct {
	PlaceSummary
	Website   *string   `json:"website,omitempty"`
	Amenities []PlaceAmenity `json:"amenities"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginationMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type PlaceListResponse struct {
	Data []PlaceSummary `json:"data"`
	Meta PaginationMeta `json:"meta"`
	City City           `json:"city"`
}

type Profile struct {
	ID          string     `json:"id"`
	DisplayName *string    `json:"display_name,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CityIndexResult struct {
	City            City           `json:"city"`
	PlacesImported  int            `json:"places_imported"`
	IndexingStatus  IndexingStatus `json:"indexing_status"`
}

type ImportedPlace struct {
	GooglePlaceID string
	Name          string
	Address       string
	Lat           float64
	Lng           float64
	Phone         *string
	Website       *string
	Category      PlaceCategory
}

type ListPlacesParams struct {
	CityID   string
	Lat      *float64
	Lng      *float64
	Category *PlaceCategory
	Search   string
	Limit    int
	Offset   int
	UserID   string
}
