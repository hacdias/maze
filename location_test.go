package maze

import (
	"math"
	"reflect"
	"testing"
)

func TestDistance(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		l1, l2   *Location
		distance float64
	}
	tests := []testCase{
		{
			name:     "Same point",
			l1:       &Location{Latitude: 38.7, Longitude: -9.1},
			l2:       &Location{Latitude: 38.7, Longitude: -9.1},
			distance: 0,
		},
		{
			name:     "Quarter of the equator",
			l1:       &Location{Latitude: 0, Longitude: 0},
			l2:       &Location{Latitude: 0, Longitude: 90},
			distance: math.Pi / 2 * 6371000,
		},
		{
			name:     "Opposite points",
			l1:       &Location{Latitude: -42.15305933728451, Longitude: -119.19357816989165},
			l2:       &Location{Latitude: 42.15305933728451, Longitude: 60.80642183010835},
			distance: math.Pi * 6371000,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := tc.l1.Distance(tc.l2)
			if math.IsNaN(d) || math.Abs(d-tc.distance) > 1 {
				t.Errorf("expected %v, got %v", tc.distance, d)
			}
		})
	}
}

func TestLocation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		location *Location
		uri      string
	}
	tests := []testCase{
		{
			name: "Simple",
			location: &Location{
				Latitude:  37.786971,
				Longitude: -122.399677,
			},
			uri: "geo:37.786971,-122.399677",
		},
		{
			name: "With Altitude",
			location: &Location{
				Latitude:  37.786971,
				Longitude: -122.399677,
				Altitude:  -123.456,
			},
			uri: "geo:37.786971,-122.399677,-123.456",
		},
		{
			name: "With Custom Fields",
			location: &Location{
				Latitude:    37.786971,
				Longitude:   -122.399677,
				Altitude:    -123.456,
				Name:        "Test Location",
				Locality:    "San Francisco",
				Region:      "California",
				Country:     "United States",
				CountryCode: "US",
				PostalCode:  "94103",
				ICAO:        "KSFO",
				IATA:        "SFO",
			},
			uri: "geo:37.786971,-122.399677,-123.456?country=United+States&countryCode=US&iata=SFO&icao=KSFO&locality=San+Francisco&name=Test+Location&postalCode=94103&region=California",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			loc, err := ParseLocation(tc.uri)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if reflect.DeepEqual(loc, tc.location) == false {
				t.Errorf("expected %v, got %v", tc.location, loc)
			}

			out := loc.String()
			if out != tc.uri {
				t.Errorf("expected %s, got %s", tc.uri, out)
			}
		})
	}
}
