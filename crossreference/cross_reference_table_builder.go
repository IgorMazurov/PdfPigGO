package crossreference

import (
	"fmt"
	"sort"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CrossReferenceTableBuilder collects cross-reference table parts and builds a merged
// CrossReferenceTable from them. The table contains a one-line entry for each indirect
// object, specifying the location of that object within the body of the file.
type CrossReferenceTableBuilder struct {
	parts []*CrossReferenceTablePart
}

// NewCrossReferenceTableBuilder creates a new CrossReferenceTableBuilder.
func NewCrossReferenceTableBuilder() *CrossReferenceTableBuilder {
	return &CrossReferenceTableBuilder{
		parts: make([]*CrossReferenceTablePart, 0),
	}
}

// Parts returns the collected cross-reference table parts.
func (b *CrossReferenceTableBuilder) Parts() []*CrossReferenceTablePart {
	return b.parts
}

// Add appends a cross-reference table part to the builder.
func (b *CrossReferenceTableBuilder) Add(part *CrossReferenceTablePart) error {
	if part == nil {
		return fmt.Errorf("part must not be nil")
	}

	b.parts = append(b.parts, part)
	return nil
}

// Build merges all collected parts into a single CrossReferenceTable.
// firstCrossReferenceOffset is the byte position of the starting xref entry from StartXref.
// offsetCorrection accounts for any trailer length discrepancy between expected and actual positions.
// isLenientParsing controls strictness when parsing the trailer dictionary.
func (b *CrossReferenceTableBuilder) Build(
	firstCrossReferenceOffset int64,
	offsetCorrection int64,
	isLenientParsing bool,
	log logging.Log,
) (*CrossReferenceTable, error) {
	fileType := Table
	trailerDictionary := &tokens.DictionaryToken{}
	objectOffsets := make(map[core.IndirectReference]int64)

	xrefPartToBytePositionOrder := make([]int64, 0)

	currentPart := findPart(b.parts, firstCrossReferenceOffset)

	if currentPart == nil {
		log.Warn(fmt.Sprintf("Did not find an XRef object at the specified startxref position %d", firstCrossReferenceOffset))

		offsets := make([]int64, len(b.parts))
		for i, p := range b.parts {
			offsets[i] = p.Offset()
		}
		sort.Slice(offsets, func(i, j int) bool { return offsets[i] < offsets[j] })
		xrefPartToBytePositionOrder = append(xrefPartToBytePositionOrder, offsets...)
	} else {
		fileType = currentPart.Type()

		xrefPartToBytePositionOrder = append(xrefPartToBytePositionOrder, firstCrossReferenceOffset)

		for currentPart != nil && currentPart.Dictionary() != nil {
			activePart := currentPart
			for _, dependent := range b.parts {
				if tied := dependent.TiedToXrefAtOffset(); tied != nil && *tied == activePart.Offset() {
					xrefPartToBytePositionOrder = append(xrefPartToBytePositionOrder, dependent.Offset())
				}
			}

			prevBytePos := currentPart.GetPreviousOffset()
			if prevBytePos == -1 {
				break
			}

			currentPart = findPart(b.parts, prevBytePos)
			if currentPart == nil {
				corrected := prevBytePos + offsetCorrection
				currentPart = findPart(b.parts, corrected)
			}

			if currentPart == nil {
				log.Warn(fmt.Sprintf("Did not found XRef object pointed to by 'Prev' key at position %d", prevBytePos))
				break
			}

			xrefPartToBytePositionOrder = append(xrefPartToBytePositionOrder, prevBytePos)

			if len(xrefPartToBytePositionOrder) >= len(b.parts) {
				break
			}
		}

		for i, j := 0, len(xrefPartToBytePositionOrder)-1; i < j; i, j = i+1, j-1 {
			xrefPartToBytePositionOrder[i], xrefPartToBytePositionOrder[j] = xrefPartToBytePositionOrder[j], xrefPartToBytePositionOrder[i]
		}
	}

	for _, bPos := range xrefPartToBytePositionOrder {
		currentObject := findPart(b.parts, bPos)
		if currentObject == nil {
			corrected := bPos + offsetCorrection
			currentObject = findPart(b.parts, corrected)
		}

		if currentObject != nil && currentObject.Dictionary() != nil {
			for key, value := range currentObject.Dictionary().Data() {
				keyLower := toLower(key)
				isSizeKey := keyLower == "size"
				hasSizeEntry := trailerDictionary.ContainsKey(tokens.Size)

				if !isSizeKey || !hasSizeEntry {
					nameToken := tokens.MustCreate(key)
					trailerDictionary = trailerDictionary.With(nameToken, value)
				}
			}
		}

		if currentObject != nil {
			for ref, offset := range currentObject.ObjectOffsets() {
				objectOffsets[ref] = offset
			}
		}
	}

	offsetsList := make([]CrossReferenceOffset, len(b.parts))
	for i, p := range b.parts {
		prev := p.GetPreviousOffset()
		var prevPtr *int64
		if prev >= 0 {
			prevPtr = &prev
		}
		offsetsList[i] = NewCrossReferenceOffset(p.Offset(), prevPtr)
	}

	trailer, err := NewTrailerDictionary(trailerDictionary, isLenientParsing)
	if err != nil {
		return nil, fmt.Errorf("failed to build trailer dictionary: %w", err)
	}

	return NewCrossReferenceTable(fileType, objectOffsets, trailer, offsetsList)
}

func findPart(parts []*CrossReferenceTablePart, offset int64) *CrossReferenceTablePart {
	for _, p := range parts {
		if p.Offset() == offset {
			return p
		}
	}
	return nil
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + ('a' - 'A')
		} else {
			result[i] = c
		}
	}
	return string(result)
}
