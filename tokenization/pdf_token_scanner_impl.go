package tokenization

import (
	"fmt"
	"log"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// encryptionHandler defines the interface for decrypting PDF tokens.
type encryptionHandler interface {
	Decrypt(reference core.IndirectReference, token tokens.Token) tokens.Token
}

// StreamDecoder decodes a stream's compressed data into raw bytes.
type StreamDecoder interface {
	DecodeStream(stream *tokens.StreamToken) []byte
}

// streamDecoder is the unexported alias used internally.
type streamDecoder = StreamDecoder

var endstreamBytes = []byte("endstream")

// PdfTokenScannerImpl scans PDF file bytes for indirect objects.
// It wraps a CoreTokenScanner and handles object resolution, stream reading,
// encryption decryption, and brute-force fallback when xref data is missing.
type PdfTokenScannerImpl struct {
	inputBytes           core.InputBytes
	objectLocationProvider ObjectLocationProvider
	encryptionHandler    encryptionHandler
	streamDecoder        streamDecoder
	fileHeaderOffset     int // byte offset of the "%PDF-" header in the file
	parsingOptions       *scannerParsingOptions
	stackDepthGuard      *core.StackDepthGuard

	coreScanner       *CoreTokenScanner
	currentToken      tokens.Token
	isDisposed        bool
	isBruteForcing    bool
	callingObject     *core.IndirectReference
	overwrittenTokens map[core.IndirectReference]*tokens.ObjectToken

	readTokens             []tokens.Token
	previousTokens         [3]tokens.Token
	previousTokenPositions [3]int64
}

// ScannerParsingOptions holds lenient parsing flag used by the scanner.
type ScannerParsingOptions struct {
	UseLenientParsing bool
}

// scannerParsingOptions is the unexported alias for internal use.
type scannerParsingOptions = ScannerParsingOptions

// NewPdfTokenScanner creates a fully-featured PdfTokenScanner that scans for
// indirect objects in a PDF file, handling streams, encryption, and xref lookups.
func NewPdfTokenScanner(
	inputBytes core.InputBytes,
	objectLocationProvider ObjectLocationProvider,
	encryptionHandler encryptionHandler,
	streamDecoder streamDecoder,
	fileHeaderOffset int,
	parsingOptions *scannerParsingOptions,
	stackDepthGuard *core.StackDepthGuard,
) *PdfTokenScannerImpl {
	if parsingOptions == nil {
		parsingOptions = &scannerParsingOptions{}
	}

	return &PdfTokenScannerImpl{
		inputBytes:           inputBytes,
		objectLocationProvider: objectLocationProvider,
		encryptionHandler:  encryptionHandler,
		streamDecoder:      streamDecoder,
		fileHeaderOffset:   fileHeaderOffset,
		parsingOptions:     parsingOptions,
		stackDepthGuard:    stackDepthGuard,
		coreScanner:        NewCoreTokenScanner(inputBytes, true, stackDepthGuard, ScannerScopeNone, nil, parsingOptions.UseLenientParsing, false),
		overwrittenTokens:  make(map[core.IndirectReference]*tokens.ObjectToken),
	}
}

// UpdateEncryptionHandler replaces the encryption handler used for decrypting tokens.
func (s *PdfTokenScannerImpl) UpdateEncryptionHandler(handler encryptionHandler) {
	if handler == nil {
		panic(fmt.Errorf("encryption handler cannot be nil"))
	}
	s.encryptionHandler = handler
}

// Advance moves to the next indirect object in the PDF file.
// Returns true if an object was successfully read, false otherwise.
func (s *PdfTokenScannerImpl) Advance() bool {
	if s.isDisposed {
		panic(fmt.Errorf("PdfTokenScanner has been disposed"))
	}

	tokensRead := 0
	for s.coreScanner.Advance() && !s.isOperatorToken(s.coreScanner.Current(), tokens.StartObject) {
		if _, ok := s.coreScanner.Current().(*tokens.CommentToken); ok {
			continue
		}

		tokensRead++

		s.previousTokens[0] = s.previousTokens[1]
		s.previousTokenPositions[0] = s.previousTokenPositions[1]

		s.previousTokens[1] = s.previousTokens[2]
		s.previousTokenPositions[1] = s.previousTokenPositions[2]

		s.previousTokens[2] = s.coreScanner.Current()
		s.previousTokenPositions[2] = s.coreScanner.CurrentTokenStart()
	}

	if tokensRead < 2 {
		return false
	}

	startPosition := s.previousTokenPositions[1]
	objectNumber, ok1 := s.previousTokens[1].(*tokens.NumericToken)
	generation, ok2 := s.previousTokens[2].(*tokens.NumericToken)

	if !ok1 || !ok2 {
		if _, okGen := s.previousTokens[2].(*tokens.NumericToken); okGen {
			if opTok, okOp := s.previousTokens[1].(*tokens.OperatorToken); okOp {
				numStr := extractTrailingNumber(opTok.Data())
				if numStr != "" {
					var number int64
					for i := 0; i < len(numStr); i++ {
						number = number*10 + int64(numStr[i]-'0')
					}
					idx := strings.Index(opTok.Data(), numStr)
					startPosition = s.previousTokenPositions[1] + int64(idx)
					objectNumber = tokens.NewNumericTokenFromInt(int(number))
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	}

	readStream := false

	for s.coreScanner.Advance() && s.isToken(s.coreScanner, tokens.EndObject) == nil {
		if _, ok := s.coreScanner.Current().(*tokens.CommentToken); ok {
			continue
		}

		if s.isOperatorToken(s.coreScanner.Current(), tokens.StartObject) {
			if readStream && len(s.readTokens) > 0 {
				if streamRead, ok := s.readTokens[0].(*tokens.StreamToken); ok {
					s.readTokens = s.readTokens[:0]
					s.readTokens = append(s.readTokens, streamRead)
					s.seekTo(s.previousTokenPositions[0])
					break
				}
			}

			if len(s.readTokens) == 3 {
				if _, okA := s.readTokens[1].(*tokens.NumericToken); okA {
					if _, okB := s.readTokens[2].(*tokens.NumericToken); okB {
						// C# uses original objectNumber and generation, not the extra ones from readTokens.
						ref, _ := core.NewIndirectReference(objectNumber.LongVal(), generation.IntVal())
						actualToken := s.encryptionHandler.Decrypt(ref, s.readTokens[0])
						s.currentToken = tokens.NewObjectToken(core.File(startPosition), ref, actualToken)
						s.readTokens = s.readTokens[:0]
						s.seekTo(s.previousTokenPositions[0])
						return true
					}
				}
			}

			return false
		}

		if s.isToken(s.coreScanner, tokens.OpXref) != nil || s.isToken(s.coreScanner, tokens.StartXref) != nil {
			if readStream && len(s.readTokens) > 0 {
				if streamRead, ok := s.readTokens[0].(*tokens.StreamToken); ok {
					s.readTokens = s.readTokens[:0]
					s.readTokens = append(s.readTokens, streamRead)
					s.seekTo(s.previousTokenPositions[2])
					break
				}
			}

			if len(s.readTokens) == 1 {
				ref, _ := core.NewIndirectReference(objectNumber.LongVal(), generation.IntVal())
				actualToken := s.encryptionHandler.Decrypt(ref, s.readTokens[0])
				s.currentToken = tokens.NewObjectToken(core.File(startPosition), ref, actualToken)
				s.readTokens = s.readTokens[:0]
				s.seekTo(s.previousTokenPositions[2])
				return true
			}

			return false
		}

		actualStartStreamPos := s.isToken(s.coreScanner, tokens.StartStream)
		if actualStartStreamPos != nil {
			streamIdentifier, _ := core.NewIndirectReference(objectNumber.LongVal(), generation.IntVal())

			getLengthFromFile := !s.isBruteForcing && !(s.callingObject != nil && *s.callingObject == streamIdentifier)

			outerCallingObject := s.callingObject
			s.callingObject = &streamIdentifier

			stream, err := s.tryReadStream(*actualStartStreamPos, getLengthFromFile)
			s.callingObject = outerCallingObject

			if err == nil && stream != nil {
				s.readTokens = s.readTokens[:0]
				s.readTokens = append(s.readTokens, stream)
				readStream = true
				break
			}
		} else {
			s.readTokens = append(s.readTokens, s.coreScanner.Current())
		}

		s.previousTokens[0] = s.previousTokens[1]
		s.previousTokenPositions[0] = s.previousTokenPositions[1]

		s.previousTokens[1] = s.previousTokens[2]
		s.previousTokenPositions[1] = s.previousTokenPositions[2]

		s.previousTokens[2] = s.coreScanner.Current()
		s.previousTokenPositions[2] = s.coreScanner.CurrentTokenStart()
	}

	if !readStream && s.isToken(s.coreScanner, tokens.EndObject) == nil {
		s.readTokens = s.readTokens[:0]
		return false
	}

	ref, _ := core.NewIndirectReference(objectNumber.LongVal(), generation.IntVal())

	var token tokens.Token
	token = s.resolveReadTokens()

	token = s.encryptionHandler.Decrypt(ref, token)

	s.currentToken = tokens.NewObjectToken(core.File(startPosition), ref, token)

	s.objectLocationProvider.UpdateOffset(ref, core.File(startPosition))

	s.readTokens = s.readTokens[:0]
	return true
}

func (s *PdfTokenScannerImpl) resolveReadTokens() tokens.Token {
	// Handle the case: objNum genNum R -> IndirectReferenceToken
	if len(s.readTokens) == 3 {
		objNum, okA := s.readTokens[0].(*tokens.NumericToken)
		genNum, okB := s.readTokens[1].(*tokens.NumericToken)
		rTok, okC := s.readTokens[2].(*tokens.OperatorToken)
		if okA && okB && okC && rTok.Data() == "R" {
			ref, _ := core.NewIndirectReference(objNum.LongVal(), genNum.IntVal())
			return tokens.NewIndirectReferenceToken(ref)
		}
	}

	if len(s.readTokens) > 1 {
		log.Printf("Found more than 1 token in an object.")

		trimmed := make([]tokens.Token, 0, len(s.readTokens))
		for _, t := range s.readTokens {
			if op, ok := t.(*tokens.OperatorToken); !ok || (op.Data() != ">" && op.Data() != "]") {
				trimmed = append(trimmed, t)
			}
		}

		if len(trimmed) == 1 {
			return trimmed[0]
		}

		if str, ok := s.readTokens[0].(*tokens.StreamToken); ok {
			allEndStream := true
			for i := 1; i < len(s.readTokens); i++ {
				op, okOp := s.readTokens[i].(*tokens.OperatorToken)
				if !okOp || op.Data() != "endstream" {
					allEndStream = false
					break
				}
			}
			if allEndStream {
				return str
			}
		}

		return s.readTokens[len(s.readTokens)-1]
	}

	if len(s.readTokens) == 0 {
		return nil
	}

	return s.readTokens[len(s.readTokens)-1]
}

// Current returns the current token.
func (s *PdfTokenScannerImpl) Current() tokens.Token {
	return s.currentToken
}

// StackDepthGuard returns the stack depth guard.
func (s *PdfTokenScannerImpl) StackDepthGuard() *core.StackDepthGuard {
	return s.stackDepthGuard
}

// CurrentPosition returns the current byte position.
func (s *PdfTokenScannerImpl) CurrentPosition() int64 {
	return s.coreScanner.CurrentPosition()
}

// Length returns the total length of the input.
func (s *PdfTokenScannerImpl) Length() int64 {
	return s.coreScanner.Length()
}

// Seek moves to an absolute position in the input.
func (s *PdfTokenScannerImpl) Seek(offset int64, whence int) (int64, error) {
	return s.coreScanner.Seek(offset, whence)
}

// seekTo moves to an absolute position internally.
func (s *PdfTokenScannerImpl) seekTo(position int64) {
	s.coreScanner.Seek(position, 0) // io.SeekStart
}

// Close releases resources held by the scanner.
func (s *PdfTokenScannerImpl) Close() error {
	var err error
	if s.inputBytes != nil {
		err = s.inputBytes.Close()
	}
	s.isDisposed = true
	return err
}

// Get returns the tokenized object with the given indirect reference.
func (s *PdfTokenScannerImpl) Get(reference core.IndirectReference) *tokens.ObjectToken {
	navSet := make([]int, 7)
	return s.getRecursive(reference, navSet, 0)
}

func (s *PdfTokenScannerImpl) getRecursive(reference core.IndirectReference, navSet []int, depth byte) *tokens.ObjectToken {
	if depth >= byte(len(navSet)) {
		chain := make([]string, len(navSet))
		for i, v := range navSet {
			chain[i] = fmt.Sprintf("%d", v)
		}
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Deep object chain detected when looking for %s: %s.", reference, strings.Join(chain, ", "))))
	}

	for i := 0; i < int(depth); i++ {
		if navSet[i] == int(reference.ObjectNumber()) {
			chain := make([]string, len(navSet))
			for j, v := range navSet {
				chain[j] = fmt.Sprintf("%d", v)
			}
			panic(core.NewPdfDocumentFormatException(
				fmt.Sprintf("Circular reference encountered when looking for object %s. Involved objects were: %s", reference, strings.Join(chain, ", "))))
		}
	}

	navSet[depth] = int(reference.ObjectNumber())
	depth++

	if s.isDisposed {
		panic(fmt.Errorf("PdfTokenScanner has been disposed"))
	}

	if value, ok := s.overwrittenTokens[reference]; ok {
		return value
	}

	if objectToken, ok := s.objectLocationProvider.TryGetCached(reference); ok {
		return objectToken
	}

	offset, found := s.objectLocationProvider.TryGetOffset(reference)
	if !found {
		return nil
	}

	if offset.Type == core.XrefEntryTypeObjectStream {
		if offset.Value1 == reference.ObjectNumber() {
			panic(core.NewPdfDocumentFormatException(
				fmt.Sprintf("Object stream cannot contain itself, looking for object %s in %d", reference, offset.Value1)))
		}

		return s.getObjectFromStream(reference, offset, navSet, depth)
	}

	s.seekTo(offset.Value1)

	if !s.Advance() {
		bfResult := s.tryBruteForceFileToFindReference(reference)
		return bfResult
	}

	foundObj, ok := s.currentToken.(*tokens.ObjectToken)
	if !ok {
		return s.tryBruteForceFileToFindReference(reference)
	}

	if foundObj.Number().Equals(reference) {
		return foundObj
	}

	return s.tryBruteForceFileToFindReference(reference)
}

// ReplaceToken adds a token to the internal cache.
func (s *PdfTokenScannerImpl) ReplaceToken(reference core.IndirectReference, token tokens.Token) {
	s.overwrittenTokens[reference] = tokens.NewObjectToken(core.File(0), reference, token)
}

// RegisterCustomTokenizer adds support for a custom tokenizer.
func (s *PdfTokenScannerImpl) RegisterCustomTokenizer(firstByte byte, t Tokenizer) {
	s.coreScanner.RegisterCustomTokenizer(firstByte, t)
}

// DeregisterCustomTokenizer removes a previously registered custom tokenizer.
func (s *PdfTokenScannerImpl) DeregisterCustomTokenizer(t Tokenizer) {
	s.coreScanner.DeregisterCustomTokenizer(t)
}

// isOperatorToken checks if token equals the given operator singleton by data comparison.
func (s *PdfTokenScannerImpl) isOperatorToken(token tokens.Token, op *tokens.OperatorToken) bool {
	if token == nil || op == nil {
		return false
	}
	tok, ok := token.(*tokens.OperatorToken)
	return ok && tok.Data() == op.Data()
}

// isToken checks if the scanner's current token matches the given operator.
// Returns a pointer to the actual start position if matched, nil otherwise.
func (s *PdfTokenScannerImpl) isToken(scanner *CoreTokenScanner, token *tokens.OperatorToken) *int64 {
	current := scanner.Current()
	if opTok, ok := current.(*tokens.OperatorToken); ok && opTok.Data() == token.Data() {
		pos := scanner.CurrentTokenStart()
		return &pos
	}

	if s.parsingOptions.UseLenientParsing {
		if opTok, ok := current.(*tokens.OperatorToken); ok {
			tokenData := token.Data()
			opData := opTok.Data()
			if strings.HasSuffix(opData, tokenData) {
				pos := scanner.CurrentTokenStart() + int64(len(opData)-len(tokenData))
				return &pos
			}
		}
	}

	return nil
}

// tryReadStream attempts to read a PDF stream starting at the given offset.
func (s *PdfTokenScannerImpl) tryReadStream(startStreamTokenOffset int64, getLength bool) (*tokens.StreamToken, error) {
	streamDict := s.getStreamDictionary()

	var length *int64
	if getLength {
		length = s.getStreamLength(streamDict)
	}

	if !getLength {
		lenTok, ok := streamDict.TryGet(tokens.Create("Length"))
		if ok {
			if numTok, okNum := lenTok.(*tokens.NumericToken); okNum {
				l := numTok.LongVal()
				length = &l
			}
		}
	}

	if !readStreamTokenStart(s.inputBytes, startStreamTokenOffset) {
		return nil, fmt.Errorf("stream token start not found")
	}

	for {
		if !s.inputBytes.MoveNext() {
			return nil, fmt.Errorf("unexpected end of data while reading stream")
		}

		c := s.inputBytes.CurrentByte()
		if c == '\r' {
			if peek, ok := s.inputBytes.Peek(); ok && peek == '\n' {
				s.inputBytes.MoveNext()
			}
			break
		}
		if c == '\n' {
			break
		}
	}

	startDataOffset := s.inputBytes.CurrentOffset()

	read := int64(0)
	endObjPosition := 0
	endStreamPosition := 0
	commonPartPosition := 0

	const endWordPart = "end"
	const streamPart = "stream"
	const objPart = "obj"

	if data, ok := tryReadUsingLength(s.inputBytes, length, startDataOffset); ok {
		stream, _ := tokens.NewStreamToken(streamDict, data)
		return stream, nil
	}

	streamDataStart := s.inputBytes.CurrentOffset()

	var possibleEndLocation *PossibleStreamEndLocation

	for s.inputBytes.MoveNext() {
		if length != nil && read == *length {
			// continue reading past expected length to verify end marker
		}

		currentByte := s.inputBytes.CurrentByte()

		if commonPartPosition < len(endWordPart) && currentByte == endWordPart[commonPartPosition] {
			commonPartPosition++
		} else if commonPartPosition == len(endWordPart) {
			if currentByte == streamPart[endStreamPosition] {
				endObjPosition = 0
				endStreamPosition++

				if endStreamPosition == len(streamPart) {
					nextOk := s.inputBytes.MoveNext()
					isEnd := !nextOk || core.IsWhitespace(s.inputBytes.CurrentByte())
					if isEnd {
						offset := s.inputBytes.CurrentOffset() - int64(len(tokens.EndStream.Data()))
						possibleEndLocation = &PossibleStreamEndLocation{
							Offset: offset,
							Type:   tokens.EndStream,
						}

						if length != nil && read > *length {
							break
						}

						endStreamPosition = 0
						endObjPosition = 0
					} else {
						// "endstream" followed by non-whitespace — not a valid end marker.
						// Reset only endStreamPosition; keep commonPartPosition for potential next match.
						endStreamPosition = 0
						endObjPosition = 0
					}
				}
			} else if currentByte == objPart[endObjPosition] {
				endStreamPosition = 0
				endObjPosition++

				if endObjPosition == len(objPart) {
					endStreamPosition = 0
					commonPartPosition = 0

					if possibleEndLocation != nil {
						lastEnd := *possibleEndLocation
						s.inputBytes.Seek(lastEnd.Offset+int64(len(lastEnd.Type.Data()))+1, 0)
						break
					}

					offset := s.inputBytes.CurrentOffset() - int64(len(tokens.EndObject.Data()))
					possibleEndLocation = &PossibleStreamEndLocation{
						Offset: offset,
						Type:   tokens.EndObject,
					}

					if read > 0 && length != nil && read > *length {
						break
					}
				}
			} else {
				// "end" followed by neither "stream" nor "obj" — reset everything.
				endStreamPosition = 0
				endObjPosition = 0
				commonPartPosition = 0
			}
		} else {
			endStreamPosition = 0
			endObjPosition = 0
			if currentByte == endWordPart[0] {
				commonPartPosition = 1
			} else {
				commonPartPosition = 0
			}
		}

		read++
	}

	if possibleEndLocation == nil {
		return nil, fmt.Errorf("no stream end location found")
	}

	lastEnd := *possibleEndLocation

	dataLength := lastEnd.Offset - startDataOffset

	if dataLength <= 0 {
		return nil, fmt.Errorf("invalid stream data length: %d", dataLength)
	}

	s.inputBytes.Seek(lastEnd.Offset-3, 0)
	s.inputBytes.MoveNext()

	if s.inputBytes.CurrentByte() == '\r' {
		dataLength -= 3
	} else {
		dataLength -= 2
	}

	if dataLength < 0 {
		return nil, fmt.Errorf("invalid stream data length after newline adjustment: %d", dataLength)
	}

	data := make([]byte, dataLength)

	s.inputBytes.Seek(streamDataStart, 0)
	count, _ := s.inputBytes.Read(data)
	if count != len(data) {
		return nil, fmt.Errorf("failed to read stream data: wanted %d, got %d", len(data), count)
	}

	streamDataEnd := s.inputBytes.CurrentOffset() + 1
	s.inputBytes.Seek(streamDataEnd, 0)

	stream, _ := tokens.NewStreamToken(streamDict, data)
	return stream, nil
}

// tryReadUsingLength attempts to read stream data using the known length value.
func tryReadUsingLength(inputBytes core.InputBytes, length *int64, startDataOffset int64) ([]byte, bool) {
	if length == nil || *length+startDataOffset >= inputBytes.Length() {
		return nil, false
	}

	readBuffer := make([]byte, len(endstreamBytes))

	newlineCount := 0

	inputBytes.Seek(*length+startDataOffset, 0)

	if next, ok := inputBytes.Peek(); ok {
		if core.IsEndOfLineByte(next) {
			newlineCount++
			inputBytes.MoveNext()

			if next2, ok2 := inputBytes.Peek(); ok2 {
				if core.IsEndOfLineByte(next2) {
					newlineCount++
					inputBytes.MoveNext()
				}
			}
		}
	}

	readLen, _ := inputBytes.Read(readBuffer)

	if readLen != len(readBuffer) {
		return nil, false
	}

	for i := 0; i < len(endstreamBytes); i++ {
		if readBuffer[i] != endstreamBytes[i] {
			inputBytes.Seek(startDataOffset, 0)
			return nil, false
		}
	}

	inputBytes.Seek(startDataOffset, 0)

	data := make([]byte, *length)

	countRead, _ := inputBytes.Read(data)

	if countRead != len(data) {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Reading using the stream length failed to read as many bytes as the stream specified. Wanted %d, got %d at %d.",
				*length, countRead, startDataOffset+1)))
	}

	inputBytes.Read(readBuffer)

	for i := 0; i < newlineCount; i++ {
		if !inputBytes.MoveNext() {
			inputBytes.Seek(startDataOffset, 0)
			return nil, false
		}
	}

	inputBytes.MoveNext() // skip past 'm' in 'endstream'

	return data, true
}

// getStreamDictionary retrieves the dictionary token that precedes a stream.
func (s *PdfTokenScannerImpl) getStreamDictionary() *tokens.DictionaryToken {
	if firstDict, ok := s.previousTokens[2].(*tokens.DictionaryToken); ok {
		return firstDict
	}
	if secondDict, ok := s.previousTokens[1].(*tokens.DictionaryToken); ok {
		return secondDict
	}

	panic(core.NewPdfDocumentFormatException(
		fmt.Sprintf("No dictionary token was found prior to the 'stream' operator. Previous tokens were: %v and %v.",
			s.previousTokens[2], s.previousTokens[1])))
}

// getStreamLength determines the stream length from the dictionary.
func (s *PdfTokenScannerImpl) getStreamLength(dictionary *tokens.DictionaryToken) *int64 {
	lengthValue, found := dictionary.TryGet(tokens.Create("Length"))
	if !found {
		return nil
	}

if numeric, ok := lengthValue.(*tokens.NumericToken); ok {
			l := numeric.LongVal()
			return &l
		}

	currentOffset := s.inputBytes.CurrentOffset()

	if lengthReference, ok := lengthValue.(*tokens.IndirectReferenceToken); ok {
		offset, found := s.objectLocationProvider.TryGetOffset(lengthReference.Data())
		if !found {
			return nil
		}

		if offset.Type == core.XrefEntryTypeObjectStream {
			navSet := make([]int, 7)
			result := s.getObjectFromStream(lengthReference.Data(), offset, navSet, 0)
			if result == nil {
				panic(core.NewPdfDocumentFormatException(
					fmt.Sprintf("Could not locate the length object with offset %+v which should have been in a stream.", offset)))
			}
			if streamLengthToken, ok := result.Data().(*tokens.NumericToken); ok {
				l := streamLengthToken.LongVal()
				return &l
			}
			panic(core.NewPdfDocumentFormatException(
				fmt.Sprintf("Could not locate the length object with offset %+v which should have been in a stream. Found: %v.", offset, result.Data())))
		}

		s.seekTo(offset.Value1)

		oldReadTokens := make([]tokens.Token, len(s.readTokens))
		copy(oldReadTokens, s.readTokens)
		s.readTokens = s.readTokens[:0]

		oldPrevTokens := [3]tokens.Token{}
		copy(oldPrevTokens[:], s.previousTokens[:])
		oldPrevPositions := [3]int64{}
		copy(oldPrevPositions[:], s.previousTokenPositions[:])
		oldCurrentToken := s.currentToken

		var length *int64
		if s.Advance() {
			if objTok, ok := s.currentToken.(*tokens.ObjectToken); ok {
				if lengthToken, okNum := objTok.Data().(*tokens.NumericToken); okNum {
					l := lengthToken.LongVal()
					length = &l
				}
			}
		}

		s.readTokens = append(s.readTokens, oldReadTokens...)
		s.previousTokens = oldPrevTokens
		s.previousTokenPositions = oldPrevPositions
		s.currentToken = oldCurrentToken
		s.seekTo(currentOffset)

		return length
	}

	return nil
}

// readStreamTokenStart verifies that the "stream" operator exists at the given offset.
func readStreamTokenStart(input core.InputBytes, tokenStart int64) bool {
	input.Seek(tokenStart, 0)

	streamData := tokens.StartStream.Data()
	for i := 0; i < len(streamData); i++ {
		if !input.MoveNext() || input.CurrentByte() != streamData[i] {
			input.Seek(tokenStart, 0)
			return false
		}
	}

	return true
}

// tryBruteForceFileToFindReference scans the entire file to find a reference.
func (s *PdfTokenScannerImpl) tryBruteForceFileToFindReference(reference core.IndirectReference) *tokens.ObjectToken {
	s.isBruteForcing = true
	defer func() { s.isBruteForcing = false }()

	s.seekTo(int64(s.fileHeaderOffset))

	for s.Advance() {
		if objTok, ok := s.currentToken.(*tokens.ObjectToken); ok {
			s.objectLocationProvider.Cache(objTok, true)
		}
	}

	objectToken, found := s.objectLocationProvider.TryGetCached(reference)
	if !found {
		return nil
	}

	return objectToken
}

// getObjectFromStream retrieves an object from within a compressed object stream.
func (s *PdfTokenScannerImpl) getObjectFromStream(
	reference core.IndirectReference, offset core.XrefLocation, navSet []int, depth byte,
) *tokens.ObjectToken {
	streamObjectNumber := offset.Value1

	streamRef, _ := core.NewIndirectReference(streamObjectNumber, 0)
	streamObject := s.getRecursive(streamRef, navSet, depth)

	if streamObject == nil {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Requested a stream object by reference but the requested stream object was not found: %s", reference)))
	}

	stream, ok := streamObject.Data().(*tokens.StreamToken)
	if !ok {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Requested a stream object by reference but the requested stream object was not a stream: %s, %v.",
				reference, streamObject.Data())))
	}

	objects := s.parseObjectStream(stream, offset)

	for _, o := range objects {
		s.objectLocationProvider.Cache(o, false)
	}

	result, found := s.objectLocationProvider.TryGetCached(reference)
	if !found {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Could not find the object %s in the stream %d.", reference, streamObjectNumber)))
	}

	return result
}

