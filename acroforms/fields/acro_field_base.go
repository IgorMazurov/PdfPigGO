package fields

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroFieldBase is the base type for a field in an interactive AcroForm.
type AcroFieldBase struct {
	// Dictionary is the raw PDF dictionary for this field.
	Dictionary *tokens.DictionaryToken

	// RawFieldType is the string representing the type of this field in PDF format.
	RawFieldType string

	// FieldType is the actual AcroFieldType represented by this field.
	FieldType AcroFieldType

	// FieldFlags specifies various characteristics of the field.
	FieldFlags uint32

	// Information is the optional information common to all types of field.
	Information *AcroFieldCommonInformation

	// PageNumber is the page number of the page containing this form field, if known.
	PageNumber *int

	// Bounds is the placement rectangle of this form field on the page given by PageNumber, if known.
	Bounds *core.PdfRectangle
}

// NewAcroFieldBase creates a new AcroFieldBase.
func NewAcroFieldBase(
	dictionary *tokens.DictionaryToken,
	rawFieldType string,
	fieldFlags uint32,
	fieldType AcroFieldType,
	information *AcroFieldCommonInformation,
	pageNumber *int,
	bounds *core.PdfRectangle,
) (*AcroFieldBase, error) {
	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}
	if rawFieldType == "" {
		return nil, fmt.Errorf("rawFieldType cannot be empty")
	}
	if information == nil {
		information = NewAcroFieldCommonInformation(nil, "", "", "")
	}

	return &AcroFieldBase{
		Dictionary:  dictionary,
		RawFieldType: rawFieldType,
		FieldFlags:  fieldFlags,
		FieldType:   fieldType,
		Information: information,
		PageNumber:  pageNumber,
		Bounds:      bounds,
	}, nil
}

// String returns a string representation of the field.
func (f *AcroFieldBase) String() string {
	return fmt.Sprintf("%v", f.FieldType)
}

// GetDictionary returns the raw PDF dictionary for this field.
func (f *AcroFieldBase) GetDictionary() *tokens.DictionaryToken {
	return f.Dictionary
}

// GetRawFieldType returns the string representing the type of this field in PDF format.
func (f *AcroFieldBase) GetRawFieldType() string {
	return f.RawFieldType
}

// GetFieldType returns the actual AcroFieldType represented by this field.
func (f *AcroFieldBase) GetFieldType() AcroFieldType {
	return f.FieldType
}

// GetFieldFlags returns various characteristics of the field.
func (f *AcroFieldBase) GetFieldFlags() uint32 {
	return f.FieldFlags
}

// GetInformation returns the optional information common to all types of field.
func (f *AcroFieldBase) GetInformation() *AcroFieldCommonInformation {
	return f.Information
}

// GetPageNumber returns the page number of the page containing this form field, if known.
func (f *AcroFieldBase) GetPageNumber() *int {
	return f.PageNumber
}

// GetBounds returns the placement rectangle of this form field on the page, if known.
func (f *AcroFieldBase) GetBounds() *core.PdfRectangle {
	return f.Bounds
}
