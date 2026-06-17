package googleplaces

import (
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
