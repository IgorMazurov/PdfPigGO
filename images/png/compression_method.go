package png

// CompressionMethod indicates the compression algorithm used for PNG image data.
type CompressionMethod byte

const (
	// DeflateWithSlidingWindow is the only valid method: DEFLATE with a 32K sliding window.
	DeflateWithSlidingWindow CompressionMethod = 0
)
