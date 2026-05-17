package filestructure

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

// xrefTableReadMode tracks the current parsing state inside an xref table.
type xrefTableReadMode int

const (
	xrefSubsectionHeader xrefTableReadMode = 2
	xrefModeEntry        xrefTableReadMode = 3
)

// TryReadTableAtOffset attempts to parse a traditional xref table starting at the given byte offset.
// Returns nil if no valid table is found or if parsing fails.
func TryReadTableAtOffset(fileHeaderOffset FileHeaderOffset, xrefOffset int64, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, log logging.Log) *XrefTable {
	if xrefOffset >= bytes.Length() || xrefOffset < 0 {
		return nil
	}

	if xrefOffset > 0 {
		scanner.Seek(xrefOffset, 0) // SeekStart
	}

	correctionType := XrefOffsetCorrectionNone
	correction := int64(0)

	if !tryReadXrefToken(scanner) {
		log.Debug("Xref not found at " + itoa(xrefOffset) + ", trying to recover")
		recovered := tryRecoverOffsetTable(fileHeaderOffset, xrefOffset, bytes, scanner)
		if recovered == nil {
			return nil
		}

		scanner.Seek(recovered.correctOffset, 0) // SeekStart
		if !tryReadXrefToken(scanner) {
			return nil
		}

		correctionType = recovered.correctionType
		correction = recovered.correctOffset - xrefOffset
		xrefOffset = recovered.correctOffset
	}

	const objRowSentinel = int64(-1)
	const freeSentinel = int64(0)
	const occupiedSentinel = int64(1)

	readNums := make([]int64, 0)

	var trailer *tokens.DictionaryToken
	readInLine := 0
	clearReadLine := false
	expectedEntryCount := 0
	mode := xrefSubsectionHeader

	for scanner.Advance() {
		if mode == xrefModeEntry && expectedEntryCount <= 0 {
			mode = xrefSubsectionHeader
		}

		readInLine++
		token := scanner.Current()

		switch t := token.(type) {
		case *tokens.NumericToken:
			readNums = append(readNums, t.LongVal())

			if mode == xrefSubsectionHeader && readInLine == 2 {
				mode = xrefModeEntry
				expectedEntryCount = int(t.LongVal())
				clearReadLine = true
			} else if mode == xrefModeEntry && readInLine > 2 {
				if clearReadLine {
					clearReadLine = false
					readInLine = 1
				} else {
					return nil
				}
			}

		case *tokens.OperatorToken:
			data := t.Data()
			lowerData := strings.ToLower(data)

			if lowerData == "f" && readInLine == 3 {
				readNums = append(readNums, freeSentinel)
				readInLine = 0
				expectedEntryCount--
				readNums = insertAt(readNums, len(readNums)-3, objRowSentinel)
			} else if lowerData == "n" && readInLine == 3 {
				readNums = append(readNums, occupiedSentinel)
				readInLine = 0
				expectedEntryCount--
				readNums = insertAt(readNums, len(readNums)-3, objRowSentinel)
			} else if lowerData == "trailer" {
				if !scanner.Advance() {
					return nil
				}
				var dictOk bool
				trailer, dictOk = scanner.Current().(*tokens.DictionaryToken)
				if !dictOk {
					return nil
				}
				goto doneParsing
			} else if mode == xrefSubsectionHeader {
				if lowerData == "obj" {
					readNums = removeRange(readNums, len(readNums)-2, 2)
				}
				goto doneParsing
			} else {
				return nil
			}

		case *tokens.CommentToken:
			readInLine--

		default:
			if _, isComment := token.(*tokens.CommentToken); !isComment {
				goto doneParsing
			}
		}
	}

doneParsing:
	if len(readNums) == 0 {
		if trailer != nil {
			return NewXrefTable(xrefOffset, make(map[core.IndirectReference]core.XrefLocation), trailer, correctionType, correction)
		}
		return nil
	}

	offsets := make(map[core.IndirectReference]core.XrefLocation)
	buff := make([]int64, 4)
	objNum := int64(-1)
	ix := 0

	tryReadBuff := func(length int) bool {
		for i := 0; i < length; i++ {
			if ix >= len(readNums) {
				return false
			}
			buff[i] = readNums[ix]
			ix++
		}
		return true
	}

	for ix < len(readNums) {
		if !tryReadBuff(2) {
			return nil
		}

		first := buff[0]
		second := buff[1]

		if first != objRowSentinel {
			objNum = first
		} else {
			if objNum == -1 {
				return nil
			}
			second = 1
			ix -= 2
		}

		for i := int64(0); i < second; i++ {
			if !tryReadBuff(4) {
				return nil
			}

			sentinel := buff[0]
			objOffset := buff[1]
			gen := buff[2]
			entryType := buff[3]

			if sentinel != objRowSentinel {
				return nil
			}

			if entryType == occupiedSentinel {
				indirectRef, err := core.NewIndirectReference(objNum, int(gen))
				if err != nil {
					return nil
				}
				offsets[indirectRef] = core.File(objOffset)
			}

			objNum++
		}
	}

	return NewXrefTable(xrefOffset, offsets, trailer, correctionType, correction)
}

