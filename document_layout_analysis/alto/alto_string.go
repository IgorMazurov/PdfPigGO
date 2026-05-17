package alto

// AltoString represents a sequence of characters in ALTO format.
// Strings are separated by white spaces or hyphenation chars.
type AltoString struct {
	AltoPositionedElement

	Shape            *AltoShape       `xml:"Shape"`
	Alternative      []AltoAlternative `xml:"ALTERNATIVE"`
	Glyphs           []AltoGlyph     `xml:"Glyph"`
	Id               string          `xml:"ID,attr"`
	StyleRefs        string          `xml:"STYLEREFS,attr"`
	TagRefs          string          `xml:"TAGREFS,attr"`
	ProcessingRefs   string          `xml:"PROCESSINGREFS,attr"`
	Content          string          `xml:"CONTENT,attr"`
	Style            *AltoFontStyles `xml:"STYLE,attr"`
	SubsType         *AltoSubsType   `xml:"SUBS_TYPE,attr"`
	SubsContent      string          `xml:"SUBS_CONTENT,attr"`
	Wc               *float32        `xml:"WC,attr"`
	Cc               string          `xml:"CC,attr"`
	CorrectionStatus *bool           `xml:"CS,attr"`
	Language         string          `xml:"LANG,attr"`
}

// String returns the string content for display.
func (s AltoString) String() string {
	return s.Content
}
