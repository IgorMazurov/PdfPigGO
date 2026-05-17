package fields

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroSignatureField represents a digital signature field.
type AcroSignatureField struct {
	AcroFieldBase
}

// NewAcroSignatureField creates a new AcroSignatureField.
func NewAcroSignatureField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags uint32,
	information *AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
) (*AcroSignatureField, error) {
	base, err := NewAcroFieldBase(
		dictionary, fieldType, fieldFlags, AcroTypeSignature,
		information, pageNumber, bounds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature field base: %w", err)
	}

	return &AcroSignatureField{
		AcroFieldBase: *base,
	}, nil
}
