package tokenization

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// CoreTokenScanner is the default TokenScanner for reading PostScript/PDF style data.
type CoreTokenScanner struct {
	inputBytes        core.InputBytes
	usePdfDocEncoding bool
	stackDepthGuard   *core.StackDepthGuard
	stringTokenizer   *StringTokenizer
	arrayTokenizer    *ArrayTokenizer
	dictionaryTokenizer *DictionaryTokenizer
	scope             ScannerScope
	namedDictionaryRequiredKeys map[*tokens.NameToken][]*tokens.NameToken
	useLenientParsing bool
	customTokenizers  []customTokenizerEntry
	currentTokenStart int64
	currentToken      tokens.Token
	hasBytePreRead    bool
	isInInlineImage   bool
	isStream          bool
}

type customTokenizerEntry struct {
	firstByte  byte
	tokenizer  Tokenizer
}

var (
	commentTokenizerInstance = &CommentTokenizer{}
	hexTokenizerInstance     = &HexTokenizer{}
	nameTokenizerInstance    = &NameTokenizer{}
	plainTokenizerInstance   = &PlainTokenizer{}
	numericTokenizerInstance = &NumericTokenizer{}
)

var _ SeekableTokenScanner = (*CoreTokenScanner)(nil)

// NewCoreTokenScanner creates a new CoreTokenScanner from the input.
func NewCoreTokenScanner(
	inputBytes core.InputBytes,
	usePdfDocEncoding bool,
	stackDepthGuard *core.StackDepthGuard,
	scope ScannerScope,
	namedDictionaryRequiredKeys map[*tokens.NameToken][]*tokens.NameToken,
	useLenientParsing bool,
	isStream bool,
) *CoreTokenScanner {
	return &CoreTokenScanner{
		inputBytes:                inputBytes,
		usePdfDocEncoding:         usePdfDocEncoding,
		stackDepthGuard:           stackDepthGuard,
		stringTokenizer:           NewStringTokenizer(usePdfDocEncoding),
		arrayTokenizer:            NewArrayTokenizer(usePdfDocEncoding, stackDepthGuard, useLenientParsing),
		dictionaryTokenizer:       NewDictionaryTokenizer(usePdfDocEncoding, stackDepthGuard, nil, useLenientParsing),
		scope:                     scope,
		namedDictionaryRequiredKeys: namedDictionaryRequiredKeys,
		useLenientParsing:         useLenientParsing,
		isStream:                  isStream,
	}
}

// CurrentTokenStart returns the offset in the input data at which the current token starts.
func (s *CoreTokenScanner) CurrentTokenStart() int64 {
	return s.currentTokenStart
}

// Advance moves the scanner to the next token in the input.
// Returns true if a token was successfully read, false otherwise.
func (s *CoreTokenScanner) Advance() bool {
	if err := s.stackDepthGuard.Enter(); err != nil {
		panic(err)
	}
	defer s.stackDepthGuard.Exit()
	return s.moveToNextInternal()
}

// Current returns the most recently scanned token.
func (s *CoreTokenScanner) Current() tokens.Token {
	return s.currentToken
}

// StackDepthGuard returns the guard object used to track and limit stack depth.
func (s *CoreTokenScanner) StackDepthGuard() *core.StackDepthGuard {
	return s.stackDepthGuard
}

// Seek moves the scanner to the specified position with the given whence.
func (s *CoreTokenScanner) Seek(offset int64, whence int) (int64, error) {
	return s.inputBytes.Seek(offset, whence)
}

// CurrentPosition returns the current byte offset in the input.
func (s *CoreTokenScanner) CurrentPosition() int64 {
	return s.inputBytes.CurrentOffset()
}

// Length returns the total length of the data represented by this scanner.
func (s *CoreTokenScanner) Length() int64 {
	return s.inputBytes.Length()
}

// RegisterCustomTokenizer adds support for a custom tokenizer identified by its first matching byte.
func (s *CoreTokenScanner) RegisterCustomTokenizer(firstByte byte, t Tokenizer) {
	s.customTokenizers = append(s.customTokenizers, customTokenizerEntry{firstByte: firstByte, tokenizer: t})
}

// DeregisterCustomTokenizer removes a previously registered custom tokenizer.
func (s *CoreTokenScanner) DeregisterCustomTokenizer(t Tokenizer) {
	result := make([]customTokenizerEntry, 0, len(s.customTokenizers))
	for _, entry := range s.customTokenizers {
		if entry.tokenizer != t {
			result = append(result, entry)
		}
	}
	s.customTokenizers = result
}

// RecoverFromIncorrectEndImage handles the situation where "EI" was encountered in the inline image data
// but was not the end of the image.
func (s *CoreTokenScanner) RecoverFromIncorrectEndImage(lastEndImageOffset int64) []byte {
	s.inputBytes.Seek(lastEndImageOffset, 0) // io.SeekStart

	if !s.inputBytes.MoveNext() || s.inputBytes.CurrentByte() != 'E' {
		msg := fmt.Sprintf("Failed to recover the image data stream for an inline image at offset %d. "+
			"Expected to read byte 'E' instead got %d.", lastEndImageOffset, s.inputBytes.CurrentByte())
		panic(core.NewPdfDocumentFormatException(msg))
	}

	data := []byte{s.inputBytes.CurrentByte()}

	if !s.inputBytes.MoveNext() || s.inputBytes.CurrentByte() != 'I' {
		msg := fmt.Sprintf("Failed to recover the image data stream for an inline image at offset %d. "+
			"Expected to read second byte 'I' following 'E' instead got %d.", lastEndImageOffset, s.inputBytes.CurrentByte())
		panic(core.NewPdfDocumentFormatException(msg))
	}

	data = append(data, s.inputBytes.CurrentByte())

	data = append(data, s.readUntilEndImage(lastEndImageOffset)...)

	s.inputBytes.MoveNext()

	return data
}

