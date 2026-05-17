package filters

import (
	"io"
)

// DecodePredictorRow decodes a single line of data in-place using PNG-style
// prediction. The actline buffer is modified directly with decoded values.
// lastline holds the previously decoded row; for the first line it should be
// an empty byte slice of the same length as actline.
func DecodePredictorRow(predictor, colors, bitsPerComponent, columns int, actline, lastline []byte) {
	if predictor == 1 {
		return
	}

	bitsPerPixel := colors * bitsPerComponent
	bytesPerPixel := (bitsPerPixel + 7) / 8
	rowLength := len(actline)

	switch predictor {
	case 2:
		decodeSubTIFF(bitsPerComponent, colors, bytesPerPixel, columns, rowLength, actline)

	case 10:
		// PRED NONE — do nothing

	case 11:
		// PRED SUB
		for p := bytesPerPixel; p < rowLength; p++ {
			sub := int(actline[p])
			left := int(actline[p-bytesPerPixel])
			actline[p] = byte(sub + left)
		}

	case 12:
		// PRED UP
		for p := 0; p < rowLength; p++ {
			up := int(actline[p]) & 0xff
			prior := int(lastline[p]) & 0xff
			actline[p] = byte((up + prior) & 0xff)
		}

	case 13:
		// PRED AVG
		for p := 0; p < rowLength; p++ {
			avg := int(actline[p]) & 0xff
			left := 0
			if p-bytesPerPixel >= 0 {
				left = int(actline[p-bytesPerPixel]) & 0xff
			}
			up := int(lastline[p]) & 0xff
			actline[p] = byte((avg + (left+up)/2) & 0xff)
		}

	case 14:
		// PRED PAETH
		for p := 0; p < rowLength; p++ {
			paeth := int(actline[p]) & 0xff
			a := 0
			if p-bytesPerPixel >= 0 {
				a = int(actline[p-bytesPerPixel]) & 0xff
			}
			b := int(lastline[p]) & 0xff
			c := 0
			if p-bytesPerPixel >= 0 {
				c = int(lastline[p-bytesPerPixel]) & 0xff
			}

			value := a + b - c
			absa := absInt(value - a)
			absb := absInt(value - b)
			absc := absInt(value - c)

			chosen := c
			if absa <= absb && absa <= absc {
				chosen = a
			} else if absb <= absc {
				chosen = b
			}

			actline[p] = byte((paeth + chosen) & 0xff)
		}
	}
}

// decodeSubTIFF handles predictor value 2 (TIFF-style SUB prediction).
func decodeSubTIFF(bitsPerComponent, colors, bytesPerPixel, columns, rowLength int, actline []byte) {
	if bitsPerComponent == 8 {
		for p := bytesPerPixel; p < rowLength; p++ {
			sub := int(actline[p]) & 0xff
			left := int(actline[p-bytesPerPixel]) & 0xff
			actline[p] = byte(sub + left)
		}
	} else if bitsPerComponent == 16 {
		for p := bytesPerPixel; p < rowLength-1; p += 2 {
			sub := (int(actline[p])&0xff)<<8 | int(actline[p+1])&0xff
			left := (int(actline[p-bytesPerPixel])&0xff)<<8 | int(actline[p-bytesPerPixel+1])&0xff
			sum := sub + left
			actline[p] = byte((sum >> 8) & 0xff)
			actline[p+1] = byte(sum & 0xff)
		}
	} else if bitsPerComponent == 1 && colors == 1 {
		for p := 0; p < rowLength; p++ {
			for bit := 7; bit >= 0; bit-- {
				sub := (int(actline[p]) >> bit) & 1

				if p == 0 && bit == 7 {
					continue
				}

				left := 0
				if bit == 7 {
					left = int(actline[p-1]) & 1
				} else {
					left = (int(actline[p]) >> (bit + 1)) & 1
				}

				if (sub+left)&1 == 0 {
					actline[p] &= byte(^(1 << bit))
				} else {
					actline[p] |= byte(1 << bit)
				}
			}
		}
	} else {
		elements := columns * colors
		for p := colors; p < elements; p++ {
			bytePosSub := p*bitsPerComponent / 8
			bitPosSub := 8 - p*bitsPerComponent%8 - bitsPerComponent
			bytePosLeft := (p - colors) * bitsPerComponent / 8
			bitPosLeft := 8 - (p-colors)*bitsPerComponent%8 - bitsPerComponent

			sub := getBitSeq(int(actline[bytePosSub]), bitPosSub, bitsPerComponent)
			left := getBitSeq(int(actline[bytePosLeft]), bitPosLeft, bitsPerComponent)
			actline[bytePosSub] = byte(calcSetBitSeq(int(actline[bytePosSub]), bitPosSub, bitsPerComponent, sub+left))
		}
	}
}

