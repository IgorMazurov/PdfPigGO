// Package main demonstrates how to extract images from a PDF document.
package main

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	pdfpig "github.com/uglytoad/pdfpig/go"
)

// ExtractImages extracts all images from the PDF at filePath and prints their details.
func ExtractImages(filePath string) error {
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

		images := page.GetImages()
		for _, image := range images {
			b, ok := image.TryGetBytes()
			if !ok {
				b = image.RawMemory()
			}

			imageType := ""
			switch image.(type) {
			case *content.XObjectImage:
				imageType = "XObject"
			case *content.InlineImage:
				imageType = "Inline"
			}

			fmt.Printf("Image with %d bytes of type '%s' on page %d. Location: %v.\n",
				len(b), imageType, page.Number(), image.Bounds())
		}
	}

	return nil
}


