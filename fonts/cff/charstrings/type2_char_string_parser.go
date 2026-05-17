package charstrings

import (
	"fmt"
	"math"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/cff"
	"github.com/uglytoad/pdfpig/go/fonts/cff_charset"
)

// Type2CharStringParser decodes the commands and numbers making up a Type 2 CharString.
// A Type 2 charstring program is a sequence of unsigned 8-bit bytes that encode numbers and operators.
const (
	hstemByte    = byte(1)
	vstemByte    = byte(3)
	hstemhmByte  = byte(18)
	hintmaskByte = byte(19)
	cntrmaskByte = byte(20)
	vstemhmByte  = byte(23)
)

var hintingCommandBytes = map[byte]bool{
	hstemByte:   true,
	vstemByte:   true,
	hstemhmByte: true,
	vstemhmByte: true,
}

var singleByteCommandStore = map[byte]*LazyType2Command{
	hstemByte: NewLazyType2Command("hstem", 2, func(ctx *Type2BuildCharContext) {
		numHints := ctx.Stack().Length() / 2
		hints := make([]core.PdfRange, numHints)

		firstStartY, _ := ctx.Stack().PopBottom()
		width, _ := ctx.Stack().PopBottom()
		endY := firstStartY + width

		hints[0] = core.NewPdfRange([]float64{firstStartY, endY})

		currentY := endY

		for i := 1; i < numHints; i++ {
			dyStart, _ := ctx.Stack().PopBottom()
			dyEnd, _ := ctx.Stack().PopBottom()

			start := currentY + dyStart
			end := start + dyEnd
			hints[i] = core.NewPdfRange([]float64{start, end})
			currentY = end
		}

		ctx.AddHorizontalStemHints(hints)
		ctx.Stack().Clear()
	}),

	vstemByte: NewLazyType2Command("vstem", 2, func(ctx *Type2BuildCharContext) {
		numHints := ctx.Stack().Length() / 2
		hints := make([]core.PdfRange, numHints)

		firstStartX, _ := ctx.Stack().PopBottom()
		width, _ := ctx.Stack().PopBottom()
		endX := firstStartX + width

		hints[0] = core.NewPdfRange([]float64{firstStartX, endX})

		currentX := endX

		for i := 1; i < numHints; i++ {
			dxStart, _ := ctx.Stack().PopBottom()
			dxEnd, _ := ctx.Stack().PopBottom()

			start := currentX + dxStart
			end := start + dxEnd
			hints[i] = core.NewPdfRange([]float64{start, end})
			currentX = end
		}

		ctx.AddVerticalStemHints(hints)
		ctx.Stack().Clear()
	}),

	4: NewLazyType2Command("vmoveto", 1, func(ctx *Type2BuildCharContext) {
		dy, _ := ctx.Stack().PopBottom()
		ctx.AddVerticalMoveTo(dy)
		ctx.Stack().Clear()
	}),

	5: NewLazyType2Command("rlineto", 2, func(ctx *Type2BuildCharContext) {
		numLines := ctx.Stack().Length() / 2
		for i := 0; i < numLines; i++ {
			dxa, _ := ctx.Stack().PopBottom()
			dya, _ := ctx.Stack().PopBottom()
			addRelativeLineToCtx(ctx, dxa, dya)
		}
		ctx.Stack().Clear()
	}),

	6: NewLazyType2Command("hlineto", 1, func(ctx *Type2BuildCharContext) {
		isOdd := ctx.Stack().Length()%2 != 0
		numAdditional := ctx.Stack().Length() - ternaryInt(isOdd, 1, 0)

		if isOdd {
			dx1, _ := ctx.Stack().PopBottom()
			ctx.AddRelativeHorizontalLine(dx1)

			for i := 0; i < numAdditional; i += 2 {
				dy, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeVerticalLine(dy)
				dx, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeHorizontalLine(dx)
			}
		} else {
			for i := 0; i < numAdditional; i += 2 {
				dx, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeHorizontalLine(dx)
				dy, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeVerticalLine(dy)
			}
		}

		ctx.Stack().Clear()
	}),

	7: NewLazyType2Command("vlineto", 1, func(ctx *Type2BuildCharContext) {
		isOdd := ctx.Stack().Length()%2 != 0
		numAdditional := ctx.Stack().Length() - ternaryInt(isOdd, 1, 0)

		if isOdd {
			dy1, _ := ctx.Stack().PopBottom()
			ctx.AddRelativeVerticalLine(dy1)

			for i := 0; i < numAdditional; i += 2 {
				dx, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeHorizontalLine(dx)
				dy, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeVerticalLine(dy)
			}
		} else {
			for i := 0; i < numAdditional; i += 2 {
				dy, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeVerticalLine(dy)
				dx, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeHorizontalLine(dx)
			}
		}

		ctx.Stack().Clear()
	}),

	8: NewLazyType2Command("rrcurveto", 6, func(ctx *Type2BuildCharContext) {
		curveCount := ctx.Stack().Length() / 6
		for i := 0; i < curveCount; i++ {
			dx1, _ := ctx.Stack().PopBottom()
			dy1, _ := ctx.Stack().PopBottom()
			dx2, _ := ctx.Stack().PopBottom()
			dy2, _ := ctx.Stack().PopBottom()
			dx3, _ := ctx.Stack().PopBottom()
			dy3, _ := ctx.Stack().PopBottom()
			ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, dy3)
		}
		ctx.Stack().Clear()
	}),

	10: NewLazyType2Command("callsubr", 1, func(ctx *Type2BuildCharContext) {}),
	11: NewLazyType2Command("return", 0, func(ctx *Type2BuildCharContext) {}),
	14: NewLazyType2Command("endchar", 0, func(ctx *Type2BuildCharContext) {
		ctx.Stack().Clear()
	}),

	hstemhmByte: NewLazyType2Command("hstemhm", 2, func(ctx *Type2BuildCharContext) {
		numHints := ctx.Stack().Length() / 2
		hints := make([]core.PdfRange, numHints)

		firstStartY, _ := ctx.Stack().PopBottom()
		width, _ := ctx.Stack().PopBottom()
		endY := firstStartY + width

		hints[0] = core.NewPdfRange([]float64{firstStartY, endY})

		currentY := endY

		for i := 1; i < numHints; i++ {
			dyStart, _ := ctx.Stack().PopBottom()
			dyEnd, _ := ctx.Stack().PopBottom()

			start := currentY + dyStart
			end := start + dyEnd
			hints[i] = core.NewPdfRange([]float64{start, end})
			currentY = end
		}

		ctx.AddHorizontalStemHints(hints)
		ctx.Stack().Clear()
	}),

	hintmaskByte: NewLazyType2Command("hintmask", 0, func(ctx *Type2BuildCharContext) {
		ctx.Stack().Clear()
	}),

	cntrmaskByte: NewLazyType2Command("cntrmask", 0, func(ctx *Type2BuildCharContext) {
		ctx.Stack().Clear()
	}),

	21: NewLazyType2Command("rmoveto", 2, func(ctx *Type2BuildCharContext) {
		dx, _ := ctx.Stack().PopBottom()
		dy, _ := ctx.Stack().PopBottom()
		ctx.AddRelativeMoveTo(dx, dy)
		ctx.Stack().Clear()
	}),

	22: NewLazyType2Command("hmoveto", 1, func(ctx *Type2BuildCharContext) {
		dx, _ := ctx.Stack().PopBottom()
		ctx.AddHorizontalMoveTo(dx)
		ctx.Stack().Clear()
	}),

	vstemhmByte: NewLazyType2Command("vstemhm", 2, func(ctx *Type2BuildCharContext) {
		numHints := ctx.Stack().Length() / 2
		hints := make([]core.PdfRange, numHints)

		firstStartX, _ := ctx.Stack().PopBottom()
		width, _ := ctx.Stack().PopBottom()
		endX := firstStartX + width

		hints[0] = core.NewPdfRange([]float64{firstStartX, endX})

		currentX := endX

		for i := 1; i < numHints; i++ {
			dxStart, _ := ctx.Stack().PopBottom()
			dxEnd, _ := ctx.Stack().PopBottom()

			start := currentX + dxStart
			end := start + dxEnd
			hints[i] = core.NewPdfRange([]float64{start, end})
			currentX = end
		}

		ctx.AddVerticalStemHints(hints)
		ctx.Stack().Clear()
	}),

	24: NewLazyType2Command("rcurveline", 8, func(ctx *Type2BuildCharContext) {
		numCurves := (ctx.Stack().Length() - 2) / 6
		for i := 0; i < numCurves; i++ {
			dx1, _ := ctx.Stack().PopBottom()
			dy1, _ := ctx.Stack().PopBottom()
			dx2, _ := ctx.Stack().PopBottom()
			dy2, _ := ctx.Stack().PopBottom()
			dx3, _ := ctx.Stack().PopBottom()
			dy3, _ := ctx.Stack().PopBottom()
			ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, dy3)
		}

		dxa, _ := ctx.Stack().PopBottom()
		dya, _ := ctx.Stack().PopBottom()
		addRelativeLineToCtx(ctx, dxa, dya)
		ctx.Stack().Clear()
	}),

	25: NewLazyType2Command("rlinecurve", 8, func(ctx *Type2BuildCharContext) {
		numLines := (ctx.Stack().Length() - 6) / 2
		for i := 0; i < numLines; i++ {
			dxa, _ := ctx.Stack().PopBottom()
			dya, _ := ctx.Stack().PopBottom()
			addRelativeLineToCtx(ctx, dxa, dya)
		}

		dx1, _ := ctx.Stack().PopBottom()
		dy1, _ := ctx.Stack().PopBottom()
		dx2, _ := ctx.Stack().PopBottom()
		dy2, _ := ctx.Stack().PopBottom()
		dx3, _ := ctx.Stack().PopBottom()
		dy3, _ := ctx.Stack().PopBottom()
		ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, dy3)

		ctx.Stack().Clear()
	}),

	26: NewLazyType2Command("vvcurveto", 4, func(ctx *Type2BuildCharContext) {
		hasDxFirst := ctx.Stack().Length()%4 != 0
		numCurves := ctx.Stack().Length() / 4

		for i := 0; i < numCurves; i++ {
			dx1 := 0.0
			if i == 0 && hasDxFirst {
				dx1, _ = ctx.Stack().PopBottom()
			}

			dy1, _ := ctx.Stack().PopBottom()
			dx2, _ := ctx.Stack().PopBottom()
			dy2, _ := ctx.Stack().PopBottom()
			dy3, _ := ctx.Stack().PopBottom()

			ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, 0, dy3)
		}

		ctx.Stack().Clear()
	}),

	27: NewLazyType2Command("hhcurveto", 4, func(ctx *Type2BuildCharContext) {
		hasDyFirst := ctx.Stack().Length()%4 != 0

		if hasDyFirst {
			dy1, _ := ctx.Stack().PopBottom()
			dx1, _ := ctx.Stack().PopBottom()
			dx2, _ := ctx.Stack().PopBottom()
			dy2, _ := ctx.Stack().PopBottom()
			dx3, _ := ctx.Stack().PopBottom()

			ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, 0)
		}

		numCurves := ctx.Stack().Length() / 4
		for i := 0; i < numCurves; i++ {
			dx1, _ := ctx.Stack().PopBottom()
			dx2, _ := ctx.Stack().PopBottom()
			dy2, _ := ctx.Stack().PopBottom()
			dx3, _ := ctx.Stack().PopBottom()

			ctx.AddRelativeBezierCurve(dx1, 0, dx2, dy2, dx3, 0)
		}

		ctx.Stack().Clear()
	}),

	29: NewLazyType2Command("callgsubr", 1, func(ctx *Type2BuildCharContext) {}),

	30: NewLazyType2Command("vhcurveto", 4, func(ctx *Type2BuildCharContext) {
		remainder := ctx.Stack().Length() % 8

		if remainder <= 1 {
			numCurves := (ctx.Stack().Length() - remainder) / 8
			for i := 0; i < numCurves; i++ {
				dy1, _ := ctx.Stack().PopBottom()
				dx2, _ := ctx.Stack().PopBottom()
				dy2, _ := ctx.Stack().PopBottom()
				dx3, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeBezierCurve(0, dy1, dx2, dy2, dx3, 0)

				sdx1, _ := ctx.Stack().PopBottom()
				sdx2, _ := ctx.Stack().PopBottom()
				sdy2, _ := ctx.Stack().PopBottom()
				sdy3, _ := ctx.Stack().PopBottom()
				var sdx3 float64 = 0
				if i == numCurves-1 && remainder == 1 {
					sdx3, _ = ctx.Stack().PopBottom()
				}
				ctx.AddRelativeBezierCurve(sdx1, 0, sdx2, sdy2, sdx3, sdy3)
			}
		} else if remainder == 4 || remainder == 5 {
			numCurves := (ctx.Stack().Length() - remainder) / 8

			dy1, _ := ctx.Stack().PopBottom()
			dx2, _ := ctx.Stack().PopBottom()
			dy2, _ := ctx.Stack().PopBottom()
			dx3, _ := ctx.Stack().PopBottom()
			var dy3 float64 = 0
			if ctx.Stack().Length() == 1 {
				dy3, _ = ctx.Stack().PopBottom()
			}
			ctx.AddRelativeBezierCurve(0, dy1, dx2, dy2, dx3, dy3)

			for i := 0; i < numCurves; i++ {
				dx1, _ := ctx.Stack().PopBottom()
				dxa, _ := ctx.Stack().PopBottom()
				dya, _ := ctx.Stack().PopBottom()
				dyb, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeBezierCurve(dx1, 0, dxa, dya, 0, dyb)

				dy1, _ := ctx.Stack().PopBottom()
				dx2, _ := ctx.Stack().PopBottom()
				dy2, _ := ctx.Stack().PopBottom()
				dx3, _ := ctx.Stack().PopBottom()
				var dy3b float64 = 0
				if i == numCurves-1 && remainder == 5 {
					dy3b, _ = ctx.Stack().PopBottom()
				}
				ctx.AddRelativeBezierCurve(0, dy1, dx2, dy2, dx3, dy3b)
			}
		} else {
			panic(fmt.Sprintf("Unexpected number of arguments for vhcurveto: %d", ctx.Stack().Length()))
		}

		ctx.Stack().Clear()
	}),

	31: NewLazyType2Command("hvcurveto", 4, func(ctx *Type2BuildCharContext) {
		remainder := ctx.Stack().Length() % 8

		if remainder <= 1 {
			numCurves := (ctx.Stack().Length() - remainder) / 8
			for i := 0; i < numCurves; i++ {
				dx1, _ := ctx.Stack().PopBottom()
				dxa, _ := ctx.Stack().PopBottom()
				dya, _ := ctx.Stack().PopBottom()
				dyb, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeBezierCurve(dx1, 0, dxa, dya, 0, dyb)

				dy1, _ := ctx.Stack().PopBottom()
				dx2, _ := ctx.Stack().PopBottom()
				dy2, _ := ctx.Stack().PopBottom()
				dx3, _ := ctx.Stack().PopBottom()
				var dy3 float64 = 0
				if i == numCurves-1 && remainder == 1 {
					dy3, _ = ctx.Stack().PopBottom()
				}
				ctx.AddRelativeBezierCurve(0, dy1, dx2, dy2, dx3, dy3)
			}
		} else if remainder == 4 || remainder == 5 {
			numCurves := (ctx.Stack().Length() - remainder) / 8

			dx1, _ := ctx.Stack().PopBottom()
			dxa, _ := ctx.Stack().PopBottom()
			dya, _ := ctx.Stack().PopBottom()
			dyb, _ := ctx.Stack().PopBottom()
			var dx3 float64 = 0
			if ctx.Stack().Length() == 1 {
				dx3, _ = ctx.Stack().PopBottom()
			}
			ctx.AddRelativeBezierCurve(dx1, 0, dxa, dya, dx3, dyb)

			for i := 0; i < numCurves; i++ {
				dy1, _ := ctx.Stack().PopBottom()
				dx2, _ := ctx.Stack().PopBottom()
				dy2, _ := ctx.Stack().PopBottom()
				dx3, _ := ctx.Stack().PopBottom()
				ctx.AddRelativeBezierCurve(0, dy1, dx2, dy2, dx3, 0)

				dx1, _ := ctx.Stack().PopBottom()
				dxa, _ := ctx.Stack().PopBottom()
				dya, _ := ctx.Stack().PopBottom()
				dyb, _ := ctx.Stack().PopBottom()
				var dx3b float64 = 0
				if i == numCurves-1 && remainder == 5 {
					dx3b, _ = ctx.Stack().PopBottom()
				}
				ctx.AddRelativeBezierCurve(dx1, 0, dxa, dya, dx3b, dyb)
			}
		} else {
			panic(fmt.Sprintf("Unexpected number of arguments for hvcurveto: %d", ctx.Stack().Length()))
		}

		ctx.Stack().Clear()
	}),

	255: NewLazyType2Command("unknown", -1, func(ctx *Type2BuildCharContext) {}),
}

