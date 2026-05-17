// Package colors provides types for PDF color handling.
package colors

import (
	"github.com/uglytoad/pdfpig/go/tokens"
)

const srgbIec61966OutputCondition = "sRGB IEC61966-2.1"
const registryName = "http://www.color.org"

// ObjectWriter is the callback type for writing a token and obtaining an indirect reference.
type ObjectWriter func(token tokens.Token) *tokens.IndirectReferenceToken

// CompressFunc is a callback that compresses raw bytes using FlateDecode.
type CompressFunc func([]byte) []byte

// GetOutputIntentsArray builds the /OutputIntents array token required for PDF/A compliance.
func GetOutputIntentsArray(objectWriter ObjectWriter, compress CompressFunc) *tokens.ArrayToken {
	rgbColorCondition := tokens.NewStringToken(srgbIec61966OutputCondition, tokens.Iso88591)

	profileBytes := ProfileStreamReaderGetSRgb2014()
	compressedBytes := compress(profileBytes)

	profileStreamDictionary := map[*tokens.NameToken]tokens.Token{
		tokens.Length: tokens.NewNumericTokenFromInt(len(compressedBytes)),
		tokens.N:      tokens.Three,
		tokens.Filter: tokens.FlateDecode,
	}

	dictToken, _ := tokens.NewDictionary(profileStreamDictionary)
	stream, _ := tokens.NewStreamToken(dictToken, compressedBytes)

	written := objectWriter(stream)

	outputIntentDict := map[*tokens.NameToken]tokens.Token{
		tokens.Type:                      tokens.OutputIntent,
		tokens.S:                         tokens.GtsPdfa1,
		tokens.OutputCondition:           rgbColorCondition,
		tokens.OutputConditionIdentifier: rgbColorCondition,
		tokens.RegistryName:              tokens.NewStringToken(registryName, tokens.Iso88591),
		tokens.Info:                      rgbColorCondition,
		tokens.DestOutputProfile:         written,
	}

	dict, _ := tokens.NewDictionary(outputIntentDict)

	return tokens.NewArrayToken([]tokens.Token{dict})
}
