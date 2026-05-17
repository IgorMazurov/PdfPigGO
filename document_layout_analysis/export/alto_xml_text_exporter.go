package export

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"time"
	"unicode/utf8"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/alto"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
)

// AltoXmlTextExporter exports page text as ALTO 4.1 XML.
// See https://github.com/altoxml/schema
type AltoXmlTextExporter struct {
	wordExtractor           content.WordExtractor
	pageSegmenter           page_segmenter.PageSegmenter
	invalidCharacterHandler InvalidCharHandler
	scale                   float64
	indentChar              string

	pageCount             int
	pageSpaceCount        int
	graphicalElementCount int
	illustrationCount     int
	textBlockCount        int
	textLineCount         int
	stringCount           int
	glyphCount            int

	InvalidCharStrategy
}

// NewAltoXmlTextExporter creates a new ALTO XML text exporter with the given word extractor,
// page segmenter, and optional configuration. Defaults: scale=1, indent="\t",
// strategy=DoNotCheck.
func NewAltoXmlTextExporter(
	wordExtractor content.WordExtractor,
	pageSegmenter page_segmenter.PageSegmenter,
	scale float64,
	indentChar string,
	strategy InvalidCharStrategy,
) *AltoXmlTextExporter {
	if indentChar == "" {
		indentChar = "\t"
	}

	return &AltoXmlTextExporter{
		wordExtractor:           wordExtractor,
		pageSegmenter:           pageSegmenter,
		scale:                   scale,
		indentChar:              indentChar,
		invalidCharacterHandler: GetXmlInvalidCharHandler(strategy),
		InvalidCharStrategy:     strategy,
	}
}

// NewAltoXmlTextExporterWithHandler creates a new ALTO XML text exporter with a custom
// invalid character handler function.
func NewAltoXmlTextExporterWithHandler(
	wordExtractor content.WordExtractor,
	pageSegmenter page_segmenter.PageSegmenter,
	scale float64,
	indentChar string,
	handler InvalidCharHandler,
) *AltoXmlTextExporter {
	if indentChar == "" {
		indentChar = "\t"
	}

	return &AltoXmlTextExporter{
		wordExtractor:           wordExtractor,
		pageSegmenter:           pageSegmenter,
		scale:                   scale,
		indentChar:              indentChar,
		invalidCharacterHandler: handler,
		InvalidCharStrategy:     Custom,
	}
}

// Get returns the ALTO XML string for a single page.
func (e *AltoXmlTextExporter) Get(page *content.Page) string {
	return e.getPageXml(page, false)
}

// GetDocument returns the ALTO XML string for all pages in the document.
func (e *AltoXmlTextExporter) GetDocument(doc *content.PdfDocument, includePaths bool) (string, error) {
	altoDoc := e.createAltoDocument("unknown")

	pagesAny, err := doc.GetPages()
	if err != nil {
		return "", fmt.Errorf("failed to get pages: %w", err)
	}

	altoPages := make([]*alto.AltoPage, 0, len(pagesAny))
	for _, pg := range pagesAny {
		page, ok := pg.(*content.Page)
		if !ok {
			continue
		}
		altoPages = append(altoPages, e.toAltoPage(page, includePaths))
	}

	altoDoc.Layout.Pages = altoPages
	return e.serialize(altoDoc), nil
}

// Deserialize reads an ALTO XML file and returns the parsed AltoDocument.
func Deserialize(xmlPath string) (*alto.AltoDocument, error) {
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	charMap := make(map[rune]rune)
	processedData := processInvalidCharsForUnmarshal(data, &charMap)

	doc := &alto.AltoDocument{}
	if err := xml.Unmarshal(processedData, doc); err != nil {
		return nil, fmt.Errorf("failed to deserialize ALTO document: %w", err)
	}

	restoreCharsInDoc(doc, charMap)

	return doc, nil
}

func (e *AltoXmlTextExporter) getPageXml(page *content.Page, includePaths bool) string {
	altoDoc := e.createAltoDocument("unknown")
	altoDoc.Layout.Pages = []*alto.AltoPage{e.toAltoPage(page, includePaths)}
	return e.serialize(altoDoc)
}