// parseObjectStream parses objects stored within an object stream.
func (s *PdfTokenScannerImpl) parseObjectStream(stream *tokens.StreamToken, offset core.XrefLocation) []*tokens.ObjectToken {
	nTok, foundN := stream.StreamDictionary.TryGet(tokens.Create("N"))
	if !foundN {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Object stream dictionary did not provide number of objects %v.", stream.StreamDictionary)))
	}

	numberOfObjects, okN := nTok.(*tokens.NumericToken)
	if !okN {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Object stream dictionary did not provide number of objects %v.", stream.StreamDictionary)))
	}

	firstTok, foundFirst := stream.StreamDictionary.TryGet(tokens.Create("First"))
	if !foundFirst {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Object stream dictionary did not provide first object offset %v.", stream.StreamDictionary)))
	}

	firstTokenNum, okFirst := firstTok.(*tokens.NumericToken)
	if !okFirst {
		panic(core.NewPdfDocumentFormatException(
			fmt.Sprintf("Object stream dictionary did not provide first object offset %v.", stream.StreamDictionary)))
	}

	firstTokenOffset := firstTokenNum.LongVal()

	// Decode the stream data using the filter provider (C#: stream.Decode(filterProvider, this)).
	var decodedBytes []byte
	if s.streamDecoder != nil {
		decodedBytes = s.streamDecoder.DecodeStream(stream)
	} else {
		// Fallback to raw bytes if no decoder is available.
		decodedBytes = stream.Data()
	}

	bytes := core.NewMemoryInputBytes(decodedBytes)

	scanner := NewCoreTokenScanner(bytes, true, s.stackDepthGuard, ScannerScopeNone, nil, s.parsingOptions.UseLenientParsing, true)

	type objEntry struct {
		objectNumber int64
		byteOffset   int64
	}
	objects := make([]objEntry, 0, numberOfObjects.IntVal())

	for i := 0; i < numberOfObjects.IntVal(); i++ {
		scanner.Advance()
		objectNumber := scanner.Current().(*tokens.NumericToken)
		scanner.Advance()
		byteOffset := scanner.Current().(*tokens.NumericToken)

		objects = append(objects, objEntry{
			objectNumber: objectNumber.LongVal(),
			byteOffset:   firstTokenOffset + byteOffset.LongVal(),
		})
	}

	results := make([]*tokens.ObjectToken, 0, len(objects))

	for _, obj := range objects {
		isBetween := ((obj.byteOffset - (scanner.CurrentPosition() - 1)) | ((scanner.CurrentPosition() + 1) - obj.byteOffset)) >= 0
		if !isBetween {
			scanner.Seek(obj.byteOffset, 0)
		}

		scanner.Advance()

		token := scanner.Current()

		if op, ok := token.(*tokens.OperatorToken); ok && op.Data() == "endobj" {
			scanner.Advance()
			token = scanner.Current()
		}

		ref, _ := core.NewIndirectReference(obj.objectNumber, 0)
		results = append(results, tokens.NewObjectToken(offset, ref, token))
	}

	return results
}

