package writer

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer/interfaces"
)

const defaultVersion = 1.2

// PdfStreamWriter is an alias to the interfaces package type to break import cycles.
type PdfStreamWriter = interfaces.PdfStreamWriter

// positionTrackingWriter wraps an io.Writer and tracks the total bytes written
// so Position() returns the current byte offset.
type positionTrackingWriter struct {
	writer  io.Writer
	written int64
}

func (p *positionTrackingWriter) Write(b []byte) (int, error) {
	n, err := p.writer.Write(b)
	p.written += int64(n)
	return n, err
}

func (p *positionTrackingWriter) Position() int64 {
	return p.written
}

// DebugPosition prints the current position to stderr (for debugging only).
func (p *positionTrackingWriter) DebugPosition(msg string) {
}

// Seek supports io.Seeker so that WriteCrossReferenceTable can determine the
// current stream position for startxref. Only SeekCurrent (whence=1) and
// SeekStart (whence=0 at offset 0) are meaningfully supported; others are
// no-ops that return the tracked position.
func (p *positionTrackingWriter) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekCurrent:
		return p.written + offset, nil
	case io.SeekStart:
		return offset, nil
	case io.SeekEnd:
		// We don't know the final size; return current position.
		return p.written, nil
	default:
		return p.written, fmt.Errorf("invalid whence: %d", whence)
	}
}

// Read delegates to the inner writer if it implements io.Reader.
func (p *positionTrackingWriter) Read(b []byte) (int, error) {
	if r, ok := p.writer.(io.Reader); ok {
		return r.Read(b)
	}
	return 0, io.ErrNoProgress
}

// defaultPdfStreamWriter is the concrete implementation of PdfStreamWriter that
// lazily flushes all tokens to the underlying stream, tracking byte offsets for
// each indirect reference so the cross-reference table can be written later.
type defaultPdfStreamWriter struct {
	stream        *positionTrackingWriter
	disposeStream bool
	tokenWriter   TokenWriter
	recordVersion func(float64)
	offsets       map[core.IndirectReference]int64
	currentNumber int64
	attemptDedup  bool
	writingPage   bool
	initialized   bool
}

// NewPdfStreamWriter creates a new PdfStreamWriter that writes to the given
// underlying writer. If disposeStream is true, Close will also close the
// underlying writer if it implements io.Closer. If tokenWriter is nil, a
// default TokenWriter is used. recordVersion is an optional callback invoked
// when the PDF version is initialized.
func NewPdfStreamWriter(
	baseStream io.Writer,
	disposeStream bool,
	tokenWriter TokenWriter,
	recordVersion func(float64),
) PdfStreamWriter {
	if baseStream == nil {
		panic("baseStream must not be nil")
	}

	if tokenWriter == nil {
		tokenWriter = NewTokenWriter()
	}

	return &defaultPdfStreamWriter{
		stream:        &positionTrackingWriter{writer: baseStream},
		disposeStream: disposeStream,
		tokenWriter:   tokenWriter,
		recordVersion: recordVersion,
		offsets:       make(map[core.IndirectReference]int64),
		currentNumber: 1,
		attemptDedup:  true,
	}
}

func (w *defaultPdfStreamWriter) AttemptDeduplication() bool {
	return w.attemptDedup
}

func (w *defaultPdfStreamWriter) SetAttemptDeduplication(v bool) {
	w.attemptDedup = v
}

func (w *defaultPdfStreamWriter) Stream() io.Writer {
	return w.stream
}

func (w *defaultPdfStreamWriter) WritingPageContents() bool {
	return w.tokenWriter.WritingPageContents()
}

func (w *defaultPdfStreamWriter) SetWritingPageContents(v bool) {
	w.tokenWriter.SetWritingPageContents(v)
}

// WriteToken writes a token to the stream, reserving a new object number and
// recording its byte offset. Returns the indirect reference token for the written object.
func (w *defaultPdfStreamWriter) WriteToken(token tokens.Token) *tokens.IndirectReferenceToken {
	if !w.initialized {
		w.InitializePdf(defaultVersion)
	}

	ir := w.ReserveObjectNumber()
	pos := w.stream.Position()
	ref := ir.Data()
	w.offsets[ref] = pos

	obj := tokens.NewObjectToken(core.File(pos), ref, token)
	w.tokenWriter.WriteToken(obj, w.stream) //nolint:errcheck
	return ir
}

// WriteTokenAt writes a token at the given reserved indirect reference and
// records its byte offset.
func (w *defaultPdfStreamWriter) WriteTokenAt(token tokens.Token, ref *tokens.IndirectReferenceToken) *tokens.IndirectReferenceToken {
	if !w.initialized {
		w.InitializePdf(defaultVersion)
	}

	pos := w.stream.Position()
	refData := ref.Data()
	w.offsets[refData] = pos

	obj := tokens.NewObjectToken(core.File(pos), refData, token)
	w.tokenWriter.WriteToken(obj, w.stream) //nolint:errcheck
	return ref
}

// ReserveObjectNumber reserves the next object number and returns an indirect
// reference token for it.
func (w *defaultPdfStreamWriter) ReserveObjectNumber() *tokens.IndirectReferenceToken {
	ref, err := core.NewIndirectReference(w.currentNumber, 0)
	if err != nil {
		panic(fmt.Sprintf("failed to create indirect reference: %v", err))
	}
	w.currentNumber++
	return tokens.NewIndirectReferenceToken(ref)
}

// InitializePdf writes the PDF header to the stream with the given version number.
func (w *defaultPdfStreamWriter) InitializePdf(version float64) {
	if w.recordVersion != nil {
		w.recordVersion(version)
	}

	fmt.Fprintf(w.stream, "%%PDF-%.1f\n", version)
	w.stream.Write([]byte{37, 169, 205, 196, 210}) //nolint:errcheck
	w.stream.Write([]byte{'\n'})                   //nolint:errcheck
	w.initialized = true
}

// CompletePdf writes the cross-reference table and trailer dictionary to finalize the PDF.
func (w *defaultPdfStreamWriter) CompletePdf(catalogRef *tokens.IndirectReferenceToken, docInfoRef *tokens.IndirectReferenceToken) {
	var infoRef *core.IndirectReference
	if docInfoRef != nil {
		d := docInfoRef.Data()
		infoRef = &d
	}
	
	w.tokenWriter.WriteCrossReferenceTable(w.offsets, catalogRef.Data(), w.stream, infoRef) //nolint:errcheck
}

// Close disposes of resources held by the stream writer. If disposeStream is
// true and the underlying writer implements io.Closer, it will be closed.
func (w *defaultPdfStreamWriter) Close() error {
	if w.disposeStream {
		if closer, ok := w.stream.writer.(io.Closer); ok {
			return closer.Close()
		}
	}
	return nil
}
