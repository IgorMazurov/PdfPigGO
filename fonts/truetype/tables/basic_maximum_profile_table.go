package tables

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// BasicMaximumProfileTable establishes the memory requirements for a TrueType font.
type BasicMaximumProfileTable struct {
	directoryTable truetype.TrueTypeHeaderTable
	version        float32
	numberOfGlyphs int
}

// NewBasicMaximumProfileTable creates a new BasicMaximumProfileTable.
func NewBasicMaximumProfileTable(
	directoryTable truetype.TrueTypeHeaderTable,
	version float32,
	numberOfGlyphs int,
) BasicMaximumProfileTable {
	return BasicMaximumProfileTable{
		directoryTable: directoryTable,
		version:        version,
		numberOfGlyphs: numberOfGlyphs,
	}
}

// Tag returns the 4-letter identifier for this table.
func (m BasicMaximumProfileTable) Tag() string {
	return truetype.Maxp
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (m BasicMaximumProfileTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return m.directoryTable
}

// Version returns the table version number. CFF fonts use version 0.5, TrueType uses version 1.
func (m BasicMaximumProfileTable) Version() float32 {
	return m.version
}

// NumberOfGlyphs returns the number of glyphs in the font.
func (m BasicMaximumProfileTable) NumberOfGlyphs() int {
	return m.numberOfGlyphs
}

// IsCompressedFontFormat returns true if the table version is 0.5 (CFF fonts).
func (m BasicMaximumProfileTable) IsCompressedFontFormat() bool {
	return m.version == 0.5
}

// MaximumProfileTable extends BasicMaximumProfileTable with additional memory
// requirements for TrueType fonts (version 1.0).
type MaximumProfileTable struct {
	BasicMaximumProfileTable
	maximumPoints             int
	maximumContours           int
	maximumCompositePoints    int
	maximumCompositeContours  int
	maximumZones              int
	maximumTwilightPoints     int
	maximumStorage            int
	maximumFunctionDefinitions int
	maximumInstructionDefs    int
	maximumStackElements      int
	maximumSizeOfInstructions int
	maximumComponentElements  int
	maximumComponentDepth     int
}

// NewMaximumProfileTable creates a new MaximumProfileTable.
func NewMaximumProfileTable(
	directoryTable truetype.TrueTypeHeaderTable,
	version float32,
	numberOfGlyphs int,
	maximumPoints int,
	maximumContours int,
	maximumCompositePoints int,
	maximumCompositeContours int,
	maximumZones int,
	maximumTwilightPoints int,
	maximumStorage int,
	maximumFunctionDefinitions int,
	maximumInstructionDefs int,
	maximumStackElements int,
	maximumSizeOfInstructions int,
	maximumComponentElements int,
	maximumComponentDepth int,
) MaximumProfileTable {
	return MaximumProfileTable{
		BasicMaximumProfileTable: NewBasicMaximumProfileTable(directoryTable, version, numberOfGlyphs),
		maximumPoints:            maximumPoints,
		maximumContours:          maximumContours,
		maximumCompositePoints:   maximumCompositePoints,
		maximumCompositeContours: maximumCompositeContours,
		maximumZones:             maximumZones,
		maximumTwilightPoints:    maximumTwilightPoints,
		maximumStorage:           maximumStorage,
		maximumFunctionDefinitions: maximumFunctionDefinitions,
		maximumInstructionDefs:   maximumInstructionDefs,
		maximumStackElements:     maximumStackElements,
		maximumSizeOfInstructions: maximumSizeOfInstructions,
		maximumComponentElements: maximumComponentElements,
		maximumComponentDepth:    maximumComponentDepth,
	}
}

// MaximumPoints returns the maximum number of points in a non-composite glyph.
func (m MaximumProfileTable) MaximumPoints() int {
	return m.maximumPoints
}

// MaximumContours returns the maximum number of contours in a non-composite glyph.
func (m MaximumProfileTable) MaximumContours() int {
	return m.maximumContours
}

// MaximumCompositePoints returns the maximum number of points in a composite glyph.
func (m MaximumProfileTable) MaximumCompositePoints() int {
	return m.maximumCompositePoints
}

// MaximumCompositeContours returns the maximum number of contours in a composite glyph.
func (m MaximumProfileTable) MaximumCompositeContours() int {
	return m.maximumCompositeContours
}

// MaximumZones returns the number of zones used by instructions (1 or 2).
func (m MaximumProfileTable) MaximumZones() int {
	return m.maximumZones
}

// MaximumTwilightPoints returns the maximum number of points in Z0 (twilight zone).
func (m MaximumProfileTable) MaximumTwilightPoints() int {
	return m.maximumTwilightPoints
}

// MaximumStorage returns the maximum number of storage area locations.
func (m MaximumProfileTable) MaximumStorage() int {
	return m.maximumStorage
}

// MaximumFunctionDefinitions returns the maximum number of function definitions.
func (m MaximumProfileTable) MaximumFunctionDefinitions() int {
	return m.maximumFunctionDefinitions
}

// MaximumInstructionDefinitions returns the maximum number of instruction definitions.
func (m MaximumProfileTable) MaximumInstructionDefinitions() int {
	return m.maximumInstructionDefs
}

// MaximumStackElements returns the maximum stack depth.
func (m MaximumProfileTable) MaximumStackElements() int {
	return m.maximumStackElements
}

// MaximumSizeOfInstructions returns the maximum byte count for glyph instructions.
func (m MaximumProfileTable) MaximumSizeOfInstructions() int {
	return m.maximumSizeOfInstructions
}

// MaximumComponentElements returns the maximum number of components at the top level for a composite glyph.
func (m MaximumProfileTable) MaximumComponentElements() int {
	return m.maximumComponentElements
}

// MaximumComponentDepth returns the maximum level of recursion (1 for simple components).
func (m MaximumProfileTable) MaximumComponentDepth() int {
	return m.maximumComponentDepth
}

var _ TrueTypeTable = MaximumProfileTable{}
