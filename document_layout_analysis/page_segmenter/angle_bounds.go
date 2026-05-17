package page_segmenter

import "errors"

// AngleBounds defines the inclusive lower and upper bound (in degrees) for an angle.
type AngleBounds struct {
	Lower float64
	Upper float64
}

// NewAngleBounds creates a new AngleBounds with the given lower and upper bounds.
// Returns an error if lower is greater than or equal to upper.
func NewAngleBounds(lower, upper float64) (AngleBounds, error) {
	if lower >= upper {
		return AngleBounds{}, errors.New("DocstrumBoundingBoxes: the lower bound should be smaller than the upper bound")
	}
	return AngleBounds{Lower: lower, Upper: upper}, nil
}

// Contains returns whether the given angle falls within the bounds (inclusive).
func (ab AngleBounds) Contains(angle float64) bool {
	return angle >= ab.Lower && angle <= ab.Upper
}
