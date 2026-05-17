package parser_test

import (
	"strings"
	"testing"

	cmapresources "github.com/uglytoad/pdfpig/go/fonts/resources/cmap"
	"github.com/uglytoad/pdfpig/go/fonts/parser"
)

// TestCanParseAllPredefinedCMaps verifies that all embedded CMap files can be parsed.
func TestCanParseAllPredefinedCMaps(t *testing.T) {
	p := parser.NewCMapParser()

	entries, err := cmapresources.Files.ReadDir(".")
	if err != nil {
		t.Fatalf("failed to read embedded CMap directory: %v", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".go") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			cmap, found := p.TryParseExternal(name)
			if !found {
				t.Fatalf("embedded CMap not found: %s", name)
			}
			if cmap == nil {
				t.Fatalf("CMap is nil for: %s", name)
			}
		})
	}
}

// TestCanParseIdentityHorizontalCMap verifies the Identity-H CMap parses correctly.
func TestCanParseIdentityHorizontalCMap(t *testing.T) {
	p := parser.NewCMapParser()

	cmap, found := p.TryParseExternal("Identity-H")
	if !found {
		t.Fatal("embedded CMap not found: Identity-H")
	}

	ranges := cmap.CodespaceRanges()
	if len(ranges) != 1 {
		t.Fatalf("expected 1 codespace range, got %d", len(ranges))
	}

	range0 := ranges[0]
	if range0.StartInt != 0 {
		t.Errorf("StartInt = %d, want 0", range0.StartInt)
	}
	if range0.EndInt != 65535 {
		t.Errorf("EndInt = %d, want 65535", range0.EndInt)
	}
	if range0.CodeLength != 2 {
		t.Errorf("CodeLength = %d, want 2", range0.CodeLength)
	}

	cidRanges := cmap.CidRanges()
	if len(cidRanges) != 256 {
		t.Errorf("CidRanges count = %d, want 256", len(cidRanges))
	}

	if cmap.Version() == nil || *cmap.Version() != "10.003" {
		t.Errorf("Version = %v, want \"10.003\"", cmap.Version())
	}
}
