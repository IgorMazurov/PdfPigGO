package destinations

import (
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// DestinationProvider provides methods for retrieving PDF destinations.
type DestinationProvider struct{}

// TryGetDestination attempts to get an explicit or named destination from a dictionary entry.
// It first tries to interpret the value as an array (explicit destination), then as a string
// (named destination key). Indirect references are resolved through the scanner.
func (DestinationProvider) TryGetDestination(
	dictionary *tokens.DictionaryToken,
	destinationToken *tokens.NameToken,
	namedDestinations *NamedDestinations,
	scanner tokenization.PdfTokenScanner,
	log logging.Log,
	isRemoteDestination bool,
) (ExplicitDestination, bool) {
	token, ok := dictionary.TryGet(destinationToken)
	if !ok {
		var zero ExplicitDestination
		return zero, false
	}

	token = unwrapIndirectRef(token, scanner)

	if arrayToken, ok := token.(*tokens.ArrayToken); ok {
		if namedDestinations != nil {
			return namedDestinations.TryGetExplicitDestination(arrayToken, log, isRemoteDestination)
		}
		// Fallback: try parsing explicit destination without NamedDestinations.
		// Only works for numeric page references (not indirect refs).
		if data := arrayToken.Data(); len(data) > 0 {
			if _, isIndirect := data[0].(*tokens.IndirectReferenceToken); isIndirect {
				return ExplicitDestination{}, false
			}
			return TryGetExplicitDestination(arrayToken, nil, log, isRemoteDestination)
		}
		return ExplicitDestination{}, false
	}

	if stringToken, ok := token.(interface{ Data() string }); ok {
		if namedDestinations == nil {
			var zero ExplicitDestination
			return zero, false
		}
		return namedDestinations.TryGetName(stringToken.Data())
	}

	var zero ExplicitDestination
	return zero, false
}

// unwrapIndirectRef resolves an IndirectReferenceToken through the scanner.
func unwrapIndirectRef(token tokens.Token, scanner tokenization.PdfTokenScanner) tokens.Token {
	if token == nil || scanner == nil {
		return token
	}

	refToken, ok := token.(*tokens.IndirectReferenceToken)
	if !ok {
		return token
	}

	obj := scanner.Get(refToken.Data())
	if obj == nil {
		return token
	}

	return obj.Data()
}
