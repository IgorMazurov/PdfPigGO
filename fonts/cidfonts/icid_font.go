// Package cidfonts provides types for handling CID (Character Identifier) fonts in PDFs.
package cidfonts

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CidFont represents a CID font containing glyph descriptions accessed by
// CID (character identifier) as character selectors. A CID font contains
// information about a CIDFont program but is not itself a font. It can only
// be a descendant of a Type 0 font.
type CidFont interface {
	// Type returns the font type name token (typically /Font).
	Type() *tokens.NameToken

	// SubType returns either Type0 (Adobe Type 1) or Type2 (TrueType).
	SubType() *tokens.NameToken

	// BaseFont returns the PostScript name of the CIDFont.
	BaseFont() *tokens.NameToken

	// SystemInfo returns the character collection definition for the font.
	SystemInfo() pdffonts.CidFontSystemInfo

	// Details returns the font details.
	Details() fonts.FontDetails

	// FontMatrix returns the transformation matrix for the font.
	FontMatrix() core.TransformationMatrix

	// CidFontType returns whether this is Type0 or Type2.
	CidFontType() CidFontType

	// Descriptor returns the font descriptor.
	Descriptor() *fonts.FontDescriptor

	// GetWidthFromDictionary gets the glyph width from the font dictionary for the given CID.
	GetWidthFromDictionary(cid int) float64

	// GetWidthFromFont gets the glyph width from the underlying font program for the given character identifier.
	GetWidthFromFont(characterIdentifier int) float64

	// GetBoundingBox returns the bounding box for the given character identifier.
	GetBoundingBox(characterIdentifier int) (core.PdfRectangle, error)

	// GetPositionVector returns the position vector for the given character identifier.
	GetPositionVector(characterIdentifier int) geometry.PdfVector

	// GetDisplacementVector returns the displacement vector for the given character identifier.
	GetDisplacementVector(characterIdentifier int) geometry.PdfVector

	// GetFontMatrix returns the font matrix for the given character identifier.
	GetFontMatrix(characterIdentifier int) core.TransformationMatrix

	// GetDescent returns the font descent value.
	GetDescent() float64

	// GetAscent returns the font ascent value.
	GetAscent() float64

	// TryGetPath attempts to get the glyph path for the given character code.
	// Returns true and the path if successful, false and nil otherwise.
	TryGetPath(characterCode int) ([]core.PdfSubpath, bool)

	// TryGetPathWithMapper attempts to get the glyph path using a custom
	// character-code-to-glyph-ID mapper function.
	TryGetPathWithMapper(characterCode int, mapper func(int) *int) ([]core.PdfSubpath, bool)

	// TryGetNormalisedPath attempts to get the normalised glyph path for the given character code.
	// Returns true and the path if successful, false and nil otherwise.
	TryGetNormalisedPath(characterCode int) ([]core.PdfSubpath, bool)

	// TryGetNormalisedPathWithMapper attempts to get the normalised glyph path using a custom
	// character-code-to-glyph-ID mapper function.
	TryGetNormalisedPathWithMapper(characterCode int, mapper func(int) *int) ([]core.PdfSubpath, bool)
}
