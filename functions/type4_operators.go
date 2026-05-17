package functions

import (
	"fmt"
	"math"
)

// --- Arithmetic Operators ---

type arithmeticAbsOp struct{}

func (arithmeticAbsOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	switch v := num.(type) {
	case int:
		if v < 0 {
			v = -v
		}
		c.Push(v)
	default:
		f := toFloat64(num)
		c.Push(math.Abs(f))
	}
}

var arithmeticAbs Operator = arithmeticAbsOp{}

type arithmeticAddOp struct{}

func (arithmeticAddOp) Execute(c *ExecutionContext) {
	num2 := c.Pop()
	num1 := c.Pop()
	if i1, ok1 := num1.(int); ok1 {
		if i2, ok2 := num2.(int); ok2 {
			sum := int64(i1) + int64(i2)
			if sum < math.MinInt32 || sum > math.MaxInt32 {
				c.Push(float64(sum))
			} else {
				c.Push(int(sum))
			}
			return
		}
	}
	c.Push(toFloat64(num1) + toFloat64(num2))
}

var arithmeticAdd Operator = arithmeticAddOp{}

type arithmeticAtanOp struct{}

func (arithmeticAtanOp) Execute(c *ExecutionContext) {
	den := c.PopReal()
	num := c.PopReal()
	atan := math.Atan2(num, den)
	atan = math.Mod(toDegrees(atan), 360)
	if atan < 0 {
		atan += 360
	}
	c.Push(atan)
}

var arithmeticAtan Operator = arithmeticAtanOp{}

type arithmeticCeilingOp struct{}

func (arithmeticCeilingOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	if v, ok := num.(int); ok {
		c.Push(v)
		return
	}
	c.Push(math.Ceil(toFloat64(num)))
}

var arithmeticCeiling Operator = arithmeticCeilingOp{}

type arithmeticCosOp struct{}

func (arithmeticCosOp) Execute(c *ExecutionContext) {
	angle := c.PopReal()
	c.Push(math.Cos(toRadians(angle)))
}

var arithmeticCos Operator = arithmeticCosOp{}

type arithmeticCviOp struct{}

func (arithmeticCviOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	c.Push(int(truncate(toFloat64(num))))
}

var arithmeticCvi Operator = arithmeticCviOp{}

type arithmeticCvrOp struct{}

func (arithmeticCvrOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	c.Push(toFloat64(num))
}

var arithmeticCvr Operator = arithmeticCvrOp{}

type arithmeticDivOp struct{}

func (arithmeticDivOp) Execute(c *ExecutionContext) {
	num2 := toFloat64(c.Pop())
	num1 := toFloat64(c.Pop())
	c.Push(num1 / num2)
}

var arithmeticDiv Operator = arithmeticDivOp{}

type arithmeticExpOp struct{}

func (arithmeticExpOp) Execute(c *ExecutionContext) {
	exp := toFloat64(c.Pop())
	base := toFloat64(c.Pop())
	c.Push(math.Pow(base, exp))
}

var arithmeticExp Operator = arithmeticExpOp{}

type arithmeticFloorOp struct{}

func (arithmeticFloorOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	if v, ok := num.(int); ok {
		c.Push(v)
		return
	}
	c.Push(math.Floor(toFloat64(num)))
}

var arithmeticFloor Operator = arithmeticFloorOp{}

type arithmeticIDivOp struct{}

func (arithmeticIDivOp) Execute(c *ExecutionContext) {
	num2 := c.PopInt()
	num1 := c.PopInt()
	c.Push(num1 / num2)
}

var arithmeticIDiv Operator = arithmeticIDivOp{}

type arithmeticLnOp struct{}

func (arithmeticLnOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	c.Push(math.Log(toFloat64(num)))
}

var arithmeticLn Operator = arithmeticLnOp{}

type arithmeticLogOp struct{}

func (arithmeticLogOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	c.Push(math.Log10(toFloat64(num)))
}

