// Package graphics provides tests for SetMiterLimit operation.
package graphics

import "testing"

func TestSetMiterLimitRunSetsMiterLimitOfCurrentState(t *testing.T) {
	ctx := NewTestOperationContext()

	limit := NewSetMiterLimit(25)
	limit.Run(ctx)

	if got := ctx.GetCurrentState().MiterLimit; got != 25 {
		t.Errorf("expected miter limit 25, got %f", got)
	}
}

func TestSetMiterLimitSymbolCorrect(t *testing.T) {
	limit := NewSetMiterLimit(10)
	if got := limit.Operator(); got != "M" {
		t.Errorf("expected symbol 'M', got '%s'", got)
	}
}
