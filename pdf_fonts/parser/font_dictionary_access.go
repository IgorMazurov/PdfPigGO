// Package parser provides types for parsing PDF font structures.
package parser

import (
	"errors"
	"fmt"
	"math"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

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

	descriptor, err := FontDescriptorFactoryGenerate(parsed, scanner)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
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

	if descriptor != nil && descriptor.FontName != nil {
		return descriptor.FontName, nil
	}

	return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("Could not find a name for this font %v.", dictionary))
}

// GetFirstCharacter returns the first character code from the font dictionary's
// /FirstChar entry. Returns an error if the entry is missing or not numeric.
func GetFirstCharacter(dictionary *tokens.DictionaryToken) (int, error) {
	token, ok := dictionary.TryGet(tokens.FirstChar)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("No first character entry was found in the font dictionary: %v.", dictionary))
	}

	number, ok := token.(*tokens.NumericToken)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("First character entry is not numeric: %v.", dictionary))
	}

	return number.IntVal(), nil
}

// GetLastCharacter returns the last character code from the font dictionary's
// /LastChar entry. Returns an error if the entry is missing or not numeric.
func GetLastCharacter(dictionary *tokens.DictionaryToken) (int, error) {
	token, ok := dictionary.TryGet(tokens.LastChar)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("No last character entry was found in the font dictionary: %v.", dictionary))
	}

	number, ok := token.(*tokens.NumericToken)
	if !ok {
		return 0, fonts.NewInvalidFontFormatException(fmt.Sprintf("Last character entry is not numeric: %v.", dictionary))
	}

	return number.IntVal(), nil
}

// FontDescriptorFactoryGenerate creates a FontDescriptor from the given dictionary token and scanner.
func FontDescriptorFactoryGenerate(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) (*fonts.FontDescriptor, error) {
	if dictionary == nil {
		return nil, errors.New("dictionary cannot be null")
	}

	name := fdgGetFontName(dictionary, scanner)
	family := fdgGetFontFamily(dictionary)
	stretch := fdgGetFontStretch(dictionary)
	flags := fdgGetFlags(dictionary)
	bounding := fdgGetBoundingBox(dictionary, scanner)
	charSet := fdgGetCharSet(dictionary)
	fontFile, err := fdgGetFontFile(dictionary)
	if err != nil {
		return nil, err
	}

	builder := fonts.NewFontDescriptorBuilder(name, flags)
	builder.FontFamily = family
	builder.Stretch = stretch
	builder.FontWeight = fdgGetDoubleOrDefault(dictionary, tokens.FontWeight)
	builder.BoundingBox = bounding
	builder.ItalicAngle = fdgGetDoubleOrDefault(dictionary, tokens.ItalicAngle)
	builder.Ascent = fdgGetDoubleOrDefault(dictionary, tokens.Ascent)
	builder.Descent = fdgGetDoubleOrDefault(dictionary, tokens.Descent)
	builder.Leading = fdgGetDoubleOrDefault(dictionary, tokens.Leading)
	builder.CapHeight = math.Abs(fdgGetDoubleOrDefault(dictionary, tokens.CapHeight))
	builder.XHeight = math.Abs(fdgGetDoubleOrDefault(dictionary, tokens.Xheight))
	builder.StemVertical = fdgGetDoubleOrDefault(dictionary, tokens.StemV)
	builder.StemHorizontal = fdgGetDoubleOrDefault(dictionary, tokens.StemH)
	builder.AverageWidth = fdgGetDoubleOrDefault(dictionary, tokens.AvgWidth)
	builder.MaxWidth = fdgGetDoubleOrDefault(dictionary, tokens.MaxWidth)
	builder.MissingWidth = fdgGetDoubleOrDefault(dictionary, tokens.MissingWidth)
	builder.FontFile = fontFile
	builder.CharSet = charSet

	return builder.Build(), nil
}

func fdgGetDoubleOrDefault(dictionary *tokens.DictionaryToken, name *tokens.NameToken) float64 {
	token, ok := dictionary.TryGet(name)
	if !ok {
		return 0
	}

	number, ok := token.(*tokens.NumericToken)
	if !ok {
		return 0
	}

	return number.DoubleVal()
}

func fdgGetFontName(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) *tokens.NameToken {
	token, ok := dictionary.TryGet(tokens.FontName)
	if !ok {
		return tokens.Create("")
	}

	name, resolved := parts.TryGet[*tokens.NameToken](token, scanner)
	if !resolved {
		return tokens.Create("")
	}

	return name
}