// extractTrailingNumber extracts trailing digits from a string.
func extractTrailingNumber(s string) string {
	i := len(s) - 1
	for i >= 0 && s[i] >= '0' && s[i] <= '9' {
		i--
	}
	if i == len(s)-1 {
		return ""
	}
	return s[i+1:]
}

// NewPdfTokenScannerImpl creates a scanner from an existing SeekableTokenScanner.
// This is a backward-compatible adapter that wraps the pre-made scanner.
func NewPdfTokenScannerImpl(
	scanner SeekableTokenScanner,
	locationProvider ObjectLocationProvider,
	filterProvider any,
	encryptionHandler encryptionHandler,
) *PdfTokenScannerImpl {
	impl := &PdfTokenScannerImpl{
		objectLocationProvider: locationProvider,
		encryptionHandler:      encryptionHandler,
		coreScanner:            nil,
		overwrittenTokens:      make(map[core.IndirectReference]*tokens.ObjectToken),
		parsingOptions:         &scannerParsingOptions{},
	}

	if cs, ok := scanner.(*CoreTokenScanner); ok {
		impl.coreScanner = cs
	}

	// Try to extract a stream decoder from the filter provider.
	if fd, ok := filterProvider.(streamDecoder); ok {
		impl.streamDecoder = fd
	}

	return impl
}

var _ PdfTokenScanner = (*PdfTokenScannerImpl)(nil)
