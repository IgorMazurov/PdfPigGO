package content

import (
	"errors"
	"math"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
)

// Word represents a word extracted from letters on a page.
type Word struct {
	Text            string
	TextOrientation TextOrientation
	fontName        string
	Letters         []*Letter
	boundingBox     core.PdfRectangle
}

// NewWord creates a new Word from an ordered list of letters.
func NewWord(letters []*Letter) (*Word, error) {
	if len(letters) == 0 {
		return nil, errors.New("empty letters provided")
	}

	tempTextOrientation := letters[0].TextOrientation
	if tempTextOrientation != OtherTextOrientation {
		for _, letter := range letters {
			if letter.TextOrientation != tempTextOrientation {
				tempTextOrientation = OtherTextOrientation
				break
			}
		}
	}

	var text string
	var bb core.PdfRectangle

	switch tempTextOrientation {
	case HorizontalTextOrientation:
		text, bb = getBoundingBoxH(letters)
	case Rotate180TextOrientation:
		text, bb = getBoundingBox180(letters)
	case Rotate90TextOrientation:
		text, bb = getBoundingBox90(letters)
	case Rotate270TextOrientation:
		text, bb = getBoundingBox270(letters)
	default:
		text, bb = getBoundingBoxOther(letters)
	}

	return &Word{
		Text:            text,
		boundingBox:     bb,
		fontName:        letters[0].FontName(),
		TextOrientation: tempTextOrientation,
		Letters:         letters,
	}, nil
}

// BoundingBox returns the rectangle completely containing the word.
func (w *Word) BoundingBox() core.PdfRectangle {
	return w.boundingBox
}

// FontName returns the name of the font for the word.
func (w *Word) FontName() string {
	return w.fontName
}

func getBoundingBoxH(letters []*Letter) (string, core.PdfRectangle) {
	var builder strings.Builder

	blX := math.MaxFloat64
	trX := -math.MaxFloat64
	blY := math.MaxFloat64
	trY := -math.MaxFloat64

	for _, letter := range letters {
		builder.WriteString(letter.Value)

		if letter.StartBaseLine.X < blX {
			blX = letter.StartBaseLine.X
		}

		if letter.StartBaseLine.Y < blY {
			blY = letter.StartBaseLine.Y
		}

		right := letter.StartBaseLine.X + math.Max(letter.Width, letter.BoundingBox.Width)
		if right > trX {
			trX = right
		}

		if letter.BoundingBox.TopLeft.Y > trY {
			trY = letter.BoundingBox.TopLeft.Y
		}
	}

	return builder.String(), core.NewPdfRectangleFloat(blX, blY, trX, trY)
}

func getBoundingBox180(letters []*Letter) (string, core.PdfRectangle) {
	var builder strings.Builder

	blX := -math.MaxFloat64
	blY := -math.MaxFloat64
	trX := math.MaxFloat64
	trY := math.MaxFloat64

	for _, letter := range letters {
		builder.WriteString(letter.Value)

		if letter.StartBaseLine.X > blX {
			blX = letter.StartBaseLine.X
		}

		if letter.StartBaseLine.Y > blY {
			blY = letter.StartBaseLine.Y
		}

		right := letter.StartBaseLine.X - math.Max(letter.Width, letter.BoundingBox.Width)
		if right < trX {
			trX = right
		}

		if letter.BoundingBox.TopRight.Y < trY {
			trY = letter.BoundingBox.TopRight.Y
		}
	}

	return builder.String(), core.NewPdfRectangleFloat(blX, blY, trX, trY)
}

func getBoundingBox90(letters []*Letter) (string, core.PdfRectangle) {
	var builder strings.Builder

	b := math.MaxFloat64
	r := math.MaxFloat64
	t := -math.MaxFloat64
	l := -math.MaxFloat64

	for _, letter := range letters {
		builder.WriteString(letter.Value)

		if letter.StartBaseLine.X < b {
			b = letter.StartBaseLine.X
		}

		if letter.EndBaseLine.Y < r {
			r = letter.EndBaseLine.Y
		}

		right := letter.StartBaseLine.X + letter.BoundingBox.Height
		if right > t {
			t = right
		}

		if letter.BoundingBox.BottomLeft.Y > l {
			l = letter.BoundingBox.BottomLeft.Y
		}
	}

	return builder.String(), core.NewPdfRectangleFromCorners(
		core.NewPdfPoint(t, l), core.NewPdfPoint(t, r),
		core.NewPdfPoint(b, l), core.NewPdfPoint(b, r))
}

