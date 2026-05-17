package alto

// AltoTextBlockTextLineHyp represents a hyphenation character at the end of a text line.
type AltoTextBlockTextLineHyp struct {
	AltoPositionedElement
	Content string `xml:"CONTENT,attr"`
}
