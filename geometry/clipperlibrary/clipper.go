// Package clipperlibrary provides polygon clipping operations based on the Clipper library.
//
// Copyright (c) Angus Johnson 2010-2017, modified for PdfPig.
// Licensed under Boost Software License - Version 1.0.
// See: http://www.boost.org/LICENSE_1_0.txt
//
// This is a Go port of the C# Clipper library (itself a translation from Delphi),
// implementing Bala Vatti's clipping algorithm for polygon operations including
// intersection, union, difference, and XOR.
package clipperlibrary

import (
	"math"
	"sort"
)

const (
	ioReverseSolution   = 1
	ioStrictlySimple    = 2
	ioPreserveCollinear = 4
)

// ClipperClipType defines the type of clipping operation.
type ClipperClipType byte

const (
	ClipperIntersection ClipperClipType = iota
	ClipperUnion
	ClipperDifference
	ClipperXor
)

// ClipperDirection indicates horizontal edge direction.
type ClipperDirection byte

const (
	ClipperRightToLeft ClipperDirection = iota
	ClipperLeftToRight
)

// ClipperEdgeSide indicates which side of a solution polygon an edge belongs to.
type ClipperEdgeSide byte

const (
	ClipperLeft ClipperEdgeSide = iota
	ClipperRight
)

// ClipperEndType defines the type of path end for offset operations.
type ClipperEndType byte

const (
	ClipperClosedPolygon ClipperEndType = iota
	ClipperClosedLine
	ClipperOpenButt
	ClipperOpenSquare
	ClipperOpenRound
)

// ClipperJoinType defines the type of join for offset operations.
type ClipperJoinType byte

const (
	ClipperSquare ClipperJoinType = iota
	ClipperRound
	ClipperMiter
)

// ClipperPolyFillType defines the winding rule for polygon filling.
type ClipperPolyFillType byte

const (
	ClipperEvenOdd ClipperPolyFillType = iota
	ClipperNonZero
)

// ClipperPolyType distinguishes subject from clip polygons.
type ClipperPolyType byte

const (
	ClipperSubject ClipperPolyType = iota
	ClipperClip
)

// ClipperIntPoint represents a point with integer coordinates.
type ClipperIntPoint struct {
	X int64
	Y int64
}

func NewClipperIntPoint(x, y int64) ClipperIntPoint {
	return ClipperIntPoint{X: x, Y: y}
}

// ClipperDoublePoint represents a point with floating-point coordinates.
type ClipperDoublePoint struct {
	X float64
	Y float64
}

func NewClipperDoublePoint(x, y float64) ClipperDoublePoint {
	return ClipperDoublePoint{X: x, Y: y}
}

// ClipperIntRect represents an axis-aligned rectangle with integer coordinates.
type ClipperIntRect struct {
	Left   int64
	Top    int64
	Right  int64
	Bottom int64
}

func NewClipperIntRect(left, top, right, bottom int64) ClipperIntRect {
	return ClipperIntRect{Left: left, Top: top, Right: right, Bottom: bottom}
}

// ClipperException is the error type for clipping operations.
type ClipperException struct {
	msg string
}

func (e *ClipperException) Error() string {
	return e.msg
}

func NewClipperException(msg string) *ClipperException {
	return &ClipperException{msg: msg}
}

// ClipperTEdge represents a trading edge in the active edge list.
type ClipperTEdge struct {
	Bot       ClipperIntPoint
	Curr      ClipperIntPoint // current position (updated for every new scanbeam)
	Top       ClipperIntPoint
	Delta     ClipperIntPoint
	Dx        float64
	PolyTyp   ClipperPolyType
	Side      ClipperEdgeSide
	WindDelta int
	WindCnt   int
	WindCnt2  int // winding count of the opposite polytype
	OutIdx    int

	// Linked list pointers
	Next       *ClipperTEdge
	Prev       *ClipperTEdge
	NextInLML  *ClipperTEdge // next in local minima list
	NextInAEL  *ClipperTEdge // next in active edge list
	PrevInAEL  *ClipperTEdge
	NextInSEL  *ClipperTEdge // next in sorted edge list (horizontals)
	PrevInSEL  *ClipperTEdge
}

// ClipperLocalMinima represents a local minimum vertex pair.
type ClipperLocalMinima struct {
	Y          int64
	LeftBound  *ClipperTEdge
	RightBound *ClipperTEdge
	Next       *ClipperLocalMinima
}

// ClipperScanbeam represents a scanline in the scanbeam priority list.
type ClipperScanbeam struct {
	Y    int64
	Next *ClipperScanbeam
}

// ClipperMaxima represents a local maximum X coordinate (doubly linked).
type ClipperMaxima struct {
	X        int64
	Next     *ClipperMaxima
	Previous *ClipperMaxima
}

// ClipperOutPt represents an output polygon point (doubly linked circular list).
type ClipperOutPt struct {
	Index int
	Pt    ClipperIntPoint
	Next  *ClipperOutPt
	Prev  *ClipperOutPt
}

// ClipperOutRec holds a path in the clipping solution.
type ClipperOutRec struct {
	Idx      int
	IsHole   bool
	IsOpen   bool
	FirstLeft *ClipperOutRec // hole state reference
	Pts       *ClipperOutPt  // linked list of output points
	BottomPt  *ClipperOutPt
	PolyNode  *ClipperPolyNode
}

// ClipperJoin represents a join between two output polygon vertices.
type ClipperJoin struct {
	OutPt1 *ClipperOutPt
	OutPt2 *ClipperOutPt
	OffPt  ClipperIntPoint
}

// ClipperIntersectNode represents an edge intersection for sorting/processing.
type ClipperIntersectNode struct {
	Edge1 *ClipperTEdge
	Edge2 *ClipperTEdge
	Pt    ClipperIntPoint
}

// ClipperPolyNode is a node in the polygon tree hierarchy.
type ClipperPolyNode struct {
	Parent   *ClipperPolyNode
	Polygon  []ClipperIntPoint
	Index    int
	JoinType ClipperJoinType
	EndType  ClipperEndType
	Children []*ClipperPolyNode
	IsOpen   bool
}

// IsHole returns true if this node represents a hole.
func (pn *ClipperPolyNode) IsHole() bool {
	result := true
	node := pn.Parent
	for node != nil {
		result = !result
		node = node.Parent
	}
	return result
}

// ChildCount returns the number of child nodes.
func (pn *ClipperPolyNode) ChildCount() int {
	return len(pn.Children)
}

// Contour returns the polygon contour points.
func (pn *ClipperPolyNode) Contour() []ClipperIntPoint {
	return pn.Polygon
}

// AddChild adds a child node to this poly node.
func (pn *ClipperPolyNode) AddChild(child *ClipperPolyNode) {
	cnt := len(pn.Children)
	pn.Children = append(pn.Children, child)
	child.Parent = pn
	child.Index = cnt
}

// GetNext returns the next sibling or parent's next sibling.
func (pn *ClipperPolyNode) GetNext() *ClipperPolyNode {
	if len(pn.Children) > 0 {
		return pn.Children[0]
	}
	return pn.getNextSiblingUp()
}

func (pn *ClipperPolyNode) getNextSiblingUp() *ClipperPolyNode {
	if pn.Parent == nil {
		return nil
	}
	if pn.Index == len(pn.Parent.Children)-1 {
		return pn.Parent.getNextSiblingUp()
	}
	return pn.Parent.Children[pn.Index+1]
}

// ClipperPolyTree is the root of a polygon hierarchy tree.
type ClipperPolyTree struct {
	*ClipperPolyNode
	AllPolys []*ClipperPolyNode
}

func NewClipperPolyTree() *ClipperPolyTree {
	return &ClipperPolyTree{
		ClipperPolyNode: &ClipperPolyNode{},
		AllPolys:        make([]*ClipperPolyNode, 0),
	}
}

// Clear removes all polygons from the tree.
func (pt *ClipperPolyTree) Clear() {
	for i := range pt.AllPolys {
		pt.AllPolys[i] = nil
	}
	pt.AllPolys = pt.AllPolys[:0]
	pt.Children = pt.Children[:0]
}

// GetFirst returns the first child node or nil.
func (pt *ClipperPolyTree) GetFirst() *ClipperPolyNode {
	if len(pt.Children) > 0 {
		return pt.Children[0]
	}
	return nil
}

// Total returns the total number of polygons in the tree.
func (pt *ClipperPolyTree) Total() int {
	result := len(pt.AllPolys)
	if result > 0 && len(pt.Children) > 0 && pt.Children[0] != pt.AllPolys[0] {
		result--
	}
	return result
}

// Base constants and helpers.
const (
	Horizontal = -3.4e+38
	Skip       = -2
	Unassigned = -1
	Tolerance  = 1.0e-20
	loRange    = 0x3FFFFFFF
	hiRange    = 0x3FFFFFFFFFFFFFFF
)

func nearZero(val float64) bool {
	return val > -Tolerance && val < Tolerance
}

// Clipper performs polygon clipping operations using the Vatti algorithm.
type Clipper struct {
	mClipType         ClipperClipType
	mMaxima           *ClipperMaxima
	mSortedEdges      *ClipperTEdge
	mIntersectList    []*ClipperIntersectNode
	executeLocked     bool
	mClipFillType     ClipperPolyFillType
	mSubjFillType     ClipperPolyFillType
	mJoins            []*ClipperJoin
	mGhostJoins       []*ClipperJoin
	mUsingPolyTree    bool

	// Inherited from ClipperBase
	minimaList   *ClipperLocalMinima
	currentLM    *ClipperLocalMinima
	edges        [][]*ClipperTEdge
	scanbeam     *ClipperScanbeam
	polyOuts     []*ClipperOutRec
	activeEdges  *ClipperTEdge
	useFullRange bool
	hasOpenPaths bool

	// Public properties
	ReverseSolution   bool
	StrictlySimple    bool
	PreserveCollinear bool
}

// NewClipper creates a new Clipper instance with optional init flags.
func NewClipper(initOptions int) *Clipper {
	return &Clipper{
		mIntersectList:    make([]*ClipperIntersectNode, 0),
		polyOuts:          make([]*ClipperOutRec, 0),
		mJoins:            make([]*ClipperJoin, 0),
		mGhostJoins:       make([]*ClipperJoin, 0),
		edges:             make([][]*ClipperTEdge, 0),
		ReverseSolution:   (initOptions&ioReverseSolution) != 0,
		StrictlySimple:    (initOptions&ioStrictlySimple) != 0,
		PreserveCollinear: (initOptions&ioPreserveCollinear) != 0,
	}
}

// Clear removes all paths and resets the clipper state.
func (c *Clipper) Clear() {
	c.minimaList = nil
	c.currentLM = nil
	for i := range c.edges {
		for j := range c.edges[i] {
			c.edges[i][j] = nil
		}
		c.edges[i] = c.edges[i][:0]
	}
	c.edges = c.edges[:0]
	c.useFullRange = false
	c.hasOpenPaths = false
}

// Execute performs the clipping operation and returns the solution polygons.
// The solution parameter is accepted for API compatibility but results are returned
// as the function return value since Go slices are passed by value. Callers should
// use: result := clipper.Execute(ct, nil, ft) or ignore the parameter with:
// clipper.Execute(ct, nil, ft).
func (c *Clipper) Execute(clipType ClipperClipType, solution [][]ClipperIntPoint, fillTypes ...ClipperPolyFillType) [][]ClipperIntPoint {
	if c.executeLocked {
		return nil
	}
	if c.hasOpenPaths {
		return nil // PolyTree struct needed for open path clipping
	}

	var subjFT, clipFT ClipperPolyFillType
	if len(fillTypes) >= 2 {
		subjFT = fillTypes[0]
		clipFT = fillTypes[1]
	} else if len(fillTypes) == 1 {
		subjFT = fillTypes[0]
		clipFT = fillTypes[0]
	} else {
		subjFT = ClipperEvenOdd
		clipFT = ClipperEvenOdd
	}

	c.executeLocked = true
	c.mSubjFillType = subjFT
	c.mClipFillType = clipFT
	c.mClipType = clipType
	c.mUsingPolyTree = false

	succeeded := c.executeInternal()
	var result [][]ClipperIntPoint
	if succeeded {
		result = c.buildResult()
	}
	c.disposeAllPolyPts()
	c.executeLocked = false

	return result
}

// ExecutePolyTree performs the clipping operation, returning a polygon tree.
func (c *Clipper) ExecutePolyTree(clipType ClipperClipType, polytree *ClipperPolyTree, fillTypes ...ClipperPolyFillType) bool {
	if c.executeLocked {
		return false
	}

	var subjFT, clipFT ClipperPolyFillType
	if len(fillTypes) >= 2 {
		subjFT = fillTypes[0]
		clipFT = fillTypes[1]
	} else if len(fillTypes) == 1 {
		subjFT = fillTypes[0]
		clipFT = fillTypes[0]
	} else {
		subjFT = ClipperEvenOdd
		clipFT = ClipperEvenOdd
	}

	c.executeLocked = true
	c.mSubjFillType = subjFT
	c.mClipFillType = clipFT
	c.mClipType = clipType
	c.mUsingPolyTree = true

	succeeded := c.executeInternal()
	if succeeded {
		c.buildResult2(polytree)
	}
	c.disposeAllPolyPts()
	c.executeLocked = false
	return succeeded
}

