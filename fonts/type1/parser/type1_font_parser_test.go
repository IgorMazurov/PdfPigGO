package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/type1/parser"
)

func getFileBytes(t *testing.T, name string) []byte {
	t.Helper()

	testDataDir := "testdata"
	files, err := os.ReadDir(testDataDir)
	if err != nil {
		t.Fatalf("could not read testdata directory: %v", err)
	}

	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(name)) {
			path := filepath.Join(testDataDir, f.Name())
			buf, readErr := os.ReadFile(path)
			if readErr != nil {
				t.Fatalf("could not read test file %q: %v", path, readErr)
			}
			return buf
		}
	}

	t.Fatalf("could not find test file %q in folder %q", name, testDataDir)
	return nil
}

func TestCanReadHexEncryptedPortion(t *testing.T) {
	bytes := getFileBytes(t, "AdobeUtopia.pfa")

	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}

	_, parseErr := parser.Parse(core.NewMemoryInputBytes(bytes), 0, 0, guard)
	if parseErr != nil {
		t.Errorf("Parse failed: %v", parseErr)
	}
}

func TestCanReadBinaryEncryptedPortionOfFullPfb(t *testing.T) {
	bytes := getFileBytes(t, "Raleway-Black.pfb")

	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}

	_, parseErr := parser.Parse(core.NewMemoryInputBytes(bytes), 0, 0, guard)
	if parseErr != nil {
		t.Errorf("Parse failed: %v", parseErr)
	}
}

func TestCanReadCharStrings(t *testing.T) {
	bytes := getFileBytes(t, "CMBX10.pfa")

	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}

	_, parseErr := parser.Parse(core.NewMemoryInputBytes(bytes), 0, 0, guard)
	if parseErr != nil {
		t.Errorf("Parse failed: %v", parseErr)
	}
}

func TestCanReadEncryptedPortion(t *testing.T) {
	bytes := getFileBytes(t, "CMCSC10")

	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}

	_, parseErr := parser.Parse(core.NewMemoryInputBytes(bytes), 0, 0, guard)
	if parseErr != nil {
		t.Errorf("Parse failed: %v", parseErr)
	}
}

func TestCanReadAsciiPart(t *testing.T) {
	bytes := getFileBytes(t, "CMBX12")

	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}

	_, parseErr := parser.Parse(core.NewMemoryInputBytes(bytes), 0, 0, guard)
	if parseErr != nil {
		t.Errorf("Parse failed: %v", parseErr)
	}
}

func TestOutputCmbx10Svgs(t *testing.T) {
	bytes := getFileBytes(t, "CMBX10")

	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}

	result, parseErr := parser.Parse(core.NewMemoryInputBytes(bytes), 0, 0, guard)
	if parseErr != nil {
		t.Fatalf("Parse failed: %v", parseErr)
	}

	for charName := range result.CharStrings.CharStrings() {
		_, ok := result.CharStrings.TryGenerate(charName)
		if !ok {
			t.Errorf("TryGenerate failed for character %q", charName)
		}
	}
}

func TestCanReadFontWithCommentsInOtherSubrs(t *testing.T) {
	bytes := getFileBytes(t, "CMR10")

	guard, err := core.NewStackDepthGuard(256)
	if err != nil {
		t.Fatalf("unexpected error creating StackDepthGuard: %v", err)
	}

	_, parseErr := parser.Parse(core.NewMemoryInputBytes(bytes), 0, 0, guard)
	if parseErr != nil {
		t.Errorf("Parse failed: %v", parseErr)
	}
}
