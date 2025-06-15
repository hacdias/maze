package maze

import (
	"net/http"
)

type Maze struct {
	httpClient *http.Client
}

func NewMaze(client *http.Client) *Maze {
	if client == nil {
		client = &http.Client{}
	}

	return &Maze{
		httpClient: client,
	}
}

func (l *Maze) Reverse(lang string, lon, lat float64) (*Location, error) {
	return l.photonReverse(lang, lon, lat)
}

func (l *Maze) ReverseGeoURI(lang, geoUri string) (*Location, error) {
	geo, err := ParseGeoURI(geoUri)
	if err != nil {
		return nil, err
	}

	return l.Reverse(lang, geo.Longitude, geo.Latitude)
}

func (l *Maze) Search(lang, query string) (*Location, error) {
	return l.photonSearch(lang, query)
}

func (l *Maze) Airport(query string) (*Location, error) {
	return l.aviowikiSearch(query)
}
