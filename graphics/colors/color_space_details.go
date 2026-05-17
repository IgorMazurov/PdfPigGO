package colors

import (
	"fmt"
	"math"
	"sync"

	"github.com/uglytoad/pdfpig/go/functions"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ColorSpaceDetails defines the interface for color space operations.
type ColorSpaceDetails interface {
	// Type returns the color space type.
	Type() ColorSpace

	// NumberOfColorComponents returns the number of components for the color space.
	NumberOfColorComponents() int

	// BaseType returns the underlying color space type, usually equal to Type
	// unless Indexed or DeviceN.
	BaseType() ColorSpace

	// BaseNumberOfColorComponents returns the number of components for the underlying color space.
	BaseNumberOfColorComponents() int

	// GetColor creates a Color from the given component values.
	GetColor(values ...float64) Color

	// Process transforms color component values without caching or validation.
	Process(values ...float64) []float64

	// GetInitializeColor returns the initial color when this color space is activated.
	GetInitializeColor() Color

	// Transform applies color space transformation to raw image bytes.
	Transform(decoded []byte) []byte
}

// convertToByte converts a double component value [0,1] to a byte [0,255].
func convertToByte(componentValue float64) byte {
	rounded := math.Round(componentValue * 255)
	if rounded < 0 {
		rounded = 0
	}
	if rounded > 255 {
		rounded = 255
	}
	return byte(rounded)
}

// deviceGrayColorSpaceDetails implements DeviceGray color space.
type deviceGrayColorSpaceDetails struct{}

// DeviceGrayColorSpaceDetails is the singleton instance for DeviceGray.
var DeviceGrayColorSpaceDetails = ColorSpaceDetails(deviceGrayColorSpaceDetails{})

func (deviceGrayColorSpaceDetails) Type() ColorSpace {
	return DeviceGray
}

func (deviceGrayColorSpaceDetails) NumberOfColorComponents() int {
	return 1
}

func (deviceGrayColorSpaceDetails) BaseType() ColorSpace {
	return DeviceGray
}

func (deviceGrayColorSpaceDetails) BaseNumberOfColorComponents() int {
	return 1
}

func (deviceGrayColorSpaceDetails) Process(values ...float64) []float64 {
	return values
}

func (deviceGrayColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) == 0 {
		return GrayBlack
	}
	gray := values[0]
	if gray == 0 {
		return GrayBlack
	}
	if gray == 1 {
		return GrayWhite
	}
	return Color(NewGrayColor(gray))
}

func (deviceGrayColorSpaceDetails) GetInitializeColor() Color {
	return GrayBlack
}

func (deviceGrayColorSpaceDetails) Transform(decoded []byte) []byte {
	return decoded
}

// deviceRgbColorSpaceDetails implements DeviceRGB color space.
type deviceRgbColorSpaceDetails struct{}

// DeviceRgbColorSpaceDetails is the singleton instance for DeviceRGB.
var DeviceRgbColorSpaceDetails = ColorSpaceDetails(deviceRgbColorSpaceDetails{})

func (deviceRgbColorSpaceDetails) Type() ColorSpace {
	return DeviceRGB
}

func (deviceRgbColorSpaceDetails) NumberOfColorComponents() int {
	return 3
}

func (deviceRgbColorSpaceDetails) BaseType() ColorSpace {
	return DeviceRGB
}

func (deviceRgbColorSpaceDetails) BaseNumberOfColorComponents() int {
	return 3
}

func (deviceRgbColorSpaceDetails) Process(values ...float64) []float64 {
	return values
}

func (deviceRgbColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) < 3 {
		return RGBBlack
	}
	r, g, b := values[0], values[1], values[2]
	if r == 0 && g == 0 && b == 0 {
		return RGBBlack
	}
	if r == 1 && g == 1 && b == 1 {
		return RGBWhite
	}
	return Color(NewRGBColor(r, g, b))
}

func (deviceRgbColorSpaceDetails) GetInitializeColor() Color {
	return RGBBlack
}

func (deviceRgbColorSpaceDetails) Transform(decoded []byte) []byte {
	return decoded
}

// deviceCmykColorSpaceDetails implements DeviceCMYK color space.
type deviceCmykColorSpaceDetails struct{}

