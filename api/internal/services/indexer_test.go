package services

import (
	"testing"

	"github.com/dw-souza/pcity/api/internal/models"
)

func TestSearchGrid(t *testing.T) {
	center := models.GeoPoint{Lat: -20.5386, Lng: -47.4008}
	cells := searchGrid(center)
	if len(cells) != 9 {
		t.Fatalf("len(cells) = %d, want 9", len(cells))
	}
	if cells[4].Lat != center.Lat || cells[4].Lng != center.Lng {
		t.Errorf("center cell = %+v, want center %+v", cells[4], center)
	}
}

func TestWithinCityRadius(t *testing.T) {
	center := models.GeoPoint{Lat: -20.5386, Lng: -47.4008}
	near := models.ImportedPlace{Lat: -20.54, Lng: -47.40}
	far := models.ImportedPlace{Lat: -20.70, Lng: -47.40}

	if !withinCityRadius(near, center, cityMaxRadiusM) {
		t.Error("near place should be within city radius")
	}
	if withinCityRadius(far, center, cityMaxRadiusM) {
		t.Error("far place should be outside city radius")
	}
}
