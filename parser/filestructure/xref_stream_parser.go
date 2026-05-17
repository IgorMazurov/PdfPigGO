package filestructure

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/util"
)

// TryReadStreamAtOffset attempts to parse an xref stream starting at the given byte offset.
// Returns nil if no valid stream is found or if parsing fails.
func TryReadStreamAtOffset(fileHeaderOffset FileHeaderOffset, xrefOffset int64, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, log logging.Log) *XrefStream {
	if xrefOffset >= bytes.Length() || xrefOffset < 0 {
		return nil
	}

	offsetCorrectionType := XrefOffsetCorrectionNone
	offsetCorrection := int64(0)

	bytes.Seek(xrefOffset, 0) // SeekStart

	dictToken, ok := tryReadStreamObjAt(xrefOffset, scanner)
	if !ok || dictToken == nil {
		log.Debug(fmt.Sprintf("Did not find the stream at %d attempting correction", xrefOffset))
		recovered := tryRecoverOffset(fileHeaderOffset, xrefOffset, scanner)

		if recovered == nil {
			return nil
		}

		streamDict, ok2 := tryReadStreamObjAt(recovered.correctOffset, scanner)
		if !ok2 || streamDict == nil {
			return nil
		}

		dictToken = streamDict
		offsetCorrection = recovered.correctOffset - xrefOffset
		offsetCorrectionType = recovered.correctionType
		xrefOffset = recovered.correctOffset
	}

	// Check Type == /XRef
	typeToken, ok := dictToken.TryGet(tokens.Type)
	if !ok {
		return nil
	}
	nameToken, ok := typeToken.(*tokens.NameToken)
	if !ok || nameToken != tokens.Xref {
		return nil
	}

	// Get W array (field sizes)
	wToken, ok := dictToken.TryGet(tokens.W)
	if !ok {
		return nil
	}
	wArray, ok := wToken.(*tokens.ArrayToken)
	if !ok {
		return nil
	}

	streamData, err := readStreamTolerant(bytes)
	if err != nil || streamData.to == nil {
		return nil
	}

	dataLen := *streamData.to - streamData.from
	if dataLen <= 0 {
		return nil
	}

	bytes.Seek(streamData.from, 0) // SeekStart

	data := make([]byte, dataLen)
	readCount, readErr := bytes.Read(data)
	if readErr != nil || int64(readCount) != dataLen {
		return nil
	}

	streamToken, err := tokens.NewStreamToken(dictToken, data)
	if err != nil {
		return nil
	}

	decoded, err := decodeStream(streamToken, filters.Instance)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to parse the XRef stream at %d: %v", xrefOffset, err))
		return nil
	}

	fieldSizes, err := newXrefFieldSize(wArray)
	if err != nil {
		return nil
	}

	lineCount := len(decoded) / fieldSizes.lineLength
	objectNumbers := getObjectNumbers(dictToken)

	numbers := make([]xrefEntry, 0, len(objectNumbers))
	lineBuffer := make([]byte, fieldSizes.lineLength)

	for i, objNum := range objectNumbers {
		if i >= lineCount {
			break
		}

		byteOffset := i * fieldSizes.lineLength
		copy(lineBuffer, decoded[byteOffset:byteOffset+fieldSizes.lineLength])

		entryType := 0
		if fieldSizes.field1Size == 0 {
			entryType = 1
		} else {
			for j := 0; j < fieldSizes.field1Size; j++ {
				val := int(lineBuffer[j]) & 0x00ff
				entryType |= val << ((fieldSizes.field1Size - j - 1) * 8)
			}
		}

		readNextStreamObject(entryType, objNum, fieldSizes, &numbers, lineBuffer)
	}

	resultMap := make(map[core.IndirectReference]core.XrefLocation, len(numbers))
	for _, entry := range numbers {
		ref, err := core.NewIndirectReference(entry.obj, entry.gen)
		if err != nil {
			return nil
		}
		resultMap[ref] = entry.location
	}

	return NewXrefStream(xrefOffset, resultMap, dictToken, offsetCorrectionType, offsetCorrection)
}

// xrefEntry holds a parsed cross-reference stream entry.
type xrefEntry struct {
	obj      int64
	gen      int
	location core.XrefLocation
}

// streamBoundaries represents the start and optional end position of stream data.
type streamBoundaries struct {
	from int64
	to   *int64
}

