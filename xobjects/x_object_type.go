package xobjects

// XObjectType indicates the kind of external object stored in a PDF XObject
// dictionary.
type XObjectType byte

const (
	// Image represents an image XObject.
	Image XObjectType = iota

	// Form represents a form XObject (reusable content stream).
	Form

	// PostScript represents a PostScript XObject (rarely used in modern PDFs).
	PostScript
)

var xObjectTypeMap = map[string]XObjectType{
	"Image":    Image,
	"Form":     Form,
	"PS":       PostScript,
}

// ParseXObjectType converts a PDF name string to the corresponding XObjectType.
// Returns (typ, true) if recognized, or (0, false) otherwise.
func ParseXObjectType(s string) (XObjectType, bool) {
	typ, ok := xObjectTypeMap[s]
	return typ, ok
}
