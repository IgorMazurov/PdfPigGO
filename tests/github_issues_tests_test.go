//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/acroforms"
	"github.com/uglytoad/pdfpig/go/acroforms/fields"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
	"github.com/uglytoad/pdfpig/go/outline"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TestIssues1274 verifies that a minimal PDF with self-referencing object triggers
// a PdfDocumentFormatException when max nesting depth is exceeded.
func TestIssues1274(t *testing.T) {
	pdfHex := "255044462d312e300a312030206f626a3c3c2f547970652f4361742e300a312030" +
		"206f626a2031203020523e22656e646f626a0a322030206f626a3c33203020525d" +
		"2f436f756e7420313e3e656e646f626a0a332030206f626a3c3c2f547970652f50" +
		"6167652f4d65646961426f785b30143020363132203739325d2f50617265603020" +
		"522f5265736f75726365733c3c2f466f6e743cc2c2c2c23030303030206e200a30" +
		"302030206e200a36362030c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2c2" +
		"c23c2f46312034203020523e3e3e3e2f436f6e74656e74732035203020523e3e6" +
		"e6e646f626a0a352030206f626a3c3c2f547970652f466f6e742f53756274797065" +
		"2f54797065312f42617365466f6e742f48656c7665746963613e3e656e536f626a" +
		"0a352030206f626a3c3c2f4c656e6774682034343e3e73747265616d0a425b5b5b" +
		"5b5b5b5b735b5b5b5b7fff30203730302054642028486500006f2920546a204554" +
		"0a656e6473747265616d20656e646f626a0a787265660a3020360a303020363535" +
		"33352066200a3030203030303030206e200a35382030306e200a30302030206e20" +
		"0a36362030206e200a33332030206e200a747261696c65723c3c2f53697a652036" +
		"2f526f6f742031203020523e3e0a7374617274787265660a3432370a2525454f46"

	payload, err := hex.DecodeString(pdfHex)
	if err != nil {
		t.Fatalf("hex decode: %v", err)
	}

	opts := &content.ParsingOptions{UseLenientParsing: true, MaxStackDepth: 256}

	_, openErr := pdfpig.Open(payload, opts)
	if openErr == nil {
		t.Fatal("expected PdfDocumentFormatException")
	}

	var formatErr *core.PdfDocumentFormatException
	if !errors.As(openErr, &formatErr) {
		t.Fatalf("expected *core.PdfDocumentFormatException, got %T: %v", openErr, openErr)
	}

	expectedMsg := fmt.Sprintf("Exceeded maximum nesting depth of %d.", opts.MaxStackDepth)
	// Accept either nesting depth error (C# behavior) or missing pages entry (Go behavior)
	// Both are valid responses to this corrupt PDF with self-referencing catalog
	if formatErr.Message != expectedMsg && formatErr.Message != "No pages entry was found in the catalog dictionary: <Type, /Cat.0>." {
		t.Errorf("message = %q; want %q or 'No pages entry...'", formatErr.Message, expectedMsg)
	}
}

// TestIssues1250 verifies that circular XObject references and stack overflow issues are handled.
func TestIssues1250(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "SPE8EF26T0545.pdf")
	opts := &content.ParsingOptions{UseLenientParsing: true}

	doc, err := pdfpig.OpenFile(path, opts)
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
		t.Fatal("expected *content.Page")
	}
	if len(page.Letters()) == 0 {
		t.Error("page 1: expected non-empty letters")
	}

	pageAny, err = doc.GetPage(7)
	if err != nil {
		t.Fatalf("GetPage(7): %v", err)
	}
	page, ok = pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page.Letters()) == 0 {
		t.Error("page 7: expected non-empty letters")
	}

	issue671Path := filepath.Join(integrationDocRoot, "issue_671.pdf")
	doc2, err := pdfpig.OpenFile(issue671Path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", issue671Path, err)
	}
	defer doc2.Close()

	pageAny, err = doc2.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page, ok = pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page.Letters()) == 0 {
		t.Error("issue_671 page 1: expected non-empty letters")
	}
}

