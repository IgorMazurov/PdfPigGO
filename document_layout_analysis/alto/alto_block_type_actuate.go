package alto

// AltoBlockTypeActuate defines the actuation type for an ALTO block reference
// as specified in the XLink namespace (http://www.w3.org/1999/xlink).
type AltoBlockTypeActuate string

const (
	// AltoBlockTypeActuateOnLoad indicates the link is activated on load.
	AltoBlockTypeActuateOnLoad AltoBlockTypeActuate = "onLoad"
	// AltoBlockTypeActuateOnRequest indicates the link is activated on request.
	AltoBlockTypeActuateOnRequest AltoBlockTypeActuate = "onRequest"
	// AltoBlockTypeActuateOther indicates some other actuation type.
	AltoBlockTypeActuateOther AltoBlockTypeActuate = "other"
	// AltoBlockTypeActuateNone indicates no actuation.
	AltoBlockTypeActuateNone AltoBlockTypeActuate = "none"
)

// String returns the XML string representation of this actuation type.
func (a AltoBlockTypeActuate) String() string {
	return string(a)
}
