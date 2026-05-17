//go:build integration

package pdfpig_test

import (
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/testpages"
	"github.com/uglytoad/pdfpig/go/testutil"
)


func TestSimpleFactory1(t *testing.T) {
	path := testutil.GetDocumentPath("ICML03-081", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	doc.AddPageFactory(testpages.NewSimplePageFactory(nil, nil, nil, nil, &content.ParsingOptions{}))

	for p := 1; p < doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", p, pageAny)
		}

		simplePage, err := content.GetTypedPage[testpages.SimplePage](doc, p)
		if err != nil {
			t.Fatalf("GetPage[*testpages.SimplePage](%d): %v", p, err)
		}

		if page.Number() != simplePage.Number {
			t.Errorf("page %d: Number mismatch: default=%d, simple=%d", p, page.Number(), simplePage.Number)
		}

		if page.Rotation().Value != simplePage.Rotation {
			t.Errorf("page %d: Rotation mismatch: default=%d, simple=%d", p, page.Rotation().Value, simplePage.Rotation)
		}

		if page.MediaBox().Bounds != simplePage.MediaBox.Bounds {
			t.Errorf("page %d: MediaBox mismatch: default=%v, simple=%v", p, page.MediaBox().Bounds, simplePage.MediaBox.Bounds)
		}
	}
}

func TestSimpleFactory2(t *testing.T) {
	path := testutil.GetDocumentPath("cat-genetics", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	doc.AddPageFactory(testpages.NewSimplePageFactory(nil, nil, nil, nil, &content.ParsingOptions{}))

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	simplePage, err := content.GetTypedPage[testpages.SimplePage](doc, 1)
	if err != nil {
		t.Fatalf("GetPage[*testpages.SimplePage](1): %v", err)
	}

	if page.Number() != simplePage.Number {
		t.Errorf("Number mismatch: default=%d, simple=%d", page.Number(), simplePage.Number)
	}

	if page.Rotation().Value != simplePage.Rotation {
		t.Errorf("Rotation mismatch: default=%d, simple=%d", page.Rotation().Value, simplePage.Rotation)
	}

	if page.MediaBox().Bounds != simplePage.MediaBox.Bounds {
		t.Errorf("MediaBox mismatch: default=%v, simple=%v", page.MediaBox().Bounds, simplePage.MediaBox.Bounds)
	}

	simplePage2, err := content.GetTypedPage[testpages.SimplePage](doc, 1)
	if err != nil {
		t.Fatalf("GetPage[*testpages.SimplePage](1) second call: %v", err)
	}

	if page.Number() != simplePage2.Number {
		t.Errorf("Number mismatch (2nd call): default=%d, simple=%d", page.Number(), simplePage2.Number)
	}

	if page.Rotation().Value != simplePage2.Rotation {
		t.Errorf("Rotation mismatch (2nd call): default=%d, simple=%d", page.Rotation().Value, simplePage2.Rotation)
	}

	if page.MediaBox().Bounds != simplePage2.MediaBox.Bounds {
		t.Errorf("MediaBox mismatch (2nd call): default=%v, simple=%v", page.MediaBox().Bounds, simplePage2.MediaBox.Bounds)
	}
}

func TestInformationFactory(t *testing.T) {
	path := testutil.GetDocumentPath("Gamebook", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	doc.AddPageFactory(testpages.NewPageInformationFactory(nil, nil, nil, nil, &content.ParsingOptions{UseLenientParsing: true}))

	for p := 1; p < doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", p, pageAny)
		}

		pageInfo, err := content.GetTypedPage[testpages.PageInformation](doc, p)
		if err != nil {
			t.Fatalf("GetPage[testpages.PageInformation](%d): %v", p, err)
		}

		if page.Number() != pageInfo.Number {
			t.Errorf("page %d: Number mismatch: default=%d, info=%d", p, page.Number(), pageInfo.Number)
		}

		if page.Rotation() != pageInfo.Rotation {
			t.Errorf("page %d: Rotation mismatch: default=%v, info=%v", p, page.Rotation(), pageInfo.Rotation)
		}

		if page.Width() != pageInfo.Width {
			t.Errorf("page %d: Width mismatch: default=%.4f, info=%.4f", p, page.Width(), pageInfo.Width)
		}

		if page.Height() != pageInfo.Height {
			t.Errorf("page %d: Height mismatch: default=%.4f, info=%.4f", p, page.Height(), pageInfo.Height)
		}

		pageInfo2, err := content.GetTypedPage[testpages.PageInformation](doc, p)
		if err != nil {
			t.Fatalf("GetPage[testpages.PageInformation](%d) second call: %v", p, err)
		}

		if page.Number() != pageInfo2.Number {
			t.Errorf("page %d: Number mismatch (2nd): default=%d, info=%d", p, page.Number(), pageInfo2.Number)
		}

		if page.Rotation() != pageInfo2.Rotation {
			t.Errorf("page %d: Rotation mismatch (2nd): default=%v, info=%v", p, page.Rotation(), pageInfo2.Rotation)
		}

		if page.Width() != pageInfo2.Width {
			t.Errorf("page %d: Width mismatch (2nd): default=%.4f, info=%.4f", p, page.Width(), pageInfo2.Width)
		}

		if page.Height() != pageInfo2.Height {
			t.Errorf("page %d: Height mismatch (2nd): default=%.4f, info=%.4f", p, page.Height(), pageInfo2.Height)
		}
	}
}