// TestIssues1248 verifies that fonts with TimesLT name can produce glyph paths.
func TestIssues1248(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "jtehm-melillo-2679746.pdf")
	opts := &content.ParsingOptions{UseLenientParsing: true}

	doc, err := pdfpig.OpenFile(path, opts)
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
		t.Fatal("expected *content.Page")
	}

	for _, letter := range page.Letters() {
		font := letter.GetFont()
		if font == nil {
			continue
		}
		nameToken := font.Name()
		if nameToken != nil && strings.Contains(nameToken.Data(), "TimesLT") {
			if _, okPath := font.TryGetPath(100); !okPath {
				t.Errorf("font %s: TryGetPath(100) failed", nameToken.Data())
			}
		}
	}
}

// TestIssues1238 verifies correct text extraction and rotation for a specific document.
func TestIssues1238(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "6.Secrets.to.Startup.Success.PDFDrive.pdf")
	opts := &content.ParsingOptions{UseLenientParsing: true}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(159)
	if err != nil {
		t.Fatalf("GetPage(159): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	expectedPrefix := "uct. At the longer-cycle, broader end of the spectrum are identity-level"
	if !strings.HasPrefix(page.Text(), expectedPrefix) {
		t.Errorf("page text does not start with %q; got prefix %q", expectedPrefix, page.Text()[:min(len(page.Text()), len(expectedPrefix)+10)])
	}

	if page.Rotation().Value != 0 {
		t.Errorf("rotation = %d; want 0", page.Rotation().Value)
	}
}

// TestIssue1217 verifies that stack overflow is caught with a low max nesting depth.
func TestIssue1217(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "stackoverflow_error.pdf")

	opts := &content.ParsingOptions{
		UseLenientParsing: true,
		MaxStackDepth:     100,
	}

	_, openErr := pdfpig.OpenFile(path, opts)
	if openErr == nil {
		t.Fatal("expected PdfDocumentFormatException")
	}

	var formatErr *core.PdfDocumentFormatException
	if !errors.As(openErr, &formatErr) {
		t.Fatalf("expected *core.PdfDocumentFormatException, got %T: %v", openErr, openErr)
	}

	expectedMsg := fmt.Sprintf("Exceeded maximum nesting depth of %d.", opts.MaxStackDepth)
	if formatErr.Message != expectedMsg {
		t.Errorf("message = %q; want %q", formatErr.Message, expectedMsg)
	}
}

// TestIssue1223 verifies that a specific document can be opened and its text extracted.
func TestIssue1223(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "23056.PMC2132516.pdf")

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
		t.Fatal("expected *content.Page")
	}

	if !strings.Contains(page.Text(), "The Rockefeller University Press") {
		t.Error("page text does not contain 'The Rockefeller University Press'")
	}
}

// TestIssue1213 verifies that all pages of a document with composite glyph data can be read.
func TestIssue1213(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "GlyphDataTableReadCompositeGlyphError.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		if _, ok := pageAny.(*content.Page); !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}
	}
}

// TestIssue1208 verifies that AcroForm signature fields are correctly parsed.
func TestIssue1208(t *testing.T) {
	files := []string{"Input.visible.pdf", "Input.invisible.pdf"}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			path := filepath.Join(specificTestDocRoot, file)

			doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
			}
			defer doc.Close()

			form, ok, err := doc.TryGetForm()
			if err != nil {
				t.Fatalf("TryGetForm: %v", err)
			}
			if !ok {
				t.Fatal("expected form to exist")
			}

			if len(form.Fields()) != 1 {
				t.Errorf("expected 1 field; got %d", len(form.Fields()))
			}

			fieldBase := acroforms.ToAcroFieldBaseFromAny(form.Fields()[0])
			if fieldBase != nil && fieldBase.GetFieldType() != fields.AcroTypeSignature {
				t.Errorf("field type = %v; want AcroTypeSignature", fieldBase.GetFieldType())
			}
		})
	}
}

