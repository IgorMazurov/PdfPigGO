package alto

// AltoVariant represents an alternative (combined) character for a glyph,
// outlined by OCR engine or similar recognition processes. In case the variant
// are two combining characters, two characters are outlined in one Variant element.
// E.g. a Glyph with CONTENT="m" can have a Variant with content "rn".
type AltoVariant struct {
	Content string  `xml:"CONTENT,attr"`
	Vc      *float32 `xml:"VC,attr"`
}
