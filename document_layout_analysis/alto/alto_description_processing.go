package alto

// AltoDescriptionProcessing describes general processing steps in ALTO format.
type AltoDescriptionProcessing struct {
	AltoProcessingStep
	Id string `xml:"ID,attr"`
}
