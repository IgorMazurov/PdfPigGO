//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/export"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
)

func TestIssue655NoCheckStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewAltoXmlTextExporter(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
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

	xmlFile := "issue655.nocheck.altoxml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.Deserialize(xmlFile)
	if err != nil {
		t.Fatalf("Deserialize(%q): %v", xmlFile, err)
	}

	textBlocks := pageXml.Layout.Pages[0].PrintSpace.TextBlocks
	if len(textBlocks) != 1 {
		t.Fatalf("expected 1 text block, got %d", len(textBlocks))
	}

	textLines := textBlocks[0].TextLines
	if len(textLines) != 1 {
		t.Fatalf("expected 1 text line, got %d", len(textLines))
	}

	strings := textLines[0].Strings
	if len(strings) != 2 {
		t.Fatalf("expected 2 strings, got %d", len(strings))
	}

	if strings[0].Content != "TM" {
		t.Errorf("strings[0].Content = %q, want %q", strings[0].Content, "TM")
	}

	expected := "1\u00062345\u0006678\u0006ABC"
	if strings[1].Content != expected {
		t.Errorf("strings[1].Content = %q, want %q (no check strategy, contains invalid xml chars)", strings[1].Content, expected)
	}
}

func TestIssue655RemoveStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewAltoXmlTextExporter(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
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

	xmlFile := "issue655.remove.altoxml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.Deserialize(xmlFile)
	if err != nil {
		t.Fatalf("Deserialize(%q): %v", xmlFile, err)
	}

	textBlocks := pageXml.Layout.Pages[0].PrintSpace.TextBlocks
	if len(textBlocks) != 1 {
		t.Fatalf("expected 1 text block, got %d", len(textBlocks))
	}

	textLines := textBlocks[0].TextLines
	if len(textLines) != 1 {
		t.Fatalf("expected 1 text line, got %d", len(textLines))
	}

	strings := textLines[0].Strings
	if len(strings) != 2 {
		t.Fatalf("expected 2 strings, got %d", len(strings))
	}

	if strings[0].Content != "TM" {
		t.Errorf("strings[0].Content = %q, want %q", strings[0].Content, "TM")
	}

	expected := "12345678ABC"
	if strings[1].Content != expected {
		t.Errorf("strings[1].Content = %q, want %q", strings[1].Content, expected)
	}
}

func TestIssue655ConvertToHexadecimalStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewAltoXmlTextExporter(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
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

	xmlFile := "issue655.hex.altoxml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.Deserialize(xmlFile)
	if err != nil {
		t.Fatalf("Deserialize(%q): %v", xmlFile, err)
	}

	textBlocks := pageXml.Layout.Pages[0].PrintSpace.TextBlocks
	if len(textBlocks) != 1 {
		t.Fatalf("expected 1 text block, got %d", len(textBlocks))
	}

	textLines := textBlocks[0].TextLines
	if len(textLines) != 1 {
		t.Fatalf("expected 1 text line, got %d", len(textLines))
	}

	strings := textLines[0].Strings
	if len(strings) != 2 {
		t.Fatalf("expected 2 strings, got %d", len(strings))
	}

	if strings[0].Content != "TM" {
		t.Errorf("strings[0].Content = %q, want %q", strings[0].Content, "TM")
	}

	expected := "10x0623450x066780x06ABC"
	if strings[1].Content != expected {
		t.Errorf("strings[1].Content = %q, want %q", strings[1].Content, expected)
	}
}

func TestIssue655CustomStrategy(t *testing.T) {
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

	exporter := export.NewAltoXmlTextExporterWithHandler(
		content.Instance,
		page_segmenter.DefaultRecursiveXYCutInstance,
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

	xmlFile := "issue655.custom.altoxml.xml"
	if err := os.WriteFile(xmlFile, []byte(xmlStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", xmlFile, err)
	}
	defer os.Remove(xmlFile)

	pageXml, err := export.Deserialize(xmlFile)
	if err != nil {
		t.Fatalf("Deserialize(%q): %v", xmlFile, err)
	}

	textBlocks := pageXml.Layout.Pages[0].PrintSpace.TextBlocks
	if len(textBlocks) != 1 {
		t.Fatalf("expected 1 text block, got %d", len(textBlocks))
	}

	textLines := textBlocks[0].TextLines
	if len(textLines) != 1 {
		t.Fatalf("expected 1 text line, got %d", len(textLines))
	}

	strings := textLines[0].Strings
	if len(strings) != 2 {
		t.Fatalf("expected 2 strings, got %d", len(strings))
	}

	if strings[0].Content != "TM" {
		t.Errorf("strings[0].Content = %q, want %q", strings[0].Content, "TM")
	}

	expected := "1!?2345!?678!?ABC"
	if strings[1].Content != expected {
		t.Errorf("strings[1].Content = %q, want %q", strings[1].Content, expected)
	}
}

// isXmlChar checks whether the given rune is a valid XML 1.0 character,
// matching the behavior of C# XmlConvert.IsXmlChar.
func isXmlChar(r rune) bool {
	return r == 0x0009 ||
		r == 0x000A ||
		r == 0x000D ||
		(r >= 0x0020 && r <= 0xD7FF) ||
		(r >= 0xE000 && r <= 0xFFFD) ||
		(r >= 0x10000 && r <= 0x10FFFF)
}
