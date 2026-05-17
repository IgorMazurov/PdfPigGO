package page

// PageXmlTextRegion represents a pure text region in PAGE XML documents.
// This includes drop capitals, but practically ornate text may be considered as a graphic.
type PageXmlTextRegion struct {
	// AlternativeImages are alternative region images (e.g., black-and-white).
	AlternativeImages []AlternativeImage `xml:"AlternativeImage"`

	// Coords defines the polygon outline of the region.
	Coords *PageXmlCoords `xml:"Coords,omitempty"`

	// TextLines are the text line elements contained within this text region.
	TextLines []PageXmlTextLine `xml:"TextLine"`

	// TextEquivs store text content in various encodings and formats.
	TextEquivs []PageXmlTextEquiv `xml:"TextEquiv"`

	// TextStyle contains font family, size, colour, and typographic attributes.
	TextStyle *PageXmlTextStyle `xml:"TextStyle,omitempty"`

	// UserDefined holds structured custom data defined by name, type, and value.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// Labels contain semantic labels/tags with optional external model references.
	Labels []PageXmlLabels `xml:"Labels"`

	// Roles describe the roles this region takes in context of a parent region.
	Roles *PageXmlRoles `xml:"Roles,omitempty"`

	// Items are child regions nested within this region.
	Items []PageXmlRegionItem `xml:"AdvertRegion|ChartRegion|ChemRegion|CustomRegion|GraphicRegion|ImageRegion|LineDrawingRegion|MathsRegion|MusicRegion|NoiseRegion|SeparatorRegion|TableRegion|TextRegion|UnknownRegion"`

	// Id is the unique identifier for this region.
	Id string `xml:"id,attr,omitempty"`

	// Custom is a free-form attribute for custom data.
	Custom string `xml:"custom,attr,omitempty"`

	// Comments contains optional annotation text about this region.
	Comments string `xml:"comments,attr,omitempty"`

	// Continuation indicates whether this region is a continuation of another region
	// in a previous column or page.
	Continuation bool `xml:"continuation,attr,omitempty"`

	// Orientation is the angle the rectangle encapsulating the region has to be rotated
	// clockwise to correct skew (negative values indicate anti-clockwise rotation).
	// Range: -179.999 to 180.
	Orientation *float32 `xml:"orientation,attr,omitempty"`

	// Type is the nature of the text in the region.
	Type *PageXmlTextSimpleType `xml:"type,attr,omitempty"`

	// Leading is the degree of space in points between the lines of text (line spacing).
	Leading *int `xml:"leading,attr,omitempty"`

	// ReadingDirection specifies the direction in which text within lines should be read
	// (order of words and characters), in addition to TextLineOrder.
	ReadingDirection *PageXmlReadingDirectionSimpleType `xml:"readingDirection,attr,omitempty"`

	// TextLineOrder defines the order of text lines within the block,
	// in addition to ReadingDirection.
	TextLineOrder *PageXmlTextLineOrderSimpleType `xml:"textLineOrder,attr,omitempty"`

	// ReadingOrientation is the angle the baseline of text within the region has to be
	// rotated (relative to the rectangle encapsulating the region) clockwise to correct
	// skew, in addition to Orientation (negative values indicate anti-clockwise rotation).
	// Range: -179.999 to 180.
	ReadingOrientation *float32 `xml:"readingOrientation,attr,omitempty"`

	// Indented defines whether a region of text is indented or not.
	Indented *bool `xml:"indented,attr,omitempty"`

	// Align specifies the text alignment.
	Align *Align `xml:"align,attr,omitempty"`

	// PrimaryLanguage indicates the primary language used in the region.
	PrimaryLanguage *PageXmlLanguageSimpleType `xml:"primaryLanguage,attr,omitempty"`

	// SecondaryLanguage indicates the secondary language used in the region.
	SecondaryLanguage *PageXmlLanguageSimpleType `xml:"secondaryLanguage,attr,omitempty"`

	// PrimaryScript indicates the primary script used in the region.
	PrimaryScript *PageXmlScriptSimpleType `xml:"primaryScript,attr,omitempty"`

	// SecondaryScript indicates the secondary script used in the region.
	SecondaryScript *PageXmlScriptSimpleType `xml:"secondaryScript,attr,omitempty"`

	// Production indicates how the text was produced (printed, handwritten, etc.).
	Production *PageXmlProductionSimpleType `xml:"production,attr,omitempty"`
}

// regionItem implements PageXmlRegionItem.
func (*PageXmlTextRegion) regionItem() {}
