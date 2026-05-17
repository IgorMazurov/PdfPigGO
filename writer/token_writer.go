package writer

import (
	"crypto/rand"
	"fmt"
	"io"
	"sort"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TokenWriter writes any type of Token to the corresponding PDF document
// format output.
type TokenWriter interface {
	// WriteToken writes the given token to the output stream with the correct
	// PDF format and encoding, including whitespace and line breaks as applicable.
	WriteToken(token tokens.Token, w io.Writer) error

	// WriteObject writes pre-serialized data as an object token to the output
	// stream using the given object number and generation.
	WriteObject(objectNumber int64, generation int, data []byte, w io.Writer) error

	// WriteCrossReferenceTable writes a valid single-section cross-reference
	// (xref) table plus trailer dictionary to the output for the set of object
	// offsets. The catalogRef is referenced from the trailer dictionary.
	// docInfoRef is the document information dictionary reference, or nil if
	// absent.
	WriteCrossReferenceTable(
		objectOffsets map[core.IndirectReference]int64,
		catalogRef core.IndirectReference,
		w io.Writer,
		docInfoRef *core.IndirectReference,
	) error

	// WritingPageContents reports whether the writer is currently writing page
	// contents.
	WritingPageContents() bool

	// SetWritingPageContents sets whether the writer is currently writing page
	// contents.
	SetWritingPageContents(v bool)
}

// tokenWriterImpl is the default implementation of TokenWriter that writes all
// tokens without any transformation.
type tokenWriterImpl struct {
	writingPageContents bool
}

// Instance is the singleton instance of the default TokenWriter.
var Instance TokenWriter = NewTokenWriter()

// NewTokenWriter creates a new default TokenWriter.
func NewTokenWriter() *tokenWriterImpl {
	return &tokenWriterImpl{}
}

// WriteToken writes the given token to the output stream with the correct PDF
// format and encoding, including whitespace and line breaks as applicable.
func (w *tokenWriterImpl) WriteToken(token tokens.Token, out io.Writer) error {
	if token == nil {
		return writeNullToken(out)
	}

	switch t := token.(type) {
	case *tokens.ArrayToken:
		return writeArray(t, out)
	case *tokens.BooleanToken:
		return writeBoolean(t, out)
	case *tokens.CommentToken:
		return writeComment(t, out)
	case *tokens.DictionaryToken:
		return writeDictionary(t, out)
	case *tokens.HexToken:
		return writeHex(t, out)
	case *tokens.IndirectReferenceToken:
		return writeIndirectReference(t, out)
	case *tokens.NameToken:
		return writeName(t, out)
	case *tokens.NullToken:
		return writeNullToken(out)
	case *tokens.NumericToken:
		return writeNumber(t, out)
	case *tokens.ObjectToken:
		return writeObjectToken(t, out)
	case *tokens.StreamToken:
		return writeStream(t, out)
	case *tokens.StringToken:
		return writeString(t, out)
	default:
		return fmt.Errorf("attempted to write unknown token type %T", token)
	}
}

// WriteObject writes pre-serialized data as an object token to the output
// stream using the given object number and generation.
func (w *tokenWriterImpl) WriteObject(objectNumber int64, generation int, data []byte, out io.Writer) error {
	if err := writeLong(objectNumber, out); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if err := writeInt(generation, out); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if _, err := out.Write(streamObjStart); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if _, err := out.Write(data); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if _, err := out.Write(streamObjEnd); err != nil {
		return err
	}
	return writeLineBreak(out)
}

// WriteCrossReferenceTable writes a valid single-section cross-reference table
// plus trailer dictionary to the output for the set of object offsets.
func (w *tokenWriterImpl) WriteCrossReferenceTable(
	objectOffsets map[core.IndirectReference]int64,
	catalogRef core.IndirectReference,
	out io.Writer,
	docInfoRef *core.IndirectReference,
) error {
	if len(objectOffsets) == 0 {
		return fmt.Errorf("cannot write empty cross reference table")
	}

	if err := writeLineBreak(out); err != nil {
		return err
	}
	position := int64(0)
	if seeker, ok := out.(interface{ Seek(int64, int) (int64, error) }); ok {
		position, _ = seeker.Seek(0, io.SeekCurrent)
	}
	if _, err := out.Write(streamXref); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}

	sorted := make([]core.IndirectReference, 0, len(objectOffsets))
	for ref := range objectOffsets {
		sorted = append(sorted, ref)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ObjectNumber() < sorted[j].ObjectNumber()
	})

	sets := buildXrefSets(sorted, objectOffsets)

	for _, set := range sets {
		if err := writeLong(set.first, out); err != nil {
			return err
		}
		if err := writeWhiteSpace(out); err != nil {
			return err
		}
		if err := writeInt(len(set.entries), out); err != nil {
			return err
		}
		if err := writeWhiteSpace(out); err != nil {
			return err
		}
		if err := writeLineBreak(out); err != nil {
			return err
		}

		for _, entry := range set.entries {
			if entry == nil {
				if err := writeFirstXrefEmptyEntry(out); err != nil {
					return err
				}
			} else {
				paddedOffset := fmt.Sprintf("%010d", entry.offset)
				if _, err := out.Write([]byte(paddedOffset)); err != nil {
					return err
				}
				if err := writeWhiteSpace(out); err != nil {
					return err
				}
				generationStr := fmt.Sprintf("%05d", entry.generation)
				if _, err := out.Write([]byte(generationStr)); err != nil {
					return err
				}
				if err := writeWhiteSpace(out); err != nil {
					return err
				}
				if _, err := out.Write([]byte{inUseEntry}); err != nil {
					return err
				}
				if err := writeWhiteSpace(out); err != nil {
					return err
				}
				if err := writeLineBreak(out); err != nil {
					return err
				}
			}
		}
	}

	if _, err := out.Write(streamTrailer); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}

	trailerData := map[*tokens.NameToken]tokens.Token{
		tokens.Size: tokens.NewNumericTokenFromInt(len(objectOffsets) + 1),
		tokens.Root: tokens.NewIndirectReferenceToken(catalogRef),
		tokens.Id:   newRandomIdArray(),
	}

	if docInfoRef != nil {
		trailerData[tokens.Info] = tokens.NewIndirectReferenceToken(*docInfoRef)
	}

	trailerDict, _ := tokens.NewDictionary(trailerData)
	if err := writeDictionary(trailerDict, out); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}

	if _, err := out.Write(streamStartXref); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if err := writeLong(position, out); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}

	if _, err := out.Write(streamEof); err != nil {
		return err
	}

	return nil
}

