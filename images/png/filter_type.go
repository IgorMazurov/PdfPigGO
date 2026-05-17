package png

// FilterType indicates the filter applied to a row of PNG image data.
type FilterType byte

const (
	// FilterTypeNone means no filtering is applied; raw bytes are used as-is.
	FilterTypeNone FilterType = iota // 0

	// FilterTypeSub predicts each byte from the left neighbor.
	FilterTypeSub // 1

	// FilterTypeUp predicts each byte from the byte directly above.
	FilterTypeUp // 2

	// FilterTypeAverage predicts each byte from the average of left and above neighbors.
	FilterTypeAverage // 3

	// FilterTypePaeth predicts each byte using a heuristic based on left, above, and upper-left neighbors.
	FilterTypePaeth // 4
)
