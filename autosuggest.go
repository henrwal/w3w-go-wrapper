package what3words

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// AutoSuggestResponse contains suggestions.
type AutoSuggestResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

// Suggestion represents a suggested 3 word address result from the AutoSuggest API.
type Suggestion struct {
	// Country is the full name of the country containing the suggested 3 word address.
	Country string `json:"country,omitempty"`

	// NearestPlace is the nearest named place associated with the suggested 3 word address.
	NearestPlace string `json:"nearestPlace,omitempty"`

	// Words is the suggested 3 word address.
	Words string `json:"words,omitempty"`

	// DistanceToFocusKm is the distance in kilometers from the suggested 3 word address to the focus specified in the AutoSuggest input.
	DistanceToFocusKm int `json:"distanceToFocusKm,omitempty"`

	// Rank is a score indicating the quality of the suggested 3 word address result.
	Rank int `json:"rank,omitempty"`

	// Language is the language in which the suggested 3 word address is given.
	Language string `json:"language,omitempty"`
}

// AutoSuggestInput contains the required and optional parameters for performing an AutoSuggestion request
type AutoSuggestInput struct {
	// Restrict AutoSuggest results to a bounding box, specified by coordinates
	ClipToBoundingBox *BoundingBox

	// clipToCircle Restrict AutoSuggest results to a circle,
	// specified by lat,lng,kilometres, where kilometres in the radius of the circle
	ClipToCircle *CoordinateRadius

	// clipToCountry: Restricts AutoSuggest to only return results inside the countries specified
	// by comma-separated list of uppercase ISO 3166-1 alpha-2 country codes.
	ClipToCountry []string

	// clipToPolygon: Restrict AutoSuggest results to a polygon, specified by a comma-separated list of lat,lng pairs.
	// The polygon should be closed, i.e. the first element should be repeated as the last
	// element; also the list should contain at least 4 entries.
	// The API is currently limited to accepting up to 25 pairs.
	ClipToPolygon PolygonCoordinates

	// focus: This is a location, specified as latitude,longitude.
	// If specified, the results will be weighted to give preference to those near the focus.
	// For convenience, longitude is allowed to wrap around the 180 line, so 361 is equivalent to 1.
	Focus *Coordinates

	// Language: For normal text input, specifies a fallback language,
	// which will help guide AutoSuggest if the input is particularly messy
	Language string

	// PreferLand: Makes AutoSuggest prefer results on land to those in the sea.
	// This setting is on by default. Use false to disable this setting and receive more suggestions in the sea.
	PreferLand *bool

	// Words: The full or partial 3 word address to obtain suggestions for.
	// At minimum this must be the first two complete words plus at least one character from the third word.
	// This is a required field for auto suggest
	Words string
}

// AutoSuggest Returns a list of 3 word addresses based on user input and other parameters.
func (w *w3w) AutoSuggest(ctx context.Context, input *AutoSuggestInput) (*AutoSuggestResponse, error) {
	u := w.endpoint.JoinPath("/autosuggest")
	query := u.Query()
	query.Set("input", input.Words)

	if input.Language != "" {
		query.Set("language", w.language)
	} else {
		query.Set("language", input.Language)
	}

	if input.Focus != nil {
		query.Set("focus", input.Focus.ToString())
	}

	if input.ClipToBoundingBox != nil {
		query.Set("clip-to-bounding-box", input.ClipToBoundingBox.ToString())
	}

	if input.ClipToCircle != nil {
		query.Set("clip-to-circle", input.ClipToCircle.ToString())
	}

	if input.ClipToPolygon != nil {
		if len(input.ClipToPolygon) > 24 {
			return nil, fmt.Errorf("clip to polygon is limited to 25 coordinate pairs")
		}
		query.Set("clip-to-polygon", input.ClipToPolygon.ToString())
	}

	if input.ClipToCountry != nil {
		query.Set("clip-to-country", strings.Join(input.ClipToCountry, ","))
	}

	if input.PreferLand == nil {
		preferLand := true
		input.PreferLand = &preferLand
	}

	query.Set("prefer-land", strconv.FormatBool(*input.PreferLand))

	u.RawQuery = query.Encode()

	var resp AutoSuggestResponse
	if err := w.request(ctx, u.String(), &resp); err != nil {
		return nil, fmt.Errorf("retrieving auto suggestion: %w", err)
	}

	return &resp, nil
}
