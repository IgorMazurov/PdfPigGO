// Package fonts provides writing font implementations for PDF generation.
package fonts

import (
	"log"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/adobe_font_metrics"
	"github.com/uglytoad/pdfpig/go/fonts/encodings"
	fontsPkg "github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer/interfaces"
)

// Standard14WritingFont represents one of the 14 standard PDF fonts for writing.
type Standard14WritingFont struct {
	metrics adobe_font_metrics.AdobeFontMetrics
}

// NewStandard14WritingFont creates a new Standard14WritingFont from Adobe font metrics.
func NewStandard14WritingFont(metrics adobe_font_metrics.AdobeFontMetrics) *Standard14WritingFont {
	return &Standard14WritingFont{
		metrics: metrics,
	}
}

// HasWidths reports whether the font has explicit character widths.
// Standard 14 fonts rely on built-in width data so this always returns false.
func (f *Standard14WritingFont) HasWidths() bool {
	return false
}

// Name returns the PostScript name of the font.
func (f *Standard14WritingFont) Name() string {
	return f.metrics.FontName
}

// TryGetBoundingBox attempts to get the bounding box for a character.
// Returns the bounding box and true if found, or nil and false otherwise.
func (f *Standard14WritingFont) TryGetBoundingBox(character rune) (*core.PdfRectangle, bool) {
	code := f.codeMapIfUnicode(character)
	if code == -1 {
		log.Printf("Font '%s' does NOT have character %c (0x%X).", f.metrics.FontName, character, character)
		return nil, false
	}

	charMetric := f.findCharacterMetric(code)
	if charMetric == nil {
		log.Printf("Font '%s' does NOT have character %c (0x%X).", f.metrics.FontName, character, character)
		return nil, false
	}

	bbox := charMetric.BoundingBox
	rect := core.NewPdfRectangleFloat(
		bbox.Left(),
		bbox.Bottom(),
		bbox.Left()+charMetric.Width.X,
		bbox.Top(),
	)

	return &rect, true
}

// TryGetAdvanceWidth attempts to get the advance width for a character.
// Returns the width and true if found, or 0 and false otherwise.
func (f *Standard14WritingFont) TryGetAdvanceWidth(character rune) (float64, bool) {
	bbox, ok := f.TryGetBoundingBox(character)
	if !ok {
		return 0, false
	}

	return bbox.Width, true
}

// GetFontMatrix returns the transformation matrix for the font.
// Standard 14 fonts use a scale factor of 1/1000 to convert from font units.
func (f *Standard14WritingFont) GetFontMatrix() core.TransformationMatrix {
	return core.FromValues(1/1000.0, 0, 0, 1/1000.0, 0, 0)
}

// WriteFont writes the font dictionary to the stream and returns an indirect reference.
// If reservedIndirect is non-nil, the font is written at that pre-reserved object number.
func (f *Standard14WritingFont) WriteFont(w interfaces.PdfStreamWriter, reservedIndirect *tokens.IndirectReferenceToken) *tokens.IndirectReferenceToken {
	encoding := tokens.StandardEncoding
	fontNameLower := f.metrics.FontName
	if fontNameLower == "Symbol" || fontNameLower == "symbol" ||
		fontNameLower == "ZapfDingbats" || fontNameLower == "zapfdingbats" {
		encoding = tokens.MustCreate("FontSpecific")
	}

	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Type:     tokens.Font,
		tokens.Subtype:  tokens.Type1,
		tokens.BaseFont: tokens.MustCreate(f.metrics.FontName),
		tokens.Encoding: encoding,
	})

	if reservedIndirect != nil {
		return w.WriteTokenAt(dict, reservedIndirect)
	}

	return w.WriteToken(dict)
}

// GetValueForCharacter returns the byte value for a character in the font's encoding.
// Panics if the character is not available in the font.
func (f *Standard14WritingFont) GetValueForCharacter(character rune) byte {
	characterCode := f.codeMapIfUnicode(character)

	charMetric := f.findCharacterMetric(characterCode)
	if charMetric == nil {
		panic(fontsPkg.NewInvalidFontFormatException(
			"Font '" + f.metrics.FontName + "' does NOT have character " + string(character) + " (0x" + formatHex(int(character)) + ")."))
	}

	return byte(charMetric.CharacterCode)
}

