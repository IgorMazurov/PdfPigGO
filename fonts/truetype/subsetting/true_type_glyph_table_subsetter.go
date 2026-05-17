package subsetting

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/glyphs"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

// OldToNewGlyphIndex maps a glyph index in the original font to its new index
// in the subset file.
type OldToNewGlyphIndex struct {
	OldIndex uint16

	NewIndex byte

	Represents rune
}

// String returns a human-readable description of the mapping.
func (m OldToNewGlyphIndex) String() string {
	return fmt.Sprintf("%c: From %d To %d.", m.Represents, m.OldIndex, m.NewIndex)
}

var errFontBytesNil = fmt.Errorf("font bytes must not be nil")

var errMappingNil = fmt.Errorf("mapping must not be nil")

// glyphType distinguishes between empty, simple, and composite glyphs.
type glyphType int

const (
	glyphEmpty     glyphType = iota
	glyphSimple    glyphType = iota
	glyphComposite glyphType = iota
)

// glyphRecord holds metadata about a single glyph in the original font.
type glyphRecord struct {
	offset            int
	typ               glyphType
	dataLength        int
	dependencyIndices []compositeGlyphIndexReference
}

// compositeGlyphIndexReference marks a glyph index referenced by a composite glyph.
type compositeGlyphIndexReference struct {
	index                   uint32
	offsetOfIndexWithinData uint32
}

// SubsetGlyphTable creates a new glyph table from the input font which contains
// only the glyphs required by the input mapping. The tables are extracted from
// the TableRegister to avoid circular imports.
func SubsetGlyphTable(
	fontBytes []byte,
	mapping []OldToNewGlyphIndex,
	horizontalMetricsTable tables.HorizontalMetricsTable,
	indexToLocationTable tables.IndexToLocationTable,
	glyphDataTable *tables.GlyphDataTable,
) (*TrueTypeSubsetGlyphTable, error) {
	if fontBytes == nil {
		return nil, errFontBytesNil
	}

	if mapping == nil {
		return nil, errMappingNil
	}

	if horizontalMetricsTable.Tag() == "" {
		return nil, fmt.Errorf("font did not contain a horizontal metrics table, cannot subset")
	}

	data := truetypeparser.NewTrueTypeDataBytes(fontBytes)

	existingGlyphs := getGlyphRecordsInFont(data, indexToLocationTable, glyphDataTable)

	glyphsToCopy := make([]glyphRecord, 0, len(mapping))
	glyphsToCopyOriginalIndex := make([]int, 0, len(mapping))

	for _, m := range mapping {
		record := existingGlyphs[m.OldIndex]
		glyphsToCopy = append(glyphsToCopy, record)
		glyphsToCopyOriginalIndex = append(glyphsToCopyOriginalIndex, int(m.OldIndex))
	}

	glyphLocations := make([]uint32, 0)
	advanceWidths := make([]glyphs.HorizontalMetric, 0)

	writer := core.NewArrayPoolBufferWriter()
	defer writer.Dispose()

	for i := 0; i < len(glyphsToCopy); i++ {
		compositeIndicesToReplace := [][2]uint32{}

		newRecord := glyphsToCopy[i]

		if newRecord.typ == glyphComposite {
			for j := 0; j < len(newRecord.dependencyIndices); j++ {
				dependency := newRecord.dependencyIndices[j]

				newDependencyIndex := getAlreadyCopiedDependencyIndex(dependency.index, glyphsToCopyOriginalIndex)

				if newDependencyIndex < 0 {
					actualDependencyRecord := existingGlyphs[dependency.index]
					newDependencyIndex = len(glyphsToCopy)
					glyphsToCopy = append(glyphsToCopy, actualDependencyRecord)
					glyphsToCopyOriginalIndex = append(glyphsToCopyOriginalIndex, int(dependency.index))
				}

				withinGlyphDataIndexOffset := dependency.offsetOfIndexWithinData - uint32(newRecord.offset)

				compositeIndicesToReplace = append(compositeIndicesToReplace, [2]uint32{withinGlyphDataIndexOffset, uint32(newDependencyIndex)})
			}
		}

		glyphLocations = append(glyphLocations, uint32(writer.WrittenCount()))

		advanceWidth := horizontalMetricsTable.GetHorizontalMetric(glyphsToCopyOriginalIndex[i])
		advanceWidths = append(advanceWidths, advanceWidth)

		if newRecord.typ == glyphEmpty {
			continue
		}

		data.Seek(int64(newRecord.offset), 0) // io.SeekStart

		glyphBytes := data.ReadByteArray(newRecord.dataLength)

		for _, pair := range compositeIndicesToReplace {
			offset, newIndex := pair[0], pair[1]
			glyphBytes[offset] = byte(newIndex >> 8)
			glyphBytes[offset+1] = byte(newIndex)
		}

		writer.WriteBytes(glyphBytes)

		remainder := len(glyphBytes) % 4
		bytesToPad := 0
		if remainder != 0 {
			bytesToPad = 4 - remainder
		}
		for j := 0; j < bytesToPad; j++ {
			writer.WriteSingle(0)
		}
	}

	output := writer.WrittenSpan()
	result := make([]byte, len(output))
	copy(result, output)

	glyphLocations = append(glyphLocations, uint32(len(result)))

	return NewTrueTypeSubsetGlyphTable(result, glyphLocations, advanceWidths), nil
}

