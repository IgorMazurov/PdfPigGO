package writer

import (
	"fmt"
	"slices"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/core"
	standard14fonts "github.com/uglytoad/pdfpig/go/fonts/standard14_fonts"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/encodings"
)

func TestZapfDingbatsFontAddText(t *testing.T) {
	pdfBuilder := NewPdfDocumentBuilder()
	f1, err := pdfBuilder.AddStandard14Font(standard14fonts.ZapfDingbats)
	if err != nil {
		t.Fatalf("AddStandard14Font(ZapfDingbats): %v", err)
	}

	encodingTable := getSortedEncodingTable(encodings.ZapfDingbatsEncodingValue.Encoding)
	zapfGL, glErr := fonts.ZapfDingbats()
	if glErr != nil {
		t.Fatalf("ZapfDingbats glyph list: %v", glErr)
	}
	unicodeChars := getUnicodeCharacters(encodingTable, zapfGL)

	f2, err := pdfBuilder.AddStandard14Font(standard14fonts.TimesRoman)
	if err != nil {
		t.Fatalf("AddStandard14Font(TimesRoman): %v", err)
	}

	page, err := pdfBuilder.AddPageWithSize(content.MediaBoxA4.Bounds.TopRight.X, content.MediaBoxA4.Bounds.TopRight.Y)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	topPageY := page.PageSize().Bounds.TopRight.Y - 50
	cm := (page.PageSize().Bounds.TopRight.X / 8.5) / 2.54
	leftX := cm

	point := core.NewPdfPoint(leftX, topPageY)
	dateTimeStampPage(pdfBuilder, page, point, cm)

	letters, _ := page.AddText("Adobe Standard Font ZapfDingbats", 21, point, f2)
	maxH := maxLetterHeight(letters)
	newY := topPageY - maxH*1.2
	point = core.NewPdfPoint(leftX, newY)

	letters, _ = page.AddText("Font Specific encoding in Black (octal) and Unicode in Blue (hex)", 10, point, f2)
	maxH = maxLetterHeight(letters)
	newY -= maxH * 3
	point = core.NewPdfPoint(leftX, newY)

	charDetails := getCharacterDetails(t, page, f1, 12, unicodeChars)
	ctx := testContext{font: f1, page: page, fontName: "F1", fontLabel: f2, maxHeight: charDetails.maxH, maxWidth: charDetails.maxW}

	// Font specific character codes (in black)
	page.SetTextAndFillColor(0, 0, 0)
	for _, entry := range encodingTable {
		ch := rune(entry.code)
		point = addLetterWithContext(t, point, string(ch), ctx, true, false)
	}

	// Second set of rows for unicode characters
	newY -= charDetails.maxH * 1.2
	point = core.NewPdfPoint(leftX, newY)

	page.SetTextAndFillColor(0, 0, 200)
	for _, ch := range unicodeChars {
		point = addLetterWithContext(t, point, string(ch), ctx, false, true)
	}

	pdfBytes, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	doc, openErr := pdfpig.Open(pdfBytes, nil)
	if openErr != nil {
		t.Fatalf("Open: %v", openErr)
	}
	defer doc.Close()

	pageAny, pgErr := doc.GetPage(1)
	if pgErr != nil {
		t.Fatalf("GetPage(1): %v", pgErr)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	allLetters := pg.Letters()

	zapfBlackLetters := filterLetters(allLetters, func(l *content.Letter) bool {
		fn := l.FontName()
		return fn == "ZapfDingbats" && l.Color.ToRGBValues().B == 0
	})

	if len(zapfBlackLetters) != len(encodingTable) {
		t.Errorf("Expected %d ZapfDingbats black letters, got %d", len(encodingTable), len(zapfBlackLetters))
	}

	for i, letter := range zapfBlackLetters {
		if i >= len(encodingTable) {
			break
		}
		entry := encodingTable[i]
		unicodeStr, _ := zapfGL.NameToUnicode(entry.name)
		if len(unicodeStr) == 0 || len(letter.Value) == 0 {
			continue
		}
		letterRunes := []rune(letter.Value)
		expectedRunes := []rune(unicodeStr)
		if len(letterRunes) > 0 && len(expectedRunes) > 0 && letterRunes[0] != expectedRunes[0] {
			t.Errorf("Letter[%d]: expected %q (%c), got %q (%c)", i, unicodeStr, expectedRunes[0], letter.Value, letterRunes[0])
		}
	}

	zapfBlueLetters := filterLetters(allLetters, func(l *content.Letter) bool {
		fn := l.FontName()
		return fn == "ZapfDingbats" && l.Color.ToRGBValues().B > 0.78
	})

	if len(zapfBlueLetters) != len(unicodeChars) {
		t.Errorf("Expected %d ZapfDingbats blue letters, got %d", len(unicodeChars), len(zapfBlueLetters))
	}

	for i, letter := range zapfBlueLetters {
		if i >= len(unicodeChars) {
			break
		}
		expectedCh := unicodeChars[i]
		if len(letter.Value) == 0 {
			continue
		}
		letterRunes := []rune(letter.Value)
		if len(letterRunes) > 0 && letterRunes[0] != expectedCh {
			t.Errorf("Unicode letter[%d]: expected %q (%c), got %q (%c)", i, string(expectedCh), expectedCh, letter.Value, letterRunes[0])
		}
	}
}

func TestZapfDingbatsFontErrorResponseAddingInvalidText(t *testing.T) {
	pdfBuilder := NewPdfDocumentBuilder()
	f1, err := pdfBuilder.AddStandard14Font(standard14fonts.ZapfDingbats)
	if err != nil {
		t.Fatalf("AddStandard14Font(ZapfDingbats): %v", err)
	}

	encodingTable := getSortedEncodingTable(encodings.ZapfDingbatsEncodingValue.Encoding)

	page, err := pdfBuilder.AddPageWithSize(content.MediaBoxA4.Bounds.TopRight.X, content.MediaBoxA4.Bounds.TopRight.Y)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	cm := (page.PageSize().Bounds.TopRight.X / 8.5) / 2.54
	point := core.NewPdfPoint(cm, page.PageSize().Bounds.TopRight.Y-cm)

	codesFromTable := make(map[rune]bool)
	for _, entry := range encodingTable {
		codesFromTable[rune(entry.code)] = true
	}

	var invalidChars []rune
	for c := rune(0); c < 256; c++ {
		if !codesFromTable[c] {
			invalidChars = append(invalidChars, c)
		}
	}

	t.Logf("Number of invalid under-255 characters: %d", len(invalidChars))

	for _, ch := range invalidChars[:min(10, len(invalidChars))] {
		_, addErr := page.AddText(string(ch), 12, point, f1)
		if addErr == nil {
			t.Errorf("Expected error for character U+%04X in AddText", ch)
		} else if !containsStr(addErr.Error(), "does not contain a character") &&
			!containsStr(addErr.Error(), "does NOT have character") {
			t.Logf("AddText(U+%04X) error: %v", ch, addErr)
		}

		_, measureErr := page.MeasureText(string(ch), 12, point, f1)
		if measureErr == nil {
			t.Errorf("Expected error for character U+%04X in MeasureText", ch)
		} else if !containsStr(measureErr.Error(), "does not contain a character") &&
			!containsStr(measureErr.Error(), "does NOT have character") {
			t.Logf("MeasureText(U+%04X) error: %v", ch, measureErr)
		}
	}

	zapfGL, glErr := fonts.ZapfDingbats()
	if glErr != nil {
		t.Fatalf("ZapfDingbats glyph list: %v", glErr)
	}
	unicodeChars := getUnicodeCharacters(encodingTable, zapfGL)

	unicodeSet := make(map[rune]bool)
	for _, c := range unicodeChars {
		unicodeSet[c] = true
	}

	var invalidUnicode []rune
	for c := rune(0x2700); c <= 0x27BF; c++ {
		if !unicodeSet[c] {
			invalidUnicode = append(invalidUnicode, c)
		}
	}

	t.Logf("Number of invalid unicode dingbat characters: %d", len(invalidUnicode))

	for _, ch := range invalidUnicode[:min(10, len(invalidUnicode))] {
		_, addErr := page.AddText(string(ch), 12, point, f1)
		if addErr == nil {
			t.Logf("Character U+%04X was accepted (mapped through glyph list)", ch)
		} else if !containsStr(addErr.Error(), "does not contain a character") &&
			!containsStr(addErr.Error(), "does NOT have character") {
			t.Logf("AddText(U+%04X) error: %v", ch, addErr)
		}

		_, measureErr := page.MeasureText(string(ch), 12, point, f1)
		if measureErr == nil {
			t.Logf("Character U+%04X was accepted in MeasureText (mapped through glyph list)", ch)
		} else if !containsStr(measureErr.Error(), "does not contain a character") &&
			!containsStr(measureErr.Error(), "does NOT have character") {
			t.Logf("MeasureText(U+%04X) error: %v", ch, measureErr)
		}
	}
}

func TestSymbolFontAddText(t *testing.T) {
	pdfBuilder := NewPdfDocumentBuilder()
	f1, err := pdfBuilder.AddStandard14Font(standard14fonts.Symbol)
	if err != nil {
		t.Fatalf("AddStandard14Font(Symbol): %v", err)
	}

	encodingTable := getSortedEncodingTable(encodings.SymbolEncodingValue.Encoding)
	adobeGL, glErr := fonts.AdobeGlyphList()
	if glErr != nil {
		t.Fatalf("AdobeGlyphList: %v", glErr)
	}
	unicodeChars := getUnicodeCharacters(encodingTable, adobeGL)

	f2, err := pdfBuilder.AddStandard14Font(standard14fonts.TimesRoman)
	if err != nil {
		t.Fatalf("AddStandard14Font(TimesRoman): %v", err)
	}

	page, err := pdfBuilder.AddPageWithSize(content.MediaBoxA4.Bounds.TopRight.X, content.MediaBoxA4.Bounds.TopRight.Y)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	topPageY := page.PageSize().Bounds.TopRight.Y - 50
	cm := (page.PageSize().Bounds.TopRight.X / 8.5) / 2.54
	leftX := cm

	point := core.NewPdfPoint(leftX, topPageY)
	dateTimeStampPage(pdfBuilder, page, point, cm)

	letters, _ := page.AddText("Adobe Standard Font Symbol", 21, point, f2)
	maxH := maxLetterHeight(letters)
	newY := topPageY - maxH*1.2
	point = core.NewPdfPoint(leftX, newY)

	letters, _ = page.AddText("Font Specific encoding in Black (octal), Unicode in Blue (hex)", 10, point, f2)
	maxH = maxLetterHeight(letters)
	newY -= maxH * 3

	charDetails := getCharacterDetails(t, page, f1, 12, unicodeChars)
	ctx := testContext{font: f1, page: page, fontName: "F1", fontLabel: f2, maxHeight: charDetails.maxH, maxWidth: charDetails.maxW}

	// Font specific character codes
	newY -= charDetails.maxH
	point = core.NewPdfPoint(leftX, newY)

	page.SetTextAndFillColor(0, 0, 0)
	isBlack := true
	for _, entry := range encodingTable {
		code := entry.code

		switch code {
		case 0xac:
			code = 0x2190
		case 0xf7:
			code = 0xf8f7
		case 0xb5:
			code = 0x221d
		case 0xd7:
			code = 0x22c5
		}

		if code != entry.code && isBlack {
			page.SetTextAndFillColor(200, 0, 0)
			isBlack = false
		}
		if code == entry.code && !isBlack {
			page.SetTextAndFillColor(0, 0, 0)
			isBlack = true
		}

		ch := rune(code)
		point = addLetterWithContext(t, point, string(ch), ctx, isBlack, false)
	}

	newY -= charDetails.maxH * 1.2
	point = core.NewPdfPoint(leftX, newY)

	page.SetTextAndFillColor(0, 0, 200)
	for _, ch := range unicodeChars {
		point = addLetterWithContext(t, point, string(ch), ctx, false, true)
	}

	pdfBytes, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	doc, openErr := pdfpig.Open(pdfBytes, nil)
	if openErr != nil {
		t.Fatalf("Open: %v", openErr)
	}
	defer doc.Close()

	pageAny, pgErr := doc.GetPage(1)
	if pgErr != nil {
		t.Fatalf("GetPage(1): %v", pgErr)
	}
	pg, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	allLetters := pg.Letters()

	symbolLetters := filterLetters(allLetters, func(l *content.Letter) bool {
		fn := l.FontName()
		r := l.Color.ToRGBValues()
		return fn == "Symbol" && (r.B == 0 || r.R == 200)
	})

	if len(symbolLetters) != len(encodingTable) {
		t.Errorf("Expected %d Symbol letters, got %d", len(encodingTable), len(symbolLetters))
	}

	for i, letter := range symbolLetters {
		if i >= len(encodingTable) {
			break
		}
		entry := encodingTable[i]
		unicodeStr, _ := adobeGL.NameToUnicode(entry.name)
		if len(unicodeStr) == 0 || len(letter.Value) == 0 {
			continue
		}
		letterRunes := []rune(letter.Value)
		expectedRunes := []rune(unicodeStr)
		if len(letterRunes) > 0 && len(expectedRunes) > 0 && letterRunes[0] != expectedRunes[0] {
			t.Errorf("Symbol letter[%d]: expected %q (%c), got %q (%c)", i, unicodeStr, expectedRunes[0], letter.Value, letterRunes[0])
		}
	}

	symbolBlueLetters := filterLetters(allLetters, func(l *content.Letter) bool {
		fn := l.FontName()
		return fn == "Symbol" && l.Color.ToRGBValues().B > 0.78
	})

	if len(symbolBlueLetters) != len(unicodeChars) {
		t.Errorf("Expected %d Symbol blue letters, got %d", len(unicodeChars), len(symbolBlueLetters))
	}

	for i, letter := range symbolBlueLetters {
		if i >= len(unicodeChars) {
			break
		}
		expectedCh := unicodeChars[i]
		if len(letter.Value) == 0 {
			continue
		}
		letterRunes := []rune(letter.Value)
		if len(letterRunes) > 0 && letterRunes[0] != expectedCh {
			t.Errorf("Symbol unicode[%d]: expected %q (%c), got %q (%c)", i, string(expectedCh), expectedCh, letter.Value, letterRunes[0])
		}
	}
}

func TestSymbolFontErrorResponseAddingInvalidText(t *testing.T) {
	pdfBuilder := NewPdfDocumentBuilder()
	f1, err := pdfBuilder.AddStandard14Font(standard14fonts.Symbol)
	if err != nil {
		t.Fatalf("AddStandard14Font(Symbol): %v", err)
	}

	encodingTable := getSortedEncodingTable(encodings.SymbolEncodingValue.Encoding)

	page, err := pdfBuilder.AddPageWithSize(content.MediaBoxA4.Bounds.TopRight.X, content.MediaBoxA4.Bounds.TopRight.Y)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	cm := (page.PageSize().Bounds.TopRight.X / 8.5) / 2.54
	point := core.NewPdfPoint(cm, page.PageSize().Bounds.TopRight.Y-cm)

	codesFromTable := make(map[rune]bool)
	for _, entry := range encodingTable {
		codesFromTable[rune(entry.code)] = true
	}

	var invalidChars []rune
	for c := rune(0); c < 256; c++ {
		if !codesFromTable[c] {
			invalidChars = append(invalidChars, c)
		}
	}

	t.Logf("Number of invalid under-255 characters for Symbol: %d", len(invalidChars))

	for _, ch := range invalidChars[:min(10, len(invalidChars))] {
		_, addErr := page.AddText(string(ch), 12, point, f1)
		if addErr == nil {
			t.Errorf("Expected error for character U+%04X in AddText", ch)
		} else if !containsStr(addErr.Error(), "does not contain a character") &&
			!containsStr(addErr.Error(), "does NOT have character") {
			t.Logf("AddText(U+%04X) error: %v", ch, addErr)
		}

		_, measureErr := page.MeasureText(string(ch), 12, point, f1)
		if measureErr == nil {
			t.Errorf("Expected error for character U+%04X in MeasureText", ch)
		} else if !containsStr(measureErr.Error(), "does not contain a character") &&
			!containsStr(measureErr.Error(), "does NOT have character") {
			t.Logf("MeasureText(U+%04X) error: %v", ch, measureErr)
		}
	}

	adobeGL, glErr := fonts.AdobeGlyphList()
	if glErr != nil {
		t.Fatalf("AdobeGlyphList: %v", glErr)
	}
	unicodeChars := getUnicodeCharacters(encodingTable, adobeGL)

	unicodeSet := make(map[int]bool)
	for _, c := range unicodeChars {
		unicodeSet[int(c)] = true
	}

	var randomInvalid []rune
	count := 0
	for c := rune(1); count < 5 && c < 0x10ffff; c++ {
		if !unicodeSet[int(c)] && int(c) >= 0xd800 && int(c) <= 0xdfff {
			continue
		}
		if !unicodeSet[int(c)] {
			randomInvalid = append(randomInvalid, c)
			count++
		}
	}

	for _, ch := range randomInvalid {
		if int(ch) > 0x10ffff || (int(ch) >= 0xd800 && int(ch) <= 0xdfff) {
			continue
		}
		_, addErr := page.AddText(string(ch), 12, point, f1)
		if addErr == nil {
			t.Errorf("Expected error for random character U+%04X in AddText", ch)
		} else if !containsStr(addErr.Error(), "does not contain a character") &&
			!containsStr(addErr.Error(), "does NOT have character") {
			t.Logf("AddText(U+%04X) error: %v", ch, addErr)
		}

		_, measureErr := page.MeasureText(string(ch), 12, point, f1)
		if measureErr == nil {
			t.Errorf("Expected error for random character U+%04X in MeasureText", ch)
		} else if !containsStr(measureErr.Error(), "does not contain a character") &&
			!containsStr(measureErr.Error(), "does NOT have character") {
			t.Logf("MeasureText(U+%04X) error: %v", ch, measureErr)
		}
	}
}

func TestStandardFontsAddText(t *testing.T) {
	fontTypes := []standard14fonts.Standard14Font{
		standard14fonts.TimesRoman,
		standard14fonts.TimesBold,
		standard14fonts.TimesItalic,
		standard14fonts.TimesBoldItalic,
		standard14fonts.Helvetica,
		standard14fonts.HelveticaBold,
		standard14fonts.HelveticaOblique,
		standard14fonts.HelveticaBoldOblique,
		standard14fonts.Courier,
		standard14fonts.CourierBold,
		standard14fonts.CourierOblique,
		standard14fonts.CourierBoldOblique,
	}

	pdfBuilder := NewPdfDocumentBuilder()
	fontsAdded := make([]*AddedFont, len(fontTypes))
	for i, ft := range fontTypes {
		f, err := pdfBuilder.AddStandard14Font(ft)
		if err != nil {
			t.Fatalf("AddStandard14Font(%v): %v", ft, err)
		}
		fontsAdded[i] = f
	}

	encodingTable := getSortedEncodingTable(encodings.StandardEncodingValue.Encoding)
	adobeGL, glErr := fonts.AdobeGlyphList()
	if glErr != nil {
		t.Fatalf("AdobeGlyphList: %v", glErr)
	}
	unicodeChars := getUnicodeCharacters(encodingTable, adobeGL)

	for fontIdx, font := range fontsAdded {
		storedFont := pdfBuilder.Fonts()[font.Id]
		fontName := storedFont.FontProgram.Name()

		page, err := pdfBuilder.AddPageWithSize(content.MediaBoxA4.Bounds.TopRight.X, content.MediaBoxA4.Bounds.TopRight.Y)
		if err != nil {
			t.Fatalf("AddPageWithSize for %s: %v", fontName, err)
		}

		topPageY := page.PageSize().Bounds.TopRight.Y - 50
		cm := (page.PageSize().Bounds.TopRight.X / 8.5) / 2.54
		leftX := cm

		point := core.NewPdfPoint(leftX, topPageY)
		dateTimeStampPage(pdfBuilder, page, point, cm)

		letters, _ := page.AddText(fmt.Sprintf("Adobe Standard Font %s", fontName), 21, point, fontsAdded[0])
		maxH := maxLetterHeight(letters)
		newY := topPageY - maxH*1.2
		point = core.NewPdfPoint(leftX, newY)

		letters, _ = page.AddText("Font Specific encoding in Black, Unicode in Blue", 10, point, fontsAdded[0])
		maxH = maxLetterHeight(letters)
		newY -= maxH * 3
		point = core.NewPdfPoint(leftX, newY)

		charDetails := getCharacterDetails(t, page, fontsAdded[0], 12, unicodeChars)
		ctx := testContext{font: font, page: page, fontName: fmt.Sprintf("F%d", fontIdx+1), fontLabel: fontsAdded[0], maxHeight: charDetails.maxH, maxWidth: charDetails.maxW}

		page.SetTextAndFillColor(0, 0, 0)
		isBlack := true
		for _, entry := range encodingTable {
			code := entry.code

			switch code {
			case 0xc6:
				code = 0x02d8
			case 0xb4:
				code = 0x00b7
			case 0xb7:
				code = 0x2022
			case 0xb8:
				code = 0x201a
			case 0xa4:
				code = 0x2044
			case 0xa8:
				code = 0x00a4
			case 0x60:
				code = 0x2018
			case 0xaf:
				code = 0xfb02
			case 0xaa:
				code = 0x201c
			case 0xba:
				code = 0x201d
			case 0xf8:
				code = 0x0142
			case 0x27:
				code = 0x2019
			}

			if code != entry.code && isBlack {
				page.SetTextAndFillColor(200, 0, 0)
				isBlack = false
			}
			if code == entry.code && !isBlack {
				page.SetTextAndFillColor(0, 0, 0)
				isBlack = true
			}

			ch := rune(code)
			point = addLetterWithContext(t, point, string(ch), ctx, isBlack, false)
		}

	newY -= charDetails.maxH * 1.2
		point = core.NewPdfPoint(leftX, newY)

		page.SetTextAndFillColor(0, 0, 200)
		for _, ch := range unicodeChars {
			point = addLetterWithContext(t, point, string(ch), ctx, false, true)
		}
	}

	pdfBytes, err := pdfBuilder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	doc, openErr := pdfpig.Open(pdfBytes, nil)
	if openErr != nil {
		t.Fatalf("Open: %v", openErr)
	}
	defer doc.Close()

	numPages := doc.NumberOfPages()
	if numPages != len(fontTypes) {
		t.Errorf("Expected %d pages, got %d", len(fontTypes), numPages)
	}

	for pgNum := 1; pgNum <= numPages; pgNum++ {
		pageAny, pgErr := doc.GetPage(pgNum)
		if pgErr != nil {
			t.Errorf("GetPage(%d): %v", pgNum, pgErr)
			continue
		}
		pg, ok := pageAny.(*content.Page)
		if !ok {
			t.Errorf("expected *content.Page for page %d, got %T", pgNum, pageAny)
			continue
		}

		allLetters := pg.Letters()
		expectedFontName := ""
		for _, l := range allLetters {
			if l.FontSize == 12 {
				expectedFontName = l.FontName()
				break
			}
		}

		fontSpecificLetters := filterLetters(allLetters, func(l *content.Letter) bool {
			return l.FontName() == expectedFontName && l.FontSize == 12 &&
				(l.Color.ToRGBValues().B == 0 || l.Color.ToRGBValues().R == 200)
		})

		if len(fontSpecificLetters) != len(encodingTable) {
			t.Errorf("Page %d (%s): expected %d font-specific letters, got %d", pgNum, expectedFontName, len(encodingTable), len(fontSpecificLetters))
		}

		for i, letter := range fontSpecificLetters {
			if i >= len(encodingTable) {
				break
			}
			entry := encodingTable[i]
			unicodeStr, _ := adobeGL.NameToUnicode(entry.name)
			if len(unicodeStr) == 0 || len(letter.Value) == 0 {
				continue
			}
			letterRunes := []rune(letter.Value)
			expectedRunes := []rune(unicodeStr)
			if len(letterRunes) > 0 && len(expectedRunes) > 0 && letterRunes[0] != expectedRunes[0] {
				t.Errorf("Page %d letter[%d]: expected %q (%c), got %q (%c)", pgNum, i, unicodeStr, expectedRunes[0], letter.Value, letterRunes[0])
			}
		}

		unicodeLetters := filterLetters(allLetters, func(l *content.Letter) bool {
			return l.FontName() == expectedFontName && l.FontSize == 12 &&
				l.Color.ToRGBValues().B > 0.78
		})

		if len(unicodeLetters) != len(unicodeChars) {
			t.Errorf("Page %d: expected %d unicode letters, got %d", pgNum, len(unicodeChars), len(unicodeLetters))
		}

		for i, letter := range unicodeLetters {
			if i >= len(unicodeChars) {
				break
			}
			expectedCh := unicodeChars[i]
			if len(letter.Value) == 0 {
				continue
			}
			letterRunes := []rune(letter.Value)
			if len(letterRunes) > 0 && letterRunes[0] != expectedCh {
				t.Errorf("Page %d unicode[%d]: expected %q (%c), got %q (%c)", pgNum, i, string(expectedCh), expectedCh, letter.Value, letterRunes[0])
			}
		}
	}
}

func TestStandardFontErrorResponseAddingInvalidText(t *testing.T) {
	fontTypes := []standard14fonts.Standard14Font{
		standard14fonts.TimesRoman,
		standard14fonts.TimesBold,
		standard14fonts.TimesItalic,
		standard14fonts.TimesBoldItalic,
		standard14fonts.Helvetica,
		standard14fonts.HelveticaBold,
		standard14fonts.HelveticaOblique,
		standard14fonts.HelveticaBoldOblique,
		standard14fonts.Courier,
		standard14fonts.CourierBold,
		standard14fonts.CourierOblique,
		standard14fonts.CourierBoldOblique,
	}

	pdfBuilder := NewPdfDocumentBuilder()
	page, err := pdfBuilder.AddPageWithSize(content.MediaBoxA4.Bounds.TopRight.X, content.MediaBoxA4.Bounds.TopRight.Y)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	cm := (page.PageSize().Bounds.TopRight.X / 8.5) / 2.54
	point := core.NewPdfPoint(cm, page.PageSize().Bounds.TopRight.Y-cm)

	fontsAdded := make([]*AddedFont, len(fontTypes))
	for i, ft := range fontTypes {
		f, addErr := pdfBuilder.AddStandard14Font(ft)
		if addErr != nil {
			t.Fatalf("AddStandard14Font(%v): %v", ft, addErr)
		}
		fontsAdded[i] = f
	}

	encodingTable := getSortedEncodingTable(encodings.StandardEncodingValue.Encoding)

	codesFromTable := make(map[rune]bool)
	for _, entry := range encodingTable {
		codesFromTable[rune(entry.code)] = true
	}

	var invalidChars []rune
	for c := rune(0); c < 256; c++ {
		if !codesFromTable[c] {
			invalidChars = append(invalidChars, c)
		}
	}

	t.Logf("Number of invalid under-255 characters for standard fonts: %d", len(invalidChars))

	adobeGL, glErr := fonts.AdobeGlyphList()
	if glErr != nil {
		t.Fatalf("AdobeGlyphList: %v", glErr)
	}
	unicodeChars := getUnicodeCharacters(encodingTable, adobeGL)

	unicodeSet := make(map[int]bool)
	for _, c := range unicodeChars {
		unicodeSet[int(c)] = true
	}

	var randomInvalid []rune
	count := 0
	for c := rune(1); count < 5 && c < 0x10ffff; c++ {
		if !unicodeSet[int(c)] && !(int(c) >= 0xd800 && int(c) <= 0xdfff) {
			randomInvalid = append(randomInvalid, c)
			count++
		}
	}

	for fontIdx, font := range fontsAdded {
		storedFont := pdfBuilder.Fonts()[font.Id]
		fontName := storedFont.FontProgram.Name()

		for _, ch := range invalidChars[:min(5, len(invalidChars))] {
			_, addErr := page.AddText(string(ch), 12, point, font)
			if addErr == nil {
				t.Errorf("Font %s: expected error for character U+%04X in AddText", fontName, ch)
			} else if !containsStr(addErr.Error(), "does not contain a character") &&
				!containsStr(addErr.Error(), "does NOT have character") {
				t.Logf("Font %s: AddText(U+%04X) error: %v", fontName, ch, addErr)
			}

			_, measureErr := page.MeasureText(string(ch), 12, point, font)
			if measureErr == nil {
				t.Errorf("Font %s: expected error for character U+%04X in MeasureText", fontName, ch)
			} else if !containsStr(measureErr.Error(), "does not contain a character") &&
				!containsStr(measureErr.Error(), "does NOT have character") {
				t.Logf("Font %s: MeasureText(U+%04X) error: %v", fontName, ch, measureErr)
			}
		}

		for _, ch := range randomInvalid {
			if int(ch) > 0x10ffff || (int(ch) >= 0xd800 && int(ch) <= 0xdfff) {
				continue
			}
			_, addErr := page.AddText(string(ch), 12, point, font)
			if addErr == nil {
				t.Errorf("Font %s: expected error for random character U+%04X in AddText", fontName, ch)
			} else if !containsStr(addErr.Error(), "does not contain a character") &&
				!containsStr(addErr.Error(), "does NOT have character") {
				t.Logf("Font %s: AddText(U+%04X) error: %v", fontName, ch, addErr)
			}

			_, measureErr := page.MeasureText(string(ch), 12, point, font)
			if measureErr == nil {
				t.Errorf("Font %s: expected error for random character U+%04X in MeasureText", fontName, ch)
			} else if !containsStr(measureErr.Error(), "does not contain a character") &&
				!containsStr(measureErr.Error(), "does NOT have character") {
				t.Logf("Font %s: MeasureText(U+%04X) error: %v", fontName, ch, measureErr)
			}
		}

		_ = fontIdx
	}
}

type encodingEntry struct {
	code int
	name string
}

func getSortedEncodingTable(enc *encodings.Encoding) []encodingEntry {
	result := make([]encodingEntry, 0, len(enc.CodeToName))
	for code, name := range enc.CodeToName {
		result = append(result, encodingEntry{code: code, name: name})
	}
	slices.SortFunc(result, func(a, b encodingEntry) int {
		return a.code - b.code
	})
	return result
}

func getUnicodeCharacters(encodingTable []encodingEntry, glyphList *fonts.GlyphList) []rune {
	result := make([]rune, 0, len(encodingTable))
	for _, entry := range encodingTable {
		unicodeStr, err := glyphList.NameToUnicode(entry.name)
		if err != nil || unicodeStr == "" {
			continue
		}
		runes := []rune(unicodeStr)
		if len(runes) > 0 {
			result = append(result, runes[0])
		}
	}
	return result
}

type testContext struct {
	font      *AddedFont
	page      *PdfPageBuilder
	fontName  string
	fontLabel *AddedFont
	maxHeight float64
	maxWidth  float64
}

func addLetterWithContext(t *testing.T, point core.PdfPoint, textToAdd string, ctx testContext, isOctalLabel, isHexLabel bool) core.PdfPoint {
	return addLetter(t, ctx.page, point, textToAdd, ctx.font, ctx.fontName, ctx.fontLabel, ctx.maxHeight, ctx.maxWidth, isOctalLabel, isHexLabel)
}

func addLetter(t *testing.T, page *PdfPageBuilder, point core.PdfPoint, textToAdd string, font *AddedFont, fontName string, fontLabel *AddedFont, maxCharHeight, maxCharWidth float64, isOctalLabel, isHexLabel bool) core.PdfPoint {
	runes := []rune(textToAdd)
	if len(runes) == 0 || len(runes) > 1 {
		t.Fatalf("text to add must be a single letter, got %q (runes: %d)", textToAdd, len(runes))
	}

	letters, err := page.AddText(textToAdd, 12, point, font)
	if err != nil {
		t.Fatalf("AddText(%q): %v", textToAdd, err)
	}

	if len(letters) == 0 {
		t.Fatalf("expected at least one letter for %q", textToAdd)
	}

	if isOctalLabel {
		labelPointSize := float64(5)
		octalStr := fmt.Sprintf("%03o", int(runes[0]))
		codeMidPoint := point.X + letters[0].BoundingBox.Width/2
		measuredLetters, _ := page.MeasureText(octalStr, labelPointSize, point, fontLabel)
		labelMaxH := maxLetterHeight(measuredLetters)
		totalW := 0.0
		for _, ml := range measuredLetters {
			totalW += ml.BoundingBox.Width
		}
		labelY := point.Y + labelMaxH*0.1 + maxCharHeight
		xLabel := codeMidPoint - (totalW / 2)
		labelPoint := core.NewPdfPoint(xLabel, labelY)
		page.AddText(octalStr, labelPointSize, labelPoint, fontLabel)
	}

	if isHexLabel {
		labelPointSize := float64(3)
		hexStr := fmt.Sprintf("0x%04X", int(runes[0]))
		codeMidPoint := point.X + letters[0].BoundingBox.Width/2
		measuredLetters, _ := page.MeasureText(hexStr, labelPointSize, point, fontLabel)
		labelMaxH := maxLetterHeight(measuredLetters)
		totalW := 0.0
		for _, ml := range measuredLetters {
			totalW += ml.BoundingBox.Width
		}
		labelY := point.Y - labelMaxH*2.5
		xLabel := codeMidPoint - (totalW / 2)
		labelPoint := core.NewPdfPoint(xLabel, labelY)
		page.AddText(hexStr, labelPointSize, labelPoint, fontLabel)
	}

	cm := (page.PageSize().Bounds.TopRight.X / 8.5) / 2.54
	newX := point.X + maxCharWidth*1.1
	newY := point.Y

	if newX > page.PageSize().Bounds.TopRight.X-cm {
		return core.PdfPoint{
			X: cm,
			Y: newY - maxCharHeight*5,
		}
	}

	return core.NewPdfPoint(newX, newY)
}

func dateTimeStampPage(pdfBuilder *PdfDocumentBuilder, page *PdfPageBuilder, point core.PdfPoint, cm float64) {
	courierFont, err := pdfBuilder.AddStandard14Font(standard14fonts.Courier)
	if err != nil {
		return
	}

	fontSize := float64(7)
	indentFromLeft := page.PageSize().Bounds.TopRight.X - cm

	_, _ = page.MeasureText("UTC: ", fontSize, point, courierFont)
	_, _ = page.MeasureText("Local: ", fontSize, point, courierFont)

	point = core.NewPdfPoint(indentFromLeft, point.Y)
	letters, _ := page.AddText(fmt.Sprintf("UTC: %d", 0), fontSize, point, courierFont)
	if len(letters) > 0 {
		maxHeight := maxLetterHeight(letters)
		point = core.NewPdfPoint(indentFromLeft, point.Y-maxHeight*1.2)
	}

	page.AddText(fmt.Sprintf("Local: %d", 0), fontSize, point, courierFont)
}

func getCharacterDetails(t *testing.T, page *PdfPageBuilder, font *AddedFont, fontSize float64, unicodeChars []rune) struct {
	maxH float64
	maxW float64
} {
	point := core.PdfPoint{X: 10, Y: 10}
	var maxH, maxW float64
	for _, ch := range unicodeChars {
		letters, err := page.MeasureText(string(ch), fontSize, point, font)
		if err != nil || len(letters) == 0 {
			continue
		}
		bb := letters[0].BoundingBox
		if bb.Height > maxH {
			maxH = bb.Height
		}
		if bb.Width > maxW {
			maxW = bb.Width
		}
	}
	return struct{ maxH, maxW float64 }{maxH: maxH, maxW: maxW}
}

func maxLetterHeight(letters []*content.Letter) float64 {
	var maxH float64
	for _, l := range letters {
		if l.BoundingBox.Height > maxH {
			maxH = l.BoundingBox.Height
		}
	}
	return maxH
}

func filterLetters(all []*content.Letter, predicate func(*content.Letter) bool) []*content.Letter {
	result := make([]*content.Letter, 0)
	for _, l := range all {
		if predicate(l) {
			result = append(result, l)
		}
	}
	return result
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
