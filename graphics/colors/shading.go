package colors

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/functions"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ShadingType defines the type of shading as specified by the PDF specification.
type ShadingType byte

const (
	// FunctionBased represents function-based shading (type 1).
	FunctionBased ShadingType = 1
	// Axial represents axial shading (type 2).
	Axial ShadingType = 2
	// Radial represents radial shading (type 3).
	Radial ShadingType = 3
	// FreeFormGouraud represents free-form Gouraud-shaded triangle mesh (type 4).
	FreeFormGouraud ShadingType = 4
	// LatticeFormGouraud represents lattice-form Gouraud-shaded triangle mesh (type 5).
	LatticeFormGouraud ShadingType = 5
	// CoonsPatch represents Coons patch mesh (type 6).
	CoonsPatch ShadingType = 6
	// TensorProductPatch represents tensor-product patch mesh (type 7).
	TensorProductPatch ShadingType = 7
)

// String returns the name of the shading type.
func (st ShadingType) String() string {
	switch st {
	case FunctionBased:
		return "FunctionBased"
	case Axial:
		return "Axial"
	case Radial:
		return "Radial"
	case FreeFormGouraud:
		return "FreeFormGouraud"
	case LatticeFormGouraud:
		return "LatticeFormGouraud"
	case CoonsPatch:
		return "CoonsPatch"
	case TensorProductPatch:
		return "TensorProductPatch"
	default:
		return ""
	}
}

// Shading specifies details of a particular gradient fill, including the type
// of shading to be used, the geometry of the area to be shaded, and the geometry of the
// gradient fill. Various shading types are available, depending on the value of the
// dictionary's ShadingType entry.
type Shading interface {
	// ShadingDictionary returns the dictionary defining the shading.
	ShadingDictionary() *tokens.DictionaryToken

	// ShadingType returns the type of this shading.
	ShadingType() ShadingType

	// ColorSpace returns the colour space in which colour values shall be expressed.
	ColorSpace() ColorSpaceDetails

	// Background returns the background colour value, or nil if not specified.
	Background() []float64

	// BBox returns the shading's bounding box, or nil if not specified.
	BBox() *core.PdfRectangle

	// AntiAlias returns whether anti-aliasing is enabled for this shading.
	AntiAlias() bool

	// Functions returns the shading's function(s), or nil if none.
	Functions() []functions.PdfFunction

	// Eval converts the input values using the functions of the shading.
	Eval(input ...float64) []float64

	// Equals reports whether this shading equals another.
	Equals(other Shading) bool
}

// FunctionBasedShading defines the colour of every point in the domain using a
// mathematical function (not necessarily smooth or continuous).
type FunctionBasedShading struct {
	shadingDictionary *tokens.DictionaryToken
	colorSpace        ColorSpaceDetails
	background        []float64
	bBox              *core.PdfRectangle
	antiAlias         bool
	domain            []float64
	matrix            core.TransformationMatrix
	functions         []functions.PdfFunction
}

// NewFunctionBasedShading creates a new FunctionBasedShading.
func NewFunctionBasedShading(
	shadingDictionary *tokens.DictionaryToken,
	colorSpace ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	domain []float64,
	matrix core.TransformationMatrix,
	functions []functions.PdfFunction,
) FunctionBasedShading {
	return FunctionBasedShading{
		shadingDictionary: shadingDictionary,
		colorSpace:        colorSpace,
		background:        background,
		bBox:              bBox,
		antiAlias:         antiAlias,
		domain:            domain,
		matrix:            matrix,
		functions:         functions,
	}
}

// ShadingDictionary returns the dictionary defining the shading.
func (f FunctionBasedShading) ShadingDictionary() *tokens.DictionaryToken {
	return f.shadingDictionary
}

// ShadingType returns FunctionBased for this shading.
func (FunctionBasedShading) ShadingType() ShadingType {
	return FunctionBased
}

// ColorSpace returns the colour space of the shading.
func (f FunctionBasedShading) ColorSpace() ColorSpaceDetails {
	return f.colorSpace
}

