// Package core provides the Hershey vector font system.
// Hershey fonts are stroke-based fonts represented as a series of line segments,
// making them ideal for vector graphics and 3D rendering.
package core

import (
	"github.com/chazu/hershey-go"
)

// Stroke represents a single line segment in a Hershey font glyph.
// Each character is made up of multiple strokes that together form the shape of the letter.
type Stroke struct {
	From Vec2 // Starting point of the stroke
	To   Vec2 // Ending point of the stroke
}

// HersheyGlyph represents a single character in the Hershey font.
// Each glyph contains metadata about its size and the strokes that form its shape.
type HersheyGlyph struct {
	Width     int      // Visual width of the glyph for rendering
	RealWidth int      // Actual width used for spacing calculations
	Size      int      // Number of strokes in the glyph
	Strokes   []Stroke // Collection of line segments that form the glyph
}

// HersheyFont represents a complete Hershey font with all its glyphs.
// It provides methods for calculating text dimensions and accessing glyph data.
type HersheyFont struct {
	Glyphs   map[int]HersheyGlyph // Map of ASCII values (minus 31) to glyphs
	Height   int                  // Standard height of the font
	FontName string               // Name of the font from hershey-go library
}

// NewHersheyFont creates a new empty Hershey font with default settings.
func NewHersheyFont() *HersheyFont {
	return &HersheyFont{
		Glyphs:   make(map[int]HersheyGlyph),
		Height:   32,
		FontName: "Simplex",
	}
}

// GetGlyph returns the glyph for a character, or nil if not found.
func (hf *HersheyFont) GetGlyph(char rune) *HersheyGlyph {
	glyph, exists := hf.Glyphs[int(char)-31]
	if !exists {
		return nil
	}
	return &glyph
}

// MeasureText calculates the width of a text string at the given scale.
func (hf *HersheyFont) MeasureText(text string, scale float32) float32 {
	totalWidth := float32(0)

	for _, char := range text {
		if char < 32 || char > 126 {
			continue
		}

		glyph, exists := hf.Glyphs[int(char)-31]
		if !exists {
			totalWidth += 8 * scale
			continue
		}

		if glyph.RealWidth > 0 {
			spacing := float32(glyph.RealWidth)
			if spacing < 5 {
				spacing = 5
			}
			totalWidth += spacing * scale
		} else {
			totalWidth += float32(glyph.Width) * scale
		}

		totalWidth += 1.0 * scale // Character spacing
	}

	return totalWidth
}

// loadHersheyGlyph loads a single glyph from the hershey-go package.
func loadHersheyGlyph(fontName string, char rune) HersheyGlyph {
	// Special case for space character
	if char == ' ' {
		return HersheyGlyph{
			Width:     16,
			RealWidth: 16,
			Size:      0,
			Strokes:   []Stroke{},
		}
	}

	var strokes []Stroke
	var vectorX, vectorY []int
	var penPositions []bool

	moveFn := func(s ...interface{}) {
		x := *s[0].(*int)
		y := *s[1].(*int)
		vectorX = append(vectorX, x)
		vectorY = append(vectorY, y)
		penPositions = append(penPositions, false)
	}

	lineFn := func(s ...interface{}) {
		x := *s[0].(*int)
		y := *s[1].(*int)
		vectorX = append(vectorX, x)
		vectorY = append(vectorY, y)
		penPositions = append(penPositions, true)
	}

	minX, _, maxX, _, err := hershey.StringBounds(fontName, 1, 0, 0, string(char))
	if err != nil {
		return HersheyGlyph{Width: 16, RealWidth: 16, Size: 0, Strokes: []Stroke{}}
	}

	width := maxX - minX

	drawX, drawY := 0, 0
	err = hershey.DrawChar(char, fontName, 1, &drawX, &drawY, moveFn, lineFn)
	if err != nil {
		return HersheyGlyph{Width: 16, RealWidth: 16, Size: 0, Strokes: []Stroke{}}
	}

	if len(vectorX) >= 2 {
		for i := 1; i < len(vectorX); i++ {
			if penPositions[i] {
				stroke := Stroke{
					From: Vec2{X: float32(vectorX[i-1]), Y: float32(vectorY[i-1])},
					To:   Vec2{X: float32(vectorX[i]), Y: float32(vectorY[i])},
				}
				strokes = append(strokes, stroke)
			}
		}
	}

	// Create marker for missing glyphs
	if len(strokes) == 0 && char != ' ' && char != '\t' && char != '\n' && char != '\r' {
		size := float32(8)
		center := float32(4)
		strokes = []Stroke{
			{From: Vec2{X: 0, Y: 0}, To: Vec2{X: size, Y: size}},
			{From: Vec2{X: 0, Y: size}, To: Vec2{X: size, Y: 0}},
			{From: Vec2{X: center - 2, Y: center - 2}, To: Vec2{X: center + 2, Y: center - 2}},
			{From: Vec2{X: center + 2, Y: center - 2}, To: Vec2{X: center + 2, Y: center + 2}},
			{From: Vec2{X: center + 2, Y: center + 2}, To: Vec2{X: center - 2, Y: center + 2}},
			{From: Vec2{X: center - 2, Y: center + 2}, To: Vec2{X: center - 2, Y: center - 2}},
		}
	}

	if width <= 0 {
		width = 16
	}

	return HersheyGlyph{
		Width:     width,
		RealWidth: drawX,
		Size:      len(strokes),
		Strokes:   strokes,
	}
}

// LoadHersheyFontData loads the complete Hershey font data from the hershey-go package.
func LoadHersheyFontData() *HersheyFont {
	font := NewHersheyFont()
	font.FontName = "Simplex"

	for i := 32; i < 127; i++ {
		glyph := loadHersheyGlyph(font.FontName, rune(i))
		font.Glyphs[i-31] = glyph
	}

	return font
}

// LoadHersheyFontByName loads a Hershey font by name.
// Available fonts: Simplex, Complex, ComplexSmall, Duplex, Gothic,
// GothicItalic, Gothic-German, GothicItalic-German, GothicEnglish,
// Italic, Italic-Complex, Script, Script-Complex, Roman, Roman-Complex
func LoadHersheyFontByName(fontName string) *HersheyFont {
	font := NewHersheyFont()
	font.FontName = fontName

	for i := 32; i < 127; i++ {
		glyph := loadHersheyGlyph(fontName, rune(i))
		font.Glyphs[i-31] = glyph
	}

	return font
}
