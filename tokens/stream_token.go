package tokens

import (
	"bytes"
	"fmt"
)

// StreamToken represents a PDF stream, consisting of a dictionary followed by
// zero or more bytes bracketed between the keywords stream and endstream. The
// bytes may be compressed by application of zero or more filters specified in
// the StreamDictionary.
type StreamToken struct {
	StreamDictionary *DictionaryToken
	data             []byte
}

var _ Token = (*StreamToken)(nil)
var _ DataToken[[]byte] = (*StreamToken)(nil)

// NewStreamToken creates a new StreamToken with the given dictionary and data.
func NewStreamToken(dict *DictionaryToken, data []byte) (*StreamToken, error) {
	if dict == nil {
		return nil, fmt.Errorf("stream dictionary cannot be nil")
	}

	if data == nil {
		return nil, fmt.Errorf("stream data cannot be nil")
	}

	return &StreamToken{
		StreamDictionary: dict,
		data:             data,
	}, nil
}

// Data returns the compressed byte data of the stream.
func (s *StreamToken) Data() []byte {
	return s.data
}

// Equals reports whether other is a StreamToken with equivalent dictionary and data.
func (s *StreamToken) Equals(other Token) bool {
	o, ok := other.(*StreamToken)
	if !ok || o == nil {
		return false
	}

	if s == o {
		return true
	}

	if !s.StreamDictionary.Equals(o.StreamDictionary) {
		return false
	}

	return bytes.Equal(s.data, o.data)
}

// String returns the string representation of the stream token.
func (s *StreamToken) String() string {
	return fmt.Sprintf("Length: %d, Dictionary: %s", len(s.data), s.StreamDictionary)
}
