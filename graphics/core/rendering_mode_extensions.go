package core

import pdfpigcore "github.com/uglytoad/pdfpig/go/core"

// IsFill returns true if the rendering mode includes filling text.
func IsFill(mode pdfpigcore.TextRenderingMode) bool {
	return mode == pdfpigcore.FillText ||
		mode == pdfpigcore.FillThenStrokeText ||
		mode == pdfpigcore.FillAndClip ||
		mode == pdfpigcore.FillThenStrokeAndClip
}

// IsStroke returns true if the rendering mode includes stroking text.
func IsStroke(mode pdfpigcore.TextRenderingMode) bool {
	return mode == pdfpigcore.StrokeText ||
		mode == pdfpigcore.FillThenStrokeText ||
		mode == pdfpigcore.StrokeAndClip ||
		mode == pdfpigcore.FillThenStrokeAndClip
}

// IsClip returns true if the rendering mode adds text to the clipping path.
func IsClip(mode pdfpigcore.TextRenderingMode) bool {
	return mode == pdfpigcore.FillAndClip ||
		mode == pdfpigcore.StrokeAndClip ||
		mode == pdfpigcore.FillThenStrokeAndClip ||
		mode == pdfpigcore.NeitherFillNorStrokeButClip
}