// Background returns the background colour value, or nil if not specified.
func (f FunctionBasedShading) Background() []float64 {
	return f.background
}

// BBox returns the shading's bounding box, or nil if not specified.
func (f FunctionBasedShading) BBox() *core.PdfRectangle {
	return f.bBox
}

// AntiAlias returns whether anti-aliasing is enabled.
func (f FunctionBasedShading) AntiAlias() bool {
	return f.antiAlias
}

// Functions returns the shading functions.
func (f FunctionBasedShading) Functions() []functions.PdfFunction {
	return f.functions
}

// Domain returns the rectangular domain of coordinates over which the colour function(s) are defined.
// Default value: [0.0, 1.0, 0.0, 1.0].
func (f FunctionBasedShading) Domain() []float64 {
	return f.domain
}

// Matrix returns the transformation matrix mapping the coordinate space specified by the Domain entry
// into the shading's target coordinate space. Default value: the identity matrix [1 0 0 1 0 0].
func (f FunctionBasedShading) Matrix() core.TransformationMatrix {
	return f.matrix
}

// Eval converts the input values using the functions of the shading.
func (f FunctionBasedShading) Eval(input ...float64) []float64 {
	return evalShading(f.functions, input)
}

// Equals reports whether this FunctionBasedShading equals another Shading.
func (f FunctionBasedShading) Equals(other Shading) bool {
	if other, ok := other.(FunctionBasedShading); ok {
		return f.shadingDictionary.Equals(other.shadingDictionary) &&
			f.colorSpace == other.colorSpace &&
			float64SliceEqual(f.background, other.background) &&
			pdfRectEqual(f.bBox, other.bBox) &&
			f.antiAlias == other.antiAlias &&
			float64SliceEqual(f.domain, other.domain) &&
			f.matrix.Equals(other.matrix) &&
			pdfFunctionSliceEqual(f.functions, other.functions)
	}
	return false
}

// AxialShading defines a colour blend along a line between two points, optionally
// extended beyond the boundary points by continuing the boundary colours.
type AxialShading struct {
	shadingDictionary *tokens.DictionaryToken
	colorSpace        ColorSpaceDetails
	background        []float64
	bBox              *core.PdfRectangle
	antiAlias         bool
	coords            []float64
	domain            []float64
	functions         []functions.PdfFunction
	extend            []bool
}

// NewAxialShading creates a new AxialShading.
func NewAxialShading(
	shadingDictionary *tokens.DictionaryToken,
	colorSpace ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	coords []float64,
	domain []float64,
	functions []functions.PdfFunction,
	extend []bool,
) AxialShading {
	return AxialShading{
		shadingDictionary: shadingDictionary,
		colorSpace:        colorSpace,
		background:        background,
		bBox:              bBox,
		antiAlias:         antiAlias,
		coords:            coords,
		domain:            domain,
		functions:         functions,
		extend:            extend,
	}
}

// ShadingDictionary returns the dictionary defining the shading.
func (a AxialShading) ShadingDictionary() *tokens.DictionaryToken {
	return a.shadingDictionary
}

// ShadingType returns Axial for this shading.
func (AxialShading) ShadingType() ShadingType {
	return Axial
}

// ColorSpace returns the colour space of the shading.
func (a AxialShading) ColorSpace() ColorSpaceDetails {
	return a.colorSpace
}

// Background returns the background colour value, or nil if not specified.
func (a AxialShading) Background() []float64 {
	return a.background
}

// BBox returns the shading's bounding box, or nil if not specified.
func (a AxialShading) BBox() *core.PdfRectangle {
	return a.bBox
}

// AntiAlias returns whether anti-aliasing is enabled.
func (a AxialShading) AntiAlias() bool {
	return a.antiAlias
}

// Functions returns the shading functions.
func (a AxialShading) Functions() []functions.PdfFunction {
	return a.functions
}

// Coords returns the starting and ending coordinates of the axis [x0 y0 x1 y1].
func (a AxialShading) Coords() []float64 {
	return a.coords
}

