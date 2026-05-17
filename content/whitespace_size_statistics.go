package content

// GetExpectedWhitespace returns the average whitespace size expected for a given letter.
func GetExpectedWhitespace(letter *Letter) float64 {
	return letter.PointSize * 0.27
}

// IsProbablyWhitespace checks if the measured gap is probably big enough to be a
// whitespace character based on the letter.
func IsProbablyWhitespace(gap float64, letter *Letter) bool {
	return gap > (GetExpectedWhitespace(letter)-(letter.PointSize*0.05))
}
