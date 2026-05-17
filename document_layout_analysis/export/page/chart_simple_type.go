package page

// ChartSimpleType represents the type of chart in a PAGE XML document.
type ChartSimpleType byte

const (
	// Bar indicates a bar chart.
	Bar ChartSimpleType = iota

	// Line indicates a line chart.
	Line

	// Pie indicates a pie chart.
	Pie

	// Scatter indicates a scatter plot.
	Scatter

	// Surface indicates a surface chart.
	Surface

	// Other indicates an unrecognized chart type.
	Other
)
