package cff

import "fmt"

// CompactFontFormatHeader is the header table for the binary data of a
// Compact Font Format file.
type CompactFontFormatHeader struct {
	// MajorVersion is the major version of this font format, starting at 1.
	MajorVersion byte

	// MinorVersion is the minor version of this font format, starting at 0.
	// It indicates extensions to the format which are undetectable by readers
	// that do not support them.
	MinorVersion byte

	// SizeInBytes indicates the size of this header in bytes so that future
	// changes to the format may include extra data after the OffsetSize field.
	SizeInBytes byte

	// OffsetSize specifies the size of all offsets relative to the start of
	// the data in the font.
	OffsetSize byte
}

// NewCompactFontFormatHeader creates a new CompactFontFormatHeader.
func NewCompactFontFormatHeader(majorVersion, minorVersion, sizeInBytes, offsetSize byte) CompactFontFormatHeader {
	return CompactFontFormatHeader{
		MajorVersion: majorVersion,
		MinorVersion: minorVersion,
		SizeInBytes:  sizeInBytes,
		OffsetSize:   offsetSize,
	}
}

// String returns a string representation of the header.
func (h CompactFontFormatHeader) String() string {
	return fmt.Sprintf("Major: %d, Minor: %d, Header Size: %d, Offset: %d",
		h.MajorVersion, h.MinorVersion, h.SizeInBytes, h.OffsetSize)
}
