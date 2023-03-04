package w3w_go_wrapper

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestW3w_AutoSuggest(t *testing.T) {
	tests := map[string]struct {
		coordinates   Coordinates
		expectedError string
		words         string
	}{
		"something happened": {
			coordinates: Coordinates{
				Lng: 51.521251,
				Lat: 0.203586,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			w := NewClient("XXX")
			got, err := w.ConvertTo3wa(context.Background(), tt.coordinates)
			assert.NoError(t, err)
			log.Println(got)

		})
	}
}
