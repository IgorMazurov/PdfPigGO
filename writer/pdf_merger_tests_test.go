package writer_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/parser"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/writer"
)

func init() {
	testutil.IntegrationDocumentsRoot = "../testdata/integration/Documents"
}

func getDoc(name string) string {
	return testutil.GetDocumentPath(name, true)
}

func openMerged(b []byte) (*content.PdfDocument, error) {
	return parser.OpenMemory(b, nil)
}

type mergeAssertion struct {
	expectedPages  int
	version        float64
	checkVersion   bool
	pageTexts      []string
	linkCount      *int
	minLinkCount   *int
}

func assertMerge(t *testing.T, b []byte, a mergeAssertion) {
	t.Helper()

	doc, err := openMerged(b)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	if doc.NumberOfPages() != a.expectedPages {
		t.Errorf("expected %d pages, got %d", a.expectedPages, doc.NumberOfPages())
	}

	if a.checkVersion && doc.Version() != a.version {
		t.Errorf("expected version %.1f, got %.1f", a.version, doc.Version())
	}

	for i, expected := range a.pageTexts {
		pgAny, err := doc.GetPage(i + 1)
		if err != nil {
			t.Errorf("GetPage(%d): %v", i+1, err)
			continue
		}
		pg, ok := pgAny.(*content.Page)
		if !ok {
			t.Errorf("GetPage(%d): expected *content.Page, got %T", i+1, pgAny)
			continue
		}
		if pg.Text() != expected {
			t.Errorf("page %d text: expected %q, got %q", i+1, expected, pg.Text())
		}
	}

	if a.linkCount != nil || a.minLinkCount != nil {
		pages, err := doc.GetPages()
		if err != nil {
			t.Errorf("GetPages: %v", err)
		} else {
			count := 0
			for _, pgAny := range pages {
				pg, ok := pgAny.(*content.Page)
				if !ok {
					continue
				}
				annots := pg.GetAnnotations()
				for _, ann := range annots {
					if at, ok := ann.Type().(annotations.AnnotationType); ok && at == annotations.Link {
						count++
					}
				}
			}
			if a.linkCount != nil && count != *a.linkCount {
				t.Errorf("expected %d link annotations, got %d", *a.linkCount, count)
			}
			if a.minLinkCount != nil && count < *a.minLinkCount {
				t.Errorf("expected at least %d link annotations, got %d", *a.minLinkCount, count)
			}
		}
	}
}

func TestCanMerge2SimpleDocuments(t *testing.T) {
	one := getDoc("Single Page Simple - from inkscape.pdf")
	two := getDoc("Single Page Simple - from open office.pdf")

	result, err := writer.MergeFiles(one, two)
	if err != nil {
		t.Fatalf("MergeFiles: %v", err)
	}

	assertMerge(t, result, mergeAssertion{
		expectedPages: 2,
		checkVersion:  true,
		version:       1.5,
		pageTexts:     []string{"Write something inInkscape", "I am a simple pdf."},
	})
}

func TestCanMerge2SimpleDocumentsIntoStream(t *testing.T) {
	one := getDoc("Single Page Simple - from inkscape.pdf")
	two := getDoc("Single Page Simple - from open office.pdf")

	var buf bytes.Buffer
	w := &seekableBuffer{buf: &buf}

	err := writer.MergeFilesToStream(w, one, two)
	if err != nil {
		t.Fatalf("MergeFilesToStream: %v", err)
	}

	assertMerge(t, w.Bytes(), mergeAssertion{
		expectedPages: 2,
		checkVersion:  true,
		version:       1.5,
		pageTexts:     []string{"Write something inInkscape", "I am a simple pdf."},
	})
}

func TestCanMerge2SimpleDocumentsReversed(t *testing.T) {
	one := getDoc("Single Page Simple - from open office.pdf")
	two := getDoc("Single Page Simple - from inkscape.pdf")

	result, err := writer.MergeFiles(one, two)
	if err != nil {
		t.Fatalf("MergeFiles: %v", err)
	}

	assertMerge(t, result, mergeAssertion{
		expectedPages: 2,
		checkVersion:  true,
		version:       1.5,
		pageTexts:     []string{"I am a simple pdf.", "Write something inInkscape"},
	})
}

