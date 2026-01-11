package core

import (
	"testing"
)

func TestNewHexCoord(t *testing.T) {
	h := NewHexCoord(3, -2)
	if h.Q != 3 || h.R != -2 {
		t.Errorf("NewHexCoord(3, -2) = %v, want {Q:3, R:-2}", h)
	}
}

func TestHexCoord_S(t *testing.T) {
	tests := []struct {
		h    HexCoord
		want int
	}{
		{HexCoord{Q: 0, R: 0}, 0},
		{HexCoord{Q: 1, R: -1}, 0},
		{HexCoord{Q: 3, R: -2}, -1},
		{HexCoord{Q: -1, R: 2}, -1},
	}

	for _, tt := range tests {
		got := tt.h.S()
		if got != tt.want {
			t.Errorf("%v.S() = %d, want %d", tt.h, got, tt.want)
		}
	}
}

func TestHexCoord_Add(t *testing.T) {
	a := HexCoord{Q: 1, R: 2}
	b := HexCoord{Q: 3, R: -1}
	got := a.Add(b)
	want := HexCoord{Q: 4, R: 1}
	if got != want {
		t.Errorf("%v.Add(%v) = %v, want %v", a, b, got, want)
	}
}

func TestHexCoord_Sub(t *testing.T) {
	a := HexCoord{Q: 3, R: 2}
	b := HexCoord{Q: 1, R: 4}
	got := a.Sub(b)
	want := HexCoord{Q: 2, R: -2}
	if got != want {
		t.Errorf("%v.Sub(%v) = %v, want %v", a, b, got, want)
	}
}

func TestHexCoord_Scale(t *testing.T) {
	h := HexCoord{Q: 2, R: -3}
	got := h.Scale(3)
	want := HexCoord{Q: 6, R: -9}
	if got != want {
		t.Errorf("%v.Scale(3) = %v, want %v", h, got, want)
	}
}

func TestHexCoord_Neighbor(t *testing.T) {
	origin := HexCoord{Q: 0, R: 0}

	tests := []struct {
		dir  HexDirection
		want HexCoord
	}{
		{HexDirE, HexCoord{Q: 1, R: 0}},
		{HexDirNE, HexCoord{Q: 1, R: -1}},
		{HexDirNW, HexCoord{Q: 0, R: -1}},
		{HexDirW, HexCoord{Q: -1, R: 0}},
		{HexDirSW, HexCoord{Q: -1, R: 1}},
		{HexDirSE, HexCoord{Q: 0, R: 1}},
	}

	for _, tt := range tests {
		got := origin.Neighbor(tt.dir)
		if got != tt.want {
			t.Errorf("origin.Neighbor(%d) = %v, want %v", tt.dir, got, tt.want)
		}
	}
}

func TestHexCoord_Neighbors(t *testing.T) {
	h := HexCoord{Q: 1, R: 1}
	neighbors := h.Neighbors()

	if len(neighbors) != 6 {
		t.Errorf("Neighbors() returned %d hexes, want 6", len(neighbors))
	}

	// Check that all neighbors are at distance 1
	for i, n := range neighbors {
		dist := h.Distance(n)
		if dist != 1 {
			t.Errorf("Neighbor[%d] %v is at distance %d, want 1", i, n, dist)
		}
	}
}

