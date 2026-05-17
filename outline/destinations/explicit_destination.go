package destinations

// ExplicitDestination represents a destination location within the same PDF file.
type ExplicitDestination struct {
	// PageNumber is the 1-based page number of the destination. A value of 0
	// means no page destination was available (i.e., an invalid explicit
	// destination).
	PageNumber int

	// Type is the display type of the destination.
	Type ExplicitDestinationType

	// Coordinates holds the display coordinates of the destination.
	Coordinates *ExplicitDestinationCoordinates
}

// NewExplicitDestination creates a new ExplicitDestination with the given page
// number, display type, and coordinates.
func NewExplicitDestination(pageNumber int, destType ExplicitDestinationType, coordinates *ExplicitDestinationCoordinates) ExplicitDestination {
	return ExplicitDestination{
		PageNumber:  pageNumber,
		Type:        destType,
		Coordinates: coordinates,
	}
}
