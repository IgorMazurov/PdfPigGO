package parser

import (
	"errors"
	"math"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// Generate creates a FontDescriptor from the given dictionary token and scanner.
func Generate(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) (*fonts.FontDescriptor, error) {
	if dictionary == nil {
		return nil, errors.New("dictionary cannot be null")
	}

	name := getFontName(dictionary, scanner)
	family := getFontFamily(dictionary)
	stretch := getFontStretch(dictionary)
	flags := getFlags(dictionary)
	bounding := getBoundingBox(dictionary, scanner)
	charSet := getCharSet(dictionary)
	fontFile, err := getFontFile(dictionary)
	if err != nil {
		return nil, err
	}

	builder := fonts.NewFontDescriptorBuilder(name, flags)
	builder.FontFamily = family
	builder.Stretch = stretch
	builder.FontWeight = getDoubleOrDefault(dictionary, tokens.FontWeight)
	builder.BoundingBox = bounding
	builder.ItalicAngle = getDoubleOrDefault(dictionary, tokens.ItalicAngle)
	builder.Ascent = getDoubleOrDefault(dictionary, tokens.Ascent)
	builder.Descent = getDoubleOrDefault(dictionary, tokens.Descent)
	builder.Leading = getDoubleOrDefault(dictionary, tokens.Leading)
	builder.CapHeight = math.Abs(getDoubleOrDefault(dictionary, tokens.CapHeight))
	builder.XHeight = math.Abs(getDoubleOrDefault(dictionary, tokens.Xheight))
	builder.StemVertical = getDoubleOrDefault(dictionary, tokens.StemV)
	builder.StemHorizontal = getDoubleOrDefault(dictionary, tokens.StemH)
	builder.AverageWidth = getDoubleOrDefault(dictionary, tokens.AvgWidth)
	builder.MaxWidth = getDoubleOrDefault(dictionary, tokens.MaxWidth)
	builder.MissingWidth = getDoubleOrDefault(dictionary, tokens.MissingWidth)
	builder.FontFile = fontFile
	builder.CharSet = charSet

	return builder.Build(), nil
}

func getDoubleOrDefault(dictionary *tokens.DictionaryToken, name *tokens.NameToken) float64 {
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

func getFontName(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) *tokens.NameToken {
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

func getFontFamily(dictionary *tokens.DictionaryToken) string {
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

func getFontStretch(dictionary *tokens.DictionaryToken) fonts.FontStretch {
	token, ok := dictionary.TryGet(tokens.FontStretch)
	if !ok {
		return fonts.Normal
	}

	stretchName, ok := token.(*tokens.NameToken)
	if !ok {
		return fonts.Normal
	}

	result := pdffonts.ConvertToFontStretch(stretchName)
	if result == fonts.Unknown {
		return fonts.Normal
	}

	return result
}

func getFlags(dictionary *tokens.DictionaryToken) fonts.FontDescriptorFlags {
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

func getBoundingBox(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) core.PdfRectangle {
	token, ok := dictionary.TryGet(tokens.FontBbox)
	if !ok {
		return core.NewPdfRectangleFloat(0, 0, 0, 0)
	}

	boxArray, ok := token.(*tokens.ArrayToken)
	if !ok {
		return core.NewPdfRectangleFloat(0, 0, 0, 0)
	}

	rect := arrayToRectangle(boxArray, scanner)
	if rect == nil {
		return core.NewPdfRectangleFloat(0, 0, 0, 0)
	}

	return *rect
}

func getCharSet(dictionary *tokens.DictionaryToken) string {
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

func getFontFile(dictionary *tokens.DictionaryToken) (*fonts.DescriptorFontFile, error) {
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

func arrayToRectangle(boxArray *tokens.ArrayToken, scanner tokenization.PdfTokenScanner) *core.PdfRectangle {
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
