package graphics

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

// FontState holds the text-related graphics state parameters.
type FontState struct {
	FontName                    *tokens.NameToken
	FontSize                    float64
	HorizontalScaling           float64
	CharacterSpacing            float64
	WordSpacing                 float64
	Leading                     float64
	TextRenderingMode           core.TextRenderingMode
	Rise                        float64
	Knockout                    bool
	FromExtendedGraphicsState   bool
}

// DeepClone creates a deep copy of the FontState.
func (f *FontState) DeepClone() *FontState {
	return &FontState{
		FontName:                    f.FontName,
		FontSize:                    f.FontSize,
		HorizontalScaling:           f.HorizontalScaling,
		CharacterSpacing:            f.CharacterSpacing,
		WordSpacing:                 f.WordSpacing,
		Leading:                     f.Leading,
		TextRenderingMode:           f.TextRenderingMode,
		Rise:                        f.Rise,
		Knockout:                    f.Knockout,
		FromExtendedGraphicsState:   f.FromExtendedGraphicsState,
	}
}

// CurrentGraphicsState holds the complete graphics state for PDF rendering.
type CurrentGraphicsState struct {
	CurrentTransformationMatrix core.TransformationMatrix
	CurrentClippingPath         *geometry.PdfPath
	ColorSpaceContext           *ColorSpaceContext
	FontState                   FontState
	LineWidth                   float64
	CapStyle                    graphiccore.LineCapStyle
	JoinStyle                   graphiccore.LineJoinStyle
	MiterLimit                  float64
	LineDashPattern             graphiccore.LineDashPattern
	Flatness                    float64
	BlendMode                   graphiccore.BlendMode
	SoftMask                    *SoftMask
	AlphaSource                 bool
	AlphaConstantStroking       float64
	AlphaConstantNonStroking    float64
	Overprint                   bool
	NonStrokingOverprint        bool
	OverprintMode               int
	StrokeAdjustment            bool
	Smoothness                  float64
	RenderingIntent             graphiccore.RenderingIntent
	CurrentStrokingColor        colors.Color
	CurrentNonStrokingColor     colors.Color
}

// DeepClone creates a deep copy of the CurrentGraphicsState.
func (s *CurrentGraphicsState) DeepClone() *CurrentGraphicsState {
	clone := &CurrentGraphicsState{
		CurrentTransformationMatrix: s.CurrentTransformationMatrix,
		FontState:                   *s.FontState.DeepClone(),
		LineWidth:                   s.LineWidth,
		CapStyle:                    s.CapStyle,
		JoinStyle:                   s.JoinStyle,
		MiterLimit:                  s.MiterLimit,
		LineDashPattern:             s.LineDashPattern,
		Flatness:                    s.Flatness,
		BlendMode:                   s.BlendMode,
		AlphaSource:                 s.AlphaSource,
		AlphaConstantStroking:       s.AlphaConstantStroking,
		AlphaConstantNonStroking:    s.AlphaConstantNonStroking,
		Overprint:                   s.Overprint,
		NonStrokingOverprint:        s.NonStrokingOverprint,
		OverprintMode:               s.OverprintMode,
StrokeAdjustment:            s.StrokeAdjustment,
	Smoothness:                  s.Smoothness,
	RenderingIntent:             s.RenderingIntent,
		CurrentStrokingColor:        s.CurrentStrokingColor,
		CurrentNonStrokingColor:     s.CurrentNonStrokingColor,
	}

	if s.CurrentClippingPath != nil {
		cloned := geometry.NewPdfPath()
		for _, subpath := range s.CurrentClippingPath.Subpaths() {
			cloned.Add(subpath)
		}
		if s.CurrentClippingPath.IsClipping() {
			cloned.SetClipping(s.CurrentClippingPath.FillingRule())
		} else if s.CurrentClippingPath.IsFilled() {
			cloned.SetFilled(s.CurrentClippingPath.FillingRule())
			cloned.SetFillColor(s.CurrentClippingPath.FillColor())
		}
		if s.CurrentClippingPath.IsStroked() {
			cloned.SetStroked()
			cloned.SetLineWidth(s.CurrentClippingPath.LineWidth())
			cloned.SetStrokeColor(s.CurrentClippingPath.StrokeColor())
			cloned.SetLineDashPattern(s.CurrentClippingPath.LineDashPattern())
			cloned.SetLineCapStyle(s.CurrentClippingPath.LineCapStyle())
			cloned.SetLineJoinStyle(s.CurrentClippingPath.LineJoinStyle())
		}
		clone.CurrentClippingPath = cloned
	}

	if s.ColorSpaceContext != nil {
		clone.ColorSpaceContext = s.ColorSpaceContext.DeepClone()
	}

	if s.SoftMask != nil {
		sm := *s.SoftMask
		clone.SoftMask = &sm
	}

	return clone
}

