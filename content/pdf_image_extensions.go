package content

// NeedsReverseDecode reports whether the image colors need to be reversed
// based on the Decode array and color space. It returns true when the Decode
// array contains at least two elements with values [1, 0], indicating a
// reverse mapping from the default range.
func NeedsReverseDecode(img PdfImage) bool {
	decode := img.Decode()
	return len(decode) >= 2 && decode[0] == 1 && decode[1] == 0
}
