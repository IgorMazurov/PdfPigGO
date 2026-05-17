package geometry

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry/clipperlibrary"
	"github.com/uglytoad/pdfpig/go/logging"
)

// Factor is the scaling factor used to convert PDF coordinates to integer coordinates for Clipper.
const Factor = 10_000.0

const linesInCurve = 10

// Clip applies a clipping path to another path and returns the clipped result.
func Clip(clipping *PdfPath, subject *PdfPath, log logging.Log) (*PdfPath, error) {
	if clipping == nil {
		return nil, fmt.Errorf("Clip: the clipping path cannot be nil")
	}

	if !clipping.IsClipping() {
		return nil, fmt.Errorf("Clip: the clipping path does not have the IsClipping flag set to true")
	}

	if subject == nil {
		return nil, fmt.Errorf("Clip: the subject path cannot be nil")
	}

	if subject.Len() == 0 {
		return subject, nil
	}

	clipper := clipperlibrary.NewClipper(0)

	for _, subPathClipping := range clipping.Subpaths() {
		if len(subPathClipping.Commands()) == 0 {
			continue
		}

		if !subPathClipping.IsClosed() {
			subPathClipping.CloseSubpath()
		}

		polygon := SubpathToClipperPolygon(subPathClipping)
		if ok := clipper.AddPath(polygon, clipperlibrary.ClipperClip, true); !ok {
			log.Error("ClippingExtensions.Clip(): failed to add clipping subpath.")
		}
	}

	subjectClose := subject.IsFilled() || subject.IsClipping()

	for _, subPathSubject := range subject.Subpaths() {
		if len(subPathSubject.Commands()) == 0 {
			continue
		}

		lineCount := countLines(subPathSubject)
		curveCount := countCurves(subPathSubject)

		if subjectClose && !subPathSubject.IsClosed() && lineCount < 2 && curveCount == 0 {
			subjectClose = false
		}

		if subjectClose && !subPathSubject.IsClosed() {
			subPathSubject.CloseSubpath()
		}

		polygon := SubpathToClipperPolygon(subPathSubject)
		if ok := clipper.AddPath(polygon, clipperlibrary.ClipperSubject, subjectClose); !ok {
			log.Error("ClippingExtensions.Clip(): failed to add subject subpath for clipping.")
		}
	}

	clippingFillType := toClipperFillType(clipping.FillingRule())
	subjectFillType := toClipperFillType(subject.FillingRule())

	if !subjectClose {
		return clipOpenPath(clipper, subject, subjectFillType, clippingFillType)
	}

	return clipClosedPath(clipper, subject, subjectFillType, clippingFillType)
}

func countLines(subpath *core.PdfSubpath) int {
	count := 0
	for _, cmd := range subpath.Commands() {
		if _, ok := cmd.(*core.Line); ok {
			count++
		}
	}
	return count
}

func countCurves(subpath *core.PdfSubpath) int {
	count := 0
	for _, cmd := range subpath.Commands() {
		if _, ok := cmd.(*core.QuadraticBezierCurve); ok {
			count++
		} else if _, ok := cmd.(*core.CubicBezierCurve); ok {
			count++
		}
	}
	return count
}

func toClipperFillType(rule core.FillingRule) clipperlibrary.ClipperPolyFillType {
	if rule == core.FillingRuleNonZeroWinding {
		return clipperlibrary.ClipperNonZero
	}
	return clipperlibrary.ClipperEvenOdd
}

func clipOpenPath(clipper *clipperlibrary.Clipper, subject *PdfPath, subjFT, clipFT clipperlibrary.ClipperPolyFillType) (*PdfPath, error) {
	clippedPath := subject.CloneEmpty()

	polytree := clipperlibrary.NewClipperPolyTree()
	if !clipper.ExecutePolyTree(clipperlibrary.ClipperIntersection, polytree, subjFT, clipFT) {
		return nil, nil
	}

	for i := 0; i < polytree.ChildCount(); i++ {
		solution := polytree.Children[i]
		contour := solution.Contour()
		if len(contour) == 0 {
			continue
		}

		clippedSubpath := core.NewPdfSubpath()
		clippedSubpath.MoveTo(float64(contour[0].X)/Factor, float64(contour[0].Y)/Factor)

		for i := 1; i < len(contour); i++ {
			clippedSubpath.LineTo(float64(contour[i].X)/Factor, float64(contour[i].Y)/Factor)
		}
		clippedPath.Add(clippedSubpath)
	}

	if clippedPath.Len() > 0 {
		return clippedPath, nil
	}

	return nil, nil
}