var arithmeticLog Operator = arithmeticLogOp{}

type arithmeticModOp struct{}

func (arithmeticModOp) Execute(c *ExecutionContext) {
	int2 := c.PopInt()
	int1 := c.PopInt()
	c.Push(int1 % int2)
}

var arithmeticMod Operator = arithmeticModOp{}

type arithmeticMulOp struct{}

func (arithmeticMulOp) Execute(c *ExecutionContext) {
	num2 := c.Pop()
	num1 := c.Pop()
	if i1, ok1 := num1.(int); ok1 {
		if i2, ok2 := num2.(int); ok2 {
			result := int64(i1) * int64(i2)
			if result >= math.MinInt32 && result <= math.MaxInt32 {
				c.Push(int(result))
			} else {
				c.Push(float64(result))
			}
			return
		}
	}
	c.Push(toFloat64(num1) * toFloat64(num2))
}

var arithmeticMul Operator = arithmeticMulOp{}

type arithmeticNegOp struct{}

func (arithmeticNegOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	if v, ok := num.(int); ok {
		if v == math.MinInt32 {
			c.Push(-toFloat64(v))
		} else {
			c.Push(-v)
		}
		return
	}
	c.Push(-toFloat64(num))
}

var arithmeticNeg Operator = arithmeticNegOp{}

type arithmeticRoundOp struct{}

func (arithmeticRoundOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	if v, ok := num.(int); ok {
		c.Push(v)
		return
	}
	value := toFloat64(num)
	// PostScript round: half-values rounded toward positive infinity (floor(x + 0.5))
	c.Push(math.Floor(value + 0.5))
}

var arithmeticRound Operator = arithmeticRoundOp{}

type arithmeticSinOp struct{}

func (arithmeticSinOp) Execute(c *ExecutionContext) {
	angle := c.PopReal()
	c.Push(math.Sin(toRadians(angle)))
}

var arithmeticSin Operator = arithmeticSinOp{}

type arithmeticSqrtOp struct{}

func (arithmeticSqrtOp) Execute(c *ExecutionContext) {
	num := c.PopReal()
	if num < 0 {
		panic(fmt.Errorf("argument must be nonnegative"))
	}
	c.Push(math.Sqrt(num))
}

var arithmeticSqrt Operator = arithmeticSqrtOp{}

type arithmeticSubOp struct{}

func (arithmeticSubOp) Execute(c *ExecutionContext) {
	num2 := c.Pop()
	num1 := c.Pop()
	if i1, ok1 := num1.(int); ok1 {
		if i2, ok2 := num2.(int); ok2 {
			result := int64(i1) - int64(i2)
			if result < math.MinInt32 || result > math.MaxInt32 {
				c.Push(float64(result))
			} else {
				c.Push(int(result))
			}
			return
		}
	}
	c.Push(toFloat64(num1) - toFloat64(num2))
}

var arithmeticSub Operator = arithmeticSubOp{}

type arithmeticTruncateOp struct{}

func (arithmeticTruncateOp) Execute(c *ExecutionContext) {
	num := c.Pop()
	if v, ok := num.(int); ok {
		c.Push(v)
		return
	}
	c.Push(truncate(toFloat64(num)))
}

var arithmeticTruncate Operator = arithmeticTruncateOp{}

// --- Bitwise / Boolean Operators ---

type bitwiseAndOp struct{}

func (bitwiseAndOp) Execute(c *ExecutionContext) {
	op2 := c.Pop()
	op1 := c.Pop()
	if b1, ok1 := op1.(bool); ok1 {
		if b2, ok2 := op2.(bool); ok2 {
			c.Push(b1 && b2)
			return
		}
	}
	if i1, ok1 := op1.(int); ok1 {
		if i2, ok2 := op2.(int); ok2 {
			c.Push(i1 & i2)
			return
		}
	}
	panic(fmt.Errorf("operands must be bool/bool or int/int"))
}

var bitwiseAnd Operator = bitwiseAndOp{}