var twoByteCommandStore = map[byte]*LazyType2Command{
	3: NewLazyType2Command("and", 2, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		b, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(ternaryBool(a != 0 && b != 0, 1, 0))
	}),

	4: NewLazyType2Command("or", 2, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		b, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(ternaryBool(a != 0 || b != 0, 1, 0))
	}),

	5: NewLazyType2Command("not", 1, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(ternaryBool(a == 0, 1, 0))
	}),

	9: NewLazyType2Command("abs", 1, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(math.Abs(a))
	}),

	10: NewLazyType2Command("add", 2, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		b, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(a + b)
	}),

	11: NewLazyType2Command("sub", 2, func(ctx *Type2BuildCharContext) {
		num1, _ := ctx.Stack().PopTop()
		num2, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(num2 - num1)
	}),

	12: NewLazyType2Command("div", 2, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		b, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(a / b)
	}),

	14: NewLazyType2Command("neg", 1, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(-1 * math.Abs(a))
	}),

	15: NewLazyType2Command("eq", 2, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		b, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(ternaryBool(a == b, 1, 0))
	}),

	18: NewLazyType2Command("drop", 1, func(ctx *Type2BuildCharContext) {
		ctx.Stack().PopTop()
	}),

	20: NewLazyType2Command("put", 2, func(ctx *Type2BuildCharContext) {
		val, _ := ctx.Stack().PopTop()
		idx, _ := ctx.Stack().PopTop()
		ctx.AddToTransientArray(val, int(idx))
	}),

	21: NewLazyType2Command("get", 1, func(ctx *Type2BuildCharContext) {
		idx, _ := ctx.Stack().PopTop()
		val, _ := ctx.GetFromTransientArray(int(idx))
		ctx.Stack().Push(val)
	}),

	22: NewLazyType2Command("ifelse", 4, func(ctx *Type2BuildCharContext) {}),

	23: NewLazyType2Command("random", 0, func(ctx *Type2BuildCharContext) {
		ctx.Stack().Push(0.5)
	}),

	24: NewLazyType2Command("mul", 2, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		b, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(a * b)
	}),

	26: NewLazyType2Command("sqrt", 1, func(ctx *Type2BuildCharContext) {
		a, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(math.Sqrt(a))
	}),

	27: NewLazyType2Command("dup", 1, func(ctx *Type2BuildCharContext) {
		val, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(val)
		ctx.Stack().Push(val)
	}),

	28: NewLazyType2Command("exch", 2, func(ctx *Type2BuildCharContext) {
		num1, _ := ctx.Stack().PopTop()
		num2, _ := ctx.Stack().PopTop()
		ctx.Stack().Push(num1)
		ctx.Stack().Push(num2)
	}),

	29: NewLazyType2Command("index", 2, func(ctx *Type2BuildCharContext) {
		idx, _ := ctx.Stack().PopTop()
		val := ctx.Stack().CopyElementAt(int(idx))
		ctx.Stack().Push(val)
	}),

	30: NewLazyType2Command("roll", 3, func(ctx *Type2BuildCharContext) {
	}),

	34: NewLazyType2Command("hflex", 7, func(ctx *Type2BuildCharContext) {
		dx1, _ := ctx.Stack().PopBottom()
		dx2, _ := ctx.Stack().PopBottom()
		dy2, _ := ctx.Stack().PopBottom()
		dx3, _ := ctx.Stack().PopBottom()
		dx4, _ := ctx.Stack().PopBottom()
		dx5, _ := ctx.Stack().PopBottom()
		dx6, _ := ctx.Stack().PopBottom()

		ctx.AddRelativeBezierCurve(dx1, 0, dx2, dy2, dx3, 0)
		ctx.AddRelativeBezierCurve(dx4, 0, dx5, -dy2, dx6, 0)

		ctx.Stack().Clear()
	}),

	35: NewLazyType2Command("flex", 13, func(ctx *Type2BuildCharContext) {
		dx1, _ := ctx.Stack().PopBottom()
		dy1, _ := ctx.Stack().PopBottom()
		dx2, _ := ctx.Stack().PopBottom()
		dy2, _ := ctx.Stack().PopBottom()
		dx3, _ := ctx.Stack().PopBottom()
		dy3, _ := ctx.Stack().PopBottom()

		ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, dy3)

		dx4, _ := ctx.Stack().PopBottom()
		dy4, _ := ctx.Stack().PopBottom()
		dx5, _ := ctx.Stack().PopBottom()
		dy5, _ := ctx.Stack().PopBottom()
		dx6, _ := ctx.Stack().PopBottom()
		dy6, _ := ctx.Stack().PopBottom()

		ctx.AddRelativeBezierCurve(dx4, dy4, dx5, dy5, dx6, dy6)

		ctx.Stack().PopBottom()
		ctx.Stack().Clear()
	}),

	36: NewLazyType2Command("hflex1", 9, func(ctx *Type2BuildCharContext) {
		dx1, _ := ctx.Stack().PopBottom()
		dy1, _ := ctx.Stack().PopBottom()
		dx2, _ := ctx.Stack().PopBottom()
		dy2, _ := ctx.Stack().PopBottom()
		dx3, _ := ctx.Stack().PopBottom()
		dx4, _ := ctx.Stack().PopBottom()
		dx5, _ := ctx.Stack().PopBottom()
		dy5, _ := ctx.Stack().PopBottom()
		dx6, _ := ctx.Stack().PopBottom()

		ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, 0)
		ctx.AddRelativeBezierCurve(dx4, 0, dx5, dy5, dx6, -dy5-dy2-dy1)

		ctx.Stack().Clear()
	}),

	37: NewLazyType2Command("flex1", 11, func(ctx *Type2BuildCharContext) {
		dx1, _ := ctx.Stack().PopBottom()
		dy1, _ := ctx.Stack().PopBottom()
		dx2, _ := ctx.Stack().PopBottom()
		dy2, _ := ctx.Stack().PopBottom()
		dx3, _ := ctx.Stack().PopBottom()
		dy3, _ := ctx.Stack().PopBottom()

		dx4, _ := ctx.Stack().PopBottom()
		dy4, _ := ctx.Stack().PopBottom()
		dx5, _ := ctx.Stack().PopBottom()
		dy5, _ := ctx.Stack().PopBottom()
		d6, _ := ctx.Stack().PopBottom()

		dxSum := dx1 + dx2 + dx3 + dx4 + dx5
		dySum := dy1 + dy2 + dy3 + dy4 + dy5

		lastPointIsX := math.Abs(dxSum) > math.Abs(dySum)

		ctx.AddRelativeBezierCurve(dx1, dy1, dx2, dy2, dx3, dy3)
		if lastPointIsX {
			ctx.AddRelativeBezierCurve(dx4, dy4, dx5, dy5, d6, 0)
		} else {
			ctx.AddRelativeBezierCurve(dx4, dy4, dx5, dy5, 0, d6)
		}
		ctx.Stack().Clear()
	}),
}