// ColorSpaceContext manages the current color space state.
type ColorSpaceContext struct {
	CurrentStrokingColorSpace   colors.ColorSpaceDetails
	CurrentNonStrokingColorSpace colors.ColorSpaceDetails
	getCurrentState             func() *CurrentGraphicsState
	resourceStore               content.ResourceStore
}

// NewColorSpaceContext creates a new ColorSpaceContext.
func NewColorSpaceContext(getCurrentState func() *CurrentGraphicsState, resourceStore content.ResourceStore) *ColorSpaceContext {
	return &ColorSpaceContext{
		CurrentStrokingColorSpace:  colors.DeviceGrayColorSpaceDetails,
		CurrentNonStrokingColorSpace: colors.DeviceGrayColorSpaceDetails,
		getCurrentState:            getCurrentState,
		resourceStore:              resourceStore,
	}
}

// DeepClone creates a deep copy of the ColorSpaceContext.
func (c *ColorSpaceContext) DeepClone() *ColorSpaceContext {
	if c == nil {
		return nil
	}
	return &ColorSpaceContext{
		CurrentStrokingColorSpace:  c.CurrentStrokingColorSpace,
		CurrentNonStrokingColorSpace: c.CurrentNonStrokingColorSpace,
		getCurrentState:            c.getCurrentState,
		resourceStore:              c.resourceStore,
	}
}

// SetStrokingColorspace sets the current stroking color space.
func (c *ColorSpaceContext) SetStrokingColorspace(colorspace *tokens.NameToken, dictionary *tokens.DictionaryToken) {
	c.CurrentStrokingColorSpace = c.resourceStore.GetColorSpaceDetails(colorspace, dictionary)
	if colors.IsUnsupported(c.CurrentStrokingColorSpace) {
		return
	}
	state := c.getCurrentState()
	if state != nil {
		state.CurrentStrokingColor = c.CurrentStrokingColorSpace.GetInitializeColor()
	}
}

// SetStrokingColor sets the current stroking color using operands and optional pattern name.
func (c *ColorSpaceContext) SetStrokingColor(operands []float64, patternName *tokens.NameToken) {
	if colors.IsUnsupported(c.CurrentStrokingColorSpace) {
		return
	}
	state := c.getCurrentState()
	if state == nil {
		return
	}

	if patternName != nil && c.CurrentStrokingColorSpace.Type() == colors.Pattern {
		if pc, ok := c.CurrentStrokingColorSpace.(interface{ GetColorByName(*tokens.NameToken) colors.Color }); ok {
			state.CurrentStrokingColor = pc.GetColorByName(patternName)
			return
		}
	}

	state.CurrentStrokingColor = c.CurrentStrokingColorSpace.GetColor(operands...)
}

// SetStrokingColorGray sets the current stroking color to a gray value.
func (c *ColorSpaceContext) SetStrokingColorGray(gray float64) {
	c.CurrentStrokingColorSpace = colors.DeviceGrayColorSpaceDetails
	state := c.getCurrentState()
	if state != nil {
		state.CurrentStrokingColor = c.CurrentStrokingColorSpace.GetColor(gray)
	}
}

// SetStrokingColorRgb sets the current stroking color to an RGB value.
func (c *ColorSpaceContext) SetStrokingColorRgb(r, g, b float64) {
	c.CurrentStrokingColorSpace = colors.DeviceRgbColorSpaceDetails
	state := c.getCurrentState()
	if state != nil {
		state.CurrentStrokingColor = c.CurrentStrokingColorSpace.GetColor(r, g, b)
	}
}

