package export

import (
	"fmt"
	"math"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	graphiccore "github.com/uglytoad/pdfpig/go/graphics/core"
)

var fontMap = map[string]string{
	"ArialMT": "Arial Rounded MT Bold",
}

// SvgTextExporter exports a page as an SVG document containing text and paths.
type SvgTextExporter struct {
	invalidCharacterHandler InvalidCharHandler
	rounding                int
	InvalidCharStrategy
}

// NewSvgTextExporter creates a new SVG text exporter with the given invalid character strategy.
// Defaults to DoNotCheck if not specified.
func NewSvgTextExporter(strategy InvalidCharStrategy) *SvgTextExporter {
	return &SvgTextExporter{
		invalidCharacterHandler: GetXmlInvalidCharHandler(strategy),
		rounding:                4,
		InvalidCharStrategy:     strategy,
	}
}

// NewSvgTextExporterWithHandler creates a new SVG text exporter with a custom
// invalid character handler function.
func NewSvgTextExporterWithHandler(handler InvalidCharHandler) *SvgTextExporter {
	return &SvgTextExporter{
		invalidCharacterHandler: handler,
		rounding:                4,
		InvalidCharStrategy:     Custom,
	}
}

// Get returns the page contents as an SVG string.
func (e *SvgTextExporter) Get(page *content.Page) string {
	var b strings.Builder
	fmt.Fprintf(&b, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width='%s' height='%s'>\n<g transform=\"scale(1, 1) translate(0, 0)\">\n",
		formatRound(page.Width(), e.rounding), formatRound(page.Height(), e.rounding))

	for _, path := range page.Paths() {
		if svgPath := e.pathToSvg(path, page.Height()); svgPath != "" {
			b.WriteString(svgPath)
			b.WriteByte('\n')
		}
	}

	for _, letter := range page.Letters() {
		b.WriteString(e.letterToSvg(letter, page.Height()))
	}

	b.WriteString("</g></svg>")
	return b.String()
}

func (e *SvgTextExporter) letterToSvg(l *content.Letter, height float64) string {
	fontFamily, style, weight := getFontFamily(l.FontName())

	rotation := ""
	if l.BoundingBox.Rotation() != 0 {
		bottomLeftX := math.Round(l.BoundingBox.BottomLeft.X*10000) / 10000
		topLeftY := math.Round(height-l.BoundingBox.TopLeft.Y*10000) / 10000
		rotAngle := math.Round(-l.BoundingBox.Rotation()*10000) / 10000
		rotation = fmt.Sprintf(" transform='rotate(%s %s,%s)'", formatRound(rotAngle, e.rounding), formatRound(bottomLeftX, e.rounding), formatRound(topLeftY, e.rounding))
	}

	fontSize := ""
	if l.FontSize != 1 {
		fontSize = fmt.Sprintf("font-size='%d'", int(l.FontSize))
	} else {
		fontSize = fmt.Sprintf("style='font-size:%spx'", formatRound(l.BoundingBox.Height, 2))
	}

	safeValue := xmlEscape(e.invalidCharacterHandler(l.Value))
	x := math.Round(l.StartBaseLine.X*10000) / 10000
	y := math.Round(height-l.StartBaseLine.Y*10000) / 10000

	return fmt.Sprintf("<text x='%s' y='%s'%s font-family='%s' font-style='%s' font-weight='%s' %s fill='%s'>%s</text>\n",
		formatRound(x, e.rounding), formatRound(y, e.rounding), rotation, fontFamily, style, weight, fontSize, colorToSvg(l.Color), safeValue)
}

func getFontFamily(fontName string) (family, style, weight string) {
	style = "normal"
	weight = "normal"

	if strings.Contains(fontName, "+") {
		if len(fontName) > 7 && fontName[6] == '+' {
			parts := strings.SplitN(fontName, "+", 2)
			if len(parts) == 2 && isAllUpper(parts[0]) {
				fontName = parts[1]
			}
		}
	}

	if strings.Contains(fontName, "-") {
		infos := strings.Split(fontName, "-")
		fontName = infos[0]

		for i := 1; i < len(infos); i++ {
			infoLower := strings.ToLower(infos[i])
			switch {
			case strings.Contains(infoLower, "light"):
				weight = "lighter"
			case strings.Contains(infoLower, "bolder"):
				weight = "bolder"
			case strings.Contains(infoLower, "bold"):
				weight = "bold"
			}

			switch {
			case strings.Contains(infoLower, "italic"):
				style = "italic"
			case strings.Contains(infoLower, "oblique"):
				style = "oblique"
			}
		}
	}

	if mapped, ok := fontMap[fontName]; ok {
		fontName = mapped
	}

	return fontName, style, weight
}

