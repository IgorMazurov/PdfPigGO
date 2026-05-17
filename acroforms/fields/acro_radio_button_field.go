package fields

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroRadioButtonField represents a single radio button within an interactive AcroForm.
type AcroRadioButtonField struct {
	AcroFieldBase

	// Flags defines the behaviour of this button type.
	Flags AcroButtonFieldFlags

	// CurrentValue is the current value of this radio button.
	CurrentValue *tokens.NameToken

	// IsSelected indicates whether the radio button is currently on/active.
	IsSelected bool
}

// NewAcroRadioButtonField creates a new AcroRadioButtonField.
func NewAcroRadioButtonField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroButtonFieldFlags,
	information *AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
	currentValue *tokens.NameToken,
	isSelected bool,
) (*AcroRadioButtonField, error) {
	base, err := NewAcroFieldBase(
		dictionary, fieldType, uint32(fieldFlags), AcroTypeRadioButton,
		information, pageNumber, bounds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create radio button field base: %w", err)
	}

	return &AcroRadioButtonField{
		AcroFieldBase: *base,
		Flags:         fieldFlags,
		CurrentValue:  currentValue,
		IsSelected:    isSelected,
	}, nil
}
