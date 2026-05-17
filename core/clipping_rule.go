package core

// FillingRule defines rules for determining which points lie inside/outside a path.
type FillingRule byte

const (
	// FillingRuleNone means no rule is applied.
	FillingRuleNone FillingRule = 0

	// FillingRuleEvenOdd determines whether a point is inside a path by drawing a ray from that point in
	// any direction and counting the number of path segments crossing the ray. If odd, the point is inside; if even, outside.
	FillingRuleEvenOdd FillingRule = 1

	// FillingRuleNonZeroWinding determines whether a point is inside by drawing a ray to infinity and tracking crossings:
	// +1 for left-to-right, -1 for right-to-left. Nonzero result means the point is inside.
	FillingRuleNonZeroWinding FillingRule = 2
)
