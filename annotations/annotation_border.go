package annotations

import "fmt"

// AnnotationBorder represents a border for a PDF annotation object.
type AnnotationBorder struct {
	horizontalCornerRadius float64
	verticalCornerRadius   float64
	borderWidth            float64
	lineDashPattern        []float64
}

// Default is the default border style if not specified.
var Default = &AnnotationBorder{
	horizontalCornerRadius: 0,
	verticalCornerRadius:   0,
	borderWidth:            1,
	lineDashPattern:        nil,
}

// HorizontalCornerRadius returns the horizontal corner radius in user space units.
func (b AnnotationBorder) HorizontalCornerRadius() float64 {
	return b.horizontalCornerRadius
}

// VerticalCornerRadius returns the vertical corner radius in user space units.
func (b AnnotationBorder) VerticalCornerRadius() float64 {
	return b.verticalCornerRadius
}

// BorderWidth returns the width of the border in user space units.
func (b AnnotationBorder) BorderWidth() float64 {
	return b.borderWidth
}

// LineDashPattern returns the dash pattern for the border lines if provided. Returns nil if no dash pattern is set.
func (b AnnotationBorder) LineDashPattern() []float64 {
	return b.lineDashPattern
}

// NewAnnotationBorder creates a new AnnotationBorder.
// Pass nil for lineDashPattern if no dash pattern is needed.
func NewAnnotationBorder(horizontalCornerRadius, verticalCornerRadius, borderWidth float64, lineDashPattern []float64) *AnnotationBorder {
	return &AnnotationBorder{
		horizontalCornerRadius: horizontalCornerRadius,
		verticalCornerRadius:   verticalCornerRadius,
		borderWidth:            borderWidth,
		lineDashPattern:        lineDashPattern,
	}
}

// String returns a string representation of the annotation border.
func (b AnnotationBorder) String() string {
	return fmt.Sprintf("%g %g %g", b.horizontalCornerRadius, b.verticalCornerRadius, b.borderWidth)
}
