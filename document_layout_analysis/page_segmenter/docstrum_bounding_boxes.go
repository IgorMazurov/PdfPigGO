package page_segmenter

import (
	"math"
	"runtime"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	document_layout_analysis "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
)

// DocstrumBoundingBoxes implements the Document Spectrum (Docstrum) algorithm,
// a bottom-up page segmentation technique based on nearest-neighbourhood clustering
// of connected components extracted from the document. This implementation leverages
// bounding boxes and does not exactly replicate the original algorithm.
// See "The document spectrum for page layout analysis." by L. O'Gorman.
type DocstrumBoundingBoxes struct {
	options DocstrumBoundingBoxesOptions
}

// Instance is a pre-created DocstrumBoundingBoxes with default options.
var Instance = NewDocstrumBoundingBoxes(DefaultDocstrumBoundingBoxesOptions())

// NewDocstrumBoundingBoxes creates a new DocstrumBoundingBoxes page segmenter
// with the given options.
func NewDocstrumBoundingBoxes(options DocstrumBoundingBoxesOptions) *DocstrumBoundingBoxes {
	return &DocstrumBoundingBoxes{options: options}
}

// GetBlocks returns text blocks generated from the given words using the document spectrum method.
func (d *DocstrumBoundingBoxes) GetBlocks(words []*content.Word) []*document_layout_analysis.TextBlock {
	if len(words) == 0 {
		return nil
	}

	filtered := make([]*content.Word, 0, len(words))
	for _, w := range words {
		if strings.TrimSpace(w.Text) != "" {
			filtered = append(filtered, w)
		}
	}
	if len(filtered) == 0 {
		return nil
	}

	return getBlocks(
		filtered,
		d.options.WithinLineBounds, d.options.WithinLineMultiplier, d.options.WithinLineBinSize,
		d.options.BetweenLineBounds, d.options.BetweenLineMultiplier, d.options.BetweenLineBinSize,
		d.options.AngularDifferenceBounds,
		d.options.Epsilon,
		d.options.WordSeparator, d.options.LineSeparator,
		d.options.MaxDegreeOfParallelism,
	)
}

func getBlocks(
	words []*content.Word,
	wlBounds AngleBounds, wlMultiplier float64, wlBinSize int,
	blBounds AngleBounds, blMultiplier float64, blBinSize int,
	angularDifferenceBounds AngleBounds,
	epsilon float64,
	wordSeparator, lineSeparator string,
	maxDegreeOfParallelism int,
) []*document_layout_analysis.TextBlock {

	var withinLineDistance, betweenLineDistance float64
	ok := GetSpacingEstimation(words, wlBounds, wlBinSize, blBounds, blBinSize, maxDegreeOfParallelism, &withinLineDistance, &betweenLineDistance)
	if !ok {
		if math.IsNaN(withinLineDistance) {
			withinLineDistance = 0
		}
		if math.IsNaN(betweenLineDistance) {
			betweenLineDistance = 0
		}
	}

	maxWithinLineDistance := wlMultiplier * withinLineDistance
	lines := GetLines(words, maxWithinLineDistance, wlBounds, wordSeparator, maxDegreeOfParallelism)

	maxBetweenLineDistance := blMultiplier * betweenLineDistance
	return GetStructuralBlocks(lines, maxBetweenLineDistance, angularDifferenceBounds, epsilon, lineSeparator, maxDegreeOfParallelism)
}

