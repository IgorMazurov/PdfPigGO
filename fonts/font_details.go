package fonts

import "fmt"

// DefaultWeight is the normal weight for a font.
const DefaultWeight = 500

// BoldWeight is the bold weight for a font.
const BoldWeight = 700

// FontDetails holds summary details of the font used to draw a glyph.
type FontDetails struct {
	// Name is the font name.
	Name string

	// IsBold indicates whether the font is bold.
	IsBold bool

	// Weight is the font weight; values above 500 represent bold.
	Weight int

	// IsItalic indicates whether the font is italic.
	IsItalic bool

	bold *FontDetails
}

// NewFontDetails creates a new FontDetails instance.
func NewFontDetails(name string, isBold bool, weight int, isItalic bool) *FontDetails {
	if name == "" {
		name = ""
	}
	fd := &FontDetails{
		Name:     name,
		IsBold:   isBold,
		Weight:   weight,
		IsItalic: isItalic,
	}
	if isBold {
		fd.bold = fd
	} else {
		fd.bold = &FontDetails{
			Name:     name,
			IsBold:   true,
			Weight:   weight,
			IsItalic: isItalic,
		}
	}
	return fd
}

// AsBold returns a FontDetails with the same properties as the current instance,
// but with IsBold set to true. If already bold, returns the same instance.
func (fd *FontDetails) AsBold() *FontDetails {
	return fd.bold
}

// GetDefault returns a default FontDetails instance with the given name.
func GetDefault(name string) *FontDetails {
	if name == "" {
		name = ""
	}
	return NewFontDetails(name, false, DefaultWeight, false)
}

// WithName returns a new FontDetails with the specified name, or the same
// instance if name is empty.
func (fd *FontDetails) WithName(name string) *FontDetails {
	if name == "" {
		return fd
	}
	return NewFontDetails(name, fd.IsBold, fd.Weight, fd.IsItalic)
}

// String returns a human-readable representation of the font details.
func (fd *FontDetails) String() string {
	s := fd.Name
	if fd.IsBold {
		s += " (bold)"
	}
	if fd.IsItalic {
		s += " (italic)"
	}
	return s
}

var _ fmt.Stringer = (*FontDetails)(nil)
