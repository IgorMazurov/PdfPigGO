package colors

import (
	"errors"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// GetFamily returns the ColorSpaceFamily for the given ColorSpace.
func GetFamily(cs ColorSpace) ColorSpaceFamily {
	switch cs {
	case DeviceGray, DeviceRGB, DeviceCMYK:
		return Device
	case CalGray, CalRGB, Lab, ICCBased:
		return CIEBased
	case Indexed, Pattern, Separation, DeviceN:
		return Special
	default:
		return 0
	}
}

// TryMapToColorSpace maps a NameToken to the corresponding ColorSpace.
// It returns true if the name token matches a known color space name.
func TryMapToColorSpace(name *tokens.NameToken) (ColorSpace, bool) {
	if name == nil {
		return 0, false
	}

	data := name.Data()

	switch data {
	case tokens.Devicegray.Data(), "G":
		return DeviceGray, true
	case tokens.Devicergb.Data(), "RGB":
		return DeviceRGB, true
	case tokens.Devicecmyk.Data(), "CMYK":
		return DeviceCMYK, true
	case tokens.Calgray.Data():
		return CalGray, true
	case tokens.Calrgb.Data():
		return CalRGB, true
	case tokens.Lab.Data():
		return Lab, true
	case tokens.Iccbased.Data():
		return ICCBased, true
	case tokens.Indexed.Data(), "I":
		return Indexed, true
	case tokens.Pattern.Data():
		return Pattern, true
	case tokens.Separation.Data():
		return Separation, true
	case tokens.Devicen.Data():
		return DeviceN, true
	default:
		return 0, false
	}
}

// ErrUnknownColorSpace is returned when an unrecognized ColorSpace value is encountered.
var ErrUnknownColorSpace = errors.New("unrecognized colorspace")

// ToNameToken returns the NameToken corresponding to the given ColorSpace.
func ToNameToken(cs ColorSpace) (*tokens.NameToken, error) {
	switch cs {
	case DeviceGray:
		return tokens.Devicegray, nil
	case DeviceRGB:
		return tokens.Devicergb, nil
	case DeviceCMYK:
		return tokens.Devicecmyk, nil
	case CalGray:
		return tokens.Calgray, nil
	case CalRGB:
		return tokens.Calrgb, nil
	case Lab:
		return tokens.Lab, nil
	case ICCBased:
		return tokens.Iccbased, nil
	case Indexed:
		return tokens.Indexed, nil
	case Pattern:
		return tokens.Pattern, nil
	case Separation:
		return tokens.Separation, nil
	case DeviceN:
		return tokens.Devicen, nil
	default:
		return nil, ErrUnknownColorSpace
	}
}
