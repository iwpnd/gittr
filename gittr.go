package gittr

import (
	"math"

	"github.com/iwpnd/piper"
	geojson "github.com/paulmach/go.geojson"
)

const earthRadius = 6371008.8

// SpatialRelation to enumerate the spatial relation between extent and feature
type SpatialRelation string

// Spatial relations
const (
	Touches SpatialRelation = "touches"
	Within  SpatialRelation = "within"
)

// Extent contains all 4 Positions of a bounding box
type Extent struct {
	s        float64
	e        float64
	n        float64
	w        float64
	relation SpatialRelation
}

// Feature extends geojson.Feature
type Feature geojson.Feature

func radToDegree(rad float64) float64 {
	return rad * 180 / math.Pi
}

func degreeToRad(degree float64) float64 {
	return degree * math.Pi / 180
}

func distanceToRadians(distance float64) float64 {
	const r = earthRadius

	return distance / r
}

// terminal calculates the terminal position travelling a distance
// from a given origin
// see https://www.movable-type.co.uk/scripts/latlong.html
func terminal(start []float64, distance, bearing float64) []float64 {
	latRad1 := degreeToRad(start[1])
	lonRad1 := degreeToRad(start[0])

	bearingRad := degreeToRad(bearing)
	distanceRad := distanceToRadians(distance)

	latRad2 := math.Asin(
		math.Sin(latRad1)*
			math.Cos(distanceRad) +
			math.Cos(latRad1)*
				math.Sin(distanceRad)*
				math.Cos(bearingRad))

	lonRad2 := lonRad1 + math.Atan2(
		math.Sin(bearingRad)*
			math.Sin(distanceRad)*
			math.Cos(latRad1),
		math.Cos(distanceRad)-
			math.Sin(latRad1)*
				math.Sin(latRad2))

	// cap decimals at .00000001 degree ~= 1.11mm
	lon := math.Round(radToDegree(lonRad2)*100000000) / 100000000
	lat := math.Round(radToDegree(latRad2)*100000000) / 100000000

	return []float64{lon, lat}
}

func bearing(start, end []float64) float64 {
	lat1 := degreeToRad(start[1])
	lat2 := degreeToRad(end[1])
	lng1 := degreeToRad(start[0])
	lng2 := degreeToRad(end[0])

	a := math.Sin(lng2-lng1) * math.Cos(lat2)

	b := (math.Cos(lat1)*
		math.Sin(lat2) -
		math.Sin(lat1)*
			math.Cos(lat2)*
			math.Cos(lng2-lng1))

	o := math.Atan2(a, b)

	return math.Mod((o*180/math.Pi + 360), 360.0)
}

func haversine(start, end []float64) float64 {
	lat1 := degreeToRad(start[1])
	lat2 := degreeToRad(end[1])
	dlat := degreeToRad(end[1] - start[1])
	dlng := degreeToRad(end[0] - start[0])

	a := (math.Sin(dlat/2)*
		math.Sin(dlat/2) +
		math.Cos(lat1)*
			math.Cos(lat2)*
			math.Sin(dlng/2)*
			math.Sin(dlng/2))

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// Extent creates bounding box for input Feature
// if a bounding box is present, it returns early
// if theres no bounding box it'll be created an attached
func (f *Feature) Extent() (*Extent, error) {
	if f.BoundingBox != nil && len(f.BoundingBox) != 0 {

		w := f.BoundingBox[0]
		s := f.BoundingBox[1]
		e := f.BoundingBox[2]
		n := f.BoundingBox[3]

		return &Extent{s: s, e: e, n: n, w: w}, nil
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

		f.BoundingBox = []float64{s, e, n, w}

		return &Extent{s: s, e: e, n: n, w: w}, nil
	default:
		return &Extent{}, ErrUnsupportedGeometry{Type: string(f.Geometry.Type)}
	}
}

func (e *Extent) intersects(polygon [][][]float64) bool {
	counter := 0
	pts := [][]float64{{e.w, e.s}, {e.w, e.n}, {e.e, e.n}, {e.e, e.s}}

	for i := range pts {
		if piper.Pip(pts[i], polygon) {
			counter++
		}
	}

	if counter >= 1 && counter < 4 {
		e.relation = Touches
		return true
	} else if counter == 4 {
		e.relation = Within
		return true
	}

	return false
}

// ToGrid cuts features bounding box into a nxn grid
func (f *Feature) ToGrid(distance float64) (*geojson.FeatureCollection, error) {
	e, err := f.Extent()

	if err != nil {
		return &geojson.FeatureCollection{}, err
	}

	columns := pointsOnLine([]float64{e.w, e.n}, []float64{e.e, e.n}, distance)
	rows := pointsOnLine([]float64{e.w, e.n}, []float64{e.w, e.s}, distance)

	c := &geojson.FeatureCollection{}

	for i := 0; i < len(columns)-1; i++ {
		for j := 0; j < len(rows)-1; j++ {
			s := rows[j+1][1]
			w := columns[i][0]
			n := rows[j][1]
			e := columns[i+1][0]

			ext := Extent{s: s, e: e, w: w, n: n}

			if ext.intersects(f.Geometry.Polygon) {
				extf := ext.toFeature()
				extf.BoundingBox = []float64{e, s, w, n}
				extf.Properties = f.Properties
				c.AddFeature(extf)
			}
		}
	}

	return c, nil
}

func (e Extent) contains(p []float64) bool {
	lon, lat := p[0], p[1]
	return (((e.w <= lon) && (lon <= e.w)) ||
		((e.e <= lon) && (lon <= e.w)) ||
		((e.s <= lat) && (lat <= e.n)) ||
		((e.n <= lat) && (lat <= e.s)))
}

func (e Extent) toFeature() *geojson.Feature {
	p := [][][]float64{
		{{e.w, e.s}, {e.w, e.n}, {e.e, e.n}, {e.e, e.s}, {e.w, e.s}},
	}
	f := geojson.NewPolygonFeature(p)
	return f
}

// CreatePointsOnEdge creates points on and along a
// line spanning from {start} to {end} every {distance} meters
// if input {distance} is bigger than the haversine distance
// between {start} and {end} it creates the last point
// {distance}m from {start} overshooting {end}
func pointsOnLine(start, end []float64, distance float64) [][]float64 {
	b := bearing(start, end)
	d := haversine(start, end)

	pts := [][]float64{start}

	// if desired distance overshoots the endpoint
	// use terminal of overshot and return early
	if distance > d {
		pts = append(pts, terminal(start, distance, b))
		return pts
	}

	// append points until travelled > distance between start and end
	// last point in pts array is the new start for terminal
	for t := 0.0; t < d; t += distance {
		s := pts[len(pts)-1]
		p := terminal(s, distance, b)
		pts = append(pts, p)
	}

	return pts
}