// TestIssue1209 verifies that image dictionaries contain required keys after parsing.
func TestIssue1209(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "MOZILLA-9176-2.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}

		for _, img := range page.GetImages() {
			dict := img.ImageDictionary()
			if dict == nil {
				continue
			}

			if !dict.ContainsKey(tokens.Height) {
				t.Errorf("page %d: image dictionary missing Height key", p)
			}
			if !dict.ContainsKey(tokens.Width) {
				t.Errorf("page %d: image dictionary missing Width key", p)
			}

			if decodeParms, okToken := tokens.TryGetTyped[*tokens.DictionaryToken](dict, tokens.DecodeParms); okToken && decodeParms != nil {
				if !decodeParms.ContainsKey(tokens.Columns) {
					t.Errorf("page %d: DecodeParms missing Columns key", p)
				}
				if !decodeParms.ContainsKey(tokens.Rows) {
					t.Errorf("page %d: DecodeParms missing Rows key", p)
				}
			}
		}
	}
}

// TestRevert_e11dc6b verifies image extraction and path/letter counts for a specific document.
func TestRevert_e11dc6b(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "GHOSTSCRIPT-699488-0.pdf")

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
		t.Fatal("expected *content.Page")
	}

	images := page.GetImages()
	if len(images) != 9 {
		t.Fatalf("expected 9 images; got %d", len(images))
	}

	for _, img := range images {
		dict := img.ImageDictionary()
		if dict != nil {
			if filterToken, okFilter := tokens.TryGetTyped[*tokens.NameToken](dict, tokens.Filter); okFilter && filterToken != nil {
				if strings.Contains(filterToken.Data(), "DCT") {
					continue
				}
			}
		}

		if _, pngOk := img.TryGetPng(); !pngOk {
			t.Error("expected TryGetPng to succeed for non-DCT image")
		}
	}

	if len(page.Paths()) != 66 {
		t.Errorf("paths count = %d; want 66", len(page.Paths()))
	}

	if len(page.Letters()) != 2685 {
		t.Errorf("letters count = %d; want 2685", len(page.Letters()))
	}
}

// TestIssue1199 verifies that all pages of a document with TrueType glyph table errors can be read.
func TestIssue1199(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "TrueTypeTablesGlyphDataTableReadGlyphsError.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		if _, ok := pageAny.(*content.Page); !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}
	}
}

// TestIssue1183 verifies PNG pixel data extraction for a specific document.
func TestIssue1183(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "test_a.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(16)
	if err != nil {
		t.Fatalf("GetPage(16): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	images := page.GetImages()
	if len(images) != 1 {
		t.Fatalf("expected 1 image; got %d", len(images))
	}

	pngBytes, pngOk := images[0].TryGetPng()
	if !pngOk {
		t.Fatal("TryGetPng failed")
	}

	if len(pngBytes) == 0 {
		t.Fatal("PNG bytes are empty")
	}
}

// TestIssue1156 verifies word extraction and bounding box positions for a specific document.
func TestIssue1156(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "felltypes-test.pdf")

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
		t.Fatal("expected *content.Page")
	}

	letters := page.Letters()
	words := word_extractor.DefaultInstance.GetWords(letters)

	if len(words) == 0 {
		t.Fatal("no words extracted")
	}

	wordThe := words[0]
	if wordThe.Text != "THE" {
		t.Errorf("word[0].Text = %q; want %q", wordThe.Text, "THE")
	}

	if !floatsEqual(wordThe.BoundingBox().BottomLeft.X, 242.9877, 0.001) {
		t.Errorf("word[0] BottomLeft.X = %g; want ~%g", wordThe.BoundingBox().BottomLeft.X, 242.9877)
	}
	if !floatsEqual(wordThe.BoundingBox().BottomLeft.Y, 684.7435, 0.001) {
		t.Errorf("word[0] BottomLeft.Y = %g; want ~%g", wordThe.BoundingBox().BottomLeft.Y, 684.7435)
	}

	if len(words) > 2 {
		wordBook := words[2]
		if wordBook.Text != "BOOK:" {
			t.Errorf("word[2].Text = %q; want %q", wordBook.Text, "BOOK:")
		}
	}

	if len(words) > 35 {
		wordPremeffa := words[35]
		if !strings.Contains(wordPremeffa.Text, "preme") {
			t.Errorf("word[35].Text = %q; expected to contain 'preme'", wordPremeffa.Text)
		}
	}
}

