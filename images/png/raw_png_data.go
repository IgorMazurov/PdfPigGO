package png

import (
	"errors"
)

// RawPngData holds raw decoded PNG image data and provides pixel access.
type RawPngData struct {
	data         []byte
	bytesPerPixel int
	width        int
	palette      *Palette
	colorType    ColorType
	rowOffset    int
	bitDepth     byte
}

// NewRawPngData creates a new RawPngData from raw image bytes and header info.
func NewRawPngData(data []byte, bytesPerPixel int, palette *Palette, header ImageHeader) (*RawPngData, error) {
	if data == nil {
		return nil, errors.New("data must not be nil")
	}

	rowOffset := 1
	if header.InterlaceMethod() == InterlaceMethodAdam7 {
		rowOffset = 0
	}

	return &RawPngData{
		data:         data,
		bytesPerPixel: bytesPerPixel,
		width:        header.Width(),
		palette:      palette,
		colorType:    header.ColorType(),
		rowOffset:    rowOffset,
		bitDepth:     header.BitDepth(),
	}, nil
}

// toSingleByte converts two bytes representing a 16-bit big-endian value into a single
// byte by scaling the range [0, 65535] down to [0, 255].
func toSingleByte(first, second byte) byte {
	us := uint16(first)<<8 | uint16(second)
	return byte((255*int(us) + 32768) / 65536)
}

// GetPixel returns the Pixel at the given (x, y) coordinate.
func (r *RawPngData) GetPixel(x, y int) Pixel {
	if r.palette != nil {
		pixelsPerByte := 8 / int(r.bitDepth)
		bytesInRow := 1 + (r.width / pixelsPerByte)
		byteIndexInRow := x / pixelsPerByte
		paletteIndex := (1 + (y * bytesInRow)) + byteIndexInRow
		b := r.data[paletteIndex]

		if r.bitDepth == 8 {
			return r.palette.GetPixel(int(b))
		}

		withinByteIndex := x % pixelsPerByte
		rightShift := 8 - ((withinByteIndex+1)*int(r.bitDepth))
		indexActual := (b >> rightShift) & ((1 << r.bitDepth) - 1)
		return r.palette.GetPixel(int(indexActual))
	}

	rowStartPixel := (r.rowOffset + (r.rowOffset * y)) + (r.bytesPerPixel * r.width * y)
	pixelStartIndex := rowStartPixel + (r.bytesPerPixel * x)
	first := r.data[pixelStartIndex]

	switch r.bytesPerPixel {
	case 1:
		return NewGrayPixel(first)

	case 2:
		if r.colorType == ColorTypeNone {
			value := toSingleByte(first, r.data[pixelStartIndex+1])
			return NewGrayPixel(value)
		}
		return NewPixel(first, first, first, r.data[pixelStartIndex+1], true)

	case 3:
		return NewPixel(first, r.data[pixelStartIndex+1], r.data[pixelStartIndex+2], 255, false)

	case 4:
		if r.colorType == ColorTypeNone || r.colorType == ColorTypeAlphaChannelUsed {
			gray := toSingleByte(first, r.data[pixelStartIndex+1])
			alpha := toSingleByte(r.data[pixelStartIndex+2], r.data[pixelStartIndex+3])
			return NewPixel(gray, gray, gray, alpha, true)
		}
		return NewPixel(first, r.data[pixelStartIndex+1], r.data[pixelStartIndex+2], r.data[pixelStartIndex+3], false)

	case 6:
		return NewPixel(first, r.data[pixelStartIndex+2], r.data[pixelStartIndex+4], 255, false)

	case 8:
		pos := pixelStartIndex
		R := toSingleByte(r.data[pos], r.data[pos+1])
		G := toSingleByte(r.data[pos+2], r.data[pos+3])
		B := toSingleByte(r.data[pos+4], r.data[pos+5])
		A := toSingleByte(r.data[pos+6], r.data[pos+7])
		return NewPixel(R, G, B, A, false)

	default:
		panic("unrecognized number of bytes per pixel")
	}
}
