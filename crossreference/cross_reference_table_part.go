package crossreference

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CrossReferenceTablePart represents a single part of the PDF cross-reference table.
// An in-use entry has the format: nnnnnnnnnn gggggg n eol, where nnnnnnnnnn is a
// 10-digit byte offset (zero-padded), gggggg is a 5-digit generation number,
// n is the literal keyword 'n', and eol is a 2-character end-of-line sequence.
type CrossReferenceTablePart struct {
	objectOffsets      map[core.IndirectReference]int64
	offset             int64
	previous           int64
	dictionary         *tokens.DictionaryToken
	fileType           CrossReferenceType
	tiedToXrefAtOffset *int64
}

// NewCrossReferenceTablePart creates a new CrossReferenceTablePart.
func NewCrossReferenceTablePart(
	objectOffsets map[core.IndirectReference]int64,
	offset int64,
	previous int64,
	dictionary *tokens.DictionaryToken,
	fileType CrossReferenceType,
	tiedToXrefAtOffset *int64,
) *CrossReferenceTablePart {
	return &CrossReferenceTablePart{
		objectOffsets:      objectOffsets,
		offset:             offset,
		previous:           previous,
		dictionary:         dictionary,
		fileType:           fileType,
		tiedToXrefAtOffset: tiedToXrefAtOffset,
	}
}

// ObjectOffsets returns the mapping from indirect references to their byte offsets.
func (p *CrossReferenceTablePart) ObjectOffsets() map[core.IndirectReference]int64 {
	return p.objectOffsets
}

// Offset returns the byte offset of this cross-reference table part in the file.
func (p *CrossReferenceTablePart) Offset() int64 {
	return p.offset
}

// Previous returns the offset of the previous cross-reference section, or -1 if none.
func (p *CrossReferenceTablePart) Previous() int64 {
	return p.previous
}

// Dictionary returns the dictionary token associated with this table part.
func (p *CrossReferenceTablePart) Dictionary() *tokens.DictionaryToken {
	return p.dictionary
}

// Type returns whether this is a cross-reference table or stream section.
func (p *CrossReferenceTablePart) Type() CrossReferenceType {
	return p.fileType
}

// TiedToXrefAtOffset returns the offset of an XRef stream that should be used
// together with this table part when constructing the final cross-reference table,
// or nil if there is no tied XRef.
func (p *CrossReferenceTablePart) TiedToXrefAtOffset() *int64 {
	return p.tiedToXrefAtOffset
}

// FixOffset updates the offset value and patches the Prev entry in the dictionary.
func (p *CrossReferenceTablePart) FixOffset(offset int64) {
	p.offset = offset
	p.dictionary = p.dictionary.With(tokens.Prev, tokens.NewNumericTokenFromInt(int(offset)))
}

// GetPreviousOffset reads the previous xref offset from the dictionary's Prev entry.
// Returns -1 if no Prev entry exists or it is not a numeric token.
func (p *CrossReferenceTablePart) GetPreviousOffset() int64 {
	if p.dictionary == nil {
		return -1
	}

	token, ok := p.dictionary.TryGet(tokens.Prev)
	if !ok {
		return -1
	}

	if numeric, ok := token.(*tokens.NumericToken); ok {
		return numeric.LongVal()
	}

	return -1
}
