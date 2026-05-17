package png

// passToScanlineGridIndex maps a pass number (1-indexed) to the scanline grid indices
// included in that pass within an 8x8 repeating block.
var passToScanlineGridIndex = [][]int{
	{0},
	{0},
	{4},
	{0, 4},
	{2, 6},
	{0, 2, 4, 6},
	{1, 3, 5, 7},
}

// passToScanlineColumnIndex maps a pass number (1-indexed) to the column indices
// included in that pass within an 8x8 repeating block.
var passToScanlineColumnIndex = [][]int{
	{0},
	{4},
	{0, 4},
	{2, 6},
	{0, 2, 4, 6},
	{1, 3, 5, 7},
	{0, 1, 2, 3, 4, 5, 6, 7},
}

// getNumberOfScanlinesInPass returns the number of scanlines in the given pass for the image.
func getNumberOfScanlinesInPass(header *ImageHeader, pass int) int {
	indices := passToScanlineGridIndex[pass]

	mod := header.Height() % 8

	if mod == 0 {
		return len(indices) * (header.Height() / 8)
	}

	additionalLines := 0
	for _, idx := range indices {
		if idx < mod {
			additionalLines++
		}
	}

	return len(indices)*(header.Height()/8) + additionalLines
}

// getPixelsPerScanlineInPass returns the number of pixels per scanline in the given pass for the image.
func getPixelsPerScanlineInPass(header *ImageHeader, pass int) int {
	indices := passToScanlineColumnIndex[pass]

	mod := header.Width() % 8

	if mod == 0 {
		return len(indices) * (header.Width() / 8)
	}

	additionalColumns := 0
	for _, idx := range indices {
		if idx < mod {
			additionalColumns++
		}
	}

	return len(indices)*(header.Width()/8) + additionalColumns
}

// getPixelIndexForScanlineInPass returns the actual (x, y) pixel coordinates for a given
// scanline index and position within that scanline during the specified Adam7 pass.
func getPixelIndexForScanlineInPass(header *ImageHeader, pass, scanlineIndex, indexInScanline int) (x, y int) {
	columnIndices := passToScanlineColumnIndex[pass]
	rows := passToScanlineGridIndex[pass]

	actualRow := scanlineIndex % len(rows)
	actualCol := indexInScanline % len(columnIndices)
	precedingRows := 8 * (scanlineIndex/len(rows))
	precedingCols := 8 * (indexInScanline/len(columnIndices))

	return precedingCols + columnIndices[actualCol], precedingRows + rows[actualRow]
}