// addRelativeLineToCtx is a package-level helper that calls the unexported
// addRelativeLine method on Type2BuildCharContext. Used by command handlers
// that need to push both dx and dy onto the path.
func addRelativeLineToCtx(ctx *Type2BuildCharContext, dx, dy float64) {
	dest := core.NewPdfPoint(ctx.CurrentLocation().X+dx, ctx.CurrentLocation().Y+dy)
	ctx.Path()[len(ctx.Path())-1].LineTo(dest.X, dest.Y)
	ctx.SetCurrentLocation(dest)
}

// GetCommand returns the command definition for the given identifier.
func GetCommand(identifier CommandIdentifier) *LazyType2Command {
	if identifier.IsMultiByteCommand {
		return twoByteCommandStore[identifier.CommandId]
	}

	return singleByteCommandStore[identifier.CommandId]
}

// Parse decodes all Type 2 CharStrings from the given byte data using the provided
// subroutine selector and charset. Returns a Type2CharStrings containing the parsed
// command sequences keyed by glyph name.
func Parse(charStringBytes [][]byte, subroutinesSelector *CompactFontFormatSubroutinesSelector, charset cffcharset.CompactFontFormatCharset) (*Type2CharStrings, error) {
	if charStringBytes == nil {
		return nil, fmt.Errorf("charStringBytes cannot be nil")
	}
	if subroutinesSelector == nil {
		return nil, fmt.Errorf("subroutinesSelector cannot be nil")
	}

	charStrings := make(map[string]CommandSequence)
	for i := range charStringBytes {
		cs := charStringBytes[i]
		name := charset.GetNameByGlyphId(i)
		globalSubs, localSubs := subroutinesSelector.GetSubroutines(i)
		bytesCopy := make([]byte, len(cs))
		copy(bytesCopy, cs)
		
		seq := parseSingle(bytesCopy, localSubs, globalSubs)
		charStrings[name] = seq
	}

	return &Type2CharStrings{charStrings: charStrings, glyphs: make(map[string]*Type2Glyph)}, nil
}

