package core

import (
	"errors"
	"io"
)

var _ InputBytes = (*StreamInputBytes)(nil)

// StreamInputBytes provides input bytes backed by an io.ReadSeeker.
type StreamInputBytes struct {
	reader      io.ReadSeeker
	length      int64
	shouldClose bool
	offset      int64 // -1 means before first byte; matches MemoryInputBytes semantics
	peekByte    byte
	hasPeek     bool
	isAtEnd     bool
	currentByte byte
	pendingByte bool
}

// NewStreamInputBytes creates a new StreamInputBytes from the given reader.
// The shouldClose parameter controls whether Close() will also close the underlying reader.
func NewStreamInputBytes(reader io.ReadSeeker, shouldClose bool) (*StreamInputBytes, error) {
	if reader == nil {
		return nil, errors.New("reader must not be nil")
	}

	length, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}

	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &StreamInputBytes{
		reader:      reader,
		length:      length,
		shouldClose: shouldClose,
		offset:      -1, // Before first byte, matches MemoryInputBytes semantics
	}, nil
}

// CurrentOffset returns the current byte offset in the stream.
// Matches MemoryInputBytes semantics: returns offset+1 (0-based position of current byte).
func (s *StreamInputBytes) CurrentOffset() int64 {
	return s.offset + 1
}

// MoveNext advances to the next byte. Returns true if successful, false if at end.
func (s *StreamInputBytes) MoveNext() bool {
	if s.pendingByte {
		s.pendingByte = false
		s.offset++
		var buf [1]byte
		n, err := s.reader.Read(buf[:])
		if n == 0 || err != nil {
			s.isAtEnd = true
			s.currentByte = 0
			return false
		}
		s.currentByte = buf[0]
		return true
	}

	var buf [1]byte
	n, err := s.reader.Read(buf[:])
	if n == 0 || err != nil {
		s.isAtEnd = true
		s.currentByte = 0
		return false
	}
	s.offset++
	s.currentByte = buf[0]
	return true
}

// CurrentByte returns the current byte at the cursor position.
func (s *StreamInputBytes) CurrentByte() byte {
	return s.currentByte
}

// Length returns the total length of the underlying data in bytes.
func (s *StreamInputBytes) Length() int64 {
	return s.length
}

// Peek returns the next byte without advancing the cursor.
// The second return value is false if at end of stream.
func (s *StreamInputBytes) Peek() (byte, bool) {
	if s.pendingByte {
		s.pendingByte = false
	}
	if !s.hasPeek {
		var buf [1]byte
		n, err := s.reader.Read(buf[:])
		if n == 0 || err != nil {
			return 0, false
		}
		s.peekByte = buf[0]
		s.hasPeek = true
	}

	return s.peekByte, true
}

// IsAtEnd reports whether the reader has reached the end of the data.
func (s *StreamInputBytes) IsAtEnd() bool {
	return s.isAtEnd
}

// Seek moves the cursor to the given position with the specified whence.
// Position N means "at byte index N" (0-based). After seeking:
// - CurrentOffset() returns N
// - CurrentByte() returns the byte at index N-1 (for N > 0), or 0 (for N == 0)
// This matches MemoryInputBytes semantics where offset is the index of current byte.
func (s *StreamInputBytes) Seek(offset int64, whence int) (int64, error) {
	s.isAtEnd = false
	s.hasPeek = false
	s.pendingByte = false

	pos, err := s.reader.Seek(offset, whence)
	if err != nil {
		return 0, err
	}

	if pos == 0 {
		s.offset = -1 // Before first byte; CurrentOffset() will return 0
		s.currentByte = 0
	} else if pos <= s.length {
		// Position to read byte at index (pos-1), then set offset=pos-1 so
		// MoveNext will advance to pos. CurrentOffset() returns offset+1 = pos.
		var buf [1]byte
		_, seekErr := s.reader.Seek(pos-1, io.SeekStart)
		if seekErr != nil {
			return pos, seekErr
		}
		n, readErr := s.reader.Read(buf[:])
		if n > 0 {
			s.currentByte = buf[0]
			s.pendingByte = true
		}
		s.offset = pos - 1 // CurrentOffset() returns pos; MoveNext will set it to pos
		if readErr != nil {
			return pos, readErr
		}
	} else {
		// Past end of file
		s.offset = pos - 1
		s.currentByte = 0
		s.isAtEnd = true
	}

	return pos, nil
}

// Read fills the buffer with bytes starting from the current position.
// Returns the number of bytes read and any error encountered.
func (s *StreamInputBytes) Read(buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, nil
	}

	if s.hasPeek {
		buffer[0] = s.peekByte
		s.hasPeek = false

		n, err := s.Read(buffer[1:])
		s.offset++
		return n + 1, err
	}

	n, err := s.reader.Read(buffer)
	if n > 0 {
		s.currentByte = buffer[n-1]
		s.offset += int64(n)
	}

	s.isAtEnd = s.offset >= s.length-1 && s.length > 0

	return n, err
}

// Close releases resources held by the input bytes.
func (s *StreamInputBytes) Close() error {
	if s.shouldClose {
		if closer, ok := s.reader.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}
