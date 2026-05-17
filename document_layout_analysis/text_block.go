package document_layout_analysis

import (
	"errors"
	"math"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
)

// TextBlock represents a block of text composed of one or more text lines.
type TextBlock struct {
	Separator       string
	Text            string
	TextOrientation content.TextOrientation
	BoundingBox     core.PdfRectangle
	TextLines       []*TextLine
	ReadingOrder    int
}

// NewTextBlock creates a new TextBlock from an ordered list of text lines.
func NewTextBlock(lines []*TextLine, separator string) (*TextBlock, error) {
	if len(lines) == 0 {
		return nil, errors.New("empty lines provided")
	}

	if separator == "" {
		separator = "\n"
	}

	var tb TextBlock
	tb.Separator = separator
	tb.ReadingOrder = -1
	tb.TextLines = lines

	if len(lines) == 1 {
		tb.BoundingBox = lines[0].BoundingBox()
		tb.Text = lines[0].Text
		tb.TextOrientation = lines[0].TextOrientation
	} else {
		tempOrientation := lines[0].TextOrientation
		if tempOrientation != content.OtherTextOrientation {
			for _, line := range lines {
				if line.TextOrientation != tempOrientation {
					tempOrientation = content.OtherTextOrientation
					break
				}
			}
		}

		switch tempOrientation {
		case content.HorizontalTextOrientation:
			tb.BoundingBox = getBlockBoundingBoxH(lines)
		case content.Rotate180TextOrientation:
			tb.BoundingBox = getBlockBoundingBox180(lines)
		case content.Rotate90TextOrientation:
			tb.BoundingBox = getBlockBoundingBox90(lines)
		case content.Rotate270TextOrientation:
			tb.BoundingBox = getBlockBoundingBox270(lines)
		default:
			var err error
			tb.BoundingBox, err = getBlockBoundingBoxOther(lines)
			if err != nil {
				return nil, err
			}
		}

		var parts []string
		for _, l := range lines {
			parts = append(parts, l.Text)
		}
		tb.Text = strings.Join(parts, separator)
		tb.TextOrientation = tempOrientation
	}

	return &tb, nil
}

func getBlockBoundingBoxH(lines []*TextLine) core.PdfRectangle {
	blX := math.Inf(1)
	trX := math.Inf(-1)
	blY := math.Inf(1)
	trY := math.Inf(-1)

	for _, line := range lines {
		bb := line.BoundingBox()
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

func getBlockBoundingBox180(lines []*TextLine) core.PdfRectangle {
	blX := math.Inf(-1)
	blY := math.Inf(-1)
	trX := math.Inf(1)
	trY := math.Inf(1)

	for _, line := range lines {
		bb := line.BoundingBox()
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

func getBlockBoundingBox90(lines []*TextLine) core.PdfRectangle {
	b := math.Inf(1)
	r := math.Inf(1)
	t := math.Inf(-1)
	l := math.Inf(-1)

	for _, line := range lines {
		bb := line.BoundingBox()
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

func getBlockBoundingBox270(lines []*TextLine) core.PdfRectangle {
	t := math.Inf(1)
	b := math.Inf(-1)
	l := math.Inf(1)
	r := math.Inf(-1)

	for _, line := range lines {
		bb := line.BoundingBox()
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

func getBlockBoundingBoxOther(lines []*TextLine) (core.PdfRectangle, error) {
	var points []core.PdfPoint
	for _, l := range lines {
		bb := l.BoundingBox()
		points = append(points, bb.BottomLeft, bb.BottomRight, bb.TopLeft, bb.TopRight)
	}

	obb, err := geometry.MinimumAreaRectangle(points)
	if err != nil {
		return core.PdfRectangle{}, err
	}

	obb1 := core.NewPdfRectangleFromCorners(obb.BottomLeft, obb.TopLeft, obb.BottomRight, obb.TopRight)
	obb2 := core.NewPdfRectangleFromCorners(obb.BottomRight, obb.BottomLeft, obb.TopRight, obb.TopLeft)
	obb3 := core.NewPdfRectangleFromCorners(obb.TopRight, obb.BottomRight, obb.TopLeft, obb.BottomLeft)

	lastLine := lines[len(lines)-1]
	baseLineAngle := BoundAngle180(Angle(lastLine.BoundingBox().BottomLeft, lastLine.BoundingBox().BottomRight))

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

	return bestOBB, nil
}

// SetReadingOrder sets the TextBlock's reading order.
func (tb *TextBlock) SetReadingOrder(readingOrder int) error {
	if readingOrder < -1 {
		return errors.New("the reading order should be more or equal to -1. A value of -1 means the block is not ordered")
	}
	tb.ReadingOrder = readingOrder
	return nil
}

// String returns the text of the block.
func (tb *TextBlock) String() string {
	return tb.Text
}
