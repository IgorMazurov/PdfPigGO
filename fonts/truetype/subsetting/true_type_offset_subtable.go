package subsetting

import (
	"io"
	"math"

	"github.com/uglytoad/pdfpig/go/core"
)

// TrueTypeOffsetSubtable represents the offset table of a TrueType font subset.
type TrueTypeOffsetSubtable struct {
	numberOfTables byte
}

// NewTrueTypeOffsetSubtable creates a new TrueTypeOffsetSubtable with the given number of tables.
func NewTrueTypeOffsetSubtable(numberOfTables byte) *TrueTypeOffsetSubtable {
	return &TrueTypeOffsetSubtable{
		numberOfTables: numberOfTables,
	}
}

// Write writes the offset subtable header to the output stream.
func (s *TrueTypeOffsetSubtable) Write(w io.Writer) error {
	versionHeader := []byte{0, 1, 0, 0}

	if _, err := w.Write(versionHeader); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, uint16(s.numberOfTables)); err != nil {
		return err
	}

	maximumPowerOf2 := getHighestPowerOf2(int(s.numberOfTables))
	searchRange := uint16(maximumPowerOf2 * 16)
	entrySelector := uint16(math.Log2(float64(maximumPowerOf2)))
	rangeShift := uint16(int(s.numberOfTables)*16 - int(searchRange))

	if _, err := core.WriteUShort(w, searchRange); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, entrySelector); err != nil {
		return err
	}

	if _, err := core.WriteUShort(w, rangeShift); err != nil {
		return err
	}

	return nil
}

// getHighestPowerOf2 returns the highest power of 2 less than or equal to numberOfTables.
func getHighestPowerOf2(numberOfTables int) uint16 {
	var result uint16 = 1
	for i := 0; i < 8*2; i++ {
		power := uint16(1 << i)
		if power > uint16(numberOfTables) {
			break
		}
		result = power
	}
	return result
}
