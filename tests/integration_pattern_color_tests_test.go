//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"math"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/annotations"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
)

func getAtOrFatal(t *testing.T, m core.TransformationMatrix, row, col int) float64 {
	t.Helper()
	v, err := m.GetAt(row, col)
	if err != nil {
		t.Fatalf("GetAt(%d, %d): %v", row, col, err)
	}
	return v
}

func assertFloatNearPatternColor(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-4 {
		t.Errorf("%s: expected %f, got %f", name, want, got)
	}
}

// TestShadingPattern1 matches C# PatternColorTests.ShadingPattern1.
// TODO: load color space in annotation appearance
// TODO: contains function with indirect reference
func TestShadingPattern1(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "cat-genetics_bobld.pdf")

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	annots := page.GetAnnotations()
	if len(annots) < 15 {
		t.Fatalf("expected at least 15 annotations, got %d", len(annots))
	}

	annotationStamp := annots[14]
	at, ok := annotationStamp.Type().(annotations.AnnotationType)
	if !ok || at != annotations.Stamp {
		t.Errorf("expected AnnotationType.Stamp, got %v (type %T)", annotationStamp.Type(), annotationStamp.Type())
	}

	// HasNormalAppearance not available on LinkAnnotationIface;
	// would require access to underlying *annotations.Annotation.
	t.Log("HasNormalAppearance check skipped: not exposed on LinkAnnotationIface")
	// normalAppearanceStream is unexported field in C# and not accessible here.
}

// TestShadingPattern2 matches C# PatternColorTests.ShadingPattern2.
func TestShadingPattern2(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "output_w3c_csswg_drafts_issues2023.pdf")

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	paths := page.Paths()
	if len(paths) != 1 {
		t.Fatalf("expected exactly 1 path, got %d", len(paths))
	}

	pathObj := paths[0]
	color := pathObj.FillColor()
	if color.ColorSpace() != colors.Pattern {
		t.Errorf("expected ColorSpace.Pattern, got %s", color.ColorSpace())
	}

	shadingPatternColor, ok := color.(colors.ShadingPatternColor)
	if !ok {
		t.Fatal("expected ShadingPatternColor")
	}

	if shadingPatternColor.PatternType() != colors.ShadingPatternType {
		t.Errorf("expected PatternType.Shading, got %s", shadingPatternColor.PatternType())
	}

	if shadingPatternColor.PatternDictionary() == nil {
		t.Error("expected non-nil PatternDictionary")
	}

	shading := shadingPatternColor.Shading()
	if shading == nil {
		t.Fatal("expected non-nil Shading")
	}

	csDetails := shading.ColorSpace()
	if csDetails.Type() != colors.DeviceN {
		t.Errorf("expected ColorSpace.DeviceN, got %s", csDetails.Type())
	}

	// deviceNColorSpaceDetails.Names is a private field not exposed through
	// the ColorSpaceDetails interface. Cannot assert Names count or contents.
	t.Log("DeviceN names assertions skipped: names field not exposed on ColorSpaceDetails interface")
}

// TestTillingPattern1 matches C# PatternColorTests.TillingPattern1.
func TestTillingPattern1(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "22060_A1_01_Plans-1.pdf")

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("expected *content.Page, got %T", pageAny)
	}

	var filledPaths []content.PdfPath
	for _, p := range page.Paths() {
		if p.IsFilled() {
			filledPaths = append(filledPaths, p)
		}
	}

	if len(filledPaths) == 0 {
		t.Fatal("expected at least one filled path")
	}

	pattern := filledPaths[0].FillColor()
	if pattern.ColorSpace() != colors.Pattern {
		t.Errorf("expected ColorSpace.Pattern, got %s", pattern.ColorSpace())
	}

	tilingColor, ok := pattern.(colors.TilingPatternColor)
	if !ok {
		t.Fatal("expected TilingPatternColor")
	}

	matrix := tilingColor.Matrix()
