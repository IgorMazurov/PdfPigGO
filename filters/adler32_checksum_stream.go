package filters

import (
	"io"
)

// Adler32ChecksumStream wraps an io.WriteCloser and computes an Adler-32
// checksum over all data written through it.
// This mirrors the C# internal sealed class Adler32ChecksumStream : Stream
// used in FlateFilter.Encode and PngBuilder for computing the zlib ADLER32
// checksum of raw (uncompressed) data while writing through a DeflateStream.
type Adler32ChecksumStream struct {
	underlying io.WriteCloser
	checksum   uint32
}

// NewAdler32ChecksumStream creates a new Adler32ChecksumStream wrapping the
// given WriteCloser. The underlying stream must not be nil.
func NewAdler32ChecksumStream(underlying io.WriteCloser) *Adler32ChecksumStream {
	if underlying == nil {
		panic("underlying stream must not be nil")
	}
	return &Adler32ChecksumStream{
		underlying: underlying,
		checksum:   1,
	}
}

// Write writes to the underlying stream and updates the Adler-32 checksum.
func (a *Adler32ChecksumStream) Write(p []byte) (int, error) {
	n, err := a.underlying.Write(p)
	if n > 0 {
		a.updateAdler(p[:n])
	}
	return n, err
}

// Close closes the underlying stream.
func (a *Adler32ChecksumStream) Close() error {
	return a.underlying.Close()
}

// Checksum returns the current Adler-32 checksum computed over all data
// that has passed through this stream.
func (a *Adler32ChecksumStream) Checksum() uint32 {
	return a.checksum
}

const adlerMod = 65521

func (a *Adler32ChecksumStream) updateAdler(data []byte) {
	sa := a.checksum & 0xFFFF
	sb := (a.checksum >> 16) & 0xFFFF

	for _, c := range data {
		sa = (sa + uint32(c)) % adlerMod
		sb = (sb + sa) % adlerMod
	}

	a.checksum = (sb << 16) | sa
}
