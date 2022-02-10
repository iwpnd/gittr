package gittr

import (
	"encoding/json"
	"math"
	"testing"
)

// test helper to appoximate float64 equality
func approxEqual(want, got, tolerance float64) bool {
	diff := math.Abs(want - got)
	mean := math.Abs(want+got) / 2

	if math.IsNaN(diff / mean) {
		return true
	}

	return (diff / mean) < tolerance
}

func approxEqualPosition(a, b []float64, tolerance float64) bool {
	if len(a) != len(b) {
		return false
	}

	if !approxEqual(a[0], b[0], tolerance) || !approxEqual(a[1], b[1], tolerance) {
		return false
	}

	return true
}

func TestExtent(t *testing.T) {
	tests := []struct {
		tcase       string
		input       []byte
		expected    Extent
		expectedErr ErrUnsupportedGeometry
	}{
		{
			tcase: "feature has bbox",
			input: []byte(`{
                "type": "Feature",
                "properties": {},
                "bbox": [-74.004862, 40.726251, -73.999586, 40.730316],
                "geometry": {
                    "type": "Polygon",
                    "coordinates": [
                        [
                            [-74.004862, 40.726251],
                            [-73.999586, 40.726251],
                            [-73.999586, 40.730316],
                            [-74.004862, 40.730316],
                            [-74.004862, 40.726251]
                        ]
                    ]
                }
            }`),
			expected: Extent{
				w: -74.004862,
				s: 40.726251,
				e: -73.999586,
				n: 40.730316,
			},
		},
		{
			tcase: "feature doesnt have bbox",
			input: []byte(`{
                "type": "Feature",
                "properties": {},
                "geometry": {
                    "type": "Polygon",
                    "coordinates": [
                        [
                            [-74.004862, 40.726251],
                            [-73.999586, 40.726251],
                            [-73.999586, 40.730316],
                            [-74.004862, 40.730316],
                            [-74.004862, 40.726251]
                        ]
                    ]
                }
            }`),
			expected: Extent{
				w: -74.004862,
				s: 40.726251,
				e: -73.999586,
				n: 40.730316,
			},
		},
		{
			tcase: "unsupported geometry type",
			input: []byte(`{
                "type": "Feature",
                "properties": {},
                "geometry": {
                    "type": "Point",
                    "coordinates": [1,1]
                }
            }`),
			expectedErr: ErrUnsupportedGeometry{Type: "Point"},
		},
	}

	for _, test := range tests {
		var f Feature
		err := json.Unmarshal(test.input, &f)
		if err != nil {
			t.Fatalf("cannot unmarshal test feature: %v", err)
		}

		got, err := f.Extent()
		if err != nil && err != test.expectedErr {
			t.Fatal("something went wrong")
		}

		// err is expected and equal expectedErr
		if err == test.expectedErr {
			break
		}

		if f.BoundingBox == nil {
			t.Error("bounding box cannot be nil after .Extent()")
		}

		if got.w != test.expected.w {
			t.Error("case:", test.tcase, "- W", got.w, "does not match expected ", test.expected.w)
		}
		if got.s != test.expected.s {
			t.Error("case:", test.tcase, "- S", got.s, "does not match expected ", test.expected.s)
		}
		if got.e != test.expected.e {
			t.Error("case:", test.tcase, "- E", got.e, "does not match expected ", test.expected.e)
		}
		if got.n != test.expected.n {
			t.Error("case:", test.tcase, "- N", got.n, "does not match expected ", test.expected.n)
		}
	}
}

func TestExtentContains(t *testing.T) {
	test := []struct {
		e        Extent
		p        []float64
		expected bool
	}{
		// inside
		{
			e:        Extent{s: 0, w: 0, e: 1, n: 1},
			p:        []float64{0.5, 0.5},
			expected: true,
		},
		// outside
		{
			e:        Extent{s: 0, w: 0, e: 1, n: 1},
			p:        []float64{2, 2},
			expected: false,
		},
		// touches
		{
			e:        Extent{s: 0, w: 0, e: 1, n: 1},
			p:        []float64{0, 0},
			expected: true,
		},
	}

	for _, test := range test {
		got := test.e.contains(test.p)

		if got != test.expected {
			t.Errorf("expected %+v, got: %+v", test.expected, got)
		}
	}
}

