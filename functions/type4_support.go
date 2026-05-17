package functions

import (
	"fmt"
	"strconv"
)

// ExecutionContext holds the state for executing a Type 4 PostScript function.
// It provides access to an operand stack and the operator set.
type ExecutionContext struct {
	stack     []any
	operators *Operators
}

// NewExecutionContext creates a new execution context with the given operator set.
func NewExecutionContext(operatorSet *Operators) *ExecutionContext {
	return &ExecutionContext{
		stack:     make([]any, 0),
		operators: operatorSet,
	}
}

// Stack returns the operand stack as a slice for direct access.
// The top of the stack is the last element in the slice.
func (c *ExecutionContext) Stack() []any {
	return c.stack
}

// Push pushes a value onto the operand stack.
func (c *ExecutionContext) Push(v any) {
	c.stack = append(c.stack, v)
}

// Pop removes and returns the top value from the operand stack.
func (c *ExecutionContext) Pop() any {
	if len(c.stack) == 0 {
		panic(fmt.Errorf("stack underflow"))
	}
	v := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	return v
}

// Peek returns the top value from the operand stack without removing it.
func (c *ExecutionContext) Peek() any {
	if len(c.stack) == 0 {
		panic(fmt.Errorf("stack underflow"))
	}
	return c.stack[len(c.stack)-1]
}

// Count returns the number of items on the operand stack.
func (c *ExecutionContext) Count() int {
	return len(c.stack)
}

// GetOperators returns the operator set used by this execution context.
func (c *ExecutionContext) GetOperators() *Operators {
	return c.operators
}

// PopReal pops a number from the stack and returns it as a float64.
func (c *ExecutionContext) PopReal() float64 {
	v := c.Pop()
	return toFloat64(v)
}

// PopInt pops an integer value from the stack.
// Panics if the value is not an integer type.
func (c *ExecutionContext) PopInt() int {
	v := c.Pop()
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	default:
		panic("PopInt cannot be done as the value is not integer")
	}
}

// AddAllToStack pushes all values from the given slice onto the stack,
// followed by the current stack contents (preserving order).
func (c *ExecutionContext) AddAllToStack(values []any) {
	newStack := make([]any, 0, len(values)+len(c.stack))
	newStack = append(newStack, values...)
	newStack = append(newStack, c.stack...)
	c.stack = newStack
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	case float32:
		return float64(val)
	default:
		return 0
	}
}

func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	case float32:
		return int(val)
	default:
		return 0
	}
}

// Operator defines a PostScript operator that can be executed within an ExecutionContext.
type Operator interface {
	Execute(context *ExecutionContext)
}

// Operators provides a registry of all supported PostScript operators for Type 4 functions.
type Operators struct {
	operatorMap map[string]Operator
}

// NewOperators creates a new Operators instance with the default set of operators.
func NewOperators() *Operators {
	op := &Operators{
		operatorMap: make(map[string]Operator),
	}

	// Arithmetic operators
	op.operatorMap["add"] = arithmeticAdd
	op.operatorMap["abs"] = arithmeticAbs
	op.operatorMap["atan"] = arithmeticAtan
	op.operatorMap["ceiling"] = arithmeticCeiling
	op.operatorMap["cos"] = arithmeticCos
	op.operatorMap["cvi"] = arithmeticCvi
	op.operatorMap["cvr"] = arithmeticCvr
	op.operatorMap["div"] = arithmeticDiv
	op.operatorMap["exp"] = arithmeticExp
	op.operatorMap["floor"] = arithmeticFloor
	op.operatorMap["idiv"] = arithmeticIDiv
	op.operatorMap["ln"] = arithmeticLn
	op.operatorMap["log"] = arithmeticLog
	op.operatorMap["mod"] = arithmeticMod
	op.operatorMap["mul"] = arithmeticMul
	op.operatorMap["neg"] = arithmeticNeg
	op.operatorMap["round"] = arithmeticRound
	op.operatorMap["sin"] = arithmeticSin
	op.operatorMap["sqrt"] = arithmeticSqrt
	op.operatorMap["sub"] = arithmeticSub
	op.operatorMap["truncate"] = arithmeticTruncate

	// Relational, boolean and bitwise operators
	op.operatorMap["and"] = bitwiseAnd
	op.operatorMap["bitshift"] = bitwiseBitshift
	op.operatorMap["eq"] = relationalEq
	op.operatorMap["false"] = bitwiseFalse
	op.operatorMap["ge"] = relationalGe
	op.operatorMap["gt"] = relationalGt
	op.operatorMap["le"] = relationalLe
	op.operatorMap["lt"] = relationalLt
	op.operatorMap["ne"] = relationalNe
	op.operatorMap["not"] = bitwiseNot
	op.operatorMap["or"] = bitwiseOr
	op.operatorMap["true"] = bitwiseTrue
	op.operatorMap["xor"] = bitwiseXor

	// Conditional operators
	op.operatorMap["if"] = conditionalIf
	op.operatorMap["ifelse"] = conditionalIfElse

	// Stack operators
	op.operatorMap["copy"] = stackCopy
	op.operatorMap["dup"] = stackDup
	op.operatorMap["exch"] = stackExch
	op.operatorMap["index"] = stackIndex
	op.operatorMap["pop"] = stackPop
	op.operatorMap["roll"] = stackRoll

	return op
}