func (s *CoreTokenScanner) moveToNextInternal() bool {
	endAngleBracesRead := 0
	isSkippingLine := false
	isSkippingSymbol := false

	for s.hasBytePreRead && !s.inputBytes.IsAtEnd() || s.inputBytes.MoveNext() {
		s.hasBytePreRead = false
		currentByte := s.inputBytes.CurrentByte()
		c := rune(currentByte)

		if isSkippingLine {
			if core.IsEndOfLineChar(c) {
				isSkippingLine = false
			}
			continue
		}

		var tokenizer Tokenizer

		for _, entry := range s.customTokenizers {
			if currentByte == entry.firstByte {
				tokenizer = entry.tokenizer
				break
			}
		}

		if tokenizer == nil {
			if core.IsWhitespace(currentByte) || isControlChar(c) {
				isSkippingSymbol = false
				continue
			}

			if currentByte == '%' && s.isStream {
				isSkippingLine = true
				continue
			}

			if isSkippingSymbol && c != '>' {
				continue
			}

			switch c {
			case '(':
				tokenizer = s.stringTokenizer
			case '<':
				followingByte, hasFollowing := s.inputBytes.Peek()
				if hasFollowing && followingByte == '<' {
					isSkippingSymbol = true
					tokenizer = s.dictionaryTokenizer

					if s.namedDictionaryRequiredKeys != nil {
						if nameToken, ok := s.currentToken.(*tokens.NameToken); ok {
							if requiredKeys, exists := s.namedDictionaryRequiredKeys[nameToken]; exists {
								tokenizer = NewDictionaryTokenizer(s.usePdfDocEncoding, s.stackDepthGuard, requiredKeys, s.useLenientParsing)
							}
						}
					}
				} else {
					tokenizer = hexTokenizerInstance
				}
			case '>':
				if s.scope == ScannerScopeDictionary {
					endAngleBracesRead++
					if endAngleBracesRead == 2 {
						return false
					}
				}
			case '[':
				tokenizer = s.arrayTokenizer
			case ']':
				if s.scope == ScannerScopeArray {
					return false
				}
			case '/':
				tokenizer = nameTokenizerInstance
			case '%':
				tokenizer = commentTokenizerInstance
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '+', '.':
				tokenizer = numericTokenizerInstance
			default:
				tokenizer = plainTokenizerInstance
			}
		}

		s.currentTokenStart = s.inputBytes.CurrentOffset() - 1

		if tokenizer == nil {
			isSkippingSymbol = true
			s.hasBytePreRead = false
			continue
		}

		token, ok := tokenizer.Tokenize(currentByte, s.inputBytes)
		if !ok {
			isSkippingSymbol = true
			s.hasBytePreRead = false
			continue
		}

		if opToken, ok := token.(*tokens.OperatorToken); ok {
			if opToken.Data() == "BI" {
				s.isInInlineImage = true
				s.hasBytePreRead = false
				continue
			} else if s.isInInlineImage && opToken.Data() == "ID" {
				imageData := s.readInlineImageData()
				s.isInInlineImage = false
				s.currentToken = tokens.NewInlineImageDataToken(imageData)
				s.hasBytePreRead = false
				return true
			}
		}

		s.currentToken = token
		s.hasBytePreRead = tokenizer.ReadsNextByte()

		return true
	}

	return false
}

func (s *CoreTokenScanner) readInlineImageData() []byte {
	if !core.IsWhitespace(s.inputBytes.CurrentByte()) {
		msg := fmt.Sprintf("No whitespace character following the image data (ID) operator. Position: %d.",
			s.inputBytes.CurrentOffset())
		panic(core.NewPdfDocumentFormatException(msg))
	}

	startsAt := s.inputBytes.CurrentOffset() - 2

	return s.readUntilEndImage(startsAt)
}

func (s *CoreTokenScanner) readUntilEndImage(startsAt int64) []byte {
	const lastPlainText byte = 127
	const space byte = 32

	imageData := make([]byte, 0)
	var prevByte byte = 0

	for s.inputBytes.MoveNext() {
		currentByte := s.inputBytes.CurrentByte()

		if currentByte == 'I' && prevByte == 'E' {
			buffer := make([]byte, 6)
			currentOffset := s.inputBytes.CurrentOffset()
			read, _ := s.inputBytes.Read(buffer)

			isEnd := true

			if read == len(buffer) {
				containsWhitespace := false
				for i := 0; i < len(buffer); i++ {
					b := buffer[i]

					if core.IsWhitespace(b) {
						containsWhitespace = true
						continue
					}

					if b > lastPlainText {
						isEnd = false
						break
					}

					if b < space && b != '\r' && b != '\n' && b != '\t' {
						isEnd = false
						break
					}
				}

				if !containsWhitespace {
					isEnd = false
				}
			}

			s.inputBytes.Seek(currentOffset, 0) // io.SeekStart

			if isEnd {
				imageData = imageData[:len(imageData)-1]
				return imageData
			}
		}

		imageData = append(imageData, currentByte)
		prevByte = currentByte
	}

	if s.useLenientParsing {
		return imageData
	}

	msg := fmt.Sprintf("No end of inline image data (EI) was found for image data at position %d.", startsAt)
	panic(core.NewPdfDocumentFormatException(msg))
}

func isControlChar(c rune) bool {
	return (c >= 0 && c <= 31) || c == 127
}
