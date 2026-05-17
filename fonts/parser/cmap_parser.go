package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	cmapresources "github.com/uglytoad/pdfpig/go/fonts/resources/cmap"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var (
	baseFontRangeParser   = &BaseFontRangeParser{}
	baseFontCharacterParser = &BaseFontCharacterParser{}
	cidRangeParser        = &CidRangeParser{}
	cidFontNameParser     = &CidFontNameParser{}
	codespaceRangeParser  = &CodespaceRangeParser{}
	cidCharacterParser    = &CidCharacterParser{}
)

// CMapParser parses PostScript CMap data into a character map.
type CMapParser struct{}

// NewCMapParser creates a new CMapParser with the default sub-parsers.
func NewCMapParser() *CMapParser {
	return &CMapParser{}
}

// Parse reads CMap data from the input and builds a character map.
func (p *CMapParser) Parse(inputBytes core.InputBytes) (*cmap.CMap, error) {
	scanner := tokenization.NewCoreTokenScanner(
		inputBytes,
		false,
		core.Infinite,
		tokenization.ScannerScopeNone,
		map[*tokens.NameToken][]*tokens.NameToken{
			tokens.CidSystemInfo: {tokens.Registry, tokens.Ordering, tokens.Supplement},
		},
		false,
		false,
	)

	builder := cmap.NewCharacterMapBuilder()

	var previousToken tokens.Token

	for scanner.Advance() {
		token := scanner.Current()

		if operatorToken, ok := token.(*tokens.OperatorToken); ok {
			switch operatorToken.Data() {
			case "usecmap":
				name, ok := previousToken.(*tokens.NameToken)
				if !ok {
					return nil, fmt.Errorf("unexpected token preceding external cmap call: %v", previousToken)
				}
				external, found := p.TryParseExternal(name.Data())
				if !found {
					return nil, fmt.Errorf("unable to find external CMap: %s", name.Data())
				}
				builder.UseCMap(external)

			case "begincodespacerange":
				numeric, ok := previousToken.(*tokens.NumericToken)
				if !ok {
					return nil, fmt.Errorf("unexpected token preceding start of codespace range: %v", previousToken)
				}
				if err := codespaceRangeParser.Parse(numeric, scanner, builder); err != nil {
					return nil, err
				}

			case "beginbfchar":
				numeric, ok := previousToken.(*tokens.NumericToken)
				if !ok {
					return nil, fmt.Errorf("unexpected token preceding start of base font characters: %v", previousToken)
				}
				if err := baseFontCharacterParser.Parse(numeric, scanner, builder); err != nil {
					return nil, err
				}

			case "beginbfrange":
				numeric, ok := previousToken.(*tokens.NumericToken)
				if !ok {
					return nil, fmt.Errorf("unexpected token preceding start of base font character ranges: %v", previousToken)
				}
				if err := baseFontRangeParser.Parse(numeric, scanner, builder); err != nil {
					return nil, err
				}

			case "begincidchar":
				numeric, ok := previousToken.(*tokens.NumericToken)
				if !ok {
					return nil, fmt.Errorf("unexpected token preceding start of Cid character mapping: %v", previousToken)
				}
				if err := cidCharacterParser.Parse(numeric, scanner, builder); err != nil {
					return nil, err
				}

			case "begincidrange":
				numeric, ok := previousToken.(*tokens.NumericToken)
				if !ok {
					return nil, fmt.Errorf("unexpected token preceding start of Cid ranges: %v", previousToken)
				}
				if err := cidRangeParser.Parse(numeric, scanner, builder); err != nil {
					return nil, err
				}
			}
		} else if nameToken, ok := token.(*tokens.NameToken); ok {
			if err := cidFontNameParser.Parse(nameToken, scanner, builder); err != nil {
				return nil, err
			}
		}

		previousToken = token
	}

	cmapResult, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return cmapResult, nil
}

// TryParseExternal attempts to load a built-in CMap by name from embedded resources.
func (p *CMapParser) TryParseExternal(name string) (*cmap.CMap, bool) {
	data, err := cmapresources.Files.ReadFile(name)
	if err != nil {
		return nil, false
	}

	inputBytes := core.NewMemoryInputBytes(data)
	cmapResult, parseErr := p.Parse(inputBytes)
	if parseErr != nil {
		return nil, false
	}

	return cmapResult, true
}
