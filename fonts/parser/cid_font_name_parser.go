package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CidFontNameParser parses CID font name directives (WMode, CMapName, etc.) in CMap data.
type CidFontNameParser struct{}

// NewCidFontNameParser creates a new CidFontNameParser.
func NewCidFontNameParser() *CidFontNameParser {
	return &CidFontNameParser{}
}

// Parse handles name token directives in the CMap and updates the builder accordingly.
func (p *CidFontNameParser) Parse(name *tokens.NameToken, scanner tokenization.TokenScanner, builder *cmap.CharacterMapBuilder) error {
	switch name.Data() {
	case "WMode":
		if numeric := tryReadNumeric(scanner); numeric != nil {
			builder.SetWMode(numeric.IntVal())
		}

	case "CMapName":
		if scanner.Advance() {
			if n, ok := scanner.Current().(*tokens.NameToken); ok {
				builder.SetName(n.Data())
			}
		}

	case "CMapVersion":
		if !scanner.Advance() {
			return nil
		}
		next := scanner.Current()
		if num, ok := next.(*tokens.NumericToken); ok {
			v := fmt.Sprintf("%g", num.Data())
			builder.SetVersion(&v)
		} else if str, ok := next.(*tokens.StringToken); ok {
			v := str.Data()
			builder.SetVersion(&v)
		}

	case "CMapType":
		if numeric := tryReadNumeric(scanner); numeric != nil {
			builder.SetType(numeric.IntVal())
		}

	case "Registry":
		if scanner.Advance() {
			if str, ok := scanner.Current().(*tokens.StringToken); ok {
				builder.SystemInfoBuilder().SetRegistry(str.Data())
			}
		}

	case "Ordering":
		if scanner.Advance() {
			if str, ok := scanner.Current().(*tokens.StringToken); ok {
				builder.SystemInfoBuilder().SetOrdering(str.Data())
			}
		}

	case "Supplement":
		if numeric := tryReadNumeric(scanner); numeric != nil {
			builder.SystemInfoBuilder().SetSupplement(numeric.IntVal())
		}

	case "CIDSystemInfo":
		if scanner.Advance() {
			if dict, ok := scanner.Current().(*tokens.DictionaryToken); ok {
				builder.SetCharacterIdentifierSystemInfo(getCharacterIdentifier(dict))
			}
		}
	}

	return nil
}

func tryReadNumeric(scanner tokenization.TokenScanner) *tokens.NumericToken {
	scanner.Advance()
	if num, ok := scanner.Current().(*tokens.NumericToken); ok {
		return num
	}
	return nil
}

func getCharacterIdentifier(dictionary *tokens.DictionaryToken) pdffonts.CharacterIdentifierSystemInfo {
	registry := "Adobe"
	if val, found := dictionary.TryGet(tokens.Registry); found {
		if s, ok := val.(*tokens.StringToken); ok {
			registry = s.Data()
		}
	}

	ordering := ""
	if val, found := dictionary.TryGet(tokens.Ordering); found {
		if s, ok := val.(*tokens.StringToken); ok {
			ordering = s.Data()
		}
	}

	supplement := 0
	if val, found := dictionary.TryGet(tokens.Supplement); found {
		if n, ok := val.(*tokens.NumericToken); ok {
			supplement = n.IntVal()
		}
	}

	return pdffonts.NewCharacterIdentifierSystemInfo(registry, ordering, supplement)
}

