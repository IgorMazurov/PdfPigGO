package cff

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/cfftypes"
	cffcharset "github.com/uglytoad/pdfpig/go/fonts/cff_charset"
)

// CompactFontFormatSubroutinesSelector is an alias for the shared type.
type CompactFontFormatSubroutinesSelector = cfftypes.CompactFontFormatSubroutinesSelector



// CompactFontFormatIndividualFontParser parses a single CFF font from the raw data,
// reading top-level and private dictionaries, charset, encoding, and charstrings.
type CompactFontFormatIndividualFontParser struct {
	topLevelDictionaryReader TopLevelDictionaryReader
	privateDictionaryReader  PrivateDictionaryReader
}

// NewCompactFontFormatIndividualFontParser creates a new parser with the given
// dictionary readers.
func NewCompactFontFormatIndividualFontParser(
	topLevelDictionaryReader TopLevelDictionaryReader,
	privateDictionaryReader PrivateDictionaryReader,
) *CompactFontFormatIndividualFontParser {
	return &CompactFontFormatIndividualFontParser{
		topLevelDictionaryReader: topLevelDictionaryReader,
		privateDictionaryReader:  privateDictionaryReader,
	}
}

// ParseResult holds the intermediate data produced by parsing a CFF font before
// Type2CharStrings are generated. The caller uses this to construct the final
// CompactFontFormatFont or CompactFontFormatCidFont via BuildFont/BuildCidFont.
type ParseResult struct {
	TopDictionary        *CompactFontFormatTopLevelDictionary
	PrivateDictionary    CompactFontFormatPrivateDictionary
	Charset              cffcharset.CompactFontFormatCharset
	CharStringIndex      *CompactFontFormatIndex
	GlobalSubroutines    *CompactFontFormatIndex
	LocalSubroutines     *CompactFontFormatIndex
	Encoding             encodingSource
	IsCidFont            bool
	CidFontDictionaries  []*CompactFontFormatTopLevelDictionary
	CidPrivateDicts      []CompactFontFormatPrivateDictionary
	CidLocalSubroutines  []*CompactFontFormatIndex
	FdSelect             FdSelect
	NumberOfGlyphs       int
}

// Parse reads a single CFF font from the data using the given indices and returns
// intermediate parsing results. The caller then constructs the final font by calling
// BuildFont or BuildCidFont with parsed Type2CharStrings. This two-step approach
// avoids an import cycle between cff and charstrings packages.
func (p *CompactFontFormatIndividualFontParser) Parse(
	data *CompactFontFormatData,
	name string,
	topDictionaryIndex []byte,
	stringIndex []string,
	globalSubroutineIndex *CompactFontFormatIndex,
) (*ParseResult, error) {
	individualData := NewCompactFontFormatData(topDictionaryIndex)

	topDictionary := p.topLevelDictionaryReader.Read(individualData, stringIndex)

	privateDictionary := DefaultCompactFontFormatPrivateDictionary

	if topDictionary.PrivateDictionaryLocation != nil && topDictionary.PrivateDictionaryLocation.Size > 0 {
		privateDictBytes, err := data.SnapshotPortion(
			topDictionary.PrivateDictionaryLocation.Offset,
			topDictionary.PrivateDictionaryLocation.Size,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to snapshot private dictionary: %w", err)
		}

		privateDictionary = p.privateDictionaryReader.Read(privateDictBytes, stringIndex)
	}

	if topDictionary.CharStringsOffset < 0 {
		return nil, fmt.Errorf("expected CFF to contain a CharString offset")
	}

	localSubroutines := (*CompactFontFormatIndex)(nil)
	if privateDictionary.LocalSubroutineOffset != nil && topDictionary.PrivateDictionaryLocation != nil {
		data.Seek(*privateDictionary.LocalSubroutineOffset + topDictionary.PrivateDictionaryLocation.Offset)
		var err error
		localSubroutines, err = ReadDictionaryData(data)
		if err != nil {
			return nil, fmt.Errorf("failed to read local subroutines index: %w", err)
		}
	}

	data.Seek(topDictionary.CharStringsOffset)

	charStringIndex, err := ReadDictionaryData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read charstring index: %w", err)
	}

	charset := p.readCharset(data, topDictionary, charStringIndex, stringIndex)

	if topDictionary.IsCidFont {
		return p.readCidFontData(
			data, topDictionary, charStringIndex.Count(), stringIndex,
			privateDictionary, charset, globalSubroutineIndex, localSubroutines,
			charStringIndex,
		)
	}

	fontEncoding := (encodingSource)(nil)
	encodingOffset := topDictionary.EncodingOffset
	if encodingOffset != UnsetOffset {
		if encodingOffset == 0 {
			fontEncoding = CFFStandardEncoding
		} else if encodingOffset == 1 {
			fontEncoding = CFFExpertEncoding
		} else {
			data.Seek(encodingOffset)
			cffEncoding, readErr := ReadEncoding(data, charset, stringIndex)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read encoding: %w", readErr)
			}
			fontEncoding = cffEncoding
		}
	}

	return &ParseResult{
		TopDictionary:     topDictionary,
		PrivateDictionary: privateDictionary,
		Charset:           charset,
		CharStringIndex:   charStringIndex,
		GlobalSubroutines: globalSubroutineIndex,
		LocalSubroutines:  localSubroutines,
		Encoding:          fontEncoding,
	}, nil
}

