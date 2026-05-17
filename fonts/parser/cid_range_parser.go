package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CidRangeParser parses begincidrange/endcidrange blocks in CMap data.
type CidRangeParser struct{}

// NewCidRangeParser creates a new CidRangeParser.
func NewCidRangeParser() *CidRangeParser {
	return &CidRangeParser{}
}

// Parse reads CID range mappings from the token scanner and adds them to the builder.
func (p *CidRangeParser) Parse(numeric *tokens.NumericToken, scanner tokenization.TokenScanner, builder *cmap.CharacterMapBuilder) error {
	count := numeric.IntVal()

	for i := 0; i < count; i++ {
		if !scanner.Advance() {
			return fmt.Errorf("could not find the starting hex token for the CIDRange in this font")
		}

		startHexToken, ok := scanner.Current().(*tokens.HexToken)
		if !ok {
			return fmt.Errorf("could not find the starting hex token for the CIDRange in this font")
		}

		if !scanner.Advance() {
			return fmt.Errorf("could not find the end hex token for the CIDRange in this font")
		}

		endHexToken, ok := scanner.Current().(*tokens.HexToken)
		if !ok {
			return fmt.Errorf("could not find the end hex token for the CIDRange in this font")
		}

		if !scanner.Advance() {
			return fmt.Errorf("could not find the starting CID numeric token for the CIDRange in this font")
		}

		mappedCode, ok := scanner.Current().(*tokens.NumericToken)
		if !ok {
			return fmt.Errorf("could not find the starting CID numeric token for the CIDRange in this font")
		}

		start := tokens.ConvertHexBytesToInt(startHexToken)
		end := tokens.ConvertHexBytesToInt(endHexToken)

		rangeObj, err := cmap.NewCidRange(start, end, mappedCode.IntVal())
		if err != nil {
			return fmt.Errorf("cid range creation failed: %w", err)
		}

		builder.AddCidRange(rangeObj)
	}

	return nil
}

