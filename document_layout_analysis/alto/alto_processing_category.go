package alto

// AltoProcessingCategory represents the category of a processing step in ALTO format.
type AltoProcessingCategory uint32

const (
	AltoProcessingCategoryContentGeneration  AltoProcessingCategory = 1
	AltoProcessingCategoryContentModification AltoProcessingCategory = 2
	AltoProcessingCategoryPreOperation       AltoProcessingCategory = 4
	AltoProcessingCategoryPostOperation      AltoProcessingCategory = 8
	AltoProcessingCategoryOther              AltoProcessingCategory = 16
)

// String returns the XML string representation of this processing category.
func (c AltoProcessingCategory) String() string {
	var parts []string
	if c&AltoProcessingCategoryContentGeneration != 0 {
		parts = append(parts, "contentGeneration")
	}
	if c&AltoProcessingCategoryContentModification != 0 {
		parts = append(parts, "contentModification")
	}
	if c&AltoProcessingCategoryPreOperation != 0 {
		parts = append(parts, "preOperation")
	}
	if c&AltoProcessingCategoryPostOperation != 0 {
		parts = append(parts, "postOperation")
	}
	if c&AltoProcessingCategoryOther != 0 {
		parts = append(parts, "other")
	}
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += "|" + parts[i]
	}
	return result
}