// BuildFont constructs a CompactFontFormatFont from the parse result and already-parsed
// Type2CharStrings. This is called after Parse returns for non-CID fonts.
func (p *ParseResult) BuildFont(type2CharStrings Type2CharStringsProvider) *CompactFontFormatFont {
	return NewCompactFontFormatFont(
		p.TopDictionary, p.PrivateDictionary, p.Charset, nil, type2CharStrings, p.Encoding,
	)
}

// BuildCidFont constructs a CompactFontFormatCidFont from the parse result and already-parsed
// Type2CharStrings. This is called after Parse returns for CID fonts.
func (p *ParseResult) BuildCidFont(type2CharStrings Type2CharStringsProvider) *CompactFontFormatCidFont {
	return NewCompactFontFormatCidFont(
		p.TopDictionary, p.PrivateDictionary, p.Charset, nil, type2CharStrings,
		p.CidFontDictionaries, p.CidPrivateDicts, p.FdSelect,
	)
}

// GetSubroutinesSelector creates a CompactFontFormatSubroutinesSelector from the parse result.
// For CID fonts it uses FDSelect and per-dictionary subroutines; for non-CID fonts it uses
// the top-level local subroutines. The returned selector satisfies the charstrings.FdSelect
// interface when IsCidFont is true.
func (p *ParseResult) GetSubroutinesSelector() CompactFontFormatSubroutinesSelector {
	if p.IsCidFont {
		return CompactFontFormatSubroutinesSelector{
			GlobalSubroutines:    p.GlobalSubroutines,
			LocalSubroutines:     p.LocalSubroutines,
			FdSelect:             p.FdSelect,
			FontLocalSubroutines: p.CidLocalSubroutines,
		}
	}
	return CompactFontFormatSubroutinesSelector{
		GlobalSubroutines: p.GlobalSubroutines,
		LocalSubroutines:  p.LocalSubroutines,
	}
}

// readCharset determines the charset for a CFF font based on the top-level dictionary's
// charset offset. Standard predefined charsets (0, 1, 2) are used when applicable;
// otherwise the charset is read from the data stream.
func (p *CompactFontFormatIndividualFontParser) readCharset(
	data *CompactFontFormatData,
	topDictionary *CompactFontFormatTopLevelDictionary,
	charStringIndex *CompactFontFormatIndex,
	stringIndex []string,
) cffcharset.CompactFontFormatCharset {
	if topDictionary.CharSetOffset >= 0 {
		charsetId := topDictionary.CharSetOffset
		if !topDictionary.IsCidFont && charsetId == 0 {
			return cffcharset.IsoAdobeValue
		} else if !topDictionary.IsCidFont && charsetId == 1 {
			return cffcharset.ExpertValue
		} else if !topDictionary.IsCidFont && charsetId == 2 {
			return cffcharset.ExpertSubsetValue
		} else {
			return p.readCharsetData(data, topDictionary, charStringIndex, stringIndex)
		}
	}

	if topDictionary.IsCidFont {
		return cffcharset.NewCompactFontFormatEmptyCharset(charStringIndex.Count())
	}

	return cffcharset.IsoAdobeValue
}

