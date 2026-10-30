package maze

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// fakeTransport answers every request with status and body, without using the
// network.
type fakeTransport struct {
	status int
	body   string
}

func (f *fakeTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: f.status,
		Body:       io.NopCloser(strings.NewReader(f.body)),
	}, nil
}

type lookupTestCase struct {
	name     string
	status   int // 200 when zero
	body     string
	location *Location // nil when an error is expected
}

func runLookupTests(t *testing.T, tests []lookupTestCase, lookup func(*Maze) (*Location, error)) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			status := tc.status
			if status == 0 {
				status = http.StatusOK
			}

			loc, err := lookup(NewMaze(&http.Client{Transport: &fakeTransport{status: status, body: tc.body}}))
			if tc.location == nil {
				if err == nil {
					t.Errorf("expected error, got %+v", *loc)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(loc, tc.location) {
				t.Errorf("expected %+v, got %+v", *tc.location, *loc)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	t.Parallel()

	runLookupTests(t, []lookupTestCase{
		{
			name: "Full",
			body: `{"features":[{"properties":{"city":"Lisboa","state":"Lisboa","country":"Portugal","countrycode":"PT","postcode":"1000-001"}}]}`,
			location: &Location{
				Latitude:    38.7,
				Longitude:   -9.1,
				Locality:    "Lisboa",
				Region:      "Lisboa",
				Country:     "Portugal",
				CountryCode: "PT",
				PostalCode:  "1000-001",
			},
		},
		{
			name: "County as region",
			body: `{"features":[{"properties":{"county":"Cascais","country":"Portugal"}}]}`,
			location: &Location{
				Latitude:  38.7,
				Longitude: -9.1,
				Region:    "Cascais",
				Country:   "Portugal",
			},
		},
		{
			name: "No features",
			body: `{"features":[]}`,
		},
		{
			name: "No useful information",
			body: `{"features":[{"properties":{"name":"Atlantic Ocean"}}]}`,
		},
		{
			name:   "Server error",
			status: http.StatusServiceUnavailable,
			body:   `{"features":[{"properties":{"city":"Lisboa"}}]}`,
		},
	}, func(m *Maze) (*Location, error) {
		return m.Reverse("en", -9.1, 38.7)
	})
}

func TestSearch(t *testing.T) {
	t.Parallel()

	runLookupTests(t, []lookupTestCase{
		{
			name: "First feature",
			body: `{"features":[
				{"geometry":{"type":"Point","coordinates":[-9.1,38.7]},"properties":{"city":"Lisboa","country":"Portugal","countrycode":"PT","postcode":"1000-001"}},
				{"geometry":{"type":"Point","coordinates":[-8.6,41.1]},"properties":{"city":"Porto","country":"Portugal"}}
			]}`,
			location: &Location{
				Latitude:    38.7,
				Longitude:   -9.1,
				Name:        "lisbon",
				Locality:    "Lisboa",
				Country:     "Portugal",
				CountryCode: "PT",
				PostalCode:  "1000-001",
			},
		},
		{
			name: "No features",
			body: `{"features":[]}`,
		},
		{
			name:   "Bad request",
			status: http.StatusBadRequest,
			body:   `{"features":[{"properties":{"city":"Lisboa"}}]}`,
		},
	}, func(m *Maze) (*Location, error) {
		return m.Search("en", "lisbon")
	})
}

func TestAirport(t *testing.T) {
	t.Parallel()

	runLookupTests(t, []lookupTestCase{
		{
			name: "Full",
			body: `{"content":[{"icao":"LPPT","iata":"LIS","name":"Lisbon Humberto Delgado","coordinates":{"latitude":38.77,"longitude":-9.13},"country":{"iso2":"PT","name":"Portugal"},"servedCity":"Lisboa, Portugal"}]}`,
			location: &Location{
				Latitude:    38.77,
				Longitude:   -9.13,
				Name:        "Lisbon Humberto Delgado",
				Locality:    "Lisboa",
				Country:     "Portugal",
				CountryCode: "PT",
				ICAO:        "LPPT",
				IATA:        "LIS",
			},
		},
		{
			name: "No content",
			body: `{"content":[]}`,
		},
		{
			name: "Null",
			body: `null`,
		},
		{
			name:   "Rate limited",
			status: http.StatusTooManyRequests,
			body:   `{"content":[{"name":"Lisbon Humberto Delgado"}]}`,
		},
	}, func(m *Maze) (*Location, error) {
		return m.Airport("LIS")
	})
}
