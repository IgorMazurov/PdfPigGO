package alto

// AltoStyles represents style definitions in an ALTO document.
type AltoStyles struct {
	TextStyle      []AltoTextStyle      `xml:"TextStyle"`
	ParagraphStyle []AltoParagraphStyle `xml:"ParagraphStyle"`
}
