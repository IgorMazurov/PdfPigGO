// Package main demonstrates advanced text extraction using word extraction,
// page segmentation, and reading order detection.
package main

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/text/unicode/norm"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/page_segmenter"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/word_extractor"
	pdfpig "github.com/uglytoad/pdfpig/go"
)

// AdvancedTextExtraction performs advanced text extraction on the PDF at filePath.
// It extracts words using a nearest-neighbour word extractor with custom filtering,
// segments the page into blocks using the Docstrum algorithm, determines reading
// order with an unsupervised detector, and prints the normalized extracted text.
func AdvancedTextExtraction(filePath string) error {
	doc, err := pdfpig.OpenFile(filePath, nil)
	if err != nil {
		return fmt.Errorf("opening PDF: %w", err)
	}

	pagesAny, err := doc.GetPages()
	if err != nil {
		return fmt.Errorf("getting pages: %w", err)
	}

	var sb strings.Builder

	for _, p := range pagesAny {
		page, ok := p.(*content.Page)
		if !ok {
			continue
		}

		// 0. Preprocessing
		letters := page.Letters() // no preprocessing

		// 1. Extract words with custom filter options
		opts := word_extractor.DefaultNearestNeighbourWordExtractorOptions()
		opts.SetFilter(func(pivot, candidate *content.Letter) bool {
			// Check if whitespace (default implementation of Filter)
			if strings.TrimSpace(candidate.Value) == "" {
				return false
			}

			// Check for height difference
			maxHeight := math.Max(pivot.PointSize, candidate.PointSize)
			minHeight := math.Min(pivot.PointSize, candidate.PointSize)
			if minHeight != 0 && maxHeight/minHeight > 2.0 {
				return false
			}

			// Check for colour difference
			pivotRGB := pivot.Color.ToRGBValues()
			candidateRGB := candidate.Color.ToRGBValues()
			if pivotRGB.R != candidateRGB.R || pivotRGB.G != candidateRGB.G || pivotRGB.B != candidateRGB.B {
				return false
			}

			return true
		})

		extractor := word_extractor.NewNearestNeighbourWordExtractorWithOptions(opts)
		words := extractor.GetWords(letters)

		// 2. Segment page into text blocks
		segmenterOpts := page_segmenter.DefaultDocstrumBoundingBoxesOptions()
		segmenter := page_segmenter.NewDocstrumBoundingBoxes(segmenterOpts)
		textBlocks := segmenter.GetBlocks(words)

		if len(textBlocks) == 0 {
			continue
		}

		// 3. Determine reading order
		orderedTextBlocks := reading_order_detector.UnsupervisedInstance.Get(textBlocks)

		// 4. Extract text from ordered blocks
		for _, block := range orderedTextBlocks {
			sb.WriteString(norm.NFKC.String(block.Text)) // normalize text
			sb.WriteString("\n")
		}

		sb.WriteString("\n")
	}

	fmt.Print(sb.String())
	return nil
}


