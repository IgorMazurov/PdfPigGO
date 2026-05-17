package alto

// AltoDescriptionOcrProcessing describes OCR processing steps in ALTO format.
// Deprecated: use AltoProcessing instead.
type AltoDescriptionOcrProcessing struct {
	AltoOcrProcessing
	Id string `xml:"Id,attr"`
}
