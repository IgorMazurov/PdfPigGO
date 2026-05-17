package acroforms

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/acroforms/fields"
)

// GetFields returns all leaf fields in the form by recursively flattening
// the field hierarchy. Non-terminal fields are expanded into their children,
// and only fields with a known type are included.
func GetFields(form *AcroForm) []*fields.AcroFieldBase {
	if form == nil {
		return nil
	}

	var result []*fields.AcroFieldBase
	for _, fAny := range form.Fields() {
		f := toAcroFieldBaseFromAny(fAny)
		if f != nil {
			result = append(result, flattenField(f)...)
		}
	}
	return result
}

// flattenField recursively collects all fields from a single root field.
// If the field has a known type (not Unknown), it is included in the result.
// If the field is non-terminal, its children are also flattened recursively.
func flattenField(field *fields.AcroFieldBase) []*fields.AcroFieldBase {
	if field == nil {
		return nil
	}

	var result []*fields.AcroFieldBase

	if field.FieldType != fields.AcroTypeUnknown {
		result = append(result, field)
	}

	if nonTerminal, ok := any(field).(*fields.AcroNonTerminalField); ok {
		for _, childAny := range nonTerminal.Children() {
			child := toAcroFieldBaseFromAny(childAny)
			if child != nil {
				result = append(result, flattenField(child)...)
			}
		}
	}

	return result
}

// GetFieldValue returns the field's partial name and its current value as a
// string pair. For text fields the value is the entered text; for checkboxes
// it is "true" or "false"; for all other types the value is an empty string.
func GetFieldValue(field *fields.AcroFieldBase) (string, string) {
	if field == nil || field.Information == nil {
		return "", ""
	}

	name := field.Information.PartialName

	switch f := any(field).(type) {
	case *fields.AcroTextField:
		return name, f.Value
	case *fields.AcroCheckboxField:
		return name, fmt.Sprintf("%t", f.IsChecked)
	default:
		return name, ""
	}
}