// offsetRecovery holds the result of an xref offset recovery attempt.
type offsetRecovery struct {
	correctOffset  int64
	correctionType XrefOffsetCorrection
}

// tryRecoverOffset attempts to find a valid xref stream near the given offset.
// The provided offset can frequently be close but not quite correct.
func tryRecoverOffset(fileHeaderOffset FileHeaderOffset, xrefOffset int64, scanner tokenization.SeekableTokenScanner) *offsetRecovery {
	if fileHeaderOffset.Value > 0 {
		if _, ok := tryReadStreamObjAt(xrefOffset+int64(fileHeaderOffset.Value), scanner); ok {
			return &offsetRecovery{
				correctOffset:  xrefOffset + int64(fileHeaderOffset.Value),
				correctionType: XrefOffsetCorrectionFileHeaderOffset,
			}
		}
	}

	return nil
}

// tryReadStreamObjAt attempts to read a stream object (obj num gen obj << ... >>) at the given offset.
func tryReadStreamObjAt(offset int64, scanner tokenization.SeekableTokenScanner) (*tokens.DictionaryToken, bool) {
	scanner.Seek(offset, 0) // SeekStart

	if !scanner.Advance() {
		return nil, false
	}
	if _, ok := scanner.Current().(*tokens.NumericToken); !ok {
		return nil, false
	}

	if !scanner.Advance() {
		return nil, false
	}
	if _, ok := scanner.Current().(*tokens.NumericToken); !ok {
		return nil, false
	}

	if !scanner.Advance() {
		return nil, false
	}
	opToken, ok := scanner.Current().(*tokens.OperatorToken)
	if !ok || opToken != tokens.StartObject {
		return nil, false
	}

	if !scanner.Advance() {
		return nil, false
	}
	dictToken, ok := scanner.Current().(*tokens.DictionaryToken)
	if !ok {
		return nil, false
	}

	return dictToken, true
}

// readNextStreamObject parses a single xref stream entry and appends it to results.
func readNextStreamObject(entryType int, objectNumber int64, fieldSizes *xrefFieldSize, results *[]xrefEntry, lineBuffer []byte) {
	switch entryType {
	case 0:
		// Free objects are ignored.
	case 1:
		offset := readUnsigned(lineBuffer, fieldSizes.field1Size, fieldSizes.field2Size)
		genNum := readUnsigned(lineBuffer, fieldSizes.field1Size+fieldSizes.field2Size, fieldSizes.field3Size)

		if offset < 0 {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Location with negative offset %d found for object %d", offset, objectNumber)))
		}

		*results = append(*results, xrefEntry{objectNumber, int(genNum), core.File(offset)})

	case 2:
		objectStreamNumber := readUnsigned(lineBuffer, fieldSizes.field1Size, fieldSizes.field2Size)
		streamIndex := readUnsigned(lineBuffer, fieldSizes.field1Size+fieldSizes.field2Size, fieldSizes.field3Size)

		if objectStreamNumber < 0 {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Location with negative or zero object stream number %d found for object %d", objectStreamNumber, objectNumber)))
		}

		if streamIndex < 0 {
			panic(core.NewPdfDocumentFormatException(fmt.Sprintf("Location with negative stream index %d found for object %d in stream %d", streamIndex, objectNumber, objectStreamNumber)))
		}

		*results = append(*results, xrefEntry{objectNumber, 0, core.Stream(objectStreamNumber, int(streamIndex))})
	}
}

// readUnsigned reads an unsigned integer from the buffer starting at start with the given byte width.
func readUnsigned(buffer []byte, start, width int) int64 {
	value := int64(0)
	for i := 0; i < width; i++ {
		value <<= 8
		value |= int64(buffer[start+i])
	}
	return value
}

// xrefFieldSize represents the size of the fields in a cross reference stream.
type xrefFieldSize struct {
	// field1Size is the type of the entry (0, 1, or 2).
	field1Size int
	// field2Size: Type 0 and 2 is object number; Type 1 is byte offset from beginning of file.
	field2Size int
	// field3Size: For types 0 and 1 this is generation number; for type 2 it is stream index.
	field3Size int
	// lineLength is the total bytes in a single entry line.
	lineLength int
}

