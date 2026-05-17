package charstrings

import (
	"fmt"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/cfftypes"
)

// notDefined is the standard glyph name for undefined characters in CFF fonts.
const notDefined = ".notdef"

// CommandIdentifier identifies a command within a Type 2 CharString command sequence.
type CommandIdentifier struct {
	// CommandIndex is the index of the numeric operand preceding this command.
	CommandIndex int

	// IsMultiByteCommand indicates whether this is a two-byte command (opcode 12 + second byte).
	IsMultiByteCommand bool

	// CommandId is the raw opcode identifying the command.
	CommandId byte
}

// CommandSequence holds the decoded operands and commands for a single Type 2 CharString or subroutine.
type CommandSequence struct {
	// Values are the numeric operands in order of appearance.
	Values []float64

	// CommandIdentifiers are the parsed commands with their positions.
	CommandIdentifiers []CommandIdentifier
}

// GetCommandsAt returns all command identifiers that occur at the given operand index.
func (cs *CommandSequence) GetCommandsAt(index int) []CommandIdentifier {
	result := make([]CommandIdentifier, 0)
	for _, id := range cs.CommandIdentifiers {
		if id.CommandIndex == index {
			result = append(result, id)
		}
	}
	return result
}

// String returns a human-readable representation of the command sequence, listing each
// numeric value and command name in order of evaluation. Implements fmt.Stringer.
func (cs *CommandSequence) String() string {
	var sb strings.Builder
	for i := -1; i < len(cs.Values); i++ {
		if i >= 0 {
			sb.WriteString(fmt.Sprintf("%.1f\n", cs.Values[i]))
		}

		for _, identifier := range cs.GetCommandsAt(i + 1) {
			cmd := GetCommand(identifier)
			sb.WriteString(cmd.Name())
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// Type2CharStrings stores the decoded command sequences for all glyphs in a CFF font,
// as well as the local and global subroutines. Glyphs are lazily evaluated and cached.
type Type2CharStrings struct {
	charStrings map[string]CommandSequence
	glyphs      map[string]*Type2Glyph
	mu          sync.Mutex
}

// NewType2CharStrings creates a new Type2CharStrings with the given parsed command sequences.
func NewType2CharStrings(charStrings map[string]CommandSequence) *Type2CharStrings {
	if charStrings == nil {
		panic("charStrings cannot be nil")
	}

	return &Type2CharStrings{
		charStrings: charStrings,
		glyphs:      make(map[string]*Type2Glyph),
	}
}

// CharacterNames returns the names of all characters that have a registered charstring.
func (t *Type2CharStrings) CharacterNames() []string {
	names := make([]string, 0, len(t.charStrings))
	for name := range t.charStrings {
		names = append(names, name)
	}
	return names
}

// Generate evaluates the CharString for the character with the given name and returns
// the path constructed for the glyph. Results are cached so repeated calls for the same
// name return the cached glyph without re-evaluation.
func (t *Type2CharStrings) Generate(name string, defaultWidthX, nominalWidthX float64) (cfftypes.Type2GlyphResult, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if glyph, ok := t.glyphs[name]; ok {
		return glyph, nil
	}

	seq, ok := t.charStrings[name]
	if !ok {
		var found bool
		seq, found = t.charStrings[notDefined]
		if !found {
			return nil, fmt.Errorf("no charstring sequence with the name /%s in this font", name)
		}
	}

	glyph, err := runCharString(seq, defaultWidthX, nominalWidthX)
	if err != nil {
		return nil, fmt.Errorf("failed to interpret charstring for symbol with name: %s. Commands: %s. Error: %w", name, seq.String(), err)
	}

	t.glyphs[name] = glyph
	return glyph, nil
}

// runCharString interprets a single command sequence and produces a glyph path and width.
func runCharString(sequence CommandSequence, defaultWidthX, nominalWidthX float64) (*Type2Glyph, error) {
	ctx := NewType2BuildCharContext()

	hasRunStackClearingCommand := false

	for i := -1; i < len(sequence.Values); i++ {
		if i >= 0 {
			ctx.Stack().Push(sequence.Values[i])
		}

		for _, identifier := range sequence.GetCommandsAt(i + 1) {
			cmd := GetCommand(identifier)

			isOnlyCommand := len(sequence.Values)+len(sequence.CommandIdentifiers) == 1

			if !hasRunStackClearingCommand {
				hasRunStackClearingCommand = true
				switch cmd.Name() {
				case "hstem", "hstemhm", "vstemhm", "vstem":
					if ctx.Stack().Length()%2 != 0 {
						bottom, err := ctx.Stack().PopBottom()
						if err != nil {
							return nil, fmt.Errorf("stack underflow on width extraction: %w", err)
						}
						ctx.SetWidth(nominalWidthX + bottom)
					}

				case "hmoveto", "vmoveto":
					setWidthFromArgsIfPresent(ctx, nominalWidthX, 1)

				case "rmoveto":
					setWidthFromArgsIfPresent(ctx, nominalWidthX, 2)

				case "cntrmask", "hintmask":
					setWidthFromArgsIfPresent(ctx, nominalWidthX, 0)

				case "endchar":
					if isOnlyCommand {
						ctx.SetWidth(defaultWidthX)
					} else {
						setWidthFromArgsIfPresent(ctx, nominalWidthX, 0)
					}

				default:
					hasRunStackClearingCommand = false
				}
			}

			cmd.Run(ctx)
		}
	}

	width := ctx.Width()
	return NewType2Glyph(ctx.Path(), width), nil
}

// setWidthFromArgsIfPresent sets the glyph width from the stack if there are enough
// arguments remaining after accounting for the expected argument length of the command.
func setWidthFromArgsIfPresent(ctx *Type2BuildCharContext, nominalWidthX float64, expectedArgLength int) {
	if ctx.Stack().Length() > expectedArgLength {
		bottom, err := ctx.Stack().PopBottom()
		if err == nil {
			ctx.SetWidth(nominalWidthX + bottom)
		}
	}
}

// Type2Glyph holds the rendered path and width for a single glyph decoded from a
// Type 2 CharString. The width may be expressed as an absolute value or as a delta
// from the font's nominal width X.
type Type2Glyph struct {
	path  []*core.PdfSubpath
	width *float64
}

// NewType2Glyph creates a new Type2Glyph with the given path and optional width.
// Pass nil for width if no width was defined in the charstring.
func NewType2Glyph(path []*core.PdfSubpath, width *float64) *Type2Glyph {
	if path == nil {
		panic("path cannot be nil")
	}

	return &Type2Glyph{
		path:  path,
		width: width,
	}
}

// Path returns the list of subpaths that make up the glyph outline.
func (g *Type2Glyph) Path() []*core.PdfSubpath {
	return g.path
}

// Width returns the character width if one was defined in the charstring, or nil otherwise.
func (g *Type2Glyph) Width() *float64 {
	return g.width
}
