package cff

import (
	"fmt"
	"io"

	"golang.org/x/text/encoding"
)

// CompactFontFormatData provides access to the raw bytes of a Compact Font Format
// file with utility methods for reading data types from it.
type CompactFontFormatData struct {
	dataBytes []byte
	position  int
}

// NewCompactFontFormatData creates a new CompactFontFormatData from the given byte slice.
func NewCompactFontFormatData(dataBytes []byte) *CompactFontFormatData {
	return &CompactFontFormatData{
		dataBytes: dataBytes,
		position:  -1,
	}
}

// Position returns the current position in the data.
func (d *CompactFontFormatData) Position() int {
	return d.position
}

// Length returns the length of the underlying data.
func (d *CompactFontFormatData) Length() int {
	return len(d.dataBytes)
}

// ReadString reads a string of the specified length using the given encoding.
func (d *CompactFontFormatData) ReadString(length int, enc encoding.Encoding) string {
	span := d.readSpan(length)
	result, _ := enc.NewDecoder().Bytes(span)
	return string(result)
}

// ReadCard8 reads an 8-bit unsigned integer.
func (d *CompactFontFormatData) ReadCard8() byte {
	b, _ := d.ReadByte()
	return b
}

// ReadCard16 reads a 16-bit big-endian unsigned integer.
func (d *CompactFontFormatData) ReadCard16() uint16 {
	b1, _ := d.ReadByte()
	b2, _ := d.ReadByte()
	return uint16(b1)<<8 | uint16(b2)
}

// ReadOffsize reads an offset size byte.
func (d *CompactFontFormatData) ReadOffsize() byte {
	b, _ := d.ReadByte()
	return b
}

// ReadOffset reads an offset value of the given byte size in big-endian order.
func (d *CompactFontFormatData) ReadOffset(offsetSize int) int {
	value := 0
	for i := 0; i < offsetSize; i++ {
		b, _ := d.ReadByte()
		value = (value << 8) | int(b)
	}
	return value
}

// ReadByte reads a single byte and advances the position.
func (d *CompactFontFormatData) ReadByte() (byte, error) {
	d.position++
	if d.position >= len(d.dataBytes) {
		return 0, io.EOF
	}
	return d.dataBytes[d.position], nil
}

// Peek returns the next byte without advancing the position.
func (d *CompactFontFormatData) Peek() byte {
	return d.dataBytes[d.position+1]
}

// CanRead reports whether there is more data to read.
func (d *CompactFontFormatData) CanRead() bool {
	return d.position < len(d.dataBytes)-1
}

// Seek moves the position to the given offset from the beginning.
func (d *CompactFontFormatData) Seek(offset int) {
	d.position = offset - 1
}

// ReadLong reads a 32-bit big-endian signed integer as two Card16 values.
func (d *CompactFontFormatData) ReadLong() int32 {
	high := d.ReadCard16()
	low := d.ReadCard16()
	return int32((uint32(high)<<16) | uint32(low))
}

// ReadSid reads a string identifier (SID) as two bytes.
func (d *CompactFontFormatData) ReadSid() int {
	b1, _ := d.ReadByte()
	b2, _ := d.ReadByte()
	return int(b1)<<8 | int(b2)
}

// ReadBytes reads a byte array of the given length.
func (d *CompactFontFormatData) ReadBytes(length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i], _ = d.ReadByte()
	}
	return result
}

// SnapshotPortion creates a new CompactFontFormatData from a portion of this data.
func (d *CompactFontFormatData) SnapshotPortion(startLocation, length int) (*CompactFontFormatData, error) {
	if length == 0 {
		return NewCompactFontFormatData([]byte{}), nil
	}

	if startLocation > len(d.dataBytes)-1 || startLocation+length > len(d.dataBytes) {
		return nil, fmt.Errorf(
			"attempted to create a snapshot of an invalid portion of the data: length=%d, start=%d, requestedLength=%d",
			len(d.dataBytes), startLocation, length,
		)
	}

	newData := make([]byte, length)
	copy(newData, d.dataBytes[startLocation:startLocation+length])

	return NewCompactFontFormatData(newData), nil
}

// readSpan returns a slice of count bytes starting at position+1 and advances
// the position. This is the internal equivalent of C# ReadSpan.
func (d *CompactFontFormatData) readSpan(count int) []byte {
	if d.position+count >= len(d.dataBytes) {
		panic(fmt.Sprintf("cannot read past end of data: attempted to read to %d when the underlying data is %d bytes long", d.position+count, len(d.dataBytes)))
	}

	result := d.dataBytes[d.position+1 : d.position+1+count]
	d.position += count
	return result
}
