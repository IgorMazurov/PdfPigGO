package actions

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

// ActionProvider provides methods for retrieving PDF actions from dictionary entries.
type ActionProvider struct{}

// TryGetAction retrieves an action (/A) entry from a dictionary and constructs the corresponding Action.
// For GoTo, GoToR, and GoToE actions it also resolves the destination. The returned Action preserves
// typed information (e.g., Destination on GoTo actions) so callers can type-assert to concrete types.
// Returns true with a non-nil result on success.
func (ActionProvider) TryGetAction(
	dictionary *tokens.DictionaryToken,
	namedDestinations *destinations.NamedDestinations,
	scanner tokenization.PdfTokenScanner,
	log logging.Log,
) (Action, bool, error) {
	actionDict := resolveDictionary(dictionary, tokens.A, scanner)
	if actionDict == nil {
		return nil, false, nil
	}

	actionType := resolveName(actionDict, tokens.S, scanner)
	if actionType == nil {
		return nil, false, core.NewPdfDocumentFormatException(
			fmt.Sprintf("No action type (/S) specified for action: %v.", actionDict))
	}

	provider := destinations.DestinationProvider{}

	switch actionType.Data() {
	case tokens.GoTo.Data():
		dest, ok := provider.TryGetDestination(
			actionDict, tokens.D, namedDestinations, scanner, log, false)
		if ok {
			return NewGoToAction(dest), true, nil
		}

	case tokens.GoToR.Data():
		filename, filenameOk := util.TryGetOptionalStringDirect(actionDict, tokens.F, scanner)
		if !filenameOk {
			return nil, false, nil
		}
		dest, ok := provider.TryGetDestination(
			actionDict, tokens.D, namedDestinations, scanner, log, true)
		if ok {
			return NewGoToRAction(dest, filename), true, nil
		}

	case tokens.GoToE.Data():
		fileSpec, fileSpecOk := util.TryGetOptionalStringDirect(actionDict, tokens.F, scanner)
		if !fileSpecOk {
			fileSpec = ""
		}
		dest, ok := provider.TryGetDestination(
			actionDict, tokens.D, namedDestinations, scanner, log, true)
		if ok {
			return NewGoToEAction(dest, fileSpec), true, nil
		}

	case tokens.Uri.Data():
		uri, uriOk := util.TryGetOptionalStringDirect(actionDict, tokens.Uri, scanner)
		if !uriOk {
			uri = ""
		}
		return NewUriAction(uri), true, nil
	}

return nil, false, nil
}

// resolveDictionary looks up a name key in the dictionary and resolves any indirect reference.
func resolveDictionary(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	token := dictEntry(dict, name, scanner)
	if token == nil {
		return nil
	}
	result, ok := token.(*tokens.DictionaryToken)
	if !ok {
		return nil
	}
	return result
}

// resolveName looks up a name key in the dictionary and returns the resolved NameToken.
func resolveName(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) *tokens.NameToken {
	token := dictEntry(dict, name, scanner)
	if token == nil {
		return nil
	}
	result, ok := token.(*tokens.NameToken)
	if !ok {
		return nil
	}
	return result
}

// dictEntry looks up a key in the dictionary and resolves indirect references via scanner.
func dictEntry(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) tokens.Token {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil
	}
	return unwrapIndirect(token, scanner)
}

// unwrapIndirect unwraps an IndirectReferenceToken by reading through the scanner.
func unwrapIndirect(token tokens.Token, scanner tokenization.PdfTokenScanner) tokens.Token {
	indirect, isIndirect := token.(*tokens.IndirectReferenceToken)
	if !isIndirect || scanner == nil {
		return token
	}
	obj := scanner.Get(indirect.Data())
	if obj != nil {
		return obj.Data()
	}
	return token
}
