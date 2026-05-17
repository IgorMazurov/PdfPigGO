package filestructure

import (
	"strconv"

	"github.com/uglytoad/pdfpig/go/core"
)

const minimumSearchOffset = 6

// GetObjectLocations finds the offset of every object in the document by searching
// the entire document contents using brute force scanning.
func GetObjectLocations(bytes core.InputBytes) (map[core.IndirectReference]core.XrefLocation, error) {
	if bytes == nil {
		return nil, &core.PdfDocumentFormatException{Message: "bytes cannot be nil"}
	}

	loopProtection := 0

	lastEndOfFile := getLastEndOfFileMarker(bytes)

	results := make(map[core.IndirectReference]core.XrefLocation)

	var generationRunes []rune
	var objectNumberRunes []rune

	originPosition := bytes.CurrentOffset()

	currentOffset := int64(minimumSearchOffset)

	currentlyInObject := false

	objBuffer := make([]byte, 4)

	for currentOffset < lastEndOfFile && !bytes.IsAtEnd() {
		if loopProtection > 10_000_000 {
			return nil, core.NewPdfDocumentFormatException("Failed to brute-force search the file due to an infinite loop.")
		}

		loopProtection++

		if currentlyInObject {
			if bytes.IsAtEnd() {
				break
			}

			if bytes.CurrentByte() == 'e' {
				next, hasNext := bytes.Peek()

				if hasNext && next == 'n' {
					if core.IsStringStr(bytes, "endobj") {
						currentlyInObject = false
						loopProtection = 0

						for i := 0; i < len("endobj"); i++ {
							bytes.MoveNext()
							currentOffset++
						}
					} else {
						bytes.MoveNext()
						currentOffset++
					}
				} else {
					bytes.MoveNext()
					currentOffset++
				}
			} else if core.IsWhitespace(bytes.CurrentByte()) {
				next, hasNext := bytes.Peek()
				if hasNext && next == 'o' {
					if core.IsStringStr(bytes, " obj") {
						currentlyInObject = false
						currentOffset--
						loopProtection = 0
					} else {
						bytes.MoveNext()
						currentOffset++
						loopProtection = 0
					}
				} else {
					bytes.MoveNext()
					currentOffset++
					loopProtection = 0
				}
			} else {
				bytes.MoveNext()
				currentOffset++
				loopProtection = 0
			}

			continue
		}

		bytes.Seek(currentOffset, 0) // io.SeekStart

		bytes.Read(objBuffer)

		if !isStartObjMarker(objBuffer) {
			currentOffset++
			continue
		}

		offset := currentOffset + 1

		bytes.Seek(offset, 0) // io.SeekStart

		for core.IsWhitespace(bytes.CurrentByte()) && offset >= minimumSearchOffset {
			offset--
			bytes.Seek(int64(offset), 0) // io.SeekStart
		}

		for core.IsDigit(bytes.CurrentByte()) && offset >= minimumSearchOffset {
			generationRunes = append(generationRunes, rune(bytes.CurrentByte()))
			offset--
			bytes.Seek(int64(offset), 0) // io.SeekStart
		}

		if !core.IsWhitespace(bytes.CurrentByte()) {
			generationRunes = generationRunes[:0]
			objectNumberRunes = objectNumberRunes[:0]
			currentOffset++
			continue
		}

		for core.IsWhitespace(bytes.CurrentByte()) {
			offset--
			bytes.Seek(int64(offset), 0) // io.SeekStart
		}

		for core.IsDigit(bytes.CurrentByte()) && offset >= minimumSearchOffset {
			objectNumberRunes = append(objectNumberRunes, rune(bytes.CurrentByte()))
			offset--
			bytes.Seek(int64(offset), 0) // io.SeekStart
		}

		if len(objectNumberRunes) == 0 || len(generationRunes) == 0 {
			generationRunes = generationRunes[:0]
			objectNumberRunes = objectNumberRunes[:0]
			currentOffset++
			continue
		}

		obj, err := strconv.ParseInt(reverseRunes(objectNumberRunes), 10, 64)
		if err != nil {
			generationRunes = generationRunes[:0]
			objectNumberRunes = objectNumberRunes[:0]
			currentOffset++
			continue
		}

		generation, err := strconv.ParseInt(reverseRunes(generationRunes), 10, 32)
		if err != nil {
			generationRunes = generationRunes[:0]
			objectNumberRunes = objectNumberRunes[:0]
			currentOffset++
			continue
		}

		ref, err := core.NewIndirectReference(obj, int(generation))
		if err != nil {
			generationRunes = generationRunes[:0]
			objectNumberRunes = objectNumberRunes[:0]
			currentOffset++
			continue
		}

		results[ref] = core.File(bytes.CurrentOffset())

		generationRunes = generationRunes[:0]
		objectNumberRunes = objectNumberRunes[:0]

		currentlyInObject = true

		currentOffset += int64(len(objBuffer))

		bytes.Seek(currentOffset, 0) // io.SeekStart
		loopProtection = 0
	}

	bytes.Seek(originPosition, 0) // io.SeekStart

	return results, nil
}

func getLastEndOfFileMarker(bytes core.InputBytes) int64 {
	originalOffset := bytes.CurrentOffset()

	searchTerm := "%%EOF"

	minimumEndOffset := bytes.Length() - int64(len(searchTerm)) + 1

	bytes.Seek(minimumEndOffset, 0) // io.SeekStart

	for bytes.CurrentOffset() > 0 {
		if core.IsStringStr(bytes, searchTerm) {
			position := bytes.CurrentOffset()

			bytes.Seek(originalOffset, 0) // io.SeekStart

			return position
		}

		minimumEndOffset--
		if minimumEndOffset < 1 {
			break
		}
		bytes.Seek(minimumEndOffset, 0) // io.SeekStart
	}

	bytes.Seek(originalOffset, 0) // io.SeekStart
	return int64(^uint64(0) >> 1) // math.MaxInt64
}

func isStartObjMarker(data []byte) bool {
	if !core.IsWhitespace(data[0]) {
		return false
	}

	return (data[1] == 'o' || data[1] == 'O') &&
		(data[2] == 'b' || data[2] == 'B') &&
		(data[3] == 'j' || data[3] == 'J')
}

// reverseRunes reverses a rune slice in-place and returns it as a string.
func reverseRunes(runes []rune) string {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
