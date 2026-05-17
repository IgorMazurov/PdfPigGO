package png

// InterlaceMethod indicates the interlacing scheme used for a PNG image.
type InterlaceMethod byte

const (
	// InterlaceMethodNone means no interlacing; rows are stored in order from top to bottom.
	InterlaceMethodNone InterlaceMethod = iota // 0

	// InterlaceMethodAdam7 is the Adam7 interlacing scheme with 7 passes.
	InterlaceMethodAdam7 // 1
)
