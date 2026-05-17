package png

import "fmt"

// HeaderValidationResult holds the result of validating a PNG file header.
type HeaderValidationResult struct {
	byte1 int
	byte2 int
	byte3 int
	byte4 int
	byte5 int
	byte6 int
	byte7 int
	byte8 int
	valid bool
}

// expectedHeader is the 8-byte signature that identifies a valid PNG file.
var expectedHeader = []byte{137, 80, 78, 71, 13, 10, 26, 10}

// NewHeaderValidationResult creates a validation result from the first eight bytes of the file.
func NewHeaderValidationResult(b1, b2, b3, b4, b5, b6, b7, b8 byte) HeaderValidationResult {
	return HeaderValidationResult{
		byte1: int(b1),
		byte2: int(b2),
		byte3: int(b3),
		byte4: int(b4),
		byte5: int(b5),
		byte6: int(b6),
		byte7: int(b7),
		byte8: int(b8),
		valid: b1 == expectedHeader[0] && b2 == expectedHeader[1] &&
			b3 == expectedHeader[2] && b4 == expectedHeader[3] &&
			b5 == expectedHeader[4] && b6 == expectedHeader[5] &&
			b7 == expectedHeader[6] && b8 == expectedHeader[7],
	}
}

// Byte1 returns the first byte of the header.
func (r HeaderValidationResult) Byte1() int { return r.byte1 }

// Byte2 returns the second byte of the header.
func (r HeaderValidationResult) Byte2() int { return r.byte2 }

// Byte3 returns the third byte of the header.
func (r HeaderValidationResult) Byte3() int { return r.byte3 }

// Byte4 returns the fourth byte of the header.
func (r HeaderValidationResult) Byte4() int { return r.byte4 }

// Byte5 returns the fifth byte of the header.
func (r HeaderValidationResult) Byte5() int { return r.byte5 }

// Byte6 returns the sixth byte of the header.
func (r HeaderValidationResult) Byte6() int { return r.byte6 }

// Byte7 returns the seventh byte of the header.
func (r HeaderValidationResult) Byte7() int { return r.byte7 }

// Byte8 returns the eighth byte of the header.
func (r HeaderValidationResult) Byte8() int { return r.byte8 }

// IsValid returns true if the eight bytes match the PNG magic number.
func (r HeaderValidationResult) IsValid() bool { return r.valid }

// String returns a space-separated string of the eight header bytes.
func (r HeaderValidationResult) String() string {
	return fmt.Sprintf("%d %d %d %d %d %d %d %d", r.byte1, r.byte2, r.byte3, r.byte4, r.byte5, r.byte6, r.byte7, r.byte8)
}
