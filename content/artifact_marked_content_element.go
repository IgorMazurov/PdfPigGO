package content

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ArtifactType defines the type of artifact marked content element.
type ArtifactType int

const (
	// Unknown represents an unknown artifact type.
	Unknown ArtifactType = iota
	// Pagination represents ancillary page features such as running heads and folios.
	Pagination
	// Layout represents purely cosmetic typographical or design elements.
	Layout
	// Page represents production aids extraneous to the document itself.
	PageArtifact
	// Background represents images, patterns or colored blocks serving as background.
	Background
)

// PdfPath represents a PDF drawing path.
// Embeds geometry.PdfPath to provide all path operations (subpaths, fill/stroke state, etc.).
type PdfPath struct {
	*geometry.PdfPath
}

// ArtifactMarkedContentElement represents graphics objects not part of the author's original
// content but generated during pagination, layout, or other mechanical processes. Artifacts may
// also describe graphical backgrounds that enhance visual experience without being required for
// understanding the document content (PDF 32000-1:2008, Section 14.8.2.2).
type ArtifactMarkedContentElement struct {
	MarkedContentElement

	// Type is the artifact type: Pagination, Layout, Page, or Background.
	Type ArtifactType
	// SubType is the artifact subtype. Standard values are Header, Footer, and Watermark.
	SubType *string
	// AttributeOwners holds the artifact's attribute owners.
	AttributeOwners *string
	// BoundingBox is the artifact's bounding box.
	BoundingBox *core.PdfRectangle
	// Attached holds the names of regions this element is attached to.
	Attached []*tokens.NameToken
}

// NewArtifactMarkedContentElement creates a new ArtifactMarkedContentElement.
func NewArtifactMarkedContentElement(
	markedContentIdentifier int,
	tag *tokens.NameToken,
	properties *tokens.DictionaryToken,
	language *string,
	actualText *string,
	alternateDescription *string,
	expandedForm *string,
	artifactType ArtifactType,
	subType *string,
	attributeOwners *string,
	boundingBox *core.PdfRectangle,
	attached []*tokens.NameToken,
	children []MarkedContentElement,
	letters []Letter,
	paths []PdfPath,
	images []PdfImage,
	index int,
) ArtifactMarkedContentElement {
	if attached == nil {
		attached = make([]*tokens.NameToken, 0)
	}

	return ArtifactMarkedContentElement{
		MarkedContentElement: MarkedContentElement{
			MarkedContentIdentifier: markedContentIdentifier,
			Index:                   index,
			Tag:                     tag.Data(),
			Properties:              properties,
			IsArtifact:              true,
			Children:                children,
			Letters:                 letters,
			Paths:                   paths,
			Images:                  images,
			Language:                language,
			ActualText:              actualText,
			AlternateDescription:    alternateDescription,
			ExpandedForm:            expandedForm,
		},
		Type:            artifactType,
		SubType:         subType,
		AttributeOwners: attributeOwners,
		BoundingBox:     boundingBox,
		Attached:        attached,
	}
}

// IsTopAttached reports whether the artifact is attached to the top edge.
func (a ArtifactMarkedContentElement) IsTopAttached() bool {
	return a.isAttached(tokens.Top)
}

// IsBottomAttached reports whether the artifact is attached to the bottom edge.
func (a ArtifactMarkedContentElement) IsBottomAttached() bool {
	return a.isAttached(tokens.Bottom)
}

// IsLeftAttached reports whether the artifact is attached to the left edge.
func (a ArtifactMarkedContentElement) IsLeftAttached() bool {
	return a.isAttached(tokens.Left)
}

// IsRightAttached reports whether the artifact is attached to the right edge.
func (a ArtifactMarkedContentElement) IsRightAttached() bool {
	return a.isAttached(tokens.Right)
}

// isAttached checks if the given edge name token is in the Attached list.
func (a ArtifactMarkedContentElement) isAttached(edge *tokens.NameToken) bool {
	for _, name := range a.Attached {
		if name == edge {
			return true
		}
	}
	return false
}
