package graphics

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// InlineImageBuilder accumulates inline image data between BI and EI operators.
type InlineImageBuilder struct {
	Properties map[*tokens.NameToken]tokens.Token
	Bytes      []byte
}

// NewInlineImageBuilder creates a new InlineImageBuilder.
func NewInlineImageBuilder() *InlineImageBuilder {
	return &InlineImageBuilder{
		Properties: make(map[*tokens.NameToken]tokens.Token),
	}
}

// CreateInlineImage builds an InlineImage from the accumulated data.
// Corresponds to C# InlineImageBuilder.CreateInlineImage().
func (b *InlineImageBuilder) CreateInlineImage(
	ctm core.TransformationMatrix,
	filterProvider content.LookupFilterProvider,
	scanner tokenization.PdfTokenScanner,
	renderingIntent graphiccore.RenderingIntent,
	resourceStore content.ResourceStore,
) *content.InlineImage {
	if b == nil || b.Properties == nil {
		return nil
	}

	bounds := ctm.TransformRect(core.NewPdfRectangle(core.Origin, core.NewPdfPoint(1, 1)))

	widthTok := b.getNumeric(tokens.Width, tokens.W, true)
	width := 0
	if widthTok != nil {
		width = widthTok.IntVal()
	}

	heightTok := b.getNumeric(tokens.Height, tokens.H, true)
	height := 0
	if heightTok != nil {
		height = heightTok.IntVal()
	}

	maskToken := b.getBoolean(tokens.ImageMask, tokens.Im, false)
	isMask := maskToken != nil && maskToken.Data()

	bitsPerComponent := 1
	if !isMask {
		if bpc := b.getNumeric(tokens.BitsPerComponent, tokens.Bpc, false); bpc != nil {
			bitsPerComponent = bpc.IntVal()
		}
	}

	imgDic := resolveDictionary(b.Properties, scanner)

	softMaskImage := b.resolveSoftMask(imgDic, scanner, filterProvider, renderingIntent, resourceStore)

	var colorSpaceName *tokens.NameToken
	if !isMask {
		csName := b.getName(tokens.ColorSpace, tokens.Cs, false)
		if csName != nil {
			colorSpaceName = csName
		} else {
			csArray := b.getArray(tokens.ColorSpace, tokens.Cs, true)
			if csArray == nil {
				b.logError("Missing required ColorSpace for inline image.")
				return nil
			}
			data := csArray.Data()
			if len(data) == 0 {
				b.logError("Empty ColorSpace array defined for inline image.")
				return nil
			}
			firstName, ok := data[0].(*tokens.NameToken)
			if !ok {
				b.logError(fmt.Sprintf("Invalid ColorSpace array defined for inline image: %v.", csArray))
				return nil
			}
			colorSpaceName = firstName
		}
	}

	csd := resourceStore.GetColorSpaceDetails(colorSpaceName, imgDic)
	details := &csd

	intentToken := b.getName(tokens.Intent, nil, false)
	intent := renderingIntent
	if intentToken != nil {
		intent = graphiccore.ParseRenderingIntent(intentToken.Data())
	}

	filterNames := b.extractFilterNames()

	decode := b.extractDecodeArray()

	interpolate := false
	if interp := b.getBoolean(tokens.Interpolate, tokens.I, false); interp != nil {
		interpolate = interp.Data()
	}

	return content.NewInlineImage(
		bounds,
		width,
		height,
		bitsPerComponent,
		isMask,
		intent,
		interpolate,
		decode,
		b.Bytes,
		toFilterProvider(filterProvider),
		filterNames,
		imgDic,
		details,
		softMaskImage,
	)
}

