package reading_order_detector

import "github.com/uglytoad/pdfpig/go/core"

// IntervalRelations represents Allen's interval thirteen relations for
// qualitative spatial reasoning between document objects on a page.
// See https://en.wikipedia.org/wiki/Allen%27s_interval_algebra
type IntervalRelations int

const (
	// Unknown indicates an undetermined interval relation.
	Unknown IntervalRelations = iota

	// Precedes means X takes place before Y.
	Precedes

	// Meets means the end of X coincides with the start of Y.
	Meets

	// Overlaps means X and Y partially overlap, each starting before the other ends.
	Overlaps

	// Starts means X and Y share the same start but Y extends further.
	Starts

	// During means X is entirely contained within Y.
	During

	// Finishes means X and Y share the same end but X starts later.
	Finishes

	// PrecedesI is the inverse of Precedes (Y precedes X).
	PrecedesI

	// MeetsI is the inverse of Meets (Y meets X).
	MeetsI

	// OverlapsI is the inverse of Overlaps (Y overlaps X).
	OverlapsI

	// StartsI is the inverse of Starts (Y starts X).
	StartsI

	// DuringI is the inverse of During (Y contains X).
	DuringI

	// FinishesI is the inverse of Finishes (Y finishes X).
	FinishesI

	// Equals means X and Y are identical in extent.
	Equals
)

// GetRelationX computes the Thick Boundary Rectangle Relations (TBRR) for the
// X coordinate between two rectangles. For every pair of document objects a and b,
// one X interval relation holds. If the pair is considered in reversed order,
// the inverse interval relation holds. T is the tolerance: coordinates closer
// than T are considered equal.
func GetRelationX(a, b core.PdfRectangle, T float64) IntervalRelations {
	aLeft := a.Left()
	aRight := a.Right()
	bLeft := b.Left()
	bRight := b.Right()

	if bLeft-T <= aLeft && aLeft <= bLeft+T && bRight-T <= aRight && aRight <= bRight+T {
		return Equals
	}

	if bLeft-T <= aRight && aRight <= bLeft+T {
		return Meets
	} else if aLeft-T <= bRight && bRight <= aLeft+T {
		return MeetsI
	}

	if bLeft-T <= aLeft && aLeft <= bLeft+T && aRight < bRight-T {
		return Starts
	} else if aLeft-T <= bLeft && bLeft <= aLeft+T && bRight < aRight-T {
		return StartsI
	}

	if aLeft > bLeft+T && bRight-T <= aRight && aRight <= bRight+T {
		return Finishes
	} else if bLeft > aLeft+T && aRight-T <= bRight && bRight <= aRight+T {
		return FinishesI
	}

	if aLeft > bLeft+T && aRight < bRight-T {
		return During
	} else if bLeft > aLeft+T && bRight < aRight-T {
		return DuringI
	}

	if aLeft < bLeft-T && bLeft+T < aRight && aRight < bRight-T {
		return Overlaps
	} else if bLeft < aLeft-T && aLeft+T < bRight && bRight < aRight-T {
		return OverlapsI
	}

	if aRight < bLeft-T {
		return Precedes
	} else if bRight < aLeft-T {
		return PrecedesI
	}

	return Unknown
}

// GetRelationY computes the Thick Boundary Rectangle Relations (TBRR) for the
// Y coordinate between two rectangles. For every pair of document objects a and b,
// one Y interval relation holds. If the pair is considered in reversed order,
// the inverse interval relation holds. T is the tolerance: coordinates closer
// than T are considered equal.
func GetRelationY(a, b core.PdfRectangle, T float64) IntervalRelations {
	aTop := a.Top()
	aBottom := a.Bottom()
	bTop := b.Top()
	bBottom := b.Bottom()

	if bTop-T <= aTop && aTop <= bTop+T && bBottom-T <= aBottom && aBottom <= bBottom+T {
		return Equals
	}

	if aTop-T <= bBottom && bBottom <= aTop+T {
		return MeetsI
	} else if bTop-T <= aBottom && aBottom <= bTop+T {
		return Meets
	}

	if bTop-T <= aTop && aTop <= bTop+T && aBottom < bBottom-T {
		return StartsI
	} else if aTop-T <= bTop && bTop <= aTop+T && bBottom < aBottom-T {
		return Starts
	}

	if aTop > bTop+T && bBottom-T <= aBottom && aBottom <= bBottom+T {
		return FinishesI
	} else if bTop > aTop+T && aBottom-T <= bBottom && bBottom <= aBottom+T {
		return Finishes
	}

	if aTop > bTop+T && aBottom < bBottom-T {
		return DuringI
	} else if bTop > aTop+T && bBottom < aBottom-T {
		return During
	}

	if aTop < bTop-T && bBottom+T < aTop && aBottom < bBottom-T {
		return OverlapsI
	} else if bTop < aTop-T && aBottom+T < bTop && bBottom < aBottom-T {
		return Overlaps
	}

	if aBottom < bTop-T {
		return PrecedesI
	} else if bBottom < aTop-T {
		return Precedes
	}

	return Unknown
}
