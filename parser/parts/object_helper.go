package parts

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
)

const objectNumberThreshold int64 = 10_000_000_000

const generationNumberThreshold int64 = 65_535

// ReadObjectNumber reads and validates a PDF object number from the input.
func ReadObjectNumber(bytes core.InputBytes) (int64, error) {
	result, err := core.ReadLong(bytes)
	if err != nil {
		return 0, err
	}
	if result < 0 || result >= objectNumberThreshold {
		return 0, fmt.Errorf("object number '%d' has more than 10 digits or is negative", result)
	}

	return result, nil
}

// ReadGenerationNumber reads and validates a PDF generation number from the input.
func ReadGenerationNumber(bytes core.InputBytes) (int64, error) {
	result, err := core.ReadInt(bytes)
	if err != nil {
		return 0, err
	}
	if result < 0 || result > generationNumberThreshold {
		return 0, fmt.Errorf("generation number '%d' has more than 5 digits", result)
	}

	return result, nil
}

// CreateObjectString creates the header string for an indirect PDF object.
func CreateObjectString(objectID, genID int64) string {
	return fmt.Sprintf("%d %d obj", objectID, genID)
}
