package acroforms

// SignatureFlags specifies document level characteristics for any signature
// fields in the document's AcroForm. Multiple flags can be combined using
// bitwise OR.
type SignatureFlags int

const (
	// SignaturesExist indicates that the document contains at least one signature field.
	SignaturesExist SignatureFlags = 1 << 0

	// AppendOnly indicates that the document contains signatures which may be
	// invalidated if the file is saved in a way which alters its previous content
	// rather than simply appending new content.
	AppendOnly = 1 << 1
)
