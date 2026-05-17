package functions

import (
	"math"
	"testing"
)

// Type4Tester is a fluent test helper for verifying Type 4 PostScript function execution.
// It wraps an ExecutionContext and provides assertion methods to pop values from the stack
// and compare them against expected results.
type Type4Tester struct {
	context *ExecutionContext
}

// Create creates a new Type4Tester by parsing the given PostScript text into an instruction
// sequence, executing it in a fresh ExecutionContext with default operators, and returning
// the tester for chained assertions.
func CreateType4Tester(text string) *Type4Tester {
	instructions, err := ParseInstructionSequence(text)
	if err != nil {
		panic(err)
	}

	context := NewExecutionContext(NewOperators())
	if err := instructions.Execute(context); err != nil {
		panic(err)
	}

	return &Type4Tester{context: context}
}

// PopBool pops a boolean value from the stack and asserts it equals the expected value.
func (t *Type4Tester) PopBool(tester *testing.T, expected bool) *Type4Tester {
	value := t.context.Pop()
	v, ok := value.(bool)
	if !ok {
		tester.Fatalf("expected bool on stack top, got %T", value)
	}
	if v != expected {
		tester.Errorf("expected bool %v, got %v", expected, v)
	}
	return t
}

// PopReal pops a numeric value from the stack and asserts it equals the expected float64
// within the default tolerance of 1e-7.
func (t *Type4Tester) PopReal(tester *testing.T, expected float64) *Type4Tester {
	return t.PopRealWithDelta(tester, expected, 1e-7)
}

// PopRealWithDelta pops a numeric value from the stack and asserts it equals the expected
// float64 within the given delta tolerance.
func (t *Type4Tester) PopRealWithDelta(tester *testing.T, expected, delta float64) *Type4Tester {
	value := t.context.Pop()
	v := toFloat64(value)
	if math.Abs(expected-v) >= delta {
		tester.Errorf("expected %g, got %g (delta=%g)", expected, v, delta)
	}
	return t
}

// PopInt pops an integer value from the stack and asserts it equals the expected int.
func (t *Type4Tester) PopInt(tester *testing.T, expected int) *Type4Tester {
	value := t.context.PopInt()
	if value != expected {
		tester.Errorf("expected %d, got %d", expected, value)
	}
	return t
}

// Pop pops a numeric value from the stack and asserts it equals the expected float64
// within the default tolerance of 1e-7. This is the general-purpose numeric pop, handling
// both integer and float stack values.
func (t *Type4Tester) Pop(tester *testing.T, expected float64) *Type4Tester {
	return t.PopWithDelta(tester, expected, 1e-7)
}

// PopWithDelta pops a numeric value from the stack and asserts it equals the expected
// float64 within the given delta tolerance. This is the general-purpose numeric pop,
// handling both integer and float stack values.
func (t *Type4Tester) PopWithDelta(tester *testing.T, expected, delta float64) *Type4Tester {
	value := t.context.Pop()
	v := toFloat64(value)
	if math.Abs(expected-v) >= delta {
		tester.Errorf("expected %g, got %g (delta=%g)", expected, v, delta)
	}
	return t
}

// IsEmpty asserts that the operand stack is empty.
func (t *Type4Tester) IsEmpty(tester *testing.T) *Type4Tester {
	if t.context.Count() != 0 {
		tester.Errorf("expected empty stack, got %d items", t.context.Count())
	}
	return t
}

// Context returns the underlying ExecutionContext for custom inspection.
func (t *Type4Tester) Context() *ExecutionContext {
	return t.context
}
