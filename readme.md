# gittr

Gittr cuts an input polygon into equal sized grid cells. Grid cells are only created
if they either touch or are within the input polygon. It also accounts for polygon holes.

## Installation

```
go get -u github.com/iwpnd/gittr
```

## Usage

```go
package main

import (
  "encoding/json"
  "fmt"

  "github.com/iwpnd/gittr"
  )

func main() {
  raw := []byte(`{
              "type": "Feature",
              "properties": {"id": 3},
              "geometry": {
                  "type": "Polygon",
                  "coordinates": [
                      [
                          [13.398424, 52.481369],
                          [13.395336, 52.480271],
                          [13.391905, 52.478598],
                          [13.390017, 52.478076],
                          [13.3871, 52.477971],
                          [13.386929, 52.471384],
                          [13.389503, 52.470338],
                          [13.392076, 52.46924],
                          [13.397395, 52.466626],
                          [13.402971, 52.465527],
                          [13.40992, 52.465789],
                          [13.417384, 52.469867],
                          [13.420558, 52.470338],
                          [13.419185, 52.473946],
                          [13.418499, 52.477501],
                          [13.418756, 52.479592],
                          [13.413781, 52.48048],
                          [13.406489, 52.482153],
                          [13.401856, 52.482728],
                          [13.398424, 52.481369]
                      ],
                      [
                          [13.391561, 52.475096],
                          [13.414552, 52.476194],
                          [13.415324, 52.47107],
                          [13.391647, 52.470756],
                          [13.391561, 52.475096]
                      ]
                  ]
              }
          }`)

  var f gittr.Feature
  err := json.Unmarshal(raw,&f)
  if err != nil {
      panic("something went wrong")
  }

  // create grid with a cell size of 100m
  grid := f.ToGrid(100)

  fmt.Println("cell count: %v", len(grid.Features))
}
```

## License

MIT

## Maintainer

Benjamin Ramser - [@iwpnd](https://github.com/iwpnd)

Project Link: [https://github.com/iwpnd/gittr](https://github.com/iwpnd/gittr)

## Acknowledgement

Paul Mach - [go.geojson](https://github.com/paulmach/go.geojson)