// findCharacterMetric searches the font metrics for a character with the given code.
func (f *Standard14WritingFont) findCharacterMetric(code int) *adobe_font_metrics.AdobeFontMetricsIndividualCharacterMetric {
	for _, metric := range f.metrics.CharacterMetrics {
		if metric.CharacterCode == code {
			return &metric
		}
	}
	return nil
}

// unicodeToSymbolCode maps a Unicode character to its Symbol font encoding code.
func (f *Standard14WritingFont) unicodeToSymbolCode(character rune) int {
	gl, err := fontsPkg.AdobeGlyphList()
	if err != nil || gl == nil {
		return -1
	}

	name := gl.UnicodeCodePointToName(int(character))
	if name == fontsPkg.NotDefined {
		return -1
	}

	code := encodings.SymbolEncodingValue.GetCode(name)
	if code == -1 {
		log.Printf("Found Unicode point %c (0x%X) but glyph name '%s' not found in font '%s' [Symbol] (StandardEncoding).",
			character, character, name, f.metrics.FontName)
	}
	return code
}

// unicodeToZapfDingbats maps a Unicode character to its ZapfDingbats encoding code.
func (f *Standard14WritingFont) unicodeToZapfDingbats(character rune) int {
	gl, err := fontsPkg.ZapfDingbats()
	if err != nil || gl == nil {
		log.Printf("Failed to find Unicode character %c (0x%X).", character, character)
		return -1
	}

	name := gl.UnicodeCodePointToName(int(character))
	if name == fontsPkg.NotDefined {
		log.Printf("Failed to find Unicode character %c (0x%X).", character, character)
		return -1
	}

	code := encodings.ZapfDingbatsEncodingValue.GetCode(name)
	if code == -1 {
		log.Printf("Found Unicode point %c (0x%X) but glyph name '%s' not found in font '%s' (font specific encoding: ZapfDingbats).",
			character, character, name, f.metrics.FontName)
	}
	return code
}

// unicodeToStandardEncoding maps a Unicode character to its StandardEncoding code.
func (f *Standard14WritingFont) unicodeToStandardEncoding(character rune) int {
	gl, err := fontsPkg.AdobeGlyphList()
	if err != nil || gl == nil {
		log.Printf("Failed to find Unicode character %c (0x%X).", character, character)
		return -1
	}

	name := gl.UnicodeCodePointToName(int(character))
	if name == fontsPkg.NotDefined {
		log.Printf("Failed to find Unicode character %c (0x%X).", character, character)
		return -1
	}

	code := encodings.StandardEncodingValue.GetCode(name)
	if code == -1 {
		nameCapitalisedChange := toggleFirstCase(name)
		code = encodings.StandardEncodingValue.GetCode(nameCapitalisedChange)
		if code == -1 {
			log.Printf("Found Unicode point %c (0x%X) but glyph name '%s' not found in font '%s' (StandardEncoding).",
				character, character, name, f.metrics.FontName)
		}
	}
	return code
}

