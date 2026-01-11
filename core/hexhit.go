// Package core provides hex input hit testing for the Spectrex framework.
package core

import "math"

// HexHitType indicates what type of element was hit.
type HexHitType int

const (
	HexHitNone HexHitType = iota // No hit (outside grid or no valid target)
	HexHitCell                   // Hit a cell (interior)
	HexHitEdge                   // Hit near an edge
)

// HexHitResult contains the result of a hit test.
type HexHitResult struct {
	Type     HexHitType // What was hit
	Cell     HexCoord   // The cell coordinate (valid for both Cell and Edge hits)
	Edge     HexEdge    // The edge (only valid when Type == HexHitEdge)
	Distance float32    // Distance from hit point to the element
}

// PointToSegmentDistance calculates the minimum distance from a point to a line segment.
// Returns the distance and the parameter t (0-1) indicating where on the segment the closest point is.
func PointToSegmentDistance(px, py, x1, y1, x2, y2 float32) (distance float32, t float32) {
	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 && dy == 0 {
		// Segment is a point
		return float32(math.Sqrt(float64((px-x1)*(px-x1) + (py-y1)*(py-y1)))), 0
	}

	// Calculate parameter t for the projection of point onto the line
	// t = ((P - A) . (B - A)) / |B - A|^2
	lengthSq := dx*dx + dy*dy
	t = ((px-x1)*dx + (py-y1)*dy) / lengthSq

	// Clamp t to [0, 1] to stay within segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// Calculate closest point on segment
	closestX := x1 + t*dx
	closestY := y1 + t*dy

	// Return distance
	distX := px - closestX
	distY := py - closestY
	distance = float32(math.Sqrt(float64(distX*distX + distY*distY)))
	return distance, t
}

// HexHitTester performs hit testing on hex grids.
type HexHitTester struct {
	Layout       HexLayout
	HexRadius    float32
	EdgeThreshold float32 // Maximum distance to consider an edge "hit"
}

// NewHexHitTester creates a new hit tester with the given parameters.
func NewHexHitTester(layout HexLayout, hexRadius, edgeThreshold float32) *HexHitTester {
	return &HexHitTester{
		Layout:       layout,
		HexRadius:    hexRadius,
		EdgeThreshold: edgeThreshold,
	}
}

// HitTest performs a complete hit test at the given pixel coordinates.
// It returns the cell containing the point and optionally the nearest edge
// if within the edge threshold.
func (h *HexHitTester) HitTest(px, py float32) HexHitResult {
	// First, find which cell contains this point
	cell := h.Layout.FromPixel(Vec2{X: px, Y: py})

	// Get the vertices for this cell
	vertices := HexVertices(h.Layout, cell, h.HexRadius)

	// Find the nearest edge
	minDist := float32(math.MaxFloat32)
	var nearestEdge HexEdge
	var nearestDir HexDirection

	for dir := HexDirE; dir <= HexDirSE; dir++ {
		v1, v2 := HexEdgeVertices(vertices, dir)
		dist, _ := PointToSegmentDistance(px, py, v1.X, v1.Y, v2.X, v2.Y)
		if dist < minDist {
			minDist = dist
			nearestDir = dir
		}
	}

	// Normalize edge direction to canonical form (E, NE, NW only)
	// to match the edge representation used elsewhere
	nearestEdge = normalizeEdge(cell, nearestDir)

	// Determine hit type based on distance
	if minDist <= h.EdgeThreshold {
		return HexHitResult{
			Type:     HexHitEdge,
			Cell:     cell,
			Edge:     nearestEdge,
			Distance: minDist,
		}
	}

	return HexHitResult{
		Type:     HexHitCell,
		Cell:     cell,
		Edge:     nearestEdge, // Still provide nearest edge info
		Distance: minDist,
	}
}

// HitTestCell returns only the cell at the given pixel coordinates.
// This is a lightweight version of HitTest when you only need the cell.
func (h *HexHitTester) HitTestCell(px, py float32) HexCoord {
	return h.Layout.FromPixel(Vec2{X: px, Y: py})
}

// HitTestEdge finds the nearest edge to the given pixel coordinates.
// Returns the edge and its distance from the point.
func (h *HexHitTester) HitTestEdge(px, py float32) (HexEdge, float32) {
	cell := h.Layout.FromPixel(Vec2{X: px, Y: py})
	vertices := HexVertices(h.Layout, cell, h.HexRadius)

	minDist := float32(math.MaxFloat32)
	var nearestDir HexDirection

	for dir := HexDirE; dir <= HexDirSE; dir++ {
		v1, v2 := HexEdgeVertices(vertices, dir)
		dist, _ := PointToSegmentDistance(px, py, v1.X, v1.Y, v2.X, v2.Y)
		if dist < minDist {
			minDist = dist
			nearestDir = dir
		}
	}

	return normalizeEdge(cell, nearestDir), minDist
}

// normalizeEdge converts an edge to its canonical form.
// Each edge is shared by two hexes. We canonicalize by always using
// the edge from the hex where the direction is E, NE, or NW.
func normalizeEdge(coord HexCoord, dir HexDirection) HexEdge {
	if dir <= HexDirNW {
		// Already canonical (E=0, NE=1, NW=2)
		return HexEdge{Coord: coord, Dir: dir}
	}

	// Convert to the neighbor's perspective
	neighbor := coord.Neighbor(dir)
	oppositeDir := (dir + 3) % 6 // Opposite direction: W->E, SW->NE, SE->NW

	return HexEdge{Coord: neighbor, Dir: oppositeDir}
}

// HitTestInGrid performs a hit test and checks if the result is within the grid.
func (h *HexHitTester) HitTestInGrid(px, py float32, gridRadius int) HexHitResult {
	result := h.HitTest(px, py)

	// Check if the cell is within the grid
	if result.Cell.Length() > gridRadius {
		return HexHitResult{Type: HexHitNone}
	}

	return result
}

// EdgeVertices returns the pixel coordinates of an edge's endpoints.
func (h *HexHitTester) EdgeVertices(edge HexEdge) (Vec2, Vec2) {
	vertices := HexVertices(h.Layout, edge.Coord, h.HexRadius)
	return HexEdgeVertices(vertices, edge.Dir)
}

// EdgeMidpoint returns the pixel coordinates of an edge's midpoint.
func (h *HexHitTester) EdgeMidpoint(edge HexEdge) Vec2 {
	v1, v2 := h.EdgeVertices(edge)
	return Vec2{
		X: (v1.X + v2.X) / 2,
		Y: (v1.Y + v2.Y) / 2,
	}
}

// CellCenter returns the pixel coordinates of a cell's center.
func (h *HexHitTester) CellCenter(coord HexCoord) Vec2 {
	return h.Layout.ToPixel(coord)
}
