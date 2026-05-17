package tables

import (
	"fmt"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/glyphs"
)

// GlyphDataTable contains the data that defines the appearance of the glyphs in
// the font. This includes specification of the points that describe the contours
// that make up a glyph outline and the instructions that grid-fit that glyph.
type GlyphDataTable struct {
	directoryTable truetype.TrueTypeHeaderTable
	glyphOffsets   []uint32
	maxGlyphBounds core.PdfRectangle
	tableBytes     []byte
	glyphsOnce     sync.Once
	glyphsValue    []glyphs.GlyphDescription
}

// NewGlyphDataTable creates a new GlyphDataTable.
func NewGlyphDataTable(
	directoryTable truetype.TrueTypeHeaderTable,
	glyphOffsets []uint32,
	maxGlyphBounds core.PdfRectangle,
	tableBytes []byte,
) GlyphDataTable {
	return GlyphDataTable{
		directoryTable: directoryTable,
		glyphOffsets:   glyphOffsets,
		maxGlyphBounds: maxGlyphBounds,
		tableBytes:     tableBytes,
	}
}

// Tag returns the 4-letter identifier for this table.
func (g *GlyphDataTable) Tag() string {
	return truetype.Glyf
}

// DirectoryTable returns the directory entry from the font's offset subtable.
func (g *GlyphDataTable) DirectoryTable() truetype.TrueTypeHeaderTable {
	return g.directoryTable
}

// Glyphs returns the lazily-evaluated list of glyph descriptions.
func (g *GlyphDataTable) Glyphs() []glyphs.GlyphDescription {
	g.glyphsOnce.Do(func() {
		g.glyphsValue = g.readGlyphs()
	})
	return g.glyphsValue
}

// TryGetGlyphBounds attempts to get the bounding rectangle for a glyph at the given index.
func (g *GlyphDataTable) TryGetGlyphBounds(glyphIndex int) (core.PdfRectangle, bool) {
	if glyphIndex < 0 || glyphIndex >= len(g.glyphOffsets)-1 {
		return core.PdfRectangle{}, false
	}

	offset := g.glyphOffsets[glyphIndex]
	nextOffset := g.glyphOffsets[glyphIndex+1]

	if nextOffset <= offset {
		return core.NewPdfRectangleFromInt(0, 0, 0, 0), true
	}

	allGlyphs := g.Glyphs()
	return allGlyphs[glyphIndex].Bounds(), true
}

// TryGetGlyphPath attempts to get the subpaths for a glyph at the given index.
func (g *GlyphDataTable) TryGetGlyphPath(glyphIndex int) ([]*core.PdfSubpath, bool) {
	allGlyphs := g.Glyphs()

	if glyphIndex < 0 || glyphIndex >= len(allGlyphs) {
		return nil, false
	}

	return allGlyphs[glyphIndex].TryGetGlyphPath()
}

// --- private helpers (mirroring C# ReadGlyphs and related methods) ---

type temporaryCompositeLocation struct {
	position int64
	bounds   core.PdfRectangle
}

type compositeComponent struct {
	index          int
	transformation glyphs.CompositeTransformMatrix3By2
}

func (g *GlyphDataTable) readGlyphs() []glyphs.GlyphDescription {
	if g.tableBytes == nil {
		return nil
	}

	data := NewTrueTypeDataBytesLocal(g.tableBytes)
	offsets := g.glyphOffsets
	entryCount := len(offsets)
	glyphCount := entryCount - 1

	result := make([]glyphs.GlyphDescription, glyphCount)
	emptyGlyph := glyphs.EmptyGlyph(g.maxGlyphBounds)
	compositeLocations := make(map[int]temporaryCompositeLocation)

	for i := 0; i < glyphCount; i++ {
		offset := offsets[i]

		if offsets[i+1] <= offset {
			result[i] = emptyGlyph
			continue
		}

		if int64(offset) >= data.Length() {
			result[i] = emptyGlyph
			continue
		}

		data.Seek(int64(offset), 0) // io.SeekStart
		contourCount := data.ReadSignedShort()
		minX := data.ReadSignedShort()
		minY := data.ReadSignedShort()
		maxX := data.ReadSignedShort()
		maxY := data.ReadSignedShort()

		bounds := core.NewPdfRectangleFromInt(int(minX), int(minY), int(maxX), int(maxY))

		if contourCount >= 0 {
			result[i] = readSimpleGlyph(data, contourCount, bounds)
		} else {
			compositeLocations[i] = temporaryCompositeLocation{
				position: data.Position(),
				bounds:   bounds,
			}
		}
	}

	for idx, loc := range compositeLocations {
		result[idx] = readCompositeGlyph(data, loc, compositeLocations, result, emptyGlyph)
	}

	g.tableBytes = nil
	return result
}

