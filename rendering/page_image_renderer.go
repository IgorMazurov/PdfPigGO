package rendering

import "github.com/uglytoad/pdfpig/go/content"

// PageImageRenderer renders a PDF page as an image.
type PageImageRenderer interface {
	// Render renders the given page to an image byte slice at the specified scale and format.
	Render(page *content.Page, scale float64, imageFormat PdfRendererImageFormat) ([]byte, error)
}
