package filters

import (
	"fmt"
	"io"

	"github.com/uglytoad/pdfpig/go/iostream"
)

// CcittFaxDecoderStream decodes CCITT Modified Huffman RLE, Group 3 (T4) and
// Group 4 (T6) fax compression. Ported from Apache PDFBox CCITTFaxDecoderStream.
type CcittFaxDecoderStream struct {
	*iostream.StreamWrapper

	columns int
	decodedRow []byte

	optionByteAligned bool
	compressionType   CcittFaxCompressionType

	decodedLength  int
	decodedPos     int

	changesReferenceRow  []int
	changesCurrentRow    []int
	changesReferenceRowCount int
	changesCurrentRowCount  int

	lastChangingElement int

	buffer    int
	bufferPos int
}

// NewCcittFaxDecoderStream creates a CCITT fax decoder stream. Use this for
// CCITT streams embedded in PDF files that use EncodedByteAlign.
func NewCcittFaxDecoderStream(stream *iostream.StreamWrapper, columns int, compressionType CcittFaxCompressionType, byteAligned bool) *CcittFaxDecoderStream {
	return &CcittFaxDecoderStream{
		StreamWrapper:      stream,
		columns:            columns,
		compressionType:    compressionType,
		optionByteAligned:  byteAligned,
		decodedRow:         make([]byte, (columns+7)/8),
		changesReferenceRow: make([]int, columns+2),
		changesCurrentRow:  make([]int, columns+2),
		buffer:             -1,
		bufferPos:          -1,
	}
}

func (d *CcittFaxDecoderStream) fetch() {
	if d.decodedPos >= d.decodedLength {
		d.decodedLength = 0

		func() {
			defer func() {
				if r := recover(); r != nil {
					if d.decodedLength == 0 {
						d.decodedLength = -1
					} else {
						panic(r)
					}
				}
			}()

			if err := d.decodeRow(); err != nil {
				if d.decodedLength == 0 {
					d.decodedLength = -1
				} else {
					panic(err)
				}
			}
		}()

		d.decodedPos = 0
	}
}

func (d *CcittFaxDecoderStream) decode1D() {
	index := 0
	white := true
	d.changesCurrentRowCount = 0

	for index < d.columns {
		var completeRun int
		if white {
			completeRun = d.decodeRun(WhiteRunTree)
		} else {
			completeRun = d.decodeRun(BlackRunTree)
		}
		index += completeRun
		d.changesCurrentRow[d.changesCurrentRowCount] = index
		d.changesCurrentRowCount++

		white = !white
	}
}

func (d *CcittFaxDecoderStream) decode2D() {
	d.changesReferenceRowCount = d.changesCurrentRowCount
	tmp := d.changesCurrentRow
	d.changesCurrentRow = d.changesReferenceRow
	d.changesReferenceRow = tmp

	white := true
	index := 0
	d.changesCurrentRowCount = 0

	for index < d.columns {
		node := CodeTree.root

		modeFound := false
		for !modeFound {
			bit := d.readBit()
			if bit {
				node = node.right
			} else {
				node = node.left
			}

			if node == nil {
				modeFound = true
				continue
			}

			if !node.isLeaf {
				continue
			}

			switch node.value {
			case valueHMode:
				var runLength int
				if white {
					runLength = d.decodeRun(WhiteRunTree)
				} else {
					runLength = d.decodeRun(BlackRunTree)
				}
				index += runLength
				d.changesCurrentRow[d.changesCurrentRowCount] = index
				d.changesCurrentRowCount++

				if white {
					runLength = d.decodeRun(BlackRunTree)
				} else {
					runLength = d.decodeRun(WhiteRunTree)
				}
				index += runLength
				d.changesCurrentRow[d.changesCurrentRowCount] = index
				d.changesCurrentRowCount++

			case valuePassMode:
				pChangingElement := d.getNextChangingElement(index, white) + 1

				if pChangingElement >= d.changesReferenceRowCount {
					index = d.columns
				} else {
					index = d.changesReferenceRow[pChangingElement]
				}

			default:
				vChangingElement := d.getNextChangingElement(index, white)

				if vChangingElement >= d.changesReferenceRowCount || vChangingElement == -1 {
					index = d.columns + node.value
				} else {
					index = d.changesReferenceRow[vChangingElement] + node.value
				}

				d.changesCurrentRow[d.changesCurrentRowCount] = index
				d.changesCurrentRowCount++
				white = !white
			}

			modeFound = true
		}
	}
}

