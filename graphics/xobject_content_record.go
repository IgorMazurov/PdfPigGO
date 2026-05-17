package graphics

import (
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// XObjectContentRecord holds metadata about an XObject referenced in a page content stream.
type XObjectContentRecord struct {
	Type                    xobjects.XObjectType
	Stream                  *tokens.StreamToken
	AppliedTransformation   core.TransformationMatrix
	DefaultRenderingIntent  graphiccore.RenderingIntent
	DefaultColorSpace       colors.ColorSpaceDetails
}

// NewXObjectContentRecord creates a new XObjectContentRecord.
func NewXObjectContentRecord(
	type_ xobjects.XObjectType,
	stream *tokens.StreamToken,
	appliedTransformation core.TransformationMatrix,
	defaultRenderingIntent graphiccore.RenderingIntent,
	defaultColorSpace colors.ColorSpaceDetails,
) *XObjectContentRecord {
	if stream == nil {
		panic("stream cannot be nil")
	}

	return &XObjectContentRecord{
		Type:                   type_,
		Stream:                 stream,
		AppliedTransformation:  appliedTransformation,
		DefaultRenderingIntent: defaultRenderingIntent,
		DefaultColorSpace:      defaultColorSpace,
	}
}
