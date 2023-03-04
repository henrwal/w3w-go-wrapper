package w3w_go_wrapper

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := map[string]struct {
		apiKey   string
		expected *w3w
	}{
		"Client is created using default endpoint and language": {
			apiKey: "fakeAPIKey",
			expected: &w3w{
				apiKey:   "fakeAPIKey",
				language: "en",
				endpoint: &url.URL{
					Scheme: "https",
					Path:   "/v3",
					Host:   "api.what3words.com",
				},
				http: http.DefaultClient,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := NewClient(tt.apiKey)
			assert.Equal(t, tt.expected, got)
		})
	}
}
