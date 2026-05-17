package fields

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroListBoxField is a scrollable list box field.
type AcroListBoxField struct {
	AcroFieldBase

	Flags                 AcroChoiceFieldFlags
	Options               []*AcroChoiceOption
	SelectedOptions       []string
	SelectedOptionIndices []int
	TopIndex              int
}

// NewAcroListBoxField creates a new AcroListBoxField.
func NewAcroListBoxField(
	dictionary *tokens.DictionaryToken,
	fieldType string,
	fieldFlags AcroChoiceFieldFlags,
	information *AcroFieldCommonInformation,
	options []*AcroChoiceOption,
	selectedOptions []string,
	selectedOptionIndices []int,
	topIndex *int,
	pageNumber *int,
	bounds *core.PdfRectangle,
) (*AcroListBoxField, error) {
	if dictionary == nil {
		return nil, errors.New("dictionary cannot be nil")
	}
	if options == nil {
		return nil, errors.New("options cannot be nil")
	}
	if selectedOptions == nil {
		return nil, errors.New("selectedOptions cannot be nil")
	}

	ti := 0
	if topIndex != nil {
		ti = *topIndex
	}

	base := AcroFieldBase{
		Dictionary:   dictionary,
		RawFieldType: fieldType,
		FieldFlags:   uint32(fieldFlags),
		FieldType:    AcroTypeListBox,
		Information:  information,
		PageNumber:   pageNumber,
		Bounds:       bounds,
	}

	return &AcroListBoxField{
		AcroFieldBase:         base,
		Flags:                 fieldFlags,
		Options:               options,
		SelectedOptions:       selectedOptions,
		SelectedOptionIndices: selectedOptionIndices,
		TopIndex:              ti,
	}, nil
}

// SupportsMultiSelect reports whether the field allows multiple selections.
func (f *AcroListBoxField) SupportsMultiSelect() bool {
	return f.Flags&AcroChoiceMultiSelect != 0
}