func (e *AltoXmlTextExporter) createAltoDocument(fileName string) *alto.AltoDocument {
	return &alto.AltoDocument{
		SchemaVersion: "4",
		Description:   e.getAltoDescription(fileName),
		Layout:        &alto.AltoLayout{},
	}
}

func (e *AltoXmlTextExporter) toAltoPage(page *content.Page, includePaths bool) *alto.AltoPage {
	e.pageCount = page.Number()
	e.pageSpaceCount++

	height := float32(math.Round(page.Height()*e.scale))
	width := float32(math.Round(page.Width()*e.scale))

	printSpaceId := fmt.Sprintf("P%d_PS%05d", e.pageCount, e.pageSpaceCount)

	altoPage := &alto.AltoPage{
		Id:            fmt.Sprintf("P%d", e.pageCount),
		PhysicalImgNr: float32(page.Number()),
		Quality:       ptrAltoQuality(alto.AltoQualityOK),
		Position:      ptrAltoPosition(alto.AltoPositionCover),
		PrintSpace: &alto.AltoPageSpace{
			AltoPositionedElement: alto.AltoPositionedElement{
				Height:           &height,
				Width:            &width,
				HorizontalPosition: floatPtr(0),
				VerticalPosition:   floatPtr(0),
			},
			Id: printSpaceId,
		},
	}

	h := page.Height()
	words := page.GetWordsWithExtractor(e.wordExtractor)
	blocks := e.pageSegmenter.GetBlocks(words)

	textBlocks := make([]alto.AltoTextBlock, 0, len(blocks))
	for _, block := range blocks {
		textBlocks = append(textBlocks, *e.toAltoTextBlock(block, h))
	}
	altoPage.PrintSpace.TextBlocks = textBlocks

	images := page.GetImages()
	illustrations := make([]alto.AltoIllustration, 0, len(images))
	for _, img := range images {
		illustrations = append(illustrations, *e.toAltoIllustration(img, h))
	}
	altoPage.PrintSpace.Illustrations = illustrations

	if includePaths {
		graphicalElements := make([]alto.AltoGraphicalElement, 0)
		for _, path := range page.Paths() {
			if ge := e.toAltoGraphicalElement(path, h); ge != nil {
				graphicalElements = append(graphicalElements, *ge)
			}
		}
		altoPage.PrintSpace.GraphicalElements = graphicalElements
	}

	return altoPage
}

func (e *AltoXmlTextExporter) toAltoGraphicalElement(path content.PdfPath, height float64) *alto.AltoGraphicalElement {
	e.graphicalElementCount++

	boundingRect := getBoundingRectangle(path)
	if boundingRect == nil {
		return nil
	}

	rect := *boundingRect

	vpos := float32(math.Round((height - rect.TopLeft.Y) * e.scale))
	hpos := float32(math.Round(rect.BottomLeft.X * e.scale))
	h := float32(math.Round(rect.Height * e.scale))
	w := float32(math.Round(rect.Width * e.scale))

	return &alto.AltoGraphicalElement{
		AltoBlock: alto.AltoBlock{
			AltoPositionedElement: alto.AltoPositionedElement{
				Height:             &h,
				Width:              &w,
				HorizontalPosition: &hpos,
				VerticalPosition:   &vpos,
			},
			Id:       fmt.Sprintf("P%d_GE%05d", e.pageCount, e.graphicalElementCount),
			Rotation: floatPtr(0),
		},
	}
}

func (e *AltoXmlTextExporter) toAltoIllustration(img content.PdfImage, height float64) *alto.AltoIllustration {
	e.illustrationCount++
	rect := img.BoundingBox()

	vpos := float32(math.Round((height - rect.TopLeft.Y) * e.scale))
	hpos := float32(math.Round(rect.BottomLeft.X * e.scale))
	h := float32(math.Round(rect.Height * e.scale))
	w := float32(math.Round(rect.Width * e.scale))

	return &alto.AltoIllustration{
		AltoBlock: alto.AltoBlock{
			AltoPositionedElement: alto.AltoPositionedElement{
				Height:             &h,
				Width:              &w,
				HorizontalPosition: &hpos,
				VerticalPosition:   &vpos,
			},
			Id:       fmt.Sprintf("P%d_I%05d", e.pageCount, e.illustrationCount),
			Rotation: floatPtr(0),
		},
		FileId: "",
	}
}

