package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// BaseFontCharacterParser parses beginbfchar/endbfchar blocks in CMap data.
type BaseFontCharacterParser struct{}

// NewBaseFontCharacterParser creates a new BaseFontCharacterParser.
func NewBaseFontCharacterParser() *BaseFontCharacterParser {
	return &BaseFontCharacterParser{}
}

// Parse reads base font character mappings from the token scanner and adds them to the builder.
func (p *BaseFontCharacterParser) Parse(numeric *tokens.NumericToken, scanner tokenization.TokenScanner, builder *cmap.CharacterMapBuilder) error {
	count := numeric.IntVal()

	for i := 0; i < count; i++ {
		if !scanner.Advance() {
			return fmt.Errorf("base font characters definition contains invalid item at index %d: scanner exhausted", i)
		}

		current := scanner.Current()

		inputCode, ok := current.(*tokens.HexToken)
		if !ok {
			if op, isOp := current.(*tokens.OperatorToken); isOp {
				opData := op.Data()
				if opData == "endbfchar" || opData == "Endbfchar" || opData == "ENDBFCHAR" ||
					opData == "endcmap" || opData == "Endcmap" || opData == "ENDCMAP" {
					return nil
				}
			}
			return fmt.Errorf("base font characters definition contains invalid item at index %d: %v", i, current)
		}

		if !scanner.Advance() {
			return fmt.Errorf("base font characters definition contains invalid item at index %d: scanner exhausted", i)
		}

		current = scanner.Current()

		if characterName, ok := current.(*tokens.NameToken); ok {
			builder.AddBaseFontCharacter(inputCode.Bytes(), characterName.Data())
		} else if characterCode, ok := current.(*tokens.HexToken); ok {
			builder.AddBaseFontCharacterBytes(inputCode.Bytes(), characterCode.Bytes())
		} else {
			return fmt.Errorf("base font characters definition contains invalid item at index %d: %v", i, current)
		}
	}

	return nil
}

