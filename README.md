# Spectrex

A vector-based UI framework for games, designed for retro vector aesthetics with modern rendering flexibility.

## Features

- **Hershey Vector Fonts** - Stroke-based fonts rendered as line segments
- **Backend Agnostic** - Core logic separated from rendering implementation
- **Text Layout System** - TextScreen, TextRegion, and TextDocument abstractions
- **3D Text Rendering** - Virtual screens positioned in 3D space
- **Animation Support** - Tweening system with easing functions

## Installation

```bash
go get github.com/chazu/spectrex
```

## Quick Start

```go
package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"

    "github.com/chazu/spectrex/backends/raylib"
    "github.com/chazu/spectrex/core"
)

func main() {
    rl.InitWindow(1280, 720, "Spectrex Demo")
    defer rl.CloseWindow()
    rl.SetTargetFPS(60)

    renderer := raylib.NewRenderer(1280, 720)
    fontRenderer := raylib.NewFontRenderer()
    font := core.LoadHersheyFontData()

    camera := core.NewDefaultCamera()

    for !rl.WindowShouldClose() {
        renderer.BeginFrame()
        renderer.Begin3D(camera)

        fontRenderer.DrawText(font, "Hello, Vector World!",
            core.Vec3{X: 0, Y: 0, Z: 0},
            core.ColorGreen, 2.0)

        renderer.End3D()
        renderer.EndFrame()
    }
}
```

## Architecture

```
spectrex/
├── core/           # Backend-agnostic types and logic
├── backends/
│   └── raylib/     # raylib rendering implementation
└── examples/
    └── demo/       # Demo application
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed design documentation and roadmap.

## Available Fonts

Spectrex uses the Hershey font system. Available fonts:

- Simplex (default)
- Complex, ComplexSmall
- Duplex
- Gothic, GothicItalic, Gothic-German, GothicItalic-German, GothicEnglish
- Italic, Italic-Complex
- Script, Script-Complex
- Roman, Roman-Complex

Load alternative fonts with:

```go
font := core.LoadHersheyFontByName("Gothic")
```

## Status

**Alpha** - Extracted from [Tartan](https://github.com/chazu/tartan) prototype. API may change.

## License

MIT