func (e *AltoXmlTextExporter) toAltoTextBlock(block *document_layout_analysis.TextBlock, height float64) *alto.AltoTextBlock {
	e.textBlockCount++

	vpos := float32(math.Round((height - block.BoundingBox.TopLeft.Y) * e.scale))
	hpos := float32(math.Round(block.BoundingBox.BottomLeft.X * e.scale))
	h := float32(math.Round(block.BoundingBox.Height * e.scale))
	w := float32(math.Round(block.BoundingBox.Width * e.scale))

	textLines := make([]alto.AltoTextBlockTextLine, 0, len(block.TextLines))
	for _, line := range block.TextLines {
		textLines = append(textLines, *e.toAltoTextLine(line, height))
	}

	id := fmt.Sprintf("P%d_TB%05d", e.pageCount, e.textBlockCount)

	return &alto.AltoTextBlock{
		AltoBlock: alto.AltoBlock{
			AltoPositionedElement: alto.AltoPositionedElement{
				Height:           &h,
				Width:            &w,
				HorizontalPosition: &hpos,
				VerticalPosition:   &vpos,
			},
			Id:       id,
			Rotation: floatPtr(0),
		},
		TextLines: textLines,
	}
}

func (e *AltoXmlTextExporter) toAltoTextLine(line *document_layout_analysis.TextLine, height float64) *alto.AltoTextBlockTextLine {
	e.textLineCount++

	strings := make([]alto.AltoString, 0, len(line.Words))
	for _, word := range line.Words {
		strings = append(strings, *e.toAltoString(word, height))
	}

	vpos := float32(math.Round((height - line.BoundingBox().TopLeft.Y) * e.scale))
	hpos := float32(math.Round(line.BoundingBox().BottomLeft.X * e.scale))
	h := float32(math.Round(line.BoundingBox().Height * e.scale))
	w := float32(math.Round(line.BoundingBox().Width * e.scale))

	id := fmt.Sprintf("P%d_TL%05d", e.pageCount, e.textLineCount)

	return &alto.AltoTextBlockTextLine{
		AltoPositionedElement: alto.AltoPositionedElement{
			Height:           &h,
			Width:            &w,
			HorizontalPosition: &hpos,
			VerticalPosition:   &vpos,
		},
		Strings: strings,
		Id:      id,
	}
}

func (e *AltoXmlTextExporter) toAltoString(word *content.Word, height float64) *alto.AltoString {
	e.stringCount++

	glyphs := make([]alto.AltoGlyph, 0, len(word.Letters))
	for _, letter := range word.Letters {
		glyphs = append(glyphs, *e.toAltoGlyph(letter, height))
	}

	vpos := float32(math.Round((height - word.BoundingBox().TopLeft.Y) * e.scale))
	hpos := float32(math.Round(word.BoundingBox().BottomLeft.X * e.scale))
	h := float32(math.Round(word.BoundingBox().Height * e.scale))
	w := float32(math.Round(word.BoundingBox().Width * e.scale))

	ccParts := make([]string, 0, len(glyphs))
	for _, g := range glyphs {
		if g.Gc != nil {
			ccParts = append(ccParts, fmt.Sprintf("%.1f", 9.0*(1.0-*g.Gc)))
		}
	}
	cc := ""
	for i, part := range ccParts {
		if i > 0 {
			cc += " "
		}
		cc += part
	}

	id := fmt.Sprintf("P%d_ST%05d", e.pageCount, e.stringCount)

	return &alto.AltoString{
		AltoPositionedElement: alto.AltoPositionedElement{
			Height:           &h,
			Width:            &w,
			HorizontalPosition: &hpos,
			VerticalPosition:   &vpos,
		},
		Glyphs:  glyphs,
		Cc:      cc,
		Content: e.invalidCharacterHandler(word.Text),
		Id:      id,
	}
}

