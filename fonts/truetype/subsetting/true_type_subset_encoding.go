package subsetting

// TrueTypeSubsetEncoding represents a new encoding to create for the subsetted TrueType file.
type TrueTypeSubsetEncoding struct {
	characters []rune
}

// NewTrueTypeSubsetEncoding creates a new TrueTypeSubsetEncoding with the given characters.
// The characters are included in the subset in order where index is the character code.
func NewTrueTypeSubsetEncoding(characters []rune) *TrueTypeSubsetEncoding {
	return &TrueTypeSubsetEncoding{
		characters: characters,
	}
}

// Characters returns the characters to include in the subset in order where index is the character code.
func (e *TrueTypeSubsetEncoding) Characters() []rune {
	return e.characters
}