// AddPath adds a single path to the clipper.
func (c *Clipper) AddPath(pg []ClipperIntPoint, polyType ClipperPolyType, closed bool) bool {
	if !closed && polyType == ClipperClip {
		return false // Open paths must be subject
	}

	highI := len(pg) - 1
	if closed {
		for highI > 0 && pg[highI] == pg[0] {
			highI--
		}
	}
	for highI > 0 && pg[highI] == pg[highI-1] {
		highI--
	}
	if (closed && highI < 2) || (!closed && highI < 1) {
		return false
	}

	edges := make([]*ClipperTEdge, highI+1)
	for i := range edges {
		edges[i] = &ClipperTEdge{}
	}

	isFlat := true

	// Basic edge initialization
	edges[1].Curr = pg[1]
	c.rangeTest(pg[0])
	c.rangeTest(pg[highI])
	c.initEdge(edges[0], edges[1], edges[highI], pg[0])
	c.initEdge(edges[highI], edges[0], edges[highI-1], pg[highI])
	for i := highI - 1; i >= 1; i-- {
		c.rangeTest(pg[i])
		c.initEdge(edges[i], edges[i+1], edges[i-1], pg[i])
	}
	eStart := edges[0]

	// Remove duplicate vertices and collinear edges
	e := eStart
	eLoopStop := eStart
	for {
		if e.Curr == e.Next.Curr && (closed || e.Next != eStart) {
			if e == e.Next {
				break
			}
			if e == eStart {
				eStart = e.Next
			}
			e = c.removeEdge(e)
			eLoopStop = e
			continue
		}
		if e.Prev == e.Next {
			break
		} else if closed && slopesEqual3(e.Prev.Curr, e.Curr, e.Next.Curr, c.useFullRange) &&
			(!c.PreserveCollinear || !c.pt2BetweenPt1AndPt3(e.Prev.Curr, e.Curr, e.Next.Curr)) {
			if e == eStart {
				eStart = e.Next
			}
			e = c.removeEdge(e)
			e = e.Prev
			eLoopStop = e
			continue
		}
		e = e.Next
		if (e == eLoopStop) || (!closed && e.Next == eStart) {
			break
		}
	}

	if ((!closed && e == e.Next) || (closed && e.Prev == e.Next)) {
		return false
	}

	if !closed {
		c.hasOpenPaths = true
		eStart.Prev.OutIdx = Skip
	}

	// Second stage initialization
	e = eStart
	for {
		c.initEdge2(e, polyType)
		e = e.Next
		if isFlat && e.Curr.Y != eStart.Curr.Y {
			isFlat = false
		}
		if e == eStart || (!closed && e.Next == eStart) {
			break
		}
	}

	// Handle flat paths
	if isFlat {
		if closed {
			return false
		}
		e.Prev.OutIdx = Skip
		locMin := &ClipperLocalMinima{
			Y:          e.Bot.Y,
			LeftBound:  nil,
			RightBound: e,
		}
		locMin.RightBound.Side = ClipperRight
		locMin.RightBound.WindDelta = 0
		for {
			if e.Bot.X != e.Prev.Top.X {
				c.reverseHorizontal(e)
			}
			if e.Next.OutIdx == Skip {
				break
			}
			e.NextInLML = e.Next
			e = e.Next
		}
		c.insertLocalMinima(locMin)
		c.edges = append(c.edges, edges)
		return true
	}

	c.edges = append(c.edges, edges)

	// Workaround for open paths with matching start/end points
	if e.Prev.Bot == e.Prev.Top {
		e = e.Next
	}

	emin := (*ClipperTEdge)(nil)
	for {
		e = c.findNextLocMin(e)
		if emin != nil && e == emin {
			break
		}
		if emin == nil {
			emin = e
		}

		locMin := &ClipperLocalMinima{Y: e.Bot.Y}
		leftBoundIsForward := true
		if e.Dx < e.Prev.Dx {
			locMin.LeftBound = e.Prev
			locMin.RightBound = e
			leftBoundIsForward = false
		} else {
			locMin.LeftBound = e
			locMin.RightBound = e.Prev
		}

		locMin.LeftBound.Side = ClipperLeft
		locMin.RightBound.Side = ClipperRight

		if !closed {
			locMin.LeftBound.WindDelta = 0
		} else if locMin.LeftBound.Next == locMin.RightBound {
			locMin.LeftBound.WindDelta = -1
		} else {
			locMin.LeftBound.WindDelta = 1
		}
		locMin.RightBound.WindDelta = -locMin.LeftBound.WindDelta

		e = c.processBound(locMin.LeftBound, leftBoundIsForward)
		if e.OutIdx == Skip {
			e = c.processBound(e, leftBoundIsForward)
		}

		e2 := c.processBound(locMin.RightBound, !leftBoundIsForward)
		if e2.OutIdx == Skip {
			e2 = c.processBound(e2, !leftBoundIsForward)
		}

		if locMin.LeftBound.OutIdx == Skip {
			locMin.LeftBound = nil
		} else if locMin.RightBound.OutIdx == Skip {
			locMin.RightBound = nil
		}
		c.insertLocalMinima(locMin)

		if !leftBoundIsForward {
			e = e2
		}
	}

	return true
}

// AddPaths adds multiple paths to the clipper.
func (c *Clipper) AddPaths(ppg [][]ClipperIntPoint, polyType ClipperPolyType, closed bool) bool {
	result := false
	for _, pg := range ppg {
		if c.AddPath(pg, polyType, closed) {
			result = true
		}
	}
	return result
}

// ReversePaths reverses the winding direction of all polygons.
func ReversePaths(polys [][]ClipperIntPoint) {
	for _, poly := range polys {
		reverseSlice(poly)
	}
}

// Orientation returns true if the polygon has clockwise orientation (area >= 0).
func Orientation(poly []ClipperIntPoint) bool {
	return AreaPoly(poly) >= 0
}

// AreaPoly calculates the signed area of a polygon.
func AreaPoly(poly []ClipperIntPoint) float64 {
	cnt := len(poly)
	if cnt < 3 {
		return 0
	}
	a := 0.0
	for i, j := 0, cnt-1; i < cnt; i++ {
		a += (float64(poly[j].X) + float64(poly[i].X)) * (float64(poly[j].Y) - float64(poly[i].Y))
		j = i
	}
	return -a * 0.5
}

// PointInPolygon tests if a point is inside, outside, or on the boundary of a polygon.
// Returns +1 if inside, -1 if on boundary, 0 if outside.
func PointInPolygon(pt ClipperIntPoint, path []ClipperIntPoint) int {
	cnt := len(path)
	if cnt < 3 {
		return 0
	}
	result := 0
	ip := path[0]
	for i := 1; i <= cnt; i++ {
		ipNext := path[0]
		if i < cnt {
			ipNext = path[i]
		}
		if ipNext.Y == pt.Y {
			if ipNext.X == pt.X || (ip.Y == pt.Y && ((ipNext.X > pt.X) == (ip.X < pt.X))) {
				return -1
			}
		}
		if (ip.Y < pt.Y) != (ipNext.Y < pt.Y) {
			if ip.X >= pt.X {
				if ipNext.X > pt.X {
					result = 1 - result
				} else {
					d := float64(ip.X-pt.X)*float64(ipNext.Y-pt.Y) - float64(ipNext.X-pt.X)*float64(ip.Y-pt.Y)
					if d == 0 {
						return -1
					}
					if (d > 0) == (ipNext.Y > ip.Y) {
						result = 1 - result
					}
				}
			} else if ipNext.X > pt.X {
				d := float64(ip.X-pt.X)*float64(ipNext.Y-pt.Y) - float64(ipNext.X-pt.X)*float64(ip.Y-pt.Y)
				if d == 0 {
					return -1
				}
				if (d > 0) == (ipNext.Y > ip.Y) {
					result = 1 - result
				}
			}
		}
		ip = ipNext
	}
	return result
}

// SimplifyPolygon converts a self-intersecting polygon into simple polygons.
func SimplifyPolygon(poly []ClipperIntPoint, fillType ClipperPolyFillType) [][]ClipperIntPoint {
	c := NewClipper(0)
	c.StrictlySimple = true
	c.AddPath(poly, ClipperSubject, true)
	return c.Execute(ClipperUnion, nil, fillType, fillType)
}

// SimplifyPolygons converts self-intersecting polygons into simple polygons.
func SimplifyPolygons(polys [][]ClipperIntPoint, fillType ClipperPolyFillType) [][]ClipperIntPoint {
	c := NewClipper(0)
	c.StrictlySimple = true
	c.AddPaths(polys, ClipperSubject, true)
	return c.Execute(ClipperUnion, nil, fillType, fillType)
}

// CleanPolygon removes redundant vertices from a polygon.
func CleanPolygon(path []ClipperIntPoint, distance ...float64) []ClipperIntPoint {
	dist := 1.415
	if len(distance) > 0 {
		dist = distance[0]
	}
	cnt := len(path)
	if cnt == 0 {
		return make([]ClipperIntPoint, 0)
	}

	outPts := make([]*ClipperOutPt, cnt)
	for i := range outPts {
		outPts[i] = &ClipperOutPt{}
	}
	for i := 0; i < cnt; i++ {
		outPts[i].Pt = path[i]
		outPts[i].Next = outPts[(i+1)%cnt]
		outPts[i].Next.Prev = outPts[i]
		outPts[i].Index = 0
	}

	distSqrd := dist * dist
	op := outPts[0]
	for op.Index == 0 && op.Next != op.Prev {
		if pointsAreClose(op.Pt, op.Prev.Pt, distSqrd) {
			op = excludeOp(op)
			cnt--
		} else if pointsAreClose(op.Prev.Pt, op.Next.Pt, distSqrd) {
			excludeOp(op.Next)
			op = excludeOp(op)
			cnt -= 2
		} else if slopesNearCollinear(op.Prev.Pt, op.Pt, op.Next.Pt, distSqrd) {
			op = excludeOp(op)
			cnt--
		} else {
			op.Index = 1
			op = op.Next
		}
	}

	if cnt < 3 {
		cnt = 0
	}
	result := make([]ClipperIntPoint, 0, cnt)
	for i := 0; i < cnt; i++ {
		result = append(result, op.Pt)
		op = op.Next
	}
	return result
}

// CleanPolygons cleans multiple polygons.
func CleanPolygons(polys [][]ClipperIntPoint, distance ...float64) [][]ClipperIntPoint {
	dist := 1.415
	if len(distance) > 0 {
		dist = distance[0]
	}
	result := make([][]ClipperIntPoint, 0, len(polys))
	for _, poly := range polys {
		result = append(result, CleanPolygon(poly, dist))
	}
	return result
}

// MinkowskiSum computes the Minkowski sum of a pattern and path.
func MinkowskiSum(pattern []ClipperIntPoint, path interface{}, pathIsClosed bool) [][]ClipperIntPoint {
	paths := minkowskiSumInternal(pattern, path, pathIsClosed)
	c := NewClipper(0)
	c.AddPaths(paths, ClipperSubject, true)
	return c.Execute(ClipperUnion, nil, ClipperNonZero, ClipperNonZero)
}

// MinkowskiSumPaths computes the Minkowski sum of a pattern with multiple paths,
// optionally using clip regions for closed paths.
func MinkowskiSumPaths(pattern []ClipperIntPoint, paths [][]ClipperIntPoint, pathIsClosed bool) [][]ClipperIntPoint {
	c := NewClipper(0)
	for i := 0; i < len(paths); i++ {
		tmp := minkowskiInternal(pattern, paths[i], true, pathIsClosed)
		c.AddPaths(tmp, ClipperSubject, true)
		if pathIsClosed {
			path := translatePath(paths[i], pattern[0])
			c.AddPath(path, ClipperClip, true)
		}
	}
	return c.Execute(ClipperUnion, nil, ClipperNonZero, ClipperNonZero)
}

// MinkowskiDiff computes the Minkowski difference of two polygons.
func MinkowskiDiff(poly1 []ClipperIntPoint, poly2 []ClipperIntPoint) [][]ClipperIntPoint {
	paths := minkowskiInternal(poly1, poly2, false, true)
	c := NewClipper(0)
	c.AddPaths(paths, ClipperSubject, true)
	return c.Execute(ClipperUnion, nil, ClipperNonZero, ClipperNonZero)
}

// PolyTreeToPaths converts a polygon tree to a list of paths.
func PolyTreeToPaths(polytree *ClipperPolyTree) [][]ClipperIntPoint {
	result := make([][]ClipperIntPoint, 0, polytree.Total())
	addPolyNodeToPaths(polytree.ClipperPolyNode, ntAny, &result)
	return result
}

// OpenPathsFromPolyTree extracts open paths from a polygon tree.
func OpenPathsFromPolyTree(polytree *ClipperPolyTree) [][]ClipperIntPoint {
	result := make([][]ClipperIntPoint, 0, polytree.ChildCount())
	for _, child := range polytree.Children {
		if child.IsOpen {
			result = append(result, child.Polygon)
		}
	}
	return result
}

// ClosedPathsFromPolyTree extracts closed paths from a polygon tree.
func ClosedPathsFromPolyTree(polytree *ClipperPolyTree) [][]ClipperIntPoint {
	result := make([][]ClipperIntPoint, 0, polytree.Total())
	addPolyNodeToPaths(polytree.ClipperPolyNode, ntClosed, &result)
	return result
}

type nodeType int

const (
	ntAny nodeType = iota
	ntOpen
	ntClosed
)

func addPolyNodeToPaths(polynode *ClipperPolyNode, nt nodeType, paths *[][]ClipperIntPoint) {
	switch nt {
	case ntOpen:
		return
	case ntClosed:
		if polynode.IsOpen {
			for _, pn := range polynode.Children {
				addPolyNodeToPaths(pn, nt, paths)
			}
			return
		}
	}

	if len(polynode.Polygon) > 0 {
		*paths = append(*paths, polynode.Polygon)
	}
	for _, pn := range polynode.Children {
		addPolyNodeToPaths(pn, nt, paths)
	}
}

func minkowskiSumInternal(pattern []ClipperIntPoint, path interface{}, pathIsClosed bool) [][]ClipperIntPoint {
	switch p := path.(type) {
	case []ClipperIntPoint:
		return minkowskiInternal(pattern, p, true, pathIsClosed)
	case [][]ClipperIntPoint:
		var result [][]ClipperIntPoint
		for _, singlePath := range p {
			tmp := minkowskiInternal(pattern, singlePath, true, pathIsClosed)
			result = append(result, tmp...)
		}
		return result
	}
	return nil
}

func minkowskiInternal(pattern []ClipperIntPoint, pathList []ClipperIntPoint, isSum bool, pathIsClosed bool) [][]ClipperIntPoint {
	delta := 0
	if pathIsClosed {
		delta = 1
	}
	polyCnt := len(pattern)
	pathCnt := len(pathList)

	result := make([][]ClipperIntPoint, 0, pathCnt)
	for i := 0; i < pathCnt; i++ {
		p := make([]ClipperIntPoint, 0, polyCnt)
		for _, ip := range pattern {
			if isSum {
				p = append(p, NewClipperIntPoint(pathList[i].X+ip.X, pathList[i].Y+ip.Y))
			} else {
				p = append(p, NewClipperIntPoint(pathList[i].X-ip.X, pathList[i].Y-ip.Y))
			}
		}
		result = append(result, p)
	}

	quads := make([][]ClipperIntPoint, 0, (pathCnt+delta)*(polyCnt+1))
	for i := 0; i < pathCnt-1+delta; i++ {
		for j := 0; j < polyCnt; j++ {
			quad := []ClipperIntPoint{
				result[i%pathCnt][j%polyCnt],
				result[(i+1)%pathCnt][j%polyCnt],
				result[(i+1)%pathCnt][(j+1)%polyCnt],
				result[i%pathCnt][(j+1)%polyCnt],
			}
			if !Orientation(quad) {
				reverseSlice(quad)
			}
			quads = append(quads, quad)
		}
	}
	return quads
}

func translatePath(path []ClipperIntPoint, delta ClipperIntPoint) []ClipperIntPoint {
	outPath := make([]ClipperIntPoint, len(path))
	for i := 0; i < len(path); i++ {
		outPath[i] = NewClipperIntPoint(path[i].X+delta.X, path[i].Y+delta.Y)
	}
	return outPath
}

