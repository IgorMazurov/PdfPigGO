package reading_order_detector

import (
	"fmt"
	"math"
	"sort"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

// OrderWordsByReadingOrder orders words by reading order in a line.
// Assumes left-to-right and accounts for rotation.
func OrderWordsByReadingOrder(words []*content.Word) ([]*content.Word, error) {
	if len(words) <= 1 {
		return words, nil
	}

	textOrientation := words[0].TextOrientation
	if textOrientation != content.OtherTextOrientation {
		for _, w := range words {
			if w.TextOrientation != textOrientation {
				textOrientation = content.OtherTextOrientation
				break
			}
		}
	}

	switch textOrientation {
	case content.HorizontalTextOrientation:
		sort.SliceStable(words, func(i, j int) bool {
			return words[i].BoundingBox().BottomLeft.X < words[j].BoundingBox().BottomLeft.X
		})
		return words, nil

	case content.Rotate180TextOrientation:
		sort.SliceStable(words, func(i, j int) bool {
			return words[i].BoundingBox().BottomLeft.X > words[j].BoundingBox().BottomLeft.X
		})
		return words, nil

	case content.Rotate90TextOrientation:
		sort.SliceStable(words, func(i, j int) bool {
			return words[i].BoundingBox().BottomLeft.Y > words[j].BoundingBox().BottomLeft.Y
		})
		return words, nil

	case content.Rotate270TextOrientation:
		sort.SliceStable(words, func(i, j int) bool {
			return words[i].BoundingBox().BottomLeft.Y < words[j].BoundingBox().BottomLeft.Y
		})
		return words, nil

	default:
		var sumAngle float64
		for _, w := range words {
			sumAngle += w.BoundingBox().Rotation()
		}
		avgAngle := sumAngle / float64(len(words))

		if math.IsNaN(avgAngle) {
			return nil, fmt.Errorf("OrderByReadingOrder: NaN bounding box rotation found when ordering words")
		}

		if 0 < avgAngle && avgAngle <= 90 {
			sortTwoKeyAscDesc(words,
				func(w *content.Word) float64 { return w.BoundingBox().BottomLeft.X },
				func(w *content.Word) float64 { return w.BoundingBox().BottomLeft.Y },
				true, true)
			return words, nil
		} else if 90 < avgAngle && avgAngle <= 180 {
			sortTwoKeyAscDesc(words,
				func(w *content.Word) float64 { return -w.BoundingBox().BottomLeft.X },
				func(w *content.Word) float64 { return w.BoundingBox().BottomLeft.Y },
				true, true)
			return words, nil
		} else if -180 < avgAngle && avgAngle <= -90 {
			sortTwoKeyAscDesc(words,
				func(w *content.Word) float64 { return -w.BoundingBox().BottomLeft.X },
				func(w *content.Word) float64 { return -w.BoundingBox().BottomLeft.Y },
				true, true)
			return words, nil
		} else if -90 < avgAngle && avgAngle <= 0 {
			sortTwoKeyAscDesc(words,
				func(w *content.Word) float64 { return w.BoundingBox().BottomLeft.X },
				func(w *content.Word) float64 { return -w.BoundingBox().BottomLeft.Y },
				true, true)
			return words, nil
		}

		return nil, fmt.Errorf("OrderByReadingOrder: unknown bounding box rotation found when ordering words")
	}
}

// OrderLinesByReadingOrder orders lines by reading order in a block.
// Assumes top-to-bottom and accounts for rotation.
func OrderLinesByReadingOrder(lines []*document_layout_analysis.TextLine) ([]*document_layout_analysis.TextLine, error) {
	if len(lines) <= 1 {
		return lines, nil
	}

	textOrientation := lines[0].TextOrientation
	if textOrientation != content.OtherTextOrientation {
		for _, l := range lines {
			if l.TextOrientation != textOrientation {
				textOrientation = content.OtherTextOrientation
				break
			}
		}
	}

	switch textOrientation {
	case content.HorizontalTextOrientation:
		sort.SliceStable(lines, func(i, j int) bool {
			return lines[i].BoundingBox().BottomLeft.Y > lines[j].BoundingBox().BottomLeft.Y
		})
		return lines, nil

	case content.Rotate180TextOrientation:
		sort.SliceStable(lines, func(i, j int) bool {
			return lines[i].BoundingBox().BottomLeft.Y < lines[j].BoundingBox().BottomLeft.Y
		})
		return lines, nil

	case content.Rotate90TextOrientation:
		sort.SliceStable(lines, func(i, j int) bool {
			return lines[i].BoundingBox().BottomLeft.X > lines[j].BoundingBox().BottomLeft.X
		})
		return lines, nil

	case content.Rotate270TextOrientation:
		sort.SliceStable(lines, func(i, j int) bool {
			return lines[i].BoundingBox().BottomLeft.X < lines[j].BoundingBox().BottomLeft.X
		})
		return lines, nil

	default:
		var sumAngle float64
		for _, l := range lines {
			sumAngle += l.BoundingBox().Rotation()
		}
		avgAngle := sumAngle / float64(len(lines))

		if math.IsNaN(avgAngle) {
			return nil, fmt.Errorf("OrderByReadingOrder: NaN bounding box rotation found when ordering lines")
		}

		if 0 < avgAngle && avgAngle <= 90 {
			sortTwoKeyAsc(lines,
				func(l *document_layout_analysis.TextLine) float64 { return -l.BoundingBox().BottomLeft.Y },
				func(l *document_layout_analysis.TextLine) float64 { return l.BoundingBox().BottomLeft.X })
			return lines, nil
		} else if 90 < avgAngle && avgAngle <= 180 {
			sortTwoKeyAsc(lines,
				func(l *document_layout_analysis.TextLine) float64 { return l.BoundingBox().BottomLeft.X },
				func(l *document_layout_analysis.TextLine) float64 { return l.BoundingBox().BottomLeft.Y })
			return lines, nil
		} else if -180 < avgAngle && avgAngle <= -90 {
			sortTwoKeyAsc(lines,
				func(l *document_layout_analysis.TextLine) float64 { return l.BoundingBox().BottomLeft.Y },
				func(l *document_layout_analysis.TextLine) float64 { return -l.BoundingBox().BottomLeft.X })
			return lines, nil
		} else if -90 < avgAngle && avgAngle <= 0 {
			sortTwoKeyAsc(lines,
				func(l *document_layout_analysis.TextLine) float64 { return -l.BoundingBox().BottomLeft.X },
				func(l *document_layout_analysis.TextLine) float64 { return -l.BoundingBox().BottomLeft.Y })
			return lines, nil
		}

		return nil, fmt.Errorf("OrderByReadingOrder: unknown bounding box rotation found when ordering lines")
	}
}

// sortTwoKeyAsc sorts elements by primary key ascending, then secondary key ascending for ties.
func sortTwoKeyAsc[T any](slice []T, primary func(T) float64, secondary func(T) float64) {
	sort.SliceStable(slice, func(i, j int) bool {
		a, b := primary(slice[i]), primary(slice[j])
		if a != b {
			return a < b
		}
		return secondary(slice[i]) < secondary(slice[j])
	})
}

// sortTwoKeyAscDesc sorts elements by primary key (ascending if asc=true), then secondary key
// (ascending if asc2=true) for ties. Negation of keys is used to achieve descending order.
func sortTwoKeyAscDesc[T any](slice []T, primary func(T) float64, secondary func(T) float64, asc, asc2 bool) {
	sort.SliceStable(slice, func(i, j int) bool {
		a, b := primary(slice[i]), primary(slice[j])
		if a != b {
			if asc {
				return a < b
			}
			return a > b
		}
		sa, sb := secondary(slice[i]), secondary(slice[j])
		if asc2 {
			return sa < sb
		}
		return sa > sb
	})
}