// DeviceCmykColorSpaceDetails is the singleton instance for DeviceCMYK.
var DeviceCmykColorSpaceDetails = ColorSpaceDetails(deviceCmykColorSpaceDetails{})

func (deviceCmykColorSpaceDetails) Type() ColorSpace {
	return DeviceCMYK
}

func (deviceCmykColorSpaceDetails) NumberOfColorComponents() int {
	return 4
}

func (deviceCmykColorSpaceDetails) BaseType() ColorSpace {
	return DeviceCMYK
}

func (deviceCmykColorSpaceDetails) BaseNumberOfColorComponents() int {
	return 4
}

func (deviceCmykColorSpaceDetails) Process(values ...float64) []float64 {
	return values
}

func (deviceCmykColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) < 4 {
		return CMYKBlack
	}
	c, m, y, k := values[0], values[1], values[2], values[3]
	if c == 0 && m == 0 && y == 0 && k == 1 {
		return CMYKBlack
	}
	if c == 0 && m == 0 && y == 0 && k == 0 {
		return CMYKWhite
	}
	return Color(NewCMYKColor(c, m, y, k))
}

func (deviceCmykColorSpaceDetails) GetInitializeColor() Color {
	return CMYKBlack
}

func (deviceCmykColorSpaceDetails) Transform(decoded []byte) []byte {
	return decoded
}

// indexedColorSpaceDetails implements Indexed color space.
type indexedColorSpaceDetails struct {
	baseColorSpace      ColorSpaceDetails
	hiVal               byte
	colorTable          []byte
	mu                  sync.RWMutex
	cache               map[float64]Color
}

// NewIndexedColorSpaceDetails creates a new Indexed color space details.
func NewIndexedColorSpaceDetails(baseColorSpace ColorSpaceDetails, hiVal byte, colorTable []byte) ColorSpaceDetails {
	return &indexedColorSpaceDetails{
		baseColorSpace: baseColorSpace,
		hiVal:          hiVal,
		colorTable:     colorTable,
		cache:          make(map[float64]Color),
	}
}

// StencilIndexedColorSpace creates an indexed color space for stencil masks with a black-and-white palette.
func StencilIndexedColorSpace(baseColorSpace ColorSpaceDetails) ColorSpaceDetails {
	return NewIndexedColorSpaceDetails(baseColorSpace, 1, []byte{0, 255})
}

func (*indexedColorSpaceDetails) Type() ColorSpace {
	return Indexed
}

func (*indexedColorSpaceDetails) NumberOfColorComponents() int {
	return 1
}

func (i *indexedColorSpaceDetails) BaseType() ColorSpace {
	return i.baseColorSpace.BaseType()
}

func (i *indexedColorSpaceDetails) BaseNumberOfColorComponents() int {
	return i.baseColorSpace.BaseNumberOfColorComponents()
}

// BaseColorSpace returns the underlying base color space for this indexed color space.
func (i *indexedColorSpaceDetails) BaseColorSpace() ColorSpaceDetails {
	return i.baseColorSpace
}

// unwrapIndexedBytes maps index values through the color table to base color space bytes.
func (i *indexedColorSpaceDetails) unwrapIndexedBytes(input []byte) []byte {
	baseType := i.baseColorSpace.Type()
	switch baseType {
	case DeviceRGB, CalRGB, Lab:
		result := make([]byte, len(input)*3)
		idx := 0
		for _, x := range input {
			for j := 0; j < 3; j++ {
				result[idx] = i.colorTable[int(x)*3+j]
				idx++
			}
		}
		return result
	case DeviceCMYK:
		result := make([]byte, len(input)*4)
		idx := 0
		for _, x := range input {
			for j := 0; j < 4; j++ {
				result[idx] = i.colorTable[int(x)*4+j]
				idx++
			}
		}
		return result
	case DeviceGray, CalGray, Separation:
		result := make([]byte, len(input))
		for idx, b := range input {
			result[idx] = i.colorTable[b]
		}
		return result
	case DeviceN, ICCBased:
		numComps := i.baseColorSpace.NumberOfColorComponents()
		if numComps == 1 {
			result := make([]byte, len(input))
			for idx, b := range input {
				result[idx] = i.colorTable[b]
			}
			return result
		}
		result := make([]byte, len(input)*numComps)
		idx := 0
		for _, x := range input {
			for j := 0; j < numComps; j++ {
				result[idx] = i.colorTable[int(x)*numComps+j]
				idx++
			}
		}
		return result
	}
	return input
}

