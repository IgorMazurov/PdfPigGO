package truetype

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
)

const (
	headerTableTag            = "head"
	checksumAdjustmentPosition = 8
)

// CalculateWholeFontChecksum calculates the checksum for the whole font by setting
// checksum adjustment in the head table to 0.
func CalculateWholeFontChecksum(bytes core.InputBytes, headerTable TrueTypeHeaderTable) (uint32, error) {
	if bytes == nil {
		return 0, fmt.Errorf("bytes must not be nil")
	}

	if !isHeadTable(headerTable) {
		return 0, fmt.Errorf("can only calculate checksum for the whole font when the head table is provided. Got: %s", headerTable.Tag)
	}

	bytes.Seek(0, 0) // io.SeekStart

	skipped := toChecksumSkippedBytes(bytes, headerTable)

	return CalculateBytes(skipped), nil
}

// CalculateTable calculates the checksum for the specific table.
func CalculateTable(bytes core.InputBytes, table TrueTypeHeaderTable) (uint32, error) {
	bytes.Seek(int64(table.Offset), 0) // io.SeekStart

	if isHeadTable(table) {
		fullTableBytes := make([]byte, table.Length)
		read, err := bytes.Read(fullTableBytes)
		if err != nil {
			return 0, err
		}
		if read != int(table.Length) {
			return 0, fmt.Errorf("could not read full table: expected %d bytes, got %d", table.Length, read)
		}

		fullTableBytes[checksumAdjustmentPosition] = 0
		fullTableBytes[checksumAdjustmentPosition+1] = 0
		fullTableBytes[checksumAdjustmentPosition+2] = 0
		fullTableBytes[checksumAdjustmentPosition+3] = 0

		return CalculateBytes(fullTableBytes), nil
	}

	var result uint32

	endAt := int64(table.Offset + table.Length)

	for {
		next, ok := tryReadUIntInput(bytes, endAt)
		if !ok {
			break
		}
		result += next
	}

	return result, nil
}

// CalculateBytes calculates the TrueType checksum for the provided bytes.
func CalculateBytes(bytes []byte) uint32 {
	var result uint32

	for i := 0; i+3 < len(bytes); i += 4 {
		next := uint32(bytes[i])<<24 | uint32(bytes[i+1])<<16 | uint32(bytes[i+2])<<8 | uint32(bytes[i+3])
		result += next
	}

	if remainder := len(bytes) % 4; remainder != 0 {
		start := len(bytes) - remainder
		var padded [4]byte
		copy(padded[:], bytes[start:])
		next := uint32(padded[0])<<24 | uint32(padded[1])<<16 | uint32(padded[2])<<8 | uint32(padded[3])
		result += next
	}

	return result
}

func isHeadTable(table TrueTypeHeaderTable) bool {
	return strings.EqualFold(headerTableTag, table.Tag)
}

func tryReadUIntInput(input core.InputBytes, endAt int64) (uint32, bool) {
	if input.CurrentOffset() >= endAt {
		return 0, false
	}

	readNext := func() byte {
		if input.CurrentOffset() == endAt || !input.MoveNext() {
			return 0
		}
		return input.CurrentByte()
	}

	top := readNext()
	three := readNext()
	two := readNext()
	one := readNext()

	result := uint32(top)<<24 | uint32(three)<<16 | uint32(two)<<8 | uint32(one)

	return result, true
}

func toChecksumSkippedBytes(bytes core.InputBytes, table TrueTypeHeaderTable) []byte {
	var result []byte
	tableStartAdj := int64(table.Offset + checksumAdjustmentPosition)
	tableEndAdj := tableStartAdj + 4

	for bytes.MoveNext() {
		offset := bytes.CurrentOffset()
		if offset > tableStartAdj && offset <= tableEndAdj {
			continue
		}
		result = append(result, bytes.CurrentByte())
	}

	return result
}
