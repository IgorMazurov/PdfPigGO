package content

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// MarkedContentElement represents a marked content segment in a PDF content stream.
// Marked content can provide application-specific data; interpretation is outside
// the PDF specification. See ISO 32000-1:2008, Section 14.7.
type MarkedContentElement struct {
	// MarkedContentIdentifier is the marked-content identifier (tag MCID).
	MarkedContentIdentifier int

	// Index is the position of this element in the page's marked content set.
	// Child elements share the same index as their parent.
	Index int

	// Tag indicates the role or significance of the marked content point.
	Tag string

	// Properties holds the property dictionary for this element.
	Properties *tokens.DictionaryToken

	// IsArtifact reports whether this is an artifact marked content element.
	IsArtifact bool

	// Children are nested marked content elements.
	Children []MarkedContentElement

	// Letters are the text letters contained in this marked content.
	Letters []Letter

	// Paths are the drawing paths contained in this marked content.
	Paths []PdfPath

	// Images are the images contained in this marked content.
	Images []PdfImage

	// Language is the natural language specification (e.g., "en-US").
	Language *string

	// ActualText is the replacement text for accessibility tools.
	ActualText *string

	// AlternateDescription provides an alternate description of the content.
	AlternateDescription *string

	// ExpandedForm is the abbreviation expansion text.
	ExpandedForm *string

	// Artifact-specific fields (only meaningful when IsArtifact is true).
	// ArtifactType indicates the type of artifact.
	ArtifactType ArtifactType
	// SubType is the artifact subtype (e.g., "Header", "Footer").
	SubType *string
	// BoundingBox is the artifact's bounding box.
	BoundingBox *core.PdfRectangle
	// IsTopAttached reports whether attached to top edge.
	IsTopAttached bool
	// IsBottomAttached reports whether attached to bottom edge.
	IsBottomAttached bool
	// IsLeftAttached reports whether attached to left edge.
	IsLeftAttached bool
	// IsRightAttached reports whether attached to right edge.
	IsRightAttached bool
}

// NewMarkedContentElement creates a new MarkedContentElement with validation.
func NewMarkedContentElement(
	markedContentIdentifier int,
	tag *tokens.NameToken,
	properties *tokens.DictionaryToken,
	language *string,
	actualText *string,
	alternateDescription *string,
	expandedForm *string,
	isArtifact bool,
	children []MarkedContentElement,
	letters []Letter,
	paths []PdfPath,
	images []PdfImage,
	index int,
) (*MarkedContentElement, error) {
	if tag == nil {
		return nil, fmt.Errorf("tag cannot be nil")
	}

	if properties == nil {
		var err error
		properties, err = tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
		if err != nil {
			return nil, fmt.Errorf("failed to create empty properties: %w", err)
		}
	}

	if children == nil {
		return nil, fmt.Errorf("children cannot be nil")
	}
	if letters == nil {
		return nil, fmt.Errorf("letters cannot be nil")
	}
	if paths == nil {
		return nil, fmt.Errorf("paths cannot be nil")
	}
	if images == nil {
		return nil, fmt.Errorf("images cannot be nil")
	}

	return &MarkedContentElement{
		MarkedContentIdentifier: markedContentIdentifier,
		Index:                   index,
		Tag:                     tag.Data(),
		Language:                language,
		ActualText:              actualText,
		AlternateDescription:    alternateDescription,
		ExpandedForm:            expandedForm,
		Properties:              properties,
		IsArtifact:              isArtifact,
		Children:                children,
		Letters:                 letters,
		Paths:                   paths,
		Images:                  images,
	}, nil
}

// String returns a summary representation of the marked content element.
func (m *MarkedContentElement) String() string {
	return fmt.Sprintf("Id=%d, MCID=%d, Tag=%s, Properties=%v, Children=%d",
		m.Index, m.MarkedContentIdentifier, m.Tag, m.Properties, len(m.Children))
}
