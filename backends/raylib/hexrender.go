// Package raylib provides hex grid rendering for the raylib backend.
package raylib

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/spectrex/core"
)

// HexRenderer renders hex grids using raylib.
type HexRenderer struct {
	Config core.HexRenderConfig

	// Style overrides by coordinate
	cellStyles map[core.HexCoord]core.HexCellStyle
	edgeStyles map[core.HexEdge]core.HexEdgeStyle
}

// NewHexRenderer creates a new hex renderer with the given configuration.
func NewHexRenderer(config core.HexRenderConfig) *HexRenderer {
	return &HexRenderer{
		Config:     config,
		cellStyles: make(map[core.HexCoord]core.HexCellStyle),
		edgeStyles: make(map[core.HexEdge]core.HexEdgeStyle),
	}
}

// SetCellStyle sets a custom style for a specific cell.
func (r *HexRenderer) SetCellStyle(coord core.HexCoord, style core.HexCellStyle) {
	r.cellStyles[coord] = style
}

// ClearCellStyle removes the custom style for a cell, reverting to default.
func (r *HexRenderer) ClearCellStyle(coord core.HexCoord) {
	delete(r.cellStyles, coord)
}

// SetEdgeStyle sets a custom style for a specific edge.
func (r *HexRenderer) SetEdgeStyle(edge core.HexEdge, style core.HexEdgeStyle) {
	r.edgeStyles[edge] = style
}

// ClearEdgeStyle removes the custom style for an edge, reverting to default.
func (r *HexRenderer) ClearEdgeStyle(edge core.HexEdge) {
	delete(r.edgeStyles, edge)
}

// ClearAllStyles removes all custom styles.
func (r *HexRenderer) ClearAllStyles() {
	r.cellStyles = make(map[core.HexCoord]core.HexCellStyle)
	r.edgeStyles = make(map[core.HexEdge]core.HexEdgeStyle)
}

// getCellStyle returns the style for a cell, using override if set.
func (r *HexRenderer) getCellStyle(coord core.HexCoord) core.HexCellStyle {
	if style, ok := r.cellStyles[coord]; ok {
		return style
	}
	return r.Config.DefaultCell
}

// getEdgeStyle returns the style for an edge, using override if set.
func (r *HexRenderer) getEdgeStyle(edge core.HexEdge) core.HexEdgeStyle {
	if style, ok := r.edgeStyles[edge]; ok {
		return style
	}
	return r.Config.DefaultEdge
}

// DrawGrid renders the entire hex grid.
func (r *HexRenderer) DrawGrid(data core.HexGridRenderData) {
	// Draw cells first (so edges appear on top)
	if r.Config.DrawCells {
		for i, coord := range data.Cells {
			style := r.getCellStyle(coord)
			if style.FillColor.A > 0 {
				r.drawCellFill(data.Vertices[i], style.FillColor)
			}
		}
	}

	// Draw edges
	if r.Config.DrawEdges {
		// Draw all edges (interior and boundary)
		r.drawEdges(data.AllEdges, data)
	}
}

// DrawGridBoundaryOnly renders only the boundary edges of the grid.
func (r *HexRenderer) DrawGridBoundaryOnly(data core.HexGridRenderData) {
	r.drawEdges(data.BoundaryEdges, data)
}

// DrawCell renders a single hex cell at the given coordinate.
func (r *HexRenderer) DrawCell(coord core.HexCoord, style core.HexCellStyle) {
	vertices := core.HexVertices3D(r.Config.Layout, coord, r.Config.HexRadius)
	if style.FillColor.A > 0 {
		r.drawCellFill(vertices, style.FillColor)
	}
}

// DrawCellEdges renders all edges of a single cell.
func (r *HexRenderer) DrawCellEdges(coord core.HexCoord, style core.HexEdgeStyle) {
	vertices := core.HexVertices3D(r.Config.Layout, coord, r.Config.HexRadius)
	for dir := core.HexDirE; dir <= core.HexDirSE; dir++ {
		v1, v2 := core.HexEdgeVertices3D(vertices, dir)
		r.drawEdgeLine(v1, v2, style)
	}
}

