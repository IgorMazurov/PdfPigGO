package filestructure

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// XrefStream represents an xref stream object in a PDF file.
type XrefStream struct {
	offset          int64
	objectOffsets   map[core.IndirectReference]core.XrefLocation
	dictionary      *tokens.DictionaryToken
	correctionType  XrefOffsetCorrection
	offsetCorrection int64
}

var _ XrefSection = (*XrefStream)(nil)

// NewXrefStream creates a new XrefStream instance.
func NewXrefStream(
	offset int64,
	objectOffsets map[core.IndirectReference]core.XrefLocation,
	streamDictionary *tokens.DictionaryToken,
	correctionType XrefOffsetCorrection,
	offsetCorrection int64,
) *XrefStream {
	return &XrefStream{
		offset:          offset,
		objectOffsets:   objectOffsets,
		dictionary:      streamDictionary,
		correctionType:  correctionType,
		offsetCorrection: offsetCorrection,
	}
}

// Offset returns the byte offset of this xref in the file.
func (x *XrefStream) Offset() int64 {
	return x.offset
}

// ObjectOffsets returns the byte offsets of the objects keyed by indirect reference.
func (x *XrefStream) ObjectOffsets() map[core.IndirectReference]core.XrefLocation {
	return x.objectOffsets
}

// Dictionary returns the stream dictionary token.
func (x *XrefStream) Dictionary() *tokens.DictionaryToken {
	return x.dictionary
}

// CorrectionType indicates how this xref was located if a correction was applied.
func (x *XrefStream) CorrectionType() XrefOffsetCorrection {
	return x.correctionType
}

// OffsetCorrection returns the number of bytes from the original location
// when applying a correction, or zero if no correction was needed.
func (x *XrefStream) OffsetCorrection() int64 {
	return x.offsetCorrection
}

// GetPrevious returns the offset of the previous xref section, or nil if none exists.
//go:noinline
func (x *XrefStream) GetPrevious() *int64 {
	if x == nil || x.dictionary == nil {
		return nil
	}
	if token, ok := x.dictionary.TryGet(tokens.Prev); ok {
		if numeric, ok := token.(*tokens.NumericToken); ok {
			val := numeric.LongVal()
			return &val
		}
	}
	return nil
}
