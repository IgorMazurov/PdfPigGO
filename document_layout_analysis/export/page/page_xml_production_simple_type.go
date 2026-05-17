package page

// PageXmlProductionSimpleType represents the text production type for PAGE XML export.
type PageXmlProductionSimpleType byte

const (
	// Printed indicates printed text.
	Printed PageXmlProductionSimpleType = iota

	// Typewritten indicates typewritten text.
	Typewritten

	// HandwrittenCursive indicates handwritten cursive text.
	HandwrittenCursive

	// HandwrittenPrintscript indicates handwritten print script text.
	HandwrittenPrintscript

	// MedievalManuscript indicates medieval manuscript text.
	MedievalManuscript

	// OtherProduction indicates an unrecognized production type.
	OtherProduction
)
