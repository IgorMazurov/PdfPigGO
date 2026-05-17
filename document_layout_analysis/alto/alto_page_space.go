package alto

// AltoPageSpace represents a region on an ALTO page.
type AltoPageSpace struct {
	AltoPositionedElement
	Shape           *AltoShape            `xml:"Shape"`
	TextBlocks      []AltoTextBlock       `xml:"TextBlock"`
	Illustrations   []AltoIllustration    `xml:"Illustration"`
	GraphicalElements []AltoGraphicalElement `xml:"GraphicalElement"`
	ComposedBlocks  []AltoComposedBlock   `xml:"ComposedBlock"`
	Id              string                `xml:"ID,attr"`
	StyleRefs       string                `xml:"STYLEREFS,attr"`
	ProcessingRefs  string                `xml:"PROCESSINGREFS,attr"`
}
