package annotations

// AnnotationType represents the standard annotation types in PDF documents.
type AnnotationType int

const (
	// Text is a 'sticky note' style annotation displaying some text with open/closed pop-up state.
	Text AnnotationType = iota

	// Link is a link to elsewhere in the document or an external application/web link.
	Link

	// FreeText displays text on the page. Unlike Text there is no associated pop-up.
	FreeText

	// Line displays a single straight line on the page with optional line ending styles.
	Line

	// Square displays a rectangle on the page.
	Square

	// Circle displays an ellipse on the page.
	Circle

	// Polygon displays a closed polygon on the page.
	Polygon

	// PolyLine displays a set of connected lines on the page which is not a closed polygon.
	PolyLine

	// Highlight is a highlight for text or content with associated annotation text.
	Highlight

	// Underline is an underline under text with associated annotation text.
	Underline

	// Squiggly is a jagged squiggly line under text with associated annotation text.
	Squiggly

	// StrikeOut is a strikeout through some text with associated annotation text.
	StrikeOut

	// Stamp is text or graphics intended to display as if inserted by a rubber stamp.
	Stamp

	// Caret is a visual symbol indicating the presence of text edits.
	Caret

	// Ink is a freehand 'scribble' formed by one or more paths.
	Ink

	// Popup displays text in a pop-up window for entry or editing.
	Popup

	// FileAttachment represents a file.
	FileAttachment

	// Sound represents a sound to be played through speakers.
	Sound

	// Movie embeds a movie from a file in a PDF document.
	Movie

	// Widget is used by interactive forms to represent field appearance and manage user interactions.
	Widget

	// Screen specifies a page region for media clips to be played and actions to be triggered from.
	Screen

	// PrinterMark represents a symbol used during the physical printing process to maintain output quality,
	// e.g. color bars or cut marks.
	PrinterMark

	// TrapNet is used during the physical printing process to prevent colors mixing.
	TrapNet

	// Watermark adds a watermark at a fixed size and position irrespective of page size.
	Watermark

	// Artwork3D represents a 3D model/artwork, for example from CAD, in a PDF document.
	Artwork3D

	// Other is a custom annotation type.
	Other
)

// String returns the name of the annotation type.
func (t AnnotationType) String() string {
	return annotationTypeName[t]
}

var annotationTypeName = map[AnnotationType]string{
	Text:           "Text",
	Link:           "Link",
	FreeText:       "FreeText",
	Line:           "Line",
	Square:         "Square",
	Circle:         "Circle",
	Polygon:        "Polygon",
	PolyLine:       "PolyLine",
	Highlight:      "Highlight",
	Underline:      "Underline",
	Squiggly:       "Squiggly",
	StrikeOut:      "StrikeOut",
	Stamp:          "Stamp",
	Caret:          "Caret",
	Ink:            "Ink",
	Popup:          "Popup",
	FileAttachment: "FileAttachment",
	Sound:          "Sound",
	Movie:          "Movie",
	Widget:         "Widget",
	Screen:         "Screen",
	PrinterMark:    "PrinterMark",
	TrapNet:        "TrapNet",
	Watermark:      "Watermark",
	Artwork3D:      "3D",
	Other:          "Other",
}
