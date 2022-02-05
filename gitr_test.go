package gitr

import (
	"encoding/json"
	"testing"
)

func equalPositions(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}

	if a[0] != b[0] || a[1] != b[1] {
		return false
	}

	return true
}

func TestBbox(t *testing.T) {
	tests := []struct {
		tcase       string
		input       []byte
		expected    Bbox
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
			expected: Bbox{
				sw: []float64{-74.004862, 40.726251},
				ne: []float64{-73.999586, 40.730316},
				se: []float64{-73.999586, 40.726251},
				nw: []float64{-74.004862, 40.730316},
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
			expected: Bbox{
				sw: []float64{-74.004862, 40.726251},
				ne: []float64{-73.999586, 40.730316},
				se: []float64{-73.999586, 40.726251},
				nw: []float64{-74.004862, 40.730316},
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

		got, err := f.ToBbox()
		if err != nil && err != test.expectedErr {
			t.Fatal("something went wrong")
		}

		// err is expected and equal expectedErr
		if err == test.expectedErr {
			break
		}

		if !equalPositions(got.sw, test.expected.sw) {
			t.Fatal("case:", test.tcase, "- SW", got.sw, "does not match expected ", test.expected.sw)
		}
		if !equalPositions(got.ne, test.expected.ne) {
			t.Fatal("case:", test.tcase, "- NE", got.ne, "does not match expected ", test.expected.ne)
		}
		if !equalPositions(got.se, test.expected.se) {
			t.Fatal("case:", test.tcase, "- SE", got.se, "does not match expected ", test.expected.se)
		}
		if !equalPositions(got.nw, test.expected.nw) {
			t.Fatal("case:", test.tcase, "- NW", got.nw, "does not match expected ", test.expected.nw)
		}
	}
}