// WritingPageContents reports whether the writer is currently writing page contents.
func (w *tokenWriterImpl) WritingPageContents() bool {
	return w.writingPageContents
}

// SetWritingPageContents sets whether the writer is currently writing page contents.
func (w *tokenWriterImpl) SetWritingPageContents(v bool) {
	w.writingPageContents = v
}

// writeStream writes a stream token without any transformation.
func writeStream(streamToken *tokens.StreamToken, out io.Writer) error {
	if err := writeDictionary(streamToken.StreamDictionary, out); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if _, err := out.Write(streamStreamStart); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if _, err := out.Write(streamToken.Data()); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if _, err := out.Write(streamStreamEnd); err != nil {
		return err
	}
	return nil
}

// newRandomIdArray creates an array token with two random hex values for the PDF ID entry.
func newRandomIdArray() *tokens.ArrayToken {
	id1 := make([]byte, 16)
	rand.Read(id1)
	id2 := make([]byte, 16)
	rand.Read(id2)
	return tokens.NewArrayToken([]tokens.Token{
		tokens.NewHexTokenFromRaw(id1),
		tokens.NewHexTokenFromRaw(id2),
	})
}

// --- Shared constants for token writing ---

var (
	streamObjStart    = []byte("obj")
	streamObjEnd      = []byte("endobj")
	streamXref        = []byte("xref")
	streamTrailer     = []byte("trailer")
	streamStartXref   = []byte("startxref")
	streamEof         = []byte("%%EOF")
	streamStreamStart = []byte("stream")
	streamStreamEnd   = []byte("endstream")
)

const (
	arrayStartByte = '['
	arrayEndByte   = ']'
	commentByte    = '%'
	inUseEntry     = 'n'
	nameStartByte  = '/'
	nullTokenBytes = "null"
	rByte          = 'R'
	stringStart    = '('
	stringEnd      = ')'
	hexStart       = '<'
	hexEnd         = '>'
	whitespace     = ' '
)

var (
	dictionaryStart = []byte("<<")
	dictionaryEnd   = []byte(">>")
	trueBytes       = []byte("true")
	falseBytes      = []byte("false")
	newlineByte     = []byte("\n")
)

// --- Shared helper functions for writing tokens ---

func writeArray(arr *tokens.ArrayToken, out io.Writer) error {
	if _, err := out.Write([]byte{arrayStartByte}); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	for _, value := range arr.Data() {
		if err := writeTokenSimple(value, out); err != nil {
			return err
		}
	}
	if _, err := out.Write([]byte{arrayEndByte}); err != nil {
		return err
	}
	return writeWhiteSpace(out)
}

