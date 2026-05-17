package alto

// AltoParagraphStyleAlign represents the alignment of a paragraph style in ALTO documents.
type AltoParagraphStyleAlign string

func (a AltoParagraphStyleAlign) String() string { return string(a) }

const (
	AltoParagraphStyleAlignLeft   AltoParagraphStyleAlign = "Left"
	AltoParagraphStyleAlignRight  AltoParagraphStyleAlign = "Right"
	AltoParagraphStyleAlignCenter AltoParagraphStyleAlign = "Center"
	AltoParagraphStyleAlignBlock  AltoParagraphStyleAlign = "Block"
)