// GetBounds returns the bounding rectangle of a set of paths.
func GetBounds(paths [][]ClipperIntPoint) ClipperIntRect {
	i, cnt := 0, len(paths)
	for i < cnt && len(paths[i]) == 0 {
		i++
	}
	if i == cnt {
		return NewClipperIntRect(0, 0, 0, 0)
	}
	result := ClipperIntRect{
		Left:   paths[i][0].X,
		Right:  paths[i][0].X,
		Top:    paths[i][0].Y,
		Bottom: paths[i][0].Y,
	}
	for ; i < cnt; i++ {
		for j := 0; j < len(paths[i]); j++ {
			if paths[i][j].X < result.Left {
				result.Left = paths[i][j].X
			} else if paths[i][j].X > result.Right {
				result.Right = paths[i][j].X
			}
			if paths[i][j].Y < result.Top {
				result.Top = paths[i][j].Y
			} else if paths[i][j].Y > result.Bottom {
				result.Bottom = paths[i][j].Y
			}
		}
	}
	return result
}

// --- Internal execution methods ---

func (c *Clipper) executeInternal() bool {
	c.reset()
	c.mSortedEdges = nil
	c.mMaxima = nil

	botY, ok := c.popScanbeam()
	if !ok {
		return false
	}
	c.insertLocalMinimaIntoAEL(botY)

	for {
		topY, hasTop := c.popScanbeam()
		if !hasTop && !c.localMinimaPending() {
			break
		}
		c.processHorizontals()
		c.mGhostJoins = c.mGhostJoins[:0]
		if !c.processIntersections(topY) {
			return false
		}
		c.processEdgesAtTopOfScanbeam(topY)
		botY = topY
		c.insertLocalMinimaIntoAEL(botY)
	}

	// Fix orientations
	for _, outRec := range c.polyOuts {
		if outRec.Pts == nil || outRec.IsOpen {
			continue
		}
		if (outRec.IsHole != c.ReverseSolution) == (c.areaOutPt(outRec.Pts) > 0) {
			c.reversePolyPtLinks(outRec.Pts)
		}
	}

	c.joinCommonEdges()

	for _, outRec := range c.polyOuts {
		if outRec.Pts == nil {
			continue
		} else if outRec.IsOpen {
			c.fixupOutPolyline(outRec)
		} else {
			c.fixupOutPolygon(outRec)
		}
	}

	if c.StrictlySimple {
		c.doSimplePolygons()
	}

	c.mJoins = c.mJoins[:0]
	c.mGhostJoins = c.mGhostJoins[:0]
	return true
}

func (c *Clipper) disposeAllPolyPts() {
	for i := range c.polyOuts {
		c.disposeOutRec(i)
	}
	c.polyOuts = c.polyOuts[:0]
}

func (c *Clipper) insertLocalMinimaIntoAEL(botY int64) {
	var lm *ClipperLocalMinima
	for c.popLocalMinima(botY, &lm) {
		lb := lm.LeftBound
		rb := lm.RightBound

		var op1 *ClipperOutPt
		if lb == nil {
			c.insertEdgeIntoAEL(rb, nil)
			c.setWindingCount(rb)
			if c.isContributing(rb) {
				op1 = c.addOutPt(rb, rb.Bot)
			}
		} else if rb == nil {
			c.insertEdgeIntoAEL(lb, nil)
			c.setWindingCount(lb)
			if c.isContributing(lb) {
				op1 = c.addOutPt(lb, lb.Bot)
			}
			c.insertScanbeam(lb.Top.Y)
		} else {
			c.insertEdgeIntoAEL(lb, nil)
			c.insertEdgeIntoAEL(rb, lb)
			c.setWindingCount(lb)
			rb.WindCnt = lb.WindCnt
			rb.WindCnt2 = lb.WindCnt2
			if c.isContributing(lb) {
				op1 = c.addLocalMinPoly(lb, rb, lb.Bot)
			}
			c.insertScanbeam(lb.Top.Y)
		}

		if rb != nil {
			if isHorizontal(rb) {
				if rb.NextInLML != nil {
					c.insertScanbeam(rb.NextInLML.Top.Y)
				}
				c.addEdgeToSEL(rb)
			} else {
				c.insertScanbeam(rb.Top.Y)
			}
		}

		if lb == nil || rb == nil {
			continue
		}

		// If output polygons share an edge with a horizontal rb, they'll need joining later
		if op1 != nil && isHorizontal(rb) && len(c.mGhostJoins) > 0 && rb.WindDelta != 0 {
			for i := 0; i < len(c.mGhostJoins); i++ {
				j := c.mGhostJoins[i]
				if horzSegmentsOverlap(j.OutPt1.Pt.X, j.OffPt.X, rb.Bot.X, rb.Top.X) {
					c.addJoin(j.OutPt1, op1, j.OffPt)
				}
			}
		}

		if lb.OutIdx >= 0 && lb.PrevInAEL != nil &&
			lb.PrevInAEL.Curr.X == lb.Bot.X &&
			lb.PrevInAEL.OutIdx >= 0 &&
			slopesEqual4(lb.PrevInAEL.Curr, lb.PrevInAEL.Top, lb.Curr, lb.Top, c.useFullRange) &&
			lb.WindDelta != 0 && lb.PrevInAEL.WindDelta != 0 {
			op2 := c.addOutPt(lb.PrevInAEL, lb.Bot)
			c.addJoin(op1, op2, lb.Top)
		}

		if lb.NextInAEL != rb {
			if rb.OutIdx >= 0 && rb.PrevInAEL.OutIdx >= 0 &&
				slopesEqual4(rb.PrevInAEL.Curr, rb.PrevInAEL.Top, rb.Curr, rb.Top, c.useFullRange) &&
				rb.WindDelta != 0 && rb.PrevInAEL.WindDelta != 0 {
				op2 := c.addOutPt(rb.PrevInAEL, rb.Bot)
				c.addJoin(op1, op2, rb.Top)
			}

			e := lb.NextInAEL
			for e != nil && e != rb {
				c.intersectEdges(rb, e, lb.Curr)
				e = e.NextInAEL
			}
		}
	}
}

func (c *Clipper) insertEdgeIntoAEL(edge, startEdge *ClipperTEdge) {
	if c.activeEdges == nil {
		edge.PrevInAEL = nil
		edge.NextInAEL = nil
		c.activeEdges = edge
	} else if startEdge == nil && e2InsertsBeforeE1(c.activeEdges, edge) {
		edge.PrevInAEL = nil
		edge.NextInAEL = c.activeEdges
		c.activeEdges.PrevInAEL = edge
		c.activeEdges = edge
	} else {
		if startEdge == nil {
			startEdge = c.activeEdges
		}
		for startEdge.NextInAEL != nil && !e2InsertsBeforeE1(startEdge.NextInAEL, edge) {
			startEdge = startEdge.NextInAEL
		}
		edge.NextInAEL = startEdge.NextInAEL
		if startEdge.NextInAEL != nil {
			startEdge.NextInAEL.PrevInAEL = edge
		}
		edge.PrevInAEL = startEdge
		startEdge.NextInAEL = edge
	}
}

func e2InsertsBeforeE1(e1, e2 *ClipperTEdge) bool {
	if e2.Curr.X == e1.Curr.X {
		if e2.Top.Y > e1.Top.Y {
			return e2.Top.X < topX(e1, e2.Top.Y)
		}
		return e1.Top.X > topX(e2, e1.Top.Y)
	}
	return e2.Curr.X < e1.Curr.X
}

func (c *Clipper) isEvenOddFillType(edge *ClipperTEdge) bool {
	if edge.PolyTyp == ClipperSubject {
		return c.mSubjFillType == ClipperEvenOdd
	}
	return c.mClipFillType == ClipperEvenOdd
}

func (c *Clipper) isEvenOddAltFillType(edge *ClipperTEdge) bool {
	if edge.PolyTyp == ClipperSubject {
		return c.mClipFillType == ClipperEvenOdd
	}
	return c.mSubjFillType == ClipperEvenOdd
}

func (c *Clipper) isContributing(edge *ClipperTEdge) bool {
var pft ClipperPolyFillType
	if edge.PolyTyp == ClipperSubject {
		pft = c.mSubjFillType
	} else {
		pft = c.mClipFillType
	}

	switch pft {
	default:
	case ClipperEvenOdd:
		if edge.WindDelta == 0 && edge.WindCnt != 1 {
			return false
		}
	case ClipperNonZero:
		if int(math.Abs(float64(edge.WindCnt))) != 1 {
			return false
		}
	}

	switch c.mClipType {
	case ClipperIntersection:
		return edge.WindCnt2 != 0
	case ClipperUnion:
		return edge.WindCnt2 == 0
	case ClipperDifference:
		if edge.PolyTyp == ClipperSubject {
			return edge.WindCnt2 == 0
		}
		return edge.WindCnt2 != 0
	case ClipperXor:
		if edge.WindDelta == 0 {
			return edge.WindCnt2 == 0
		}
		return true
	}
	return true
}

func (c *Clipper) setWindingCount(edge *ClipperTEdge) {
	e := edge.PrevInAEL
	for e != nil && ((e.PolyTyp != edge.PolyTyp) || (e.WindDelta == 0)) {
		e = e.PrevInAEL
	}
	if e == nil {
		var pft ClipperPolyFillType
		if edge.PolyTyp == ClipperSubject {
			pft = c.mSubjFillType
		} else {
			pft = c.mClipFillType
		}
		_ = pft
		if edge.WindDelta == 0 {
			edge.WindCnt = 1
		} else {
			edge.WindCnt = edge.WindDelta
		}
		edge.WindCnt2 = 0
		e = c.activeEdges
	} else if edge.WindDelta == 0 && c.mClipType != ClipperUnion {
		edge.WindCnt = 1
		edge.WindCnt2 = e.WindCnt2
		e = e.NextInAEL
	} else if c.isEvenOddFillType(edge) {
		if edge.WindDelta == 0 {
			inside := true
			e2 := e.PrevInAEL
			for e2 != nil {
				if e2.PolyTyp == e.PolyTyp && e2.WindDelta != 0 {
					inside = !inside
				}
				e2 = e2.PrevInAEL
			}
			if inside {
				edge.WindCnt = 0
			} else {
				edge.WindCnt = 1
			}
		} else {
			edge.WindCnt = edge.WindDelta
		}
		edge.WindCnt2 = e.WindCnt2
		e = e.NextInAEL
	} else {
		if e.WindCnt*e.WindDelta < 0 {
			if int(math.Abs(float64(e.WindCnt))) > 1 {
				if e.WindDelta*edge.WindDelta < 0 {
					edge.WindCnt = e.WindCnt
				} else {
					edge.WindCnt = e.WindCnt + edge.WindDelta
				}
			} else {
				if edge.WindDelta == 0 {
					edge.WindCnt = 1
				} else {
					edge.WindCnt = edge.WindDelta
				}
			}
		} else {
			if edge.WindDelta == 0 {
				if e.WindCnt < 0 {
					edge.WindCnt = e.WindCnt - 1
				} else {
					edge.WindCnt = e.WindCnt + 1
				}
			} else if e.WindDelta*edge.WindDelta < 0 {
				edge.WindCnt = e.WindCnt
			} else {
				edge.WindCnt = e.WindCnt + edge.WindDelta
			}
		}
		edge.WindCnt2 = e.WindCnt2
		e = e.NextInAEL
	}

	if c.isEvenOddAltFillType(edge) {
		for e != edge {
			if e.WindDelta != 0 {
				if edge.WindCnt2 == 0 {
					edge.WindCnt2 = 1
				} else {
					edge.WindCnt2 = 0
				}
			}
			e = e.NextInAEL
		}
	} else {
		for e != edge {
			edge.WindCnt2 += e.WindDelta
			e = e.NextInAEL
		}
	}
}

func (c *Clipper) addEdgeToSEL(edge *ClipperTEdge) {
	if c.mSortedEdges == nil {
		c.mSortedEdges = edge
		edge.PrevInSEL = nil
		edge.NextInSEL = nil
	} else {
		edge.NextInSEL = c.mSortedEdges
		edge.PrevInSEL = nil
		c.mSortedEdges.PrevInSEL = edge
		c.mSortedEdges = edge
	}
}

func (c *Clipper) popEdgeFromSEL() (*ClipperTEdge, bool) {
	e := c.mSortedEdges
	if e == nil {
		return nil, false
	}
	c.mSortedEdges = e.NextInSEL
	if c.mSortedEdges != nil {
		c.mSortedEdges.PrevInSEL = nil
	}
	e.NextInSEL = nil
	e.PrevInSEL = nil
	return e, true
}

func (c *Clipper) copyAELToSEL() {
	e := c.activeEdges
	c.mSortedEdges = e
	for e != nil {
		e.PrevInSEL = e.PrevInAEL
		e.NextInSEL = e.NextInAEL
		e = e.NextInAEL
	}
}

func (c *Clipper) swapPositionsInSEL(edge1, edge2 *ClipperTEdge) {
	if edge1.NextInSEL == nil && edge1.PrevInSEL == nil {
		return
	}
	if edge2.NextInSEL == nil && edge2.PrevInSEL == nil {
		return
	}

	if edge1.NextInSEL == edge2 {
		next := edge2.NextInSEL
		if next != nil {
			next.PrevInSEL = edge1
		}
		prev := edge1.PrevInSEL
		if prev != nil {
			prev.NextInSEL = edge2
		}
		edge2.PrevInSEL = prev
		edge2.NextInSEL = edge1
		edge1.PrevInSEL = edge2
		edge1.NextInSEL = next
	} else if edge2.NextInSEL == edge1 {
		next := edge1.NextInSEL
		if next != nil {
			next.PrevInSEL = edge2
		}
		prev := edge2.PrevInSEL
		if prev != nil {
			prev.NextInSEL = edge1
		}
		edge1.PrevInSEL = prev
		edge1.NextInSEL = edge2
		edge2.PrevInSEL = edge1
		edge2.NextInSEL = next
	} else {
		next := edge1.NextInSEL
		prev := edge1.PrevInSEL
		edge1.NextInSEL = edge2.NextInSEL
		if edge1.NextInSEL != nil {
			edge1.NextInSEL.PrevInSEL = edge1
		}
		edge1.PrevInSEL = edge2.PrevInSEL
		if edge1.PrevInSEL != nil {
			edge1.PrevInSEL.NextInSEL = edge1
		}
		edge2.NextInSEL = next
		if edge2.NextInSEL != nil {
			edge2.NextInSEL.PrevInSEL = edge2
		}
		edge2.PrevInSEL = prev
		if edge2.PrevInSEL != nil {
			edge2.PrevInSEL.NextInSEL = edge2
		}
	}

	if edge1.PrevInSEL == nil {
		c.mSortedEdges = edge1
	} else if edge2.PrevInSEL == nil {
		c.mSortedEdges = edge2
	}
}

// --- Output point management ---

func (c *Clipper) addLocalMaxPoly(e1, e2 *ClipperTEdge, pt ClipperIntPoint) {
	c.addOutPt(e1, pt)
	if e2.WindDelta == 0 {
		c.addOutPt(e2, pt)
	}
	if e1.OutIdx == e2.OutIdx {
		e1.OutIdx = Unassigned
		e2.OutIdx = Unassigned
	} else if e1.OutIdx < e2.OutIdx {
		c.appendPolygon(e1, e2)
	} else {
		c.appendPolygon(e2, e1)
	}
}