func parseSingle(bytes []byte, localSubroutines, globalSubroutines *cff.CompactFontFormatIndex) CommandSequence {
	values := make([]float64, 0)
	commandIdentifiers := make([]CommandIdentifier, 0)

	for i := 0; i < len(bytes); i++ {
		b := bytes[i]
		if b <= 31 && b != 28 {
			var cmd *CommandIdentifier
			cmd, values = getCommand(b, &bytes, values, commandIdentifiers, localSubroutines, globalSubroutines, &i)
			if cmd != nil {
				commandIdentifiers = append(commandIdentifiers, *cmd)
			}
		} else {
			num := interpretNumber(b, bytes, &i)
			values = append(values, num)
		}
	}

	return CommandSequence{
		Values:             values,
		CommandIdentifiers: commandIdentifiers,
	}
}

// interpretNumber decodes a numeric operand from the charstring byte stream.
// Byte 28 encodes a 16-bit two's complement number. Bytes 32-246 encode single-byte
// numbers in range -107..1131. Bytes 247-254 encode two-byte positive/negative numbers.
// Byte 255 encodes a four-byte fixed-point number with 16 bits of fraction.
func interpretNumber(b byte, bytes []byte, i *int) float64 {
	switch {
	case b == 28:
		if *i+2 >= len(bytes) {
			return 0
		}
		*i++
		num := int(bytes[*i])<<8 | int(bytes[*i+1])
		*i++
		return float64(int16(num))

	case b >= 32 && b <= 246:
		return float64(int(b) - 139)

	case b >= 247 && b <= 250:
		if *i+1 >= len(bytes) {
			return 0
		}
		*i++
		w := bytes[*i]
		return float64((int(b)-247)*256) + float64(w) + 108

	case b >= 251 && b <= 254:
		if *i+1 >= len(bytes) {
			return 0
		}
		*i++
		w := bytes[*i]
		return -float64((int(b)-251)*256) - float64(w) - 108

	default:
		// Byte 255: four-byte fixed-point number (16.16 format)
		if *i+4 >= len(bytes) {
			return 0
		}
		*i++
		lead := int(int16(bytes[*i])<<8) + int(bytes[*i+1])
		*i += 2
		fractionalPart := int(bytes[*i])<<8 | int(bytes[*i+1])
		*i++
		return float64(lead) + float64(fractionalPart)/65535.0
	}
}

