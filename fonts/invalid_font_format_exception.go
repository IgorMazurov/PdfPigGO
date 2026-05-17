package fonts

import "fmt"

// InvalidFontFormatException is returned when an error is encountered parsing
// a font from the PDF document, where the format of the font program or
// dictionary does not meet the specification.
type InvalidFontFormatException struct {
	Message string
	Inner   error
}

// Error implements the error interface.
func (e *InvalidFontFormatException) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Inner)
	}
	return e.Message
}

// Unwrap returns the inner error for errors.Is/errors.As support.
func (e *InvalidFontFormatException) Unwrap() error {
	return e.Inner
}

// NewInvalidFontFormatException creates a new InvalidFontFormatException with the given message.
func NewInvalidFontFormatException(message string) *InvalidFontFormatException {
	return &InvalidFontFormatException{Message: message}
}

// NewInvalidFontFormatExceptionWithInner creates a new InvalidFontFormatException
// with a message and an inner error.
func NewInvalidFontFormatExceptionWithInner(message string, inner error) *InvalidFontFormatException {
	return &InvalidFontFormatException{Message: message, Inner: inner}
}
