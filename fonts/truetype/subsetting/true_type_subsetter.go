package subsetting

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/glyphs"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
	cmapsubtables "github.com/uglytoad/pdfpig/go/fonts/truetype/tables/cmap_sub_tables"
)

var requiredTags = []string{
	truetype.Cmap,
	truetype.Glyf,
	truetype.Head,
	truetype.Hhea,
	truetype.Hmtx,
	truetype.Loca,
	truetype.Maxp,
}

var optionalTags = []string{
	truetype.Cvt,
	truetype.Fpgm,
	truetype.Prep,
	truetype.Name,
}

var paddingBytes = []byte{0, 0, 0, 0}

// mutableBuffer provides an append-only byte buffer with seek support for
// patching previously written data.
type mutableBuffer struct {
	data []byte
	pos  int64
}

func newMutableBuffer() *mutableBuffer {
	return &mutableBuffer{data: make([]byte, 0, 4096)}
}

func (b *mutableBuffer) Write(p []byte) (int, error) {
	if b.pos >= int64(len(b.data)) {
		b.data = append(b.data, p...)
	} else {
		remaining := len(b.data) - int(b.pos)
		if remaining >= len(p) {
			copy(b.data[b.pos:], p)
		} else {
			copy(b.data[b.pos:], p[:remaining])
			b.data = append(b.data, p[remaining:]...)
		}
	}
	b.pos += int64(len(p))
	return len(p), nil
}

func (b *mutableBuffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.pos = offset
	case io.SeekCurrent:
		b.pos += offset
	case io.SeekEnd:
		b.pos = int64(len(b.data)) + offset
	}
	return b.pos, nil
}

func (b *mutableBuffer) Bytes() []byte {
	return b.data
}

