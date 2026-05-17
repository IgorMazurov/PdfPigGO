package core

import "fmt"

var _ InputBytes = (*MemoryInputBytes)(nil)

// MemoryInputBytes provides input bytes backed by a byte slice in memory.
type MemoryInputBytes struct {
	memory      []byte
	upperBound  int
	currentOffset int
	currentByte   byte
}

// NewMemoryInputBytes creates a new MemoryInputBytes from the given byte slice.
func NewMemoryInputBytes(data []byte) *MemoryInputBytes {
	copied := make([]byte, len(data))
	copy(copied, data)

	return &MemoryInputBytes{
		memory:        copied,
		upperBound:    len(copied) - 1,
		currentOffset: -1,
	}
}

// CurrentOffset returns the current offset (1-based position).
func (m *MemoryInputBytes) CurrentOffset() int64 {
	return int64(m.currentOffset + 1)
}

// MoveNext advances to the next byte. Returns true if successful, false if at end.
func (m *MemoryInputBytes) MoveNext() bool {
	if m.currentOffset == m.upperBound {
		return false
	}

	m.currentOffset++
	m.currentByte = m.memory[m.currentOffset]
	return true
}

// CurrentByte returns the byte at the current cursor position.
func (m *MemoryInputBytes) CurrentByte() byte {
	return m.currentByte
}

// Length returns the total number of bytes in the underlying data.
func (m *MemoryInputBytes) Length() int64 {
	return int64(len(m.memory))
}

// Peek returns the next byte without advancing the cursor.
// The second return value is false if at end.
func (m *MemoryInputBytes) Peek() (byte, bool) {
	if m.currentOffset == m.upperBound {
		return 0, false
	}

	return m.memory[m.currentOffset+1], true
}

// IsAtEnd reports whether the cursor is at the last byte.
func (m *MemoryInputBytes) IsAtEnd() bool {
	return m.currentOffset == m.upperBound
}

// Seek moves the cursor to the given position with the specified whence.
func (m *MemoryInputBytes) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0: // io.SeekStart
		m.currentOffset = int(offset) - 1
	case 1: // io.SeekCurrent
		m.currentOffset += int(offset)
	case 2: // io.SeekEnd
		m.currentOffset = len(m.memory) + int(offset) - 1
	default:
		return 0, fmt.Errorf("invalid whence: %d", whence)
	}
	if m.currentOffset < 0 {
		m.currentByte = 0
	} else if m.currentOffset >= len(m.memory) {
		m.currentOffset = len(m.memory) - 1
		m.currentByte = m.memory[m.currentOffset]
	} else {
		m.currentByte = m.memory[m.currentOffset]
	}
	return int64(m.currentOffset + 1), nil
}

// Read fills the buffer with bytes starting from the current position.
// Returns the number of bytes read and nil error.
func (m *MemoryInputBytes) Read(buffer []byte) (int, error) {
	if len(buffer) == 0 {
		return 0, nil
	}

	viableLength := len(m.memory) - m.currentOffset - 1
	readLength := viableLength
	if readLength > len(buffer) {
		readLength = len(buffer)
	}
	startFrom := m.currentOffset + 1

	copy(buffer, m.memory[startFrom:startFrom+readLength])

	if readLength > 0 {
		m.currentOffset += readLength
		m.currentByte = buffer[readLength-1]
	}

	return readLength, nil
}

// Close releases resources. No-op for in-memory data.
func (m *MemoryInputBytes) Close() error {
	return nil
}
