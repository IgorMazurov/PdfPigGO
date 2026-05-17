package page

import "time"

// PageXmlMetadataItem represents a single metadata entry in PAGE XML documents,
// containing semantic labels, type classification, name/value pairs, and optional dates.
type PageXmlMetadataItem struct {
	// Labels contains semantic label/tag elements associated with this metadata item.
	Labels []PageXmlLabels `xml:"Labels"`

	// Type indicates the kind of metadata (author, image properties, processing step, etc.).
	Type *PageXmlMetadataItemType `xml:"type,attr,omitempty"`

	// Name is a descriptive identifier (e.g., "imagePhotometricInterpretation").
	Name string `xml:"name,attr,omitempty"`

	// Value is the metadata value as a string (e.g., "RGB").
	Value string `xml:"value,attr,omitempty"`

	// Date is an optional timestamp associated with this metadata item.
	Date *time.Time `xml:"date,attr,omitempty"`
}
