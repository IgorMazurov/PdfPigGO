package truetype

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/core"
)

// TrueTypeHeaderTable represents a table directory entry from the TrueType font file.
// It indicates the position of the corresponding table data in the TrueType font.
type TrueTypeHeaderTable struct {
	Tag      string
	CheckSum uint32
	Offset   uint32
	Length   uint32
}

// Required table tags.

const (
	// Cmap is the character to glyph mapping table tag.
	Cmap = "cmap"
	// Glyf is the glyph data table tag.
	Glyf = "glyf"
	// Head is the font header table tag.
	Head = "head"
	// Hhea is the horizontal header table tag.
	Hhea = "hhea"
	// Hmtx is the horizontal metrics table tag.
	Hmtx = "hmtx"
	// Loca is the index to location table tag.
	Loca = "loca"
	// Maxp is the maximum profile table tag.
	Maxp = "maxp"
	// Name is the naming table tag.
	Name = "name"
	// Post is the PostScript information table tag.
	Post = "post"
	// Os2 is the OS/2 and Windows specific metrics table tag.
	Os2 = "OS/2"
)

// Optional table tags.

const (
	// Cvt is the control value table tag.
	Cvt = "cvt "
	// Ebdt is the embedded bitmap data table tag.
	Ebdt = "EBDT"
	// Eblc is the embedded bitmap location data table tag.
	Eblc = "EBLC"
	// Ebsc is the embedded bitmap scaling data table tag.
	Ebsc = "EBSC"
	// Fpgm is the font program table tag.
	Fpgm = "fpgm"
	// Gasp is the grid-fitting and scan conversion procedure (grayscale) table tag.
	Gasp = "gasp"
	// Hdmx is the horizontal device metrics table tag.
	Hdmx = "hdmx"
	// Kern is the kerning table tag.
	Kern = "kern"
	// Ltsh is the linear threshold title table tag.
	Ltsh = "LTSH"
	// Prep is the CVT program table tag.
	Prep = "prep"
	// Pclt is the PCL5 table tag.
	Pclt = "PCLT"
	// Vdmx is the vertical device metrics table tag.
	Vdmx = "VDMX"
	// Vhea is the vertical metrics header table tag.
	Vhea = "vhea"
	// Vmtx is the vertical metrics table tag.
	Vmtx = "vmtx"
)

// PostScript table tags.

const (
	// Cff is the compact font format table tag. It contains a Compact Font Format
	// font representation (also known as a PostScript Type 1, or CIDFont).
	Cff = "CFF "
)

// NewTrueTypeHeaderTable creates a new TrueTypeHeaderTable with validation.
func NewTrueTypeHeaderTable(tag string, checkSum, offset, length uint32) (TrueTypeHeaderTable, error) {
	if tag == "" {
		return TrueTypeHeaderTable{}, fmt.Errorf("tag must not be empty")
	}

	if len(tag) != 4 {
		return TrueTypeHeaderTable{}, fmt.Errorf("a TrueType table tag must be a uint32, 4 bytes long, instead got: %s", tag)
	}

	return TrueTypeHeaderTable{
		Tag:      tag,
		CheckSum: checkSum,
		Offset:   offset,
		Length:   length,
	}, nil
}

// GetEmptyHeaderTable returns an empty header table with the non-tag values set to zero.
func GetEmptyHeaderTable(tag string) (TrueTypeHeaderTable, error) {
	return NewTrueTypeHeaderTable(tag, 0, 0, 0)
}

// Write writes the TrueType header table to the output stream.
func (h TrueTypeHeaderTable) Write(w io.Writer) error {
	for i := 0; i < len(h.Tag); i++ {
		if _, err := w.Write([]byte{h.Tag[i]}); err != nil {
			return err
		}
	}

	if _, err := core.WriteUInt(w, h.CheckSum); err != nil {
		return err
	}

	if _, err := core.WriteUInt(w, h.Offset); err != nil {
		return err
	}

	if _, err := core.WriteUInt(w, h.Length); err != nil {
		return err
	}

	return nil
}

// String returns a string representation of the TrueType header table.
func (h TrueTypeHeaderTable) String() string {
	return fmt.Sprintf("%s - Offset: %d Length: %d Checksum: %d", h.Tag, h.Offset, h.Length, h.CheckSum)
}
