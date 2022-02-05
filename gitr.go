package gitr

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

// ToBbox creates bounding box for input Feature
func (f Feature) ToBbox() *Bbox {
	if f.BoundingBox != nil && len(f.BoundingBox) != 0 {

		sw := []float64{f.BoundingBox[0], f.BoundingBox[1]}
		ne := []float64{f.BoundingBox[2], f.BoundingBox[3]}

		return &Bbox{
			sw: sw,
			ne: ne,
			se: []float64{ne[0], sw[1]},
			nw: []float64{sw[0], ne[1]}}
	}

	return &Bbox{}
}
