//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/export"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/export/page"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
)

func TestWhenReadingOrderContainsReadingOrderXmlElements(t *testing.T) {
	exporter := export.NewPageXmlTextExporter(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
		reading_order_detector.UnsupervisedInstance,
		1.0,
		"\t",
		export.DoNotCheck,
	)

	xmlStr := getPageXml(exporter)

	if !strings.Contains(xmlStr, "<ReadingOrder>") {
		t.Error("XML does not contain <ReadingOrder>")
	}
	if !strings.Contains(xmlStr, "</OrderedGroup>") {
		t.Error("XML does not contain </OrderedGroup>")
	}
}

func TestPageHeightAndWidthArePresent(t *testing.T) {
	xmlStr := getPageXml(nil)

	expected := `<Page imageFilename="unknown" imageWidth="595" imageHeight="842">`
	if !strings.Contains(xmlStr, expected) {
		t.Errorf("XML does not contain expected page dimensions")
	}
}

func TestContainsExpectedNumberOfTextRegions(t *testing.T) {
	xmlStr := getPageXml(nil)

	re := regexp.MustCompile(`</TextRegion>`)
	count := len(re.FindAllString(xmlStr, -1))

	if count != 22 {
		t.Errorf("expected 22 TextRegion elements, got %d", count)
	}
}

func TestContainsExpectedText(t *testing.T) {
	xmlStr := getPageXml(nil)

	if !strings.Contains(xmlStr, "2006 Swedish Touring Car Championship") {
		t.Error("XML does not contain expected text '2006 Swedish Touring Car Championship'")
	}

	expectedCoords := `<Coords points="35,77 35,62 397,62 397,77" />`
	if !strings.Contains(xmlStr, expectedCoords) {
		t.Errorf("XML does not contain expected coords %q", expectedCoords)
	}
}

func TestNoPointsAreOnThePageBoundary(t *testing.T) {
	pageWidth := float64(100)
	pageHeight := float64(200)

	tests := []struct {
		name     string
		point    core.PdfPoint
		expected string
	}{
		{
			name:     "topLeftPagePoint",
			point:    core.NewPdfPoint(0, 0),
			expected: "1,199",
		},
		{
			name:     "bottomLeftPagePoint",
			point:    core.NewPdfPoint(0, pageHeight),
			expected: "1,1",
		},
		{
			name:     "bottomRightPagePoint",
			point:    core.NewPdfPoint(pageWidth, pageHeight),
			expected: "99,1",
		},
		{
			name:     "normalPoint",
			point:    core.NewPdfPoint(60, 60),
			expected: "60,140",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := export.PointToString(tc.point, pageWidth, pageHeight, 1.0)
			if result != tc.expected {
				t.Errorf("PointToString(%v, %f, %f, 1.0) = %q, want %q", tc.point, pageWidth, pageHeight, result, tc.expected)
			}
		})
	}
}

func TestPageXmlIssue655NoCheckStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewPageXmlTextExporter(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
		reading_order_detector.UnsupervisedInstance,
		1.0,
		"\t",
		export.DoNotCheck,
	)

	if exporter.InvalidCharStrategy != export.DoNotCheck {
		t.Errorf("expected InvalidCharStrategy DoNotCheck, got %v", exporter.InvalidCharStrategy)
	}

	doc, err := pdfpig.OpenFile(hexPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", hexPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	xmlStr := exporter.Get(page)

	xmlFile := "issue655.nocheck.pagexml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.DeserializePageXml(xmlFile)
	if err != nil {
		t.Fatalf("DeserializePageXml(%q): %v", xmlFile, err)
	}

	textRegions := filterTextRegions(pageXml.Page.Items)
	if len(textRegions) != 1 {
		t.Fatalf("expected 1 text region, got %d", len(textRegions))
	}

	textEquivs := textRegions[0].TextEquivs
	if len(textEquivs) != 1 {
		t.Fatalf("expected 1 TextEquiv, got %d", len(textEquivs))
	}

	unicode := textEquivs[0].Unicode
	expected := "TM 1\u00062345\u0006678\u0006ABC"
	if unicode != expected {
		t.Errorf("unicode = %q, want %q (no check strategy, contains invalid xml chars)", unicode, expected)
	}
}

