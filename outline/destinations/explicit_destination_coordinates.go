package destinations

// ExplicitDestinationCoordinates holds the coordinates of the region to display
// for an explicit destination.
type ExplicitDestinationCoordinates struct {
	// Left is the left side of the region to display.
	Left *float64

	// Top is the top edge of the region to display.
	Top *float64

	// Right is the right side of the region to display.
	Right *float64

	// Bottom is the bottom edge of the region to display.
	Bottom *float64
}

// Empty is an empty set of coordinates where no values have been set.
var Empty = &ExplicitDestinationCoordinates{}

// floatPtr returns a pointer to the given float64 value.
func floatPtr(f float64) *float64 {
	return &f
}

// NewExplicitDestinationCoordinates creates a new ExplicitDestinationCoordinates
// with only the left coordinate set.
func NewExplicitDestinationCoordinates(left float64) *ExplicitDestinationCoordinates {
	return &ExplicitDestinationCoordinates{
		Left: floatPtr(left),
	}
}

// NewExplicitDestinationCoordinatesWithTop creates a new ExplicitDestinationCoordinates
// with the left and top coordinates set.
func NewExplicitDestinationCoordinatesWithTop(left, top float64) *ExplicitDestinationCoordinates {
	return &ExplicitDestinationCoordinates{
		Left: floatPtr(left),
		Top:  floatPtr(top),
	}
}

// NewExplicitDestinationCoordinatesRect creates a new ExplicitDestinationCoordinates
// with all four coordinates set defining a rectangular region.
func NewExplicitDestinationCoordinatesRect(left, top, right, bottom float64) *ExplicitDestinationCoordinates {
	return &ExplicitDestinationCoordinates{
		Left:   floatPtr(left),
		Top:    floatPtr(top),
		Right:  floatPtr(right),
		Bottom: floatPtr(bottom),
	}
}
