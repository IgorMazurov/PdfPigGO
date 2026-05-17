//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testutil"
)

func TestNoMarkedOptionalContent(t *testing.T) {
	path := testutil.GetDocumentPath("AcroFormsBasicFields.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
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
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	oc := page.GetOptionalContents()

	if len(oc) != 0 {
		t.Errorf("expected no optional contents, got %d groups", len(oc))
	}
}

func TestMarkedOptionalContent(t *testing.T) {
	path := testutil.GetDocumentPath("odwriteex.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
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
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	oc := page.GetOptionalContents()

	if len(oc) != 3 {
		t.Errorf("expected 3 optional content groups, got %d", len(oc))
	}

	expectedKeys := []string{"0", "Dimentions", "Text"}
	for _, key := range expectedKeys {
		if _, exists := oc[key]; !exists {
			t.Errorf("missing expected optional content group: %q", key)
		}
	}

	if elems, exists := oc["0"]; !exists || len(elems) != 1 {
		t.Errorf("expected oc[\"0\"] to have 1 element, got %d", lenSafe(oc, "0"))
	}

	if elems, exists := oc["Dimentions"]; !exists || len(elems) != 2 {
		t.Errorf("expected oc[\"Dimentions\"] to have 2 elements, got %d", lenSafe(oc, "Dimentions"))
	}

	if elems, exists := oc["Text"]; !exists || len(elems) != 1 {
		t.Errorf("expected oc[\"Text\"] to have 1 element, got %d", lenSafe(oc, "Text"))
	}
}

func TestMarkedOptionalContentRecursion(t *testing.T) {
	path := testutil.GetDocumentPath("Layer pdf - 322_High_Holborn_building_Brochure.pdf", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	page1Any, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page1, ok := page1Any.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", page1Any)
	}

	oc1 := page1.GetOptionalContents()

	if len(oc1) != 16 {
		t.Errorf("page 1: expected 16 optional content groups, got %d", len(oc1))
	}

	if _, exists := oc1["NEW ARRANGEMENT"]; !exists {
		t.Error("page 1: missing expected group \"NEW ARRANGEMENT\"")
	}

	page2Any, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}

	page2, ok := page2Any.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(2): expected *content.Page, got %T", page2Any)
	}

	oc2 := page2.GetOptionalContents()

	if len(oc2) != 15 {
		t.Errorf("page 2: expected 15 optional content groups, got %d", len(oc2))
	}

	if _, exists := oc2["NEW ARRANGEMENT"]; exists {
		t.Error("page 2: should not contain group \"NEW ARRANGEMENT\"")
	}

	if _, exists := oc2["WDL Shell text"]; !exists {
		t.Error("page 2: missing expected group \"WDL Shell text\"")
	}

	if elems, exists := oc2["WDL Shell text"]; !exists || len(elems) != 2 {
		t.Errorf("page 2: expected oc[\"WDL Shell text\"] to have 2 elements, got %d", lenSafe(oc2, "WDL Shell text"))
	}

	page3Any, err := doc.GetPage(3)
	if err != nil {
		t.Fatalf("GetPage(3): %v", err)
	}

	page3, ok := page3Any.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(3): expected *content.Page, got %T", page3Any)
	}

	oc3 := page3.GetOptionalContents()

	if len(oc3) != 15 {
		t.Errorf("page 3: expected 15 optional content groups, got %d", len(oc3))
	}

	if _, exists := oc3["WDL Shell text"]; !exists {
		t.Error("page 3: missing expected group \"WDL Shell text\"")
	}

	if elems, exists := oc3["WDL Shell text"]; !exists || len(elems) != 2 {
		t.Errorf("page 3: expected oc[\"WDL Shell text\"] to have 2 elements, got %d", lenSafe(oc3, "WDL Shell text"))
	}
}

func lenSafe(oc map[string][]*content.OptionalContentGroupElement, key string) int {
	elems, exists := oc[key]
	if !exists {
		return -1
	}
	return len(elems)
}
