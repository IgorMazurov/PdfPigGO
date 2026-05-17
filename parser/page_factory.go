package parser

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// PageFactory creates concrete Page instances from PDF page dictionaries.
// It extends the base factory logic with annotation extraction and content stream processing.
type PageFactory struct {
	ParsingOptions    *content.ParsingOptions
	PdfScanner        tokenization.PdfTokenScanner
	ResourceStore     content.ResourceStore
	FilterProvider    content.LookupFilterProvider
	PageContentParser content.PageContentParser
}

// NewPageFactory creates a new PageFactory with the given dependencies.
func NewPageFactory(
	pdfScanner tokenization.PdfTokenScanner,
	resourceStore content.ResourceStore,
	filterProvider content.LookupFilterProvider,
	pageContentParser content.PageContentParser,
	parsingOptions *content.ParsingOptions,
) *PageFactory {
	return &PageFactory{
		PdfScanner:        pdfScanner,
		ResourceStore:     resourceStore,
		FilterProvider:    filterProvider,
		PageContentParser: pageContentParser,
		ParsingOptions:    parsingOptions,
	}
}

// Create constructs a Page from the given dictionary and context.
func (f *PageFactory) Create(
	number int,
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *content.PageTreeMembers,
	namedDestinations *destinations.NamedDestinations,
) (*content.Page, error) {
	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}

	var result *content.Page
	var createErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				createErr = fmt.Errorf("failed to create page %d: %v\n%s", number, r, string(stack))
			}
		}()

		pgType := dictGetOrEmpty(dictionary, tokens.Type)

		if pgType != nil && pgType.Data() != tokens.Page.Data() {
			f.ParsingOptions.Logger.Error(fmt.Sprintf("Page %d had its type specified as %s rather than 'Page'.", number, pgType.Data()))
		}

		rotation := newRotationFromPtr(pageTreeMembers.Rotation)

		if rotateToken, found := dictionary.TryGet(tokens.Rotate); found {
			if numericTok, ok := rotateToken.(*tokens.NumericToken); ok {
				var err error
				rotation, err = content.NewPageRotationDegrees(numericTok.IntVal())
				if err != nil {
					f.ParsingOptions.Logger.Error(fmt.Sprintf("Invalid rotation value for page %d: %v", number, err))
				}
			}
		}

		stackDepth := 0

		for len(pageTreeMembers.ParentResources) > 0 {
			resource := pageTreeMembers.ParentResources[0]
			pageTreeMembers.ParentResources = pageTreeMembers.ParentResources[1:]

			if dict, ok := resource.(*tokens.DictionaryToken); ok {
				f.ResourceStore.LoadResourceDictionary(dict)
			} else if indRef, ok := resource.(*tokens.IndirectReferenceToken); ok {
				resolved := f.PdfScanner.Get(indRef.Data())
				if resolved != nil {
					if resDict, ok := resolved.Data().(*tokens.DictionaryToken); ok {
						f.ResourceStore.LoadResourceDictionary(resDict)
					}
				}
			}
			stackDepth++
		}

		if resourcesRaw, found := dictionary.TryGet(tokens.Resources); found {
			resDict := resolveToDictionary(resourcesRaw, f.PdfScanner)
			if resDict != nil {
				f.ResourceStore.LoadResourceDictionary(resDict)
				stackDepth++
			}
		}

		userSpaceUnit := content.GetUserSpaceUnits(dictionary)

		mediaBox := content.GetMediaBox(number, dictionary, pageTreeMembers, f.ParsingOptions.Logger)
		cropBox := content.GetCropBox(dictionary, pageTreeMembers, mediaBox, f.PdfScanner, f.ParsingOptions.Logger)

		initialMatrix := content.GetInitialMatrix(userSpaceUnit, mediaBox, cropBox, rotation, f.ParsingOptions.Logger)

		content.ApplyTransformNormalise(initialMatrix, &mediaBox, &cropBox)

		var operations []content.GraphicsStateOperation

		contentsRaw, hasContents := dictionary.TryGet(tokens.Contents)

		if !hasContents {
			operations = make([]content.GraphicsStateOperation, 0)
		} else {
			contentBytes := f.extractContentBytes(contentsRaw)

			if contentBytes == nil {
				createErr = fmt.Errorf("failed to parse the content for the page: %d", number)
				return
			}

			inputBytes := core.NewMemoryInputBytes(contentBytes)
			operations = f.PageContentParser.Parse(number, inputBytes, f.ParsingOptions.Logger)
		}

		page, err := f.processPage(
			number,
			dictionary,
			namedDestinations,
			mediaBox,
			cropBox,
			userSpaceUnit,
			rotation,
			initialMatrix,
			operations,
		)

		for i := 0; i < stackDepth; i++ {
			f.ResourceStore.UnloadResourceDictionary()
		}

		result = page
		createErr = err
	}()

	return result, createErr
}

