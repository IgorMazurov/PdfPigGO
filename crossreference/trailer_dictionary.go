package crossreference

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TrailerDictionary contains information for interpreting the cross-reference table.
type TrailerDictionary struct {
	size                        int
	previousCrossReferenceOffset *int64
	root                        core.IndirectReference
	info                        tokens.Token
	identifier                  []tokens.DataToken[string]
	encryptionToken             tokens.Token
}

// NewTrailerDictionary creates a new TrailerDictionary from the parsed dictionary token.
func NewTrailerDictionary(dictionary *tokens.DictionaryToken, isLenientParsing bool) (*TrailerDictionary, error) {
	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}

	size, err := getInt(dictionary, tokens.Size)
	if err != nil {
		return nil, err
	}

	prevOffset := getLongOrDefault(dictionary, tokens.Prev)

	rootRef, ok, err := getIndirectReferenceToken(dictionary, tokens.Root)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("no root token was found in the trailer dictionary: %v", dictionary))
	}

	var info tokens.Token
	if infoToken, found := dictionary.TryGet(tokens.Info); found {
		if !isLenientParsing {
			if _, ok := infoToken.(*tokens.IndirectReferenceToken); !ok {
				return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("The info token in the trailer dictionary should only contain indirect references, instead got: %v.", infoToken))
			}
		}
		info = infoToken
	}

	identifier := make([]tokens.DataToken[string], 0)
	if arrToken, found := dictionary.TryGet(tokens.Id); found {
		if arr, ok := arrToken.(*tokens.ArrayToken); ok {
			for _, token := range arr.Data() {
				switch t := token.(type) {
				case *tokens.StringToken:
					identifier = append(identifier, t)
				case *tokens.HexToken:
					identifier = append(identifier, t)
				}
			}
		}
	}

	var encryptionToken tokens.Token
	if encToken, found := dictionary.TryGet(tokens.Encrypt); found {
		encryptionToken = encToken
	}

	return &TrailerDictionary{
		size:                        size,
		previousCrossReferenceOffset: prevOffset,
		root:                        rootRef,
		info:                        info,
		identifier:                  identifier,
		encryptionToken:             encryptionToken,
	}, nil
}

// Size returns the total number of object entries across both the original cross-reference table
// and in any incremental updates. Any object in a cross-reference section whose number is greater
// than this value is ignored and considered missing.
func (t *TrailerDictionary) Size() int {
	return t.size
}

// PreviousCrossReferenceOffset returns the offset in bytes to the previous cross-reference table
// or stream if the document has more than one cross-reference section. Returns nil if not present.
func (t *TrailerDictionary) PreviousCrossReferenceOffset() *int64 {
	return t.previousCrossReferenceOffset
}

// Root returns the object reference for the document's catalog dictionary.
func (t *TrailerDictionary) Root() core.IndirectReference {
	return t.root
}

// Info returns the object reference for the document's information dictionary if it contains one.
func (t *TrailerDictionary) Info() tokens.Token {
	return t.info
}

// Identifier returns the list of string tokens which act as file identifiers.
func (t *TrailerDictionary) Identifier() []tokens.DataToken[string] {
	return t.identifier
}

// EncryptionToken returns the document's encryption dictionary.
func (t *TrailerDictionary) EncryptionToken() tokens.Token {
	return t.encryptionToken
}

// String returns the string representation of the trailer dictionary.
func (t *TrailerDictionary) String() string {
	return fmt.Sprintf("Size: %d, Root: %v", t.size, t.root)
}

func getInt(dictionary *tokens.DictionaryToken, name *tokens.NameToken) (int, error) {
	token, found := dictionary.TryGet(name)
	if !found {
		return 0, core.NewPdfDocumentFormatException(fmt.Sprintf("the dictionary did not contain a number with the key %v", name))
	}
	num, ok := token.(*tokens.NumericToken)
	if !ok {
		return 0, core.NewPdfDocumentFormatException(fmt.Sprintf("the dictionary did not contain a number with the key %v. Dictionary was: %v", name, dictionary))
	}
	return num.IntVal(), nil
}

func getLongOrDefault(dictionary *tokens.DictionaryToken, name *tokens.NameToken) *int64 {
	token, found := dictionary.TryGet(name)
	if !found {
		return nil
	}
	num, ok := token.(*tokens.NumericToken)
	if !ok {
		return nil
	}
	v := num.LongVal()
	return &v
}

func getIndirectReferenceToken(dictionary *tokens.DictionaryToken, name *tokens.NameToken) (core.IndirectReference, bool, error) {
	token, found := dictionary.TryGet(name)
	if !found {
		return core.IndirectReference{}, false, nil
	}
	ref, ok := token.(*tokens.IndirectReferenceToken)
	if !ok {
		return core.IndirectReference{}, false, nil
	}
	return ref.Data(), true, nil
}