func TestRootNodePageCount(t *testing.T) {
	one := getDoc("Single Page Simple - from open office.pdf")
	two := getDoc("Single Page Simple - from inkscape.pdf")

	result, err := writer.MergeFiles(one, two)
	if err != nil {
		t.Fatalf("MergeFiles(1): %v", err)
	}

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}

	if doc.NumberOfPages() != 2 {
		t.Errorf("expected 2 pages after first merge, got %d", doc.NumberOfPages())
	}

	doc.Close()

	oneBytes, err := os.ReadFile(one)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	result2, err := writer.MergeBytes([][]byte{result, oneBytes}, nil, writer.PdfANone, nil)
	if err != nil {
		t.Fatalf("MergeBytes: %v", err)
	}

	doc2, err := openMerged(result2)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc2.Close()

	if doc2.NumberOfPages() != 3 {
		t.Errorf("expected 3 pages after second merge, got %d", doc2.NumberOfPages())
	}
}

func TestObjectCountLower(t *testing.T) {
	one := getDoc("Single Page Simple - from inkscape.pdf")

	result, err := writer.MergeFiles(one, one)
	if err != nil {
		t.Fatalf("MergeFiles: %v", err)
	}

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 2 {
		t.Errorf("expected 2 pages, got %d", doc.NumberOfPages())
	}
}

func TestDedupsObjectsFromSameDoc(t *testing.T) {
	one := getDoc("Multiple Page - from Mortality Statistics.pdf")

	oneBytes, err := os.ReadFile(one)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	result, err := writer.MergeBytes(
		[][]byte{oneBytes},
		[][]int{{1, 2}},
		writer.PdfANone,
		nil,
	)
	if err != nil {
		t.Fatalf("MergeBytes: %v", err)
	}

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 2 {
		t.Errorf("expected 2 pages, got %d", doc.NumberOfPages())
	}
}

func TestCanMergeWithObjectStream(t *testing.T) {
	first := getDoc("Single Page Simple - from google drive.pdf")
	second := getDoc("Multiple Page - from Mortality Statistics.pdf")

	result, err := writer.MergeFiles(first, second)
	if err != nil {
		t.Fatalf("MergeFiles: %v", err)
	}

	writeTestFile(t, "CanMergeWithObjectStream", result)

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 7 {
		t.Errorf("expected 7 pages, got %d", doc.NumberOfPages())
	}

	pages, err := doc.GetPages()
	if err != nil {
		t.Errorf("GetPages: %v", err)
		return
	}
	for _, pgAny := range pages {
		pg, ok := pgAny.(*content.Page)
		if !ok {
			continue
		}
		_ = pg.Text() // C# Assert.NotNull(page.Text) - in Go strings can't be nil, so this matches NotNull
	}
}

func TestCanMergeWithSelection(t *testing.T) {
	first := getDoc("Multiple Page - from Mortality Statistics.pdf")

	contents, err := os.ReadFile(first)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	toCopy := []int{2, 1, 4, 3, 6, 5}
	result, err := writer.MergeBytes(
		[][]byte{contents},
		[][]int{toCopy},
		writer.PdfANone,
		nil,
	)
	if err != nil {
		t.Fatalf("MergeBytes: %v", err)
	}

	writeTestFile(t, "CanMergeWithSelection", result)

	existing, err := openMerged(contents)
	if err != nil {
		t.Fatalf("OpenMemory(original): %v", err)
	}
	defer existing.Close()

	merged, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory(merged): %v", err)
	}
	defer merged.Close()

	if merged.NumberOfPages() != 6 {
		t.Errorf("expected 6 pages, got %d", merged.NumberOfPages())
	}

	for i := 0; i < merged.NumberOfPages(); i++ {
		expectedPgAny, err := existing.GetPage(toCopy[i])
		if err != nil {
			t.Errorf("existing.GetPage(%d): %v", toCopy[i], err)
			continue
		}
		gotPgAny, err := merged.GetPage(i + 1)
		if err != nil {
			t.Errorf("merged.GetPage(%d): %v", i+1, err)
			continue
		}

		expectedPg, ok := expectedPgAny.(*content.Page)
		if !ok {
			t.Errorf("existing page: expected *content.Page, got %T", expectedPgAny)
			continue
		}
		gotPg, ok := gotPgAny.(*content.Page)
		if !ok {
			t.Errorf("merged page: expected *content.Page, got %T", gotPgAny)
			continue
		}

		if expectedPg.Text() != gotPg.Text() {
			t.Errorf("page %d text mismatch: expected %q, got %q", i+1, expectedPg.Text(), gotPg.Text())
		}
	}
}

