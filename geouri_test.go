package maze

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParseGeoURI(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			name    string
			input   string
			output  *GeoURI
			encoded string
		}
		tests := []testCase{
			{
				name:    "Simple",
				input:   "geo:37.786971,-122.399677",
				encoded: "geo:37.786971,-122.399677",
				output: &GeoURI{
					Latitude:   37.786971,
					Longitude:  -122.399677,
					Parameters: map[string][]string{},
					Query:      url.Values{},
				},
			},
			{
				name:    "Altitude",
				input:   "geo:37.786971,-122.399677,-123.456",
				encoded: "geo:37.786971,-122.399677,-123.456",
				output: &GeoURI{
					Latitude:   37.786971,
					Longitude:  -122.399677,
					Altitude:   -123.456,
					Parameters: map[string][]string{},
					Query:      url.Values{},
				},
			},
			{
				name:    "Indigenous",
				input:   "geo:51.5258325,-0.1359825,0.0;name=london;url=https://hwclondon.co.uk",
				encoded: "geo:51.5258325,-0.1359825;name=london;url=https%3A%2F%2Fhwclondon.co.uk",
				output: &GeoURI{
					Latitude:  51.5258325,
					Longitude: -0.1359825,
					Altitude:  0.0,
					Parameters: map[string][]string{
						"name": {"london"},
						"url":  {"https://hwclondon.co.uk"},
					},
					Query: url.Values{},
				},
			},
			{
				name:    "Indigenous with encoded URL",
				input:   "geo:51.525832,-0.135983;name=london;url=https%3A%2F%2Fhwclondon.co.uk",
				encoded: "geo:51.525832,-0.135983;name=london;url=https%3A%2F%2Fhwclondon.co.uk",
				output: &GeoURI{
					Latitude:  51.525832,
					Longitude: -0.135983,
					Altitude:  0.0,
					Parameters: map[string][]string{
						"name": {"london"},
						"url":  {"https://hwclondon.co.uk"},
					},
					Query: url.Values{},
				},
			},
			{
				name:    "Parameter without value",
				input:   "geo:51.5258325,-0.1359825,0.0;name",
				encoded: "geo:51.5258325,-0.1359825;name",
				output: &GeoURI{
					Latitude:  51.5258325,
					Longitude: -0.1359825,
					Altitude:  0.0,
					Parameters: map[string][]string{
						"name": {},
					},
					Query: url.Values{},
				},
			},
			{
				name:    "Multiple parameters without value (order)",
				input:   "geo:51.5258325,-0.1359825,0.0;name;location",
				encoded: "geo:51.5258325,-0.1359825;location;name",
				output: &GeoURI{
					Latitude:  51.5258325,
					Longitude: -0.1359825,
					Altitude:  0.0,
					Parameters: map[string][]string{
						"name":     {},
						"location": {},
					},
					Query: url.Values{},
				},
			},
			{
				name:    "Multiple parameters with same name",
				input:   "geo:51.5258325,-0.1359825,0.0;name=a;name=b;name=c",
				encoded: "geo:51.5258325,-0.1359825;name=a;name=b;name=c",
				output: &GeoURI{
					Latitude:  51.5258325,
					Longitude: -0.1359825,
					Altitude:  0.0,
					Parameters: map[string][]string{
						"name": {"a", "b", "c"},
					},
					Query: url.Values{},
				},
			},
			{
				name:    "Query parameters",
				input:   "geo:51.5258325,-0.1359825,0.0;name=a;name=b;name=c?a=b&c=d",
				encoded: "geo:51.5258325,-0.1359825;name=a;name=b;name=c?a=b&c=d",
				output: &GeoURI{
					Latitude:  51.5258325,
					Longitude: -0.1359825,
					Altitude:  0.0,
					Parameters: map[string][]string{
						"name": {"a", "b", "c"},
					},
					Query: url.Values{
						"a": {"b"},
						"c": {"d"},
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseGeoURI(tt.input)

				if err != nil {
					t.Errorf("Parse() error = %v", err)
				}
				if !reflect.DeepEqual(got, tt.output) {
					t.Errorf("Parse() = %+v, want %+v", got, tt.output)
				}

				v := got.String()
				if v != tt.encoded {
					t.Errorf("String() = %v, want %v", v, tt.encoded)
				}
			})
		}

	})

	t.Run("Failures", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			name  string
			input string
		}
		tests := []testCase{
			// Negative
			{
				name:  "Missing scheme",
				input: "37.786971,-122.399677",
			},
			{
				name:  "Missing path",
				input: "geo:",
			},
			{
				name:  "1 coordinate",
				input: "geo:37.786971",
			},
			{
				name:  "4 coordinates",
				input: "geo:123,123,123,123",
			},
			{
				name:  "Malformed latitude",
				input: "geo:12x,123",
			},
			{
				name:  "Malformed longitude",
				input: "geo:123,12x",
			},
			{
				name:  "Malformed altitude",
				input: "geo:123,123,12x",
			},
			{
				name:  "Malformed parameter",
				input: "geo:123,123;a=a=a",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := ParseGeoURI(tt.input)
				if err == nil {
					t.Errorf("Parse() expected error, got nil")
					return
				}
			})
		}
	})

}