type bitwiseBitshiftOp struct{}

func (bitwiseBitshiftOp) Execute(c *ExecutionContext) {
	shift := toInt(c.Pop())
	int1 := toInt(c.Pop())
	if shift < 0 {
		c.Push(int1 >> uint(math.Abs(float64(shift))))
	} else {
		c.Push(int1 << uint(shift))
	}
}

var bitwiseBitshift Operator = bitwiseBitshiftOp{}

type bitwiseFalseOp struct{}

func (bitwiseFalseOp) Execute(c *ExecutionContext) {
	c.Push(false)
}

var bitwiseFalse Operator = bitwiseFalseOp{}

type bitwiseNotOp struct{}

func (bitwiseNotOp) Execute(c *ExecutionContext) {
	op1 := c.Pop()
	if b, ok := op1.(bool); ok {
		c.Push(!b)
	} else if i, ok := op1.(int); ok {
		c.Push(-i)
	} else {
		panic(fmt.Errorf("operand must be bool or int"))
	}
}

var bitwiseNot Operator = bitwiseNotOp{}

type bitwiseOrOp struct{}

func (bitwiseOrOp) Execute(c *ExecutionContext) {
	op2 := c.Pop()
	op1 := c.Pop()
	if b1, ok1 := op1.(bool); ok1 {
		if b2, ok2 := op2.(bool); ok2 {
			c.Push(b1 || b2)
			return
		}
	}
	if i1, ok1 := op1.(int); ok1 {
		if i2, ok2 := op2.(int); ok2 {
			c.Push(i1 | i2)
			return
		}
	}
	panic(fmt.Errorf("operands must be bool/bool or int/int"))
}

var bitwiseOr Operator = bitwiseOrOp{}

type bitwiseTrueOp struct{}

func (bitwiseTrueOp) Execute(c *ExecutionContext) {
	c.Push(true)
}

var bitwiseTrue Operator = bitwiseTrueOp{}

type bitwiseXorOp struct{}

func (bitwiseXorOp) Execute(c *ExecutionContext) {
	op2 := c.Pop()
	op1 := c.Pop()
	if b1, ok1 := op1.(bool); ok1 {
		if b2, ok2 := op2.(bool); ok2 {
			c.Push(b1 != b2)
			return
		}
	}
	if i1, ok1 := op1.(int); ok1 {
		if i2, ok2 := op2.(int); ok2 {
			c.Push(i1 ^ i2)
			return
		}
	}
	panic(fmt.Errorf("operands must be bool/bool or int/int"))
}

var bitwiseXor Operator = bitwiseXorOp{}

// truncate returns the integer part of x, rounding towards zero (matching C# Math.Truncate).
func truncate(x float64) float64 {
	if x >= 0 {
		return math.Floor(x)
	}
	return math.Ceil(x)
}

// --- Relational Operators ---

type relationalEqOp struct{ negated bool }

func (r relationalEqOp) Execute(c *ExecutionContext) {
	op2 := c.Pop()
	op1 := c.Pop()
	result := isEqual(op1, op2)
	if r.negated {
		result = !result
	}
	c.Push(result)
}

var relationalEq Operator = relationalEqOp{negated: false}
var relationalNe Operator = relationalEqOp{negated: true}

func isEqual(op1, op2 any) bool {
	if f1, ok1 := op1.(float64); ok1 {
		if f2, ok2 := op2.(float64); ok2 {
			return f1 == f2
		}
	}
	return op1 == op2
}

type numberComparisonOp struct {
	compare func(float64, float64) bool
}

func (n numberComparisonOp) Execute(c *ExecutionContext) {
	op2 := c.Pop()
	op1 := c.Pop()
	num1 := toFloat64(op1)
	num2 := toFloat64(op2)
	c.Push(n.compare(num1, num2))
}

