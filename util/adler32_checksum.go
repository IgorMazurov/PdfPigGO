package util

// adlerModulus is the prime modulus used in Adler-32 checksum calculation (RFC 1950).
const adlerModulus = 65521

// Calculate computes the Adler-32 checksum for the given data,
// as defined by RFC 1950: ZLIB Compressed Data Format Specification.
func Calculate(data []byte) int {
	s1 := 1
	s2 := 0

	for _, b := range data {
		s1 = (s1 + int(b)) % adlerModulus
		s2 = (s1 + s2) % adlerModulus
	}

	return s2*65536 + s1
}