func TestTerminal(t *testing.T) {
	test := []struct {
		origin, expected  []float64
		bearing, distance float64
	}{
		{
			origin:   []float64{13.35, 52.45},
			distance: 1112.758,
			bearing:  90,
			expected: []float64{13.3664, 52.45},
		},
		{
			origin:   []float64{0.0, 0.0},
			distance: 10000,
			bearing:  180,
			expected: []float64{0.0, -0.089932},
		},
		{
			origin:   []float64{13.35, -52.45},
			distance: 10000,
			bearing:  180,
			expected: []float64{13.35, -52.539932},
		},
	}

	for _, test := range test {
		got := terminal(test.origin, test.distance, test.bearing)
		lon, lat := got[0], got[1]

		if !approxEqual(test.expected[1], lat, 0.00001) {
			t.Errorf("expected %+v, got: %+v", test.expected[1], lat)
		}

		if !approxEqual(test.expected[0], lon, 0.00001) {
			t.Errorf("expected %+v, got: %+v", test.expected[0], lon)
		}
	}
}

func TestBearing(t *testing.T) {
	test := []struct {
		start, end []float64
		expected   float64
	}{
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{0.5, 0},
			expected: 180.0,
		},
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{0, 0},
			expected: 225.0,
		},
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{0, 0.5},
			expected: 270.0,
		},
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{0, 1},
			expected: 315.0,
		},
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{0.5, 1},
			expected: 0.0,
		},
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{1, 1},
			expected: 45.0,
		},
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{1, 0.5},
			expected: 90.0,
		},
		{
			start:    []float64{0.5, 0.5},
			end:      []float64{1, 0},
			expected: 135.0,
		},
	}

	for _, test := range test {
		got := bearing(test.start, test.end)

		if !approxEqual(got, test.expected, 0.001) {
			t.Errorf("expected %+v, got: %+v", test.expected, got)
		}
	}
}

func TestCreatePointsOnEdge(t *testing.T) {
	test := []struct {
		tcase          string
		start, end     []float64
		distance       float64
		expectedLine   [][]float64
		expectedLength int
	}{
		{
			// distance between start and end approx 680m
			// since distance is overshooting, so expected[len(expected-1)]
			// and start are supposed to be ~1000m apart
			tcase:          "overshooting",
			start:          []float64{13.37, 52.25},
			end:            []float64{13.36, 52.25},
			distance:       1000,
			expectedLine:   [][]float64{{13.37, 52.25}, {13.35531, 52.25}},
			expectedLength: 2,
		},
		{
			// distance between start and end approx 680m
			// since distance is NOT overshooting right from the start
			// it should return an array of points.
			tcase:          "should return array with 8 pts overshooting end",
			start:          []float64{13.37, 52.25},
			end:            []float64{13.36, 52.25},
			distance:       100,
			expectedLine:   [][]float64{{13.37, 52.25}, {13.35971728, 52.25}},
			expectedLength: 8,
		},
	}

	for _, test := range test {
		got := CreatePointsOnEdge(test.start, test.end, test.distance)
		l := len(got)

		start := got[0]
		end := got[l-1]

		if l != test.expectedLength {
			t.Errorf("%s - expected length %+v, got: %+v", test.tcase, test.expectedLength, got)
		}

		if !approxEqualPosition(start, test.expectedLine[0], 0.00001) {
			t.Errorf("%s - first position: expected %+v, got: %+v", test.tcase, test.expectedLine, got)
		}

		if !approxEqualPosition(end, test.expectedLine[len(test.expectedLine)-1], 0.00001) {
			t.Errorf("%s - last position: expected %+v, got: %+v", test.tcase, test.expectedLine, got)
		}
	}
}

func TestBuildGrid(t *testing.T) {
	tests := []struct {
		tcase          string
		input          []byte
		distance       float64
		expectedLength int
	}{
		{
			tcase: "feature has bbox",
			input: []byte(`{
              "type": "Feature",
              "properties": {"id": 1},
              "geometry": {
                "type": "Polygon",
                "coordinates": [
                      [
                        [0,0],
                        [1,0],
                        [1,1],
                        [0,1],
                        [0,0]
                      ]
                    ]
              }
            }`),
			distance:       27798,
			expectedLength: 20,
		},
	}

	for _, test := range tests {
		var f Feature
		err := json.Unmarshal(test.input, &f)
		if err != nil {
			t.Fatalf("cannot unmarshal test feature: %v", err)
		}

		fc, err := f.ToGrid(test.distance)
		if err != nil {
			t.Fatalf("failed to build grid: %v", err)
		}

		if len(fc.Features) != test.expectedLength {
			t.Errorf("%s - expected features to have length: %+v, got: %+v",
				test.tcase,
				test.expectedLength,
				len(fc.Features),
			)
		}

	}
}