// Domain returns the limiting values of the parametric variable t [t0 t1].
// Default value: [0.0, 1.0].
func (a AxialShading) Domain() []float64 {
	return a.domain
}

// Extend returns whether to extend the shading beyond the starting and ending points.
// Default value: [false, false].
func (a AxialShading) Extend() []bool {
	return a.extend
}

// Eval converts the input values using the functions of the shading.
func (a AxialShading) Eval(input ...float64) []float64 {
	return evalShading(a.functions, input)
}

// Equals reports whether this AxialShading equals another Shading.
func (a AxialShading) Equals(other Shading) bool {
	if other, ok := other.(AxialShading); ok {
		return a.shadingDictionary.Equals(other.shadingDictionary) &&
			a.colorSpace == other.colorSpace &&
			float64SliceEqual(a.background, other.background) &&
			pdfRectEqual(a.bBox, other.bBox) &&
			a.antiAlias == other.antiAlias &&
			float64SliceEqual(a.coords, other.coords) &&
			float64SliceEqual(a.domain, other.domain) &&
			boolSliceEqual(a.extend, other.extend) &&
			pdfFunctionSliceEqual(a.functions, other.functions)
	}
	return false
}

// RadialShading defines a blend between two circles, optionally extended beyond the
// boundary circles by continuing the boundary colours.
type RadialShading struct {
	shadingDictionary *tokens.DictionaryToken
	colorSpace        ColorSpaceDetails
	background        []float64
	bBox              *core.PdfRectangle
	antiAlias         bool
	coords            []float64
	domain            []float64
	functions         []functions.PdfFunction
	extend            []bool
}

// NewRadialShading creates a new RadialShading.
func NewRadialShading(
	shadingDictionary *tokens.DictionaryToken,
	colorSpace ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	coords []float64,
	domain []float64,
	functions []functions.PdfFunction,
	extend []bool,
) RadialShading {
	return RadialShading{
		shadingDictionary: shadingDictionary,
		colorSpace:        colorSpace,
		background:        background,
		bBox:              bBox,
		antiAlias:         antiAlias,
		coords:            coords,
		domain:            domain,
		functions:         functions,
		extend:            extend,
	}
}

// ShadingDictionary returns the dictionary defining the shading.
func (r RadialShading) ShadingDictionary() *tokens.DictionaryToken {
	return r.shadingDictionary
}

// ShadingType returns Radial for this shading.
func (RadialShading) ShadingType() ShadingType {
	return Radial
}

// ColorSpace returns the colour space of the shading.
func (r RadialShading) ColorSpace() ColorSpaceDetails {
	return r.colorSpace
}

// Background returns the background colour value, or nil if not specified.
func (r RadialShading) Background() []float64 {
	return r.background
}

// BBox returns the shading's bounding box, or nil if not specified.
func (r RadialShading) BBox() *core.PdfRectangle {
	return r.bBox
}

// AntiAlias returns whether anti-aliasing is enabled.
func (r RadialShading) AntiAlias() bool {
	return r.antiAlias
}

// Functions returns the shading functions.
func (r RadialShading) Functions() []functions.PdfFunction {
	return r.functions
}

// Coords returns the centres and radii of the starting and ending circles [x0 y0 r0 x1 y1 r1].
func (r RadialShading) Coords() []float64 {
	return r.coords
}

// Domain returns the limiting values of the parametric variable t [t0 t1].
// Default value: [0.0, 1.0].
func (r RadialShading) Domain() []float64 {
	return r.domain
}

// Extend returns whether to extend the shading beyond the boundary circles.
// Default value: [false, false].
func (r RadialShading) Extend() []bool {
	return r.extend
}

// Eval converts the input values using the functions of the shading.
func (r RadialShading) Eval(input ...float64) []float64 {
	return evalShading(r.functions, input)
}

