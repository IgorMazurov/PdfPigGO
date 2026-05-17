package core

import "io"

// Writeable indicates that a data structure can be written to an output stream.
type Writeable interface {
	// Write writes the data to the output stream.
	Write(stream io.Writer) error
}
