package charstrings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands/arithmetic"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands/hint"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands/pathconstruction"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings/commands/startfinishoutline"
)

// SourceType indicates the origin of a Type 1 charstring byte sequence.
type SourceType int

const (
	// Subroutine indicates the bytes come from a subroutine definition.
	Subroutine SourceType = iota
	// Charstring indicates the bytes come from a character string definition.
	Charstring
)

// Type1CharstringDecryptedBytes holds the decrypted byte sequence for a single
// Type 1 charstring or subroutine along with its metadata.
type Type1CharstringDecryptedBytes struct {
	bytes   []byte
	Index   int
	Name    string
	Source  SourceType
}

// Bytes returns the decrypted bytes as a read-only span.
func (d *Type1CharstringDecryptedBytes) Bytes() []byte {
	return d.bytes
}

// NewType1CharstringDecryptedBytes creates a new instance for a subroutine.
func NewType1CharstringDecryptedBytes(bytes []byte, index int) *Type1CharstringDecryptedBytes {
	if bytes == nil {
		panic("bytes cannot be nil")
	}
	return &Type1CharstringDecryptedBytes{
		bytes:  bytes,
		Index:  index,
		Name:   ".notdef",
		Source: Subroutine,
	}
}

// NewType1CharstringWithName creates a new instance for a named charstring.
func NewType1CharstringWithName(name string, bytes []byte, index int) *Type1CharstringDecryptedBytes {
	if bytes == nil {
		panic("bytes cannot be nil")
	}
	if name == "" {
		name = strconv.Itoa(index)
	}
	return &Type1CharstringDecryptedBytes{
		bytes:  bytes,
		Index:  index,
		Name:   name,
		Source: Charstring,
	}
}

// Type1CharStringParser decodes a set of CharStrings to their corresponding Type 1 BuildChar operations.
// A charstring is an encrypted sequence of unsigned 8-bit bytes that encode integers and commands.
// Type 1 BuildChar, when interpreting a charstring, will first decrypt it and then decode
// its bytes one at a time in sequence. The value in a byte indicates a command, a number, or
// subsequent bytes that are to be interpreted in a special way. Once the bytes are decoded into
// numbers and commands, the execution proceeds similarly to PostScript language operation using
// its own operand stack (distinct from the PostScript interpreter operand stack). This stack holds
// up to 24 numeric entries; a number is pushed onto it, and a command expects its arguments in order.
type Type1CharStringParser struct{}

// Parse decodes charstrings and subroutines into their command sequences.
func (p *Type1CharStringParser) Parse(charStrings, subroutines []*Type1CharstringDecryptedBytes) (*Type1CharStrings, error) {
	if charStrings == nil {
		return nil, errors.New("charStrings cannot be nil")
	}
	if subroutines == nil {
		return nil, errors.New("subroutines cannot be nil")
	}

	charStringResults := make(map[string]*CommandSequence, len(charStrings))
	charStringIndexToName := make(map[int]string, len(charStrings))

	for i, cs := range charStrings {
		commandSequence := parseSingle(cs.bytes)
		charStringResults[cs.Name] = &CommandSequence{Commands: commandSequence}
		charStringIndexToName[i] = cs.Name
	}

	subroutineResults := make(map[int]*CommandSequence, len(subroutines))
	for _, sub := range subroutines {
		commandSequence := parseSingle(sub.bytes)
		subroutineResults[sub.Index] = &CommandSequence{Commands: commandSequence}
	}

	return &Type1CharStrings{
		charStringIndexToName: charStringIndexToName,
		charStrings:           charStringResults,
		subroutines:           subroutineResults,
		glyphs:                make(map[string][]*core.PdfSubpath),
	}, nil
}

// Type1CharStrings holds the decoded command sequences for Type 1 CharStrings.
type Type1CharStrings struct {
	charStrings           map[string]*CommandSequence
	charStringIndexToName map[int]string
	subroutines           map[int]*CommandSequence
	glyphs                map[string][]*core.PdfSubpath
	glyphsMu              sync.Mutex
}

// CharStrings returns the mapping from character names to their command sequences.
func (t *Type1CharStrings) CharStrings() map[string]*CommandSequence {
	return t.charStrings
}

// Subroutines returns the mapping from subroutine indices to their command sequences.
func (t *Type1CharStrings) Subroutines() map[int]*CommandSequence {
	return t.subroutines
}

// GetCharStringNameByIndex looks up a character name by its index in the charstrings array.
func (t *Type1CharStrings) GetCharStringNameByIndex(index int) (string, bool) {
	name, ok := t.charStringIndexToName[index]
	return name, ok
}

// TryGenerate attempts to generate the glyph path for the given character name.
// Results are cached for subsequent calls. Returns the path and true on success,
// nil and false if the character has no charstring definition or generation fails.
func (t *Type1CharStrings) TryGenerate(name string) ([]*core.PdfSubpath, bool) {
	t.glyphsMu.Lock()
	defer t.glyphsMu.Unlock()

	if path, ok := t.glyphs[name]; ok {
		return path, true
	}

	sequence, ok := t.charStrings[name]
	if !ok {
		return nil, false
	}

	path := t.run(sequence)
	t.glyphs[name] = path

	return path, true
}

