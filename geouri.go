package maze

import (
	"errors"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

// GeoURI represents a geo URI according to [rfc5870].
//
// [rfc5870]: https://datatracker.ietf.org/doc/html/rfc5870
type GeoURI struct {
	Latitude, Longitude, Altitude float64
	Parameters                    map[string][]string
	Query                         url.Values
}

const geoScheme = "geo"

func ParseGeoURI(uri string) (*GeoURI, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if u.Scheme != geoScheme {
		return nil, errors.New("invalid geo URI scheme")
	}

	g := &GeoURI{
		Parameters: map[string][]string{},
		Query:      u.Query(),
	}

	coordinates, parameters, _ := strings.Cut(u.Opaque, ";")
	if len(strings.TrimSpace(coordinates)) < 1 {
		return nil, errors.New("empty path")
	}

	g.Latitude, g.Longitude, g.Altitude, err = parseCoordinates(coordinates)
	if err != nil {
		return nil, fmt.Errorf("cannot parse coordinates: %w", err)
	}

	g.Parameters, err = parseParameters(parameters)
	if err != nil {
		return nil, fmt.Errorf("cannot parse parameters: %w", err)
	}

	return g, nil
}

func (g *GeoURI) String() string {
	u := &url.URL{
		Scheme:   geoScheme,
		Opaque:   encodeCoordinates(g.Latitude, g.Longitude, g.Altitude),
		RawQuery: g.Query.Encode(),
	}

	if parameters := encodeParameters(g.Parameters); parameters != "" {
		u.Opaque += ";" + parameters
	}

	return u.String()
}

func parseCoordinates(coordinates string) (lat float64, lon float64, alt float64, err error) {
	coords := strings.Split(coordinates, ",")

	if l := len(coords); l < 2 || l > 3 {
		return 0, 0, 0, errors.New("invalid number of coordinates, expected 2 or 3")
	}

	lat, err = strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot parse latitude: %w", err)
	}

	lon, err = strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot parse longitude: %w", err)
	}

	if len(coords) == 3 {
		alt, err = strconv.ParseFloat(coords[2], 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("cannot parse altitude: %w", err)
		}
	}

	return lat, lon, alt, nil
}

func encodeCoordinates(lat, lon, alt float64) string {
	latStr := strconv.FormatFloat(lat, 'f', -1, 64)
	lonStr := strconv.FormatFloat(lon, 'f', -1, 64)

	if alt == 0 {
		return fmt.Sprintf("%s,%s", latStr, lonStr)
	}

	altStr := strconv.FormatFloat(alt, 'f', -1, 64)
	return fmt.Sprintf("%s,%s,%s", latStr, lonStr, altStr)
}

func parseParameters(parameters string) (url.Values, error) {
	p := make(url.Values)

	for parameters != "" {
		var key string
		key, parameters, _ = strings.Cut(parameters, ";")
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		if strings.Contains(value, "=") {
			return nil, fmt.Errorf("invalid equal separator in value")
		}
		key, err := url.QueryUnescape(key)
		if err != nil {
			return nil, err
		}
		if value == "" {
			p[key] = []string{}
			continue
		}
		value, err = url.QueryUnescape(value)
		if err != nil {
			return nil, err
		}
		p[key] = append(p[key], value)
	}

	return p, nil
}

func encodeParameters(p map[string][]string) string {
	if len(p) == 0 {
		return ""
	}
	var buf strings.Builder
	for _, k := range slices.Sorted(maps.Keys(p)) {
		vs := p[k]
		keyEscaped := url.QueryEscape(k)

		// If there are no values for the key, we still need to write the key.
		if len(vs) == 0 {
			if buf.Len() > 0 {
				buf.WriteByte(';')
			}
			buf.WriteString(keyEscaped)
			continue
		}

		// If there are multiple values, we write them all.
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte(';')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}
