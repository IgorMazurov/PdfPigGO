package cff

import "fmt"

// ReadDictionaryData reads the index from data and extracts each indexed entry as a byte slice.
func ReadDictionaryData(data *CompactFontFormatData) (*CompactFontFormatIndex, error) {
	index := ReadIndex(data)

	if len(index) == 0 {
		return NewCompactFontFormatIndex(nil), nil
	}

	count := len(index) - 1
	results := make([][]byte, count)

	for i := 0; i < count; i++ {
		length := index[i+1] - index[i]

		if length < 0 {
			return nil, fmt.Errorf("negative object length %d at %d. Current position: %d", length, i, data.Position())
		}

		if length > data.Length() {
			return nil, fmt.Errorf("attempted to read data of length %d in data array of length %d", length, data.Length())
		}

		results[i] = data.ReadBytes(length)
	}

	return NewCompactFontFormatIndex(results), nil
}

// ReadIndex reads the count and offsets from the data stream.
func ReadIndex(data *CompactFontFormatData) []int {
	count := int(data.ReadCard16())

	if count == 0 {
		return []int{}
	}

	offsetSize := int(data.ReadOffsize())

	offsets := make([]int, count+1)

	for i := 0; i < len(offsets); i++ {
		offsets[i] = data.ReadOffset(offsetSize)
	}

	return offsets
}
