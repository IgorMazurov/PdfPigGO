package kerning

// KerningSubTable holds the data for a single sub-table within a TrueType kerning table.
type KerningSubTable struct {
	// Version is the version of this kerning sub-table (0 or 1).
	Version int

	// Coverage specifies what type of kerning this sub-table covers.
	Coverage KernCoverage

	// Pairs is the list of kern pairs in this sub-table.
	Pairs []KernPair
}

// NewKerningSubTable creates a new KerningSubTable with the given version, coverage, and pairs.
func NewKerningSubTable(version int, coverage KernCoverage, pairs []KernPair) *KerningSubTable {
	return &KerningSubTable{
		Version:  version,
		Coverage: coverage,
		Pairs:    pairs,
	}
}