var relationalGe Operator = numberComparisonOp{compare: func(a, b float64) bool { return a >= b }}
var relationalGt Operator = numberComparisonOp{compare: func(a, b float64) bool { return a > b }}
var relationalLe Operator = numberComparisonOp{compare: func(a, b float64) bool { return a <= b }}
var relationalLt Operator = numberComparisonOp{compare: func(a, b float64) bool { return a < b }}

// --- Conditional Operators ---

type conditionalIfOp struct{}

func (conditionalIfOp) Execute(c *ExecutionContext) {
	proc := c.Pop().(*InstructionSequence)
	condition := c.Pop().(bool)
	if condition {
		proc.Execute(c) // ignore error; C# also swallows here implicitly via void return
	}
}

var conditionalIf Operator = conditionalIfOp{}

type conditionalIfElseOp struct{}

func (conditionalIfElseOp) Execute(c *ExecutionContext) {
	proc2 := c.Pop().(*InstructionSequence)
	proc1 := c.Pop().(*InstructionSequence)
	condition := c.Pop().(bool)
	if condition {
		proc1.Execute(c)
	} else {
		proc2.Execute(c)
	}
}

var conditionalIfElse Operator = conditionalIfElseOp{}

// --- Stack Operators ---

type stackCopyOp struct{}

func (stackCopyOp) Execute(c *ExecutionContext) {
	n := toInt(c.Pop())
	if n > 0 && len(c.stack) >= n {
		topN := make([]any, n)
		copy(topN, c.stack[len(c.stack)-n:])
		below := c.stack[:len(c.stack)-n]
		newStack := make([]any, 0, len(below)+2*n)
		newStack = append(newStack, below...)
		newStack = append(newStack, topN...)
		newStack = append(newStack, topN...)
		c.stack = newStack
	}
}

var stackCopy Operator = stackCopyOp{}

type stackDupOp struct{}

func (stackDupOp) Execute(c *ExecutionContext) {
	c.Push(c.Peek())
}

var stackDup Operator = stackDupOp{}

type stackExchOp struct{}

func (stackExchOp) Execute(c *ExecutionContext) {
	any2 := c.Pop()
	any1 := c.Pop()
	c.Push(any2)
	c.Push(any1)
}

var stackExch Operator = stackExchOp{}

type stackIndexOp struct{}

func (stackIndexOp) Execute(c *ExecutionContext) {
	n := toInt(c.Pop())
	if n < 0 || n >= len(c.stack) {
		panic(fmt.Errorf("rangecheck: %d", n))
	}
	c.Push(c.stack[len(c.stack)-1-n])
}

var stackIndex Operator = stackIndexOp{}

type stackPopOp struct{}

func (stackPopOp) Execute(c *ExecutionContext) {
	c.Pop()
}

var stackPop Operator = stackPopOp{}

type stackRollOp struct{}

func (stackRollOp) Execute(c *ExecutionContext) {
	j := toInt(c.Pop())
	n := toInt(c.Pop())
	if j == 0 {
		return
	}
	if n < 0 {
		panic(fmt.Errorf("rangecheck: %d", n))
	}

	rolled := make([]any, 0)
	moved := make([]any, 0)

	if j < 0 {
		n1 := n + j
		for i := 0; i < n1; i++ {
			moved = append(moved, c.Pop())
		}
		for i := j; i < 0; i++ {
			rolled = append(rolled, c.Pop())
		}
		c.AddAllToStack(moved)
		c.AddAllToStack(rolled)
	} else {
		n1 := n - j
		for i := j; i > 0; i-- {
			rolled = append(rolled, c.Pop())
		}
		for i := 0; i < n1; i++ {
			moved = append(moved, c.Pop())
		}
		c.AddAllToStack(rolled)
		c.AddAllToStack(moved)
	}
}

var stackRoll Operator = stackRollOp{}

// --- Helper functions ---

func toRadians(val float64) float64 {
	return (math.Pi / 180.0) * val
}

func toDegrees(val float64) float64 {
	return (180.0 / math.Pi) * val
}
