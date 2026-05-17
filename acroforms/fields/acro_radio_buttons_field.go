package fields

import (
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroRadioButtonsField represents a set of radio buttons in an AcroForm.
type AcroRadioButtonsField struct {
	AcroNonTerminalField

	// Flags defines the behaviour of this button type.
	Flags AcroButtonFieldFlags
}

// NewAcroRadioButtonsField creates a new AcroRadioButtonsField.
func NewAcroRadioButtonsField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroButtonFieldFlags,
	information *AcroFieldCommonInformation,
	children []any,
) (*AcroRadioButtonsField, error) {
	nonTerminal, err := NewAcroNonTerminalField(
		dictionary, fieldType, uint32(fieldFlags), information, AcroTypeRadioButtons, children,
	)
	if err != nil {
		return nil, err
	}

	return &AcroRadioButtonsField{
		AcroNonTerminalField: *nonTerminal,
		Flags:                fieldFlags,
	}, nil
}