// GetOperator returns the operator for the given name, or nil if not found.
func (o *Operators) GetOperator(name string) Operator {
	if op, ok := o.operatorMap[name]; ok {
		return op
	}
	return nil
}

// InstructionSequence holds a sequence of parsed instructions for Type 4 function execution.
// Each instruction is either an operator name (string), a numeric literal (int or float64),
// a boolean, or a nested proc (another InstructionSequence).
type InstructionSequence struct {
	instructions []any
}

// NewInstructionSequence creates an empty instruction sequence.
func NewInstructionSequence() *InstructionSequence {
	return &InstructionSequence{
		instructions: make([]any, 0),
	}
}

// AddName adds an operator name to the sequence.
func (s *InstructionSequence) AddName(name string) {
	s.instructions = append(s.instructions, name)
}

// AddInteger adds an integer literal to the sequence.
func (s *InstructionSequence) AddInteger(value int) {
	s.instructions = append(s.instructions, value)
}

// AddReal adds a real (floating-point) literal to the sequence.
func (s *InstructionSequence) AddReal(value float64) {
	s.instructions = append(s.instructions, value)
}

// AddBoolean adds a boolean literal to the sequence.
func (s *InstructionSequence) AddBoolean(value bool) {
	s.instructions = append(s.instructions, value)
}

// AddProc adds a nested proc (sub-sequence of instructions) to the sequence.
func (s *InstructionSequence) AddProc(child *InstructionSequence) {
	s.instructions = append(s.instructions, child)
}

// Execute runs all instructions in this sequence against the given execution context.
func (s *InstructionSequence) Execute(context *ExecutionContext) error {
	for _, inst := range s.instructions {
		switch v := inst.(type) {
		case string:
			cmd := context.GetOperators().GetOperator(v)
			if cmd != nil {
				cmd.Execute(context)
			} else {
				return fmt.Errorf("unknown operator or name: %s", v)
			}
		default:
			context.Push(inst)
		}
	}

	// Handle top-level procs that simply need to be executed
	for context.Count() > 0 {
		top := context.Peek()
		if seq, ok := top.(*InstructionSequence); ok {
			context.Pop()
			if err := seq.Execute(context); err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

// ParseInstructionSequence parses Type 4 function stream text into an executable instruction sequence.
func ParseInstructionSequence(text string) (*InstructionSequence, error) {
	builder := newInstructionSequenceBuilder()
	Parse(text, builder)
	return builder.GetInstructionSequence(), nil
}

type instructionSequenceBuilder struct {
	DefaultSyntaxHandler
	mainSeq *InstructionSequence
	seqStack []*InstructionSequence
}

func newInstructionSequenceBuilder() *instructionSequenceBuilder {
	main := NewInstructionSequence()
	b := &instructionSequenceBuilder{
		mainSeq:  main,
		seqStack: []*InstructionSequence{main},
	}
	return b
}

func (b *instructionSequenceBuilder) GetInstructionSequence() *InstructionSequence {
	return b.mainSeq
}

func (b *instructionSequenceBuilder) getCurrentSequence() *InstructionSequence {
	return b.seqStack[len(b.seqStack)-1]
}

func (b *instructionSequenceBuilder) Token(token string) {
	switch token {
	case "{":
		child := NewInstructionSequence()
		b.getCurrentSequence().AddProc(child)
		b.seqStack = append(b.seqStack, child)
	case "}":
		b.seqStack = b.seqStack[:len(b.seqStack)-1]
	default:
		if val, ok := parseInt(token); ok {
			b.getCurrentSequence().AddInteger(val)
			return
		}
		if val, ok := parseFloat(token); ok {
			b.getCurrentSequence().AddReal(val)
			return
		}
		b.getCurrentSequence().AddName(token)
	}
}

// parseInt tries to parse the token as an integer. Returns (value, true) on success.
func parseInt(s string) (int, bool) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return int(v), true
}

// parseFloat tries to parse the token as a float64. Returns (value, true) on success.
func parseFloat(s string) (float64, bool) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
