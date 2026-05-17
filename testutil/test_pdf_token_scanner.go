package testutil

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var _ tokenization.SeekableTokenScanner = (*TestPdfTokenScanner)(nil)

// TestPdfTokenScanner is a test implementation of SeekableTokenScanner.
type TestPdfTokenScanner struct {
	currentToken tokens.Token
	position     int64
	Objects      map[core.IndirectReference]*tokens.ObjectToken
}

// NewTestPdfTokenScanner creates a new TestPdfTokenScanner with an initialized objects map.
func NewTestPdfTokenScanner() *TestPdfTokenScanner {
	return &TestPdfTokenScanner{
		Objects: make(map[core.IndirectReference]*tokens.ObjectToken),
	}
}

// Advance moves the scanner to the next token (not implemented for test helper).
func (s *TestPdfTokenScanner) Advance() bool {
	panic("not implemented")
}

// Current returns the most recently scanned token.
func (s *TestPdfTokenScanner) Current() tokens.Token {
	return s.currentToken
}

// CurrentToken returns the current token for test manipulation.
func (s *TestPdfTokenScanner) CurrentToken() tokens.Token {
	return s.currentToken
}

// SetCurrentToken sets the current token directly for test setup.
func (s *TestPdfTokenScanner) SetCurrentToken(t tokens.Token) {
	s.currentToken = t
}

// StackDepthGuard returns an infinite depth guard.
func (s *TestPdfTokenScanner) StackDepthGuard() *core.StackDepthGuard {
	return core.Infinite
}

// Seek moves the scanner to the specified position with the given whence.
func (s *TestPdfTokenScanner) Seek(offset int64, whence int) (int64, error) {
	panic("not implemented")
}

// CurrentPosition returns the current byte offset in the input.
func (s *TestPdfTokenScanner) CurrentPosition() int64 {
	return s.position
}

// SetCurrentPosition sets the current position directly for test setup.
func (s *TestPdfTokenScanner) SetCurrentPosition(pos int64) {
	s.position = pos
}

// Length returns the total length of the data represented by this scanner.
func (s *TestPdfTokenScanner) Length() int64 {
	return 10
}

// RegisterCustomTokenizer adds support for a custom tokenizer.
func (s *TestPdfTokenScanner) RegisterCustomTokenizer(firstByte byte, t tokenization.Tokenizer) {
	panic("not implemented")
}

// DeregisterCustomTokenizer removes a previously registered custom tokenizer.
func (s *TestPdfTokenScanner) DeregisterCustomTokenizer(t tokenization.Tokenizer) {
	panic("not implemented")
}

// Get retrieves an ObjectToken by its indirect reference from the objects cache.
func (s *TestPdfTokenScanner) Get(reference core.IndirectReference) *tokens.ObjectToken {
	return s.Objects[reference]
}

// ReplaceToken updates or inserts an ObjectToken for the given reference in the cache.
func (s *TestPdfTokenScanner) ReplaceToken(reference core.IndirectReference, token tokens.Token) {
	if obj, ok := token.(*tokens.ObjectToken); ok {
		s.Objects[reference] = obj
	}
}

// Close releases any resources held by the scanner (no-op for test implementation).
func (s *TestPdfTokenScanner) Close() error {
	return nil
}
