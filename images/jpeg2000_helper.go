package images

import (
	"encoding/binary"
	"errors"
)

// GetBitsPerComponent gets bits per component values for Jp2 (Jpx) encoded images (first component).
func GetBitsPerComponent(jp2Bytes []byte) (byte, error) {
	if len(jp2Bytes) < 12 {
		return 0, errors.New("input is too short to be a valid JP2 file")
	}

	length := binary.BigEndian.Uint32(jp2Bytes[0:4])
	if length == 0xFF4FFF51 {
		return parseCodestream(jp2Bytes)
	}

	type_ := binary.BigEndian.Uint32(jp2Bytes[4:8])
	magic := binary.BigEndian.Uint32(jp2Bytes[8:12])
	if length == 0x0000000C && type_ == 0x6A502020 && magic == 0x0D0A870A {
		return parseBoxes(jp2Bytes[12:])
	}

	return 0, errors.New("invalid JP2 or J2K signature")
}

func parseBoxes(jp2Bytes []byte) (byte, error) {
	offset := 0
	for offset < len(jp2Bytes) {
		if offset+8 > len(jp2Bytes) {
			return 0, errors.New("invalid JP2 or J2K box structure")
		}

		boxLength := binary.BigEndian.Uint32(jp2Bytes[offset : offset+4])
		boxType := binary.BigEndian.Uint32(jp2Bytes[offset+4 : offset+8])

		if boxType == 0x6A703263 {
			return parseCodestream(jp2Bytes[offset+8:])
		}

		step := int(boxLength)
		if boxLength == 0 {
			step = 8
		}
		offset += step
	}

	return 0, errors.New("codestream box not found in JP2 or J2K file")
}

func parseCodestream(codestream []byte) (byte, error) {
	offset := 0
	for offset+2 <= len(codestream) {
		marker := binary.BigEndian.Uint16(codestream[offset : offset+2])

		if marker == 0xFF51 {
			if offset+38 > len(codestream) {
				return 0, errors.New("invalid SIZ marker structure")
			}

			offset += 38

			numComponents := binary.BigEndian.Uint16(codestream[offset : offset+2])
			offset += 2

			if numComponents < 1 {
				return 0, errors.New("invalid number of components in SIZ marker")
			}

			bitsPerComponent := codestream[offset]
			return bitsPerComponent + 1, nil
		}

		offset += 2
	}

	return 0, errors.New("SIZ marker not found in JPEG2000 codestream")
}
