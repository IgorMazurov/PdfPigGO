package fonts

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// FontDescriptor specifies metrics and attributes of a simple font or CID Font
// for the whole font rather than per-glyph. It provides information to enable
// consumer applications to find a substitute font when the font is unavailable.
type FontDescriptor struct {
	// FontName is the PostScript name for the font. Required.
	FontName *tokens.NameToken

	// FontFamily is the preferred font family. Optional.
	FontFamily string

	// Stretch is the font stretch value. Optional.
	Stretch FontStretch

	// FontWeight is the weight/thickness of the font. Values: 100, 200, 300,
	// 500 (normal), 600, 700, 800, 900. Optional.
	FontWeight float64

	// Flags defines various font characteristics. See FontDescriptorFlags. Required.
	Flags FontDescriptorFlags

	// BoundingBox is a rectangle in glyph coordinates representing the smallest
	// rectangle containing all glyphs of the font. Required (except Type 3).
	BoundingBox core.PdfRectangle

	// ItalicAngle is the angle in degrees counter-clockwise from vertical of the
	// vertical lines of the font. Negative for fonts sloping right (italic).
	// Required.
	ItalicAngle float64

	// Ascent is the maximum height above the baseline for any glyph except accents.
	// Required (except Type 3).
	Ascent float64

	// Descent is the maximum depth below the baseline for any glyph. Negative value.
	// Required (except Type 3).
	Descent float64

	// Leading is the spacing between consecutive lines of text. Default 0. Optional.
	Leading float64

	// CapHeight is the vertical distance of the top of flat capital letters from
	// the baseline. Required where Latin characters are present (except Type 3).
	CapHeight float64

	// XHeight is the vertical distance of the top of flat non-ascending lowercase
	// letters (e.g., x) from the baseline. Default 0. Optional.
	XHeight float64

	// StemVertical is the horizontal thickness of vertical stems of glyphs.
	// Required (except Type 3).
	StemVertical float64

	// StemHorizontal is the vertical thickness of horizontal stems of glyphs.
	// Default 0. Optional.
	StemHorizontal float64

	// AverageWidth is the average glyph width in the font. Default 0. Optional.
	AverageWidth float64

	// MaxWidth is the maximum glyph width in the font. Default 0. Optional.
	MaxWidth float64

	// MissingWidth is the width for character codes whose widths are not present
	// in the Widths array of the font dictionary. Default 0. Optional.
	MissingWidth float64

	// FontFile holds the bytes of the font program. Optional.
	FontFile *DescriptorFontFile

	// CharSet is the character names defined in a font subset. Optional.
	CharSet string
}

// NewFontDescriptor creates a new FontDescriptor from the given builder.
func NewFontDescriptor(builder *FontDescriptorBuilder) *FontDescriptor {
	return &FontDescriptor{
		FontName:       builder.FontName,
		FontFamily:     builder.FontFamily,
		Stretch:        builder.Stretch,
		FontWeight:     builder.FontWeight,
		Flags:          builder.Flags,
		BoundingBox:    builder.BoundingBox,
		ItalicAngle:    builder.ItalicAngle,
		Ascent:         builder.Ascent,
		Descent:        builder.Descent,
		Leading:        builder.Leading,
		CapHeight:      builder.CapHeight,
		XHeight:        builder.XHeight,
		StemVertical:   builder.StemVertical,
		StemHorizontal: builder.StemHorizontal,
		AverageWidth:   builder.AverageWidth,
		MaxWidth:       builder.MaxWidth,
		MissingWidth:   builder.MissingWidth,
		FontFile:       builder.FontFile,
		CharSet:        builder.CharSet,
	}
}

// ToDetails returns a FontDetails instance derived from this descriptor.
// If name is empty, the font's own FontName is used.
func (fd *FontDescriptor) ToDetails(name string) *FontDetails {
	if name == "" && fd.FontName != nil {
		name = fd.FontName.Data()
	}
	isItalic := fd.Flags.HasFlag(Italic) || fd.ItalicAngle != 0
	return NewFontDetails(name, fd.FontWeight > 500, int(fd.FontWeight), isItalic)
}

// HasFlag reports whether the FontDescriptorFlags value contains the given flag.
func (f FontDescriptorFlags) HasFlag(flag FontDescriptorFlags) bool {
	return f&flag != 0
}

// FontDescriptorBuilder provides a mutable way to construct a FontDescriptor.
type FontDescriptorBuilder struct {
	FontName       *tokens.NameToken
	FontFamily     string
	Stretch        FontStretch
	FontWeight     float64
	Flags          FontDescriptorFlags
	BoundingBox    core.PdfRectangle
	ItalicAngle    float64
	Ascent         float64
	Descent        float64
	Leading        float64
	CapHeight      float64
	XHeight        float64
	StemVertical   float64
	StemHorizontal float64
	AverageWidth   float64
	MaxWidth       float64
	MissingWidth   float64
	FontFile       *DescriptorFontFile
	CharSet        string
}

// NewFontDescriptorBuilder creates a new builder with the required fields.
func NewFontDescriptorBuilder(fontName *tokens.NameToken, flags FontDescriptorFlags) *FontDescriptorBuilder {
	return &FontDescriptorBuilder{
		FontName:  fontName,
		Stretch:   Normal,
		FontWeight: 400,
		Flags:     flags,
	}
}

// Build creates the FontDescriptor with values from this builder.
func (b *FontDescriptorBuilder) Build() *FontDescriptor {
	return NewFontDescriptor(b)
}
