package writer

import "fmt"

// PdfAStandard represents the PDF/A compliance standard for generated documents.
type PdfAStandard int

const (
	// PdfANone means no PDF/A compliance.
	PdfANone PdfAStandard = iota

	// PdfA1B is compliance with PDF/A-1b (basic). Level B conformance covers
	// standards necessary for the reliable reproduction of a document's visual appearance.
	PdfA1B

	// PdfA1A is compliance with PDF/A-1a (accessible). Includes PDF/A-1b
	// standards plus features intended to improve document accessibility.
	PdfA1A

	// PdfA2B is compliance with PDF/A-2b (basic). Level B conformance covers
	// standards necessary for the reliable reproduction of a document's visual appearance.
	PdfA2B

	// PdfA2A is compliance with PDF/A-2a (accessible). Includes PDF/A-2b
	// standards plus features intended to improve document accessibility.
	PdfA2A

	// PdfA3B is compliance with PDF/A-3b (basic). Extends PDF/A-2b with support
	// for embedded files.
	PdfA3B

	// PdfA3A is compliance with PDF/A-3a (accessible). Includes PDF/A-3b
	// standards plus features intended to improve document accessibility.
	PdfA3A
)

// String returns the name of the PdfAStandard constant.
func (s PdfAStandard) String() string {
	switch s {
	case PdfANone:
		return "None"
	case PdfA1B:
		return "A1B"
	case PdfA1A:
		return "A1A"
	case PdfA2B:
		return "A2B"
	case PdfA2A:
		return "A2A"
	case PdfA3B:
		return "A3B"
	case PdfA3A:
		return "A3A"
	default:
		return fmt.Sprintf("PdfAStandard(%d)", s)
	}
}
