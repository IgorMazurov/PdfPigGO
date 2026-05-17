package cff

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// CompactFontFormatTopLevelDictionaryReader reads a top-level CFF dictionary from raw bytes.
type CompactFontFormatTopLevelDictionaryReaderImpl struct{}

// NewCompactFontFormatTopLevelDictionaryReaderImpl creates a new top-level dictionary reader.
func NewCompactFontFormatTopLevelDictionaryReaderImpl() *CompactFontFormatTopLevelDictionaryReaderImpl {
	return &CompactFontFormatTopLevelDictionaryReaderImpl{}
}

// Read parses the top-level dictionary from data using the given string index.
func (r *CompactFontFormatTopLevelDictionaryReaderImpl) Read(
	data *CompactFontFormatData,
	stringIndex []string,
) *CompactFontFormatTopLevelDictionary {
	dictionary := NewCompactFontFormatTopLevelDictionary()

	readDictionary(dictionary, data, stringIndex, r.applyOperation)

	return dictionary
}

func (r *CompactFontFormatTopLevelDictionaryReaderImpl) applyOperation(
	builder any,
	operands []Operand,
	key OperandKey,
	stringIndex []string,
) {
	dictionary := builder.(*CompactFontFormatTopLevelDictionary)

	switch key.Byte0 {
	case 0:
		if s, err := GetString(operands, stringIndex); err == nil {
			dictionary.Version = s
		}
	case 1:
		if s, err := GetString(operands, stringIndex); err == nil {
			dictionary.Notice = s
		}
	case 2:
		if s, err := GetString(operands, stringIndex); err == nil {
			dictionary.FullName = s
		}
	case 3:
		if s, err := GetString(operands, stringIndex); err == nil {
			dictionary.FamilyName = s
		}
	case 4:
		if s, err := GetString(operands, stringIndex); err == nil {
			dictionary.Weight = s
		}
	case 5:
		dictionary.FontBoundingBox = GetBoundingBox(operands)
	case 12:
		if key.Byte1 == nil {
			panic("a single byte sequence beginning with 12 was found in the CFF top-level dictionary")
		}

		switch *key.Byte1 {
		case 0:
			if s, err := GetString(operands, stringIndex); err == nil {
				dictionary.Copyright = s
			}
		case 1:
			dictionary.IsFixedPitch = operands[0].Double == 1
		case 2:
			dictionary.ItalicAngle = operands[0].Double
		case 3:
			dictionary.UnderlinePosition = operands[0].Double
		case 4:
			dictionary.UnderlineThickness = operands[0].Double
		case 5:
			dictionary.PaintType = operands[0].Double
		case 6:
			dictionary.CharStringType = CompactFontFormatCharStringType(GetIntOrDefault(operands, 2))
		case 7:
			array := ToArray(operands)
			if len(array) == 4 || len(array) == 6 {
				matrix, err := core.FromArray(array)
				if err == nil {
					dictionary.FontMatrix = &matrix
				}
			} else {
				panic(fmt.Sprintf("expected four or six values for the font matrix, instead got: %d", len(array)))
			}
		case 8:
			dictionary.StrokeWidth = operands[0].Double
		case 20:
			dictionary.SyntheticBaseFontIndex = GetIntOrDefault(operands, -1)
		case 21:
			if s, err := GetString(operands, stringIndex); err == nil {
				dictionary.PostScript = s
			}
		case 22:
			if s, err := GetString(operands, stringIndex); err == nil {
				dictionary.BaseFontName = s
			}
		case 23:
			dictionary.BaseFontBlend = ReadDeltaToArray(operands)
		case 30:
			registry, _ := GetString(operands, stringIndex)
			operands = operands[1:]
			ordering, _ := GetString(operands, stringIndex)
			operands = operands[1:]
			supplement := GetIntOrDefault(operands, -1)
			dictionary.CidFontOperators.Ros = RegistryOrderingSupplement{
				Registry:   registry,
				Ordering:   ordering,
				Supplement: float64(supplement),
			}
			dictionary.IsCidFont = true
		case 31:
			dictionary.CidFontOperators.Version = GetIntOrDefault(operands, -1)
		case 32:
			dictionary.CidFontOperators.Revision = GetIntOrDefault(operands, -1)
		case 33:
			dictionary.CidFontOperators.Type = GetIntOrDefault(operands, -1)
		case 34:
			dictionary.CidFontOperators.Count = GetIntOrDefault(operands, -1)
		case 35:
			dictionary.CidFontOperators.UidBase = float64(GetIntOrDefault(operands, -1))
		case 36:
			dictionary.CidFontOperators.FontDictionaryArray = GetIntOrDefault(operands, -1)
		case 37:
			dictionary.CidFontOperators.FontDictionarySelect = GetIntOrDefault(operands, -1)
		case 38:
			if s, err := GetString(operands, stringIndex); err == nil {
				dictionary.CidFontOperators.FontName = s
			}
		}
	case 13:
		if len(operands) > 0 {
			dictionary.UniqueId = operands[0].Double
		} else {
			dictionary.UniqueId = 0
		}
	case 14:
		dictionary.Xuid = ToArray(operands)
	case 15:
		dictionary.CharSetOffset = GetIntOrDefault(operands, UnsetOffset)
	case 16:
		dictionary.EncodingOffset = GetIntOrDefault(operands, UnsetOffset)
	case 17:
		dictionary.CharStringsOffset = GetIntOrDefault(operands, -1)
	case 18:
		size := GetIntOrDefault(operands, -1)
		operands = operands[1:]
		offset := GetIntOrDefault(operands, -1)
		dictionary.PrivateDictionaryLocation = &SizeAndOffset{
			Size:   size,
			Offset: offset,
		}
	}
}
