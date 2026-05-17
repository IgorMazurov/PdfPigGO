package cmapsubtables

// TrueTypeCMapPlatform identifies the platform for a CMap (character mapping)
// table in a TrueType font, determining how glyph indices are mapped to character codes.
type TrueTypeCMapPlatform uint16

const (
	// Unicode is the Unicode platform identifier for CMap tables.
	Unicode TrueTypeCMapPlatform = 0

	// Macintosh is the Apple Macintosh platform identifier for CMap tables.
	Macintosh TrueTypeCMapPlatform = 1

	// Reserved2 is an unused reserved value in the CMap platform enumeration.
	Reserved2 TrueTypeCMapPlatform = 2

	// Windows is the Microsoft Windows platform identifier for CMap tables.
	Windows TrueTypeCMapPlatform = 3
)
