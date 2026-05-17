package writer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/parser"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/writer"
)

func init() {
	testutil.IntegrationDocumentsRoot = "../testdata/integration/Documents"
}

func getTextRemoverDoc(name string) string {
	return testutil.GetDocumentPath(name, true)
}

func TestTextRemoverRemovesText(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"Two Page Text Only", "Two Page Text Only - from libre office.pdf"},
		{"Cat Genetics", "cat-genetics.pdf"},
		{"Motor Insurance Claim Form", "Motor Insurance claim form.pdf"},
		{"Single Page Images", "Single Page Images - from libre office.pdf"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filePath := getTextRemoverDoc(tc.file)

			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("ReadFile(%q): %v", filePath, err)
			}

			doc, err := parser.OpenMemory(data, nil)
			if err != nil {
				t.Fatalf("OpenMemory(original): %v", err)
			}
			defer doc.Close()

			withoutText, err := writer.RemoveTextBytes(filePath, nil)
			if err != nil {
				t.Fatalf("RemoveTextBytes: %v", err)
			}

			outputDir := "Writer"
			err = os.MkdirAll(outputDir, 0755)
			if err != nil {
				t.Fatalf("MkdirAll(%q): %v", outputDir, err)
			}
			outPath := filepath.Join(outputDir, "TextRemoverRemovesText_"+tc.file)
			err = os.WriteFile(outPath, withoutText, 0644)
			if err != nil {
				t.Logf("WriteFile(%q): %v (non-fatal)", outPath, err)
			}

			docWithoutText, err := parser.OpenMemory(withoutText, nil)
			if err != nil {
				t.Fatalf("OpenMemory(stripped): %v", err)
			}
			defer docWithoutText.Close()

			if doc.NumberOfPages() != docWithoutText.NumberOfPages() {
				t.Errorf("expected %d pages, got %d", doc.NumberOfPages(), docWithoutText.NumberOfPages())
			}

			for i := 1; i <= docWithoutText.NumberOfPages(); i++ {
				origPageAny, err := doc.GetPage(i)
				if err != nil {
					t.Errorf("GetPage(original, %d): %v", i, err)
					continue
				}

				strippedPageAny, err := docWithoutText.GetPage(i)
				if err != nil {
					t.Errorf("GetPage(stripped, %d): %v", i, err)
					continue
				}

				origPage, ok := origPageAny.(interface{ Text() string })
				if !ok {
					t.Errorf("page %d: original page does not implement Text()", i)
					continue
				}

				strippedPage, ok := strippedPageAny.(interface{ Text() string })
				if !ok {
					t.Errorf("page %d: stripped page does not implement Text()", i)
					continue
				}

				origText := origPage.Text()
				strippedText := strippedPage.Text()

				if origText == "" {
					t.Errorf("page %d: original document text is empty", i)
				}

				if strippedText != "" {
					t.Errorf("page %d: expected empty text after removal, got %q", i, strippedText)
				}
			}
		})
	}
}
