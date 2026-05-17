// Package functions provides tests for Type 4 PostScript function operators.
package functions

import (
	"fmt"
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	CreateType4Tester("5 6 add").Pop(t, 11).IsEmpty(t)
	CreateType4Tester("5 0.23 add").Pop(t, 5.23).IsEmpty(t)

	bigValue := int(math.MaxInt32) - 2
	text := fmt.Sprintf("%d %d add", bigValue, bigValue)
	ctx := CreateType4Tester(text)
	floatResult := toFloat64(ctx.Context().Pop())
	expected := float64(2*int64(math.MaxInt32) - 4)
	if math.Abs(expected-floatResult) >= 1 {
		t.Errorf("expected %g, got %g", expected, floatResult)
	}
	ctx.IsEmpty(t)
}

func TestAbs(t *testing.T) {
	CreateType4Tester("-3 abs 2.1 abs -2.1 abs -7.5 abs").
		Pop(t, 7.5).Pop(t, 2.1).Pop(t, 2.1).Pop(t, 3).IsEmpty(t)
}

func TestAnd(t *testing.T) {
	CreateType4Tester("true true and true false and").
		PopBool(t, false).PopBool(t, true).IsEmpty(t)
	CreateType4Tester("99 1 and 52 7 and").
		PopInt(t, 4).PopInt(t, 1).IsEmpty(t)
}

func TestAtan(t *testing.T) {
	CreateType4Tester("0 1 atan").Pop(t, 0.0).IsEmpty(t)
	CreateType4Tester("1 0 atan").Pop(t, 90.0).IsEmpty(t)
	CreateType4Tester("-100 0 atan").Pop(t, 270.0).IsEmpty(t)
	CreateType4Tester("4 4 atan").Pop(t, 45.0).IsEmpty(t)
}

func TestCeiling(t *testing.T) {
	CreateType4Tester("3.2 ceiling -4.8 ceiling 99 ceiling").
		Pop(t, 99.0).Pop(t, -4.0).Pop(t, 4.0).IsEmpty(t)
}

func TestCos(t *testing.T) {
	CreateType4Tester("0 cos").PopReal(t, 1).IsEmpty(t)
	CreateType4Tester("90 cos").PopReal(t, 0).IsEmpty(t)
}

func TestCvi(t *testing.T) {
	CreateType4Tester("-47.8 cvi").PopInt(t, -47).IsEmpty(t)
	CreateType4Tester("520.9 cvi").PopInt(t, 520).IsEmpty(t)
}

func TestCvr(t *testing.T) {
	CreateType4Tester("-47.8 cvr").PopReal(t, -47.8).IsEmpty(t)
	CreateType4Tester("520.9 cvr").PopReal(t, 520.9).IsEmpty(t)
	CreateType4Tester("77 cvr").PopReal(t, 77).IsEmpty(t)

	ctx := CreateType4Tester("77 77 cvr")
	top := ctx.Context().Pop()
	if _, ok := top.(float64); !ok {
		t.Error("expected a float64 as the result of 'cvr'")
	}
	bottom := ctx.Context().Pop()
	if _, ok := bottom.(int); !ok {
		t.Error("expected an int from an int literal")
	}
	ctx.IsEmpty(t)
}

func TestDiv(t *testing.T) {
	CreateType4Tester("3 2 div").PopReal(t, 1.5).IsEmpty(t)
	CreateType4Tester("4 2 div").PopReal(t, 2.0).IsEmpty(t)
}

func TestExp(t *testing.T) {
	CreateType4Tester("9 0.5 exp").PopRealWithDelta(t, 3.0, 1e-7).IsEmpty(t)
	CreateType4Tester("-9 -1 exp").PopRealWithDelta(t, -0.111111, 0.000001).IsEmpty(t)
}

func TestFloor(t *testing.T) {
	CreateType4Tester("3.2 floor -4.8 floor 99 floor").
		Pop(t, 99.0).Pop(t, -5.0).Pop(t, 3.0).IsEmpty(t)
}

func TestIDiv(t *testing.T) {
	CreateType4Tester("3 2 idiv").PopInt(t, 1).IsEmpty(t)
	CreateType4Tester("4 2 idiv").PopInt(t, 2).IsEmpty(t)
	CreateType4Tester("-5 2 idiv").PopInt(t, -2).IsEmpty(t)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for float operand in idiv")
		}
	}()
	CreateType4Tester("4.4 2 idiv")
}

func TestLn(t *testing.T) {
	CreateType4Tester("10 ln").PopRealWithDelta(t, 2.30259, 0.00001).IsEmpty(t)
	CreateType4Tester("100 ln").PopRealWithDelta(t, 4.60517, 0.00001).IsEmpty(t)
}

func TestLog(t *testing.T) {
	CreateType4Tester("10 log").PopReal(t, 1.0).IsEmpty(t)
	CreateType4Tester("100 log").PopReal(t, 2.0).IsEmpty(t)
}

