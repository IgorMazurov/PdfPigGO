// Package graphics provides tests for Pop (Q) operation.
package graphics

import "testing"

func TestPopSymbolCorrect(t *testing.T) {
	if got := InstancePop.Operator(); got != "Q" {
		t.Errorf("expected symbol 'Q', got '%s'", got)
	}
}

func TestCannotPopWithNoFrames(t *testing.T) {
	ctx := NewTestOperationContext()

	// Remove the initial frame so stack is empty.
	ctx.StateStack = nil

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when popping with no frames")
		}
	}()
	InstancePop.Run(ctx)
}

func TestPopsTopFrame(t *testing.T) {
	ctx := NewTestOperationContext()

	ctx.StateStack = append(ctx.StateStack, &CurrentGraphicsState{
		LineWidth: 23,
	})

	InstancePop.Run(ctx)

	if got := ctx.StackSize(); got != 1 {
		t.Errorf("expected stack size 1 after pop, got %d", got)
	}
	if got := ctx.GetCurrentState().LineWidth; got != 1 {
		t.Errorf("expected line width 1 (initial state default per PDF spec), got %f", got)
	}
}
