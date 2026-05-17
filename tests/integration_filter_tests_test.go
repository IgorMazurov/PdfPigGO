//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// noFilter is a filter that reports itself as unsupported and always fails decoding.
type noFilter struct{}

func (f *noFilter) IsSupported() bool {
	return false
}

func (f *noFilter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider filters.FilterProvider, filterIndex int) ([]byte, error) {
	return nil, nil
}

var _ filters.Filter = (*noFilter)(nil)

// myFilterProvider is a custom filter provider that disables image-decoding filters
// (CcittFax, Dct, Jbig2, Jpx) while keeping structural filters functional.
type myFilterProvider struct {
	*filters.BaseFilterProvider
}

var myFilterProviderInstance = newMyFilterProvider()

func newMyFilterProvider() *myFilterProvider {
	no := &noFilter{}
	ascii85 := filters.NewAscii85Filter()
	asciiHex := filters.NewAsciiHexDecodeFilter()
	flate := filters.NewFlateFilter()
	runLength := filters.NewRunLengthFilter()
	lzw := filters.NewLzwFilter()

	dict := map[string]filters.Filter{
		tokens.Ascii85Decode.Data():              ascii85,
		tokens.Ascii85DecodeAbbreviation.Data():  ascii85,
		tokens.AsciiHexDecode.Data():             asciiHex,
		tokens.AsciiHexDecodeAbbreviation.Data(): asciiHex,
		tokens.CcittfaxDecode.Data():             no,
		tokens.CcittfaxDecodeAbbreviation.Data(): no,
		tokens.DctDecode.Data():                  no,
		tokens.DctDecodeAbbreviation.Data():      no,
		tokens.FlateDecode.Data():                flate,
		tokens.FlateDecodeAbbreviation.Data():    flate,
		tokens.Jbig2Decode.Data():                no,
		tokens.JpxDecode.Data():                  no,
		tokens.RunLengthDecode.Data():            runLength,
		tokens.RunLengthDecodeAbbreviation.Data(): runLength,
		tokens.LzwDecode.Data():                  lzw,
		tokens.LzwDecodeAbbreviation.Data():      lzw,
	}

	return &myFilterProvider{
		BaseFilterProvider: filters.NewBaseFilterProvider(dict),
	}
}

// TestNoImageDecoding verifies that images encoded with unsupported filters cannot be
// decoded to PNG when using a filter provider that only supports FlateDecode and LzwDecode
// for image data. This matches the C# FilterTests.NoImageDecoding test.
func TestNoImageDecoding(t *testing.T) {
	documents := getAllDocuments(t)

	for _, documentName := range documents {
		t.Run(documentName, func(t *testing.T) {
			fullPath := resolveFilterTestDocumentPath(t, documentName)

			opts := &content.ParsingOptions{
				UseLenientParsing: true,
				FilterProvider:    myFilterProviderInstance,
			}

			doc, err := pdfpig.OpenFile(fullPath, opts)
			if err != nil {
				t.Fatalf("pdfpig.OpenFile(%q): %v", fullPath, err)
			}
			defer doc.Close()

			for i := 0; i < doc.NumberOfPages(); i++ {
				pageAny, err := doc.GetPage(i + 1)
				if err != nil {
					t.Fatalf("GetPage(%d): %v", i+1, err)
				}

				page, ok := pageAny.(*content.Page)
				if !ok {
					t.Fatalf("GetPage(%d): expected *content.Page, got %T", i+1, pageAny)
				}

				for _, pdfImage := range page.GetImages() {
					filterToken, found := pdfImage.ImageDictionary().TryGet(tokens.Filter)
					if !found {
						continue
					}

					nameToken, isName := filterToken.(*tokens.NameToken)
					if !isName {
						continue
					}

					filterData := nameToken.Data()
					if filterData == tokens.FlateDecode.Data() ||
						filterData == tokens.FlateDecodeAbbreviation.Data() ||
						filterData == tokens.LzwDecode.Data() ||
						filterData == tokens.LzwDecodeAbbreviation.Data() {
						continue
					}

					if _, ok := pdfImage.TryGetPng(); ok {
						t.Errorf("expected image with filter %q to fail PNG decoding, but it succeeded in document %s", filterData, documentName)
					}
				}
			}
		})
	}
}

func getAllDocuments(t *testing.T) []string {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(integrationDocRoot, "*.pdf"))
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}

	ignored := map[string]bool{
		"issue_671.pdf":            true,
		"GHOSTSCRIPT-698363-0.pdf": true,
		"ErcotFacts.pdf":           true,
	}

	result := make([]string, 0, len(files))
	for _, f := range files {
		name := filepath.Base(f)
		if !ignored[name] {
			result = append(result, name)
		}
	}

	return result
}

func resolveFilterTestDocumentPath(t *testing.T, documentName string) string {
	t.Helper()

	path := filepath.Join(integrationDocRoot, documentName)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("document not found: %s", path)
	}

	return path
}
