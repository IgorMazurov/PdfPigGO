package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/adobe_font_metrics"
)

func TestCanParseAfmFile(t *testing.T) {
	courierAfmSnippet := `StartFontMetrics 4.1

Comment Copyright (c) 1989, 1990, 1991, 1992, 1993, 1997 Adobe Systems Incorporated.  All Rights Reserved.

Comment Creation Date: Thu May  1 17:27:09 1997

Comment UniqueID 43050

Comment VMusage 39754 50779

FontName Courier

FullName Courier

FamilyName Courier

Weight Medium

ItalicAngle 0

IsFixedPitch true

CharacterSet ExtendedRoman

FontBBox -23 -250 715 805 

UnderlinePosition -100

UnderlineThickness 50

Version 003.000

Notice Copyright (c) 1989, 1990, 1991, 1992, 1993, 1997 Adobe Systems Incorporated.  All Rights Reserved.

EncodingScheme AdobeStandardEncoding

CapHeight 562

XHeight 426

Ascender 629

Descender -157

StdHW 51

StdVW 51

StartCharMetrics 6

C 32 ; WX 600 ; N space ; B 0 0 0 0 ;

C 33 ; WX 600 ; N exclam ; B 236 -15 364 572 ;

C 34 ; WX 600 ; N quotedbl ; B 187 328 413 562 ;

C 35 ; WX 600 ; N numbersign ; B 93 -32 507 639 ;

C 36 ; WX 600 ; N dollar ; B 105 -126 496 662 ;

C 37 ; WX 600 ; N percent ; B 81 -15 518 622 ;

EndCharMetrics`

	input := core.NewMemoryInputBytes([]byte(courierAfmSnippet))

	metrics, err := adobe_font_metrics.Parse(input, false)
	if err != nil {
		t.Fatalf("unexpected error parsing AFM: %v", err)
	}

	if metrics.FontName == "" {
		t.Error("expected non-empty FontName")
	}
}

func TestCanParseHelveticaAfmFile(t *testing.T) {
	helveticaPath := filepath.Join("testdata", "adobe_font_metrics", "Helvetica.afm")

	data, err := os.ReadFile(helveticaPath)
	if err != nil {
		t.Fatalf("failed to read Helvetica.afm: %v", err)
	}

	input := core.NewMemoryInputBytes(data)

	metrics, err := adobe_font_metrics.Parse(input, false)
	if err != nil {
		t.Fatalf("unexpected error parsing Helvetica AFM: %v", err)
	}

	if metrics.FontName == "" {
		t.Error("expected non-empty FontName")
	}
}
