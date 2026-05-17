package png

import (
	"fmt"
)

// GetBytesAndSamplesPerPixel calculates the bytes per pixel and samples per pixel
// from the image header's bit depth and color type.
func GetBytesAndSamplesPerPixel(header *ImageHeader) (bytesPerPixel byte, samplesPerPixel byte) {
	bitDepthCorrected := byte((int(header.BitDepth()) + 7) / 8)

	samplesPerPixel = spPerPixel(header)

	return samplesPerPixel * bitDepthCorrected, samplesPerPixel
}

// Decode reverses the PNG row filters and returns the raw pixel data.
// For non-interlaced images it operates in-place on decompressedData.
// For Adam7 interlaced images it allocates a new byte slice with pixels in raster order.
func Decode(decompressedData []byte, header *ImageHeader) ([]byte, error) {
	bytesPerPixel, samples := GetBytesAndSamplesPerPixel(header)

	switch header.InterlaceMethod() {
	case InterlaceMethodNone:
		return decodeNonInterlaced(decompressedData, header, bytesPerPixel, samples), nil

	case InterlaceMethodAdam7:
		return decodeAdam7(decompressedData, header, bytesPerPixel, samples), nil

	default:
		return nil, fmt.Errorf("invalid interlace method: %v", header.InterlaceMethod())
	}
}

func decodeNonInterlaced(data []byte, header *ImageHeader, bytesPerPixel, samplesPerPixel byte) []byte {
	bytesPerScanline := calcBytesPerScanline(header, samplesPerPixel)

	currentRowStartByteAbsolute := 1
	for rowIndex := 0; rowIndex < header.Height(); rowIndex++ {
		filterType := FilterType(data[currentRowStartByteAbsolute-1])

		previousRowStartByteAbsolute := rowIndex + bytesPerScanline*(rowIndex-1)

		end := currentRowStartByteAbsolute + bytesPerScanline
		for currentByteAbsolute := currentRowStartByteAbsolute; currentByteAbsolute < end; currentByteAbsolute++ {
			reverseFilter(data, filterType, previousRowStartByteAbsolute, currentRowStartByteAbsolute, currentByteAbsolute, currentByteAbsolute-currentRowStartByteAbsolute, int(bytesPerPixel))
		}

		currentRowStartByteAbsolute += bytesPerScanline + 1
	}

	return data
}

func decodeAdam7(data []byte, header *ImageHeader, bytesPerPixel byte, samplesPerPixel byte) []byte {
	pixelsPerRow := header.Width() * int(bytesPerPixel)
	newBytes := make([]byte, header.Height()*pixelsPerRow)
	i := 0
	previousStartRowByteAbsolute := -1

	for pass := 0; pass < 7; pass++ {
		numberOfScanlines := getNumberOfScanlinesInPass(header, pass)
		numberOfPixelsPerScanline := getPixelsPerScanlineInPass(header, pass)

		if numberOfScanlines <= 0 || numberOfPixelsPerScanline <= 0 {
			continue
		}

		for scanlineIndex := 0; scanlineIndex < numberOfScanlines; scanlineIndex++ {
			filterType := FilterType(data[i])
			i++
			rowStartByte := i

			for j := 0; j < numberOfPixelsPerScanline; j++ {
				pixelX, pixelY := getPixelIndexForScanlineInPass(header, pass, scanlineIndex, j)
				for k := byte(0); k < bytesPerPixel; k++ {
					byteLineNumber := (j * int(bytesPerPixel)) + int(k)
					reverseFilter(data, filterType, previousStartRowByteAbsolute, rowStartByte, i, byteLineNumber, int(bytesPerPixel))
					i++
				}

				start := pixelsPerRow*pixelY + pixelX*int(bytesPerPixel)
				copy(newBytes[start:start+int(bytesPerPixel)], data[rowStartByte+j*int(bytesPerPixel):rowStartByte+(j+1)*int(bytesPerPixel)])
			}

			previousStartRowByteAbsolute = rowStartByte
		}
	}

	return newBytes
}

// spPerPixel returns the number of color channels for the given ImageHeader's ColorType.
func spPerPixel(header *ImageHeader) byte {
	switch header.ColorType() {
	case ColorTypeNone:
		return 1
	case ColorTypePaletteUsed:
		return 1
	case ColorTypeColorUsed:
		return 3
	case ColorTypeAlphaChannelUsed:
		return 2
	case ColorTypeColorUsed | ColorTypeAlphaChannelUsed:
		return 4
	default:
		return 0
	}
}

func calcBytesPerScanline(header *ImageHeader, samplesPerPixel byte) int {
	width := header.Width()

	switch header.BitDepth() {
	case 1:
		return (width + 7) / 8
	case 2:
		return (width + 3) / 4
	case 4:
		return (width + 1) / 2
	case 8, 16:
		return width * int(samplesPerPixel) * (int(header.BitDepth()) / 8)
	default:
		return 0
	}
}

func reverseFilter(data []byte, filterType FilterType, previousRowStartByteAbsolute, rowStartByteAbsolute, byteAbsolute, rowByteIndex, bytesPerPixel int) {
	getLeft := func() byte {
		leftIndex := rowByteIndex - bytesPerPixel
		if leftIndex >= 0 {
			return data[rowStartByteAbsolute+leftIndex]
		}
		return 0
	}

	getAbove := func() byte {
		upIndex := previousRowStartByteAbsolute + rowByteIndex
		if upIndex >= 0 {
			return data[upIndex]
		}
		return 0
	}

	getAboveLeft := func() byte {
		index := previousRowStartByteAbsolute + rowByteIndex - bytesPerPixel
		if index < previousRowStartByteAbsolute || previousRowStartByteAbsolute < 0 {
			return 0
		}
		return data[index]
	}

	switch filterType {
	case FilterTypeNone:
		return

	case FilterTypeUp:
		above := previousRowStartByteAbsolute + rowByteIndex
		if above < 0 {
			return
		}
		data[byteAbsolute] += data[above]
		return

	case FilterTypeSub:
		leftIndex := rowByteIndex - bytesPerPixel
		if leftIndex < 0 {
			return
		}
		data[byteAbsolute] += data[rowStartByteAbsolute+leftIndex]
		return

	case FilterTypeAverage:
		data[byteAbsolute] += (getLeft() + getAbove()) / 2
		return

	case FilterTypePaeth:
		a := getLeft()
		b := getAbove()
		c := getAboveLeft()
		data[byteAbsolute] += paeth(a, b, c)
		return
	}
}

// paeth returns the Paeth predictor value based on left (a), above (b), and upper-left (c) neighbors.
func paeth(a, b, c byte) byte {
	p := int(a) + int(b) - int(c)
	pa := abs(p - int(a))
	pb := abs(p - int(b))
	pc := abs(p - int(c))

	if pa <= pb && pa <= pc {
		return a
	}

	if pb <= pc {
		return b
	}

	return c
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