func (d *CcittFaxDecoderStream) getNextChangingElement(a0 int, white bool) int {
	start := (d.lastChangingElement & 0xFFFF_FFFE)
	if white {
		start += 0
	} else {
		start += 1
	}
	if start > 2 {
		start -= 2
	}

	if a0 == 0 {
		return start
	}

	for i := start; i < d.changesReferenceRowCount; i += 2 {
		if a0 < d.changesReferenceRow[i] {
			d.lastChangingElement = i
			return i
		}
	}

	return -1
}

func (d *CcittFaxDecoderStream) decodeRowType2() {
	if d.optionByteAligned {
		d.resetBuffer()
	}
	d.decode1D()
}

func (d *CcittFaxDecoderStream) decodeRowType4() {
	if d.optionByteAligned {
		d.resetBuffer()
	}

foundEOL:
	for {
		node := EolOnlyTree.root

		for {
			bit := d.readBit()
			if bit {
				node = node.right
			} else {
				node = node.left
			}

			if node == nil {
				continue foundEOL
			}

			if node.isLeaf {
				break foundEOL
			}
		}
	}

	if d.compressionType == Group3_1D || d.readBit() {
		d.decode1D()
	} else {
		d.decode2D()
	}
}

func (d *CcittFaxDecoderStream) decodeRowType6() {
	if d.optionByteAligned {
		d.resetBuffer()
	}
	d.decode2D()
}

func (d *CcittFaxDecoderStream) decodeRow() error {
	switch d.compressionType {
	case ModifiedHuffman:
		d.decodeRowType2()
	case Group3_1D, Group3_2D:
		d.decodeRowType4()
	case Group4_2D:
		d.decodeRowType6()
	default:
		return fmt.Errorf("%v is not a supported compression type", d.compressionType)
	}

	index := 0
	white := true

	d.lastChangingElement = 0
	for i := 0; i <= d.changesCurrentRowCount; i++ {
		nextChange := d.columns

		if i != d.changesCurrentRowCount {
			nextChange = d.changesCurrentRow[i]
		}

		if nextChange > d.columns {
			nextChange = d.columns
		}

		byteIndex := index / 8

		for index%8 != 0 && nextChange-index > 0 {
			if !white {
				d.decodedRow[byteIndex] |= 1 << (7 - index%8)
			}
			index++
		}

		if index%8 == 0 {
			byteIndex = index / 8
			var value byte
			if white {
				value = 0x00
			} else {
				value = 0xff
			}

			for nextChange-index > 7 {
				d.decodedRow[byteIndex] = value
				index += 8
				byteIndex++
			}
		}

		for nextChange-index > 0 {
			if index%8 == 0 {
				d.decodedRow[byteIndex] = 0
			}

			if !white {
				d.decodedRow[byteIndex] |= 1 << (7 - index%8)
			}
			index++
		}

		white = !white
	}

	if index != d.columns {
		return fmt.Errorf("sum of run-lengths does not equal scan line width: %d > %d", index, d.columns)
	}

	d.decodedLength = (index + 7) / 8
	return nil
}

func (d *CcittFaxDecoderStream) decodeRun(tree *tree) int {
	total := 0

	node := tree.root

	for {
		bit := d.readBit()
		if bit {
			node = node.right
		} else {
			node = node.left
		}

		if node == nil {
			return d.columns
		}

		if !node.isLeaf {
			continue
		}

		total += node.value
		if node.value >= 64 {
			node = tree.root
		} else if node.value >= 0 {
			return total
		} else {
			return d.columns
		}
	}
}

func (d *CcittFaxDecoderStream) resetBuffer() {
	d.bufferPos = -1
}

func (d *CcittFaxDecoderStream) readBit() bool {
	if d.bufferPos < 0 || d.bufferPos > 7 {
		var buf [1]byte
		n, err := d.StreamWrapper.Read(buf[:])
		if err != nil || n == 0 {
			d.buffer = -1
		} else {
			d.buffer = int(buf[0])
		}

		if d.buffer == -1 {
			panic(fmt.Errorf("Unexpected end of Huffman RLE stream"))
		}

		d.bufferPos = 0
	}

	isSet := (d.buffer>>(7-d.bufferPos))&1 == 1

	d.bufferPos++

	if d.bufferPos > 7 {
		d.bufferPos = -1
	}

	return isSet
}

