//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/uglytoad/pdfpig/go/acroforms"
	"github.com/uglytoad/pdfpig/go/acroforms/fields"
)

const acroFormsBasicFieldsDocName = "AcroFormsBasicFields.pdf"

func resolveAcroFormsBasicFieldsPath(t *testing.T) string {
	t.Helper()

	path := filepath.Join(integrationDocRoot, acroFormsBasicFieldsDocName)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("document not found: %s", path)
	}

	return path
}

func TestTryGetFormNotNull(t *testing.T) {
	fullPath := resolveAcroFormsBasicFieldsPath(t)

	doc, err := pdfpig.OpenFile(fullPath, &pdfpig.LenientParsingOff)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
	}
	defer doc.Close()

	form, ok, err := doc.TryGetForm()
	if err != nil {
		t.Fatalf("TryGetForm: %v", err)
	}

	if !ok || form == nil {
		t.Fatal("expected form to be not nil")
	}
}

func TestTryGetFormDisposedReturnsError(t *testing.T) {
	fullPath := resolveAcroFormsBasicFieldsPath(t)

	doc, err := pdfpig.OpenFile(fullPath, &pdfpig.LenientParsingOff)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
	}

	doc.Close()

	_, _, err = doc.TryGetForm()
	if err == nil {
		t.Fatal("expected error when calling TryGetForm on disposed document")
	}
}

func TestTryGetGetsAllFormFields(t *testing.T) {
	fullPath := resolveAcroFormsBasicFieldsPath(t)

	doc, err := pdfpig.OpenFile(fullPath, &pdfpig.LenientParsingOff)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
	}
	defer doc.Close()

	form, ok, err := doc.TryGetForm()
	if err != nil {
		t.Fatalf("TryGetForm: %v", err)
	}

	if !ok || form == nil {
		t.Fatal("expected form to be not nil")
	}

	expectedCount := 18
	if len(form.Fields()) != expectedCount {
		t.Errorf("expected %d fields, got %d", expectedCount, len(form.Fields()))
	}
}

func TestTryGetFormFieldsByPage(t *testing.T) {
	fullPath := resolveAcroFormsBasicFieldsPath(t)

	doc, err := pdfpig.OpenFile(fullPath, &pdfpig.LenientParsingOff)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
	}
	defer doc.Close()

	form, ok, err := doc.TryGetForm()
	if err != nil {
		t.Fatalf("TryGetForm: %v", err)
	}

	if !ok || form == nil {
		t.Fatal("expected form to be not nil")
	}

	pageFields, err := form.GetFieldsForPage(1)
	if err != nil {
		t.Fatalf("GetFieldsForPage(1): %v", err)
	}

	expectedCount := 18
	if len(pageFields) != expectedCount {
		t.Errorf("expected %d fields on page 1, got %d", expectedCount, len(pageFields))
	}
}

func TestTryGetGetsRadioButtonState(t *testing.T) {
	fullPath := resolveAcroFormsBasicFieldsPath(t)

	doc, err := pdfpig.OpenFile(fullPath, &pdfpig.LenientParsingOff)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
	}
	defer doc.Close()

	form, ok, err := doc.TryGetForm()
	if err != nil {
		t.Fatalf("TryGetForm: %v", err)
	}

	if !ok || form == nil {
		t.Fatal("expected form to be not nil")
	}

	var radioButtons []*fields.AcroRadioButtonsField
	for _, f := range form.Fields() {
		if rb, ok := any(f).(*fields.AcroRadioButtonsField); ok {
			radioButtons = append(radioButtons, rb)
		}
	}

	expectedCount := 2
	if len(radioButtons) != expectedCount {
		t.Fatalf("expected %d radio button groups, got %d", expectedCount, len(radioButtons))
	}

	sort.Slice(radioButtons, func(i, j int) bool {
		minLeft := func(rb *fields.AcroRadioButtonsField) float64 {
			children := rb.Children()
			if len(children) == 0 {
				return 0
			}
			firstChild := acroforms.ToAcroFieldBaseFromAny(children[0])
			if firstChild == nil || firstChild.GetBounds() == nil {
				return 0
			}
			left := firstChild.GetBounds().Left()
			for k := 1; k < len(children); k++ {
				childK := acroforms.ToAcroFieldBaseFromAny(children[k])
				if childK != nil && childK.GetBounds() != nil && childK.GetBounds().Left() < left {
					left = childK.GetBounds().Left()
				}
			}
			return left
		}
		return minLeft(radioButtons[i]) < minLeft(radioButtons[j])
	})

	left := radioButtons[0]

	if len(left.Children()) != 2 {
		t.Errorf("expected left radio group to have 2 children, got %d", len(left.Children()))
	}

	for _, child := range left.Children() {
		button, ok := any(child).(*fields.AcroRadioButtonField)
		if !ok {
			t.Errorf("expected *AcroRadioButtonField, got %T", child)
			continue
		}
		if button.IsSelected {
			t.Error("expected left radio group children to not be selected")
		}
	}

	right := radioButtons[1]

	if len(right.Children()) != 2 {
		t.Errorf("expected right radio group to have 2 children, got %d", len(right.Children()))
	}

	buttonOn, ok := any(right.Children()[0]).(*fields.AcroRadioButtonField)
	if !ok {
		t.Fatalf("expected *AcroRadioButtonField for first child of right group, got %T", right.Children()[0])
	}
	if !buttonOn.IsSelected {
		t.Error("expected first child of right radio group to be selected")
	}

	buttonOff, ok := any(right.Children()[1]).(*fields.AcroRadioButtonField)
	if !ok {
		t.Fatalf("expected *AcroRadioButtonField for second child of right group, got %T", right.Children()[1])
	}
	if buttonOff.IsSelected {
		t.Error("expected second child of right radio group to not be selected")
	}
}
