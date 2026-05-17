package text_extractor

import (
	"math"

	"github.com/uglytoad/pdfpig/go/content"
)

var replaceableWhitespace = map[string]bool{
	"\t": true,
	"\v": true,
	"\f": true,
}

// Options controls the text generation algorithm.
type Options struct {
	// SeparateParagraphsWithDoubleNewline includes a double new-line when the text is likely to be a new paragraph.
	SeparateParagraphsWithDoubleNewline bool
	// ReplaceWhitespaceWithSpace replaces all whitespace characters (except line breaks) with single space ' '.
	ReplaceWhitespaceWithSpace bool
	// NegativeGapAsWhitespace better predicts spaces between words in tables containing multiple lines or merged cells.
	NegativeGapAsWhitespace bool
}

// GetText returns a human readable representation of the text from the page based on
// the letter order of the original PDF document.
func GetText(page *content.Page, addDoubleNewline bool) string {
	return GetTextWithOptions(page, Options{
		SeparateParagraphsWithDoubleNewline: addDoubleNewline,
	})
}

// GetTextWithOptions returns a human readable representation of the text from the page based on
// the letter order of the original PDF document with customizable options.
func GetTextWithOptions(page *content.Page, opts Options) string {
	letters := page.Letters()
	if letters == nil || len(letters) == 0 {
		return ""
	}

	var runeSlice []rune
	appendRunes := func(s string) {
		runeSlice = append(runeSlice, []rune(s)...)
	}
	removeLastRune := func() {
		if len(runeSlice) > 0 {
			runeSlice = runeSlice[:len(runeSlice)-1]
		}
	}

	var previous *content.Letter
	hasJustAddedWhitespace := false

	for i := 0; i < len(letters); i++ {
		letter := letters[i]

		if letter.Value == "" {
			continue
		}

		if opts.ReplaceWhitespaceWithSpace && replaceableWhitespace[letter.Value] {
			letter = copyLetter(letter, " ")
		}

		if letter.Value == " " && !hasJustAddedWhitespace {
			nextNonWS := getNextNonWhitespace(letters, i)
			if nextNonWS != nil && isLineBreak(previous, nextNonWS) {
				continue
			}

			appendRunes(" ")
			previous = letter
			hasJustAddedWhitespace = true
			continue
		}

		hasJustAddedWhitespace = false

		if previous != nil && letter.Value != " " {
			nwPrevious := getNonWhitespacePrevious(letters, i)

			isDoubleNewline := false
			if nwPrevious != nil {
				isDoubleNewline = isNewlineDouble(nwPrevious, letter)
			}

			if isLineBreak(nwPrevious, letter) {
				if previous.Value == " " {
					removeLastRune()
				}

				appendRunes("\n")
				if opts.SeparateParagraphsWithDoubleNewline && isDoubleNewline {
					appendRunes("\n")
				}

				hasJustAddedWhitespace = true
			} else if previous.Value != " " {
				gap := letter.StartBaseLine.X - previous.EndBaseLine.X

				if opts.NegativeGapAsWhitespace {
					gap = math.Abs(gap)
				}

				if content.IsProbablyWhitespace(gap, previous) {
					appendRunes(" ")
					hasJustAddedWhitespace = true
				}
			}
		}

		appendRunes(letter.Value)
		previous = letter
	}

	return string(runeSlice)
}

func getNextNonWhitespace(letters []*content.Letter, index int) *content.Letter {
	for i := index + 1; i < len(letters); i++ {
		letter := letters[i]
		if letter.Value != "" && !isWhitespaceOnly(letter.Value) {
			return letter
		}
	}

	return nil
}

func getNonWhitespacePrevious(letters []*content.Letter, index int) *content.Letter {
	for i := index - 1; i >= 0; i-- {
		letter := letters[i]
		if letter.Value != "" && !isWhitespaceOnly(letter.Value) {
			return letter
		}
	}

	return nil
}

func isWhitespaceOnly(value string) bool {
	for _, r := range value {
		if !isWhitespaceRune(r) {
			return false
		}
	}
	return true
}

func isWhitespaceRune(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\v' || r == '\f'
}

// isLineBreak checks if there's a line break between previous and letter.
// Returns true when the vertical gap exceeds 0.9x the minimum point size
// and the previous baseline is below (higher Y value in PDF coords).
func isLineBreak(previous, letter *content.Letter) bool {
	if previous == nil {
		return false
	}

	ptSizePrevious := int(math.Round(previous.PointSize))
	ptSizeLetter := int(math.Round(letter.PointSize))
	minPtSize := ptSizePrevious
	if ptSizeLetter < ptSizePrevious {
		minPtSize = ptSizeLetter
	}

	gap := math.Abs(previous.StartBaseLine.Y - letter.StartBaseLine.Y)

	return gap > float64(minPtSize)*0.9 && previous.StartBaseLine.Y > letter.StartBaseLine.Y
}

// isNewlineDouble checks if the line break qualifies as a double newline (paragraph break).
// Returns true when the vertical gap exceeds 1.7x the minimum point size.
func isNewlineDouble(previous, letter *content.Letter) bool {
	if previous == nil {
		return false
	}

	ptSizePrevious := int(math.Round(previous.PointSize))
	ptSizeLetter := int(math.Round(letter.PointSize))
	minPtSize := ptSizePrevious
	if ptSizeLetter < ptSizePrevious {
		minPtSize = ptSizeLetter
	}

	gap := math.Abs(previous.StartBaseLine.Y - letter.StartBaseLine.Y)

	return gap > float64(minPtSize)*1.7 && previous.StartBaseLine.Y > letter.StartBaseLine.Y
}

// copyLetter creates a new Letter with the same properties but a different value.
func copyLetter(original *content.Letter, newValue string) *content.Letter {
	return &content.Letter{
		Value:               newValue,
		BoundingBox:         original.BoundingBox,
		GlyphRectangleLoose: original.GlyphRectangleLoose,
		StartBaseLine:       original.StartBaseLine,
		EndBaseLine:         original.EndBaseLine,
		Width:               original.Width,
		FontSize:            original.FontSize,
		FontDetails:         original.FontDetails,
		RenderingMode:       original.RenderingMode,
		StrokeColor:         original.StrokeColor,
		FillColor:           original.FillColor,
		PointSize:           original.PointSize,
		TextSequence:        original.TextSequence,
		TextOrientation:     original.TextOrientation,
		Color:               original.Color,
	}
}
