package core

import "testing"

func TestNewHexGrid(t *testing.T) {
	tests := []struct {
		radius   int
		wantSize int
	}{
		{0, 1},
		{1, 7},
		{2, 19},
		{3, 37},
		{4, 61},
		{-1, 1}, // negative radius treated as 0
	}

	for _, tt := range tests {
		grid := NewHexGrid[int](tt.radius)
		if got := grid.Size(); got != tt.wantSize {
			t.Errorf("NewHexGrid(%d).Size() = %d, want %d", tt.radius, got, tt.wantSize)
		}
	}
}

func TestHexGridRadius(t *testing.T) {
	grid := NewHexGrid[int](4)
	if grid.Radius() != 4 {
		t.Errorf("Radius() = %d, want 4", grid.Radius())
	}
}

func TestHexGridIsValid(t *testing.T) {
	grid := NewHexGrid[int](2)

	tests := []struct {
		coord HexCoord
		want  bool
	}{
		{HexCoord{0, 0}, true},   // center
		{HexCoord{1, 0}, true},   // distance 1
		{HexCoord{0, 1}, true},   // distance 1
		{HexCoord{-1, 0}, true},  // distance 1
		{HexCoord{2, 0}, true},   // distance 2
		{HexCoord{1, 1}, true},   // distance 2
		{HexCoord{-2, 1}, true},  // distance 2
		{HexCoord{3, 0}, false},  // distance 3 (out of radius 2)
		{HexCoord{2, 1}, false},  // distance 3
		{HexCoord{-1, -2}, false}, // distance 3
	}

	for _, tt := range tests {
		if got := grid.IsValid(tt.coord); got != tt.want {
			t.Errorf("IsValid(%v) = %v, want %v (distance=%d)",
				tt.coord, got, tt.want, tt.coord.Length())
		}
	}
}

func TestHexGridSetGet(t *testing.T) {
	grid := NewHexGrid[string](2)

	// Set value at valid coordinate
	if ok := grid.Set(HexCoord{1, 0}, "hello"); !ok {
		t.Error("Set at valid coord should return true")
	}

	// Get value
	if got := grid.Get(HexCoord{1, 0}); got != "hello" {
		t.Errorf("Get(%v) = %q, want %q", HexCoord{1, 0}, got, "hello")
	}

	// Get unset value returns zero
	if got := grid.Get(HexCoord{0, 1}); got != "" {
		t.Errorf("Get unset coord = %q, want empty string", got)
	}

	// Set at invalid coordinate fails
	if ok := grid.Set(HexCoord{5, 0}, "nope"); ok {
		t.Error("Set at invalid coord should return false")
	}

	// Get at invalid coordinate returns zero
	if got := grid.Get(HexCoord{5, 0}); got != "" {
		t.Errorf("Get invalid coord = %q, want empty string", got)
	}
}

func TestHexGridGetOk(t *testing.T) {
	grid := NewHexGrid[int](2)
	grid.Set(HexCoord{0, 0}, 42)

	// Set value
	val, ok := grid.GetOk(HexCoord{0, 0})
	if !ok || val != 42 {
		t.Errorf("GetOk set value = (%d, %v), want (42, true)", val, ok)
	}

	// Unset but valid coordinate
	val, ok = grid.GetOk(HexCoord{1, 0})
	if ok || val != 0 {
		t.Errorf("GetOk unset = (%d, %v), want (0, false)", val, ok)
	}

	// Invalid coordinate
	val, ok = grid.GetOk(HexCoord{5, 0})
	if ok || val != 0 {
		t.Errorf("GetOk invalid = (%d, %v), want (0, false)", val, ok)
	}
}

func TestHexGridDelete(t *testing.T) {
	grid := NewHexGrid[int](2)
	grid.Set(HexCoord{1, 1}, 99)

	// Delete existing
	if ok := grid.Delete(HexCoord{1, 1}); !ok {
		t.Error("Delete valid coord should return true")
	}
	if _, ok := grid.GetOk(HexCoord{1, 1}); ok {
		t.Error("After delete, GetOk should return false")
	}

	// Delete invalid coordinate
	if ok := grid.Delete(HexCoord{5, 0}); ok {
		t.Error("Delete invalid coord should return false")
	}
}

