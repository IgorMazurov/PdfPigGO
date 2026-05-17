package document_layout_analysis

import (
	"errors"
	"math"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
)

// TextLine represents a line of text composed of one or more words.
type TextLine struct {
	Separator         string
	Text              string
	TextOrientation   content.TextOrientation
	boundingBox       core.PdfRectangle
	Words             []*content.Word
}

// NewTextLine creates a new TextLine from an ordered list of words.
func NewTextLine(words []*content.Word, separator string) (*TextLine, error) {
	if len(words) == 0 {
		return nil, errors.New("empty words provided")
	}

	if separator == "" {
		separator = " "
	}

	tempOrientation := words[0].TextOrientation
	if tempOrientation != content.OtherTextOrientation {
		for _, w := range words {
			if w.TextOrientation != tempOrientation {
				tempOrientation = content.OtherTextOrientation
				break
			}
		}
	}

	var bb core.PdfRectangle
	var text string

	if len(words) == 1 {
		bb = words[0].BoundingBox()
		text = words[0].Text
	} else {
		switch tempOrientation {
		case content.HorizontalTextOrientation:
			bb = getLineBoundingBoxH(words)
		case content.Rotate180TextOrientation:
			bb = getLineBoundingBox180(words)
		case content.Rotate90TextOrientation:
			bb = getLineBoundingBox90(words)
		case content.Rotate270TextOrientation:
			bb = getLineBoundingBox270(words)
		default:
			bb = getLineBoundingBoxOther(words)
		}

		var parts []string
		for _, w := range words {
			if w.Text != "" {
				parts = append(parts, w.Text)
			}
		}
		text = strings.Join(parts, separator)
	}

	return &TextLine{
		Separator:       separator,
		Text:            text,
		TextOrientation: tempOrientation,
		boundingBox:     bb,
		Words:           words,
	}, nil
}

// BoundingBox returns the rectangle completely containing the line.
func (tl *TextLine) BoundingBox() core.PdfRectangle {
	return tl.boundingBox
}

func getLineBoundingBoxH(words []*content.Word) core.PdfRectangle {
	blX := math.Inf(1)
	trX := math.Inf(-1)
	blY := math.Inf(1)
	trY := math.Inf(-1)

	for _, w := range words {
		bb := w.BoundingBox()
		if bb.BottomLeft.X < blX {
			blX = bb.BottomLeft.X
		}
		if bb.BottomLeft.Y < blY {
			blY = bb.BottomLeft.Y
		}
		right := bb.BottomLeft.X + bb.Width
		if right > trX {
			trX = right
		}
		if bb.TopLeft.Y > trY {
			trY = bb.TopLeft.Y
		}
	}

	return core.NewPdfRectangleFloat(blX, blY, trX, trY)
}

func getLineBoundingBox180(words []*content.Word) core.PdfRectangle {
	blX := math.Inf(-1)
	blY := math.Inf(-1)
	trX := math.Inf(1)
	trY := math.Inf(1)

	for _, w := range words {
		bb := w.BoundingBox()
		if bb.BottomLeft.X > blX {
			blX = bb.BottomLeft.X
		}
		if bb.BottomLeft.Y > blY {
			blY = bb.BottomLeft.Y
		}
		right := bb.BottomLeft.X - bb.Width
		if right < trX {
			trX = right
		}
		if bb.TopRight.Y < trY {
			trY = bb.TopRight.Y
		}
	}

	return core.NewPdfRectangleFloat(blX, blY, trX, trY)
}

func getLineBoundingBox90(words []*content.Word) core.PdfRectangle {
	b := math.Inf(1)
	r := math.Inf(1)
	t := math.Inf(-1)
	l := math.Inf(-1)

	for _, w := range words {
		bb := w.BoundingBox()
		if bb.BottomLeft.X < b {
			b = bb.BottomLeft.X
		}
		if bb.BottomRight.Y < r {
			r = bb.BottomRight.Y
		}
		right := bb.BottomLeft.X + bb.Height
		if right > t {
			t = right
		}
		if bb.BottomLeft.Y > l {
			l = bb.BottomLeft.Y
		}
	}

	return core.NewPdfRectangleFromCorners(
		core.NewPdfPoint(t, l), core.NewPdfPoint(t, r),
		core.NewPdfPoint(b, l), core.NewPdfPoint(b, r))
}

