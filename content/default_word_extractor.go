package content

import (
	"math"
	"sort"
	"strings"
)

// WordExtractor extracts words from a list of letters.
type WordExtractor interface {
	// GetWords returns the words extracted from the given letters.
	GetWords(letters []*Letter) []*Word
}

// DefaultWordExtractor is the default implementation of WordExtractor.
type DefaultWordExtractor struct{}

var _ WordExtractor = (*DefaultWordExtractor)(nil)

// Instance returns the shared DefaultWordExtractor singleton.
var Instance WordExtractor = &DefaultWordExtractor{}

// GetWords returns the words extracted from the given letters.
func (e *DefaultWordExtractor) GetWords(letters []*Letter) []*Word {
	lettersOrder := make([]*Letter, len(letters))
	copy(lettersOrder, letters)
	sort.Slice(lettersOrder, func(i, j int) bool {
		yi := lettersOrder[i].Location().Y
		yj := lettersOrder[j].Location().Y
		if yi != yj {
			return yi > yj
		}
		return lettersOrder[i].Location().X < lettersOrder[j].Location().X
	})

	lettersSoFar := make([]*Letter, 0, 10)
	gapCountsSoFarByFontSize := make(map[float64]map[float64]int)
	gapInsertionOrderByFontSize := make(map[float64][]float64)
	var words []*Word
	var y float64
	var ySet bool
	var lastX float64
	var lastXSet bool
	var lastLetter *Letter

	for _, letter := range lettersOrder {
		if !ySet {
			y = letter.Location().Y
			ySet = true
		}

		if !lastXSet {
			lastX = letter.Location().X
			lastXSet = true
		}

		if lastLetter == nil {
			if strings.TrimSpace(letter.Value) == "" {
				continue
			}

			lettersSoFar = append(lettersSoFar, letter)
			lastLetter = letter
			y = letter.Location().Y
			lastX = letter.Location().X
			continue
		}

		if letter.Location().Y < y-0.5 {
			if len(lettersSoFar) > 0 {
				if word := generateWord(lettersSoFar); word != nil {
					words = append(words, word)
				}
				lettersSoFar = make([]*Letter, 0, 10)
			}

			if strings.TrimSpace(letter.Value) != "" {
				lettersSoFar = append(lettersSoFar, letter)
			}

			y = letter.Location().Y
			lastX = letter.Location().X
			lastLetter = letter

			continue
		}

		letterHeight := math.Max(lastLetter.BoundingBox.Height, letter.BoundingBox.Height)
		gap := letter.Location().X - (lastLetter.Location().X+lastLetter.Width)
		nextToLeft := letter.Location().X < lastX-1
		nextBigSpace := gap > letterHeight*0.39
		nextIsWhiteSpace := strings.TrimSpace(letter.Value) == ""
		nextFontDiffers := !strings.EqualFold(letter.FontName(), lastLetter.FontName()) && gap > letter.Width*0.1
		nextFontSizeDiffers := math.Abs(letter.FontSize-lastLetter.FontSize) > 0.1
		nextTextOrientationDiffers := letter.TextOrientation != lastLetter.TextOrientation

		suspectGap := false

		if !nextFontSizeDiffers && letter.FontSize > 0 && gap >= 0 {
			fontSize := roundHalfToEven(letter.FontSize)
			gapCounts, ok := gapCountsSoFarByFontSize[fontSize]
			if !ok {
				gapCounts = make(map[float64]int)
				gapCountsSoFarByFontSize[fontSize] = gapCounts
			}
			insertionOrder := gapInsertionOrderByFontSize[fontSize]

			gapRounded := roundHalfToEvenDigits(gap, 2)
			if gapCounts[gapRounded] == 0 {
				insertionOrder = append(insertionOrder, gapRounded)
				gapInsertionOrderByFontSize[fontSize] = insertionOrder
			}
			gapCounts[gapRounded]++

			if len(gapCounts) > 1 && gap > letterHeight*0.16 {
				mostCommonKey, mostCommonValue := findMostCommonGap(gapCounts, insertionOrder)

				if gap > mostCommonKey*5 && mostCommonValue > 1 {
					suspectGap = true
				}
			}
		}

		if nextToLeft || nextBigSpace || nextIsWhiteSpace || nextFontDiffers || nextFontSizeDiffers || nextTextOrientationDiffers || suspectGap {
			if len(lettersSoFar) > 0 {
				if word := generateWord(lettersSoFar); word != nil {
					words = append(words, word)
				}
				lettersSoFar = make([]*Letter, 0, 10)
			}
		}

		if strings.TrimSpace(letter.Value) != "" {
			lettersSoFar = append(lettersSoFar, letter)
		}

		lastLetter = letter
		lastX = letter.Location().X
	}

	if len(lettersSoFar) > 0 {
		if word := generateWord(lettersSoFar); word != nil {
			words = append(words, word)
		}
	}

	return words
}

func findMostCommonGap(gapCounts map[float64]int, insertionOrder []float64) (float64, int) {
	mostCommonKey := float64(0)
	mostCommonValue := 0
	for _, k := range insertionOrder {
		v := gapCounts[k]
		if v > mostCommonValue {
			mostCommonKey = k
			mostCommonValue = v
		}
	}
	return mostCommonKey, mostCommonValue
}

// roundHalfToEven rounds x to the nearest integer using banker's rounding
// (round half to even), matching C#'s Math.Round(double) behavior.
func roundHalfToEven(x float64) float64 {
	i := int64(math.Floor(x))
	frac := x - float64(i)
	if math.Abs(frac-0.5) < 1e-9 {
		if i%2 == 0 {
			return float64(i)
		}
		return float64(i + 1)
	}
	if frac < 0.5 {
		return float64(i)
	}
	return float64(i + 1)
}

// roundHalfToEvenDigits rounds x to the given number of decimal places
// using banker's rounding, matching C#'s Math.Round(double, int) behavior.
func roundHalfToEvenDigits(x float64, digits int) float64 {
	factor := math.Pow(10, float64(digits))
	scaled := x * factor
	i := int64(math.Floor(scaled))
	frac := scaled - float64(i)
	if math.Abs(frac-0.5) < 1e-9 {
		if i%2 == 0 {
			return float64(i) / factor
		}
		return float64(i+1) / factor
	}
	if frac < 0.5 {
		return float64(i) / factor
	}
	return float64(i+1) / factor
}

func generateWord(letters []*Letter) *Word {
	word, _ := NewWord(letters)
	return word
}
