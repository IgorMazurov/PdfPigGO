package images

import (
	"errors"
	"io"
)

// jpegReader combines read, seek, and byte-read capabilities needed for JPEG parsing.
type jpegReader interface {
	io.Reader
	io.Seeker
	io.ByteReader
}

var (
	ErrNotJpeg               = errors.New("the input stream did not start with the expected JPEG header [ 255 216 ]")
	ErrJpegDimensionsMissing = errors.New("file was a valid JPEG but the width and height could not be determined")
	ErrInvalidMarker         = errors.New("invalid JPEG marker sequence")
	ErrReadShortFailed       = errors.New("failed to read a short where expected in the JPEG stream")
)

const (
	markerStart    byte = 255
	startOfImage   byte = 216
)

// GetInformation reads JPEG header information from the provided stream.
func GetInformation(stream jpegReader) (*JpegInformation, error) {
	if stream == nil {
		return nil, errors.New("stream cannot be nil")
	}

	if !hasRecognizedHeader(stream) {
		return nil, ErrNotJpeg
	}

	marker := JpegMarker(startOfImage)
	shortBuf := make([]byte, 2)

	for marker != EndOfImage {
		switch marker {
		case StartOfImage, Restart0, Restart1, Restart2, Restart3, Restart4, Restart5, Restart6, Restart7:
			// No length markers
		case StartOfBaselineDctFrame, StartOfProgressiveDctFrame:
			if _, err := io.ReadFull(stream, shortBuf); err != nil {
				return nil, err
			}
			_ = binaryU16(shortBuf) // length

			bppByte, err := stream.ReadByte()
			if err != nil {
				return nil, err
			}
			bpp := int(bppByte)

			if _, err := io.ReadFull(stream, shortBuf); err != nil {
				return nil, err
			}
			height := int(binaryU16(shortBuf))

			if _, err := io.ReadFull(stream, shortBuf); err != nil {
				return nil, err
			}
			width := int(binaryU16(shortBuf))

			numCompByte, err := stream.ReadByte()
			if err != nil {
				return nil, err
			}
			numComp := int(numCompByte)

			return NewJpegInformation(width, height, bpp, numComp), nil

		default:
			if _, err := io.ReadFull(stream, shortBuf); err != nil {
				return nil, err
			}
			length := int(binaryU16(shortBuf))
			if _, err := stream.Seek(int64(length-2), io.SeekCurrent); err != nil {
				return nil, err
			}
		}

		nextMarker, err := readSegmentMarker(stream, true)
		if err != nil {
			return nil, err
		}
		marker = JpegMarker(nextMarker)
	}

	return nil, ErrJpegDimensionsMissing
}

func hasRecognizedHeader(stream jpegReader) bool {
	buf := make([]byte, 2)
	n, _ := io.ReadAtLeast(stream, buf, 2)
	if n != 2 {
		return false
	}
	return buf[0] == markerStart && buf[1] == startOfImage
}

func readSegmentMarker(stream jpegReader, skipData bool) (byte, error) {
	var previous *byte

	for {
		b, err := stream.ReadByte()
		if err != nil {
			return 0, ErrInvalidMarker
		}

		if !skipData {
			if previous == nil && b != markerStart {
				return 0, ErrInvalidMarker
			}

			if b != markerStart {
				return b, nil
			}
		}

		if previous != nil && *previous == markerStart && b != markerStart {
			return b, nil
		}

		prev := b
		previous = &prev
	}
}

func binaryU16(buf []byte) uint16 {
	return uint16(buf[0])<<8 | uint16(buf[1])
}