// Subset generates a subset of the input font containing only the data required
// for the glyphs specified in the encoding.
func Subset(fontBytes []byte, newEncoding *TrueTypeSubsetEncoding) ([]byte, error) {
	if fontBytes == nil {
		return nil, fmt.Errorf("font bytes must not be nil")
	}

	if newEncoding == nil {
		return nil, fmt.Errorf("encoding must not be nil")
	}

	font := truetypeparser.Parse(truetypeparser.NewTrueTypeDataBytes(fontBytes))

	indexMapping, err := getIndexMapping(fontBytes, font, newEncoding)
	if err != nil {
		return nil, err
	}

	buf := newMutableBuffer()

	copiedTableTags := make(map[string]bool)
	for _, tag := range requiredTags {
		header, ok := font.TableHeaders()[tag]
		if !ok || header.Tag == "" {
			return nil, fmt.Errorf("font does not contain table required for subsetting: %s", tag)
		}
		copiedTableTags[tag] = true
	}

	for _, tag := range optionalTags {
		header, ok := font.TableHeaders()[tag]
		if !ok || header.Tag == "" {
			continue
		}
		copiedTableTags[tag] = true
	}

	sortedTags := make([]string, 0, len(copiedTableTags))
	for tag := range copiedTableTags {
		sortedTags = append(sortedTags, tag)
	}
	sort.Strings(sortedTags)

	offsetSubtable := NewTrueTypeOffsetSubtable(byte(len(sortedTags)))
	if err := offsetSubtable.Write(buf); err != nil {
		return nil, fmt.Errorf("failed to write offset subtable: %w", err)
	}

	directoryEntries := make([]*directoryEntry, len(sortedTags))

	for i, tag := range sortedTags {
		entry := newDirectoryEntry(tag, int64(len(buf.data)), font.TableHeaders()[tag])
		if err := entry.DummyHeader.Write(buf); err != nil {
			return nil, fmt.Errorf("failed to write dummy directory entry for %s: %w", tag, err)
		}
		directoryEntries[i] = entry
	}

	hmtxTable := parseHorizontalMetrics(fontBytes, font.TableHeaders()[truetype.Hmtx])
	indexToLocaTable := parseIndexToLocation(fontBytes, font.TableHeaders()[truetype.Loca], font.TableHeaders()[truetype.Head])

	trueTypeSubsetGlyphTable, err := SubsetGlyphTable(
		fontBytes, indexMapping,
		hmtxTable, indexToLocaTable, font.TableRegister().GlyphTable,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to subset glyph table: %w", err)
	}

	for i := 0; i < len(directoryEntries); i++ {
		entry := directoryEntries[i]
		entry.OutputTableOffset = int64(len(buf.data))

		switch entry.Tag {
		case truetype.Cmap:
			cmapTable := getCMapTable(entry, indexMapping)
			if err := cmapTable.Write(buf); err != nil {
				return nil, fmt.Errorf("failed to write cmap table: %w", err)
			}
		case truetype.Glyf:
			if _, err := buf.Write(trueTypeSubsetGlyphTable.Bytes); err != nil {
				return nil, fmt.Errorf("failed to write glyf table: %w", err)
			}
		case truetype.Hmtx:
			if err := writeHMTX(buf, trueTypeSubsetGlyphTable); err != nil {
				return nil, fmt.Errorf("failed to write hmtx table: %w", err)
			}
		case truetype.Loca:
			if err := writeLoca(buf, trueTypeSubsetGlyphTable); err != nil {
				return nil, fmt.Errorf("failed to write loca table: %w", err)
			}
		case truetype.Head:
			headBytes := getRawInputTableBytes(fontBytes, entry)
			writeUShort(headBytes, len(headBytes)-4, 1)
			if _, err := buf.Write(headBytes); err != nil {
				return nil, fmt.Errorf("failed to write head table: %w", err)
			}
		case truetype.Hhea:
			hheaBytes := getRawInputTableBytes(fontBytes, entry)
			writeUShort(hheaBytes, len(hheaBytes)-2, uint16(len(trueTypeSubsetGlyphTable.HorizontalMetrics)))
			if _, err := buf.Write(hheaBytes); err != nil {
				return nil, fmt.Errorf("failed to write hhea table: %w", err)
			}
		case truetype.Maxp:
			maxpBytes := getRawInputTableBytes(fontBytes, entry)
			writeUShort(maxpBytes, 4, uint16(trueTypeSubsetGlyphTable.GlyphCount()))
			if _, err := buf.Write(maxpBytes); err != nil {
				return nil, fmt.Errorf("failed to write maxp table: %w", err)
			}
		default:
			buffer := getRawInputTableBytes(fontBytes, entry)
			if _, err := buf.Write(buffer); err != nil {
				return nil, fmt.Errorf("failed to write table %s: %w", entry.Tag, err)
			}
		}

		entry.Length = uint32(int64(len(buf.data)) - entry.OutputTableOffset)

		remainder := len(buf.data) % 4
		if remainder > 0 {
			toAppend := 4 - remainder
			if _, err := buf.Write(paddingBytes[:toAppend]); err != nil {
				return nil, fmt.Errorf("failed to write padding after table %s: %w", entry.Tag, err)
			}
		}
	}

	resultBytes := buf.Bytes()
	inputBytes, err := core.NewStreamInputBytes(bytes.NewReader(resultBytes), false)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream input bytes: %w", err)
	}
	defer inputBytes.Close()

	for i := 0; i < len(directoryEntries); i++ {
		entry := directoryEntries[i]

		actualHeaderExceptChecksum, _ := truetype.NewTrueTypeHeaderTable(
			entry.Tag, 0, uint32(entry.OutputTableOffset), entry.Length)

		checksum, err := truetype.CalculateTable(inputBytes, actualHeaderExceptChecksum)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate checksum for table %s: %w", entry.Tag, err)
		}

		actualHeader, _ := truetype.NewTrueTypeHeaderTable(
			entry.Tag, checksum, uint32(entry.OutputTableOffset), entry.Length)

		if _, err := buf.Seek(entry.OutputEntryOffset, io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed to seek for updating directory entry %s: %w", entry.Tag, err)
		}

		if err := actualHeader.Write(buf); err != nil {
			return nil, fmt.Errorf("failed to write updated directory entry for %s: %w", entry.Tag, err)
		}
	}

	return buf.Bytes(), nil
}

