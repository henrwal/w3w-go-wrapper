package w3w_go_wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GridLine contains start and end coordinates of a line
type GridLine struct {
	Start Coordinates `json:"start"`
	End   Coordinates `json:"end"`
}

// GridSection contains horizontal and vertical lines covering a grid area
type GridSection struct {
	Lines []GridLine `json:"lines"`
}

// BoundingBox defines the bounds of the grid square
type BoundingBox struct {
	southLat float64
	westLng  float64
	northLat float64
	eastLng  float64
}

// ToString outputs the bounding box as a comma separated string
func (b BoundingBox) ToString() string {
	return fmt.Sprintf("%f,%f,%f,%f", b.southLat, b.westLng, b.northLat, b.eastLng)
}

// LocationResponse contains the country, the bounds of the grid square, a nearest place (such as a local town)
// and a link to the What3Words' map site
type LocationResponse struct {
	Coordinates  Coordinates `json:"coordinates"`
	Country      string      `json:"country"`
	Language     string      `json:"language"`
	Map          string      `json:"map"`
	NearestPlace string      `json:"nearestPlace"`
	Square       BoundingBox `json:"square"`
	Words        string      `json:"words"`
}

// Coordinates contain latitude and longitude which are encoded according to the World Geodetic System (WGS84).
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// ToString outputs the coordinates in the format "Latitude, Longitude"
func (c Coordinates) ToString() string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Lng)
}

// Language contains a language's ISO 639-1 2-letter code, english name and native name.
type Language struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	NativeName string `json:"nativeName"`
}

// AvailableLanguages contains a slice of Languages.
type AvailableLanguages struct {
	Languages []Language `json:"languages"`
}

// ConvertTo3wa This function will convert a latitude and longitude to a 3 word address, in the language of your choice.
// It also returns country, the bounds of the grid square, a nearby place (such as a local town) and a link to our map site.
func (w *w3w) ConvertTo3wa(ctx context.Context, coordinates Coordinates) (*LocationResponse, error) {
	u := w.endpoint.JoinPath("/convert-to-3wa")
	query := u.Query()
	query.Set("coordinates", coordinates.ToString())
	query.Set("language", w.language)
	u.RawQuery = query.Encode()

	var resp LocationResponse
	if err := w.request(ctx, u.String(), &resp); err != nil {
		return nil, fmt.Errorf("converting coordinates to 3 Word Address: %w", err)
	}

	return &resp, nil
}

// GridSection returns a section of the What3Words 3m x 3m grid as a set of horizontal and vertical lines
// covering the requested area, which can then be drawn onto a map.
func (w *w3w) GridSection(ctx context.Context, b BoundingBox) (*GridSection, error) {
	u := w.endpoint.JoinPath("/grid-section")
	query := u.Query()
	query.Set("bounding-box", b.ToString())
	query.Set("language", w.language)
	u.RawQuery = query.Encode()

	var resp GridSection
	if err := w.request(ctx, u.String(), &resp); err != nil {
		return nil, fmt.Errorf("retrieving grid section: %w", err)
	}

	return &resp, nil
}

// ConvertToCoordinates converts a 3 word address to a latitude and longitude. It also returns country,
// the bounds of the grid square, the nearest place (such as a local town) and a link to the What3Words map site.
func (w *w3w) ConvertToCoordinates(ctx context.Context, words string) (*LocationResponse, error) {
	u := w.endpoint.JoinPath("/convert-to-coordinates")
	query := u.Query()
	query.Set("words", words)
	query.Set("language", w.language)
	u.RawQuery = query.Encode()

	var resp LocationResponse
	if err := w.request(ctx, u.String(), &resp); err != nil {
		return nil, fmt.Errorf("converting w3w to coordinates: %w", err)
	}

	return &resp, nil
}

// AvailableLanguages Retrieves a list of all available 3 word address languages,
// including the ISO 3166-1 alpha-2 2-letter code, english name and native name.
func (w *w3w) AvailableLanguages(ctx context.Context) ([]Language, error) {
	u := w.endpoint.JoinPath("/available-languages").String()

	var resp AvailableLanguages
	if err := w.request(ctx, u, &resp); err != nil {
		return nil, fmt.Errorf("retrieving available languages: %w", err)
	}

	return resp.Languages, nil
}

func (w *w3w) request(ctx context.Context, url string, out interface{}) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", w.apiKey)

	resp, err := w.http.Do(request)
	if err != nil {
		return fmt.Errorf("sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request for %s returned unexpected status %s", url, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding response body into output: %w", err)
	}

	return nil
}
