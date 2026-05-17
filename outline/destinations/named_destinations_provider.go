package destinations

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Read reads named destinations from the PDF catalog dictionary.
// It supports both PDF 1.1 style (/Dests in catalog) and PDF 1.2+ style
// (name tree via /Names -> /Dests).
func Read(
	catalogDictionary *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	pages Pages,
	log logging.Log,
) *NamedDestinations {
	destinationsByName := make(map[string]ExplicitDestination)

	// PDF 1.1 style: /Dests entry in catalog dictionary.
	// Each key is a destination name and the value is either an array
	// defining the destination or a dictionary with a /D entry.
	if destsToken, ok := catalogDictionary.TryGet(tokens.Dests); ok {
		destsDict := resolveToDictionary(destsToken, scanner)
		if destsDict != nil {
			for key, value := range destsDict.Data() {
				if dest, ok := tryReadExplicitDestination(value, scanner, pages, log, false); ok {
					destinationsByName[key] = dest
				}
			}
		}
	} else if namesToken, ok := catalogDictionary.TryGet(tokens.Names); ok {
		// PDF 1.2+ style: name tree via /Names -> /Dests.
		namesDict := resolveToDictionary(namesToken, scanner)
		if namesDict != nil {
			if destsToken, ok := namesDict.TryGet(tokens.Dests); ok {
				destsNode := resolveToDictionary(destsToken, scanner)
				if destsNode != nil {
					flattenNameTree(destsNode, scanner, func(value tokens.Token) *ExplicitDestination {
						if dest, ok := tryReadExplicitDestination(value, scanner, pages, log, false); ok {
							return &dest
						}
						return nil
					}, destinationsByName)
				}
			}
		}
	}

	return NewNamedDestinations(destinationsByName, pages)
}

// tryReadExplicitDestination attempts to read an explicit destination from a token.
// The token may be an array (direct destination), or a dictionary with a /D entry
// whose value is such an array.
func tryReadExplicitDestination(
	value tokens.Token,
	scanner tokenization.PdfTokenScanner,
	pages Pages,
	log logging.Log,
	isRemoteDestination bool,
) (ExplicitDestination, bool) {
	// Try as array directly.
	if valueArray := resolveToArray(value, scanner); valueArray != nil {
		return TryGetExplicitDestination(valueArray, pages, log, isRemoteDestination)
	}

	// Try as dictionary with /D entry.
	if valueDict := resolveToDictionary(value, scanner); valueDict != nil {
		dToken, ok := valueDict.TryGet(tokens.D)
		if !ok {
			return ExplicitDestination{}, false
		}
		valueArray := resolveToArray(dToken, scanner)
		if valueArray == nil {
			return ExplicitDestination{}, false
		}
		return TryGetExplicitDestination(valueArray, pages, log, isRemoteDestination)
	}

	return ExplicitDestination{}, false
}

