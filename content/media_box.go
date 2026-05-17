package content

import (
	"github.com/uglytoad/pdfpig/go/core"
)

// PointsPerInch is the number of PDF user space units per inch.
const PointsPerInch = 72.0

// PointsPerMm is the number of PDF user space units per millimeter.
const PointsPerMm = PointsPerInch / (10 * 2.54)

// MediaBox defines the boundaries of the physical medium on which the page shall be displayed or printed.
// See table 3.27 from the PDF specification version 1.7.
type MediaBox struct {
	Bounds core.PdfRectangle
}

// NewMediaBox creates a new MediaBox with the given bounds.
func NewMediaBox(bounds core.PdfRectangle) *MediaBox {
	return &MediaBox{Bounds: bounds}
}

// MediaBoxUSLetter is a MediaBox representing U.S. Letter size (8.5" x 11").
var MediaBoxUSLetter = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 8.5*PointsPerInch, 11*PointsPerInch)}

// MediaBoxUSLegal is a MediaBox representing U.S. Legal size (8.5" x 14").
var MediaBoxUSLegal = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 8.5*PointsPerInch, 14*PointsPerInch)}

// MediaBoxA0 is a MediaBox representing A0 paper size (841mm x 1189mm).
var MediaBoxA0 = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 841*PointsPerMm, 1189*PointsPerMm)}

// MediaBoxA1 is a MediaBox representing A1 paper size (594mm x 841mm).
var MediaBoxA1 = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 594*PointsPerMm, 841*PointsPerMm)}

// MediaBoxA2 is a MediaBox representing A2 paper size (420mm x 594mm).
var MediaBoxA2 = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 420*PointsPerMm, 594*PointsPerMm)}

// MediaBoxA3 is a MediaBox representing A3 paper size (297mm x 420mm).
var MediaBoxA3 = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 297*PointsPerMm, 420*PointsPerMm)}

// MediaBoxA4 is a MediaBox representing A4 paper size (210mm x 297mm).
var MediaBoxA4 = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 210*PointsPerMm, 297*PointsPerMm)}

// MediaBoxA5 is a MediaBox representing A5 paper size (148mm x 210mm).
var MediaBoxA5 = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 148*PointsPerMm, 210*PointsPerMm)}

// MediaBoxA6 is a MediaBox representing A6 paper size (105mm x 148mm).
var MediaBoxA6 = &MediaBox{Bounds: core.NewPdfRectangleFloat(0, 0, 105*PointsPerMm, 148*PointsPerMm)}
