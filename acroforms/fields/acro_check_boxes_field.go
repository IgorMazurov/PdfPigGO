package fields

import (
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroCheckboxesField represents a set of related checkboxes that share the same export value.
type AcroCheckboxesField struct {
	AcroNonTerminalField
}

// NewAcroCheckboxesField creates a new AcroCheckboxesField.
func NewAcroCheckboxesField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroButtonFieldFlags,
	information *AcroFieldCommonInformation,
	children []any,
) (*AcroCheckboxesField, error) {
	nonTerminal, err := NewAcroNonTerminalField(
		dictionary, fieldType, uint32(fieldFlags), information, AcroTypeCheckboxes, children,
	)
	if err != nil {
		return nil, err
	}

	return &AcroCheckboxesField{
		AcroNonTerminalField: *nonTerminal,
	}, nil
}
