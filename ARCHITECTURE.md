# Spectrex Architecture & Roadmap

Spectrex is a vector-based UI framework for games, designed for retro vector aesthetics with modern rendering flexibility.

## Origin

Spectrex was extracted from the [Tartan](https://github.com/chazu/tartan) project's `pkg/spex` prototype. The original code demonstrated a working vector text system using Hershey fonts rendered in 3D space via raylib.

## Current State

The codebase has been restructured into a backend-agnostic architecture:

```
spectrex/
├── core/                     # Rendering-agnostic core
│   ├── types.go              # Vec2, Vec3, Color, Matrix
│   ├── font.go               # HersheyFont, Glyph, Stroke
│   ├── layout.go             # TextScreen, TextRegion
│   ├── document.go           # TextDocument, TextSection
│   ├── animation.go          # Animation primitives
│   ├── geometry.go           # Polygon generation
│   └── renderer.go           # Renderer interface definition
│
├── backends/
│   └── raylib/               # raylib implementation
│       ├── renderer.go       # Core renderer impl
│       ├── convert.go        # Type conversions
│       ├── font.go           # Font rendering
│       └── textscreen.go     # TextScreen rendering
│
└── examples/
    └── demo/                 # Demo application
        └── main.go
```

## Architecture Principles

### 1. Backend Agnostic Core

The `core/` package defines all types, interfaces, and logic without importing any rendering library. This enables:

- Multiple rendering backends (raylib, SDL, OpenGL, terminal)
- Easier testing of core logic
- Clean separation of concerns

### 2. Renderer Interface

```go
type Renderer interface {
    BeginFrame()
    EndFrame()
    Begin3D(camera Camera)
    End3D()
    DrawLine3D(start, end Vec3, color Color)
    DrawTriangle3D(v1, v2, v3 Vec3, color Color)
    DrawGrid(slices int, spacing float32)
    DrawFPS(x, y int32)
    DrawText2D(text string, x, y, fontSize int32, color Color)
    GetScreenWidth() int32
    GetScreenHeight() int32
}
```

Backends implement this interface to provide actual rendering.

### 3. Text Layout System

The text system follows a hierarchical model:

```
TextDocument (high-level document with sections)
    └── TextSection (titled content blocks)
            └── TextRegion (positioned text areas)
                    └── HersheyFont (vector glyph data)
                            └── Stroke (line segments)
```

Each level adds abstraction for common UI patterns.

## Refactoring Tasks

### Phase 1: Cleanup & Testing (Current)

- [ ] Add unit tests for `core/` package
  - [ ] types.go - Matrix multiplication, vector operations
  - [ ] font.go - MeasureText, glyph lookup
  - [ ] layout.go - WrapText, TruncateLineToFit
  - [ ] geometry.go - MakePoly, TransformPoly
- [ ] Remove dead code paths
- [ ] Consistent error handling
- [ ] Add godoc comments to all public APIs

### Phase 2: Interface Refinement

- [ ] Review Renderer interface for completeness
- [ ] Add FontRenderer interface to backends
- [ ] Consider adding `TextScreenRenderer` to core interface
- [ ] Evaluate whether Scene belongs in core or backends

### Phase 3: Additional Backends

Potential backends to implement:

- [ ] **Terminal/TUI** - ASCII art rendering using box-drawing characters
- [ ] **SDL2** - Alternative to raylib with different tradeoffs
- [ ] **SVG** - Static image export
- [ ] **Canvas/WASM** - Browser-based rendering

### Phase 4: Feature Expansion

- [ ] **Animation System**
  - Keyframe animations
  - Animation curves (bezier, spring)
  - Property binding

- [ ] **Input Handling**
  - Focus management
  - Hit testing for regions
  - Keyboard navigation

- [ ] **Theming**
  - Color schemes
  - Font variations
  - Style inheritance

- [ ] **Effects**
  - Glow/bloom for vector lines
  - CRT scanline effect
  - Color cycling

### Phase 5: Performance

- [ ] Glyph caching
- [ ] Dirty region tracking
- [ ] Batch rendering
- [ ] Memory pooling for frequent allocations

## Quality Issues to Address

### From Original Codebase

1. **No tests** - Zero test coverage needs immediate attention
2. **Magic numbers** - Font heights, spacing values should be configurable
3. **Package comments** - Several files had incorrect package descriptions
4. **Unused code** - `drawDebugBounds` was never called
5. **Error handling** - Font loading silently fails on missing glyphs

### New Architecture Concerns

1. **Matrix math** - Verify matrix operations match raylib's conventions
2. **Coordinate systems** - Document the expected coordinate orientation
3. **Memory allocation** - Type conversions create allocations; consider pooling

## API Examples

### Basic Text Display

```go
package main

import (
    "github.com/chazu/spectrex/backends/raylib"
    "github.com/chazu/spectrex/core"
)

func main() {
    renderer := raylib.NewRenderer(1280, 720)
    fontRenderer := raylib.NewFontRenderer()
    font := core.LoadHersheyFontData()

    camera := core.NewDefaultCamera()

    for running {
        renderer.BeginFrame()
        renderer.Begin3D(camera)

        fontRenderer.DrawText(font, "Hello, Vector World!",
            core.Vec3{X: 0, Y: 0, Z: 0},
            core.ColorGreen,
            2.0)

        renderer.End3D()
        renderer.EndFrame()
    }
}
```

### Text Layout with Regions

```go
screen := core.NewTextScreen(
    core.Vec3{X: 0, Y: 50, Z: 200},
    800, 600, 1.0,
)

region := screen.AddRegion(20, 20, 360, 200)
region.SetContent("Wrapped text goes here...", font, core.ColorWhite)
region.SetAlignment(core.AlignLeft, core.AlignTop)
region.WordWrap = true

textRenderer := raylib.NewTextScreenRenderer()
textRenderer.DrawTextScreen(screen)
```

### Document with Sections

```go
doc := core.NewTextDocument(screen, 2, 20) // 2 columns, 20px padding
doc.PageStyle.Font = font

section := doc.AddSection("Title", "Content goes here")
section.SetStyle(core.TextStyle{
    Font:   font,
    Color:  core.ColorYellow,
    Scale:  1.5,
    HAlign: core.AlignCenter,
})

textRenderer.DrawTextDocument(doc)
```

## Contributing

When adding new features:

1. Add core types/logic to `core/` without backend dependencies
2. Define interfaces in `core/renderer.go` if needed
3. Implement in `backends/raylib/` (and other backends)
4. Add tests for core logic
5. Update examples to demonstrate usage

## Dependencies

- `github.com/chazu/hershey-go` - Hershey font data
- `github.com/gen2brain/raylib-go/raylib` - raylib backend (optional)

## License

MIT (to be confirmed)