func writeBoolean(b *tokens.BooleanToken, out io.Writer) error {
	var bts []byte
	if b.Data() {
		bts = trueBytes
	} else {
		bts = falseBytes
	}
	if _, err := out.Write(bts); err != nil {
		return err
	}
	return writeWhiteSpace(out)
}

func writeComment(c *tokens.CommentToken, out io.Writer) error {
	if _, err := out.Write([]byte{commentByte}); err != nil {
		return err
	}
	if _, err := out.Write([]byte(c.Data())); err != nil {
		return err
	}
	return writeLineBreak(out)
}

func writeDictionary(dict *tokens.DictionaryToken, out io.Writer) error {
	if _, err := out.Write(dictionaryStart); err != nil {
		return err
	}
	for keyStr, value := range dict.Data() {
		nameKey := tokens.Create(keyStr)
		if err := writeName(nameKey, out); err != nil {
			return err
		}
		if value == nil {
			if err := writeNullToken(out); err != nil {
				return err
			}
		} else {
			if err := writeTokenSimple(value, out); err != nil {
				return err
			}
		}
	}
	_, err := out.Write(dictionaryEnd)
	return err
}

func writeHex(h *tokens.HexToken, out io.Writer) error {
	if _, err := out.Write([]byte{hexStart}); err != nil {
		return err
	}
	if _, err := out.Write([]byte(h.GetHexString())); err != nil {
		return err
	}
	_, err := out.Write([]byte{hexEnd})
	return err
}

func writeIndirectReference(ref *tokens.IndirectReferenceToken, out io.Writer) error {
	data := ref.Data()
	if err := writeLong(data.ObjectNumber(), out); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if err := writeInt(data.Generation(), out); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if _, err := out.Write([]byte{rByte}); err != nil {
		return err
	}
	return writeWhiteSpace(out)
}

func writeName(name *tokens.NameToken, out io.Writer) error {
	if _, err := out.Write([]byte{nameStartByte}); err != nil {
		return err
	}
	nameStr := name.Data()
	for i := 0; i < len(nameStr); i++ {
		c := nameStr[i]
		if c < 33 || c > 126 || isDelimiter(c) {
			hex := fmt.Sprintf("#%02X", c)
			if _, err := out.Write([]byte(hex)); err != nil {
				return err
			}
		} else {
			if _, err := out.Write([]byte{c}); err != nil {
				return err
			}
		}
	}
	return writeWhiteSpace(out)
}

func isDelimiter(c byte) bool {
	switch c {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	default:
		return false
	}
}

func writeNumber(n *tokens.NumericToken, out io.Writer) error {
	if n.HasDecimalPlaces() {
		formatted := fmt.Sprintf("%.9f", n.Data())
		formatted = stripTrailingZeros(formatted)
		if _, err := out.Write([]byte(formatted)); err != nil {
			return err
		}
	} else {
		if err := writeInt(n.IntVal(), out); err != nil {
			return err
		}
	}
	return writeWhiteSpace(out)
}

func stripTrailingZeros(s string) string {
	for len(s) > 0 && s[len(s)-1] == '0' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '.' {
		s = s[:len(s)-1]
	}
	return s
}

func writeObjectToken(obj *tokens.ObjectToken, out io.Writer) error {
	num := obj.Number()
	if err := writeLong(num.ObjectNumber(), out); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if err := writeInt(num.Generation(), out); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if _, err := out.Write(streamObjStart); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if err := writeTokenSimple(obj.Data(), out); err != nil {
		return err
	}
	if err := writeLineBreak(out); err != nil {
		return err
	}
	if _, err := out.Write(streamObjEnd); err != nil {
		return err
	}
	return writeLineBreak(out)
}

// escapeNeeded and escapedChars map control characters to their PDF escape sequences.
var (
	escapeNeeded = []rune{'\r', '\n', '\t', '\b', '\f', '\\'}
	escapedChars = []rune{'r', 'n', 't', 'b', 'f', '\\'}
)

