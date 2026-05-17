package core

// TextRenderingMode determines whether showing text causes glyph outlines to be stroked, filled,
// used as a clipping boundary, or some combination of the three.
type TextRenderingMode byte

const (
	// FillText fills the entire letter region.
	FillText TextRenderingMode = 0

	// StrokeText draws the border/outline of the letter.
	StrokeText TextRenderingMode = 1

	// FillThenStrokeText fills then strokes the text.
	FillThenStrokeText TextRenderingMode = 2

	// NeitherFillNorStroke makes text invisible by neither filling nor stroking.
	NeitherFillNorStroke TextRenderingMode = 3

	// FillAndClip fills the text and adds it to the clipping path.
	FillAndClip TextRenderingMode = 4

	// StrokeAndClip strokes the text and adds it to the clipping path.
	StrokeAndClip TextRenderingMode = 5

	// FillThenStrokeAndClip fills then strokes the text and adds it to the clipping path.
	FillThenStrokeAndClip TextRenderingMode = 6

	// NeitherFillNorStrokeButClip makes no fill nor stroke but adds to the clipping path.
	NeitherFillNorStrokeButClip TextRenderingMode = 7
)
