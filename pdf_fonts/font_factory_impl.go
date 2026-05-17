package pdffonts

import (
	"errors"
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokens"

	handlers "github.com/uglytoad/pdfpig/go/pdf_fonts/parser/handlers"
)

// fontFactoryImpl is a concrete implementation of FontFactory that delegates
// to type-specific handlers based on the /Subtype entry in the dictionary.
type fontFactoryImpl struct {
	log      logging.Log
	handlers map[*tokens.NameToken]handlers.FontHandler
}

// NewFontFactoryImpl creates a new FontFactory with the given type-specific handlers.
func NewFontFactoryImpl(
	log logging.Log,
	type0Handler handlers.FontHandler,
	trueTypeHandler handlers.FontHandler,
	type1Handler handlers.FontHandler,
	type3Handler handlers.FontHandler,
) FontFactory {
	return &fontFactoryImpl{
		log: log,
		handlers: map[*tokens.NameToken]handlers.FontHandler{
			tokens.Type0:     type0Handler,
			tokens.TrueType:  trueTypeHandler,
			tokens.Type1:     type1Handler,
			tokens.MmType1:   type1Handler,
			tokens.Type3:     type3Handler,
		},
	}
}

// Get returns the font described by the given dictionary.
func (f *fontFactoryImpl) Get(dictionary *tokens.DictionaryToken) (fonts.Font, error) {
	if dictionary == nil {
		return nil, nil
	}

	typeToken := getNameOrDefault(dictionary, tokens.Type)

	if typeToken != nil && typeToken != tokens.Font {
		msg := fmt.Sprintf("The font dictionary did not have type 'Font'. %v", dictionary)
		f.log.Error(msg)
	}

	subtype := getNameOrDefault(dictionary, tokens.Subtype)

	if subtype != nil {
		if handler, ok := f.handlers[subtype]; ok {
			result := handler.Generate(dictionary)
			return result, nil
		}
	}

	fallbacks := []*tokens.NameToken{tokens.Type1, tokens.TrueType}
	for _, fallback := range fallbacks {
		handler, ok := f.handlers[fallback]
		if !ok {
			continue
		}

		result := handler.Generate(dictionary)
		if result != nil {
			return result, nil
		}

		f.log.Error(fmt.Sprintf("Tried to parse font as fallback type: %v", fallback))
	}

	msg := "Parsing not implemented for fonts of type: "
	if subtype != nil {
		msg += fmt.Sprintf("%v", subtype)
	} else {
		msg += "(unknown)"
	}
	msg += ", please submit a pull request or an issue."

	return nil, errors.New(msg)
}

// getNameOrDefault returns the NameToken value for the given key, or nil if not found.
func getNameOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken) *tokens.NameToken {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil
	}
	if n, ok := token.(*tokens.NameToken); ok {
		return n
	}
	return nil
}