func readSimpleGlyph(data *TrueTypeDataBytesLocal, contourCount int16, bounds core.PdfRectangle) glyphs.GlyphDescription {
	if contourCount == 0 {
		return glyphs.NewGlyph(true, nil, nil, nil, bounds)
	}

	endPointsOfContours := data.ReadUnsignedShortArray(int(contourCount))
	instructionLength := data.ReadUnsignedShort()
	instructions := data.ReadByteArray(int(instructionLength))

	pointCount := 0
	if contourCount > 0 {
		pointCount = int(endPointsOfContours[contourCount-1]) + 1
	}

	flags := readFlags(data, pointCount)
	xCoordinates := readCoordinates(data, pointCount, flags, glyphs.XSingleByte, glyphs.ThisXIsTheSame)
	yCoordinates := readCoordinates(data, pointCount, flags, glyphs.YSingleByte, glyphs.ThisYIsTheSame)

	endPtIndex := len(endPointsOfContours) - 1
	endPtOfContourIndex := -1
	points := make([]glyphs.GlyphPoint, pointCount)

	for i := pointCount - 1; i >= 0; i-- {
		if endPtOfContourIndex == -1 {
			endPtOfContourIndex = int(endPointsOfContours[endPtIndex])
		}
		isEndPt := endPtOfContourIndex == i
		if isEndPt && endPtIndex > 0 {
			endPtIndex--
			endPtOfContourIndex = -1
		}

		isOnCurve := (flags[i]&glyphs.OnCurve) == glyphs.OnCurve
		points[i] = glyphs.NewGlyphPoint(xCoordinates[i], yCoordinates[i], isOnCurve, isEndPt)
	}

	return glyphs.NewGlyph(true, instructions, endPointsOfContours, points, bounds)
}

func readCompositeGlyph(
	data *TrueTypeDataBytesLocal,
	compositeLoc temporaryCompositeLocation,
	compositeLocations map[int]temporaryCompositeLocation,
	glyphsArr []glyphs.GlyphDescription,
	emptyGlyph glyphs.GlyphDescription,
) glyphs.GlyphDescription {
	hasFlag := func(value, target glyphs.CompositeGlyphFlags) bool {
		return (value & target) == target
	}

	data.Seek(compositeLoc.position, 0) // io.SeekStart
	components := make([]compositeComponent, 0)

	var flagsVal glyphs.CompositeGlyphFlags
	for {
		flagsVal = glyphs.CompositeGlyphFlags(data.ReadUnsignedShort())
		glyphIndex := int(data.ReadUnsignedShort())

		if glyphIndex >= len(glyphsArr) {
			continue
		}

		childGlyph := glyphsArr[glyphIndex]

		if childGlyph == nil {
			missingComposite, ok := compositeLocations[glyphIndex]
			if !ok {
				return emptyGlyph
			}

			position := data.Position()
			childGlyph = readCompositeGlyph(data, missingComposite, compositeLocations, glyphsArr, emptyGlyph)
			data.Seek(position, 0) // io.SeekStart
			glyphsArr[glyphIndex] = childGlyph
		}

		var arg1, arg2 int16
		if hasFlag(flagsVal, glyphs.Args1And2AreWords) {
			arg1 = data.ReadSignedShort()
			arg2 = data.ReadSignedShort()
		} else {
			b1, _ := data.ReadByte()
			b2, _ := data.ReadByte()
			arg1 = int16(b1)
			arg2 = int16(b2)
		}

		xscale := 1.0
		scale01 := 0.0
		scale10 := 0.0
		yscale := 1.0

		if hasFlag(flagsVal, glyphs.WeHaveAScale) {
			xscale = readTwoFourteenFormat(data)
			yscale = xscale
		} else if hasFlag(flagsVal, glyphs.WeHaveAnXAndYScale) {
			xscale = readTwoFourteenFormat(data)
			yscale = readTwoFourteenFormat(data)
		} else if hasFlag(flagsVal, glyphs.WeHaveATwoByTwo) {
			xscale = readTwoFourteenFormat(data)
			scale01 = readTwoFourteenFormat(data)
			scale10 = readTwoFourteenFormat(data)
			yscale = readTwoFourteenFormat(data)
		}

		if hasFlag(flagsVal, glyphs.ArgsAreXAndYValues) {
			transform := glyphs.NewCompositeTransformMatrix3By2(xscale, scale01, scale10, yscale, float64(arg1), float64(arg2))
			components = append(components, compositeComponent{index: glyphIndex, transformation: transform})
		}

		if !hasFlag(flagsVal, glyphs.MoreComponents) {
			break
		}
	}

	var builderGlyph glyphs.GlyphDescription
	for _, comp := range components {
		glyph := glyphsArr[comp.index]
		transformed := glyph.Transform(comp.transformation)

		if builderGlyph == nil {
			builderGlyph = transformed
		} else {
			builderGlyph = builderGlyph.Merge(transformed)
		}
	}

	if builderGlyph == nil {
		builderGlyph = emptyGlyph
	}

	return glyphs.NewGlyph(
		false,
		builderGlyph.Instructions(),
		builderGlyph.EndPointsOfContours(),
		builderGlyph.Points(),
		compositeLoc.bounds,
	)
}

