package actions

import "github.com/uglytoad/pdfpig/go/outline/destinations"

// GoToRAction represents a PDF action that goes to a destination in another document (PDF reference 8.5.1).
type GoToRAction struct {
	*AbstractGoToAction

	// Filename is the name of the remote PDF file.
	Filename string
}

// NewGoToRAction creates a new GoToRAction with the given destination and filename.
func NewGoToRAction(dest destinations.ExplicitDestination, filename string) *GoToRAction {
	return &GoToRAction{
		AbstractGoToAction: NewAbstractGoToAction(GoToR, dest),
		Filename:           filename,
	}
}
