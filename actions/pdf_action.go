package actions

// Action is the common interface for all PDF action types. It allows callers
// to preserve typed information (e.g., Destination on GoTo actions) when the
// action flows through generic *PdfAction parameters.
type Action interface {
	// Pdf returns the base PdfAction that carries the ActionType.
	Pdf() *PdfAction
}

// PdfAction represents a PDF action as defined in PDF reference 8.5.
type PdfAction struct {
	// Type is the kind of action this instance represents.
	Type ActionType
}

// Pdf returns the receiver, satisfying the Action interface for bare PdfAction values.
func (a *PdfAction) Pdf() *PdfAction {
	return a
}

// NewPdfAction creates a new PdfAction with the given type.
func NewPdfAction(actionType ActionType) *PdfAction {
	return &PdfAction{
		Type: actionType,
	}
}
