package ccitffax

// CcittFaxCompressionType specifies the compression type to use with CCITT fax
// decoder streams as defined in PDF reference Table 7.1.
type CcittFaxCompressionType byte

const (
	// ModifiedHuffman represents Modified Huffman (MH) — Group 3 variation (T2).
	ModifiedHuffman CcittFaxCompressionType = iota

	// Group3_1D represents Modified Huffman (MH) — Group 3 (T4).
	Group3_1D

	// Group3_2D represents Modified Read (MR) — Group 3 (T4).
	Group3_2D

	// Group4_2D represents Modified Modified Read (MMR) — Group 4 (T6).
	Group4_2D
)
