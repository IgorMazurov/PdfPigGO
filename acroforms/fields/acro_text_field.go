package fields

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroTextField is a box or space in which the user can enter text from the keyboard.
// The text may be restricted to a single line or may be permitted to span multiple lines.
type AcroTextField struct {
	AcroFieldBase

	// Flags specifies the behaviour of this field.
	Flags AcroTextFieldFlags

	// Value is the text value of this text field. Empty string if no value has been set.
	Value string

	// MaxLength is the optional maximum length of the text field.
	MaxLength *int

	// IsRichText indicates whether the field supports rich text content.
	IsRichText bool

	// IsMultiline indicates whether the field allows multiline text.
	IsMultiline bool
}

// NewAcroTextField creates a new AcroTextField.
func NewAcroTextField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroTextFieldFlags,
	information *AcroFieldCommonInformation,
	value *string,
	maxLength *int,
	pageNumber *int,
	bounds *core.PdfRectangle,
) (*AcroTextField, error) {
	base, err := NewAcroFieldBase(
		dictionary, fieldType, uint32(fieldFlags), AcroTypeText,
		information, pageNumber, bounds,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create text field base: %w", err)
	}

	textValue := ""
	if value != nil {
		textValue = *value
	}

	return &AcroTextField{
		AcroFieldBase: *base,
		Flags:         fieldFlags,
		Value:         textValue,
		MaxLength:     maxLength,
		IsRichText:    fieldFlags.HasFlag(AcroTextRichText),
		IsMultiline:   fieldFlags.HasFlag(AcroTextMultiline),
	}, nil
}

// HasFlag returns true if the specified flag is set.
func (f AcroTextFieldFlags) HasFlag(flag AcroTextFieldFlags) bool {
	return f&flag != 0
}

// String returns a string representation of the text field.
func (tf *AcroTextField) String() string {
	return fmt.Sprintf("%v: %s", tf.FieldType, tf.Value)
}
