//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/export"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
)

func TestHOcrIssue655NoCheckStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewHOcrTextExporter(
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

	hocrStr := exporter.Get(page)

	expected := "1\u00062345\u0006678\u0006ABC"
	if !strings.Contains(hocrStr, expected) {
		t.Errorf("hOCR output does not contain expected text %q (no check strategy, contains invalid xml chars)", expected)
	}
}

func TestHOcrIssue655RemoveStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewHOcrTextExporter(
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

	hocrStr := exporter.Get(page)

	expected := "12345678ABC"
	if !strings.Contains(hocrStr, expected) {
		t.Errorf("hOCR output does not contain expected text %q (remove strategy)", expected)
	}
}

func TestHOcrIssue655ConvertToHexadecimalStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewHOcrTextExporter(
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

	hocrStr := exporter.Get(page)

	expected := "10x0623450x066780x06ABC"
	if !strings.Contains(hocrStr, expected) {
		t.Errorf("hOCR output does not contain expected text %q (hexadecimal strategy)", expected)
	}
}

func TestHOcrIssue655CustomStrategy(t *testing.T) {
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

	exporter := export.NewHOcrTextExporterWithHandler(
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

	hocrStr := exporter.Get(page)

	expected := "1!?2345!?678!?ABC"
	if !strings.Contains(hocrStr, expected) {
		t.Errorf("hOCR output does not contain expected text %q (custom strategy)", expected)
	}
}
