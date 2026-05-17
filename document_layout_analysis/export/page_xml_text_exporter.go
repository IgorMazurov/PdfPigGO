package export

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/export/page"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
)

// postProcessUnmarshal restores escaped invalid XML chars in all string fields of the document.
func postProcessUnmarshal(doc *page.PageXmlDocument) {
	restoreTextEquivs := func(tes []page.PageXmlTextEquiv) {
		for i, te := range tes {
			tes[i].Unicode = restoreEscapedInvalidXmlCharsInString(te.Unicode)
		}
	}

	restoreTextRegion := func(tr *page.PageXmlTextRegion) {
		restoreTextEquivs(tr.TextEquivs)
		for _, line := range tr.TextLines {
			restoreTextEquivs(line.TextEquivs)
			for _, word := range line.Words {
				restoreTextEquivs(word.TextEquivs)
			}
		}
	}

	if doc.Page != nil {
		for _, item := range doc.Page.Items {
			if tr, ok := item.(*page.PageXmlTextRegion); ok {
				restoreTextRegion(tr)
			}
		}
	}
}

// PageXmlTextExporter exports page text as PAGE-XML 2019-07-15 format.
// See https://github.com/PRImA-Research-Lab/PAGE-XML
type PageXmlTextExporter struct {
	wordExtractor           content.WordExtractor
	pageSegmenter           page_segmenter.PageSegmenter
	readingOrderDetector    reading_order_detector.ReadingOrderDetector
	invalidCharacterHandler InvalidCharHandler
	scale                   float64
	indentChar              string

	InvalidCharStrategy
}

// NewPageXmlTextExporter creates a new PAGE-XML text exporter with the given word extractor,
// page segmenter, optional reading order detector, and configuration.
// Defaults: scale=1.0, indent="\t", strategy=DoNotCheck.
func NewPageXmlTextExporter(
	wordExtractor content.WordExtractor,
	pageSegmenter page_segmenter.PageSegmenter,
	readingOrderDetector reading_order_detector.ReadingOrderDetector,
	scale float64,
	indentChar string,
	strategy InvalidCharStrategy,
) *PageXmlTextExporter {
	if indentChar == "" {
		indentChar = "\t"
	}

	return &PageXmlTextExporter{
		wordExtractor:           wordExtractor,
		pageSegmenter:           pageSegmenter,
		readingOrderDetector:    readingOrderDetector,
		scale:                   scale,
		indentChar:              indentChar,
		invalidCharacterHandler: GetXmlInvalidCharHandler(strategy),
		InvalidCharStrategy:     strategy,
	}
}

// NewPageXmlTextExporterWithHandler creates a new PAGE-XML text exporter with a custom
// invalid character handler function.
func NewPageXmlTextExporterWithHandler(
	wordExtractor content.WordExtractor,
	pageSegmenter page_segmenter.PageSegmenter,
	readingOrderDetector reading_order_detector.ReadingOrderDetector,
	scale float64,
	indentChar string,
	handler InvalidCharHandler,
) *PageXmlTextExporter {
	if indentChar == "" {
		indentChar = "\t"
	}

	return &PageXmlTextExporter{
		wordExtractor:           wordExtractor,
		pageSegmenter:           pageSegmenter,
		readingOrderDetector:    readingOrderDetector,
		scale:                   scale,
		indentChar:              indentChar,
		invalidCharacterHandler: handler,
		InvalidCharStrategy:     Custom,
	}
}

// Get returns the PAGE-XML string for a single page. Excludes paths.
func (e *PageXmlTextExporter) Get(pageObj *content.Page) string {
	return e.getPageXml(pageObj, false)
}

// PointToString converts a PDF point to a PAGE-XML coordinate string.
// The Y axis is flipped so that the origin is at the top-left corner.
func PointToString(point core.PdfPoint, pageWidth, pageHeight, scaleToApply float64) string {
	x := math.Round(point.X * scaleToApply)
	y := math.Round((pageHeight-point.Y)*scaleToApply)

	if x <= 1 {
		x = 1
	}
	if y <= 1 {
		y = 1
	}

	if x >= pageWidth-1 {
		x = pageWidth - 1
	}
	if y >= pageHeight-1 {
		y = pageHeight - 1
	}

	return strconv.FormatFloat(x, 'f', -1, 64) + "," + strconv.FormatFloat(y, 'f', -1, 64)
}

