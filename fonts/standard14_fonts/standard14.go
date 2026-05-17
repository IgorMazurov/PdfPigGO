// Package standard14fonts provides access to the 14 standard Type 1 fonts required by PDF spec.
package standard14fonts

import (
	"embed"
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts/adobe_font_metrics"
)

//go:embed resources/adobe_font_metrics/*.afm
var afmFS embed.FS

var (
	standard14Names      map[string]struct{}
	standard14Mapping    map[string]string
	builderTypesToNames  map[Standard14Font]string
	standard14Cache      map[string]adobe_font_metrics.AdobeFontMetrics
)

func init() {
	standard14Names = make(map[string]struct{})
	standard14Mapping = make(map[string]string, 34)
	builderTypesToNames = make(map[Standard14Font]string, 14)
	standard14Cache = make(map[string]adobe_font_metrics.AdobeFontMetrics, 34)

	addAdobeFontMetrics("Courier-Bold", CourierBold)
	addAdobeFontMetrics("Courier-BoldOblique", CourierBoldOblique)
	addAdobeFontMetrics("Courier", Courier)
	addAdobeFontMetrics("Courier-Oblique", CourierOblique)
	addAdobeFontMetrics("Helvetica", Helvetica)
	addAdobeFontMetrics("Helvetica-Bold", HelveticaBold)
	addAdobeFontMetrics("Helvetica-BoldOblique", HelveticaBoldOblique)
	addAdobeFontMetrics("Helvetica-Oblique", HelveticaOblique)
	addAdobeFontMetrics("Symbol", Symbol)
	addAdobeFontMetrics("Times-Bold", TimesBold)
	addAdobeFontMetrics("Times-BoldItalic", TimesBoldItalic)
	addAdobeFontMetrics("Times-Italic", TimesItalic)
	addAdobeFontMetrics("Times-Roman", TimesRoman)
	addAdobeFontMetrics("ZapfDingbats", ZapfDingbats)

	// Alternative names from Adobe Supplement to the ISO 32000
	addAdobeFontMetricsAlias("CourierCourierNew", "Courier")
	addAdobeFontMetricsAlias("CourierNew", "Courier")
	addAdobeFontMetricsAlias("CourierNew,Italic", "Courier-Oblique")
	addAdobeFontMetricsAlias("CourierNew,Bold", "Courier-Bold")
	addAdobeFontMetricsAlias("CourierNew,BoldItalic", "Courier-BoldOblique")
	addAdobeFontMetricsAlias("Arial", "Helvetica")
	addAdobeFontMetricsAlias("Arial,Italic", "Helvetica-Oblique")
	addAdobeFontMetricsAlias("Arial,Bold", "Helvetica-Bold")
	addAdobeFontMetricsAlias("Arial,BoldItalic", "Helvetica-BoldOblique")
	addAdobeFontMetricsAlias("TimesNewRoman", "Times-Roman")
	addAdobeFontMetricsAlias("TimesNewRoman,Italic", "Times-Italic")
	addAdobeFontMetricsAlias("TimesNewRoman,Bold", "Times-Bold")
	addAdobeFontMetricsAlias("TimesNewRoman,BoldItalic", "Times-BoldItalic")

	// Acrobat treats these fonts as standard 14 too
	addAdobeFontMetricsAlias("Symbol,Italic", "Symbol")
	addAdobeFontMetricsAlias("Symbol,Bold", "Symbol")
	addAdobeFontMetricsAlias("Symbol,BoldItalic", "Symbol")
	addAdobeFontMetricsAlias("Times", "Times-Roman")
	addAdobeFontMetricsAlias("Times,Italic", "Times-Italic")
	addAdobeFontMetricsAlias("Times,Bold", "Times-Bold")
	addAdobeFontMetricsAlias("Times,BoldItalic", "Times-BoldItalic")

	// Additional names for Arial
	addAdobeFontMetricsAlias("ArialMT", "Helvetica")
	addAdobeFontMetricsAlias("Arial-ItalicMT", "Helvetica-Oblique")
	addAdobeFontMetricsAlias("Arial-BoldMT", "Helvetica-Bold")
	addAdobeFontMetricsAlias("Arial-BoldMT,Bold", "Helvetica-Bold")
	addAdobeFontMetricsAlias("Arial-BoldItalicMT", "Helvetica-BoldOblique")
}

func addAdobeFontMetrics(fontName string, fontType Standard14Font) {
	addAdobeFontMetricsInternal(fontName, fontName, &fontType)
}

func addAdobeFontMetricsAlias(fontName, afmName string) {
	addAdobeFontMetricsInternal(fontName, afmName, nil)
}

func addAdobeFontMetricsInternal(fontName, afmName string, fontType *Standard14Font) {
	standard14Names[fontName] = struct{}{}
	standard14Mapping[fontName] = afmName

	if fontType != nil {
		builderTypesToNames[*fontType] = afmName
	}

	if metrics, ok := standard14Cache[afmName]; ok {
		standard14Cache[fontName] = metrics
		return
	}

	data, err := afmFS.ReadFile("resources/adobe_font_metrics/" + afmName + ".afm")
	if err != nil {
		panic(fmt.Sprintf("could not find AFM resource: %s.afm", afmName))
	}

	bytes := core.NewMemoryInputBytes(data)
	metrics, err := adobe_font_metrics.Parse(bytes, true)
	if err != nil {
		panic(fmt.Sprintf("could not load %s from the AFM files: %v", fontName, err))
	}

	standard14Cache[fontName] = metrics
}

// GetAdobeFontMetrics returns the Adobe Font Metrics for a font by name.
// Returns (zero value, false) if the font is not found in the standard 14 set.
func GetAdobeFontMetrics(baseName string) (adobe_font_metrics.AdobeFontMetrics, bool) {
	metrics, ok := standard14Cache[baseName]
	return metrics, ok
}

// GetAdobeFontMetricsByType returns the Adobe Font Metrics for a Standard14 font type.
func GetAdobeFontMetricsByType(fontType Standard14Font) adobe_font_metrics.AdobeFontMetrics {
	afmName := builderTypesToNames[fontType]
	return standard14Cache[afmName]
}

// IsFontInStandard14 reports whether the given font name belongs to the standard 14 set.
func IsFontInStandard14(baseName string) bool {
	_, ok := standard14Names[baseName]
	return ok
}

// GetNames returns a copy of all known standard 14 font names, including aliases.
func GetNames() []string {
	names := make([]string, 0, len(standard14Names))
	for name := range standard14Names {
		names = append(names, name)
	}
	return names
}

// GetMappedFontName returns the official Standard 14 font name that the given name maps to.
// Returns empty string if the name is not in the mapping.
func GetMappedFontName(baseName string) string {
	name, ok := standard14Mapping[baseName]
	if !ok {
		return ""
	}
	return name
}