// Equals reports whether this RadialShading equals another Shading.
func (r RadialShading) Equals(other Shading) bool {
	if other, ok := other.(RadialShading); ok {
		return r.shadingDictionary.Equals(other.shadingDictionary) &&
			r.colorSpace == other.colorSpace &&
			float64SliceEqual(r.background, other.background) &&
			pdfRectEqual(r.bBox, other.bBox) &&
			r.antiAlias == other.antiAlias &&
			float64SliceEqual(r.coords, other.coords) &&
			float64SliceEqual(r.domain, other.domain) &&
			boolSliceEqual(r.extend, other.extend) &&
			pdfFunctionSliceEqual(r.functions, other.functions)
	}
	return false
}

// FreeFormGouraudShading defines a common construct used by many three-dimensional
// applications to represent complex coloured and shaded shapes. Vertices are specified
// in free-form geometry.
type FreeFormGouraudShading struct {
	shadingDictionary *tokens.DictionaryToken
	colorSpace        ColorSpaceDetails
	background        []float64
	bBox              *core.PdfRectangle
	antiAlias         bool
	bitsPerCoordinate int
	bitsPerComponent  int
	bitsPerFlag       int
	decode            []float64
	functions         []functions.PdfFunction
	data              []byte
}

// NewFreeFormGouraudShading creates a new FreeFormGouraudShading.
func NewFreeFormGouraudShading(
	shadingStream *tokens.StreamToken,
	colorSpace ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	bitsPerCoordinate int,
	bitsPerComponent int,
	bitsPerFlag int,
	decode []float64,
	functions []functions.PdfFunction,
) FreeFormGouraudShading {
	return FreeFormGouraudShading{
		shadingDictionary: shadingStream.StreamDictionary,
		colorSpace:        colorSpace,
		background:        background,
		bBox:              bBox,
		antiAlias:         antiAlias,
		bitsPerCoordinate: bitsPerCoordinate,
		bitsPerComponent:  bitsPerComponent,
		bitsPerFlag:       bitsPerFlag,
		decode:            decode,
		functions:         functions,
		data:              shadingStream.Data(),
	}
}

// ShadingDictionary returns the dictionary defining the shading.
func (f FreeFormGouraudShading) ShadingDictionary() *tokens.DictionaryToken {
	return f.shadingDictionary
}

// ShadingType returns FreeFormGouraud for this shading.
func (FreeFormGouraudShading) ShadingType() ShadingType {
	return FreeFormGouraud
}

// ColorSpace returns the colour space of the shading.
func (f FreeFormGouraudShading) ColorSpace() ColorSpaceDetails {
	return f.colorSpace
}

// Background returns the background colour value, or nil if not specified.
func (f FreeFormGouraudShading) Background() []float64 {
	return f.background
}

// BBox returns the shading's bounding box, or nil if not specified.
func (f FreeFormGouraudShading) BBox() *core.PdfRectangle {
	return f.bBox
}

// AntiAlias returns whether anti-aliasing is enabled.
func (f FreeFormGouraudShading) AntiAlias() bool {
	return f.antiAlias
}

// Functions returns the shading functions, or nil if none.
func (f FreeFormGouraudShading) Functions() []functions.PdfFunction {
	return f.functions
}

// BitsPerCoordinate returns the number of bits used to represent each vertex coordinate.
func (f FreeFormGouraudShading) BitsPerCoordinate() int {
	return f.bitsPerCoordinate
}

// BitsPerComponent returns the number of bits used to represent each colour component.
func (f FreeFormGouraudShading) BitsPerComponent() int {
	return f.bitsPerComponent
}

// BitsPerFlag returns the number of bits used to represent the edge flag for each vertex.
func (f FreeFormGouraudShading) BitsPerFlag() int {
	return f.bitsPerFlag
}

// Decode returns the decode array for mapping vertex coordinates and colour components.
func (f FreeFormGouraudShading) Decode() []float64 {
	return f.decode
}

// Data returns the decoded stream data containing descriptive data characterizing the shading.
func (f FreeFormGouraudShading) Data() []byte {
	return f.data
}

