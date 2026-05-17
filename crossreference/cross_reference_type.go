package crossreference

// CrossReferenceType represents the type of a cross-reference section in a PDF document.
type CrossReferenceType byte

const (
	// Table represents a cross-reference table.
	Table CrossReferenceType = 0

	// Stream represents a cross-reference stream.
	Stream CrossReferenceType = 1
)