func getIndexMapping(fontBytes []byte, font *truetypeparser.TrueTypeFont, newEncoding *TrueTypeSubsetEncoding) ([]OldToNewGlyphIndex, error) {
	characters := newEncoding.Characters()
	result := make([]OldToNewGlyphIndex, len(characters)+1)

	result[0] = OldToNewGlyphIndex{OldIndex: 0, NewIndex: 0, Represents: 0}

	cmapHeader, ok := font.TableHeaders()[truetype.Cmap]
	if !ok || cmapHeader.Tag == "" {
		return nil, fmt.Errorf("cannot subset font due to missing cmap table")
	}

	charToGlyph, err := parseCMapSubtablesForLookup(fontBytes, cmapHeader)
	if err != nil {
		return nil, fmt.Errorf("cannot subset font due to missing cmap subtables: %w", err)
	}

	for i, character := range characters {
		oldIndex := uint16(charToGlyph[int(character)])
		result[i+1] = OldToNewGlyphIndex{OldIndex: oldIndex, NewIndex: byte(i + 1), Represents: character}
	}

	return result, nil
}

func getCMapTable(entry *directoryEntry, encoding []OldToNewGlyphIndex) *byteEncodingCMapTable {
	data := make([]byte, 256)
	for i := 0; i < len(data); i++ {
		if i < len(encoding) {
			data[i] = encoding[i].NewIndex
		} else {
			break
		}
	}

	return &byteEncodingCMapTable{
		platformID:   cmapsubtables.Macintosh,
		encodingID:   0,
		formatOffset: 4,
		data:         data,
	}
}

func getRawInputTableBytes(font []byte, entry *directoryEntry) []byte {
	buffer := make([]byte, entry.InputHeader.Length)
	copy(buffer, font[entry.InputHeader.Offset:entry.InputHeader.Offset+entry.InputHeader.Length])
	return buffer
}

func writeUShort(arr []byte, offset int, value uint16) {
	arr[offset] = byte(value >> 8)
	arr[offset+1] = byte(value)
}

type directoryEntry struct {
	Tag               string
	OutputEntryOffset int64
	InputHeader       truetype.TrueTypeHeaderTable
	DummyHeader       truetype.TrueTypeHeaderTable
	OutputTableOffset int64
	Length            uint32
}

func newDirectoryEntry(tag string, outputEntryOffset int64, inputHeader truetype.TrueTypeHeaderTable) *directoryEntry {
	dummyHeader, _ := truetype.GetEmptyHeaderTable(tag)
	return &directoryEntry{
		Tag:               tag,
		OutputEntryOffset: outputEntryOffset,
		InputHeader:       inputHeader,
		DummyHeader:       dummyHeader,
	}
}

// byteEncodingCMapTable represents a format 0 cmap subtable for writing
// character code to glyph index mappings in the subset font.
type byteEncodingCMapTable struct {
	platformID   cmapsubtables.TrueTypeCMapPlatform
	encodingID   uint16
	formatOffset int
	data         []byte
}

func (b *byteEncodingCMapTable) Write(w io.Writer) error {
	if _, err := w.Write([]byte{0, 0}); err != nil {
		return fmt.Errorf("failed to write cmap version: %w", err)
	}

	if _, err := w.Write([]byte{0, 1}); err != nil {
		return fmt.Errorf("failed to write cmap numTables: %w", err)
	}

	if _, err := core.WriteUShort(w, uint16(b.platformID)); err != nil {
		return fmt.Errorf("failed to write platformID: %w", err)
	}

	if _, err := core.WriteUShort(w, b.encodingID); err != nil {
		return fmt.Errorf("failed to write encodingID: %w", err)
	}

	if _, err := core.WriteUInt(w, uint32(b.formatOffset)); err != nil {
		return fmt.Errorf("failed to write format offset: %w", err)
	}

	if _, err := w.Write([]byte{0, 0}); err != nil {
		return fmt.Errorf("failed to write format: %w", err)
	}

	if _, err := core.WriteUShort(w, 0); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	if _, err := w.Write([]byte{0}); err != nil {
		return fmt.Errorf("failed to write language: %w", err)
	}

	if _, err := w.Write(b.data); err != nil {
		return fmt.Errorf("failed to write glyphIDArray: %w", err)
	}

	return nil
}

