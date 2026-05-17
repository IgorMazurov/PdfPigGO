package graphics

import "github.com/uglytoad/pdfpig/go/core"

// TextMatrices manages the text matrix (Tm), text line matrix (Tlm) and is used to generate the text rendering matrix (Trm).
type TextMatrices struct {
	// TextMatrix is the current text matrix (Tm).
	TextMatrix core.TransformationMatrix

	// TextLineMatrix captures the value of the TextMatrix at the beginning of a line of text.
	// This is convenient for aligning evenly spaced lines of text.
	TextLineMatrix core.TransformationMatrix
}

// NewTextMatrices creates a new TextMatrices instance with identity matrices.
func NewTextMatrices() *TextMatrices {
	return &TextMatrices{
		TextMatrix:     core.Identity,
		TextLineMatrix: core.Identity,
	}
}
