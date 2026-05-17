package util

// StringSplitter iteratively splits a string by a single separator byte,
// returning one token at a time without allocating intermediate slices for
// the full split result. This mirrors C# ref struct StringSplitter.
type StringSplitter struct {
	text      []byte
	separator byte
	position  int
}

// NewStringSplitter creates a new StringSplitter that splits text by separator.
func NewStringSplitter(text string, separator byte) StringSplitter {
	return StringSplitter{
		text:      []byte(text),
		separator: separator,
		position:  0,
	}
}

// TryRead reads the next token from the splitter. It returns the token and true
// if a token was read, or an empty string and false if no more tokens remain.
func (s *StringSplitter) TryRead() (string, bool) {
	if s.IsEof() {
		return "", false
	}

	idx := findByte(s.text[s.position:], s.separator)

	if idx >= 0 {
		result := string(s.text[s.position : s.position+idx])
		s.position += idx + 1
		return result, true
	}

	result := string(s.text[s.position:])
	s.position = len(s.text)
	return result, true
}

// Read reads the next token from the splitter. It panics if no more tokens remain.
func (s *StringSplitter) Read() string {
	if s.IsEof() {
		panic("end of stream")
	}

	start := s.position
	idx := findByte(s.text[s.position:], s.separator)

	if idx >= 0 {
		s.position += idx + 1
		return string(s.text[start : start+idx])
	}

	s.position = len(s.text)
	return string(s.text[start:])
}

// IsEof reports whether the splitter has reached the end of its input.
func (s *StringSplitter) IsEof() bool {
	return s.position == len(s.text)
}

// findByte returns the index of the first occurrence of b in data, or -1 if not found.
func findByte(data []byte, b byte) int {
	for i := range data {
		if data[i] == b {
			return i
		}
	}
	return -1
}