// readCharsetData reads a charset from the data stream at the position indicated by
// the top-level dictionary's CharSetOffset. Supports format 0 (individual SIDs),
// format 1 (ranges with byte count), and format 2 (ranges with 16-bit count).
func (p *CompactFontFormatIndividualFontParser) readCharsetData(
	data *CompactFontFormatData,
	topDictionary *CompactFontFormatTopLevelDictionary,
	charStringIndex *CompactFontFormatIndex,
	stringIndex []string,
) cffcharset.CompactFontFormatCharset {
	data.Seek(topDictionary.CharSetOffset)

	format := data.ReadCard8()

	switch format {
	case 0:
		return p.readFormat0Charset(data, charStringIndex.Count(), stringIndex)
	case 1, 2:
		return p.readFormatRangeCharset(data, format == 1, charStringIndex.Count(), stringIndex)
	default:
		panic(fmt.Sprintf("unrecognized format for the Charset table in a CFF font. Got: %d.", format))
	}
}

// readFormat0Charset reads charset format 0 where each glyph has an individual SID.
func (p *CompactFontFormatIndividualFontParser) readFormat0Charset(
	data *CompactFontFormatData, charStringCount int, stringIndex []string,
) cffcharset.CompactFontFormatCharset {
	entries := make([]struct {
		GlyphId  int
		StringId int
		Name     string
	}, 0, charStringCount-1)

	for glyphId := 1; glyphId < charStringCount; glyphId++ {
		stringId := data.ReadSid()
		entries = append(entries, struct {
			GlyphId  int
			StringId int
			Name     string
		}{glyphId, stringId, readCharsetString(stringId, stringIndex)})
	}

	return cffcharset.NewCompactFontFormatFormat0Charset(entries)
}

// readFormatRangeCharset reads charset formats 1 or 2 where glyphs are grouped in ranges.
// isFormat1 determines whether range counts are single bytes (format 1) or two bytes (format 2).
func (p *CompactFontFormatIndividualFontParser) readFormatRangeCharset(
	data *CompactFontFormatData, isFormat1 bool, charStringCount int, stringIndex []string,
) cffcharset.CompactFontFormatCharset {
	entries := make([]struct {
		GlyphId  int
		StringId int
		Name     string
	}, 0, charStringCount-1)

	glyphId := 1
	for glyphId < charStringCount {
		firstSid := data.ReadSid()
		var numberInRange int
		if isFormat1 {
			numberInRange = int(data.ReadCard8())
		} else {
			numberInRange = int(data.ReadCard16())
		}

		entries = append(entries, struct {
			GlyphId  int
			StringId int
			Name     string
		}{glyphId, firstSid, readCharsetString(firstSid, stringIndex)})

		for i := 0; i < numberInRange; i++ {
			glyphId++
			sid := firstSid + i + 1
			entries = append(entries, struct {
				GlyphId  int
				StringId int
				Name     string
			}{glyphId, sid, readCharsetString(sid, stringIndex)})
		}
		glyphId++
	}

	if isFormat1 {
		return cffcharset.NewCompactFontFormatFormat1Charset(entries)
	}
	return cffcharset.NewCompactFontFormatFormat2Charset(entries)
}

// readString resolves a SID to its string name. SIDs 0-390 come from the CFF standard
// strings table; higher SIDs are looked up in the font's string index; out-of-range
// values produce a synthetic "SID{n}" name.
func readCharsetString(index int, stringIndex []string) string {
	if index >= 0 && index <= 390 {
		return GetName(index)
	}
	if index-391 < len(stringIndex) {
		return stringIndex[index-391]
	}
	return fmt.Sprintf("SID%d", index)
}

