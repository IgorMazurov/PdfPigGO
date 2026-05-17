package alto

import "strings"

// AltoTextBlockTextLine represents a single line of text in ALTO format.
type AltoTextBlockTextLine struct {
	AltoPositionedElement

	Shape            *AltoShape               `xml:"Shape"`
	Strings          []AltoString             `xml:"String"`
	Sp               []AltoSP                 `xml:"SP"`
	Hyp              *AltoTextBlockTextLineHyp `xml:"HYP"`
	Id               string                   `xml:"ID,attr"`
	StyleRefs        string                   `xml:"STYLEREFS,attr"`
	TagRefs          string                   `xml:"TAGREFS,attr"`
	ProcessingRefs   string                   `xml:"PROCESSINGREFS,attr"`
	BaseLine         *float32                 `xml:"BASELINE,attr,omitempty"`
	Language         string                   `xml:"LANG,attr"`
	CorrectionStatus *bool                    `xml:"CS,attr,omitempty"`
}

// String returns the text content of all strings in the line joined by spaces.
func (l *AltoTextBlockTextLine) String() string {
	parts := make([]string, 0, len(l.Strings))
	for _, s := range l.Strings {
		parts = append(parts, s.String())
	}
	return strings.Join(parts, " ")
}