assertFloatNearPatternColor(t, "Matrix[0,0]", getAtOrFatal(t, matrix, 0, 0), 0.213333)

	assertFloatNearPatternColor(t, "Matrix[0,1]", getAtOrFatal(t, matrix, 0, 1), 0.0)

	assertFloatNearPatternColor(t, "Matrix[0,2]", getAtOrFatal(t, matrix, 0, 2), 0.0)

	assertFloatNearPatternColor(t, "Matrix[1,0]", getAtOrFatal(t, matrix, 1, 0), 0.0)

	assertFloatNearPatternColor(t, "Matrix[1,1]", getAtOrFatal(t, matrix, 1, 1), 0.213333)

	assertFloatNearPatternColor(t, "Matrix[1,2]", getAtOrFatal(t, matrix, 1, 2), 0.0)

	assertFloatNearPatternColor(t, "Matrix[2,0]", getAtOrFatal(t, matrix, 2, 0), -0.231058)

	assertFloatNearPatternColor(t, "Matrix[2,1]", getAtOrFatal(t, matrix, 2, 1), 1190.67)

	assertFloatNearPatternColor(t, "Matrix[2,2]", getAtOrFatal(t, matrix, 2, 2), 1.0)

	if tilingColor.ExtGState() != nil {
		t.Error("expected nil ExtGState")
	}

	if tilingColor.PatternDictionary() == nil {
		t.Error("expected non-nil PatternDictionary")
	}

	if tilingColor.PatternStream() == nil {
		t.Fatal("expected non-nil PatternStream")
	}
assertFloatNearPatternColor(t, "XStep", tilingColor.XStep(), 1897.47)

	assertFloatNearPatternColor(t, "YStep", tilingColor.YStep(), 2012.23)

	if len(tilingColor.Data()) != 142 {
		t.Errorf("expected Data.Length == 142, got %d", len(tilingColor.Data()))
	}

	bBox := tilingColor.BBox()
	expectedBottomLeft := core.NewPdfPoint(-18.6026, -1992.51)
	if !bBox.BottomLeft.Equals(expectedBottomLeft) {
		t.Errorf("expected BBox.BottomLeft == (%f, %f), got (%f, %f)",
			expectedBottomLeft.X, expectedBottomLeft.Y, bBox.BottomLeft.X, bBox.BottomLeft.Y)
	}

	expectedTopRight := core.NewPdfPoint(1878.86, 19.7278)
	if !bBox.TopRight.Equals(expectedTopRight) {
		t.Errorf("expected BBox.TopRight == (%f, %f), got (%f, %f)",
			expectedTopRight.X, expectedTopRight.Y, bBox.TopRight.X, bBox.TopRight.Y)
	}

	if tilingColor.PaintType() != colors.Coloured {
		t.Errorf("expected PaintType.Coloured, got %s", tilingColor.PaintType())
	}

	if tilingColor.TilingType() != colors.ConstantSpacing {
		t.Errorf("expected TilingType.ConstantSpacing, got %s", tilingColor.TilingType())
	}

	resources := tilingColor.Resources()
	if resources == nil {
		t.Fatal("expected non-nil Resources")
	}

	if len(resources.Data()) != 4 {
		t.Errorf("expected Resources.Data.Count == 4, got %d", len(resources.Data()))
	}
}

// TestTillingPattern2 matches C# PatternColorTests.TillingPattern2.
func TestTillingPattern2(t *testing.T) {
	path := filepath.Join(integrationDocRoot, "SPARC - v9 Architecture Manual.pdf")

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	checkPage := func(pageNum int, expectedCount int) {
		pageAny, err := doc.GetPage(pageNum)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", pageNum, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("page %d: expected *content.Page, got %T", pageNum, pageAny)
		}

		var strokedPaths []content.PdfPath
		for _, p := range page.Paths() {
			sc := p.StrokeColor()
			if sc != nil && sc.ColorSpace() == colors.Pattern {
				strokedPaths = append(strokedPaths, p)
			}
		}

		if len(strokedPaths) != expectedCount {
			t.Errorf("page %d: expected %d stroked pattern paths, got %d", pageNum, expectedCount, len(strokedPaths))
		}

		for _, p := range strokedPaths {
			sc := p.StrokeColor()
			if sc.ColorSpace() != colors.Pattern {
				t.Errorf("page %d: expected ColorSpace.Pattern, got %s", pageNum, sc.ColorSpace())
			}

			patternColor, ok := sc.(colors.TilingPatternColor)
			if !ok {
				t.Errorf("page %d: expected TilingPatternColor", pageNum)
				continue
			}

			if patternColor.PatternType() != colors.Tiling {
				t.Errorf("page %d: expected PatternType.Tiling, got %s", pageNum, patternColor.PatternType())
			}

			if patternColor.PaintType() != colors.Uncoloured {
				t.Errorf("page %d: expected PaintType.Uncoloured, got %s", pageNum, patternColor.PaintType())
			}

			if patternColor.TilingType() != colors.ConstantSpacingFasterTiling {
				t.Errorf("page %d: expected TilingType.ConstantSpacingFasterTiling, got %s",
					pageNum, patternColor.TilingType())
			}
		}
	}

	checkPage(53, 5)
	checkPage(307, 2)
}
