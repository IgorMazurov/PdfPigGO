package alto

// AltoQuality represents the quality status of an ALTO element.
type AltoQuality string

const (
	AltoQualityOK               AltoQuality = "OK"
	AltoQualityMissing           AltoQuality = "Missing"
	AltoQualityMissingInOriginal AltoQuality = "Missing in original"
	AltoQualityDamaged           AltoQuality = "Damaged"
	AltoQualityRetained          AltoQuality = "Retained"
	AltoQualityTarget            AltoQuality = "Target"
	AltoQualityAsInOriginal      AltoQuality = "As in original"
)

// String returns the XML string representation of this quality.
func (q AltoQuality) String() string {
	return string(q)
}
