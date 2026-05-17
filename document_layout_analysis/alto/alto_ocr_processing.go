package alto

// AltoOcrProcessing describes OCR processing steps in ALTO format.
type AltoOcrProcessing struct {
	PreProcessingSteps  []AltoProcessingStep `xml:"preProcessingStep"`
	OcrProcessingStep   *AltoProcessingStep  `xml:"OcrProcessingStep"`
	PostProcessingSteps []AltoProcessingStep `xml:"postProcessingStep"`
}
