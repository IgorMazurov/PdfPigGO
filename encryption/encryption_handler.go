package encryption

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// EncryptionHandler manages decryption of tokens in a PDF document where encryption is used.
type EncryptionHandler interface {
	// Decrypt decrypts the given token using the provided indirect reference.
	Decrypt(reference core.IndirectReference, token tokens.Token) tokens.Token
}
