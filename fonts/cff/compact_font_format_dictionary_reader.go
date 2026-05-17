package cff

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
)

// Operand represents a single operand value in a CFF dictionary.
// An operand is either an integer or a floating-point number.
type Operand struct {
	Int    *int
	Double float64
}

// NewOperandInt creates an Operand holding an integer value.
func NewOperandInt(v int) Operand {
	return Operand{Int: &v, Double: float64(v)}
}

// NewOperandDouble creates an Operand holding a floating-point value.
func NewOperandDouble(v float64) Operand {
	return Operand{Double: v}
}

// OperandKey represents the operator key in a CFF dictionary entry.
// Keys are either single-byte (0-21, or 12 with an extension byte).
type OperandKey struct {
	Byte0 byte
	Byte1 *byte
}

// NewOperandKey creates a single-byte operand key.
func NewOperandKey(b0 byte) OperandKey {
	return OperandKey{Byte0: b0}
}

// NewOperandKeyTwoByte creates a two-byte operand key (operator 12 with extension).
func NewOperandKeyTwoByte(b0, b1 byte) OperandKey {
	return OperandKey{Byte0: b0, Byte1: &b1}
}

// applyOperationFunc is the signature for the callback that handles individual
// dictionary operators during parsing. It receives the builder (as any), the
// collected operands, the operator key, and the string index.
type applyOperationFunc func(builder any, operands []Operand, key OperandKey, stringIndex []string)

// readDictionary parses a CFF dictionary from data until CanRead returns false.
// For each operator encountered, it calls applyOp with the accumulated operands.
func readDictionary(
	builder any,
	data *CompactFontFormatData,
	stringIndex []string,
	applyOp applyOperationFunc,
) any {
	var operands []Operand

	for data.CanRead() {
		operands = operands[:0]
		infiniteLoopProtection := 0

		for {
			infiniteLoopProtection++
			if infiniteLoopProtection > 256 {
				panic("got caught in an infinite loop trying to read a CFF dictionary")
			}

			byte0, _ := data.ReadByte()

			if byte0 <= 21 {
				var key OperandKey
				if byte0 == 12 {
					b1, _ := data.ReadByte()
					key = NewOperandKeyTwoByte(byte0, b1)
				} else {
					key = NewOperandKey(byte0)
				}
				applyOp(builder, operands, key, stringIndex)
				break
			}

			switch {
			case byte0 == 28:
				b1, _ := data.ReadByte()
				b2, _ := data.ReadByte()
				v := int(b1)<<8 | int(b2)
				operands = append(operands, NewOperandInt(v))
			case byte0 == 29:
				b1, _ := data.ReadByte()
				b2, _ := data.ReadByte()
				b3, _ := data.ReadByte()
				b4, _ := data.ReadByte()
				v := int(b1)<<24 | int(b2)<<16 | int(b3)<<8 | int(b4)
				operands = append(operands, NewOperandInt(v))
			case byte0 == 30:
				d, err := readRealNumber(data)
				if err != nil {
					panic(fmt.Sprintf("failed to parse real number in CFF dictionary: %v", err))
				}
				operands = append(operands, NewOperandDouble(d))
			case byte0 >= 32 && byte0 <= 246:
				v := int(byte0) - 139
				operands = append(operands, NewOperandInt(v))
			case byte0 >= 247 && byte0 <= 250:
				b1, _ := data.ReadByte()
				v := (int(byte0)-247)*256 + int(b1) + 108
				operands = append(operands, NewOperandInt(v))
			case byte0 >= 251 && byte0 <= 254:
				b1, _ := data.ReadByte()
				v := -(int(byte0)-251)*256 - int(b1) - 108
				operands = append(operands, NewOperandInt(v))
			default:
				panic(fmt.Sprintf("the first dictionary byte was not in the range 29-254. Got %d.", byte0))
			}
		}
	}

	return builder
}

