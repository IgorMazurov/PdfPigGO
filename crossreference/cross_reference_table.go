package crossreference

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

// CrossReferenceTable contains information that enables random access to PDF objects
// within the file by object number so that specific objects can be located directly
// without having to scan the whole document. A PDF document may contain multiple cross
// reference tables; this struct provides access to the merged result with the latest
// offset for each object. The offsets of the original cross reference tables or streams
// merged into this result are available in the CrossReferenceOffsets slice.
type CrossReferenceTable struct {
	objectOffsets       map[core.IndirectReference]int64
	fileType            CrossReferenceType
	trailer             *TrailerDictionary
	crossReferenceOffs  []CrossReferenceOffset
}

// NewCrossReferenceTable creates a new CrossReferenceTable.
func NewCrossReferenceTable(
	fileType CrossReferenceType,
	objectOffsets map[core.IndirectReference]int64,
	trailer *TrailerDictionary,
	crossReferenceOffs []CrossReferenceOffset,
) (*CrossReferenceTable, error) {
	if objectOffsets == nil {
		return nil, fmt.Errorf("objectOffsets must not be nil")
	}

	if trailer == nil {
		return nil, fmt.Errorf("trailer must not be nil")
	}

	result := make(map[core.IndirectReference]int64, len(objectOffsets))
	for k, v := range objectOffsets {
		result[k] = v
	}

	return &CrossReferenceTable{
		objectOffsets:       result,
		fileType:            fileType,
		trailer:             trailer,
		crossReferenceOffs:  crossReferenceOffs,
	}, nil
}

// ObjectOffsets returns the corresponding byte offset for each keyed object in this document.
func (t *CrossReferenceTable) ObjectOffsets() map[core.IndirectReference]int64 {
	return t.objectOffsets
}

// Type returns the type of the first cross-reference table located in this document.
func (t *CrossReferenceTable) Type() CrossReferenceType {
	return t.fileType
}

// Trailer returns the trailer dictionary.
func (t *CrossReferenceTable) Trailer() *TrailerDictionary {
	return t.trailer
}

// CrossReferenceOffsets returns the byte offsets of each cross-reference table or stream
// in this document and the previous table or stream they link to if applicable.
func (t *CrossReferenceTable) CrossReferenceOffsets() []CrossReferenceOffset {
	return t.crossReferenceOffs
}

// CrossReferenceOffset represents the offset of a cross-reference table or stream in the document.
type CrossReferenceOffset struct {
	Current  int64
	Previous *int64
}

// NewCrossReferenceOffset creates a new CrossReferenceOffset.
func NewCrossReferenceOffset(current int64, previous *int64) CrossReferenceOffset {
	return CrossReferenceOffset{
		Current:  current,
		Previous: previous,
	}
}

// String returns the string representation of the offset.
func (o CrossReferenceOffset) String() string {
	if o.Previous != nil {
		return fmt.Sprintf("%d %d", o.Current, *o.Previous)
	}
	return fmt.Sprintf("%d", o.Current)
}
