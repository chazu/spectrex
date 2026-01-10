// Example demo application for Spectrex.
package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/spectrex/backends/raylib"
	"github.com/chazu/spectrex/core"
)

func main() {
	// Initialize window
	screenWidth := int32(1280)
	screenHeight := int32(720)

	rl.InitWindow(screenWidth, screenHeight, "Spectrex Demo")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// Create renderer
	renderer := raylib.NewRenderer(screenWidth, screenHeight)
	textScreenRenderer := raylib.NewTextScreenRenderer()

	// Load Hershey font
	font := core.LoadHersheyFontData()

	// Create a text screen
	textScreen := core.NewTextScreen(
		core.Vec3{X: 0, Y: 50, Z: 200},
		600, 400, 1.0,
	)
	textScreen.SetDebug(true)

	// Create a text document
	doc := core.NewTextDocument(textScreen, 2, 20)
	doc.PageStyle.Font = font

	// Add title
	titleSection := doc.AddSection("", "SPECTREX")
	titleSection.SetStyle(core.TextStyle{
		Font:        font,
		Color:       core.ColorGreen,
		Scale:       2.5,
		LineSpacing: 1.2,
		HAlign:      core.AlignCenter,
		VAlign:      core.AlignTop,
		WordWrap:    true,
	})

	// Add subtitle
	subtitleSection := doc.AddSection("", "Vector UI Framework")
	subtitleSection.SetStyle(core.TextStyle{
		Font:        font,
		Color:       core.ColorSkyBlue,
		Scale:       1.5,
		LineSpacing: 1.2,
		HAlign:      core.AlignCenter,
		VAlign:      core.AlignTop,
		WordWrap:    true,
	})

	// Add content
	contentSection := doc.AddSection("Features", "Hershey vector fonts\nBackend-agnostic design\nText layout system\nAnimation support")
	contentSection.SetStyle(core.TextStyle{
		Font:        font,
		Color:       core.ColorYellow,
		Scale:       1.0,
		LineSpacing: 1.5,
		HAlign:      core.AlignLeft,
		VAlign:      core.AlignTop,
		WordWrap:    true,
	})
	contentSection.SetTitleStyle(core.TextStyle{
		Font:        font,
		Color:       core.ColorOrange,
		Scale:       1.3,
		LineSpacing: 1.2,
		HAlign:      core.AlignCenter,
		VAlign:      core.AlignTop,
		WordWrap:    true,
	})

	// Camera setup
	camera := core.Camera{
		Position:   core.Vec3{X: 0, Y: 100, Z: -300},
		Target:     core.Vec3{X: 0, Y: 0, Z: 100},
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

		// Animate camera
		camera.Position.Y = 100 + float32(math.Sin(float64(totalTime*0.5)))*20
		camera.Position.X = float32(math.Sin(float64(totalTime*0.2))) * 50
		camera.Position.Z = -300 + float32(math.Cos(float64(totalTime*0.2)))*50

		renderer.BeginFrame()
		renderer.Begin3D(camera)

		// Draw text document
		textScreenRenderer.DrawTextDocument(doc)

		// Draw grid
		renderer.DrawGrid(10, 10.0)

		renderer.End3D()

		renderer.DrawFPS(10, 10)
		renderer.DrawText2D("Spectrex Demo - Press ESC to exit", 10, 40, 20, core.ColorWhite)

		renderer.EndFrame()
	}
}
