package what3words

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const _defaultEndpoint = "https://api.what3words.com/v3"

// What3Words interface defines a set of methods that can be used to interact with the What3Words API.
// The methods provide functionality for auto-suggesting 3 word addresses, retrieving available languages,
// converting between 3 word addresses and coordinates, and retrieving a grid section for a given bounding box.
type What3Words interface {
	AutoSuggest(ctx context.Context, input *AutoSuggestInput) (*AutoSuggestResponse, error)
	AvailableLanguages(ctx context.Context) ([]Language, error)
	ConvertTo3wa(ctx context.Context, coordinates *Coordinates) (*LocationResponse, error)
	ConvertToCoordinates(ctx context.Context, words string) (*LocationResponse, error)
	GridSection(ctx context.Context, boundingBox *BoundingBox) (*GridSection, error)
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

// w3w contains the api key, language and endpoint for making requests to the w3w API.
type w3w struct {
	http     *http.Client
	apiKey   string
	language string
	endpoint *url.URL
}

// Option is an optional function parameter for the w3w struct
type Option func(*w3w)

// WithLanguage is a Functional Option for setting the w3w client language.
func WithLanguage(language string) Option {
	return func(w *w3w) {
		w.language = language
	}
}

// WithEndpoint is a Functional Option for setting the w3w client endpoint.
func WithEndpoint(endpoint *url.URL) Option {
	return func(w *w3w) {
		w.endpoint = endpoint
	}
}

// WithHTTPClient is a Functional Option for setting the w3w HTTP client.
// This option allows you to pass in various http client implementations such as hashicorp/go-retryablehttp
// which provides automatic retries and exponential backoff.
func WithHTTPClient(client *http.Client) Option {
	return func(w *w3w) {
		w.http = client
	}
}

// NewClient creates a new client which can be used to interact with the what3words API.
func NewClient(apiKey string, opts ...Option) What3Words {
	w3wURL, _ := url.Parse(_defaultEndpoint)

	w := &w3w{
		apiKey:   apiKey,
		language: "en",
		endpoint: w3wURL,
		http:     http.DefaultClient,
	}
	for _, opt := range opts {
		opt(w)
	}

	return w
}

// ConvertTo3wa This function will convert a latitude and longitude to a 3 word address, in the language of your choice.
// It also returns country, the bounds of the grid square, a nearby place (such as a local town) and a link to our map site.
func (w *w3w) ConvertTo3wa(ctx context.Context, coordinates *Coordinates) (*LocationResponse, error) {
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
func (w *w3w) GridSection(ctx context.Context, box *BoundingBox) (*GridSection, error) {
	u := w.endpoint.JoinPath("/grid-section")
	query := u.Query()
	query.Set("bounding-box", box.ToString())
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
