package rendering

// PdfRendererImageFormat represents the output image format of a page renderer.
type PdfRendererImageFormat byte

const (
	// Bmp is the bitmap image format.
	Bmp PdfRendererImageFormat = iota
	// Jpeg is the JPEG/JPG image format.
	Jpeg
	// Png is the PNG image format.
	Png
	// Tiff is the TIFF image format.
	Tiff
	// Gif is the GIF image format.
	Gif
)