// ReadByte reads a single decoded byte from the CCITT fax stream.
func (d *CcittFaxDecoderStream) ReadByte() (byte, error) {
	if d.decodedLength < 0 {
		return 0x0, io.EOF
	}

	if d.decodedPos >= d.decodedLength {
		d.fetch()

		if d.decodedLength < 0 {
			return 0x0, io.EOF
		}
	}

	b := d.decodedRow[d.decodedPos] & 0xff
	d.decodedPos++
	return b, nil
}

// Read reads up to len(b) decoded bytes into the provided slice at offset off.
func (d *CcittFaxDecoderStream) Read(b []byte, off, length int) (int, error) {
	if d.decodedLength < 0 {
		for i := off; i < off+length; i++ {
			b[i] = 0x0
		}
		return length, nil
	}

	if d.decodedPos >= d.decodedLength {
		d.fetch()

		if d.decodedLength < 0 {
			for i := off; i < off+length; i++ {
				b[i] = 0x0
			}
			return length, nil
		}
	}

	read := d.decodedLength - d.decodedPos
	if read > length {
		read = length
	}

	copy(b[off:], d.decodedRow[d.decodedPos:d.decodedPos+read])
	d.decodedPos += read

	if d.decodedLength < 0 && read == 0 {
		return 0, io.EOF
	}

	return read, nil
}

// --- Huffman code tree structures ---

type node struct {
	left     *node
	right    *node
	value    int
	canBeFill bool
	isLeaf   bool
}

func (n *node) set(next bool, child *node) {
	if !next {
		n.left = child
	} else {
		n.right = child
	}
}

type tree struct {
	root *node
}

func newTree() *tree {
	return &tree{root: &node{}}
}

// fillValue fills the tree with a numeric leaf value at the given depth/path.
func (t *tree) fillValue(depth, path, value int) {
	current := t.root

	for i := 0; i < depth; i++ {
		bitPos := depth - 1 - i
		isSet := (path>>bitPos)&1 == 1
		var next *node
		if isSet {
			next = current.right
		} else {
			next = current.left
		}

		if next == nil {
			next = &node{}

			if i == depth-1 {
				next.value = value
				next.isLeaf = true
			}

			if path == 0 {
				next.canBeFill = true
			}

			current.set(isSet, next)
		} else if next.isLeaf {
			return
		}

		current = next
	}
}

// fillNode fills the tree with an existing node at the given depth/path.
func (t *tree) fillNode(depth, path int, child *node) {
	current := t.root

	for i := 0; i < depth; i++ {
		bitPos := depth - 1 - i
		isSet := (path>>bitPos)&1 == 1
		var next *node
		if isSet {
			next = current.right
		} else {
			next = current.left
		}

		if next == nil {
			if i == depth-1 {
				next = child
			} else {
				next = &node{}
			}

			if path == 0 {
				next.canBeFill = true
			}

			current.set(isSet, next)
		} else if next.isLeaf {
			return
		}

		current = next
	}
}

// --- Static Huffman code tables (TIFF 6.0 Specification, Section 10) ---

var blackCodes = [][]int{
	{0x2, 0x3},
	{0x2, 0x3},
	{0x2, 0x3},
	{0x3},
	{0x4, 0x5},
	{0x4, 0x5, 0x7},
	{0x4, 0x7},
	{0x18},
	{0x17, 0x18, 0x37, 0x8, 0xf},
	{0x17, 0x18, 0x28, 0x37, 0x67, 0x68, 0x6c, 0x8, 0xc, 0xd},
	{0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x1c, 0x1d, 0x1e, 0x1f, 0x24, 0x27, 0x28, 0x2b, 0x2c, 0x33,
		0x34, 0x35, 0x37, 0x38, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x64, 0x65,
		0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xd2, 0xd3,
		0xd4, 0xd5, 0xd6, 0xd7, 0xda, 0xdb},
	{0x4a, 0x4b, 0x4c, 0x4d, 0x52, 0x53, 0x54, 0x55, 0x5a, 0x5b, 0x64, 0x65, 0x6c, 0x6d, 0x72, 0x73,
		0x74, 0x75, 0x76, 0x77},
}

var blackRunLengths = [][]int{
	{3, 2},
	{1, 4},
	{6, 5},
	{7},
	{9, 8},
	{10, 11, 12},
	{13, 14},
	{15},
	{16, 17, 0, 18, 64},
	{24, 25, 23, 22, 19, 20, 21, 1792, 1856, 1920},
	{1984, 2048, 2112, 2176, 2240, 2304, 2368, 2432, 2496, 2560, 52, 55, 56, 59, 60, 320, 384, 448, 53,
		54, 50, 51, 44, 45, 46, 47, 57, 58, 61, 256, 48, 49, 62, 63, 30, 31, 32, 33, 40, 41, 128, 192, 26,
		27, 28, 29, 34, 35, 36, 37, 38, 39, 42, 43},
	{640, 704, 768, 832, 1280, 1344, 1408, 1472, 1536, 1600, 1664, 1728, 512, 576, 896, 960, 1024, 1088,
		1152, 1216},
}

