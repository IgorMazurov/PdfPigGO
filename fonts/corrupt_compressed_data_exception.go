package fonts

import "fmt"

// CorruptCompressedDataException is returned when a PDF contains an invalid
// compressed data stream.
type CorruptCompressedDataException struct {
	Message string
	Inner   error
}

// Error implements the error interface.
func (e *CorruptCompressedDataException) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Inner)
	}
	return e.Message
}

// Unwrap returns the inner error for errors.Is/errors.As support.
func (e *CorruptCompressedDataException) Unwrap() error {
	return e.Inner
}

// NewCorruptCompressedDataException creates a new CorruptCompressedDataException with the given message.
func NewCorruptCompressedDataException(message string) *CorruptCompressedDataException {
	return &CorruptCompressedDataException{Message: message}
}

// NewCorruptCompressedDataExceptionWithInner creates a new CorruptCompressedDataException
// with a message and an inner error.
func NewCorruptCompressedDataExceptionWithInner(message string, inner error) *CorruptCompressedDataException {
	return &CorruptCompressedDataException{Message: message, Inner: inner}
}
