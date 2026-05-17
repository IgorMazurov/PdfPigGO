package filestructure

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// XrefOffsetCorrection indicates how an xref section was located in the file.
type XrefOffsetCorrection byte

const (
	// XrefOffsetCorrectionNone means the xref was found at exactly the
	// specified byte offset in the file.
	XrefOffsetCorrectionNone XrefOffsetCorrection = iota
	// XrefOffsetCorrectionFileHeaderOffset means the xref was shifted by
	// the offset of the version header start comment in the file.
	XrefOffsetCorrectionFileHeaderOffset
	// XrefOffsetCorrectionRandom means the xref was not at the correct
	// location but was found nearby.
	XrefOffsetCorrectionRandom
)

// XrefSection represents a cross-reference section in a PDF file, either
// a traditional xref table or an xref stream object.
type XrefSection interface {
	// Offset returns the byte offset of this xref in the file. For tables
	// this is the position of the "xref" operator; for stream objects it
	// is the start of the object number marker, e.g., "14 0 obj".
	Offset() int64

	// ObjectOffsets returns the byte offsets of the objects in this xref.
	ObjectOffsets() map[core.IndirectReference]core.XrefLocation

	// Dictionary returns the dictionary for this xref. For the trailer
	// xref this is the trailer dictionary; for streams it is the stream
	// dictionary. May be nil.
	Dictionary() *tokens.DictionaryToken

	// GetPrevious returns the offset of the previous xref section, or nil
	// if there is none.
	GetPrevious() *int64

	// CorrectionType indicates how this xref was located if a correction
	// had to be applied.
	CorrectionType() XrefOffsetCorrection

	// OffsetCorrection returns how many bytes from the original location
	// we had to move when applying a correction, or 0 if no correction
	// was needed.
	OffsetCorrection() int64
}