// resolveSoftMask processes the /SMask entry if present in the image dictionary.
func (b *InlineImageBuilder) resolveSoftMask(
	imgDic *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	filterProvider content.LookupFilterProvider,
	renderingIntent graphiccore.RenderingIntent,
	resourceStore content.ResourceStore,
) content.PdfImage {
	if imgDic == nil {
		return nil
	}

	smaskToken, ok := imgDic.TryGet(tokens.Smask)
	if !ok {
		return nil
	}

	stream, ok := resolveStream(smaskToken, scanner)
	if !ok || stream == nil {
		return nil
	}

	subtype, ok := stream.StreamDictionary.TryGet(tokens.Subtype)
	if !ok || !nameEquals(subtype, tokens.Image) {
		b.logError("The SMask dictionary does not contain a 'Subtype' entry, or its value is not 'Image'.")
		return nil
	}

	colorSpace, ok := stream.StreamDictionary.TryGet(tokens.ColorSpace)
	if !ok || !nameEquals(colorSpace, tokens.Devicegray) {
		b.logError("The SMask dictionary does not contain a 'ColorSpace' entry, or its value is not 'DeviceGray'.")
		return nil
	}

	if stream.StreamDictionary.ContainsKey(tokens.Mask) || stream.StreamDictionary.ContainsKey(tokens.Smask) {
		b.logError("The SMask dictionary contains a 'Mask' or 'Smask' entry.")
		return nil
	}

	lp, ok := filterProvider.(filters.LookupFilterProvider)
	if !ok {
		return nil
	}

	softMaskRecord := xobjects.NewXObjectContentRecord(
		nil,
		stream,
		xobjects.Image,
		core.Identity,
		renderingIntent,
		colors.DeviceGrayColorSpaceDetails,
	)

	img, err := content.ReadImage(softMaskRecord, scanner, lp, resourceStore)
	if err != nil {
		b.logError(fmt.Sprintf("Failed to read soft mask image: %v", err))
		return nil
	}

	return img
}

// extractFilterNames extracts filter names from /Filter or /F property entry.
func (b *InlineImageBuilder) extractFilterNames() []*tokens.NameToken {
	var result []*tokens.NameToken

	filterName := b.getName(tokens.Filter, tokens.F, false)
	if filterName != nil {
		result = append(result, filterName)
		return result
	}

	filterArray := b.getArray(tokens.Filter, tokens.F, false)
	if filterArray == nil {
		return result
	}

	for _, item := range filterArray.Data() {
		if name, ok := item.(*tokens.NameToken); ok {
			result = append(result, name)
		}
	}

	return result
}

// extractDecodeArray parses the /Decode or /D array property into []float64.
func (b *InlineImageBuilder) extractDecodeArray() []float64 {
	decodeRaw := b.getArray(tokens.Decode, tokens.D, false)
	if decodeRaw == nil {
		return []float64{}
	}

	var result []float64
	for _, item := range decodeRaw.Data() {
		if numeric, ok := item.(*tokens.NumericToken); ok {
			result = append(result, numeric.DoubleVal())
		}
	}

	return result
}

// getNumeric retrieves a NumericToken from Properties by trying two key names.
func (b *InlineImageBuilder) getNumeric(name1, name2 *tokens.NameToken, required bool) *tokens.NumericToken {
	if name1 != nil {
		if val, ok := b.Properties[name1]; ok {
			if t, ok := val.(*tokens.NumericToken); ok {
				return t
			}
		}
	}
	if name2 != nil {
		if val, ok := b.Properties[name2]; ok {
			if t, ok := val.(*tokens.NumericToken); ok {
				return t
			}
		}
	}
	if required {
		b.logError(fmt.Sprintf("Inline image dictionary missing required entry %v/%v.", name1, name2))
	}
	return nil
}

// getBoolean retrieves a BooleanToken from Properties by trying two key names.
func (b *InlineImageBuilder) getBoolean(name1, name2 *tokens.NameToken, required bool) *tokens.BooleanToken {
	if name1 != nil {
		if val, ok := b.Properties[name1]; ok {
			if t, ok := val.(*tokens.BooleanToken); ok {
				return t
			}
		}
	}
	if name2 != nil {
		if val, ok := b.Properties[name2]; ok {
			if t, ok := val.(*tokens.BooleanToken); ok {
				return t
			}
		}
	}
	if required {
		b.logError(fmt.Sprintf("Inline image dictionary missing required entry %v/%v.", name1, name2))
	}
	return nil
}

