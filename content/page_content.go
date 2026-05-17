package content

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// PageContent wraps content parsed from a page content stream for access.
// It contains a replayable stack of drawing instructions for page content
// from a content stream in addition to lazily evaluated state such as text on the page or images.
type PageContent struct {
	graphicsStateOperations []GraphicsStateOperation
	letters                 []*Letter
	paths                   []PdfPath
	images                  []core.Union[xobjects.XObjectContentRecord, *InlineImage]
	markedContents          []MarkedContentElement
	pdfScanner              tokenization.PdfTokenScanner
	filterProvider          filters.FilterProvider
	resourceStore           ResourceStore
}

// NewPageContent creates a new PageContent.
func NewPageContent(
	graphicsStateOperations []GraphicsStateOperation,
	letters []*Letter,
	paths []PdfPath,
	images []core.Union[xobjects.XObjectContentRecord, *InlineImage],
	markedContents []MarkedContentElement,
	pdfScanner tokenization.PdfTokenScanner,
	filterProvider filters.FilterProvider,
	resourceStore ResourceStore,
) (*PageContent, error) {
	if pdfScanner == nil {
		return nil, core.NewPdfDocumentFormatException("pdfScanner cannot be null")
	}

	if filterProvider == nil {
		return nil, core.NewPdfDocumentFormatException("filterProvider cannot be null")
	}

	if resourceStore == nil {
		return nil, core.NewPdfDocumentFormatException("resourceStore cannot be null")
	}

	return &PageContent{
		graphicsStateOperations: graphicsStateOperations,
		letters:                 letters,
		paths:                   paths,
		images:                  images,
		markedContents:          markedContents,
		pdfScanner:              pdfScanner,
		filterProvider:          filterProvider,
		resourceStore:           resourceStore,
	}, nil
}

// GraphicsStateOperations returns the parsed graphics state operations in the content stream.
func (pc *PageContent) GraphicsStateOperations() []GraphicsStateOperation {
	if pc == nil {
		return nil
	}
	return pc.graphicsStateOperations
}

// Letters returns the set of letters drawn by the PDF content.
func (pc *PageContent) Letters() []*Letter {
	if pc == nil {
		return nil
	}
	return pc.letters
}

// Paths returns the set of PdfPaths drawn by the PDF content.
func (pc *PageContent) Paths() []PdfPath {
	if pc == nil {
		return nil
	}
	return pc.paths
}

// NumberOfImages returns the number of images on this page.
func (pc *PageContent) NumberOfImages() int {
	if pc == nil {
		return 0
	}
	return len(pc.images)
}

// GetImages returns any images on the page, evaluating XObject references lazily.
func (pc *PageContent) GetImages() []PdfImage {
	if pc == nil {
		return nil
	}

	result := make([]PdfImage, 0, len(pc.images))

	for _, img := range pc.images {
		if record, ok := img.TryGetFirst(); ok {
			pdfImg, err := ReadImage(&record, pc.pdfScanner, pc.filterProvider, pc.resourceStore)
			if err != nil || pdfImg == nil {
				continue
			}
			result = append(result, pdfImg)
		} else if inlineImg, ok := img.TryGetSecond(); ok {
			result = append(result, inlineImg)
		}
	}

	return result
}

// GetMarkedContents returns any marked content on the page.
func (pc *PageContent) GetMarkedContents() []MarkedContentElement {
	if pc == nil {
		return nil
	}
	return pc.markedContents
}
