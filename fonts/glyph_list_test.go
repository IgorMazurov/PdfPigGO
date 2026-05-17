package fonts

import (
	"testing"
)

func TestCanLoadAdobeGlyphList(t *testing.T) {
	list, err := AdobeGlyphList()
	if err != nil {
		t.Fatalf("unexpected error loading Adobe Glyph List: %v", err)
	}

	result, err := list.NameToUnicode("Acute")
	if err != nil {
		t.Fatalf("NameToUnicode(Acute) error: %v", err)
	}

	expected := "\uF6C9"
	if result != expected {
		t.Errorf("expected NameToUnicode(Acute) = %q, got %q", expected, result)
	}
}

func TestCanLoadZapfDingbatsGlyphList(t *testing.T) {
	list, err := ZapfDingbats()
	if err != nil {
		t.Fatalf("unexpected error loading Zapf Dingbats: %v", err)
	}

	result, err := list.NameToUnicode("a69")
	if err != nil {
		t.Fatalf("NameToUnicode(a69) error: %v", err)
	}

	expected := "\u274A"
	if result != expected {
		t.Errorf("expected NameToUnicode(a69) = %q, got %q", expected, result)
	}
}

func TestUnicodeToNameWorks(t *testing.T) {
	list, err := AdobeGlyphList()
	if err != nil {
		t.Fatalf("unexpected error loading Adobe Glyph List: %v", err)
	}

	name := list.UnicodeCodePointToName(79)

	expected := "O"
	if name != expected {
		t.Errorf("expected UnicodeCodePointToName(79) = %q, got %q", expected, name)
	}
}

func TestUnicodeToNameNotDefined(t *testing.T) {
	list := NewGlyphList(map[string]string{})

	name := list.UnicodeCodePointToName(120)

	expected := NotDefined
	if name != expected {
		t.Errorf("expected UnicodeCodePointToName(120) = %q, got %q", expected, name)
	}
}

func TestNameToUnicodeEmpty(t *testing.T) {
	list := NewGlyphList(map[string]string{})

	result, err := list.NameToUnicode("")
	if err != nil {
		t.Fatalf("NameToUnicode(empty) error: %v", err)
	}

	if result != "" {
		t.Errorf("expected NameToUnicode(\"\") = \"\", got %q", result)
	}
}

func TestNameToUnicodeRemovesSuffix(t *testing.T) {
	list := NewGlyphList(map[string]string{
		"Boris": "B",
	})

	result, err := list.NameToUnicode("Boris.Special")
	if err != nil {
		t.Fatalf("NameToUnicode(Boris.Special) error: %v", err)
	}

	expected := "B"
	if result != expected {
		t.Errorf("expected NameToUnicode(Boris.Special) = %q, got %q", expected, result)
	}
}

func TestNameToUnicodeConvertsHexAndUsesHexValue(t *testing.T) {
	list := NewGlyphList(map[string]string{
		"B": "X",
	})

	result, err := list.NameToUnicode("uni0042")
	if err != nil {
		t.Fatalf("NameToUnicode(uni0042) error: %v", err)
	}

	expected := "B"
	if result != expected {
		t.Errorf("expected NameToUnicode(uni0042) = %q, got %q", expected, result)
	}
}

func TestNameToUnicodeConvertsShortHexAndUsesHexValue(t *testing.T) {
	list := NewGlyphList(map[string]string{
		"E": "\u00C6",
	})

	result, err := list.NameToUnicode("u0045")
	if err != nil {
		t.Fatalf("NameToUnicode(u0045) error: %v", err)
	}

	expected := "E"
	if result != expected {
		t.Errorf("expected NameToUnicode(u0045) = %q, got %q", expected, result)
	}
}

func TestNameToUnicodeConvertAglSpecification(t *testing.T) {
	list := NewGlyphList(map[string]string{
		"Lcommaaccent": "\u013B",
	})

	result, err := list.NameToUnicode("Lcommaaccent_uni20AC0308_u1040C.alternate")
	if err != nil {
		t.Fatalf("NameToUnicode error: %v", err)
	}

	expected := "\u013B\u20AC\u0308\U0001040C"
	if result != expected {
		t.Errorf("expected NameToUnicode = %q (len=%d), got %q (len=%d)", expected, len(expected), result, len(result))
	}

	runes := []rune(result)
	if len(runes) != 4 {
		t.Fatalf("expected 4 runes, got %d", len(runes))
	}

	tests := []struct {
		idx      int
		expected rune
	}{
		{0, 0x013B},
		{1, 0x20AC},
		{2, 0x0308},
		{3, 0x1040C},
	}

	for _, tc := range tests {
		if runes[tc.idx] != tc.expected {
			t.Errorf("expected rune[%d] = U+%04X, got U+%04X", tc.idx, tc.expected, runes[tc.idx])
		}
	}
}
