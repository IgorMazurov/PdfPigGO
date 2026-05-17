package fonts

// FontStretch represents the width variation of a font face.
type FontStretch int

const (
	// Unknown indicates that the stretch value is not specified or not recognized.
	Unknown FontStretch = -1

	// UltraCondensed indicates an ultra condensed font width.
	UltraCondensed FontStretch = iota

	// ExtraCondensed indicates an extra condensed font width.
	ExtraCondensed

	// Condensed indicates a condensed font width.
	Condensed

	// SemiCondensed indicates a semi condensed font width.
	SemiCondensed

	// Normal indicates a normal (regular) font width.
	Normal

	// SemiExpanded indicates a semi expanded font width.
	SemiExpanded

	// Expanded indicates an expanded font width.
	Expanded

	// ExtraExpanded indicates an extra expanded font width.
	ExtraExpanded

	// UltraExpanded indicates an ultra expanded font width.
	UltraExpanded
)
