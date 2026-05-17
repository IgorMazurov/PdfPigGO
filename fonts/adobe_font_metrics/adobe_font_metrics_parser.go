package adobe_font_metrics

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
)

const (
	Comment          = "Comment"
	StartFontMetrics = "StartFontMetrics"
	EndFontMetrics   = "EndFontMetrics"
	FontName         = "FontName"
	FullName         = "FullName"
	FamilyName       = "FamilyName"
	Weight           = "Weight"
	FontBbox         = "FontBBox"
	Version          = "Version"
	Notice           = "Notice"
	EncodingScheme   = "EncodingScheme"
	MappingScheme    = "MappingScheme"
	EscChar          = "EscChar"
	CharacterSet     = "CharacterSet"
	Characters       = "Characters"
	IsBaseFont       = "IsBaseFont"
	VVector          = "VVector"
	IsFixedV         = "IsFixedV"
	CapHeight        = "CapHeight"
	XHeight          = "XHeight"
	Ascender         = "Ascender"
	Descender        = "Descender"
	UnderlinePosition  = "UnderlinePosition"
	UnderlineThickness = "UnderlineThickness"
	ItalicAngle      = "ItalicAngle"
	CharWidth        = "CharWidth"
	IsFixedPitch     = "IsFixedPitch"
	StartCharMetrics = "StartCharMetrics"
	EndCharMetrics   = "EndCharMetrics"
	CharmetricsC     = "C"
	CharmetricsCh    = "CH"
	CharmetricsWx    = "WX"
	CharmetricsW0X   = "W0X"
	CharmetricsW1X   = "W1X"
	CharmetricsWy    = "WY"
	CharmetricsW0Y   = "W0Y"
	CharmetricsW1Y   = "W1Y"
	CharmetricsW     = "W"
	CharmetricsW0    = "W0"
	CharmetricsW1    = "W1"
	CharmetricsVv    = "VV"
	CharmetricsN     = "N"
	CharmetricsB     = "B"
	CharmetricsL     = "L"
	StdHw            = "StdHW"
	StdVw            = "StdVW"
	StartTrackKern   = "StartTrackKern"
	EndTrackKern     = "EndTrackKern"
	StartKernData    = "StartKernData"
	EndKernData      = "EndKernData"
	StartKernPairs   = "StartKernPairs"
	EndKernPairs     = "EndKernPairs"
	StartKernPairs0  = "StartKernPairs0"
	StartKernPairs1  = "StartKernPairs1"
	StartComposites  = "StartComposites"
	EndComposites    = "EndComposites"
	Cc               = "CC"
	Pcc              = "PCC"
	KernPairKp       = "KP"
	KernPairKph      = "KPH"
	KernPairKpx      = "KPX"
	KernPairKpy      = "KPY"
)

var (
	characterNamesMu sync.Mutex
	characterNames   = make(map[string]string)
)

