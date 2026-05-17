package png

import "io"

// ChunkVisitor visits each chunk in a PNG file during parsing.
type ChunkVisitor interface {
	Visit(stream io.Reader, header ImageHeader, chunkHeader ChunkHeader, data []byte, crc []byte)
}