func isAllUpper(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return len(s) > 0
}

func xmlEscape(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '\'':
			b.WriteString("&apos;")
		case '"':
			b.WriteString("&quot;")
		default:
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
				fmt.Fprintf(&b, "&#x%X;", r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

func colorToSvg(c colors.Color) string {
	if c == nil {
		return ""
	}
	r := uint8(math.Round(c.ToRGBValues().R * 255))
	g := uint8(math.Round(c.ToRGBValues().G * 255))
	bv := uint8(math.Round(c.ToRGBValues().B * 255))
	return fmt.Sprintf("rgb(%d,%d,%d)", r, g, bv)
}

// pathInfo defines the interface needed from a PDF path for SVG export.
type pathInfo interface {
	IsClipping() bool
	Subpaths() []*core.PdfSubpath
	IsStroked() bool
	StrokeColor() colors.Color
	LineWidth() float64
	LineDashPattern() graphiccore.LineDashPattern
	LineCapStyle() graphiccore.LineCapStyle
	LineJoinStyle() graphiccore.LineJoinStyle
	IsFilled() bool
	FillColor() colors.Color
}

func (e *SvgTextExporter) pathToSvg(p content.PdfPath, height float64) string {
	pathObj, ok := any(&p).(pathInfo)
	if !ok {
		return ""
	}

	if pathObj.IsClipping() {
		return ""
	}

	var builder strings.Builder

	for _, subpath := range pathObj.Subpaths() {
		for _, command := range subpath.Commands() {
			command.WriteSvg(&builder, height)
		}
	}

	if builder.Len() == 0 {
		return ""
	}

	glyph := builder.String()
	if len(glyph) > 0 && glyph[len(glyph)-1] == ' ' {
		glyph = glyph[:len(glyph)-1]
	}

	dashArray := ""
	capStyle := ""
	jointStyle := ""
	strokeColor := " stroke='none'"
	strokeWidth := ""

	if pathObj.IsStroked() {
		strokeColor = fmt.Sprintf(" stroke='%s'", colorToSvg(pathObj.StrokeColor()))
		strokeWidth = fmt.Sprintf(" stroke-width='%g'", pathObj.LineWidth())

		dp := pathObj.LineDashPattern()
		if len(dp.Array) > 0 {
			parts := make([]string, len(dp.Array))
			for i, v := range dp.Array {
				parts[i] = formatNumber(v)
			}
			dashArray = fmt.Sprintf(" stroke-dasharray='%s'", strings.Join(parts, " "))
		}

		if pathObj.LineCapStyle() != graphiccore.LineCapButting {
			switch pathObj.LineCapStyle() {
			case graphiccore.LineCapRound:
				capStyle = " stroke-linecap='round'"
			default:
				capStyle = " stroke-linecap='square'"
			}
		}

		if pathObj.LineJoinStyle() != graphiccore.LineJoinMiter {
			switch pathObj.LineJoinStyle() {
			case graphiccore.LineJoinRound:
				jointStyle = " stroke-linejoin='round'"
			default:
				jointStyle = " stroke-linejoin='bevel'"
			}
		}
	}

	fillColor := " fill='none'"
	if pathObj.IsFilled() {
		fillColor = fmt.Sprintf(" fill='%s'", colorToSvg(pathObj.FillColor()))
	}

	return fmt.Sprintf("<path d='%s'%s%s%s%s%s%s></path>", glyph, fillColor, "", strokeColor, strokeWidth, dashArray, capStyle+jointStyle)
}

func formatRound(v float64, precision int) string {
	return fmt.Sprintf("%."+fmt.Sprint(precision)+"f", math.Round(v*math.Pow(10, float64(precision)))/math.Pow(10, float64(precision)))
}

func formatNumber(v float64) string {
	if v == float64(int(v)) {
		return fmt.Sprintf("%d", int(v))
	}
	return fmt.Sprintf("%g", v)
}

var _ TextExporter = (*SvgTextExporter)(nil)
