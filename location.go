package maze

import "math"

type Location struct {
	Latitude  float64 `json:"latitude,omitempty" yaml:"latitude,omitempty" csv:"latitude"`
	Longitude float64 `json:"longitude,omitempty" yaml:"longitude,omitempty" csv:"longitude"`
	Altitude  float64 `json:"altitude,omitempty" yaml:"altitude,omitempty" csv:"altitude"`
	Name      string  `json:"name,omitempty" yaml:"name,omitempty" csv:"name"`
	Locality  string  `json:"locality,omitempty" yaml:"locality,omitempty" csv:"locality"`
	Region    string  `json:"region,omitempty" yaml:"region,omitempty" csv:"region"`
	Country   string  `json:"country,omitempty" yaml:"country,omitempty" csv:"country"`
	ICAO      string  `json:"icao,omitempty" yaml:"icao,omitempty" csv:"icao"`
	IATA      string  `json:"iata,omitempty" yaml:"iata,omitempty" csv:"iata"`
}

// ParseLocation parses a geo URI string and returns a [Location] struct. Query
// parameters are used to fill in additional fields, such as name, locality, etc.
func ParseLocation(geoUri string) (*Location, error) {
	geo, err := ParseGeoURI(geoUri)
	if err != nil {
		return nil, err
	}

	loc := &Location{
		Latitude:  geo.Latitude,
		Longitude: geo.Longitude,
		Altitude:  geo.Altitude,
	}

	if name, ok := geo.Query["name"]; ok && len(name) > 0 {
		loc.Name = name[0]
	}
	if locality, ok := geo.Query["locality"]; ok && len(locality) > 0 {
		loc.Locality = locality[0]
	}
	if region, ok := geo.Query["region"]; ok && len(region) > 0 {
		loc.Region = region[0]
	}
	if country, ok := geo.Query["country"]; ok && len(country) > 0 {
		loc.Country = country[0]
	}
	if icao, ok := geo.Query["icao"]; ok && len(icao) > 0 {
		loc.ICAO = icao[0]
	}
	if iata, ok := geo.Query["iata"]; ok && len(iata) > 0 {
		loc.IATA = iata[0]
	}

	return loc, nil
}

// Distance returns the distance, in meters, between l1 and l2.
func (l1 *Location) Distance(l2 *Location) float64 {
	if l2 == nil || l1 == nil {
		return 0
	}

	lat1 := l1.Latitude * (math.Pi / 180)
	lon1 := l1.Longitude * (math.Pi / 180)
	lat2 := l2.Latitude * (math.Pi / 180)
	lon2 := l2.Longitude * (math.Pi / 180)

	dlon := lon2 - lon1
	dlat := lat2 - lat1

	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Asin(math.Sqrt(a))
	r := float64(6371)

	return c * r * 1000
}

// String returns a geo URI representation of the [Location] struct, including
// query parameters for the fields that are not present in the base specification.
func (l *Location) String() string {
	geo := GeoURI{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Altitude:  l.Altitude,
		Query:     map[string][]string{},
	}

	if l.Name != "" {
		geo.Query["name"] = []string{l.Name}
	}
	if l.Locality != "" {
		geo.Query["locality"] = []string{l.Locality}
	}
	if l.Region != "" {
		geo.Query["region"] = []string{l.Region}
	}
	if l.Country != "" {
		geo.Query["country"] = []string{l.Country}
	}
	if l.ICAO != "" {
		geo.Query["icao"] = []string{l.ICAO}
	}
	if l.IATA != "" {
		geo.Query["iata"] = []string{l.IATA}
	}

	return geo.String()
}