// Eval converts the input values using the functions of the shading.
func (f FreeFormGouraudShading) Eval(input ...float64) []float64 {
	return evalShading(f.functions, input)
}

// Equals reports whether this FreeFormGouraudShading equals another Shading.
func (f FreeFormGouraudShading) Equals(other Shading) bool {
	if other, ok := other.(FreeFormGouraudShading); ok {
		return f.shadingDictionary.Equals(other.shadingDictionary) &&
			f.colorSpace == other.colorSpace &&
			float64SliceEqual(f.background, other.background) &&
			pdfRectEqual(f.bBox, other.bBox) &&
			f.antiAlias == other.antiAlias &&
			f.bitsPerCoordinate == other.bitsPerCoordinate &&
			f.bitsPerComponent == other.bitsPerComponent &&
			f.bitsPerFlag == other.bitsPerFlag &&
			float64SliceEqual(f.decode, other.decode) &&
			byteSliceEqual(f.data, other.data) &&
			pdfFunctionSliceEqual(f.functions, other.functions)
	}
	return false
}

// LatticeFormGouraudShading is based on the same geometrical construct as type 4 but with
// vertices specified as a pseudorectangular lattice.
type LatticeFormGouraudShading struct {
	shadingDictionary *tokens.DictionaryToken
	colorSpace        ColorSpaceDetails
	background        []float64
	bBox              *core.PdfRectangle
	antiAlias         bool
	bitsPerCoordinate int
	bitsPerComponent  int
	verticesPerRow    int
	decode            []float64
	functions         []functions.PdfFunction
	data              []byte
}

// NewLatticeFormGouraudShading creates a new LatticeFormGouraudShading.
func NewLatticeFormGouraudShading(
	shadingStream *tokens.StreamToken,
	colorSpace ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	bitsPerCoordinate int,
	bitsPerComponent int,
	verticesPerRow int,
	decode []float64,
	functions []functions.PdfFunction,
) LatticeFormGouraudShading {
	return LatticeFormGouraudShading{
		shadingDictionary: shadingStream.StreamDictionary,
		colorSpace:        colorSpace,
		background:        background,
		bBox:              bBox,
		antiAlias:         antiAlias,
		bitsPerCoordinate: bitsPerCoordinate,
		bitsPerComponent:  bitsPerComponent,
		verticesPerRow:    verticesPerRow,
		decode:            decode,
		functions:         functions,
		data:              shadingStream.Data(),
	}
}

// ShadingDictionary returns the dictionary defining the shading.
func (l LatticeFormGouraudShading) ShadingDictionary() *tokens.DictionaryToken {
	return l.shadingDictionary
}

// ShadingType returns LatticeFormGouraud for this shading.
func (LatticeFormGouraudShading) ShadingType() ShadingType {
	return LatticeFormGouraud
}

// ColorSpace returns the colour space of the shading.
func (l LatticeFormGouraudShading) ColorSpace() ColorSpaceDetails {
	return l.colorSpace
}

// Background returns the background colour value, or nil if not specified.
func (l LatticeFormGouraudShading) Background() []float64 {
	return l.background
}

// BBox returns the shading's bounding box, or nil if not specified.
func (l LatticeFormGouraudShading) BBox() *core.PdfRectangle {
	return l.bBox
}

// AntiAlias returns whether anti-aliasing is enabled.
func (l LatticeFormGouraudShading) AntiAlias() bool {
	return l.antiAlias
}

// Functions returns the shading functions, or nil if none.
func (l LatticeFormGouraudShading) Functions() []functions.PdfFunction {
	return l.functions
}

// BitsPerCoordinate returns the number of bits used to represent each vertex coordinate.
func (l LatticeFormGouraudShading) BitsPerCoordinate() int {
	return l.bitsPerCoordinate
}

// BitsPerComponent returns the number of bits used to represent each colour component.
func (l LatticeFormGouraudShading) BitsPerComponent() int {
	return l.bitsPerComponent
}