func (i *indexedColorSpaceDetails) Process(values ...float64) []float64 {
	if len(values) == 0 {
		return nil
	}
	csBytes := i.unwrapIndexedBytes([]byte{byte(values[0])})
	scaled := make([]float64, len(csBytes))
	for idx, b := range csBytes {
		scaled[idx] = float64(b) / 255.0
	}
	return i.baseColorSpace.Process(scaled...)
}

func (i *indexedColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) == 0 {
		return GrayBlack
	}

	i.mu.RLock()
	cached, ok := i.cache[values[0]]
	i.mu.RUnlock()
	if ok {
		return cached
	}

	csBytes := i.unwrapIndexedBytes([]byte{byte(values[0])})
	scaled := make([]float64, len(csBytes))
	for idx, b := range csBytes {
		scaled[idx] = float64(b) / 255.0
	}
	color := i.baseColorSpace.GetColor(scaled...)

	i.mu.Lock()
	i.cache[values[0]] = color
	i.mu.Unlock()

	return color
}

func (i *indexedColorSpaceDetails) GetInitializeColor() Color {
	return i.GetColor(0)
}

func (i *indexedColorSpaceDetails) Transform(decoded []byte) []byte {
	unwrapped := i.unwrapIndexedBytes(decoded)
	return i.baseColorSpace.Transform(unwrapped)
}

// deviceNColorSpaceDetails implements DeviceN color space.
type deviceNColorSpaceDetails struct {
	names              []*tokens.NameToken
	alternateColorSpace ColorSpaceDetails
	tintFunction       functions.PdfFunction
	attributes         *DeviceNColorSpaceAttributes
}

// NewDeviceNColorSpaceDetails creates a new DeviceN color space details.
func NewDeviceNColorSpaceDetails(names []*tokens.NameToken, alternateColorSpace ColorSpaceDetails, tintFunction functions.PdfFunction, attributes *DeviceNColorSpaceAttributes) ColorSpaceDetails {
	return &deviceNColorSpaceDetails{
		names:              names,
		alternateColorSpace: alternateColorSpace,
		tintFunction:       tintFunction,
		attributes:         attributes,
	}
}

// DeviceNColorSpaceAttributes holds optional attributes for a DeviceN color space.
type DeviceNColorSpaceAttributes struct {
	Subtype      *tokens.NameToken
	Colorants    *tokens.DictionaryToken
	Process      *tokens.DictionaryToken
	MixingHints  *tokens.DictionaryToken
}

func (*deviceNColorSpaceDetails) Type() ColorSpace {
	return DeviceN
}

func (d *deviceNColorSpaceDetails) NumberOfColorComponents() int {
	return len(d.names)
}

func (d *deviceNColorSpaceDetails) BaseType() ColorSpace {
	return d.alternateColorSpace.BaseType()
}

func (d *deviceNColorSpaceDetails) BaseNumberOfColorComponents() int {
	return d.alternateColorSpace.NumberOfColorComponents()
}

// AlternateColorSpace returns the alternate color space for this DeviceN color space.
func (d *deviceNColorSpaceDetails) AlternateColorSpace() ColorSpaceDetails {
	return d.alternateColorSpace
}

func (d *deviceNColorSpaceDetails) Process(values ...float64) []float64 {
	evaled := d.tintFunction.Eval(values...)
	return d.alternateColorSpace.Process(evaled...)
}

func (d *deviceNColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) != len(d.names) {
		return GrayBlack
	}
	evaled := d.tintFunction.Eval(values...)
	return d.alternateColorSpace.GetColor(evaled...)
}

func (d *deviceNColorSpaceDetails) GetInitializeColor() Color {
	initValues := make([]float64, len(d.names))
	for i := range initValues {
		initValues[i] = 1.0
	}
	return d.GetColor(initValues...)
}

func (d *deviceNColorSpaceDetails) Transform(decoded []byte) []byte {
	numComps := d.NumberOfColorComponents()
	cache := make(map[int][]float64)
	var transformed []byte
	for i := 0; i < len(decoded); i += numComps {
		key := 0
		comps := make([]float64, numComps)
		for n := 0; n < numComps; n++ {
			b := decoded[i+n]
			key = key*31 ^ int(b)
			comps[n] = float64(b) / 255.0
		}
		colors, ok := cache[key]
		if !ok {
			colors = d.Process(comps...)
			cache[key] = colors
		}
		for _, c := range colors {
			transformed = append(transformed, convertToByte(c))
		}
	}
	return transformed
}

