//go:build benchmark

// Package benchmarks provides system font finder benchmarks.
package benchmarks

import (
	"path/filepath"
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/fonts/systemfonts"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
)

func BenchmarkARVE2745540212Open(b *testing.B) {
	path := filepath.Join(benchmarkDocRoot, "iizieileamidagi.ARVE_2745540212.pdf")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
		if err != nil {
			b.Fatalf("OpenFile(%q): %v", path, err)
		}

		var letters []*content.Letter
		for j := 0; j < doc.NumberOfPages(); j++ {
			pageAny, err := doc.GetPage(j + 1)
			if err != nil {
				b.Fatalf("GetPage(%d): %v", j+1, err)
			}

			page, ok := pageAny.(*content.Page)
			if !ok {
				b.Fatalf("GetPage(%d): expected *content.Page, got %T", j+1, pageAny)
			}

			letters = append(letters, page.Letters()...)
		}
		doc.Close()
	}
}

func BenchmarkARVE2745540212GetTrueTypeFont(b *testing.B) {
	path := filepath.Join(benchmarkDocRoot, "iizieileamidagi.ARVE_2745540212.pdf")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{UseLenientParsing: false})
		if err != nil {
			b.Fatalf("OpenFile(%q): %v", path, err)
		}

		var fonts []*truetypeparser.TrueTypeFont
		for j := 0; j < doc.NumberOfPages(); j++ {
			pageAny, err := doc.GetPage(j + 1)
			if err != nil {
				b.Fatalf("GetPage(%d): %v", j+1, err)
			}

			page, ok := pageAny.(*content.Page)
			if !ok {
				b.Fatalf("GetPage(%d): expected *content.Page, got %T", j+1, pageAny)
			}

			for _, letter := range page.Letters() {
				font := systemfonts.Instance.GetTrueTypeFont(letter.FontName())
				if font != nil {
					fonts = append(fonts, font)
				}
			}
		}
		doc.Close()
	}
}
