// Package core provides hex coordinate utilities for the Spectrex framework.
package core

import "math"

// HexCoord represents a hex coordinate in axial coordinate system.
// Uses Q (column) and R (row) coordinates, where the third cube coordinate
// S can be derived as S = -Q - R.
type HexCoord struct {
	Q, R int
}

// HexDirection represents the six cardinal directions on a hex grid.
type HexDirection int

const (
	HexDirE  HexDirection = iota // East (+Q)
	HexDirNE                     // Northeast (+Q, -R)
	HexDirNW                     // Northwest (-R)
	HexDirW                      // West (-Q)
	HexDirSW                     // Southwest (-Q, +R)
	HexDirSE                     // Southeast (+R)
)

// hexDirectionVectors maps each direction to its axial coordinate offset.
// These are for "pointy-top" hex orientation.
var hexDirectionVectors = [6]HexCoord{
	{Q: +1, R: 0},  // E
	{Q: +1, R: -1}, // NE
	{Q: 0, R: -1},  // NW
	{Q: -1, R: 0},  // W
	{Q: -1, R: +1}, // SW
	{Q: 0, R: +1},  // SE
}

// NewHexCoord creates a new hex coordinate.
func NewHexCoord(q, r int) HexCoord {
	return HexCoord{Q: q, R: r}
}

// S returns the third cube coordinate (derived from Q and R).
func (h HexCoord) S() int {
	return -h.Q - h.R
}

// Add returns the sum of two hex coordinates.
func (h HexCoord) Add(other HexCoord) HexCoord {
	return HexCoord{Q: h.Q + other.Q, R: h.R + other.R}
}

// Sub returns the difference of two hex coordinates.
func (h HexCoord) Sub(other HexCoord) HexCoord {
	return HexCoord{Q: h.Q - other.Q, R: h.R - other.R}
}

// Scale returns the hex coordinate scaled by a factor.
func (h HexCoord) Scale(k int) HexCoord {
	return HexCoord{Q: h.Q * k, R: h.R * k}
}

// Neighbor returns the adjacent hex in the given direction.
func (h HexCoord) Neighbor(dir HexDirection) HexCoord {
	return h.Add(hexDirectionVectors[dir])
}

// Neighbors returns all six adjacent hex coordinates.
func (h HexCoord) Neighbors() [6]HexCoord {
	var result [6]HexCoord
	for i := 0; i < 6; i++ {
		result[i] = h.Add(hexDirectionVectors[i])
	}
	return result
}

// Distance returns the distance between two hex coordinates.
// This is equivalent to the number of steps needed to move from one to the other.
func (h HexCoord) Distance(other HexCoord) int {
	dq := h.Q - other.Q
	dr := h.R - other.R
	ds := h.S() - other.S()
	return (abs(dq) + abs(dr) + abs(ds)) / 2
}

// Length returns the distance from the origin (0, 0).
func (h HexCoord) Length() int {
	return (abs(h.Q) + abs(h.R) + abs(h.S())) / 2
}

// Equal returns true if two hex coordinates are the same.
func (h HexCoord) Equal(other HexCoord) bool {
	return h.Q == other.Q && h.R == other.R
}

// HexCubeCoord represents a hex coordinate in cube coordinate system.
// Cube coordinates satisfy the constraint Q + R + S = 0.
type HexCubeCoord struct {
	Q, R, S int
}

// ToCube converts axial coordinates to cube coordinates.
func (h HexCoord) ToCube() HexCubeCoord {
	return HexCubeCoord{Q: h.Q, R: h.R, S: h.S()}
}

// ToAxial converts cube coordinates to axial coordinates.
func (c HexCubeCoord) ToAxial() HexCoord {
	return HexCoord{Q: c.Q, R: c.R}
}

// HexLayout defines the orientation and size for converting hex to pixel coordinates.
type HexLayout struct {
	Size   Vec2 // Size of each hex (width/2 and height/2 for pointy-top)
	Origin Vec2 // Pixel coordinate of hex (0, 0)
}

// NewHexLayout creates a new hex layout with the given size and origin.
func NewHexLayout(size, origin Vec2) HexLayout {
	return HexLayout{Size: size, Origin: origin}
}

// ToPixel converts a hex coordinate to pixel coordinates (center of hex).
// Uses pointy-top orientation.
func (l HexLayout) ToPixel(h HexCoord) Vec2 {
	// Pointy-top orientation matrix
	x := l.Size.X * (sqrt3*float32(h.Q) + sqrt3/2*float32(h.R))
	y := l.Size.Y * (3.0 / 2.0 * float32(h.R))
	return Vec2{X: x + l.Origin.X, Y: y + l.Origin.Y}
}

// FromPixel converts pixel coordinates to the nearest hex coordinate.
// Uses pointy-top orientation.
func (l HexLayout) FromPixel(p Vec2) HexCoord {
	// Inverse of pointy-top orientation matrix
	px := (p.X - l.Origin.X) / l.Size.X
	py := (p.Y - l.Origin.Y) / l.Size.Y

	q := sqrt3/3*px - 1.0/3*py
	r := 2.0 / 3 * py

	return hexRound(float64(q), float64(r))
}

// hexRound rounds fractional hex coordinates to the nearest integer hex coordinate.
func hexRound(q, r float64) HexCoord {
	s := -q - r

	rq := math.Round(q)
	rr := math.Round(r)
	rs := math.Round(s)

	qDiff := math.Abs(rq - q)
	rDiff := math.Abs(rr - r)
	sDiff := math.Abs(rs - s)

	// Reset the component with the largest difference
	if qDiff > rDiff && qDiff > sDiff {
		rq = -rr - rs
	} else if rDiff > sDiff {
		rr = -rq - rs
	}

	return HexCoord{Q: int(rq), R: int(rr)}
}

// HexRing returns all hex coordinates at exactly the given radius from center.
func HexRing(center HexCoord, radius int) []HexCoord {
	if radius <= 0 {
		return []HexCoord{center}
	}

	results := make([]HexCoord, 0, 6*radius)

	// Start at the hex radius steps in the SW direction
	hex := center.Add(hexDirectionVectors[HexDirSW].Scale(radius))

	// Walk around the ring
	for i := 0; i < 6; i++ {
		for j := 0; j < radius; j++ {
			results = append(results, hex)
			hex = hex.Neighbor(HexDirection(i))
		}
	}

	return results
}

// HexSpiral returns all hex coordinates within the given radius, starting from center.
func HexSpiral(center HexCoord, radius int) []HexCoord {
	results := []HexCoord{center}

	for r := 1; r <= radius; r++ {
		results = append(results, HexRing(center, r)...)
	}

	return results
}

// HexLine returns the hex coordinates on a line between two hexes.
func HexLine(a, b HexCoord) []HexCoord {
	n := a.Distance(b)
	if n == 0 {
		return []HexCoord{a}
	}

	results := make([]HexCoord, n+1)

	for i := 0; i <= n; i++ {
		t := float64(i) / float64(n)
		q := lerp(float64(a.Q), float64(b.Q), t)
		r := lerp(float64(a.R), float64(b.R), t)
		results[i] = hexRound(q, r)
	}

	return results
}

// Helper functions

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const sqrt3 = float32(1.7320508075688772) // math.Sqrt(3)

func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t
}

// Scale multiplies the direction vector by k.
func (h HexCoord) ScaleDir(k int) HexCoord {
	return HexCoord{Q: h.Q * k, R: h.R * k}
}

// DirectionVector returns the unit vector for a direction.
func DirectionVector(dir HexDirection) HexCoord {
	return hexDirectionVectors[dir]
}