func getCommand(b byte, bytes *[]byte, precedingValues []float64, precedingCommands []CommandIdentifier, localSubroutines, globalSubroutines *cff.CompactFontFormatIndex, i *int) (*CommandIdentifier, []float64) {
	const returnCommand = byte(11)

	if b == 12 {
		if *i+1 >= len(*bytes) {
			return &CommandIdentifier{
				CommandIndex:       len(precedingValues),
				IsMultiByteCommand: false,
				CommandId:          255,
			}, precedingValues
		}
		*i++
		b2 := (*bytes)[*i]
		if _, ok := twoByteCommandStore[b2]; ok {
			return &CommandIdentifier{
				CommandIndex:       len(precedingValues),
				IsMultiByteCommand: true,
				CommandId:          b2,
			}, precedingValues
		}
		return &CommandIdentifier{
			CommandIndex:       len(precedingValues),
			IsMultiByteCommand: false,
			CommandId:          255,
		}, precedingValues
	}

	if b == 10 || b == 29 {
		isLocal := b == 10
		if len(precedingValues) == 0 {
			return nil, precedingValues
		}
		precedingNumber := int(precedingValues[len(precedingValues)-1])

		subCount := localSubroutines.Count()
		if !isLocal {
			subCount = globalSubroutines.Count()
		}
		bias := CountToBias(subCount)
		index := precedingNumber + bias
		var subroutineBytes []byte
		if isLocal {
			subroutineBytes = localSubroutines.Get(index)
		} else {
			subroutineBytes = globalSubroutines.Get(index)
		}

		oldI := *i
		*bytes = append((*bytes)[:oldI-1], append(subroutineBytes, (*bytes)[oldI+1:]...)...)

		precedingValues = precedingValues[:len(precedingValues)-1]
		*i -= 2
		return nil, precedingValues
	}

	if b == 19 || b == 20 {
		minFullBytes := calculatePrecedingHintBytes(precedingValues, precedingCommands)
		*i += minFullBytes
	}

	if _, ok := singleByteCommandStore[b]; ok {
		if b == returnCommand {
			return nil, precedingValues
		}
		return &CommandIdentifier{
			CommandIndex:       len(precedingValues),
			IsMultiByteCommand: false,
			CommandId:          b,
		}, precedingValues
	}

	return &CommandIdentifier{
		CommandIndex:       len(precedingValues),
		IsMultiByteCommand: false,
		CommandId:          255,
	}, precedingValues
}

