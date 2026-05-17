package png

// PngOpenerSettings holds configuration options for opening and parsing PNG files.
type PngOpenerSettings struct {
	ChunkVisitor       ChunkVisitor
	DisallowTrailingData bool
}