func (e *PageXmlTextExporter) getPageXml(pageObj *content.Page, includePaths bool) string {
	data := newPageXmlData()

	now := time.Now().UTC()

	pageDoc := &page.PageXmlDocument{
		Metadata: &page.PageXmlMetadata{
			Created:    now,
			LastChange: now,
			Creator:    "PdfPig",
			Comments:   fmt.Sprintf("%T|%T", e.pageSegmenter, e.wordExtractor),
		},
		PcGtsId: "pc-" + strconv.Itoa(pageObj.Number()),
	}

	pageDoc.Page = e.toPageXmlPage(pageObj, data, includePaths)

	return e.serialize(pageDoc)
}

func (e *PageXmlTextExporter) toPoints(points []core.PdfPoint, pageWidth, pageHeight float64) string {
	parts := make([]string, len(points))
	for i, p := range points {
		parts[i] = PointToString(p, pageWidth, pageHeight, e.scale)
	}
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}
	return result
}

func (e *PageXmlTextExporter) toPointsRect(rect core.PdfRectangle, pageWidth, pageHeight float64) string {
	return e.toPoints([]core.PdfPoint{
		rect.BottomLeft,
		rect.TopLeft,
		rect.TopRight,
		rect.BottomRight,
	}, pageWidth, pageHeight)
}

func (e *PageXmlTextExporter) toCoords(rect core.PdfRectangle, pageWidth, pageHeight float64) *page.PageXmlCoords {
	return &page.PageXmlCoords{
		Points: e.toPointsRect(rect, pageWidth, pageHeight),
	}
}

// toRgbEncoded converts a color to PAGE-XML RGB encoded format.
// Format: (red value) + (256 x green value) + (65536 x blue value).
func (e *PageXmlTextExporter) toRgbEncoded(color colors.Color) string {
	if color == nil {
		return ""
	}
	rgb := color.ToRGBValues()
	red := byte(math.Round(255.0 * rgb.R))
	green := byte(math.Round(255.0 * rgb.G))
	blue := byte(math.Round(255.0 * rgb.B))
	sum := int(red) + 256*int(green) + 65536*int(blue)
	return strconv.Itoa(sum)
}

func (e *PageXmlTextExporter) toPageXmlPage(pageObj *content.Page, data *pageXmlData, includePaths bool) *page.PageXmlPage {
	xmlPage := &page.PageXmlPage{
		ImageFilename: "unknown",
		ImageHeight:   int(math.Round(pageObj.Height() * e.scale)),
		ImageWidth:    int(math.Round(pageObj.Width() * e.scale)),
	}

	regions := make([]page.PageXmlRegionItem, 0)

	words := pageObj.GetWordsWithExtractor(e.wordExtractor)
	if len(words) > 0 {
		blocks := e.pageSegmenter.GetBlocks(words)

		if e.readingOrderDetector != nil {
			blocks = e.readingOrderDetector.Get(blocks)
		}

		for _, block := range blocks {
			regions = append(regions, e.toPageXmlTextRegion(block, data, pageObj.Width(), pageObj.Height()))
		}

		if len(data.OrderedRegions) > 0 {
			data.GroupOrdersCount++
			orderedItems := make([]any, len(data.OrderedRegions))
			for i, ref := range data.OrderedRegions {
				orderedItems[i] = ref
			}
			xmlPage.ReadingOrder = &page.PageXmlReadingOrder{
				Item: &page.PageXmlOrderedGroup{
					Items: orderedItems,
					Id:    "g" + strconv.Itoa(data.GroupOrdersCount),
				},
			}
		}
	}

	images := pageObj.GetImages()
	if len(images) > 0 {
		for _, img := range images {
			regions = append(regions, e.toPageXmlImageRegion(img, data, pageObj.Width(), pageObj.Height()))
		}
	}

	if includePaths {
		for _, path := range pageObj.Paths() {
			if graphElem := e.toPageXmlLineDrawingRegion(path, data, pageObj.Width(), pageObj.Height()); graphElem != nil {
				regions = append(regions, graphElem)
			}
		}
	}

	xmlPage.Items = regions
	return xmlPage
}