func (c *Clipper) addLocalMinPoly(e1, e2 *ClipperTEdge, pt ClipperIntPoint) *ClipperOutPt {
	var result *ClipperOutPt
	var e, prevE *ClipperTEdge

	if isHorizontal(e2) || e1.Dx > e2.Dx {
		result = c.addOutPt(e1, pt)
		e2.OutIdx = e1.OutIdx
		e1.Side = ClipperLeft
		e2.Side = ClipperRight
		e = e1
		if e.PrevInAEL == e2 {
			prevE = e2.PrevInAEL
		} else {
			prevE = e.PrevInAEL
		}
	} else {
		result = c.addOutPt(e2, pt)
		e1.OutIdx = e2.OutIdx
		e1.Side = ClipperRight
		e2.Side = ClipperLeft
		e = e2
		if e.PrevInAEL == e1 {
			prevE = e1.PrevInAEL
		} else {
			prevE = e.PrevInAEL
		}
	}

	if prevE != nil && prevE.OutIdx >= 0 && prevE.Top.Y < pt.Y && e.Top.Y < pt.Y {
		xPrev := topX(prevE, pt.Y)
		xE := topX(e, pt.Y)
		if xPrev == xE && e.WindDelta != 0 && prevE.WindDelta != 0 &&
			slopesEqual4(NewClipperIntPoint(xPrev, pt.Y), prevE.Top, NewClipperIntPoint(xE, pt.Y), e.Top, c.useFullRange) {
			outPt := c.addOutPt(prevE, pt)
			c.addJoin(result, outPt, e.Top)
		}
	}
	return result
}

func (c *Clipper) addOutPt(e *ClipperTEdge, pt ClipperIntPoint) *ClipperOutPt {
	if e.OutIdx < 0 {
		outRec := c.createOutRec()
		outRec.IsOpen = (e.WindDelta == 0)
		newOp := &ClipperOutPt{}
		outRec.Pts = newOp
		newOp.Index = outRec.Idx
		newOp.Pt = pt
		newOp.Next = newOp
		newOp.Prev = newOp
		if !outRec.IsOpen {
			c.setHoleState(e, outRec)
		}
		e.OutIdx = outRec.Idx
		return newOp
	}

	outRec := c.polyOuts[e.OutIdx]
	op := outRec.Pts
	toFront := (e.Side == ClipperLeft)
	if toFront && pt == op.Pt {
		return op
	}
	if !toFront && pt == op.Prev.Pt {
		return op.Prev
	}

	newOp := &ClipperOutPt{
		Index: outRec.Idx,
		Pt:    pt,
		Next:  op,
		Prev:  op.Prev,
	}
	newOp.Prev.Next = newOp
	op.Prev = newOp
	if toFront {
		outRec.Pts = newOp
	}
	return newOp
}

func (c *Clipper) getLastOutPt(e *ClipperTEdge) *ClipperOutPt {
	outRec := c.polyOuts[e.OutIdx]
	if e.Side == ClipperLeft {
		return outRec.Pts
	}
	return outRec.Pts.Prev
}

func (c *Clipper) setHoleState(e *ClipperTEdge, outRec *ClipperOutRec) {
	e2 := e.PrevInAEL
	var eTmp *ClipperTEdge
	for e2 != nil {
		if e2.OutIdx >= 0 && e2.WindDelta != 0 {
			if eTmp == nil {
				eTmp = e2
			} else if eTmp.OutIdx == e2.OutIdx {
				eTmp = nil
			}
		}
		e2 = e2.PrevInAEL
	}

	if eTmp == nil {
		outRec.FirstLeft = nil
		outRec.IsHole = false
	} else {
		outRec.FirstLeft = c.polyOuts[eTmp.OutIdx]
		outRec.IsHole = !outRec.FirstLeft.IsHole
	}
}

// --- Intersection processing ---

func (c *Clipper) intersectEdges(e1, e2 *ClipperTEdge, pt ClipperIntPoint) {
	e1Contributing := e1.OutIdx >= 0
	e2Contributing := e2.OutIdx >= 0

	if e1.WindDelta == 0 || e2.WindDelta == 0 {
		if e1.WindDelta == 0 && e2.WindDelta == 0 {
			return
		} else if e1.PolyTyp == e2.PolyTyp && e1.WindDelta != e2.WindDelta && c.mClipType == ClipperUnion {
			if e1.WindDelta == 0 {
				if e2Contributing {
					c.addOutPt(e1, pt)
					if e1Contributing {
						e1.OutIdx = Unassigned
					}
				}
			} else {
				if e1Contributing {
					c.addOutPt(e2, pt)
					if e2Contributing {
						e2.OutIdx = Unassigned
					}
				}
			}
		} else if e1.PolyTyp != e2.PolyTyp {
			if e1.WindDelta == 0 && int(math.Abs(float64(e2.WindCnt))) == 1 &&
				(c.mClipType != ClipperUnion || e2.WindCnt2 == 0) {
				c.addOutPt(e1, pt)
				if e1Contributing {
					e1.OutIdx = Unassigned
				}
			} else if e2.WindDelta == 0 && int(math.Abs(float64(e1.WindCnt))) == 1 &&
				(c.mClipType != ClipperUnion || e1.WindCnt2 == 0) {
				c.addOutPt(e2, pt)
				if e2Contributing {
					e2.OutIdx = Unassigned
				}
			}
		}
		return
	}

	if e1.PolyTyp == e2.PolyTyp {
		if c.isEvenOddFillType(e1) {
			oldE1WindCnt := e1.WindCnt
			e1.WindCnt = e2.WindCnt
			e2.WindCnt = oldE1WindCnt
		} else {
			if e1.WindCnt+e2.WindDelta == 0 {
				e1.WindCnt = -e1.WindCnt
			} else {
				e1.WindCnt += e2.WindDelta
			}
			if e2.WindCnt-e1.WindDelta == 0 {
				e2.WindCnt = -e2.WindCnt
			} else {
				e2.WindCnt -= e1.WindDelta
			}
		}
	} else {
		if !c.isEvenOddFillType(e2) {
			e1.WindCnt2 += e2.WindDelta
		} else {
			if e1.WindCnt2 == 0 {
				e1.WindCnt2 = 1
			} else {
				e1.WindCnt2 = 0
			}
		}
		if !c.isEvenOddFillType(e1) {
			e2.WindCnt2 -= e1.WindDelta
		} else {
			if e2.WindCnt2 == 0 {
				e2.WindCnt2 = 1
			} else {
				e2.WindCnt2 = 0
			}
		}
	}

	e1Wc := int(math.Abs(float64(e1.WindCnt)))
	e2Wc := int(math.Abs(float64(e2.WindCnt)))

	if e1Contributing && e2Contributing {
		if (e1Wc != 0 && e1Wc != 1) || (e2Wc != 0 && e2Wc != 1) ||
			(e1.PolyTyp != e2.PolyTyp && c.mClipType != ClipperXor) {
			c.addLocalMaxPoly(e1, e2, pt)
		} else {
			c.addOutPt(e1, pt)
			c.addOutPt(e2, pt)
			swapSides(e1, e2)
			swapPolyIndexes(e1, e2)
		}
	} else if e1Contributing {
		if e2Wc == 0 || e2Wc == 1 {
			c.addOutPt(e1, pt)
			swapSides(e1, e2)
			swapPolyIndexes(e1, e2)
		}
	} else if e2Contributing {
		if e1Wc == 0 || e1Wc == 1 {
			c.addOutPt(e2, pt)
			swapSides(e1, e2)
			swapPolyIndexes(e1, e2)
		}
	} else if (e1Wc == 0 || e1Wc == 1) && (e2Wc == 0 || e2Wc == 1) {
		e1Wc2 := int(math.Abs(float64(e1.WindCnt2)))
		e2Wc2 := int(math.Abs(float64(e2.WindCnt2)))

		if e1.PolyTyp != e2.PolyTyp {
			c.addLocalMinPoly(e1, e2, pt)
		} else if e1Wc == 1 && e2Wc == 1 {
			switch c.mClipType {
			case ClipperIntersection:
				if e1Wc2 > 0 && e2Wc2 > 0 {
					c.addLocalMinPoly(e1, e2, pt)
				}
			case ClipperUnion:
				if e1Wc2 <= 0 && e2Wc2 <= 0 {
					c.addLocalMinPoly(e1, e2, pt)
				}
			case ClipperDifference:
				if (e1.PolyTyp == ClipperClip && e1Wc2 > 0 && e2Wc2 > 0) ||
					(e1.PolyTyp == ClipperSubject && e1Wc2 <= 0 && e2Wc2 <= 0) {
					c.addLocalMinPoly(e1, e2, pt)
				}
			case ClipperXor:
				c.addLocalMinPoly(e1, e2, pt)
			}
		} else {
			swapSides(e1, e2)
		}
	}
}

func (c *Clipper) processIntersections(topY int64) bool {
	if c.activeEdges == nil {
		return true
	}
	c.buildIntersectList(topY)
	if len(c.mIntersectList) == 0 {
		return true
	}
	if len(c.mIntersectList) == 1 || c.fixupIntersectionOrder() {
		c.processIntersectList()
	} else {
		return false
	}
	c.mSortedEdges = nil
	return true
}

func (c *Clipper) buildIntersectList(topY int64) {
	if c.activeEdges == nil {
		return
	}

	e := c.activeEdges
	c.mSortedEdges = e
	for e != nil {
		e.PrevInSEL = e.PrevInAEL
		e.NextInSEL = e.NextInAEL
		e.Curr.X = topX(e, topY)
		e = e.NextInAEL
	}

	isModified := true
	for isModified && c.mSortedEdges != nil {
		isModified = false
		e = c.mSortedEdges
		for e.NextInSEL != nil {
			eNext := e.NextInSEL
			if e.Curr.X > eNext.Curr.X {
				var pt ClipperIntPoint
				c.intersectPoint(e, eNext, &pt)
				if pt.Y < topY {
					pt = NewClipperIntPoint(topX(e, topY), topY)
				}
				c.mIntersectList = append(c.mIntersectList, &ClipperIntersectNode{
					Edge1: e, Edge2: eNext, Pt: pt,
				})
				c.swapPositionsInSEL(e, eNext)
				isModified = true
			} else {
				e = eNext
			}
		}
		if e.PrevInSEL != nil {
			e.PrevInSEL.NextInSEL = nil
		} else {
			break
		}
	}
	c.mSortedEdges = nil
}

func edgesAdjacent(inode *ClipperIntersectNode) bool {
	return inode.Edge1.NextInSEL == inode.Edge2 || inode.Edge1.PrevInSEL == inode.Edge2
}

func (c *Clipper) fixupIntersectionOrder() bool {
	sort.Slice(c.mIntersectList, func(i, j int) bool {
		return c.mIntersectList[i].Pt.Y > c.mIntersectList[j].Pt.Y
	})

	c.copyAELToSEL()
	cnt := len(c.mIntersectList)
	for i := 0; i < cnt; i++ {
		if !edgesAdjacent(c.mIntersectList[i]) {
			j := i + 1
			for j < cnt && !edgesAdjacent(c.mIntersectList[j]) {
				j++
			}
			if j == cnt {
				return false
			}
			c.mIntersectList[i], c.mIntersectList[j] = c.mIntersectList[j], c.mIntersectList[i]
		}
		c.swapPositionsInSEL(c.mIntersectList[i].Edge1, c.mIntersectList[i].Edge2)
	}
	return true
}

func (c *Clipper) processIntersectList() {
	for _, iNode := range c.mIntersectList {
		c.intersectEdges(iNode.Edge1, iNode.Edge2, iNode.Pt)
		c.swapPositionsInAEL(iNode.Edge1, iNode.Edge2)
	}
	c.mIntersectList = c.mIntersectList[:0]
}

// Round rounds a float64 to the nearest int64.
func Round(value float64) int64 {
	if value < 0 {
		return int64(value - 0.5)
	}
	return int64(value + 0.5)
}

func topX(edge *ClipperTEdge, currentY int64) int64 {
	if currentY == edge.Top.Y {
		return edge.Top.X
	}
	return edge.Bot.X + Round(edge.Dx*(float64(currentY)-float64(edge.Bot.Y)))
}

func (c *Clipper) intersectPoint(edge1, edge2 *ClipperTEdge, ip *ClipperIntPoint) {
	ip.X = 0
	ip.Y = 0

	if edge1.Dx == edge2.Dx {
		ip.Y = edge1.Curr.Y
		ip.X = topX(edge1, ip.Y)
		return
	}

	if edge1.Delta.X == 0 {
		ip.X = edge1.Bot.X
		if isHorizontal(edge2) {
			ip.Y = edge2.Bot.Y
		} else {
			b2 := float64(edge2.Bot.Y) - float64(edge2.Bot.X)/edge2.Dx
			ip.Y = Round(float64(ip.X)/edge2.Dx + b2)
		}
	} else if edge2.Delta.X == 0 {
		ip.X = edge2.Bot.X
		if isHorizontal(edge1) {
			ip.Y = edge1.Bot.Y
		} else {
			b1 := float64(edge1.Bot.Y) - float64(edge1.Bot.X)/edge1.Dx
			ip.Y = Round(float64(ip.X)/edge1.Dx + b1)
		}
	} else {
		b1 := float64(edge1.Bot.X) - float64(edge1.Bot.Y)*edge1.Dx
		b2 := float64(edge2.Bot.X) - float64(edge2.Bot.Y)*edge2.Dx
		q := (b2 - b1) / (edge1.Dx - edge2.Dx)
		ip.Y = Round(q)
		if math.Abs(edge1.Dx) < math.Abs(edge2.Dx) {
			ip.X = Round(edge1.Dx*q + b1)
		} else {
			ip.X = Round(edge2.Dx*q + b2)
		}
	}

	if ip.Y < edge1.Top.Y || ip.Y < edge2.Top.Y {
		if edge1.Top.Y > edge2.Top.Y {
			ip.Y = edge1.Top.Y
		} else {
			ip.Y = edge2.Top.Y
		}
		if math.Abs(edge1.Dx) < math.Abs(edge2.Dx) {
			ip.X = topX(edge1, ip.Y)
		} else {
			ip.X = topX(edge2, ip.Y)
		}
	}

	if ip.Y > edge1.Curr.Y {
		ip.Y = edge1.Curr.Y
		if math.Abs(edge1.Dx) > math.Abs(edge2.Dx) {
			ip.X = topX(edge2, ip.Y)
		} else {
			ip.X = topX(edge1, ip.Y)
		}
	}
}

// --- Horizontal edge and maxima processing ---

func (c *Clipper) processHorizontals() {
	for {
		horzEdge, ok := c.popEdgeFromSEL()
		if !ok {
			break
		}
		c.processHorizontal(horzEdge)
	}
}

func (c *Clipper) getHorzDirection(horzEdge *ClipperTEdge) (dir ClipperDirection, left, right int64) {
	if horzEdge.Bot.X < horzEdge.Top.X {
		left = horzEdge.Bot.X
		right = horzEdge.Top.X
		dir = ClipperLeftToRight
	} else {
		left = horzEdge.Top.X
		right = horzEdge.Bot.X
		dir = ClipperRightToLeft
	}
	return
}

