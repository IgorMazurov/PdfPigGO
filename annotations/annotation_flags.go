package annotations

// AnnotationFlags specifies characteristics of an annotation in a PDF or FDF
// document. Multiple flags can be combined using bitwise OR.
type AnnotationFlags int

const (
	// Invisible indicates not to display the annotation if it is not one of the standard annotation types.
	Invisible AnnotationFlags = 1 << 0

	// Hidden indicates not to display or print the annotation irrespective of type and do not allow interaction.
	Hidden = 1 << 1

	// Print indicates that the annotation should be included when the document is physically printed.
	Print = 1 << 2

	// NoZoom indicates not to zoom/scale the annotation as the zoom of the document is changed.
	NoZoom = 1 << 3

	// NoRotate indicates not to rotate the annotation as the page is rotated.
	NoRotate = 1 << 4

	// NoView indicates not to display the annotation in viewer applications, however allow the annotation
	// to be printed if Print is set.
	NoView = 1 << 5

	// ReadOnly allows the annotation to be displayed/printed if applicable but does not respond to user interaction.
	ReadOnly = 1 << 6

	// Locked indicates not to allow deleting the annotation or changing size/position but allows the contents to be modified.
	Locked = 1 << 7

	// ToggleNoView inverts the meaning of the NoView flag.
	ToggleNoView = 1 << 8

	// LockedContents allows the annotation to be deleted, resized, moved or restyled but disallows changes
	// to the annotation contents. Opposite to Locked.
	LockedContents = 1 << 9
)
