package colors

// chromaticAdaptationMethod defines the algorithm for chromatic adaptation.
type chromaticAdaptationMethod int

const (
	chromaticXYZScaling chromaticAdaptationMethod = iota
	chromaticBradford
	chromaticVonKries
)

// chromaticAdaptation encapsulates the algorithm for chromatic adaptation.
// Reference: http://www.brucelindbloom.com/index.html?Eqn_ChromAdapt.html
type chromaticAdaptation struct {
	adaptationMatrix mat3
}

// newChromaticAdaptation creates a chromatic adapter from source to destination white point.
func newChromaticAdaptation(
	sourceReferenceWhite xyzTriplet,
	destinationReferenceWhite xyzTriplet,
	method chromaticAdaptationMethod,
) chromaticAdaptation {
	coneResponseDomain := getConeResponseDomain(method)
	inverseConeResponseDomain := mat3Inverse(coneResponseDomain)

	srcCone := mat3MultiplyVec(coneResponseDomain, sourceReferenceWhite)
	dstCone := mat3MultiplyVec(coneResponseDomain, destinationReferenceWhite)

	rhoS, gammaS, betaS := srcCone.X, srcCone.Y, srcCone.Z
	rhoD, gammaD, betaD := dstCone.X, dstCone.Y, dstCone.Z

	scale := newMat3(
		rhoD/rhoS, 0, 0,
		0, gammaD/gammaS, 0,
		0, 0, betaD/betaS,
	)

	adaptationMatrix := mat3Multiply(mat3Multiply(inverseConeResponseDomain, scale), coneResponseDomain)

	return chromaticAdaptation{
		adaptationMatrix: adaptationMatrix,
	}
}

// transform adapts the source color from one white point to another.
func (ca chromaticAdaptation) transform(sourceColor xyzTriplet) xyzTriplet {
	return mat3MultiplyVec(ca.adaptationMatrix, sourceColor)
}

// getConeResponseDomain returns the cone response domain matrix for the given method.
func getConeResponseDomain(method chromaticAdaptationMethod) mat3 {
	switch method {
	case chromaticXYZScaling:
		return newMat3(
			1.0000000, 0.0000000, 0.0000000,
			0.0000000, 1.0000000, 0.0000000,
			0.0000000, 0.0000000, 1.0000000)

	case chromaticBradford:
		return newMat3(
			0.8951000, 0.2664000, -0.1614000,
			-0.7502000, 1.7135000, 0.0367000,
			0.0389000, -0.0685000, 1.0296000)

	case chromaticVonKries:
		return newMat3(
			0.4002400, 0.7076000, -0.0808100,
			-0.2263000, 1.1653200, 0.0457000,
			0.0000000, 0.0000000, 0.9182200)

	default:
		return getConeResponseDomain(chromaticBradford)
	}
}
