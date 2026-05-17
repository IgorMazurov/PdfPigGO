package truetype_test

import (
	"os"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// getFileBytes reads a test data file from the testdata directory.
func getFileBytes(t *testing.T, name string) []byte {
	t.Helper()
	path := "testdata/" + name
	buf, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read test file %q: %v", path, err)
	}
	return buf
}

// readUint32 reads a big-endian uint32 from bytes at the given offset.
func readUint32(b []byte, off int) uint32 {
	return uint32(b[off])<<24 | uint32(b[off+1])<<16 | uint32(b[off+2])<<8 | uint32(b[off+3])
}

// readUint16 reads a big-endian uint16 from bytes at the given offset.
func readUint16(b []byte, off int) uint16 {
	return uint16(b[off])<<8 | uint16(b[off+1])
}

// readTag reads a 4-byte table tag from bytes at the given offset.
func readTag(b []byte, off int) string {
	return string(b[off : off+4])
}

// readTableDirectory reads the TrueType font directory and returns the table headers
// along with the checksum adjustment value from the head table (if present).
func readTableDirectory(bytes []byte) (map[string]truetype.TrueTypeHeaderTable, uint32) {
	numTables := int(readUint16(bytes, 4))

	tablesMap := make(map[string]truetype.TrueTypeHeaderTable, numTables)

	for i := 0; i < numTables; i++ {
		entryOffset := 12 + i*16

		tag := readTag(bytes, entryOffset)
		checksum := readUint32(bytes, entryOffset+4)
		offset := readUint32(bytes, entryOffset+8)
		length := readUint32(bytes, entryOffset+12)

		tablesMap[tag] = truetype.TrueTypeHeaderTable{
			Tag:      tag,
			CheckSum: checksum,
			Offset:   offset,
			Length:   length,
		}
	}

	var checkSumAdjustment uint32
	if headHeader, ok := tablesMap["head"]; ok {
		checkSumAdjustment = readUint32(bytes, int(headHeader.Offset)+8)
	}

	return tablesMap, checkSumAdjustment
}

// runChecksumValidation parses the font bytes and validates checksums for all table headers.
func runChecksumValidation(t *testing.T, bytes []byte, checkHeaderChecksum, checkWholeFileChecksum bool) {
	t.Helper()

	tableHeaders, checkSumAdjustment := readTableDirectory(bytes)

	inputBytes := core.NewMemoryInputBytes(bytes)

	for tag, header := range tableHeaders {
		if tag == "head" {
			if checkHeaderChecksum {
				headerChecksum, err := truetype.CalculateTable(inputBytes, header)
				if err != nil {
					t.Errorf("CalculateTable for head table: %v", err)
					continue
				}

				if header.CheckSum != headerChecksum {
					t.Errorf("head checksum mismatch: expected 0x%08X, got 0x%08X", header.CheckSum, headerChecksum)
				}
			}
			continue
		}

		tableBytes := bytes[header.Offset : header.Offset+header.Length]
		checksum := truetype.CalculateBytes(tableBytes)

		if header.CheckSum != checksum {
			t.Errorf("%s checksum mismatch (CalculateBytes): expected 0x%08X, got 0x%08X", tag, header.CheckSum, checksum)
		}

		checksumByTable, err := truetype.CalculateTable(inputBytes, header)
		if err != nil {
			t.Errorf("%s CalculateTable error: %v", tag, err)
			continue
		}

		if header.CheckSum != checksumByTable {
			t.Errorf("%s checksum mismatch (CalculateTable): expected 0x%08X, got 0x%08X", tag, header.CheckSum, checksumByTable)
		}
	}

	if checkWholeFileChecksum {
		headHeader := tableHeaders["head"]
		wholeFontChecksum, err := truetype.CalculateWholeFontChecksum(inputBytes, headHeader)
		if err != nil {
			t.Fatalf("CalculateWholeFontChecksum: %v", err)
		}

		adjustment := uint32(0xB1B0AFBA) - wholeFontChecksum

		if checkSumAdjustment != adjustment {
			t.Errorf("checksum adjustment mismatch: expected 0x%08X, got 0x%08X", checkSumAdjustment, adjustment)
		}

		expectedWholeFontChecksum := uint32(0xB1B0AFBA) - checkSumAdjustment

		if expectedWholeFontChecksum != wholeFontChecksum {
			t.Errorf("whole font checksum mismatch: expected 0x%08X, got 0x%08X", expectedWholeFontChecksum, wholeFontChecksum)
		}
	}
}

// Roboto-Regular.ttf has wrong checksums in the file.
func TestCalculatedChecksumsMatchRoboto(t *testing.T) {
	fileBytes := getFileBytes(t, "Roboto-Regular.ttf")
	runChecksumValidation(t, fileBytes, false, false)
}

// Andada-Regular.ttf has correct checksums throughout.
func TestCalculatedChecksumsMatchAndada(t *testing.T) {
	fileBytes := getFileBytes(t, "Andada-Regular.ttf")
	runChecksumValidation(t, fileBytes, true, true)
}

// google-simple-doc.ttf has wrong checksum adjustment.
func TestCalculatedChecksumsMatchGoogleDoc(t *testing.T) {
	fileBytes := getFileBytes(t, "google-simple-doc.ttf")
	runChecksumValidation(t, fileBytes, true, false)
}

// PMingLiU.ttf has wrong checksum adjustment.
func TestCalculatedChecksumsMatchPMing(t *testing.T) {
	fileBytes := getFileBytes(t, "PMingLiU.ttf")
	runChecksumValidation(t, fileBytes, true, false)
}
