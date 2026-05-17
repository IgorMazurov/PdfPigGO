package crossreference

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CrossReferenceTablePartBuilder accumulates object offset entries and metadata
// for a single part of the PDF cross-reference table, then produces a
// CrossReferenceTablePart via Build.
type CrossReferenceTablePartBuilder struct {
	objects                map[core.IndirectReference]int64
	offset                 int64
	previous               int64
	dictionary             *tokens.DictionaryToken
	xrefType               CrossReferenceType
	tiedToPreviousAtOffset *int64
}

// NewCrossReferenceTablePartBuilder creates a new builder with an empty object map.
func NewCrossReferenceTablePartBuilder() *CrossReferenceTablePartBuilder {
	return &CrossReferenceTablePartBuilder{
		objects: make(map[core.IndirectReference]int64),
	}
}

// Offset gets or sets the byte offset of this cross-reference table part in the file.
func (b *CrossReferenceTablePartBuilder) Offset() int64 {
	return b.offset
}

// SetOffset sets the byte offset of this cross-reference table part.
func (b *CrossReferenceTablePartBuilder) SetOffset(offset int64) {
	b.offset = offset
}

// Previous gets or sets the offset of the previous cross-reference section.
func (b *CrossReferenceTablePartBuilder) Previous() int64 {
	return b.previous
}

// SetPrevious sets the offset of the previous cross-reference section.
func (b *CrossReferenceTablePartBuilder) SetPrevious(previous int64) {
	b.previous = previous
}

// Dictionary gets or sets the dictionary token associated with this table part.
func (b *CrossReferenceTablePartBuilder) Dictionary() *tokens.DictionaryToken {
	return b.dictionary
}

// SetDictionary sets the dictionary token for this table part.
func (b *CrossReferenceTablePartBuilder) SetDictionary(dict *tokens.DictionaryToken) {
	b.dictionary = dict
}

// XRefType gets or sets whether this is a cross-reference table or stream section.
func (b *CrossReferenceTablePartBuilder) XRefType() CrossReferenceType {
	return b.xrefType
}

// SetXRefType sets the type of this cross-reference section.
func (b *CrossReferenceTablePartBuilder) SetXRefType(xrefType CrossReferenceType) {
	b.xrefType = xrefType
}

// TiedToPreviousAtOffset gets or sets the offset of an XRef stream tied to this part,
// or nil if there is no tied XRef.
func (b *CrossReferenceTablePartBuilder) TiedToPreviousAtOffset() *int64 {
	return b.tiedToPreviousAtOffset
}

// SetTiedToPreviousAtOffset sets the offset of a tied XRef stream.
func (b *CrossReferenceTablePartBuilder) SetTiedToPreviousAtOffset(offset *int64) {
	b.tiedToPreviousAtOffset = offset
}

// Add registers an object at the given byte offset. If generationNumber exceeds
// ushort.MaxValue (65535), the entry is silently skipped as invalid. Duplicate
// keys are ignored — the first offset wins.
func (b *CrossReferenceTablePartBuilder) Add(objectID int64, generationNumber int, offset int64) {
	if generationNumber > 0xFFFF {
		return
	}

	objKey, err := core.NewIndirectReference(objectID, generationNumber)
	if err != nil {
		return
	}

	if _, exists := b.objects[objKey]; !exists {
		b.objects[objKey] = offset
	}
}

// Build produces a CrossReferenceTablePart from the accumulated data.
func (b *CrossReferenceTablePartBuilder) Build() *CrossReferenceTablePart {
	return NewCrossReferenceTablePart(
		b.objects,
		b.offset,
		b.previous,
		b.dictionary,
		b.xrefType,
		b.tiedToPreviousAtOffset,
	)
}
