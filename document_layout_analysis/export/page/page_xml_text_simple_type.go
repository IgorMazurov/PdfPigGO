package page

// PageXmlTextSimpleType represents text type classification in PAGE XML documents.
type PageXmlTextSimpleType byte

const (
	// TextParagraph represents paragraph text.
	TextParagraph PageXmlTextSimpleType = iota

	// Heading represents heading text.
	Heading

	// Caption represents caption text.
	Caption

	// Header represents header text.
	Header

	// Footer represents footer text.
	Footer

	// PageNumber represents page number text.
	PageNumber

	// DropCapital represents a drop capital, a letter at the beginning of a word that is bigger than the usual character size. Usually to start a chapter.
	DropCapital

	// Credit represents credit text.
	Credit

	// Floating represents floating text.
	Floating

	// SignatureMark represents signature mark text.
	SignatureMark

	// CatchWord represents catch word text.
	CatchWord

	// Marginalia represents marginalia text.
	Marginalia

	// FootNote represents foot note text.
	FootNote

	// FootNoteContinued represents continued foot note text.
	FootNoteContinued

	// EndNote represents end note text.
	EndNote

	// TocEntry represents table of content entry text.
	TocEntry

	// LisLabel represents list label text.
	LisLabel

	// TextOther represents other text type.
	TextOther
)