// SetStrokingColorCmyk sets the current stroking color to a CMYK value.
func (c *ColorSpaceContext) SetStrokingColorCmyk(cyan, magenta, yellow, black float64) {
	c.CurrentStrokingColorSpace = colors.DeviceCmykColorSpaceDetails
	state := c.getCurrentState()
	if state != nil {
		state.CurrentStrokingColor = c.CurrentStrokingColorSpace.GetColor(cyan, magenta, yellow, black)
	}
}

// SetNonStrokingColorspace sets the current non-stroking color space.
func (c *ColorSpaceContext) SetNonStrokingColorspace(colorspace *tokens.NameToken, dictionary *tokens.DictionaryToken) {
	c.CurrentNonStrokingColorSpace = c.resourceStore.GetColorSpaceDetails(colorspace, dictionary)
	if colors.IsUnsupported(c.CurrentNonStrokingColorSpace) {
		return
	}
	state := c.getCurrentState()
	if state != nil {
		state.CurrentNonStrokingColor = c.CurrentNonStrokingColorSpace.GetInitializeColor()
	}
}

// SetNonStrokingColor sets the current non-stroking color using operands and optional pattern name.
func (c *ColorSpaceContext) SetNonStrokingColor(operands []float64, patternName *tokens.NameToken) {
	if colors.IsUnsupported(c.CurrentNonStrokingColorSpace) {
		return
	}
	state := c.getCurrentState()
	if state == nil {
		return
	}

	if patternName != nil && c.CurrentNonStrokingColorSpace.Type() == colors.Pattern {
		if pc, ok := c.CurrentNonStrokingColorSpace.(interface{ GetColorByName(*tokens.NameToken) colors.Color }); ok {
			state.CurrentNonStrokingColor = pc.GetColorByName(patternName)
			return
		}
	}

	state.CurrentNonStrokingColor = c.CurrentNonStrokingColorSpace.GetColor(operands...)
}

// SetNonStrokingColorGray sets the current non-stroking color to a gray value.
func (c *ColorSpaceContext) SetNonStrokingColorGray(gray float64) {
	c.CurrentNonStrokingColorSpace = colors.DeviceGrayColorSpaceDetails
	state := c.getCurrentState()
	if state != nil {
		state.CurrentNonStrokingColor = c.CurrentNonStrokingColorSpace.GetColor(gray)
	}
}

// SetNonStrokingColorRgb sets the current non-stroking color to an RGB value.
func (c *ColorSpaceContext) SetNonStrokingColorRgb(r, g, b float64) {
	c.CurrentNonStrokingColorSpace = colors.DeviceRgbColorSpaceDetails
	state := c.getCurrentState()
	if state != nil {
		state.CurrentNonStrokingColor = c.CurrentNonStrokingColorSpace.GetColor(r, g, b)
	}
}

// SetNonStrokingColorCmyk sets the current non-stroking color to a CMYK value.
func (c *ColorSpaceContext) SetNonStrokingColorCmyk(cyan, magenta, yellow, black float64) {
	c.CurrentNonStrokingColorSpace = colors.DeviceCmykColorSpaceDetails
	state := c.getCurrentState()
	if state != nil {
		state.CurrentNonStrokingColor = c.CurrentNonStrokingColorSpace.GetColor(cyan, magenta, yellow, black)
	}
}

// InlineImageBuilder has been moved to inline_image_builder.go for better organization.
// See graphics/inline_image_builder.go for the full implementation including:
// - Soft mask (Smask) processing with validation
// - Color space resolution from NameToken or ArrayToken
// - Filter name extraction from /Filter or /F property
// - Decode array parsing from /Decode or /D property
// - Rendering intent extraction from /Intent property
// - imgDic resolution via token scanner

// IVerticalWritingSupported is implemented by fonts that support vertical text layout.
type IVerticalWritingSupported interface {
	GetPositionVector(code int) core.PdfPoint
	GetDisplacementVector(code int) core.PdfPoint
}

// resolveToDictionary resolves a token to its dictionary form, following indirect references.
func resolveToDictionary(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	if dict, ok := token.(*tokens.DictionaryToken); ok {
		return dict
	}
	if indRef, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		obj := scanner.Get(indRef.Data())
		if obj != nil {
			if d, isArray := obj.Data().(*tokens.DictionaryToken); isArray {
				return d
			}
		}
	}
	return nil
}
