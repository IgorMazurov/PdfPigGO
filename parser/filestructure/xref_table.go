package filestructure

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// XrefTable represents a traditional cross-reference table in a PDF file.
type XrefTable struct {
	offset           int64
	objectOffsets    map[core.IndirectReference]core.XrefLocation
	dictionary       *tokens.DictionaryToken
	correctionType   XrefOffsetCorrection
	offsetCorrection int64
}

var _ XrefSection = (*XrefTable)(nil)

// NewXrefTable creates a new XrefTable instance.
func NewXrefTable(
	offset int64,
	objectOffsets map[core.IndirectReference]core.XrefLocation,
	trailer *tokens.DictionaryToken,
	correctionType XrefOffsetCorrection,
	offsetCorrection int64,
) *XrefTable {
	return &XrefTable{
		offset:           offset,
		objectOffsets:    objectOffsets,
		dictionary:       trailer,
		correctionType:   correctionType,
		offsetCorrection: offsetCorrection,
	}
}

// Offset returns the byte offset of the "xref" operator in the file.
func (x *XrefTable) Offset() int64 {
	return x.offset
}

// ObjectOffsets returns the keyed object-to-byte-offset mapping for this document.
func (x *XrefTable) ObjectOffsets() map[core.IndirectReference]core.XrefLocation {
	return x.objectOffsets
}

// Dictionary returns the trailer dictionary token, or nil if absent.
func (x *XrefTable) Dictionary() *tokens.DictionaryToken {
	return x.dictionary
}

// CorrectionType indicates how this xref was located if a correction was applied.
func (x *XrefTable) CorrectionType() XrefOffsetCorrection {
	return x.correctionType
}

// OffsetCorrection returns the number of bytes from the original location
// when applying a correction, or zero if no correction was needed.
func (x *XrefTable) OffsetCorrection() int64 {
	return x.offsetCorrection
}

// GetPrevious returns the offset of the previous xref section, or nil if none exists.
//go:noinline
func (x *XrefTable) GetPrevious() *int64 {
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

// GetXRefStm returns the byte offset of the xref stream, or nil if absent.
//go:noinline
func (x *XrefTable) GetXRefStm() *int64 {
	if x == nil || x.dictionary == nil {
		return nil
	}
	if token, ok := x.dictionary.TryGet(tokens.XrefStm); ok {
		if numeric, ok := token.(*tokens.NumericToken); ok {
			val := numeric.LongVal()
			return &val
		}
	}
	return nil
}
