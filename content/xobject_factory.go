package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	pdfimages "github.com/uglytoad/pdfpig/go/images"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// ReadImage reads an XObject image from the given content record.
func ReadImage(
	xObject *xobjects.XObjectContentRecord,
	pdfScanner tokenization.PdfTokenScanner,
	filterProvider filters.FilterProvider,
	resourceStore ResourceStore,
) (*XObjectImage, error) {
	if xObject == nil {
		return nil, core.NewPdfDocumentFormatException("xObject cannot be null")
	}

	if xObject.Type != xobjects.Image {
		return nil, fmt.Errorf("cannot create an image from an XObject with type: %v", xObject.Type)
	}

	dictionary := xObject.Stream.StreamDictionary

	bounds := xObject.AppliedTransformation.TransformRect(
		core.NewPdfRectangle(core.Origin, core.NewPdfPoint(1, 1)),
	)

	width, err := dictGetInt(dictionary, tokens.Width, pdfScanner)
	if err != nil {
		return nil, err
	}

	height, err := dictGetInt(dictionary, tokens.Height, pdfScanner)
	if err != nil {
		return nil, err
	}

	isImageMask := false
	if maskToken, ok := dictionary.TryGet(tokens.ImageMask); ok {
		maskToken = resolveTokenIndirect(maskToken, pdfScanner)
		if boolToken, ok := maskToken.(*tokens.BooleanToken); ok {
			isImageMask = boolToken.Data()
		}
	}

	var softMaskImage PdfImage

	if smaskToken, ok := tryGetStream(dictionary, tokens.Smask, pdfScanner); ok {
		subtype, ok := smaskToken.StreamDictionary.TryGet(tokens.Subtype)
		if !ok || !nameEquals(subtype, tokens.Image) {
			return nil, fmt.Errorf("the SMask dictionary does not contain a 'Subtype' entry, or its value is not 'Image'")
		}

		colorSpace, ok := smaskToken.StreamDictionary.TryGet(tokens.ColorSpace)
		if !ok || !nameEquals(colorSpace, tokens.Devicegray) {
			return nil, fmt.Errorf("the SMask dictionary does not contain a 'ColorSpace' entry, or its value is not 'DeviceGray'")
		}

		if smaskToken.StreamDictionary.ContainsKey(tokens.Mask) || smaskToken.StreamDictionary.ContainsKey(tokens.Smask) {
			return nil, fmt.Errorf("the SMask dictionary contains a 'Mask' or 'Smask' entry")
		}

		softMaskRecord := xobjects.NewXObjectContentRecord(
			nil,
			smaskToken,
			xobjects.Image,
			core.Identity,
			xObject.DefaultRenderingIntent,
			colors.DeviceGrayColorSpaceDetails,
		)

		softMaskImage, err = ReadImage(softMaskRecord, pdfScanner, filterProvider, resourceStore)
		if err != nil {
			return nil, err
		}
	} else if maskStream, ok := tryGetStream(dictionary, tokens.Mask, pdfScanner); ok {
		maskRecord := xobjects.NewXObjectContentRecord(
			nil,
			maskStream,
			xobjects.Image,
			core.Identity,
			xObject.DefaultRenderingIntent,
			nil,
		)

		softMaskImage, err = ReadImage(maskRecord, pdfScanner, filterProvider, resourceStore)
		if err != nil {
			return nil, err
		}
	}

	isJpxDecode := false
	if filterToken, ok := dictionary.TryGet(tokens.Filter); ok {
		filterToken = resolveTokenIndirect(filterToken, pdfScanner)
		if nameToken, ok := filterToken.(*tokens.NameToken); ok {
			isJpxDecode = nameToken.Data() == tokens.JpxDecode.Data()
		}
	}

	var bitsPerComponent int

	if isImageMask {
		bitsPerComponent = 1
	} else if isJpxDecode {
		if bpcToken, ok := dictionary.TryGet(tokens.BitsPerComponent); ok {
			bpcToken = resolveTokenIndirect(bpcToken, pdfScanner)
			if numeric, ok := bpcToken.(*tokens.NumericToken); ok {
				bitsPerComponent = numeric.IntVal()
			}
		} else {
			bpcByte, err := pdfimages.GetBitsPerComponent(xObject.Stream.Data())
			if err != nil {
				return nil, fmt.Errorf("failed to get bits per component from JPX data: %w", err)
			}
			bitsPerComponent = int(bpcByte)
		}
	} else {
		bpcToken, ok := dictionary.TryGet(tokens.BitsPerComponent)
		if !ok {
			return nil, core.NewPdfDocumentFormatException(
				fmt.Sprintf("no bits per component defined for image: %v", dictionary))
		}
		bpcToken = resolveTokenIndirect(bpcToken, pdfScanner)
		if numeric, ok := bpcToken.(*tokens.NumericToken); ok {
			bitsPerComponent = numeric.IntVal()
		} else {
			return nil, core.NewPdfDocumentFormatException(
				fmt.Sprintf("bits per component is not a number: %v", bpcToken))
		}
	}

	intent := xObject.DefaultRenderingIntent
	if intentToken, ok := dictionary.TryGet(tokens.Intent); ok {
		intentToken = resolveTokenIndirect(intentToken, pdfScanner)
		if name, ok := intentToken.(*tokens.NameToken); ok {
			intent = graphiccore.ParseRenderingIntent(name.Data())
		}
	}

	interpolate := false
	if interpToken, ok := dictionary.TryGet(tokens.Interpolate); ok {
		interpToken = resolveTokenIndirect(interpToken, pdfScanner)
		if boolTok, ok := interpToken.(*tokens.BooleanToken); ok {
			interpolate = boolTok.Data()
		}
	}

	fl, err := getFiltersWithScanner(filterProvider, dictionary, pdfScanner)
	if err != nil {
		return nil, err
	}

	supportsFilters := true
	for _, f := range fl {
		if !f.IsSupported() {
			supportsFilters = false
			break
		}
	}

	var varDecode []float64
	if decodeArrayToken, ok := dictionary.TryGet(tokens.Decode); ok {
		decodeArrayToken = resolveTokenIndirect(decodeArrayToken, pdfScanner)
		if arr, ok := decodeArrayToken.(*tokens.ArrayToken); ok {
			data := arr.Data()
			varDecode = make([]float64, 0, len(data))
			for _, token := range data {
				if numeric, ok := token.(*tokens.NumericToken); ok {
					varDecode = append(varDecode, numeric.DoubleVal())
				}
			}
		}
	}
	if varDecode == nil {
		varDecode = []float64{}
	}

	var details colors.ColorSpaceDetails

	if isImageMask {
		details = resourceStore.GetColorSpaceDetails(nil, dictionary)
	} else if csToken, ok := dictionary.TryGet(tokens.ColorSpace); ok {
		csToken = resolveTokenIndirect(csToken, pdfScanner)
		if name, ok := csToken.(*tokens.NameToken); ok {
			details = resourceStore.GetColorSpaceDetails(name, dictionary)
		} else if arr, ok := csToken.(*tokens.ArrayToken); ok {
			data := arr.Data()
			if len(data) > 0 {
				first := data[0]
				if ref, ok := first.(*tokens.IndirectReferenceToken); ok && pdfScanner != nil {
					obj := pdfScanner.Get(ref.Data())
					if obj != nil && obj.Data() != nil {
						first = obj.Data()
					}
				}
				if first, ok := first.(*tokens.NameToken); ok {
					details = resourceStore.GetColorSpaceDetails(first, dictionary.With(tokens.ColorSpace, arr))
				}
			}
		}
	} else if !isJpxDecode {
		details = xObject.DefaultColorSpace
	}

	var streamForDecode *tokens.StreamToken
	if supportsFilters {
		streamForDecode = xObject.Stream
	}

	img := NewXObjectImage(
		bounds,
		width,
		height,
		bitsPerComponent,
		isJpxDecode,
		isImageMask,
		intent,
		interpolate,
		varDecode,
		dictionary,
		xObject.Stream.Data(),
		details,
		softMaskImage,
		filterProvider,
		streamForDecode,
		supportsFilters,
	)

	img.SetPdfScanner(pdfScanner)
	return img, nil
}