func TestHexGridClear(t *testing.T) {
	grid := NewHexGrid[int](2)
	grid.Set(HexCoord{0, 0}, 1)
	grid.Set(HexCoord{1, 0}, 2)
	grid.Set(HexCoord{0, 1}, 3)

	if grid.Count() != 3 {
		t.Errorf("Before clear, Count() = %d, want 3", grid.Count())
	}

	grid.Clear()

	if grid.Count() != 0 {
		t.Errorf("After clear, Count() = %d, want 0", grid.Count())
	}
}

func TestHexGridCount(t *testing.T) {
	grid := NewHexGrid[int](2)

	if grid.Count() != 0 {
		t.Errorf("Empty grid Count() = %d, want 0", grid.Count())
	}

	grid.Set(HexCoord{0, 0}, 1)
	grid.Set(HexCoord{1, 0}, 2)

	if grid.Count() != 2 {
		t.Errorf("After 2 sets, Count() = %d, want 2", grid.Count())
	}
}

func TestHexGridAll(t *testing.T) {
	grid := NewHexGrid[int](2)
	all := grid.All()

	if len(all) != 19 {
		t.Errorf("All() returned %d coords, want 19", len(all))
	}

	// First should be center
	if !all[0].Equal(HexCoord{0, 0}) {
		t.Errorf("All()[0] = %v, want (0, 0)", all[0])
	}

	// All should be valid
	for _, coord := range all {
		if !grid.IsValid(coord) {
			t.Errorf("All() returned invalid coord %v", coord)
		}
	}

	// Check uniqueness
	seen := make(map[HexCoord]bool)
	for _, coord := range all {
		if seen[coord] {
			t.Errorf("All() returned duplicate coord %v", coord)
		}
		seen[coord] = true
	}
}

func TestHexGridRing(t *testing.T) {
	grid := NewHexGrid[int](3)

	// Ring at distance 0 is just center
	ring0 := grid.Ring(0)
	if len(ring0) != 1 || !ring0[0].Equal(HexCoord{0, 0}) {
		t.Errorf("Ring(0) = %v, want [(0,0)]", ring0)
	}

	// Ring at distance 1 has 6 hexes
	ring1 := grid.Ring(1)
	if len(ring1) != 6 {
		t.Errorf("Ring(1) length = %d, want 6", len(ring1))
	}

	// Ring at distance 2 has 12 hexes
	ring2 := grid.Ring(2)
	if len(ring2) != 12 {
		t.Errorf("Ring(2) length = %d, want 12", len(ring2))
	}

	// Ring at distance 3 has 18 hexes
	ring3 := grid.Ring(3)
	if len(ring3) != 18 {
		t.Errorf("Ring(3) length = %d, want 18", len(ring3))
	}

	// Ring beyond radius returns nil
	ring4 := grid.Ring(4)
	if ring4 != nil {
		t.Errorf("Ring(4) = %v, want nil", ring4)
	}
}

func TestHexGridForEach(t *testing.T) {
	grid := NewHexGrid[int](1)
	grid.Set(HexCoord{0, 0}, 10)
	grid.Set(HexCoord{1, 0}, 20)

	count := 0
	sum := 0
	grid.ForEach(func(coord HexCoord, value int) {
		count++
		sum += value
	})

	if count != 7 {
		t.Errorf("ForEach visited %d coords, want 7", count)
	}
	if sum != 30 {
		t.Errorf("ForEach sum = %d, want 30", sum)
	}
}

func TestHexGridForEachSet(t *testing.T) {
	grid := NewHexGrid[int](2)
	grid.Set(HexCoord{0, 0}, 10)
	grid.Set(HexCoord{1, 0}, 20)
	grid.Set(HexCoord{-1, 1}, 30)

	count := 0
	sum := 0
	grid.ForEachSet(func(coord HexCoord, value int) {
		count++
		sum += value
	})

	if count != 3 {
		t.Errorf("ForEachSet visited %d coords, want 3", count)
	}
	if sum != 60 {
		t.Errorf("ForEachSet sum = %d, want 60", sum)
	}
}

