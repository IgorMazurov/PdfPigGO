package fonts

// AdobeStylePrivateDictionary holds common properties between Adobe Type 1 and
// Compact Font Format private dictionaries.
type AdobeStylePrivateDictionary struct {
	// BlueValues is an array containing an even number of integers. The first pair
	// is the baseline overshoot position and the baseline. All following pairs
	// describe top-zones. Required.
	BlueValues []int

	// OtherBlues are pairs of integers similar to BlueValues describing bottom zones. Optional.
	OtherBlues []int

	// FamilyBlues are integer pairs similar to BlueValues used to enforce consistency
	// across a font family when there are small differences (<1px) in font alignment. Optional.
	FamilyBlues []int

	// FamilyOtherBlues are integer pairs similar to OtherBlues used to enforce consistency
	// across a font family with small differences in alignment similarly to FamilyBlues. Optional.
	FamilyOtherBlues []int

	// BlueScale is the point size at which overshoot suppression stops. The value is related
	// to the number of pixels tall that one character space unit will be before overshoot
	// suppression is switched off. Default: 0.039625.
	BlueScale float64

	// BlueShift is the character space distance beyond the flat position of alignment zones
	// at which overshoot enforcement occurs. Default: 7.
	BlueShift int

	// BlueFuzz is the number of character space units to extend an alignment zone on a
	// horizontal stem. If the top or bottom of a horizontal stem is within BlueFuzz units
	// outside a top-zone then the stem top/bottom is treated as if it were within the zone.
	// Default: 1.
	BlueFuzz int

	// StandardHorizontalWidth is the dominant width of horizontal stems vertically in character space units. Optional.
	StandardHorizontalWidth *float64

	// StandardVerticalWidth is the dominant width of vertical stems horizontally in character space units. Optional.
	StandardVerticalWidth *float64

	// StemSnapHorizontalWidths are up to 12 numbers with the most common widths for horizontal stems vertically. Optional.
	StemSnapHorizontalWidths []float64

	// StemSnapVerticalWidths are up to 12 numbers with the most common widths for vertical stems horizontally. Optional.
	StemSnapVerticalWidths []float64

	// ForceBold controls whether bold characters should appear thicker using special techniques
	// at small sizes at low resolutions. Optional.
	ForceBold bool

	// LanguageGroup determines the script family: 0 includes Latin, Greek and Cyrillic;
	// 1 includes Chinese, Japanese Kanji and Korean Hangul. Default: 0.
	LanguageGroup int

	// ExpansionFactor is the limit for changing the size of a character bounding box for
	// LanguageGroup 1 counters during font processing. Optional.
	ExpansionFactor float64
}

// DefaultBlueScale is the default value of BlueScale (0.039625).
const DefaultBlueScale = 0.039625

// DefaultExpansionFactor is the default value of ExpansionFactor (0.06).
const DefaultExpansionFactor = 0.06

// DefaultBlueFuzz is the default value of BlueFuzz (1).
const DefaultBlueFuzz = 1

// DefaultBlueShift is the default value of BlueShift (7).
const DefaultBlueShift = 7

// DefaultLanguageGroup is the default value of LanguageGroup (0).
const DefaultLanguageGroup = 0

// AdobeStylePrivateDictionaryBuilder is a mutable builder for constructing an
// AdobeStylePrivateDictionary. It performs no validation.
type AdobeStylePrivateDictionaryBuilder struct {
	BlueValues               []int
	OtherBlues               []int
	FamilyBlues              []int
	FamilyOtherBlues         []int
	BlueScale                *float64
	BlueShift                *int
	BlueFuzz                 *int
	StandardHorizontalWidth  *float64
	StandardVerticalWidth    *float64
	StemSnapHorizontalWidths []float64
	StemSnapVerticalWidths   []float64
	ForceBold                *bool
	LanguageGroup            *int
	ExpansionFactor          *float64
}

// NewAdobeStylePrivateDictionary creates a new AdobeStylePrivateDictionary from the given builder.
func NewAdobeStylePrivateDictionary(builder *AdobeStylePrivateDictionaryBuilder) AdobeStylePrivateDictionary {
	if builder == nil {
		builder = &AdobeStylePrivateDictionaryBuilder{}
	}

	return AdobeStylePrivateDictionary{
		BlueValues:               builder.BlueValues,
		OtherBlues:               builder.OtherBlues,
		FamilyBlues:              builder.FamilyBlues,
		FamilyOtherBlues:         builder.FamilyOtherBlues,
		BlueScale:                valueOrDefault(builder.BlueScale, DefaultBlueScale),
		BlueShift:                intValOrDefault(builder.BlueShift, DefaultBlueShift),
		BlueFuzz:                 intValOrDefault(builder.BlueFuzz, DefaultBlueFuzz),
		StandardHorizontalWidth:  builder.StandardHorizontalWidth,
		StandardVerticalWidth:    builder.StandardVerticalWidth,
		StemSnapHorizontalWidths: builder.StemSnapHorizontalWidths,
		StemSnapVerticalWidths:   builder.StemSnapVerticalWidths,
		ForceBold:                boolValOrDefault(builder.ForceBold, false),
		LanguageGroup:            intValOrDefault(builder.LanguageGroup, DefaultLanguageGroup),
		ExpansionFactor:          valueOrDefault(builder.ExpansionFactor, DefaultExpansionFactor),
	}
}

func valueOrDefault(v *float64, def float64) float64 {
	if v == nil {
		return def
	}
	return *v
}

func intValOrDefault(v *int, def int) int {
	if v == nil {
		return def
	}
	return *v
}

func boolValOrDefault(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}
