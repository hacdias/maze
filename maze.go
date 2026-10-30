package maze

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Maze struct {
	httpClient *http.Client
}

// NewMaze returns a [Maze] that uses client. If client is nil, a client with a
// 10 second timeout is used.
func NewMaze(client *http.Client) *Maze {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	return &Maze{
		httpClient: client,
	}
}

func (l *Maze) ReverseGeoURI(lang, geoUri string) (*Location, error) {
	p, err := ParseLocation(geoUri)
	if err != nil {
		return nil, err
	}

	loc, err := l.Reverse(lang, p.Longitude, p.Latitude)
	if err != nil {
		return nil, err
	}

	loc.Name = p.Name
	loc.Altitude = p.Altitude
	return loc, nil
}

// getJSON gets u and decodes its JSON body into v. service names the API in
// errors.
func (l *Maze) getJSON(service, u string, v any) error {
	res, err := l.httpClient.Get(u)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 200))
		return fmt.Errorf("%s: unexpected status %d: %s", service, res.StatusCode, b)
	}

	return json.NewDecoder(res.Body).Decode(v)
}
