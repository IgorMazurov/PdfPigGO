package cff

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// CompactFontFormatFontCollection holds a CFF font program that may contain
// multiple fonts, achieving compression by sharing details between fonts in the set.
type CompactFontFormatFontCollection struct {
	header CompactFontFormatHeader
	fonts  map[string]*CompactFontFormatFont
	first  *CompactFontFormatFont
}

// NewCompactFontFormatFontCollection creates a new font collection with the given
// header and set of fonts keyed by name. The fontSet must not be nil.
func NewCompactFontFormatFontCollection(header CompactFontFormatHeader, fontSet map[string]*CompactFontFormatFont) (*CompactFontFormatFontCollection, error) {
	if fontSet == nil {
		return nil, fmt.Errorf("font set cannot be nil")
	}

	coll := &CompactFontFormatFontCollection{
		header: header,
		fonts:  fontSet,
	}

	for _, f := range fontSet {
		coll.first = f
		break
	}

	return coll, nil
}

// Header returns the decoded header table for this font collection.
func (c *CompactFontFormatFontCollection) Header() CompactFontFormatHeader {
	return c.header
}

// Fonts returns the individual fonts contained in this collection, keyed by name.
func (c *CompactFontFormatFontCollection) Fonts() map[string]*CompactFontFormatFont {
	return c.fonts
}

// FirstFont returns the first font contained in the collection.
func (c *CompactFontFormatFontCollection) FirstFont() *CompactFontFormatFont {
	return c.first
}

// GetFirstTransformationMatrix returns the transformation matrix of the first font
// in the collection, or the identity matrix if the collection is empty.
func (c *CompactFontFormatFontCollection) GetFirstTransformationMatrix() core.TransformationMatrix {
	for _, f := range c.fonts {
		return f.FontMatrix()
	}

	return core.Identity
}

// GetCharacterBoundingBox returns the bounding box for a character if the font
// contains a corresponding glyph, or nil otherwise.
func (c *CompactFontFormatFontCollection) GetCharacterBoundingBox(characterName string) *core.PdfRectangle {
	if c.first == nil {
		return nil
	}
	return c.first.GetCharacterBoundingBox(characterName)
}

// GetCharacterName returns the name for the character with the given character code.
// If isCid is true, it performs a CID lookup. Returns ".notdef" if no name is found.
func (c *CompactFontFormatFontCollection) GetCharacterName(characterCode int, isCid bool) string {
	if c.first == nil {
		return notDefined
	}

	name := c.first.GetCharacterName(characterCode, isCid)
	if name == "" {
		return notDefined
	}
	return name
}

const notDefined = ".notdef"
