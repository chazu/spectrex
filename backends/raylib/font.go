// Package raylib provides font rendering for the Spectrex framework using raylib.
package raylib

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/spectrex/core"
)

// FontRenderer implements core.FontRenderer using raylib.
type FontRenderer struct{}

// NewFontRenderer creates a new raylib font renderer.
func NewFontRenderer() *FontRenderer {
	return &FontRenderer{}
}

// DrawGlyph draws a single glyph at the specified position.
func (fr *FontRenderer) DrawGlyph(font *core.HersheyFont, char int, position core.Vec3, color core.Color, scale float32) {
	glyph, exists := font.Glyphs[char-31]
	if !exists || len(glyph.Strokes) == 0 {
		return
	}

	rlColor := coreToRlColor(color)
	rlPos := coreToRlVec3(position)

	for _, stroke := range glyph.Strokes {
		start := rl.Vector3{
			X: rlPos.X - stroke.From.X*scale,
			Y: rlPos.Y + stroke.From.Y*scale,
			Z: rlPos.Z,
		}
		end := rl.Vector3{
			X: rlPos.X - stroke.To.X*scale,
			Y: rlPos.Y + stroke.To.Y*scale,
			Z: rlPos.Z,
		}
		rl.DrawLine3D(start, end, rlColor)
	}
}

// DrawText draws a complete text string centered at the position.
func (fr *FontRenderer) DrawText(font *core.HersheyFont, text string, position core.Vec3, color core.Color, scale float32) {
	totalWidth := font.MeasureText(text, scale)
	startX := totalWidth / 2.0
	xOffset := float32(0)

	for _, char := range text {
		if char < 32 || char > 126 {
			continue
		}

		glyph, exists := font.Glyphs[int(char)-31]
		if !exists {
			xOffset += 8 * scale
			continue
		}

		glyphPos := core.Vec3{
			X: position.X + startX - xOffset,
			Y: position.Y,
			Z: position.Z,
		}

		fr.DrawGlyph(font, int(char), glyphPos, color, scale)

		if glyph.RealWidth > 0 {
			spacing := float32(glyph.RealWidth)
			if spacing < 5 {
				spacing = 5
			}
			xOffset += spacing * scale
		} else {
			xOffset += float32(glyph.Width) * scale
		}

		xOffset += 1.0 * scale
	}
}

// DrawGlyphTransformed draws a glyph with a transformation matrix applied.
func (fr *FontRenderer) DrawGlyphTransformed(font *core.HersheyFont, char int, position core.Vec3, color core.Color, scale float32, transform rl.Matrix) {
	glyph, exists := font.Glyphs[char-31]
	if !exists || len(glyph.Strokes) == 0 {
		return
	}

	rlColor := coreToRlColor(color)

	for _, stroke := range glyph.Strokes {
		start := rl.Vector3{
			X: position.X - stroke.From.X*scale,
			Y: position.Y + stroke.From.Y*scale,
			Z: position.Z,
		}
		end := rl.Vector3{
			X: position.X - stroke.To.X*scale,
			Y: position.Y + stroke.To.Y*scale,
			Z: position.Z,
		}

		// Transform positions
		start = rl.Vector3Transform(start, transform)
		end = rl.Vector3Transform(end, transform)

		rl.DrawLine3D(start, end, rlColor)
	}
}