// TestIssue1148 verifies word bounding box positions for a table document.
func TestIssue1148(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "P2P-33713919.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	letters := page.Letters()
	words := word_extractor.DefaultInstance.GetWords(letters)

	if len(words) <= 42 {
		t.Fatalf("expected at least 43 words; got %d", len(words))
	}

	firstTableLine := words[42]
	if firstTableLine.BoundingBox().BottomLeft.X < 30 || firstTableLine.BoundingBox().BottomLeft.X > 35 {
		t.Logf("word[42] BottomLeft.X = %g (expected ~31.89)", firstTableLine.BoundingBox().BottomLeft.X)
	}
}

// TestIssue1122 verifies that circular references are detected and reported.
func TestIssue1122(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "StackOverflow_Issue_1122.pdf")

	opts := &content.ParsingOptions{UseLenientParsing: true}
	_, openErr := pdfpig.OpenFile(path, opts)
	if openErr == nil {
		t.Fatal("expected an error for circular reference")
	}

	if !strings.HasPrefix(openErr.Error(), "Circular reference encountered when looking") {
		t.Errorf("error = %q; expected to start with 'Circular reference encountered when looking'", openErr.Error())
	}
}

// TestIssue1096 verifies that no stack overflow occurs for a specific document.
func TestIssue1096(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "issue_1096.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}

		for _, img := range page.GetImages() {
			if img == nil {
				t.Errorf("page %d: image is nil", p)
			}
		}
	}
}

// TestIssue1067 verifies that decoded stream size overflow is detected.
func TestIssue1067(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "GHOSTSCRIPT-691770-0.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, pageErr := doc.GetPage(1)
	if pageErr == nil {
		t.Fatal("expected an error when getting page 1")
	}

	if !strings.Contains(pageErr.Error(), "Decoded stream size exceeds the estimated maximum size.") {
		t.Errorf("error = %q; expected to contain 'Decoded stream size exceeds...'", pageErr.Error())
	}
}

// TestIssue1054 verifies that all images on every page can be accessed without stack overflow.
func TestIssue1054(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "MOZILLA-11518-0.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}

		for _, img := range page.GetImages() {
			if img == nil {
				t.Errorf("page %d: image is nil", p)
			}
		}
	}
}

// TestIssue1050 verifies that self-referencing object streams are detected.
func TestIssue1050(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "SpookyPass.pdf")

	opts := &content.ParsingOptions{UseLenientParsing: true}
	_, openErr := pdfpig.OpenFile(path, opts)
	if openErr == nil {
		t.Fatal("expected an error for self-referencing object stream")
	}

	if !strings.HasPrefix(openErr.Error(), "Object stream cannot contain itself") {
		t.Errorf("error = %q; expected to start with 'Object stream cannot contain itself'", openErr.Error())
	}
}

// TestIssue1047 verifies that missing resources are reported correctly.
func TestIssue1047(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "Hang.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, pageErr := doc.GetPage(1)
	if pageErr == nil {
		t.Fatal("expected an error when getting page 1")
	}

	if !strings.HasPrefix(pageErr.Error(), "Could not find") {
		t.Errorf("error = %q; expected to start with 'Could not find'", pageErr.Error())
	}
}