// GetSpacingEstimation estimates within-line and between-line spacing using the Docstrum algorithm's first step.
// Returns false if either distance estimate is NaN (no valid data). The estimated distances are written to
// the provided pointer arguments.
func GetSpacingEstimation(
	words []*content.Word,
	wlBounds AngleBounds, wlBinSize int,
	blBounds AngleBounds, blBinSize int,
	maxDegreeOfParallelism int,
	withinLineDistance, betweenLineDistance *float64,
) bool {
	if len(words) == 0 {
		*withinLineDistance = math.NaN()
		*betweenLineDistance = math.NaN()
		return false
	}

	bottomLefts := make([]core.PdfPoint, len(words))
	for i, w := range words {
		bottomLefts[i] = w.BoundingBox().BottomLeft
	}

	kdTree, err := document_layout_analysis.NewKdTree(bottomLefts)
	if err != nil {
		*withinLineDistance = math.NaN()
		*betweenLineDistance = math.NaN()
		return false
	}

	var withinLineDistances, betweenLineDistances []float64
	var mu sync.Mutex

	sem := make(chan struct{}, resolveParallelism(maxDegreeOfParallelism))
	var wg sync.WaitGroup

	for i := 0; i < len(words); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			word := words[idx]

			// Within-line distance
			// C# uses KdTree<Word>.FindNearestNeighbours(word, 2, w => w.BoundingBox.BottomRight, Euclidean)
			// which ranks neighbors by horizontal gap (BottomLeft_candidate to BottomRight_query).
			// Go's non-generic KdTree ranks by spatial distance from query point to stored BottomLeft.
			// Fetch more candidates and sort by horizontal gap to match C# semantics.
			fetchK := 2 * 10 // fetch ~10x to ensure we cover the true nearest by horizontal gap
			if fetchK > len(words) {
				fetchK = len(words)
			}
			wlNeighbours := kdTree.FindNearestNeighbours(word.BoundingBox().BottomRight, fetchK, document_layout_analysis.Euclidean)
			sortByHorizontalGap(wlNeighbours, word.BoundingBox().BottomRight, bottomLefts)
			if len(wlNeighbours) > 2 {
				wlNeighbours = wlNeighbours[:2]
			}
			for _, n := range wlNeighbours {
				neighbourWord := words[n.Index]
				angle := angleWL(word, neighbourWord)
				if wlBounds.Contains(angle) {
					dist := document_layout_analysis.Euclidean(word.BoundingBox().BottomRight, neighbourWord.BoundingBox().BottomLeft)
					mu.Lock()
					withinLineDistances = append(withinLineDistances, dist)
					mu.Unlock()
				}
			}

			// Between-line distance
			// Same approach: fetch more and sort by spatial gap (TopLeft_candidate to TopLeft_query).
			blNeighbours := kdTree.FindNearestNeighbours(word.BoundingBox().TopLeft, fetchK, document_layout_analysis.Euclidean)
			sortBySpatialGap(blNeighbours, word.BoundingBox().TopLeft)
			if len(blNeighbours) > 2 {
				blNeighbours = blNeighbours[:2]
			}
			for _, n := range blNeighbours {
				neighbourWord := words[n.Index]
				angle := angleBL(word, neighbourWord)
				if blBounds.Contains(angle) {
					hypotenuse := document_layout_analysis.Euclidean(word.BoundingBox().Centroid(), neighbourWord.BoundingBox().Centroid())

					a := angle
					if a > 90 {
						a -= 180
					}

					dist := math.Abs(hypotenuse*math.Cos((90-a)*math.Pi/180)) - word.BoundingBox().Height/2.0 - neighbourWord.BoundingBox().Height/2.0

					if dist >= 0 {
						mu.Lock()
						betweenLineDistances = append(betweenLineDistances, dist)
						mu.Unlock()
					}
				}
			}
		}(i)
	}

	wg.Wait()

	withinLinePeak := peakAverageDistance(withinLineDistances, wlBinSize)
	betweenLinePeak := peakAverageDistance(betweenLineDistances, blBinSize)

	if withinLinePeak != nil {
		*withinLineDistance = *withinLinePeak
	} else {
		*withinLineDistance = math.NaN()
	}

	if betweenLinePeak != nil {
		*betweenLineDistance = *betweenLinePeak
	} else {
		*betweenLineDistance = math.NaN()
	}

	return withinLinePeak != nil && betweenLinePeak != nil
}