// writeHMTX writes the horizontal metrics table bytes for the subsetted glyphs.
func writeHMTX(w io.Writer, glyphTable *TrueTypeSubsetGlyphTable) error {
	for _, m := range glyphTable.HorizontalMetrics {
		if _, err := core.WriteUShort(w, m.AdvanceWidth); err != nil {
			return fmt.Errorf("failed to write advance width: %w", err)
		}
		if _, err := core.WriteShort(w, m.LeftSideBearing); err != nil {
			return fmt.Errorf("failed to write left side bearing: %w", err)
		}
	}
	return nil
}

// writeLoca writes the index to location table bytes for the subsetted glyphs
// using long format (32-bit offsets).
func writeLoca(w io.Writer, glyphTable *TrueTypeSubsetGlyphTable) error {
	for _, offset := range glyphTable.GlyphOffsets {
		if _, err := core.WriteUInt(w, offset); err != nil {
			return fmt.Errorf("failed to write loca entry: %w", err)
		}
	}
	return nil
}

func parseCMapSubtablesForLookup(fontBytes []byte, cmapHeader truetype.TrueTypeHeaderTable) (map[int]uint16, error) {
	if cmapHeader.Length < 4 {
		return nil, fmt.Errorf("cmap table too small")
	}

	offset := int(cmapHeader.Offset)
	numTables := int(uint16(fontBytes[offset+2])<<8 | uint16(fontBytes[offset+3]))

	for i := 0; i < numTables; i++ {
		subtableOffset := offset + 4 + i*8
		if subtableOffset+8 > len(fontBytes) {
			continue
		}

		platformID := uint16(fontBytes[subtableOffset])<<8 | uint16(fontBytes[subtableOffset+1])
		encodingID := uint16(fontBytes[subtableOffset+2])<<8 | uint16(fontBytes[subtableOffset+3])
		subtableAddr := uint32(fontBytes[subtableOffset+4])<<24 |
			uint32(fontBytes[subtableOffset+5])<<16 |
			uint32(fontBytes[subtableOffset+6])<<8 |
			uint32(fontBytes[subtableOffset+7])

		subtableOffsetInFile := int(cmapHeader.Offset + subtableAddr)
		if subtableOffsetInFile+2 > len(fontBytes) {
			continue
		}

		format := uint16(fontBytes[subtableOffsetInFile])<<8 | uint16(fontBytes[subtableOffsetInFile+1])

		if format == 4 &&
			((platformID == 3 && encodingID == 1) ||
				(platformID == 3 && encodingID == 0) ||
				(platformID == 1 && encodingID == 0)) {
			return parseFormat4CMap(fontBytes, subtableOffsetInFile), nil
		}
	}

	return nil, fmt.Errorf("no suitable cmap subtable found")
}

