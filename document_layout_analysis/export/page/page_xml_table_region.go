package page

// PageXmlTableRegion represents tabular data in PAGE XML documents.
// Rows and columns may or may not have separator lines; these lines are not separator regions.
type PageXmlTableRegion struct {
	// Grid is the table grid (visible or virtual grid lines).
	Grid []PageXmlGridPoints `xml:"Grid>GridPoints"`

	// Coords defines the polygon outline of the region.
	Coords *PageXmlCoords `xml:"Coords,omitempty"`

	// AlternativeImages are alternative region images (e.g., black-and-white).
	AlternativeImages []AlternativeImage `xml:"AlternativeImage"`

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

	// Orientation is the angle the rectangle encapsulating a region has to be rotated
	// clockwise to correct skew (negative values indicate anti-clockwise rotation).
	// Range: -179.999 to 180.
	Orientation *float32 `xml:"orientation,attr,omitempty"`

	// Rows is the number of rows present in the table.
	Rows *int `xml:"rows,attr,omitempty"`

	// Columns is the number of columns present in the table.
	Columns *int `xml:"columns,attr,omitempty"`

	// LineColour is the colour of the lines used in the region.
	LineColour *PageXmlColourSimpleType `xml:"lineColour,attr,omitempty"`

	// BgColour is the background colour of the region.
	BgColour *PageXmlColourSimpleType `xml:"bgColour,attr,omitempty"`

	// LineSeparators specifies the presence of line separators.
	LineSeparators *bool `xml:"lineSeparators,attr,omitempty"`

	// EmbText specifies whether the region also contains text.
	EmbText *bool `xml:"embText,attr,omitempty"`
}

// regionItem implements PageXmlRegionItem.
func (*PageXmlTableRegion) regionItem() {}
