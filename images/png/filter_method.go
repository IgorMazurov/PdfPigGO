package png

// FilterMethod indicates the filtering method used for PNG image data.
type FilterMethod byte

const (
	// AdaptiveFiltering means each row may use a different filter, indicated by a filter type byte.
	AdaptiveFiltering FilterMethod = 0
)
