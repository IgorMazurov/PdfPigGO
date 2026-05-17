package filestructure

import "fmt"

// FileHeaderOffset represents how many bytes precede the "%PDF-" version
// header in the file. In some files this junk can offset all following
// offset bytes.
type FileHeaderOffset struct {
	Value int
}

// NewFileHeaderOffset creates a new FileHeaderOffset with the given value.
func NewFileHeaderOffset(value int) FileHeaderOffset {
	return FileHeaderOffset{Value: value}
}

// String returns the string representation of the offset value.
func (f FileHeaderOffset) String() string {
	return fmt.Sprintf("%d", f.Value)
}

// Equals reports whether f and other represent the same offset.
func (f FileHeaderOffset) Equals(other FileHeaderOffset) bool {
	return f.Value == other.Value
}
