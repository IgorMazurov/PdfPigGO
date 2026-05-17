package alto

// AltoDescription represents the Description element in ALTO format.
type AltoDescription struct {
	MeasurementUnit   AltoMeasurementUnit          `xml:"MeasurementUnit"`
	SourceImageInfo   *AltoSourceImageInformation  `xml:"sourceImageInformation"`
	OcrProcessings    []AltoDescriptionOcrProcessing `xml:"OCRProcessing"`
	Processings       []AltoDescriptionProcessing  `xml:"Processing"`
}
