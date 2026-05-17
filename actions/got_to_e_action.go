package actions

import "github.com/uglytoad/pdfpig/go/outline/destinations"

// GoToEAction represents a PDF action that goes to a destination in an embedded file (PDF 1.6, reference 8.5.1).
type GoToEAction struct {
	*AbstractGoToAction

	// FileSpecification is the specification of the embedded file.
	FileSpecification string
}

// NewGoToEAction creates a new GoToEAction with the given destination and file specification.
func NewGoToEAction(dest destinations.ExplicitDestination, fileSpecification string) *GoToEAction {
	return &GoToEAction{
		AbstractGoToAction:  NewAbstractGoToAction(GoToE, dest),
		FileSpecification:   fileSpecification,
	}
}