// VerticesPerRow returns the number of vertices in each row of the lattice.
func (l LatticeFormGouraudShading) VerticesPerRow() int {
	return l.verticesPerRow
}

// Decode returns the decode array for mapping vertex coordinates and colour components.
func (l LatticeFormGouraudShading) Decode() []float64 {
	return l.decode
}

// Data returns the decoded stream data containing descriptive data characterizing the shading.
func (l LatticeFormGouraudShading) Data() []byte {
	return l.data
}

// Eval converts the input values using the functions of the shading.
func (l LatticeFormGouraudShading) Eval(input ...float64) []float64 {
	return evalShading(l.functions, input)
}

// Equals reports whether this LatticeFormGouraudShading equals another Shading.
func (l LatticeFormGouraudShading) Equals(other Shading) bool {
	if other, ok := other.(LatticeFormGouraudShading); ok {
		return l.shadingDictionary.Equals(other.shadingDictionary) &&
			l.colorSpace == other.colorSpace &&
			float64SliceEqual(l.background, other.background) &&
			pdfRectEqual(l.bBox, other.bBox) &&
			l.antiAlias == other.antiAlias &&
			l.bitsPerCoordinate == other.bitsPerCoordinate &&
			l.bitsPerComponent == other.bitsPerComponent &&
			l.verticesPerRow == other.verticesPerRow &&
			float64SliceEqual(l.decode, other.decode) &&
			byteSliceEqual(l.data, other.data) &&
			pdfFunctionSliceEqual(l.functions, other.functions)
	}
	return false
}

// CoonsPatchMeshesShading constructs a shading from one or more colour patches, each
// bounded by four cubic Bézier curves.
type CoonsPatchMeshesShading struct {
	shadingDictionary *tokens.DictionaryToken
	colorSpace        ColorSpaceDetails
	background        []float64
	bBox              *core.PdfRectangle
	antiAlias         bool
	bitsPerCoordinate int
	bitsPerComponent  int
	bitsPerFlag       int
	decode            []float64
	functions         []functions.PdfFunction
	data              []byte
}

// NewCoonsPatchMeshesShading creates a new CoonsPatchMeshesShading.
func NewCoonsPatchMeshesShading(
	shadingStream *tokens.StreamToken,
	colorSpace ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	bitsPerCoordinate int,
	bitsPerComponent int,
	bitsPerFlag int,
	decode []float64,
	functions []functions.PdfFunction,
) CoonsPatchMeshesShading {
	return CoonsPatchMeshesShading{
		shadingDictionary: shadingStream.StreamDictionary,
		colorSpace:        colorSpace,
		background:        background,
		bBox:              bBox,
		antiAlias:         antiAlias,
		bitsPerCoordinate: bitsPerCoordinate,
		bitsPerComponent:  bitsPerComponent,
		bitsPerFlag:       bitsPerFlag,
		decode:            decode,
		functions:         functions,
		data:              shadingStream.Data(),
	}
}

// ShadingDictionary returns the dictionary defining the shading.
func (c CoonsPatchMeshesShading) ShadingDictionary() *tokens.DictionaryToken {
	return c.shadingDictionary
}

// ShadingType returns CoonsPatch for this shading.
func (CoonsPatchMeshesShading) ShadingType() ShadingType {
	return CoonsPatch
}

// ColorSpace returns the colour space of the shading.
func (c CoonsPatchMeshesShading) ColorSpace() ColorSpaceDetails {
	return c.colorSpace
}

// Background returns the background colour value, or nil if not specified.
func (c CoonsPatchMeshesShading) Background() []float64 {
	return c.background
}

// BBox returns the shading's bounding box, or nil if not specified.
func (c CoonsPatchMeshesShading) BBox() *core.PdfRectangle {
	return c.bBox
}

// AntiAlias returns whether anti-aliasing is enabled.
func (c CoonsPatchMeshesShading) AntiAlias() bool {
	return c.antiAlias
}

// Functions returns the shading functions, or nil if none.
func (c CoonsPatchMeshesShading) Functions() []functions.PdfFunction {
	return c.functions
}

