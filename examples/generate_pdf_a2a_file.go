// Package main demonstrates generating a PDF/A-2a compliant document with text and an image.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/writer"
)

// GeneratePdfA2AFile generates a PDF/A-2a document containing text and a JPEG image.
func GeneratePdfA2AFile(trueTypeFontPath, jpgImagePath string) error {
	fontBytes, err := os.ReadFile(trueTypeFontPath)
	if err != nil {
		return fmt.Errorf("reading font file: %w", err)
	}

	jpgBytes, err := os.ReadFile(jpgImagePath)
	if err != nil {
		return fmt.Errorf("reading image file: %w", err)
	}

	builder := writer.NewPdfDocumentBuilder()
	builder.SetArchiveStandard(writer.PdfA2A)

	font, err := builder.AddTrueTypeFont(fontBytes)
	if err != nil {
		return fmt.Errorf("adding TrueType font: %w", err)
	}

	pageSize := content.PageSizeA4
	dims, ok := pageSize.TryGetPdfRectangle()
	if !ok {
		return fmt.Errorf("unsupported page size")
	}

	page, err := builder.AddPageWithSize(dims.Width, dims.Height)
	if err != nil {
		return fmt.Errorf("adding page: %w", err)
	}

	pageTop := core.NewPdfPoint(0, page.PageSize().Bounds.Top())

	letters, err := page.AddText(
		"This is some text added to the output file near the top of the page.",
		12,
		pageTop.Translate(20, -25),
		font,
	)
	if err != nil {
		return fmt.Errorf("adding text: %w", err)
	}

	bottomOfText := letters[0].GlyphRectangle().Bottom()
	for _, l := range letters[1:] {
		b := l.GlyphRectangle().Bottom()
		if b < bottomOfText {
			bottomOfText = b
		}
	}

	imagePlacement := core.NewPdfRectangle(
		core.NewPdfPoint(50, bottomOfText-200),
		core.NewPdfPoint(150, bottomOfText),
	)

	_, err = page.AddJpeg(jpgBytes, imagePlacement)
	if err != nil {
		return fmt.Errorf("adding JPEG: %w", err)
	}

	fileBytes, err := builder.Build()
	if err != nil {
		return fmt.Errorf("building document: %w", err)
	}

	outputPath := filepath.Join(".", "outputOfPdfA2A.pdf")
	err = os.WriteFile(outputPath, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("File output to: %s\n", outputPath)
	return nil
}