func peakAverageDistance(distances []float64, binLength int) *float64 {
	if len(distances) == 0 {
		return nil
	}

	if binLength <= 0 {
		panic("DocstrumBoundingBoxes: the bin length must be positive when computing peak average distance")
	}

	maxDbl := distances[0]
	for i := 1; i < len(distances); i++ {
		if distances[i] > maxDbl {
			maxDbl = distances[i]
		}
	}

	maxDbl = math.Ceil(maxDbl)
	if maxDbl > float64(math.MaxInt32) {
		panic("DocstrumBoundingBoxes: error while casting maximum distance to integer")
	}

	max := int(maxDbl)
	if max == 0 {
		max = binLength
	} else if binLength > max {
		binLength = max
	}

	binCount := int(math.Ceil(float64(max)/float64(binLength))) + 1
	bins := make([][]float64, binCount)
	for i := range bins {
		bins[i] = make([]float64, 0)
	}

	for _, distance := range distances {
		bin := int(math.Floor(distance / float64(binLength)))
		if bin < 0 {
			panic("DocstrumBoundingBoxes: negative distance found while computing peak average distance")
		}
		if bin >= binCount {
			bin = binCount - 1
		}
		bins[bin] = append(bins[bin], distance)
	}

	var best []float64
	for _, bin := range bins {
		if len(bin) > len(best) {
			best = bin
		}
	}

	if len(best) == 0 {
		return nil
	}

	sum := 0.0
	for _, v := range best {
		sum += v
	}
	avg := sum / float64(len(best))
	return &avg
}

// GetLines groups words into text lines using nearest-neighbour clustering.
// This is the Docstrum algorithm's second step.
func GetLines(
	words []*content.Word,
	maxWLDistance float64,
	wlBounds AngleBounds,
	wordSeparator string,
	maxDegreeOfParallelism int,
) []*document_layout_analysis.TextLine {
	if len(words) == 0 {
		return nil
	}

	groupedWords := document_layout_analysis.NearestNeighboursPointK(
		words,
		2,
		document_layout_analysis.Euclidean,
		func(_, _ *content.Word) float64 { return maxWLDistance },
		func(pivot *content.Word) core.PdfPoint { return pivot.BoundingBox().BottomRight },
		func(candidate *content.Word) core.PdfPoint { return candidate.BoundingBox().BottomLeft },
		func(_ *content.Word) bool { return true },
		func(pivot, candidate *content.Word) bool {
			return wlBounds.Contains(angleWL(pivot, candidate))
		},
		maxDegreeOfParallelism,
	)

	lines := make([]*document_layout_analysis.TextLine, 0, len(groupedWords))
	for _, g := range groupedWords {
		ordered, err := reading_order_detector.OrderWordsByReadingOrder(g)
		if err != nil {
			ordered = g
		}
		line, err := document_layout_analysis.NewTextLine(ordered, wordSeparator)
		if err == nil && line != nil {
			lines = append(lines, line)
		}
	}
	return lines
}

func angleWL(pivot, candidate *content.Word) float64 {
	angle := document_layout_analysis.BoundAngle180(
		document_layout_analysis.Angle(pivot.BoundingBox().BottomRight, candidate.BoundingBox().BottomLeft) - pivot.BoundingBox().Rotation())

	if angle > 90 {
		angle -= 180
	} else if angle < -90 {
		angle += 180
	}

	return angle
}

