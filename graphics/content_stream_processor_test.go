package graphics

import (
	"math"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/logging"
)

func TestInitialMatrixHandlesDefaultCase(t *testing.T) {
	mediaBox := core.NewPdfRectangleFromInt(0, 0, 595, 842)
	cropBox := core.NewPdfRectangleFromInt(0, 0, 595, 842)

	glyph := core.NewPdfRectangleFloat(cropBox.Left(), cropBox.Top()-20, cropBox.Left()+10, cropBox.Top())

	rotation, _ := content.NewPageRotationDegrees(0)
	initialMatrix := GetInitialMatrix(geometry.Default, content.NewMediaBox(mediaBox), content.NewCropBox(cropBox), rotation, logging.NoopLog)
	inverseMatrix := initialMatrix.Inverse()

	transformedGlyph := initialMatrix.TransformRect(glyph)
	inverseTransformedGlyph := inverseMatrix.TransformRect(transformedGlyph)

	assertPdfRectangleEqual(t, glyph, transformedGlyph)
	assertPdfRectangleEqual(t, glyph, inverseTransformedGlyph)
}

func TestInitialMatrixHandlesCropBoxOutsideMediaBox(t *testing.T) {
	mediaBox := core.NewPdfRectangleFromInt(0, 0, 595, 842)
	cropBox := core.NewPdfRectangleFromInt(400, 400, 1000, 1000)

	pointInsideViewBox := core.NewPdfPoint(500, 500)
	pointBelowViewBox := core.NewPdfPoint(500, 100)
	pointLeftOfViewBox := core.NewPdfPoint(200, 500)
	pointAboveViewBox := core.NewPdfPoint(500, 1000)
	pointRightOfViewBox := core.NewPdfPoint(1000, 500)

	rotation, _ := content.NewPageRotationDegrees(0)
	initialMatrix := GetInitialMatrix(geometry.Default, content.NewMediaBox(mediaBox), content.NewCropBox(cropBox), rotation, logging.NoopLog)
	inverseMatrix := initialMatrix.Inverse()

	pt := initialMatrix.TransformPoint(pointInsideViewBox)
	p0 := inverseMatrix.TransformPoint(pt)
	assertPdfPointEqual(t, pointInsideViewBox, p0)
	if !(pt.X > 0 && pt.X < 195 && pt.Y > 0 && pt.Y < 442) {
		t.Errorf("point inside view box transformed to (%g, %g), expected X in (0, 195) and Y in (0, 442)", pt.X, pt.Y)
	}

	pt = initialMatrix.TransformPoint(pointBelowViewBox)
	p0 = inverseMatrix.TransformPoint(pt)
	assertPdfPointEqual(t, pointBelowViewBox, p0)
	if !(pt.X > 0 && pt.X < 195 && pt.Y < 0) {
		t.Errorf("point below view box transformed to (%g, %g), expected X in (0, 195) and Y < 0", pt.X, pt.Y)
	}

	pt = initialMatrix.TransformPoint(pointLeftOfViewBox)
	p0 = inverseMatrix.TransformPoint(pt)
	assertPdfPointEqual(t, pointLeftOfViewBox, p0)
	if !(pt.X < 0 && pt.Y > 0 && pt.Y < 442) {
		t.Errorf("point left of view box transformed to (%g, %g), expected X < 0 and Y in (0, 442)", pt.X, pt.Y)
	}

	rotation180, _ := content.NewPageRotationDegrees(180)
	initialMatrix180 := GetInitialMatrix(geometry.Default, content.NewMediaBox(mediaBox), content.NewCropBox(cropBox), rotation180, logging.NoopLog)
	inverseMatrix180 := initialMatrix180.Inverse()

	pt = initialMatrix180.TransformPoint(pointInsideViewBox)
	p0 = inverseMatrix180.TransformPoint(pt)
	assertPdfPointEqual(t, pointInsideViewBox, p0)
	if !(pt.X > 0 && pt.X < 195 && pt.Y > 0 && pt.Y < 442) {
		t.Errorf("point inside view box (180) transformed to (%g, %g), expected X in (0, 195) and Y in (0, 442)", pt.X, pt.Y)
	}

	pt = initialMatrix180.TransformPoint(pointAboveViewBox)
	p0 = inverseMatrix180.TransformPoint(pt)
	assertPdfPointEqual(t, pointAboveViewBox, p0)
	if !(pt.X > 0 && pt.X < 195 && pt.Y < 0) {
		t.Errorf("point above view box (180) transformed to (%g, %g), expected X in (0, 195) and Y < 0", pt.X, pt.Y)
	}

	pt = initialMatrix180.TransformPoint(pointRightOfViewBox)
	p0 = inverseMatrix180.TransformPoint(pt)
	assertPdfPointEqual(t, pointRightOfViewBox, p0)
	if !(pt.X < 0 && pt.Y > 0 && pt.Y < 442) {
		t.Errorf("point right of view box (180) transformed to (%g, %g), expected X < 0 and Y in (0, 442)", pt.X, pt.Y)
	}
}

