package cmap

// WritingMode defines the text writing direction used in CMap parsing.
type WritingMode int

const (
	// Horizontal represents left-to-right horizontal text flow.
	Horizontal WritingMode = iota

	// Vertical represents top-to-bottom vertical text flow.
	Vertical
)
