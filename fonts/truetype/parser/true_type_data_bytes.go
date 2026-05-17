// Package truetypeparser provides types and helpers for TrueType font parsing.
// TrueTypeDataBytes and GetNameTable are defined here; full Parse logic is in true_type_font_parser.go (row 193).
package truetypeparser

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/uglytoad/pdfpig/go/fonts/truetype"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

// NameTableStub is a minimal name table representation used during font discovery.
type NameTableStub struct {
	fontName         string
	postscriptName   string
}

// FontName returns the font name from the name table.
func (n *NameTableStub) FontName() string {
	return n.fontName
}

// GetPostscriptName returns the PostScript name for the font.
func (n *NameTableStub) GetPostscriptName() string {
	return n.postscriptName
}

// TrueTypeDataBytes wraps raw TrueType font bytes for sequential reading of
// TrueType-specific data types. Ported from src/UglyToad.PdfPig.Fonts/TrueType/TrueTypeDataBytes.cs.
type TrueTypeDataBytes struct {
	data   []byte
	offset int64
}

// NewTrueTypeDataBytes creates a TrueTypeDataBytes from raw font bytes.
func NewTrueTypeDataBytes(bytes []byte) *TrueTypeDataBytes {
	return &TrueTypeDataBytes{data: bytes, offset: 0}
}

// Position returns the current position in the data.
func (d *TrueTypeDataBytes) Position() int64 {
	return d.offset
}

// Length returns the length of the data in bytes.
func (d *TrueTypeDataBytes) Length() int64 {
	return int64(len(d.data))
}

// Seek moves to the specified position with the given whence.
func (d *TrueTypeDataBytes) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0: // io.SeekStart
		d.offset = offset
	case 1: // io.SeekCurrent
		d.offset += offset
	case 2: // io.SeekEnd
		d.offset = int64(len(d.data)) + offset
	default:
		return 0, fmt.Errorf("invalid whence: %d", whence)
	}
	if d.offset < 0 {
		d.offset = 0
	}
	if d.offset > int64(len(d.data)) {
		d.offset = int64(len(d.data))
	}
	return d.offset, nil
}

// Read32Fixed reads a 16.16 fixed-point value and returns it as float32.
func (d *TrueTypeDataBytes) Read32Fixed() float32 {
	integral := float32(d.ReadSignedShort())
	fractional := float32(d.ReadUnsignedShort()) / 65536.0
	return integral + fractional
}

// ReadSignedShort reads a big-endian int16 from the current position.
func (d *TrueTypeDataBytes) ReadSignedShort() int16 {
	b0 := d.data[d.offset]
	b1 := d.data[d.offset+1]
	d.offset += 2
	return int16(uint16(b0)<<8 | uint16(b1))
}

// ReadUnsignedShort reads a big-endian uint16 from the current position.
func (d *TrueTypeDataBytes) ReadUnsignedShort() uint16 {
	if d.offset+1 >= int64(len(d.data)) {
		return 0
	}
	b0 := d.data[d.offset]
	b1 := d.data[d.offset+1]
	d.offset += 2
	return uint16(b0)<<8 | uint16(b1)
}

// ReadByte reads a single byte from the current position and advances the offset.
func (d *TrueTypeDataBytes) ReadByte() (byte, error) {
	if d.offset >= int64(len(d.data)) {
		return 0, io.EOF
	}
	b := d.data[d.offset]
	d.offset++
	return b, nil
}

// ReadTag reads a 4-character tag from the TrueType file at the current position.
func (d *TrueTypeDataBytes) ReadTag() string {
	if result, ok := d.TryReadString(4, encoding.Nop); ok {
		return result
	}
	return ""
}

// TryReadString reads bytesToRead bytes from the current position and decodes them
// using the given encoding. Returns (string, true) on success or ("", false) on failure.
func (d *TrueTypeDataBytes) TryReadString(bytesToRead int, enc encoding.Encoding) (string, bool) {
	if enc == nil || bytesToRead <= 0 || d.offset+int64(bytesToRead) > int64(len(d.data)) {
		return "", false
	}

	start := d.offset
	d.offset += int64(bytesToRead)

	if enc == encoding.Nop {
		return string(d.data[start:d.offset]), true
	}

	src := d.data[start:d.offset]

	// Fast path: ISO-8859-1 maps bytes 0x00-0xFF directly to runes U+0000-U+00FF.
	// Convert in-place without decoder allocation.
	if enc == iso88591Encoding {
		r := make([]rune, bytesToRead)
		for i, b := range src {
			r[i] = rune(b)
		}
		return string(r), true
	}

	// Use pre-created decoders for common TTF encodings to avoid allocation.
	var dec *encoding.Decoder
	switch enc {
	case utf16BEEncoding:
		dec = utf16BEDecoder
	default:
		dec = enc.NewDecoder()
	}

	// Transform directly on source bytes to avoid intermediate string allocation.
	dst := ttfDstPool.Get().([]byte)
	poolCap := cap(dst)
	if poolCap < bytesToRead*2 {
		dst = make([]byte, bytesToRead*2)
	} else {
		dst = dst[:bytesToRead*2]
	}

	n, _, err := dec.Transform(dst, src, true)
	if err != nil {
		ttfDstPool.Put(dst)
		return "", false
	}

	result := string(dst[:n])
	ttfDstPool.Put(dst)
	return result, true
}