func (c *Clipper) processHorizontal(horzEdge *ClipperTEdge) {
	isOpen := horzEdge.WindDelta == 0
	dir, horzLeft, horzRight := c.getHorzDirection(horzEdge)

	eLastHorz := horzEdge
	for eLastHorz.NextInLML != nil && isHorizontal(eLastHorz.NextInLML) {
		eLastHorz = eLastHorz.NextInLML
	}
	var eMaxPair *ClipperTEdge
	if eLastHorz.NextInLML == nil {
		eMaxPair = c.getMaximaPair(eLastHorz)
	}

	currMax := c.mMaxima
	if currMax != nil {
		if dir == ClipperLeftToRight {
			for currMax != nil && currMax.X <= horzEdge.Bot.X {
				currMax = currMax.Next
			}
			if currMax != nil && currMax.X >= eLastHorz.Top.X {
				currMax = nil
			}
		} else {
			for currMax.Next != nil && currMax.Next.X < horzEdge.Bot.X {
				currMax = currMax.Next
			}
			if currMax.X <= eLastHorz.Top.X {
				currMax = nil
			}
		}
	}

	var op1 *ClipperOutPt
	for {
		isLastHorz := horzEdge == eLastHorz
		e := c.getNextInAEL(horzEdge, dir)
		for e != nil {
			if currMax != nil {
				if dir == ClipperLeftToRight {
					for currMax != nil && currMax.X < e.Curr.X {
						if horzEdge.OutIdx >= 0 && !isOpen {
							c.addOutPt(horzEdge, NewClipperIntPoint(currMax.X, horzEdge.Bot.Y))
						}
						currMax = currMax.Next
					}
				} else {
					for currMax != nil && currMax.X > e.Curr.X {
						if horzEdge.OutIdx >= 0 && !isOpen {
							c.addOutPt(horzEdge, NewClipperIntPoint(currMax.X, horzEdge.Bot.Y))
						}
						currMax = currMax.Previous
					}
				}
			}

			if (dir == ClipperLeftToRight && e.Curr.X > horzRight) ||
				(dir == ClipperRightToLeft && e.Curr.X < horzLeft) {
				break
			}
			if e.Curr.X == horzEdge.Top.X && horzEdge.NextInLML != nil && e.Dx < horzEdge.NextInLML.Dx {
				break
			}

			if horzEdge.OutIdx >= 0 && !isOpen {
				op1 = c.addOutPt(horzEdge, e.Curr)
				eNextHorz := c.mSortedEdges
				for eNextHorz != nil {
					if eNextHorz.OutIdx >= 0 &&
						horzSegmentsOverlap(horzEdge.Bot.X, horzEdge.Top.X, eNextHorz.Bot.X, eNextHorz.Top.X) {
						op2 := c.getLastOutPt(eNextHorz)
						c.addJoin(op2, op1, eNextHorz.Top)
					}
					eNextHorz = eNextHorz.NextInSEL
				}
				c.addGhostJoin(op1, horzEdge.Bot)
			}

			if e == eMaxPair && isLastHorz {
				if horzEdge.OutIdx >= 0 {
					c.addLocalMaxPoly(horzEdge, eMaxPair, horzEdge.Top)
				}
				c.deleteFromAEL(horzEdge)
				c.deleteFromAEL(eMaxPair)
				return
			}

			var pt ClipperIntPoint
			if dir == ClipperLeftToRight {
				pt = NewClipperIntPoint(e.Curr.X, horzEdge.Curr.Y)
				c.intersectEdges(horzEdge, e, pt)
			} else {
				pt = NewClipperIntPoint(e.Curr.X, horzEdge.Curr.Y)
				c.intersectEdges(e, horzEdge, pt)
			}

			eNext := c.getNextInAEL(e, dir)
			c.swapPositionsInAEL(horzEdge, e)
			e = eNext
		}

		if horzEdge.NextInLML == nil || !isHorizontal(horzEdge.NextInLML) {
			break
		}
		horzEdge = c.updateEdgeIntoAEL(horzEdge)
		if horzEdge.OutIdx >= 0 {
			c.addOutPt(horzEdge, horzEdge.Bot)
		}
		dir, horzLeft, horzRight = c.getHorzDirection(horzEdge)
	}

	if horzEdge.OutIdx >= 0 && op1 == nil {
		op1 = c.getLastOutPt(horzEdge)
		eNextHorz := c.mSortedEdges
		for eNextHorz != nil {
			if eNextHorz.OutIdx >= 0 &&
				horzSegmentsOverlap(horzEdge.Bot.X, horzEdge.Top.X, eNextHorz.Bot.X, eNextHorz.Top.X) {
				op2 := c.getLastOutPt(eNextHorz)
				c.addJoin(op2, op1, eNextHorz.Top)
			}
			eNextHorz = eNextHorz.NextInSEL
		}
		c.addGhostJoin(op1, horzEdge.Top)
	}

	if horzEdge.NextInLML != nil {
		if horzEdge.OutIdx >= 0 {
			op1 = c.addOutPt(horzEdge, horzEdge.Top)
			horzEdge = c.updateEdgeIntoAEL(horzEdge)
			if horzEdge.WindDelta == 0 {
				return
			}
			ePrev := horzEdge.PrevInAEL
			eNext := horzEdge.NextInAEL
			if ePrev != nil && ePrev.Curr.X == horzEdge.Bot.X &&
				ePrev.Curr.Y == horzEdge.Bot.Y && ePrev.WindDelta != 0 &&
				ePrev.OutIdx >= 0 && ePrev.Curr.Y > ePrev.Top.Y &&
				slopesEqual2h(horzEdge, ePrev, c.useFullRange) {
				op2 := c.addOutPt(ePrev, horzEdge.Bot)
				c.addJoin(op1, op2, horzEdge.Top)
			} else if eNext != nil && eNext.Curr.X == horzEdge.Bot.X &&
				eNext.Curr.Y == horzEdge.Bot.Y && eNext.WindDelta != 0 &&
				eNext.OutIdx >= 0 && eNext.Curr.Y > eNext.Top.Y &&
				slopesEqual2h(horzEdge, eNext, c.useFullRange) {
				op2 := c.addOutPt(eNext, horzEdge.Bot)
				c.addJoin(op1, op2, horzEdge.Top)
			}
		} else {
			horzEdge = c.updateEdgeIntoAEL(horzEdge)
		}
	} else {
		if horzEdge.OutIdx >= 0 {
			c.addOutPt(horzEdge, horzEdge.Top)
		}
		c.deleteFromAEL(horzEdge)
	}
}

func (c *Clipper) getNextInAEL(e *ClipperTEdge, direction ClipperDirection) *ClipperTEdge {
	if direction == ClipperLeftToRight {
		return e.NextInAEL
	}
	return e.PrevInAEL
}

func isMaxima(e *ClipperTEdge, y int64) bool {
	return e != nil && e.Top.Y == y && e.NextInLML == nil
}

func isIntermediate(e *ClipperTEdge, y int64) bool {
	return e.Top.Y == y && e.NextInLML != nil
}

func (c *Clipper) getMaximaPair(e *ClipperTEdge) *ClipperTEdge {
	if e.Next.Top == e.Top && e.Next.NextInLML == nil {
		return e.Next
	} else if e.Prev.Top == e.Top && e.Prev.NextInLML == nil {
		return e.Prev
	}
	return nil
}

func (c *Clipper) getMaximaPairEx(e *ClipperTEdge) *ClipperTEdge {
	result := c.getMaximaPair(e)
	if result == nil || result.OutIdx == Skip ||
		(result.NextInAEL == result.PrevInAEL && !isHorizontal(result)) {
		return nil
	}
	return result
}

func (c *Clipper) processEdgesAtTopOfScanbeam(topY int64) {
	e := c.activeEdges
	for e != nil {
		isMaximaEdge := isMaxima(e, topY)
		if isMaximaEdge {
			eMaxPair := c.getMaximaPairEx(e)
			isMaximaEdge = eMaxPair == nil || !isHorizontal(eMaxPair)
		}

		if isMaximaEdge {
			if c.StrictlySimple {
				c.insertMaxima(e.Top.X)
			}
			ePrev := e.PrevInAEL
			c.doMaxima(e)
			if ePrev == nil {
				e = c.activeEdges
			} else {
				e = ePrev.NextInAEL
			}
		} else {
			if isIntermediate(e, topY) && isHorizontal(e.NextInLML) {
				e = c.updateEdgeIntoAEL(e)
				if e.OutIdx >= 0 {
					c.addOutPt(e, e.Bot)
				}
				c.addEdgeToSEL(e)
			} else {
				e.Curr.X = topX(e, topY)
				e.Curr.Y = topY
			}

			if c.StrictlySimple {
				ePrev := e.PrevInAEL
				if e.OutIdx >= 0 && e.WindDelta != 0 && ePrev != nil &&
					ePrev.OutIdx >= 0 && ePrev.Curr.X == e.Curr.X && ePrev.WindDelta != 0 {
					ip := NewClipperIntPoint(e.Curr.X, e.Curr.Y)
					op := c.addOutPt(ePrev, ip)
					op2 := c.addOutPt(e, ip)
					c.addJoin(op, op2, ip)
				}
			}

			e = e.NextInAEL
		}
	}

	c.processHorizontals()
	c.mMaxima = nil

	e = c.activeEdges
	for e != nil {
		if isIntermediate(e, topY) {
			var op *ClipperOutPt
			if e.OutIdx >= 0 {
				op = c.addOutPt(e, e.Top)
			}
			e = c.updateEdgeIntoAEL(e)

			ePrev := e.PrevInAEL
			eNext := e.NextInAEL
			if ePrev != nil && ePrev.Curr.X == e.Bot.X &&
				ePrev.Curr.Y == e.Bot.Y && op != nil &&
				ePrev.OutIdx >= 0 && ePrev.Curr.Y > ePrev.Top.Y &&
				slopesEqual4(e.Curr, e.Top, ePrev.Curr, ePrev.Top, c.useFullRange) &&
				e.WindDelta != 0 && ePrev.WindDelta != 0 {
				op2 := c.addOutPt(ePrev, e.Bot)
				c.addJoin(op, op2, e.Top)
			} else if eNext != nil && eNext.Curr.X == e.Bot.X &&
				eNext.Curr.Y == e.Bot.Y && op != nil &&
				eNext.OutIdx >= 0 && eNext.Curr.Y > eNext.Top.Y &&
				slopesEqual4(e.Curr, e.Top, eNext.Curr, eNext.Top, c.useFullRange) &&
				e.WindDelta != 0 && eNext.WindDelta != 0 {
				op2 := c.addOutPt(eNext, e.Bot)
				c.addJoin(op, op2, e.Top)
			}
		}
		e = e.NextInAEL
	}
}

func (c *Clipper) doMaxima(e *ClipperTEdge) {
	eMaxPair := c.getMaximaPairEx(e)
	if eMaxPair == nil {
		if e.OutIdx >= 0 {
			c.addOutPt(e, e.Top)
		}
		c.deleteFromAEL(e)
		return
	}

	eNext := e.NextInAEL
	for eNext != nil && eNext != eMaxPair {
		c.intersectEdges(e, eNext, e.Top)
		c.swapPositionsInAEL(e, eNext)
		eNext = e.NextInAEL
	}

	if e.OutIdx == Unassigned && eMaxPair.OutIdx == Unassigned {
		c.deleteFromAEL(e)
		c.deleteFromAEL(eMaxPair)
	} else if e.OutIdx >= 0 && eMaxPair.OutIdx >= 0 {
		if e.OutIdx >= 0 {
			c.addLocalMaxPoly(e, eMaxPair, e.Top)
		}
		c.deleteFromAEL(e)
		c.deleteFromAEL(eMaxPair)
	} else if e.WindDelta == 0 {
		if e.OutIdx >= 0 {
			c.addOutPt(e, e.Top)
			e.OutIdx = Unassigned
		}
		c.deleteFromAEL(e)
		if eMaxPair.OutIdx >= 0 {
			c.addOutPt(eMaxPair, e.Top)
			eMaxPair.OutIdx = Unassigned
		}
		c.deleteFromAEL(eMaxPair)
	}
}

func (c *Clipper) insertMaxima(x int64) {
	newMax := &ClipperMaxima{X: x}
	if c.mMaxima == nil {
		c.mMaxima = newMax
	} else if x < c.mMaxima.X {
		newMax.Next = c.mMaxima
		c.mMaxima.Previous = newMax
		c.mMaxima = newMax
	} else {
		m := c.mMaxima
		for m.Next != nil && x >= m.Next.X {
			m = m.Next
		}
		if x == m.X {
			return
		}
		newMax.Next = m.Next
		newMax.Previous = m
		if m.Next != nil {
			m.Next.Previous = newMax
		}
		m.Next = newMax
	}
}

func (c *Clipper) addJoin(op1, op2 *ClipperOutPt, offPt ClipperIntPoint) {
	c.mJoins = append(c.mJoins, &ClipperJoin{OutPt1: op1, OutPt2: op2, OffPt: offPt})
}

func (c *Clipper) addGhostJoin(op *ClipperOutPt, offPt ClipperIntPoint) {
	c.mGhostJoins = append(c.mGhostJoins, &ClipperJoin{OutPt1: op, OffPt: offPt})
}

func horzSegmentsOverlap(seg1a, seg1b, seg2a, seg2b int64) bool {
	if seg1a > seg1b {
		seg1a, seg1b = seg1b, seg1a
	}
	if seg2a > seg2b {
		seg2a, seg2b = seg2b, seg2a
	}
	return seg1a < seg2b && seg2a < seg1b
}

// --- Result building ---

func (c *Clipper) pointCount(pts *ClipperOutPt) int {
	if pts == nil {
		return 0
	}
	result := 0
	p := pts
	for {
		result++
		p = p.Next
		if p == pts {
			break
		}
	}
	return result
}

func (c *Clipper) buildResult() [][]ClipperIntPoint {
	result := make([][]ClipperIntPoint, 0)
	for i := 0; i < len(c.polyOuts); i++ {
		outRec := c.polyOuts[i]
		if outRec.Pts == nil {
			continue
		}
		p := outRec.Pts.Prev
		cnt := c.pointCount(p)
		if cnt < 2 {
			continue
		}
		pg := make([]ClipperIntPoint, 0, cnt)
		for j := 0; j < cnt; j++ {
			pg = append(pg, p.Pt)
			p = p.Prev
		}
		result = append(result, pg)
	}
	return result
}

