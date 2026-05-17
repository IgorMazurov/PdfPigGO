package png

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// OpenPng reads a PNG image from an io.Reader stream. chunkVisitor is optional and may be nil.
func OpenPng(stream io.Reader, chunkVisitor ChunkVisitor) (*Png, error) {
	return OpenWithSettings(stream, &PngOpenerSettings{
		ChunkVisitor: chunkVisitor,
	})
}

// OpenWithSettings reads a PNG image from an io.Reader stream with the given settings.
func OpenWithSettings(stream io.Reader, settings *PngOpenerSettings) (*Png, error) {
	if stream == nil {
		return nil, errors.New("stream must not be nil")
	}

	readSeeker, ok := stream.(io.ReadSeeker)
	if !ok {
		return nil, fmt.Errorf("the provided stream of type %T does not support seeking", stream)
	}

	validHeader := hasValidHeader(readSeeker)
	if !validHeader.IsValid() {
		return nil, fmt.Errorf("the provided stream did not start with the PNG header. Got %v", validHeader)
	}

	crc := make([]byte, 4)
	imageHeader, err := readImageHeader(readSeeker, crc)
	if err != nil {
		return nil, err
	}

	hasEncounteredImageEnd := false
	var palette *Palette

	memStream := &bytes.Buffer{}

	for {
		header, ok := tryReadChunkHeader(readSeeker)
		if !ok {
			break
		}

		if hasEncounteredImageEnd {
			if settings != nil && settings.DisallowTrailingData {
				return nil, fmt.Errorf("found another chunk %v after already reading the IEND chunk", header)
			}
			break
		}

		bytesChunk := make([]byte, header.Length())
		read, err := io.ReadFull(readSeeker, bytesChunk)
		if err != nil {
			return nil, fmt.Errorf("did not read %d bytes for the %v header: %w", header.Length(), header, err)
		}
		if read != len(bytesChunk) {
			return nil, fmt.Errorf("did not read %d bytes for the %v header, only found: %d", header.Length(), header, read)
		}

		if header.IsCritical() {
			switch header.Name() {
			case "PLTE":
				if header.Length()%3 != 0 {
					return nil, fmt.Errorf("palette data must be multiple of 3, got %d", header.Length())
				}

				if (imageHeader.ColorType() & ColorTypePaletteUsed) != 0 {
					palette = NewPalette(bytesChunk)
				}

			case "IDAT":
				memStream.Write(bytesChunk)

			case "IEND":
				hasEncounteredImageEnd = true

			default:
				return nil, fmt.Errorf("encountered critical header %v which was not recognised", header)
			}
		} else {
			switch header.Name() {
			case "tRNS":
				if palette != nil {
					palette.SetAlphaValues(bytesChunk)
				}
			}
		}

		crcBytes := make([]byte, 4)
		read, err = io.ReadFull(readSeeker, crcBytes)
		if err != nil {
			return nil, fmt.Errorf("did not read 4 bytes for the CRC: %w", err)
		}
		if read != 4 {
			return nil, fmt.Errorf("did not read 4 bytes for the CRC, only found: %d", read)
		}

		chunkNameBytes := []byte(header.Name())
		result := Crc32CalculateTwo(chunkNameBytes, bytesChunk)
		crcActual := uint32(crcBytes[0])<<24 | uint32(crcBytes[1])<<16 | uint32(crcBytes[2])<<8 | uint32(crcBytes[3])

		if result != crcActual {
			return nil, fmt.Errorf("CRC calculated %d did not match file %d for chunk: %s", result, crcActual, header.Name())
		}

		if settings != nil && settings.ChunkVisitor != nil {
			settings.ChunkVisitor.Visit(readSeeker, *imageHeader, header, bytesChunk, crcBytes)
		}
	}

	idatData := memStream.Bytes()

	zlibReader, err := zlib.NewReader(bytes.NewReader(idatData))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader for IDAT data: %w", err)
	}

	outputBuf := &bytes.Buffer{}
	_, err = io.Copy(outputBuf, zlibReader)
	zlibReader.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to decompress IDAT data: %w", err)
	}

	bytesOut := outputBuf.Bytes()
	bytesPerPixel, _ := GetBytesAndSamplesPerPixel(imageHeader)

	decoded, err := Decode(bytesOut, imageHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG data: %w", err)
	}

	rawData, err := NewRawPngData(decoded, int(bytesPerPixel), palette, *imageHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to create raw PNG data: %w", err)
	}

	hasAlpha := false
	if palette != nil {
		hasAlpha = palette.HasAlphaValues()
	}

	return NewPng(imageHeader, rawData, hasAlpha)
}

