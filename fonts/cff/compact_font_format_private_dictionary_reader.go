package cff

// CompactFontFormatPrivateDictionaryReader reads a private CFF dictionary from raw bytes.
type CompactFontFormatPrivateDictionaryReader struct{}

// NewCompactFontFormatPrivateDictionaryReader creates a new private dictionary reader.
func NewCompactFontFormatPrivateDictionaryReader() *CompactFontFormatPrivateDictionaryReader {
	return &CompactFontFormatPrivateDictionaryReader{}
}

// Read parses the private dictionary from data using the given string index.
func (r *CompactFontFormatPrivateDictionaryReader) Read(
	data *CompactFontFormatData,
	stringIndex []string,
) CompactFontFormatPrivateDictionary {
	builder := &CompactFontFormatPrivateDictionaryBuilder{}

	readDictionary(builder, data, stringIndex, r.applyOperation)

	return NewCompactFontFormatPrivateDictionary(builder)
}

func (r *CompactFontFormatPrivateDictionaryReader) applyOperation(
	builder any,
	operands []Operand,
	key OperandKey,
	stringIndex []string,
) {
	dictionary := builder.(*CompactFontFormatPrivateDictionaryBuilder)

	switch key.Byte0 {
	case 6:
		dictionary.BlueValues = ReadDeltaToIntArray(operands)
	case 7:
		dictionary.OtherBlues = ReadDeltaToIntArray(operands)
	case 8:
		dictionary.FamilyBlues = ReadDeltaToIntArray(operands)
	case 9:
		dictionary.FamilyOtherBlues = ReadDeltaToIntArray(operands)
	case 10:
		v := operands[0].Double
		dictionary.StandardHorizontalWidth = &v
	case 11:
		v := operands[0].Double
		dictionary.StandardVerticalWidth = &v
	case 12:
		if key.Byte1 == nil {
			panic("in the CFF private dictionary, got the operation key 12 without a second byte")
		}

		switch *key.Byte1 {
		case 9:
			v := operands[0].Double
			dictionary.BlueScale = &v
		case 10:
			v := operands[0].Int
			if v != nil {
				dictionary.BlueShift = v
			}
		case 11:
			v := operands[0].Int
			if v != nil {
				dictionary.BlueFuzz = v
			}
		case 12:
			dictionary.StemSnapHorizontalWidths = ReadDeltaToArray(operands)
		case 13:
			dictionary.StemSnapVerticalWidths = ReadDeltaToArray(operands)
		case 14:
			v := operands[0].Double == 1
			dictionary.ForceBold = &v
		case 17:
			v := operands[0].Int
			if v != nil {
				dictionary.LanguageGroup = v
			}
		case 18:
			v := operands[0].Double
			dictionary.ExpansionFactor = &v
		case 19:
			dictionary.InitialRandomSeed = operands[0].Double
		}
	case 19:
		v := GetIntOrDefault(operands, -1)
		dictionary.LocalSubroutineOffset = &v
	case 20:
		dictionary.DefaultWidthX = operands[0].Double
	case 21:
		dictionary.NominalWidthX = operands[0].Double
	}
}
