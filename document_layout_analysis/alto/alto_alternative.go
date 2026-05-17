package alto

// AltoAlternative represents an alternative text interpretation in ALTO format.
type AltoAlternative struct {
	Purpose string `xml:"PURPOSE,attr"`
	Value   string `xml:",chardata"`
}
