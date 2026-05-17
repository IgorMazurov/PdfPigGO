//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/content"
)

func getAnnotationCommentsPath() string {
	return filepath.Join(integrationDocRoot, "annotation-comments.pdf")
}

// TestHasCorrectNumberOfAnnotations verifies that the annotation-comments document
// contains exactly 4 annotations with types Text, Popup, Text, Popup.
// This matches C# AnnotationReplyToTests.HasCorrectNumberOfAnnotations.
func TestHasCorrectNumberOfAnnotations(t *testing.T) {
	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(getAnnotationCommentsPath(), opts)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	linkAnnotations := page.GetAnnotations()
	if linkAnnotations == nil {
		t.Fatal("GetAnnotations returned nil")
	}

	if len(linkAnnotations) != 4 {
		t.Errorf("expected 4 annotations, got %d", len(linkAnnotations))
		return
	}

	expectedTypes := []annotations.AnnotationType{
		annotations.Text,
		annotations.Popup,
		annotations.Text,
		annotations.Popup,
	}

	for i, expected := range expectedTypes {
		gotAny := linkAnnotations[i].Type()
		if gotAny == nil {
			t.Errorf("annotation[%d]: Type is nil", i)
			continue
		}

		gotType, ok := gotAny.(annotations.AnnotationType)
		if !ok {
			t.Errorf("annotation[%d]: expected annotations.AnnotationType, got %T", i, gotAny)
			continue
		}

		if gotType != expected {
			t.Errorf("annotation[%d]: expected type %v, got %v", i, expected, gotType)
		}
	}
}

// TestSecondTextReplyToFirst verifies that the third annotation (index 2) has its
// InReplyTo field pointing to the first annotation (index 0).
// This matches C# AnnotationReplyToTests.SecondTextReplyToFirst.
func TestSecondTextReplyToFirst(t *testing.T) {
	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(getAnnotationCommentsPath(), opts)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
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

	linkAnnotations := page.GetAnnotations()
	if linkAnnotations == nil {
		t.Fatal("GetAnnotations returned nil")
	}

	if len(linkAnnotations) < 3 {
		t.Fatalf("expected at least 3 annotations for InReplyTo test, got %d", len(linkAnnotations))
	}

	inReplyTo := linkAnnotations[2].InReplyTo()
	if inReplyTo == nil {
		t.Fatal("annotations[2].InReplyTo is nil")
	}

	if linkAnnotations[0] != inReplyTo {
		t.Errorf("expected annotations[0] to equal annotations[2].InReplyTo")
	}
}