func (c *Clipper) buildResult2(polytree *ClipperPolyTree) {
	polytree.Clear()
	polytree.AllPolys = make([]*ClipperPolyNode, 0, len(c.polyOuts))

	for i := 0; i < len(c.polyOuts); i++ {
		outRec := c.polyOuts[i]
		cnt := c.pointCount(outRec.Pts)
		if (outRec.IsOpen && cnt < 2) || (!outRec.IsOpen && cnt < 3) {
			continue
		}
		c.fixHoleLinkage(outRec)
		pn := &ClipperPolyNode{Polygon: make([]ClipperIntPoint, 0)}
		polytree.AllPolys = append(polytree.AllPolys, pn)
		outRec.PolyNode = pn
		pn.Polygon = make([]ClipperIntPoint, 0, cnt)
		op := outRec.Pts.Prev
		for j := 0; j < cnt; j++ {
			pn.Polygon = append(pn.Polygon, op.Pt)
			op = op.Prev
		}
	}

	polytree.Children = make([]*ClipperPolyNode, 0, len(c.polyOuts))
	for i := 0; i < len(c.polyOuts); i++ {
		outRec := c.polyOuts[i]
		if outRec.PolyNode == nil {
			continue
		} else if outRec.IsOpen {
			outRec.PolyNode.IsOpen = true
			polytree.AddChild(outRec.PolyNode)
		} else if outRec.FirstLeft != nil && outRec.FirstLeft.PolyNode != nil {
			outRec.FirstLeft.PolyNode.AddChild(outRec.PolyNode)
		} else {
			polytree.AddChild(outRec.PolyNode)
		}
	}
}

func (c *Clipper) fixHoleLinkage(outRec *ClipperOutRec) {
	if outRec.FirstLeft == nil ||
		(outRec.IsHole != outRec.FirstLeft.IsHole && outRec.FirstLeft.Pts != nil) {
		return
	}
	orfl := outRec.FirstLeft
	for orfl != nil && ((orfl.IsHole == outRec.IsHole) || orfl.Pts == nil) {
		orfl = orfl.FirstLeft
	}
	outRec.FirstLeft = orfl
}

func (c *Clipper) fixupOutPolyline(outrec *ClipperOutRec) {
	pp := outrec.Pts
	lastPP := pp.Prev
	for pp != lastPP {
		pp = pp.Next
		if pp.Pt == pp.Prev.Pt {
			if pp == lastPP {
				lastPP = pp.Prev
			}
			tmpPP := pp.Prev
			tmpPP.Next = pp.Next
			pp.Next.Prev = tmpPP
			pp = tmpPP
		}
	}
	if pp == pp.Prev {
		outrec.Pts = nil
	}
}

func (c *Clipper) fixupOutPolygon(outRec *ClipperOutRec) {
	var lastOK *ClipperOutPt
	outRec.BottomPt = nil
	pp := outRec.Pts
	preserveCol := c.PreserveCollinear || c.StrictlySimple
	for {
		if pp.Prev == pp || pp.Prev == pp.Next {
			outRec.Pts = nil
			return
		}
		if (pp.Pt == pp.Next.Pt) || (pp.Pt == pp.Prev.Pt) ||
			(slopesEqual3(pp.Prev.Pt, pp.Pt, pp.Next.Pt, c.useFullRange) &&
				(!preserveCol || !c.pt2BetweenPt1AndPt3(pp.Prev.Pt, pp.Pt, pp.Next.Pt))) {
			lastOK = nil
			pp.Prev.Next = pp.Next
			pp.Next.Prev = pp.Prev
			pp = pp.Prev
		} else if pp == lastOK {
			break
		} else {
			if lastOK == nil {
				lastOK = pp
			}
			pp = pp.Next
		}
	}
	outRec.Pts = pp
}

func (c *Clipper) reversePolyPtLinks(pp *ClipperOutPt) {
	if pp == nil {
		return
	}
	pp1 := pp
	for {
		pp2 := pp1.Next
		pp1.Next = pp1.Prev
		pp1.Prev = pp2
		pp1 = pp2
		if pp1 == pp {
			break
		}
	}
}

func swapSides(edge1, edge2 *ClipperTEdge) {
	side := edge1.Side
	edge1.Side = edge2.Side
	edge2.Side = side
}

func swapPolyIndexes(edge1, edge2 *ClipperTEdge) {
	outIdx := edge1.OutIdx
	edge1.OutIdx = edge2.OutIdx
	edge2.OutIdx = outIdx
}

// --- Join and polygon manipulation ---

func (c *Clipper) appendPolygon(e1, e2 *ClipperTEdge) {
	outRec1 := c.polyOuts[e1.OutIdx]
	outRec2 := c.polyOuts[e2.OutIdx]

	var holeStateRec *ClipperOutRec
	if c.outRec1RightOfOutRec2(outRec1, outRec2) {
		holeStateRec = outRec2
	} else if c.outRec1RightOfOutRec2(outRec2, outRec1) {
		holeStateRec = outRec1
	} else {
		holeStateRec = c.getLowermostRec(outRec1, outRec2)
	}

	p1Lft := outRec1.Pts
	p1Rt := p1Lft.Prev
	p2Lft := outRec2.Pts
	p2Rt := p2Lft.Prev

	if e1.Side == ClipperLeft {
		if e2.Side == ClipperLeft {
			c.reversePolyPtLinks(p2Lft)
			p2Lft.Next = p1Lft
			p1Lft.Prev = p2Lft
			p1Rt.Next = p2Rt
			p2Rt.Prev = p1Rt
			outRec1.Pts = p2Rt
		} else {
			p2Rt.Next = p1Lft
			p1Lft.Prev = p2Rt
			p2Lft.Prev = p1Rt
			p1Rt.Next = p2Lft
			outRec1.Pts = p2Lft
		}
	} else {
		if e2.Side == ClipperRight {
			c.reversePolyPtLinks(p2Lft)
			p1Rt.Next = p2Rt
			p2Rt.Prev = p1Rt
			p2Lft.Next = p1Lft
			p1Lft.Prev = p2Lft
		} else {
			p1Rt.Next = p2Lft
			p2Lft.Prev = p1Rt
			p1Lft.Prev = p2Rt
			p2Rt.Next = p1Lft
		}
	}

	outRec1.BottomPt = nil
	if holeStateRec == outRec2 && outRec2.FirstLeft != outRec1 {
		outRec1.FirstLeft = outRec2.FirstLeft
	}
	if holeStateRec == outRec2 {
		outRec1.IsHole = outRec2.IsHole
	}
	outRec2.Pts = nil
	outRec2.BottomPt = nil
	outRec2.FirstLeft = outRec1

	okIdx := e1.OutIdx
	obsoleteIdx := e2.OutIdx
	e1.OutIdx = Unassigned
	e2.OutIdx = Unassigned

	e := c.activeEdges
	for e != nil {
		if e.OutIdx == obsoleteIdx {
			e.OutIdx = okIdx
			e.Side = e1.Side
			break
		}
		e = e.NextInAEL
	}
	outRec2.Idx = outRec1.Idx
}

func (c *Clipper) getDx(pt1, pt2 ClipperIntPoint) float64 {
	if pt1.Y == pt2.Y {
		return Horizontal
	}
	return float64(pt2.X-pt1.X) / float64(pt2.Y-pt1.Y)
}

func (c *Clipper) firstIsBottomPt(btmPt1, btmPt2 *ClipperOutPt) bool {
	p := btmPt1.Prev
	for p.Pt == btmPt1.Pt && p != btmPt1 {
		p = p.Prev
	}
	dx1p := math.Abs(c.getDx(btmPt1.Pt, p.Pt))
	p = btmPt1.Next
	for p.Pt == btmPt1.Pt && p != btmPt1 {
		p = p.Next
	}
	dx1n := math.Abs(c.getDx(btmPt1.Pt, p.Pt))

	p = btmPt2.Prev
	for p.Pt == btmPt2.Pt && p != btmPt2 {
		p = p.Prev
	}
	dx2p := math.Abs(c.getDx(btmPt2.Pt, p.Pt))
	p = btmPt2.Next
	for p.Pt == btmPt2.Pt && p != btmPt2 {
		p = p.Next
	}
	dx2n := math.Abs(c.getDx(btmPt2.Pt, p.Pt))

	if math.Max(dx1p, dx1n) == math.Max(dx2p, dx2n) && math.Min(dx1p, dx1n) == math.Min(dx2p, dx2n) {
		return c.areaOutPt(btmPt1) > 0
	}
	return (dx1p >= dx2p && dx1p >= dx2n) || (dx1n >= dx2p && dx1n >= dx2n)
}

func (c *Clipper) getBottomPt(pp *ClipperOutPt) *ClipperOutPt {
	var dups *ClipperOutPt
	p := pp.Next
	for p != pp {
		if p.Pt.Y > pp.Pt.Y {
			pp = p
			dups = nil
		} else if p.Pt.Y == pp.Pt.Y && p.Pt.X <= pp.Pt.X {
			if p.Pt.X < pp.Pt.X {
				dups = nil
				pp = p
			} else if p.Next != pp && p.Prev != pp {
				dups = p
			}
		}
		p = p.Next
	}
	if dups != nil {
		for dups != p {
			if !c.firstIsBottomPt(p, dups) {
				pp = dups
			}
			dups = dups.Next
			for dups.Pt != pp.Pt {
				dups = dups.Next
			}
		}
	}
	return pp
}

func (c *Clipper) getLowermostRec(outRec1, outRec2 *ClipperOutRec) *ClipperOutRec {
	if outRec1.BottomPt == nil {
		outRec1.BottomPt = c.getBottomPt(outRec1.Pts)
	}
	if outRec2.BottomPt == nil {
		outRec2.BottomPt = c.getBottomPt(outRec2.Pts)
	}
	bPt1 := outRec1.BottomPt
	bPt2 := outRec2.BottomPt
	if bPt1.Pt.Y > bPt2.Pt.Y {
		return outRec1
	} else if bPt1.Pt.Y < bPt2.Pt.Y {
		return outRec2
	} else if bPt1.Pt.X < bPt2.Pt.X {
		return outRec1
	} else if bPt1.Pt.X > bPt2.Pt.X {
		return outRec2
	} else if bPt1.Next == bPt1 {
		return outRec2
	} else if bPt2.Next == bPt2 {
		return outRec1
	} else if c.firstIsBottomPt(bPt1, bPt2) {
		return outRec1
	}
	return outRec2
}

func (c *Clipper) outRec1RightOfOutRec2(outRec1, outRec2 *ClipperOutRec) bool {
	for {
		outRec1 = outRec1.FirstLeft
		if outRec1 == outRec2 {
			return true
		}
		if outRec1 == nil {
			break
		}
	}
	return false
}

func (c *Clipper) getOutRec(idx int) *ClipperOutRec {
	outrec := c.polyOuts[idx]
	for outrec != c.polyOuts[outrec.Idx] {
		outrec = c.polyOuts[outrec.Idx]
	}
	return outrec
}

func (c *Clipper) dupOutPt(outPt *ClipperOutPt, insertAfter bool) *ClipperOutPt {
	result := &ClipperOutPt{Pt: outPt.Pt, Index: outPt.Index}
	if insertAfter {
		result.Next = outPt.Next
		result.Prev = outPt
		outPt.Next.Prev = result
		outPt.Next = result
	} else {
		result.Prev = outPt.Prev
		result.Next = outPt
		outPt.Prev.Next = result
		outPt.Prev = result
	}
	return result
}

func getOverlap(a1, a2, b1, b2 int64) (left, right int64, ok bool) {
	var la, ra int64
	if a1 < a2 {
		if b1 < b2 {
			la = maxInt64(a1, b1)
			ra = minInt64(a2, b2)
		} else {
			la = maxInt64(a1, b2)
			ra = minInt64(a2, b1)
		}
	} else {
		if b1 < b2 {
			la = maxInt64(a2, b1)
			ra = minInt64(a1, b2)
		} else {
			la = maxInt64(a2, b2)
			ra = minInt64(a1, b1)
		}
	}
	return la, ra, la < ra
}

func (c *Clipper) joinHorz(op1, op1b, op2, op2b *ClipperOutPt, pt ClipperIntPoint, discardLeft bool) bool {
	dir1 := ClipperRightToLeft
	if op1.Pt.X <= op1b.Pt.X {
		dir1 = ClipperLeftToRight
	}
	dir2 := ClipperRightToLeft
	if op2.Pt.X <= op2b.Pt.X {
		dir2 = ClipperLeftToRight
	}
	if dir1 == dir2 {
		return false
	}

	if dir1 == ClipperLeftToRight {
		for op1.Next.Pt.X <= pt.X && op1.Next.Pt.X >= op1.Pt.X && op1.Next.Pt.Y == pt.Y {
			op1 = op1.Next
		}
		if discardLeft && op1.Pt.X != pt.X {
			op1 = op1.Next
		}
		op1b = c.dupOutPt(op1, !discardLeft)
		if op1b.Pt != pt {
			op1 = op1b
			op1.Pt = pt
			op1b = c.dupOutPt(op1, !discardLeft)
		}
	} else {
		for op1.Next.Pt.X >= pt.X && op1.Next.Pt.X <= op1.Pt.X && op1.Next.Pt.Y == pt.Y {
			op1 = op1.Next
		}
		if !discardLeft && op1.Pt.X != pt.X {
			op1 = op1.Next
		}
		op1b = c.dupOutPt(op1, discardLeft)
		if op1b.Pt != pt {
			op1 = op1b
			op1.Pt = pt
			op1b = c.dupOutPt(op1, discardLeft)
		}
	}

	if dir2 == ClipperLeftToRight {
		for op2.Next.Pt.X <= pt.X && op2.Next.Pt.X >= op2.Pt.X && op2.Next.Pt.Y == pt.Y {
			op2 = op2.Next
		}
		if discardLeft && op2.Pt.X != pt.X {
			op2 = op2.Next
		}
		op2b = c.dupOutPt(op2, !discardLeft)
		if op2b.Pt != pt {
			op2 = op2b
			op2.Pt = pt
			op2b = c.dupOutPt(op2, !discardLeft)
		}
	} else {
		for op2.Next.Pt.X >= pt.X && op2.Next.Pt.X <= op2.Pt.X && op2.Next.Pt.Y == pt.Y {
			op2 = op2.Next
		}
		if !discardLeft && op2.Pt.X != pt.X {
			op2 = op2.Next
		}
		op2b = c.dupOutPt(op2, discardLeft)
		if op2b.Pt != pt {
			op2 = op2b
			op2.Pt = pt
			op2b = c.dupOutPt(op2, discardLeft)
		}
	}

	if (dir1 == ClipperLeftToRight) == discardLeft {
		op1.Prev = op2
		op2.Next = op1
		op1b.Next = op2b
		op2b.Prev = op1b
	} else {
		op1.Next = op2
		op2.Prev = op1
		op1b.Prev = op2b
		op2b.Next = op1b
	}
	return true
}

