package pdfpig

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TryGet attempts to get an entry from the dictionary by name. If the token is
// already of type T it returns directly; if it's an indirect reference it resolves
// through the scanner recursively. Returns false if the key doesn't exist or the
// resolved token isn't of type T.
func TryGet[T tokens.Token](dictionary *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (T, bool) {
	token, ok := dictionary.TryGet(name)
	if !ok {
		var zero T
		return zero, false
	}

	if result, ok := any(token).(T); ok {
		return result, true
	}

	ref, ok := any(token).(*tokens.IndirectReferenceToken)
	if !ok {
		var zero T
		return zero, false
	}

	return parts.TryGet[T](ref, scanner)
}

// Get gets an entry from the dictionary by name, resolving indirect references
// through the scanner if necessary. Returns an error if the key doesn't exist or
// the resolved token isn't of type T.
func Get[T tokens.Token](dictionary *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (T, error) {
	token, ok := dictionary.TryGet(name)
	if !ok {
		var zero T
		return zero, core.NewPdfDocumentFormatException(
			fmt.Sprintf("dictionary does not contain token with name %v of type %T", name, zero))
	}

	if result, ok := any(token).(T); ok {
		return result, nil
	}

	ref, ok := any(token).(*tokens.IndirectReferenceToken)
	if !ok {
		var zero T
		return zero, core.NewPdfDocumentFormatException(
			fmt.Sprintf("dictionary does not contain token with name %v of type %T", name, zero))
	}

	return parts.GetByToken[T](ref, scanner)
}

// Decode decodes the stream data by applying each filter in sequence. Size limits
// are enforced between filters to prevent malicious decompression bombs.
func Decode(stream *tokens.StreamToken, provider filters.FilterProvider) ([]byte, error) {
	fl, err := provider.GetFilters(stream.StreamDictionary)
	if err != nil {
		return nil, err
	}

	totalMaxEstSize := float64(len(stream.Data())) * 100

	transform := stream.Data()

	for i, filter := range fl {
		totalMaxEstSize *= getEstimatedSizeMultiplier(filter)

		transform, err = filter.Decode(transform, stream.StreamDictionary, provider, i)
		if err != nil {
			return nil, err
		}

		if i < len(fl)-1 && float64(len(transform)) > totalMaxEstSize {
			return nil, core.NewPdfDocumentFormatException(
				fmt.Sprintf("decoded stream size exceeds the estimated maximum size. Current decoded stream length: %d, %d filters applied out of %d",
					len(transform), i+1, len(fl)))
		}
	}

	return transform, nil
}

// DecodeWithScanner decodes the stream data using a scanner-aware filter provider
// that can resolve indirect references in the /Filter or /F dictionary entries.
func DecodeWithScanner(stream *tokens.StreamToken, provider filters.LookupFilterProvider, scanner tokenization.PdfTokenScanner) ([]byte, error) {
	fl, err := provider.GetFiltersWithScanner(stream.StreamDictionary, scanner)
	if err != nil {
		return nil, err
	}

	totalMaxEstSize := float64(len(stream.Data())) * 100

	transform := stream.Data()

	for i, filter := range fl {
		totalMaxEstSize *= getEstimatedSizeMultiplier(filter)

		transform, err = filter.Decode(transform, stream.StreamDictionary, provider, i)
		if err != nil {
			return nil, err
		}

		if i < len(fl)-1 && float64(len(transform)) > totalMaxEstSize {
			return nil, core.NewPdfDocumentFormatException(
				fmt.Sprintf("decoded stream size exceeds the estimated maximum size. Current decoded stream length: %d, %d filters applied out of %d",
					len(transform), i+1, len(fl)))
		}
	}

	return transform, nil
}

func getEstimatedSizeMultiplier(filter filters.Filter) float64 {
	switch filter.(type) {
	case *filters.AsciiHexDecodeFilter:
		return 0.5
	case *filters.Ascii85Filter:
		return 0.8
	case *filters.RunLengthFilter:
		return 1.5
	case *filters.LzwFilter:
		return 2
	case *filters.FlateFilter:
		return 10
	default:
		return 1000
	}
}

// Resolve recursively resolves all indirect references within a token tree.
// For dictionaries and arrays it traverses child entries; for streams it resolves
// the stream dictionary. The visited set prevents infinite loops from cyclic
// references. Pass nil to use a fresh visited set.
func Resolve(token tokens.Token, scanner tokenization.PdfTokenScanner, visited map[core.IndirectReference]bool) (tokens.Token, error) {
	if visited == nil {
		visited = make(map[core.IndirectReference]bool)
	}

	resolved, err := resolveInternal(token, scanner, visited)
	if resolved == nil {
		return nil, err
	}
	return *resolved, err
}

func resolveInternal(token tokens.Token, scanner tokenization.PdfTokenScanner, visited map[core.IndirectReference]bool) (*tokens.Token, error) {
	if token == nil {
		var zero tokens.Token
		return &zero, nil
	}

	if stream, ok := any(token).(*tokens.StreamToken); ok {
		resolvedDict, err := resolveInternal(stream.StreamDictionary, scanner, visited)
		if err != nil {
			return nil, err
		}
		dict, _ := (*resolvedDict).(*tokens.DictionaryToken)
		newStream, err := tokens.NewStreamToken(dict, stream.Data())
		if err != nil {
			return nil, err
		}
		t := tokens.Token(newStream)
		return &t, nil
	}

	if dict, ok := any(token).(*tokens.DictionaryToken); ok {
		resolvedItems := make(map[string]tokens.Token)
		for k, v := range dict.Data() {
			value := v
			if ref, ok := any(v).(*tokens.IndirectReferenceToken); ok {
				if visited[ref.Data()] {
					continue
				}
				obj := scanner.Get(ref.Data())
				if obj != nil {
					value = obj.Data()
				}
				visited[ref.Data()] = true
			}
			resolvedVal, err := resolveInternal(value, scanner, visited)
			if err != nil {
				return nil, err
			}
			if resolvedVal != nil {
				resolvedItems[k] = *resolvedVal
			}
		}

		if len(resolvedItems) > len(dict.Data()) {
			return nil, fmt.Errorf("resolved more items than were present in the original dictionary")
		}

		if len(resolvedItems) < len(dict.Data()) {
			for missing := range dict.Data() {
				if _, exists := resolvedItems[missing]; !exists {
					if ref, ok := any(dict.Data()[missing]).(*tokens.IndirectReferenceToken); ok {
						resolvedVal, err := resolveInternal(ref, scanner, visited)
						if err != nil {
							return nil, err
						}
						if resolvedVal != nil {
							resolvedItems[missing] = *resolvedVal
						}
					}
				}
			}
		}

		nameKeyed := make(map[*tokens.NameToken]tokens.Token, len(resolvedItems))
		for k, v := range resolvedItems {
			nameKeyed[tokens.Create(k)] = v
		}
		newDict, err := tokens.NewDictionary(nameKeyed)
		if err != nil {
			return nil, err
		}
		t := tokens.Token(newDict)
		return &t, nil
	}

	if arr, ok := any(token).(*tokens.ArrayToken); ok {
		resolvedItems := make([]tokens.Token, 0, arr.Length())
		data := arr.Data()
		for i := 0; i < len(data); i++ {
			value := data[i]
			if ref, ok := any(value).(*tokens.IndirectReferenceToken); ok {
				obj := scanner.Get(ref.Data())
				if obj != nil {
					value = obj.Data()
				}
			}
			resolvedVal, err := resolveInternal(value, scanner, visited)
			if err != nil {
				return nil, err
			}
			if resolvedVal != nil {
				resolvedItems = append(resolvedItems, *resolvedVal)
			}
		}
		newArr := tokens.NewArrayToken(resolvedItems)
		t := tokens.Token(newArr)
		return &t, nil
	}

	if ref, ok := any(token).(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(ref.Data())
		if obj != nil {
			t := obj.Data()
			return &t, nil
		}
		var zero tokens.Token
		return &zero, nil
	}

	t := token
	return &t, nil
}
