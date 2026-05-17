package word_extractor

import (
	"math"
	"runtime"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

// NearestNeighbourWordExtractor groups letters into words using a nearest-neighbour
// algorithm based on bounding-box proximity. Letters whose baseline endpoints fall within
// a configurable maximum distance are considered part of the same word.
type NearestNeighbourWordExtractor struct {
	options *NearestNeighbourWordExtractorOptions
}

var _ content.WordExtractor = (*NearestNeighbourWordExtractor)(nil)

// DefaultInstance returns a shared NearestNeighbourWordExtractor with default options.
var DefaultInstance content.WordExtractor = NewNearestNeighbourWordExtractor()

// NewNearestNeighbourWordExtractor creates an extractor using the default options.
func NewNearestNeighbourWordExtractor() *NearestNeighbourWordExtractor {
	return NewNearestNeighbourWordExtractorWithOptions(DefaultNearestNeighbourWordExtractorOptions())
}

// NewNearestNeighbourWordExtractorWithOptions creates an extractor with custom options.
func NewNearestNeighbourWordExtractorWithOptions(options *NearestNeighbourWordExtractorOptions) *NearestNeighbourWordExtractor {
	if options == nil {
		options = DefaultNearestNeighbourWordExtractorOptions()
	}
	return &NearestNeighbourWordExtractor{options: options}
}

// GetWords groups the page letters into words using the nearest-neighbour method.
func (e *NearestNeighbourWordExtractor) GetWords(letters []*content.Letter) []*content.Word {
	if len(letters) == 0 {
		return nil
	}

	if e.options.GroupByOrientation() {
		buckets := [5][]*content.Letter{}
		for i := range buckets {
			buckets[i] = make([]*content.Letter, 0)
		}

		for _, l := range letters {
			switch l.TextOrientation {
			case content.HorizontalTextOrientation:
				buckets[0] = append(buckets[0], l)
			case content.Rotate270TextOrientation:
				buckets[1] = append(buckets[1], l)
			case content.Rotate180TextOrientation:
				buckets[2] = append(buckets[2], l)
			case content.Rotate90TextOrientation:
				buckets[3] = append(buckets[3], l)
			default:
				buckets[4] = append(buckets[4], l)
			}
		}

		var mu sync.Mutex
		results := make([]*content.Word, 0, len(letters))

		sem := make(chan struct{}, resolveParallelism(e.options.MaxDegreeOfParallelism()))
		var wg sync.WaitGroup

		for i := 0; i < len(buckets); i++ {
			if len(buckets[i]) == 0 {
				continue
			}
			wg.Add(1)
			go func(bucketIdx int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				measure := e.options.DistanceMeasureAA()
				if bucketIdx == 4 {
					measure = e.options.DistanceMeasure()
				}

				words := extractWords(
					buckets[bucketIdx],
					e.options.MaximumDistance(),
					measure,
					e.options.FilterPivot(),
					e.options.Filter(),
					e.options.MaxDegreeOfParallelism(),
				)

				mu.Lock()
				results = append(results, words...)
				mu.Unlock()
			}(i)
		}

		wg.Wait()
		return results
	}

	return extractWords(
		letters,
		e.options.MaximumDistance(),
		e.options.DistanceMeasure(),
		e.options.FilterPivot(),
		e.options.Filter(),
		e.options.MaxDegreeOfParallelism(),
	)
}

func extractWords(
	letters []*content.Letter,
	maxDistance func(*content.Letter, *content.Letter) float64,
	distMeasure func(core.PdfPoint, core.PdfPoint) float64,
	filterPivot func(*content.Letter) bool,
	filter func(*content.Letter, *content.Letter) bool,
	maxDegreeOfParallelism int,
) []*content.Word {
	if len(letters) == 0 {
		return nil
	}

	grouped := document_layout_analysis.NearestNeighboursPoint(
		letters,
		distMeasure,
		maxDistance,
		func(l *content.Letter) core.PdfPoint { return l.EndBaseLine },
		func(l *content.Letter) core.PdfPoint { return l.StartBaseLine },
		filterPivot,
		filter,
		maxDegreeOfParallelism,
	)

	words := make([]*content.Word, 0, len(grouped))
	for _, g := range grouped {
		if word, err := content.NewWord(g); err == nil {
			words = append(words, word)
		}
	}
	return words
}

// NearestNeighbourWordExtractorOptions holds configuration for the nearest-neighbour
// word extractor. Implements WordExtractorOptions and DlaOptions via embedding.
type NearestNeighbourWordExtractorOptions struct {
	maxDegreeOfParallelism int
	maximumDistance        func(*content.Letter, *content.Letter) float64
	distanceMeasure        func(core.PdfPoint, core.PdfPoint) float64
	distanceMeasureAA      func(core.PdfPoint, core.PdfPoint) float64
	filter                 func(*content.Letter, *content.Letter) bool
	filterPivot            func(*content.Letter) bool
	groupByOrientation     bool
}

var _ WordExtractorOptions = (*NearestNeighbourWordExtractorOptions)(nil)

// DefaultNearestNeighbourWordExtractorOptions returns options with the recommended defaults.
func DefaultNearestNeighbourWordExtractorOptions() *NearestNeighbourWordExtractorOptions {
	return &NearestNeighbourWordExtractorOptions{
		maxDegreeOfParallelism: -1,
		maximumDistance: func(l1, l2 *content.Letter) float64 {
			maxDist := math.Max(
				math.Max(
					math.Max(
						math.Max(
							math.Max(
								math.Abs(l1.BoundingBox.Width),
								math.Abs(l2.BoundingBox.Width)),
							math.Abs(l1.Width)),
						math.Abs(l2.Width)),
					l1.PointSize),
				l2.PointSize,
			) * 0.2

			if l1.TextOrientation == content.OtherTextOrientation || l2.TextOrientation == content.OtherTextOrientation {
				return 2.0 * maxDist
			}
			return maxDist
		},
		distanceMeasure:    document_layout_analysis.Euclidean,
		distanceMeasureAA:  document_layout_analysis.Manhattan,
		filter:             func(_ *content.Letter, l2 *content.Letter) bool { return !isWhitespace(l2.Value) },
		filterPivot:        func(l *content.Letter) bool { return !isWhitespace(l.Value) },
		groupByOrientation: true,
	}
}

// MaxDegreeOfParallelism returns the maximum number of concurrent tasks enabled.
func (o *NearestNeighbourWordExtractorOptions) MaxDegreeOfParallelism() int {
	return o.maxDegreeOfParallelism
}

// SetMaxDegreeOfParallelism sets the maximum number of concurrent tasks.
func (o *NearestNeighbourWordExtractorOptions) SetMaxDegreeOfParallelism(value int) {
	o.maxDegreeOfParallelism = value
}

// MaximumDistance returns the function that determines the maximum distance between two letters.
func (o *NearestNeighbourWordExtractorOptions) MaximumDistance() func(*content.Letter, *content.Letter) float64 {
	return o.maximumDistance
}

// SetMaximumDistance sets the maximum distance function.
func (o *NearestNeighbourWordExtractorOptions) SetMaximumDistance(fn func(*content.Letter, *content.Letter) float64) {
	if fn != nil {
		o.maximumDistance = fn
	}
}

// DistanceMeasure returns the distance measure between two baseline points.
func (o *NearestNeighbourWordExtractorOptions) DistanceMeasure() func(core.PdfPoint, core.PdfPoint) float64 {
	return o.distanceMeasure
}

// SetDistanceMeasure sets the default distance measure function.
func (o *NearestNeighbourWordExtractorOptions) SetDistanceMeasure(fn func(core.PdfPoint, core.PdfPoint) float64) {
	if fn != nil {
		o.distanceMeasure = fn
	}
}

// DistanceMeasureAA returns the distance measure for axis-aligned text orientation.
func (o *NearestNeighbourWordExtractorOptions) DistanceMeasureAA() func(core.PdfPoint, core.PdfPoint) float64 {
	return o.distanceMeasureAA
}

// SetDistanceMeasureAA sets the axis-aligned distance measure function.
func (o *NearestNeighbourWordExtractorOptions) SetDistanceMeasureAA(fn func(core.PdfPoint, core.PdfPoint) float64) {
	if fn != nil {
		o.distanceMeasureAA = fn
	}
}

// Filter returns the function used to filter connections between letters.
func (o *NearestNeighbourWordExtractorOptions) Filter() func(*content.Letter, *content.Letter) bool {
	return o.filter
}

// SetFilter sets the connection filter function.
func (o *NearestNeighbourWordExtractorOptions) SetFilter(fn func(*content.Letter, *content.Letter) bool) {
	if fn != nil {
		o.filter = fn
	}
}

// FilterPivot returns the function used to decide whether a letter should search for neighbours.
func (o *NearestNeighbourWordExtractorOptions) FilterPivot() func(*content.Letter) bool {
	return o.filterPivot
}

// SetFilterPivot sets the pivot filter function.
func (o *NearestNeighbourWordExtractorOptions) SetFilterPivot(fn func(*content.Letter) bool) {
	if fn != nil {
		o.filterPivot = fn
	}
}

// GroupByOrientation returns whether letters are grouped by text orientation before processing.
func (o *NearestNeighbourWordExtractorOptions) GroupByOrientation() bool {
	return o.groupByOrientation
}

// SetGroupByOrientation sets whether to group letters by text orientation.
func (o *NearestNeighbourWordExtractorOptions) SetGroupByOrientation(value bool) {
	o.groupByOrientation = value
}

func resolveParallelism(maxDegreeOfParallelism int) int {
	if maxDegreeOfParallelism <= 0 {
		return runtime.GOMAXPROCS(0)
	}
	return maxDegreeOfParallelism
}

func isWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}
