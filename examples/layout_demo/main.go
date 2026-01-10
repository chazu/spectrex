// Layout demo for Spectrex showing text screens and spinning hexagons.
package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/spectrex/backends/raylib"
	"github.com/chazu/spectrex/core"
)

const loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.`

const loremShort = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore.`

func main() {
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "Spectrex Layout Demo")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	renderer := raylib.NewRenderer(screenWidth, screenHeight)
	textScreenRenderer := raylib.NewTextScreenRenderer()

	font := core.LoadHersheyFontData()

	// Top-left text screen (simple text, no background)
	// In 3D: higher Y = higher on screen, lower X = left side
	topLeftScreen := core.NewTextScreen(
		core.Vec3{X: -50, Y: 50, Z: 200},
		250, 120, 1.0,
	)
	topLeftScreen.SetDebug(false)
	topLeftRegion := topLeftScreen.AddRegion(0, 0, 250, 120)
	topLeftRegion.SetContent(loremShort, font, core.ColorGreen)
	topLeftRegion.Scale = 0.8
	topLeftRegion.SetAlignment(core.AlignLeft, core.AlignTop)
	topLeftRegion.LineSpacing = 1.3

	// Bottom-right transparent text screen (with visible border)
	bottomRightScreen := core.NewTextScreen(
		core.Vec3{X: 280, Y: -120, Z: 200},
		280, 150, 1.0,
	)
	bottomRightScreen.SetTransparency(true)
	bottomRightScreen.SetBorder(true, core.ColorSkyBlue)
	bottomRightRegion := bottomRightScreen.AddRegion(10, 10, 260, 130)
	bottomRightRegion.SetContent(loremIpsum, font, core.ColorSkyBlue)
	bottomRightRegion.Scale = 0.7
	bottomRightRegion.SetAlignment(core.AlignLeft, core.AlignTop)
	bottomRightRegion.LineSpacing = 1.2

	// Top-right opaque text screen (with dark background)
	topRightScreen := core.NewTextScreen(
		core.Vec3{X: 280, Y: 50, Z: 200},
		280, 120, 1.0,
	)
	topRightScreen.SetTransparency(false)
	topRightScreen.SetBackground(core.Color{R: 30, G: 30, B: 50, A: 255})
	topRightScreen.SetBorder(true, core.ColorOrange)
	topRightRegion := topRightScreen.AddRegion(10, 10, 260, 100)
	topRightRegion.SetContent(loremShort, font, core.ColorOrange)
	topRightRegion.Scale = 0.8
	topRightRegion.SetAlignment(core.AlignLeft, core.AlignTop)
	topRightRegion.LineSpacing = 1.3

	// Create hexagons for center spinning
	hexRadius := float32(40)
	hex1 := core.MakePoly(6, hexRadius, 0)
	hex2 := core.MakePoly(6, hexRadius*0.6, math.Pi/6)
	hex3 := core.MakePoly(6, hexRadius*1.4, 0)

	camera := core.Camera{
		Position:   core.Vec3{X: 0, Y: 50, Z: -400},
		Target:     core.Vec3{X: 0, Y: 50, Z: 100},
		Up:         core.Vec3{X: 0, Y: 1, Z: 0},
		Fovy:       45.0,
		Projection: 0,
	}

	totalTime := float32(0)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEscape) {
			break
		}

		deltaTime := rl.GetFrameTime()
		totalTime += deltaTime

		renderer.BeginFrame()
		renderer.Begin3D(camera)

		// Draw text screens
		textScreenRenderer.DrawTextScreen(topLeftScreen)
		textScreenRenderer.DrawTextScreen(bottomRightScreen)
		textScreenRenderer.DrawTextScreen(topRightScreen)

		// Draw spinning hexagons in center
		centerPos := core.Vec3{X: -50, Y: 50, Z: 200}

		// Hexagon 1: spinning on Y axis
		rot1 := core.Vec3{X: 0, Y: totalTime * 60, Z: 0}
		transformed1 := core.TransformPoly(hex1, centerPos, rot1)
		drawPolygon(renderer, transformed1, core.ColorYellow)

		// Hexagon 2: spinning on X axis, offset position
		pos2 := core.Vec3{X: centerPos.X, Y: centerPos.Y, Z: centerPos.Z}
		rot2 := core.Vec3{X: totalTime * 45, Y: 0, Z: totalTime * 30}
		transformed2 := core.TransformPoly(hex2, pos2, rot2)
		drawPolygon(renderer, transformed2, core.ColorLime)

		// Hexagon 3: slow spin, outer ring
		rot3 := core.Vec3{X: 0, Y: -totalTime * 20, Z: totalTime * 10}
		transformed3 := core.TransformPoly(hex3, centerPos, rot3)
		drawPolygon(renderer, transformed3, core.ColorRed)

		renderer.DrawGrid(10, 10.0)
		renderer.End3D()

		renderer.DrawFPS(10, 10)
		renderer.DrawText2D("Layout Demo - Press ESC to exit", 10, 40, 20, core.ColorWhite)

		renderer.EndFrame()
	}
}

// drawPolygon draws a polygon as connected line segments.
func drawPolygon(renderer *raylib.Renderer, vertices []core.Vec3, color core.Color) {
	n := len(vertices)
	for i := 0; i < n; i++ {
		renderer.DrawLine3D(vertices[i], vertices[(i+1)%n], color)
	}
}
