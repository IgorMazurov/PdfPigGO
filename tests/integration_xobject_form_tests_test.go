//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
)

// TestCanReadDocumentWithoutStackOverflowIssue671 verifies that a document with a
// self-referencing XObject form can be read without stack overflow when using
// lenient parsing (the default). This matches C# XObjectFormTests.CanReadDocumentWithoutStackOverflowIssue671.
func TestCanReadDocumentWithoutStackOverflowIssue671(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "issue_671.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}
	_ = page // successfully retrieved without stack overflow
}

// TestCanReadDocumentThrowsIssue671 verifies that a document with a self-referencing
// XObject form causes an error when lenient parsing is disabled. This matches
// C# XObjectFormTests.CanReadDocumentThrowsIssue671.
func TestCanReadDocumentThrowsIssue671(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "issue_671.pdf")

	doc, err := pdfpig.OpenFile(path, &pdfpig.LenientParsingOff)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, err = doc.GetPage(1)
	if err == nil {
		t.Fatal("expected error for self-referencing XObject form with lenient parsing off")
	}

	if !strings.Contains(err.Error(), "is referencing itself which can cause unexpected behaviour") {
		t.Errorf("error message does not contain expected text; got: %s", err.Error())
	}
}

// TestCanReadDocumentMOZILLA3136_0 verifies that a document without circular references
// reads successfully even with lenient parsing disabled. This matches
// C# XObjectFormTests.CanReadDocumentMOZILLA_3136_0.
func TestCanReadDocumentMOZILLA3136_0(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "MOZILLA-3136-0.pdf")

	doc, err := pdfpig.OpenFile(path, &pdfpig.LenientParsingOff)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}
	_ = page // successfully retrieved without error
}

// TestXObjectFormXClipping verifies that XObject form clipping produces the expected
// number of paths and text content. This matches C# XObjectFormTests.XObjectFormXClipping.
func TestXObjectFormXClipping(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "ICML03-081.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{ClipPaths: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(4)
	if err != nil {
		t.Fatalf("GetPage(4): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	if len(page.Paths()) <= 3 {
		t.Errorf("expected more than 3 paths, got %d", len(page.Paths()))
	}

	extractor := word_extractor.NewNearestNeighbourWordExtractor()
	words := extractor.GetWords(page.Letters())

	segmenter := page_segmenter.Instance
	blocks := segmenter.GetBlocks(words)

	// Flatten all text lines from blocks.
	var allLines []*dla.TextLine
	for _, block := range blocks {
		allLines = append(allLines, block.TextLines...)
	}

	trainingExamplesCount := 0
	classificationWeightCount := 0
	for _, line := range allLines {
		if line.Text == "Training Examples per Class" {
			trainingExamplesCount++
		}
		if line.Text == "Classification Weight" {
			classificationWeightCount++
		}
	}

	if trainingExamplesCount != 2 {
		t.Errorf("expected 2 lines with 'Training Examples per Class', got %d", trainingExamplesCount)
	}
	if classificationWeightCount != 2 {
		t.Errorf("expected 2 lines with 'Classification Weight', got %d", classificationWeightCount)
	}
}
