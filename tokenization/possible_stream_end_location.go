package tokenization

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// PossibleStreamEndLocation stores an observed occurrence of 'endobj' or
// 'endstream' while reading streams, used internally by PdfTokenScanner.
type PossibleStreamEndLocation struct {
	// Offset is the byte position at which the token started in the file.
	Offset int64

	// Type is one of tokens.EndObject or tokens.EndStream.
	Type *tokens.OperatorToken
}

// NewPossibleStreamEndLocation creates a new PossibleStreamEndLocation.
func NewPossibleStreamEndLocation(offset int64, typ *tokens.OperatorToken) PossibleStreamEndLocation {
	return PossibleStreamEndLocation{
		Offset: offset,
		Type:   typ,
	}
}

// String returns the string representation of the possible stream end location.
func (p PossibleStreamEndLocation) String() string {
	return fmt.Sprintf("%d: %v", p.Offset, p.Type)
}
