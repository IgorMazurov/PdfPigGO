package tables

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// TrueTypeTable represents a table in a TrueType font.
type TrueTypeTable interface {
	// Tag returns the 4-letter identifier for this table.
	Tag() string

	// DirectoryTable returns the directory entry from the font's offset
	// subtable, indicating the length, offset, type and checksum of the table.
	DirectoryTable() truetype.TrueTypeHeaderTable
}
