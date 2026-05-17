package fields

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroComboBoxField is a combo box consisting of a drop-down list optionally accompanied
// by an editable text box in which the user can type a value other than predefined choices.
type AcroComboBoxField struct {
	AcroFieldBase

	Flags               AcroChoiceFieldFlags
	Options             []*AcroChoiceOption
	SelectedOptions     []string
	SelectedOptionIndices []int
}

// NewAcroComboBoxField creates a new AcroComboBoxField.
func NewAcroComboBoxField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroChoiceFieldFlags,
	information *AcroFieldCommonInformation,
	options []*AcroChoiceOption,
	selectedOptions []string,
	selectedOptionIndices []int,
	pageNumber *int,
	bounds *core.PdfRectangle,
) (*AcroComboBoxField, error) {
	if dictionary == nil {
		return nil, errors.New("dictionary cannot be nil")
	}
	if options == nil {
		return nil, errors.New("options cannot be nil")
	}
	if selectedOptions == nil {
		return nil, errors.New("selectedOptions cannot be nil")
	}

	base := AcroFieldBase{
		Dictionary:   dictionary,
		RawFieldType: fieldType,
		FieldFlags:   uint32(fieldFlags),
		FieldType:    AcroTypeComboBox,
		Information:  information,
		PageNumber:   pageNumber,
		Bounds:       bounds,
	}

	return &AcroComboBoxField{
		AcroFieldBase:         base,
		Flags:                 fieldFlags,
		Options:               options,
		SelectedOptions:       selectedOptions,
		SelectedOptionIndices: selectedOptionIndices,
	}, nil
}
