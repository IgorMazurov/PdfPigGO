package alto

// AltoTag represents a tag element in ALTO format used to classify and group information.
type AltoTag struct {
	XmlData *AltoTagXmlData `xml:"XmlData"`

	Id          string `xml:"ID,attr"`
	Type        string `xml:"TYPE,attr"`
	Label       string `xml:"LABEL,attr"`
	Description string `xml:"DESCRIPTION,attr"`
	Uri         string `xml:"URI,attr"`
}

// AltoTagXmlData contains XML encoded metadata within a tag element.
type AltoTagXmlData struct {
	Any []byte `xml:",any"`
}
