package alto

// AltoTextStyle defines font properties of text in ALTO format.
type AltoTextStyle struct {
	// Id is the unique identifier for this text style.
	Id string `xml:"ID,attr"`

	// FontFamily is the font name.
	FontFamily string `xml:"FONTFAMILY,attr,omitempty"`

	// FontType is the type of font (serif or sans-serif).
	// Nil means not specified.
	FontType *AltoFontType `xml:"FONTTYPE,attr,omitempty"`

	// FontWidth is the width characteristic of the font.
	// Nil means not specified.
	FontWidth *AltoFontWidth `xml:"FONTWIDTH,attr,omitempty"`

	// FontSize is the font size in points (1/72 of an inch).
	FontSize float32 `xml:"FONTSIZE,attr"`

	// FontColor is the font color as an RGB value.
	FontColor []byte `xml:"FONTCOLOR,attr,omitempty"`

	// FontStyle describes the style of the font (bold, italics, etc.).
	// Nil means not specified.
	FontStyle *AltoFontStyles `xml:"FONTSTYLE,attr,omitempty"`
}
