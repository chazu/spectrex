// Package core provides hex grid storage utilities for the Spectrex framework.
package core

// HexGrid is a generic container for storing values at hex coordinates.
// It uses a radius-based layout where all hexes within the specified radius
// from the origin (0,0) are valid positions.
//
// For radius R, the grid contains (3*R*R + 3*R + 1) cells:
// - Radius 0: 1 cell (just the center)
// - Radius 1: 7 cells
// - Radius 2: 19 cells
// - Radius 3: 37 cells
// - Radius 4: 61 cells
type HexGrid[T any] struct {
	radius int
	data   map[HexCoord]T
}

// NewHexGrid creates a new hex grid with the given radius.
// The grid will contain all hexes within the radius from the origin.
func NewHexGrid[T any](radius int) *HexGrid[T] {
	if radius < 0 {
		radius = 0
	}
	return &HexGrid[T]{
		radius: radius,
		data:   make(map[HexCoord]T),
	}
}

// Radius returns the radius of the grid.
func (g *HexGrid[T]) Radius() int {
	return g.radius
}

// Size returns the total number of valid cells in the grid.
// Formula: 3*r*r + 3*r + 1 for radius r.
func (g *HexGrid[T]) Size() int {
	r := g.radius
	return 3*r*r + 3*r + 1
}

// IsValid returns true if the coordinate is within the grid's radius.
func (g *HexGrid[T]) IsValid(coord HexCoord) bool {
	return coord.Length() <= g.radius
}

// Get returns the value at the given coordinate.
// Returns the zero value if the coordinate is invalid or not set.
func (g *HexGrid[T]) Get(coord HexCoord) T {
	var zero T
	if !g.IsValid(coord) {
		return zero
	}
	return g.data[coord]
}

// GetOk returns the value at the given coordinate and whether it was found.
// Returns (zero, false) if the coordinate is invalid.
// Returns (zero, false) if valid but never set.
// Returns (value, true) if the coordinate has been set.
func (g *HexGrid[T]) GetOk(coord HexCoord) (T, bool) {
	var zero T
	if !g.IsValid(coord) {
		return zero, false
	}
	val, ok := g.data[coord]
	return val, ok
}

// Set stores a value at the given coordinate.
// Returns false if the coordinate is outside the grid's radius.
func (g *HexGrid[T]) Set(coord HexCoord, value T) bool {
	if !g.IsValid(coord) {
		return false
	}
	g.data[coord] = value
	return true
}

// Delete removes the value at the given coordinate.
// Returns false if the coordinate is outside the grid's radius.
func (g *HexGrid[T]) Delete(coord HexCoord) bool {
	if !g.IsValid(coord) {
		return false
	}
	delete(g.data, coord)
	return true
}

// Clear removes all values from the grid.
func (g *HexGrid[T]) Clear() {
	g.data = make(map[HexCoord]T)
}

// Count returns the number of cells that have been set.
func (g *HexGrid[T]) Count() int {
	return len(g.data)
}

// All returns all valid coordinates in the grid, starting from center.
// Uses spiral ordering (center first, then expanding rings).
func (g *HexGrid[T]) All() []HexCoord {
	return HexSpiral(HexCoord{Q: 0, R: 0}, g.radius)
}

// Ring returns all coordinates at exactly the given distance from center.
// Returns nil if distance is greater than the grid's radius.
func (g *HexGrid[T]) Ring(distance int) []HexCoord {
	if distance > g.radius {
		return nil
	}
	return HexRing(HexCoord{Q: 0, R: 0}, distance)
}

// ForEach calls the function for each valid coordinate in the grid.
// Iterates in spiral order (center first, then expanding rings).
func (g *HexGrid[T]) ForEach(fn func(coord HexCoord, value T)) {
	for _, coord := range g.All() {
		fn(coord, g.data[coord])
	}
}

// ForEachSet calls the function for each coordinate that has a value set.
// Order is not guaranteed.
func (g *HexGrid[T]) ForEachSet(fn func(coord HexCoord, value T)) {
	for coord, value := range g.data {
		fn(coord, value)
	}
}

// ForEachRing calls the function for each coordinate at the given distance.
// Returns false if distance is greater than the grid's radius.
func (g *HexGrid[T]) ForEachRing(distance int, fn func(coord HexCoord, value T)) bool {
	if distance > g.radius {
		return false
	}
	for _, coord := range g.Ring(distance) {
		fn(coord, g.data[coord])
	}
	return true
}

// Neighbors returns valid neighbors of the given coordinate.
// Only returns neighbors that are within the grid's radius.
func (g *HexGrid[T]) Neighbors(coord HexCoord) []HexCoord {
	all := coord.Neighbors()
	result := make([]HexCoord, 0, 6)
	for _, n := range all {
		if g.IsValid(n) {
			result = append(result, n)
		}
	}
	return result
}

// Fill sets all valid coordinates to the given value.
func (g *HexGrid[T]) Fill(value T) {
	for _, coord := range g.All() {
		g.data[coord] = value
	}
}

// Clone creates a deep copy of the grid.
func (g *HexGrid[T]) Clone() *HexGrid[T] {
	clone := &HexGrid[T]{
		radius: g.radius,
		data:   make(map[HexCoord]T, len(g.data)),
	}
	for k, v := range g.data {
		clone.data[k] = v
	}
	return clone
}
