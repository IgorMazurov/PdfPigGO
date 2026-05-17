package alto

// AltoProcessingSoftware describes a software tool used in ALTO processing.
type AltoProcessingSoftware struct {
	SoftwareCreator        string `xml:"softwareCreator,attr"`
	SoftwareName           string `xml:"softwareName,attr"`
	SoftwareVersion        string `xml:"softwareVersion,attr"`
	ApplicationDescription string `xml:"applicationDescription,attr"`
}
