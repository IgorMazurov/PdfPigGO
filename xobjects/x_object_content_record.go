package xobjects

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// XObjectContentRecord holds metadata about an XObject referenced in a page content stream.
type XObjectContentRecord struct {
	// Name is the PDF name by which this XObject is referenced in the resource dictionary.
	Name *tokens.NameToken

	// Stream is the raw stream token containing the XObject data.
	Stream *tokens.StreamToken

	// Type indicates the kind of XObject (Image, Form, PostScript).
	Type XObjectType

	// AppliedTransformation is the transformation matrix applied to this XObject.
	AppliedTransformation core.TransformationMatrix

	// DefaultRenderingIntent is the default color rendering intent for this XObject.
	DefaultRenderingIntent graphiccore.RenderingIntent

	// DefaultColorSpace is the default color space for this XObject.
	DefaultColorSpace colors.ColorSpaceDetails
}

// NewXObjectContentRecord creates a new XObjectContentRecord.
func NewXObjectContentRecord(
	name *tokens.NameToken,
	stream *tokens.StreamToken,
	typ XObjectType,
	appliedTransformation core.TransformationMatrix,
	defaultRenderingIntent graphiccore.RenderingIntent,
	defaultColorSpace colors.ColorSpaceDetails,
) *XObjectContentRecord {
	return &XObjectContentRecord{
		Name:                   name,
		Stream:                 stream,
		Type:                   typ,
		AppliedTransformation:  appliedTransformation,
		DefaultRenderingIntent: defaultRenderingIntent,
		DefaultColorSpace:      defaultColorSpace,
	}
}
