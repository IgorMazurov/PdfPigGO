package page_segmenter

import (
	"sort"

	"github.com/uglytoad/pdfpig/go/content"
	document_layout_analysis "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/document_layout_analysis/reading_order_detector"
)

// DefaultPageSegmenter puts all words into a single text block.
type DefaultPageSegmenter struct {
	options DefaultPageSegmenterOptions
}

// DefaultInstance is a pre-created DefaultPageSegmenter with default options.
var DefaultInstance = NewDefaultPageSegmenter(DefaultDefaultPageSegmenterOptions())

// NewDefaultPageSegmenter creates a new DefaultPageSegmenter with the given options.
func NewDefaultPageSegmenter(options DefaultPageSegmenterOptions) *DefaultPageSegmenter {
	return &DefaultPageSegmenter{options: options}
}

// GetBlocks returns all words grouped into one text block.
func (d *DefaultPageSegmenter) GetBlocks(words []*content.Word) []*document_layout_analysis.TextBlock {
	if len(words) == 0 {
		return nil
	}

	lines := getLines(words, d.options.WordSeparator())
	if len(lines) == 0 {
		return nil
	}

	tb, err := document_layout_analysis.NewTextBlock(lines, d.options.LineSeparator())
	if err != nil {
		return nil
	}

	return []*document_layout_analysis.TextBlock{tb}
}

// getLines groups words by their bottom Y coordinate, orders each group by reading
// order to form a line, then orders the resulting lines by reading order.
func getLines(words []*content.Word, wordSeparator string) []*document_layout_analysis.TextLine {
	groups := make(map[float64][]*content.Word)
	for _, w := range words {
		bb := w.BoundingBox()
		y := bb.BottomLeft.Y
		groups[y] = append(groups[y], w)
	}

	lines := make([]*document_layout_analysis.TextLine, 0, len(groups))
	for _, groupWords := range groups {
		ordered, err := reading_order_detector.OrderWordsByReadingOrder(groupWords)
		if err != nil {
			continue
		}

		line, err := document_layout_analysis.NewTextLine(ordered, wordSeparator)
		if err != nil {
			continue
		}
		lines = append(lines, line)
	}

	ordered, err := reading_order_detector.OrderLinesByReadingOrder(lines)
	if err == nil {
		return ordered
	}

	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i].BoundingBox().BottomLeft.Y > lines[j].BoundingBox().BottomLeft.Y
	})

	return lines
}

// DefaultPageSegmenterOptions holds configuration for the default page segmenter.
type DefaultPageSegmenterOptions struct {
	maxDegreeOfParallelism int
	wordSeparator          string
	lineSeparator          string
}

// DefaultDefaultPageSegmenterOptions returns options with standard defaults.
func DefaultDefaultPageSegmenterOptions() DefaultPageSegmenterOptions {
	return DefaultPageSegmenterOptions{
		maxDegreeOfParallelism: -1,
		wordSeparator:          " ",
		lineSeparator:          "\n",
	}
}

// MaxDegreeOfParallelism returns the maximum number of concurrent tasks enabled.
func (o DefaultPageSegmenterOptions) MaxDegreeOfParallelism() int {
	return o.maxDegreeOfParallelism
}

// SetMaxDegreeOfParallelism sets the maximum number of concurrent tasks enabled.
func (o *DefaultPageSegmenterOptions) SetMaxDegreeOfParallelism(value int) {
	o.maxDegreeOfParallelism = value
}

// WordSeparator returns the separator used between words when building lines.
func (o DefaultPageSegmenterOptions) WordSeparator() string {
	return o.wordSeparator
}

// SetWordSeparator sets the separator used between words when building lines.
func (o *DefaultPageSegmenterOptions) SetWordSeparator(value string) {
	o.wordSeparator = value
}

// LineSeparator returns the separator used between lines when building blocks.
func (o DefaultPageSegmenterOptions) LineSeparator() string {
	return o.lineSeparator
}

// SetLineSeparator sets the separator used between lines when building blocks.
func (o *DefaultPageSegmenterOptions) SetLineSeparator(value string) {
	o.lineSeparator = value
}

var _ PageSegmenter = (*DefaultPageSegmenter)(nil)
var _ PageSegmenterOptions = (*DefaultPageSegmenterOptions)(nil)