// getName retrieves a NameToken from Properties by trying two key names.
func (b *InlineImageBuilder) getName(name1, name2 *tokens.NameToken, required bool) *tokens.NameToken {
	if name1 != nil {
		if val, ok := b.Properties[name1]; ok {
			if t, ok := val.(*tokens.NameToken); ok {
				return t
			}
		}
	}
	if name2 != nil {
		if val, ok := b.Properties[name2]; ok {
			if t, ok := val.(*tokens.NameToken); ok {
				return t
			}
		}
	}
	if required {
		b.logError(fmt.Sprintf("Inline image dictionary missing required entry %v/%v.", name1, name2))
	}
	return nil
}

// getArray retrieves an ArrayToken from Properties by trying two key names.
func (b *InlineImageBuilder) getArray(name1, name2 *tokens.NameToken, required bool) *tokens.ArrayToken {
	if name1 != nil {
		if val, ok := b.Properties[name1]; ok {
			if t, ok := val.(*tokens.ArrayToken); ok {
				return t
			}
		}
	}
	if name2 != nil {
		if val, ok := b.Properties[name2]; ok {
			if t, ok := val.(*tokens.ArrayToken); ok {
				return t
			}
		}
	}
	if required {
		b.logError(fmt.Sprintf("Inline image dictionary missing required entry %v/%v.", name1, name2))
	}
	return nil
}

func (b *InlineImageBuilder) logError(msg string) {
	// Inline images don't have direct access to logger; errors surface via nil returns.
	_ = msg
}

// resolveDictionary converts the properties map into a DictionaryToken and resolves indirect references.
func resolveDictionary(properties map[*tokens.NameToken]tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	if len(properties) == 0 {
		return nil
	}

	dictData := make(map[string]tokens.Token, len(properties))
	for k, v := range properties {
		dictData[k.Data()] = v
	}

	dict, err := tokens.WithMap(dictData)
	if err != nil {
		return nil
	}

	resolved := resolveDictionaryToken(dict, scanner)
	if resolved != nil {
		return resolved
	}

	return dict
}

// resolveDictionaryToken follows an indirect reference if needed.
func resolveDictionaryToken(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	if dict, ok := token.(*tokens.DictionaryToken); ok {
		return dict
	}
	if indRef, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		obj := scanner.Get(indRef.Data())
		if obj != nil && obj.Data() != nil {
			if d, isArray := obj.Data().(*tokens.DictionaryToken); isArray {
				return d
			}
		}
	}
	return nil
}

// resolveStream resolves a token to a StreamToken, following indirect references if necessary.
func resolveStream(token tokens.Token, scanner tokenization.PdfTokenScanner) (*tokens.StreamToken, bool) {
	if stream, ok := token.(*tokens.StreamToken); ok {
		return stream, true
	}
	if ref, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		obj := scanner.Get(ref.Data())
		if obj != nil && obj.Data() != nil {
			if stream, ok := obj.Data().(*tokens.StreamToken); ok {
				return stream, true
			}
		}
	}
	return nil, false
}

// nameEquals checks if a token is a NameToken equal to the expected name.
func nameEquals(token tokens.Token, expected *tokens.NameToken) bool {
	if name, ok := token.(*tokens.NameToken); ok {
		return name.Data() == expected.Data()
	}
	return false
}

// toFilterProvider casts content.LookupFilterProvider to filters.FilterProvider when possible.
func toFilterProvider(fp content.LookupFilterProvider) filters.FilterProvider {
	if fp == nil {
		return nil
	}
	if provider, ok := fp.(filters.FilterProvider); ok {
		return provider
	}
	return nil
}