// BitsPerCoordinate returns the number of bits used to represent each vertex coordinate.
func (c CoonsPatchMeshesShading) BitsPerCoordinate() int {
	return c.bitsPerCoordinate
}

// BitsPerComponent returns the number of bits used to represent each colour component.
func (c CoonsPatchMeshesShading) BitsPerComponent() int {
	return c.bitsPerComponent
}

// BitsPerFlag returns the number of bits used to represent the edge flag for each patch.
func (c CoonsPatchMeshesShading) BitsPerFlag() int {
	return c.bitsPerFlag
}

// Decode returns the decode array for mapping coordinates and colour components.
func (c CoonsPatchMeshesShading) Decode() []float64 {
	return c.decode
}

// Data returns the decoded stream data containing descriptive data characterizing the shading.
func (c CoonsPatchMeshesShading) Data() []byte {
	return c.data
}

// Eval converts the input values using the functions of the shading.
func (c CoonsPatchMeshesShading) Eval(input ...float64) []float64 {
	return evalShading(c.functions, input)
}

// Equals reports whether this CoonsPatchMeshesShading equals another Shading.
func (c CoonsPatchMeshesShading) Equals(other Shading) bool {
	if other, ok := other.(CoonsPatchMeshesShading); ok {
		return c.shadingDictionary.Equals(other.shadingDictionary) &&
			c.colorSpace == other.colorSpace &&
			float64SliceEqual(c.background, other.background) &&
			pdfRectEqual(c.bBox, other.bBox) &&
			c.antiAlias == other.antiAlias &&
			c.bitsPerCoordinate == other.bitsPerCoordinate &&
			c.bitsPerComponent == other.bitsPerComponent &&
			c.bitsPerFlag == other.bitsPerFlag &&
			float64SliceEqual(c.decode, other.decode) &&
			byteSliceEqual(c.data, other.data) &&
			pdfFunctionSliceEqual(c.functions, other.functions)
	}
	return false
}

// TensorProductPatchMeshesShading is similar to type 6 but with additional control
// points in each patch, affording greater control over colour mapping.
type TensorProductPatchMeshesShading struct {
	shadingDictionary *tokens.DictionaryToken
	colorSpace        ColorSpaceDetails
	background        []float64
	bBox              *core.PdfRectangle
	antiAlias         bool
	bitsPerCoordinate int
	bitsPerComponent  int
	bitsPerFlag       int
	decode            []float64
	functions         []functions.PdfFunction
	data              []byte
}

// NewTensorProductPatchMeshesShading creates a new TensorProductPatchMeshesShading.
func NewTensorProductPatchMeshesShading(
	shadingStream *tokens.StreamToken,
	colorSpace ColorSpaceDetails,
	background []float64,
	bBox *core.PdfRectangle,
	antiAlias bool,
	bitsPerCoordinate int,
	bitsPerComponent int,
	bitsPerFlag int,
	decode []float64,
	functions []functions.PdfFunction,
) TensorProductPatchMeshesShading {
	return TensorProductPatchMeshesShading{
		shadingDictionary: shadingStream.StreamDictionary,
		colorSpace:        colorSpace,
		background:        background,
		bBox:              bBox,
		antiAlias:         antiAlias,
		bitsPerCoordinate: bitsPerCoordinate,
		bitsPerComponent:  bitsPerComponent,
		bitsPerFlag:       bitsPerFlag,
		decode:            decode,
		functions:         functions,
		data:              shadingStream.Data(),
	}
}

// ShadingDictionary returns the dictionary defining the shading.
func (t TensorProductPatchMeshesShading) ShadingDictionary() *tokens.DictionaryToken {
	return t.shadingDictionary
}

// ShadingType returns TensorProductPatch for this shading.
func (TensorProductPatchMeshesShading) ShadingType() ShadingType {
	return TensorProductPatch
}

// ColorSpace returns the colour space of the shading.
func (t TensorProductPatchMeshesShading) ColorSpace() ColorSpaceDetails {
	return t.colorSpace
}

