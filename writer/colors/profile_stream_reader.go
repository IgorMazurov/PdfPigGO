// Package colors provides types for PDF color handling.
package colors

import _ "embed"

//go:embed resources/icc/sRGB2014.icc
var sRgb2014ICC []byte

// ProfileStreamReaderGetSRgb2014 returns the embedded sRGB 2014 ICC profile bytes.
func ProfileStreamReaderGetSRgb2014() []byte {
	return sRgb2014ICC
}
