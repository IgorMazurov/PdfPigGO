package core

import (
	"encoding/binary"
	"io"
)

// WriteUInt writes a int64 value to the writer as a big-endian uint32.
func WriteUIntInt64(w io.Writer, value int64) (int, error) {
	return WriteUInt(w, uint32(value))
}

// WriteUInt writes a uint32 value to the writer in big-endian byte order.
func WriteUInt(w io.Writer, value uint32) (int, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, value)
	return w.Write(buf)
}

// WriteUShort writes a int32 value to the writer as a big-endian uint16.
func WriteUShortInt32(w io.Writer, value int32) (int, error) {
	return WriteUShort(w, uint16(value))
}

// WriteUShort writes a uint16 value to the writer in big-endian byte order.
func WriteUShort(w io.Writer, value uint16) (int, error) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, value)
	return w.Write(buf)
}

// WriteShort writes a uint16 value to the writer as a big-endian int16.
func WriteShortUint16(w io.Writer, value uint16) (int, error) {
	return WriteShort(w, int16(value))
}

// WriteShort writes an int16 value to the writer in big-endian byte order.
func WriteShort(w io.Writer, value int16) (int, error) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(value))
	return w.Write(buf)
}