// readRealNumber parses a CFF real number encoded as packed nibbles.
func readRealNumber(data *CompactFontFormatData) (float64, error) {
	var sb strings.Builder
	done := false
	exponentMissing := false
	hasExponent := false

	for !done {
		b, _ := data.ReadByte()
		nibble1 := b / 16
		nibble2 := b % 16

		for i := 0; i < 2; i++ {
			var nibble byte
			if i == 0 {
				nibble = nibble1
			} else {
				nibble = nibble2
			}

			switch nibble {
			case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9:
				sb.WriteByte(byte(nibble) + '0')
				exponentMissing = false
			case 0xa:
				sb.WriteByte('.')
			case 0xb:
				if hasExponent {
					break
				}
				sb.WriteByte('E')
				exponentMissing = true
				hasExponent = true
			case 0xc:
				if hasExponent {
					break
				}
				sb.WriteString("E-")
				exponentMissing = true
				hasExponent = true
			case 0xd:
				// End of number marker — ignored inside the nibble loop;
				// the outer loop terminates on 0xf.
			case 0xe:
				sb.WriteByte('-')
			case 0xf:
				done = true
			default:
				return 0, fmt.Errorf("unexpected nibble value in CFF real number: %d", nibble)
			}

			if done {
				break
			}
		}
	}

	if exponentMissing {
		sb.WriteByte('0')
	}

	str := sb.String()
	if len(str) == 0 {
		return 0, nil
	}

	if hasExponent {
		return strconv.ParseFloat(str, 64)
	}
	return strconv.ParseFloat(str, 64)
}

// GetString resolves a string identifier (SID) from the first operand.
// SIDs 0-390 map to CFF standard strings; 391+ index into stringIndex.
func GetString(operands []Operand, stringIndex []string) (string, error) {
	if len(operands) == 0 {
		return "", errors.New("cannot read a string from an empty operands array")
	}

	if operands[0].Int == nil {
		return "", fmt.Errorf("the first operand for reading a string was not an integer. Got: %g", operands[0].Double)
	}

	index := *operands[0].Int

	if index >= 0 && index <= 390 {
		return GetName(index), nil
	}

	stringIndexIdx := index - 391
	if stringIndexIdx >= 0 && stringIndexIdx < len(stringIndex) {
		return stringIndex[stringIndexIdx], nil
	}

	return fmt.Sprintf("SID%d", index), nil
}

// GetBoundingBox constructs a PdfRectangle from four numeric operands.
func GetBoundingBox(operands []Operand) core.PdfRectangle {
	if len(operands) != 4 {
		return core.NewPdfRectangleFromInt(0, 0, 0, 0)
	}

	return core.NewPdfRectangleFloat(
		operands[0].Double, operands[1].Double,
		operands[2].Double, operands[3].Double,
	)
}

// ToArray converts all operand values to a float64 slice.
func ToArray(operands []Operand) []float64 {
	result := make([]float64, len(operands))
	for i := range operands {
		result[i] = operands[i].Double
	}
	return result
}

// GetIntOrDefault returns the integer value of the first operand or defaultValue.
func GetIntOrDefault(operands []Operand, defaultValue int) int {
	if len(operands) == 0 {
		return defaultValue
	}

	if operands[0].Int != nil {
		return *operands[0].Int
	}

	return defaultValue
}

// ReadDeltaToIntArray converts operands to a cumulative-sum integer array.
// The first element is taken directly; subsequent elements are added cumulatively.
func ReadDeltaToIntArray(operands []Operand) []int {
	results := make([]int, len(operands))

	if len(operands) == 0 {
		return results
	}

	results[0] = int(operands[0].Double)

	for i := 1; i < len(operands); i++ {
		results[i] = results[i-1] + int(operands[i].Double)
	}

	return results
}

// ReadDeltaToArray converts operands to a cumulative-sum float64 array.
// The first element is taken directly; subsequent elements are added cumulatively.
func ReadDeltaToArray(operands []Operand) []float64 {
	results := make([]float64, len(operands))

	if len(operands) == 0 {
		return results
	}

	results[0] = operands[0].Double

	for i := 1; i < len(operands); i++ {
		results[i] = results[i-1] + operands[i].Double
	}

	return results
}
