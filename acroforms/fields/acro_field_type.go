package fields

// AcroFieldType indicates the type of field for an AcroFieldBase.
type AcroFieldType int

const (
	// AcroTypePushButton is a button that immediately responds to user input without retaining state.
	AcroTypePushButton AcroFieldType = iota

	// AcroTypeCheckboxes is a set of checkboxes.
	AcroTypeCheckboxes

	// AcroTypeCheckbox is a checkbox which toggles between on and off states.
	AcroTypeCheckbox

	// AcroTypeRadioButtons is a set of radio buttons.
	AcroTypeRadioButtons

	// AcroTypeRadioButton is a single radio button, as part of a set or on its own.
	AcroTypeRadioButton

	// AcroTypeText is a textbox allowing user input through the keyboard.
	AcroTypeText

	// AcroTypeComboBox is a dropdown list of options with optional user-editable textbox.
	AcroTypeComboBox

	// AcroTypeListBox is a list of options for the user to select from.
	AcroTypeListBox

	// AcroTypeSignature is a field containing a digital signature.
	AcroTypeSignature

	// AcroTypeUnknown indicates that the field type wasn't specified.
	AcroTypeUnknown
)
