// Package writer provides types for writing PDF documents.
package writer

import (
	"bytes"
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// NoTextTokenWriter is a TokenWriter that removes ShowText, ShowTextsWithPositioning,
// and MoveToNextLineShowText operations from page content streams.
type NoTextTokenWriter struct {
	writingPageContents bool
	page                int
}

// NewNoTextTokenWriter creates a new NoTextTokenWriter.
func NewNoTextTokenWriter() *NoTextTokenWriter {
	return &NoTextTokenWriter{}
}

// Page gets or sets the current page number for log messages.
func (w *NoTextTokenWriter) Page() int {
	return w.page
}

// SetPage sets the current page number.
func (w *NoTextTokenWriter) SetPage(p int) {
	w.page = p
}

// WritingPageContents reports whether the writer is currently writing page contents.
func (w *NoTextTokenWriter) WritingPageContents() bool {
	return w.writingPageContents
}

// SetWritingPageContents sets whether the writer is currently writing page contents.
func (w *NoTextTokenWriter) SetWritingPageContents(v bool) {
	w.writingPageContents = v
}

// WriteToken writes the given token to the output stream, filtering text operations
// from page content streams.
func (w *NoTextTokenWriter) WriteToken(token tokens.Token, out io.Writer) error {
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
		if streamTok, ok := t.Data().(*tokens.StreamToken); ok {
			return w.writeObjectWithStreamFilter(t, streamTok, out)
		}
		return writeObjectToken(t, out)
	case *tokens.StreamToken:
		return w.writeStream(t, out)
	case *tokens.StringToken:
		return writeString(t, out)
	default:
		return fmt.Errorf("attempted to write unknown token type %T", token)
	}
}

