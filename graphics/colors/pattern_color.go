package colors

import (
	"bytes"
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PatternType defines the type of pattern color.
type PatternType byte

const (
	// Tiling represents a tiling pattern type.
	Tiling PatternType = 1
	// Shading represents a shading pattern type.
	ShadingPatternType PatternType = 2
)

// String returns the name of the pattern type.
func (pt PatternType) String() string {
	switch pt {
	case Tiling:
		return "Tiling"
	case ShadingPatternType:
		return "Shading"
	default:
		return fmt.Sprintf("Unknown(%d)", pt)
	}
}

// PatternPaintType determines how the colour of the pattern cell shall be specified.
type PatternPaintType byte

const (
	// Coloured indicates a coloured pattern paint type.
	Coloured PatternPaintType = 1
	// Uncoloured indicates an uncoloured pattern paint type.
	Uncoloured PatternPaintType = 2
)

// String returns the name of the pattern paint type.
func (ppt PatternPaintType) String() string {
	switch ppt {
	case Coloured:
		return "Coloured"
	case Uncoloured:
		return "Uncoloured"
	default:
		return fmt.Sprintf("Unknown(%d)", ppt)
	}
}

// PatternTilingType controls adjustments to the spacing of tiles relative to the device pixel grid.
type PatternTilingType byte

const (
	// ConstantSpacing indicates constant spacing pattern tiling type.
	ConstantSpacing PatternTilingType = 1
	// NoDistortion indicates no distortion pattern tiling type.
	NoDistortion PatternTilingType = 2
	// ConstantSpacingFasterTiling indicates constant spacing faster tiling pattern tiling type.
	ConstantSpacingFasterTiling PatternTilingType = 3
)

// String returns the name of the pattern tiling type.
func (ptt PatternTilingType) String() string {
	switch ptt {
	case ConstantSpacing:
		return "ConstantSpacing"
	case NoDistortion:
		return "NoDistortion"
	case ConstantSpacingFasterTiling:
		return "ConstantSpacingFasterTiling"
	default:
		return fmt.Sprintf("Unknown(%d)", ptt)
	}
}

// PatternColor represents a pattern-based color (tiling or shading).
// Base interface for TilingPatternColor and ShadingPatternColor.
type PatternColor interface {
	// colorInterface marks this as implementing the Color interface methods below.
	colorInterface()

	// PatternType returns 1 for tiling, 2 for shading.
	PatternType() PatternType

	// PatternDictionary returns the dictionary defining the pattern.
	PatternDictionary() *tokens.DictionaryToken

	// ExtGState returns the graphics state parameter dictionary containing graphics state
	// parameters to be put into effect temporarily while the shading pattern is painted.
	ExtGState() *tokens.DictionaryToken

	// Matrix returns the pattern matrix. Default value: the identity matrix [1 0 0 1 0 0].
	Matrix() core.TransformationMatrix
}

// TilingPatternColor represents a tiling pattern color.
// A tiling pattern consists of a small graphical figure called a pattern cell. Painting with the pattern
// replicates the cell at fixed horizontal and vertical intervals to fill an area.
type TilingPatternColor struct {
	patternDictionary *tokens.DictionaryToken
	extGState         *tokens.DictionaryToken
	matrix            core.TransformationMatrix
	patternStream     *tokens.StreamToken
	paintType         PatternPaintType
	tilingType        PatternTilingType
	bBox              core.PdfRectangle
	xStep             float64
	yStep             float64
	resources         *tokens.DictionaryToken
	data              []byte
}

// NewTilingPatternColor creates a new TilingPatternColor.
func NewTilingPatternColor(
	matrix core.TransformationMatrix,
	extGState *tokens.DictionaryToken,
	patternStream *tokens.StreamToken,
	paintType PatternPaintType,
	tilingType PatternTilingType,
	bBox core.PdfRectangle,
	xStep, yStep float64,
	resources *tokens.DictionaryToken,
	data []byte,
) TilingPatternColor {
	return TilingPatternColor{
		patternDictionary: patternStream.StreamDictionary,
		extGState:         extGState,
		matrix:            matrix,
		patternStream:     patternStream,
		paintType:         paintType,
		tilingType:        tilingType,
		bBox:              bBox,
		xStep:             xStep,
		yStep:             yStep,
		resources:         resources,
		data:              data,
	}
}

func (TilingPatternColor) colorInterface() {}

// PatternType returns Tiling for this pattern.
func (t TilingPatternColor) PatternType() PatternType {
	return Tiling
}

// PatternDictionary returns the dictionary defining the pattern.
func (t TilingPatternColor) PatternDictionary() *tokens.DictionaryToken {
	return t.patternDictionary
}

// ExtGState returns the graphics state parameter dictionary.
func (t TilingPatternColor) ExtGState() *tokens.DictionaryToken {
	return t.extGState
}

// Matrix returns the pattern matrix.
func (t TilingPatternColor) Matrix() core.TransformationMatrix {
	return t.matrix
}

// PatternStream returns the content stream containing the painting operators needed to paint one instance of the cell.
func (t TilingPatternColor) PatternStream() *tokens.StreamToken {
	return t.patternStream
}

// PaintType returns the code that determines how the colour of the pattern cell shall be specified.
func (t TilingPatternColor) PaintType() PatternPaintType {
	return t.paintType
}

// TilingType returns the code that controls adjustments to the spacing of tiles relative to the device pixel grid.
func (t TilingPatternColor) TilingType() PatternTilingType {
	return t.tilingType
}

// BBox returns the pattern cell's bounding box used to clip the pattern cell.
func (t TilingPatternColor) BBox() core.PdfRectangle {
	return t.bBox
}

// XStep returns the desired horizontal spacing between pattern cells, measured in the pattern coordinate system.
func (t TilingPatternColor) XStep() float64 {
	return t.xStep
}

// YStep returns the desired vertical spacing between pattern cells, measured in the pattern coordinate system.
func (t TilingPatternColor) YStep() float64 {
	return t.yStep
}

// Resources returns the resource dictionary containing all named resources required by the pattern's content stream.
func (t TilingPatternColor) Resources() *tokens.DictionaryToken {
	return t.resources
}

// Data returns the raw content bytes of the painting operators needed to paint one instance of the cell.
func (t TilingPatternColor) Data() []byte {
	return t.data
}

// ColorSpace returns Pattern for this color.
func (TilingPatternColor) ColorSpace() ColorSpace {
	return Pattern
}

// ToRGBValues cannot be called for a pattern color and panics.
func (TilingPatternColor) ToRGBValues() RGBValues {
	panic("cannot call ToRGBValues on a Pattern color")
}

// Equals reports whether this TilingPatternColor equals another.
func (t TilingPatternColor) Equals(other TilingPatternColor) bool {
	return t.patternDictionary.Equals(other.patternDictionary) &&
		t.matrix.Equals(other.matrix) &&
		(t.extGState == nil && other.extGState == nil || (t.extGState != nil && t.extGState.Equals(other.extGState))) &&
		t.paintType == other.paintType &&
		t.tilingType == other.tilingType &&
		t.bBox.Equals(other.bBox) &&
		t.xStep == other.xStep &&
		t.yStep == other.yStep &&
		t.resources.Equals(other.resources) &&
		bytes.Equal(t.data, other.data)
}

// String returns a string representation of the tiling pattern color.
func (t TilingPatternColor) String() string {
	return fmt.Sprintf("Pattern: (%s)[%s][%s]", t.PatternType(), t.paintType, t.tilingType)
}

// ShadingPatternColor represents a shading pattern color.
// Shading patterns provide a smooth transition between colours across an area to be painted, independent of
// the resolution of any particular output device and without specifying the number of steps in the colour transition.
type ShadingPatternColor struct {
	patternDictionary *tokens.DictionaryToken
	extGState         *tokens.DictionaryToken
	matrix            core.TransformationMatrix
	shading           Shading
}

// NewShadingPatternColor creates a new ShadingPatternColor.
func NewShadingPatternColor(
	matrix core.TransformationMatrix,
	extGState *tokens.DictionaryToken,
	patternDictionary *tokens.DictionaryToken,
	shading Shading,
) ShadingPatternColor {
	return ShadingPatternColor{
		patternDictionary: patternDictionary,
		extGState:         extGState,
		matrix:            matrix,
		shading:           shading,
	}
}

func (ShadingPatternColor) colorInterface() {}

// PatternType returns Shading for this pattern.
func (s ShadingPatternColor) PatternType() PatternType {
	return ShadingPatternType
}

// PatternDictionary returns the dictionary defining the pattern.
func (s ShadingPatternColor) PatternDictionary() *tokens.DictionaryToken {
	return s.patternDictionary
}

// ExtGState returns the graphics state parameter dictionary.
func (s ShadingPatternColor) ExtGState() *tokens.DictionaryToken {
	return s.extGState
}

// Matrix returns the pattern matrix.
func (s ShadingPatternColor) Matrix() core.TransformationMatrix {
	return s.matrix
}

// Shading returns the shading object defining the shading pattern's gradient fill.
func (s ShadingPatternColor) Shading() Shading {
	return s.shading
}

// ColorSpace returns Pattern for this color.
func (ShadingPatternColor) ColorSpace() ColorSpace {
	return Pattern
}

// ToRGBValues cannot be called for a pattern color and panics.
func (ShadingPatternColor) ToRGBValues() RGBValues {
	panic("cannot call ToRGBValues on a Pattern color")
}

// Equals reports whether this ShadingPatternColor equals another.
func (s ShadingPatternColor) Equals(other ShadingPatternColor) bool {
	return s.patternDictionary.Equals(other.patternDictionary) &&
		s.matrix.Equals(other.matrix) &&
		(s.extGState == nil && other.extGState == nil || (s.extGState != nil && s.extGState.Equals(other.extGState))) &&
		s.shading.Equals(other.shading)
}

// String returns a string representation of the shading pattern color.
func (s ShadingPatternColor) String() string {
	return fmt.Sprintf("Pattern: (%s)", s.PatternType())
}