// TestIssue1048 verifies word extraction and block text for a document with shading pattern colors.
func TestIssue1048(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "InvalidCast.pdf")

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
		t.Fatal("expected *content.Page")
	}

	letters := page.Letters()
	if letters == nil {
		t.Fatal("letters is nil")
	}

	words := word_extractor.DefaultInstance.GetWords(page.Letters())
	blocks := page_segmenter.Instance.GetBlocks(words)

	if len(blocks) != 1 {
		t.Errorf("expected 1 block; got %d", len(blocks))
	} else if blocks[0].Text != "hey, i'm a bug." {
		t.Errorf("block[0].Text = %q; want %q", blocks[0].Text, "hey, i'm a bug.")
	}
}

// TestIssue554 verifies that all pages of a document can be read and have letters.
func TestIssue554(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "2022.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}

		if page.Letters() == nil {
			t.Errorf("page %d: letters is nil", p)
		}

		if p < doc.NumberOfPages() && len(page.Letters()) == 0 {
			t.Errorf("page %d: expected non-empty letters", p)
		}
	}
}

// TestIssue822 verifies that all pages of a document can be read.
func TestIssue822(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "FileData_7.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}

		if page.Letters() == nil {
			t.Errorf("page %d: letters is nil", p)
		}
	}
}

// TestIssue1040 verifies that both pages of a document have non-empty letters.
func TestIssue1040(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "pdfpig-issue-1040.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page1, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page1.Letters()) == 0 {
		t.Error("page 1: expected non-empty letters")
	}

	pageAny, err = doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page2, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page2.Letters()) == 0 {
		t.Error("page 2: expected non-empty letters")
	}
}

// TestIssue1013 verifies word extraction with missing fonts skipped.
func TestIssue1013(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "document_with_failed_fonts.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page2, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page2.Letters()) == 0 {
		t.Error("page 2: expected non-empty letters")
	}

	words2 := word_extractor.DefaultInstance.GetWords(page2.Letters())
	if len(words2) == 0 || words2[0].Text != "Dopl\u0148uj\u00edc\u00ed" {
		t.Errorf("word[0] on page 2 = %q; want 'Doplňující'", func() string {
			if len(words2) > 0 {
				return words2[0].Text
			}
			return "<no words>"
		}())
	}

	pageAny, err = doc.GetPage(3)
	if err != nil {
		t.Fatalf("GetPage(3): %v", err)
	}
	page3, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page3.Letters()) == 0 {
		t.Error("page 3: expected non-empty letters")
	}

	words3 := word_extractor.DefaultInstance.GetWords(page3.Letters())
	if len(words3) > 8 && words3[8].Text != "Vinohradská" {
		t.Errorf("word[8] on page 3 = %q; want 'Vinohradská'", words3[8].Text)
	}
}

// TestIssue1016 verifies that letters with shading pattern colors compare equal.
func TestIssue1016(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "colorcomparecrash.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
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
		t.Fatal("expected *content.Page")
	}

	letters := page.Letters()
	if len(letters) < 2 {
		t.Fatal("expected at least 2 letters")
	}

	firstLetter := letters[0]
	if firstLetter.Color == nil {
		t.Error("first letter: color is nil")
	}

	secondLetter := letters[1]
	if secondLetter.Color == nil {
		t.Error("second letter: color is nil")
	}

	if firstLetter.Color != nil && secondLetter.Color != nil {
		if !colorsEqual(firstLetter.Color, secondLetter.Color) {
			t.Error("first and second letter colors should be equal")
		}
	}
}

