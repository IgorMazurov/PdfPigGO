package acroforms

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/acroforms/fields"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// AcroForm is a collection of interactive fields for gathering data from a user
// through dropdowns, textboxes, checkboxes, etc. Each PdfDocument with form
// functionality contains a single AcroForm spread across one or more pages.
type AcroForm struct {
	dictionary           *tokens.DictionaryToken
	signatureFlags       SignatureFlags
	needAppearances      bool
	fieldsWithReferences map[core.IndirectReference]any
	rawFields            []any
}

// NewAcroForm creates a new AcroForm.
func NewAcroForm(
	dictionary *tokens.DictionaryToken,
	signatureFlags SignatureFlags,
	needAppearances bool,
	fieldsWithReferences map[core.IndirectReference]any,
) (*AcroForm, error) {
	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}
	if fieldsWithReferences == nil {
		return nil, fmt.Errorf("fieldsWithReferences cannot be nil")
	}

	fieldList := make([]any, 0, len(fieldsWithReferences))
	for _, f := range fieldsWithReferences {
		fieldList = append(fieldList, f)
	}

	return &AcroForm{
		dictionary:           dictionary,
		signatureFlags:       signatureFlags,
		needAppearances:      needAppearances,
		fieldsWithReferences: fieldsWithReferences,
		rawFields:            fieldList,
	}, nil
}

// Dictionary returns the raw PDF dictionary which is the root form object.
func (a *AcroForm) Dictionary() *tokens.DictionaryToken {
	return a.dictionary
}

// SignatureFlags returns document-level characteristics related to signature fields.
func (a *AcroForm) SignatureFlags() SignatureFlags {
	return a.signatureFlags
}

// NeedAppearances reports whether all widget annotations need appearance dictionaries and streams.
func (a *AcroForm) NeedAppearances() bool {
	return a.needAppearances
}

// Fields returns all root fields in this form as their concrete types.
// Each element can be type-asserted to its specific field type, e.g.:
//   if rb, ok := any(f).(*fields.AcroRadioButtonsField); ok { ... }
func (a *AcroForm) Fields() []any {
	return a.rawFields
}

// GetFieldsForPage returns the set of fields which appear on the given page number.
// Page numbers are 1-based. An error is returned if pageNumber is less than 1.
func (a *AcroForm) GetFieldsForPage(pageNumber int) ([]*fields.AcroFieldBase, error) {
	if pageNumber <= 0 {
		return nil, fmt.Errorf("page number starts at 1, instead got %d", pageNumber)
	}

	result := make([]*fields.AcroFieldBase, 0)

	for _, fieldAny := range a.rawFields {
		field := toAcroFieldBase(fieldAny)
		if field == nil {
			continue
		}

		pageNum := field.GetPageNumber()
		if pageNum != nil && *pageNum == pageNumber {
			result = append(result, field)
			continue
		}

		if nonTerminal, ok := unwrapAcroNonTerminalFromAny(fieldAny); ok {
			children := nonTerminal.Children()
			for _, childAny := range children {
				child := toAcroFieldBaseFromAny(childAny)
				if child == nil {
					continue
				}
				childPageNum := child.GetPageNumber()
				if childPageNum != nil && *childPageNum == pageNumber {
					result = append(result, field)
					break
				}
			}
		}
	}

	return result, nil
}

// unwrapAcroNonTerminal checks whether the given *AcroFieldBase pointer was derived from
// an AcroNonTerminalField and returns it. The second return value is true if the
// conversion succeeded. Note: this works on base pointers extracted via toAcroFieldBase.
func unwrapAcroNonTerminal(field *fields.AcroFieldBase) (*fields.AcroNonTerminalField, bool) {
	nt, ok := any(field).(*fields.AcroNonTerminalField)
	return nt, ok
}

// unwrapAcroNonTerminalFromAny checks whether the given interface{} holds an
// AcroNonTerminalField (or a type embedding it) and returns it.
func unwrapAcroNonTerminalFromAny(fieldAny any) (*fields.AcroNonTerminalField, bool) {
	if nt, ok := fieldAny.(*fields.AcroNonTerminalField); ok {
		return nt, true
	}
	if rb, ok := fieldAny.(*fields.AcroRadioButtonsField); ok {
		return &rb.AcroNonTerminalField, true
	}
	if cb, ok := fieldAny.(*fields.AcroCheckboxesField); ok {
		return &cb.AcroNonTerminalField, true
	}
	return nil, false
}

// ToAcroFieldBaseFromAny extracts a pointer to the embedded AcroFieldBase from any field type.
func ToAcroFieldBaseFromAny(fieldAny any) *fields.AcroFieldBase {
	if fieldAny == nil {
		return nil
	}
	switch v := fieldAny.(type) {
	case *fields.AcroFieldBase:
		return v
	case *fields.AcroNonTerminalField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroTextField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroCheckboxField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroRadioButtonField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroPushButtonField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroSignatureField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroRadioButtonsField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroCheckboxesField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroComboBoxField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	case *fields.AcroListBoxField:
		if v == nil {
			return nil
		}
		return &v.AcroFieldBase
	}
	return nil
}

// toAcroFieldBaseFromAny is the internal alias for ToAcroFieldBaseFromAny.
func toAcroFieldBaseFromAny(fieldAny any) *fields.AcroFieldBase {
	return ToAcroFieldBaseFromAny(fieldAny)
}

// String returns the string representation of the AcroForm.
func (a *AcroForm) String() string {
	if a == nil || a.dictionary == nil {
		return ""
	}
	return a.dictionary.String()
}
