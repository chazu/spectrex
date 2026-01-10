// Package core provides rendering-agnostic types and logic for the Spectrex
// vector UI framework. These types are designed to be portable across different
// rendering backends (raylib, SDL, OpenGL, terminal, etc.).
package core

import "math"

// Vec2 represents a 2D vector.
type Vec2 struct {
	X, Y float32
}

// Vec3 represents a 3D vector.
type Vec3 struct {
	X, Y, Z float32
}

// Add returns the sum of two vectors.
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// Sub returns the difference of two vectors.
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// Scale returns the vector scaled by a factor.
func (v Vec3) Scale(s float32) Vec3 {
	return Vec3{X: v.X * s, Y: v.Y * s, Z: v.Z * s}
}

// Color represents an RGBA color.
type Color struct {
	R, G, B, A uint8
}

// Common colors
var (
	ColorWhite   = Color{255, 255, 255, 255}
	ColorBlack   = Color{0, 0, 0, 255}
	ColorRed     = Color{255, 0, 0, 255}
	ColorGreen   = Color{0, 255, 0, 255}
	ColorBlue    = Color{0, 0, 255, 255}
	ColorYellow  = Color{255, 255, 0, 255}
	ColorOrange  = Color{255, 165, 0, 255}
	ColorSkyBlue = Color{135, 206, 235, 255}
	ColorLime    = Color{50, 205, 50, 255}
)

// Matrix represents a 4x4 transformation matrix.
type Matrix [16]float32

// MatrixIdentity returns an identity matrix.
func MatrixIdentity() Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// MatrixTranslate creates a translation matrix.
func MatrixTranslate(x, y, z float32) Matrix {
	return Matrix{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	}
}

// MatrixRotateX creates a rotation matrix around the X axis.
func MatrixRotateX(angle float32) Matrix {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	return Matrix{
		1, 0, 0, 0,
		0, c, s, 0,
		0, -s, c, 0,
		0, 0, 0, 1,
	}
}

// MatrixRotateY creates a rotation matrix around the Y axis.
func MatrixRotateY(angle float32) Matrix {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	return Matrix{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

// MatrixRotateZ creates a rotation matrix around the Z axis.
func MatrixRotateZ(angle float32) Matrix {
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))
	return Matrix{
		c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// Multiply multiplies two matrices.
func (m Matrix) Multiply(other Matrix) Matrix {
	var result Matrix
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := float32(0)
			for k := 0; k < 4; k++ {
				sum += m[i*4+k] * other[k*4+j]
			}
			result[i*4+j] = sum
		}
	}
	return result
}

// TransformVec3 transforms a Vec3 by this matrix.
func (m Matrix) TransformVec3(v Vec3) Vec3 {
	return Vec3{
		X: v.X*m[0] + v.Y*m[4] + v.Z*m[8] + m[12],
		Y: v.X*m[1] + v.Y*m[5] + v.Z*m[9] + m[13],
		Z: v.X*m[2] + v.Y*m[6] + v.Z*m[10] + m[14],
	}
}

// Pi is the mathematical constant.
const Pi = float32(math.Pi)

// DegToRad converts degrees to radians.
func DegToRad(degrees float32) float32 {
	return degrees * (Pi / 180.0)
}

// RadToDeg converts radians to degrees.
func RadToDeg(radians float32) float32 {
	return radians * (180.0 / Pi)
}