// separationColorSpaceDetails implements Separation color space.
type separationColorSpaceDetails struct {
	name              *tokens.NameToken
	alternateColorSpace ColorSpaceDetails
	tintFunction      functions.PdfFunction
	mu                sync.RWMutex
	cache             map[float64]Color
}

// NewSeparationColorSpaceDetails creates a new Separation color space details.
func NewSeparationColorSpaceDetails(name *tokens.NameToken, alternateColorSpace ColorSpaceDetails, tintFunction functions.PdfFunction) ColorSpaceDetails {
	return &separationColorSpaceDetails{
		name:              name,
		alternateColorSpace: alternateColorSpace,
		tintFunction:      tintFunction,
		cache:             make(map[float64]Color),
	}
}

func (*separationColorSpaceDetails) Type() ColorSpace {
	return Separation
}

func (*separationColorSpaceDetails) NumberOfColorComponents() int {
	return 1
}

func (s *separationColorSpaceDetails) BaseType() ColorSpace {
	return s.alternateColorSpace.BaseType()
}

func (s *separationColorSpaceDetails) BaseNumberOfColorComponents() int {
	return s.alternateColorSpace.NumberOfColorComponents()
}

// AlternateColorSpace returns the alternate color space for this Separation color space.
func (s *separationColorSpaceDetails) AlternateColorSpace() ColorSpaceDetails {
	return s.alternateColorSpace
}

func (s *separationColorSpaceDetails) Process(values ...float64) []float64 {
	if len(values) == 0 {
		return nil
	}
	evaled := s.tintFunction.Eval(values[0])
	return s.alternateColorSpace.Process(evaled...)
}

func (s *separationColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) == 0 {
		return GrayBlack
	}

	s.mu.RLock()
	cached, ok := s.cache[values[0]]
	s.mu.RUnlock()
	if ok {
		return cached
	}

	evaled := s.tintFunction.Eval(values[0])
	color := s.alternateColorSpace.GetColor(evaled...)

	s.mu.Lock()
	s.cache[values[0]] = color
	s.mu.Unlock()

	return color
}

func (s *separationColorSpaceDetails) GetInitializeColor() Color {
	return s.GetColor(1.0)
}

func (s *separationColorSpaceDetails) Transform(decoded []byte) []byte {
	cache := make(map[int][]float64, len(decoded))
	var transformed []byte
	for i := 0; i < len(decoded); i++ {
		b := decoded[i]
		key := int(b)
		colors, ok := cache[key]
		if !ok {
			colors = s.Process(float64(b) / 255.0)
			cache[key] = colors
		}
		for _, c := range colors {
			transformed = append(transformed, convertToByte(c))
		}
	}
	return transformed
}

// calGrayColorSpaceDetails implements CalGray CIE-based color space.
type calGrayColorSpaceDetails struct {
	whitePoint []float64
	blackPoint []float64
	gamma      float64
	transformer cieBasedColorSpaceTransformer
}

// NewCalGrayColorSpaceDetails creates a new CalGray color space details.
func NewCalGrayColorSpaceDetails(whitePoint []float64, blackPoint []float64, gamma float64) (ColorSpaceDetails, error) {
	if len(whitePoint) != 3 {
		return nil, fmt.Errorf("WhitePoint must consist of exactly three numbers, but was passed %d", len(whitePoint))
	}
	if blackPoint == nil || len(blackPoint) == 0 {
		blackPoint = []float64{0, 0, 0}
	}
	if len(blackPoint) != 3 {
		return nil, fmt.Errorf("BlackPoint must consist of exactly three numbers, but was passed %d", len(blackPoint))
	}
	if gamma == 0 {
		gamma = 1.0
	}

	t := newCIEBasedColorSpaceTransformer(xyzTriplet{whitePoint[0], whitePoint[1], whitePoint[2]}, sRGB)
	t.WithDecoderABC(func(c xyzTriplet) xyzTriplet {
		return xyzTriplet{
			X: math.Pow(c.X, gamma),
			Y: math.Pow(c.Y, gamma),
			Z: math.Pow(c.Z, gamma),
		}
	})
	t.WithMatrixABC(newMat3(whitePoint[0], 0, 0, 0, whitePoint[1], 0, 0, 0, whitePoint[2]))

	return &calGrayColorSpaceDetails{
		whitePoint:  whitePoint,
		blackPoint:  blackPoint,
		gamma:       gamma,
		transformer: t,
	}, nil
}

