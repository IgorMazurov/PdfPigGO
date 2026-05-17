package filestructure

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

const junkTokensTolerance = 30
const versionLength = 8
const bufferLength = 64
const toDoubleStartLength = 4

// Parse retrieves the version header from the PDF file.
// The first line of a PDF file should be a header consisting of %PDF- followed by
// a version number of the form 1.N, where N is a digit between 0 and 7.
func Parse(scanner tokenization.SeekableTokenScanner, inputBytes core.InputBytes, isLenientParsing bool, log logging.Log) (*content.HeaderVersion, error) {
	if scanner == nil {
		return nil, fmt.Errorf("scanner cannot be nil")
	}

	startPosition := scanner.CurrentPosition()

	attempts := 0
	var comment *tokens.CommentToken

	for {
		if attempts == junkTokensTolerance || !scanner.Advance() {
			version, ok := tryBruteForceVersionLocation(startPosition, inputBytes)
			if !ok {
				return nil, core.NewPdfDocumentFormatException("Could not find the version header comment at the start of the document.")
			}

			scanner.Seek(startPosition, 0) // io.SeekStart
			return version, nil
		}

		if ct, ok := scanner.Current().(*tokens.CommentToken); ok {
			comment = ct
		} else {
			comment = nil
		}

		attempts++

		if comment != nil {
			break
		}
	}

	return getHeaderVersionAndResetScanner(comment, scanner, isLenientParsing, log)
}

func getHeaderVersionAndResetScanner(comment *tokens.CommentToken, scanner tokenization.SeekableTokenScanner, isLenientParsing bool, log logging.Log) (*content.HeaderVersion, error) {
	data := comment.Data()

	if !strings.HasPrefix(strings.ToLower(data), "pdf-1.") && !strings.HasPrefix(strings.ToLower(data), "fdf-1.") {
		return handleMissingVersion(comment, isLenientParsing, log)
	}

	versionStr := data[toDoubleStartLength:]
	version, err := strconv.ParseFloat(versionStr, 64)
	if err != nil {
		return handleMissingVersion(comment, isLenientParsing, log)
	}

	atEnd := scanner.CurrentPosition() == scanner.Length()
	rewind := int64(2)
	if atEnd {
		rewind = 1
	}

	commentOffset := scanner.CurrentPosition() - int64(len(data)) - rewind

	scanner.Seek(0, 0) // io.SeekStart

	return content.NewHeaderVersion(version, data, commentOffset)
}

func tryBruteForceVersionLocation(startPosition int64, inputBytes core.InputBytes) (*content.HeaderVersion, bool) {
	inputBytes.Seek(startPosition, 0) // io.SeekStart

	buffer := make([]byte, bufferLength)

	currentOffset := startPosition

	for {
		n, err := inputBytes.Read(buffer)
		if err != nil || n == 0 {
			return nil, false
		}
		readLength := n

		contentBytes := buffer[:readLength]
		lowerContent := asciiToLower(contentBytes)

		pdfIndex := bytes.Index(lowerContent, []byte("%pdf-"))
		fdfIndex := bytes.Index(lowerContent, []byte("%fdf-"))

		actualIndex := -1
		if pdfIndex >= 0 {
			actualIndex = pdfIndex
		} else if fdfIndex >= 0 {
			actualIndex = fdfIndex
		}

		if actualIndex >= 0 && len(contentBytes)-actualIndex >= versionLength {
			numberPart := string(contentBytes[actualIndex+5 : actualIndex+5+3])
			version, err := strconv.ParseFloat(numberPart, 64)
			if err == nil {
				afterCommentSymbolIndex := actualIndex + 1
				versionString := string(contentBytes[afterCommentSymbolIndex : afterCommentSymbolIndex+versionLength-1])

				headerVersion, hvErr := content.NewHeaderVersion(version, versionString, currentOffset+int64(actualIndex))
				if hvErr != nil {
					inputBytes.Seek(startPosition, 0) // io.SeekStart
					return nil, false
				}

				inputBytes.Seek(startPosition, 0) // io.SeekStart
				return headerVersion, true
			}
		}

		currentOffset += int64(readLength - versionLength)
		inputBytes.Seek(currentOffset, 0) // io.SeekStart

		if readLength != bufferLength {
			break
		}
	}

	return nil, false
}

// asciiToLower converts ASCII uppercase letters to lowercase in a byte slice,
// leaving all other bytes unchanged. This avoids UTF-8 encoding issues with
// non-ASCII byte values that would occur with strings.ToLower or bytes.ToLower.
func asciiToLower(b []byte) []byte {
	result := make([]byte, len(b))
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			result[i] = b[i] + ('a' - 'A')
		} else {
			result[i] = b[i]
		}
	}
	return result
}

func handleMissingVersion(comment *tokens.CommentToken, isLenientParsing bool, log logging.Log) (*content.HeaderVersion, error) {
	if isLenientParsing {
		log.Warn(fmt.Sprintf("Did not find a version header of the correct format, defaulting to 1.4 since lenient. Header was: %s.", comment.Data()))

		return content.NewHeaderVersion(1.4, "PDF-1.4", 0)
	}

	return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("The comment which should have provided the version was in the wrong format: %s.", comment.Data()))
}