// readCidFontData parses the CID font data structures (font dictionaries, private dicts,
// FDSelect) and returns them in a ParseResult. The caller then builds the final font
// with parsed Type2CharStrings via BuildCidFont.
func (p *CompactFontFormatIndividualFontParser) readCidFontData(
	data *CompactFontFormatData,
	topLevelDictionary *CompactFontFormatTopLevelDictionary,
	numberOfGlyphs int,
	stringIndex []string,
	privateDictionary CompactFontFormatPrivateDictionary,
	charset cffcharset.CompactFontFormatCharset,
	globalSubroutines *CompactFontFormatIndex,
	localSubroutinesTop *CompactFontFormatIndex,
	charStringIndex *CompactFontFormatIndex,
) (*ParseResult, error) {
	offset := topLevelDictionary.CidFontOperators.FontDictionaryArray

	data.Seek(offset)

	fontDict, err := ReadDictionaryData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read font dictionary index: %w", err)
	}

	privateDictionaries := make([]CompactFontFormatPrivateDictionary, 0, fontDict.Count())
	fontDictionaries := make([]*CompactFontFormatTopLevelDictionary, 0, fontDict.Count())
	fontLocalSubroutines := make([]*CompactFontFormatIndex, 0, fontDict.Count())

	for i := 0; i < fontDict.Count(); i++ {
		indexBytes := fontDict.Get(i)
		topLevelDictionaryCid := p.topLevelDictionaryReader.Read(
			NewCompactFontFormatData(indexBytes), stringIndex,
		)

		if topLevelDictionaryCid.PrivateDictionaryLocation == nil {
			return nil, fonts.NewInvalidFontFormatException(
				"The CID keyed Compact Font Format font did not contain a private dictionary for the font dictionary.",
			)
		}

		privateDictionaryBytes, snapErr := data.SnapshotPortion(
			topLevelDictionaryCid.PrivateDictionaryLocation.Offset,
			topLevelDictionaryCid.PrivateDictionaryLocation.Size,
		)
		if snapErr != nil {
			return nil, fmt.Errorf("failed to snapshot CID private dictionary: %w", snapErr)
		}

		privateDictionaryCid := p.privateDictionaryReader.Read(privateDictionaryBytes, stringIndex)

		if privateDictionaryCid.LocalSubroutineOffset != nil && *privateDictionaryCid.LocalSubroutineOffset > 0 {
			data.Seek(topLevelDictionaryCid.PrivateDictionaryLocation.Offset + *privateDictionaryCid.LocalSubroutineOffset)
			localSubroutines, readErr := ReadDictionaryData(data)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read CID local subroutines: %w", readErr)
			}
			fontLocalSubroutines = append(fontLocalSubroutines, localSubroutines)
		} else {
			fontLocalSubroutines = append(fontLocalSubroutines, nil)
		}

		fontDictionaries = append(fontDictionaries, topLevelDictionaryCid)
		privateDictionaries = append(privateDictionaries, privateDictionaryCid)
	}

	data.Seek(topLevelDictionary.CidFontOperators.FontDictionarySelect)

	format := data.ReadCard8()

	var fdSelect FdSelect
	switch format {
	case 0:
		fdSelect = readFormat0FdSelect(data, numberOfGlyphs, topLevelDictionary.CidFontOperators.Ros)
	case 3:
		fdSelect = readFormat3FdSelect(data, topLevelDictionary.CidFontOperators.Ros)
	default:
		return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("invalid Font Dictionary Select format: %d.", format))
	}

	return &ParseResult{
		TopDictionary:       topLevelDictionary,
		PrivateDictionary:   privateDictionary,
		Charset:             charset,
		CharStringIndex:     charStringIndex,
		GlobalSubroutines:   globalSubroutines,
		LocalSubroutines:    localSubroutinesTop,
		IsCidFont:           true,
		CidFontDictionaries: fontDictionaries,
		CidPrivateDicts:     privateDictionaries,
		CidLocalSubroutines: fontLocalSubroutines,
		FdSelect:            fdSelect,
		NumberOfGlyphs:      numberOfGlyphs,
	}, nil
}

// readFormat0FdSelect reads a format 0 FDSelect where each glyph has an explicit
// font dictionary index byte.
func readFormat0FdSelect(data *CompactFontFormatData, numberOfGlyphs int, ros RegistryOrderingSupplement) FdSelect {
	dictionaries := make([]int, numberOfGlyphs)

	for i := 0; i < numberOfGlyphs; i++ {
		dictionaries[i] = int(data.ReadCard8())
	}

	return NewCff0FdSelect(ros, dictionaries)
}

// readFormat3FdSelect reads a format 3 FDSelect where font dictionary assignments
// are encoded as ranges. Each range specifies the first glyph ID and the dictionary index.
func readFormat3FdSelect(data *CompactFontFormatData, ros RegistryOrderingSupplement) FdSelect {
	numberOfRanges := int(data.ReadCard16())
	ranges := make([]Cff3FdRange, numberOfRanges)

	for i := 0; i < numberOfRanges; i++ {
		first := int(data.ReadCard16())
		dictionary := int(data.ReadCard8())
		ranges[i] = Cff3FdRange{First: first, FontDictionary: dictionary}
	}

	sentinel := int(data.ReadCard16())

	return NewCff3FdSelect(ros, ranges, sentinel)
}

