package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

// ParsingOptions holds configuration options for PDF parsing.
type ParsingOptions struct {
	Logger          logging.Log
	UseLenientParsing bool
	SkipMissingFonts  bool
	ClipPaths         bool
	Password          string
	Passwords         []string
	FilterProvider    any
	MaxStackDepth     int
}

// LookupFilterProvider provides filter decoding for compressed PDF streams.
type LookupFilterProvider interface {
	DecodeStream(stream *tokens.StreamToken, scanner tokenization.PdfTokenScanner) []byte
}

// PageContentParser parses raw page content bytes into graphics operations.
type PageContentParser interface {
	Parse(pageNumber int, inputBytes core.InputBytes, logger logging.Log) []GraphicsStateOperation
}

// GraphicsStateOperation represents a single graphics operation in a PDF content stream.
type GraphicsStateOperation interface {
	IsGraphicsOp()
}

// PageTreeMembers holds inherited properties from ancestor page tree nodes.
type PageTreeMembers struct {
	ParentResources []tokens.Token
	Rotation        *int
	MediaBoxValue   *MediaBox
	CropBoxValue    *CropBox
}

// GetCropBox returns the inherited crop box, or nil if none is defined.
func (p *PageTreeMembers) GetCropBox() *CropBox {
	return p.CropBoxValue
}

// PageFactory defines the contract for creating page instances from PDF dictionaries.
type PageFactory[T any] interface {
	Create(number int, dictionary *tokens.DictionaryToken, pageTreeMembers *PageTreeMembers, namedDestinations *destinations.NamedDestinations) (T, error)
}

// ProcessPageHook is a function type for virtual dispatch of ProcessPage.
// External packages can set this field to override page construction behavior
// without relying on Go's embedding-based method resolution.
type ProcessPageHook[T any] func(
	pageNumber int,
	dictionary *tokens.DictionaryToken,
	namedDestinations *destinations.NamedDestinations,
	mediaBox *MediaBox,
	cropBox *CropBox,
	userSpaceUnit geometry.UserSpaceUnit,
	rotation PageRotationDegrees,
	initialMatrix core.TransformationMatrix,
	operations []GraphicsStateOperation,
) T

// BasePageFactory provides shared logic for constructing pages from their PDF dictionary representations.
// It handles resource loading, box extraction, content stream decoding, and operation parsing.
type BasePageFactory[T any] struct {
	ParsingOptions    *ParsingOptions
	PdfScanner        tokenization.PdfTokenScanner
	ResourceStore     ResourceStore
	FilterProvider    LookupFilterProvider
	PageContentParser PageContentParser
	ProcessPageHook   ProcessPageHook[T]
}

var _ PageFactory[any] = (*BasePageFactory[any])(nil)

// NewBasePageFactory creates a new BasePageFactory with the given dependencies.
func NewBasePageFactory[T any](
	pdfScanner tokenization.PdfTokenScanner,
	resourceStore ResourceStore,
	filterProvider LookupFilterProvider,
	pageContentParser PageContentParser,
	parsingOptions *ParsingOptions,
) *BasePageFactory[T] {
	return &BasePageFactory[T]{
		PdfScanner:        pdfScanner,
		ResourceStore:     resourceStore,
		FilterProvider:    filterProvider,
		PageContentParser: pageContentParser,
		ParsingOptions:    parsingOptions,
	}
}