func TestHexCoord_Distance(t *testing.T) {
	tests := []struct {
		a, b HexCoord
		want int
	}{
		{HexCoord{0, 0}, HexCoord{0, 0}, 0},
		{HexCoord{0, 0}, HexCoord{1, 0}, 1},
		{HexCoord{0, 0}, HexCoord{0, 1}, 1},
		{HexCoord{0, 0}, HexCoord{1, -1}, 1},
		{HexCoord{0, 0}, HexCoord{2, 0}, 2},
		{HexCoord{0, 0}, HexCoord{2, -1}, 2},
		{HexCoord{0, 0}, HexCoord{3, -3}, 3},
		{HexCoord{1, 2}, HexCoord{4, -1}, 3},
	}

	for _, tt := range tests {
		got := tt.a.Distance(tt.b)
		if got != tt.want {
			t.Errorf("%v.Distance(%v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
		// Distance should be symmetric
		gotReverse := tt.b.Distance(tt.a)
		if gotReverse != tt.want {
			t.Errorf("%v.Distance(%v) = %d, want %d (reverse)", tt.b, tt.a, gotReverse, tt.want)
		}
	}
}

func TestHexCoord_Length(t *testing.T) {
	tests := []struct {
		h    HexCoord
		want int
	}{
		{HexCoord{0, 0}, 0},
		{HexCoord{1, 0}, 1},
		{HexCoord{0, -1}, 1},
		{HexCoord{2, -1}, 2},
		{HexCoord{3, -3}, 3},
	}

	for _, tt := range tests {
		got := tt.h.Length()
		if got != tt.want {
			t.Errorf("%v.Length() = %d, want %d", tt.h, got, tt.want)
		}
	}
}

func TestHexCoord_Equal(t *testing.T) {
	a := HexCoord{Q: 1, R: 2}
	b := HexCoord{Q: 1, R: 2}
	c := HexCoord{Q: 1, R: 3}

	if !a.Equal(b) {
		t.Errorf("%v.Equal(%v) = false, want true", a, b)
	}
	if a.Equal(c) {
		t.Errorf("%v.Equal(%v) = true, want false", a, c)
	}
}

func TestHexCoord_ToCube(t *testing.T) {
	h := HexCoord{Q: 2, R: -1}
	cube := h.ToCube()

	if cube.Q != 2 || cube.R != -1 || cube.S != -1 {
		t.Errorf("%v.ToCube() = %v, want {Q:2, R:-1, S:-1}", h, cube)
	}

	// Cube constraint: Q + R + S = 0
	if cube.Q+cube.R+cube.S != 0 {
		t.Errorf("Cube constraint violated: %d + %d + %d != 0", cube.Q, cube.R, cube.S)
	}
}

func TestHexCubeCoord_ToAxial(t *testing.T) {
	cube := HexCubeCoord{Q: 2, R: -1, S: -1}
	axial := cube.ToAxial()

	if axial.Q != 2 || axial.R != -1 {
		t.Errorf("%v.ToAxial() = %v, want {Q:2, R:-1}", cube, axial)
	}
}

func TestHexLayout_ToPixel(t *testing.T) {
	layout := NewHexLayout(Vec2{X: 10, Y: 10}, Vec2{X: 0, Y: 0})

	// Origin hex should be at origin pixel (plus any offset)
	origin := HexCoord{Q: 0, R: 0}
	p := layout.ToPixel(origin)
	if p.X != 0 || p.Y != 0 {
		t.Errorf("ToPixel(%v) = %v, want {0, 0}", origin, p)
	}
}

func TestHexLayout_FromPixel_Roundtrip(t *testing.T) {
	layout := NewHexLayout(Vec2{X: 20, Y: 20}, Vec2{X: 100, Y: 100})

	testHexes := []HexCoord{
		{0, 0},
		{1, 0},
		{0, 1},
		{-1, 1},
		{2, -1},
		{-3, 2},
	}

	for _, h := range testHexes {
		pixel := layout.ToPixel(h)
		recovered := layout.FromPixel(pixel)
		if !h.Equal(recovered) {
			t.Errorf("Roundtrip failed: %v -> pixel %v -> %v", h, pixel, recovered)
		}
	}
}

func TestHexRing(t *testing.T) {
	center := HexCoord{Q: 0, R: 0}

	// Ring at radius 0 should just be center
	ring0 := HexRing(center, 0)
	if len(ring0) != 1 || !ring0[0].Equal(center) {
		t.Errorf("HexRing(center, 0) = %v, want [center]", ring0)
	}

	// Ring at radius 1 should have 6 hexes
	ring1 := HexRing(center, 1)
	if len(ring1) != 6 {
		t.Errorf("HexRing(center, 1) has %d hexes, want 6", len(ring1))
	}

	// All hexes in ring 1 should be at distance 1 from center
	for _, h := range ring1 {
		dist := center.Distance(h)
		if dist != 1 {
			t.Errorf("HexRing hex %v is at distance %d, want 1", h, dist)
		}
	}

	// Ring at radius 2 should have 12 hexes
	ring2 := HexRing(center, 2)
	if len(ring2) != 12 {
		t.Errorf("HexRing(center, 2) has %d hexes, want 12", len(ring2))
	}

	// Ring at radius 3 should have 18 hexes
	ring3 := HexRing(center, 3)
	if len(ring3) != 18 {
		t.Errorf("HexRing(center, 3) has %d hexes, want 18", len(ring3))
	}
}

func TestHexSpiral(t *testing.T) {
	center := HexCoord{Q: 0, R: 0}

	// Spiral with radius 0 should just be center
	spiral0 := HexSpiral(center, 0)
	if len(spiral0) != 1 || !spiral0[0].Equal(center) {
		t.Errorf("HexSpiral(center, 0) = %v, want [center]", spiral0)
	}

	// Spiral with radius 1: 1 (center) + 6 (ring1) = 7
	spiral1 := HexSpiral(center, 1)
	if len(spiral1) != 7 {
		t.Errorf("HexSpiral(center, 1) has %d hexes, want 7", len(spiral1))
	}

	// Spiral with radius 2: 1 + 6 + 12 = 19
	spiral2 := HexSpiral(center, 2)
	if len(spiral2) != 19 {
		t.Errorf("HexSpiral(center, 2) has %d hexes, want 19", len(spiral2))
	}
}

func TestHexLine(t *testing.T) {
	// Line from origin to self
	origin := HexCoord{Q: 0, R: 0}
	line0 := HexLine(origin, origin)
	if len(line0) != 1 {
		t.Errorf("HexLine to self has %d hexes, want 1", len(line0))
	}

	// Line to adjacent hex
	adj := HexCoord{Q: 1, R: 0}
	line1 := HexLine(origin, adj)
	if len(line1) != 2 {
		t.Errorf("HexLine to adjacent has %d hexes, want 2", len(line1))
	}

	// Line across 3 hexes
	far := HexCoord{Q: 3, R: 0}
	line3 := HexLine(origin, far)
	if len(line3) != 4 {
		t.Errorf("HexLine across 3 steps has %d hexes, want 4", len(line3))
	}

	// Verify line is contiguous (each step is distance 1)
	for i := 1; i < len(line3); i++ {
		dist := line3[i-1].Distance(line3[i])
		if dist != 1 {
			t.Errorf("HexLine step %d to %d: distance = %d, want 1", i-1, i, dist)
		}
	}
}

func TestDirectionVector(t *testing.T) {
	// Verify each direction vector has length 1
	for dir := HexDirE; dir <= HexDirSE; dir++ {
		vec := DirectionVector(dir)
		if vec.Length() != 1 {
			t.Errorf("DirectionVector(%d) = %v has length %d, want 1", dir, vec, vec.Length())
		}
	}
}

func TestHexCoord_DirectionOpposites(t *testing.T) {
	// E and W should be opposites
	e := DirectionVector(HexDirE)
	w := DirectionVector(HexDirW)
	if e.Add(w) != (HexCoord{0, 0}) {
		t.Errorf("E + W = %v, want {0, 0}", e.Add(w))
	}

	// NE and SW should be opposites
	ne := DirectionVector(HexDirNE)
	sw := DirectionVector(HexDirSW)
	if ne.Add(sw) != (HexCoord{0, 0}) {
		t.Errorf("NE + SW = %v, want {0, 0}", ne.Add(sw))
	}

	// NW and SE should be opposites
	nw := DirectionVector(HexDirNW)
	se := DirectionVector(HexDirSE)
	if nw.Add(se) != (HexCoord{0, 0}) {
		t.Errorf("NW + SE = %v, want {0, 0}", nw.Add(se))
	}
}