func (t *Type1CharStrings) run(sequence *CommandSequence) []*core.PdfSubpath {
	subroutinesAny := make(map[int]any, len(t.subroutines))
	for k, v := range t.subroutines {
		subroutinesAny[k] = v.Commands
	}

	context := commands.NewType1BuildCharContext(
		subroutinesAny,
		func(i int) []*core.PdfSubpath {
			name, ok := t.charStringIndexToName[i]
			if !ok {
				panic(fmt.Sprintf("tried to retrieve Type 1 charstring by index %d which did not exist", i))
			}

			t.glyphsMu.Lock()
			if result, ok := t.glyphs[name]; ok {
				t.glyphsMu.Unlock()
				return result
			}
			t.glyphsMu.Unlock()

			charstring, ok := t.charStrings[name]
			if !ok {
				panic(fmt.Sprintf("tried to retrieve Type 1 charstring by index %d which mapped to name %s but was not found in the charstrings", i, name))
			}

			path := t.run(charstring)
			t.glyphsMu.Lock()
			t.glyphs[name] = path
			t.glyphsMu.Unlock()

			return path
		},
		func(s string) []*core.PdfSubpath {
			t.glyphsMu.Lock()
			if result, ok := t.glyphs[s]; ok {
				t.glyphsMu.Unlock()
				return result
			}
			t.glyphsMu.Unlock()

			charstring, ok := t.charStrings[s]
			if !ok {
				panic(fmt.Sprintf("tried to retrieve Type 1 charstring by name %s but it was not found in the charstrings", s))
			}

			path := t.run(charstring)
			t.glyphsMu.Lock()
			t.glyphs[s] = path
			t.glyphsMu.Unlock()

			return path
		},
	)

	for _, command := range sequence.Commands {
		if num, ok := command.TryGetFirst(); ok {
			context.Stack().Push(num)
		} else if lazyCmd, ok := command.TryGetSecond(); ok {
			lazyCmd.Run(context)
		}
	}

	return context.Path()
}

// CommandSequence holds the ordered list of numbers and commands for a Type 1 charstring or subroutine.
type CommandSequence struct {
	Commands []core.Union[float64, *commands.LazyType1Command]
}

// String returns a comma-separated string representation of all commands in the sequence.
func (cs *CommandSequence) String() string {
	parts := make([]string, len(cs.Commands))
	for i, cmd := range cs.Commands {
		parts[i] = cmd.Match(
			func(f float64) any { return fmt.Sprintf("%g", f) },
			func(l *commands.LazyType1Command) any { return l.String() },
		).(string)
	}
	return strings.Join(parts, ", ")
}

func parseSingle(charStringBytes []byte) []core.Union[float64, *commands.LazyType1Command] {
	interpreted := make([]core.Union[float64, *commands.LazyType1Command], 0, len(charStringBytes))

	for i := 0; i < len(charStringBytes); i++ {
		b := charStringBytes[i]

		if b <= 31 {
			command := GetCommand(b, charStringBytes, &i)
			if command == nil {
				continue
			}
			interpreted = append(interpreted, core.UnionTwo[float64, *commands.LazyType1Command](command))
		} else {
			val := interpretNumber(b, charStringBytes, &i)
			interpreted = append(interpreted, core.UnionOne[float64, *commands.LazyType1Command](float64(val)))
		}
	}

	return interpreted
}

func interpretNumber(b byte, bytes []byte, i *int) int {
	if b >= 32 && b <= 246 {
		return int(b) - 139
	}

	if b >= 247 && b <= 250 {
		*i++
		w := bytes[*i]
		return (int(b)-247)*256 + int(w) + 108
	}

	if b >= 251 && b <= 254 {
		*i++
		w := bytes[*i]
		return -((int(b) - 251) * 256) - int(w) - 108
	}

	*i++
	b1 := bytes[*i]
	*i++
	b2 := bytes[*i]
	*i++
	b3 := bytes[*i]
	*i++
	b4 := bytes[*i]

	return int(b1)<<24 | int(b2)<<16 | int(b3)<<8 | int(b4)
}

// GetCommand maps a byte value to the corresponding Type 1 charstring command.
// For two-byte opcodes (byte value 12), it reads the next byte from the stream
// and increments i accordingly. Returns nil for unrecognized opcode sequences.
func GetCommand(v byte, bytes []byte, i *int) *commands.LazyType1Command {
	switch v {
	case 1:
		return hint.HStemLazy
	case 3:
		return hint.VStemLazy
	case 4:
		return pathconstruction.VMoveToLazy
	case 5:
		return pathconstruction.RLineToLazy
	case 6:
		return pathconstruction.HLineToLazy
	case 7:
		return pathconstruction.VLineToLazy
	case 8:
		return pathconstruction.Lazy // relative_r_curve_to_command exports Lazy
	case 9:
		return pathconstruction.ClosePathLazy
	case 10:
		return arithmetic.CallSubrLazy
	case 11:
		return arithmetic.ReturnLazy
	case 13:
		return startfinishoutline.HsbwLazy
	case 14:
		return startfinishoutline.EndCharLazy
	case 21:
		return pathconstruction.RMoveToLazy
	case 22:
		return pathconstruction.HMoveToLazy
	case 30:
		return pathconstruction.VhCurveToLazy
	case 31:
		return pathconstruction.HvCurveToLazy
	case 12:
		*i++
		next := bytes[*i]
		switch next {
		case 0:
			return hint.DotSectionLazy
		case 1:
			return hint.VStem3Lazy
		case 2:
			return hint.HStem3Lazy
		case 6:
			return startfinishoutline.SeacLazy
		case 7:
			return startfinishoutline.SbwLazy
		case 12:
			return arithmetic.DivLazy
		case 16:
			return arithmetic.Lazy // call_other_subr_command exports Lazy
		case 17:
			return arithmetic.PopLazy
		case 33:
			return arithmetic.SetCurrentPointLazy
		}
	}

	return nil
}