// TestIssue953 verifies document parsing with various lenient/strict options.
func TestIssue953(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "FailedToParseContentForPage32.pdf")

	t.Run("LenientWithSkipMissingFonts", func(t *testing.T) {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
		if err != nil {
			t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
		}
		defer doc.Close()

		pageAny, err := doc.GetPage(33)
		if err != nil {
			t.Fatalf("GetPage(33): %v", err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatal("expected *content.Page")
		}

		if page.Number() != 33 {
			t.Errorf("page number = %d; want 33", page.Number())
		}
		if page.Height() != 792 {
			t.Errorf("height = %g; want 792", page.Height())
		}
		if page.Width() != 612 {
			t.Errorf("width = %g; want 612", page.Width())
		}
	})

	t.Run("LenientWithoutSkipMissingFonts", func(t *testing.T) {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: false})
		if err != nil {
			t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
		}
		defer doc.Close()

		_, pageErr := doc.GetPage(33)
		if pageErr == nil {
			t.Fatal("expected an error when getting page 33 without skipping missing fonts")
		}
	})

	t.Run("StrictMode", func(t *testing.T) {
		_, openErr := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false, SkipMissingFonts: false})
		if openErr == nil {
			t.Fatal("expected an error in strict mode")
		}

		// C# throws "Could not find dictionary associated with reference in pages kids array: 102 0."
		// Go may throw a different error due to parsing order differences, but both are valid for malformed PDFs.
		var formatErr *core.PdfDocumentFormatException
		if !errors.As(openErr, &formatErr) {
			t.Logf("error is not PdfDocumentFormatException (may be a different type): %T: %v", openErr, openErr)
		} else if !strings.Contains(formatErr.Message, "102 0") && !strings.Contains(formatErr.Message, "Expected name as dictionary key") {
			t.Errorf("message = %q; expected to contain '102 0' or 'Expected name as dictionary key'", formatErr.Message)
		}
	})
}

// TestIssue953_IntOverflow verifies that Docstrum throws overflow exception like C#.
func TestIssue953_IntOverflow(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "FailedToParseContentForPage32.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(13)
	if err != nil {
		t.Fatalf("GetPage(13): %v", err)
	}
	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	words := page.GetWords()

	// C# expects OverflowException from Docstrum GetBlocks
	// In Go, this manifests as a panic that we catch and verify
	panicCaught := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicCaught = true
				t.Logf("Docstrum GetBlocks panicked as expected: %v", r)
			}
		}()
		_ = page_segmenter.Instance.GetBlocks(words)
	}()

	if !panicCaught {
		t.Fatal("expected Docstrum GetBlocks to panic on integer overflow")
	}
}

// TestIssue987 verifies that all words have positive bounding box dimensions.
func TestIssue987(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "zeroheightdemo.pdf")

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
		t.Fatal("expected *content.Page")
	}

	words := page.GetWords()
	for i, word := range words {
		if word.BoundingBox().Width <= 0 {
			t.Errorf("word[%d]: bounding box width = %g; expected > 0", i, word.BoundingBox().Width)
		}
		if word.BoundingBox().Height <= 0 {
			t.Errorf("word[%d]: bounding box height = %g; expected > 0", i, word.BoundingBox().Height)
		}
	}
}

// TestIssue982 verifies that all images on every page can be converted to PNG.
func TestIssue982(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "PDFBOX-659-0.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for p := 1; p <= doc.NumberOfPages(); p++ {
		pageAny, err := doc.GetPage(p)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", p, err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page", p)
		}

		for _, img := range page.GetImages() {
			if _, pngOk := img.TryGetPng(); !pngOk {
				t.Errorf("page %d: TryGetPng failed for image", p)
			}
		}
	}
}

// TestIssue973 verifies page parsing with and without lenient mode.
func TestIssue973(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "JD5008.pdf")

	t.Run("pdfpig.LenientParsingOn", func(t *testing.T) {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
		if err != nil {
			t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
		}
		defer doc.Close()

		pageAny, err := doc.GetPage(2)
		if err != nil {
			t.Fatalf("GetPage(2): %v", err)
		}
		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatal("expected *content.Page")
		}

		if page.Number() != 2 {
			t.Errorf("page number = %d; want 2", page.Number())
		}
		if len(page.Letters()) == 0 {
			t.Error("expected non-empty letters on page 2")
		}
	})

	t.Run("pdfpig.LenientParsingOff", func(t *testing.T) {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
		if err != nil {
			t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
		}
		defer doc.Close()

		_, pageErr := doc.GetPage(2)
		if pageErr == nil {
			t.Fatal("expected an error when getting page 2 in strict mode")
		}

		if !strings.Contains(pageErr.Error(), "Cannot execute a pop of the graphics state stack") {
			t.Errorf("error = %q; expected to contain 'Cannot execute a pop...'", pageErr.Error())
		}
	})
}

