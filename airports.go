package maze

import (
	"errors"
	"net/url"
	"strings"
)

type aviowikiResponse struct {
	Content []*aviowikiContent `json:"content"`
}

type aviowikiContent struct {
	ICAO        string              `json:"icao"`
	IATA        string              `json:"iata"`
	Name        string              `json:"name"`
	Coordinates aviowikiCoordinates `json:"coordinates"`
	Country     aviowikiCountry     `json:"country"`
	City        string              `json:"servedCity"`
}

type aviowikiCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type aviowikiCountry struct {
	ISO2 string `json:"iso2"`
	Name string `json:"name"`
}

func (l *Maze) Airport(query string) (*Location, error) {
	uv := url.Values{}
	uv.Set("query", query)
	uv.Set("size", "1")

	var avioRes *aviowikiResponse
	err := l.getJSON("aviowiki", "https://api.aviowiki.com/free/airports/search?"+uv.Encode(), &avioRes)
	if err != nil {
		return nil, err
	}

	if avioRes == nil || len(avioRes.Content) == 0 {
		return nil, errors.New("no airport found")
	}

	f := avioRes.Content[0]

	loc := &Location{
		Name:        f.Name,
		Latitude:    f.Coordinates.Latitude,
		Longitude:   f.Coordinates.Longitude,
		Country:     f.Country.Name,
		CountryCode: f.Country.ISO2,
		Locality:    strings.TrimSpace(strings.Split(f.City, ",")[0]),
		IATA:        f.IATA,
		ICAO:        f.ICAO,
	}

	return loc, nil
}
