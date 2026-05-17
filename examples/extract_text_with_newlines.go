// Package main demonstrates how to extract text from a PDF document with newlines,
// preserving line breaks based on letter positioning in the original document.
package main

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/text_extractor"
)

// ExtractTextWithNewlines extracts text from the PDF at filePath using content-order extraction
// with paragraph separation, printing each page's text to stdout.
func ExtractTextWithNewlines(filePath string) error {
	doc, err := pdfpig.OpenFile(filePath, nil)
	if err != nil {
		return fmt.Errorf("opening PDF: %w", err)
	}
	defer doc.Close()

	pagesAny, err := doc.GetPages()
	if err != nil {
		return fmt.Errorf("getting pages: %w", err)
	}

	for _, p := range pagesAny {
		page, ok := p.(*content.Page)
		if !ok {
			continue
		}

		text := text_extractor.GetText(page, true)
		fmt.Println(text)
	}

	return nil
}


