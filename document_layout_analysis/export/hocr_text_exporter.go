package export

import (
	"fmt"
	"math"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
)

const (
	xmlHeader = `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"` + "\n" +
		"\t" + `"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">` + "\n"

	hocrjs = "<script src='https://unpkg.com/hocrjs'></script>\n"
)

// HOcrTextExporter exports page text as hOCR v1.2 HTML.
// See http://kba.cloud/hocr-spec/1.2/
type HOcrTextExporter struct {
	wordExtractor           content.WordExtractor
	pageSegmenter           page_segmenter.PageSegmenter
	invalidCharacterHandler InvalidCharHandler
	scale                   float64
	indentChar              string

	pageCount  int
	areaCount  int
	lineCount  int
	wordCount  int
	pathCount  int
	paraCount  int
	imageCount int

	InvalidCharStrategy
}

// NewHOcrTextExporter creates a new hOCR text exporter with the given word extractor,
// page segmenter, and optional configuration. Defaults: scale=1, indent="\t",
// strategy=DoNotCheck.
func NewHOcrTextExporter(
	wordExtractor content.WordExtractor,
	pageSegmenter page_segmenter.PageSegmenter,
	scale float64,
	indentChar string,
	strategy InvalidCharStrategy,
) *HOcrTextExporter {
	if indentChar == "" {
		indentChar = "\t"
	}

	return &HOcrTextExporter{
		wordExtractor:           wordExtractor,
		pageSegmenter:           pageSegmenter,
		scale:                   scale,
		indentChar:              indentChar,
		InvalidCharStrategy:     strategy,
		invalidCharacterHandler: GetXmlInvalidCharHandler(strategy),
	}
}

// NewHOcrTextExporterWithHandler creates a new hOCR text exporter with a custom
// invalid character handler function.
func NewHOcrTextExporterWithHandler(
	wordExtractor content.WordExtractor,
	pageSegmenter page_segmenter.PageSegmenter,
	scale float64,
	indentChar string,
	handler InvalidCharHandler,
) *HOcrTextExporter {
	if indentChar == "" {
		indentChar = "\t"
	}

	return &HOcrTextExporter{
		wordExtractor:           wordExtractor,
		pageSegmenter:           pageSegmenter,
		scale:                   scale,
		indentChar:              indentChar,
		InvalidCharStrategy:     Custom,
		invalidCharacterHandler: handler,
	}
}

// Get returns the hOCR HTML string for a single page. Excludes paths.
func (e *HOcrTextExporter) Get(page *content.Page) string {
	return e.getPageCode(page, false, "unknown", false)
}

// GetDocument returns the hOCR HTML string for an entire document.
// includePaths controls whether PDF drawing paths are rendered as linedrawing elements.
// useHocrjs adds a reference to the hocrjs script before the closing body tag.
func (e *HOcrTextExporter) GetDocument(doc interface{ GetPages() ([]any, error) }, includePaths, useHocrjs bool) string {
	allPages, err := doc.GetPages()
	if err != nil || len(allPages) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(e.getHead())
	b.WriteString(e.indentChar)
	b.WriteString("<body>\n")

	for _, pg := range allPages {
		page, ok := pg.(*content.Page)
		if !ok {
			continue
		}
		b.WriteString(e.getCodePage(page, includePaths, "unknown"))
		b.WriteString("\n")
	}

	if useHocrjs {
		b.WriteString(e.indentChar)
		b.WriteString(e.indentChar)
		b.WriteString(hocrjs)
	}

	b.WriteString(e.indentChar)
	b.WriteString("</body>")

	body := b.String()
	return xmlHeader + e.addHtmlHeader(body)
}

func (e *HOcrTextExporter) getPageCode(page *content.Page, includePaths bool, imageName string, useHocrjs bool) string {
	var b strings.Builder
	b.WriteString(e.getHead())
	b.WriteString(e.indentChar)
	b.WriteString("<body>\n")

	b.WriteString(e.getCodePage(page, includePaths, imageName))
	b.WriteString("\n")

	if useHocrjs {
		b.WriteString(e.indentChar)
		b.WriteString(e.indentChar)
		b.WriteString(hocrjs)
	}

	b.WriteString(e.indentChar)
	b.WriteString("</body>")

	body := b.String()
	return xmlHeader + e.addHtmlHeader(body)
}