func (*calGrayColorSpaceDetails) Type() ColorSpace {
	return CalGray
}

func (*calGrayColorSpaceDetails) NumberOfColorComponents() int {
	return 1
}

func (*calGrayColorSpaceDetails) BaseType() ColorSpace {
	return CalGray
}

func (*calGrayColorSpaceDetails) BaseNumberOfColorComponents() int {
	return 1
}

func (c *calGrayColorSpaceDetails) Process(values ...float64) []float64 {
	if len(values) == 0 {
		return nil
	}
	r, _, _ := c.transformer.transformToRGB(xyzTriplet{values[0], values[0], values[0]})
	return []float64{r}
}

func (c *calGrayColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) == 0 {
		return RGBBlack
	}
	r, g, b := c.transformer.transformToRGB(xyzTriplet{values[0], values[0], values[0]})
	return Color(NewRGBColor(r, g, b))
}

func (c *calGrayColorSpaceDetails) GetInitializeColor() Color {
	return c.GetColor(0)
}

func (c *calGrayColorSpaceDetails) Transform(decoded []byte) []byte {
	transformed := make([]byte, len(decoded))
	for i := 0; i < len(decoded); i++ {
		component := float64(decoded[i]) / 255.0
		rgbPixel := c.Process(component)
		transformed[i] = convertToByte(rgbPixel[0])
	}
	return transformed
}

// calRGBColorSpaceDetails implements CalRGB CIE-based color space.
type calRGBColorSpaceDetails struct {
	whitePoint []float64
	blackPoint []float64
	gamma      []float64
	matrix     []float64
	transformer cieBasedColorSpaceTransformer
}

// NewCalRGBColorSpaceDetails creates a new CalRGB color space details.
func NewCalRGBColorSpaceDetails(whitePoint []float64, blackPoint []float64, gamma []float64, matrix []float64) (ColorSpaceDetails, error) {
	if len(whitePoint) != 3 {
		return nil, fmt.Errorf("WhitePoint must consist of exactly three numbers, but was passed %d", len(whitePoint))
	}
	if blackPoint == nil || len(blackPoint) == 0 {
		blackPoint = []float64{0, 0, 0}
	}
	if len(blackPoint) != 3 {
		return nil, fmt.Errorf("BlackPoint must consist of exactly three numbers, but was passed %d", len(blackPoint))
	}
	if gamma == nil || len(gamma) == 0 {
		gamma = []float64{1, 1, 1}
	}
	if len(gamma) != 3 {
		return nil, fmt.Errorf("Gamma must consist of exactly three numbers, but was passed %d", len(gamma))
	}
	if matrix == nil || len(matrix) == 0 {
		matrix = []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}
	}
	if len(matrix) != 9 {
		return nil, fmt.Errorf("Matrix must consist of exactly nine numbers, but was passed %d", len(matrix))
	}

	t := newCIEBasedColorSpaceTransformer(xyzTriplet{whitePoint[0], whitePoint[1], whitePoint[2]}, sRGB)
	t.WithDecoderABC(func(c xyzTriplet) xyzTriplet {
		return xyzTriplet{
			X: math.Pow(c.X, gamma[0]),
			Y: math.Pow(c.Y, gamma[1]),
			Z: math.Pow(c.Z, gamma[2]),
		}
	})
	t.WithMatrixABC(newMat3(
		matrix[0], matrix[3], matrix[6],
		matrix[1], matrix[4], matrix[7],
		matrix[2], matrix[5], matrix[8]))

	return &calRGBColorSpaceDetails{
		whitePoint:  whitePoint,
		blackPoint:  blackPoint,
		gamma:       gamma,
		matrix:      matrix,
		transformer: t,
	}, nil
}

func (*calRGBColorSpaceDetails) Type() ColorSpace {
	return CalRGB
}

func (*calRGBColorSpaceDetails) NumberOfColorComponents() int {
	return 3
}

