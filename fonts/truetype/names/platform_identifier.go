package names

// TrueTypePlatformIdentifier identifies the platform for which a TrueType/OpenType
// font name table entry is intended, allowing for platform-specific implementations.
type TrueTypePlatformIdentifier uint16

const (
	// Unicode is the Unicode platform identifier.
	Unicode TrueTypePlatformIdentifier = 0

	// Macintosh is the Macintosh platform identifier.
	Macintosh TrueTypePlatformIdentifier = 1

	// Iso was originally for ISO 10646 but is now deprecated as it and Unicode
	// have identical character code assignments.
	Iso TrueTypePlatformIdentifier = 2

	// Windows is the Microsoft Windows platform identifier.
	Windows TrueTypePlatformIdentifier = 3
)