func (e *HOcrTextExporter) getHead() string {
	var b strings.Builder
	b.WriteString(e.indentChar)
	b.WriteString("<head>")
	b.WriteString("\n")
	b.WriteString(e.getIndent(2))
	b.WriteString("<title></title>")
	b.WriteString("\n")
	b.WriteString(e.getIndent(2))
	b.WriteString(`<meta http-equiv='Content-Type' content='text/html;charset=utf-8' />`)
	b.WriteString("\n")
	b.WriteString(e.getIndent(2))
	b.WriteString(`<meta name='ocr-system' content='` + e.pageSegmenterName() + "|" + e.wordExtractorName() + `'>`)
	b.WriteString("\n")
	b.WriteString(e.getIndent(2))
	b.WriteString(`<meta name='ocr-capabilities' content='ocr_page ocr_carea ocr_par ocr_line ocrx_word ocr_linedrawing' />`)
	b.WriteString("\n")
	b.WriteString(e.indentChar)
	b.WriteString("</head>")
	b.WriteString("\n")
	return b.String()
}

func (e *HOcrTextExporter) addHtmlHeader(content string) string {
	return `<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">` + "\n" + content + "\n</html>"
}

func (e *HOcrTextExporter) getIndent(level int) string {
	if level <= 0 || e.indentChar == "" {
		return ""
	}
	buf := make([]byte, 0, len(e.indentChar)*level)
	for i := 0; i < level; i++ {
		buf = append(buf, e.indentChar...)
	}
	return string(buf)
}

func (e *HOcrTextExporter) getCodePage(page *content.Page, includePaths bool, imageName string) string {
	e.pageCount++
	level := 2

	bbRight := clampZero(int(math.Round(page.Width() * e.scale)))
	bbBottom := clampZero(int(math.Round(page.Height() * e.scale)))

	var b strings.Builder
	indent := e.getIndent(level)
	b.WriteString(indent)
	fmt.Fprintf(&b, `<div class='ocr_page' id='page_%d' title='image "%s"; bbox 0 0 %d %d; ppageno %d'>`,
		page.Number(), imageName, bbRight, bbBottom, page.Number()-1)

	if includePaths {
		for _, path := range page.Paths() {
			pathObj, ok := any(path).(interface{ GetBoundingRectangle() *core.PdfRectangle })
			if !ok {
				continue
			}
			bb := pathObj.GetBoundingRectangle()
			if bb == nil {
				continue
			}

			e.areaCount++
			b.WriteString("\n")
			areaIndent := e.getIndent(level + 1)
			fmt.Fprintf(&b, `%s<div class='ocr_carea' id='block_%d_%d' title='%s'>`,
				areaIndent, e.pageCount, e.areaCount, e.getCodeBBox(*bb, page.Height()))

			subpathsObj, ok2 := any(path).(interface{ Subpaths() []*core.PdfSubpath })
			if ok2 {
				for _, subPath := range subpathsObj.Subpaths() {
					subBb := subPath.GetBoundingRectangle()
					if subBb == nil {
						continue
					}
					e.pathCount++
					b.WriteString("\n")
					fmt.Fprintf(&b, `%s<span class='ocr_linedrawing' id='drawing_%d_%d' title='%s' />`,
						e.getIndent(level+2), e.pageCount, e.pathCount, e.getCodeBBox(*subBb, page.Height()))
				}
			}

			b.WriteString("\n")
			b.WriteString(areaIndent)
			b.WriteString("</div>")
		}
	}

	for _, img := range page.GetImages() {
		b.WriteString("\n")
		b.WriteString(e.getCodeImage(img, page.Height(), level+1))
	}

	words := page.GetWordsWithExtractor(e.wordExtractor)
	if len(words) > 0 {
		for _, block := range e.pageSegmenter.GetBlocks(words) {
			b.WriteString("\n")
			b.WriteString(e.getCodeArea(block, page.Height(), level+1))
		}
	}

	b.WriteString("\n")
	b.WriteString(indent)
	b.WriteString("</div>")
	return b.String()
}

func (e *HOcrTextExporter) getCodeImage(img content.PdfImage, pageHeight float64, level int) string {
	e.imageCount++
	bb := img.BoundingBox()
	var b strings.Builder
	fmt.Fprintf(&b, `%s<span class='ocr_image' id='image_%d_%d' title='%s' />`,
		e.getIndent(level), e.pageCount, e.imageCount, e.getCodeBBox(bb, pageHeight))
	return b.String()
}

