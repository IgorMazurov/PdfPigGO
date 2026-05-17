package content

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// GetHyperlinksFromPage extracts all hyperlinks from the given page by inspecting link annotations
// with /URI actions and collecting the letters that fall within each link's bounding rectangle.
func GetHyperlinksFromPage(
	pageLetters []*Letter,
	scanner tokenization.PdfTokenScanner,
	annotations []LinkAnnotationIface,
) []*Hyperlink {
	result := make([]*Hyperlink, 0)

	for _, annotation := range annotations {
		if !isLinkAnnotation(annotation) {
			continue
		}

		actionDict := resolveActionDictionary(annotation.AnnotationDictionary(), scanner)
		if actionDict == nil {
			continue
		}

		actionType := resolveNameToken(actionDict, tokens.S, scanner)
		if actionType == nil || actionType.Data() != "URI" {
			continue
		}

		uri := resolveUri(actionDict, scanner)
		if uri == "" {
			continue
		}

		bounds := annotation.Rectangle()

		tolerantBounds := core.NewPdfRectangle(
			bounds.BottomLeft.Translate(-0.5, -0.5),
			bounds.TopRight.Translate(0.5, 0.5),
		)

		linkLetters := collectLetters(pageLetters, tolerantBounds)

		presentationText := buildPresentationText(linkLetters)

		result = append(result, NewHyperlink(bounds, linkLetters, presentationText, uri, annotation))
	}

	return result
}

// isLinkAnnotation checks whether the given annotation represents a /Link annotation.
func isLinkAnnotation(annotation LinkAnnotationIface) bool {
	t := annotation.Type()
	switch v := t.(type) {
	case AnnotationType:
		return v == AnnotationTypeLink
	case string:
		return v == "Link" || v == "/Link"
	case fmt.Stringer:
		s := v.String()
		return s == "Link" || s == "1"
	default:
		return false
	}
}

// resolveActionDictionary resolves the /A entry from the annotation dictionary to a DictionaryToken.
func resolveActionDictionary(dict tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	dt, ok := dict.(*tokens.DictionaryToken)
	if !ok {
		return nil
	}

	token, ok := dt.TryGet(tokens.A)
	if !ok {
		return nil
	}

	resolved := unwrapIndirect(token, scanner)
	actionDict, ok := resolved.(*tokens.DictionaryToken)
	if !ok {
		return nil
	}

	return actionDict
}

// resolveNameToken resolves a named entry from a dictionary and returns it as a NameToken, or nil.
func resolveNameToken(dict *tokens.DictionaryToken, key *tokens.NameToken, scanner tokenization.PdfTokenScanner) *tokens.NameToken {
	token, ok := dict.TryGet(key)
	if !ok {
		return nil
	}

	resolved := unwrapIndirect(token, scanner)
	name, ok := resolved.(*tokens.NameToken)
	if !ok {
		return nil
	}

	return name
}

// resolveUri resolves the /URI entry from an action dictionary and returns its string value.
func resolveUri(dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) string {
	token, ok := dict.TryGet(tokens.Uri)
	if !ok {
		return ""
	}

	resolved := unwrapIndirect(token, scanner)

	switch t := resolved.(type) {
	case *tokens.StringToken:
		return t.Data()
	case *tokens.HexToken:
		return t.Data()
	}

	return ""
}

// collectLetters gathers all letters whose location falls within the given rectangle.
func collectLetters(letters []*Letter, bounds core.PdfRectangle) []*Letter {
	var result []*Letter
	for _, letter := range letters {
		if bounds.Contains(letter.Location(), true) {
			result = append(result, letter)
		}
	}
	return result
}

// buildPresentationText builds a space-separated string from the words extracted from the given letters.
func buildPresentationText(letters []*Letter) string {
	extractor := DefaultWordExtractor{}
	words := extractor.GetWords(letters)

	parts := make([]string, 0, len(words))
	for _, word := range words {
		parts = append(parts, word.Text)
	}

	return strings.Join(parts, " ")
}

// unwrapIndirect resolves an IndirectReferenceToken through the scanner.
func unwrapIndirect(token tokens.Token, scanner tokenization.PdfTokenScanner) tokens.Token {
	if indRef, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		ref := indRef.Data()
		obj := scanner.Get(ref)
		if obj != nil {
			return obj.Data()
		}
	}
	return token
}