func TestHexGridForEachRing(t *testing.T) {
	grid := NewHexGrid[int](2)
	grid.Fill(1) // All cells = 1

	ring1Sum := 0
	ok := grid.ForEachRing(1, func(coord HexCoord, value int) {
		ring1Sum += value
	})
	if !ok {
		t.Error("ForEachRing(1) should return true")
	}
	if ring1Sum != 6 {
		t.Errorf("ForEachRing(1) sum = %d, want 6", ring1Sum)
	}

	// Ring beyond radius
	ok = grid.ForEachRing(3, func(coord HexCoord, value int) {})
	if ok {
		t.Error("ForEachRing(3) should return false for radius 2 grid")
	}
}

func TestHexGridNeighbors(t *testing.T) {
	grid := NewHexGrid[int](1)

	// Center has all 6 neighbors
	neighbors := grid.Neighbors(HexCoord{0, 0})
	if len(neighbors) != 6 {
		t.Errorf("Center neighbors = %d, want 6", len(neighbors))
	}

	// Edge cell has fewer neighbors
	neighbors = grid.Neighbors(HexCoord{1, 0})
	if len(neighbors) != 3 {
		t.Errorf("Edge neighbors = %d, want 3", len(neighbors))
	}

	// Corner cell (1,-1) has neighbors at (0,-1), (0,0), and (1,0) - all at distance ≤1
	neighbors = grid.Neighbors(HexCoord{1, -1})
	if len(neighbors) != 3 {
		t.Errorf("Corner neighbors = %d, want 3", len(neighbors))
	}
}

func TestHexGridFill(t *testing.T) {
	grid := NewHexGrid[string](1)
	grid.Fill("x")

	if grid.Count() != 7 {
		t.Errorf("After Fill, Count() = %d, want 7", grid.Count())
	}

	for _, coord := range grid.All() {
		if v := grid.Get(coord); v != "x" {
			t.Errorf("After Fill, Get(%v) = %q, want 'x'", coord, v)
		}
	}
}

func TestHexGridClone(t *testing.T) {
	grid := NewHexGrid[int](2)
	grid.Set(HexCoord{0, 0}, 100)
	grid.Set(HexCoord{1, 1}, 200)

	clone := grid.Clone()

	// Same values
	if clone.Get(HexCoord{0, 0}) != 100 || clone.Get(HexCoord{1, 1}) != 200 {
		t.Error("Clone has different values")
	}

	// Same radius
	if clone.Radius() != grid.Radius() {
		t.Error("Clone has different radius")
	}

	// Independent - modifying clone doesn't affect original
	clone.Set(HexCoord{0, 0}, 999)
	if grid.Get(HexCoord{0, 0}) != 100 {
		t.Error("Modifying clone affected original")
	}
}

func TestHexGridRadius4Size(t *testing.T) {
	// Verify the specific radius 4 = 61 hexes requirement
	grid := NewHexGrid[int](4)

	if grid.Size() != 61 {
		t.Errorf("Radius 4 grid size = %d, want 61", grid.Size())
	}

	all := grid.All()
	if len(all) != 61 {
		t.Errorf("Radius 4 All() returned %d, want 61", len(all))
	}
}

func TestHexGridWithStructType(t *testing.T) {
	type Cell struct {
		Value int
		Name  string
	}

	grid := NewHexGrid[Cell](1)
	grid.Set(HexCoord{0, 0}, Cell{Value: 42, Name: "center"})

	got := grid.Get(HexCoord{0, 0})
	if got.Value != 42 || got.Name != "center" {
		t.Errorf("Get struct = %+v, want {Value:42 Name:center}", got)
	}
}

func TestHexGridWithPointerType(t *testing.T) {
	type Data struct {
		Value int
	}

	grid := NewHexGrid[*Data](1)
	d := &Data{Value: 99}
	grid.Set(HexCoord{0, 0}, d)

	got := grid.Get(HexCoord{0, 0})
	if got != d || got.Value != 99 {
		t.Error("Pointer type not stored correctly")
	}

	// Nil for unset
	if grid.Get(HexCoord{1, 0}) != nil {
		t.Error("Unset pointer should be nil")
	}
}
