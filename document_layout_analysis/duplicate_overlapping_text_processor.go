package document_layout_analysis

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
)

// Get removes duplicate overlapping letters from the given slice.
// Two letters are considered duplicates if they share the same Value and FontName,
// and their bounding box bottom-left corners fall within a tolerance region.
// When an overlap is detected, the existing letter is replaced with its bold variant.
func Get(letters []*content.Letter) []*content.Letter {
	if len(letters) == 0 {
		return letters
	}

	duplicateIndex := make(map[string][]int)
	cleanLetters := make([]*content.Letter, 0, len(letters))

	for _, letter := range letters {
		addLetter := true
		duplicatesOverlappingIndex := -1

		key := makeKey(letter.Value, letter.FontName())

		if candidateIndices, ok := duplicateIndex[key]; ok {
			valueLen := len(letter.Value)
			if valueLen == 0 {
				valueLen = 1
			}
			tolerance := letter.BoundingBox.Width / float64(valueLen) / 3.0
			minX := letter.BoundingBox.BottomLeft.X - tolerance
			maxX := letter.BoundingBox.BottomLeft.X + tolerance
			minY := letter.BoundingBox.BottomLeft.Y - tolerance
			maxY := letter.BoundingBox.BottomLeft.Y + tolerance

			for _, idx := range candidateIndices {
				l := cleanLetters[idx]
				if minX <= l.BoundingBox.BottomLeft.X &&
					maxX >= l.BoundingBox.BottomLeft.X &&
					minY <= l.BoundingBox.BottomLeft.Y &&
					maxY >= l.BoundingBox.BottomLeft.Y {
					addLetter = false
					duplicatesOverlappingIndex = idx
					break
				}
			}
		}

		if addLetter {
			newIndex := len(cleanLetters)
			cleanLetters = append(cleanLetters, letter)

			list, ok := duplicateIndex[key]
			if !ok {
				list = make([]int, 0, 2)
			}
			duplicateIndex[key] = append(list, newIndex)
		} else if duplicatesOverlappingIndex != -1 {
			cleanLetters[duplicatesOverlappingIndex] = letter.AsBold()
		}
	}

	return cleanLetters
}

func makeKey(value, fontName string) string {
	return fmt.Sprintf("%s\x00%s", value, fontName)
}
