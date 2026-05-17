package testutil

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

// StringBytesConvertResult holds the result of converting a string to input bytes for testing.
type StringBytesConvertResult struct {
	First byte
	Bytes core.InputBytes
}

// ConvertStringToBytes converts a string to UTF-8 encoded MemoryInputBytes.
// If readFirst is true, it advances the cursor once and captures the first byte.
func ConvertStringToBytes(s string, readFirst bool) StringBytesConvertResult {
	input := core.NewMemoryInputBytes([]byte(s))

	var initialByte byte
	if readFirst {
		input.MoveNext()
		initialByte = input.CurrentByte()
	}

	return StringBytesConvertResult{
		First: initialByte,
		Bytes: input,
	}
}

// CreateStringScanner creates a CoreTokenScanner from a Latin-1 encoded string.
func CreateStringScanner(s string) (*tokenization.CoreTokenScanner, core.InputBytes) {
	inputBytes := core.NewMemoryInputBytes(core.StringAsLatin1Bytes(s))
	guard, _ := core.NewStackDepthGuard(256)
	scanner := tokenization.NewCoreTokenScanner(inputBytes, true, guard, tokenization.ScannerScopeNone, nil, false, false)

	return scanner, inputBytes
}
