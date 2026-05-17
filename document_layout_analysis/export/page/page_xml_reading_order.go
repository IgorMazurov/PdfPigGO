package page

// PageXmlReadingOrder represents the reading order of elements on a PAGE XML page,
// containing either an ordered or unordered group with optional confidence value.
type PageXmlReadingOrder struct {
	// Item holds the child element — either a PageXmlOrderedGroup or a PageXmlUnorderedGroup.
	Item any `xml:"OrderedGroup,UnorderedGroup"`

	// Conf is the confidence value between 0 and 1. Nil indicates no confidence was set.
	Conf *float32 `xml:"conf,attr,omitempty"`
}
