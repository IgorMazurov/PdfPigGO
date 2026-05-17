package page

// PageXmlPrintSpace determines the effective area on the paper of a printed page.
// Its size is equal for all pages of a book (exceptions: titlepage, multipage pictures).
// It contains all living elements (except marginals) like body type, footnotes,
// headings, running titles. It does not contain pagenumber (if not part of running
// title), marginals, signature mark, preview words.
type PageXmlPrintSpace struct {
	// Coords defines the polygon outline of the print space area.
	Coords *PageXmlCoords `xml:"coords,omitempty"`
}
