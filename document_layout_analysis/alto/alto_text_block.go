package alto

import "strings"

// AltoTextBlock represents a block of text in ALTO format.
type AltoTextBlock struct {
	AltoBlock

	TextLines []AltoTextBlockTextLine `xml:"TextLine"`
	Language  string                  `xml:"language,attr,omitempty"`
	Lang      string                  `xml:"LANG,attr,omitempty"`
}

// NewAltoTextBlock creates a new AltoTextBlock with default values.
func NewAltoTextBlock() *AltoTextBlock {
	return &AltoTextBlock{}
}

// String returns the text content of all lines joined by spaces.
func (b *AltoTextBlock) String() string {
	parts := make([]string, 0, len(b.TextLines))
	for _, line := range b.TextLines {
		parts = append(parts, line.String())
	}
	return strings.Join(parts, " ")
}