// drawCellFill renders a filled hex using triangles.
func (r *HexRenderer) drawCellFill(vertices [6]core.Vec3, color core.Color) {
	rlColor := coreToRlColor(color)
	center := core.Vec3{
		X: (vertices[0].X + vertices[3].X) / 2,
		Y: (vertices[0].Y + vertices[3].Y) / 2,
		Z: (vertices[0].Z + vertices[3].Z) / 2,
	}

	// Draw 6 triangles from center to each edge
	for i := 0; i < 6; i++ {
		next := (i + 1) % 6
		rl.DrawTriangle3D(
			coreToRlVec3(center),
			coreToRlVec3(vertices[i]),
			coreToRlVec3(vertices[next]),
			rlColor,
		)
	}
}

// drawEdges renders a list of edges.
func (r *HexRenderer) drawEdges(edges []core.HexEdge, data core.HexGridRenderData) {
	// Build a coordinate to index map for fast lookup
	coordIndex := make(map[core.HexCoord]int, len(data.Cells))
	for i, coord := range data.Cells {
		coordIndex[coord] = i
	}

	for _, edge := range edges {
		style := r.getEdgeStyle(edge)

		// Find the vertices for this edge
		idx, ok := coordIndex[edge.Coord]
		if !ok {
			// Compute vertices on the fly if not in pre-computed data
			vertices := core.HexVertices3D(r.Config.Layout, edge.Coord, r.Config.HexRadius)
			v1, v2 := core.HexEdgeVertices3D(vertices, edge.Dir)
			r.drawEdgeLine(v1, v2, style)
		} else {
			v1, v2 := core.HexEdgeVertices3D(data.Vertices[idx], edge.Dir)
			r.drawEdgeLine(v1, v2, style)
		}
	}
}

// drawEdgeLine renders a single edge line with the given style.
func (r *HexRenderer) drawEdgeLine(v1, v2 core.Vec3, style core.HexEdgeStyle) {
	rlColor := coreToRlColor(style.Color)

	if style.Dashed {
		r.drawDashedLine3D(v1, v2, r.Config.DashLength, r.Config.DashGap, rlColor)
	} else {
		rl.DrawLine3D(coreToRlVec3(v1), coreToRlVec3(v2), rlColor)
	}
}

// drawDashedLine3D draws a dashed line between two points.
func (r *HexRenderer) drawDashedLine3D(start, end core.Vec3, dashLen, gapLen float32, color rl.Color) {
	dx := end.X - start.X
	dy := end.Y - start.Y
	dz := end.Z - start.Z
	totalLen := float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))

	if totalLen == 0 {
		return
	}

	// Normalize direction
	dx /= totalLen
	dy /= totalLen
	dz /= totalLen

	segmentLen := dashLen + gapLen
	pos := float32(0)

	for pos < totalLen {
		// Start of this dash
		p1 := core.Vec3{
			X: start.X + dx*pos,
			Y: start.Y + dy*pos,
			Z: start.Z + dz*pos,
		}

		// End of this dash (clamped to line end)
		dashEnd := pos + dashLen
		if dashEnd > totalLen {
			dashEnd = totalLen
		}
		p2 := core.Vec3{
			X: start.X + dx*dashEnd,
			Y: start.Y + dy*dashEnd,
			Z: start.Z + dz*dashEnd,
		}

		rl.DrawLine3D(coreToRlVec3(p1), coreToRlVec3(p2), color)

		pos += segmentLen
	}
}

// DrawGridWithCallback renders the grid, calling the callback for each cell
// to get customized styles. This is useful for dynamic styling.
func (r *HexRenderer) DrawGridWithCallback(
	data core.HexGridRenderData,
	cellStyleFn func(coord core.HexCoord) *core.HexCellStyle,
	edgeStyleFn func(edge core.HexEdge) *core.HexEdgeStyle,
) {
	// Draw cells
	if r.Config.DrawCells && cellStyleFn != nil {
		for i, coord := range data.Cells {
			if style := cellStyleFn(coord); style != nil && style.FillColor.A > 0 {
				r.drawCellFill(data.Vertices[i], style.FillColor)
			}
		}
	}

	// Draw edges
	if r.Config.DrawEdges && edgeStyleFn != nil {
		for _, edge := range data.AllEdges {
			if style := edgeStyleFn(edge); style != nil {
				idx := -1
				for i, coord := range data.Cells {
					if coord == edge.Coord {
						idx = i
						break
					}
				}
				if idx >= 0 {
					v1, v2 := core.HexEdgeVertices3D(data.Vertices[idx], edge.Dir)
					r.drawEdgeLine(v1, v2, *style)
				}
			}
		}
	}
}
