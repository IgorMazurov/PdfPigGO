package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// EmbeddedFile represents a file embedded in a PDF document for document references.
type EmbeddedFile struct {
	// Name is the name given to this embedded file in the document's name tree.
	Name string

	// FileSpecification is the specification of the path to the file.
	FileSpecification string

	// Memory holds the decrypted memory of the file.
	Memory []byte

	// Stream is the underlying embedded file stream.
	Stream *tokens.StreamToken
}

// NewEmbeddedFile creates a new EmbeddedFile with the given parameters.
func NewEmbeddedFile(name, fileSpecification string, bytes []byte, stream *tokens.StreamToken) (*EmbeddedFile, error) {
	if name == "" {
		return nil, fmt.Errorf("embedded file name cannot be empty")
	}

	if stream == nil {
		return nil, fmt.Errorf("embedded file stream cannot be nil")
	}

	return &EmbeddedFile{
		Name:            name,
		FileSpecification: fileSpecification,
		Memory:          bytes,
		Stream:          stream,
	}, nil
}

// Bytes returns the decrypted bytes of the file.
func (e *EmbeddedFile) Bytes() []byte {
	return e.Memory
}

// String returns the string representation of the embedded file.
func (e *EmbeddedFile) String() string {
	return fmt.Sprintf("%s: %v.", e.Name, e.Stream.StreamDictionary)
}

var _ fmt.Stringer = (*EmbeddedFile)(nil)