func newXrefFieldSize(wArray *tokens.ArrayToken) (*xrefFieldSize, error) {
	data := wArray.Data()
	if len(data) < 3 {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("There must be at least 3 entries in a W entry for a stream dictionary: %v", wArray))
	}

	numeric0, ok0 := wArray.Get(0).(*tokens.NumericToken)
	numeric1, ok1 := wArray.Get(1).(*tokens.NumericToken)
	numeric2, ok2 := wArray.Get(2).(*tokens.NumericToken)

	if !ok0 || !ok1 || !ok2 {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("W array entries must be numeric: %v", wArray))
	}

	field1 := numeric0.IntVal()
	field2 := numeric1.IntVal()
	field3 := numeric2.IntVal()

	return &xrefFieldSize{
		field1Size: int(field1),
		field2Size: int(field2),
		field3Size: int(field3),
		lineLength: int(field1 + field2 + field3),
	}, nil
}

// readStreamTolerant scans forward from the current position to find stream data boundaries.
func readStreamTolerant(bytes core.InputBytes) (streamBoundaries, error) {
	buffer := util.NewCircularByteBuffer(len("endstream "))

	startMarker := bytes.CurrentOffset()
	var endMarker *int64

	for bytes.CurrentByte() == '>' && bytes.MoveNext() {
	}

	isStreamWhitespace := func() bool {
		b := bytes.CurrentByte()
		return b == ' ' || b == '\r' || b == '\n'
	}

	isWhitespaceActive := isStreamWhitespace()

	for {
		if isStreamWhitespace() {
			buffer.Add(' ')
			if isWhitespaceActive {
				startMarker = bytes.CurrentOffset()
			}
		} else {
			buffer.Add(bytes.CurrentByte())
			isWhitespaceActive = false
		}

		if buffer.EndsWith("endstream ") {
			val := bytes.CurrentOffset() - int64(len("endstream "))
			endMarker = &val
			break
		}

		if buffer.EndsWith("stream ") {
			startMarker = bytes.CurrentOffset()
			isWhitespaceActive = isStreamWhitespace()
		} else if buffer.EndsWith("endobj ") {
			val := bytes.CurrentOffset() - int64(len("endobj "))
			endMarker = &val
			break
		}

		if !bytes.MoveNext() {
			break
		}
	}

	return streamBoundaries{from: startMarker, to: endMarker}, nil
}

// getObjectNumbers returns the list of object numbers covered by this xref stream.
func getObjectNumbers(dictionary *tokens.DictionaryToken) []int64 {
	sizeToken, ok := dictionary.TryGet(tokens.Size)
	if !ok {
		panic(core.NewPdfDocumentFormatException(fmt.Sprintf("The stream dictionary must contain a numeric size value: %v", dictionary)))
	}
	sizeNumeric, ok := sizeToken.(*tokens.NumericToken)
	if !ok {
		panic(core.NewPdfDocumentFormatException(fmt.Sprintf("The stream dictionary must contain a numeric size value: %v", dictionary)))
	}

	objNums := make([]int64, 0)

	indexToken, hasIndex := dictionary.TryGet(tokens.Index)
	if hasIndex {
		if indexArray, ok := indexToken.(*tokens.ArrayToken); ok {
			for i := 0; i < indexArray.Length(); i += 2 {
				firstObjNumeric, ok1 := indexArray.Get(i).(*tokens.NumericToken)
				sizeNumericInner, ok2 := indexArray.Get(i+1).(*tokens.NumericToken)
				if !ok1 || !ok2 {
					continue
				}
				firstObjectNumber := int64(firstObjNumeric.IntVal())
				subsetSize := int64(sizeNumericInner.IntVal())
				for j := int64(0); j < subsetSize; j++ {
					objNums = append(objNums, firstObjectNumber+j)
				}
			}
		}
	} else {
		for i := int64(0); i < int64(sizeNumeric.IntVal()); i++ {
			objNums = append(objNums, i)
		}
	}

	return objNums
}

// decodeStream decodes the given stream token using the provided filter provider.
func decodeStream(stream *tokens.StreamToken, provider filters.FilterProvider) ([]byte, error) {
	filterList, err := provider.GetFilters(stream.StreamDictionary)
	if err != nil {
		return nil, err
	}

	transform := stream.Data()

	for i, filter := range filterList {
		decoded, decodeErr := filter.Decode(transform, stream.StreamDictionary, provider, i)
		if decodeErr != nil {
			return nil, fmt.Errorf("filter %d failed: %w", i, decodeErr)
		}
		transform = decoded
	}

	return transform, nil
}
