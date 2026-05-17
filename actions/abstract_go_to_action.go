package actions

import "github.com/uglytoad/pdfpig/go/outline/destinations"

// AbstractGoToAction represents a GoTo-type action (GoTo, GoToE, GoToR) that has
// an explicit destination within the PDF.
type AbstractGoToAction struct {
	*PdfAction

	// Destination is the target location for the GoTo-type action.
	Destination destinations.ExplicitDestination
}

// Pdf returns the base PdfAction, satisfying the Action interface.
func (a *AbstractGoToAction) Pdf() *PdfAction {
	return a.PdfAction
}

// NewAbstractGoToAction creates a new AbstractGoToAction with the given type and
// destination.
func NewAbstractGoToAction(actionType ActionType, dest destinations.ExplicitDestination) *AbstractGoToAction {
	return &AbstractGoToAction{
		PdfAction:   NewPdfAction(actionType),
		Destination: dest,
	}
}
