package gittr

import (
	"encoding/json"
	"testing"
)

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

		if got.w != test.expected.w {
			t.Fatal("case:", test.tcase, "- W", got.w, "does not match expected ", test.expected.w)
		}
		if got.s != test.expected.s {
			t.Fatal("case:", test.tcase, "- S", got.s, "does not match expected ", test.expected.s)
		}
		if got.e != test.expected.e {
			t.Fatal("case:", test.tcase, "- E", got.e, "does not match expected ", test.expected.e)
		}
		if got.n != test.expected.n {
			t.Fatal("case:", test.tcase, "- N", got.n, "does not match expected ", test.expected.n)
		}
	}
}