// CalculateRowLength returns the number of bytes in a single row given the
// color components, bits per component, and column count.
func CalculateRowLength(colors, bitsPerComponent, columns int) int {
	bitsPerPixel := colors * bitsPerComponent
	return (columns*bitsPerPixel + 7) / 8
}

// getBitSeq extracts a bit-sized value starting at startBit from the given byte.
func getBitSeq(by, startBit, bitSize int) int {
	mask := (1 << bitSize) - 1
	return (by >> startBit) & mask
}

// calcSetBitSeq sets a bit-sized value at startBit position in the given byte
// and returns the modified byte.
func calcSetBitSeq(by, startBit, bitSize, val int) int {
	mask := (1 << bitSize) - 1
	truncatedVal := val & mask
	mask = ^(mask << startBit)
	return (by & mask) | (truncatedVal << startBit)
}

// absInt returns the absolute value of an integer.
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// WrapPredictor wraps an io.Writer with a predictor-decoding writer if
// predictor > 1. When predictor is 1 or less the original writer is returned
// unchanged. The returned writer also exposes a Flush() error method for
// draining any buffered partial rows.
func WrapPredictor(out io.Writer, predictor, colors, bitsPerComponent, columns int) io.Writer {
	if predictor > 1 {
		return NewPredictorOutputStream(out, predictor, colors, bitsPerComponent, columns)
	}
	return out
}

// PredictorOutputStream buffers incoming data until a complete row is available,
// decodes it using PNG-style prediction, then writes the result to the base
// stream. The previous decoded row is retained for use in decoding subsequent
// rows. Implements io.Writer and Flush() error.
type PredictorOutputStream struct {
	baseStream        io.Writer
	predictor         int
	colors            int
	bitsPerComponent  int
	columns           int
	rowLength         int
	predictorPerRow   bool
	currentRow        []byte
	lastRow           []byte
	currentRowData    int
	predictorRead     bool
}

// NewPredictorOutputStream creates a new PredictorOutputStream that decodes
// data using the given PNG prediction parameters before writing to baseStream.
func NewPredictorOutputStream(baseStream io.Writer, predictor, colors, bitsPerComponent, columns int) *PredictorOutputStream {
	rowLength := CalculateRowLength(colors, bitsPerComponent, columns)
	return &PredictorOutputStream{
		baseStream:       baseStream,
		predictor:        predictor,
		colors:           colors,
		bitsPerComponent: bitsPerComponent,
		columns:          columns,
		rowLength:        rowLength,
		predictorPerRow:  predictor >= 10,
		currentRow:       make([]byte, rowLength),
		lastRow:          make([]byte, rowLength),
	}
}

// Write buffers input bytes and decodes complete rows. Returns the number of
// bytes consumed from buffer (which is always len(buffer)).
func (p *PredictorOutputStream) Write(buffer []byte) (int, error) {
	currentOffset := 0
	totalLen := len(buffer)

	for currentOffset < totalLen {
		if p.predictorPerRow && p.currentRowData == 0 && !p.predictorRead {
			p.predictor = int(buffer[currentOffset]) + 10
			currentOffset++
			p.predictorRead = true
		} else {
			toRead := p.rowLength - p.currentRowData
			if remaining := totalLen - currentOffset; remaining < toRead {
				toRead = remaining
			}
			copy(p.currentRow[p.currentRowData:], buffer[currentOffset:currentOffset+toRead])
			p.currentRowData += toRead
			currentOffset += toRead

			if p.currentRowData == len(p.currentRow) {
				p.decodeAndWriteRow()
			}
		}
	}

	return totalLen, nil
}

// Flush finalizes any partially filled row by padding with zeros, decoding it,
// and writing to the underlying stream.
func (p *PredictorOutputStream) Flush() error {
	if p.currentRowData > 0 {
		for i := p.currentRowData; i < p.rowLength; i++ {
			p.currentRow[i] = 0
		}
		p.decodeAndWriteRow()
	}
	return nil
}

func (p *PredictorOutputStream) decodeAndWriteRow() {
	DecodePredictorRow(p.predictor, p.colors, p.bitsPerComponent, p.columns, p.currentRow, p.lastRow)
	p.baseStream.Write(p.currentRow)
	p.flipRows()
}

func (p *PredictorOutputStream) flipRows() {
	p.lastRow, p.currentRow = p.currentRow, p.lastRow
	p.currentRowData = 0
	p.predictorRead = false
}
