package alto

// AltoProcessingStep describes a single processing step in ALTO format.
type AltoProcessingStep struct {
	ProcessingCategory      *AltoProcessingCategory `xml:"processingCategory,attr"`
	ProcessingDateTime      string                  `xml:"processingDateTime,attr"`
	ProcessingAgency        string                  `xml:"processingAgency,attr"`
	ProcessingStepDescription []string              `xml:"processingStepDescription"`
	ProcessingStepSettings  string                  `xml:"processingStepSettings,attr"`
	ProcessingSoftware      *AltoProcessingSoftware `xml:"processingSoftware"`
}
