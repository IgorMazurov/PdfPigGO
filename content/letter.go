package content

import (
	"fmt"
	"math"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
)

// letterPool caches Letter objects to reduce GC pressure during repeated document parsing.
var letterPool = sync.Pool{
	New: func() any {
		return &Letter{}
	},
}

// Letter represents a glyph or combination of glyphs drawn by a PDF content stream.
type Letter struct {
	Value             string
	TextOrientation   TextOrientation
	StartBaseLine     core.PdfPoint
	EndBaseLine       core.PdfPoint
	Width             float64
	BoundingBox       core.PdfRectangle
	GlyphRectangleLoose core.PdfRectangle
	FontSize          float64
	FontDetails       *fonts.FontDetails
	font              fonts.Font
	RenderingMode     core.TextRenderingMode
	Color             colors.Color
	StrokeColor       colors.Color
	FillColor         colors.Color
	PointSize         float64
	TextSequence      int
	// Code is the raw character code this letter was decoded from. It is the value
	// to pass to the font's TryGetPath/TryGetNormalisedPath to obtain the glyph outline.
	Code int
	// GlyphTransform maps the font's normalised glyph path (em units, as returned by
	// TryGetNormalisedPath) into page space, including font size, text matrix and CTM.
	// Apply it with TransformPath to obtain the on-page glyph outline.
	GlyphTransform core.TransformationMatrix
}

// Location returns the placement position of the character in PDF space.
func (l *Letter) Location() core.PdfPoint {
	return l.StartBaseLine
}

// FontName returns the name of the font, or empty string if unavailable.
func (l *Letter) FontName() string {
	if l.FontDetails == nil {
		return ""
	}
	return l.FontDetails.Name
}

// GlyphRectangle is deprecated; use BoundingBox instead.
func (l *Letter) GlyphRectangle() core.PdfRectangle {
	return l.BoundingBox
}

// Font is deprecated; use FontDetails instead.
func (l *Letter) Font() *fonts.FontDetails {
	return l.FontDetails
}

// NewLetter creates a new Letter from an IFont source.
func NewLetter(
	value string,
	glyphRectangle core.PdfRectangle,
	glyphRectangleLoose core.PdfRectangle,
	startBaseLine core.PdfPoint,
	endBaseLine core.PdfPoint,
	width float64,
	fontSize float64,
	font fonts.Font,
	renderingMode core.TextRenderingMode,
	strokeColor colors.Color,
	fillColor colors.Color,
	pointSize float64,
	textSequence int,
) *Letter {
	return newLetter(
		value,
		glyphRectangle,
		glyphRectangleLoose,
		startBaseLine,
		endBaseLine,
		width,
		fontSize,
		font.Details(),
		font,
		renderingMode,
		strokeColor,
		fillColor,
		pointSize,
		textSequence,
	)
}

// NewLetterWithDetails creates a new Letter from FontDetails (no full font).
func NewLetterWithDetails(
	value string,
	glyphRectangle core.PdfRectangle,
	glyphRectangleLoose core.PdfRectangle,
	startBaseLine core.PdfPoint,
	endBaseLine core.PdfPoint,
	width float64,
	fontSize float64,
	fontDetails *fonts.FontDetails,
	renderingMode core.TextRenderingMode,
	strokeColor colors.Color,
	fillColor colors.Color,
	pointSize float64,
	textSequence int,
) *Letter {
	return newLetter(
		value,
		glyphRectangle,
		glyphRectangleLoose,
		startBaseLine,
		endBaseLine,
		width,
		fontSize,
		fontDetails,
		nil,
		renderingMode,
		strokeColor,
		fillColor,
		pointSize,
		textSequence,
	)
}

func newLetter(
	value string,
	glyphRectangle core.PdfRectangle,
	glyphRectangleLoose core.PdfRectangle,
	startBaseLine core.PdfPoint,
	endBaseLine core.PdfPoint,
	width float64,
	fontSize float64,
	fontDetails *fonts.FontDetails,
	font fonts.Font,
	renderingMode core.TextRenderingMode,
	strokeColor colors.Color,
	fillColor colors.Color,
	pointSize float64,
	textSequence int,
) *Letter {
	var color, sc, fc colors.Color
	if renderingMode == core.StrokeText {
		sc = strokeColor
		if sc == nil {
			sc = colors.GrayBlack
		}
		fc = fillColor
		color = sc
	} else {
		fc = fillColor
		if fc == nil {
			fc = colors.GrayBlack
		}
		sc = strokeColor
		color = fc
	}

	l := letterPool.Get().(*Letter)
	*l = Letter{
		Value:             value,
		BoundingBox:       glyphRectangle,
		GlyphRectangleLoose: glyphRectangleLoose,
		StartBaseLine:     startBaseLine,
		EndBaseLine:       endBaseLine,
		Width:             width,
		FontSize:          fontSize,
		FontDetails:       fontDetails,
		font:              font,
		RenderingMode:     renderingMode,
		Color:             color,
		StrokeColor:       sc,
		FillColor:         fc,
		PointSize:         pointSize,
		TextSequence:      textSequence,
		TextOrientation:   getTextOrientation(startBaseLine, endBaseLine, glyphRectangle),
	}
	return l
}

// AsBold returns a new Letter with the same properties but bold font details.
func (l *Letter) AsBold() *Letter {
	bold := newLetter(
		l.Value,
		l.BoundingBox,
		l.GlyphRectangleLoose,
		l.StartBaseLine,
		l.EndBaseLine,
		l.Width,
		l.FontSize,
		l.FontDetails.AsBold(),
		l.font,
		l.RenderingMode,
		l.StrokeColor,
		l.FillColor,
		l.PointSize,
		l.TextSequence,
	)
	bold.Code = l.Code
	bold.GlyphTransform = l.GlyphTransform
	return bold
}

// GetFont returns the font associated with this letter, or nil if unavailable.
func (l *Letter) GetFont() fonts.Font {
	return l.font
}

func getTextOrientation(startBL, endBL core.PdfPoint, bb core.PdfRectangle) TextOrientation {
	if math.Abs(startBL.Y-endBL.Y) < 10e-5 {
		if math.Abs(startBL.X-endBL.X) < 10e-5 {
			return getTextOrientationRot(bb)
		}
		if startBL.X > endBL.X {
			return Rotate180TextOrientation
		}
		return HorizontalTextOrientation
	}

	if math.Abs(startBL.X-endBL.X) < 10e-5 {
		if math.Abs(startBL.Y-endBL.Y) < 10e-5 {
			return getTextOrientationRot(bb)
		}
		if startBL.Y > endBL.Y {
			return Rotate90TextOrientation
		}
		return Rotate270TextOrientation
	}

	return OtherTextOrientation
}

func getTextOrientationRot(bb core.PdfRectangle) TextOrientation {
	rotation := bb.Rotation()
	if math.Abs(math.Mod(rotation, 90)) >= 10e-5 {
		return OtherTextOrientation
	}

	rotInt := int(math.Round(rotation))
	switch rotInt {
	case 0:
		return HorizontalTextOrientation
	case -90:
		return Rotate90TextOrientation
	case 180, -180:
		return Rotate180TextOrientation
	case 90:
		return Rotate270TextOrientation
	default:
		return OtherTextOrientation
	}
}

// String returns a string representation of the letter and its position.
func (l *Letter) String() string {
	return fmt.Sprintf("%s %s %s %g", l.Value, l.Location(), l.FontName(), l.PointSize)
}
