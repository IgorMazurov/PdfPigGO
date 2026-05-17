package fonts

// FontDescriptorFlags specifies various characteristics of a font.
type FontDescriptorFlags uint32

const (
	// None indicates that no flags are set.
	None FontDescriptorFlags = 0

	// FixedPitch indicates that all glyphs have the same width.
	FixedPitch FontDescriptorFlags = 1

	// Serif indicates that glyphs have serifs.
	Serif = 1 << 1

	// Symbolic indicates that there are glyphs outside the Adobe standard Latin set.
	Symbolic = 1 << 2

	// Script indicates that the glyphs resemble cursive handwriting.
	Script = 1 << 3

	// NonSymbolic indicates that font uses a (sub)set of the Adobe standard Latin set.
	// Cannot be set at the same time as Symbolic.
	NonSymbolic = 1 << 5

	// Italic indicates that font is italic.
	Italic = 1 << 6

	// AllCap indicates that font contains only uppercase letters.
	AllCap = 1 << 16

	// SmallCap indicates that lowercase letters are smaller versions of the uppercase equivalent.
	SmallCap = 1 << 17

	// ForceBold forces small bold text to be rendered bold.
	ForceBold = 1 << 18
)