// Create constructs a page from the given dictionary and context.
func (f *BasePageFactory[T]) Create(
	number int,
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *PageTreeMembers,
	namedDestinations *destinations.NamedDestinations,
) (T, error) {
	var zero T

	if dictionary == nil {
		return zero, fmt.Errorf("dictionary cannot be nil")
	}

	pgType, _ := util.GetNameOrDefault(dictionary, tokens.Type)

	if pgType != nil && pgType.Data() != tokens.Page.Data() {
		f.ParsingOptions.Logger.Error(fmt.Sprintf("Page %d had its type specified as %s rather than 'Page'.", number, pgType.Data()))
	}

	rotation := newRotationFromPtr(pageTreeMembers.Rotation)

	if rotateRaw, found := dictionary.TryGet(tokens.Rotate); found {
		if numericTok, ok := parts.TryGet[*tokens.NumericToken](rotateRaw, f.PdfScanner); ok {
			newRot, err := NewPageRotationDegrees(numericTok.IntVal())
			if err == nil {
				rotation = newRot
			} else {
				f.ParsingOptions.Logger.Error(fmt.Sprintf("Invalid rotation value for page %d: %v", number, err))
			}
		}
	}

	stackDepth := 0

	for len(pageTreeMembers.ParentResources) > 0 {
		resource := pageTreeMembers.ParentResources[0]
		pageTreeMembers.ParentResources = pageTreeMembers.ParentResources[1:]

		resDict, ok := parts.TryGet[*tokens.DictionaryToken](resource, f.PdfScanner)
		if ok {
			f.ResourceStore.LoadResourceDictionary(resDict)
		}
		stackDepth++
	}

	if resourcesRaw, found := dictionary.TryGet(tokens.Resources); found {
		resDict, ok := parts.TryGet[*tokens.DictionaryToken](resourcesRaw, f.PdfScanner)
		if ok {
			f.ResourceStore.LoadResourceDictionary(resDict)
			stackDepth++
		}
	}

	userSpaceUnit := GetUserSpaceUnits(dictionary)

	mediaBox := GetMediaBox(number, dictionary, pageTreeMembers, f.ParsingOptions.Logger)
	cropBox := GetCropBox(dictionary, pageTreeMembers, mediaBox, f.PdfScanner, f.ParsingOptions.Logger)

	initialMatrix := GetInitialMatrix(userSpaceUnit, mediaBox, cropBox, rotation, f.ParsingOptions.Logger)

	ApplyTransformNormalise(initialMatrix, &mediaBox, &cropBox)

	var page T

	contentsRaw, hasContents := dictionary.TryGet(tokens.Contents)

	if !hasContents {
		page = f.processPageInternal(number, dictionary, namedDestinations, mediaBox, cropBox, userSpaceUnit, rotation, initialMatrix, nil)
	} else {
		contentBytes, err := f.extractContentBytes(contentsRaw, number)
		if err != nil {
			return zero, err
		}

		page = f.processPageInternal(number, dictionary, namedDestinations, mediaBox, cropBox, userSpaceUnit, rotation, initialMatrix, contentBytes)
	}

	for i := 0; i < stackDepth; i++ {
		f.ResourceStore.UnloadResourceDictionary()
	}

	return page, nil
}

func (f *BasePageFactory[T]) extractContentBytes(contentsRaw tokens.Token, pageNumber int) ([]byte, error) {
	arrayTok, isArray := parts.TryGet[*tokens.ArrayToken](contentsRaw, f.PdfScanner)

	if isArray {
		data := arrayTok.Data()
		buf := make([]byte, 0, 1024*64)

		for i, item := range data {
			obj, ok := parts.TryGet[*tokens.IndirectReferenceToken](item, f.PdfScanner)
			if !ok {
				return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("the contents contained something which was not an indirect reference: %v.", item))
			}

			streamTok, err := parts.GetByToken[*tokens.StreamToken](obj, f.PdfScanner)
			if err != nil || streamTok == nil {
				return nil, fmt.Errorf("could not find the contents for object %v", obj)
			}

			decoded := f.FilterProvider.DecodeStream(streamTok, f.PdfScanner)
			buf = append(buf, decoded...)

			if i < len(data)-1 {
				buf = append(buf, '\n')
			}
		}

		return buf, nil
	}

	streamTok, err := parts.GetByToken[*tokens.StreamToken](contentsRaw, f.PdfScanner)
	if err != nil || streamTok == nil {
		return nil, fmt.Errorf("failed to parse the content for the page: %d", pageNumber)
	}

	return f.FilterProvider.DecodeStream(streamTok, f.PdfScanner), nil
}

func (f *BasePageFactory[T]) processPageInternal(
	pageNumber int,
	dictionary *tokens.DictionaryToken,
	namedDestinations *destinations.NamedDestinations,
	mediaBox *MediaBox,
	cropBox *CropBox,
	userSpaceUnit geometry.UserSpaceUnit,
	rotation PageRotationDegrees,
	initialMatrix core.TransformationMatrix,
	contentBytes []byte,
) T {
	var operations []GraphicsStateOperation

	if len(contentBytes) == 0 {
		operations = make([]GraphicsStateOperation, 0)
	} else {
		inputBytes := core.NewMemoryInputBytes(contentBytes)
		operations = f.PageContentParser.Parse(pageNumber, inputBytes, f.ParsingOptions.Logger)
	}

	if f.ProcessPageHook != nil {
		return f.ProcessPageHook(
			pageNumber,
			dictionary,
			namedDestinations,
			mediaBox,
			cropBox,
			userSpaceUnit,
			rotation,
			initialMatrix,
			operations,
		)
	}

	return f.ProcessPage(
		pageNumber,
		dictionary,
		namedDestinations,
		mediaBox,
		cropBox,
		userSpaceUnit,
		rotation,
		initialMatrix,
		operations,
	)
}

