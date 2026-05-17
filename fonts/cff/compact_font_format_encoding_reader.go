// Package cff provides types for parsing and working with Compact Font Format (CFF) data in PDF files.
package cff

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/cff_charset"
	"github.com/uglytoad/pdfpig/go/fonts/encodings"
)

// CompactFontFormatFormat1Encoding represents a CFF format 1 encoding table.
// Format 1 uses ranges to map character codes compactly.
type CompactFontFormatFormat1Encoding struct {
	*CFFBuiltInEncoding
	NumberOfRanges int
	Entries        []CFFFormat0CodeEntry
	Supplements    []CFFBuiltInEncodingSupplement
}

// NewCompactFontFormatFormat1Encoding creates a new format 1 encoding.
func NewCompactFontFormatFormat1Encoding(numberOfRanges int, entries []CFFFormat0CodeEntry, supplements []CFFBuiltInEncodingSupplement) *CompactFontFormatFormat1Encoding {
	if supplements == nil {
		supplements = make([]CFFBuiltInEncodingSupplement, 0)
	}
	return &CompactFontFormatFormat1Encoding{
		CFFBuiltInEncoding: newCFFBuiltInEncodingFromEntries(entries, supplements),
		NumberOfRanges:     numberOfRanges,
		Entries:            entries,
		Supplements:        supplements,
	}
}

// newCFFBuiltInEncodingFromEntries builds a CFFBuiltInEncoding from code entries and supplements.
func newCFFBuiltInEncodingFromEntries(entries []CFFFormat0CodeEntry, supplements []CFFBuiltInEncodingSupplement) *CFFBuiltInEncoding {
	base := NewCompactFontFormatBaseEncoding()
	base.AddWithName(0, 0, encodings.NotDefined)

	e := &CFFBuiltInEncoding{
		CompactFontFormatBaseEncoding: base,
		Supplements:                   supplements,
	}

	for _, s := range supplements {
		e.AddWithName(s.Code, s.Sid, s.Name)
	}

	for _, entry := range entries {
		e.AddWithName(entry.Code, entry.Sid, entry.Str)
	}

	return e
}

// ReadEncoding reads a CFF encoding table from the given data using the provided charset and string index.
func ReadEncoding(data *CompactFontFormatData, charset cffcharset.CompactFontFormatCharset, stringIndex []string) (*CFFBuiltInEncoding, error) {
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	format := data.ReadCard8()
	baseFormat := format & 0x7f

	switch baseFormat {
	case 0:
		return readFormat0Encoding(data, charset, stringIndex, format)
	case 1:
		return readFormat1Encoding(data, charset, stringIndex, format)
	default:
		return nil, fonts.NewInvalidFontFormatException(fmt.Sprintf("the provided format %d for this Compact Font Format encoding was invalid", format))
	}
}

func readFormat0Encoding(data *CompactFontFormatData, charset cffcharset.CompactFontFormatCharset, stringIndex []string, format byte) (*CFFBuiltInEncoding, error) {
	numberOfCodes := int(data.ReadCard8())

	entries := make([]CFFFormat0CodeEntry, 0, numberOfCodes)
	for i := 1; i <= numberOfCodes; i++ {
		code := int(data.ReadCard8())
		sid := charset.GetStringIdByGlyphId(i)
		str := readString(sid, stringIndex)
		entries = append(entries, CFFFormat0CodeEntry{code, sid, str})
	}

	var supplements []CFFBuiltInEncodingSupplement
	if hasSupplement(format) {
		supplements = readSupplement(data, stringIndex)
	}

	enc := NewCompactFontFormatFormat0Encoding(entries, supplements)
	return enc.CFFBuiltInEncoding, nil
}

func readFormat1Encoding(data *CompactFontFormatData, charset cffcharset.CompactFontFormatCharset, stringIndex []string, format byte) (*CFFBuiltInEncoding, error) {
	numberOfRanges := int(data.ReadCard8())

	fromRanges := make([]CFFFormat0CodeEntry, 0)

	gid := 1
	for i := 0; i < numberOfRanges; i++ {
		rangeFirst := int(data.ReadCard8())
		rangeLeft := int(data.ReadCard8())
		for j := 0; j < 1+rangeLeft; j++ {
			sid := charset.GetStringIdByGlyphId(gid)
			code := rangeFirst + j
			str := readString(sid, stringIndex)
			fromRanges = append(fromRanges, CFFFormat0CodeEntry{code, sid, str})
			gid++
		}
	}

	var supplements []CFFBuiltInEncodingSupplement
	if hasSupplement(format) {
		supplements = readSupplement(data, stringIndex)
	}

	enc := NewCompactFontFormatFormat1Encoding(numberOfRanges, fromRanges, supplements)
	return enc.CFFBuiltInEncoding, nil
}

func readSupplement(data *CompactFontFormatData, stringIndex []string) []CFFBuiltInEncodingSupplement {
	numberOfSupplements := int(data.ReadCard8())
	supplements := make([]CFFBuiltInEncodingSupplement, numberOfSupplements)

	for i := 0; i < numberOfSupplements; i++ {
		code := int(data.ReadCard8())
		sid := data.ReadSid()
		name := readString(sid, stringIndex)
		supplements[i] = NewCFFBuiltInEncodingSupplement(code, sid, name)
	}

	return supplements
}

func readString(index int, stringIndex []string) string {
	if index >= 0 && index <= 390 {
		return GetName(index)
	}
	if index-391 < len(stringIndex) {
		return stringIndex[index-391]
	}

	return "SID" + fmt.Sprintf("%d", index)
}

func hasSupplement(format byte) bool {
	return (format & 0x80) != 0
}