func getFiltersWithScanner(fp filters.FilterProvider, dict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) ([]filters.Filter, error) {
	if lp, ok := fp.(filters.LookupFilterProvider); ok && scanner != nil {
		return lp.GetFiltersWithScanner(dict, scanner)
	}
	return fp.GetFilters(dict)
}

func dictGetInt(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (int, error) {
	token, ok := dict.TryGet(name)
	if !ok {
		return 0, core.NewPdfDocumentFormatException(
			fmt.Sprintf("dictionary does not contain key %v", name))
	}
	token = resolveTokenIndirect(token, scanner)
	if numeric, ok := token.(*tokens.NumericToken); ok {
		return numeric.IntVal(), nil
	}
	return 0, core.NewPdfDocumentFormatException(
		fmt.Sprintf("value for key %v is not a number: %v", name, token))
}

func resolveTokenIndirect(token tokens.Token, scanner tokenization.PdfTokenScanner) tokens.Token {
	if ref, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		obj := scanner.Get(ref.Data())
		if obj != nil && obj.Data() != nil {
			return obj.Data()
		}
	}
	return token
}

func tryGetStream(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (*tokens.StreamToken, bool) {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil, false
	}

	if stream, ok := token.(*tokens.StreamToken); ok {
		return stream, true
	}

	if ref, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		obj := scanner.Get(ref.Data())
		if obj == nil || obj.Data() == nil {
			return nil, false
		}
		if stream, ok := obj.Data().(*tokens.StreamToken); ok {
			return stream, true
		}
	}

	return nil, false
}

func nameEquals(token tokens.Token, expected *tokens.NameToken) bool {
	if name, ok := token.(*tokens.NameToken); ok {
		return name.Data() == expected.Data()
	}
	return false
}