func getLineBoundingBox270(words []*content.Word) core.PdfRectangle {
	t := math.Inf(1)
	b := math.Inf(-1)
	l := math.Inf(1)
	r := math.Inf(-1)

	for _, w := range words {
		bb := w.BoundingBox()
		if bb.BottomLeft.X > b {
			b = bb.BottomLeft.X
		}
		if bb.BottomLeft.Y < l {
			l = bb.BottomLeft.Y
		}
		right := bb.BottomLeft.X - bb.Height
		if right < t {
			t = right
		}
		if bb.BottomRight.Y > r {
			r = bb.BottomRight.Y
		}
	}

	return core.NewPdfRectangleFromCorners(
		core.NewPdfPoint(t, l), core.NewPdfPoint(t, r),
		core.NewPdfPoint(b, l), core.NewPdfPoint(b, r))
}

func getLineBoundingBoxOther(words []*content.Word) core.PdfRectangle {
	baseLinePoints := make([]core.PdfPoint, 0, len(words)*2)
	for _, w := range words {
		bb := w.BoundingBox()
		baseLinePoints = append(baseLinePoints, bb.BottomLeft, bb.BottomRight)
	}

	x0 := averagePointX(baseLinePoints)
	y0 := averagePointY(baseLinePoints)

	var sumProduct, sumDiffSquaredX float64
	for _, p := range baseLinePoints {
		xDiff := p.X - x0
		yDiff := p.Y - y0
		sumProduct += xDiff * yDiff
		sumDiffSquaredX += xDiff * xDiff
	}

	cos := 0.0
	sin := 1.0
	if sumDiffSquaredX > 1e-3 {
		angleRad := math.Atan(sumProduct / sumDiffSquaredX)
		cos = math.Cos(angleRad)
		sin = math.Sin(angleRad)
	}

	inverseRotation := core.NewTransformationMatrix(
		cos, -sin, 0,
		sin, cos, 0,
		0, 0, 1)

	seen := make(map[core.PdfPoint]bool)
	var transformedPoints []core.PdfPoint
	for _, w := range words {
		bb := w.BoundingBox()
		candidates := []core.PdfPoint{bb.BottomLeft, bb.BottomRight, bb.TopLeft, bb.TopRight}
		for _, p := range candidates {
			if !seen[p] {
				seen[p] = true
				transformedPoints = append(transformedPoints, inverseRotation.TransformPoint(p))
			}
		}
	}

	minX, minY, maxX, maxY := transformedPoints[0].X, transformedPoints[0].Y, transformedPoints[0].X, transformedPoints[0].Y
	for _, p := range transformedPoints[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	aabb := core.NewPdfRectangleFloat(minX, minY, maxX, maxY)

	rotateBack := core.NewTransformationMatrix(
		cos, sin, 0,
		-sin, cos, 0,
		0, 0, 1)

	obb := rotateBack.TransformRect(aabb)
	obb1 := core.NewPdfRectangleFromCorners(obb.BottomLeft, obb.TopLeft, obb.BottomRight, obb.TopRight)
	obb2 := core.NewPdfRectangleFromCorners(obb.BottomRight, obb.BottomLeft, obb.TopRight, obb.TopLeft)
	obb3 := core.NewPdfRectangleFromCorners(obb.TopRight, obb.BottomRight, obb.TopLeft, obb.BottomLeft)

	firstWord := words[0]
	lastWord := words[len(words)-1]
	baseLineAngle := Angle(firstWord.BoundingBox().BottomLeft, lastWord.BoundingBox().BottomRight)

	deltaAngle := math.Abs(BoundAngle180(obb.Rotation()-baseLineAngle))
	bestOBB := obb

	if d := math.Abs(BoundAngle180(obb1.Rotation()-baseLineAngle)); d < deltaAngle {
		deltaAngle = d
		bestOBB = obb1
	}
	if d := math.Abs(BoundAngle180(obb2.Rotation()-baseLineAngle)); d < deltaAngle {
		deltaAngle = d
		bestOBB = obb2
	}
	if d := math.Abs(BoundAngle180(obb3.Rotation()-baseLineAngle)); d < deltaAngle {
		bestOBB = obb3
	}

	return bestOBB
}

func averagePointX(points []core.PdfPoint) float64 {
	var sum float64
	for _, p := range points {
		sum += p.X
	}
	return sum / float64(len(points))
}

func averagePointY(points []core.PdfPoint) float64 {
	var sum float64
	for _, p := range points {
		sum += p.Y
	}
	return sum / float64(len(points))
}

// String returns the text of the line.
func (tl *TextLine) String() string {
	return tl.Text
}
