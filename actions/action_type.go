package actions

// ActionType represents PDF action types as defined in PDF reference 8.5.3.
type ActionType int

const (
	// GoTo goes to a destination in the current document.
	GoTo ActionType = iota

	// GoToR goes to a destination in another document (go-to remote).
	GoToR

	// GoToE goes to a destination in an embedded file (PDF 1.6, go-to embedded).
	GoToE

	// Launch launches an application, usually to open a file.
	Launch

	// Thread begins reading an article thread.
	Thread

	// URI resolves a uniform resource identifier.
	URI

	// Sound plays a sound (PDF 1.2).
	Sound

	// Movie plays a movie (PDF 1.2).
	Movie

	// Hide sets an annotation's Hidden flag (PDF 1.2).
	Hide

	// Named executes an action predefined by the viewer application (PDF 1.2).
	Named

	// SubmitForm sends data to a uniform resource locator (PDF 1.2).
	SubmitForm

	// ResetForm sets fields to their default values (PDF 1.2).
	ResetForm

	// ImportData imports field values from a file (PDF 1.2).
	ImportData

	// JavaScript executes a JavaScript script (PDF 1.3).
	JavaScript

	// SetOCGState sets the states of optional content groups (PDF 1.5).
	SetOCGState

	// Rendition controls the playing of multimedia content (PDF 1.5).
	Rendition

	// Trans updates the display of a document using a transition dictionary (PDF 1.5).
	Trans

	// GoTo3DView sets the current view of a 3D annotation (PDF 1.6).
	GoTo3DView
)
