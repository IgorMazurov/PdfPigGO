package graphics

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/logging"
)

// GetInitialMatrix calculates the initial transformation matrix for a page,
// accounting for crop box, media box intersection, and page rotation.
func GetInitialMatrix(
	userSpaceUnit geometry.UserSpaceUnit,
	mediaBox *content.MediaBox,
	cropBox *content.CropBox,
	rotation content.PageRotationDegrees,
	log logging.Log,
) core.TransformationMatrix {
	viewBox := geometry.RectangleIntersect(mediaBox.Bounds, cropBox.Bounds)
	if viewBox == nil {
		viewBox = &cropBox.Bounds
	}

	if rotation.Value == 0 &&
		viewBox.Left() == 0 &&
		viewBox.Bottom() == 0 &&
		userSpaceUnit.PointMultiples == 1 {
		return core.Identity
	}

	t1 := core.GetTranslationMatrix(-viewBox.Left(), -viewBox.Bottom())

	if userSpaceUnit.PointMultiples != 1 {
		log.Warn("User space unit other than 1 is not implemented")
	}

	var dx, dy float64
	switch rotation.Value {
	case 0:
		return t1
	case 90:
		dx = 0
		dy = viewBox.Width
	case 180:
		dx = viewBox.Width
		dy = viewBox.Height
	case 270:
		dx = viewBox.Height
		dy = 0
	default:
		panic(fmt.Sprintf("invalid value for page rotation: %d", rotation.Value))
	}

	r := core.GetRotationMatrix(float64(-rotation.Value))

	t2 := core.GetTranslationMatrix(dx, dy)

	return t1.Multiply(r.Multiply(t2))
}
