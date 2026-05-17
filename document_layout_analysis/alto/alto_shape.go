package alto

import "encoding/xml"

// AltoShape describes the bounding shape of a block, if it is not rectangular.
type AltoShape struct {
	// Item holds one of Circle, Ellipse, or Polygon.
	Item any `xml:"Circle,Ellipse,Polygon"`
}

// MarshalXML implements xml.Marshaler for AltoShape.
func (s *AltoShape) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Space: AltoNamespace, Local: "Shape"}

	switch v := s.Item.(type) {
	case AltoCircle:
		return enc.EncodeElement(v, start)
	case AltoEllipse:
		return enc.EncodeElement(v, start)
	case AltoPolygon:
		return enc.EncodeElement(v, start)
	default:
		return enc.EncodeElement(s, start)
	}
}