func (c *Clipper) joinPoints(j *ClipperJoin, outRec1, outRec2 *ClipperOutRec) bool {
	op1 := j.OutPt1
	var op1b *ClipperOutPt
	op2 := j.OutPt2
	var op2b *ClipperOutPt

	isHorizontalJoin := j.OutPt1.Pt.Y == j.OffPt.Y

	if isHorizontalJoin && j.OffPt == j.OutPt1.Pt && j.OffPt == j.OutPt2.Pt {
		if outRec1 != outRec2 {
			return false
		}
		op1b = j.OutPt1.Next
		for op1b != op1 && op1b.Pt == j.OffPt {
			op1b = op1b.Next
		}
		reverse1 := op1b.Pt.Y > j.OffPt.Y
		op2b = j.OutPt2.Next
		for op2b != op2 && op2b.Pt == j.OffPt {
			op2b = op2b.Next
		}
		reverse2 := op2b.Pt.Y > j.OffPt.Y
		if reverse1 == reverse2 {
			return false
		}
		if reverse1 {
			op1b = c.dupOutPt(op1, false)
			op2b = c.dupOutPt(op2, true)
			op1.Prev = op2
			op2.Next = op1
			op1b.Next = op2b
			op2b.Prev = op1b
		} else {
			op1b = c.dupOutPt(op1, true)
			op2b = c.dupOutPt(op2, false)
			op1.Next = op2
			op2.Prev = op1
			op1b.Prev = op2b
			op2b.Next = op1b
		}
		j.OutPt1 = op1
		j.OutPt2 = op1b
		return true
	} else if isHorizontalJoin {
		op1b = op1
		for op1.Prev.Pt.Y == op1.Pt.Y && op1.Prev != op1b && op1.Prev != op2 {
			op1 = op1.Prev
		}
		for op1b.Next.Pt.Y == op1b.Pt.Y && op1b.Next != op1 && op1b.Next != op2 {
			op1b = op1b.Next
		}
		if op1b.Next == op1 || op1b.Next == op2 {
			return false
		}

		op2b = op2
		for op2.Prev.Pt.Y == op2.Pt.Y && op2.Prev != op2b && op2.Prev != op1b {
			op2 = op2.Prev
		}
		for op2b.Next.Pt.Y == op2b.Pt.Y && op2b.Next != op2 && op2b.Next != op1 {
			op2b = op2b.Next
		}
		if op2b.Next == op2 || op2b.Next == op1 {
			return false
		}

		left, right, ok := getOverlap(op1.Pt.X, op1b.Pt.X, op2.Pt.X, op2b.Pt.X)
		if !ok {
			return false
		}

		var pt ClipperIntPoint
		var discardLeftSide bool
		if op1.Pt.X >= left && op1.Pt.X <= right {
			pt = op1.Pt
			discardLeftSide = op1.Pt.X > op1b.Pt.X
		} else if op2.Pt.X >= left && op2.Pt.X <= right {
			pt = op2.Pt
			discardLeftSide = op2.Pt.X > op2b.Pt.X
		} else if op1b.Pt.X >= left && op1b.Pt.X <= right {
			pt = op1b.Pt
			discardLeftSide = op1b.Pt.X > op1.Pt.X
		} else {
			pt = op2b.Pt
			discardLeftSide = op2b.Pt.X > op2.Pt.X
		}
		j.OutPt1 = op1
		j.OutPt2 = op2
		return c.joinHorz(op1, op1b, op2, op2b, pt, discardLeftSide)
	} else {
		op1b = op1.Next
		for op1b.Pt == op1.Pt && op1b != op1 {
			op1b = op1b.Next
		}
		reverse1 := op1b.Pt.Y > op1.Pt.Y || !slopesEqual3(op1.Pt, op1b.Pt, j.OffPt, c.useFullRange)
		if reverse1 {
			op1b = op1.Prev
			for op1b.Pt == op1.Pt && op1b != op1 {
				op1b = op1b.Prev
			}
			if op1b.Pt.Y > op1.Pt.Y || !slopesEqual3(op1.Pt, op1b.Pt, j.OffPt, c.useFullRange) {
				return false
			}
		}

		op2b = op2.Next
		for op2b.Pt == op2.Pt && op2b != op2 {
			op2b = op2b.Next
		}
		reverse2 := op2b.Pt.Y > op2.Pt.Y || !slopesEqual3(op2.Pt, op2b.Pt, j.OffPt, c.useFullRange)
		if reverse2 {
			op2b = op2.Prev
			for op2b.Pt == op2.Pt && op2b != op2 {
				op2b = op2b.Prev
			}
			if op2b.Pt.Y > op2.Pt.Y || !slopesEqual3(op2.Pt, op2b.Pt, j.OffPt, c.useFullRange) {
				return false
			}
		}

		if (op1b == op1) || (op2b == op2) || (op1b == op2b) ||
			(outRec1 == outRec2 && reverse1 == reverse2) {
			return false
		}

		if reverse1 {
			op1b = c.dupOutPt(op1, false)
			op2b = c.dupOutPt(op2, true)
			op1.Prev = op2
			op2.Next = op1
			op1b.Next = op2b
			op2b.Prev = op1b
		} else {
			op1b = c.dupOutPt(op1, true)
			op2b = c.dupOutPt(op2, false)
			op1.Next = op2
			op2.Prev = op1
			op1b.Prev = op2b
			op2b.Next = op1b
		}
		j.OutPt1 = op1
		j.OutPt2 = op1b
		return true
	}
}

func (c *Clipper) joinCommonEdges() {
	for _, join := range c.mJoins {
		outRec1 := c.getOutRec(join.OutPt1.Index)
		outRec2 := c.getOutRec(join.OutPt2.Index)

		if outRec1.Pts == nil || outRec2.Pts == nil {
			continue
		}
		if outRec1.IsOpen || outRec2.IsOpen {
			continue
		}

		var holeStateRec *ClipperOutRec
		if outRec1 == outRec2 {
			holeStateRec = outRec1
		} else if c.outRec1RightOfOutRec2(outRec1, outRec2) {
			holeStateRec = outRec2
		} else if c.outRec1RightOfOutRec2(outRec2, outRec1) {
			holeStateRec = outRec1
		} else {
			holeStateRec = c.getLowermostRec(outRec1, outRec2)
		}

		if !c.joinPoints(join, outRec1, outRec2) {
			continue
		}

		if outRec1 == outRec2 {
			outRec1.Pts = join.OutPt1
			outRec1.BottomPt = nil
			outRec2 = c.createOutRec()
			outRec2.Pts = join.OutPt2
			c.updateOutPtIdxs(outRec2)

			if poly2ContainsPoly1(outRec2.Pts, outRec1.Pts) {
				outRec2.IsHole = !outRec1.IsHole
				outRec2.FirstLeft = outRec1
				if c.mUsingPolyTree {
					c.fixupFirstLefts2(outRec2, outRec1)
				}
				if (outRec2.IsHole != c.ReverseSolution) == (c.areaOutRec(outRec2) > 0) {
					c.reversePolyPtLinks(outRec2.Pts)
				}
			} else if poly2ContainsPoly1(outRec1.Pts, outRec2.Pts) {
				outRec2.IsHole = outRec1.IsHole
				outRec1.IsHole = !outRec2.IsHole
				outRec2.FirstLeft = outRec1.FirstLeft
				outRec1.FirstLeft = outRec2
				if c.mUsingPolyTree {
					c.fixupFirstLefts2(outRec1, outRec2)
				}
				if (outRec1.IsHole != c.ReverseSolution) == (c.areaOutRec(outRec1) > 0) {
					c.reversePolyPtLinks(outRec1.Pts)
				}
			} else {
				outRec2.IsHole = outRec1.IsHole
				outRec2.FirstLeft = outRec1.FirstLeft
				if c.mUsingPolyTree {
					c.fixupFirstLefts1(outRec1, outRec2)
				}
			}
		} else {
			outRec2.Pts = nil
			outRec2.BottomPt = nil
			outRec2.Idx = outRec1.Idx
			outRec1.IsHole = holeStateRec.IsHole
			if holeStateRec == outRec2 {
				outRec1.FirstLeft = outRec2.FirstLeft
			}
			outRec2.FirstLeft = outRec1
			if c.mUsingPolyTree {
				c.fixupFirstLefts3(outRec2, outRec1)
			}
		}
	}
}

func (c *Clipper) updateOutPtIdxs(outrec *ClipperOutRec) {
	op := outrec.Pts
	for {
		op.Index = outrec.Idx
		op = op.Prev
		if op == outrec.Pts {
			break
		}
	}
}

func pointInPolygonOp(pt ClipperIntPoint, op *ClipperOutPt) int {
	result := 0
	startOp := op
	ptx, pty := pt.X, pt.Y
	poly0x, poly0y := op.Pt.X, op.Pt.Y
	for {
		op = op.Next
		poly1x, poly1y := op.Pt.X, op.Pt.Y

		if poly1y == pty {
			if poly1x == ptx || (poly0y == pty && ((poly1x > ptx) == (poly0x < ptx))) {
				return -1
			}
		}
		if (poly0y < pty) != (poly1y < pty) {
			if poly0x >= ptx {
				if poly1x > ptx {
					result = 1 - result
				} else {
					d := float64(poly0x-ptx)*float64(poly1y-pty) - float64(poly1x-ptx)*float64(poly0y-pty)
					if d == 0 {
						return -1
					}
					if (d > 0) == (poly1y > poly0y) {
						result = 1 - result
					}
				}
			} else if poly1x > ptx {
				d := float64(poly0x-ptx)*float64(poly1y-pty) - float64(poly1x-ptx)*float64(poly0y-pty)
				if d == 0 {
					return -1
				}
				if (d > 0) == (poly1y > poly0y) {
					result = 1 - result
				}
			}
		}
		poly0x, poly0y = poly1x, poly1y
		if startOp == op {
			break
		}
	}
	return result
}

func poly2ContainsPoly1(outPt1, outPt2 *ClipperOutPt) bool {
	op := outPt1
	for {
		res := pointInPolygonOp(op.Pt, outPt2)
		if res >= 0 {
			return res > 0
		}
		op = op.Next
		if op == outPt1 {
			break
		}
	}
	return true
}

func parseFirstLeft(firstLeft *ClipperOutRec) *ClipperOutRec {
	for firstLeft != nil && firstLeft.Pts == nil {
		firstLeft = firstLeft.FirstLeft
	}
	return firstLeft
}

func (c *Clipper) fixupFirstLefts1(oldOutRec, newOutRec *ClipperOutRec) {
	for _, outRec := range c.polyOuts {
		firstLeft := parseFirstLeft(outRec.FirstLeft)
		if outRec.Pts != nil && firstLeft == oldOutRec {
			if poly2ContainsPoly1(outRec.Pts, newOutRec.Pts) {
				outRec.FirstLeft = newOutRec
			}
		}
	}
}

func (c *Clipper) fixupFirstLefts2(innerOutRec, outerOutRec *ClipperOutRec) {
	orfl := outerOutRec.FirstLeft
	for _, outRec := range c.polyOuts {
		if outRec.Pts == nil || outRec == outerOutRec || outRec == innerOutRec {
			continue
		}
		firstLeft := parseFirstLeft(outRec.FirstLeft)
		if firstLeft != orfl && firstLeft != innerOutRec && firstLeft != outerOutRec {
			continue
		}
		if poly2ContainsPoly1(outRec.Pts, innerOutRec.Pts) {
			outRec.FirstLeft = innerOutRec
		} else if poly2ContainsPoly1(outRec.Pts, outerOutRec.Pts) {
			outRec.FirstLeft = outerOutRec
		} else if outRec.FirstLeft == innerOutRec || outRec.FirstLeft == outerOutRec {
			outRec.FirstLeft = orfl
		}
	}
}

func (c *Clipper) fixupFirstLefts3(oldOutRec, newOutRec *ClipperOutRec) {
	for _, outRec := range c.polyOuts {
		firstLeft := parseFirstLeft(outRec.FirstLeft)
		if outRec.Pts != nil && firstLeft == oldOutRec {
			outRec.FirstLeft = newOutRec
		}
	}
}

func (c *Clipper) doSimplePolygons() {
	i := 0
	for i < len(c.polyOuts) {
		outrec := c.polyOuts[i]
		i++
		op := outrec.Pts
		if op == nil || outrec.IsOpen {
			continue
		}
		for {
			op2 := op.Next
			for op2 != outrec.Pts {
				if op.Pt == op2.Pt && op2.Next != op && op2.Prev != op {
					op3 := op.Prev
					op4 := op2.Prev
					op.Prev = op4
					op4.Next = op
					op2.Prev = op3
					op3.Next = op2

					outrec.Pts = op
					outrec2 := c.createOutRec()
					outrec2.Pts = op2
					c.updateOutPtIdxs(outrec2)

					if poly2ContainsPoly1(outrec2.Pts, outrec.Pts) {
						outrec2.IsHole = !outrec.IsHole
						outrec2.FirstLeft = outrec
						if c.mUsingPolyTree {
							c.fixupFirstLefts2(outrec2, outrec)
						}
					} else if poly2ContainsPoly1(outrec.Pts, outrec2.Pts) {
						outrec2.IsHole = outrec.IsHole
						outrec.IsHole = !outrec2.IsHole
						outrec2.FirstLeft = outrec.FirstLeft
						outrec.FirstLeft = outrec2
						if c.mUsingPolyTree {
							c.fixupFirstLefts2(outrec, outrec2)
						}
					} else {
						outrec2.IsHole = outrec.IsHole
						outrec2.FirstLeft = outrec.FirstLeft
						if c.mUsingPolyTree {
							c.fixupFirstLefts1(outrec, outrec2)
						}
					}
					op2 = op
				}
				op2 = op2.Next
			}
			op = op.Next
			if op == outrec.Pts {
				break
			}
		}
	}
}

func (c *Clipper) areaOutRec(outRec *ClipperOutRec) float64 {
	return c.areaOutPt(outRec.Pts)
}

func (c *Clipper) areaOutPt(op *ClipperOutPt) float64 {
	opFirst := op
	if op == nil {
		return 0
	}
	a := 0.0
	for {
		a += float64(op.Prev.Pt.X+op.Pt.X) * float64(op.Prev.Pt.Y-op.Pt.Y)
		op = op.Next
		if op == opFirst {
			break
		}
	}
	return a * 0.5
}

// --- ClipperBase inherited methods ---

func (c *Clipper) reset() {
	c.currentLM = c.minimaList
	if c.currentLM == nil {
		return
	}

	c.scanbeam = nil
	lm := c.minimaList
	for lm != nil {
		c.insertScanbeam(lm.Y)
		e := lm.LeftBound
		if e != nil {
			e.Curr = e.Bot
			e.OutIdx = Unassigned
		}
		e = lm.RightBound
		if e != nil {
			e.Curr = e.Bot
			e.OutIdx = Unassigned
		}
		lm = lm.Next
	}
	c.activeEdges = nil
}

func isHorizontal(e *ClipperTEdge) bool {
	return e.Delta.Y == 0
}

// slopesEqual2h checks slope equality for a horizontal edge and another edge.
func slopesEqual2h(e1, e2 *ClipperTEdge, useFullRange bool) bool {
	if useFullRange {
		return Int128Mul(e1.Delta.Y, e2.Delta.X) == Int128Mul(e1.Delta.X, e2.Delta.Y)
	}
	return int64(e1.Delta.Y)*e2.Delta.X == e1.Delta.X*int64(e2.Delta.Y)
}

func slopesEqual2(e1, e2 *ClipperTEdge, useFullRange bool) bool {
	if useFullRange {
		return Int128Mul(e1.Delta.Y, e2.Delta.X) == Int128Mul(e1.Delta.X, e2.Delta.Y)
	}
	return int64(e1.Delta.Y)*e2.Delta.X == e1.Delta.X*int64(e2.Delta.Y)
}

