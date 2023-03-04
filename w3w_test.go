package what3words

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func TestW3w_ConvertToCoordinates(t *testing.T) {
	type response struct {
		statusCode int
		body       []byte
	}

	tests := map[string]struct {
		expected      *LocationResponse
		expectedError string
		response      response
		words         string
	}{
		"successfully convert W3W to coordinates": {
			expected: &LocationResponse{
				Coordinates: Coordinates{
					Lat: 51.520847,
					Lng: -0.195521,
				},
				Country:      "GB",
				Language:     "en",
				Map:          "https://w3w.co/filled.count.soap",
				NearestPlace: "Bayswater, London",
				Square: Square{
					Southwest: Coordinates{
						Lat: 51.520833,
						Lng: -0.195543,
					},
					Northeast: Coordinates{
						Lat: 51.52086,
						Lng: -0.195499,
					},
				},
				Words: "filled.count.soap",
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
				"country": "GB",
				"square": {
					"southwest": {
						"lng": -0.195543,
						"lat": 51.520833
					},
					"northeast": {
						"lng": -0.195499,
						"lat": 51.52086
					}
				},
				"nearestPlace": "Bayswater, London",
				"coordinates": {
					"lng": -0.195521,
					"lat": 51.520847
				},
				"words": "filled.count.soap",
				"language": "en",
				"map": "https://w3w.co/filled.count.soap"
			}`),
			},
			words: "filled.count.soap",
		},
		"error converting invalid W3W to coordinates": {
			words:         "invalidwhat3words",
			expectedError: "converting w3w to coordinates",
			response: response{
				statusCode: http.StatusBadRequest,
				body: []byte(`{
					"error": {
						"code": "BadWords",
						"message": "words must be a valid 3 word address, such as filled.count.soap or ///filled.count.soap"
					}
				}`),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(tt.response.statusCode)
				_, err := rw.Write(tt.response.body)
				assert.NoError(t, err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.NoError(t, err)

			w := NewClient("example-api-key", WithEndpoint(u))
			got, err := w.ConvertToCoordinates(context.Background(), tt.words)
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.Equal(t, tt.expected, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestW3w_AvailableLanguages(t *testing.T) {
	type response struct {
		statusCode int
		body       []byte
	}

	tests := map[string]struct {
		expected      []Language
		expectedError string
		response      response
	}{
		"successfully retrieve available languages": {
			expected: []Language{
				{
					Code:       "de",
					Name:       "German",
					NativeName: "Deutsch",
				},
				{
					Code:       "hi",
					Name:       "Hindi",
					NativeName: "हिन्दी",
				},
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
				  "languages": [
					{
					  "nativeName": "Deutsch",
					  "code": "de",
					  "name": "German"
					},
					{
					  "nativeName": "हिन्दी",
					  "code": "hi",
					  "name": "Hindi"
					}
				  ]
				}`),
			},
		},
		"error retrieving available languages": {
			expectedError: "retrieving available languages",
			response: response{
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(tt.response.statusCode)
				_, err := rw.Write(tt.response.body)
				assert.NoError(t, err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.NoError(t, err)

			w := NewClient("example-api-key", WithEndpoint(u))
			got, err := w.AvailableLanguages(context.Background())
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.Equal(t, tt.expected, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestW3w_ConvertTo3wa(t *testing.T) {
	type response struct {
		statusCode int
		body       []byte
	}

	tests := map[string]struct {
		expected      *LocationResponse
		expectedError string
		coordinates   Coordinates
		response      response
	}{
		"successfully convert coordinates to 3 word address": {
			expected: &LocationResponse{
				Coordinates: Coordinates{
					Lat: 51.520847,
					Lng: -0.195521,
				},
				Country:      "GB",
				Language:     "en",
				Map:          "https://w3w.co/filled.count.soap",
				NearestPlace: "Bayswater, London",
				Square: Square{
					Southwest: Coordinates{
						Lat: 51.520833,
						Lng: -0.195543,
					},
					Northeast: Coordinates{
						Lat: 51.52086,
						Lng: -0.195499,
					},
				},
				Words: "filled.count.soap",
			},
			coordinates: Coordinates{
				Lat: 51.520847,
				Lng: -0.195521,
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"country": "GB",
					"square": {
						"southwest": {
							"lng": -0.195543,
							"lat": 51.520833
						},
						"northeast": {
							"lng": -0.195499,
							"lat": 51.52086
						}
					},
					"nearestPlace": "Bayswater, London",
					"coordinates": {
						"lng": -0.195521,
						"lat": 51.520847
					},
					"words": "filled.count.soap",
					"language": "en",
					"map": "https://w3w.co/filled.count.soap"
				}`),
			},
		},
		"error converting coordinates to 3wa": {
			expectedError: "converting coordinates to 3 Word Address",
			response: response{
				statusCode: http.StatusBadRequest,
				body: []byte(`{
					"error": {
						"code": "BadCoordinates",
						"message": "latitude must be >=-90 and <= 90"
					}
				}`),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(tt.response.statusCode)
				_, err := rw.Write(tt.response.body)
				assert.NoError(t, err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.NoError(t, err)

			w := NewClient("example-api-key", WithEndpoint(u))
			got, err := w.ConvertTo3wa(context.Background(), &tt.coordinates)
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.Equal(t, tt.expected, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestW3w_GridSection(t *testing.T) {
	type response struct {
		statusCode int
		body       []byte
	}

	tests := map[string]struct {
		expected      *GridSection
		expectedError string
		boundingBox   BoundingBox
		response      response
	}{
		"successfully return grid section": {
			boundingBox: BoundingBox{
				SouthLat: 52.207988,
				WestLng:  0.116126,
				NorthLat: 52.208867,
				EastLng:  0.117540,
			},
			expected: &GridSection{
				Lines: []GridLine{
					{
						Start: Coordinates{
							Lat: 52.20801,
							Lng: 0.116126,
						},
						End: Coordinates{
							Lat: 52.20801,
							Lng: 0.11754,
						},
					},
					{
						Start: Coordinates{
							Lat: 52.208037,
							Lng: 0.116126,
						},
						End: Coordinates{
							Lat: 52.208037,
							Lng: 0.11754,
						},
					},
				},
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
				  "lines": [
					{
					  "start": {
						"lng": 0.116126,
						"lat": 52.20801
					  },
					  "end": {
						"lng": 0.11754,
						"lat": 52.20801
					  }
					},
					{
					  "start": {
						"lng": 0.116126,
						"lat": 52.208037
					  },
					  "end": {
						"lng": 0.11754,
						"lat": 52.208037
					  }
					}
				  ]
				}`),
			},
		},
		"error retrieving grid section": {
			expectedError: "retrieving grid section",
			response: response{
				statusCode: http.StatusBadRequest,
				body: []byte(`{
					"error": {
						"code": "BadBoundingBox",
						"message": "Failed to parse coordinates for bounding-box. Must be valid decimal numbers."
					}
				}`),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(tt.response.statusCode)
				_, err := rw.Write(tt.response.body)
				assert.NoError(t, err)
			}))
			defer ts.Close()

			u, err := url.Parse(ts.URL)
			assert.NoError(t, err)

			w := NewClient("example-api-key", WithEndpoint(u))
			got, err := w.GridSection(context.Background(), &tt.boundingBox)
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.Equal(t, tt.expected, got)
				assert.NoError(t, err)
			}
		})
	}
}
