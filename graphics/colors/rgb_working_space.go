package colors

import (
	"math"
)

// xyzTriplet represents an XYZ color value with X, Y, Z components.
type xyzTriplet struct {
	X, Y, Z float64
}

// rgbPrimary represents an RGB primary chromaticity and luminance.
type rgbPrimary struct {
	x, y float64
	Y    float64
}

// RGBWorkingSpace defines the parameters of an RGB color working space.
// The RGB working space specifications were obtained from: http://www.brucelindbloom.com/index.html?WorkingSpaceInfo.html
type RGBWorkingSpace struct {
	GammaCorrection func(float64) float64
	ReferenceWhite  xyzTriplet
	RedPrimary      rgbPrimary
	GreenPrimary    rgbPrimary
	BluePrimary     rgbPrimary
}

// createGammaFunc creates a gamma correction function for the given gamma value.
func createGammaFunc(gamma float64) func(float64) float64 {
	return func(val float64) float64 {
		result := math.Pow(val, 1/gamma)
		if math.IsNaN(result) {
			return 0
		}
		return result
	}
}

// sRGBGammaCorrection implements the piecewise gamma correction for sRGB.
// Obtained from: http://www.brucelindbloom.com/index.html?Eqn_XYZ_to_RGB.html
func sRGBGammaCorrection(val float64) float64 {
	if val <= 0.0031308 {
		return 12.92 * val
	}
	return 1.055*math.Pow(val, 1/2.4) - 0.055
}

// Predefined RGB working spaces.

// AdobeRGB1998 is the Adobe RGB (1998) color space.
var AdobeRGB1998 = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D65,
	RedPrimary:      rgbPrimary{x: 0.6400, y: 0.3300, Y: 0.297361},
	GreenPrimary:    rgbPrimary{x: 0.2100, y: 0.7100, Y: 0.627355},
	BluePrimary:     rgbPrimary{x: 0.1500, y: 0.0600, Y: 0.075285},
}

// AppleRGB is the Apple RGB color space.
var AppleRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(1.8),
	ReferenceWhite:  referenceWhites.D65,
	RedPrimary:      rgbPrimary{x: 0.6250, y: 0.3400, Y: 0.244634},
	GreenPrimary:    rgbPrimary{x: 0.2800, y: 0.5950, Y: 0.672034},
	BluePrimary:     rgbPrimary{x: 0.1550, y: 0.0700, Y: 0.083332},
}

// BestRGB is the Best RGB color space.
var BestRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D50,
	RedPrimary:      rgbPrimary{x: 0.7347, y: 0.2653, Y: 0.228457},
	GreenPrimary:    rgbPrimary{x: 0.2150, y: 0.7750, Y: 0.737352},
	BluePrimary:     rgbPrimary{x: 0.1300, y: 0.0350, Y: 0.034191},
}

// BetaRGB is the Beta RGB color space.
var BetaRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D50,
	RedPrimary:      rgbPrimary{x: 0.6888, y: 0.3112, Y: 0.303273},
	GreenPrimary:    rgbPrimary{x: 0.1986, y: 0.7551, Y: 0.663786},
	BluePrimary:     rgbPrimary{x: 0.1265, y: 0.0352, Y: 0.032941},
}

// BruceRGB is the Bruce RGB color space.
var BruceRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D65,
	RedPrimary:      rgbPrimary{x: 0.6400, y: 0.3300, Y: 0.240995},
	GreenPrimary:    rgbPrimary{x: 0.2800, y: 0.6500, Y: 0.683554},
	BluePrimary:     rgbPrimary{x: 0.1500, y: 0.0600, Y: 0.075452},
}

// CIERGB is the CIE RGB color space.
var CIERGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.E,
	RedPrimary:      rgbPrimary{x: 0.7350, y: 0.2650, Y: 0.176204},
	GreenPrimary:    rgbPrimary{x: 0.2740, y: 0.7170, Y: 0.812985},
	BluePrimary:     rgbPrimary{x: 0.1670, y: 0.0090, Y: 0.010811},
}

// ColorMatchRGB is the Color Match RGB color space.
var ColorMatchRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(1.8),
	ReferenceWhite:  referenceWhites.D50,
	RedPrimary:      rgbPrimary{x: 0.6300, y: 0.3400, Y: 0.274884},
	GreenPrimary:    rgbPrimary{x: 0.2950, y: 0.6050, Y: 0.658132},
	BluePrimary:     rgbPrimary{x: 0.1500, y: 0.0750, Y: 0.066985},
}

// DonRGB4 is the Don RGB 4 color space.
var DonRGB4 = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D50,
	RedPrimary:      rgbPrimary{x: 0.6960, y: 0.3000, Y: 0.278350},
	GreenPrimary:    rgbPrimary{x: 0.2150, y: 0.7650, Y: 0.687970},
	BluePrimary:     rgbPrimary{x: 0.1300, y: 0.0350, Y: 0.033680},
}

