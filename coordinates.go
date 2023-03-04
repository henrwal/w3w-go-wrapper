package what3words

import (
	"fmt"
	"strings"
)

// CoordinateRadius represents a circle specified by a central point (coordinates)
// and a radius (in kilometers). This is used to restrict search results to a specific area.
type CoordinateRadius struct {
	Coordinates Coordinates
	// Radius is the kilometer distance surrounding the coordinates
	Radius int
}

// ToString method converts the CoordinateRadius to a string representation
func (r CoordinateRadius) ToString() string {
	return fmt.Sprintf("%s,%d", r.Coordinates.ToString(), r.Radius)
}

// PolygonCoordinates represents a series of coordinate points used to define a polygon.
// The maximum number of coordinate points is 25 for auto-suggestion.
type PolygonCoordinates []Coordinates

// ToString method converts the PolygonCoordinates to a string representation
func (p PolygonCoordinates) ToString() string {
	polygon := make([]string, 0, len(p))
	for _, c := range p {
		polygon = append(polygon, c.ToString())
	}
	return strings.Join(polygon, ",")
}

// GridLine contains start and end coordinates of a line
type GridLine struct {
	Start Coordinates `json:"start"`
	End   Coordinates `json:"end"`
}

// GridSection contains horizontal and vertical lines covering a grid area
type GridSection struct {
	Lines []GridLine `json:"lines"`
}

// Square represents a geographical area defined by its southwest and northeast coordinates
type Square struct {
	Southwest Coordinates `json:"southwest"`
	Northeast Coordinates `json:"northeast"`
}

// LocationResponse contains the country, the bounds of the grid square, the nearest place (such as a local town)
// and a link to the What3Words' map site
type LocationResponse struct {
	Coordinates  Coordinates `json:"coordinates"`
	Country      string      `json:"country"`
	Language     string      `json:"language"`
	Map          string      `json:"map"`
	NearestPlace string      `json:"nearestPlace"`
	Square       Square      `json:"square"`
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
