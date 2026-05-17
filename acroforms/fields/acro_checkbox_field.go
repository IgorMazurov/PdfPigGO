package fields

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroCheckboxField represents a checkbox which may be toggled on or off.
type AcroCheckboxField struct {
	AcroFieldBase

	// Flags defines the behaviour of this button type.
	Flags AcroButtonFieldFlags

	// CurrentValue is the current value of this checkbox.
	CurrentValue *tokens.NameToken

	// IsChecked indicates whether this checkbox is currently checked/on.
	IsChecked bool
}

// NewAcroCheckboxField creates a new AcroCheckboxField.
func NewAcroCheckboxField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroButtonFieldFlags,
	information *AcroFieldCommonInformation,
	currentValue *tokens.NameToken,
	isChecked bool,
	pageNumber *int,
	bounds *core.PdfRectangle,
) (*AcroCheckboxField, error) {
	base, err := NewAcroFieldBase(
		dictionary, fieldType, uint32(fieldFlags), AcroTypeCheckbox,
		information, pageNumber, bounds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkbox field base: %w", err)
	}

	return &AcroCheckboxField{
		AcroFieldBase: *base,
		Flags:         fieldFlags,
		CurrentValue:  currentValue,
		IsChecked:     isChecked,
	}, nil
}
