package encryption

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// NoOpEncryptionHandler is a no-operation encryption handler that returns tokens unchanged.
// It is used for unencrypted PDF documents where no decryption is required.
type NoOpEncryptionHandler struct{}

// NoOpInstance is the singleton instance of NoOpEncryptionHandler.
var NoOpInstance = &NoOpEncryptionHandler{}

func init() {
	var _ EncryptionHandler = (*NoOpEncryptionHandler)(nil)
}

// Decrypt returns the token unchanged, performing no decryption.
func (h *NoOpEncryptionHandler) Decrypt(reference core.IndirectReference, token tokens.Token) tokens.Token {
	return token
}
