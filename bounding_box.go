package what3words

import "fmt"

// BoundingBox defines the bounds of a rectangular area on a map or a grid.
// It is defined by four coordinates representing the southernmost latitude,
// westernmost longitude, northernmost latitude, and easternmost longitude of the area
type BoundingBox struct {
	SouthLat float64
	WestLng  float64
	NorthLat float64
	EastLng  float64
}

// ToString outputs the BoundingBox as a comma separated string
func (b BoundingBox) ToString() string {
	return fmt.Sprintf("%f,%f,%f,%f", b.SouthLat, b.WestLng, b.NorthLat, b.EastLng)
}

// NewBoundingBox constructs a BoundingBox
func NewBoundingBox(southLat, westLng, northLat, eastLng float64) *BoundingBox {
	return &BoundingBox{
		SouthLat: southLat,
		WestLng:  westLng,
		NorthLat: northLat,
		EastLng:  eastLng,
	}
}
