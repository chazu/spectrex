// Package raylib provides a raylib-based renderer for the Spectrex framework.
package raylib

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/spectrex/core"
)

// Renderer implements core.Renderer using raylib.
type Renderer struct {
	ScreenWidth  int32
	ScreenHeight int32
	camera       rl.Camera3D
}

// NewRenderer creates a new raylib renderer.
func NewRenderer(screenWidth, screenHeight int32) *Renderer {
	camera := rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 100, Z: -300},
		Target:     rl.Vector3{X: 0, Y: 0, Z: 100},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	return &Renderer{
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		camera:       camera,
	}
}

// BeginFrame begins a new frame.
func (r *Renderer) BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
}

// EndFrame ends the current frame.
func (r *Renderer) EndFrame() {
	rl.EndDrawing()
}

// Begin3D begins 3D rendering with the specified camera.
func (r *Renderer) Begin3D(camera core.Camera) {
	r.camera = coreToRlCamera(camera)
	rl.BeginMode3D(r.camera)
}

// End3D ends 3D rendering.
func (r *Renderer) End3D() {
	rl.EndMode3D()
}

// DrawLine3D draws a 3D line.
func (r *Renderer) DrawLine3D(start, end core.Vec3, color core.Color) {
	rl.DrawLine3D(coreToRlVec3(start), coreToRlVec3(end), coreToRlColor(color))
}

// DrawTriangle3D draws a 3D triangle.
func (r *Renderer) DrawTriangle3D(v1, v2, v3 core.Vec3, color core.Color) {
	rl.DrawTriangle3D(coreToRlVec3(v1), coreToRlVec3(v2), coreToRlVec3(v3), coreToRlColor(color))
}

// DrawGrid draws a reference grid.
func (r *Renderer) DrawGrid(slices int, spacing float32) {
	rl.DrawGrid(int32(slices), spacing)
}

// DrawFPS draws the current FPS.
func (r *Renderer) DrawFPS(x, y int32) {
	rl.DrawFPS(x, y)
}

// DrawText2D draws 2D text on screen.
func (r *Renderer) DrawText2D(text string, x, y int32, fontSize int32, color core.Color) {
	rl.DrawText(text, x, y, fontSize, coreToRlColor(color))
}

// GetScreenWidth returns the screen width.
func (r *Renderer) GetScreenWidth() int32 {
	return r.ScreenWidth
}

// GetScreenHeight returns the screen height.
func (r *Renderer) GetScreenHeight() int32 {
	return r.ScreenHeight
}

// GetCamera returns the current raylib camera (for advanced use).
func (r *Renderer) GetCamera() rl.Camera3D {
	return r.camera
}

// SetCamera sets the raylib camera directly (for advanced use).
func (r *Renderer) SetCamera(camera rl.Camera3D) {
	r.camera = camera
}