func hasValidHeader(stream io.ReadSeeker) HeaderValidationResult {
	b := make([]byte, 8)
	io.ReadFull(stream, b)
	return NewHeaderValidationResult(b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7])
}

func tryReadChunkHeader(stream io.ReadSeeker) (ChunkHeader, bool) {
	position, _ := stream.Seek(0, io.SeekCurrent)
	headerBytes, ok := TryReadHeaderBytes(stream)
	if !ok {
		return ChunkHeader{}, false
	}

	length := int(binary.BigEndian.Uint32(headerBytes[0:4]))
	name := string(headerBytes[4:8])

	chunkHeader, err := NewChunkHeader(position, length, name)
	if err != nil {
		return ChunkHeader{}, false
	}

	return chunkHeader, true
}

func readImageHeader(stream io.ReadSeeker, crc []byte) (*ImageHeader, error) {
	header, ok := tryReadChunkHeader(stream)
	if !ok {
		return nil, errors.New("the provided stream did not contain a single chunk")
	}

	if header.Name() != "IHDR" {
		return nil, fmt.Errorf("the first chunk was not the IHDR chunk: %v", header)
	}

	if header.Length() != 13 {
		return nil, fmt.Errorf("the first chunk did not have a length of 13 bytes: %v", header)
	}

	ihdrBytes := make([]byte, 13)
	read, err := io.ReadFull(stream, ihdrBytes)
	if err != nil {
		return nil, fmt.Errorf("did not read 13 bytes for the IHDR: %w", err)
	}
	if read != 13 {
		return nil, fmt.Errorf("did not read 13 bytes for the IHDR, only found: %d", read)
	}

	readCrc := make([]byte, 4)
	read, err = io.ReadFull(stream, readCrc)
	if err != nil {
		return nil, fmt.Errorf("did not read 4 bytes for the CRC: %w", err)
	}
	if read != 4 {
		return nil, fmt.Errorf("did not read 4 bytes for the CRC, only found: %d", read)
	}

	width := int(binary.BigEndian.Uint32(ihdrBytes[0:4]))
	height := int(binary.BigEndian.Uint32(ihdrBytes[4:8]))
	bitDepth := ihdrBytes[8]
	colorType := ColorType(ihdrBytes[9])
	compressionMethod := CompressionMethod(ihdrBytes[10])
	filterMethod := FilterMethod(ihdrBytes[11])
	interlaceMethod := InterlaceMethod(ihdrBytes[12])

	imageHeader, err := NewImageHeader(width, height, bitDepth, colorType, compressionMethod, filterMethod, interlaceMethod)
	if err != nil {
		return nil, err
	}

	return imageHeader, nil
}

// Crc32CalculateTwo computes the combined CRC-32 checksum of two byte slices.
func Crc32CalculateTwo(data, data2 []byte) uint32 {
	crc := uint32(0xFFFFFFFF)
	for i := 0; i < len(data); i++ {
		index := (crc ^ uint32(data[i])) & 0xFF
		crc = (crc >> 8) ^ crc32Lookup[index]
	}
	for i := 0; i < len(data2); i++ {
		index := (crc ^ uint32(data2[i])) & 0xFF
		crc = (crc >> 8) ^ crc32Lookup[index]
	}
	return crc ^ 0xFFFFFFFF
}