// WriteObject writes pre-serialized data as an object token.
func (w *NoTextTokenWriter) WriteObject(objectNumber int64, generation int, data []byte, out io.Writer) error {
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

// WriteCrossReferenceTable writes a cross-reference table and trailer dictionary.
func (w *NoTextTokenWriter) WriteCrossReferenceTable(
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
	sortByObjectNumber(sorted)

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
		tokens.Id:   tokens.NewArrayToken([]tokens.Token{}),
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

// sortByObjectNumber sorts indirect references by object number using bubble sort.
func sortByObjectNumber(refs []core.IndirectReference) {
	for i := 0; i < len(refs)-1; i++ {
		for j := i + 1; j < len(refs); j++ {
			if refs[i].ObjectNumber() > refs[j].ObjectNumber() {
				refs[i], refs[j] = refs[j], refs[i]
			}
		}
	}
}

// writeStream writes a stream token to the output, filtering text operations when appropriate.
func (w *NoTextTokenWriter) writeStream(streamToken *tokens.StreamToken, out io.Writer) error {
	var outputStreamToken *tokens.StreamToken

	if !w.writingPageContents && !isFormStream(streamToken) {
		outputStreamToken = streamToken
	} else {
		result, ok := w.tryGetStreamWithoutText(streamToken)
		if !ok {
			outputStreamToken = streamToken
		} else {
			outputStreamToken = result
		}
	}

	if err := writeDictionary(outputStreamToken.StreamDictionary, out); err != nil {
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
	if _, err := out.Write(outputStreamToken.Data()); err != nil {
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

// writeObjectWithStreamFilter writes an ObjectToken that wraps a StreamToken,
// using the filtering writeStream for the inner content so text operations
// are properly removed from page content streams.
func (w *NoTextTokenWriter) writeObjectWithStreamFilter(obj *tokens.ObjectToken, streamTok *tokens.StreamToken, out io.Writer) error {
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
	if err := w.writeStream(streamTok, out); err != nil {
		return err
	}
	if _, err := out.Write(streamObjEnd); err != nil {
		return err
	}
	return writeLineBreak(out)
}

// isFormStream checks whether the stream is a Form XObject.
func isFormStream(streamToken *tokens.StreamToken) bool {
	subtype, ok := streamToken.StreamDictionary.TryGet(tokens.Subtype)
	if !ok {
		return false
	}
	nameToken, ok := subtype.(*tokens.NameToken)
	if !ok {
		return false
	}
	return nameToken.Equals(tokens.Form)
}

// tryGetStreamWithoutText decodes the stream, parses its content operations,
// removes text-showing operations, re-encodes and returns a new StreamToken.
// Returns (nil, false) if no text operations were found or an error occurred.
func (w *NoTextTokenWriter) tryGetStreamWithoutText(streamToken *tokens.StreamToken) (*tokens.StreamToken, bool) {
	filterProvider := filters.NewFilterProviderWithLookup(filters.Instance)

	decodedBytes, err := decodeStreamData(streamToken, filterProvider)
	if err != nil {
		return nil, false
	}

	pageContentParser := parser.NewPageContentParser(
		graphics.GetReflectionFactory(),
		core.Infinite,
		false,
	)

	operations := pageContentParser.Parse(w.page, core.NewMemoryInputBytes(decodedBytes), logging.NoopLog)

	var outputBuffer bytes.Buffer
	haveText := false

	for i := 0; i < len(operations); i++ {
		op := operations[i]
		if isTextOperation(op) {
			haveText = true
			continue
		}
		if writer, ok := op.(interface{ Write(io.Writer) error }); ok {
			if err := writer.Write(&outputBuffer); err != nil {
				return nil, false
			}
		}
	}

	if !haveText {
		return nil, false
	}

	compressedBytes := CompressBytes(outputBuffer.Bytes())

	outputDict := map[*tokens.NameToken]tokens.Token{
		tokens.Length: tokens.NewNumericTokenFromInt(len(compressedBytes)),
		tokens.Filter: tokens.FlateDecode,
	}

	for keyStr, value := range streamToken.StreamDictionary.Data() {
		key := tokens.Create(keyStr)
		if _, exists := outputDict[key]; !exists {
			outputDict[key] = value
		}
	}

	dictToken, err := tokens.NewDictionary(outputDict)
	if err != nil {
		return nil, false
	}

	streamTokenResult, err := tokens.NewStreamToken(dictToken, compressedBytes)
	if err != nil {
		return nil, false
	}

	return streamTokenResult, true
}

// isTextOperation checks if the given operation is a text-showing operation.
func isTextOperation(op content.GraphicsStateOperation) bool {
	switch op.(type) {
	case *graphics.ShowText, graphics.ShowText:
		return true
	case *graphics.ShowTextsWithPositioning, graphics.ShowTextsWithPositioning:
		return true
	case *graphics.MoveToNextLineShowText, graphics.MoveToNextLineShowText:
		return true
	default:
		return false
	}
}

// decodeStreamData decodes the raw stream data by applying filters.
func decodeStreamData(stream *tokens.StreamToken, provider filters.LookupFilterProvider) ([]byte, error) {
	filterNames, err := getFilterNames(stream.StreamDictionary)
	if err != nil {
		return nil, err
	}

	if len(filterNames) == 0 {
		return stream.Data(), nil
	}

	filterInstances, err := provider.GetNamedFilters(filterNames)
	if err != nil {
		return nil, err
	}

	data := stream.Data()
	for i, f := range filterInstances {
		if !f.IsSupported() {
			continue
		}
		decoded, err := f.Decode(data, stream.StreamDictionary, provider, i)
		if err != nil {
			return nil, fmt.Errorf("failed to decode with filter %T at index %d: %w", f, i, err)
		}
		data = decoded
	}

	return data, nil
}

// getFilterNames extracts the filter name tokens from the stream dictionary.
func getFilterNames(dict *tokens.DictionaryToken) ([]*tokens.NameToken, error) {
	filterToken, ok := dict.TryGet(tokens.Filter)
	if !ok {
		filterToken, ok = dict.TryGet(tokens.FFilter)
		if !ok {
			return nil, nil
		}
	}

	switch t := filterToken.(type) {
	case *tokens.NameToken:
		return []*tokens.NameToken{t}, nil
	case *tokens.ArrayToken:
		var result []*tokens.NameToken
		for _, item := range t.Data() {
			if name, ok := item.(*tokens.NameToken); ok {
				result = append(result, name)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unexpected filter token type: %T", filterToken)
	}
}
