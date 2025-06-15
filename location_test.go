package maze

import (
	"reflect"
	"testing"
)

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
				Latitude:  37.786971,
				Longitude: -122.399677,
				Altitude:  -123.456,
				Name:      "Test Location",
				Locality:  "San Francisco",
			},
			uri: "geo:37.786971,-122.399677,-123.456?locality=San+Francisco&name=Test+Location",
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
