package tokenization

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// DictionaryTokenizer tokenizes PDF dictionary objects starting with '<<'.
type DictionaryTokenizer struct {
	usePdfDocEncoding bool
	requiredKeys      []*tokens.NameToken
	useLenientParsing bool
	stackDepthGuard   *core.StackDepthGuard
}

var _ Tokenizer = (*DictionaryTokenizer)(nil)

// NewDictionaryTokenizer creates a new DictionaryTokenizer.
func NewDictionaryTokenizer(usePdfDocEncoding bool, stackDepthGuard *core.StackDepthGuard, requiredKeys []*tokens.NameToken, useLenientParsing bool) *DictionaryTokenizer {
	return &DictionaryTokenizer{
		usePdfDocEncoding: usePdfDocEncoding,
		stackDepthGuard:   stackDepthGuard,
		requiredKeys:      requiredKeys,
		useLenientParsing: useLenientParsing,
	}
}

// ReadsNextByte returns false because this tokenizer does not read past the closing '>>'.
func (t *DictionaryTokenizer) ReadsNextByte() bool {
	return false
}

// Tokenize attempts to tokenize a PDF dictionary starting with '<'.
// Returns the DictionaryToken and true on success, or nil and false if currentByte is not '<'.
func (t *DictionaryTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '<' {
		return nil, false
	}

	start := input.CurrentOffset()

	var token tokens.Token
	var ok bool

	// When requiredKeys is set, catch format exceptions to fall back to required-keys mode.
	// Otherwise, let the panic propagate up (matching C# behavior).
	func() {
		if t.requiredKeys != nil {
			defer func() {
				if r := recover(); r != nil {
					if _, isFormatErr := r.(*core.PdfDocumentFormatException); isFormatErr {
						token = nil
						ok = false
						return
					}
					panic(r)
				}
			}()
		}
		token, ok = t.tokenizeInternal(currentByte, input, false)
	}()

	if !ok && t.requiredKeys == nil {
		return nil, false
	}
	if !ok {
		input.Seek(start, 0)
		token, ok = t.tokenizeInternal(currentByte, input, true)
	}

	return token, ok
}

func (t *DictionaryTokenizer) tokenizeInternal(currentByte byte, input core.InputBytes, useRequiredKeys bool) (tokens.Token, bool) {
	if currentByte != '<' {
		return nil, false
	}

	// Skip whitespace after first '<' to find second '<' of '<<'
	foundNextOpenBrace := false
	for input.MoveNext() {
		if input.CurrentByte() == '<' {
			foundNextOpenBrace = true
			break
		}
		if !core.IsWhitespace(input.CurrentByte()) {
			break
		}
	}

	if !foundNextOpenBrace {
		return nil, false
	}

	scanner := NewCoreTokenScanner(input, t.usePdfDocEncoding, t.stackDepthGuard, ScannerScopeDictionary, nil, t.useLenientParsing, false)

	allTokens := make([]tokens.Token, 0)

	for scanner.Advance() {
		if _, isComment := scanner.Current().(*tokens.CommentToken); isComment {
			continue
		}

		allTokens = append(allTokens, scanner.Current())

		// Check if we have enough key/values for each required key
		if useRequiredKeys && len(allTokens) >= len(t.requiredKeys)*2 {
			proposedEntries, err := convertToDictionary(allTokens, t.useLenientParsing)
			if err != nil {
				// Malformed dictionary during required-keys check; continue scanning.
				continue
			}

			isAcceptable := true
			for _, rk := range t.requiredKeys {
				found := false
				for _, entry := range proposedEntries {
					if entry.Key == rk.Data() {
						found = true
						break
					}
				}
				if !found {
					isAcceptable = false
					break
				}
			}

			if isAcceptable {
				dictToken, _ := tokens.NewOrderedDictionary(proposedEntries)
				return dictToken, true
			}
		}
	}

	entries, err := convertToDictionary(allTokens, t.useLenientParsing)
	if err != nil {
		return nil, false
	}

	dictToken, _ := tokens.NewOrderedDictionary(entries)
	return dictToken, true
}

// convertToDictionary converts a flat list of tokens into a key-value dictionary.
// Returns an ordered slice of entries preserving insertion order, plus a map for lookups.
// Returns an error if a non-name token is encountered where a dictionary key is expected
// (and lenient parsing is disabled). This matches C#'s PdfDocumentFormatException behavior.
func convertToDictionary(allTokens []tokens.Token, useLenientParsing bool) ([]tokens.DictEntry, error) {
	var entries []tokens.DictEntry

	var key *tokens.NameToken
	for i := 0; i < len(allTokens); i++ {
		token := allTokens[i]

		if key == nil {
			name, ok := token.(*tokens.NameToken)
			if ok {
				key = name
				continue
			}

	if useLenientParsing {
			continue
		}

		panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Expected name as dictionary key, instead got: %v", token)))
		}

		var value tokens.Token
		// Combine indirect references, e.g. 12 0 R
		if num, ok := token.(*tokens.NumericToken); ok {
			gen := peekNext(allTokens, i+1)
			r := peekNext(allTokens, i+2)

			if gen != nil && r == tokens.OpR {
				genNum := gen.(*tokens.NumericToken)
				ref, _ := core.NewIndirectReference(num.LongVal(), genNum.IntVal())
				value = tokens.NewIndirectReferenceToken(ref)
				i = i + 2
			} else {
				value = token
			}
		} else {
			value = token
		}

		entries = append(entries, tokens.DictEntry{Key: key.Data(), Value: value})

		// Skip def keyword
		if peekNext(allTokens, i) == tokens.Def {
			i++
		}

		key = nil
	}

	return entries, nil
}

func peekNext(allTokens []tokens.Token, index int) tokens.Token {
	if len(allTokens)-1 < index {
		return nil
	}
	return allTokens[index]
}
