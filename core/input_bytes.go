package core

// InputBytes represents the input bytes for a PDF document.
type InputBytes interface {
	// CurrentOffset returns the current offset in bytes.
	CurrentOffset() int64

	// MoveNext moves to the next byte if available.
	// Returns true if there is another byte, false otherwise.
	MoveNext() bool

	// CurrentByte returns the current byte at the cursor position.
	CurrentByte() byte

	// Length returns the length of the data in bytes.
	Length() int64

	// Peek returns the next byte without advancing the cursor.
	// The second return value indicates whether a byte was available.
	Peek() (byte, bool)

	// IsAtEnd reports whether we are at the end of the available data.
	IsAtEnd() bool

	// Seek moves to a given position with the specified whence.
	// Returns the new offset from the start of the file.
	Seek(offset int64, whence int) (int64, error)

	// Read fills the buffer with bytes starting from the current position.
	// Returns the number of bytes successfully read.
	Read(buffer []byte) (int, error)

	// Close releases any resources held by the input bytes.
	Close() error
}
