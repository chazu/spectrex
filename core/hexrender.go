// Package core provides hex grid rendering utilities for the Spectrex framework.
package core

import "math"

// HexCellStyle defines the visual style for a hex cell.
type HexCellStyle struct {
	FillColor Color // Fill color for the cell (use alpha 0 for transparent)
}

// HexEdgeStyle defines the visual style for hex edges.
type HexEdgeStyle struct {
	Color  Color // Edge color
	Dashed bool  // If true, render as dashed line
}

// HexEdge represents an edge between two hex cells.
// An edge is uniquely identified by the hex coordinate and the direction
// of the edge (E, NE, NW only - the other 3 are the same edges from neighbors).
type HexEdge struct {
	Coord HexCoord     // The hex this edge belongs to
	Dir   HexDirection // Direction of the edge (E=0, NE=1, NW=2 only)
}

// HexRenderConfig configures how a hex grid is rendered.
type HexRenderConfig struct {
	Layout       HexLayout    // Layout for hex-to-pixel conversion
	HexRadius    float32      // Radius of each hex (from center to vertex)
	DefaultCell  HexCellStyle // Default cell style
	DefaultEdge  HexEdgeStyle // Default edge style
	DrawCells    bool         // Whether to draw cell fills
	DrawEdges    bool         // Whether to draw edges
	DashLength   float32      // Length of dash segments for dashed edges
	DashGap      float32      // Gap between dashes
}

// DefaultHexRenderConfig returns a default hex render configuration.
func DefaultHexRenderConfig(hexRadius float32) HexRenderConfig {
	return HexRenderConfig{
		Layout:    NewHexLayout(Vec2{X: hexRadius, Y: hexRadius}, Vec2{X: 0, Y: 0}),
		HexRadius: hexRadius,
		DefaultCell: HexCellStyle{
			FillColor: Color{R: 50, G: 50, B: 80, A: 200},
		},
		DefaultEdge: HexEdgeStyle{
			Color:  ColorWhite,
			Dashed: false,
		},
		DrawCells:  true,
		DrawEdges:  true,
		DashLength: 5.0,
		DashGap:    3.0,
	}
}

// HexVertices returns the 6 vertices of a hex at the given coordinate.
// Uses pointy-top orientation (first vertex at top).
// Returns vertices in counter-clockwise order starting from top.
func HexVertices(layout HexLayout, coord HexCoord, radius float32) [6]Vec2 {
	center := layout.ToPixel(coord)
	var vertices [6]Vec2

	// Pointy-top: vertices at angles 90°, 30°, -30°, -90°, -150°, 150° (or 90, 30, 330, 270, 210, 150)
	// Start at top (90°) and go clockwise
	for i := 0; i < 6; i++ {
		angle := math.Pi/2 - float64(i)*math.Pi/3 // 90° - i*60°
		vertices[i] = Vec2{
			X: center.X + radius*float32(math.Cos(angle)),
			Y: center.Y - radius*float32(math.Sin(angle)), // Negate Y for screen coords
		}
	}

	return vertices
}

// HexVertices3D returns the 6 vertices of a hex in 3D space (on the XZ plane at Y=0).
func HexVertices3D(layout HexLayout, coord HexCoord, radius float32) [6]Vec3 {
	v2 := HexVertices(layout, coord, radius)
	var vertices [6]Vec3
	for i := 0; i < 6; i++ {
		vertices[i] = Vec3{X: v2[i].X, Y: 0, Z: v2[i].Y}
	}
	return vertices
}

