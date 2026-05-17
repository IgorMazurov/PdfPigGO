//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document"
	"github.com/uglytoad/pdfpig/go/testutil"
)

// TestEmbeddedFileAttachmentHasCorrectText verifies that every page of the embedded
// file attachment test document starts with the expected text prefix. This matches
// C# EmbeddedFileAttachmentTests.HasCorrectText.
func TestEmbeddedFileAttachmentHasCorrectText(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("embedded-file-attachment.pdf", true)

	opts := &content.ParsingOptions{}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 1; i <= doc.NumberOfPages(); i++ {
		pageAny, err := doc.GetPage(i)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", i, pageAny)
		}

		expectedPrefix := "This is a test document. It contains a file attachment."
		text := page.Text()
		if len(text) < len(expectedPrefix) || text[:len(expectedPrefix)] != expectedPrefix {
			t.Errorf("page %d: text does not start with %q, got %q", i, expectedPrefix, text)
		}
	}
}

// TestEmbeddedFileAttachmentHasEmbeddedFiles verifies that the document contains
// exactly one embedded file of the expected size. This matches C#
// EmbeddedFileAttachmentTests.HasEmbeddedFiles.
func TestEmbeddedFileAttachmentHasEmbeddedFiles(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("embedded-file-attachment.pdf", true)

	opts := &content.ParsingOptions{}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	adv, ok := doc.Advanced.(*document.AdvancedPdfDocumentAccess)
	if !ok || adv == nil {
		t.Fatal("expected Advanced to be *document.AdvancedPdfDocumentAccess")
	}

	files, found, err := adv.TryGetEmbeddedFiles()
	if err != nil {
		t.Fatalf("TryGetEmbeddedFiles: %v", err)
	}

	if !found {
		t.Fatal("expected document to contain embedded files")
	}

	if len(files) != 1 {
		t.Errorf("expected 1 embedded file, got %d", len(files))
	}

	if len(files) > 0 && len(files[0].Memory) != 20668 {
		t.Errorf("expected embedded file memory length 20668, got %d", len(files[0].Memory))
	}
}
