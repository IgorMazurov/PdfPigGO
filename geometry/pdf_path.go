package geometry

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
)

// StrokeState exposes the stroke-related fields needed by PdfPath.SetStrokeDetails.
type StrokeState struct {
	LineDashPattern       graphiccore.LineDashPattern
	CurrentStrokingColor  colors.Color
	LineWidth             float64
	CapStyle              graphiccore.LineCapStyle
	JoinStyle             graphiccore.LineJoinStyle
}

// FillState exposes the fill-related fields needed by PdfPath.SetFillDetails.
type FillState struct {
	CurrentNonStrokingColor colors.Color
}

// PdfPath represents a PDF drawing path composed of one or more disconnected subpaths.
type PdfPath struct {
	subpaths        []*core.PdfSubpath
	fillingRule     core.FillingRule
	isClipping      bool
	isFilled        bool
	isStroked       bool
	lineWidth       float64
	fillColor       colors.Color
	strokeColor     colors.Color
	lineDashPattern graphiccore.LineDashPattern
	lineCapStyle    graphiccore.LineCapStyle
	lineJoinStyle   graphiccore.LineJoinStyle
}

// NewPdfPath creates a new empty PdfPath.
func NewPdfPath() *PdfPath {
	return &PdfPath{
		subpaths: make([]*core.PdfSubpath, 0),
	}
}

// Len returns the number of subpaths in this path.
func (p *PdfPath) Len() int {
	return len(p.subpaths)
}

// Subpaths returns the subpaths that make up this path.
func (p *PdfPath) Subpaths() []*core.PdfSubpath {
	return p.subpaths
}

// IsClipping reports whether this is a clipping path.
func (p *PdfPath) IsClipping() bool {
	return p.isClipping
}

// IsFilled reports whether the path is filled.
func (p *PdfPath) IsFilled() bool {
	return p.isFilled
}

// IsStroked reports whether the path is stroked.
func (p *PdfPath) IsStroked() bool {
	return p.isStroked
}

// FillingRule returns the filling rule for this path.
func (p *PdfPath) FillingRule() core.FillingRule {
	return p.fillingRule
}

// LineWidth returns the stroke line width in user space units.
func (p *PdfPath) LineWidth() float64 {
	return p.lineWidth
}

// FillColor returns the fill color for this path.
func (p *PdfPath) FillColor() colors.Color {
	return p.fillColor
}

// SetFillColor sets the fill color for this path.
func (p *PdfPath) SetFillColor(c colors.Color) {
	p.fillColor = c
}

// StrokeColor returns the stroke color for this path.
func (p *PdfPath) StrokeColor() colors.Color {
	return p.strokeColor
}

// SetStrokeColor sets the stroke color for this path.
func (p *PdfPath) SetStrokeColor(c colors.Color) {
	p.strokeColor = c
}

// LineDashPattern returns the line dash pattern for stroking this path.
func (p *PdfPath) LineDashPattern() graphiccore.LineDashPattern {
	return p.lineDashPattern
}

// SetLineDashPattern sets the line dash pattern for stroking this path.
func (p *PdfPath) SetLineDashPattern(pattern graphiccore.LineDashPattern) {
	p.lineDashPattern = pattern
}

// LineCapStyle returns the line cap style for stroking this path.
func (p *PdfPath) LineCapStyle() graphiccore.LineCapStyle {
	return p.lineCapStyle
}

// SetLineCapStyle sets the line cap style for stroking this path.
func (p *PdfPath) SetLineCapStyle(style graphiccore.LineCapStyle) {
	p.lineCapStyle = style
}

// LineJoinStyle returns the line join style for stroking this path.
func (p *PdfPath) LineJoinStyle() graphiccore.LineJoinStyle {
	return p.lineJoinStyle
}

// SetLineJoinStyle sets the line join style for stroking this path.
func (p *PdfPath) SetLineJoinStyle(style graphiccore.LineJoinStyle) {
	p.lineJoinStyle = style
}

// SetClipping sets the clipping mode and filling rule for this path.
func (p *PdfPath) SetClipping(rule core.FillingRule) {
	p.isFilled = false
	p.isStroked = false
	p.isClipping = true
	p.fillingRule = rule
}

// SetFilled marks the path as filled with the given filling rule.
func (p *PdfPath) SetFilled(rule core.FillingRule) {
	p.isFilled = true
	p.fillingRule = rule
}

// SetStroked marks the path as stroked.
func (p *PdfPath) SetStroked() {
	p.isStroked = true
}

// SetLineWidth sets the stroke line width.
func (p *PdfPath) SetLineWidth(width float64) {
	p.lineWidth = width
}

// Add appends a subpath to this path.
func (p *PdfPath) Add(subpath *core.PdfSubpath) {
	p.subpaths = append(p.subpaths, subpath)
}

// CloneEmpty creates a clone with no subpaths but the same rendering settings.
func (p *PdfPath) CloneEmpty() *PdfPath {
	newPath := NewPdfPath()
	if p.isClipping {
		newPath.SetClipping(p.fillingRule)
	} else {
		if p.isFilled {
			newPath.SetFilled(p.fillingRule)
			newPath.fillColor = p.fillColor
		}
		if p.isStroked {
			newPath.SetStroked()
			newPath.lineWidth = p.lineWidth
			newPath.lineCapStyle = p.lineCapStyle
			newPath.lineDashPattern = p.lineDashPattern
			newPath.lineJoinStyle = p.lineJoinStyle
			newPath.strokeColor = p.strokeColor
		}
	}
	return newPath
}

// GetBoundingRectangle returns the bounding rectangle of the path, or nil if empty.
func (p *PdfPath) GetBoundingRectangle() *core.PdfRectangle {
	return core.GetBoundingRectangleForPath(p.subpaths)
}

// Clip clips this path against the given clipping path and returns the result.
// Returns nil if the clipped result is empty or an error occurs.
func (p *PdfPath) Clip(clipping *PdfPath, log logging.Log) *PdfPath {
	if p == nil || clipping == nil {
		return nil
	}
	result, _ := Clip(clipping, p, log)
	return result
}

// SetStrokeDetails captures stroke rendering details from the current graphics state.
func (p *PdfPath) SetStrokeDetails(state StrokeState) {
	p.lineDashPattern = state.LineDashPattern
	p.strokeColor = state.CurrentStrokingColor
	p.lineWidth = state.LineWidth
	p.lineCapStyle = state.CapStyle
	p.lineJoinStyle = state.JoinStyle
}

// SetFillDetails captures fill rendering details from the current graphics state.
func (p *PdfPath) SetFillDetails(state FillState) {
	p.fillColor = state.CurrentNonStrokingColor
}