func (e *PageXmlTextExporter) toPageXmlLineDrawingRegion(path content.PdfPath, data *pageXmlData, pageWidth, pageHeight float64) *page.PageXmlLineDrawingRegion {
	bbox := getBoundingRectangle(path)
	if bbox != nil {
		data.RegionsCount++
		return &page.PageXmlLineDrawingRegion{
			Coords: e.toCoords(*bbox, pageWidth, pageHeight),
			Id:     "r" + strconv.Itoa(data.RegionsCount),
		}
	}
	return nil
}

func (e *PageXmlTextExporter) toPageXmlImageRegion(img content.PdfImage, data *pageXmlData, pageWidth, pageHeight float64) *page.PageXmlImageRegion {
	data.RegionsCount++
	bbox := img.BoundingBox()
	return &page.PageXmlImageRegion{
		Coords: e.toCoords(bbox, pageWidth, pageHeight),
		Id:     "r" + strconv.Itoa(data.RegionsCount),
	}
}

func (e *PageXmlTextExporter) toPageXmlTextRegion(textBlock *document_layout_analysis.TextBlock, data *pageXmlData, pageWidth, pageHeight float64) *page.PageXmlTextRegion {
	data.RegionsCount++
	regionId := "r" + strconv.Itoa(data.RegionsCount)

	if e.readingOrderDetector != nil && textBlock.ReadingOrder > -1 {
		data.OrderedRegions = append(data.OrderedRegions, &page.PageXmlRegionRefIndexed{
			RegionRef: regionId,
			Index:     textBlock.ReadingOrder,
		})
	}

	textSimpleType := page.TextParagraph

	textLines := make([]page.PageXmlTextLine, len(textBlock.TextLines))
	for i, line := range textBlock.TextLines {
		textLines[i] = *e.toPageXmlTextLine(line, data, pageWidth, pageHeight)
	}

	return &page.PageXmlTextRegion{
		Coords: e.toCoords(textBlock.BoundingBox, pageWidth, pageHeight),
		Type:   &textSimpleType,
		TextLines: textLines,
		TextEquivs: []page.PageXmlTextEquiv{
			{Unicode: e.invalidCharacterHandler(textBlock.Text)},
		},
		Id: regionId,
	}
}

func (e *PageXmlTextExporter) toPageXmlTextLine(textLine *document_layout_analysis.TextLine, data *pageXmlData, pageWidth, pageHeight float64) *page.PageXmlTextLine {
	data.LinesCount++
	production := page.Printed

	words := make([]page.PageXmlWord, len(textLine.Words))
	for i, w := range textLine.Words {
		words[i] = *e.toPageXmlWord(w, data, pageWidth, pageHeight)
	}

	return &page.PageXmlTextLine{
		Coords:     e.toCoords(textLine.BoundingBox(), pageWidth, pageHeight),
		Production: &production,
		Words:      words,
		TextEquivs: []page.PageXmlTextEquiv{
			{Unicode: e.invalidCharacterHandler(textLine.Text)},
		},
		Id: "l" + strconv.Itoa(data.LinesCount),
	}
}

func (e *PageXmlTextExporter) toPageXmlWord(word *content.Word, data *pageXmlData, pageWidth, pageHeight float64) *page.PageXmlWord {
	data.WordsCount++

	glyphs := make([]page.PageXmlGlyph, len(word.Letters))
	for i, letter := range word.Letters {
		glyphs[i] = *e.toPageXmlGlyph(letter, data, pageWidth, pageHeight)
	}

	return &page.PageXmlWord{
		Coords: e.toCoords(word.BoundingBox(), pageWidth, pageHeight),
		Glyphs: glyphs,
		TextEquivs: []page.PageXmlTextEquiv{
			{Unicode: e.invalidCharacterHandler(word.Text)},
		},
		Id: "w" + strconv.Itoa(data.WordsCount),
	}
}