// ProcessPage is called by subclasses to handle the final page construction.
// Implementations should embed BasePageFactory[T] and override this method
// via function field injection or composition pattern. This base implementation
// returns a zero value and logs a warning if called directly.
func (f *BasePageFactory[T]) ProcessPage(
	pageNumber int,
	dictionary *tokens.DictionaryToken,
	namedDestinations *destinations.NamedDestinations,
	mediaBox *MediaBox,
	cropBox *CropBox,
	userSpaceUnit geometry.UserSpaceUnit,
	rotation PageRotationDegrees,
	initialMatrix core.TransformationMatrix,
	operations []GraphicsStateOperation,
) T {
	var zero T
	f.ParsingOptions.Logger.Error(fmt.Sprintf("ProcessPage was not overridden for page %d. Returning empty result.", pageNumber))
	return zero
}

// GetUserSpaceUnits extracts the user space unit from a page dictionary.
func GetUserSpaceUnits(dictionary *tokens.DictionaryToken) geometry.UserSpaceUnit {
	if userUnitRaw, found := dictionary.TryGet(tokens.UserUnit); found {
		if numericTok, ok := userUnitRaw.(*tokens.NumericToken); ok {
			unit, err := geometry.NewUserSpaceUnit(numericTok.IntVal())
			if err == nil {
				return unit
			}
		}
	}

	return geometry.Default
}

// GetCropBox extracts or computes the crop box for a page.
func GetCropBox(
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *PageTreeMembers,
	mediaBox *MediaBox,
	scanner tokenization.PdfTokenScanner,
	logger logging.Log,
) *CropBox {
	if cropBoxRaw, found := dictionary.TryGet(tokens.CropBox); found {
		cropArray, ok := parts.TryGet[*tokens.ArrayToken](cropBoxRaw, scanner)

		if ok {
			if cropArray.Length() != 4 {
				logger.Error(fmt.Sprintf("The CropBox was the wrong length in the dictionary: %v. Array had %d elements. Using MediaBox.", dictionary, cropArray.Length()))
				return NewCropBox(mediaBox.Bounds)
			}

			rect, err := util.ToRectangle(cropArray, scanner)
			if err == nil && rect != nil {
				return NewCropBox(*rect)
			}
		}
	}

	inherited := pageTreeMembers.GetCropBox()
	if inherited != nil {
		return inherited
	}

	return NewCropBox(mediaBox.Bounds)
}

// GetMediaBox extracts or computes the media box for a page.
func GetMediaBox(
	number int,
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *PageTreeMembers,
	logger logging.Log,
) *MediaBox {
	if mediaBoxRaw, found := dictionary.TryGet(tokens.MediaBox); found {
		mediaArray := resolveToArray(mediaBoxRaw, nil)

		if mediaArray != nil {
			data := mediaArray.Data()
			if len(data) != 4 {
				logger.Error(fmt.Sprintf("The MediaBox was the wrong length in the dictionary: %v. Array had %d elements. Defaulting to US Letter.", dictionary, len(data)))
				return MediaBoxUSLetter
			}

			rect := arrayToRectangle(data, nil)
			if rect != nil {
				return NewMediaBox(*rect)
			}
		}
	}

	if pageTreeMembers.MediaBoxValue != nil {
		return pageTreeMembers.MediaBoxValue
	}

	logger.Error(fmt.Sprintf("The MediaBox was missing for page %d. Using US Letter.", number))
	return MediaBoxUSLetter
}

// GetInitialMatrix computes the initial transformation matrix from user space units, media box, crop box, and rotation.
// Matches C# OperationContextHelper.GetInitialMatrix exactly.
func GetInitialMatrix(
	userSpaceUnit geometry.UserSpaceUnit,
	mediaBox *MediaBox,
	cropBox *CropBox,
	rotation PageRotationDegrees,
	logger logging.Log,
) core.TransformationMatrix {
	// Cater for scenario where the cropbox is larger than the mediabox.
	// If there is no intersection, fall back to the cropbox.
	viewBox := intersectRectangles(mediaBox.Bounds, cropBox.Bounds)
	if viewBox == nil {
		viewBox = &cropBox.Bounds
	}

	if rotation.Value == 0 &&
		viewBox.Left() == 0 &&
		viewBox.Bottom() == 0 &&
		userSpaceUnit.PointMultiples == 1 {
		return core.Identity
	}

	// Move points so that (0,0) is equal to the viewbox bottom left corner.
	t1 := core.GetTranslationMatrix(-viewBox.Left(), -viewBox.Bottom())

	if userSpaceUnit.PointMultiples != 1 {
		logger.Warn("User space unit other than 1 is not implemented")
	}

	if rotation.Value == 0 {
		return t1
	}

	// After rotating around the origin, our points will have negative x/y coordinates.
	// Fix this by translating them by a certain dx/dy after rotation based on the viewbox.
	var dx, dy float64
	switch rotation.Value {
	case 90:
		dx = 0
		dy = viewBox.Width
	case 180:
		dx = viewBox.Width
		dy = viewBox.Height
	case 270:
		dx = viewBox.Height
		dy = 0
	default:
		dx = 0
		dy = 0
	}

	// GetRotationMatrix uses counter clockwise angles, whereas our page rotation
	// is a clockwise angle, so flip the sign.
	r := core.GetRotationMatrix(-float64(rotation.Value))

	// Fix up negative coordinates after rotation
	t2 := core.GetTranslationMatrix(dx, dy)

	// Now get the final combined matrix T1 > R > T2
	return t1.Multiply(r.Multiply(t2))
}

