package w3w_go_wrapper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Suggestion contains information from the auto-suggestion.
type Suggestion struct {
	Country           string `json:"country"`
	NearestPlace      string `json:"nearestPlace"`
	Words             string `json:"words"`
	DistanceToFocusKm int    `json:"distanceToFocusKm"`
	Rank              int    `json:"rank"`
	Language          string `json:"language"`
}

// AutoSuggestResponse contains suggestions.
type AutoSuggestResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

// CoordinateRadius TODO
type CoordinateRadius struct {
	Coordinates Coordinates
	// Radius is the kilometer distance surrounding the coordinates
	Radius int
}

func (r CoordinateRadius) ToString() string {
	return fmt.Sprintf("%s,%d", r.Coordinates.ToString(), r.Radius)
}

// AutoSuggestInput contains the required and optional parameters for performing an AutoSuggestion request
type AutoSuggestInput struct {
	// Restrict AutoSuggest results to a bounding box, specified by coordinates
	clipToBoundingBox *BoundingBox

	// clipToCircle Restrict AutoSuggest results to a circle,
	// specified by lat,lng,kilometres, where kilometres in the radius of the circle
	clipToCircle *CoordinateRadius

	// clipToCountry: Restricts AutoSuggest to only return results inside the countries specified
	// by comma-separated list of uppercase ISO 3166-1 alpha-2 country codes.
	clipToCountry []string

	// clipToPolygon: Restrict AutoSuggest results to a polygon, specified by a comma-separated list of lat,lng pairs.
	// The polygon should be closed, i.e. the first element should be repeated as the last
	// element; also the list should contain at least 4 entries.
	// The API is currently limited to accepting up to 25 pairs.
	clipToPolygon []Coordinates

	// focus: This is a location, specified as latitude,longitude.
	// If specified, the results will be weighted to give preference to those near the focus.
	// For convenience, longitude is allowed to wrap around the 180 line, so 361 is equivalent to 1.
	focus *Coordinates

	// Language: For normal text input, specifies a fallback language,
	// which will help guide AutoSuggest if the input is particularly messy
	language string

	// preferLand: Makes AutoSuggest prefer results on land to those in the sea.
	// This setting is on by default. Use false to disable this setting and receive more suggestions in the sea.
	preferLand bool

	// words: The full or partial 3 word address to obtain suggestions for.
	// At minimum this must be the first two complete words plus at least one character from the third word.
	// This is a required field for auto suggest
	words string
}

// AutoSuggest Returns a list of 3 word addresses based on user input and other parameters.
func (w *w3w) AutoSuggest(ctx context.Context, input *AutoSuggestInput) (*AutoSuggestResponse, error) {
	if len(input.clipToPolygon) > 24 {
		return nil, fmt.Errorf("clip to polygon is limited to 25 coordinate pairs")
	}

	u := w.endpoint.JoinPath("/autosuggest")
	query := u.Query()
	query.Set("input", input.words)

	if input.language != "" {
		query.Set("language", w.language)
	} else {
		query.Set("language", input.language)
	}

	if input.focus != nil {
		query.Set("focus", input.focus.ToString())
	}

	if input.clipToBoundingBox != nil {
		query.Set("clip-to-bounding-box", input.clipToBoundingBox.ToString())
	}

	if input.clipToCircle != nil {
		query.Set("clip-to-circle", input.clipToCircle.ToString())

	}

	if input.clipToPolygon != nil {
		query.Set("clip-to-bounding-box", input.clipToBoundingBox.ToString())
	}

	if input.clipToCountry != nil {
		query.Set("clip-to-country", strings.Join(input.clipToCountry, ","))
	}

	query.Set("prefer-land", strconv.FormatBool(input.preferLand))
	u.RawQuery = query.Encode()

	var resp AutoSuggestResponse
	if err := w.request(ctx, u.String(), &resp); err != nil {
		return nil, fmt.Errorf("retrieving auto suggestion: %w", err)
	}

	return &resp, nil
}
