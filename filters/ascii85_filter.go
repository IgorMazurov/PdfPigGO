package filters

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

const (
	ascii85EmptyBlock           = byte('z')
	ascii85Offset               = byte('!')
	ascii85EmptyCharPadding     = byte('u')
	ascii85EndOfDataFirst       = '~'
	ascii85EndOfDataSecond      = '>'
)

var ascii85PowerByIndex = [5]int{1, 85, 85 * 85, 85 * 85 * 85, 85 * 85 * 85 * 85}

// Ascii85Filter decodes ASCII85 (Base85) encoded PDF stream data.
type Ascii85Filter struct{}

// NewAscii85Filter creates a new Ascii85Filter instance.
func NewAscii85Filter() *Ascii85Filter {
	return &Ascii85Filter{}
}

// IsSupported reports whether this filter is supported. Always true.
func (f *Ascii85Filter) IsSupported() bool {
	return true
}

// Decode decodes ASCII85-encoded input bytes.
func (f *Ascii85Filter) Decode(input []byte, streamDictionary *tokens.DictionaryToken, provider FilterProvider, filterIndex int) ([]byte, error) {
	asciiBuffer := [5]byte{}
	index := 0

	writer := core.NewArrayPoolBufferWriter()
	defer writer.Dispose()

	for i := 0; i < len(input); i++ {
		value := input[i]

		if isAscii85Whitespace(value) {
			continue
		}

		if value == ascii85EndOfDataFirst {
			if i == len(input)-1 || input[i+1] == ascii85EndOfDataSecond {
				if index > 0 {
					writeAscii85Data(&asciiBuffer, index, writer, true)
				}
				index = 0
				break
			}
		}

		if value == ascii85EmptyBlock {
			if index > 0 {
				writer.Dispose()
				return nil, fmt.Errorf("unexpected empty block marker at position %d", i)
			}

			for j := 0; j < 4; j++ {
				writer.WriteSingle(0)
			}
			index = 0
		} else {
			asciiBuffer[index] = value - ascii85Offset
			index++
		}

		if index == 5 {
			writeAscii85Data(&asciiBuffer, index, writer, false)
			index = 0
		}
	}

	if index > 0 {
		writeAscii85Data(&asciiBuffer, index, writer, true)
	}

	result := make([]byte, len(writer.WrittenMemory()))
	copy(result, writer.WrittenMemory())
	return result, nil
}

func writeAscii85Data(ascii *[5]byte, index int, writer *core.ArrayPoolBufferWriter, isAtEnd bool) {
	if index < 2 {
		if !isAtEnd {
			writer.Dispose()
		}
		return
	}

	for i := index; i < 5; i++ {
		ascii[i] = ascii85EmptyCharPadding - ascii85Offset
	}

	value := 0
	value += int(ascii[0]) * ascii85PowerByIndex[4]
	value += int(ascii[1]) * ascii85PowerByIndex[3]
	value += int(ascii[2]) * ascii85PowerByIndex[2]
	value += int(ascii[3]) * ascii85PowerByIndex[1]
	value += int(ascii[4]) * ascii85PowerByIndex[0]

	writer.WriteSingle(byte(value >> 24))

	if index > 2 {
		writer.WriteSingle(byte(value >> 16))
	}

	if index > 3 {
		writer.WriteSingle(byte(value >> 8))
	}

	if index > 4 {
		writer.WriteSingle(byte(value))
	}
}

func isAscii85Whitespace(b byte) bool {
	return b == '\r' || b == '\n' || b == ' '
}

var _ Filter = (*Ascii85Filter)(nil)
