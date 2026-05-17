package core

import "fmt"

// PdfDocumentFormatException is returned when PDF document contents do not match
// the specification in a way that renders the document unreadable.
type PdfDocumentFormatException struct {
	Message   string
	Inner     error
}

// Error implements the error interface.
func (e *PdfDocumentFormatException) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Inner)
	}
	return e.Message
}

// Unwrap returns the inner error for errors.Is/errors.As support.
func (e *PdfDocumentFormatException) Unwrap() error {
	return e.Inner
}

// NewPdfDocumentFormatException creates a new PdfDocumentFormatException with the given message.
func NewPdfDocumentFormatException(message string) *PdfDocumentFormatException {
	return &PdfDocumentFormatException{Message: message}
}

// NewPdfDocumentFormatExceptionWithInner creates a new PdfDocumentFormatException
// with a message and an inner error.
func NewPdfDocumentFormatExceptionWithInner(message string, inner error) *PdfDocumentFormatException {
	return &PdfDocumentFormatException{Message: message, Inner: inner}
}
