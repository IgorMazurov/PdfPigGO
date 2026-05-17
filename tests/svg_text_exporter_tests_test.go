//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/export"
)

func TestSvgDoc68_1990_01_A(t *testing.T) {
	docPath := filepath.Join(integrationDocRoot, "68-1990-01_A.pdf")

	exporter := export.NewSvgTextExporter(export.DoNotCheck)

	if exporter.InvalidCharStrategy != export.DoNotCheck {
		t.Errorf("expected InvalidCharStrategy DoNotCheck, got %v", exporter.InvalidCharStrategy)
	}

	doc, err := pdfpig.OpenFile(docPath, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", docPath, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(7)
	if err != nil {
		t.Fatalf("GetPage(7): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(7): expected *content.Page, got %T", pageAny)
	}

	svgStr := exporter.Get(page)

	svgFile := "68-1990-01_A.7.svg"
	if err := os.WriteFile(svgFile, []byte(svgStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", svgFile, err)
	}
	defer os.Remove(svgFile)

	if len(svgStr) == 0 {
		t.Error("SVG output is empty")
	}
}

func TestSvgIssue655NoCheckStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewSvgTextExporter(export.DoNotCheck)

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

	svgStr := exporter.Get(page)

	svgFile := "issue655.nocheck.svg"
	if err := os.WriteFile(svgFile, []byte(svgStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", svgFile, err)
	}
	defer os.Remove(svgFile)

	rawText, err := os.ReadFile(svgFile)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", svgFile, err)
	}

	expected := "&#x6;"
	if !strings.Contains(string(rawText), expected) {
		t.Errorf("SVG output does not contain %q (no check strategy)", expected)
	}
}

func TestSvgIssue655RemoveStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewSvgTextExporter(export.Remove)

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

	svgStr := exporter.Get(page)

	svgFile := "issue655.remove.svg"
	if err := os.WriteFile(svgFile, []byte(svgStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", svgFile, err)
	}
	defer os.Remove(svgFile)

	rawText, err := os.ReadFile(svgFile)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", svgFile, err)
	}

	unexpected := "0x06"
	if strings.Contains(string(rawText), unexpected) {
		t.Errorf("SVG output should not contain %q (remove strategy)", unexpected)
	}
}

func TestSvgIssue655ConvertToHexadecimalStrategy(t *testing.T) {
	hexPath := filepath.Join(integrationDocRoot, "hex_0x0006.pdf")

	exporter := export.NewSvgTextExporter(export.ConvertToHexadecimal)

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

	svgStr := exporter.Get(page)

	svgFile := "issue655.hex.svg"
	if err := os.WriteFile(svgFile, []byte(svgStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", svgFile, err)
	}
	defer os.Remove(svgFile)

	rawText, err := os.ReadFile(svgFile)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", svgFile, err)
	}

	expected := "0x06"
	if !strings.Contains(string(rawText), expected) {
		t.Errorf("SVG output does not contain %q (hexadecimal strategy)", expected)
	}
}

func TestSvgIssue655CustomStrategy(t *testing.T) {
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

	exporter := export.NewSvgTextExporterWithHandler(customHandler)

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

	svgStr := exporter.Get(page)

	svgFile := "issue655.custom.svg"
	if err := os.WriteFile(svgFile, []byte(svgStr), 0o644); err != nil {
		t.Fatalf("WriteFile(%q): %v", svgFile, err)
	}
	defer os.Remove(svgFile)

	rawText, err := os.ReadFile(svgFile)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", svgFile, err)
	}

	expected := "!?"
	if !strings.Contains(string(rawText), expected) {
		t.Errorf("SVG output does not contain %q (custom strategy)", expected)
	}
}