func writeString(s *tokens.StringToken, out io.Writer) error {
	if _, err := out.Write([]byte{stringStart}); err != nil {
		return err
	}

	switch s.EncodedWith() {
	case tokens.Iso88591, tokens.StringEncodingPDFDoc:
		data := []rune(s.Data())

		// If any character exceeds 255, re-encode as UTF-16BE bytes then process those.
		hasHighChars := false
		for _, r := range data {
			if r > 255 {
				hasHighChars = true
				break
			}
		}

		var values []rune
		if hasHighChars {
			utf16Bytes := tokens.NewStringToken(s.Data(), tokens.Utf16BE).GetBytes()
			values = make([]rune, len(utf16Bytes))
			for i, b := range utf16Bytes {
				values[i] = rune(b)
			}
		} else {
			values = data
		}

		for _, v := range values {
			c := int(v)
			switch c {
			case '(', ')':
				if _, err := out.Write([]byte{'\\', byte(c)}); err != nil {
					return err
				}
			default:
				ei := -1
				for i, er := range escapeNeeded {
					if int(er) == c {
						ei = i
						break
					}
				}
				if ei >= 0 {
					if _, err := out.Write([]byte{'\\', byte(escapedChars[ei])}); err != nil {
						return err
					}
				} else if c < 32 || c > 126 {
					b3 := c / 64
					b2 := (c - b3*64) / 8
					b1 := c % 8
					if _, err := out.Write([]byte{'\\', byte(b3 + '0'), byte(b2 + '0'), byte(b1 + '0')}); err != nil {
						return err
					}
				} else {
					if _, err := out.Write([]byte{byte(c)}); err != nil {
						return err
					}
				}
			}
		}

	default:
		bytes := s.GetBytes()
		if _, err := out.Write(bytes); err != nil {
			return err
		}
	}

	if _, err := out.Write([]byte{stringEnd}); err != nil {
		return err
	}
	return writeWhiteSpace(out)
}

func writeNullToken(out io.Writer) error {
	if _, err := out.Write([]byte(nullTokenBytes)); err != nil {
		return err
	}
	return writeWhiteSpace(out)
}

func writeInt(value int, out io.Writer) error {
	buf := []byte(fmt.Sprintf("%d", value))
	_, err := out.Write(buf)
	return err
}

func writeLong(value int64, out io.Writer) error {
	buf := []byte(fmt.Sprintf("%d", value))
	_, err := out.Write(buf)
	return err
}

func writeLineBreak(out io.Writer) error {
	_, err := out.Write(newlineByte)
	return err
}

func writeWhiteSpace(out io.Writer) error {
	_, err := out.Write([]byte{whitespace})
	return err
}

func writeFirstXrefEmptyEntry(out io.Writer) error {
	if _, err := out.Write([]byte("0000000000")); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if _, err := out.Write([]byte("65535")); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	if _, err := out.Write([]byte("f")); err != nil {
		return err
	}
	if err := writeWhiteSpace(out); err != nil {
		return err
	}
	return writeLineBreak(out)
}

// writeTokenSimple writes a token without handling streams (to avoid recursion).
func writeTokenSimple(token tokens.Token, out io.Writer) error {
	if token == nil {
		return writeNullToken(out)
	}

	switch t := token.(type) {
	case *tokens.ArrayToken:
		return writeArray(t, out)
	case *tokens.BooleanToken:
		return writeBoolean(t, out)
	case *tokens.CommentToken:
		return writeComment(t, out)
	case *tokens.DictionaryToken:
		return writeDictionary(t, out)
	case *tokens.HexToken:
		return writeHex(t, out)
	case *tokens.IndirectReferenceToken:
		return writeIndirectReference(t, out)
	case *tokens.NameToken:
		return writeName(t, out)
	case *tokens.NullToken:
		return writeNullToken(out)
	case *tokens.NumericToken:
		return writeNumber(t, out)
	case *tokens.ObjectToken:
		return writeObjectToken(t, out)
	case *tokens.StreamToken:
		return writeStream(t, out)
	case *tokens.StringToken:
		return writeString(t, out)
	default:
		return fmt.Errorf("attempted to write unknown token type %T", token)
	}
}

// --- Shared xref types and helpers ---

type xrefEntry struct {
	offset     int64
	generation int
}

type xrefSet struct {
	first   int64
	entries []*xrefEntry
}

func buildXrefSets(sorted []core.IndirectReference, offsets map[core.IndirectReference]int64) []xrefSet {
	var sets []xrefSet

	// Zero entry (object 0 is always free with generation 65535)
	currentEntries := []*xrefEntry{nil}
	var firstObjNum int64 = 0
	currentObjNum := int64(0)

	for _, ref := range sorted {
		step := ref.ObjectNumber() - currentObjNum
		if step == 1 {
			currentObjNum = ref.ObjectNumber()
			currentEntries = append(currentEntries, &xrefEntry{offset: offsets[ref], generation: ref.Generation()})
		} else {
			if len(currentEntries) > 0 {
				sets = append(sets, xrefSet{first: firstObjNum, entries: currentEntries})
			}
			currentEntries = []*xrefEntry{{offset: offsets[ref], generation: ref.Generation()}}
			currentObjNum = ref.ObjectNumber()
			firstObjNum = ref.ObjectNumber()
		}
	}

	if len(currentEntries) > 0 {
		sets = append(sets, xrefSet{first: firstObjNum, entries: currentEntries})
	}

	return sets
}