func clipClosedPath(clipper *clipperlibrary.Clipper, subject *PdfPath, subjFT, clipFT clipperlibrary.ClipperPolyFillType) (*PdfPath, error) {
	clippedPath := subject.CloneEmpty()

	solutions := clipper.Execute(clipperlibrary.ClipperIntersection, nil, subjFT, clipFT)
	if solutions == nil {
		return nil, nil
	}

	for _, solution := range solutions {
		if len(solution) == 0 {
			continue
		}

		clippedSubpath := core.NewPdfSubpath()
		clippedSubpath.MoveTo(float64(solution[0].X)/Factor, float64(solution[0].Y)/Factor)

		for i := 1; i < len(solution); i++ {
			clippedSubpath.LineTo(float64(solution[i].X)/Factor, float64(solution[i].Y)/Factor)
		}
		clippedSubpath.CloseSubpath()
		clippedPath.Add(clippedSubpath)
	}

	if clippedPath.Len() > 0 {
		return clippedPath, nil
	}

	return nil, nil
}

// SubpathToClipperPolygon converts a PdfSubpath to Clipper integer points.
func SubpathToClipperPolygon(pdfPath *core.PdfSubpath) []clipperlibrary.ClipperIntPoint {
	cmds := pdfPath.Commands()
	if len(cmds) == 0 {
		return nil
	}

	var result []clipperlibrary.ClipperIntPoint

	move, ok := cmds[0].(*core.Move)
	if !ok {
		return nil
	}

	movePoint := PointToClipperInt(move.Location)
	result = append(result, movePoint)

	if len(cmds) == 1 {
		return result
	}

	for i := 1; i < len(cmds); i++ {
		cmd := cmds[i]

		if _, ok := cmd.(*core.Move); ok {
			continue
		}

		switch c := cmd.(type) {
		case *core.Line:
			result = append(result, PointToClipperInt(c.From))
			result = append(result, PointToClipperInt(c.To))

		case *core.QuadraticBezierCurve:
			for _, lineB := range c.ToLines(linesInCurve) {
				result = append(result, PointToClipperInt(lineB.From))
				result = append(result, PointToClipperInt(lineB.To))
			}

		case *core.CubicBezierCurve:
			for _, lineB := range c.ToLines(linesInCurve) {
				result = append(result, PointToClipperInt(lineB.From))
				result = append(result, PointToClipperInt(lineB.To))
			}

		case *core.Close:
			result = append(result, movePoint)
		}
	}

	return result
}

// RectangleToClipperPolygon converts a PdfRectangle to Clipper integer points.
func RectangleToClipperPolygon(rectangle core.PdfRectangle) []clipperlibrary.ClipperIntPoint {
	return []clipperlibrary.ClipperIntPoint{
		PointToClipperInt(rectangle.BottomLeft),
		PointToClipperInt(rectangle.TopLeft),
		PointToClipperInt(rectangle.TopRight),
		PointToClipperInt(rectangle.BottomRight),
	}
}

// PointToClipperInt converts a PdfPoint to Clipper integer coordinates.
func PointToClipperInt(point core.PdfPoint) clipperlibrary.ClipperIntPoint {
	return clipperlibrary.NewClipperIntPoint(
		int64(point.X*Factor),
		int64(point.Y*Factor),
	)
}

// PdfLineToClipperInt converts a PdfLine to Clipper integer points.
func PdfLineToClipperInt(line core.PdfLine) []clipperlibrary.ClipperIntPoint {
	return []clipperlibrary.ClipperIntPoint{
		PointToClipperInt(line.Point1),
		PointToClipperInt(line.Point2),
	}
}
