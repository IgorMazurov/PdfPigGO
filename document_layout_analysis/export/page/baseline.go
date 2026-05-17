package page

// Baseline represents a baseline element in PAGE XML documents.
type Baseline struct {
	// Points is a string of coordinate pairs defining the baseline geometry.
	Points string `xml:"points,attr"`

	// Conf is a confidence value between 0 and 1.
	Conf *float32 `xml:"conf,attr,omitempty"`
}
