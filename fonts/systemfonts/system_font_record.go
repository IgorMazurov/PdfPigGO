package systemfonts

import (
	"errors"
	"strings"
)

// SystemFontRecord holds metadata about a discovered system font file.
type SystemFontRecord struct {
	path string
	typ  SystemFontType
}

// Path returns the file path of the system font.
func (r SystemFontRecord) Path() string {
	return r.path
}

// Type returns the font type determined from the file extension.
func (r SystemFontRecord) Type() SystemFontType {
	return r.typ
}

// NewSystemFontRecord creates a new SystemFontRecord with the given path and type.
// Returns an error if path is empty.
func NewSystemFontRecord(path string, typ SystemFontType) (SystemFontRecord, error) {
	if path == "" {
		return SystemFontRecord{}, errors.New("path must not be empty")
	}

	return SystemFontRecord{
		path: path,
		typ:  typ,
	}, nil
}

// TryCreate attempts to create a SystemFontRecord from a file path by inspecting its extension.
// Returns the record and true if the extension is recognized, or zero value and false otherwise.
func TryCreate(path string) (SystemFontRecord, bool) {
	if path == "" {
		return SystemFontRecord{}, false
	}

	lower := strings.ToLower(path)
	var fontType SystemFontType

	switch {
	case strings.HasSuffix(lower, ".ttf"):
		fontType = TrueType
	case strings.HasSuffix(lower, ".otf"):
		fontType = OpenType
	case strings.HasSuffix(lower, ".ttc"):
		fontType = TrueTypeCollection
	case strings.HasSuffix(lower, ".otc"):
		fontType = OpenTypeCollection
	case strings.HasSuffix(lower, ".pfb"):
		fontType = Type1
	default:
		return SystemFontRecord{}, false
	}

	record, _ := NewSystemFontRecord(path, fontType)
	return record, true
}
