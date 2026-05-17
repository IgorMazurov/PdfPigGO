package actions

import "github.com/uglytoad/pdfpig/go/outline/destinations"

// GoToAction represents a PDF action that goes to a destination in the current document (PDF reference 8.5.1).
type GoToAction struct {
	*AbstractGoToAction
}

// NewGoToAction creates a new GoToAction with the given destination.
func NewGoToAction(dest destinations.ExplicitDestination) *GoToAction {
	return &GoToAction{
		AbstractGoToAction: NewAbstractGoToAction(GoTo, dest),
	}
}
