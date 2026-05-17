package fonts

import (
	"errors"
	"fmt"
	"strings"
)

// CharStringStack is the stack of numeric operands currently active in a CharString.
type CharStringStack struct {
	stack []float64
}

// NewCharStringStack creates an empty CharStringStack.
func NewCharStringStack() *CharStringStack {
	return &CharStringStack{}
}

var errEmptyStack = errors.New("cannot pop from an empty stack, invalid charstring parsed")

// Length returns the current size of the stack.
func (s *CharStringStack) Length() int {
	return len(s.stack)
}

// CanPop reports whether it is possible to pop a value from either end of the stack.
func (s *CharStringStack) CanPop() bool {
	return len(s.stack) > 0
}

// PopTop removes and returns the value from the top of the stack.
func (s *CharStringStack) PopTop() (float64, error) {
	if len(s.stack) == 0 {
		return 0, errEmptyStack
	}
	result := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return result, nil
}

// PopBottom removes and returns the value from the bottom of the stack.
func (s *CharStringStack) PopBottom() (float64, error) {
	if len(s.stack) == 0 {
		return 0, errEmptyStack
	}
	result := s.stack[0]
	s.stack = s.stack[1:]
	return result, nil
}

// Push adds the value to the top of the stack.
func (s *CharStringStack) Push(value float64) {
	s.stack = append(s.stack, value)
}

// CopyElementAt copies the element at the given index. Negative indices are
// resolved relative to the end of the stack (-1 is the top element).
func (s *CharStringStack) CopyElementAt(index int) float64 {
	if index < 0 {
		return s.stack[len(s.stack)+index]
	}
	return s.stack[index]
}

// Clear removes all values from the stack.
func (s *CharStringStack) Clear() {
	s.stack = s.stack[:0]
}

// String returns a space-separated string representation of the stack contents.
func (s *CharStringStack) String() string {
	parts := make([]string, len(s.stack))
	for i, v := range s.stack {
		parts[i] = fmt.Sprintf("%g", v)
	}
	return strings.Join(parts, " ")
}
