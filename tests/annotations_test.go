//go:build integration

package pdfpig_test


import (
	"github.com/uglytoad/pdfpig/go/testutil"
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"slices"
	"testing"

	"github.com/uglytoad/pdfpig/go/actions"
	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/content"
)

var specificTestDocRoot = testutil.SpecificTestDocumentsRoot

func getTocPath() string {
	return filepath.Join(integrationDocRoot, "toc.pdf")
}

func getAppearancesPath() string {
	return filepath.Join(specificTestDocRoot, "appearances.pdf")
}

// TestAnnotationsHaveActions verifies that the toc document's first page contains
// 5 annotations, each with a GoToAction pointing to a valid page number.
// This matches C# AnnotationsTest.AnnotationsHaveActions.
func TestAnnotationsHaveActions(t *testing.T) {
	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(getTocPath(), opts)
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

	annots := page.GetAnnotations()
	if annots == nil {
		t.Fatal("GetAnnotations returned nil")
	}

	if len(annots) != 5 {
		t.Errorf("expected 5 annotations, got %d", len(annots))
		return
	}

	for i, a := range annots {
		action := a.Action()
		if action == nil {
			t.Errorf("annotation[%d]: Action is nil", i)
			continue
		}

		gotoAction, ok := action.(*actions.GoToAction)
		if !ok {
			t.Errorf("annotation[%d]: expected *actions.GoToAction, got %T", i, action)
			continue
		}

		if gotoAction.Destination.PageNumber <= 0 {
			t.Errorf("annotation[%d]: GoToAction Destination.PageNumber = %d, expected > 0", i, gotoAction.Destination.PageNumber)
		}
	}
}

// TestCheckAnnotationAppearanceStreams verifies that the appearances document's
// annotation has correct appearance stream properties: HasDownAppearance=true,
// HasNormalAppearance=true, HasRollOverAppearance=false, both non-stateless streams
// contain "Off" and "Yes" states, and the current state is "Off".
// This matches C# AnnotationsTest.CheckAnnotationAppearanceStreams.
func TestCheckAnnotationAppearanceStreams(t *testing.T) {
	opts := &content.ParsingOptions{
		UseLenientParsing: false,
	}

	doc, err := pdfpig.OpenFile(getAppearancesPath(), opts)
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

	annots := page.GetAnnotations()
	if annots == nil {
		t.Fatal("GetAnnotations returned nil")
	}

	if len(annots) != 1 {
		t.Fatalf("expected exactly 1 annotation, got %d", len(annots))
	}

	annotation := annots[0]

	if !annotation.HasDownAppearance() {
		t.Error("expected HasDownAppearance to be true")
	}

	if !annotation.HasNormalAppearance() {
		t.Error("expected HasNormalAppearance to be true")
	}

	if annotation.HasRollOverAppearance() {
		t.Error("expected HasRollOverAppearance to be false")
	}

	downStreamAny := annotation.DownAppearanceStream()
	if downStreamAny == nil {
		t.Fatal("DownAppearanceStream is nil")
	}

	downStream, ok := downStreamAny.(*annotations.AppearanceStream)
	if !ok {
		t.Fatalf("DownAppearanceStream: expected *annotations.AppearanceStream, got %T", downStreamAny)
	}

	if downStream.IsStateless() {
		t.Error("expected downAppearanceStream.IsStateless to be false")
	}

	downStates := downStream.States()
	if !slices.Contains(downStates, "Off") {
		t.Error("expected downAppearanceStream states to contain 'Off'")
	}
	if !slices.Contains(downStates, "Yes") {
		t.Error("expected downAppearanceStream states to contain 'Yes'")
	}

	normalStreamAny := annotation.NormalAppearanceStream()
	if normalStreamAny == nil {
		t.Fatal("NormalAppearanceStream is nil")
	}

	normalStream, ok := normalStreamAny.(*annotations.AppearanceStream)
	if !ok {
		t.Fatalf("NormalAppearanceStream: expected *annotations.AppearanceStream, got %T", normalStreamAny)
	}

	if normalStream.IsStateless() {
		t.Error("expected normalAppearanceStream.IsStateless to be false")
	}

	normalStates := normalStream.States()
	if !slices.Contains(normalStates, "Off") {
		t.Error("expected normalAppearanceStream states to contain 'Off'")
	}
	if !slices.Contains(normalStates, "Yes") {
		t.Error("expected normalAppearanceStream states to contain 'Yes'")
	}

	appearanceState := annotation.AppearanceState()
	if appearanceState == nil {
		t.Fatal("AppearanceState is nil")
	} else if *appearanceState != "Off" {
		t.Errorf("expected appearanceState to be 'Off', got %q", *appearanceState)
	}
}
