//go:build benchmark

// Package benchmarks provides layout analysis benchmarks.
package benchmarks

import (
	"path/filepath"
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
)

var (
	layoutLetters []*content.Letter
	layoutWords   []*content.Word
)

func init() {
	path := filepath.Join(benchmarkDocRoot, "fseprd1102849.pdf")

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	doc, err := pdfpig.OpenFile(absPath, &content.ParsingOptions{UseLenientParsing: false})
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		panic(err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		panic("expected *content.Page")
	}

	layoutLetters = page.Letters()
	layoutWords = word_extractor.DefaultInstance.GetWords(layoutLetters)
}

func BenchmarkGetWordsNearestNeighbour(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = word_extractor.DefaultInstance.GetWords(layoutLetters)
	}
}

func BenchmarkGetBlocksDocstrum(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = page_segmenter.Instance.GetBlocks(layoutWords)
	}
}

func BenchmarkDuplicateOverlappingText(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = document_layout_analysis.Get(layoutLetters)
	}
}
