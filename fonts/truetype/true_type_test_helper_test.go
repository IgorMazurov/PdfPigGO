package truetype_test

import (
	"os"
	"strings"
	"testing"
)

// GetFileBytes reads a TrueType font or text file from the testdata directory.
// If the filename does not end with .ttf or .txt, .ttf is appended automatically.
func GetFileBytes(t *testing.T, name string) []byte {
	t.Helper()

	if !strings.HasSuffix(name, ".ttf") && !strings.HasSuffix(name, ".txt") {
		name += ".ttf"
	}

	path := "testdata/" + name

	buf, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read test file %q: %v", path, err)
	}

	return buf
}
