package lzw

import (
	"errors"
)

// BitStream provides bit-level reading from a byte slice.
// Mirrors C# internal ref struct UglyToad.PdfPig.Filters.Lzw.BitStream.
type BitStream struct {
	data                     []byte
	currentWithinByteBitOffset int
	currentByteIndex         int
}

// NewBitStream creates a new BitStream over the given byte slice.
func NewBitStream(data []byte) *BitStream {
	return &BitStream{
		data: data,
	}
}

// Get reads numberOfBits bits from the stream and returns them as an integer.
// numberOfBits must be less than 32. Returns an error if the end of the stream
// is reached before all requested bits can be read.
func (b *BitStream) Get(numberOfBits int) (int32, error) {
	endWithinByteBitOffset := (numberOfBits + b.currentWithinByteBitOffset) % 8

	numberOfBytesToRead := (numberOfBits + b.currentWithinByteBitOffset) / 8

	if endWithinByteBitOffset != 0 {
		numberOfBytesToRead++
	}

	var result int32
	for i := 0; i < numberOfBytesToRead; i++ {
		if i > 0 {
			b.currentByteIndex++
		}

		if b.currentByteIndex >= len(b.data) {
			return 0, errors.New("reached the end of the bit stream while trying to read bits")
		}

		result <<= 8
		result |= int32(b.data[b.currentByteIndex])
	}

	// Trim trailing bits.
	if endWithinByteBitOffset > 0 {
		result >>= int32(8 - endWithinByteBitOffset)
	} else {
		b.currentByteIndex++
	}

	// Mask out the leading bits.
	result &= (1 << uint(numberOfBits)) - 1

	b.currentWithinByteBitOffset = endWithinByteBitOffset

	return result, nil
}
