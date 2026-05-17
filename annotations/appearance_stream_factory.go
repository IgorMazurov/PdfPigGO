package annotations

import (
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TryCreateAppearanceStream attempts to create an AppearanceStream from an appearance
// dictionary entry identified by the given name key. It handles two cases:
// 1. The entry resolves to a single stream (stateless appearance).
// 2. The entry is a sub-dictionary mapping state names to streams (stateful appearance).
func TryCreateAppearanceStream(
	appearanceDictionary *tokens.DictionaryToken,
	name *tokens.NameToken,
	scanner tokenization.PdfTokenScanner,
) (*AppearanceStream, bool) {
	if appearanceDictionary == nil || name == nil || scanner == nil {
		return nil, false
	}

	token, ok := appearanceDictionary.TryGet(name)
	if !ok {
		return nil, false
	}

	switch t := token.(type) {
	case *tokens.IndirectReferenceToken:
		ref := t.Data()
		objToken := scanner.Get(ref)
		if objToken == nil {
			return nil, false
		}

		streamToken, ok := objToken.Data().(*tokens.StreamToken)
		if !ok || streamToken == nil {
			return nil, false
		}

		appearanceStream, err := NewAppearanceStream(streamToken)
		if err != nil {
			return nil, false
		}

		return appearanceStream, true

	case *tokens.DictionaryToken:
		dict := make(map[string]*tokens.StreamToken)

		for state, value := range t.Data() {
			if indRef, ok := value.(*tokens.IndirectReferenceToken); ok {
				objToken := scanner.Get(indRef.Data())
				if objToken == nil {
					continue
				}

				if streamToken, ok := objToken.Data().(*tokens.StreamToken); ok && streamToken != nil {
					dict[state] = streamToken
				}
			}
		}

		if len(dict) > 0 {
			appearanceStream, err := NewStatefulAppearanceStream(dict)
			if err != nil {
				return nil, false
			}

			return appearanceStream, true
		}
	}

	return nil, false
}
