package w3w_go_wrapper

import (
	"context"
	"net/http"
	"net/url"
)

const _defaultEndpoint = "https://api.what3words.com/v3"

// What3Words is an interface containing methods for interacting with the What3Words API.
type What3Words interface {
	AutoSuggest(ctx context.Context, input *AutoSuggestInput) (*AutoSuggestResponse, error)
	AvailableLanguages(ctx context.Context) ([]Language, error)
	ConvertTo3wa(ctx context.Context, coordinates Coordinates) (*LocationResponse, error)
	ConvertToCoordinates(ctx context.Context, words string) (*LocationResponse, error)
	GridSection(ctx context.Context, boundingBox BoundingBox) (*GridSection, error)
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

// NewClient creates a new client for interacting with the w3w API.
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
