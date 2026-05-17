package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// GetFirstCharacter returns the first character code from the font dictionary's
// /FirstChar entry. Returns an error if the entry is missing or not numeric.
func GetFirstCharacter(dictionary *tokens.DictionaryToken) (int, error) {
	token, ok := dictionary.TryGet(tokens.FirstChar)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("No first character entry was found in the font dictionary for this TrueType font: %v.", dictionary))
	}

	number, ok := token.(*tokens.NumericToken)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("No first character entry was found in the font dictionary for this TrueType font: %v.", dictionary))
	}

	return number.IntVal(), nil
}

// GetLastCharacter returns the last character code from the font dictionary's
// /LastChar entry. Returns an error if the entry is missing or not numeric.
func GetLastCharacter(dictionary *tokens.DictionaryToken) (int, error) {
	token, ok := dictionary.TryGet(tokens.LastChar)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("No last character entry was found in the font dictionary for this TrueType font: %v.", dictionary))
	}

	number, ok := token.(*tokens.NumericToken)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("No last character entry was found in the font dictionary for this TrueType font: %v.", dictionary))
	}

	return number.IntVal(), nil
}

// GetWidths returns the array of character widths from the font dictionary's
// /Widths entry. Uses DirectObjectFinder to resolve indirect references.
func GetWidths(scanner tokenization.PdfTokenScanner, dictionary *tokens.DictionaryToken) ([]float64, error) {
	token, ok := dictionary.TryGet(tokens.Widths)
	if !ok {
		return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("No widths array found for the font: %v.", dictionary))
	}

	widthArray, err := parts.GetByToken[*tokens.ArrayToken](token, scanner)
	if err != nil {
		return nil, fonts.NewInvalidFontFormatExceptionWithInner("Could not resolve widths array", err)
	}

	result := make([]float64, widthArray.Length())
	for i := 0; i < widthArray.Length(); i++ {
		arrayElement := widthArray.Get(i)
		number, ok := arrayElement.(*tokens.NumericToken)
		if !ok {
			return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("Token which was not a number found in the widths array: %v.", arrayElement))
		}
		result[i] = number.DoubleVal()
	}

	return result, nil
}

// GetFontDescriptor resolves and generates the FontDescriptor from the font
// dictionary's /FontDescriptor entry. Uses DirectObjectFinder to resolve indirect
// references, then delegates to FontDescriptorFactory.Generate.
func GetFontDescriptor(scanner tokenization.PdfTokenScanner, dictionary *tokens.DictionaryToken) (*fonts.FontDescriptor, error) {
	token, ok := dictionary.TryGet(tokens.FontDescriptor)
	if !ok {
		return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("No font descriptor indirect reference found in the TrueType font: %v.", dictionary))
	}

	parsed, err := parts.GetByToken[*tokens.DictionaryToken](token, scanner)
	if err != nil {
		return nil, fonts.NewInvalidFontFormatExceptionWithInner("Could not resolve font descriptor", err)
	}

	return Generate(parsed, scanner)
}

// GetName resolves the font name from either the /BaseFont entry in the
// dictionary or falls back to the descriptor's FontName. Returns an error if
// neither source provides a name.
func GetName(scanner tokenization.PdfTokenScanner, dictionary *tokens.DictionaryToken, descriptor *fonts.FontDescriptor) (*tokens.NameToken, error) {
	nameBase, ok := dictionary.TryGet(tokens.BaseFont)
	if ok {
		name, err := parts.GetByToken[*tokens.NameToken](nameBase, scanner)
		if err != nil {
			return nil, fonts.NewInvalidFontFormatExceptionWithInner("Could not resolve base font name", err)
		}
		return name, nil
	}

	if descriptor.FontName != nil {
		return descriptor.FontName, nil
	}

	return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("Could not find a name for this font %v.", dictionary))
}
