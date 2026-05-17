package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// BaseFontRangeParser parses beginbfrange/endbfrange blocks in CMap data.
type BaseFontRangeParser struct{}

// NewBaseFontRangeParser creates a new BaseFontRangeParser.
func NewBaseFontRangeParser() *BaseFontRangeParser {
	return &BaseFontRangeParser{}
}

// Parse reads base font range mappings from the token scanner and adds them to the builder.
func (p *BaseFontRangeParser) Parse(numeric *tokens.NumericToken, scanner tokenization.TokenScanner, builder *cmap.CharacterMapBuilder) error {
	count := numeric.IntVal()

	for i := 0; i < count; i++ {
		if !scanner.Advance() {
			return fmt.Errorf("bfrange was missing the low source code: %v", scanner.Current())
		}

		lowSourceCode, ok := scanner.Current().(*tokens.HexToken)
		if !ok {
			if op, isOp := scanner.Current().(*tokens.OperatorToken); isOp {
				opData := op.Data()
				if opData == "endbfrange" || opData == "Endbfrange" || opData == "ENDBFRANGE" {
					break
				}
			}
			return fmt.Errorf("bfrange was missing the low source code: %v", scanner.Current())
		}

		if !scanner.Advance() {
			return fmt.Errorf("bfrange was missing the high source code: %v", scanner.Current())
		}

		highSourceCode, ok := scanner.Current().(*tokens.HexToken)
		if !ok {
			return fmt.Errorf("bfrange was missing the high source code: %v", scanner.Current())
		}

		if !scanner.Advance() {
			return fmt.Errorf("bfrange ended unexpectedly after the high source code")
		}

		var destinationBytes []byte
		var destinationArray *tokens.ArrayToken

		switch current := scanner.Current().(type) {
		case *tokens.ArrayToken:
			destinationArray = current
		case *tokens.HexToken:
			destinationBytes = make([]byte, len(current.Bytes()))
			copy(destinationBytes, current.Bytes())
		case *tokens.NumericToken:
			return fmt.Errorf("numeric destination in bfrange is not yet supported")
		default:
			return fmt.Errorf("unexpected token type in bfrange destination: %T", scanner.Current())
		}

		startCode := make([]byte, len(lowSourceCode.Bytes()))
		copy(startCode, lowSourceCode.Bytes())
		endCode := highSourceCode.Bytes()

		done := false

		if destinationArray != nil {
			arrayIndex := 0
			for !done {
				if byteCompare(startCode, endCode) >= 0 {
					done = true
				}

				dest := destinationArray.Data()[arrayIndex]

				if name, ok := dest.(*tokens.NameToken); ok {
					builder.AddBaseFontCharacter(startCode, name.Data())
				} else if hexDest, ok := dest.(*tokens.HexToken); ok {
					builder.AddBaseFontCharacterBytes(startCode, hexDest.Bytes())
				}

				incrementByteSlice(startCode)
				arrayIndex++
			}
			continue
		}

		for !done {
			if byteCompare(startCode, endCode) >= 0 {
				done = true
			}

			builder.AddBaseFontCharacterBytes(startCode, destinationBytes)

			incrementByteSlice(startCode)
			incrementByteSlice(destinationBytes)
		}
	}

	return nil
}

func incrementByteSlice(data []byte) {
	pos := len(data) - 1
	for pos > 0 && data[pos] == 255 {
		data[pos] = 0
		pos--
	}
	data[pos]++
}

func byteCompare(first, second []byte) int {
	minLen := len(first)
	if len(second) < minLen {
		minLen = len(second)
	}
	for i := 0; i < minLen; i++ {
		if first[i] == second[i] {
			continue
		}
		if first[i] < second[i] {
			return -1
		}
		return 1
	}
	if len(first) < len(second) {
		return -1
	}
	if len(first) > len(second) {
		return 1
	}
	return 0
}