func (e *PageXmlTextExporter) toPageXmlGlyph(letter *content.Letter, data *pageXmlData, pageWidth, pageHeight float64) *page.PageXmlGlyph {
	data.GlyphsCount++
	ligature := false
	production := page.Printed

	fontSize := float32(letter.FontSize)
	textStyle := &page.PageXmlTextStyle{
		FontSize:      &fontSize,
		FontFamily:    letter.FontName(),
		TextColourRgb: e.toRgbEncoded(letter.Color),
	}

	return &page.PageXmlGlyph{
		Coords:     e.toCoords(letter.BoundingBox, pageWidth, pageHeight),
		Ligature:   &ligature,
		Production: &production,
		TextStyle:  textStyle,
		TextEquivs: []page.PageXmlTextEquiv{
			{Unicode: e.invalidCharacterHandler(letter.Value)},
		},
		Id: "c" + strconv.Itoa(data.GlyphsCount),
	}
}

func (e *PageXmlTextExporter) serialize(doc *page.PageXmlDocument) string {
	// When DoNotCheck strategy is used, escape invalid XML chars before encoding,
	// then restore them after to match C#'s CheckCharacters=false behavior.
	if e.InvalidCharStrategy == DoNotCheck {
		escapeDocInvalidChars(doc)
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(&buf)

	if err := serializePageXmlDocument(enc, doc); err != nil {
		return fmt.Sprintf("<error>failed to serialize PAGE-XML document: %s</error>", err.Error())
	}

	if err := enc.Flush(); err != nil {
		return fmt.Sprintf("<error>failed to serialize PAGE-XML document: %s</error>", err.Error())
	}

	result := buf.String()
	// Convert paired empty tags to self-closing tags (matching C# XmlSerializer behavior).
	// Go's encoding/xml produces <tag></tag> for empty elements; C# uses <tag />.
	result = convertEmptyTagsToSelfClosing(result)

	// Restore escaped invalid XML chars in the serialized output.
	if e.InvalidCharStrategy == DoNotCheck {
		result = restoreEscapedInvalidXmlCharsInString(result)
	}

	return result
}

// convertEmptyTagsToSelfClosing replaces <tag attr="..."></tag> with <tag attr="..." />
// for all empty element pairs in the XML string.
func convertEmptyTagsToSelfClosing(s string) string {
	result := s
	// Convert known empty element pairs to self-closing tags.
	// Go's encoding/xml produces <tag></tag>; C# XmlSerializer uses <tag />.
	for _, tag := range []string{"Coords", "PlainText", "Unicode", "Items"} {
		// Pattern: <tag ...></tag> → <tag ... />
		pattern := regexp.MustCompile(`<` + tag + `(\s[^>]*)?>\s*</` + tag + `>`)
		result = pattern.ReplaceAllString(result, "<"+tag+"$1 />")
	}
	return result
}

// escapeDocInvalidChars replaces invalid XML chars in all TextEquiv Unicode fields with
// placeholders that survive Go's xml.Encoder (which would otherwise replace them with U+FFFD).
func escapeDocInvalidChars(doc *page.PageXmlDocument) {
	escapeTextEquivs := func(tes []page.PageXmlTextEquiv) {
		for i, te := range tes {
			tes[i].Unicode = escapeInvalidChars(te.Unicode)
		}
	}

	escapeTextRegion := func(tr *page.PageXmlTextRegion) {
		escapeTextEquivs(tr.TextEquivs)
		for _, line := range tr.TextLines {
			escapeTextEquivs(line.TextEquivs)
			for _, word := range line.Words {
				escapeTextEquivs(word.TextEquivs)
			}
		}
	}

	if doc.Page != nil {
		for _, item := range doc.Page.Items {
			if tr, ok := item.(*page.PageXmlTextRegion); ok {
				escapeTextRegion(tr)
			}
		}
	}
}

// escapeInvalidChars replaces invalid XML 1.0 characters with a placeholder pattern:
// __INV_XX__ where XX is the hex value of the byte (e.g., __INV_06__ for byte 0x06).
func escapeInvalidChars(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		r := rune(c)
		if !isValidXmlChar(r) {
			b.WriteString(fmt.Sprintf("__INV_%02X__", c))
		} else {
			// Handle multi-byte runes properly
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size <= 1 {
				b.WriteByte(c)
			} else {
				b.WriteRune(r)
			}
			i += size - 1
		}
	}
	return b.String()
}

