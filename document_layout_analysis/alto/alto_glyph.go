package alto

// AltoGlyph represents a single character or ligature in ALTO format.
// Modern OCR software stores information on glyph level. Each glyph has its own
// coordinate information and must be separately addressable as a distinct object.
// The glyph elements are in order of the word. Each glyph needs to be recorded
// to build up the whole word sequence.
type AltoGlyph struct {
	AltoPositionedElement

	Shape    *AltoShape     `xml:"Shape"`
	Variant  []AltoVariant  `xml:"Variant"`
	Id       string         `xml:"ID,attr"`
	Content  string         `xml:"CONTENT,attr"`
	Gc       *float32       `xml:"GC,attr"`
}

// String returns the glyph content for display.
func (g AltoGlyph) String() string {
	return g.Content
}
