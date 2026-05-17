package page

// PageXmlPage represents a single document page within a PAGE XML document,
// containing image metadata, regions, reading order, layers, relations, and other
// page-level properties according to the PAGE scheme (2019-07-15).
type PageXmlPage struct {
	// AlternativeImages are alternative document page images (e.g., black-and-white).
	AlternativeImages []AlternativeImage `xml:"AlternativeImage"`

	// Border defines the polygon outline of the actual page.
	Border *PageXmlBorder `xml:"Border,omitempty"`

	// PrintSpace determines the effective area on the paper of a printed page.
	PrintSpace *PageXmlPrintSpace `xml:"PrintSpace,omitempty"`

	// ReadingOrder specifies the order of blocks within the page.
	ReadingOrder *PageXmlReadingOrder `xml:"ReadingOrder,omitempty"`

	// Layers contain unassigned regions considered to be in the virtual default layer,
	// treated as below any other layers.
	Layers *PageXmlLayers `xml:"Layers,omitempty"`

	// Relations define relationships between elements on the page.
	Relations *PageXmlRelations `xml:"Relations,omitempty"`

	// TextStyle contains the default text style for the page.
	TextStyle *PageXmlTextStyle `xml:"TextStyle,omitempty"`

	// UserDefined holds structured custom data defined by name, type, and value.
	UserDefined []PageXmlUserAttribute `xml:"UserAttribute"`

	// Labels contain semantic labels/tags with optional external model references.
	Labels []PageXmlLabels `xml:"Labels"`

	// Items are the top-level regions on this page (text, table, image, etc.).
	Items []PageXmlRegionItem `xml:"AdvertRegion|ChartRegion|ChemRegion|CustomRegion|GraphicRegion|ImageRegion|LineDrawingRegion|MapRegion|MathsRegion|MusicRegion|NoiseRegion|SeparatorRegion|TableRegion|TextRegion|UnknownRegion"`

	// ImageFilename contains the image file name including the file extension.
	ImageFilename string `xml:"imageFilename,attr,omitempty"`

	// ImageWidth specifies the width of the image in pixels.
	ImageWidth int `xml:"imageWidth,attr,omitempty"`

	// ImageHeight specifies the height of the image in pixels.
	ImageHeight int `xml:"imageHeight,attr,omitempty"`

	// ImageXResolution specifies the image resolution in width.
	ImageXResolution *float32 `xml:"imageXResolution,attr,omitempty"`

	// ImageYResolution specifies the image resolution in height.
	ImageYResolution *float32 `xml:"imageYResolution,attr,omitempty"`

	// ImageResolutionUnit specifies the unit of the resolution information
	// referring to a standardised unit of measurement (pixels per inch, pixels
	// per centimeter or other).
	ImageResolutionUnit *PageXmlPageImageResolutionUnit `xml:"imageResolutionUnit,attr,omitempty"`

	// Custom is a free-form attribute for generic use.
	Custom string `xml:"custom,attr,omitempty"`

	// Orientation is the angle the rectangle encapsulating the page (or its Border)
	// has to be rotated in clockwise direction in order to correct the present skew
	// (negative values indicate anti-clockwise rotation). The rotated image can be
	// further referenced via AlternativeImage. Range: -179.999 to 180.
	Orientation *float32 `xml:"orientation,attr,omitempty"`

	// Type is the type of the page within the document (e.g., cover page).
	Type *PageXmlPageSimpleType `xml:"type,attr,omitempty"`

	// PrimaryLanguage indicates the primary language used in the page; lower-level
	// definitions override the page-level definition.
	PrimaryLanguage *PageXmlLanguageSimpleType `xml:"primaryLanguage,attr,omitempty"`

	// SecondaryLanguage indicates the secondary language used in the page; lower-level
	// definitions override the page-level definition.
	SecondaryLanguage *PageXmlLanguageSimpleType `xml:"secondaryLanguage,attr,omitempty"`

	// PrimaryScript indicates the primary script used in the page; lower-level
	// definitions override the page-level definition.
	PrimaryScript *PageXmlScriptSimpleType `xml:"primaryScript,attr,omitempty"`

	// SecondaryScript indicates the secondary script used in the page; lower-level
	// definitions override the page-level definition.
	SecondaryScript *PageXmlScriptSimpleType `xml:"secondaryScript,attr,omitempty"`

	// ReadingDirection specifies the direction in which text within lines should be
	// read (order of words and characters), in addition to TextLineOrder; lower-level
	// definitions override the page-level definition.
	ReadingDirection *PageXmlReadingDirectionSimpleType `xml:"readingDirection,attr,omitempty"`

	// TextLineOrder defines the order of text lines within a block, in addition to
	// ReadingDirection; lower-level definitions override the page-level definition.
	TextLineOrder *PageXmlTextLineOrderSimpleType `xml:"textLineOrder,attr,omitempty"`

	// Conf is a confidence value for the whole page between 0 and 1.
	Conf *float32 `xml:"conf,attr,omitempty"`
}
