package alto

// AltoLayout represents the layout structure in ALTO format.
type AltoLayout struct {
	Pages     []*AltoPage `xml:"Page"`
	StyleRefs string      `xml:"STYLEREFS,attr"`
}