// ReadUnsignedInt reads a big-endian uint32 from the current position.
func (d *TrueTypeDataBytes) ReadUnsignedInt() uint32 {
	if d.offset+3 >= int64(len(d.data)) {
		return 0
	}
	b0 := d.data[d.offset]
	b1 := d.data[d.offset+1]
	b2 := d.data[d.offset+2]
	b3 := d.data[d.offset+3]
	d.offset += 4
	return uint32(b0)<<24 | uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3)
}

// ReadSignedInt reads a big-endian int32 from the current position.
func (d *TrueTypeDataBytes) ReadSignedInt() int32 {
	return int32(d.ReadUnsignedInt())
}

// ReadLong reads a big-endian int64 as two consecutive 32-bit signed integers.
func (d *TrueTypeDataBytes) ReadLong() int64 {
	upper := int64(d.ReadSignedInt())
	lower := int64(d.ReadSignedInt())
	return (upper << 32) | (lower & 0xFFFFFFFF)
}

// ReadInternationalDate reads a TrueType date (seconds since 1904-01-01 UTC).
func (d *TrueTypeDataBytes) ReadInternationalDate() time.Time {
	secondsSince1904 := d.ReadLong()
	base := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	result := base.Add(time.Duration(secondsSince1904) * time.Second)
	if result.Before(base) && secondsSince1904 > 0 {
		return base
	}
	return result
}

// ReadByteArray reads length bytes from the current position and advances the offset.
func (d *TrueTypeDataBytes) ReadByteArray(length int) []byte {
	result := make([]byte, length)
	copy(result, d.data[d.offset:d.offset+int64(length)])
	d.offset += int64(length)
	return result
}

// ReadSignedByte reads a signed byte from the current position and advances the offset.
func (d *TrueTypeDataBytes) ReadSignedByte() int8 {
	b := d.data[d.offset]
	d.offset++
	return int8(b)
}

// ReadUnsignedShortArray reads n big-endian uint16 values from the current position.
func (d *TrueTypeDataBytes) ReadUnsignedShortArray(n int) []uint16 {
	result := make([]uint16, n)
	for i := 0; i < n; i++ {
		result[i] = d.ReadUnsignedShort()
	}
	return result
}

// ReadSignedShortArray reads n big-endian int16 values from the current position.
func (d *TrueTypeDataBytes) ReadSignedShortArray(n int) []int16 {
	result := make([]int16, n)
	for i := 0; i < n; i++ {
		result[i] = d.ReadSignedShort()
	}
	return result
}

// ReadUnsignedIntArray reads n big-endian uint32 values from the current position.
func (d *TrueTypeDataBytes) ReadUnsignedIntArray(n int) []uint32 {
	result := make([]uint32, n)
	for i := 0; i < n; i++ {
		result[i] = d.ReadUnsignedInt()
	}
	return result
}

// String returns a debug string showing current position and total length.
func (d *TrueTypeDataBytes) String() string {
	return fmt.Sprintf("@: %d of %d bytes", d.offset, len(d.data))
}

var iso88591Encoding encoding.Encoding = charmap.ISO8859_1

var asciiEncoding encoding.Encoding = encoding.Nop

var utf16BEEncoding encoding.Encoding = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)

// Pre-created decoders for common TTF encodings (stateless, safe to reuse).
var (
	iso88591Decoder = charmap.ISO8859_1.NewDecoder()
	utf16BEDecoder  = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
)

// ttfDstPool reuses byte slices for decoder output to avoid per-call allocation.
var ttfDstPool = sync.Pool{
	New: func() any { return make([]byte, 0, 256) },
}

// GetNameTable extracts the name table from TrueType font bytes by scanning the
// table directory for the "name" entry and parsing it via NameTableParser.
func GetNameTable(data *TrueTypeDataBytes) *NameTableStub {
	if data == nil || len(data.data) == 0 {
		return nil
	}

	data.Read32Fixed() // version, unused
	numTables := int(data.ReadUnsignedShort())
	data.ReadUnsignedShort() // searchRange
	data.ReadUnsignedShort() // entrySelector
	data.ReadUnsignedShort() // rangeShift

	var nameHeader truetype.TrueTypeHeaderTable
	for i := 0; i < numTables; i++ {
		tag := data.ReadTag()
		checksum := data.ReadUnsignedInt()
		offset := data.ReadUnsignedInt()
		length := data.ReadUnsignedInt()

		if tag == truetype.Name {
			nameHeader = truetype.TrueTypeHeaderTable{
				Tag:      tag,
				CheckSum: checksum,
				Offset:   offset,
				Length:   length,
			}
			break
		}
	}

	if nameHeader.Tag == "" {
		return nil
	}

	nameTable, err := ParseName(nameHeader, data, &TableRegisterBuilder{})
	if err != nil || nameTable == nil {
		return nil
	}

	psName := nameTable.GetPostscriptName()
	fontName := nameTable.FontName()

	return &NameTableStub{
		fontName:       fontName,
		postscriptName: psName,
	}
}

// Parse fully parses TrueType font data into a TrueTypeFont.
// This is a convenience wrapper that discards errors for backward compatibility.
// Use ParseFull() for proper error handling.
func Parse(data *TrueTypeDataBytes) *TrueTypeFont {
	font, _ := ParseFull(data)
	return font
}
