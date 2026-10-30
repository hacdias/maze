package maze

import (
	"errors"
	"fmt"
	"net/url"
)

type photonResponse struct {
	Features []photonFeature `json:"features"`
}

type photonFeature struct {
	Geometry   photonGeometry   `json:"geometry"`
	Properties photonProperties `json:"properties"`
}

type photonGeometry struct {
	Coordinates []float64 `json:"coordinates"`
}

type photonProperties struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	City        string `json:"city"`
	State       string `json:"state"`
	County      string `json:"county"`
	Country     string `json:"country"`
	CountryCode string `json:"countrycode"`
	PostCode    string `json:"postcode"`
}

// Reverse returns the [Location] at the given coordinates, using Photon. lang
// must be empty or a language Photon supports: "default", "de", "en" or "fr".
func (l *Maze) Reverse(lang string, lon, lat float64) (*Location, error) {
	uv := url.Values{}
	uv.Set("lat", fmt.Sprintf("%v", lat))
	uv.Set("lon", fmt.Sprintf("%v", lon))

	loc, err := l.photon("https://photon.komoot.io/reverse", lang, uv)
	if err != nil {
		return nil, err
	}

	loc.Latitude = lat
	loc.Longitude = lon
	return loc, nil
}

// Search returns the first [Location] that matches query, using Photon. lang is
// as in [Maze.Reverse].
func (l *Maze) Search(lang, query string) (*Location, error) {
	uv := url.Values{}
	uv.Set("q", query)

	loc, err := l.photon("https://photon.komoot.io/api/", lang, uv)
	if err != nil {
		return nil, err
	}

	loc.Name = query
	return loc, nil
}

// photon calls a Photon endpoint and returns its first feature as a [Location].
func (l *Maze) photon(endpoint, lang string, uv url.Values) (*Location, error) {
	// Photon rejects an empty lang.
	if lang != "" {
		uv.Set("lang", lang)
	}

	var fc photonResponse
	err := l.getJSON("photon", endpoint+"?"+uv.Encode(), &fc)
	if err != nil {
		return nil, err
	}

	if len(fc.Features) < 1 {
		return nil, errors.New("features missing from request")
	}

	p := fc.Features[0].Properties
	loc := &Location{
		Locality:    p.City,
		Region:      p.State,
		Country:     p.Country,
		CountryCode: p.CountryCode,
		PostalCode:  p.PostCode,
	}
	if loc.Region == "" {
		loc.Region = p.County
	}

	// A city, county or state does not repeat its own name in those fields.
	switch p.Type {
	case "city":
		if loc.Locality == "" {
			loc.Locality = p.Name
		}
	case "county", "state":
		if loc.Region == "" {
			loc.Region = p.Name
		}
	}

	if loc.Locality == "" && loc.Region == "" && loc.Country == "" {
		return nil, errors.New("no useful information found")
	}

	if c := fc.Features[0].Geometry.Coordinates; len(c) == 2 {
		loc.Longitude = c[0]
		loc.Latitude = c[1]
	}

	return loc, nil
}
