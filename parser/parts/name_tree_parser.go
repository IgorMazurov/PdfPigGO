package parts

import (
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// FlattenNameTreeToDictionary flattens a PDF name tree into a map keyed by name string.
// The valuesFactory function is called for each value token in the /Names array to produce
// the corresponding TResult value. Child nodes referenced via /Kids are traversed recursively.
func FlattenNameTreeToDictionary[T any](
	nameTreeNode *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	valuesFactory func(tokens.Token) T,
) map[string]T {
	result := make(map[string]T)

	FlattenNameTree(nameTreeNode, scanner, valuesFactory, result)

	return result
}

// FlattenNameTree recursively flattens a PDF name tree node into the given result map.
// It processes /Names entries (key-value pairs) and /Kids entries (child nodes).
func FlattenNameTree[T any](
	nameTreeNode *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	valuesFactory func(tokens.Token) T,
	result map[string]T,
) {
	// Process /Names array: alternating key-value pairs.
	if namesToken, ok := nameTreeNode.TryGet(tokens.Names); ok {
		if nodeNames, ok := TryGet[*tokens.ArrayToken](namesToken, scanner); ok {
			data := nodeNames.Data()
			for i := 0; i+1 < len(data); i += 2 {
				key, ok := data[i].(interface{ Data() string })
				if !ok {
					continue
				}

				value := valuesFactory(data[i+1])

				result[key.Data()] = value
			}
		}
	}

	// Process /Kids array: recursive child dictionary nodes.
	if kidsToken, ok := nameTreeNode.TryGet(tokens.Kids); ok {
		if kids, ok := TryGet[*tokens.ArrayToken](kidsToken, scanner); ok {
			for _, kid := range kids.Data() {
				if kidDict, ok := TryGet[*tokens.DictionaryToken](kid, scanner); ok {
					FlattenNameTree(kidDict, scanner, valuesFactory, result)
				}
			}
		}
	}
}
