package geometry

import "fmt"

// UserSpaceUnit represents a user space unit in a PDF page.
// By default user space units correspond to 1/72nd of an inch (a typographic point).
// The UserUnit entry in a page dictionary can define the space units as a different
// multiple of 1/72 (1 point).
type UserSpaceUnit struct {
	// PointMultiples is the number of points (1/72nd of an inch) corresponding to
	// a single unit in user space.
	PointMultiples int
}

// Default is the default user space unit with PointMultiples set to 1.
var Default = UserSpaceUnit{PointMultiples: 1}

// NewUserSpaceUnit creates a new user space unit specification for a page.
func NewUserSpaceUnit(pointMultiples int) (UserSpaceUnit, error) {
	if pointMultiples <= 0 {
		return UserSpaceUnit{}, fmt.Errorf("cannot have a zero or negative value of point multiples: %d", pointMultiples)
	}

	return UserSpaceUnit{PointMultiples: pointMultiples}, nil
}

// String returns the string representation of the user space unit.
func (u UserSpaceUnit) String() string {
	return fmt.Sprintf("%d", u.PointMultiples)
}