func (*calRGBColorSpaceDetails) BaseType() ColorSpace {
	return CalRGB
}

func (*calRGBColorSpaceDetails) BaseNumberOfColorComponents() int {
	return 3
}

func (c *calRGBColorSpaceDetails) Process(values ...float64) []float64 {
	if len(values) < 3 {
		return nil
	}
	r, g, b := c.transformer.transformToRGB(xyzTriplet{values[0], values[1], values[2]})
	return []float64{r, g, b}
}

func (c *calRGBColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) < 3 {
		return RGBBlack
	}
	r, g, b := c.transformer.transformToRGB(xyzTriplet{values[0], values[1], values[2]})
	return Color(NewRGBColor(r, g, b))
}

func (*calRGBColorSpaceDetails) GetInitializeColor() Color {
	return RGBBlack
}

func (c *calRGBColorSpaceDetails) Transform(decoded []byte) []byte {
	transformed := make([]byte, len(decoded))
	idx := 0
	for i := 0; i < len(decoded); i += 3 {
		rgbPixel := c.Process(float64(decoded[i])/255.0, float64(decoded[i+1])/255.0, float64(decoded[i+2])/255.0)
		transformed[idx] = convertToByte(rgbPixel[0])
		transformed[idx+1] = convertToByte(rgbPixel[1])
		transformed[idx+2] = convertToByte(rgbPixel[2])
		idx += 3
	}
	return transformed
}

// labColorSpaceDetails implements Lab CIE-based color space.
type labColorSpaceDetails struct {
	whitePoint []float64
	blackPoint []float64
	matrix     []float64
	transformer cieBasedColorSpaceTransformer
}

// NewLabColorSpaceDetails creates a new Lab color space details.
func NewLabColorSpaceDetails(whitePoint []float64, blackPoint []float64, matrix []float64) (ColorSpaceDetails, error) {
	if len(whitePoint) != 3 {
		return nil, fmt.Errorf("WhitePoint must consist of exactly three numbers, but was passed %d", len(whitePoint))
	}
	if blackPoint == nil || len(blackPoint) == 0 {
		blackPoint = []float64{0, 0, 0}
	}
	if len(blackPoint) != 3 {
		return nil, fmt.Errorf("BlackPoint must consist of exactly three numbers, but was passed %d", len(blackPoint))
	}
	if matrix == nil || len(matrix) == 0 {
		matrix = []float64{-100, 100, -100, 100}
	}
	if len(matrix) != 4 {
		return nil, fmt.Errorf("Matrix must consist of exactly four numbers, but was passed %d", len(matrix))
	}

	t := newCIEBasedColorSpaceTransformer(xyzTriplet{whitePoint[0], whitePoint[1], whitePoint[2]}, sRGB)

	return &labColorSpaceDetails{
		whitePoint:  whitePoint,
		blackPoint:  blackPoint,
		matrix:      matrix,
		transformer: t,
	}, nil
}

func (*labColorSpaceDetails) Type() ColorSpace {
	return Lab
}

func (*labColorSpaceDetails) NumberOfColorComponents() int {
	return 3
}

func (*labColorSpaceDetails) BaseType() ColorSpace {
	return Lab
}

func (*labColorSpaceDetails) BaseNumberOfColorComponents() int {
	return 3
}

// labG is the inverse of the Lab f function.
func labG(x float64) float64 {
	if x > 6.0/29.0 {
		return x * x * x
	}
	return 108.0 / 841.0 * (x - 4.0/29.0)
}

func (l *labColorSpaceDetails) Process(values ...float64) []float64 {
	if len(values) < 3 {
		return nil
	}

	b := functions.ClipToRange(values[1], l.matrix[0], l.matrix[1])
	c := functions.ClipToRange(values[2], l.matrix[2], l.matrix[3])

	m := (values[0] + 16.0) / 116.0
	labL := m + (b / 500.0)
	n := m - (c / 200.0)

	x := l.whitePoint[0] * labG(labL)
	y := l.whitePoint[1] * labG(m)
	z := l.whitePoint[2] * labG(n)

	r, g, bOut := l.transformer.transformToRGB(xyzTriplet{x, y, z})
	return []float64{r, g, bOut}
}

