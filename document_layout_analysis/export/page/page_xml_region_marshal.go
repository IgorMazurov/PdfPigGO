package page

import (
	"encoding/xml"
	"strconv"
)

// RegionElementName returns the correct XML element name for a PageXmlRegionItem.
func RegionElementName(item PageXmlRegionItem) string {
	switch item.(type) {
	case *PageXmlAdvertRegion:
		return "AdvertRegion"
	case *PageXmlChartRegion:
		return "ChartRegion"
	case *PageXmlChemRegion:
		return "ChemRegion"
	case *PageXmlCustomRegion:
		return "CustomRegion"
	case *PageXmlGraphicRegion:
		return "GraphicRegion"
	case *PageXmlImageRegion:
		return "ImageRegion"
	case *PageXmlLineDrawingRegion:
		return "LineDrawingRegion"
	case *PageXmlMapRegion:
		return "MapRegion"
	case *PageXmlMathsRegion:
		return "MathsRegion"
	case *PageXmlMusicRegion:
		return "MusicRegion"
	case *PageXmlNoiseRegion:
		return "NoiseRegion"
	case *PageXmlSeparatorRegion:
		return "SeparatorRegion"
	case *PageXmlTableRegion:
		return "TableRegion"
	case *PageXmlTextRegion:
		return "TextRegion"
	case *PageXmlUnknownRegion:
		return "UnknownRegion"
	default:
		return "UnknownRegion"
	}
}

// createRegionInstance creates a new instance of the appropriate region type based on element name.
func createRegionInstance(name string) PageXmlRegionItem {
	switch name {
	case "AdvertRegion":
		return &PageXmlAdvertRegion{}
	case "ChartRegion":
		return &PageXmlChartRegion{}
	case "ChemRegion":
		return &PageXmlChemRegion{}
	case "CustomRegion":
		return &PageXmlCustomRegion{}
	case "GraphicRegion":
		return &PageXmlGraphicRegion{}
	case "ImageRegion":
		return &PageXmlImageRegion{}
	case "LineDrawingRegion":
		return &PageXmlLineDrawingRegion{}
	case "MapRegion":
		return &PageXmlMapRegion{}
	case "MathsRegion":
		return &PageXmlMathsRegion{}
	case "MusicRegion":
		return &PageXmlMusicRegion{}
	case "NoiseRegion":
		return &PageXmlNoiseRegion{}
	case "SeparatorRegion":
		return &PageXmlSeparatorRegion{}
	case "TableRegion":
		return &PageXmlTableRegion{}
	case "TextRegion":
		return &PageXmlTextRegion{}
	case "UnknownRegion":
		return &PageXmlUnknownRegion{}
	default:
		return nil
	}
}

// MarshalXML implements xml.Marshaler for PageXmlPage.
func (p *PageXmlPage) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type pageMarshal PageXmlPage

	start.Name.Local = "Page"

	tmp := (*pageMarshal)(p)
	itemsCopy := tmp.Items
	tmp.Items = nil

	if err := e.EncodeElement(tmp, start); err != nil {
		return err
	}

	for _, item := range itemsCopy {
		name := RegionElementName(item)
		itemStart := xml.StartElement{Name: xml.Name{Local: name}}
		if m, ok := item.(xml.Marshaler); ok {
			if err := m.MarshalXML(e, itemStart); err != nil {
				return err
			}
		} else {
			if err := e.EncodeElement(item, itemStart); err != nil {
				return err
			}
		}
	}

	return nil
}

// UnmarshalXML implements xml.Unmarshaler for PageXmlPage.
func (p *PageXmlPage) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Parse attributes from start element
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "imageFilename":
			p.ImageFilename = attr.Value
		case "imageWidth":
			if v, err := parseAttrInt(attr.Value); err == nil {
				p.ImageWidth = v
			}
		case "imageHeight":
			if v, err := parseAttrInt(attr.Value); err == nil {
				p.ImageHeight = v
			}
		}
	}

	p.Items = make([]PageXmlRegionItem, 0)

	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			// Check if this is a region element
			if instance := createRegionInstance(t.Name.Local); instance != nil {
				if err := d.DecodeElement(instance, &t); err != nil {
					return err
				}
				p.Items = append(p.Items, instance)
				continue
			}

			// Handle known non-region child elements
			switch t.Name.Local {
			case "ReadingOrder":
				var ro PageXmlReadingOrder
				if err := d.DecodeElement(&ro, &t); err != nil {
					return err
				}
				p.ReadingOrder = &ro
			case "AlternativeImage", "Border", "PrintSpace", "Layers", "Relations",
				"TextStyle", "UserAttribute", "Labels":
				d.Skip()
			default:
				d.Skip()
			}

		case xml.EndElement:
			if t.Name.Local == "Page" {
				return nil
			}
		}
	}
}

