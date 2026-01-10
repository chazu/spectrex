// Package raylib provides conversion utilities between core and raylib types.
package raylib

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/spectrex/core"
)

// coreToRlVec3 converts core.Vec3 to rl.Vector3.
func coreToRlVec3(v core.Vec3) rl.Vector3 {
	return rl.Vector3{X: v.X, Y: v.Y, Z: v.Z}
}

// rlToCoreVec3 converts rl.Vector3 to core.Vec3.
func rlToCoreVec3(v rl.Vector3) core.Vec3 {
	return core.Vec3{X: v.X, Y: v.Y, Z: v.Z}
}

// coreToRlVec2 converts core.Vec2 to rl.Vector2.
func coreToRlVec2(v core.Vec2) rl.Vector2 {
	return rl.Vector2{X: v.X, Y: v.Y}
}

// rlToCoreVec2 converts rl.Vector2 to core.Vec2.
func rlToCoreVec2(v rl.Vector2) core.Vec2 {
	return core.Vec2{X: v.X, Y: v.Y}
}

// coreToRlColor converts core.Color to rl.Color.
func coreToRlColor(c core.Color) rl.Color {
	return rl.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// rlToCoreColor converts rl.Color to core.Color.
func rlToCoreColor(c rl.Color) core.Color {
	return core.Color{R: c.R, G: c.G, B: c.B, A: c.A}
}

// coreToRlMatrix converts core.Matrix to rl.Matrix.
func coreToRlMatrix(m core.Matrix) rl.Matrix {
	return rl.Matrix{
		M0: m[0], M1: m[1], M2: m[2], M3: m[3],
		M4: m[4], M5: m[5], M6: m[6], M7: m[7],
		M8: m[8], M9: m[9], M10: m[10], M11: m[11],
		M12: m[12], M13: m[13], M14: m[14], M15: m[15],
	}
}

// rlToCoreMatrix converts rl.Matrix to core.Matrix.
func rlToCoreMatrix(m rl.Matrix) core.Matrix {
	return core.Matrix{
		m.M0, m.M1, m.M2, m.M3,
		m.M4, m.M5, m.M6, m.M7,
		m.M8, m.M9, m.M10, m.M11,
		m.M12, m.M13, m.M14, m.M15,
	}
}

// coreToRlCamera converts core.Camera to rl.Camera3D.
func coreToRlCamera(c core.Camera) rl.Camera3D {
	projection := rl.CameraPerspective
	if c.Projection == 1 {
		projection = rl.CameraOrthographic
	}
	return rl.Camera3D{
		Position:   coreToRlVec3(c.Position),
		Target:     coreToRlVec3(c.Target),
		Up:         coreToRlVec3(c.Up),
		Fovy:       c.Fovy,
		Projection: projection,
	}
}

// rlToCoreCamera converts rl.Camera3D to core.Camera.
func rlToCoreCamera(c rl.Camera3D) core.Camera {
	projection := 0
	if c.Projection == rl.CameraOrthographic {
		projection = 1
	}
	return core.Camera{
		Position:   rlToCoreVec3(c.Position),
		Target:     rlToCoreVec3(c.Target),
		Up:         rlToCoreVec3(c.Up),
		Fovy:       c.Fovy,
		Projection: projection,
	}
}

// Vec3Transform transforms a Vec3 using a raylib matrix.
func Vec3Transform(v core.Vec3, m rl.Matrix) core.Vec3 {
	result := rl.Vector3Transform(coreToRlVec3(v), m)
	return rlToCoreVec3(result)
}
