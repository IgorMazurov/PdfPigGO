package page

// PageXmlTextLineOrderSimpleType represents the text line ordering direction in PAGE XML documents.
type PageXmlTextLineOrderSimpleType byte

const (
	// TextLineTopToBottom indicates top-to-bottom text line order.
	TextLineTopToBottom PageXmlTextLineOrderSimpleType = iota

	// TextLineBottomToTop indicates bottom-to-top text line order.
	TextLineBottomToTop

	// TextLineLeftToRight indicates left-to-right text line order.
	TextLineLeftToRight

	// TextLineRightToLeft indicates right-to-left text line order.
	TextLineRightToLeft
)
