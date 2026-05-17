package images

// JpegMarker represents a JPEG marker code.
type JpegMarker byte

const (
	// StartOfBaselineDctFrame indicates that this is a baseline DCT-based JPEG,
	// and specifies the width, height, number of components, and component subsampling.
	StartOfBaselineDctFrame JpegMarker = 0xC0

	// StartOfProgressiveDctFrame indicates that this is a progressive DCT-based JPEG,
	// and specifies the width, height, number of components, and component subsampling.
	StartOfProgressiveDctFrame JpegMarker = 0xC2

	// DefineHuffmanTable specifies one or more Huffman tables.
	DefineHuffmanTable JpegMarker = 0xC4

	// StartOfScan begins a top-to-bottom scan of the image. In baseline images,
	// there is generally a single scan. Progressive images usually contain multiple scans.
	StartOfScan JpegMarker = 0xDA

	// DefineQuantizationTable specifies one or more quantization tables.
	DefineQuantizationTable JpegMarker = 0xDB

	// DefineRestartInterval specifies the interval between RSTn markers,
	// in Minimum Coded Units (MCUs). This marker is followed by two bytes
	// indicating the fixed size so it can be treated like any other variable size segment.
	DefineRestartInterval JpegMarker = 0xDD

	// Restart0 is inserted every r macroblocks.
	Restart0 JpegMarker = 0xD0

	// Restart1 is inserted every r macroblocks.
	Restart1 JpegMarker = 0xD1

	// Restart2 is inserted every r macroblocks.
	Restart2 JpegMarker = 0xD2

	// Restart3 is inserted every r macroblocks.
	Restart3 JpegMarker = 0xD3

	// Restart4 is inserted every r macroblocks.
	Restart4 JpegMarker = 0xD4

	// Restart5 is inserted every r macroblocks.
	Restart5 JpegMarker = 0xD5

	// Restart6 is inserted every r macroblocks.
	Restart6 JpegMarker = 0xD6

	// Restart7 is inserted every r macroblocks.
	Restart7 JpegMarker = 0xD7

	// StartOfImage marks the start of a JPEG image file.
	StartOfImage JpegMarker = 0xD8

	// EndOfImage marks the end of a JPEG image file.
	EndOfImage JpegMarker = 0xD9

	// ApplicationSpecific0 is an application-specific marker (APP0).
	ApplicationSpecific0 JpegMarker = 0xE0

	// ApplicationSpecific1 is an application-specific marker (APP1).
	ApplicationSpecific1 JpegMarker = 0xE1

	// ApplicationSpecific2 is an application-specific marker (APP2).
	ApplicationSpecific2 JpegMarker = 0xE2

	// ApplicationSpecific3 is an application-specific marker (APP3).
	ApplicationSpecific3 JpegMarker = 0xE3

	// ApplicationSpecific4 is an application-specific marker (APP4).
	ApplicationSpecific4 JpegMarker = 0xE4

	// ApplicationSpecific5 is an application-specific marker (APP5).
	ApplicationSpecific5 JpegMarker = 0xE5

	// ApplicationSpecific6 is an application-specific marker (APP6).
	ApplicationSpecific6 JpegMarker = 0xE6

	// ApplicationSpecific7 is an application-specific marker (APP7).
	ApplicationSpecific7 JpegMarker = 0xE7

	// ApplicationSpecific8 is an application-specific marker (APP8).
	ApplicationSpecific8 JpegMarker = 0xE8

	// ApplicationSpecific9 is an application-specific marker (APP9).
	ApplicationSpecific9 JpegMarker = 0xE9

	// ApplicationSpecific10 is an application-specific marker (APP10).
	ApplicationSpecific10 JpegMarker = 0xEA

	// ApplicationSpecific11 is an application-specific marker (APP11).
	ApplicationSpecific11 JpegMarker = 0xEB

	// ApplicationSpecific12 is an application-specific marker (APP12).
	ApplicationSpecific12 JpegMarker = 0xEC

	// ApplicationSpecific13 is an application-specific marker (APP13).
	ApplicationSpecific13 JpegMarker = 0xED

	// ApplicationSpecific14 is an application-specific marker (APP14).
	ApplicationSpecific14 JpegMarker = 0xEE

	// ApplicationSpecific15 is an application-specific marker (APP15).
	ApplicationSpecific15 JpegMarker = 0xEF

	// Comment marks a text comment.
	Comment JpegMarker = 0xFE
)
