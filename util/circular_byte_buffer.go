package util

// CircularByteBuffer provides a fixed-size circular buffer for bytes.
type CircularByteBuffer struct {
	buffer []byte
	start  int
	count  int
}

// NewCircularByteBuffer creates a new circular byte buffer with the given size.
func NewCircularByteBuffer(size int) *CircularByteBuffer {
	return &CircularByteBuffer{
		buffer: make([]byte, size),
	}
}

// Add appends a byte to the end of the buffer. If the buffer is full,
// the oldest byte (at the start) is overwritten and start advances.
func (c *CircularByteBuffer) Add(b byte) {
	insertionPosition := (c.start + c.count) % len(c.buffer)
	c.buffer[insertionPosition] = b
	if c.count < len(c.buffer) {
		c.count++
	} else {
		c.start = (c.start + 1) % len(c.buffer)
	}
}

// AddReverse adds a byte to the start of the buffer. If the buffer is full,
// the byte at the end is overwritten.
func (c *CircularByteBuffer) AddReverse(b byte) {
	c.start = (c.start - 1 + len(c.buffer)) % len(c.buffer)
	c.buffer[c.start] = b
	if c.count < len(c.buffer) {
		c.count++
	}
}

// EndsWith checks whether the logical contents of the buffer end with the given string.
func (c *CircularByteBuffer) EndsWith(s string) bool {
	if len(s) > c.count {
		return false
	}
	for i := 0; i < len(s); i++ {
		strIdx := s[i]
		inBuffer := c.count - (len(s) - i)
		idx := indexToBufferIndex(c, inBuffer)
		if c.buffer[idx] != strIdx {
			return false
		}
	}
	return true
}

// IsCurrentlyEqual checks whether the buffer currently equals the given string.
func (c *CircularByteBuffer) IsCurrentlyEqual(s string) bool {
	if len(s) > len(c.buffer) {
		return false
	}
	for i := 0; i < len(s); i++ {
		idx := indexToBufferIndex(c, i)
		if c.buffer[idx] != s[i] {
			return false
		}
	}
	return true
}

// AsBytes returns a contiguous copy of the logical bytes in order.
func (c *CircularByteBuffer) AsBytes() []byte {
	tmp := make([]byte, c.count)
	for i := 0; i < c.count; i++ {
		idx := indexToBufferIndex(c, i)
		tmp[i] = c.buffer[idx]
	}
	return tmp
}

// String returns an ASCII string representation of the buffer contents.
func (c *CircularByteBuffer) String() string {
	return string(c.AsBytes())
}

func indexToBufferIndex(c *CircularByteBuffer, i int) int {
	return (c.start + i) % len(c.buffer)
}
