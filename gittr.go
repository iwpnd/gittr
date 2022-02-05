package gittr

import (
	geojson "github.com/paulmach/go.geojson"
)

// Bbox contains all 4 Positions of a bounding box
type Bbox struct {
	sw []float64
	nw []float64
	ne []float64
	se []float64
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

func swneToBbox(sw, ne []float64) Bbox {
	return Bbox{
		sw: sw,
		ne: ne,
		se: []float64{ne[0], sw[1]},
		nw: []float64{sw[0], ne[1]},
	}
}

// ToBbox creates bounding box for input Feature
func (f Feature) ToBbox() (Bbox, error) {
	if f.BoundingBox != nil && len(f.BoundingBox) != 0 {

		sw := []float64{f.BoundingBox[0], f.BoundingBox[1]}
		ne := []float64{f.BoundingBox[2], f.BoundingBox[3]}

		return swneToBbox(sw, ne), nil
	}

	switch f.Geometry.Type {
	case "Polygon":
		// set a starting point for the comparison
		outerRing := f.Geometry.Polygon[0]
		sw := []float64{outerRing[0][0], outerRing[0][1]}
		ne := []float64{outerRing[0][0], outerRing[0][1]}

		for _, p := range outerRing {
			if sw[0] > p[0] {
				sw[0] = p[0]
			}

			if sw[1] > p[1] {
				sw[1] = p[1]
			}

			if ne[0] < p[0] {
				ne[0] = p[0]
			}

			if ne[1] < p[1] {
				ne[1] = p[1]
			}
		}

		return swneToBbox(sw, ne), nil
	default:
		return Bbox{}, ErrUnsupportedGeometry{Type: string(f.Geometry.Type)}
	}
}