// Parse reads AFM data from the given input bytes and returns an AdobeFontMetrics instance.
func Parse(bytes core.InputBytes, useReducedDataSet bool) (AdobeFontMetrics, error) {
	var sb strings.Builder

	token := readString(bytes, &sb)

	if !strings.EqualFold(StartFontMetrics, token) {
		return AdobeFontMetrics{}, fonts.NewInvalidFontFormatException(
			fmt.Sprintf("The AFM file was not valid, it did not start with %s.", StartFontMetrics))
	}

	version := readDouble(bytes, &sb)

	builder := NewAdobeFontMetricsBuilder(version)

	for {
		token = readString(bytes, &sb)
		if token == EndFontMetrics {
			break
		}

		switch token {
		case Comment:
			line := readLine(bytes, &sb)
			builder.Comments = append(builder.Comments, line)
		case FontName:
			builder.FontName = readLine(bytes, &sb)
		case FullName:
			builder.FullName = readLine(bytes, &sb)
		case FamilyName:
			builder.FamilyName = readLine(bytes, &sb)
		case Weight:
			builder.Weight = readLine(bytes, &sb)
		case ItalicAngle:
			builder.ItalicAngle = readDouble(bytes, &sb)
		case IsFixedPitch:
			b, err := readBool(bytes, &sb)
			if err != nil {
				return AdobeFontMetrics{}, err
			}
			builder.IsFixedPitch = b
		case FontBbox:
			x1 := readDouble(bytes, &sb)
			y1 := readDouble(bytes, &sb)
			x2 := readDouble(bytes, &sb)
			y2 := readDouble(bytes, &sb)
			builder.SetBoundingBox(x1, y1, x2, y2)
		case UnderlinePosition:
			builder.UnderlinePosition = readDouble(bytes, &sb)
		case UnderlineThickness:
			builder.UnderlineThickness = readDouble(bytes, &sb)
		case Version:
			builder.Version = readLine(bytes, &sb)
		case Notice:
			builder.Notice = readLine(bytes, &sb)
		case EncodingScheme:
			builder.EncodingScheme = readLine(bytes, &sb)
		case MappingScheme:
			builder.MappingScheme = int(readDouble(bytes, &sb))
		case CharacterSet:
			builder.CharacterSet = readLine(bytes, &sb)
		case EscChar:
			builder.EscapeCharacter = int(readDouble(bytes, &sb))
		case Characters:
			builder.Characters = int(readDouble(bytes, &sb))
		case IsBaseFont:
			b, err := readBool(bytes, &sb)
			if err != nil {
				return AdobeFontMetrics{}, err
			}
			builder.IsBaseFont = b
		case CapHeight:
			builder.CapHeight = readDouble(bytes, &sb)
		case XHeight:
			builder.XHeight = readDouble(bytes, &sb)
		case Ascender:
			builder.Ascender = readDouble(bytes, &sb)
		case Descender:
			builder.Descender = readDouble(bytes, &sb)
		case StdHw:
			builder.StdHw = readDouble(bytes, &sb)
		case StdVw:
			builder.StdVw = readDouble(bytes, &sb)
		case CharWidth:
			x := readDouble(bytes, &sb)
			y := readDouble(bytes, &sb)
			builder.SetCharacterWidth(x, y)
		case VVector:
			x := readDouble(bytes, &sb)
			y := readDouble(bytes, &sb)
			builder.SetVVector(x, y)
		case IsFixedV:
			b, err := readBool(bytes, &sb)
			if err != nil {
				return AdobeFontMetrics{}, err
			}
			builder.IsFixedV = b
		case StartCharMetrics:
			count := int(readDouble(bytes, &sb))
			for i := 0; i < count; i++ {
				metric, err := readCharacterMetric(bytes, &sb)
				if err != nil {
					return AdobeFontMetrics{}, err
				}
				builder.CharacterMetrics = append(builder.CharacterMetrics, metric)
			}

			end := readString(bytes, &sb)
			if end != EndCharMetrics {
				return AdobeFontMetrics{}, fonts.NewInvalidFontFormatException(
					fmt.Sprintf("The character metrics section did not end with %s instead it was %s.", EndCharMetrics, end))
			}
		case StartKernData:
		}
	}

	return builder.Build(), nil
}

