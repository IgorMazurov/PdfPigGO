package interfaces

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// WritingFont defines the interface for a font used during PDF writing.
type WritingFont interface {
	HasWidths() bool
	Name() string
	TryGetBoundingBox(character rune) (*core.PdfRectangle, bool)
	TryGetAdvanceWidth(character rune) (float64, bool)
	GetFontMatrix() core.TransformationMatrix
	WriteFont(writer PdfStreamWriter, reservedIndirect *tokens.IndirectReferenceToken) *tokens.IndirectReferenceToken
	GetValueForCharacter(character rune) byte
}
