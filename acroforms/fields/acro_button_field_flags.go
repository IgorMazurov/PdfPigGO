package fields

// AcroButtonFieldFlags specifies various characteristics of a button type field
// in an AcroForm. Multiple flags can be combined using bitwise OR.
type AcroButtonFieldFlags uint32

const (
	// ReadOnly indicates that the user may not change the value of the field.
	ReadOnly AcroButtonFieldFlags = 1 << 0

	// Required indicates that the field must have a value before the form can be submitted.
	Required = 1 << 1

	// NoExport indicates that the field must not be exported by the submit form action.
	NoExport = 1 << 2

	// NoToggleToOff indicates that for radio buttons, one radio button must be set at all times.
	NoToggleToOff = 1 << 14

	// Radio indicates that the field is a set of radio buttons.
	Radio = 1 << 15

	// PushButton indicates that the field is a push button.
	PushButton = 1 << 16

	// RadiosInUnison indicates that for radio buttons a group of radio buttons will toggle
	// on/off at the same time based on their initial value.
	RadiosInUnison = 1 << 25
)