func readDouble(input core.InputBytes, sb *strings.Builder) float64 {
	s := readString(input, sb)
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func readBool(input core.InputBytes, sb *strings.Builder) (bool, error) {
	boolean := readString(input, sb)

	switch boolean {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fonts.NewInvalidFontFormatException(
			fmt.Sprintf("The AFM should have contained a boolean but instead contained: %s.", boolean))
	}
}

func readString(input core.InputBytes, sb *strings.Builder) string {
	sb.Reset()

	if input.IsAtEnd() {
		return EndFontMetrics
	}

	for core.IsWhitespace(input.CurrentByte()) && input.MoveNext() {
	}

	sb.WriteByte(input.CurrentByte())

	for input.MoveNext() && !core.IsWhitespace(input.CurrentByte()) {
		sb.WriteByte(input.CurrentByte())
	}

	return sb.String()
}

func readLine(input core.InputBytes, sb *strings.Builder) string {
	sb.Reset()

	for core.IsWhitespace(input.CurrentByte()) && input.MoveNext() {
	}

	sb.WriteByte(input.CurrentByte())

	for input.MoveNext() && !core.IsEndOfLineByte(input.CurrentByte()) {
		sb.WriteByte(input.CurrentByte())
	}

	return sb.String()
}

func readCharacterMetric(bytes core.InputBytes, sb *strings.Builder) (AdobeFontMetricsIndividualCharacterMetric, error) {
	line := readLine(bytes, sb)

	split := strings.Split(line, ";")

	metric := AdobeFontMetricsIndividualCharacterMetricBuilder{}

	for _, s := range split {
		parts := strings.Fields(s)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case CharmetricsC:
			code, err := strconv.Atoi(parts[1])
			if err != nil {
				return AdobeFontMetricsIndividualCharacterMetric{}, fonts.NewInvalidFontFormatException(
					fmt.Sprintf("Invalid character code '%s'.", parts[1]))
			}
			metric.CharacterCode = code
		case CharmetricsCh:
			code, err := strconv.ParseInt(parts[1], 16, 64)
			if err != nil {
				return AdobeFontMetricsIndividualCharacterMetric{}, fonts.NewInvalidFontFormatException(
					fmt.Sprintf("Invalid hex character code '%s'.", parts[1]))
			}
			metric.CharacterCode = int(code)
		case CharmetricsWx:
			v, _ := strconv.ParseFloat(parts[1], 64)
			metric.WidthX = v
		case CharmetricsW0X:
			v, _ := strconv.ParseFloat(parts[1], 64)
			metric.WidthXDirection0 = v
		case CharmetricsW1X:
			v, _ := strconv.ParseFloat(parts[1], 64)
			metric.WidthXDirection1 = v
		case CharmetricsWy:
			v, _ := strconv.ParseFloat(parts[1], 64)
			metric.WidthY = v
		case CharmetricsW0Y:
			v, _ := strconv.ParseFloat(parts[1], 64)
			metric.WidthYDirection0 = v
		case CharmetricsW1Y:
			v, _ := strconv.ParseFloat(parts[1], 64)
			metric.WidthYDirection1 = v
		case CharmetricsW:
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			metric.WidthX = x
			metric.WidthY = y
		case CharmetricsW0:
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			metric.WidthXDirection0 = x
			metric.WidthYDirection0 = y
		case CharmetricsW1:
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			metric.WidthXDirection1 = x
			metric.WidthYDirection1 = y
		case CharmetricsVv:
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			metric.VVector = AdobeFontMetricsVector{X: x, Y: y}
		case CharmetricsN:
			name := parts[1]
			cached, err := getOrCreateCharacterName(name)
			if err != nil {
				return AdobeFontMetricsIndividualCharacterMetric{}, err
			}
			metric.Name = cached
		case CharmetricsB:
			x0, _ := strconv.ParseFloat(parts[1], 64)
			y0, _ := strconv.ParseFloat(parts[2], 64)
			x1, _ := strconv.ParseFloat(parts[3], 64)
			y1, _ := strconv.ParseFloat(parts[4], 64)
			metric.BoundingBox = core.NewPdfRectangleFloat(x0, y0, x1, y1)
		case CharmetricsL:
			metric.Ligature = &AdobeFontMetricsLigature{
				Successor: parts[1],
				Value:     parts[2],
			}
		default:
			return AdobeFontMetricsIndividualCharacterMetric{}, fonts.NewInvalidFontFormatException(
				fmt.Sprintf("Unknown CharMetrics command '%s'.", parts[0]))
		}
	}

	return metric.Build(), nil
}

func getOrCreateCharacterName(name string) (string, error) {
	characterNamesMu.Lock()
	defer characterNamesMu.Unlock()

	cached, ok := characterNames[name]
	if !ok {
		cached = name
		characterNames[name] = cached
	}

	return cached, nil
}
