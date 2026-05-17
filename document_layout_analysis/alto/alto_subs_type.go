package alto

// AltoSubsType represents the type of substitution in ALTO format.
type AltoSubsType string

const (
	// AltoSubsTypeHypPart1 indicates the first part of a hyphenated word.
	AltoSubsTypeHypPart1 AltoSubsType = "HypPart1"
	// AltoSubsTypeHypPart2 indicates the second part of a hyphenated word.
	AltoSubsTypeHypPart2 AltoSubsType = "HypPart2"
	// AltoSubsTypeAbbreviation indicates an abbreviation substitution.
	AltoSubsTypeAbbreviation AltoSubsType = "Abbreviation"
)

// String returns the XML string representation of this substitution type.
func (t AltoSubsType) String() string {
	return string(t)
}
