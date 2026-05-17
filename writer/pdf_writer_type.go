package writer

// PdfWriterType represents the type of PDF writer to use when generating documents.
type PdfWriterType int

const (
	// PdfWriterDefault is the default output writer.
	PdfWriterDefault PdfWriterType = iota

	// PdfWriterObjectInMemoryDedup de-duplicates objects while writing but
	// requires keeping references in memory.
	PdfWriterObjectInMemoryDedup
)
