// ClipperOffset handles polygon offsetting (inflate/deflate) operations.
//
// Copyright (c) Angus Johnson 2010-2017, modified for PdfPig.
// Licensed under Boost Software License - Version 1.0.
// See: http://www.boost.org/LICENSE_1_0.txt
package clipperlibrary

import "math"

const defArcTolerance = 0.25

// ClipperOffset computes offset (inflated/deflated) versions of polygons.
type ClipperOffset struct {
	lowest       ClipperIntPoint
	polyNodes    ClipperPolyNode
	ArcTolerance float64
	MiterLimit   float64
}

// NewClipperOffset creates a new ClipperOffset with optional miter limit and arc tolerance.
func NewClipperOffset(opts ...func(*ClipperOffset)) *ClipperOffset {
	co := &ClipperOffset{
		lowest:       NewClipperIntPoint(-1, 0),
		ArcTolerance: defArcTolerance,
		MiterLimit:   2.0,
	}
	for _, opt := range opts {
		opt(co)
	}
	return co
}

// WithMiterLimit sets the miter limit for a ClipperOffset.
func WithMiterLimit(limit float64) func(*ClipperOffset) {
	return func(co *ClipperOffset) {
		co.MiterLimit = limit
	}
}

// WithArcTolerance sets the arc tolerance for a ClipperOffset.
func WithArcTolerance(tolerance float64) func(*ClipperOffset) {
	return func(co *ClipperOffset) {
		co.ArcTolerance = tolerance
	}
}

// Clear removes all paths and resets the offset state.
func (co *ClipperOffset) Clear() {
	co.polyNodes.Children = co.polyNodes.Children[:0]
	co.lowest = NewClipperIntPoint(-1, 0)
}

// AddPath adds a single path to the offset with the specified join and end types.
func (co *ClipperOffset) AddPath(path []ClipperIntPoint, joinType ClipperJoinType, endType ClipperEndType) {
	highI := len(path) - 1
	if highI < 0 {
		return
	}

	newNode := &ClipperPolyNode{
		JoinType: joinType,
		EndType:  endType,
		Polygon:  make([]ClipperIntPoint, 0, highI+1),
	}

	if endType == ClipperClosedLine || endType == ClipperClosedPolygon {
		for highI > 0 && path[0] == path[highI] {
			highI--
		}
	}

	newNode.Polygon = append(newNode.Polygon, path[0])
	j := 0
	k := 0
	for i := 1; i <= highI; i++ {
		if newNode.Polygon[j] != path[i] {
			j++
			newNode.Polygon = append(newNode.Polygon, path[i])
			if path[i].Y > newNode.Polygon[k].Y || (path[i].Y == newNode.Polygon[k].Y && path[i].X < newNode.Polygon[k].X) {
				k = j
			}
		}
	}

	if endType == ClipperClosedPolygon && j < 2 {
		return
	}

	co.polyNodes.AddChild(newNode)

	if endType != ClipperClosedPolygon {
		return
	}

	if co.lowest.X < 0 {
		co.lowest = NewClipperIntPoint(int64(co.polyNodes.ChildCount()-1), int64(k))
	} else {
		ip := co.polyNodes.Children[co.lowest.X].Polygon[co.lowest.Y]
		if newNode.Polygon[k].Y > ip.Y || (newNode.Polygon[k].Y == ip.Y && newNode.Polygon[k].X < ip.X) {
			co.lowest = NewClipperIntPoint(int64(co.polyNodes.ChildCount()-1), int64(k))
		}
	}
}

// AddPaths adds multiple paths to the offset with the specified join and end types.
func (co *ClipperOffset) AddPaths(paths [][]ClipperIntPoint, joinType ClipperJoinType, endType ClipperEndType) {
	for _, p := range paths {
		co.AddPath(p, joinType, endType)
	}
}

// GetUnitNormal returns the unit normal vector for the segment from pt1 to pt2.
func GetUnitNormal(pt1, pt2 ClipperIntPoint) ClipperDoublePoint {
	dx := float64(pt2.X - pt1.X)
	dy := float64(pt2.Y - pt1.Y)
	if dx == 0 && dy == 0 {
		return ClipperDoublePoint{}
	}
	f := 1.0 / math.Sqrt(dx*dx+dy*dy)
	dx *= f
	dy *= f
	return ClipperDoublePoint{X: dy, Y: -dx}
}
