package encryption

import (
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
)

func TestAesDecryptDateString(t *testing.T) {
	key := []byte{
		54, 109, 249, 186, 109, 210, 209, 44, 94, 28, 227, 232, 73, 86, 128, 186,
	}

	data := []byte{
		123, 28, 227, 85, 79, 126, 149, 28, 211, 96, 199, 192, 105, 149, 76, 231,
		8, 136, 51, 141, 139, 44, 0, 230, 228, 116, 12, 145, 132, 157, 5, 123, 235, 247, 232, 244, 36, 217, 73, 147, 157, 124, 27, 143, 255, 79, 220, 194,
	}

	output, err := AesDecrypt(data, key)
	if err != nil {
		t.Fatalf("AesDecrypt error: %v", err)
	}

	actual := core.BytesAsLatin1String(output)

	expected := "D:20180808103317-07'00'"
	if actual != expected {
		t.Errorf("decrypt mismatch:\nexpected: %q\n  actual: %q", expected, actual)
	}
}
