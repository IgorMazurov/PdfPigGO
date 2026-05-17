package page

// Align represents text alignment in PAGE XML documents.
type Align byte

const (
	// AlignLeft indicates left-aligned text.
	AlignLeft Align = iota

	// AlignCentre indicates center-aligned text.
	AlignCentre

	// AlignRight indicates right-aligned text.
	AlignRight

	// AlignJustify indicates justified text.
	AlignJustify
)
