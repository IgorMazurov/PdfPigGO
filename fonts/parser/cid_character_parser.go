package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CidCharacterParser parses begincidchar/endcidchar blocks in CMap data.
type CidCharacterParser struct{}

// NewCidCharacterParser creates a new CidCharacterParser.
func NewCidCharacterParser() *CidCharacterParser {
	return &CidCharacterParser{}
}

// Parse reads CID character mappings from the token scanner and adds them to the builder.
func (p *CidCharacterParser) Parse(numeric *tokens.NumericToken, scanner tokenization.TokenScanner, builder *cmap.CharacterMapBuilder) error {
	count := numeric.IntVal()
	results := make([]cmap.CidCharacterMapping, 0, count)

	for i := 0; i < count; i++ {
		if !scanner.Advance() {
			return fmt.Errorf("cidchar was missing the source code: %v", scanner.Current())
		}

		sourceCode, ok := scanner.Current().(*tokens.HexToken)
		if !ok {
			return fmt.Errorf("the first token in a line for Cid Characters should be a hex, instead it was: %v", scanner.Current())
		}

		if !scanner.Advance() {
			return fmt.Errorf("cidchar was missing the destination code: %v", scanner.Current())
		}

		destinationCode, ok := scanner.Current().(*tokens.NumericToken)
		if !ok {
			return fmt.Errorf("the destination token in a line for Cid Character should be an integer, instead it was: %v", scanner.Current())
		}

		sourceInteger := tokens.ConvertHexBytesToInt(sourceCode)
		mapping := cmap.NewCidCharacterMapping(sourceInteger, destinationCode.IntVal())

		results = append(results, mapping)
	}

	builder.SetCidCharacterMappings(results)

	return nil
}