func TestMod(t *testing.T) {
	CreateType4Tester("5 3 mod").PopInt(t, 2).IsEmpty(t)
	CreateType4Tester("5 2 mod").PopInt(t, 1).IsEmpty(t)
	CreateType4Tester("-5 3 mod").PopInt(t, -2).IsEmpty(t)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for float operand in mod")
		}
	}()
	CreateType4Tester("4.4 2 mod")
}

func TestMul(t *testing.T) {
	CreateType4Tester("1 2 mul").PopInt(t, 2).IsEmpty(t)
	CreateType4Tester("1.5 2 mul").PopReal(t, 3.0).IsEmpty(t)
	CreateType4Tester("1.5 2.1 mul").PopRealWithDelta(t, 3.15, 0.001).IsEmpty(t)
	CreateType4Tester("2147483644 2 mul"). // int overflow -> real
		PopRealWithDelta(t, float64(2*int64(math.MaxInt32-3)), 0.001).IsEmpty(t)
}

func TestNeg(t *testing.T) {
	CreateType4Tester("4.5 neg").PopReal(t, -4.5).IsEmpty(t)
	CreateType4Tester("-3 neg").PopInt(t, 3).IsEmpty(t)

	// Border cases: C# uses (int.MinValue + 1) = -2147483647 → neg gives int.MaxValue
	CreateType4Tester("-2147483647 neg").PopInt(t, math.MaxInt32).IsEmpty(t)
	CreateType4Tester("-2147483648 neg").PopRealWithDelta(t, -float64(math.MinInt32), 0.001).IsEmpty(t)
}

func TestRound(t *testing.T) {
	CreateType4Tester("3.2 round").PopReal(t, 3.0).IsEmpty(t)
	CreateType4Tester("6.5 round").PopReal(t, 7.0).IsEmpty(t)
	CreateType4Tester("-4.8 round").PopReal(t, -5.0).IsEmpty(t)
	CreateType4Tester("-6.5 round").PopReal(t, -6.0).IsEmpty(t)
	CreateType4Tester("99 round").PopInt(t, 99).IsEmpty(t)
}

func TestSin(t *testing.T) {
	CreateType4Tester("0 sin").PopReal(t, 0).IsEmpty(t)
	CreateType4Tester("90 sin").PopReal(t, 1).IsEmpty(t)
	CreateType4Tester("-90.0 sin").PopReal(t, -1).IsEmpty(t)
}

func TestSqrt(t *testing.T) {
	CreateType4Tester("0 sqrt").PopReal(t, 0).IsEmpty(t)
	CreateType4Tester("1 sqrt").PopReal(t, 1).IsEmpty(t)
	CreateType4Tester("4 sqrt").PopReal(t, 2).IsEmpty(t)
	CreateType4Tester("4.4 sqrt").PopRealWithDelta(t, 2.097617, 0.000001).IsEmpty(t)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for negative sqrt")
		}
	}()
	CreateType4Tester("-4.1 sqrt")
}

func TestSub(t *testing.T) {
	CreateType4Tester("5 2 sub -7.5 1 sub").
		Pop(t, -8.5).Pop(t, 3).IsEmpty(t)
}

func TestTruncate(t *testing.T) {
	CreateType4Tester("3.2 truncate").PopReal(t, 3.0).IsEmpty(t)
	CreateType4Tester("-4.8 truncate").PopReal(t, -4.0).IsEmpty(t)
	CreateType4Tester("99 truncate").PopInt(t, 99).IsEmpty(t)
}

func TestBitshift(t *testing.T) {
	CreateType4Tester("7 3 bitshift 142 -3 bitshift").
		PopInt(t, 17).PopInt(t, 56).IsEmpty(t)
}

func TestEq(t *testing.T) {
	// Stack: [true, false, false, true, false, true] → pop LIFO: true, false, true, false, false, true
	CreateType4Tester("7 7 eq 7 6 eq 7 -7 eq true true eq false true eq 7.7 7.7 eq").
		PopBool(t, true).PopBool(t, false).PopBool(t, true).PopBool(t, false).PopBool(t, false).PopBool(t, true).IsEmpty(t)
}

func TestGe(t *testing.T) {
	CreateType4Tester("5 7 ge 7 5 ge 7 7 ge -1 2 ge").
		PopBool(t, false).PopBool(t, true).PopBool(t, true).PopBool(t, false).IsEmpty(t)
}

func TestGt(t *testing.T) {
	// Stack after ops: [false, true, false, false] → pop order reversed
	CreateType4Tester("5 7 gt 7 5 gt 7 7 gt -1 2 gt").
		PopBool(t, false).PopBool(t, false).PopBool(t, true).PopBool(t, false).IsEmpty(t)
}

func TestLe(t *testing.T) {
	// Stack: [true, false, true, true] → pop order reversed
	CreateType4Tester("5 7 le 7 5 le 7 7 le -1 2 le").
		PopBool(t, true).PopBool(t, true).PopBool(t, false).PopBool(t, true).IsEmpty(t)
}