func TestSimpleAndInformationFactory(t *testing.T) {
	path := testutil.GetDocumentPath("DeviceN_CS_test", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	doc.AddPageFactory(testpages.NewPageInformationFactory(nil, nil, nil, nil, &content.ParsingOptions{}))
	doc.AddPageFactory(testpages.NewSimplePageFactory(nil, nil, nil, nil, &content.ParsingOptions{}))

	for p := 1; p < doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", p, pageAny)
		}

		pageInfo, err := content.GetTypedPage[testpages.PageInformation](doc, p)
		if err != nil {
			t.Fatalf("GetPage[testpages.PageInformation](%d): %v", p, err)
		}

		if page.Number() != pageInfo.Number {
			t.Errorf("page %d: Number mismatch (info): default=%d, info=%d", p, page.Number(), pageInfo.Number)
		}

		if page.Rotation() != pageInfo.Rotation {
			t.Errorf("page %d: Rotation mismatch (info): default=%v, info=%v", p, page.Rotation(), pageInfo.Rotation)
		}

		if page.Width() != pageInfo.Width {
			t.Errorf("page %d: Width mismatch (info): default=%.4f, info=%.4f", p, page.Width(), pageInfo.Width)
		}

		if page.Height() != pageInfo.Height {
			t.Errorf("page %d: Height mismatch (info): default=%.4f, info=%.4f", p, page.Height(), pageInfo.Height)
		}

		simplePage, err := content.GetTypedPage[testpages.SimplePage](doc, p)
		if err != nil {
			t.Fatalf("GetPage[*testpages.SimplePage](%d): %v", p, err)
		}

		if page.Number() != simplePage.Number {
			t.Errorf("page %d: Number mismatch (simple): default=%d, simple=%d", p, page.Number(), simplePage.Number)
		}

		if page.Rotation().Value != simplePage.Rotation {
			t.Errorf("page %d: Rotation mismatch (simple): default=%d, simple=%d", p, page.Rotation().Value, simplePage.Rotation)
		}

		if page.MediaBox().Bounds != simplePage.MediaBox.Bounds {
			t.Errorf("page %d: MediaBox mismatch (simple): default=%v, simple=%v", p, page.MediaBox().Bounds, simplePage.MediaBox.Bounds)
		}
	}
}

