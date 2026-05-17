package fields

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroNonTerminalField represents a non-leaf field in the form's structure.
type AcroNonTerminalField struct {
	AcroFieldBase
	children []any
}

// Children returns the child fields of this field, preserving their concrete types.
// Each element can be type-asserted to its specific field type.
func (f *AcroNonTerminalField) Children() []any {
	return f.children
}

// NewAcroNonTerminalField creates a new AcroNonTerminalField.
func NewAcroNonTerminalField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags uint32,
	information *AcroFieldCommonInformation,
	acroFieldType AcroFieldType,
	children []any,
) (*AcroNonTerminalField, error) {
	if dictionary == nil {
		return nil, errors.New("dictionary cannot be nil")
	}
	if fieldType == "" {
		return nil, errors.New("fieldType cannot be empty")
	}
	if children == nil {
		return nil, errors.New("children cannot be nil")
	}

	base := AcroFieldBase{
		Dictionary:   dictionary,
		RawFieldType: fieldType,
		FieldFlags:   fieldFlags,
		FieldType:    acroFieldType,
		Information:  information,
		PageNumber:   nil,
		Bounds:       nil,
	}

	return &AcroNonTerminalField{
		AcroFieldBase: base,
		children:      children,
	}, nil
}
