package page_segmenter

import (
	"math"
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	document_layout_analysis "github.com/uglytoad/pdfpig/go/document_layout_analysis"
	"github.com/uglytoad/pdfpig/go/geometry"
)

// RecursiveXYCut implements the recursive X-Y cut page segmentation algorithm.
// It decomposes a document recursively into rectangular blocks using bounding boxes.
// See "Recursive X-Y Cut using Bounding Boxes of Connected Components" by
// Jaekyu Ha, Robert M. Haralick and Ihsin T. Phillips.
type RecursiveXYCut struct {
	options RecursiveXYCutOptions
}

// DefaultRecursiveXYCutInstance is a pre-created instance with default options.
var DefaultRecursiveXYCutInstance = NewRecursiveXYCut(DefaultRecursiveXYCutOptions())

// NewRecursiveXYCut creates a new RecursiveXYCut with the given options.
func NewRecursiveXYCut(options RecursiveXYCutOptions) *RecursiveXYCut {
	return &RecursiveXYCut{options: options}
}

// GetBlocks returns text blocks generated from the given words on a page.
func (r *RecursiveXYCut) GetBlocks(words []*content.Word) []*document_layout_analysis.TextBlock {
	filtered := make([]*content.Word, 0, len(words))
	for _, w := range words {
		if w != nil && strings.TrimSpace(w.Text) != "" {
			filtered = append(filtered, w)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	leaf, err := NewXYLeafs(filtered)
	if err != nil {
		return nil
	}

	result := r.verticalCut(leaf, r.options.MinimumWidth(),
		r.options.DominantFontWidthFunc(), r.options.DominantFontHeightFunc(), 0)

	if result.isLeaf() {
		lines := result.leaf.GetLines(r.options.WordSeparator())
		tb, err := document_layout_analysis.NewTextBlock(lines, r.options.LineSeparator())
		if err != nil {
			return nil
		}
		return []*document_layout_analysis.TextBlock{tb}
	}

	leaves := result.getLeaves()
	blocks := make([]*document_layout_analysis.TextBlock, 0, len(leaves))
	for _, l := range leaves {
		lines := l.GetLines(r.options.WordSeparator())
		tb, err := document_layout_analysis.NewTextBlock(lines, r.options.LineSeparator())
		if err != nil {
			continue
		}
		blocks = append(blocks, tb)
	}

	return blocks
}

// xyCutNode is an internal node in the recursive X-Y cut tree.
type xyCutNode struct {
	leaf     *XYLeaf
	bbox     core.PdfRectangle
	children []*xyCutNode
}

func (n *xyCutNode) isLeaf() bool {
	return n.leaf != nil
}

// getLeaves recursively collects leaf nodes in reading order.
func (n *xyCutNode) getLeaves() []*XYLeaf {
	if n.isLeaf() {
		return []*XYLeaf{n.leaf}
	}

	var leaves []*XYLeaf
	n.collectLeaves(0, &leaves)
	return leaves
}

func (n *xyCutNode) collectLeaves(level int, result *[]*XYLeaf) {
	if n.children == nil {
		return
	}

	for _, child := range n.children {
		if child.isLeaf() {
			*result = append(*result, child.leaf)
		}
	}

	level++
	isVerticalCut := (level-1)%2 == 0

	var notLeaves []*xyCutNode
	for _, child := range n.children {
		if !child.isLeaf() {
			notLeaves = append(notLeaves, child)
		}
	}

	if isVerticalCut {
		sortByLeftAscCut(notLeaves)
	} else {
		sortByTopDescCut(notLeaves)
	}

	for _, node := range notLeaves {
		node.collectLeaves(level, result)
	}
}

func sortByLeftAscCut(nodes []*xyCutNode) {
	for i := 1; i < len(nodes); i++ {
		for j := 0; j < len(nodes)-i; j++ {
			if nodes[j].bbox.Left() > nodes[j+1].bbox.Left() {
				nodes[j], nodes[j+1] = nodes[j+1], nodes[j]
			}
		}
	}
}

func sortByTopDescCut(nodes []*xyCutNode) {
	for i := 1; i < len(nodes); i++ {
		for j := 0; j < len(nodes)-i; j++ {
			if nodes[j].bbox.Top() < nodes[j+1].bbox.Top() {
				nodes[j], nodes[j+1] = nodes[j+1], nodes[j]
			}
		}
	}
}

// projection represents a range on one axis used during the cut algorithm.
type projection struct {
	lowerBound float64
	upperBound float64
}

func (p *projection) contains(value float64) bool {
	return value >= p.lowerBound && value <= p.upperBound
}

func (r *RecursiveXYCut) verticalCut(leaf *XYLeaf, minimumWidth float64,
	dominantFontWidthFunc func([]*content.Letter) float64,
	dominantFontHeightFunc func([]*content.Letter) float64, level int) *xyCutNode {

	words := make([]*content.Word, len(leaf.Words))
	copy(words, leaf.Words)

	sortByLeftAscWords(words)

	if len(words) == 0 {
		return &xyCutNode{}
	}

	leaf, _ = NewXYLeafs(words)

	if leaf.CountWords() <= 1 || leaf.BoundingBox.Width <= minimumWidth {
		return &xyCutNode{leaf: leaf, bbox: leaf.BoundingBox}
	}

	allLetters := collectLetters(words)
	dominantFontWidth := dominantFontWidthFunc(allLetters)

	projectionProfile := make([]projection, 0)

	firstWordBound := geometry.RectangleNormalise(words[0].BoundingBox())
	currentProj := projection{firstWordBound.Left(), firstWordBound.Right()}
	wordsCount := len(words)

	for i := 1; i < wordsCount; i++ {
		currentWordBound := geometry.RectangleNormalise(words[i].BoundingBox())

		if currentProj.contains(currentWordBound.Left()) || currentProj.contains(currentWordBound.Right()) {
			if currentWordBound.Left() >= currentProj.lowerBound &&
				currentWordBound.Left() <= currentProj.upperBound &&
				currentWordBound.Right() > currentProj.upperBound {
				currentProj.upperBound = currentWordBound.Right()
			}
		} else {
			if currentWordBound.Left()-currentProj.upperBound <= dominantFontWidth {
				currentProj.upperBound = currentWordBound.Right()
			} else if currentProj.upperBound-currentProj.lowerBound < minimumWidth {
				currentProj.upperBound = currentWordBound.Right()
			} else {
				if i != wordsCount-1 {
					projectionProfile = append(projectionProfile, currentProj)
					currentProj = projection{currentWordBound.Left(), currentWordBound.Right()}
				}
			}
		}

		if i == wordsCount-1 {
			projectionProfile = append(projectionProfile, currentProj)
		}
	}

	newLeaves := make([]*XYLeaf, 0)
	var allAssignedWords []*content.Word

	for _, p := range projectionProfile {
		var projWords []*content.Word
		for _, w := range leaf.Words {
			normBB := geometry.RectangleNormalise(w.BoundingBox())
			if normBB.Left() >= p.lowerBound && normBB.Right() <= p.upperBound {
				projWords = append(projWords, w)
			}
		}
		allAssignedWords = append(allAssignedWords, projWords...)

		if len(projWords) > 0 {
			l, err := NewXYLeafs(projWords)
			if err == nil {
				newLeaves = append(newLeaves, l)
			}
		}
	}

	newNodes := make([]*xyCutNode, 0, len(newLeaves))
	for _, l := range newLeaves {
		node := r.horizontalCut(l, minimumWidth, dominantFontWidthFunc, dominantFontHeightFunc, level)
		newNodes = append(newNodes, node)
	}

	lost := findLostWords(leaf.Words, allAssignedWords)
	for _, w := range lost {
		l, err := NewXYLeafs([]*content.Word{w})
		if err == nil {
			newNodes = append(newNodes, &xyCutNode{leaf: l, bbox: l.BoundingBox})
		}
	}

	bbox := computeBBoxFromChildren(newNodes)
	return &xyCutNode{children: newNodes, bbox: bbox}
}

func (r *RecursiveXYCut) horizontalCut(leaf *XYLeaf, minimumWidth float64,
	dominantFontWidthFunc func([]*content.Letter) float64,
	dominantFontHeightFunc func([]*content.Letter) float64, level int) *xyCutNode {

	words := make([]*content.Word, len(leaf.Words))
	copy(words, leaf.Words)

	sortByBottomAscWords(words)

	if len(words) == 0 {
		return &xyCutNode{}
	}

	leaf, _ = NewXYLeafs(words)

	if leaf.CountWords() <= 1 {
		return &xyCutNode{leaf: leaf, bbox: leaf.BoundingBox}
	}

	allLetters := collectLetters(words)
	dominantFontHeight := dominantFontHeightFunc(allLetters)

	projectionProfile := make([]projection, 0)

	firstWordBound := geometry.RectangleNormalise(words[0].BoundingBox())
	currentProj := projection{firstWordBound.Bottom(), firstWordBound.Top()}
	wordsCount := len(words)

	for i := 1; i < wordsCount; i++ {
		currentWordBound := geometry.RectangleNormalise(words[i].BoundingBox())

		if currentProj.contains(currentWordBound.Bottom()) || currentProj.contains(currentWordBound.Top()) {
			if currentWordBound.Bottom() >= currentProj.lowerBound &&
				currentWordBound.Bottom() <= currentProj.upperBound &&
				currentWordBound.Top() > currentProj.upperBound {
				currentProj.upperBound = currentWordBound.Top()
			}
		} else {
			if currentWordBound.Bottom()-currentProj.upperBound <= dominantFontHeight {
				currentProj.upperBound = currentWordBound.Top()
			} else {
				if i != wordsCount-1 {
					projectionProfile = append(projectionProfile, currentProj)
					currentProj = projection{currentWordBound.Bottom(), currentWordBound.Top()}
				}
			}
		}

		if i == wordsCount-1 {
			projectionProfile = append(projectionProfile, currentProj)
		}
	}

	if len(projectionProfile) == 1 {
		if level >= 1 {
			return &xyCutNode{leaf: leaf, bbox: leaf.BoundingBox}
		}
		level++
	}

	newLeaves := make([]*XYLeaf, 0)
	var allAssignedWords []*content.Word

	for _, p := range projectionProfile {
		var projWords []*content.Word
		for _, w := range leaf.Words {
			normBB := geometry.RectangleNormalise(w.BoundingBox())
			if normBB.Bottom() >= p.lowerBound && normBB.Top() <= p.upperBound {
				projWords = append(projWords, w)
			}
		}
		allAssignedWords = append(allAssignedWords, projWords...)

		if len(projWords) > 0 {
			l, err := NewXYLeafs(projWords)
			if err == nil {
				newLeaves = append(newLeaves, l)
			}
		}
	}

	newNodes := make([]*xyCutNode, 0, len(newLeaves))
	for _, l := range newLeaves {
		node := r.verticalCut(l, minimumWidth, dominantFontWidthFunc, dominantFontHeightFunc, level)
		newNodes = append(newNodes, node)
	}

	lost := findLostWords(leaf.Words, allAssignedWords)
	for _, w := range lost {
		l, err := NewXYLeafs([]*content.Word{w})
		if err == nil {
			newNodes = append(newNodes, &xyCutNode{leaf: l, bbox: l.BoundingBox})
		}
	}

	bbox := computeBBoxFromChildren(newNodes)
	return &xyCutNode{children: newNodes, bbox: bbox}
}

func collectLetters(words []*content.Word) []*content.Letter {
	var letters []*content.Letter
	for _, w := range words {
		letters = append(letters, w.Letters...)
	}
	return letters
}

func sortByLeftAscWords(words []*content.Word) {
	sortByFunc(words, func(w *content.Word) float64 {
		return geometry.RectangleNormalise(w.BoundingBox()).Left()
	})
}

func sortByBottomAscWords(words []*content.Word) {
	sortByFunc(words, func(w *content.Word) float64 {
		return geometry.RectangleNormalise(w.BoundingBox()).Bottom()
	})
}

func sortByFunc(words []*content.Word, key func(*content.Word) float64) {
	for i := 1; i < len(words); i++ {
		for j := 0; j < len(words)-i; j++ {
			if key(words[j]) > key(words[j+1]) {
				words[j], words[j+1] = words[j+1], words[j]
			}
		}
	}
}

func findLostWords(all, assigned []*content.Word) []*content.Word {
	assignedSet := make(map[*content.Word]bool)
	for _, w := range assigned {
		assignedSet[w] = true
	}

	var lost []*content.Word
	for _, w := range all {
		if !assignedSet[w] && strings.TrimSpace(w.Text) != "" {
			lost = append(lost, w)
		}
	}
	return lost
}

func computeBBoxFromChildren(children []*xyCutNode) core.PdfRectangle {
	if len(children) == 0 {
		return core.PdfRectangle{}
	}

	minLeft := children[0].bbox.Left()
	minBottom := children[0].bbox.Bottom()
	maxRight := children[0].bbox.Right()
	maxTop := children[0].bbox.Top()

	for _, child := range children[1:] {
		left := child.bbox.Left()
		if left < minLeft {
			minLeft = left
		}
		bottom := child.bbox.Bottom()
		if bottom < minBottom {
			minBottom = bottom
		}
		right := child.bbox.Right()
		if right > maxRight {
			maxRight = right
		}
		top := child.bbox.Top()
		if top > maxTop {
			maxTop = top
		}
	}

	return core.NewPdfRectangleFloat(minLeft, minBottom, maxRight, maxTop)
}

// RecursiveXYCutOptions holds configuration for the recursive X-Y cut page segmenter.
type RecursiveXYCutOptions struct {
	maxDegreeOfParallelism int
	wordSeparator          string
	lineSeparator          string
	minimumWidth           float64
	dominantFontWidthFunc  func([]*content.Letter) float64
	dominantFontHeightFunc func([]*content.Letter) float64
}

// DefaultRecursiveXYCutOptions returns options with standard defaults.
func DefaultRecursiveXYCutOptions() RecursiveXYCutOptions {
	return RecursiveXYCutOptions{
		maxDegreeOfParallelism: -1,
		wordSeparator:          " ",
		lineSeparator:          "\n",
		minimumWidth:           1,
		dominantFontWidthFunc: defaultDominantFontWidth,
		dominantFontHeightFunc: defaultDominantFontHeight,
	}
}

// MaxDegreeOfParallelism returns the maximum number of concurrent tasks enabled.
func (o RecursiveXYCutOptions) MaxDegreeOfParallelism() int {
	return o.maxDegreeOfParallelism
}

// SetMaxDegreeOfParallelism sets the maximum number of concurrent tasks enabled.
func (o *RecursiveXYCutOptions) SetMaxDegreeOfParallelism(value int) {
	o.maxDegreeOfParallelism = value
}

// WordSeparator returns the separator used between words when building lines.
func (o RecursiveXYCutOptions) WordSeparator() string {
	return o.wordSeparator
}

// SetWordSeparator sets the separator used between words when building lines.
func (o *RecursiveXYCutOptions) SetWordSeparator(value string) {
	o.wordSeparator = value
}

// LineSeparator returns the separator used between lines when building blocks.
func (o RecursiveXYCutOptions) LineSeparator() string {
	return o.lineSeparator
}

// SetLineSeparator sets the separator used between lines when building blocks.
func (o *RecursiveXYCutOptions) SetLineSeparator(value string) {
	o.lineSeparator = value
}

// MinimumWidth returns the minimum width for a block.
func (o RecursiveXYCutOptions) MinimumWidth() float64 {
	return o.minimumWidth
}

// SetMinimumWidth sets the minimum width for a block.
func (o *RecursiveXYCutOptions) SetMinimumWidth(value float64) {
	o.minimumWidth = value
}

// DominantFontWidthFunc returns the function that determines the dominant font width.
func (o RecursiveXYCutOptions) DominantFontWidthFunc() func([]*content.Letter) float64 {
	if o.dominantFontWidthFunc == nil {
		return defaultDominantFontWidth
	}
	return o.dominantFontWidthFunc
}

// SetDominantFontWidthFunc sets the function that determines the dominant font width.
func (o *RecursiveXYCutOptions) SetDominantFontWidthFunc(value func([]*content.Letter) float64) {
	o.dominantFontWidthFunc = value
}

// DominantFontHeightFunc returns the function that determines the dominant font height.
func (o RecursiveXYCutOptions) DominantFontHeightFunc() func([]*content.Letter) float64 {
	if o.dominantFontHeightFunc == nil {
		return defaultDominantFontHeight
	}
	return o.dominantFontHeightFunc
}

// SetDominantFontHeightFunc sets the function that determines the dominant font height.
func (o *RecursiveXYCutOptions) SetDominantFontHeightFunc(value func([]*content.Letter) float64) {
	o.dominantFontHeightFunc = value
}

// defaultDominantFontWidth computes the mode of letter widths, falling back to average.
func defaultDominantFontWidth(letters []*content.Letter) float64 {
	widths := make([]float64, 0, len(letters))
	for _, l := range letters {
		w := math.Max(math.Round(l.Width*1000)/1000, math.Round(l.BoundingBox.Width*1000)/1000)
		widths = append(widths, w)
	}

	mode := document_layout_analysis.ModeFloat64(widths)
	if math.IsNaN(mode) || mode == 0 {
		var sum float64
		for _, w := range widths {
			sum += w
		}
		mode = sum / float64(len(widths))
	}
	return mode
}

// defaultDominantFontHeight computes the mode of letter heights times 1.5, falling back to average.
func defaultDominantFontHeight(letters []*content.Letter) float64 {
	heights := make([]float64, 0, len(letters))
	for _, l := range letters {
		h := math.Round(l.BoundingBox.Height*1000) / 1000
		heights = append(heights, h)
	}

	mode := document_layout_analysis.ModeFloat64(heights)
	if math.IsNaN(mode) || mode == 0 {
		var sum float64
		for _, h := range heights {
			sum += h
		}
		mode = sum / float64(len(heights))
	}
	return mode * 1.5
}

var _ PageSegmenter = (*RecursiveXYCut)(nil)
var _ PageSegmenterOptions = (*RecursiveXYCutOptions)(nil)