var whiteCodes = [][]int{
	{0x7, 0x8, 0xb, 0xc, 0xe, 0xf},
	{0x12, 0x13, 0x14, 0x1b, 0x7, 0x8},
	{0x17, 0x18, 0x2a, 0x2b, 0x3, 0x34, 0x35, 0x7, 0x8},
	{0x13, 0x17, 0x18, 0x24, 0x27, 0x28, 0x2b, 0x3, 0x37, 0x4, 0x8, 0xc},
	{0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x1a, 0x1b, 0x2, 0x24, 0x25, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d,
		0x3, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x4, 0x4a, 0x4b, 0x5, 0x52, 0x53, 0x54, 0x55, 0x58, 0x59,
		0x5a, 0x5b, 0x64, 0x65, 0x67, 0x68, 0xa, 0xb},
	{0x98, 0x99, 0x9a, 0x9b, 0xcc, 0xcd, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb},
	{},
	{0x8, 0xc, 0xd},
	{0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x1c, 0x1d, 0x1e, 0x1f},
}

var whiteRunLengths = [][]int{
	{2, 3, 4, 5, 6, 7},
	{128, 8, 9, 64, 10, 11},
	{192, 1664, 16, 17, 13, 14, 15, 1, 12},
	{26, 21, 28, 27, 18, 24, 25, 22, 256, 23, 20, 19},
	{33, 34, 35, 36, 37, 38, 31, 32, 29, 53, 54, 39, 40, 41, 42, 43, 44, 30, 61, 62, 63, 0, 320, 384, 45,
		59, 60, 46, 49, 50, 51, 52, 55, 56, 57, 58, 448, 512, 640, 576, 47, 48},
	{1472, 1536, 1600, 1728, 704, 768, 832, 896, 960, 1024, 1088, 1152, 1216, 1280, 1344, 1408},
	{},
	{1792, 1856, 1920},
	{1984, 2048, 2112, 2176, 2240, 2304, 2368, 2432, 2496, 2560},
}

var (
	eolNode     *node
	fillNode    *node
	BlackRunTree *tree
	WhiteRunTree *tree
	EolOnlyTree  *tree
	CodeTree     *tree
)

const (
	valueEOL      = -2000
	valueFill     = -1000
	valuePassMode = -3000
	valueHMode    = -4000
)

func init() {
	eolNode = &node{
		isLeaf: true,
		value:  valueEOL,
	}

	fillNode = &node{
		value: valueFill,
	}
	fillNode.left = fillNode
	fillNode.right = eolNode

	EolOnlyTree = newTree()
	EolOnlyTree.fillNode(12, 0, fillNode)
	EolOnlyTree.fillNode(12, 1, eolNode)

	BlackRunTree = newTree()
	for i := 0; i < len(blackCodes); i++ {
		for j := 0; j < len(blackCodes[i]); j++ {
			BlackRunTree.fillValue(i+2, blackCodes[i][j], blackRunLengths[i][j])
		}
	}
	BlackRunTree.fillNode(12, 0, fillNode)
	BlackRunTree.fillNode(12, 1, eolNode)

	WhiteRunTree = newTree()
	for i := 0; i < len(whiteCodes); i++ {
		for j := 0; j < len(whiteCodes[i]); j++ {
			WhiteRunTree.fillValue(i+4, whiteCodes[i][j], whiteRunLengths[i][j])
		}
	}
	WhiteRunTree.fillNode(12, 0, fillNode)
	WhiteRunTree.fillNode(12, 1, eolNode)

	CodeTree = newTree()
	CodeTree.fillValue(4, 1, valuePassMode) // pass mode
	CodeTree.fillValue(3, 1, valueHMode)    // H mode
	CodeTree.fillValue(1, 1, 0)             // V(0)
	CodeTree.fillValue(3, 3, 1)             // V_R(1)
	CodeTree.fillValue(6, 3, 2)             // V_R(2)
	CodeTree.fillValue(7, 3, 3)             // V_R(3)
	CodeTree.fillValue(3, 2, -1)            // V_L(1)
	CodeTree.fillValue(6, 2, -2)            // V_L(2)
	CodeTree.fillValue(7, 2, -3)            // V_L(3)
}
