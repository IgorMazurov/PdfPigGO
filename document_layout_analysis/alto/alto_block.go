package alto

import "math"

// AltoBlock is the base type for any kind of block on the page.
type AltoBlock struct {
	AltoPositionedElement

	Shape *AltoShape `xml:"Shape"`

	Id string `xml:"ID,attr,omitempty"`

	StyleRefs string `xml:"STYLEREFS,attr,omitempty"`

	TagRefs string `xml:"TAGREFS,attr,omitempty"`

	ProcessingRefs string `xml:"PROCESSINGREFS,attr,omitempty"`

	Rotation *float32 `xml:"ROTATION,attr,omitempty"`

	IdNext string `xml:"IDNEXT,attr,omitempty"`

	CorrectionStatus *bool `xml:"CS,attr,omitempty"`

	TypeXlink string

	Href string

	Role string

	Arcrole string

	Title string

	Show *AltoBlockTypeShow

	Actuate *AltoBlockTypeActuate
}

// NewAltoBlock creates a new AltoBlock with default values.
func NewAltoBlock() *AltoBlock {
	return &AltoBlock{}
}

// GetRotation returns the rotation value or 0 if not set.
// The rotation is in degrees counterclockwise.
func (b *AltoBlock) GetRotation() float32 {
	if b.Rotation != nil {
		return *b.Rotation
	}
	return 0
}

// SetRotation sets the rotation value and marks it as specified.
// The rotation is in degrees counterclockwise.
func (b *AltoBlock) SetRotation(value float32) {
	if !math.IsNaN(float64(value)) {
		b.Rotation = &value
	} else {
		b.Rotation = nil
	}
}

// HasRotation returns true if Rotation has been explicitly set.
func (b *AltoBlock) HasRotation() bool {
	return b.Rotation != nil
}

// GetCorrectionStatus returns the correction status value or false if not set.
// Correction status indicates whether manual correction has been done.
func (b *AltoBlock) GetCorrectionStatus() bool {
	if b.CorrectionStatus != nil {
		return *b.CorrectionStatus
	}
	return false
}

// SetCorrectionStatus sets the correction status and marks it as specified.
func (b *AltoBlock) SetCorrectionStatus(value bool) {
	b.CorrectionStatus = &value
}

// HasCorrectionStatus returns true if CorrectionStatus has been explicitly set.
func (b *AltoBlock) HasCorrectionStatus() bool {
	return b.CorrectionStatus != nil
}

// GetShow returns the show value or empty string if not set.
func (b *AltoBlock) GetShow() AltoBlockTypeShow {
	if b.Show != nil {
		return *b.Show
	}
	return ""
}

// SetShow sets the show value and marks it as specified.
func (b *AltoBlock) SetShow(value AltoBlockTypeShow) {
	b.Show = &value
}

// HasShow returns true if Show has been explicitly set.
func (b *AltoBlock) HasShow() bool {
	return b.Show != nil
}

// GetActuate returns the actuate value or empty string if not set.
func (b *AltoBlock) GetActuate() AltoBlockTypeActuate {
	if b.Actuate != nil {
		return *b.Actuate
	}
	return ""
}

// SetActuate sets the actuate value and marks it as specified.
func (b *AltoBlock) SetActuate(value AltoBlockTypeActuate) {
	b.Actuate = &value
}

// HasActuate returns true if Actuate has been explicitly set.
func (b *AltoBlock) HasActuate() bool {
	return b.Actuate != nil
}
