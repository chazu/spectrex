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
	RenderWidth  int32
	RenderHeight int32
	camera       rl.Camera3D

	// Render texture for resolution scaling
	renderTarget  rl.RenderTexture2D
	useRenderTex  bool
	windowResized bool
}

// NewRenderer creates a new raylib renderer with basic settings.
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
		RenderWidth:  screenWidth,
		RenderHeight: screenHeight,
		camera:       camera,
		useRenderTex: false,
	}
}

// NewRendererWithConfig creates a renderer using DisplayConfig.
// This handles window creation, maximization, and render texture setup.
func NewRendererWithConfig(config core.DisplayConfig) *Renderer {
	// Set window flags before creation
	if config.Resizable {
		rl.SetConfigFlags(rl.FlagWindowResizable)
	}
	if config.VSync {
		rl.SetConfigFlags(rl.FlagVsyncHint)
	}

	// Create window
	rl.InitWindow(config.WindowWidth, config.WindowHeight, config.Title)

	// Maximize if requested
	if config.Maximized {
		rl.MaximizeWindow()
	}

	// Set target FPS
	if config.TargetFPS > 0 {
		rl.SetTargetFPS(config.TargetFPS)
	}

	// Determine render size
	renderW, renderH := config.EffectiveRenderSize()

	camera := rl.Camera3D{
		Position:   rl.Vector3{X: 0, Y: 100, Z: -300},
		Target:     rl.Vector3{X: 0, Y: 0, Z: 100},
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       config.DefaultFOV,
		Projection: rl.CameraPerspective,
	}

	// Always use render texture if a specific render size is configured
	// This ensures consistent rendering regardless of window size
	useRenderTex := config.RenderWidth > 0 && config.RenderHeight > 0

	r := &Renderer{
		ScreenWidth:  int32(rl.GetScreenWidth()),
		ScreenHeight: int32(rl.GetScreenHeight()),
		RenderWidth:  renderW,
		RenderHeight: renderH,
		camera:       camera,
		useRenderTex: useRenderTex,
	}

	// Create render texture if using fixed resolution
	if r.useRenderTex {
		r.renderTarget = rl.LoadRenderTexture(renderW, renderH)
	}

	return r
}

// Close releases renderer resources.
func (r *Renderer) Close() {
	if r.useRenderTex {
		rl.UnloadRenderTexture(r.renderTarget)
	}
}

// HandleResize updates renderer when window is resized.
func (r *Renderer) HandleResize() {
	newW := int32(rl.GetScreenWidth())
	newH := int32(rl.GetScreenHeight())
	if newW != r.ScreenWidth || newH != r.ScreenHeight {
		r.ScreenWidth = newW
		r.ScreenHeight = newH
		r.windowResized = true
	}
}

// BeginFrame begins a new frame.
func (r *Renderer) BeginFrame() {
	r.HandleResize()

	if r.useRenderTex {
		rl.BeginTextureMode(r.renderTarget)
		rl.ClearBackground(rl.Black)
	} else {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
	}
}

// End3DAndBlit ends 3D rendering and blits the render texture if used.
// After this, you can draw 2D overlays directly to the screen.
func (r *Renderer) End3DAndBlit() {
	if r.useRenderTex {
		rl.EndTextureMode()

		// Update screen size in case window was resized
		r.ScreenWidth = int32(rl.GetScreenWidth())
		r.ScreenHeight = int32(rl.GetScreenHeight())

		// Draw render texture scaled to window
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Calculate scaling to fit window while maintaining aspect ratio
		srcRect := rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(r.RenderWidth),
			Height: -float32(r.RenderHeight), // Negative to flip Y
		}

		// Scale to fit window
		scale := min(
			float32(r.ScreenWidth)/float32(r.RenderWidth),
			float32(r.ScreenHeight)/float32(r.RenderHeight),
		)
		destW := float32(r.RenderWidth) * scale
		destH := float32(r.RenderHeight) * scale
		destX := (float32(r.ScreenWidth) - destW) / 2
		destY := (float32(r.ScreenHeight) - destH) / 2

		destRect := rl.Rectangle{
			X:      destX,
			Y:      destY,
			Width:  destW,
			Height: destH,
		}

		rl.DrawTexturePro(r.renderTarget.Texture, srcRect, destRect, rl.Vector2{}, 0, rl.White)
		// Don't EndDrawing yet - allow screen overlays
	}
	// If not using render tex, we're already in drawing mode
}

// EndFrame ends the current frame. Call after all drawing is complete.
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
