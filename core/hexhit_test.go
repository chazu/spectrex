package core

import (
	"math"
	"testing"
)

func TestPointToSegmentDistance(t *testing.T) {
	tests := []struct {
		name         string
		px, py       float32
		x1, y1       float32
		x2, y2       float32
		wantDist     float32
		wantTApprox  float32 // approximate t value
		tolerance    float32
	}{
		{
			name:        "point on segment start",
			px: 0, py: 0,
			x1: 0, y1: 0, x2: 10, y2: 0,
			wantDist: 0, wantTApprox: 0, tolerance: 0.001,
		},
		{
			name:        "point on segment end",
			px: 10, py: 0,
			x1: 0, y1: 0, x2: 10, y2: 0,
			wantDist: 0, wantTApprox: 1, tolerance: 0.001,
		},
		{
			name:        "point on segment middle",
			px: 5, py: 0,
			x1: 0, y1: 0, x2: 10, y2: 0,
			wantDist: 0, wantTApprox: 0.5, tolerance: 0.001,
		},
		{
			name:        "point perpendicular to middle",
			px: 5, py: 3,
			x1: 0, y1: 0, x2: 10, y2: 0,
			wantDist: 3, wantTApprox: 0.5, tolerance: 0.001,
		},
		{
			name:        "point before segment start",
			px: -3, py: 0,
			x1: 0, y1: 0, x2: 10, y2: 0,
			wantDist: 3, wantTApprox: 0, tolerance: 0.001,
		},
		{
			name:        "point after segment end",
			px: 13, py: 0,
			x1: 0, y1: 0, x2: 10, y2: 0,
			wantDist: 3, wantTApprox: 1, tolerance: 0.001,
		},
		{
			name:        "diagonal segment",
			px: 5, py: 5,
			x1: 0, y1: 0, x2: 10, y2: 10,
			wantDist: 0, wantTApprox: 0.5, tolerance: 0.001,
		},
		{
			name:        "point perpendicular to diagonal",
			px: 0, py: 10,
			x1: 0, y1: 0, x2: 10, y2: 10,
			wantDist: float32(10 / math.Sqrt(2)), wantTApprox: 0.5, tolerance: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dist, tVal := PointToSegmentDistance(tt.px, tt.py, tt.x1, tt.y1, tt.x2, tt.y2)

			if math.Abs(float64(dist-tt.wantDist)) > float64(tt.tolerance) {
				t.Errorf("distance = %v, want %v (tolerance %v)", dist, tt.wantDist, tt.tolerance)
			}
			if math.Abs(float64(tVal-tt.wantTApprox)) > float64(tt.tolerance) {
				t.Errorf("t = %v, want ~%v (tolerance %v)", tVal, tt.wantTApprox, tt.tolerance)
			}
		})
	}
}

func TestPointToSegmentDistance_ZeroLengthSegment(t *testing.T) {
	// Degenerate case: segment is a single point
	dist, tVal := PointToSegmentDistance(3, 4, 0, 0, 0, 0)

	if math.Abs(float64(dist-5)) > 0.001 {
		t.Errorf("distance to point = %v, want 5", dist)
	}
	if tVal != 0 {
		t.Errorf("t for point = %v, want 0", tVal)
	}
}

func TestHexHitTester_HitTestCell(t *testing.T) {
	// Create a hit tester with origin at (100, 100) and hex radius of 20
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 100, Y: 100})
	tester := NewHexHitTester(layout, 20, 5)

	tests := []struct {
		name   string
		px, py float32
		wantQ  int
		wantR  int
	}{
		{
			name:  "origin",
			px:    100, py: 100,
			wantQ: 0, wantR: 0,
		},
		{
			name:  "offset from origin - still in center hex",
			px:    105, py: 105,
			wantQ: 0, wantR: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cell := tester.HitTestCell(tt.px, tt.py)
			if cell.Q != tt.wantQ || cell.R != tt.wantR {
				t.Errorf("HitTestCell(%v, %v) = (%v, %v), want (%v, %v)",
					tt.px, tt.py, cell.Q, cell.R, tt.wantQ, tt.wantR)
			}
		})
	}
}

func TestHexHitTester_HitTest(t *testing.T) {
	// Create a hit tester with hex radius of 20 and edge threshold of 5
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 100, Y: 100})
	tester := NewHexHitTester(layout, 20, 5)

	// Test cell center - should return HexHitCell
	result := tester.HitTest(100, 100)
	if result.Type != HexHitCell {
		t.Errorf("hit at center: Type = %v, want HexHitCell", result.Type)
	}
	if result.Cell.Q != 0 || result.Cell.R != 0 {
		t.Errorf("hit at center: Cell = (%v, %v), want (0, 0)", result.Cell.Q, result.Cell.R)
	}
}

func TestHexHitTester_HitTestEdge(t *testing.T) {
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 100, Y: 100})
	tester := NewHexHitTester(layout, 20, 5)

	// Get an edge and test near it
	edge, dist := tester.HitTestEdge(100, 100)

	// At center, we should get some edge (the nearest one)
	if dist < 0 {
		t.Errorf("HitTestEdge returned negative distance: %v", dist)
	}

	// Edge should have valid coordinates
	if edge.Dir > HexDirNW {
		t.Errorf("HitTestEdge returned non-canonical direction: %v", edge.Dir)
	}
}