func (e *AltoXmlTextExporter) toAltoGlyph(letter *content.Letter, height float64) *alto.AltoGlyph {
	e.glyphCount++

	vpos := float32(math.Round((height - letter.BoundingBox.TopLeft.Y) * e.scale))
	hpos := float32(math.Round(letter.BoundingBox.BottomLeft.X * e.scale))
	h := float32(math.Round(letter.BoundingBox.Height * e.scale))
	w := float32(math.Round(letter.BoundingBox.Width * e.scale))

	gc := float32(1.0)

	id := fmt.Sprintf("P%d_ST%05d_G%02d", e.pageCount, e.stringCount, e.glyphCount)

	return &alto.AltoGlyph{
		AltoPositionedElement: alto.AltoPositionedElement{
			Height:           &h,
			Width:            &w,
			HorizontalPosition: &hpos,
			VerticalPosition:   &vpos,
		},
		Gc:      &gc,
		Content: e.invalidCharacterHandler(letter.Value),
		Id:      id,
	}
}

func (e *AltoXmlTextExporter) getAltoDescription(fileName string) *alto.AltoDescription {
	processingCategory := alto.AltoProcessingCategoryOther

	return &alto.AltoDescription{
		MeasurementUnit: alto.AltoMeasurementUnitPixel,
		Processings: []alto.AltoDescriptionProcessing{
			{
				AltoProcessingStep: alto.AltoProcessingStep{
					ProcessingCategory:     &processingCategory,
					ProcessingDateTime:     time.Now().UTC().Format(time.RFC3339),
					ProcessingSoftware: &alto.AltoProcessingSoftware{
					SoftwareName:           "PdfPig",
					SoftwareCreator:        "https://github.com/UglyToad/PdfPig",
					ApplicationDescription: "Read and extract text and other content from PDFs in C# (port of PdfBox)",
					SoftwareVersion:        "x.x.xx",
				},
					ProcessingStepSettings: fmt.Sprintf("%T|%T", e.pageSegmenter, e.wordExtractor),
				},
				Id: fmt.Sprintf("P%d_D1", e.pageCount),
			},
		},
		SourceImageInfo: &alto.AltoSourceImageInformation{
			FileName: fileName,
			FileIdentifiers: []alto.AltoFileIdentifier{
				{},
			},
			DocumentIdentifiers: []alto.AltoDocumentIdentifier{
				{},
			},
		},
	}
}

func (e *AltoXmlTextExporter) serialize(doc *alto.AltoDocument) string {
	var charMap map[rune]rune

	if e.InvalidCharStrategy == DoNotCheck {
		charMap = replaceInvalidCharsWithPlaceholders(doc)
	}

	data, err := xml.MarshalIndent(doc, "", e.indentChar)
	if err != nil {
		return fmt.Sprintf("<error>failed to serialize ALTO document: %s</error>", err.Error())
	}

	if e.InvalidCharStrategy == DoNotCheck && len(charMap) > 0 {
		data = restorePlaceholdersInBytes(data, charMap)
	}

	header := []byte(xml.Header)
	result := make([]byte, 0, len(header)+len(data))
	result = append(result, header...)
	result = append(result, data...)
	return string(result)
}

// replaceInvalidCharsWithPlaceholders replaces all invalid XML characters in Content
// fields of the Alto document with unique Unicode placeholders (U+F0XX where XX is
// the hex value of the original char). Returns a map from placeholder to original
// character for restoration after marshaling. This works around Go's encoding/xml
// which always replaces invalid XML chars with U+FFFD during marshaling.
func replaceInvalidCharsWithPlaceholders(doc *alto.AltoDocument) map[rune]rune {
	charMap := make(map[rune]rune)

	for _, page := range doc.Layout.Pages {
		if page == nil || page.PrintSpace == nil {
			continue
		}
		for tbIdx := range page.PrintSpace.TextBlocks {
			block := &page.PrintSpace.TextBlocks[tbIdx]
			for tlIdx := range block.TextLines {
				line := &block.TextLines[tlIdx]
				for sIdx := range line.Strings {
					str := &line.Strings[sIdx]
					str.Content = replaceInString(str.Content, &charMap)

					for gIdx := range str.Glyphs {
						glyph := &str.Glyphs[gIdx]
						glyph.Content = replaceInString(glyph.Content, &charMap)
					}
				}
			}
		}
	}

	return charMap
}

