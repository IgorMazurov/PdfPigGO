package page_segmenter

// DocstrumBoundingBoxesOptions holds configuration for the Docstrum bounding boxes
// page segmenter algorithm.
type DocstrumBoundingBoxesOptions struct {
	MaxDegreeOfParallelism int
	WordSeparator          string
	LineSeparator          string
	Epsilon                float64
	WithinLineBounds       AngleBounds
	WithinLineMultiplier   float64
	WithinLineBinSize      int
	BetweenLineBounds      AngleBounds
	BetweenLineMultiplier  float64
	BetweenLineBinSize     int
	AngularDifferenceBounds AngleBounds
}

// DefaultDocstrumBoundingBoxesOptions returns the default options for the Docstrum
// bounding boxes page segmenter.
func DefaultDocstrumBoundingBoxesOptions() DocstrumBoundingBoxesOptions {
	wl, _ := NewAngleBounds(-30, 30)
	bl, _ := NewAngleBounds(45, 135)
	ad, _ := NewAngleBounds(-30, 30)
	return DocstrumBoundingBoxesOptions{
		MaxDegreeOfParallelism: -1,
		WordSeparator:          " ",
		LineSeparator:          "\n",
		Epsilon:                1e-3,
		WithinLineBounds:       wl,
		WithinLineMultiplier:   3.0,
		WithinLineBinSize:      10,
		BetweenLineBounds:      bl,
		BetweenLineMultiplier:  1.3,
		BetweenLineBinSize:     10,
		AngularDifferenceBounds: ad,
	}
}

// MaxDegreeOfParallelism returns the maximum number of concurrent tasks enabled.
func (o DocstrumBoundingBoxesOptions) MaxDegreeOfParallelismOpt() int {
	return o.MaxDegreeOfParallelism
}

// SetMaxDegreeOfParallelism sets the maximum number of concurrent tasks enabled.
func (o *DocstrumBoundingBoxesOptions) SetMaxDegreeOfParallelism(value int) {
	o.MaxDegreeOfParallelism = value
}

// WordSeparator returns the separator used between words when building lines.
func (o DocstrumBoundingBoxesOptions) WordSeparatorOpt() string {
	return o.WordSeparator
}

// SetWordSeparator sets the separator used between words when building lines.
func (o *DocstrumBoundingBoxesOptions) SetWordSeparator(value string) {
	o.WordSeparator = value
}

// LineSeparator returns the separator used between lines when building paragraphs.
func (o DocstrumBoundingBoxesOptions) LineSeparatorOpt() string {
	return o.LineSeparator
}

// SetLineSeparator sets the separator used between lines when building paragraphs.
func (o *DocstrumBoundingBoxesOptions) SetLineSeparator(value string) {
	o.LineSeparator = value
}
