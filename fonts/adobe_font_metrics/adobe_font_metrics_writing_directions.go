package adobe_font_metrics

// AdobeFontMetricsWritingDirection indicates the meaning of the metric sets field in an AFM file.
type AdobeFontMetricsWritingDirection byte

const (
	// Direction0Only means writing direction 0 only.
	Direction0Only AdobeFontMetricsWritingDirection = iota
	// Direction1Only means writing direction 1 only.
	Direction1Only
	// Direction0And1 means writing direction 0 and 1.
	Direction0And1
)
