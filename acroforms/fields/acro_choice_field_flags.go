package fields

// AcroChoiceFieldFlags specifies various characteristics of a choice type field
// in an AcroForm. Multiple flags can be combined using bitwise OR.
type AcroChoiceFieldFlags uint32

const (
	// ReadOnly indicates that the user may not change the value of the field.
	AcroChoiceReadOnly AcroChoiceFieldFlags = 1 << 0

	// Required indicates that the field must have a value before the form can be submitted.
	AcroChoiceRequired = 1 << 1

	// NoExport indicates that the field must not be exported by the submit form action.
	AcroChoiceNoExport = 1 << 2

	// Combo indicates that the field is a combo box.
	AcroChoiceCombo = 1 << 17

	// Edit indicates that the combo box includes an editable text box.
	// Combo must be set.
	AcroChoiceEdit = 1 << 18

	// Sort indicates that the options should be sorted alphabetically;
	// this should be ignored by viewer applications.
	AcroChoiceSort = 1 << 19

	// MultiSelect indicates that the field allows multiple options to be selected.
	AcroChoiceMultiSelect = 1 << 21

	// DoNotSpellCheck indicates that the text entered in the field is not spell checked.
	// Combo and Edit must be set.
	AcroChoiceDoNotSpellCheck = 1 << 22

	// CommitOnSelectionChange indicates that any associated field action is fired when
	// the selection is changed rather than on losing focus.
	AcroChoiceCommitOnSelectionChange = 1 << 26
)