func (f *PageFactory) extractContentBytes(contentsRaw tokens.Token) []byte {
	arrayTok := resolveToArray(contentsRaw, f.PdfScanner)

	if arrayTok != nil {
		data := arrayTok.Data()
		buf := make([]byte, 0, 1024*64)

		for i, item := range data {
			obj := resolveToIndirectRef(item, f.PdfScanner)
			if obj == nil {
				return nil
			}

			streamTok := f.resolveStream(obj, f.PdfScanner)
			if streamTok == nil {
				return nil
			}

			decoded := f.FilterProvider.DecodeStream(streamTok, f.PdfScanner)
			buf = append(buf, decoded...)

			if i < len(data)-1 {
				buf = append(buf, '\n')
			}
		}

		return buf
	}

	streamTok := f.resolveStream(contentsRaw, f.PdfScanner)
	if streamTok == nil {
		return nil
	}

	return f.FilterProvider.DecodeStream(streamTok, f.PdfScanner)
}

func (f *PageFactory) resolveStream(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.StreamToken {
	if streamTok, ok := token.(*tokens.StreamToken); ok {
		return streamTok
	}
	if indirectRef, ok := token.(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(indirectRef.Data())
		if obj != nil {
			if streamTok, ok := obj.Data().(*tokens.StreamToken); ok {
				return streamTok
			}
		}
	}
	return nil
}

func (f *PageFactory) processPage(
	pageNumber int,
	dictionary *tokens.DictionaryToken,
	namedDestinations *destinations.NamedDestinations,
	mediaBox *content.MediaBox,
	cropBox *content.CropBox,
	userSpaceUnit geometry.UserSpaceUnit,
	rotation content.PageRotationDegrees,
	initialMatrix core.TransformationMatrix,
	operations []content.GraphicsStateOperation,
) (*content.Page, error) {
	annotationProvider, err := annotations.NewAnnotationProvider(
		f.PdfScanner,
		dictionary,
		initialMatrix,
		namedDestinations,
		f.ParsingOptions.Logger,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create annotation provider: %w", err)
	}

	annotationAdapter := &annotationAdapter{provider: annotationProvider}

	if operations == nil || len(operations) == 0 {
		filterProv := toFilterProvider(f.FilterProvider)

		emptyContent, err := content.NewPageContent(
			make([]content.GraphicsStateOperation, 0),
			make([]*content.Letter, 0),
			make([]content.PdfPath, 0),
			make([]core.Union[xobjects.XObjectContentRecord, *content.InlineImage], 0),
			make([]content.MarkedContentElement, 0),
			f.PdfScanner,
			filterProv,
			f.ResourceStore,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to create empty page content: %w", err)
		}

		return content.NewPage(
			pageNumber,
			dictionary,
			mediaBox,
			cropBox,
			rotation,
			emptyContent,
			annotationAdapter,
			f.PdfScanner,
		)
	}

	ctx := graphics.NewContentStreamProcessor(
		pageNumber,
		f.ResourceStore,
		f.PdfScanner,
		f.PageContentParser,
		f.FilterProvider,
		cropBox,
		userSpaceUnit,
		rotation,
		initialMatrix,
		f.ParsingOptions,
	)

	var pageContent *content.PageContent
	func() {
		defer func() {
			if r := recover(); r != nil {
				f.ParsingOptions.Logger.Error(fmt.Sprintf("Failed to process page %d content: %v", pageNumber, r))
				if !f.ParsingOptions.UseLenientParsing {
					panic(r)
				}
				errStr := fmt.Sprintf("%v", r)
				if !f.ParsingOptions.SkipMissingFonts && strings.Contains(errStr, "Could not find the font") {
					panic(r)
				}
				pageContent = nil
			}
		}()
		pageContent = ctx.Process(pageNumber, operations)
	}()

	if pageContent == nil {
		filterProv := toFilterProvider(f.FilterProvider)
		emptyContent, createErr := content.NewPageContent(
			make([]content.GraphicsStateOperation, 0),
			make([]*content.Letter, 0),
			make([]content.PdfPath, 0),
			make([]core.Union[xobjects.XObjectContentRecord, *content.InlineImage], 0),
			make([]content.MarkedContentElement, 0),
			f.PdfScanner,
			filterProv,
			f.ResourceStore,
		)
		if createErr != nil {
			return nil, fmt.Errorf("failed to create fallback empty page content for page %d: %w", pageNumber, createErr)
		}
		pageContent = emptyContent
	}

	return content.NewPage(
		pageNumber,
		dictionary,
		mediaBox,
		cropBox,
		rotation,
		pageContent,
		annotationAdapter,
		f.PdfScanner,
	)
}

// annotationAdapter wraps *annotations.AnnotationProvider to satisfy content.AnnotationProviderIface.
type annotationAdapter struct {
	provider *annotations.AnnotationProvider
}

func (a *annotationAdapter) GetAnnotations() []content.LinkAnnotationIface {
	if a.provider == nil {
		return nil
	}
	annots := a.provider.GetAnnotations()
	result := make([]content.LinkAnnotationIface, 0, len(annots))
	for _, ann := range annots {
		result = append(result, newLinkAnnotationAdapter(ann))
	}
	return result
}

func newLinkAnnotationAdapter(a *annotations.Annotation) content.LinkAnnotationIface {
	return annotations.NewLinkAnnotationAdapter(a)
}

func toFilterProvider(fp content.LookupFilterProvider) filters.FilterProvider {
	if x, ok := fp.(filters.FilterProvider); ok {
		return x
	}
	return nil
}

// dictGetOrEmpty returns the NameToken value for a dictionary key, or nil if not found.
func dictGetOrEmpty(dict *tokens.DictionaryToken, name *tokens.NameToken) *tokens.NameToken {
	token, ok := dict.TryGet(name)
	if !ok {
		return nil
	}
	if n, ok := token.(*tokens.NameToken); ok {
		return n
	}
	return nil
}

// resolveToArray resolves a token to an ArrayToken, unwrapping indirect references if needed.
func resolveToArray(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.ArrayToken {
	if arr, ok := token.(*tokens.ArrayToken); ok {
		return arr
	}
	if indRef, ok := token.(*tokens.IndirectReferenceToken); ok && scanner != nil {
		obj := scanner.Get(indRef.Data())
		if obj != nil {
			if arr, ok := obj.Data().(*tokens.ArrayToken); ok {
				return arr
			}
		}
	}
	return nil
}

// resolveToIndirectRef resolves a token to an IndirectReferenceToken.
func resolveToIndirectRef(token tokens.Token, _ tokenization.PdfTokenScanner) *tokens.IndirectReferenceToken {
	if ref, ok := token.(*tokens.IndirectReferenceToken); ok {
		return ref
	}
	return nil
}

func newRotationFromPtr(rotation *int) content.PageRotationDegrees {
	if rotation == nil {
		return content.PageRotationDegrees{Value: 0}
	}
	return content.PageRotationDegrees{Value: *rotation}
}