func calculatePrecedingHintBytes(precedingValues []float64, precedingCommands []CommandIdentifier) int {
	safeStemCount := func(counts int) int {
		if counts%2 == 0 {
			return counts / 2
		}
		return (counts - 1) / 2
	}

	stemCount := 0
	precedingNumbers := 0
	hasEncounteredInitialHintMask := false

	for vIdx := -1; vIdx < len(precedingValues); vIdx++ {
		if vIdx >= 0 {
			precedingNumbers++
		}

		for _, identifier := range precedingCommands {
			if identifier.CommandIndex != vIdx+1 {
				continue
			}

			if !identifier.IsMultiByteCommand && (identifier.CommandId == hintmaskByte || identifier.CommandId == cntrmaskByte) && !hasEncounteredInitialHintMask {
				hasEncounteredInitialHintMask = true
				stemCount += safeStemCount(precedingNumbers)
			} else if !identifier.IsMultiByteCommand && !hintingCommandBytes[identifier.CommandId] {
				precedingNumbers = 0
			} else if identifier.IsMultiByteCommand && identifier.CommandId > 35 {
				precedingNumbers = 0
			} else {
				stemCount += safeStemCount(precedingNumbers)
				precedingNumbers = 0
			}

			if hasEncounteredInitialHintMask {
				break
			}
		}

		if hasEncounteredInitialHintMask {
			break
		}
	}

	fullStemCount := stemCount
	if precedingNumbers > 0 && !hasEncounteredInitialHintMask {
		fullStemCount += safeStemCount(precedingNumbers)
	}

	minFullBytes := int(math.Ceil(float64(fullStemCount) / 8))

	return minFullBytes
}

func ternaryInt(condition bool, a, b int) int {
	if condition {
		return a
	}
	return b
}

func ternaryBool(condition bool, a, b float64) float64 {
	if condition {
		return a
	}
	return b
}
