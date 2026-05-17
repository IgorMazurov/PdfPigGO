package fonts

import (
	"strings"
	"testing"
)

func TestGlyphListFactoryGetAdobeGlyphList(t *testing.T) {
	factory := GetGlyphListFactory()
	result, err := factory.Get("glyphlist")
	if err != nil {
		t.Fatalf("unexpected error loading glyphlist: %v", err)
	}

	h, err := result.NameToUnicode("H")
	if err != nil {
		t.Fatalf("NameToUnicode(H) error: %v", err)
	}

	if h != "H" {
		t.Errorf("expected NameToUnicode(H) = %q, got %q", "H", h)
	}
}

func TestGlyphListFactoryGetMissingResource(t *testing.T) {
	factory := GetGlyphListFactory()
	_, err := factory.Get("missing resource")

	if err == nil {
		t.Fatal("expected an error for missing resource name, got none")
	}
}

func TestGlyphListFactoryReadSkipsBlankLine(t *testing.T) {
	input := "# comment\n\none;0031"
	reader := strings.NewReader(input)

	factory := GetGlyphListFactory()
	result, err := factory.Read(reader)
	if err != nil {
		t.Fatalf("unexpected error reading glyph list: %v", err)
	}

	val, err := result.NameToUnicode("one")
	if err != nil {
		t.Fatalf("NameToUnicode(one) error: %v", err)
	}

	expected := "\u0031"
	if val != expected {
		t.Errorf("expected NameToUnicode(one) = %q (U+0031), got %q", expected, val)
	}
}

func TestGlyphListFactoryReadNilStream(t *testing.T) {
	factory := GetGlyphListFactory()
	_, err := factory.Read(nil)

	if err == nil {
		t.Fatal("expected an error for nil stream, got none")
	}
}

func TestGlyphListFactoryReadInvalidFormat(t *testing.T) {
	input := "one;0031\ntwelve;"
	reader := strings.NewReader(input)

	factory := GetGlyphListFactory()
	_, err := factory.Read(reader)

	if err == nil {
		t.Fatal("expected an error for invalid glyph list format, got none")
	}
}