func cosDegrees(degrees float64) float64 {
	switch int(degrees) {
	case 0:
		return 1.0
	case 90, 270:
		return 0.0
	case 180:
		return -1.0
	default:
		return 0.0
	}
}

func sinDegrees(degrees float64) float64 {
	switch int(degrees) {
	case 0, 180:
		return 0.0
	case 90:
		return 1.0
	case 270:
		return -1.0
	default:
		return 0.0
	}
}

// ApplyTransformNormalise applies a transformation matrix to the media box and crop box,
// then normalises them so rotation=0 and width/height match on-screen dimensions.
func ApplyTransformNormalise(
	transformationMatrix core.TransformationMatrix,
	mediaBox **MediaBox,
	cropBox **CropBox,
) {
	if transformationMatrix != core.Identity {
		mb := (*mediaBox).Bounds
		transformedMB := transformationMatrix.TransformRect(mb)
		normalisedMB := geometry.RectangleNormalise(transformedMB)
		*mediaBox = NewMediaBox(normalisedMB)

		cb := (*cropBox).Bounds
		transformedCB := transformationMatrix.TransformRect(cb)
		normalisedCB := geometry.RectangleNormalise(transformedCB)
		*cropBox = NewCropBox(normalisedCB)
	}
}

func newRotationFromPtr(rotation *int) PageRotationDegrees {
	if rotation == nil {
		return PageRotationDegrees{Value: 0}
	}
	return PageRotationDegrees{Value: *rotation}
}

// resolveToArray resolves a token to an ArrayToken, unwrapping indirect references if needed.
// Used by pages.go and resource_store.go.
func resolveToArray(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.ArrayToken {
	arr, _ := parts.TryGet[*tokens.ArrayToken](token, scanner)
	return arr
}

// resolveToDictionary resolves a token to a DictionaryToken, unwrapping indirect references if needed.
// Used by pages.go and resource_store.go.
func resolveToDictionary(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	dict, _ := parts.TryGet[*tokens.DictionaryToken](token, scanner)
	return dict
}

// arrayToRectangle converts an array of 4 numeric tokens to a PdfRectangle, resolving indirect refs.
// Used by pages.go.
func arrayToRectangle(data []tokens.Token, scanner tokenization.PdfTokenScanner) *core.PdfRectangle {
	if len(data) < 4 {
		return nil
	}

	var coords [4]float64
	for i := 0; i < 4; i++ {
		item := data[i]
		switch t := item.(type) {
		case *tokens.NumericToken:
			coords[i] = t.DoubleVal()
		case *tokens.IndirectReferenceToken:
			if scanner != nil {
				resolved := scanner.Get(t.Data())
				if resolved != nil {
					if numericTok, ok := resolved.Data().(*tokens.NumericToken); ok {
						coords[i] = numericTok.DoubleVal()
					} else {
						panic(core.NewPdfDocumentFormatException(
							fmt.Sprintf("Could not find the object %v with type NumericToken instead, it was found with type %T.", item, resolved.Data())))
					}
				} else {
					panic(core.NewPdfDocumentFormatException(
						fmt.Sprintf("Could not find the object %v with type NumericToken instead, it was found with type %T.", item, item)))
				}
			} else {
				panic(core.NewPdfDocumentFormatException(
					fmt.Sprintf("Could not find the object %v with type NumericToken instead, it was found with type %T.", item, item)))
			}
		default:
			panic(core.NewPdfDocumentFormatException(
				fmt.Sprintf("Could not find the object %v with type NumericToken instead, it was found with type %T.", item, item)))
		}
	}

	rect := core.NewPdfRectangleFloat(coords[0], coords[1], coords[2], coords[3])
	return &rect
}