func readFlags(data *TrueTypeDataBytesLocal, pointCount int) []glyphs.SimpleGlyphFlags {
	result := make([]glyphs.SimpleGlyphFlags, pointCount)

	for i := 0; i < pointCount; i++ {
		b, _ := data.ReadByte()
		result[i] = glyphs.SimpleGlyphFlags(b)

		if result[i]&glyphs.Repeat != 0 {
			nr, _ := data.ReadByte()
			numberOfRepeats := int(nr)
			for j := 0; j < numberOfRepeats; j++ {
				p := i + j + 1
				if p >= len(result) {
					break
				}
				result[p] = result[i]
			}
			i += numberOfRepeats
		}
	}

	return result
}

func readCoordinates(
	data *TrueTypeDataBytesLocal,
	pointCount int,
	flags []glyphs.SimpleGlyphFlags,
	isByte glyphs.SimpleGlyphFlags,
	signOrSame glyphs.SimpleGlyphFlags,
) []int16 {
	hasFlag := func(value, target glyphs.SimpleGlyphFlags) bool {
		return (value & target) == target
	}

	xs := make([]int16, pointCount)
	x := 0

	for i := 0; i < pointCount; i++ {
		flag := flags[i]
		var dx int

		if hasFlag(flag, isByte) {
			rb, _ := data.ReadByte()
			b := int(rb)
			if hasFlag(flag, signOrSame) {
				dx = b
			} else {
				dx = -b
			}
		} else {
			if hasFlag(flag, signOrSame) {
				dx = 0
			} else {
				dx = int(data.ReadSignedShort())
			}
		}

		x += dx
		xs[i] = int16(x)
	}

	return xs
}

func readTwoFourteenFormat(data *TrueTypeDataBytesLocal) float64 {
	return float64(data.ReadSignedShort()) / (1 << 14)
}

// --- Local TrueTypeDataBytes wrapper to avoid cross-package dependency ---

// TrueTypeDataBytesLocal is a minimal byte reader for glyph table data.
type TrueTypeDataBytesLocal struct {
	data   []byte
	offset int64
}

func NewTrueTypeDataBytesLocal(data []byte) *TrueTypeDataBytesLocal {
	return &TrueTypeDataBytesLocal{data: data, offset: 0}
}

func (d *TrueTypeDataBytesLocal) Position() int64 { return d.offset }

func (d *TrueTypeDataBytesLocal) Length() int64 { return int64(len(d.data)) }

func (d *TrueTypeDataBytesLocal) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		d.offset = offset
	case 1:
		d.offset += offset
	case 2:
		d.offset = int64(len(d.data)) + offset
	}
	if d.offset < 0 {
		d.offset = 0
	}
	if d.offset > int64(len(d.data)) {
		d.offset = int64(len(d.data))
	}
	return d.offset, nil
}

func (d *TrueTypeDataBytesLocal) ReadByte() (byte, error) {
	if d.offset >= int64(len(d.data)) {
		return 0, fmt.Errorf("read past end of data at offset %d", d.offset)
	}
	b := d.data[d.offset]
	d.offset++
	return b, nil
}

func (d *TrueTypeDataBytesLocal) ReadSignedShort() int16 {
	if d.offset+2 > int64(len(d.data)) {
		d.offset = int64(len(d.data))
		return 0
	}
	b0 := d.data[d.offset]
	b1 := d.data[d.offset+1]
	d.offset += 2
	return int16(uint16(b0)<<8 | uint16(b1))
}

func (d *TrueTypeDataBytesLocal) ReadUnsignedShort() uint16 {
	if d.offset+2 > int64(len(d.data)) {
		d.offset = int64(len(d.data))
		return 0
	}
	b0 := d.data[d.offset]
	b1 := d.data[d.offset+1]
	d.offset += 2
	return uint16(b0)<<8 | uint16(b1)
}

func (d *TrueTypeDataBytesLocal) ReadUnsignedShortArray(n int) []uint16 {
	result := make([]uint16, n)
	for i := 0; i < n; i++ {
		result[i] = d.ReadUnsignedShort()
	}
	return result
}

func (d *TrueTypeDataBytesLocal) ReadByteArray(n int) []byte {
	end := d.offset + int64(n)
	if end > int64(len(d.data)) {
		end = int64(len(d.data))
	}
	result := make([]byte, end-d.offset)
	copy(result, d.data[d.offset:end])
	d.offset = end
	return result
}
