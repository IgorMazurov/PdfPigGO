package writer

import (
	"bytes"
	"io"

	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer/interfaces"
)

var _ interfaces.PdfStreamWriter = (*PdfDedupStreamWriter)(nil)

// PdfDedupStreamWriter writes PDF content to a stream with deduplication support.
// It detects identical object contents using FNV-1a hashing and reuses existing
// indirect references, reducing file size for documents containing duplicate objects.
type PdfDedupStreamWriter struct {
	*defaultPdfStreamWriter

	hashes map[string]*tokens.IndirectReferenceToken
	ms     bytes.Buffer
}

// NewPdfDedupStreamWriter creates a new PdfDedupStreamWriter that writes to the
// given underlying writer. If disposeStream is true, Close will also close the
// underlying writer if it implements io.Closer. If tokenWriter is nil, a default
// TokenWriter (NewNoTextTokenWriter) is used. recordVersion is an optional callback
// invoked when the PDF version is initialized during the first WriteToken call.
func NewPdfDedupStreamWriter(
	baseStream io.Writer,
	disposeStream bool,
	tokenWriter TokenWriter,
	recordVersion func(float64),
) *PdfDedupStreamWriter {
	return &PdfDedupStreamWriter{
		defaultPdfStreamWriter: NewPdfStreamWriter(baseStream, disposeStream, tokenWriter, recordVersion).(*defaultPdfStreamWriter),
		hashes:                 make(map[string]*tokens.IndirectReferenceToken),
	}
}

// WriteToken writes a token to the stream and returns an indirect reference.
// If deduplication is enabled and identical content was already written, the
// existing reference is returned instead of writing duplicate data.
func (w *PdfDedupStreamWriter) WriteToken(token tokens.Token) *tokens.IndirectReferenceToken {
	if !w.initialized {
		w.InitializePdf(defaultVersion)
	}

	w.ms.Reset()
	w.tokenWriter.WriteToken(token, &w.ms) //nolint:errcheck
	contents := w.ms.Bytes()
	key := string(contents)

	if w.attemptDedup {
		if ir, ok := w.hashes[key]; ok {
			return ir
		}
	}

	ir := w.ReserveObjectNumber()

	if w.attemptDedup {
		w.hashes[key] = ir
	}

	pos := w.stream.Position()
	ref := ir.Data()
	w.offsets[ref] = pos

	w.tokenWriter.WriteObject(ir.Data().ObjectNumber(), ir.Data().Generation(), contents, w.stream) //nolint:errcheck
	return ir
}

// WriteTokenAt writes a token at a pre-reserved indirect reference and records
// its byte offset in the hash map for future deduplication lookups.
func (w *PdfDedupStreamWriter) WriteTokenAt(token tokens.Token, ref *tokens.IndirectReferenceToken) *tokens.IndirectReferenceToken {
	if !w.initialized {
		w.InitializePdf(defaultVersion)
	}

	w.ms.Reset()
	w.tokenWriter.WriteToken(token, &w.ms) //nolint:errcheck
	contents := w.ms.Bytes()
	key := string(contents)

	w.hashes[key] = ref

	pos := w.stream.Position()
	refData := ref.Data()
	w.offsets[refData] = pos

	w.tokenWriter.WriteObject(ref.Data().ObjectNumber(), ref.Data().Generation(), contents, w.stream) //nolint:errcheck
	return ref
}

// Close releases resources held by the writer including clearing the deduplication
// cache. If disposeStream was set during construction and the underlying stream
// implements io.Closer, it will be closed.
func (w *PdfDedupStreamWriter) Close() error {
	w.hashes = nil
	return w.defaultPdfStreamWriter.Close()
}

// --- FNV-1a Hash ---

// fnvHashOffset is the starting point of the FNV-1a 32-bit hash.
const fnvHashOffset = 2166136261

// fnvHashPrime is the prime number used to compute the FNV-1a hash.
const fnvHashPrime = 16777619

// fnvHash computes an FNV-1a hash over the given byte slice. This matches the
// Fowler/Noll/Vo algorithm used by the original C# implementation for efficient
// byte-array comparison in the deduplication dictionary.
func fnvHash(data []byte) int {
	h := fnvHashOffset
	for i := range data {
		h ^= int(data[i])
		h *= fnvHashPrime
	}
	return h
}