// EktaSpacePS5 is the Ekta Space PS5 color space.
var EktaSpacePS5 = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D50,
	RedPrimary:      rgbPrimary{x: 0.6950, y: 0.3050, Y: 0.260629},
	GreenPrimary:    rgbPrimary{x: 0.2600, y: 0.7000, Y: 0.734946},
	BluePrimary:     rgbPrimary{x: 0.1100, y: 0.0050, Y: 0.004425},
}

// NTSCRGB is the NTSC RGB color space.
var NTSCRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.C,
	RedPrimary:      rgbPrimary{x: 0.6700, y: 0.3300, Y: 0.298839},
	GreenPrimary:    rgbPrimary{x: 0.2100, y: 0.7100, Y: 0.586811},
	BluePrimary:     rgbPrimary{x: 0.1400, y: 0.0800, Y: 0.114350},
}

// PALSECAMRGB is the PAL/SECAM RGB color space.
var PALSECAMRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D65,
	RedPrimary:      rgbPrimary{x: 0.6400, y: 0.3300, Y: 0.222021},
	GreenPrimary:    rgbPrimary{x: 0.2900, y: 0.6000, Y: 0.706645},
	BluePrimary:     rgbPrimary{x: 0.1500, y: 0.0600, Y: 0.071334},
}

// ProPhotoRGB is the ProPhoto RGB color space.
var ProPhotoRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(1.8),
	ReferenceWhite:  referenceWhites.D50,
	RedPrimary:      rgbPrimary{x: 0.7347, y: 0.2653, Y: 0.288040},
	GreenPrimary:    rgbPrimary{x: 0.1596, y: 0.8404, Y: 0.711874},
	BluePrimary:     rgbPrimary{x: 0.0366, y: 0.0001, Y: 0.000086},
}

// SMPTECRGB is the SMPTE C RGB color space.
var SMPTECRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D65,
	RedPrimary:      rgbPrimary{x: 0.6300, y: 0.3400, Y: 0.212395},
	GreenPrimary:    rgbPrimary{x: 0.3100, y: 0.5950, Y: 0.701049},
	BluePrimary:     rgbPrimary{x: 0.1550, y: 0.0700, Y: 0.086556},
}

// sRGB is the standard sRGB color space.
var sRGB = RGBWorkingSpace{
	GammaCorrection: sRGBGammaCorrection,
	ReferenceWhite:  referenceWhites.D65,
	RedPrimary:      rgbPrimary{x: 0.6400, y: 0.3300, Y: 0.212656},
	GreenPrimary:    rgbPrimary{x: 0.3000, y: 0.6000, Y: 0.715158},
	BluePrimary:     rgbPrimary{x: 0.1500, y: 0.0600, Y: 0.072186},
}

// WideGamutRGB is the Wide Gamut RGB color space.
var WideGamutRGB = RGBWorkingSpace{
	GammaCorrection: createGammaFunc(2.2),
	ReferenceWhite:  referenceWhites.D50,
	RedPrimary:      rgbPrimary{x: 0.7350, y: 0.2650, Y: 0.258187},
	GreenPrimary:    rgbPrimary{x: 0.1150, y: 0.8260, Y: 0.724938},
	BluePrimary:     rgbPrimary{x: 0.1570, y: 0.0180, Y: 0.016875},
}

// xyzReferenceWhite holds XYZ reference white point values for various illuminants.
// The reference white values were obtained from: http://www.brucelindbloom.com/index.html?Eqn_ChromAdapt.html
type xyzReferenceWhite struct {
	A   xyzTriplet
	B   xyzTriplet
	C   xyzTriplet
	D50 xyzTriplet
	D55 xyzTriplet
	D65 xyzTriplet
	D75 xyzTriplet
	E   xyzTriplet
	F2  xyzTriplet
	F7  xyzTriplet
	F11 xyzTriplet
}

// referenceWhites provides access to standard XYZ reference white points.
var referenceWhites = xyzReferenceWhite{
	A:   xyzTriplet{X: 1.09850, Y: 1.00000, Z: 0.35585},
	B:   xyzTriplet{X: 0.99072, Y: 1.00000, Z: 0.85223},
	C:   xyzTriplet{X: 0.98074, Y: 1.00000, Z: 1.18232},
	D50: xyzTriplet{X: 0.96422, Y: 1.00000, Z: 0.82521},
	D55: xyzTriplet{X: 0.95682, Y: 1.00000, Z: 0.92149},
	D65: xyzTriplet{X: 0.95047, Y: 1.00000, Z: 1.08883},
	D75: xyzTriplet{X: 0.94972, Y: 1.00000, Z: 1.22638},
	E:   xyzTriplet{X: 1.00000, Y: 1.00000, Z: 1.00000},
	F2:  xyzTriplet{X: 0.99186, Y: 1.00000, Z: 0.67393},
	F7:  xyzTriplet{X: 0.95041, Y: 1.00000, Z: 1.08747},
	F11: xyzTriplet{X: 1.00962, Y: 1.00000, Z: 0.64350},
}
