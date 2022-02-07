package gittr

import (
	geojson "github.com/paulmach/go.geojson"
)

// Extent contains all 4 Positions of a bounding box
type Extent struct {
	s float64
	e float64
	n float64
	w float64
}

// Geometry extends geojson.Geometry
type Geometry struct {
	geojson.Geometry
}

// Feature extends geojson.Feature
type Feature struct {
	geojson.Feature
}

// FeatureCollection extends geojson.FeatureCollection
type FeatureCollection struct {
	geojson.FeatureCollection
}

// Extent creates bounding box for input Feature
func (f Feature) Extent() (Extent, error) {
	if f.BoundingBox != nil && len(f.BoundingBox) != 0 {

		w := f.BoundingBox[0]
		s := f.BoundingBox[1]
		e := f.BoundingBox[2]
		n := f.BoundingBox[3]

		return Extent{s, e, n, w}, nil
	}

	switch f.Geometry.Type {
	case "Polygon":
		// set a starting point for the comparison
		outerRing := f.Geometry.Polygon[0]
		w := outerRing[0][0]
		s := outerRing[0][1]
		e := outerRing[0][0]
		n := outerRing[0][1]

		for _, p := range outerRing {
			if w > p[0] {
				w = p[0]
			}

			if s > p[1] {
				s = p[1]
			}

			if e < p[0] {
				e = p[0]
			}

			if n < p[1] {
				n = p[1]
			}
		}

		return Extent{s, e, n, w}, nil
	default:
		return Extent{}, ErrUnsupportedGeometry{Type: string(f.Geometry.Type)}
	}
}