func TestLt(t *testing.T) {
	CreateType4Tester("5 7 lt 7 5 lt 7 7 lt -1 2 lt").
		PopBool(t, true).PopBool(t, false).PopBool(t, false).PopBool(t, true).IsEmpty(t)
}

func TestNe(t *testing.T) {
	// Stack: [false, true, true, false, true, false] → pop LIFO: false, true, false, true, true, false
	CreateType4Tester("7 7 ne 7 6 ne 7 -7 ne true true ne false true ne 7.7 7.7 ne").
		PopBool(t, false).PopBool(t, true).PopBool(t, false).PopBool(t, true).PopBool(t, true).PopBool(t, false).IsEmpty(t)
}

func TestNot(t *testing.T) {
	CreateType4Tester("true not false not").
		PopBool(t, true).PopBool(t, false).IsEmpty(t)
	CreateType4Tester("52 not -37 not").
		PopInt(t, 37).PopInt(t, -52).IsEmpty(t)
}

func TestOr(t *testing.T) {
	CreateType4Tester("true true or true false or false false or").
		PopBool(t, false).PopBool(t, true).PopBool(t, true).IsEmpty(t)
	CreateType4Tester("17 5 or 1 1 or").
		PopInt(t, 1).PopInt(t, 21).IsEmpty(t)
}

func TestXor(t *testing.T) {
	CreateType4Tester("true true xor true false xor false false xor").
		PopBool(t, false).PopBool(t, true).PopBool(t, false).IsEmpty(t)
	CreateType4Tester("7 3 xor 12 3 or").
		PopInt(t, 15).PopInt(t, 4).IsEmpty(t)
}

func TestIf(t *testing.T) {
	CreateType4Tester("true { 2 1 add } if").
		PopInt(t, 3).IsEmpty(t)
	CreateType4Tester("false { 2 1 add } if").
		IsEmpty(t)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-boolean condition in if")
		}
	}()
	CreateType4Tester("0 { 2 1 add } if")
}

func TestIfElse(t *testing.T) {
	CreateType4Tester("true { 2 1 add } { 2 1 sub } ifelse").
		PopInt(t, 3).IsEmpty(t)
	CreateType4Tester("false { 2 1 add } { 2 1 sub } ifelse").
		PopInt(t, 1).IsEmpty(t)
}

func TestCopy(t *testing.T) {
	CreateType4Tester("true 1 2 3 3 copy").
		PopInt(t, 3).PopInt(t, 2).PopInt(t, 1).
		PopInt(t, 3).PopInt(t, 2).PopInt(t, 1).
		PopBool(t, true).IsEmpty(t)
}

func TestDup(t *testing.T) {
	CreateType4Tester("true 1 2 dup").
		PopInt(t, 2).PopInt(t, 2).PopInt(t, 1).
		PopBool(t, true).IsEmpty(t)
	CreateType4Tester("true dup").
		PopBool(t, true).PopBool(t, true).IsEmpty(t)
}

func TestExch(t *testing.T) {
	CreateType4Tester("true 1 exch").
		PopBool(t, true).PopInt(t, 1).IsEmpty(t)
	CreateType4Tester("1 2.5 exch").
		PopInt(t, 1).Pop(t, 2.5).IsEmpty(t)
}

func TestIndex(t *testing.T) {
	CreateType4Tester("1 2 3 4 0 index").
		PopInt(t, 4).PopInt(t, 4).PopInt(t, 3).PopInt(t, 2).PopInt(t, 1).IsEmpty(t)
	CreateType4Tester("1 2 3 4 3 index").
		PopInt(t, 1).PopInt(t, 4).PopInt(t, 3).PopInt(t, 2).PopInt(t, 1).IsEmpty(t)
}

func TestStackPop(t *testing.T) {
	CreateType4Tester("1 pop 7 2 pop").
		PopInt(t, 7).IsEmpty(t)
	CreateType4Tester("1 2 3 pop pop").
		PopInt(t, 1).IsEmpty(t)
}

func TestRoll(t *testing.T) {
	// "5 -2 roll" on [1,2,3,4,5] → stack becomes [2,1,5,4,3] → pop: 3,4,5,1,2
	CreateType4Tester("1 2 3 4 5 5 -2 roll").
		PopInt(t, 3).PopInt(t, 4).PopInt(t, 5).PopInt(t, 1).PopInt(t, 2).IsEmpty(t)
	// "5 2 roll" on [1,2,3,4,5] → stack becomes [3,2,1,5,4] → pop: 4,5,1,2,3
	CreateType4Tester("1 2 3 4 5 5 2 roll").
		PopInt(t, 4).PopInt(t, 5).PopInt(t, 1).PopInt(t, 2).PopInt(t, 3).IsEmpty(t)
	// "3 0 roll" is a no-op
	CreateType4Tester("1 2 3 3 0 roll").
		PopInt(t, 3).PopInt(t, 2).PopInt(t, 1).IsEmpty(t)
}