// restoreEscapedInvalidXmlCharsInString replaces placeholder patterns (__INV_XX__)
// with the original invalid XML character bytes.
func restoreEscapedInvalidXmlCharsInString(s string) string {
	re := regexp.MustCompile(`__INV_([0-9A-F]{2})__`)
	resultBytes := []byte(s)
	fixedBytes := re.ReplaceAllFunc(resultBytes, func(match []byte) []byte {
		hex := string(match[6:8]) // __INV_ is 6 chars, then hex digits at positions 6 and 7
		val, _ := strconv.ParseInt(hex, 16, 8)
		return []byte{byte(val)}
	})
	return string(fixedBytes)
}

// escapeInvalidXmlCharsForParsing replaces raw invalid XML bytes in the file data
// with placeholders so Go's xml.Unmarshal can parse them. postProcessUnmarshal restores them.
func escapeInvalidXmlCharsForParsing(data []byte) []byte {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		r := rune(b)
		if !isValidXmlChar(r) {
			result = append(result, []byte(fmt.Sprintf("__INV_%02X__", b))...)
		} else {
			result = append(result, b)
		}
	}
	return result
}

// serializePageXmlDocument handles the polymorphic Items field correctly.
func serializePageXmlDocument(enc *xml.Encoder, doc *page.PageXmlDocument) error {
	ns := page.PageXmlNamespace

	pcgtsStart := xml.StartElement{
		Name: xml.Name{Space: ns, Local: "PcGts"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "pcGtsId"}, Value: doc.PcGtsId},
		},
	}

	if err := enc.EncodeToken(pcgtsStart); err != nil {
		return err
	}

	// Metadata
	if doc.Metadata != nil {
		metaStart := xml.StartElement{Name: xml.Name{Local: "Metadata"}}
		if err := enc.EncodeToken(metaStart); err != nil {
			return err
		}

		if doc.Metadata.Creator != "" {
			elem := xml.StartElement{Name: xml.Name{Local: "creator"}}
			if err := enc.EncodeElement(doc.Metadata.Creator, elem); err != nil {
				return err
			}
		}
		if !doc.Metadata.Created.IsZero() {
			elem := xml.StartElement{Name: xml.Name{Local: "created"}}
			if err := enc.EncodeElement(doc.Metadata.Created.Format(time.RFC3339Nano), elem); err != nil {
				return err
			}
		}
		if !doc.Metadata.LastChange.IsZero() {
			elem := xml.StartElement{Name: xml.Name{Local: "lastChange"}}
			if err := enc.EncodeElement(doc.Metadata.LastChange.Format(time.RFC3339Nano), elem); err != nil {
				return err
			}
		}
		if doc.Metadata.Comments != "" {
			elem := xml.StartElement{Name: xml.Name{Local: "comments"}}
			if err := enc.EncodeElement(doc.Metadata.Comments, elem); err != nil {
				return err
			}
		}

		if err := enc.EncodeToken(metaStart.End()); err != nil {
			return err
		}
	}

	// Page element with Items handled polymorphically
	if doc.Page != nil {
		pageStart := xml.StartElement{
			Name: xml.Name{Local: "Page"},
			Attr: []xml.Attr{},
		}
		if doc.Page.ImageFilename != "" {
			pageStart.Attr = append(pageStart.Attr, xml.Attr{Name: xml.Name{Local: "imageFilename"}, Value: doc.Page.ImageFilename})
		}
		if doc.Page.ImageWidth != 0 {
			pageStart.Attr = append(pageStart.Attr, xml.Attr{Name: xml.Name{Local: "imageWidth"}, Value: strconv.Itoa(doc.Page.ImageWidth)})
		}
		if doc.Page.ImageHeight != 0 {
			pageStart.Attr = append(pageStart.Attr, xml.Attr{Name: xml.Name{Local: "imageHeight"}, Value: strconv.Itoa(doc.Page.ImageHeight)})
		}

		if err := enc.EncodeToken(pageStart); err != nil {
			return err
		}

		// ReadingOrder
		if doc.Page.ReadingOrder != nil {
			roStart := xml.StartElement{Name: xml.Name{Local: "ReadingOrder"}}
			if err := enc.EncodeToken(roStart); err != nil {
				return err
			}
			if doc.Page.ReadingOrder.Item != nil {
				if og, ok := doc.Page.ReadingOrder.Item.(*page.PageXmlOrderedGroup); ok {
					ogStart := xml.StartElement{
						Name: xml.Name{Local: "OrderedGroup"},
						Attr: []xml.Attr{},
					}
					if og.Id != "" {
						ogStart.Attr = append(ogStart.Attr, xml.Attr{Name: xml.Name{Local: "id"}, Value: og.Id})
					}
					if err := enc.EncodeToken(ogStart); err != nil {
						return err
					}

					for _, refItem := range og.Items {
						if ref, ok := refItem.(*page.PageXmlRegionRefIndexed); ok {
							itemStart := xml.StartElement{
								Name: xml.Name{Local: "Items"},
								Attr: []xml.Attr{},
							}
							itemStart.Attr = append(itemStart.Attr, xml.Attr{Name: xml.Name{Local: "index"}, Value: strconv.Itoa(ref.Index)})
							if ref.RegionRef != "" {
								itemStart.Attr = append(itemStart.Attr, xml.Attr{Name: xml.Name{Local: "regionRef"}, Value: ref.RegionRef})
							}
							if err := enc.EncodeToken(itemStart); err != nil {
								return err
							}
							if err := enc.EncodeToken(itemStart.End()); err != nil {
								return err
							}
						}
					}

					if err := enc.EncodeToken(ogStart.End()); err != nil {
						return err
					}
				}
			}
			if err := enc.EncodeToken(roStart.End()); err != nil {
				return err
			}
		}

		// Items (regions) - encoded with correct element names
		for _, item := range doc.Page.Items {
			name := page.RegionElementName(item)
			itemStart := xml.StartElement{Name: xml.Name{Local: name}}
			if m, ok := any(item).(xml.Marshaler); ok {
				if err := m.MarshalXML(enc, itemStart); err != nil {
					return err
				}
			} else {
				if err := enc.EncodeElement(item, itemStart); err != nil {
					return err
				}
			}
		}

		if err := enc.EncodeToken(pageStart.End()); err != nil {
			return err
		}
	}

	pcgtsEnd := xml.EndElement{Name: xml.Name{Space: ns, Local: "PcGts"}}
	return enc.EncodeToken(pcgtsEnd)
}