// Background returns the background colour value, or nil if not specified.
func (t TensorProductPatchMeshesShading) Background() []float64 {
	return t.background
}

// BBox returns the shading's bounding box, or nil if not specified.
func (t TensorProductPatchMeshesShading) BBox() *core.PdfRectangle {
	return t.bBox
}

// AntiAlias returns whether anti-aliasing is enabled.
func (t TensorProductPatchMeshesShading) AntiAlias() bool {
	return t.antiAlias
}

// Functions returns the shading functions, or nil if none.
func (t TensorProductPatchMeshesShading) Functions() []functions.PdfFunction {
	return t.functions
}

// BitsPerCoordinate returns the number of bits used to represent each vertex coordinate.
func (t TensorProductPatchMeshesShading) BitsPerCoordinate() int {
	return t.bitsPerCoordinate
}

// BitsPerComponent returns the number of bits used to represent each colour component.
func (t TensorProductPatchMeshesShading) BitsPerComponent() int {
	return t.bitsPerComponent
}

// BitsPerFlag returns the number of bits used to represent the edge flag for each patch.
func (t TensorProductPatchMeshesShading) BitsPerFlag() int {
	return t.bitsPerFlag
}

// Decode returns the decode array for mapping coordinates and colour components.
func (t TensorProductPatchMeshesShading) Decode() []float64 {
	return t.decode
}

// Data returns the decoded stream data containing descriptive data characterizing the shading.
func (t TensorProductPatchMeshesShading) Data() []byte {
	return t.data
}

// Eval converts the input values using the functions of the shading.
func (t TensorProductPatchMeshesShading) Eval(input ...float64) []float64 {
	return evalShading(t.functions, input)
}

// Equals reports whether this TensorProductPatchMeshesShading equals another Shading.
func (t TensorProductPatchMeshesShading) Equals(other Shading) bool {
	if other, ok := other.(TensorProductPatchMeshesShading); ok {
		return t.shadingDictionary.Equals(other.shadingDictionary) &&
			t.colorSpace == other.colorSpace &&
			float64SliceEqual(t.background, other.background) &&
			pdfRectEqual(t.bBox, other.bBox) &&
			t.antiAlias == other.antiAlias &&
			t.bitsPerCoordinate == other.bitsPerCoordinate &&
			t.bitsPerComponent == other.bitsPerComponent &&
			t.bitsPerFlag == other.bitsPerFlag &&
			float64SliceEqual(t.decode, other.decode) &&
			byteSliceEqual(t.data, other.data) &&
			pdfFunctionSliceEqual(t.functions, other.functions)
	}
	return false
}

// evalShading evaluates input values through shading functions and clamps results to [0, 1].
func evalShading(functions []functions.PdfFunction, input []float64) []float64 {
	if len(functions) == 0 {
		return input
	}
	if len(functions) == 1 {
		return clampToUnit(functions[0].Eval(input...))
	}

	result := make([]float64, len(functions))
	for i, fn := range functions {
		newValue := fn.Eval(input...)
		result[i] = newValue[0]
	}
	return clampToUnit(result)
}

// clampToUnit adjusts values to the nearest valid value in [0, 1].
// From the PDF spec: "If the value returned by the function for a given colour component
// is out of range, it shall be adjusted to the nearest valid value."
func clampToUnit(values []float64) []float64 {
	for i := range values {
		if values[i] < 0 {
			values[i] = 0
		} else if values[i] > 1 {
			values[i] = 1
		}
	}
	return values
}

func float64SliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func boolSliceEqual(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func byteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func pdfRectEqual(a, b *core.PdfRectangle) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil) != (b == nil) {
		return false
	}
	return a.Equals(*b)
}

func pdfFunctionSliceEqual(a, b []functions.PdfFunction) bool {
	if len(a) != len(b) {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	for i := range a {
		if (a[i] == nil) != (b[i] == nil) {
			return false
		}
		if a[i] == nil && b[i] == nil {
			continue
		}
	}
	return true
}
