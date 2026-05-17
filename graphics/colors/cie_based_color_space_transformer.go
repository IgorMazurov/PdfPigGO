package colors

// cieBasedColorSpaceTransformer transforms colors from CIEBased color spaces to RGB.
// The transformation implementation is based on:
// https://en.wikipedia.org/wiki/SRGB#The_forward_transformation_(CIE_XYZ_to_sRGB)
// http://www.brucelindbloom.com/index.html?Eqn_XYZ_to_RGB.html
type cieBasedColorSpaceTransformer struct {
	destinationWorkingSpace RGBWorkingSpace
	transformationMatrix    mat3
	chromaticAdaptation     chromaticAdaptation

	decoderABC func(xyzTriplet) xyzTriplet
	decoderLMN func(xyzTriplet) xyzTriplet
	matrixABC  mat3
	matrixLMN  mat3
}

// newCIEBasedColorSpaceTransformer creates a transformer from a source white point to a destination RGB working space.
func newCIEBasedColorSpaceTransformer(
	sourceReferenceWhite xyzTriplet,
	destinationWorkingSpace RGBWorkingSpace,
) cieBasedColorSpaceTransformer {
	chromaticAdaptation := newChromaticAdaptation(sourceReferenceWhite, destinationWorkingSpace.ReferenceWhite, chromaticBradford)

	xr := destinationWorkingSpace.RedPrimary.x
	yr := destinationWorkingSpace.RedPrimary.y

	xg := destinationWorkingSpace.GreenPrimary.x
	yg := destinationWorkingSpace.GreenPrimary.y

	xb := destinationWorkingSpace.BluePrimary.x
	yb := destinationWorkingSpace.BluePrimary.y

	Xr := xr / yr
	Yr := 1.0
	Zr := (1 - xr - yr) / yr

	Xg := xg / yg
	Yg := 1.0
	Zg := (1 - xg - yg) / yg

	Xb := xb / yb
	Yb := 1.0
	Zb := (1 - xb - yb) / yb

	mXYZ := newMat3(
		Xr, Xg, Xb,
		Yr, Yg, Yb,
		Zr, Zg, Zb)

	inverseMXYZ := mat3Inverse(mXYZ)

	S := mat3MultiplyVec(inverseMXYZ, destinationWorkingSpace.ReferenceWhite)

	M := newMat3(
		S.X*Xr, S.Y*Xg, S.Z*Xb,
		S.X*Yr, S.Y*Yg, S.Z*Yb,
		S.X*Zr, S.Y*Zg, S.Z*Zb)

	transformationMatrix := mat3Inverse(M)

	return cieBasedColorSpaceTransformer{
		destinationWorkingSpace: destinationWorkingSpace,
		transformationMatrix:    transformationMatrix,
		chromaticAdaptation:     chromaticAdaptation,
		decoderABC:              func(t xyzTriplet) xyzTriplet { return t },
		decoderLMN:              func(t xyzTriplet) xyzTriplet { return t },
		matrixABC:               identityMat3,
		matrixLMN:               identityMat3,
	}
}

// WithDecoderABC sets the decoder function for ABC color components.
func (t *cieBasedColorSpaceTransformer) WithDecoderABC(decoder func(xyzTriplet) xyzTriplet) {
	t.decoderABC = decoder
}

// WithDecoderLMN sets the decoder function for LMN color components.
func (t *cieBasedColorSpaceTransformer) WithDecoderLMN(decoder func(xyzTriplet) xyzTriplet) {
	t.decoderLMN = decoder
}

// WithMatrixABC sets the ABC transformation matrix.
func (t *cieBasedColorSpaceTransformer) WithMatrixABC(matrix mat3) {
	t.matrixABC = matrix
}

// WithMatrixLMN sets the LMN transformation matrix.
func (t *cieBasedColorSpaceTransformer) WithMatrixLMN(matrix mat3) {
	t.matrixLMN = matrix
}

// transformToRGB transforms an ABC color to RGB in the destination working space.
// A, B, C represent calibrated color values in the range 0 to 1.
func (t cieBasedColorSpaceTransformer) transformToRGB(color xyzTriplet) (float64, float64, float64) {
	xyz := t.transformToXYZ(color)

	adaptedColor := t.chromaticAdaptation.transform(xyz)
	rgb := mat3MultiplyVec(t.transformationMatrix, adaptedColor)

	r := t.destinationWorkingSpace.GammaCorrection(rgb.X)
	g := t.destinationWorkingSpace.GammaCorrection(rgb.Y)
	b := t.destinationWorkingSpace.GammaCorrection(rgb.Z)

	return clamp(r), clamp(g), clamp(b)
}

// transformToXYZ converts an ABC color to XYZ color space.
func (t cieBasedColorSpaceTransformer) transformToXYZ(color xyzTriplet) xyzTriplet {
	decodedABC := t.decoderABC(color)
	lmn := mat3MultiplyVec(t.matrixABC, decodedABC)

	decodedLMN := t.decoderLMN(lmn)
	xyz := mat3MultiplyVec(t.matrixLMN, decodedLMN)

	return xyz
}

// clamp forces a value into the range [0, 1].
func clamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}