func fdgGetFontFamily(dictionary *tokens.DictionaryToken) string {
	token, ok := dictionary.TryGet(tokens.FontFamily)
	if !ok {
		return ""
	}

	family, ok := token.(*tokens.StringToken)
	if !ok {
		return ""
	}

	return family.Data()
}

func fdgGetFontStretch(dictionary *tokens.DictionaryToken) fonts.FontStretch {
	token, ok := dictionary.TryGet(tokens.FontStretch)
	if !ok {
		return fonts.Normal
	}

	stretchName, ok := token.(*tokens.NameToken)
	if !ok {
		return fonts.Normal
	}

	result := fdgConvertToFontStretch(stretchName)
	if result == fonts.Unknown {
		return fonts.Normal
	}

	return result
}

func fdgConvertToFontStretch(name *tokens.NameToken) fonts.FontStretch {
	switch name.Data() {
	case "UltraCondensed":
		return fonts.UltraCondensed
	case "ExtraCondensed":
		return fonts.ExtraCondensed
	case "Condensed":
		return fonts.Condensed
	case "Normal":
		return fonts.Normal
	case "SemiExpanded":
		return fonts.SemiExpanded
	case "Expanded":
		return fonts.Expanded
	case "ExtraExpanded":
		return fonts.ExtraExpanded
	case "UltraExpanded":
		return fonts.UltraExpanded
	default:
		return fonts.Unknown
	}
}

func fdgGetFlags(dictionary *tokens.DictionaryToken) fonts.FontDescriptorFlags {
	token, ok := dictionary.TryGet(tokens.Flags)
	if !ok {
		return 0
	}

	number, ok := token.(*tokens.NumericToken)
	if !ok {
		return 0
	}

	val := number.IntVal()
	if val == -1 {
		return 0
	}

	return fonts.FontDescriptorFlags(val)
}

func fdgGetBoundingBox(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) core.PdfRectangle {
	token, ok := dictionary.TryGet(tokens.FontBbox)
	if !ok {
		return core.NewPdfRectangleFloat(0, 0, 0, 0)
	}

	boxArray, ok := token.(*tokens.ArrayToken)
	if !ok {
		return core.NewPdfRectangleFloat(0, 0, 0, 0)
	}

	rect := fdgArrayToRectangle(boxArray, scanner)
	if rect == nil {
		return core.NewPdfRectangleFloat(0, 0, 0, 0)
	}

	return *rect
}

func fdgGetCharSet(dictionary *tokens.DictionaryToken) string {
	token, ok := dictionary.TryGet(tokens.CharSet)
	if !ok {
		return ""
	}

	setName, ok := token.(*tokens.NameToken)
	if !ok {
		return ""
	}

	return setName.Data()
}

func fdgGetFontFile(dictionary *tokens.DictionaryToken) (*fonts.DescriptorFontFile, error) {
	token, ok := dictionary.TryGet(tokens.FontFile)
	if ok {
		obj, ok := token.(*tokens.IndirectReferenceToken)
		if !ok {
			return nil, errors.New("expected FontFile to be an indirect object reference")
		}
		return fonts.NewDescriptorFontFile(obj, fonts.Type1), nil
	}

	token, ok = dictionary.TryGet(tokens.FontFile2)
	if ok {
		obj, ok := token.(*tokens.IndirectReferenceToken)
		if !ok {
			return nil, errors.New("expected FontFile2 to be an indirect object reference")
		}
		return fonts.NewDescriptorFontFile(obj, fonts.TrueType), nil
	}

	token, ok = dictionary.TryGet(tokens.FontFile3)
	if ok {
		obj, ok := token.(*tokens.IndirectReferenceToken)
		if !ok {
			return nil, errors.New("expected FontFile3 to be an indirect object reference")
		}
		return fonts.NewDescriptorFontFile(obj, fonts.FromSubtype), nil
	}

	return nil, nil
}

func fdgArrayToRectangle(boxArray *tokens.ArrayToken, scanner tokenization.PdfTokenScanner) *core.PdfRectangle {
	if boxArray.Length() != 4 {
		return nil
	}

	var coords [4]float64
	for i := 0; i < 4; i++ {
		item := boxArray.Get(i)
		switch t := item.(type) {
		case *tokens.NumericToken:
			coords[i] = t.DoubleVal()
		case *tokens.IndirectReferenceToken:
			resolved := scanner.Get(t.Data())
			if resolved == nil {
				return nil
			}
			numericTok, ok := resolved.Data().(*tokens.NumericToken)
			if !ok {
				return nil
			}
			coords[i] = numericTok.DoubleVal()
		default:
			return nil
		}
	}

	rect := core.NewPdfRectangleFloat(coords[0], coords[1], coords[2], coords[3])
	return &rect
}
