package writer

import "github.com/uglytoad/pdfpig/go/tokens"

// Obey adds the required PDF/A-1a structure tree root and mark info entries to
// the given catalog dictionary so that it satisfies PDF/A-1a accessibility rules.
func Obey(catalog map[*tokens.NameToken]tokens.Token) {
	catalog[tokens.StructTreeRoot] = generateStructTree()

	markInfo, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Marked: tokens.True,
	})
	catalog[tokens.MarkInfo] = markInfo
}

// generateStructTree creates a minimal structure tree root dictionary with its
// Type set to /StructTreeRoot.
func generateStructTree() *tokens.DictionaryToken {
	dt, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Type: tokens.StructTreeRoot,
	})
	return dt
}
