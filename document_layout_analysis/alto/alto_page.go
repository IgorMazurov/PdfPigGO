package alto

import (
	"encoding/xml"
	"math"
)

// AltoPage represents one page of a document in ALTO format.
type AltoPage struct {
	XMLName xml.Name `xml:"http://www.loc.gov/standards/alto/ns-v4# Page"`

	// TopMargin is the area between the top line of print and the upper edge of the leaf.
	TopMargin *AltoPageSpace `xml:"TopMargin,omitempty"`

	// LeftMargin is the area between the printspace and the left border of a page.
	LeftMargin *AltoPageSpace `xml:"LeftMargin,omitempty"`

	// RightMargin is the area between the printspace and the right border of a page.
	RightMargin *AltoPageSpace `xml:"RightMargin,omitempty"`

	// BottomMargin is the area between the bottom line of letterpress or writing and the bottom edge of the leaf.
	BottomMargin *AltoPageSpace `xml:"BottomMargin,omitempty"`

	// PrintSpace is the rectangle covering the printed area of a page.
	PrintSpace *AltoPageSpace `xml:"PrintSpace,omitempty"`

	// Id is the unique identifier for this page.
	Id string `xml:"ID,attr"`

	// PageClass is any user-defined class like title page.
	PageClass string `xml:"PAGECLASS,attr"`

	// StyleRefs references style definitions applied to this page.
	StyleRefs string `xml:"STYLEREFS,attr"`

	// ProcessingRefs references processing descriptions applied to this page.
	ProcessingRefs string `xml:"PROCESSINGREFS,attr"`

	// Height is the height of the page in measurement units defined by Description.
	Height *float32 `xml:"HEIGHT,attr"`

	// Width is the width of the page in measurement units defined by Description.
	Width *float32 `xml:"WIDTH,attr"`

	// PhysicalImgNr is the number of the page within the document.
	PhysicalImgNr float32 `xml:"PHYSICAL_IMG_NR,attr"`

	// PrintedImgNr is the page number that is printed on the page.
	PrintedImgNr string `xml:"PRINTED_IMG_NR,attr"`

	// Quality indicates the quality status of this page.
	Quality *AltoQuality `xml:"QUALITY,attr"`

	// QualityDetail provides additional information about the quality.
	QualityDetail string `xml:"QUALITY_DETAIL,attr"`

	// Position describes the position of the page (left, right, foldout, etc.).
	Position *AltoPosition `xml:"POSITION,attr"`

	// Processing is a link to the processing description used for this page.
	Processing string `xml:"PROCESSING,attr"`

	// Accuracy is the estimated percentage of OCR accuracy in range from 0 to 100.
	Accuracy *float32 `xml:"ACCURACY,attr"`

	// Pc indicates whether the page contains printed content.
	Pc *float32 `xml:"PC,attr"`
}

// NewAltoPage creates a new AltoPage with default values.
func NewAltoPage() *AltoPage {
	return &AltoPage{}
}

// SetHeight sets the Height field and tracks whether it has been specified.
// A nil value is set when v is NaN, matching C# semantics where HeightSpecified
// becomes true only for non-NaN values.
func (p *AltoPage) SetHeight(v float32) {
	if !math.IsNaN(float64(v)) {
		p.Height = &v
	} else {
		p.Height = nil
	}
}

// GetHeight returns the Height value and whether it has been set.
func (p *AltoPage) GetHeight() (float32, bool) {
	if p.Height != nil {
		return *p.Height, true
	}
	return 0, false
}

// SetWidth sets the Width field and tracks whether it has been specified.
// A nil value is set when v is NaN, matching C# semantics where WidthSpecified
// becomes true only for non-NaN values.
func (p *AltoPage) SetWidth(v float32) {
	if !math.IsNaN(float64(v)) {
		p.Width = &v
	} else {
		p.Width = nil
	}
}

// GetWidth returns the Width value and whether it has been set.
func (p *AltoPage) GetWidth() (float32, bool) {
	if p.Width != nil {
		return *p.Width, true
	}
	return 0, false
}

// SetAccuracy sets the Accuracy field and tracks whether it has been specified.
// A nil value is set when v is NaN, matching C# semantics where AccuracySpecified
// becomes true only for non-NaN values.
func (p *AltoPage) SetAccuracy(v float32) {
	if !math.IsNaN(float64(v)) {
		p.Accuracy = &v
	} else {
		p.Accuracy = nil
	}
}

// GetAccuracy returns the Accuracy value and whether it has been set.
func (p *AltoPage) GetAccuracy() (float32, bool) {
	if p.Accuracy != nil {
		return *p.Accuracy, true
	}
	return 0, false
}

// SetPc sets the Pc field and tracks whether it has been specified.
// A nil value is set when v is NaN, matching C# semantics where PcSpecified
// becomes true only for non-NaN values.
func (p *AltoPage) SetPc(v float32) {
	if !math.IsNaN(float64(v)) {
		p.Pc = &v
	} else {
		p.Pc = nil
	}
}

// GetPc returns the Pc value and whether it has been set.
func (p *AltoPage) GetPc() (float32, bool) {
	if p.Pc != nil {
		return *p.Pc, true
	}
	return 0, false
}
