package interfaces

import (
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfStreamWriter defines the interface for writing PDF content to a stream.
type PdfStreamWriter interface {
	io.Closer

	AttemptDeduplication() bool
	SetAttemptDeduplication(v bool)
	Stream() io.Writer
	WritingPageContents() bool
	SetWritingPageContents(v bool)
	WriteToken(token tokens.Token) *tokens.IndirectReferenceToken
	WriteTokenAt(token tokens.Token, ref *tokens.IndirectReferenceToken) *tokens.IndirectReferenceToken
	ReserveObjectNumber() *tokens.IndirectReferenceToken
	InitializePdf(version float64)
	CompletePdf(catalogRef *tokens.IndirectReferenceToken, docInfoRef *tokens.IndirectReferenceToken)
}
