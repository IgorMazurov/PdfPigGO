package testutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
)

// AssertablePositionData holds expected position data for a letter in a PDF page,
// typically parsed from tab-separated test fixture files.
type AssertablePositionData struct {
	X        float64
	Y        float64
	Width    float64
	Text     string
	FontSize float64
	FontName string
	Height   float64
}

// ParseAssertablePositionData parses a tab-separated line into an AssertablePositionData.
// The expected format is: X\tY\tWidth\tText\tFontSize\tFontName[\tHeight]
func ParseAssertablePositionData(line string) (*AssertablePositionData, error) {
	parts := strings.Split(line, "\t")

	if len(parts) < 6 {
		return nil, fmt.Errorf("expected 6 parts to the line, instead got %d", len(parts))
	}

	parts[0] = strings.TrimPrefix(parts[0], "\ufeff")

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X: %w", err)
	}

	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Y: %w", err)
	}

	width, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Width: %w", err)
	}

	fontSize, err := strconv.ParseFloat(parts[4], 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FontSize: %w", err)
	}

	height := float64(0)
	if len(parts) >= 7 {
		h, err := strconv.ParseFloat(parts[6], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Height: %w", err)
		}
		height = h
	}

	return &AssertablePositionData{
		X:        x,
		Y:        y,
		Width:    width,
		Text:     parts[3],
		FontSize: fontSize,
		FontName: parts[5],
		Height:   height,
	}, nil
}

// AssertWithinTolerance compares the expected position data against an actual Letter.
// If includeHeight is true, Height is also compared against letter.BoundingBox.Height.
// Note: the Page parameter from the original C# method was unused and has been omitted.
func (a *AssertablePositionData) AssertWithinTolerance(t *testing.T, letter *content.Letter, includeHeight bool) {
	t.Helper()

	if a.Text != letter.Value {
		t.Errorf("Text mismatch: expected %q, got %q", a.Text, letter.Value)
	}

	if a.FontName != letter.FontName() {
		t.Errorf("FontName mismatch: expected %q, got %q", a.FontName, letter.FontName())
	}

	if !floatWithinTolerance(a.X, letter.Location().X, 1) {
		t.Errorf("X out of tolerance: expected %g, got %g", a.X, letter.Location().X)
	}

	if !floatWithinTolerance(a.Width, letter.Width, 1) {
		t.Errorf("Width out of tolerance: expected %g, got %g", a.Width, letter.Width)
	}

	if includeHeight {
		if !floatWithinTolerance(a.Height, letter.BoundingBox.Height, 1) {
			t.Errorf("Height out of tolerance: expected %g, got %g", a.Height, letter.BoundingBox.Height)
		}
	}
}

// String returns a string representation of the position data.
func (a *AssertablePositionData) String() string {
	return fmt.Sprintf("%g %g %g %s %g %s %g", a.X, a.Y, a.Width, a.Text, a.FontSize, a.FontName, a.Height)
}

// floatWithinTolerance returns true if x and y are within the given tolerance.
func floatWithinTolerance(x, y, tolerance float64) bool {
	return math.Abs(x-y) <= tolerance
}
