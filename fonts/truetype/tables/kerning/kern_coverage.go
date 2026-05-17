package kerning

// KernCoverage specifies the type of kerning covered by a kerning sub-table
// in TrueType fonts. Multiple flags can be combined using bitwise OR.
type KernCoverage int

const (
	// Horizontal indicates that the table contains horizontal kerning data.
	Horizontal KernCoverage = 1 << iota

	// Minimum indicates that the table has minimum values rather than kerning values.
	Minimum

	// CrossStream indicates that kerning is perpendicular to the flow of text.
	// If text is horizontal, kerning will be in the up/down direction.
	CrossStream

	// Override indicates that the value in this sub-table should replace the
	// currently accumulated value.
	Override
)
