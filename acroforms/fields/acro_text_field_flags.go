package fields

// AcroTextFieldFlags specifies various characteristics of a text type field
// in an AcroForm. Multiple flags can be combined using bitwise OR.
type AcroTextFieldFlags uint32

const (
	// ReadOnly indicates that the user may not change the value of the field.
	AcroTextReadOnly AcroTextFieldFlags = 1 << 0

	// Required indicates that the field must have a value before the form can be submitted.
	AcroTextRequired = 1 << 1

	// NoExport indicates that the field must not be exported by the submit form action.
	AcroTextNoExport = 1 << 2

	// Multiline indicates that the field can contain multiple lines of text.
	AcroTextMultiline = 1 << 12

	// Password indicates that the field is for a password and should not be displayed
	// as text and should not be stored to file.
	AcroTextPassword = 1 << 13

	// FileSelect indicates that the field represents a file path selection.
	AcroTextFileSelect = 1 << 20

	// DoNotSpellCheck indicates that the text entered is not spell checked.
	AcroTextDoNotSpellCheck = 1 << 22

	// DoNotScroll indicates that the field does not scroll if the text exceeds
	// the bounds of the field.
	AcroTextDoNotScroll = 1 << 23

	// Comb indicates that for a text field which is not a Password, Multiline or FileSelect,
	// the field text is evenly spaced by splitting into combs based on the MaxLen entry
	// in the field dictionary.
	AcroTextComb = 1 << 24

	// RichText indicates that the value of the field is a rich text string.
	AcroTextRichText = 1 << 25
)
