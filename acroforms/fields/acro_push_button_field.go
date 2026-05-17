package fields

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroPushButtonField represents a push button that responds immediately to user input
// without storing any state.
type AcroPushButtonField struct {
	AcroFieldBase

	// Flags defines the behaviour of this button type.
	Flags AcroButtonFieldFlags
}

// NewAcroPushButtonField creates a new AcroPushButtonField.
func NewAcroPushButtonField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroButtonFieldFlags,
	information *AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
) (*AcroPushButtonField, error) {
	base, err := NewAcroFieldBase(
		dictionary, fieldType, uint32(fieldFlags), AcroTypePushButton,
		information, pageNumber, bounds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create push button field base: %w", err)
	}

	return &AcroPushButtonField{
		AcroFieldBase: *base,
		Flags:         fieldFlags,
	}, nil
}
