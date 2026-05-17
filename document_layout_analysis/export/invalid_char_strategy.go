package export

// InvalidCharStrategy defines how to handle invalid characters during text export.
type InvalidCharStrategy byte

const (
	// Custom strategy allows caller-defined handling of invalid characters.
	Custom InvalidCharStrategy = 0

	// DoNotCheck skips validation and passes through all characters as-is.
	DoNotCheck InvalidCharStrategy = 1

	// Remove strips invalid characters from the output.
	Remove InvalidCharStrategy = 2

	// ConvertToHexadecimal replaces invalid characters with their hexadecimal representation.
	ConvertToHexadecimal InvalidCharStrategy = 3
)
