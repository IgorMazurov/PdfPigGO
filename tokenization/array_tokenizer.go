package tokenization

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ArrayTokenizer tokenizes PDF array objects starting with '['.
type ArrayTokenizer struct {
	usePdfDocEncoding bool
	stackDepthGuard   *core.StackDepthGuard
	useLenientParsing bool
}

var _ Tokenizer = (*ArrayTokenizer)(nil)

// NewArrayTokenizer creates a new ArrayTokenizer.
func NewArrayTokenizer(usePdfDocEncoding bool, stackDepthGuard *core.StackDepthGuard, useLenientParsing bool) *ArrayTokenizer {
	return &ArrayTokenizer{
		usePdfDocEncoding: usePdfDocEncoding,
		stackDepthGuard:   stackDepthGuard,
		useLenientParsing: useLenientParsing,
	}
}

// ReadsNextByte returns false because this tokenizer does not read past the closing ']'.
func (t *ArrayTokenizer) ReadsNextByte() bool {
	return false
}

// Tokenize attempts to tokenize a PDF array starting with '['.
// Returns the ArrayToken and true on success, or nil and false if currentByte is not '['.
func (t *ArrayTokenizer) Tokenize(currentByte byte, input core.InputBytes) (tokens.Token, bool) {
	if currentByte != '[' {
		return nil, false
	}

	scanner := NewCoreTokenScanner(input, t.usePdfDocEncoding, t.stackDepthGuard, ScannerScopeArray, nil, t.useLenientParsing, false)

	var rawTokens []tokens.Token
	var previousToken tokens.Token

	for !currentByteEndsCurrentArray(input, previousToken) && scanner.Advance() {
		previousToken = scanner.Current()

		if _, ok := previousToken.(*tokens.CommentToken); ok {
			continue
		}

		rawTokens = append(rawTokens, previousToken)
	}

	contents := assembleIndirectReferences(rawTokens)

	return tokens.NewArrayToken(contents), true
}

// assembleIndirectReferences combines "num gen R" token sequences into IndirectReferenceTokens.
func assembleIndirectReferences(tokensList []tokens.Token) []tokens.Token {
	result := make([]tokens.Token, 0, len(tokensList))

	for i := 0; i < len(tokensList); i++ {
		token := tokensList[i]

		if num, ok := token.(*tokens.NumericToken); ok {
			if i+2 < len(tokensList) {
				gen := tokensList[i+1]
				r := tokensList[i+2]

				if genNum, okGen := gen.(*tokens.NumericToken); okGen && r == tokens.OpR {
					ref, _ := core.NewIndirectReference(num.LongVal(), genNum.IntVal())
					result = append(result, tokens.NewIndirectReferenceToken(ref))
					i += 2
					continue
				}
			}
		}

		result = append(result, token)
	}

	return result
}

func currentByteEndsCurrentArray(input core.InputBytes, previousToken tokens.Token) bool {
	if input.CurrentByte() != ']' {
		return false
	}
	if _, ok := previousToken.(*tokens.ArrayToken); ok {
		return false
	}
	return true
}