func (e *HOcrTextExporter) getCodeArea(block *document_layout_analysis.TextBlock, pageHeight float64, level int) string {
	e.areaCount++

	var b strings.Builder
	indent := e.getIndent(level)
	fmt.Fprintf(&b, `%s<div class='ocr_carea' id='block_%d_%d' title='%s'>`,
		indent, e.pageCount, e.areaCount, e.getCodeBBox(block.BoundingBox, pageHeight))

	b.WriteString(e.getCodeParagraph(block, pageHeight, level+1))
	b.WriteString("\n")
	b.WriteString(indent)
	b.WriteString("</div>")
	return b.String()
}

func (e *HOcrTextExporter) getCodeParagraph(block *document_layout_analysis.TextBlock, pageHeight float64, level int) string {
	e.paraCount++

	var b strings.Builder
	indent := e.getIndent(level)
	fmt.Fprintf(&b, "\n%s<p class='ocr_par' id='par_%d_%d' title='%s'>",
		indent, e.pageCount, e.paraCount, e.getCodeBBox(block.BoundingBox, pageHeight))

	for _, line := range block.TextLines {
		b.WriteString("\n")
		b.WriteString(e.getCodeLine(line, pageHeight, level+1))
	}
	b.WriteString("\n")
	b.WriteString(indent)
	b.WriteString("</p>")
	return b.String()
}

func (e *HOcrTextExporter) getCodeLine(line *document_layout_analysis.TextLine, pageHeight float64, level int) string {
	e.lineCount++

	baseLine := 0.0
	if len(line.Words) > 0 && len(line.Words[0].Letters) > 0 {
		baseLine = line.Words[0].Letters[0].StartBaseLine.Y
		baseLine = line.BoundingBox().Bottom() - baseLine
	}

	var b strings.Builder
	indent := e.getIndent(level)
	fmt.Fprintf(&b, `%s<span class='ocr_line' id='line_%d_%d' title='%s; baseline 0 %g'>`,
		indent, e.pageCount, e.lineCount, e.getCodeBBox(line.BoundingBox(), pageHeight), baseLine)

	for _, word := range line.Words {
		b.WriteString("\n")
		b.WriteString(e.getCodeWord(word, pageHeight, level+1))
	}
	b.WriteString("\n")
	b.WriteString(indent)
	b.WriteString("</span>")
	return b.String()
}

func (e *HOcrTextExporter) getCodeWord(word *content.Word, pageHeight float64, level int) string {
	e.wordCount++

	var b strings.Builder
	indent := e.getIndent(level)
	fmt.Fprintf(&b, `%s<span class='ocrx_word' id='word_%d_%d' title='%s; x_wconf %d`,
		indent, e.pageCount, e.wordCount, e.getCodeBBox(word.BoundingBox(), pageHeight), 100)

	fmt.Fprintf(&b, "; x_font %s", word.FontName())

	if len(word.Letters) > 0 && word.Letters[0].FontSize != 1 {
		fmt.Fprintf(&b, "; x_fsize %g", word.Letters[0].FontSize)
	}
	b.WriteString("'")

	b.WriteString(">")
	b.WriteString(e.invalidCharacterHandler(word.Text))
	b.WriteString("</span> ")
	return b.String()
}

func (e *HOcrTextExporter) getCodeBBox(rect core.PdfRectangle, pageHeight float64) string {
	left := clampZero(int(math.Round(rect.Left() * e.scale)))
	top := clampZero(int(math.Round((pageHeight - rect.Top()) * e.scale)))
	right := clampZero(int(math.Round(rect.Right() * e.scale)))
	bottom := clampZero(int(math.Round((pageHeight - rect.Bottom()) * e.scale)))

	return fmt.Sprintf("bbox %d %d %d %d", left, top, right, bottom)
}

func (e *HOcrTextExporter) wordExtractorName() string {
	switch e.wordExtractor.(type) {
	case *content.DefaultWordExtractor:
		return "DefaultWordExtractor"
	default:
		return fmt.Sprintf("%T", e.wordExtractor)
	}
}

func (e *HOcrTextExporter) pageSegmenterName() string {
	return fmt.Sprintf("%T", e.pageSegmenter)
}

func clampZero(v int) int {
	if v > 0 {
		return v
	}
	return 0
}

var _ TextExporter = (*HOcrTextExporter)(nil)
