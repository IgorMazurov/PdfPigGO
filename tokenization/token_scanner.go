package tokenization

import "github.com/uglytoad/pdfpig/go/tokens"

// TokenScanner scans input for PostScript/PDF tokens.
type TokenScanner interface {
	// Advance moves the scanner to the next token in the input.
	// Returns true if a token was successfully read, false otherwise.
	Advance() bool

	// Current returns the most recently scanned token.
	// The result is valid only after a successful call to Advance.
	Current() tokens.Token
}