// GetStructuralBlocks groups text lines into structural blocks using nearest-neighbour clustering.
// This is the Docstrum algorithm's third and final step.
func GetStructuralBlocks(
	lines []*document_layout_analysis.TextLine,
	maxBLDistance float64,
	angularDifferenceBounds AngleBounds,
	epsilon float64,
	lineSeparator string,
	maxDegreeOfParallelism int,
) []*document_layout_analysis.TextBlock {
	if len(lines) == 0 {
		return nil
	}

	groupedLines := document_layout_analysis.NearestNeighboursLine(
		lines,
		func(l1, l2 core.PdfLine) float64 {
			return perpendicularOverlappingDistance(l1, l2, angularDifferenceBounds, epsilon)
		},
		func(_, _ *document_layout_analysis.TextLine) float64 { return maxBLDistance },
		func(pivot *document_layout_analysis.TextLine) core.PdfLine {
			return core.NewPdfLine(pivot.BoundingBox().BottomLeft, pivot.BoundingBox().BottomRight)
		},
		func(candidate *document_layout_analysis.TextLine) core.PdfLine {
			return core.NewPdfLine(candidate.BoundingBox().TopLeft, candidate.BoundingBox().TopRight)
		},
		func(_ *document_layout_analysis.TextLine) bool { return true },
		func(_, _ *document_layout_analysis.TextLine) bool { return true },
		maxDegreeOfParallelism,
	)

	blocks := make([]*document_layout_analysis.TextBlock, 0, len(groupedLines))
	for _, g := range groupedLines {
		ordered, err := reading_order_detector.OrderLinesByReadingOrder(g)
		if err != nil {
			ordered = g
		}
		block, err := document_layout_analysis.NewTextBlock(ordered, lineSeparator)
		if err == nil && block != nil {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func perpendicularOverlappingDistance(line1, line2 core.PdfLine, angularDifferenceBounds AngleBounds, epsilon float64) float64 {
	ok, theta, _, ed := GetStructuralBlockingParameters(line1, line2, epsilon)
	if !ok {
		return math.Inf(1)
	}

	if theta > 90 {
		theta -= 180
	} else if theta < -90 {
		theta += 180
	}

	if !angularDifferenceBounds.Contains(theta) {
		return math.Inf(1)
	}

	return math.Abs(ed)
}

// GetStructuralBlockingParameters computes the structural blocking parameters between two lines.
// Returns (overlap, angularDifference, normalisedOverlap, perpendicularDistance).
// overlap is true if the lines overlap horizontally.
func GetStructuralBlockingParameters(i, j core.PdfLine, epsilon float64) (overlap bool, angularDifference, normalisedOverlap, perpendicularDistance float64) {
	if almostEquals(i, j, epsilon) {
		return true, 0, 1, 0
	}

	dXi := i.Point2.X - i.Point1.X
	dYi := i.Point2.Y - i.Point1.Y
	dXj := j.Point2.X - j.Point1.X
	dYj := j.Point2.Y - j.Point1.Y

	angularDifference = document_layout_analysis.BoundAngle180((math.Atan2(dYj, dXj) - math.Atan2(dYi, dXi)) * 180 / math.Pi)

	Aj := getTranslatedPoint(i.Point1.X, i.Point1.Y, j.Point1.X, j.Point1.Y, dXi, dYi, dXj, dYj, epsilon)
	Bj := getTranslatedPoint(i.Point2.X, i.Point2.Y, j.Point2.X, j.Point2.Y, dXi, dYi, dXj, dYj, epsilon)

	if Aj == nil || Bj == nil {
		return false, angularDifference, math.NaN(), math.NaN()
	}

	ps := [4]core.PdfPoint{j.Point1, j.Point2, *Aj, *Bj}

	if !document_layout_analysis.AlmostEqualsToZero(dXj, epsilon) {
		sortByXY(ps[:])
	} else if !document_layout_analysis.AlmostEqualsToZero(dYj, epsilon) {
		sortByY(ps[:])
	}

	Cj := ps[1]
	Dj := ps[2]

	overlap = pointInLine(j.Point1, j.Point2, Cj) && pointInLine(j.Point1, j.Point2, Dj) &&
		pointInLine(*Aj, *Bj, Cj) && pointInLine(*Aj, *Bj, Dj)

	pj := document_layout_analysis.Euclidean(Cj, Dj)

	if overlap {
		normalisedOverlap = pj / j.Length()
	} else {
		normalisedOverlap = -pj / j.Length()
	}

	xMj := (Cj.X + Dj.X) / 2.0
	yMj := (Cj.Y + Dj.Y) / 2.0

	if !document_layout_analysis.AlmostEqualsToZero(dXi, epsilon) && !document_layout_analysis.AlmostEqualsToZero(dYi, epsilon) {
		perpendicularDistance = ((xMj - i.Point1.X) - (yMj-i.Point1.Y)*dXi/dYi) / math.Sqrt(dXi*dXi/(dYi*dYi)+1)
	} else if document_layout_analysis.AlmostEqualsToZero(dXi, epsilon) {
		perpendicularDistance = xMj - i.Point1.X
	} else {
		perpendicularDistance = yMj - i.Point1.Y
	}

	return overlap, angularDifference, normalisedOverlap, perpendicularDistance
}

func sortByXY(ps []core.PdfPoint) {
	for i := 1; i < len(ps); i++ {
		for j := i; j > 0; j-- {
			a, b := ps[j-1], ps[j]
			swap := false
			if a.X > b.X || (a.X == b.X && a.Y > b.Y) {
				swap = true
			}
			if swap {
				ps[j-1], ps[j] = ps[j], ps[j-1]
			}
		}
	}
}

func sortByY(ps []core.PdfPoint) {
	for i := 1; i < len(ps); i++ {
		for j := i; j > 0; j-- {
			if ps[j-1].Y > ps[j].Y {
				ps[j-1], ps[j] = ps[j], ps[j-1]
			}
		}
	}
}

func getTranslatedPoint(xPi, yPi, xPj, yPj, dXi, dYi, dXj, dYj, epsilon float64) *core.PdfPoint {
	dYidYj := dYi * dYj
	dXidXj := dXi * dXj
	denominator := dYidYj + dXidXj

	if document_layout_analysis.AlmostEqualsToZero(denominator, epsilon) {
		return nil
	}

	var xAj, yAj float64

	if !document_layout_analysis.AlmostEqualsToZero(dXj, epsilon) {
		xAj = (xPi*dXidXj + xPj*dYidYj + dXj*dYi*(yPi-yPj)) / denominator
		yAj = dYj/dXj*(xAj-xPj) + yPj
	} else {
		yAj = (yPi*dYidYj + yPj*dXidXj + dYj*dXi*(xPi-xPj)) / denominator
		xAj = xPj
	}

	p := core.NewPdfPoint(xAj, yAj)
	return &p
}

func pointInLine(pl1, pl2, point core.PdfPoint) bool {
	ax := point.X - pl1.X
	ay := point.Y - pl1.Y
	bx := pl2.X - pl1.X
	by := pl2.Y - pl1.Y

	dotProd1 := ax*bx + ay*by
	return dotProd1 >= 0 && dotProd1 <= (bx*bx+by*by)
}

func almostEquals(line1, line2 core.PdfLine, epsilon float64) bool {
	return document_layout_analysis.AlmostEqualsToZero(line1.Point1.X-line2.Point1.X, epsilon) &&
		document_layout_analysis.AlmostEqualsToZero(line1.Point1.Y-line2.Point1.Y, epsilon) &&
		document_layout_analysis.AlmostEqualsToZero(line1.Point2.X-line2.Point2.X, epsilon) &&
		document_layout_analysis.AlmostEqualsToZero(line1.Point2.Y-line2.Point2.Y, epsilon)
}

func angleBL(pivot, candidate *content.Word) float64 {
	angle := document_layout_analysis.BoundAngle180(
		document_layout_analysis.Angle(pivot.BoundingBox().Centroid(), candidate.BoundingBox().Centroid()) - pivot.BoundingBox().Rotation())

	if angle < 0 {
		angle += 180
	}

	return angle
}

func sortByHorizontalGap(results []document_layout_analysis.NearestResult, queryBR core.PdfPoint, candidateBLs []core.PdfPoint) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0; j-- {
			a, b := results[j-1], results[j]
			gapA := document_layout_analysis.Euclidean(queryBR, candidateBLs[a.Index])
			gapB := document_layout_analysis.Euclidean(queryBR, candidateBLs[b.Index])
			if gapA > gapB {
				results[j-1], results[j] = results[j], results[j-1]
			}
		}
	}
}

func sortBySpatialGap(results []document_layout_analysis.NearestResult, queryPoint core.PdfPoint) {
	for i := 1; i < len(results); i++ {
		for j := i; j > 0; j-- {
			a, b := results[j-1], results[j]
			distA := document_layout_analysis.Euclidean(queryPoint, a.Point)
			distB := document_layout_analysis.Euclidean(queryPoint, b.Point)
			if distA > distB {
				results[j-1], results[j] = results[j], results[j-1]
			}
		}
	}
}

func resolveParallelism(maxDegreeOfParallelism int) int {
	if maxDegreeOfParallelism <= 0 || maxDegreeOfParallelism == math.MaxInt32 {
		return runtime.GOMAXPROCS(0)
	}
	return maxDegreeOfParallelism
}
