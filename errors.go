package gitr

import "fmt"

// ErrUnsupportedGeometry ...
type ErrUnsupportedGeometry struct {
	Type string
}

func (e ErrUnsupportedGeometry) Error() string {
	return fmt.Sprintf("unsupported geometry type %s", e.Type)
}