// Deserialize reads a PAGE-XML file and returns the parsed PageXmlDocument.
func DeserializePageXml(xmlPath string) (*page.PageXmlDocument, error) {
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML file: %w", err)
	}

	// Escape invalid XML chars before parsing (Go's xml.Decoder rejects them).
	escaped := escapeInvalidXmlCharsForParsing(data)

	doc := &page.PageXmlDocument{}
	if err := xml.Unmarshal(escaped, doc); err != nil {
		return nil, fmt.Errorf("failed to deserialize PAGE-XML document: %w", err)
	}

	// Restore escaped invalid chars in parsed document.
	postProcessUnmarshal(doc)

	return doc, nil
}

// deserializePageXml parses PAGE-XML with proper handling of polymorphic Items.
func deserializePageXml(data []byte) (*page.PageXmlDocument, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))

	doc := &page.PageXmlDocument{}

	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("failed to read token: %w", err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "PcGts":
				for _, attr := range t.Attr {
					if attr.Name.Local == "pcGtsId" {
						doc.PcGtsId = attr.Value
					}
				}
			case "Metadata":
				doc.Metadata = deserializeMetadata(dec)
			case "Page":
				doc.Page = deserializePage(dec, t)
			}
		case xml.EndElement:
			if t.Name.Local == "PcGts" {
				goto done
			}
		}
	}