// TryGetExplicitDestination attempts to parse an explicit destination from an array token.
func TryGetExplicitDestination(
	array *tokens.ArrayToken,
	pages Pages,
	log logging.Log,
	isRemoteDestination bool,
) (ExplicitDestination, bool) {
	if array == nil || array.Length() == 0 {
		return ExplicitDestination{}, false
	}

	data := array.Data()

	getPossibleEntry := func(index int) *float64 {
		if index >= len(data) {
			return nil
		}
		if num, ok := data[index].(*tokens.NumericToken); ok {
			v := num.Data()
			return &v
		}
		return nil
	}

	pageNumber := 0
	pageToken := data[0]

	switch t := pageToken.(type) {
	case *tokens.IndirectReferenceToken:
		if isRemoteDestination {
			log.Error(fmt.Sprintf("TryGetExplicitDestination: Cannot use indirect reference for remote destination."))
			return ExplicitDestination{}, false
		}
		if pages == nil {
			return ExplicitDestination{}, false
		}
		page := pages.GetPageByReference(t.Data())
		if page == nil {
			return ExplicitDestination{}, false
		}
		if pn := page.PageNumber(); pn != nil {
			pageNumber = *pn
		} else {
			return ExplicitDestination{}, false
		}

	case *tokens.NumericToken:
		pageNumber = t.IntVal() + 1

	default:
		log.Error(fmt.Sprintf("TryGetExplicitDestination: No page number given in 'Dest': '%v', defaulting to 0.", array))
		pageNumber = 0
	}

	var destTypeToken *tokens.NameToken
	if len(data) > 1 {
		if name, ok := data[1].(*tokens.NameToken); ok {
			destTypeToken = name
		}
	}

	if destTypeToken == nil {
		log.Error(fmt.Sprintf("Missing name token as second argument to explicit destination: %v.", array))
		return NewExplicitDestination(pageNumber, FitPage, Empty), true
	}

	switch {
	case destTypeToken.Equals(tokens.XYZ):
		left := getPossibleEntry(2)
		top := getPossibleEntry(3)
		return NewExplicitDestination(pageNumber, XyzCoordinates, newCoords(left, top)), true

	case destTypeToken.Equals(tokens.Fit):
		return NewExplicitDestination(pageNumber, FitPage, Empty), true

	case destTypeToken.Equals(tokens.FitH):
		top := getPossibleEntry(2)
		return NewExplicitDestination(pageNumber, FitHorizontally, newCoords(nil, top)), true

	case destTypeToken.Equals(tokens.FitV):
		left := getPossibleEntry(2)
		return NewExplicitDestination(pageNumber, FitVertically, newCoords(left, nil)), true

	case destTypeToken.Equals(tokens.FitR):
		left := getPossibleEntry(2)
		bottom := getPossibleEntry(3)
		right := getPossibleEntry(4)
		top := getPossibleEntry(5)
		return NewExplicitDestination(pageNumber, FitRectangle, newCoordsRect(left, top, right, bottom)), true

	case destTypeToken.Equals(tokens.FitB):
		return NewExplicitDestination(pageNumber, FitBoundingBox, Empty), true

	case destTypeToken.Equals(tokens.FitBH):
		top := getPossibleEntry(2)
		return NewExplicitDestination(pageNumber, FitBoundingBoxHorizontally, newCoords(nil, top)), true

	case destTypeToken.Equals(tokens.FitBV):
		left := getPossibleEntry(2)
		return NewExplicitDestination(pageNumber, FitBoundingBoxVertically, newCoords(left, nil)), true
	}

	return ExplicitDestination{}, false
}

// resolveToArray resolves a token to an ArrayToken, following indirect references.
func resolveToArray(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.ArrayToken {
	if arr, ok := token.(*tokens.ArrayToken); ok {
		return arr
	}
	if ref, ok := token.(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(ref.Data())
		if obj == nil {
			return nil
		}
		return resolveToArray(obj.Data(), scanner)
	}
	return nil
}

// resolveToDictionary resolves a token to a DictionaryToken, following indirect references.
func resolveToDictionary(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	if dict, ok := token.(*tokens.DictionaryToken); ok {
		return dict
	}
	if ref, ok := token.(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(ref.Data())
		if obj == nil {
			return nil
		}
		return resolveToDictionary(obj.Data(), scanner)
	}
	return nil
}

// flattenNameTree recursively flattens a PDF name tree into the result map.
func flattenNameTree(
	nodeDict *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	valueFactory func(tokens.Token) *ExplicitDestination,
	result map[string]ExplicitDestination,
) {
	// Process /Names array (key-value pairs at this level).
	namesToken, ok := nodeDict.TryGet(tokens.Names)
	if ok {
		namesArr := resolveToArray(namesToken, scanner)
		if namesArr != nil {
			data := namesArr.Data()
			for i := 0; i < len(data)-1; i += 2 {
				key, ok := data[i].(interface{ Data() string })
				if !ok {
					continue
				}
				value := valueFactory(data[i+1])
				if value != nil {
					result[key.Data()] = *value
				}
			}
		}
	}

	// Process /Kids array (recursive child nodes).
	kidsToken, ok := nodeDict.TryGet(tokens.Kids)
	if ok {
		kidsArr := resolveToArray(kidsToken, scanner)
		if kidsArr != nil {
			for _, kid := range kidsArr.Data() {
				kidDict := resolveToDictionary(kid, scanner)
				if kidDict != nil {
					flattenNameTree(kidDict, scanner, valueFactory, result)
				}
			}
		}
	}
}

// newCoords creates coordinates with optional left and top values.
func newCoords(left, top *float64) *ExplicitDestinationCoordinates {
	return &ExplicitDestinationCoordinates{
		Left: left,
		Top:  top,
	}
}

// newCoordsRect creates coordinates with all four rectangle values.
func newCoordsRect(left, top, right, bottom *float64) *ExplicitDestinationCoordinates {
	return &ExplicitDestinationCoordinates{
		Left:   left,
		Top:    top,
		Right:  right,
		Bottom: bottom,
	}
}
