package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CodespaceRangeParser parses begincodespacerange/endcodespacerange blocks in CMap data.
type CodespaceRangeParser struct{}

// NewCodespaceRangeParser creates a new CodespaceRangeParser.
func NewCodespaceRangeParser() *CodespaceRangeParser {
	return &CodespaceRangeParser{}
}

// Parse reads codespace ranges from the token scanner and sets them on the builder.
func (p *CodespaceRangeParser) Parse(numeric *tokens.NumericToken, scanner tokenization.TokenScanner, builder *cmap.CharacterMapBuilder) error {
	count := numeric.IntVal()
	ranges := make([]*cmap.CodespaceRange, 0, count)

	for i := 0; i < count; i++ {
		if !scanner.Advance() {
			return fmt.Errorf("codespace range has reached an unexpected end")
		}

		if op, ok := scanner.Current().(*tokens.OperatorToken); ok && op.Data() == "endcodespacerange" {
			break
		}

		start, ok := scanner.Current().(*tokens.HexToken)
		if !ok {
			return fmt.Errorf("codespace range contains an unexpected token: %v", scanner.Current())
		}

		if !scanner.Advance() {
			return fmt.Errorf("codespace range contains an unexpected token: %v", scanner.Current())
		}

		end, ok := scanner.Current().(*tokens.HexToken)
		if !ok {
			return fmt.Errorf("codespace range contains an unexpected token: %v", scanner.Current())
		}

		rangeObj, err := cmap.NewCodespaceRange(start.Bytes(), end.Bytes())
		if err != nil {
			return fmt.Errorf("failed to create codespace range: %w", err)
		}

		ranges = append(ranges, rangeObj)
	}

	builder.SetCodespaceRanges(ranges)

	return nil
}

