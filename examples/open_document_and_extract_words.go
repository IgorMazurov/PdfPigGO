// Package main demonstrates how to open a PDF document and extract words,
// inserting newlines when consecutive words have a significant baseline gap
// (indicating a new line in the original document).
package main

import (
	"fmt"
	"math"
	"strings"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
)

// OpenDocumentAndExtractWords opens the PDF at filePath, extracts words from each page, and prints
// the reconstructed text to stdout. Words on different baselines (gap > 3 points)
// are separated by a newline; otherwise they are separated by a space.
func OpenDocumentAndExtractWords(filePath string) error {
	doc, err := pdfpig.OpenFile(filePath, nil)
	if err != nil {
		return fmt.Errorf("opening PDF: %w", err)
	}
	defer doc.Close()

	pagesAny, err := doc.GetPages()
	if err != nil {
		return fmt.Errorf("getting pages: %w", err)
	}

	var sb strings.Builder
	var previous *content.Word

	for _, p := range pagesAny {
		page, ok := p.(*content.Page)
		if !ok {
			continue
		}

		words := page.GetWords()
		for _, word := range words {
			if previous != nil {
				hasInsertedWhitespace := false
				bothNonEmpty := len(previous.Letters) > 0 && len(word.Letters) > 0

				if bothNonEmpty {
					prevLetter1 := previous.Letters[0]
					currentLetter1 := word.Letters[0]

					baselineGap := math.Abs(prevLetter1.StartBaseLine.Y - currentLetter1.StartBaseLine.Y)

					if baselineGap > 3 {
						hasInsertedWhitespace = true
						sb.WriteString("\n")
					}
				}

				if !hasInsertedWhitespace {
					sb.WriteString(" ")
				}
			}

			sb.WriteString(word.Text)
			previous = word
		}
	}

	fmt.Println(sb.String())
	return nil
}


