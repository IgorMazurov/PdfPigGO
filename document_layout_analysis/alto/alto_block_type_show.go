package alto

// AltoBlockTypeShow defines the show type for an ALTO block reference
// as specified in the XLink namespace (http://www.w3.org/1999/xlink).
type AltoBlockTypeShow string

const (
	// AltoBlockTypeShowNew indicates the content is shown in a new window.
	AltoBlockTypeShowNew AltoBlockTypeShow = "new"
	// AltoBlockTypeShowReplace indicates the content replaces the current view.
	AltoBlockTypeShowReplace AltoBlockTypeShow = "replace"
	// AltoBlockTypeShowEmbed indicates the content is embedded in place.
	AltoBlockTypeShowEmbed AltoBlockTypeShow = "embed"
	// AltoBlockTypeShowOther indicates some other show type.
	AltoBlockTypeShowOther AltoBlockTypeShow = "other"
	// AltoBlockTypeShowNone indicates no show behavior.
	AltoBlockTypeShowNone AltoBlockTypeShow = "none"
)

// String returns the XML string representation of this show type.
func (a AltoBlockTypeShow) String() string {
	return string(a)
}