// replaceInString replaces each invalid XML character in s with a unique Unicode
// placeholder (U+F0XX) and records the mapping in charMap.
func replaceInString(s string, charMap *map[rune]rune) string {
	var b []rune
	for _, r := range s {
		if !isValidXmlChar(r) && r != utf8.RuneError {
			placeholder := rune(0xF000 + uint32(r))
			(*charMap)[placeholder] = r
			b = append(b, placeholder)
		} else {
			b = append(b, r)
		}
	}
	return string(b)
}

// restorePlaceholdersInBytes replaces placeholder characters in the marshaled XML bytes
// back to their original invalid characters using the charMap.
func restorePlaceholdersInBytes(data []byte, charMap map[rune]rune) []byte {
	for placeholder, original := range charMap {
		placeholderBytes := []byte(string(placeholder))
		originalBytes := []byte(string(original))
		data = bytes.ReplaceAll(data, placeholderBytes, originalBytes)
	}
	return data
}

// processInvalidCharsForUnmarshal scans raw XML bytes and replaces each invalid XML
// character with a unique Unicode placeholder (U+F0XX). Returns the processed byte
// slice and populates charMap with the placeholder-to-original mapping for later
// restoration after unmarshaling.
func processInvalidCharsForUnmarshal(data []byte, charMap *map[rune]rune) []byte {
	result := make([]byte, 0, len(data))
	i := 0
	for i < len(data) {
		r, size := utf8.DecodeRune(data[i:])
		if !isValidXmlChar(r) && r != utf8.RuneError {
			placeholder := rune(0xF000 + uint32(r))
			(*charMap)[placeholder] = r
			result = append(result, []byte(string(placeholder))...)
		} else {
			for j := 0; j < size; j++ {
				result = append(result, data[i+j])
			}
		}
		i += size
	}
	return result
}

// restoreCharsInDoc walks the Alto document tree and restores original characters
// from placeholders using the charMap.
func restoreCharsInDoc(doc *alto.AltoDocument, charMap map[rune]rune) {
	if len(charMap) == 0 {
		return
	}

	for _, page := range doc.Layout.Pages {
		if page == nil || page.PrintSpace == nil {
			continue
		}
		for tbIdx := range page.PrintSpace.TextBlocks {
			block := &page.PrintSpace.TextBlocks[tbIdx]
			for tlIdx := range block.TextLines {
				line := &block.TextLines[tlIdx]
				for sIdx := range line.Strings {
					str := &line.Strings[sIdx]
					str.Content = restoreString(str.Content, charMap)

					for gIdx := range str.Glyphs {
						glyph := &str.Glyphs[gIdx]
						glyph.Content = restoreString(glyph.Content, charMap)
					}
				}
			}
		}
	}
}

// restoreString replaces placeholder characters with their original values using charMap.
func restoreString(s string, charMap map[rune]rune) string {
	var b []rune
	for _, r := range s {
		if original, ok := charMap[r]; ok {
			b = append(b, original)
		} else {
			b = append(b, r)
		}
	}
	return string(b)
}

// getBoundingRectangle attempts to extract the bounding rectangle from a PdfPath.
// Returns nil if the path has no bounding rectangle or is not supported.
func getBoundingRectangle(path content.PdfPath) *core.PdfRectangle {
	type boundingRectProvider interface {
		GetBoundingRectangle() *core.PdfRectangle
	}

	if provider, ok := any(&path).(boundingRectProvider); ok {
		return provider.GetBoundingRectangle()
	}
	return nil
}

// Helper functions for pointer types.
func floatPtr(v float32) *float32 { return &v }

func ptrAltoQuality(q alto.AltoQuality) *alto.AltoQuality { return &q }

func ptrAltoPosition(p alto.AltoPosition) *alto.AltoPosition { return &p }

// Compile-time interface assertion.
var _ TextExporter = (*AltoXmlTextExporter)(nil)
