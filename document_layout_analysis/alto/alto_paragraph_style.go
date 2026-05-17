package alto

// AltoParagraphStyle defines formatting properties of text blocks in ALTO format.
type AltoParagraphStyle struct {
	// Id is the unique identifier for this paragraph style.
	Id string `xml:"ID,attr"`

	// Align indicates the alignment of the paragraph (left, right, center, or justify).
	Align AltoParagraphStyleAlign `xml:"ALIGN,attr,omitempty"`

	// Left is the left indent of the paragraph in relation to the column.
	// Nil means not specified.
	Left *float32 `xml:"LEFT,attr,omitempty"`

	// Right is the right indent of the paragraph in relation to the column.
	// Nil means not specified.
	Right *float32 `xml:"RIGHT,attr,omitempty"`

	// LineSpace is the line spacing between two lines of the paragraph,
	// measured from baseline to baseline.
	// Nil means not specified.
	LineSpace *float32 `xml:"LINESPACE,attr,omitempty"`

	// FirstLine is the indent of the first line of the paragraph if different
	// from other lines. A negative value indicates an indent to the left,
	// a positive value indicates an indent to the right.
	// Nil means not specified.
	FirstLine *float32 `xml:"FIRSTLINE,attr,omitempty"`
}
