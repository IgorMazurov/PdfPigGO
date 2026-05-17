package tokenization

import "github.com/uglytoad/pdfpig/go/core"

// SeekableTokenScanner extends TokenScanner with seeking capabilities
// on the underlying input data.
type SeekableTokenScanner interface {
	TokenScanner

	// StackDepthGuard returns the guard object used to track and limit
	// stack depth during recursive operations.
	StackDepthGuard() *core.StackDepthGuard

	// Seek moves the scanner to the specified position with the given whence.
	Seek(offset int64, whence int) (int64, error)

	// CurrentPosition returns the current byte offset in the input.
	CurrentPosition() int64

	// Length returns the total length of the data represented by this scanner.
	Length() int64

	// RegisterCustomTokenizer adds support for a custom tokenizer identified
	// by its first matching byte.
	RegisterCustomTokenizer(firstByte byte, t Tokenizer)

	// DeregisterCustomTokenizer removes a previously registered custom tokenizer.
	DeregisterCustomTokenizer(t Tokenizer)
}
