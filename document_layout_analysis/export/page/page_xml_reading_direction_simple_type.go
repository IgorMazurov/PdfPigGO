package page

// PageXmlReadingDirectionSimpleType represents the reading direction for PAGE XML export.
type PageXmlReadingDirectionSimpleType byte

const (
	// LeftToRight indicates left-to-right reading direction.
	LeftToRight PageXmlReadingDirectionSimpleType = iota

	// RightToLeft indicates right-to-left reading direction.
	RightToLeft

	// TopToBottom indicates top-to-bottom reading direction.
	TopToBottom

	// BottomToTop indicates bottom-to-top reading direction.
	BottomToTop
)
