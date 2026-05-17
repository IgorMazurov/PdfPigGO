// Package fonts provides types for PDF font handling.
package fonts

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CMapProvider is an interface for types that can convert character codes to Unicode strings.
// Used to avoid import cycles between pdf_fonts and fonts/cmap packages.
type CMapProvider interface {
	// TryConvertToUnicode returns the Unicode string for the given character code.
	TryConvertToUnicode(code int) (string, bool)
}

// Font represents the base interface for a PDF font.
type Font interface {
	// Name returns the name of the font.
	Name() *tokens.NameToken

	// IsVertical reports whether the font is used for vertical text layout.
	IsVertical() bool

	// Details returns the details associated with this font.
	Details() *FontDetails

	// ReadCharacterCode reads the next character code from the given bytes
	// and returns it along with the number of bytes consumed (codeLength).
	ReadCharacterCode(bytes core.InputBytes) (int, int)

	// TryGetUnicode attempts to get the Unicode string for the given character code.
	// Returns true and the Unicode value if successful, false and an empty string otherwise.
	TryGetUnicode(characterCode int) (string, bool)

	// GetBoundingBox returns the bounding box for the given character code.
	// Returns an error if the character code is invalid or the bounding box cannot be determined.
	GetBoundingBox(characterCode int) (*CharacterBoundingBox, error)

	// GetFontMatrix returns the transformation matrix for this font.
	GetFontMatrix() core.TransformationMatrix

	// GetDescent returns the descent of the font adjusted by the font matrix.
	GetDescent() float64

	// GetAscent returns the ascent of the font adjusted by the font matrix.
	GetAscent() float64

	// TryGetPath attempts to get the glyph path for the given character code.
	// Returns true and the path if successful, false and nil otherwise.
	TryGetPath(characterCode int) ([]*core.PdfSubpath, bool)

	// TryGetNormalisedPath attempts to get the normalised glyph path for the
	// given character code. Returns true and the path if successful, false
	// and nil otherwise.
	TryGetNormalisedPath(characterCode int) ([]*core.PdfSubpath, bool)
}