func TestPageXmlIssue655RemoveStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewPageXmlTextExporter(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
		reading_order_detector.UnsupervisedInstance,
		1.0,
		"\t",
		export.Remove,
	)

	if exporter.InvalidCharStrategy != export.Remove {
		t.Errorf("expected InvalidCharStrategy Remove, got %v", exporter.InvalidCharStrategy)
	}

	doc, err := pdfpig.OpenFile(hexPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", hexPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	xmlStr := exporter.Get(page)

	xmlFile := "issue655.remove.pagexml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.DeserializePageXml(xmlFile)
	if err != nil {
		t.Fatalf("DeserializePageXml(%q): %v", xmlFile, err)
	}

	textRegions := filterTextRegions(pageXml.Page.Items)
	if len(textRegions) != 1 {
		t.Fatalf("expected 1 text region, got %d", len(textRegions))
	}

	textEquivs := textRegions[0].TextEquivs
	if len(textEquivs) != 1 {
		t.Fatalf("expected 1 TextEquiv, got %d", len(textEquivs))
	}

	unicode := textEquivs[0].Unicode
	expected := "TM 12345678ABC"
	if unicode != expected {
		t.Errorf("unicode = %q, want %q (remove strategy)", unicode, expected)
	}
}

func TestPageXmlIssue655ConvertToHexadecimalStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewPageXmlTextExporter(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
		reading_order_detector.UnsupervisedInstance,
		1.0,
		"\t",
		export.ConvertToHexadecimal,
	)

	if exporter.InvalidCharStrategy != export.ConvertToHexadecimal {
		t.Errorf("expected InvalidCharStrategy ConvertToHexadecimal, got %v", exporter.InvalidCharStrategy)
	}

	doc, err := pdfpig.OpenFile(hexPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", hexPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	xmlStr := exporter.Get(page)

	xmlFile := "issue655.hex.pagexml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.DeserializePageXml(xmlFile)
	if err != nil {
		t.Fatalf("DeserializePageXml(%q): %v", xmlFile, err)
	}

	textRegions := filterTextRegions(pageXml.Page.Items)
	if len(textRegions) != 1 {
		t.Fatalf("expected 1 text region, got %d", len(textRegions))
	}

	textEquivs := textRegions[0].TextEquivs
	if len(textEquivs) != 1 {
		t.Fatalf("expected 1 TextEquiv, got %d", len(textEquivs))
	}

	unicode := textEquivs[0].Unicode
	expected := "TM 10x0623450x066780x06ABC"
	if unicode != expected {
		t.Errorf("unicode = %q, want %q (hexadecimal strategy)", unicode, expected)
	}
}

func TestPageXmlIssue655CustomStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	customHandler := func(s string) string {
		if s == "" {
			return s
		}

		var b []rune
		for _, r := range s {
			if isXmlChar(r) {
				b = append(b, r)
			} else {
				b = append(b, '!', '?')
			}
		}
		return string(b)
	}

	exporter := export.NewPageXmlTextExporterWithHandler(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
		reading_order_detector.UnsupervisedInstance,
		1.0,
		"\t",
		customHandler,
	)

	if exporter.InvalidCharStrategy != export.Custom {
		t.Errorf("expected InvalidCharStrategy Custom, got %v", exporter.InvalidCharStrategy)
	}

	doc, err := pdfpig.OpenFile(hexPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", hexPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	xmlStr := exporter.Get(page)

	xmlFile := "issue655.custom.pagexml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.DeserializePageXml(xmlFile)
	if err != nil {
		t.Fatalf("DeserializePageXml(%q): %v", xmlFile, err)
	}

	textRegions := filterTextRegions(pageXml.Page.Items)
	if len(textRegions) != 1 {
		t.Fatalf("expected 1 text region, got %d", len(textRegions))
	}

	textEquivs := textRegions[0].TextEquivs
	if len(textEquivs) != 1 {
		t.Fatalf("expected 1 TextEquiv, got %d", len(textEquivs))
	}

	unicode := textEquivs[0].Unicode
	expected := "TM 1!?2345!?678!?ABC"
	if unicode != expected {
		t.Errorf("unicode = %q, want %q (custom strategy)", unicode, expected)
	}
}

// getPageXml opens the Swedish Touring Car Championship PDF and exports page 1 as PAGE-XML.
func getPageXml(custom *export.PageXmlTextExporter) string {
	if custom == nil {
		custom = export.NewPageXmlTextExporter(
			content.Instance,
			page_segmenter.DefaultRecursiveXYCutInstance,
			reading_order_detector.UnsupervisedInstance,
			1.0,
			"\t",
			export.DoNotCheck,
		)
	}

	docPath := filepath.Join(integrationDocRoot, "2006_Swedish_Touring_Car_Championship.pdf")

	doc, err := pdfpig.OpenFile(docPath, nil)
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		panic(err)
	}

	pageObj, ok := pageAny.(*content.Page)
	if !ok {
		panic("expected *content.Page")
	}

	return custom.Get(pageObj)
}

// filterTextRegions filters PageXmlRegionItem slice to only *PageXmlTextRegion items.
func filterTextRegions(items []page.PageXmlRegionItem) []*page.PageXmlTextRegion {
	var result []*page.PageXmlTextRegion
	for _, item := range items {
		if tr, ok := item.(*page.PageXmlTextRegion); ok {
			result = append(result, tr)
		}
	}
	return result
}