func TestNoPageFactory(t *testing.T) {
	path := testutil.GetDocumentPath("cat-genetics", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, err = content.GetTypedPage[testpages.SimplePage](doc, 1)
	if err == nil {
		t.Fatal("expected error when requesting testpages.SimplePage without registered factory")
	}

	expectedPrefix := "could not find page factory of type"
	if len(err.Error()) < len(expectedPrefix) || err.Error()[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("error message should start with %q, got: %s", expectedPrefix, err.Error())
	}
}

// TestWrongSignatureFactory is skipped because Go's type system prevents registering
// a factory with an incompatible constructor at compile time. The C# original tests
// runtime validation of BasePageFactory constructor arguments, which does not have
// a direct equivalent in the Go port since NewBasePageFactory accepts typed parameters.
func TestWrongSignatureFactory(t *testing.T) {
	t.Skip("Go type system prevents wrong-signature factory registration at compile time; no direct C# equivalent")
}

/*
TestTableSimpleFactory verifies SimplePageFactory behavior with table-driven tests
across multiple documents and pages. This supplements the individual test functions
with a more compact verification pattern.
*/
func TestTableSimpleFactory(t *testing.T) {
	tests := []struct {
		name    string
		docName string
	}{
		{"ICML03-081", "ICML03-081"},
		{"cat-genetics", "cat-genetics"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := testutil.GetDocumentPath(tc.docName, true)

			doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
			}
			defer doc.Close()

			doc.AddPageFactory(testpages.NewSimplePageFactory(nil, nil, nil, nil, &content.ParsingOptions{}))

			for p := 1; p <= doc.NumberOfPages(); p++ {
				pageAny, err := doc.GetPage(p)
				if err != nil {
					t.Fatalf("GetPage(%d): %v", p, err)
				}

				page, ok := pageAny.(*content.Page)
				if !ok {
					t.Fatalf("GetPage(%d): expected *content.Page, got %T", p, pageAny)
				}

				simplePage, err := content.GetTypedPage[testpages.SimplePage](doc, p)
				if err != nil {
					t.Fatalf("GetPage[*testpages.SimplePage](%d): %v", p, err)
				}

				if simplePage.Number != page.Number() {
					t.Errorf("page %d: Number = %d, want %d", p, simplePage.Number, page.Number())
				}
			}
		})
	}
}

/*
TestTableInformationFactory verifies PageInformationFactory behavior with table-driven tests.
*/
func TestTableInformationFactory(t *testing.T) {
	tests := []struct {
		name    string
		docName string
	}{
		{"Gamebook", "Gamebook"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := testutil.GetDocumentPath(tc.docName, true)

			doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
			}
			defer doc.Close()

			doc.AddPageFactory(testpages.NewPageInformationFactory(nil, nil, nil, nil, &content.ParsingOptions{UseLenientParsing: true}))

			for p := 1; p <= doc.NumberOfPages(); p++ {
				pageAny, err := doc.GetPage(p)
				if err != nil {
					t.Fatalf("GetPage(%d): %v", p, err)
				}

				page, ok := pageAny.(*content.Page)
				if !ok {
					t.Fatalf("GetPage(%d): expected *content.Page, got %T", p, pageAny)
				}

				pageInfo, err := content.GetTypedPage[testpages.PageInformation](doc, p)
				if err != nil {
					t.Fatalf("GetPage[testpages.PageInformation](%d): %v", p, err)
				}

				if pageInfo.Number != page.Number() {
					t.Errorf("page %d: Number = %d, want %d", p, pageInfo.Number, page.Number())
				}

				if !pageInfo.Rotation.Equals(page.Rotation()) {
					t.Errorf("page %d: Rotation = %v, want %v", p, pageInfo.Rotation, page.Rotation())
				}

				if page.Width() != pageInfo.Width {
					t.Errorf("page %d: Width = %.4f, want %.4f", p, pageInfo.Width, page.Width())
				}

				if page.Height() != pageInfo.Height {
					t.Errorf("page %d: Height = %.4f, want %.4f", p, pageInfo.Height, page.Height())
				}
			}
		})
	}
}
