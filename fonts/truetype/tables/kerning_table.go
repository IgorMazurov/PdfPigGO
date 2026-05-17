package tables

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables/kerning"
)

// KerningTable holds the kerning data for a TrueType font.
type KerningTable struct {
	kerningTables []*kerning.KerningSubTable
}

// NewKerningTable creates a new KerningTable from sub-tables, filtering out nil entries.
func NewKerningTable(kerningTables []*kerning.KerningSubTable) *KerningTable {
	notNull := make([]*kerning.KerningSubTable, 0, len(kerningTables))
	for _, kt := range kerningTables {
		if kt != nil {
			notNull = append(notNull, kt)
		}
	}
	return &KerningTable{kerningTables: notNull}
}

// KerningTables returns the list of kerning sub-tables.
func (k *KerningTable) KerningTables() []*kerning.KerningSubTable {
	return k.kerningTables
}
