package cff

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	cffcharset "github.com/uglytoad/pdfpig/go/fonts/cff_charset"
)

const (
	tagOtto    = "OTTO"
	tagTtcf    = "ttcf"
	tagTtfOnly = "\x00\x01\x00\x00"
)

// CharStringsParseCallback is a function type for parsing Type 2 charstring bytes into
// a provider that satisfies the Type2CharStringsProvider interface. The caller passes
// an adapted version of charstrings.Parse to bridge the import-cycle gap between the
// cff and charstrings packages. Returns nil and a non-nil error on failure.
type CharStringsParseCallback func(
	charStringBytes [][]byte,
	subroutinesSelector *CompactFontFormatSubroutinesSelector,
	charset cffcharset.CompactFontFormatCharset,
) (Type2CharStringsProvider, error)

// Parse reads a Compact Font Format font collection from the given raw data.
// The parseCharStrings callback is used to decode Type 2 charstring programs;
// pass a wrapper around charstrings.Parse to resolve the cff↔charstrings import cycle.
// Returns nil and an error for unsupported tag formats or parsing failures.
func Parse(
	data *CompactFontFormatData,
	parseCharStrings CharStringsParseCallback,
) (*CompactFontFormatFontCollection, error) {
	tag := readTag(data)

	switch tag {
	case tagOtto:
		return nil, fmt.Errorf("tagged CFF data is not supported")
	case tagTtcf:
		return nil, fmt.Errorf("True Type Collection fonts are not currently supported")
	case tagTtfOnly:
		return nil, fmt.Errorf("OpenType fonts containing a TrueType font are not currently supported")
	default:
		data.Seek(0)
	}

	header := readHeader(data)

	fontNames := readStringIndex(data)

	topLevelDictionaryIndex, err := ReadDictionaryData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read top-level dictionary index: %w", err)
	}

	stringIndex := readStringIndex(data)

	globalSubroutineIndex, err := ReadDictionaryData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read global subroutine index: %w", err)
	}

	fonts := make(map[string]*CompactFontFormatFont)

	topLevelReader := NewCompactFontFormatTopLevelDictionaryReader()
	privateReader := NewCompactFontFormatPrivateDictionaryReader()
	individualParser := NewCompactFontFormatIndividualFontParser(topLevelReader, privateReader)

	for i := 0; i < len(fontNames); i++ {
		fontName := fontNames[i]
		topDictEntry := topLevelDictionaryIndex.Get(i)

		result, parseErr := individualParser.Parse(data, fontName, topDictEntry, stringIndex, globalSubroutineIndex)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse font %q: %w", fontName, parseErr)
		}

		subroutinesSelector := result.GetSubroutinesSelector()
		var type2CharStrings Type2CharStringsProvider
		if parseCharStrings != nil {
			var csErr error
			type2CharStrings, csErr = parseCharStrings(
				result.CharStringIndex.GetBytes(),
				&subroutinesSelector,
				result.Charset,
			)
			if csErr != nil {
				return nil, fmt.Errorf("failed to parse charstrings for font %q: %w", fontName, csErr)
			}
		}

		fonts[fontName] = result.BuildFont(type2CharStrings)
	}

	return NewCompactFontFormatFontCollection(header, fonts)
}

// readTag reads the first 4 bytes as a tag string using ISO-8859-1 encoding.
func readTag(data *CompactFontFormatData) string {
	return data.ReadString(4, core.Iso88591)
}

// readHeader reads the CFF header fields: major version, minor version,
// header size, and offset size.
func readHeader(data *CompactFontFormatData) CompactFontFormatHeader {
	major := data.ReadCard8()
	minor := data.ReadCard8()
	headerSize := data.ReadCard8()
	offsetSize := data.ReadOffsize()

	return NewCompactFontFormatHeader(major, minor, headerSize, offsetSize)
}

// readStringIndex reads an index structure from the data stream and extracts
// each indexed entry as a string using ISO-8859-1 encoding.
func readStringIndex(data *CompactFontFormatData) []string {
	index := ReadIndex(data)

	if len(index) == 0 {
		return nil
	}

	count := len(index) - 1
	result := make([]string, count)

	for i := 0; i < count; i++ {
		length := index[i+1] - index[i]

		if length < 0 {
			panic(fmt.Sprintf("negative object length %d at %d. Current position: %d.", length, i, data.Position()))
		}

		result[i] = data.ReadString(length, core.Iso88591)
	}

	return result
}

// CompactFontFormatTopLevelDictionaryReader reads a top-level CFF dictionary from raw data.
type CompactFontFormatTopLevelDictionaryReader = CompactFontFormatTopLevelDictionaryReaderImpl

// NewCompactFontFormatTopLevelDictionaryReader creates a new top-level dictionary reader.
var NewCompactFontFormatTopLevelDictionaryReader = NewCompactFontFormatTopLevelDictionaryReaderImpl
