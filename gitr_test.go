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
		input    []byte
		expected Bbox
	}{
		{
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
	}

	for _, test := range tests {
		var f Feature
		err := json.Unmarshal(test.input, &f)
		if err != nil {
			t.Fatalf("cannot unmarshal test feature: %v", err)
		}

		got := f.ToBbox()

		if !equalPositions(got.sw, test.expected.sw) {
			t.Fatal("SW", got.sw, "does not match expected ", test.expected.sw)
		}
		if !equalPositions(got.ne, test.expected.ne) {
			t.Fatal("NE", got.ne, "does not match expected ", test.expected.ne)
		}
		if !equalPositions(got.se, test.expected.se) {
			t.Fatal("SE", got.se, "does not match expected ", test.expected.se)
		}
		if !equalPositions(got.nw, test.expected.nw) {
			t.Fatal("NW", got.nw, "does not match expected ", test.expected.nw)
		}
	}
}
