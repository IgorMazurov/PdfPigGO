package core

// LineCapStyle defines the shape used at the ends of open subpaths.
type LineCapStyle byte

const (
	// LineCapButting produces square ends perpendicular to the line segment.
	LineCapButting LineCapStyle = iota
	// LineCapRound produces round ends centered on endpoints.
	LineCapRound
	// LineCapProjecting produces projecting square ends.
	LineCapProjecting
)

// LineJoinStyle defines the shape used at the corners of paths that have adjoining segments.
type LineJoinStyle byte

const (
	// LineJoinMiter produces pointed corners via outward extension.
	LineJoinMiter LineJoinStyle = iota
	// LineJoinRound produces round corners.
	LineJoinRound
	// LineJoinBevel produces beveled corners.
	LineJoinBevel
)