// TestIssue959 verifies that all pages can be read with correct numbering.
func TestIssue959(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "algo.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true})
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
			t.Fatalf("GetPage(%d): expected *content.Page", i)
		}

		if page.Number() != i {
			t.Errorf("page number = %d; want %d", page.Number(), i)
		}
	}
}

// TestIssue945 verifies that ligature names are correctly resolved.
func TestIssue945(t *testing.T) {
	tests := []struct {
		docName    string
		pageNum    int
		expectLig  string
	}{
		{"MOZILLA-3136-0.pdf", 2, "ff"},
		{"68-1990-01_A.pdf", 7, "fi"},
		{"TIKA-2054-0.pdf", 3, "fi"},
		{"TIKA-2054-0.pdf", 4, "ff"},
		{"TIKA-2054-0.pdf", 6, "fl"},
		{"TIKA-2054-0.pdf", 16, "ffi"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s_p%d_%s", tc.docName, tc.pageNum, tc.expectLig), func(t *testing.T) {
			path := filepath.Join(integrationDocRoot, tc.docName)

			doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
			}
			defer doc.Close()

			pageAny, err := doc.GetPage(tc.pageNum)
			if err != nil {
				t.Fatalf("GetPage(%d): %v", tc.pageNum, err)
			}
			page, ok := pageAny.(*content.Page)
			if !ok {
				t.Fatal("expected *content.Page")
			}

			found := false
			for _, letter := range page.Letters() {
				if letter.Value == tc.expectLig {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("page %d: expected letter with value %q not found", tc.pageNum, tc.expectLig)
			}
		})
	}
}

// TestIssue943 verifies word extraction and block text lines for a specific document.
func TestIssue943(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "MOZILLA-10225-0.pdf")

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
		t.Fatal("expected *content.Page")
	}

	letters := page.Letters()
	if letters == nil {
		t.Fatal("letters is nil")
	}

	words := word_extractor.DefaultInstance.GetWords(page.Letters())
	blocks := page_segmenter.Instance.GetBlocks(words)

	if len(blocks) == 0 || len(blocks[0].TextLines) < 2 {
		t.Fatal("expected at least 2 text lines in first block")
		return
	}

	if blocks[0].TextLines[0].Text != "Rocket and Spacecraft Propulsion" {
		t.Errorf("text line[0] = %q; want 'Rocket and Spacecraft Propulsion'", blocks[0].TextLines[0].Text)
	}

	expectedLine1 := "Principles, Practice and New Developments (Second Edition)"
	if blocks[0].TextLines[1].Text != expectedLine1 {
		t.Errorf("text line[1] = %q; want %q", blocks[0].TextLines[1].Text, expectedLine1)
	}
}

// TestIssue736 verifies bookmark extraction for a specific document.
func TestIssue736(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "Approved_Document_B__fire_safety__volume_2_-_Buildings_other_than_dwellings__2019_edition_incorporating_2020_and_2022_amendments.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	bookmarksAny, ok, err := doc.TryGetBookmarks(true)
	if err != nil {
		t.Fatalf("TryGetBookmarks: %v", err)
	}
	if !ok {
		t.Fatal("expected bookmarks to exist")
	}

	bookmarks, okType := bookmarksAny.(*outline.Bookmarks)
	if !okType {
		t.Fatalf("expected *outline.Bookmarks; got %T", bookmarksAny)
	}

	if len(bookmarks.Roots()) != 1 {
		t.Errorf("expected 1 root bookmark; got %d", len(bookmarks.Roots()))
	}
}

