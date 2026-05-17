package png

import (
	"encoding/binary"
	"io"
)

// WriteBigEndianInt32 writes a 32-bit integer in big-endian byte order to w.
func WriteBigEndianInt32(w io.Writer, value int32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(value))
	_, err := w.Write(buf)
	return err
}

// TryReadHeaderBytes reads exactly 8 bytes from r into the returned slice.
// It returns false if fewer than 8 bytes could be read.
func TryReadHeaderBytes(r io.Reader) ([]byte, bool) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, false
	}
	return buf, true
}
