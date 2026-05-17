package filters

import (
	"errors"
	"fmt"

	"github.com/uglytoad/pdfpig/go/tokens"
)

var (
	errNilStreamDictionary = errors.New("stream dictionary cannot be nil")
	errNegativeIndex       = errors.New("index must be 0 or greater")
)

// GetFilterParameters retrieves the decode parameters for a given filter index
// from a stream dictionary. If the /Filter entry is a single name, it returns
// the /DecodeParms dictionary directly. If /Filter is an array of names, it
// returns the DictionaryToken at the corresponding index inside the /DecodeParms
// array. Returns an empty DictionaryToken when no parameters are found.
func GetFilterParameters(streamDictionary *tokens.DictionaryToken, index int) (*tokens.DictionaryToken, error) {
	if streamDictionary == nil {
		return nil, fmt.Errorf("get filter parameters: %w", errNilStreamDictionary)
	}

	if index < 0 {
		return nil, fmt.Errorf("get filter parameters: %w", errNegativeIndex)
	}

	filter := getObjectOrDefault(streamDictionary, tokens.Filter, tokens.F)
	parameters := getObjectOrDefault(streamDictionary, tokens.DecodeParms, tokens.Dp)

	switch filter.(type) {
	case *tokens.NameToken:
		if dict, ok := parameters.(*tokens.DictionaryToken); ok {
			return dict, nil
		}
	case *tokens.ArrayToken:
		if arr, ok := parameters.(*tokens.ArrayToken); ok {
			data := arr.Data()
			if index < len(data) {
				if dict, ok := data[index].(*tokens.DictionaryToken); ok {
					return dict, nil
				}
			}
		}
	}

	return tokens.WithMap(make(map[string]tokens.Token))
}

// getObjectOrDefault tries to get a token from the dictionary using the first
// name key; if not found, falls back to the second name key. Returns nil when
// neither key exists.
func getObjectOrDefault(dict *tokens.DictionaryToken, first, fallback *tokens.NameToken) tokens.Token {
	if token, ok := dict.TryGet(first); ok {
		return token
	}
	if token, ok := dict.TryGet(fallback); ok {
		return token
	}
	return nil
}
