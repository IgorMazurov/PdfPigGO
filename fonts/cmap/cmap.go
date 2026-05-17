package cmap

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/pdf_fonts"
)

// CMap maps character codes to character identifiers (CIDs).
// The set of characters which a CMap refers to is the "character set" (charset).
type CMap struct {
	info                  pdffonts.CharacterIdentifierSystemInfo
	typ                   int
	name                  string
	version               *string
	baseFontCharacterMap  map[int]string
	codespaceRanges       []*CodespaceRange
	cidRanges             []*CidRange
	cidCharacterMappings  map[int]CidCharacterMapping
	writingMode           WritingMode
	hasEmptyCodespace     bool
	minCodeLength         int
	maxCodeLength         int
}

// NewCMap creates a new CMap instance. Returns an error if required arguments are nil/empty.
func NewCMap(
	info pdffonts.CharacterIdentifierSystemInfo,
	typ int,
	wMode int,
	name string,
	version *string,
	baseFontCharacterMap map[int]string,
	codespaceRanges []*CodespaceRange,
	cidRanges []*CidRange,
	cidCharacterMappings []CidCharacterMapping,
) (*CMap, error) {
	if cidCharacterMappings == nil {
		return nil, fmt.Errorf("cidCharacterMappings must not be nil")
	}

	if baseFontCharacterMap == nil {
		return nil, fmt.Errorf("baseFontCharacterMap must not be nil")
	}

	if codespaceRanges == nil {
		return nil, fmt.Errorf("codespaceRanges must not be nil")
	}

	if cidRanges == nil {
		return nil, fmt.Errorf("cidRanges must not be nil")
	}

	minCodeLength := 4
	var maxCodeLength int
	hasEmptyCodespace := false

	if len(codespaceRanges) > 0 {
		maxCL := codespaceRanges[0].CodeLength
		minCL := codespaceRanges[0].CodeLength
		for _, r := range codespaceRanges {
			if r.CodeLength > maxCL {
				maxCL = r.CodeLength
			}
			if r.CodeLength < minCL {
				minCL = r.CodeLength
			}
			minCodeLength = minCL
			maxCodeLength = maxCL
		}
	} else {
		hasEmptyCodespace = true
	}

	characterMappings := make(map[int]CidCharacterMapping)
	for _, m := range cidCharacterMappings {
		characterMappings[m.SourceCharacterCode] = m
	}

	return &CMap{
		info:                 info,
		typ:                  typ,
		name:                 name,
		version:              version,
		baseFontCharacterMap: baseFontCharacterMap,
		codespaceRanges:      codespaceRanges,
		cidRanges:            cidRanges,
		cidCharacterMappings: characterMappings,
		writingMode:          WritingMode(wMode),
		hasEmptyCodespace:    hasEmptyCodespace,
		minCodeLength:        minCodeLength,
		maxCodeLength:        maxCodeLength,
	}, nil
}

// Info returns the CID system info.
func (c *CMap) Info() pdffonts.CharacterIdentifierSystemInfo {
	return c.info
}

// Type returns the type of the internal organization of the CMap file.
func (c *CMap) Type() int {
	return c.typ
}

// Name returns the name of the CMap file.
func (c *CMap) Name() string {
	return c.name
}

// Version returns the version number of the CIDFont file, or nil if not set.
func (c *CMap) Version() *string {
	return c.version
}

// BaseFontCharacterMap returns the base font character map.
func (c *CMap) BaseFontCharacterMap() map[int]string {
	return c.baseFontCharacterMap
}

// CodespaceRanges returns the set of valid input character codes.
func (c *CMap) CodespaceRanges() []*CodespaceRange {
	return c.codespaceRanges
}

// CidRanges returns the CID range associations.
func (c *CMap) CidRanges() []*CidRange {
	return c.cidRanges
}

// CidCharacterMappings returns the single character code to CID overrides.
func (c *CMap) CidCharacterMappings() map[int]CidCharacterMapping {
	return c.cidCharacterMappings
}

// WritingMode returns whether the font writes horizontally or vertically.
func (c *CMap) WritingMode() WritingMode {
	return c.writingMode
}

// HasCidMappings reports whether the CMap has any CID mappings.
func (c *CMap) HasCidMappings() bool {
	return len(c.cidCharacterMappings) > 0 || len(c.cidRanges) > 0
}

// HasUnicodeMappings reports whether the CMap has Unicode mappings.
func (c *CMap) HasUnicodeMappings() bool {
	return len(c.baseFontCharacterMap) > 0
}

// TryConvertToUnicode returns the Unicode string for the given character code.
// Returns true if an entry exists, false otherwise.
func (c *CMap) TryConvertToUnicode(code int) (string, bool) {
	result, found := c.baseFontCharacterMap[code]
	return result, found
}

// ConvertToCid converts a character code to a CID.
// Returns 0 if no mapping is found.
func (c *CMap) ConvertToCid(code int) int {
	if mapping, ok := c.cidCharacterMappings[code]; ok {
		return mapping.DestinationCid
	}

	for _, r := range c.cidRanges {
		if cid, ok := r.TryMap(code); ok {
			return cid
		}
	}

	return 0
}

// String returns the CMap name.
func (c *CMap) String() string {
	return c.name
}

// ReadCode reads a character code from the input bytes according to the codespace ranges.
// Returns an error if the CMap is invalid and lenient parsing is disabled.
func (c *CMap) ReadCode(bytes core.InputBytes, useLenientParsing bool) (int, error) {
	myPosition := bytes.CurrentOffset()

	if c.hasEmptyCodespace {
		data := make([]byte, c.minCodeLength)
		n, err := bytes.Read(data)
		if err != nil {
			return 0, err
		}
		return ToInt(data[:n]), nil
	}

	result := make([]byte, c.maxCodeLength)
	result[0] = bytes.CurrentByte()

	for i := 1; i < c.minCodeLength; i++ {
		if bytes.IsAtEnd() {
			break
		}
		b, err := readByte(bytes, useLenientParsing)
		if err != nil {
			return 0, err
		}
		result[i] = b
	}

	for i := c.minCodeLength - 1; i < c.maxCodeLength; i++ {
		byteCount := i + 1
		for _, r := range c.codespaceRanges {
			if r.IsFullMatch(result, byteCount) {
				return ToInt(result[:byteCount]), nil
			}
		}
		if byteCount < c.maxCodeLength {
			b, err := readByte(bytes, useLenientParsing)
			if err != nil {
				return 0, err
			}
			result[byteCount] = b
		}
	}

	if useLenientParsing {
		bytes.Seek(myPosition, 0)
		for i := 0; i < c.minCodeLength; i++ {
			b, err := readByte(bytes, useLenientParsing)
			if err != nil {
				return 0, err
			}
			result[i] = b
		}
		return ToInt(result[:c.minCodeLength]), nil
	}

	hexBytes := make([]string, len(result))
	for i, b := range result {
		hexBytes[i] = fmt.Sprintf("%02X", b)
	}
	return 0, core.NewPdfDocumentFormatException(
		fmt.Sprintf("CMap is invalid, min code length was %d, max was %d. Bytes: %s.",
			c.minCodeLength, c.maxCodeLength, joinHex(hexBytes)))
}

func readByte(bytes core.InputBytes, useLenientParsing bool) (byte, error) {
	if !bytes.MoveNext() {
		if useLenientParsing {
			return 0, nil
		}
		return 0, fmt.Errorf("read byte called on input bytes which was at end of byte set. Current offset: %d", bytes.CurrentOffset())
	}
	return bytes.CurrentByte(), nil
}

func joinHex(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += "-"
		}
		result += p
	}
	return result
}
