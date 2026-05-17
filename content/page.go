package content

import (
	"fmt"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PageSize represents standard PDF page sizes.
type PageSize int

const (
	PageSizeCustom PageSize = iota
	PageSizeA0
	PageSizeA1
	PageSizeA2
	PageSizeA3
	PageSizeA4
	PageSizeA5
	PageSizeA6
	PageSizeA7
	PageSizeA8
	PageSizeA9
	PageSizeA10
	PageSizeB4
	PageSizeB5
	PageSizeExecutive
	PageSizeFanfoldUSLandscape
	PageSizeFanfoldUSPortrait
	PageSizeFanfoldGermanLandscape
	PageSizeFanfoldGermanPortrait
	PageSizeFolio
	PageSizeLedger
	PageSizeLegal
	PageSizeLetter
	PageSizeTabloid
)

// AnnotationProviderIface abstracts annotation retrieval to avoid an import cycle
// between content and annotations packages (annotations imports content for Letter).
type AnnotationProviderIface interface {
	GetAnnotations() []LinkAnnotationIface
}

// Page represents a single page in a PDF document.
// It contains the content and provides access to methods for extracting text, images, annotations, etc.
type Page struct {
	dictionary         *tokens.DictionaryToken
	number             int
	cropBox            *CropBox
	mediaBox           *MediaBox
	content            *PageContent
	rotation           PageRotationDegrees
	width              float64
	height             float64
	size               PageSize
	annotationProvider AnnotationProviderIface
	pdfScanner         tokenization.PdfTokenScanner
	textOnce           sync.Once
	textValue          string
}

// Dictionary returns the raw PDF dictionary token for this page in the document.
func (p *Page) Dictionary() *tokens.DictionaryToken {
	return p.dictionary
}

// Number returns the page number (starting at 1).
func (p *Page) Number() int {
	return p.number
}

// CropBox returns the visible region of the page; content outside is clipped.
func (p *Page) CropBox() *CropBox {
	return p.cropBox
}

// MediaBox returns the boundaries of the physical medium on which the page shall be displayed or printed.
func (p *Page) MediaBox() *MediaBox {
	return p.mediaBox
}

// Rotation returns the rotation of the page in degrees clockwise. Valid values are 0, 90, 180, and 270.
func (p *Page) Rotation() PageRotationDegrees {
	return p.rotation
}

// Letters returns the set of letters drawn by the PDF content.
func (p *Page) Letters() []*Letter {
	if p.content == nil {
		return nil
	}
	return p.content.Letters()
}

// Text returns the full text of all characters on the page in presentation order.
func (p *Page) Text() string {
	p.textOnce.Do(func() {
		p.textValue = getText(p.content)
	})
	return p.textValue
}

// Width returns the width of the page in points.
func (p *Page) Width() float64 {
	return p.width
}

// Height returns the height of the page in points.
func (p *Page) Height() float64 {
	return p.height
}

// Size returns the standard page size, or PageSizeCustom if no matching standard size is found.
func (p *Page) Size() PageSize {
	return p.size
}

// NumberOfImages returns the number of images on this page.
func (p *Page) NumberOfImages() int {
	if p.content == nil {
		return 0
	}
	return p.content.NumberOfImages()
}

// Operations returns the parsed graphics state operations in the content stream for this page.
func (p *Page) Operations() []GraphicsStateOperation {
	if p.content == nil {
		return nil
	}
	return p.content.GraphicsStateOperations()
}

// Paths returns the set of PdfPaths drawn by the PDF content.
func (p *Page) Paths() []PdfPath {
	if p.content == nil {
		return nil
	}
	return p.content.Paths()
}

// GetWords returns the words for this page using the default word extractor.
func (p *Page) GetWords() []*Word {
	return p.GetWordsWithExtractor(Instance)
}

// GetWordsWithExtractor returns the words for this page using the provided word extractor.
// If wordExtractor is nil, the default Instance extractor is used.
func (p *Page) GetWordsWithExtractor(wordExtractor WordExtractor) []*Word {
	if wordExtractor == nil {
		wordExtractor = Instance
	}
	return wordExtractor.GetWords(p.Letters())
}

// GetHyperlinks returns the hyperlinks which link to external resources on the page.
// These are based on the annotations on the page with a type of '/Link'.
func (p *Page) GetHyperlinks() []*Hyperlink {
	return GetHyperlinksFromPage(p.Letters(), p.pdfScanner, p.GetAnnotations())
}

// GetImages returns any images on the page.
func (p *Page) GetImages() []PdfImage {
	if p.content == nil {
		return nil
	}
	return p.content.GetImages()
}

// GetMarkedContents returns any marked content on the page.
func (p *Page) GetMarkedContents() []MarkedContentElement {
	if p.content == nil {
		return nil
	}
	return p.content.GetMarkedContents()
}

// GetAnnotations returns the lazily evaluated set of annotations on this page.
// Returns nil if no annotation provider is configured.
func (p *Page) GetAnnotations() []LinkAnnotationIface {
	if p.annotationProvider == nil {
		return nil
	}
	return p.annotationProvider.GetAnnotations()
}

// PdfScanner returns the internal PDF token scanner used for resolving indirect references.
// Exposed for callers that need to resolve tokens when working with annotations and hyperlinks.
func (p *Page) PdfScanner() tokenization.PdfTokenScanner {
	return p.pdfScanner
}

// GetOptionalContents returns any optional content on the page grouped by name.
// Does not handle XObjects and annotations.
func (p *Page) GetOptionalContents() map[string][]*OptionalContentGroupElement {
	mcesOptional := make([]*OptionalContentGroupElement, 0)

	markedContents := p.GetMarkedContents()
	getOptionalContentsRecursively(markedContents, &mcesOptional)

	result := make(map[string][]*OptionalContentGroupElement)
	for _, ocg := range mcesOptional {
		name := ""
		if ocg.Name != nil {
			name = *ocg.Name
		}
		result[name] = append(result[name], ocg)
	}

	return result
}

func getOptionalContentsRecursively(markedContentElements []MarkedContentElement, mcesOptional *[]*OptionalContentGroupElement) {
	if len(markedContentElements) == 0 {
		return
	}

	for _, mce := range markedContentElements {
		if mce.Tag == "OC" {
			ocg, err := NewOptionalContentGroupElement(&mce)
			if err == nil {
				*mcesOptional = append(*mcesOptional, ocg)
			}
		} else if len(mce.Children) > 0 {
			getOptionalContentsRecursively(mce.Children, mcesOptional)
		}
	}
}

func getText(content *PageContent) string {
	letters := content.Letters()
	if letters == nil || len(letters) == 0 {
		return ""
	}

	totalLen := 0
	for _, l := range letters {
		totalLen += len(l.Value)
	}

	buf := strings.Builder{}
	buf.Grow(totalLen)
	for _, l := range letters {
		buf.WriteString(l.Value)
	}
	return buf.String()
}

// NewPage creates a new Page with the given parameters.
func NewPage(
	number int,
	dictionary *tokens.DictionaryToken,
	mediaBox *MediaBox,
	cropBox *CropBox,
	rotation PageRotationDegrees,
	content *PageContent,
	annotationProvider AnnotationProviderIface,
	pdfScanner tokenization.PdfTokenScanner,
) (*Page, error) {
	if number <= 0 {
		return nil, fmt.Errorf("page number cannot be 0 or negative")
	}

	if content == nil {
		return nil, fmt.Errorf("content cannot be nil")
	}

	if dictionary == nil {
		return nil, fmt.Errorf("dictionary cannot be nil")
	}

	if pdfScanner == nil {
		return nil, fmt.Errorf("pdfScanner cannot be nil")
	}

	viewBox := intersectRectangles(mediaBox.Bounds, cropBox.Bounds)

	if viewBox == nil {
		viewBox = &cropBox.Bounds
	}

	width := viewBox.Width
	height := viewBox.Height
	size := getPageSize(*viewBox)

	return &Page{
		dictionary:         dictionary,
		number:             number,
		cropBox:            cropBox,
		mediaBox:           mediaBox,
		content:            content,
		rotation:           rotation,
		width:              width,
		height:             height,
		size:               size,
		annotationProvider: annotationProvider,
		pdfScanner:         pdfScanner,
	}, nil
}

// intersectRectangles returns the intersection of two rectangles, or nil if they don't overlap.
func intersectRectangles(a, b core.PdfRectangle) *core.PdfRectangle {
	minX := maxOf(a.Left(), b.Left())
	minY := maxOf(a.Bottom(), b.Bottom())
	maxX := minOf(a.Right(), b.Right())
	maxY := minOf(a.Top(), b.Top())

	if minX > maxX || minY > maxY {
		return nil
	}

	rect := core.NewPdfRectangleFloat(minX, minY, maxX, maxY)
	return &rect
}

// getPageSize determines the standard page size for the given rectangle dimensions.
func getPageSize(rect core.PdfRectangle) PageSize {
	w := rect.Width
	h := rect.Height

	swap := false
	if w < h {
		w, h = h, w
		swap = true
	}

	eq := func(a, b float64) bool { return absFloat(a-b) < 2.0 }

	checkSize := func(sw, sh float64) bool {
		return eq(w, sw) && eq(h, sh) || (swap && eq(w, sh) && eq(h, sw))
	}

	ppi := 72.0

	switch {
	case checkSize(841*PointsPerMm, 1189*PointsPerMm):
		return PageSizeA0
	case checkSize(594*PointsPerMm, 841*PointsPerMm):
		return PageSizeA1
	case checkSize(420*PointsPerMm, 594*PointsPerMm):
		return PageSizeA2
	case checkSize(297*PointsPerMm, 420*PointsPerMm):
		return PageSizeA3
	case checkSize(210*PointsPerMm, 297*PointsPerMm):
		return PageSizeA4
	case checkSize(148*PointsPerMm, 210*PointsPerMm):
		return PageSizeA5
	case checkSize(105*PointsPerMm, 148*PointsPerMm):
		return PageSizeA6
	case checkSize(75*PointsPerMm, 105*PointsPerMm):
		return PageSizeA7
	case checkSize(52*PointsPerMm, 74*PointsPerMm):
		return PageSizeA8
	case checkSize(37*PointsPerMm, 52*PointsPerMm):
		return PageSizeA9
	case checkSize(26*PointsPerMm, 37*PointsPerMm):
		return PageSizeA10
	case checkSize(8.5*ppi, 11*ppi):
		return PageSizeLetter
	case checkSize(8.5*ppi, 14*ppi):
		return PageSizeLegal
	case checkSize(17*ppi, 11*ppi):
		return PageSizeLedger
	case checkSize(11*ppi, 17*ppi):
		return PageSizeTabloid
	case checkSize(7.5*ppi, 10*ppi):
		return PageSizeExecutive
	default:
		return PageSizeCustom
	}
}

func maxOf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minOf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// pageSizeRectangles maps PageSize values to their standard dimensions in points.
var pageSizeRectangles = map[PageSize]struct{ w, h float64 }{
	PageSizeA0:        {2384, 3370},
	PageSizeA1:        {1684, 2384},
	PageSizeA2:        {1190, 1684},
	PageSizeA3:        {842, 1190},
	PageSizeA4:        {595, 842},
	PageSizeA5:        {420, 595},
	PageSizeA6:        {298, 420},
	PageSizeA7:        {210, 298},
	PageSizeA8:        {147, 210},
	PageSizeA9:        {105, 147},
	PageSizeA10:       {74, 105},
	PageSizeLetter:    {612, 792},
	PageSizeLegal:     {612, 1008},
	PageSizeLedger:    {1224, 792},
	PageSizeTabloid:   {792, 1224},
	PageSizeExecutive: {540, 720},
}

// TryGetPdfRectangle returns the standard PdfRectangle for the page size.
// Returns false if the size is PageSizeCustom or unrecognized.
func (s PageSize) TryGetPdfRectangle() (core.PdfRectangle, bool) {
	dims, ok := pageSizeRectangles[s]
	if !ok || dims.w == 0 {
		var zero core.PdfRectangle
		return zero, false
	}
	return core.NewPdfRectangleFloat(0, 0, dims.w, dims.h), true
}