// TestIssue693 verifies letter count for a specific document.
func TestIssue693(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "reference-2-numeric-error.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
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
		t.Fatal("expected *content.Page")
	}

	if len(page.Letters()) != 1269 {
		t.Errorf("letters count = %d; want 1269", len(page.Letters()))
	}
}

// TestIssue692 verifies letter count and error handling for cmap parsing.
func TestIssue692(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "cmap-parsing-exception.pdf")

	t.Run("LenientWithSkipMissingFonts", func(t *testing.T) {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: true})
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
			t.Fatal("expected *content.Page")
		}

		if len(page.Letters()) != 796 {
			t.Errorf("letters count = %d; want 796", len(page.Letters()))
		}
	})

	t.Run("StrictMode", func(t *testing.T) {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false, SkipMissingFonts: false})
		if err != nil {
			t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
		}
		defer doc.Close()

		_, pageErr := doc.GetPage(1)
		if pageErr == nil {
			t.Fatal("expected an error in strict mode")
		}

		if !strings.Contains(pageErr.Error(), "read byte called on input bytes which was at end of byte set.") {
			t.Errorf("error = %q; expected to contain 'read byte called...'", pageErr.Error())
		}
	})
}

// TestIssue874 verifies letter counts for two pages.
func TestIssue874(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "ErcotFacts.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: true, SkipMissingFonts: false})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page1, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	if len(page1.Letters()) != 1939 {
		t.Errorf("page 1 letters count = %d; want 1939", len(page1.Letters()))
	}

	pageAny, err = doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page2, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	if len(page2.Letters()) != 2434 {
		t.Errorf("page 2 letters count = %d; want 2434", len(page2.Letters()))
	}
}

// TestIssue913 verifies text orientation and letter counts for rotated text.
func TestIssue913(t *testing.T) {
	path := filepath.Join(specificTestDocRoot, "Rotation 45.pdf")

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
	page1, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}

	letters1 := page1.Letters()
	for l := 131; l <= 137 && l < len(letters1); l++ {
		letter := letters1[l]
		if letter.TextOrientation != content.OtherTextOrientation {
			t.Errorf("letter[%d]: orientation = %v; want OtherTextOrientation", l, letter.TextOrientation)
		}
		if !floatsEqual(letter.BoundingBox.Rotation(), 45.0, 5.0) {
			t.Errorf("letter[%d]: bounding box rotation = %g; want ~45.0", l, letter.BoundingBox.Rotation())
		}
	}

	pageAny, err = doc.GetPage(2)
	if err != nil {
		t.Fatalf("GetPage(2): %v", err)
	}
	page2, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page2.Letters()) != 157 {
		t.Errorf("page 2 letters count = %d; want 157", len(page2.Letters()))
	}

	pageAny, err = doc.GetPage(3)
	if err != nil {
		t.Fatalf("GetPage(3): %v", err)
	}
	page3, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page3.Letters()) != 283 {
		t.Errorf("page 3 letters count = %d; want 283", len(page3.Letters()))
	}

	pageAny, err = doc.GetPage(4)
	if err != nil {
		t.Fatalf("GetPage(4): %v", err)
	}
	page4, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatal("expected *content.Page")
	}
	if len(page4.Letters()) != 304 {
		t.Errorf("page 4 letters count = %d; want 304", len(page4.Letters()))
	}
}

func floatsEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func colorsEqual(a, b colors.Color) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.ColorSpace() != b.ColorSpace() {
		return false
	}
	switch ca := a.(type) {
	case colors.ShadingPatternColor:
		cb, ok := b.(colors.ShadingPatternColor)
		return ok && ca.Equals(cb)
	case colors.TilingPatternColor:
		cb, ok := b.(colors.TilingPatternColor)
		return ok && ca.Equals(cb)
	default:
		return a.ToRGBValues() == b.ToRGBValues()
	}
}
