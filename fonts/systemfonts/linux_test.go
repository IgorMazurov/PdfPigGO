package systemfonts_test

import (
	"math"
	"testing"

	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/fonts/systemfonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/content"
)

// expectedLetterData holds the expected bounding box values for a letter.
type expectedLetterData struct {
	topLeft  core.PdfPoint
	width    float64
	height   float64
	rotation float64
}

var linuxBBoxTestData = []struct {
	name     string
	expected []expectedLetterData
}{
	{
		name: "90 180 270 rotated.pdf",
		expected: []expectedLetterData{
			{
				topLeft:  core.NewPdfPoint(53.88, 759.48),
				width:    2.495859375,
				height:   0,
				rotation: 0,
			},
			{
				topLeft:  core.NewPdfPoint(514.925312502883, 744.099765720344),
				width:    6.83203125,
				height:   7.94531249999983,
				rotation: -90,
			},
			{
				topLeft:  core.NewPdfPoint(512.505390717836, 736.603703191305),
				width:    5.1796875,
				height:   5.68945312499983,
				rotation: -90,
			},
			{
				topLeft:  core.NewPdfPoint(512.505390785898, 730.931828191305),
				width:    3.99609375,
				height:   5.52539062499994,
				rotation: -90,
			},
		},
	},
}

func TestLinuxGetCorrectBBox(t *testing.T) {
	font := systemfonts.Instance.GetTrueTypeFont("TimesNewRomanPSMT")
	if font == nil {
		t.Skip("Skipped because the font TimesNewRomanPSMT could not be found in the execution environment.")
	}

	for _, tc := range linuxBBoxTestData {
		t.Run(tc.name, func(t *testing.T) {
			docPath := dla.GetDocumentPath(tc.name, true)

			doc, err := pdfpig.OpenFile(docPath, nil)
			if err != nil {
				t.Fatalf("OpenFile(%q): %v", docPath, err)
			}
			defer doc.Close()

			pageAny, err := doc.GetPage(1)
			if err != nil {
				t.Fatalf("GetPage(1): %v", err)
			}

			page, ok := pageAny.(*content.Page)
			if !ok {
				t.Fatal("expected *content.Page")
			}

			letters := page.Letters()
			if len(letters) < len(tc.expected) {
				t.Fatalf("expected at least %d letters, got %d", len(tc.expected), len(letters))
			}

			for i, exp := range tc.expected {
				current := letters[i]
				bb := current.BoundingBox

				assertApprox(t, exp.topLeft.X, bb.TopLeft.X, 1e-6, "Letter[%d] BoundingBox.TopLeft.X", i)
				assertApprox(t, exp.topLeft.Y, bb.TopLeft.Y, 1e-6, "Letter[%d] BoundingBox.TopLeft.Y", i)
				assertApprox(t, exp.width, bb.Width, 1e-6, "Letter[%d] BoundingBox.Width", i)
				assertApprox(t, exp.height, bb.Height, 1e-6, "Letter[%d] BoundingBox.Height", i)
				assertApprox(t, exp.rotation, bb.Rotation(), 1e-3, "Letter[%d] BoundingBox.Rotation", i)

				glyphRect := current.GlyphRectangleLoose
				if !geometry.RectangleIntersectsWithRect(bb, glyphRect) {
					t.Errorf("Letter[%d]: BoundingBox does not intersect GlyphRectangleLoose", i)
				}

				assertApprox(t, bb.Rotation(), glyphRect.Rotation(), 1e-3, "Letter[%d] Rotation match between BoundingBox and GlyphRectangleLoose", i)
			}
		})
	}
}

func assertApprox(t *testing.T, expected, actual, tolerance float64, msg string, args ...any) {
	t.Helper()
	if math.Abs(expected-actual) >= tolerance {
		t.Errorf(msg+": expected %g, got %g (diff=%g, tolerance=%g)", append([]any{expected, actual, math.Abs(expected-actual), tolerance}, args...)...)
	}
}
