package alto

// AltoMeasurementUnit represents the measurement unit used in ALTO documents.
type AltoMeasurementUnit string

func (u AltoMeasurementUnit) String() string { return string(u) }

const (
	AltoMeasurementUnitPixel    AltoMeasurementUnit = "pixel"
	AltoMeasurementUnitMm10     AltoMeasurementUnit = "mm10"
	AltoMeasurementUnitInch1200 AltoMeasurementUnit = "inch1200"
)