func TestNormalizeEdge(t *testing.T) {
	tests := []struct {
		name     string
		coord    HexCoord
		dir      HexDirection
		wantDir  HexDirection // Should always be E, NE, or NW
	}{
		{
			name:    "E direction - already canonical",
			coord:   HexCoord{Q: 0, R: 0},
			dir:     HexDirE,
			wantDir: HexDirE,
		},
		{
			name:    "NE direction - already canonical",
			coord:   HexCoord{Q: 0, R: 0},
			dir:     HexDirNE,
			wantDir: HexDirNE,
		},
		{
			name:    "NW direction - already canonical",
			coord:   HexCoord{Q: 0, R: 0},
			dir:     HexDirNW,
			wantDir: HexDirNW,
		},
		{
			name:    "W direction - converts to neighbor's E",
			coord:   HexCoord{Q: 0, R: 0},
			dir:     HexDirW,
			wantDir: HexDirE,
		},
		{
			name:    "SW direction - converts to neighbor's NE",
			coord:   HexCoord{Q: 0, R: 0},
			dir:     HexDirSW,
			wantDir: HexDirNE,
		},
		{
			name:    "SE direction - converts to neighbor's NW",
			coord:   HexCoord{Q: 0, R: 0},
			dir:     HexDirSE,
			wantDir: HexDirNW,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edge := normalizeEdge(tt.coord, tt.dir)
			if edge.Dir != tt.wantDir {
				t.Errorf("normalizeEdge(%v, %v).Dir = %v, want %v",
					tt.coord, tt.dir, edge.Dir, tt.wantDir)
			}
		})
	}
}

func TestNormalizeEdge_SameEdge(t *testing.T) {
	// Two cells sharing an edge should normalize to the same HexEdge
	coord1 := HexCoord{Q: 0, R: 0}
	coord2 := HexCoord{Q: 1, R: 0} // East neighbor

	edge1 := normalizeEdge(coord1, HexDirE)
	edge2 := normalizeEdge(coord2, HexDirW)

	if edge1.Coord != edge2.Coord || edge1.Dir != edge2.Dir {
		t.Errorf("same edge from different cells: %v != %v", edge1, edge2)
	}
}

func TestHexHitTester_HitTestInGrid(t *testing.T) {
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 100, Y: 100})
	tester := NewHexHitTester(layout, 20, 5)

	// Test inside grid (radius 2)
	result := tester.HitTestInGrid(100, 100, 2)
	if result.Type == HexHitNone {
		t.Error("HitTestInGrid at origin with radius 2 returned HexHitNone")
	}

	// Test outside grid - point far from origin
	result = tester.HitTestInGrid(500, 500, 2)
	if result.Type != HexHitNone {
		t.Errorf("HitTestInGrid far from origin: Type = %v, want HexHitNone", result.Type)
	}
}

func TestHexHitTester_EdgeVertices(t *testing.T) {
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 0, Y: 0})
	tester := NewHexHitTester(layout, 20, 5)

	edge := HexEdge{Coord: HexCoord{Q: 0, R: 0}, Dir: HexDirE}
	v1, v2 := tester.EdgeVertices(edge)

	// Vertices should be different
	if v1.X == v2.X && v1.Y == v2.Y {
		t.Error("EdgeVertices returned identical points")
	}

	// Distance between vertices should be approximately hex edge length
	dx := v2.X - v1.X
	dy := v2.Y - v1.Y
	edgeLen := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// For a hex with radius 20, edge length should be 20 (each side = radius)
	if math.Abs(float64(edgeLen-20)) > 0.1 {
		t.Errorf("edge length = %v, want ~20", edgeLen)
	}
}

func TestHexHitTester_EdgeMidpoint(t *testing.T) {
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 0, Y: 0})
	tester := NewHexHitTester(layout, 20, 5)

	edge := HexEdge{Coord: HexCoord{Q: 0, R: 0}, Dir: HexDirE}
	mid := tester.EdgeMidpoint(edge)
	v1, v2 := tester.EdgeVertices(edge)

	// Midpoint should be average of vertices
	expectedX := (v1.X + v2.X) / 2
	expectedY := (v1.Y + v2.Y) / 2

	if math.Abs(float64(mid.X-expectedX)) > 0.001 || math.Abs(float64(mid.Y-expectedY)) > 0.001 {
		t.Errorf("EdgeMidpoint = (%v, %v), want (%v, %v)", mid.X, mid.Y, expectedX, expectedY)
	}
}

func TestHexHitTester_CellCenter(t *testing.T) {
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 100, Y: 100})
	tester := NewHexHitTester(layout, 20, 5)

	center := tester.CellCenter(HexCoord{Q: 0, R: 0})

	// Origin cell center should be at layout origin
	if center.X != 100 || center.Y != 100 {
		t.Errorf("CellCenter(0,0) = (%v, %v), want (100, 100)", center.X, center.Y)
	}
}

func TestHexHitType_Coverage(t *testing.T) {
	// Ensure all hit types are distinct
	types := []HexHitType{HexHitNone, HexHitCell, HexHitEdge}
	seen := make(map[HexHitType]bool)

	for _, typ := range types {
		if seen[typ] {
			t.Errorf("duplicate HexHitType value: %v", typ)
		}
		seen[typ] = true
	}
}
