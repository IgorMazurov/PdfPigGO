package actions

// UriAction represents a PDF action that opens a URI (PDF reference 8.5.1).
type UriAction struct {
	*PdfAction

	// Uri is the uniform resource identifier to open.
	Uri string
}

// NewUriAction creates a new UriAction with the given URI.
func NewUriAction(uri string) *UriAction {
	return &UriAction{
		PdfAction: NewPdfAction(URI),
		Uri:       uri,
	}
}
