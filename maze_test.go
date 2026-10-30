package maze

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// fakeTransport answers every request with status and body, without using the
// network.
type fakeTransport struct {
	status int
	body   string
	query  url.Values // query of the last request
}

func (f *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	f.query = r.URL.Query()
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

func TestReverseGeoURI(t *testing.T) {
	t.Parallel()

	for _, uri := range []string{"geo:38.7,-9.1,120?name=Home", "geo:38.7,-9.1,120;name=Home"} {
		t.Run(uri, func(t *testing.T) {
			t.Parallel()

			ft := &fakeTransport{status: http.StatusOK, body: `{"features":[{"properties":{"city":"Lisboa"}}]}`}
			loc, err := NewMaze(&http.Client{Transport: ft}).ReverseGeoURI("en", uri)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if loc.Name != "Home" {
				t.Errorf("expected name Home, got %q", loc.Name)
			}
			if loc.Altitude != 120 {
				t.Errorf("expected altitude 120, got %v", loc.Altitude)
			}
		})
	}
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
			name: "City",
			body: `{"features":[{"properties":{"type":"city","name":"Lisboa","country":"Portugal"}}]}`,
			location: &Location{
				Name:     "lisbon",
				Locality: "Lisboa",
				Country:  "Portugal",
			},
		},
		{
			name: "County",
			body: `{"features":[{"properties":{"type":"county","name":"Lisboa","country":"Portugal"}}]}`,
			location: &Location{
				Name:    "lisbon",
				Region:  "Lisboa",
				Country: "Portugal",
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

func TestSearchEmptyLang(t *testing.T) {
	t.Parallel()

	ft := &fakeTransport{status: http.StatusOK, body: `{"features":[{"properties":{"city":"Lisboa"}}]}`}
	_, err := NewMaze(&http.Client{Transport: ft}).Search("", "lisbon")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ft.query.Has("lang") {
		t.Errorf("expected no lang, got %q", ft.query.Get("lang"))
	}
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