// codeMapIfUnicode maps a Unicode character to its encoding-specific code.
func (f *Standard14WritingFont) codeMapIfUnicode(character rune) int {
	i := int(character)
	fontName := f.metrics.FontName

	if fontName == "ZapfDingbats" || fontName == "zapfdingbats" {
		return codeOrMap(i, func() int { return f.unicodeToZapfDingbats(character) })
	} else if fontName == "Symbol" || fontName == "symbol" {
		switch i {
		case 0x00AC:
			log.Println("Warning: 0x00AC used as Unicode ('¬') (logicalnot). For (arrowleft)('←') from Adobe Symbol Font Specific (0330) use Unicode 0x2190 ('←').")
			return 0x00D8
		case 0x00F7:
			log.Println("Warning: 0x00F7 used as Unicode ('÷')(divide). For (parenrightex) from Adobe Symbol Font Specific (0367) use Unicode 0xF8F7.")
			return 0x00B8
		case 0x00B5:
			log.Println("Warning: 0x00B5 used as Unicode divide ('µ')(mu). For (proportional)('∝') from Adobe Symbol Font Specific (0265) use Unicode 0x221D('∝').")
			return 0x006D
		case 0x00D7:
			log.Println("Warning: 0x00D7 used as Unicode multiply ('×')(multiply). For (dotmath)('⋅') from Adobe Symbol Font Specific (0327) use Unicode 0x22C5('⋅').")
			return 0x00B4
		}
		return codeOrMap(i, func() int { return f.unicodeToSymbolCode(character) })
	}

	switch i {
	case 0x00C6:
		log.Println("Warning: 0x00C6 used as Unicode ('Æ') (AE). For (breve)('˘') from Adobe Standard Font Specific (0306) use Unicode 0x02D8 ('˘').")
		return 0x00E1
	case 0x00B4:
		log.Println("Warning: 0x00B4 used as Unicode ('´') (acute). For (periodcentered)('·') from Adobe Standard Font Specific (0264) use Unicode 0x00B7  ('·').")
		return 0x00C2
	case 0x00B7:
		log.Println("Warning: 0x00B7 used as Unicode ('·') (periodcentered). For (bullet)('•') from Adobe Standard Font Specific (0267) use Unicode 0x2022 ('•').")
		return 0x00B4
	case 0x00B8:
		log.Println("Warning: 0x00B8 used as Unicode ('¸') (cedilla). For (quotesinglbase)('‚') from Adobe Standard Font Specific (0267) use Unicode 0x201A ('‚').")
		return 0x00CB
	case 0x00A4:
		log.Println("Warning: 0x00A4 used as Unicode (currency). For (fraction) ('⁄') from Adobe Standard Font Specific (0244) use Unicode 0x2044 ('⁄').")
		return 0x00A8
	case 0x00A8:
		log.Println("Warning: 0x00A8 used as Unicode (dieresis)('¨'). For (currency) from Adobe Standard Font Specific (0250) use Unicode 0x00A4.")
		return 0x00C8
	case 0x0060:
		log.Println("Warning: 0x0060 used as Unicode (grave)('`'). For (quoteleft)('‘') from Adobe Standard Font Specific (0140) use Unicode 0x2018.")
		return 0x00C1
	case 0x00AF:
		log.Println("Warning: 0x00AF used as Unicode (macron)('¯'). For (fl)('ﬂ') from Adobe Standard Font Specific (0257) use Unicode 0xFB02.")
		return 0x00C5
	case 0x00AA:
		log.Println(`Warning: 0x00AA used as Unicode (ordfeminine)('ª'). For (quotedblleft) ('\u201C') from Adobe Standard Font Specific (0252) use Unicode 0x201C.`)
		return 0x00E3
	case 0x00BA:
		log.Println(`Warning: 0x00BA used as Unicode (ordmasculine)('º'). For (quotedblright) ('\u201D') from Adobe Standard Font Specific (0272) use Unicode 0x201D.`)
		return 0x00EB
	case 0x00F8:
		log.Println("Warning: 0x00F8 used as Unicode (oslash)('ø'). For (lslash) ('ł') from Adobe Standard Font Specific (0370) use Unicode 0x0142.")
		return 0x00F9
	case 0x0027:
		log.Println(`Warning: 0x0027 used as Unicode (quotesingle)("'"). For (quoteright) ('\u2019') from Adobe Standard Font Specific (0047) use Unicode 0x2019.`)
		return 0x00A9
	}

	if f.characterCodeInMetrics(i) {
		return i
	}

	return f.unicodeToStandardEncoding(character)
}

// codeOrMap returns the character code directly if < 255, otherwise calls mapper.
func codeOrMap(code int, mapper func() int) int {
	if code < 255 {
		return code
	}
	return mapper()
}

// characterCodeInMetrics checks whether any character metric has the given code.
func (f *Standard14WritingFont) characterCodeInMetrics(code int) bool {
	for _, metric := range f.metrics.CharacterMetrics {
		if metric.CharacterCode == code {
			return true
		}
	}
	return false
}

// toggleFirstCase toggles the case of the first character of a string.
func toggleFirstCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'A' && r[0] <= 'Z' {
		r[0] = r[0] + 32
	} else if r[0] >= 'a' && r[0] <= 'z' {
		r[0] = r[0] - 32
	}
	return string(r)
}

// formatHex formats an integer as a two-digit uppercase hexadecimal string.
func formatHex(v int) string {
	const hexChars = "0123456789ABCDEF"
	return string(hexChars[(v>>4)&0xF]) + string(hexChars[v&0xF])
}

var _ interfaces.WritingFont = (*Standard14WritingFont)(nil)