// MarshalXML implements xml.Marshaler for PageXmlTextRegion.
func (r *PageXmlTextRegion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type textMarshal PageXmlTextRegion

	start.Name.Local = "TextRegion"

	tmp := (*textMarshal)(r)
	itemsCopy := tmp.Items
	tmp.Items = nil

	if err := e.EncodeElement(tmp, start); err != nil {
		return err
	}

	for _, item := range itemsCopy {
		name := RegionElementName(item)
		itemStart := xml.StartElement{Name: xml.Name{Local: name}}
		if m, ok := item.(xml.Marshaler); ok {
			if err := m.MarshalXML(e, itemStart); err != nil {
				return err
			}
		} else {
			if err := e.EncodeElement(item, itemStart); err != nil {
				return err
			}
		}
	}

	return nil
}

// UnmarshalXML implements xml.Unmarshaler for PageXmlTextRegion.
func (r *PageXmlTextRegion) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type textMarshal PageXmlTextRegion

	tmp := (*textMarshal)(r)
	if err := d.DecodeElement(tmp, &start); err != nil {
		return err
	}

	*r = *(*PageXmlTextRegion)(tmp)
	return nil
}

// MarshalXML implements xml.Marshaler for PageXmlImageRegion.
func (r *PageXmlImageRegion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type imageMarshal PageXmlImageRegion

	start.Name.Local = "ImageRegion"

	tmp := (*imageMarshal)(r)
	itemsCopy := tmp.Items
	tmp.Items = nil

	if err := e.EncodeElement(tmp, start); err != nil {
		return err
	}

	for _, item := range itemsCopy {
		name := RegionElementName(item)
		itemStart := xml.StartElement{Name: xml.Name{Local: name}}
		if m, ok := item.(xml.Marshaler); ok {
			if err := m.MarshalXML(e, itemStart); err != nil {
				return err
			}
		} else {
			if err := e.EncodeElement(item, itemStart); err != nil {
				return err
			}
		}
	}

	return nil
}

// UnmarshalXML implements xml.Unmarshaler for PageXmlImageRegion.
func (r *PageXmlImageRegion) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type imageMarshal PageXmlImageRegion

	tmp := (*imageMarshal)(r)
	if err := d.DecodeElement(tmp, &start); err != nil {
		return err
	}

	*r = *(*PageXmlImageRegion)(tmp)
	return nil
}

// MarshalXML implements xml.Marshaler for PageXmlLineDrawingRegion.
func (r *PageXmlLineDrawingRegion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type lineMarshal PageXmlLineDrawingRegion

	start.Name.Local = "LineDrawingRegion"

	tmp := (*lineMarshal)(r)
	itemsCopy := tmp.Items
	tmp.Items = nil

	if err := e.EncodeElement(tmp, start); err != nil {
		return err
	}

	for _, item := range itemsCopy {
		name := RegionElementName(item)
		itemStart := xml.StartElement{Name: xml.Name{Local: name}}
		if m, ok := item.(xml.Marshaler); ok {
			if err := m.MarshalXML(e, itemStart); err != nil {
				return err
			}
		} else {
			if err := e.EncodeElement(item, itemStart); err != nil {
				return err
			}
		}
	}

	return nil
}

// UnmarshalXML implements xml.Unmarshaler for PageXmlLineDrawingRegion.
func (r *PageXmlLineDrawingRegion) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type lineMarshal PageXmlLineDrawingRegion

	tmp := (*lineMarshal)(r)
	if err := d.DecodeElement(tmp, &start); err != nil {
		return err
	}

	*r = *(*PageXmlLineDrawingRegion)(tmp)
	return nil
}

func parseAttrInt(s string) (int, error) {
	return strconv.Atoi(s)
}
