package maze

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	geojson "github.com/paulmach/go.geojson"
)

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
	uv.Set("lang", lang)

	res, err := l.httpClient.Get(endpoint + "?" + uv.Encode())
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}

	if len(fc.Features) < 1 {
		return nil, errors.New("features missing from request")
	}

	f := fc.Features[0]
	loc := &Location{
		Locality:    f.PropertyMustString("city", ""),
		Region:      f.PropertyMustString("state", f.PropertyMustString("county", "")),
		Country:     f.PropertyMustString("country", ""),
		CountryCode: f.PropertyMustString("countrycode", ""),
		PostalCode:  f.PropertyMustString("postcode", ""),
	}

	if loc.Locality == "" && loc.Region == "" && loc.Country == "" {
		return nil, errors.New("no useful information found")
	}

	if f.Geometry != nil && len(f.Geometry.Point) == 2 {
		loc.Longitude = f.Geometry.Point[0]
		loc.Latitude = f.Geometry.Point[1]
	}

	return loc, nil
}