done:
	return doc, nil
}

func deserializeMetadata(dec *xml.Decoder) *page.PageXmlMetadata {
	meta := &page.PageXmlMetadata{}
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			var val string
			dec.DecodeElement(&val, &t)
			switch t.Name.Local {
			case "creator":
				meta.Creator = val
			case "created":
				if ts, err := time.Parse(time.RFC3339Nano, val); err == nil {
					meta.Created = ts
				}
			case "lastChange":
				if ts, err := time.Parse(time.RFC3339Nano, val); err == nil {
					meta.LastChange = ts
				}
			case "comments":
				meta.Comments = val
			}
		case xml.EndElement:
			if t.Name.Local == "Metadata" {
				return meta
			}
		}
	}
	return meta
}

func deserializePage(dec *xml.Decoder, start xml.StartElement) *page.PageXmlPage {
	pg := &page.PageXmlPage{}
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "imageFilename":
			pg.ImageFilename = attr.Value
		case "imageWidth":
			if v, err := strconv.Atoi(attr.Value); err == nil {
				pg.ImageWidth = v
			}
		case "imageHeight":
			if v, err := strconv.Atoi(attr.Value); err == nil {
				pg.ImageHeight = v
			}
		}
	}

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "ReadingOrder":
				pg.ReadingOrder = deserializeReadingOrder(dec)
			case "TextRegion":
				pg.Items = append(pg.Items, deserializeTextRegion(dec))
			case "ImageRegion":
				pg.Items = append(pg.Items, deserializeImageRegion(dec))
			default:
				dec.Skip()
			}
		case xml.EndElement:
			if t.Name.Local == "Page" {
				return pg
			}
		}
	}
	return pg
}

func deserializeReadingOrder(dec *xml.Decoder) *page.PageXmlReadingOrder {
	ro := &page.PageXmlReadingOrder{}
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "OrderedGroup" {
				og := &page.PageXmlOrderedGroup{Items: make([]any, 0)}
				for _, attr := range t.Attr {
					if attr.Name.Local == "id" {
						og.Id = attr.Value
					}
				}
				for {
					itok, ierr := dec.Token()
					if ierr != nil {
						break
					}
					switch it := itok.(type) {
					case xml.StartElement:
						if it.Name.Local == "Items" {
							ref := &page.PageXmlRegionRefIndexed{}
							for _, attr := range it.Attr {
								switch attr.Name.Local {
								case "index":
									if v, err := strconv.Atoi(attr.Value); err == nil {
										ref.Index = v
									}
								case "regionRef":
									ref.RegionRef = attr.Value
								}
							}
							og.Items = append(og.Items, ref)
							dec.Skip() // skip to end of Items element
						} else {
							dec.Skip()
						}
					case xml.EndElement:
						if it.Name.Local == "OrderedGroup" {
							ro.Item = og
							goto done
						}
					}
				}
			} else {
				dec.Skip()
			}
		case xml.EndElement:
			if t.Name.Local == "ReadingOrder" {
				return ro
			}
		}
	}
done:
	return ro
}

func deserializeTextRegion(dec *xml.Decoder) *page.PageXmlTextRegion {
	// Use default unmarshaling for TextRegion content
	var region page.PageXmlTextRegion
	dec.Decode(&region)
	return &region
}

func deserializeImageRegion(dec *xml.Decoder) *page.PageXmlImageRegion {
	var region page.PageXmlImageRegion
	dec.Decode(&region)
	return &region
}

// pageXmlData tracks counters for generating unique IDs during PAGE-XML export.
type pageXmlData struct {
	LinesCount       int
	WordsCount       int
	GlyphsCount      int
	RegionsCount     int
	GroupOrdersCount int
	OrderedRegions   []*page.PageXmlRegionRefIndexed
}

func newPageXmlData() *pageXmlData {
	return &pageXmlData{
		OrderedRegions: make([]*page.PageXmlRegionRefIndexed, 0),
	}
}

// Compile-time interface assertion.
var _ TextExporter = (*PageXmlTextExporter)(nil)
