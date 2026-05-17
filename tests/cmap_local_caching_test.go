//go:build integration

package pdfpig_test


import (
	"github.com/uglytoad/pdfpig/go/testutil"
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

var (
	integrationDocRoot = testutil.IntegrationDocumentsRoot
	dlaDocRoot         = filepath.Join(testutil.IntegrationDocumentsRoot, "..", "Dla", "Documents")
)

var cmapCachingDocuments = []string{
	"68-1990-01_A.pdf",
	"Type0 Font.pdf",
	"11194059_2017-11_de_s.pdf",
	"2108.11480.pdf",
	"reference-2-numeric-error.pdf",
	"MOZILLA-3136-0.pdf",
	"FICTIF_TABLE_INDEX.pdf",
	"Approved_Document_B__fire_safety__volume_2_-_Buildings_other_than_dwellings__2019_edition_incorporating_2020_and_2022_amendments.pdf",
	"dotnet-ai.pdf",
	"Old Gutnish Internet Explorer.pdf",
	"Random 2 Columns Lists Hyph - Justified.pdf",
}

func TestCMapLocalCaching_CheckText(t *testing.T) {
	for _, documentName := range cmapCachingDocuments {
		t.Run(documentName, func(t *testing.T) {
			fullPath := resolveDocumentPath(t, documentName)
			expectedPath := strings.TrimSuffix(fullPath, ".pdf") + ".txt"

			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Skipf("expected text file not found: %s", expectedPath)
			}

			opts := &content.ParsingOptions{
				UseLenientParsing: true,
			}

			doc, err := pdfpig.OpenFile(fullPath, opts)
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
			}
			defer doc.Close()

			var sb strings.Builder
			for i := 0; i < doc.NumberOfPages(); i++ {
				pageAny, err := doc.GetPage(i + 1)
				if err != nil {
					t.Fatalf("GetPage(%d): %v", i+1, err)
				}

				page, ok := pageAny.(*content.Page)
				if !ok {
					t.Fatalf("GetPage(%d): expected *content.Page, got %T", i+1, pageAny)
				}

				sb.WriteString(page.Text())
			}

			expectedBytes, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("ReadAll(%q): %v", expectedPath, err)
			}

			got := sb.String()
			expected := string(expectedBytes)

			if got != expected {
				t.Errorf("text mismatch for %s\n--- got (%d bytes):\n%s\n--- expected (%d bytes):\n%s",
					documentName, len(got), got, len(expected), expected)
			}
		})
	}
}

func resolveDocumentPath(t *testing.T, documentName string) string {
	t.Helper()

	candidates := []string{integrationDocRoot, dlaDocRoot}
	for _, dir := range candidates {
		path := filepath.Join(dir, documentName)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}

	t.Fatalf("document not found in any test data directory: %s (searched: %v)",
		documentName, candidates)
	return ""
}
