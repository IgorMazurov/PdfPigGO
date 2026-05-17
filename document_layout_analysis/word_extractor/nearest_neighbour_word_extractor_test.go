package word_extractor_test

import (
	"strings"
	"testing"

	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
	"github.com/uglytoad/pdfpig/go/content"
	pdfpig "github.com/uglytoad/pdfpig/go"
)

type wordCountData struct {
	path            string
	wordCount       int
	noSpacesWordCount int
}

func wordCountTestData() []wordCountData {
	return []wordCountData{
		{"2559 words.pdf", 5118, 2559},
		{"fseprd1102849.pdf", 12903, 11177},
		{"90 180 270 rotated.pdf", 589, 292},
		{"complex rotated.pdf", 805, 403},
		{"no horizontal distance.pdf", 4, 2},
		{"no vertical distance.pdf", 22, 10},
		{"no vertical horizontal distance.pdf", 4, 2},
		{"Random 2 Columns Lists Hyph - Justified.pdf", 1191, 607},
		{"caly-issues-56-1.pdf", 184, 156},
		{"caly-issues-58-2.pdf", 49, 49},
	}
}

func TestWordCount(t *testing.T) {
	data := wordCountTestData()

	for _, d := range data {
		t.Run(d.path, func(t *testing.T) {
			docPath := dla.GetDocumentPath(d.path, true)

			opts := &content.ParsingOptions{
				UseLenientParsing: true,
			}

			doc, err := pdfpig.OpenFile(docPath, opts)
			if err != nil {
				t.Fatalf("OpenFile(%q): %v", docPath, err)
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

			words := word_extractor.DefaultInstance.GetWords(page.Letters())

			if len(words) != d.wordCount {
				t.Errorf("expected %d words, got %d", d.wordCount, len(words))
			}

			noSpacesWords := make([]*content.Word, 0, len(words))
			for _, w := range words {
				if strings.TrimSpace(w.Text) != "" {
					noSpacesWords = append(noSpacesWords, w)
				}
			}

			if len(noSpacesWords) != d.noSpacesWordCount {
				t.Errorf("expected %d non-empty words, got %d", d.noSpacesWordCount, len(noSpacesWords))
			}
		})
	}
}