func slopesEqual3(pt1, pt2, pt3 ClipperIntPoint, useFullRange bool) bool {
	if useFullRange {
		return Int128Mul(pt1.Y-pt2.Y, pt2.X-pt3.X) == Int128Mul(pt1.X-pt2.X, pt2.Y-pt3.Y)
	}
	return (pt1.Y-pt2.Y)*(pt2.X-pt3.X)-(pt1.X-pt2.X)*(pt2.Y-pt3.Y) == 0
}

func slopesEqual4(pt1, pt2, pt3, pt4 ClipperIntPoint, useFullRange bool) bool {
	if useFullRange {
		return Int128Mul(pt1.Y-pt2.Y, pt3.X-pt4.X) == Int128Mul(pt1.X-pt2.X, pt3.Y-pt4.Y)
	}
	return (pt1.Y-pt2.Y)*(pt3.X-pt4.X)-(pt1.X-pt2.X)*(pt3.Y-pt4.Y) == 0
}

func (c *Clipper) rangeTest(pt ClipperIntPoint) {
	if c.useFullRange {
		if pt.X > hiRange || pt.Y > hiRange || -pt.X > hiRange || -pt.Y > hiRange {
			return // Silently ignore out-of-range for Go safety
		}
	} else if pt.X > loRange || pt.Y > loRange || -pt.X > loRange || -pt.Y > loRange {
		c.useFullRange = true
		c.rangeTest(pt)
	}
}

func (c *Clipper) initEdge(e, eNext, ePrev *ClipperTEdge, pt ClipperIntPoint) {
	e.Next = eNext
	e.Prev = ePrev
	e.Curr = pt
	e.OutIdx = Unassigned
}

func (c *Clipper) initEdge2(e *ClipperTEdge, polyType ClipperPolyType) {
	if e.Curr.Y >= e.Next.Curr.Y {
		e.Bot = e.Curr
		e.Top = e.Next.Curr
	} else {
		e.Top = e.Curr
		e.Bot = e.Next.Curr
	}
	c.setDx(e)
	e.PolyTyp = polyType
}

func (c *Clipper) setDx(e *ClipperTEdge) {
	e.Delta.X = e.Top.X - e.Bot.X
	e.Delta.Y = e.Top.Y - e.Bot.Y
	if e.Delta.Y == 0 {
		e.Dx = Horizontal
	} else {
		e.Dx = float64(e.Delta.X) / float64(e.Delta.Y)
	}
}

func (c *Clipper) insertLocalMinima(newLm *ClipperLocalMinima) {
	if c.minimaList == nil {
		c.minimaList = newLm
	} else if newLm.Y >= c.minimaList.Y {
		newLm.Next = c.minimaList
		c.minimaList = newLm
	} else {
		tmpLm := c.minimaList
		for tmpLm.Next != nil && newLm.Y < tmpLm.Next.Y {
			tmpLm = tmpLm.Next
		}
		newLm.Next = tmpLm.Next
		tmpLm.Next = newLm
	}
}

// popLocalMinima pops the current local minima at Y and stores it in lm.
func (c *Clipper) popLocalMinima(y int64, lm **ClipperLocalMinima) bool {
	if c.currentLM != nil && c.currentLM.Y == y {
		*lm = c.currentLM
		c.currentLM = c.currentLM.Next
		return true
	}
	return false
}

func (c *Clipper) localMinimaPending() bool {
	return c.currentLM != nil
}

func (c *Clipper) reverseHorizontal(e *ClipperTEdge) {
	e.Top.X, e.Bot.X = e.Bot.X, e.Top.X
}

func (c *Clipper) insertScanbeam(y int64) {
	if c.scanbeam == nil {
		c.scanbeam = &ClipperScanbeam{Next: nil, Y: y}
	} else if y > c.scanbeam.Y {
		c.scanbeam = &ClipperScanbeam{Y: y, Next: c.scanbeam}
	} else {
		sb2 := c.scanbeam
		for sb2.Next != nil && y <= sb2.Next.Y {
			sb2 = sb2.Next
		}
		if y == sb2.Y {
			return
		}
		sb2.Next = &ClipperScanbeam{Y: y, Next: sb2.Next}
	}
}

func (c *Clipper) popScanbeam() (int64, bool) {
	if c.scanbeam == nil {
		return 0, false
	}
	y := c.scanbeam.Y
	c.scanbeam = c.scanbeam.Next
	return y, true
}

func (c *Clipper) createOutRec() *ClipperOutRec {
	result := &ClipperOutRec{
		Idx:      Unassigned,
		IsHole:   false,
		IsOpen:   false,
		FirstLeft: nil,
		Pts:       nil,
		BottomPt:  nil,
		PolyNode:  nil,
	}
	c.polyOuts = append(c.polyOuts, result)
	result.Idx = len(c.polyOuts) - 1
	return result
}

func (c *Clipper) disposeOutRec(index int) {
	if index < len(c.polyOuts) {
		outRec := c.polyOuts[index]
		outRec.Pts = nil
		c.polyOuts[index] = nil
	}
}

func (c *Clipper) updateEdgeIntoAEL(e *ClipperTEdge) *ClipperTEdge {
	if e.NextInLML == nil {
		panic("UpdateEdgeIntoAEL: invalid call")
	}
	aelPrev := e.PrevInAEL
	aelNext := e.NextInAEL
	e.NextInLML.OutIdx = e.OutIdx
	if aelPrev != nil {
		aelPrev.NextInAEL = e.NextInLML
	} else {
		c.activeEdges = e.NextInLML
	}
	if aelNext != nil {
		aelNext.PrevInAEL = e.NextInLML
	}
	e.NextInLML.Side = e.Side
	e.NextInLML.WindDelta = e.WindDelta
	e.NextInLML.WindCnt = e.WindCnt
	e.NextInLML.WindCnt2 = e.WindCnt2
	result := e.NextInLML
	result.Curr = result.Bot
	result.PrevInAEL = aelPrev
	result.NextInAEL = aelNext
	if !isHorizontal(result) {
		c.insertScanbeam(result.Top.Y)
	}
	return result
}


func (c *Clipper) swapPositionsInAEL(edge1, edge2 *ClipperTEdge) {
	if edge1.NextInAEL == edge1.PrevInAEL || edge2.NextInAEL == edge2.PrevInAEL {
		return
	}

	if edge1.NextInAEL == edge2 {
		next := edge2.NextInAEL
		if next != nil {
			next.PrevInAEL = edge1
		}
		prev := edge1.PrevInAEL
		if prev != nil {
			prev.NextInAEL = edge2
		}
		edge2.PrevInAEL = prev
		edge2.NextInAEL = edge1
		edge1.PrevInAEL = edge2
		edge1.NextInAEL = next
	} else if edge2.NextInAEL == edge1 {
		next := edge1.NextInAEL
		if next != nil {
			next.PrevInAEL = edge2
		}
		prev := edge2.PrevInAEL
		if prev != nil {
			prev.NextInAEL = edge1
		}
		edge1.PrevInAEL = prev
		edge1.NextInAEL = edge2
		edge2.PrevInAEL = edge1
		edge2.NextInAEL = next
	} else {
		next := edge1.NextInAEL
		prev := edge1.PrevInAEL
		edge1.NextInAEL = edge2.NextInAEL
		if edge1.NextInAEL != nil {
			edge1.NextInAEL.PrevInAEL = edge1
		}
		edge1.PrevInAEL = edge2.PrevInAEL
		if edge1.PrevInAEL != nil {
			edge1.PrevInAEL.NextInAEL = edge1
		}
		edge2.NextInAEL = next
		if edge2.NextInAEL != nil {
			edge2.NextInAEL.PrevInAEL = edge2
		}
		edge2.PrevInAEL = prev
		if edge2.PrevInAEL != nil {
			edge2.PrevInAEL.NextInAEL = edge2
		}
	}

	if edge1.PrevInAEL == nil {
		c.activeEdges = edge1
	} else if edge2.PrevInAEL == nil {
		c.activeEdges = edge2
	}
}

func (c *Clipper) deleteFromAEL(e *ClipperTEdge) {
	aelPrev := e.PrevInAEL
	aelNext := e.NextInAEL
	if aelPrev == nil && aelNext == nil && e != c.activeEdges {
		return // already deleted
	}
	if aelPrev != nil {
		aelPrev.NextInAEL = aelNext
	} else {
		c.activeEdges = aelNext
	}
	if aelNext != nil {
		aelNext.PrevInAEL = aelPrev
	}
	e.NextInAEL = nil
	e.PrevInAEL = nil
}

func (c *Clipper) removeEdge(e *ClipperTEdge) *ClipperTEdge {
	e.Prev.Next = e.Next
	e.Next.Prev = e.Prev
	result := e.Next
	e.Prev = nil // flag as removed
	return result
}

func (c *Clipper) pt2BetweenPt1AndPt3(pt1, pt2, pt3 ClipperIntPoint) bool {
	if pt1 == pt3 || pt1 == pt2 || pt3 == pt2 {
		return false
	}
	if pt1.X != pt3.X {
		return (pt2.X > pt1.X) == (pt2.X < pt3.X)
	}
	return (pt2.Y > pt1.Y) == (pt2.Y < pt3.Y)
}

func (c *Clipper) findNextLocMin(e *ClipperTEdge) *ClipperTEdge {
	for {
		for e.Bot != e.Prev.Bot || e.Curr == e.Top {
			e = e.Next
		}
		if e.Dx != Horizontal && e.Prev.Dx != Horizontal {
			break
		}
		for e.Prev.Dx == Horizontal {
			e = e.Prev
		}
		e2 := e
		for e.Dx == Horizontal {
			e = e.Next
		}
		if e.Top.Y == e.Prev.Bot.Y {
			continue
		}
		if e2.Prev.Bot.X < e.Bot.X {
			e = e2
		}
		break
	}
	return e
}

func (c *Clipper) processBound(e *ClipperTEdge, leftBoundIsForward bool) *ClipperTEdge {
	result := e
	if result.OutIdx == Skip {
		if leftBoundIsForward {
			for e.Top.Y == e.Next.Bot.Y {
				e = e.Next
			}
			for e != result && e.Dx == Horizontal {
				e = e.Prev
			}
		} else {
			for e.Top.Y == e.Prev.Bot.Y {
				e = e.Prev
			}
			for e != result && e.Dx == Horizontal {
				e = e.Next
			}
		}
		if e == result {
			if leftBoundIsForward {
				result = e.Next
			} else {
				result = e.Prev
			}
		} else {
			if leftBoundIsForward {
				e = result.Next
			} else {
				e = result.Prev
			}
			locMin := &ClipperLocalMinima{Y: e.Bot.Y, LeftBound: nil, RightBound: e}
			e.WindDelta = 0
			result = c.processBound(e, leftBoundIsForward)
			c.insertLocalMinima(locMin)
		}
		return result
	}

	if e.Dx == Horizontal {
		var eStart *ClipperTEdge
		if leftBoundIsForward {
			eStart = e.Prev
		} else {
			eStart = e.Next
		}
		if eStart.Dx == Horizontal {
			if eStart.Bot.X != e.Bot.X && eStart.Top.X != e.Bot.X {
				c.reverseHorizontal(e)
			}
		} else if eStart.Bot.X != e.Bot.X {
			c.reverseHorizontal(e)
		}
	}

	eStart := e
	if leftBoundIsForward {
		for result.Top.Y == result.Next.Bot.Y && result.Next.OutIdx != Skip {
			result = result.Next
		}
		if result.Dx == Horizontal && result.Next.OutIdx != Skip {
			horz := result
			for horz.Prev.Dx == Horizontal {
				horz = horz.Prev
			}
			if horz.Prev.Top.X > result.Next.Top.X {
				result = horz.Prev
			}
		}
		for e != result {
			e.NextInLML = e.Next
			if e.Dx == Horizontal && e != eStart && e.Bot.X != e.Prev.Top.X {
				c.reverseHorizontal(e)
			}
			e = e.Next
		}
		if e.Dx == Horizontal && e != eStart && e.Bot.X != e.Prev.Top.X {
			c.reverseHorizontal(e)
		}
		result = result.Next
	} else {
		for result.Top.Y == result.Prev.Bot.Y && result.Prev.OutIdx != Skip {
			result = result.Prev
		}
		if result.Dx == Horizontal && result.Prev.OutIdx != Skip {
			horz := result
			for horz.Next.Dx == Horizontal {
				horz = horz.Next
			}
			if horz.Next.Top.X == result.Prev.Top.X || horz.Next.Top.X > result.Prev.Top.X {
				result = horz.Next
			}
		}
		for e != result {
			e.NextInLML = e.Prev
			if e.Dx == Horizontal && e != eStart && e.Bot.X != e.Next.Top.X {
				c.reverseHorizontal(e)
			}
			e = e.Prev
		}
		if e.Dx == Horizontal && e != eStart && e.Bot.X != e.Next.Top.X {
			c.reverseHorizontal(e)
		}
		result = result.Prev
	}
	return result
}

// --- Helper functions ---

func reverseSlice(s []ClipperIntPoint) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func pointsAreClose(pt1, pt2 ClipperIntPoint, distSqrd float64) bool {
	dx := float64(pt1.X - pt2.X)
	dy := float64(pt1.Y - pt2.Y)
	return dx*dx+dy*dy <= distSqrd
}

func distanceFromLineSqrd(pt, ln1, ln2 ClipperIntPoint) float64 {
	a := float64(ln1.Y - ln2.Y)
	b := float64(ln2.X - ln1.X)
	c := a*float64(ln1.X) + b*float64(ln1.Y)
	c = a*float64(pt.X) + b*float64(pt.Y) - c
	return (c * c) / (a*a + b*b)
}

func slopesNearCollinear(pt1, pt2, pt3 ClipperIntPoint, distSqrd float64) bool {
	if math.Abs(float64(pt1.X-pt2.X)) > math.Abs(float64(pt1.Y-pt2.Y)) {
		if (pt1.X > pt2.X) == (pt1.X < pt3.X) {
			return distanceFromLineSqrd(pt1, pt2, pt3) < distSqrd
		} else if (pt2.X > pt1.X) == (pt2.X < pt3.X) {
			return distanceFromLineSqrd(pt2, pt1, pt3) < distSqrd
		}
		return distanceFromLineSqrd(pt3, pt1, pt2) < distSqrd
	} else {
		if (pt1.Y > pt2.Y) == (pt1.Y < pt3.Y) {
			return distanceFromLineSqrd(pt1, pt2, pt3) < distSqrd
		} else if (pt2.Y > pt1.Y) == (pt2.Y < pt3.Y) {
			return distanceFromLineSqrd(pt2, pt1, pt3) < distSqrd
		}
		return distanceFromLineSqrd(pt3, pt1, pt2) < distSqrd
	}
}

func excludeOp(op *ClipperOutPt) *ClipperOutPt {
	result := op.Prev
	result.Next = op.Next
	op.Next.Prev = result
	result.Index = 0
	return result
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