// tryReadXrefToken attempts to read an "xref" operator token from the scanner.
func tryReadXrefToken(scanner tokenization.SeekableTokenScanner) bool {
	if !scanner.Advance() {
		return false
	}

	op, ok := scanner.Current().(*tokens.OperatorToken)
	if !ok {
		return false
	}

	data := op.Data()
	lowerData := strings.ToLower(data)

	if lowerData == "xref" {
		return true
	}

	if strings.HasPrefix(lowerData, "xref") {
		backtrack := int64(len(data)) - int64(len("xref"))
		scanner.Seek(scanner.CurrentPosition()-backtrack, 0) // SeekStart
		return true
	}

	return false
}

// offsetRecoveryTable holds the result of an xref table offset recovery attempt.
type offsetRecoveryTable struct {
	correctOffset  int64
	correctionType XrefOffsetCorrection
}

// tryRecoverOffsetTable attempts to find a valid xref table near the given offset.
func tryRecoverOffsetTable(fileHeaderOffset FileHeaderOffset, xrefOffset int64, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner) *offsetRecoveryTable {
	if fileHeaderOffset.Value > 0 {
		scanner.Seek(xrefOffset+int64(fileHeaderOffset.Value), 0) // SeekStart
		if tryReadXrefToken(scanner) {
			return &offsetRecoveryTable{
				correctOffset:  xrefOffset + int64(fileHeaderOffset.Value),
				correctionType: XrefOffsetCorrectionFileHeaderOffset,
			}
		}
	}

	buffer := make([]byte, 20)
	offset := xrefOffset - 10
	if offset < 0 {
		offset = 0
	}
	bytes.Seek(offset, 0) // SeekStart

	readCount, _ := bytes.Read(buffer)
	if readCount < len(buffer) {
		return nil
	}

	str := core.BytesAsLatin1String(buffer)
	xrefIx := strings.Index(strings.ToLower(str), "xref")
	if xrefIx < 0 {
		return nil
	}

	actualOffset := offset + int64(xrefIx)
	scanner.Seek(actualOffset, 0) // SeekStart

	if tryReadXrefToken(scanner) {
		return &offsetRecoveryTable{
			correctOffset:  actualOffset,
			correctionType: XrefOffsetCorrectionRandom,
		}
	}

	return nil
}

// insertAt inserts value at the specified index in the slice.
func insertAt(slice []int64, index int, value int64) []int64 {
	slice = append(slice, 0)
	copy(slice[index+1:], slice[index:])
	slice[index] = value
	return slice
}

// removeRange removes count elements starting at index from the slice.
func removeRange(slice []int64, index int, count int) []int64 {
	if index >= len(slice) {
		return slice
	}
	end := index + count
	if end > len(slice) {
		end = len(slice)
	}
	return append(slice[:index], slice[end:]...)
}

// itoa converts an int64 to its decimal string representation.
func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append(buf, byte('0'+n%10))
		n /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
