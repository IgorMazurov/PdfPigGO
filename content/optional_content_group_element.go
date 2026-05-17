package content

import (
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// OptionalContentGroupElement represents a collection of graphics that can be made
// visible or invisible dynamically by viewer applications. See ISO 32000-1:2008,
// Section 14.11.2 (Optional Content Groups).
type OptionalContentGroupElement struct {
	// Type is the PDF object type. Must be "OCG" for an optional content group.
	Type string

	// Name is the display name of the group, suitable for presentation in a viewer UI.
	Name *string

	// Intent specifies when the group should be visible. Default is ["View"].
	Intent []string

	// Usage describes the nature of the content controlled by the group.
	Usage map[string]tokens.Token

	// MarkedContent holds the underlying marked content element.
	MarkedContent *MarkedContentElement
}

// NewOptionalContentGroupElement creates a new OptionalContentGroupElement from a
// marked content element, parsing its property dictionary for OCG fields.
func NewOptionalContentGroupElement(
	markedContent *MarkedContentElement,
) (*OptionalContentGroupElement, error) {
	if markedContent == nil {
		return nil, fmt.Errorf("marked content cannot be nil")
	}

	result := &OptionalContentGroupElement{
		MarkedContent: markedContent,
	}

	props := markedContent.Properties
	if props == nil {
		return nil, fmt.Errorf("marked content properties cannot be nil")
	}

	// Type - Required
	typeToken, ok := props.TryGet(tokens.Type)
	if !ok {
		return nil, fmt.Errorf("cannot parse optional content Type from properties: required field missing")
	}

	switch t := typeToken.(type) {
	case *tokens.NameToken:
		result.Type = t.Data()
	case *tokens.StringToken:
		result.Type = t.Data()
	default:
		return nil, fmt.Errorf("cannot parse optional content Type from properties: expected NameToken or StringToken, got %T", typeToken)
	}

	switch result.Type {
	case "OCG":
		if err := parseOCG(result, props); err != nil {
			return nil, err
		}

	case "OCMD":
		if err := parseOCMD(result, props); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unknown optional content type %q", result.Type)
	}

	return result, nil
}

// parseOCG extracts Name and Intent fields for an OCG (optional content group).
func parseOCG(ocg *OptionalContentGroupElement, props *tokens.DictionaryToken) error {
	// Name - Per spec required, but some PDFs store layer names at the document catalog level.
	nameToken, ok := props.TryGet(tokens.Name)
	if ok {
		switch t := nameToken.(type) {
		case *tokens.NameToken:
			v := t.Data()
			ocg.Name = &v
		case *tokens.StringToken:
			v := t.Data()
			ocg.Name = &v
		case *tokens.HexToken:
			data := t.Data()
			ocg.Name = &data
		default:
			fallback := ocg.MarkedContent.Tag
			ocg.Name = &fallback
		}
	} else {
		fallback := ocg.MarkedContent.Tag
		ocg.Name = &fallback
	}

	// Intent - Optional, defaults to ["View"]
	intentToken, ok := props.TryGet(tokens.Intent)
	if !ok {
		ocg.Intent = []string{"View"}
		return nil
	}

	switch t := intentToken.(type) {
	case *tokens.NameToken:
		ocg.Intent = []string{t.Data()}
	case *tokens.StringToken:
		ocg.Intent = []string{t.Data()}
	case *tokens.ArrayToken:
		intentList := make([]string, 0, len(t.Data()))
		for _, elem := range t.Data() {
			switch et := elem.(type) {
			case *tokens.NameToken:
				intentList = append(intentList, et.Data())
			case *tokens.StringToken:
				intentList = append(intentList, et.Data())
			default:
				return fmt.Errorf("unsupported Intent array element type: %T", elem)
			}
		}
		ocg.Intent = intentList
	default:
		ocg.Intent = []string{"View"}
	}

	// Usage - Optional
	usageToken, ok := props.TryGet(tokens.Usage)
	if ok {
		if dict, ok := usageToken.(*tokens.DictionaryToken); ok {
			ocg.Usage = dict.Data()
		}
	}

	return nil
}

// parseOCMD handles OCMD (optional content membership dictionary).
func parseOCMD(ocg *OptionalContentGroupElement, props *tokens.DictionaryToken) error {
	// OCGs - Optional (not fully implemented)
	if ocgsToken, ok := props.TryGet(tokens.Ocgs); ok {
		return fmt.Errorf("unsupported OCMD OCGs type: %T", ocgsToken)
	}

	// P - Optional (not fully implemented)
	if pToken, ok := props.TryGet(tokens.P); ok {
		return fmt.Errorf("unsupported OCMD P type: %T", pToken)
	}

	// VE - Optional (not fully implemented)
	if veToken, ok := props.TryGet(tokens.VE); ok {
		return fmt.Errorf("unsupported OCMD VE type: %T", veToken)
	}

	return nil
}

// String returns a summary representation of the optional content group element.
func (o *OptionalContentGroupElement) String() string {
	intentStr := strings.Join(o.Intent, ",")
	name := ""
	if o.Name != nil {
		name = *o.Name
	}
	return fmt.Sprintf("%s - %s [%s]: %v", o.Type, name, intentStr, o.MarkedContent)
}