func (l *labColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) < 3 {
		return RGBBlack
	}
	result := l.Process(values...)
	if len(result) < 3 {
		return RGBBlack
	}
	return Color(NewRGBColor(result[0], result[1], result[2]))
}

func (l *labColorSpaceDetails) GetInitializeColor() Color {
	b := functions.ClipToRange(0, l.matrix[0], l.matrix[1])
	c := functions.ClipToRange(0, l.matrix[2], l.matrix[3])
	return l.GetColor(0, b, c)
}

func (l *labColorSpaceDetails) Transform(decoded []byte) []byte {
	transformed := make([]byte, len(decoded))
	idx := 0
	for i := 0; i < len(decoded); i += 3 {
		rgbPixel := l.Process(float64(decoded[i])/255.0, float64(decoded[i+1])/255.0, float64(decoded[i+2])/255.0)
		transformed[idx] = convertToByte(rgbPixel[0])
		transformed[idx+1] = convertToByte(rgbPixel[1])
		transformed[idx+2] = convertToByte(rgbPixel[2])
		idx += 3
	}
	return transformed
}

// iccBasedColorSpaceDetails implements ICCBased color space.
type iccBasedColorSpaceDetails struct {
	numberOfColorComponents int
	alternateColorSpace     ColorSpaceDetails
	rangeVals               []float64
}

// NewICCBasedColorSpaceDetails creates a new ICCBased color space details.
func NewICCBasedColorSpaceDetails(numberOfColorComponents int, alternateColorSpace ColorSpaceDetails, rangeVals []float64) (ColorSpaceDetails, error) {
	if numberOfColorComponents != 1 && numberOfColorComponents != 3 && numberOfColorComponents != 4 {
		return nil, fmt.Errorf("NumberOfColorComponents must be 1, 3 or 4, got %d", numberOfColorComponents)
	}

	if alternateColorSpace == nil {
		switch numberOfColorComponents {
		case 1:
			alternateColorSpace = DeviceGrayColorSpaceDetails
		case 3:
			alternateColorSpace = DeviceRgbColorSpaceDetails
		case 4:
			alternateColorSpace = DeviceCmykColorSpaceDetails
		}
	}

	if rangeVals == nil || len(rangeVals) == 0 {
		rangeVals = make([]float64, numberOfColorComponents*2)
		for i := 0; i < numberOfColorComponents; i++ {
			rangeVals[i*2] = 0.0
			rangeVals[i*2+1] = 1.0
		}
	}
	if len(rangeVals) != numberOfColorComponents*2 {
		return nil, fmt.Errorf("Range must consist of exactly %d values (2 x NumberOfColorComponents), but was passed %d", numberOfColorComponents*2, len(rangeVals))
	}

	return &iccBasedColorSpaceDetails{
		numberOfColorComponents: numberOfColorComponents,
		alternateColorSpace:     alternateColorSpace,
		rangeVals:               rangeVals,
	}, nil
}

func (*iccBasedColorSpaceDetails) Type() ColorSpace {
	return ICCBased
}

func (i *iccBasedColorSpaceDetails) NumberOfColorComponents() int {
	return i.numberOfColorComponents
}

func (i *iccBasedColorSpaceDetails) BaseType() ColorSpace {
	return i.alternateColorSpace.BaseType()
}

func (i *iccBasedColorSpaceDetails) BaseNumberOfColorComponents() int {
	return i.numberOfColorComponents
}

func (i *iccBasedColorSpaceDetails) Process(values ...float64) []float64 {
	return i.alternateColorSpace.Process(values...)
}

func (i *iccBasedColorSpaceDetails) GetColor(values ...float64) Color {
	if len(values) != i.numberOfColorComponents {
		return GrayBlack
	}
	for c := 0; c < len(values); c++ {
		idx := 2 * c
		values[c] = functions.ClipToRange(values[c], i.rangeVals[idx], i.rangeVals[idx+1])
	}
	return i.alternateColorSpace.GetColor(values...)
}

func (i *iccBasedColorSpaceDetails) GetInitializeColor() Color {
	v := functions.ClipToRange(0.0, i.rangeVals[0], i.rangeVals[1])
	init := make([]float64, i.numberOfColorComponents)
	for c := range init {
		init[c] = v
	}
	return i.GetColor(init...)
}

func (i *iccBasedColorSpaceDetails) Transform(decoded []byte) []byte {
	return i.alternateColorSpace.Transform(decoded)
}

