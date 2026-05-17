//go:build benchmark

// Package benchmarks provides brute-force letter extraction benchmarks.
package benchmarks

import (
	"path/filepath"
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
)

const benchmarkDocRoot = "../testdata/integration/Documents"

func BenchmarkOpenOfficeLetters(b *testing.B) {
	path := filepath.Join(benchmarkDocRoot, "Single Page Simple - from open office.pdf")

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

func BenchmarkInkscapeLetters(b *testing.B) {
	path := filepath.Join(benchmarkDocRoot, "Single Page Simple - from inkscape.pdf")

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

func BenchmarkAlgoLetters(b *testing.B) {
	path := filepath.Join(benchmarkDocRoot, "algo.pdf")

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

func BenchmarkPDFBOX492Letters(b *testing.B) {
	path := filepath.Join(benchmarkDocRoot, "PDFBOX-492-4.jar-8.pdf")

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
