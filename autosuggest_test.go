package what3words

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestW3w_AutoSuggest(t *testing.T) {
	type response struct {
		statusCode int
		body       []byte
	}

	tests := map[string]struct {
		expectedError string
		expected      *AutoSuggestResponse
		input         *AutoSuggestInput
		response      response
	}{
		"auto suggest": {
			input: &AutoSuggestInput{
				Words: "film.crunchy.spiri",
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"suggestions": [
						{
							"country": "ES",
							"nearestPlace": "Caldes de Montbui, Catalonia",
							"words": "film.crunchy.spirit",
							"rank": 1,
							"language": "en"
						},
						{
							"country": "TH",
							"nearestPlace": "Lam Luk Ka, Pathum Thani",
							"words": "films.crunchy.spirit",
							"rank": 2,
							"language": "en"
						}
					]
				}`),
			},
			expected: &AutoSuggestResponse{
				Suggestions: []Suggestion{
					{
						Country:      "ES",
						NearestPlace: "Caldes de Montbui, Catalonia",
						Words:        "film.crunchy.spirit",
						Rank:         1,
						Language:     "en",
					},
					{
						Country:      "TH",
						NearestPlace: "Lam Luk Ka, Pathum Thani",
						Words:        "films.crunchy.spirit",
						Rank:         2,
						Language:     "en",
					},
				},
			},
		},
		"auto suggest with focus only": {
			input: &AutoSuggestInput{
				Words: "film.crunchy.spiri",
				Focus: &Coordinates{
					Lat: 50.842404,
					Lng: 4.361177,
				},
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"suggestions": [
						{
							"country": "BE",
							"nearestPlace": "Brussels, Brussels Capital",
							"words": "film.crunchy.spirits",
							"rank": 1,
							"distanceToFocusKm": 1,
							"language": "en"
						},
						{
							"country": "ES",
							"nearestPlace": "Caldes de Montbui, Catalonia",
							"words": "film.crunchy.spirit",
							"rank": 2,
							"distanceToFocusKm": 1039,
							"language": "en"
						},
						{
							"country": "NL",
							"nearestPlace": "Weert, Limburg",
							"words": "firm.crunchy.spires",
							"rank": 3,
							"distanceToFocusKm": 105,
							"language": "en"
						}
					]
				}`),
			},
			expected: &AutoSuggestResponse{
				Suggestions: []Suggestion{
					{
						Country:           "BE",
						NearestPlace:      "Brussels, Brussels Capital",
						Words:             "film.crunchy.spirits",
						Rank:              1,
						DistanceToFocusKm: 1,
						Language:          "en",
					},
					{
						Country:           "ES",
						NearestPlace:      "Caldes de Montbui, Catalonia",
						Words:             "film.crunchy.spirit",
						Rank:              2,
						DistanceToFocusKm: 1039,
						Language:          "en",
					},
					{
						Country:           "NL",
						NearestPlace:      "Weert, Limburg",
						Words:             "firm.crunchy.spires",
						Rank:              3,
						DistanceToFocusKm: 105,
						Language:          "en",
					},
				},
			},
		},
		"auto suggest with clip to circle": {
			input: &AutoSuggestInput{
				Words: "film.crunchy.spiri",
				ClipToCircle: &CoordinateRadius{
					Coordinates: Coordinates{
						Lat: 51.424388,
						Lng: -0.347452,
					},
					Radius: 10,
				},
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"suggestions": [
						{
							"country": "GB",
							"nearestPlace": "Ashford, Surrey",
							"words": "plan.flips.dawn",
							"rank": 1,
							"language": "en"
						}
					]
				}`),
			},
			expected: &AutoSuggestResponse{
				Suggestions: []Suggestion{
					{
						Country:      "GB",
						NearestPlace: "Ashford, Surrey",
						Words:        "plan.flips.dawn",
						Rank:         1,
						Language:     "en",
					},
				},
			},
		},
		"auto suggest with focus and clip to country ": {
			input: &AutoSuggestInput{
				Words: "plan.clips.a",
				Focus: &Coordinates{
					Lat: 51.4243877,
					Lng: -0.3474524,
				},
				ClipToCountry: []string{"GB"},
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"suggestions": [
						{
							"country": "GB",
							"nearestPlace": "Brixton Hill, London",
							"words": "plan.clips.area",
							"rank": 1,
							"distanceToFocusKm": 14,
							"language": "en"
						},
						{
							"country": "GB",
							"nearestPlace": "Ashford, Surrey",
							"words": "plan.flips.dawn",
							"rank": 2,
							"distanceToFocusKm": 7,
							"language": "en"
						},
						{
							"country": "GB",
							"nearestPlace": "Borehamwood, Hertfordshire",
							"words": "plan.clips.arts",
							"rank": 3,
							"distanceToFocusKm": 27,
							"language": "en"
						}
					]
				}`),
			},
			expected: &AutoSuggestResponse{
				Suggestions: []Suggestion{
					{
						Country:           "GB",
						NearestPlace:      "Brixton Hill, London",
						Words:             "plan.clips.area",
						Rank:              1,
						DistanceToFocusKm: 14,
						Language:          "en",
					},
					{
						Country:           "GB",
						NearestPlace:      "Ashford, Surrey",
						Words:             "plan.flips.dawn",
						Rank:              2,
						DistanceToFocusKm: 7,
						Language:          "en",
					},
					{
						Country:           "GB",
						NearestPlace:      "Borehamwood, Hertfordshire",
						Words:             "plan.clips.arts",
						Rank:              3,
						DistanceToFocusKm: 27,
						Language:          "en",
					},
				},
			},
		},
		"auto suggest with clip to bounding box": {
			input: &AutoSuggestInput{
				Words: "plan.clips.a",
				ClipToBoundingBox: &BoundingBox{
					SouthLat: 50,
					WestLng:  -4,
					NorthLat: 54,
					EastLng:  2,
				},
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"suggestions": [
						{
							"country": "GB",
							"nearestPlace": "Brixton Hill, London",
							"words": "plan.clips.area",
							"rank": 1,
							"language": "en"
						},
						{
							"country": "GB",
							"nearestPlace": "Skegness, Lincolnshire",
							"words": "plan.clips.army",
							"rank": 2,
							"language": "en"
						},
						{
							"country": "GB",
							"nearestPlace": "Borehamwood, Hertfordshire",
							"words": "plan.clips.arts",
							"rank": 3,
							"language": "en"
						}
					]
				}`),
			},
			expected: &AutoSuggestResponse{
				Suggestions: []Suggestion{
					{
						Country:      "GB",
						NearestPlace: "Brixton Hill, London",
						Words:        "plan.clips.area",
						Rank:         1,
						Language:     "en",
					},
					{
						Country:      "GB",
						NearestPlace: "Skegness, Lincolnshire",
						Words:        "plan.clips.army",
						Rank:         2,
						Language:     "en",
					},
					{
						Country:      "GB",
						NearestPlace: "Borehamwood, Hertfordshire",
						Words:        "plan.clips.arts",
						Rank:         3,
						Language:     "en",
					},
				},
			},
		},
		"auto suggest with clip to polygon": {
			input: &AutoSuggestInput{
				Words: "plan.clips.a",
				ClipToPolygon: []Coordinates{
					{
						Lat: 51.521,
						Lng: -0.343,
					},
					{
						Lat: 52.6,
						Lng: 2.3324,
					},
					{
						Lat: 54.234,
						Lng: 8.343,
					},
					{
						Lat: 51.521,
						Lng: -0.343,
					},
				},
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"suggestions": [
						{
							"country": "GB",
							"nearestPlace": "High Ongar, Essex",
							"words": "plan.clip.bags",
							"rank": 1,
							"language": "en"
						},
						{
							"country": "GB",
							"nearestPlace": "Wood Green, London",
							"words": "plan.slips.cage",
							"rank": 2,
							"language": "en"
						},
						{
							"country": "GB",
							"nearestPlace": "High Ongar, Essex",
							"words": "plan.flips.ants",
							"rank": 3,
							"language": "en"
						}
					]
				}`),
			},
			expected: &AutoSuggestResponse{
				Suggestions: []Suggestion{
					{
						Country:      "GB",
						NearestPlace: "High Ongar, Essex",
						Words:        "plan.clip.bags",
						Rank:         1,
						Language:     "en",
					},
					{
						Country:      "GB",
						NearestPlace: "Wood Green, London",
						Words:        "plan.slips.cage",
						Rank:         2,
						Language:     "en",
					},
					{
						Country:      "GB",
						NearestPlace: "High Ongar, Essex",
						Words:        "plan.flips.ants",
						Rank:         3,
						Language:     "en",
					},
				},
			},
		},
		"no suggestions returned": {
			input: &AutoSuggestInput{
				Words: "plan.clips.a",
			},
			response: response{
				statusCode: http.StatusOK,
				body: []byte(`{
					"suggestions": []
				}`),
			},
			expected: &AutoSuggestResponse{
				Suggestions: []Suggestion{},
			},
		},
		"error making auto suggest": {
			input: &AutoSuggestInput{
				Words: "",
			},
			response: response{
				statusCode: http.StatusBadRequest,
				body: []byte(`{
					"error": {
						"code": "MissingInput",
						"message": "input must be specified"
					}
				}`),
			},
			expected:      nil,
			expectedError: "retrieving auto suggestion",
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
			got, err := w.AutoSuggest(context.Background(), tt.input)
			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.Equal(t, tt.expected, got)
				assert.NoError(t, err)
			}
		})
	}
}