// HexEdgeVertices returns the two vertices that form the edge in the given direction.
// For pointy-top hexes:
// - E edge: vertices 1 and 2 (right side)
// - NE edge: vertices 0 and 1 (upper right)
// - NW edge: vertices 5 and 0 (upper left)
// - W edge: vertices 4 and 5 (left side)
// - SW edge: vertices 3 and 4 (lower left)
// - SE edge: vertices 2 and 3 (lower right)
func HexEdgeVertices(vertices [6]Vec2, dir HexDirection) (Vec2, Vec2) {
	// Vertex indices for each edge direction
	edgeVertexMap := [6][2]int{
		{1, 2}, // E
		{0, 1}, // NE
		{5, 0}, // NW
		{4, 5}, // W
		{3, 4}, // SW
		{2, 3}, // SE
	}
	idx := edgeVertexMap[dir]
	return vertices[idx[0]], vertices[idx[1]]
}

// HexEdgeVertices3D returns the two vertices of an edge in 3D space.
func HexEdgeVertices3D(vertices [6]Vec3, dir HexDirection) (Vec3, Vec3) {
	edgeVertexMap := [6][2]int{
		{1, 2}, // E
		{0, 1}, // NE
		{5, 0}, // NW
		{4, 5}, // W
		{3, 4}, // SW
		{2, 3}, // SE
	}
	idx := edgeVertexMap[dir]
	return vertices[idx[0]], vertices[idx[1]]
}

// GridEdges returns all unique edges for a hex grid.
// Only includes edges where at least one endpoint is within the grid.
// To avoid duplicates, only returns edges with direction E, NE, or NW (0, 1, 2).
func GridEdges[T any](grid *HexGrid[T]) []HexEdge {
	var edges []HexEdge

	grid.ForEach(func(coord HexCoord, _ T) {
		// Only add edges in directions 0, 1, 2 (E, NE, NW) to avoid duplicates
		for dir := HexDirE; dir <= HexDirNW; dir++ {
			edges = append(edges, HexEdge{Coord: coord, Dir: dir})
		}
	})

	return edges
}

// InteriorEdges returns edges that are shared between two cells in the grid.
// These are edges where both adjacent hexes are valid grid positions.
func InteriorEdges[T any](grid *HexGrid[T]) []HexEdge {
	var edges []HexEdge

	grid.ForEach(func(coord HexCoord, _ T) {
		// Only check directions 0, 1, 2 to avoid duplicates
		for dir := HexDirE; dir <= HexDirNW; dir++ {
			neighbor := coord.Neighbor(dir)
			if grid.IsValid(neighbor) {
				edges = append(edges, HexEdge{Coord: coord, Dir: dir})
			}
		}
	})

	return edges
}

// BoundaryEdges returns edges that are on the boundary of the grid.
// These are edges where one hex is valid and the neighbor is not.
func BoundaryEdges[T any](grid *HexGrid[T]) []HexEdge {
	var edges []HexEdge

	grid.ForEach(func(coord HexCoord, _ T) {
		// Check all 6 directions for boundary edges
		for dir := HexDirE; dir <= HexDirSE; dir++ {
			neighbor := coord.Neighbor(dir)
			if !grid.IsValid(neighbor) {
				edges = append(edges, HexEdge{Coord: coord, Dir: dir})
			}
		}
	})

	return edges
}

// HexGridRenderData holds pre-computed rendering data for a hex grid.
type HexGridRenderData struct {
	Cells       []HexCoord  // All cell coordinates
	Vertices    [][6]Vec3   // Vertices for each cell (same index as Cells)
	AllEdges    []HexEdge   // All unique edges
	BoundaryEdges []HexEdge // Edges on the grid boundary
	InteriorEdges []HexEdge // Edges between cells
}

// PrepareGridRenderData computes all the rendering data for a hex grid.
func PrepareGridRenderData[T any](grid *HexGrid[T], config HexRenderConfig) HexGridRenderData {
	cells := grid.All()
	vertices := make([][6]Vec3, len(cells))

	for i, coord := range cells {
		vertices[i] = HexVertices3D(config.Layout, coord, config.HexRadius)
	}

	return HexGridRenderData{
		Cells:         cells,
		Vertices:      vertices,
		AllEdges:      GridEdges(grid),
		BoundaryEdges: BoundaryEdges(grid),
		InteriorEdges: InteriorEdges(grid),
	}
}
