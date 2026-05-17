package cmap

import (
	"encoding/binary"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/pdf_fonts"
)

// CharacterMapBuilder is a mutable builder used when parsing and generating a CMap.
type CharacterMapBuilder struct {
	characterIdentifierSystemInfo pdffonts.CharacterIdentifierSystemInfo
	systemInfoBuilder             *CharacterIdentifierSystemInfoBuilder
	wMode                         int
	name                          string
	version                       *string
	typ                           int
	codespaceRanges               []*CodespaceRange
	cidCharacterMappings          []CidCharacterMapping
	cidRanges                     []*CidRange
	baseFontCharacterMap          map[int]string
}

// NewCharacterMapBuilder creates a new CharacterMapBuilder instance.
func NewCharacterMapBuilder() *CharacterMapBuilder {
	return &CharacterMapBuilder{
		systemInfoBuilder:    NewCharacterIdentifierSystemInfoBuilder(),
		wMode:                0,
		typ:                  -1,
		baseFontCharacterMap: make(map[int]string),
	}
}

// CharacterIdentifierSystemInfo returns the CID system info.
func (b *CharacterMapBuilder) CharacterIdentifierSystemInfo() pdffonts.CharacterIdentifierSystemInfo {
	return b.characterIdentifierSystemInfo
}

// SetCharacterIdentifierSystemInfo sets the CID system info directly.
func (b *CharacterMapBuilder) SetCharacterIdentifierSystemInfo(info pdffonts.CharacterIdentifierSystemInfo) {
	b.characterIdentifierSystemInfo = info
}

// SystemInfoBuilder returns the builder for constructing CID system info incrementally.
func (b *CharacterMapBuilder) SystemInfoBuilder() *CharacterIdentifierSystemInfoBuilder {
	return b.systemInfoBuilder
}

// WMode returns the writing mode: 0 for horizontal, 1 for vertical.
func (b *CharacterMapBuilder) WMode() int {
	return b.wMode
}

// SetWMode sets the writing mode.
func (b *CharacterMapBuilder) SetWMode(mode int) {
	b.wMode = mode
}

// Name returns the PostScript name of the CMap.
func (b *CharacterMapBuilder) Name() string {
	return b.name
}

// SetName sets the PostScript name of the CMap.
func (b *CharacterMapBuilder) SetName(name string) {
	b.name = name
}

// Version returns the version string, or nil if not set.
func (b *CharacterMapBuilder) Version() *string {
	return b.version
}

// SetVersion sets the version string.
func (b *CharacterMapBuilder) SetVersion(version *string) {
	b.version = version
}

// Type returns the CMap file type number.
func (b *CharacterMapBuilder) Type() int {
	return b.typ
}

// SetType sets the CMap file type number.
func (b *CharacterMapBuilder) SetType(typ int) {
	b.typ = typ
}

// CodespaceRanges returns the codespace ranges.
func (b *CharacterMapBuilder) CodespaceRanges() []*CodespaceRange {
	return b.codespaceRanges
}

// SetCodespaceRanges sets the codespace ranges.
func (b *CharacterMapBuilder) SetCodespaceRanges(ranges []*CodespaceRange) {
	b.codespaceRanges = ranges
}

// CidCharacterMappings returns the CID character mappings.
func (b *CharacterMapBuilder) CidCharacterMappings() []CidCharacterMapping {
	return b.cidCharacterMappings
}

// SetCidCharacterMappings sets the CID character mappings.
func (b *CharacterMapBuilder) SetCidCharacterMappings(mappings []CidCharacterMapping) {
	b.cidCharacterMappings = mappings
}

// CidRanges returns the CID ranges.
func (b *CharacterMapBuilder) CidRanges() []*CidRange {
	return b.cidRanges
}

// BaseFontCharacterMap returns the base font character map.
func (b *CharacterMapBuilder) BaseFontCharacterMap() map[int]string {
	return b.baseFontCharacterMap
}

// AddBaseFontCharacter adds a base font character mapping from byte slices.
func (b *CharacterMapBuilder) AddBaseFontCharacterBytes(bytes []byte, value []byte) {
	str := createStringFromBytes(value)
	b.AddBaseFontCharacter(bytes, str)
}

// AddBaseFontCharacter adds a base font character mapping from bytes and string.
func (b *CharacterMapBuilder) AddBaseFontCharacter(bytes []byte, value string) {
	code := getCodeFromArray(bytes)
	b.baseFontCharacterMap[code] = value
}