func TestInitialMatrixHandlesCropBoxAndRotation(t *testing.T) {
	mediaBox := core.NewPdfRectangleFromInt(0, 0, 595, 842)
	cropBox := core.NewPdfRectangleFromInt(100, 200, 400, 600)
	glyph := core.NewPdfRectangleFloat(cropBox.Left(), cropBox.Top()-20, cropBox.Left()+10, cropBox.Top())

	testRotation(t, mediaBox, cropBox, glyph, 0, func(tr core.PdfRectangle) {
		assertFloatEqual(t, 0, tr.BottomLeft.X, "BL.X for 0deg")
		assertFloatEqual(t, cropBox.Height-glyph.Height, tr.BottomLeft.Y, "BL.Y for 0deg")
		assertFloatEqual(t, glyph.Width, tr.TopRight.X, "TR.X for 0deg")
		assertFloatEqual(t, cropBox.Height, tr.TopRight.Y, "TR.Y for 0deg")
	})

	testRotation(t, mediaBox, cropBox, glyph, 90, func(tr core.PdfRectangle) {
		assertFloatEqual(t, cropBox.Height-glyph.Height, tr.BottomLeft.X, "BL.X for 90deg")
		assertFloatEqual(t, cropBox.Width, tr.BottomLeft.Y, "BL.Y for 90deg")
		assertFloatEqual(t, cropBox.Height, tr.TopRight.X, "TR.X for 90deg")
		assertFloatEqual(t, cropBox.Width-glyph.Width, tr.TopRight.Y, "TR.Y for 90deg")
	})

	testRotation(t, mediaBox, cropBox, glyph, 180, func(tr core.PdfRectangle) {
		assertFloatEqual(t, cropBox.Width, tr.BottomLeft.X, "BL.X for 180deg")
		assertFloatEqual(t, glyph.Height, tr.BottomLeft.Y, "BL.Y for 180deg")
		assertFloatEqual(t, cropBox.Width-glyph.Width, tr.TopRight.X, "TR.X for 180deg")
		assertFloatEqual(t, 0, tr.TopRight.Y, "TR.Y for 180deg")
	})

	testRotation(t, mediaBox, cropBox, glyph, 270, func(tr core.PdfRectangle) {
		assertFloatEqual(t, glyph.Height, tr.BottomLeft.X, "BL.X for 270deg")
		assertFloatEqual(t, 0, tr.BottomLeft.Y, "BL.Y for 270deg")
		assertFloatEqual(t, 0, tr.TopRight.X, "TR.X for 270deg")
		assertFloatEqual(t, glyph.Width, tr.TopRight.Y, "TR.Y for 270deg")
	})
}

func testRotation(t *testing.T, mediaBox, cropBox, glyph core.PdfRectangle, degrees int, assertTransformed func(core.PdfRectangle)) {
	t.Helper()
	rotation, _ := content.NewPageRotationDegrees(degrees)
	initialMatrix := GetInitialMatrix(geometry.Default, content.NewMediaBox(mediaBox), content.NewCropBox(cropBox), rotation, logging.NoopLog)
	inverseMatrix := initialMatrix.Inverse()

	transformedGlyph := initialMatrix.TransformRect(glyph)
	inverseTransformedGlyph := inverseMatrix.TransformRect(transformedGlyph)

	assertPdfRectangleEqual(t, glyph, inverseTransformedGlyph)
	assertTransformed(transformedGlyph)
}

func assertFloatEqual(t *testing.T, expected, actual float64, msg string) {
	t.Helper()
	if math.Abs(expected-actual) > 1e-9 {
		t.Errorf("%s: expected %g, got %g", msg, expected, actual)
	}
}

func assertPdfPointEqual(t *testing.T, p1, p2 core.PdfPoint) {
	t.Helper()
	if math.Abs(p1.X-p2.X) > 1e-9 {
		t.Errorf("point X: expected %g, got %g", p1.X, p2.X)
	}
	if math.Abs(p1.Y-p2.Y) > 1e-9 {
		t.Errorf("point Y: expected %g, got %g", p1.Y, p2.Y)
	}
}

func assertPdfRectangleEqual(t *testing.T, r1, r2 core.PdfRectangle) {
	t.Helper()
	assertPdfPointEqual(t, r1.BottomLeft, r2.BottomLeft)
	assertPdfPointEqual(t, r1.TopRight, r2.TopRight)
}