func getBoundingBox270(letters []*Letter) (string, core.PdfRectangle) {
	var builder strings.Builder

	t := math.MaxFloat64
	b := -math.MaxFloat64
	l := math.MaxFloat64
	r := -math.MaxFloat64

	for _, letter := range letters {
		builder.WriteString(letter.Value)

		if letter.StartBaseLine.X > b {
			b = letter.StartBaseLine.X
		}

		if letter.StartBaseLine.Y < l {
			l = letter.StartBaseLine.Y
		}

		right := letter.StartBaseLine.X - letter.BoundingBox.Height
		if right < t {
			t = right
		}

		if letter.BoundingBox.BottomRight.Y > r {
			r = letter.BoundingBox.BottomRight.Y
		}
	}

	return builder.String(), core.NewPdfRectangleFromCorners(
		core.NewPdfPoint(t, l), core.NewPdfPoint(t, r),
		core.NewPdfPoint(b, l), core.NewPdfPoint(b, r))
}

func getBoundingBoxOther(letters []*Letter) (string, core.PdfRectangle) {
	var builder strings.Builder
	for _, letter := range letters {
		builder.WriteString(letter.Value)
	}

	if len(letters) == 1 {
		return builder.String(), letters[0].BoundingBox
	}

	baseLinePoints := make([]core.PdfPoint, 0, len(letters)*2)
	for _, l := range letters {
		baseLinePoints = append(baseLinePoints, l.StartBaseLine, l.EndBaseLine)
	}

	x0 := averageFloat(baseLinePoints, func(p core.PdfPoint) float64 { return p.X })
	y0 := averageFloat(baseLinePoints, func(p core.PdfPoint) float64 { return p.Y })

	var sumProduct, sumDiffSquaredX float64
	for _, point := range baseLinePoints {
		xDiff := point.X - x0
		yDiff := point.Y - y0
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

	points := make([]core.PdfPoint, 0)
	seen := make(map[core.PdfPoint]bool)
	for _, l := range letters {
		candidates := []core.PdfPoint{
			l.StartBaseLine, l.EndBaseLine,
			l.BoundingBox.TopLeft, l.BoundingBox.TopRight,
		}
		for _, p := range candidates {
			if !seen[p] {
				seen[p] = true
				points = append(points, inverseRotation.TransformPoint(p))
			}
		}
	}

	minX := points[0].X
	minY := points[0].Y
	maxX := points[0].X
	maxY := points[0].Y
	for _, p := range points[1:] {
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

	firstLetter := letters[0]
	lastLetter := letters[len(letters)-1]

	baseLineAngle := math.Atan2(
		lastLetter.EndBaseLine.Y-firstLetter.StartBaseLine.Y,
		lastLetter.EndBaseLine.X-firstLetter.StartBaseLine.X) * 180 / math.Pi

	deltaAngle := math.Abs(boundAngle180(obb.Rotation()-baseLineAngle))
	bestOBB := obb

	deltaAngle1 := math.Abs(boundAngle180(obb1.Rotation()-baseLineAngle))
	if deltaAngle1 < deltaAngle {
		deltaAngle = deltaAngle1
		bestOBB = obb1
	}

	deltaAngle2 := math.Abs(boundAngle180(obb2.Rotation()-baseLineAngle))
	if deltaAngle2 < deltaAngle {
		deltaAngle = deltaAngle2
		bestOBB = obb2
	}

	deltaAngle3 := math.Abs(boundAngle180(obb3.Rotation()-baseLineAngle))
	if deltaAngle3 < deltaAngle {
		bestOBB = obb3
	}

	return builder.String(), bestOBB
}

func averageFloat[T any](slice []T, extract func(T) float64) float64 {
	if len(slice) == 0 {
		return 0
	}
	var sum float64
	for _, v := range slice {
		sum += extract(v)
	}
	return sum / float64(len(slice))
}

// boundAngle180 bounds angle so that -180 <= theta <= 180.
func boundAngle180(angle float64) float64 {
	angle = math.Mod(angle+180, 360)
	if angle < 0 {
		angle += 360
	}
	return angle - 180
}

// String returns the text of the word.
func (w *Word) String() string {
	return w.Text
}

var _ BoundingBox = (*Word)(nil)
