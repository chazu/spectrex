// Package core defines the Renderer interface for the Spectrex framework.
// Implementations of this interface provide the actual drawing capabilities
// for different backends (raylib, SDL, OpenGL, terminal, etc.).
package core

// Camera represents a 3D camera for scene rendering.
type Camera struct {
	Position   Vec3
	Target     Vec3
	Up         Vec3
	Fovy       float32
	Projection int // 0 = perspective, 1 = orthographic
}

// NewDefaultCamera creates a camera with sensible defaults.
func NewDefaultCamera() Camera {
	return Camera{
		Position:   Vec3{X: 0, Y: 100, Z: -300},
		Target:     Vec3{X: 0, Y: 0, Z: 100},
		Up:         Vec3{X: 0, Y: 1, Z: 0},
		Fovy:       45.0,
		Projection: 0, // Perspective
	}
}

// Renderer defines the interface for rendering backends.
// Implementations must provide these drawing primitives.
type Renderer interface {
	// Frame management
	BeginFrame()
	EndFrame()

	// 3D mode management
	Begin3D(camera Camera)
	End3D()

	// Primitive drawing
	DrawLine3D(start, end Vec3, color Color)
	DrawTriangle3D(v1, v2, v3 Vec3, color Color)

	// Utility drawing
	DrawGrid(slices int, spacing float32)
	DrawFPS(x, y int32)
	DrawText2D(text string, x, y int32, fontSize int32, color Color)

	// Screen info
	GetScreenWidth() int32
	GetScreenHeight() int32
}

// FontRenderer defines the interface for rendering vector fonts.
type FontRenderer interface {
	// DrawGlyph draws a single glyph at the specified position.
	DrawGlyph(font *HersheyFont, char int, position Vec3, color Color, scale float32)

	// DrawText draws a complete text string centered at the position.
	DrawText(font *HersheyFont, text string, position Vec3, color Color, scale float32)
}

// TextScreenRenderer defines the interface for rendering text screens.
type TextScreenRenderer interface {
	// DrawTextScreen renders a complete text screen with all its regions.
	DrawTextScreen(screen *TextScreen)

	// DrawTextRegion renders a single text region.
	DrawTextRegion(region *TextRegion, transform Matrix, scale float32)

	// DrawTextDocument renders a complete text document.
	DrawTextDocument(doc *TextDocument)
}

// Object is an interface for all renderable objects in a scene.
type Object interface {
	Update(deltaTime float32)
	Draw(renderer Renderer)
}

// Scene represents a collection of objects to be rendered.
type Scene struct {
	Camera          Camera
	Objects         []Object
	BackgroundColor Color
}

// NewScene creates a new scene with a default camera.
func NewScene() *Scene {
	return &Scene{
		Camera:          NewDefaultCamera(),
		Objects:         make([]Object, 0),
		BackgroundColor: ColorBlack,
	}
}

// AddObject adds an object to the scene.
func (s *Scene) AddObject(obj Object) {
	s.Objects = append(s.Objects, obj)
}

// Update updates all objects in the scene.
func (s *Scene) Update(deltaTime float32) {
	for _, obj := range s.Objects {
		obj.Update(deltaTime)
	}
}

// Draw renders all objects in the scene.
func (s *Scene) Draw(renderer Renderer) {
	for _, obj := range s.Objects {
		obj.Draw(renderer)
	}
}
