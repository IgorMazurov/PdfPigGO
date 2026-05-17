// Package graphics provides tests for Push (q) operation.
package graphics

import "testing"

func TestPushSymbolCorrect(t *testing.T) {
	if got := InstancePush.Operator(); got != "q" {
		t.Errorf("expected symbol 'q', got '%s'", got)
	}
}

func TestPushAddsToStack(t *testing.T) {
	ctx := NewTestOperationContext()

	InstancePush.Run(ctx)

	if got := ctx.StackSize(); got != 2 {
		t.Errorf("expected stack size 2 after push, got %d", got)
	}
}