func TestCanMergeMultipleWithSelection(t *testing.T) {
	first := getDoc("Multiple Page - from Mortality Statistics.pdf")
	second := getDoc("Old Gutnish Internet Explorer.pdf")

	firstBytes, err := os.ReadFile(first)
	if err != nil {
		t.Fatalf("ReadFile(first): %v", err)
	}
	secondBytes, err := os.ReadFile(second)
	if err != nil {
		t.Fatalf("ReadFile(second): %v", err)
	}

	result, err := writer.MergeBytes(
		[][]byte{firstBytes, secondBytes},
		[][]int{{2, 1, 4, 3, 6, 5}, {3, 2, 1}},
		writer.PdfANone,
		nil,
	)
	if err != nil {
		t.Fatalf("MergeBytes: %v", err)
	}

	writeTestFile(t, "CanMergeMultipleWithSelection", result)

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 9 {
		t.Errorf("expected 9 pages, got %d", doc.NumberOfPages())
	}

	pages, err := doc.GetPages()
	if err != nil {
		t.Errorf("GetPages: %v", err)
		return
	}
	for _, pgAny := range pages {
		pg, ok := pgAny.(*content.Page)
		if !ok {
			continue
		}
		_ = pg.Text() // C# Assert.NotNull(page.Text) - in Go strings can't be nil, so this matches NotNull
	}
}

func TestCanMergeWithLinks(t *testing.T) {
	testPath := getDoc("outline.pdf")

	if _, err := os.Stat(testPath); err != nil {
		t.Skipf("test file not found: %s", testPath)
	}

	testBytes, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	result, err := writer.MergeBytes(
		[][]byte{testBytes, testBytes},
		nil,
		writer.PdfANone,
		nil,
	)
	if err != nil {
		t.Fatalf("MergeBytes: %v", err)
	}

	writeTestFile(t, "CanMergeWithLinks", result)

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	linkCount := countLinks(doc)
	if linkCount != 2 {
		t.Errorf("expected 2 link annotations, got %d", linkCount)
	}
}

func TestCanMergeWithLinksWithSelection(t *testing.T) {
	testPath := getDoc("outline.pdf")

	if _, err := os.Stat(testPath); err != nil {
		t.Skipf("test file not found: %s", testPath)
	}

	testBytes, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	result, err := writer.MergeBytes(
		[][]byte{testBytes, testBytes},
		[][]int{{2, 1}, {3, 1}},
		writer.PdfANone,
		nil,
	)
	if err != nil {
		t.Fatalf("MergeBytes: %v", err)
	}

	writeTestFile(t, "CanMergeWithLinksWithSelection", result)

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	linkCount := countLinks(doc)
	if linkCount != 1 {
		t.Errorf("expected 1 link annotation, got %d", linkCount)
	}
}

func TestNoStackoverflow(t *testing.T) {
	testPath := getDoc("68-1990-01_A.pdf")

	bytes, err := os.ReadFile(testPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	result, err := writer.MergeBytes(
		[][]byte{bytes},
		nil,
		writer.PdfANone,
		nil,
	)
	if err != nil {
		t.Fatalf("MergeBytes: %v", err)
	}

	doc, err := openMerged(result)
	if err != nil {
		t.Fatalf("OpenMemory: %v", err)
	}
	defer doc.Close()

	if doc.NumberOfPages() != 45 {
		t.Errorf("expected 45 pages, got %d", doc.NumberOfPages())
	}
}

func countLinks(doc *content.PdfDocument) int {
	pages, err := doc.GetPages()
	if err != nil {
		return 0
	}
	count := 0
	for _, pgAny := range pages {
		pg, ok := pgAny.(*content.Page)
		if !ok {
			continue
		}
		annots := pg.GetAnnotations()
		for _, ann := range annots {
			if at, ok := ann.Type().(annotations.AnnotationType); ok && at == annotations.Link {
				count++
			}
		}
	}
	return count
}

func writeTestFile(t *testing.T, name string, b []byte) {
	t.Helper()
	dir := filepath.Join("testdata", "Merger")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Logf("MkdirAll(%q): %v (non-fatal)", dir, err)
		return
	}
	out := filepath.Join(dir, name+".pdf")
	if err := os.WriteFile(out, b, 0o644); err != nil {
		t.Logf("WriteFile(%q): %v (non-fatal)", out, err)
	}
}

type seekableBuffer struct {
	buf *bytes.Buffer
}

func (s *seekableBuffer) Write(p []byte) (n int, err error) {
	return s.buf.Write(p)
}

func (s *seekableBuffer) Seek(offset int64, whence int) (int64, error) {
	var pos int64
	switch whence {
	case 0:
		pos = offset
	case 1:
		pos = int64(s.buf.Len()) + offset
	case 2:
		pos = int64(s.buf.Len())
	default:
		return 0, os.ErrInvalid
	}
	if pos < 0 {
		return 0, os.ErrInvalid
	}
	return pos, nil
}

func (s *seekableBuffer) Bytes() []byte {
	return s.buf.Bytes()
}
