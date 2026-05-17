package png

import (
	"fmt"
)

// ChunkHeader describes the header of a PNG chunk.
type ChunkHeader struct {
	position int64
	length   int
	name     string
}

// NewChunkHeader creates a new ChunkHeader, returning an error if length is negative.
func NewChunkHeader(position int64, length int, name string) (ChunkHeader, error) {
	if length < 0 {
		return ChunkHeader{}, fmt.Errorf("length less than zero (%d) encountered when reading chunk at position %d", length, position)
	}
	return ChunkHeader{
		position: position,
		length:   length,
		name:     name,
	}, nil
}

// Position returns the byte offset of this chunk in the file.
func (h ChunkHeader) Position() int64 { return h.position }

// Length returns the number of bytes in the chunk data (excluding header and CRC).
func (h ChunkHeader) Length() int { return h.length }

// Name returns the four-byte chunk type identifier as a string.
func (h ChunkHeader) Name() string { return h.name }

// IsCritical returns true if this is a critical chunk that must be understood by any PNG decoder.
func (h ChunkHeader) IsCritical() bool {
	return len(h.name) > 0 && h.name[0] >= 'A' && h.name[0] <= 'Z'
}

// IsPublic returns true if this is a public chunk registered with the IANA.
func (h ChunkHeader) IsPublic() bool {
	return len(h.name) > 1 && h.name[1] >= 'A' && h.name[1] <= 'Z'
}

// IsSafeToCopy returns true if this ancillary chunk can be safely copied when the image is modified.
func (h ChunkHeader) IsSafeToCopy() bool {
	return len(h.name) > 3 && h.name[3] >= 'A' && h.name[3] <= 'Z'
}

// String returns a human-readable representation of the chunk header.
func (h ChunkHeader) String() string {
	return fmt.Sprintf("%s at %d (length: %d)", h.name, h.position, h.length)
}
