package content

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var (
	errNilFilterProvider = errors.New("filter provider cannot be nil")
	errNilStream         = errors.New("stream token cannot be nil")
)

// XmpMetadata wraps an XML-based Extensible Metadata Platform (XMP) document.
// These XML documents are embedded in PDFs to provide metadata about objects
// (the entire document, images, etc.). They can be present as plain text or
// encoded/encrypted streams.
type XmpMetadata struct {
	metadataStreamToken *tokens.StreamToken
	filterProvider      LookupFilterProvider
	pdfScanner          tokenization.PdfTokenScanner
}

// NewXmpMetadata creates a new XmpMetadata from the given stream token,
// filter provider, and PDF token scanner.
func NewXmpMetadata(
	stream *tokens.StreamToken,
	filterProvider LookupFilterProvider,
	pdfScanner tokenization.PdfTokenScanner,
) (*XmpMetadata, error) {
	if filterProvider == nil {
		return nil, errNilFilterProvider
	}

	if stream == nil {
		return nil, errNilStream
	}

	return &XmpMetadata{
		metadataStreamToken: stream,
		filterProvider:      filterProvider,
		pdfScanner:          pdfScanner,
	}, nil
}

// MetadataStreamToken returns the underlying StreamToken for this metadata.
func (x *XmpMetadata) MetadataStreamToken() *tokens.StreamToken {
	return x.metadataStreamToken
}

// GetXmlBytes returns the decoded bytes for the metadata stream.
// The bytes can be interpreted as a sequence of plain-text bytes with any
// filters removed.
func (x *XmpMetadata) GetXmlBytes() []byte {
	return x.filterProvider.DecodeStream(x.metadataStreamToken, x.pdfScanner)
}

// GetXmlString returns the metadata stream as a Latin-1 decoded XML string.
// This is the Go equivalent of C#'s GetXDocument() — callers should parse
// the returned string using encoding/xml for full document access, since
// Go has no direct XDocument equivalent.
func (x *XmpMetadata) GetXmlString() string {
	return core.BytesAsLatin1String(x.GetXmlBytes())
}
