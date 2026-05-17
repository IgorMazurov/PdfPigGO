package iostream

import (
	"io"
)

// StreamWrapper wraps an io.ReadWriteSeeker and delegates all operations
// to the underlying stream. This mirrors the C# internal class StreamWrapper : Stream
// used as a base class for filter decoder streams like CcittFaxDecoderStream.
type StreamWrapper struct {
	stream io.ReadWriteSeeker
}

// NewStreamWrapper creates a new StreamWrapper around the given ReadWriteSeeker.
func NewStreamWrapper(stream io.ReadWriteSeeker) *StreamWrapper {
	return &StreamWrapper{stream: stream}
}

// Stream returns the underlying io.ReadWriteSeeker.
func (s *StreamWrapper) Stream() io.ReadWriteSeeker {
	return s.stream
}

// Read reads from the underlying stream.
func (s *StreamWrapper) Read(p []byte) (int, error) {
	return s.stream.Read(p)
}

// Write writes to the underlying stream.
func (s *StreamWrapper) Write(p []byte) (int, error) {
	return s.stream.Write(p)
}

// Seek sets the offset for the next Read or Write relative to origin.
func (s *StreamWrapper) Seek(offset int64, whence int) (int64, error) {
	return s.stream.Seek(offset, whence)
}

// Flush flushes any buffered data to the underlying stream if it supports flushing.
func (s *StreamWrapper) Flush() error {
	if flusher, ok := s.stream.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}

// SetLength sets the length of the underlying stream if it supports setting length.
func (s *StreamWrapper) SetLength(value int64) error {
	if setter, ok := s.stream.(interface{ SetLength(int64) error }); ok {
		return setter.SetLength(value)
	}
	return nil
}

// CanRead reports whether the underlying stream supports reading.
func (s *StreamWrapper) CanRead() bool {
	_, ok := s.stream.(io.Reader)
	return ok
}

// CanSeek reports whether the underlying stream supports seeking.
func (s *StreamWrapper) CanSeek() bool {
	_, ok := s.stream.(io.Seeker)
	return ok
}

// CanWrite reports whether the underlying stream supports writing.
func (s *StreamWrapper) CanWrite() bool {
	_, ok := s.stream.(io.Writer)
	return ok
}

// Length returns the length/size of the underlying stream by seeking to end and back.
func (s *StreamWrapper) Length() (int64, error) {
	current, err := s.stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	length, err := s.stream.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err2 := s.stream.Seek(current, io.SeekStart)
	if err2 != nil {
		return 0, err2
	}

	return length, nil
}

// Position returns the current read/write position of the underlying stream.
func (s *StreamWrapper) Position() (int64, error) {
	return s.stream.Seek(0, io.SeekCurrent)
}

// SetPosition sets the current read/write position of the underlying stream.
func (s *StreamWrapper) SetPosition(pos int64) error {
	_, err := s.stream.Seek(pos, io.SeekStart)
	return err
}

// Close closes the underlying stream if it implements io.Closer.
func (s *StreamWrapper) Close() error {
	if closer, ok := s.stream.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
