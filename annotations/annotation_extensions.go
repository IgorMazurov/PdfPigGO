package annotations

import (
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ToAnnotationType converts a NameToken to its corresponding AnnotationType.
func ToAnnotationType(name *tokens.NameToken) AnnotationType {
	if name == nil {
		return Other
	}

	switch name {
	case tokens.Text:
		return Text
	case tokens.Link:
		return Link
	case tokens.FreeText:
		return FreeText
	case tokens.Line:
		return Line
	case tokens.Square:
		return Square
	case tokens.Circle:
		return Circle
	case tokens.Polygon:
		return Polygon
	case tokens.PolyLine:
		return PolyLine
	case tokens.Highlight:
		return Highlight
	case tokens.Underline:
		return Underline
	case tokens.Squiggly:
		return Squiggly
	case tokens.StrikeOut:
		return StrikeOut
	case tokens.Stamp:
		return Stamp
	case tokens.Caret:
		return Caret
	case tokens.Ink:
		return Ink
	case tokens.Popup:
		return Popup
	case tokens.FileAttachment:
		return FileAttachment
	case tokens.Sound:
		return Sound
	case tokens.Movie:
		return Movie
	case tokens.Widget:
		return Widget
	case tokens.Screen:
		return Screen
	case tokens.PrinterMark:
		return PrinterMark
	case tokens.TrapNet:
		return TrapNet
	case tokens.Watermark:
		return Watermark
	case tokens.Annotation3D:
		return Artwork3D
	default:
		return Other
	}
}