// Build constructs the final CMap. Returns an error if construction fails.
func (b *CharacterMapBuilder) Build() (*CMap, error) {
	codespaceRanges := b.codespaceRanges
	if codespaceRanges == nil {
		codespaceRanges = make([]*CodespaceRange, 0)
	}

	cidRanges := b.cidRanges
	if cidRanges == nil {
		cidRanges = make([]*CidRange, 0)
	}

	cidCharacterMappings := b.cidCharacterMappings
	if cidCharacterMappings == nil {
		cidCharacterMappings = make([]CidCharacterMapping, 0)
	}

	return NewCMap(
		b.getCidSystemInfo(),
		b.typ,
		b.wMode,
		b.name,
		b.version,
		b.baseFontCharacterMap,
		codespaceRanges,
		cidRanges,
		cidCharacterMappings,
	)
}

func (b *CharacterMapBuilder) getCidSystemInfo() pdffonts.CharacterIdentifierSystemInfo {
	if b.characterIdentifierSystemInfo.Registry != "" {
		return b.characterIdentifierSystemInfo
	}

	if b.systemInfoBuilder.HasOrdering() && b.systemInfoBuilder.HasRegistry() && b.systemInfoBuilder.HasSupplement() {
		return pdffonts.NewCharacterIdentifierSystemInfo(
			b.systemInfoBuilder.Registry(),
			b.systemInfoBuilder.Ordering(),
			b.systemInfoBuilder.Supplement(),
		)
	}

	return b.characterIdentifierSystemInfo
}

// UseCMap merges another CMap's data into this builder.
func (b *CharacterMapBuilder) UseCMap(other *CMap) {
	b.codespaceRanges = combineCodespaceRanges(b.codespaceRanges, other.CodespaceRanges())

	otherMappings := make([]CidCharacterMapping, 0, len(other.CidCharacterMappings()))
	for _, m := range other.CidCharacterMappings() {
		otherMappings = append(otherMappings, m)
	}
	b.cidCharacterMappings = combineCidMappings(b.cidCharacterMappings, otherMappings)

	b.cidRanges = append(b.cidRanges, other.CidRanges()...)

	if bfcm := other.BaseFontCharacterMap(); bfcm != nil {
		for k, v := range bfcm {
			b.baseFontCharacterMap[k] = v
		}
	}
}

// AddCidRange appends a CID range to the builder.
func (b *CharacterMapBuilder) AddCidRange(r *CidRange) {
	b.cidRanges = append(b.cidRanges, r)
}

func combineCodespaceRanges(a, b []*CodespaceRange) []*CodespaceRange {
	if a == nil && b == nil {
		return make([]*CodespaceRange, 0)
	}
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	result := make([]*CodespaceRange, 0, len(a)+len(b))
	result = append(result, a...)
	result = append(result, b...)
	return result
}

func combineCidMappings(a []CidCharacterMapping, b []CidCharacterMapping) []CidCharacterMapping {
	if a == nil && b == nil {
		return make([]CidCharacterMapping, 0)
	}
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	result := make([]CidCharacterMapping, 0, len(a)+len(b))
	result = append(result, a...)
	result = append(result, b...)
	return result
}

func getCodeFromArray(data []byte) int {
	code := 0
	for _, v := range data {
		code <<= 8
		code |= int(v)
	}
	return code
}

func createStringFromBytes(bytes []byte) string {
	if len(bytes) == 1 {
		return core.BytesAsLatin1String(bytes)
	}

	if len(bytes)%2 != 0 {
		bytes = append(bytes, 0)
	}

	codeUnits := make([]uint16, len(bytes)/2)
	for i := 0; i < len(bytes); i += 2 {
		codeUnits[i/2] = binary.BigEndian.Uint16(bytes[i : i+2])
	}

	return string(runeSliceFromUTF16(codeUnits))
}

func runeSliceFromUTF16(codeUnits []uint16) []rune {
	var result []rune
	i := 0
	for i < len(codeUnits) {
		cu := codeUnits[i]
		if cu >= 0xD800 && cu <= 0xDBFF && i+1 < len(codeUnits) {
			next := codeUnits[i+1]
			if next >= 0xDC00 && next <= 0xDFFF {
				high := int(cu-0xD800) & 0x3FF
				low := int(next - 0xDC00) & 0x3FF
				runeVal := rune(0x10000 + (high << 10) + low)
				result = append(result, runeVal)
				i += 2
				continue
			}
		}
		if cu >= 0xD800 && cu <= 0xDFFF {
			result = append(result, '\uFFFD')
		} else {
			result = append(result, rune(cu))
		}
		i++
	}
	return result
}

