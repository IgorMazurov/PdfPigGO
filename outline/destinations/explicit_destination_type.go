package destinations

// ExplicitDestinationType defines the display type for opening an
// ExplicitDestination.
type ExplicitDestinationType byte

const (
	// XyzCoordinates displays the page with the given top left coordinates
	// and zoom level.
	XyzCoordinates ExplicitDestinationType = iota

	// FitPage fits the entire page within the window.
	FitPage

	// FitHorizontally fits the entire page width within the window.
	FitHorizontally

	// FitVertically fits the entire page height within the window.
	FitVertically

	// FitRectangle fits the rectangle specified by the ExplicitDestinationCoordinates
	// within the window.
	FitRectangle

	// FitBoundingBox fits the page's bounding box within the window.
	FitBoundingBox

	// FitBoundingBoxHorizontally fits the page's bounding box width within the window.
	FitBoundingBoxHorizontally

	// FitBoundingBoxVertically fits the page's bounding box height within the window.
	FitBoundingBoxVertically
)