// unsupportedColorSpaceDetails represents an unsupported color space.
type unsupportedColorSpaceDetails struct{}

// UnsupportedColorSpaceDetails is the singleton instance for unsupported color spaces.
var UnsupportedColorSpaceDetails = ColorSpaceDetails(unsupportedColorSpaceDetails{})

func (unsupportedColorSpaceDetails) Type() ColorSpace {
	return DeviceGray
}

func (unsupportedColorSpaceDetails) NumberOfColorComponents() int {
	return 1
}

func (unsupportedColorSpaceDetails) BaseType() ColorSpace {
	return DeviceGray
}

func (unsupportedColorSpaceDetails) BaseNumberOfColorComponents() int {
	return 1
}

func (unsupportedColorSpaceDetails) Process(values ...float64) []float64 {
	return values
}

func (unsupportedColorSpaceDetails) GetColor(values ...float64) Color {
	return GrayBlack
}

func (unsupportedColorSpaceDetails) GetInitializeColor() Color {
	return GrayBlack
}

func (unsupportedColorSpaceDetails) Transform(decoded []byte) []byte {
	return decoded
}

// patternColorSpaceDetails implements Pattern color space.
type patternColorSpaceDetails struct {
	patterns             map[*tokens.NameToken]Color
	underlyingColorSpace ColorSpaceDetails
}

// NewPatternColorSpaceDetails creates a new Pattern color space details.
func NewPatternColorSpaceDetails(patterns map[*tokens.NameToken]Color, underlyingColorSpace ColorSpaceDetails) ColorSpaceDetails {
	return &patternColorSpaceDetails{
		patterns:             patterns,
		underlyingColorSpace: underlyingColorSpace,
	}
}

func (*patternColorSpaceDetails) Type() ColorSpace {
	return Pattern
}

func (*patternColorSpaceDetails) NumberOfColorComponents() int {
	return -1
}

func (p *patternColorSpaceDetails) BaseType() ColorSpace {
	if p.underlyingColorSpace != nil {
		return p.underlyingColorSpace.BaseType()
	}
	return Pattern
}

func (p *patternColorSpaceDetails) BaseNumberOfColorComponents() int {
	if p.underlyingColorSpace != nil {
		return p.underlyingColorSpace.NumberOfColorComponents()
	}
	return -1
}

func (*patternColorSpaceDetails) Process(values ...float64) []float64 {
	return nil
}

// GetColorByName looks up a pattern by its NameToken, matching C# PatternColorSpaceDetails.GetColor(NameToken).
func (p *patternColorSpaceDetails) GetColorByName(name *tokens.NameToken) Color {
	if p.patterns == nil || name == nil {
		return GrayBlack
	}
	if c, ok := p.patterns[name]; ok {
		return c
	}
	for _, c := range p.patterns {
		return c
	}
	return GrayBlack
}

func (p *patternColorSpaceDetails) GetColor(values ...float64) Color {
	if p.patterns == nil || len(p.patterns) == 0 {
		return GrayBlack
	}
	for _, c := range p.patterns {
		return c
	}
	return GrayBlack
}

func (*patternColorSpaceDetails) GetInitializeColor() Color {
	return nil
}

func (p *patternColorSpaceDetails) Transform(decoded []byte) []byte {
	if p.underlyingColorSpace != nil {
		return p.underlyingColorSpace.Transform(decoded)
	}
	return decoded
}

// IsUnsupported checks if the given ColorSpaceDetails is an unsupported color space.
func IsUnsupported(details ColorSpaceDetails) bool {
	_, ok := details.(unsupportedColorSpaceDetails)
	return ok
}

// IndexedColorSpaceAccessor provides access to the base color space of an indexed color space.
type IndexedColorSpaceAccessor interface {
	ColorSpaceDetails
	BaseColorSpace() ColorSpaceDetails
}

// DeviceNColorSpaceAccessor provides access to the alternate color space of a DeviceN color space.
type DeviceNColorSpaceAccessor interface {
	ColorSpaceDetails
	AlternateColorSpace() ColorSpaceDetails
}

// SeparationColorSpaceAccessor provides access to the alternate color space of a Separation color space.
type SeparationColorSpaceAccessor interface {
	ColorSpaceDetails
	AlternateColorSpace() ColorSpaceDetails
}
