package tables

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/names"
)

// NameTable allows multilingual strings to be associated with the TrueType font.
type NameTable struct {
	directoryTable   truetype.TrueTypeHeaderTable
	fontName         string
	fontFamilyName   string
	fontSubFamilyName string
	nameRecords      []names.TrueTypeNameRecord
}

// NewNameTable creates a new NameTable.
func NewNameTable(
	directoryTable truetype.TrueTypeHeaderTable,
	fontName, fontFamilyName, fontSubFamilyName string,
	nameRecords []names.TrueTypeNameRecord,
) *NameTable {
	return &NameTable{
		directoryTable:    directoryTable,
		fontName:          fontName,
		fontFamilyName:    fontFamilyName,
		fontSubFamilyName: fontSubFamilyName,
		nameRecords:       nameRecords,
	}
}

// Tag returns the 4-letter identifier for this table.
func (n *NameTable) Tag() string {
	return truetype.Name
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (n *NameTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return n.directoryTable
}

// FontName returns the font name.
func (n *NameTable) FontName() string {
	return n.fontName
}

// FontFamilyName returns the font family name.
func (n *NameTable) FontFamilyName() string {
	return n.fontFamilyName
}

// FontSubFamilyName returns the font sub-family name.
func (n *NameTable) FontSubFamilyName() string {
	return n.fontSubFamilyName
}

// NameRecords returns the name records contained in this name table.
func (n *NameTable) NameRecords() []names.TrueTypeNameRecord {
	return n.nameRecords
}

// GetPostscriptName gets the PostScript name for the font if specified,
// preferring the Windows platform name if present.
func (n *NameTable) GetPostscriptName() string {
	var anyName string
	for _, record := range n.nameRecords {
		if record.NameId != 6 {
			continue
		}

		if record.PlatformId == names.Windows {
			return record.Value
		}

		anyName = record.Value
	}

	return anyName
}
