package page

// PageXmlGraphemeBaseCharType represents the character type within a grapheme base for PAGE XML export.
type PageXmlGraphemeBaseCharType byte

const (
	// Base indicates a base character.
	Base PageXmlGraphemeBaseCharType = iota

	// Combining indicates a combining character.
	Combining
)