// TopLevelDictionaryReader reads a top-level CFF dictionary from raw data.
type TopLevelDictionaryReader interface {
	Read(data *CompactFontFormatData, stringIndex []string) *CompactFontFormatTopLevelDictionary
}

// PrivateDictionaryReader reads a private CFF dictionary from raw bytes.
type PrivateDictionaryReader interface {
	Read(bytes *CompactFontFormatData, stringIndex []string) CompactFontFormatPrivateDictionary
}



/* ------------------------------------------------------------------ */
/*  FdSelect implementations                                           */
/* ------------------------------------------------------------------ */

// Cff0FdSelect maps each glyph ID to a font dictionary index using a flat array.
type Cff0FdSelect struct {
	registryOrderingSupplement RegistryOrderingSupplement
	fontDictionaries           []int
}

// NewCff0FdSelect creates a format 0 FDSelect from the given ROS and per-glyph dictionary indices.
func NewCff0FdSelect(ros RegistryOrderingSupplement, fontDictionaries []int) *Cff0FdSelect {
	return &Cff0FdSelect{
		registryOrderingSupplement: ros,
		fontDictionaries:           fontDictionaries,
	}
}

// GetFontDictionaryIndex returns the font dictionary index for the given glyph ID.
func (f *Cff0FdSelect) GetFontDictionaryIndex(glyphId int) int {
	if glyphId >= 0 && glyphId < len(f.fontDictionaries) {
		return f.fontDictionaries[glyphId]
	}
	return 0
}

// Cff3FdRange represents a single range entry in format 3 FDSelect.
type Cff3FdRange struct {
	First          int
	FontDictionary int
}

// Cff3FdSelect maps glyph IDs to font dictionary indices using compressed ranges.
type Cff3FdSelect struct {
	registryOrderingSupplement RegistryOrderingSupplement
	ranges                     []Cff3FdRange
	sentinel                   int
}

// NewCff3FdSelect creates a format 3 FDSelect from the given ROS, range entries, and sentinel value.
func NewCff3FdSelect(ros RegistryOrderingSupplement, ranges []Cff3FdRange, sentinel int) *Cff3FdSelect {
	return &Cff3FdSelect{
		registryOrderingSupplement: ros,
		ranges:                     ranges,
		sentinel:                   sentinel,
	}
}

// GetFontDictionaryIndex returns the font dictionary index for the given glyph ID
// by searching through the range entries. Returns -1 if no matching range is found
// and the glyph ID exceeds the sentinel value.
func (f *Cff3FdSelect) GetFontDictionaryIndex(glyphId int) int {
	for i := 0; i < len(f.ranges); i++ {
		if f.ranges[i].First <= glyphId {
			if i+1 < len(f.ranges) {
				if f.ranges[i+1].First > glyphId {
					return f.ranges[i].FontDictionary
				}
			} else {
				if f.sentinel > glyphId {
					return f.ranges[i].FontDictionary
				}
				return -1
			}
		}
	}
	return 0
}

/* ------------------------------------------------------------------ */
/*  CompactFontFormatSubroutinesSelector                               */
/* ------------------------------------------------------------------ */

// NewCompactFontFormatSubroutinesSelector creates a selector for non-CID fonts.
func NewCompactFontFormatSubroutinesSelector(global, local *CompactFontFormatIndex) CompactFontFormatSubroutinesSelector {
	return CompactFontFormatSubroutinesSelector{
		GlobalSubroutines: global,
		LocalSubroutines:  local,
	}
}

// NewCidCompactFontFormatSubroutinesSelector creates a selector for CID-keyed fonts.
func NewCidCompactFontFormatSubroutinesSelector(global, local *CompactFontFormatIndex, fdSelect FdSelect, perFontLocalSubroutines []*CompactFontFormatIndex) CompactFontFormatSubroutinesSelector {
	return CompactFontFormatSubroutinesSelector{
		GlobalSubroutines:    global,
		LocalSubroutines:     local,
		FdSelect:             fdSelect,
		FontLocalSubroutines: perFontLocalSubroutines,
	}
}
