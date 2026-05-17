package filestructure

import (
	"fmt"
	"strconv"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

// XrefBruteForcerResult holds the result of a brute-force xref search.
type XrefBruteForcerResult struct {
	// XRefParts contains all xref sections found during the scan.
	XRefParts []XrefSection
	// ObjectOffsets maps each indirect reference to its file location,
	// discovered by scanning for "obj" markers.
	ObjectOffsets map[core.IndirectReference]core.XrefLocation
	// LastTrailer is the last trailer dictionary encountered during the scan.
	LastTrailer *tokens.DictionaryToken
}

// FindAllXrefsInFileOrder scans the entire PDF byte stream looking for xref
// tables and /XRef streams, recording all object positions along the way.
func FindAllXrefsInFileOrder(bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, log logging.Log) XrefBruteForcerResult {
	results := make([]XrefSection, 0)

	xrefOffsetSeen := make(map[int64]bool)

	bruteForceObjPositions := make(map[core.IndirectReference]core.XrefLocation)

	var trailer *tokens.DictionaryToken

	bytes.Seek(0, 0)

	buffer := util.NewCircularByteBuffer(10)

	numberByteBuffer := make([]byte, 0, 20)

	inNum := false
	lastWhitespace := false
	inComment := false

	numericsQueue := [2]int64{}
	positionsQueue := [2]int64{}

	var lastObjPosition *int64

	clearQueues := func() {
		numericsQueue[0] = 0
		numericsQueue[1] = 0
		positionsQueue[0] = 0
		positionsQueue[1] = 0
	}

	addQueues := func(num int64) {
		numericsQueue[0] = numericsQueue[1]
		numericsQueue[1] = num

		positionsQueue[0] = positionsQueue[1]
		positionsQueue[1] = bytes.CurrentOffset() - int64(len(numberByteBuffer)) - 1
	}

	for bytes.MoveNext() && !bytes.IsAtEnd() {
		if bytes.CurrentByte() == '%' {
			inComment = true

			if inNum && len(numberByteBuffer) > 0 {
				num := core.BytesAsLatin1String(numberByteBuffer)
				if numLong, err := strconv.ParseInt(num, 10, 64); err == nil {
					addQueues(numLong)
				}

				numberByteBuffer = numberByteBuffer[:0]
			}

			inNum = false
			lastWhitespace = false

		}

		if core.IsWhitespace(bytes.CurrentByte()) {
			if core.IsEndOfLineByte(bytes.CurrentByte()) {
				inComment = false
			}

			buffer.Add(' ')

			if inNum && len(numberByteBuffer) > 0 {
				num := core.BytesAsLatin1String(numberByteBuffer)
				if numLong, err := strconv.ParseInt(num, 10, 64); err == nil {
					addQueues(numLong)
				}

				numberByteBuffer = numberByteBuffer[:0]
			}

			lastWhitespace = true
			inNum = false
		} else {
			buffer.Add(bytes.CurrentByte())

			if !inComment && core.IsDigit(bytes.CurrentByte()) && (inNum || lastWhitespace) {
				inNum = true
				numberByteBuffer = append(numberByteBuffer, bytes.CurrentByte())
			} else {
				inNum = false
				numberByteBuffer = numberByteBuffer[:0]
			}

			lastWhitespace = false
		}

		if buffer.EndsWith(" obj") && numericsQueue[0] > 0 {
			if ref, err := core.NewIndirectReference(numericsQueue[0], int(numericsQueue[1])); err == nil {
				bruteForceObjPositions[ref] = core.File(positionsQueue[0])
			}

			pos := positionsQueue[0]
			lastObjPosition = &pos

			clearQueues()
		} else if buffer.EndsWith(" xref") {
			clearQueues()

			potentialTableOffset := bytes.CurrentOffset() - 4

			if xrefOffsetSeen[potentialTableOffset] {
				log.Debug(fmt.Sprintf("Skipping circular xref reference at %d", potentialTableOffset))
				continue
			}
			xrefOffsetSeen[potentialTableOffset] = true

			table := TryReadTableAtOffset(
				NewFileHeaderOffset(0),
				potentialTableOffset,
				bytes,
				scanner,
				log,
			)

			if table != nil {
				results = append(results, table)
			} else {
				log.Warn(fmt.Sprintf("Found a table at %d but couldn't parse it.", potentialTableOffset))
			}
		} else if buffer.EndsWith("/XRef") {
			clearQueues()

			offset := lastObjPosition
			if offset == nil {
				log.Error("Found an /XRef without having encountered an object first")
				continue
			}

			if xrefOffsetSeen[*offset] {
				log.Debug(fmt.Sprintf("Skipping circular /XRef reference at %d", *offset))
				continue
			}
			xrefOffsetSeen[*offset] = true

			stream := TryReadStreamAtOffset(
				NewFileHeaderOffset(0),
				*offset,
				bytes,
				scanner,
				log,
			)

			if stream != nil {
				results = append(results, stream)
			}
		} else if buffer.EndsWith("trailer ") {
			clearQueues()

			if scanner.Advance() {
				if dictToken, ok := scanner.Current().(*tokens.DictionaryToken); ok {
					trailer = dictToken
				}
			}
		}
	}

	return XrefBruteForcerResult{
		XRefParts:     results,
		ObjectOffsets: bruteForceObjPositions,
		LastTrailer:   trailer,
	}
}