func getGlyphRecordsInFont(
	data *truetypeparser.TrueTypeDataBytes,
	indexToLocationTable tables.IndexToLocationTable,
	glyphDataTable *tables.GlyphDataTable,
) []glyphRecord {
	numGlyphs := len(indexToLocationTable.GlyphOffsets()) - 1

	glyphDirectory := glyphDataTable.DirectoryTable()

	data.Seek(int64(glyphDirectory.Offset), 0) // io.SeekStart

	glyphRecords := make([]glyphRecord, numGlyphs)

	for i := 0; i < numGlyphs; i++ {
		glyphOffsets := indexToLocationTable.GlyphOffsets()
		glyphOffset := int(glyphDirectory.Offset + glyphOffsets[i])

		if glyphOffsets[i+1] <= glyphOffsets[i] {
			glyphRecords[i] = glyphRecord{offset: glyphOffset, typ: glyphEmpty, dataLength: 0}
			continue
		}

		data.Seek(int64(glyphOffset), 0) // io.SeekStart

		if glyphOffset >= int(glyphDirectory.Offset+glyphDirectory.Length) {
			panic("failed to read expected number of glyphs before reaching end of input")
		}

		numberOfContours := data.ReadSignedShort()
		typ := glyphSimple
		if numberOfContours < 0 {
			typ = glyphComposite
		}

		data.ReadSignedShort() // xMin, unused
		data.ReadSignedShort() // yMin, unused
		data.ReadSignedShort() // xMax, unused
		data.ReadSignedShort() // yMax, unused

		if typ == glyphSimple {
			readSimpleGlyph(data, int(numberOfContours))
			glyphRecords[i] = glyphRecord{offset: glyphOffset, typ: typ, dataLength: int(data.Position()) - glyphOffset}
		} else {
			glyphIndices := readCompositeGlyph(data)

			next := glyphOffsets[i+1]
			data.Seek(int64(glyphDirectory.Offset+next)-1, 0) // io.SeekStart

			glyphRecords[i] = glyphRecord{offset: glyphOffset, typ: typ, dataLength: int(data.Position()) - glyphOffset, dependencyIndices: glyphIndices}
		}
	}

	return glyphRecords
}

func getAlreadyCopiedDependencyIndex(dependencyIndex uint32, copiedGlyphOriginalIndices []int) int {
	for i, originalIndex := range copiedGlyphOriginalIndices {
		if originalIndex == int(dependencyIndex) {
			return i
		}
	}
	return -1
}

func readSimpleGlyph(data *truetypeparser.TrueTypeDataBytes, numberOfContours int) {
	if numberOfContours == 0 {
		return
	}

	endPointsOfContours := make([]uint16, numberOfContours)
	for i := 0; i < numberOfContours; i++ {
		endPointsOfContours[i] = data.ReadUnsignedShort()
	}

	instructionLength := int(data.ReadUnsignedShort())

	for i := 0; i < instructionLength; i++ {
		_, _ = data.ReadByte()
	}

	lastPointIndex := endPointsOfContours[numberOfContours-1]
	pointCount := int(lastPointIndex) + 1

	perPointFlags := make([]glyphs.SimpleGlyphFlags, pointCount)
	for i := 0; i < pointCount; {
		b, _ := data.ReadByte()
		flags := glyphs.SimpleGlyphFlags(b)
		perPointFlags[i] = flags

		if flags&glyphs.Repeat == 0 {
			i++
			continue
		}

		nr, _ := data.ReadByte()
		numberOfRepeats := int(nr)
		for r := 0; r < numberOfRepeats; r++ {
			i++
			perPointFlags[i] = flags
		}
		i++
	}

	readCoordinates(perPointFlags, data, glyphs.XSingleByte, glyphs.ThisXIsTheSame)
	readCoordinates(perPointFlags, data, glyphs.YSingleByte, glyphs.ThisYIsTheSame)
}

func readCoordinates(
	flags []glyphs.SimpleGlyphFlags,
	data *truetypeparser.TrueTypeDataBytes,
	isSingleByte, isTheSameAsPrevious glyphs.SimpleGlyphFlags,
) []int16 {
	coordinates := make([]int16, len(flags))
	value := 0

	for i, flag := range flags {
		if flag&isSingleByte != 0 {
			rb, _ := data.ReadByte()
			b := int(rb)

			if flag&isTheSameAsPrevious != 0 {
				value += b
			} else {
				value -= b
			}
		} else {
			var delta int16

			if flag&isTheSameAsPrevious != 0 {
				delta = 0
			} else {
				delta = data.ReadSignedShort()
			}

			value += int(delta)
		}

		coordinates[i] = int16(value)
	}

	return coordinates
}

func readCompositeGlyph(data *truetypeparser.TrueTypeDataBytes) []compositeGlyphIndexReference {
	var glyphIndices []compositeGlyphIndexReference

	for {
		flags := glyphs.CompositeGlyphFlags(data.ReadUnsignedShort())
		indexOffset := uint32(data.Position())
		glyphIndex := data.ReadUnsignedShort()
		glyphIndices = append(glyphIndices, compositeGlyphIndexReference{index: uint32(glyphIndex), offsetOfIndexWithinData: indexOffset})

		if flags&glyphs.Args1And2AreWords != 0 {
			data.ReadSignedShort()
			data.ReadSignedShort()
		} else {
		_, _ = data.ReadByte()
		_, _ = data.ReadByte()
	}

		if flags&glyphs.WeHaveAScale != 0 {
			data.ReadSignedShort()
		} else if flags&glyphs.WeHaveAnXAndYScale != 0 {
			data.ReadSignedShort()
			data.ReadSignedShort()
		} else if flags&glyphs.WeHaveATwoByTwo != 0 {
			data.ReadSignedShort()
			data.ReadSignedShort()
			data.ReadSignedShort()
			data.ReadSignedShort()
		}

		if flags&glyphs.MoreComponents == 0 {
			break
		}
	}

	return glyphIndices
}
