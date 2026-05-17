package page

// PageXmlTextStyle represents text style information in PAGE XML documents,
// including font family, size, colour, and various typographic attributes.
type PageXmlTextStyle struct {
	// FontFamily is the name of the font (e.g., Arial, Times New Roman).
	FontFamily string `xml:"fontFamily,attr"`

	// Serif indicates whether the typeface is serif or sans-serif.
	Serif *bool `xml:"serif,attr,omitempty"`

	// Monospace indicates monospace (fixed-pitch) or proportional font.
	Monospace *bool `xml:"monospace,attr,omitempty"`

	// FontSize is the size of the characters in points.
	FontSize *float32 `xml:"fontSize,attr,omitempty"`

	// XHeight is the distance between the baseline and the mean line of lower-case letters in pixels.
	XHeight string `xml:"xHeight,attr"`

	// Kerning is the degree of space (in points) between characters.
	Kerning *int `xml:"kerning,attr,omitempty"`

	// TextColour is the text colour value.
	TextColour *PageXmlColourSimpleType `xml:"textColour,attr,omitempty"`

	// TextColourRgb is the text colour in RGB encoded format: (red) + (256 × green) + (65536 × blue).
	TextColourRgb string `xml:"textColourRgb,attr"`

	// BgColour is the background colour value.
	BgColour *PageXmlColourSimpleType `xml:"bgColour,attr,omitempty"`

	// BgColourRgb is the background colour in RGB encoded format: (red) + (256 × green) + (65536 × blue).
	BgColourRgb string `xml:"bgColourRgb,attr"`

	// ReverseVideo specifies whether the text colour appears reversed against a background colour.
	ReverseVideo *bool `xml:"reverseVideo,attr,omitempty"`

	// Bold indicates bold typeface.
	Bold *bool `xml:"bold,attr,omitempty"`

	// Italic indicates italic typeface.
	Italic *bool `xml:"italic,attr,omitempty"`

	// Underlined indicates whether text is underlined.
	Underlined *bool `xml:"underlined,attr,omitempty"`

	// UnderlineStyle specifies line style details when Underlined is true.
	UnderlineStyle *PageXmlUnderlineStyleSimpleType `xml:"underlineStyle,attr,omitempty"`

	// Subscript indicates subscript text positioning.
	Subscript *bool `xml:"subscript,attr,omitempty"`

	// Superscript indicates superscript text positioning.
	Superscript *bool `xml:"superscript,attr,omitempty"`

	// Strikethrough indicates strikethrough text decoration.
	Strikethrough *bool `xml:"strikethrough,attr,omitempty"`

	// SmallCaps indicates small capitals typeface.
	SmallCaps *bool `xml:"smallCaps,attr,omitempty"`

	// LetterSpaced indicates letter-spaced text.
	LetterSpaced *bool `xml:"letterSpaced,attr,omitempty"`
}