func parseFormat4CMap(fontBytes []byte, offset int) map[int]uint16 {
	result := make(map[int]uint16)

	if offset+14 > len(fontBytes) {
		return result
	}

	length := int(uint16(fontBytes[offset+2])<<8 | uint16(fontBytes[offset+3]))
	if offset+length > len(fontBytes) {
		return result
	}

	searchRange := int(uint16(fontBytes[offset+4])<<8 | uint16(fontBytes[offset+5]))
	entrySelector := int(uint16(fontBytes[offset+6])<<8 | uint16(fontBytes[offset+7]))
	rangeShift := int(uint16(fontBytes[offset+8])<<8 | uint16(fontBytes[offset+9]))

	numEntries := (searchRange/2+rangeShift) / 6
	if entrySelector > 0 {
		numEntries = 1 << entrySelector
	}

	endCodesOffset := offset + 10
	startCodesOffset := endCodesOffset + numEntries*2
	idDeltasOffset := startCodesOffset + numEntries*2
	idRangeOffsetsOffset := idDeltasOffset + numEntries*2

	for i := 0; i < numEntries; i++ {
		idx := endCodesOffset + i*2
		if idx+1 >= len(fontBytes) {
			break
		}
		endCode := int(uint16(fontBytes[idx])<<8 | uint16(fontBytes[idx+1]))

		idx = startCodesOffset + i*2
		startCode := int(uint16(fontBytes[idx])<<8 | uint16(fontBytes[idx+1]))

		idx = idDeltasOffset + i*2
		delta := int(int16(uint16(fontBytes[idx])<<8 | uint16(fontBytes[idx+1])))

		idx = idRangeOffsetsOffset + i*2
		rangeOffset := int(uint16(fontBytes[idx])<<8 | uint16(fontBytes[idx+1]))

		for c := startCode; c <= endCode; c++ {
			var glyphIndex uint16
			if rangeOffset == 0 {
				glyphIndex = uint16((c + delta) & 0xFFFF)
			} else {
				innerIdx := int(rangeOffset/2+c-startCode) - (numEntries - i)
				innerOffset := idRangeOffsetsOffset + rangeOffset + innerIdx*2
				if innerOffset+1 < len(fontBytes) {
					glyphID := int(int16(uint16(fontBytes[innerOffset])<<8 | uint16(fontBytes[innerOffset+1])))
					glyphIndex = uint16((glyphID + delta) & 0xFFFF)
				} else {
					glyphIndex = uint16((c + delta) & 0xFFFF)
				}
			}
			result[c] = glyphIndex
		}
	}

	return result
}

func parseHorizontalMetrics(fontBytes []byte, hmtxHeader truetype.TrueTypeHeaderTable) tables.HorizontalMetricsTable {
	offset := int(hmtxHeader.Offset)
	length := int(hmtxHeader.Length)
	metrics := make([]glyphs.HorizontalMetric, 0, length/4)

	for i := 0; i+3 < length && offset+i*4+3 < len(fontBytes); i++ {
		idx := offset + i*4
		advanceWidth := uint16(fontBytes[idx])<<8 | uint16(fontBytes[idx+1])
		lsb := int16(uint16(fontBytes[idx+2])<<8 | uint16(fontBytes[idx+3]))
		metrics = append(metrics, glyphs.NewHorizontalMetric(advanceWidth, lsb))
	}

	dummyHeader, _ := truetype.GetEmptyHeaderTable(truetype.Hmtx)
	result, _ := tables.NewHorizontalMetricsTable(dummyHeader, metrics, []int16{})
	return result
}

func parseIndexToLocation(fontBytes []byte, locaHeader truetype.TrueTypeHeaderTable, headHeader truetype.TrueTypeHeaderTable) tables.IndexToLocationTable {
	offset := int(locaHeader.Offset)
	length := int(locaHeader.Length)

	isLong := true
	if headHeader.Length >= 54 {
		idx := int(headHeader.Offset + 50)
		if idx+1 < len(fontBytes) {
			format := uint16(fontBytes[idx])<<8 | uint16(fontBytes[idx+1])
			isLong = format == 1
		}
	}

	var numGlyphs int
	if isLong {
		numGlyphs = length / 4
	} else {
		numGlyphs = length/2 + 1
	}

	glyphOffsets := make([]uint32, numGlyphs)

	stride := 4
	if !isLong {
		stride = 2
	}

	for i := 0; i < numGlyphs; i++ {
		idx := offset + i*stride
		if isLong && idx+3 < len(fontBytes) {
			glyphOffsets[i] = uint32(fontBytes[idx])<<24 |
				uint32(fontBytes[idx+1])<<16 |
				uint32(fontBytes[idx+2])<<8 |
				uint32(fontBytes[idx+3])
		} else if idx+1 < len(fontBytes) {
			glyphOffsets[i] = uint32(uint16(fontBytes[idx])<<8 | uint16(fontBytes[idx+1])) * 2
		}
	}

	dummyHeader, _ := truetype.GetEmptyHeaderTable(truetype.Loca)
	return tables.NewIndexToLocationTable(dummyHeader, tables.IndexToLocationTableLong, glyphOffsets)
}