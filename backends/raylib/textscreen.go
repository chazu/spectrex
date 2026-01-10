// Package raylib provides text screen rendering for the Spectrex framework.
package raylib

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/spectrex/core"
)

// TextScreenRenderer implements core.TextScreenRenderer using raylib.
type TextScreenRenderer struct {
	fontRenderer *FontRenderer
}

// NewTextScreenRenderer creates a new raylib text screen renderer.
func NewTextScreenRenderer() *TextScreenRenderer {
	return &TextScreenRenderer{
		fontRenderer: NewFontRenderer(),
	}
}

// DrawTextScreen renders a complete text screen with all its regions.
func (tsr *TextScreenRenderer) DrawTextScreen(screen *core.TextScreen) {
	model := tsr.calculateTransform(screen)

	// Draw background if not transparent
	if !screen.Transparent {
		tsr.drawScreenBackground(screen, model)
	}

	// Draw border if enabled
	if screen.ShowBorder || screen.Debug {
		tsr.drawScreenBorder(screen, model)
	}

	// Draw all regions
	for _, region := range screen.Regions {
		tsr.DrawTextRegion(region, model, screen.Scale)
	}
}

// DrawTextRegion renders a single text region.
func (tsr *TextScreenRenderer) DrawTextRegion(region *core.TextRegion, screenTransform rl.Matrix, screenScale float32) {
	// Draw background if not transparent
	if !region.Transparent {
		tsr.drawRegionBackground(region, screenTransform)
	}

	// Draw border if enabled
	if region.ShowBorder || region.Parent.Debug {
		tsr.drawRegionBorder(region, screenTransform)
	}

	// Skip text rendering if no text or font
	if region.Font == nil || region.Text == "" {
		return
	}

	effectiveScale := region.Scale * screenScale
	lines := region.GetLines()
	if len(lines) == 0 {
		return
	}

	totalTextHeight := region.CalculateTextHeight(lines)
	startY := region.CalculateStartY(totalTextHeight)
	lineHeight := float32(region.Font.Height) * effectiveScale

	for i, line := range lines {
		// Subtract Y because in 3D space Y increases upward, but text flows downward
		yPos := startY - float32(i)*lineHeight*region.LineSpacing

		if yPos < region.Y || yPos > region.Y+region.Height {
			continue
		}

		// Handle truncation
		if region.TruncateOverflow && region.HAlign != core.AlignJustified {
			lineWidth := region.CalculateLineWidth(line, effectiveScale)
			if lineWidth > region.Width {
				if !strings.HasSuffix(line, region.OverflowMarker) {
					markerWidth := region.CalculateLineWidth(region.OverflowMarker, effectiveScale)
					line = region.TruncateLineToFit(line, region.Width-markerWidth, effectiveScale) + region.OverflowMarker
				}
			}
		}

		// Calculate X position based on alignment
		// Note: 180° Y rotation flips X axis, so local right → world left
		lineWidth := region.CalculateLineWidth(line, effectiveScale)
		var xPos float32
		switch region.HAlign {
		case core.AlignLeft:
			// Start from local right edge (becomes world left after transform)
			xPos = region.X + region.Width
		case core.AlignCenter:
			// Center the text
			xPos = region.X + (region.Width+lineWidth)/2
		case core.AlignRight:
			// End at local left edge (becomes world right after transform)
			xPos = region.X + lineWidth
		case core.AlignJustified:
			if i < len(lines)-1 && strings.Contains(line, " ") {
				tsr.drawJustifiedLine(region, line, region.X+region.Width, yPos, effectiveScale, screenTransform)
				continue
			}
			xPos = region.X + region.Width
		}

		pos := rl.Vector3{X: xPos, Y: yPos, Z: 0}
		transformedPos := rl.Vector3Transform(pos, screenTransform)

		tsr.drawLine(region, line, transformedPos, effectiveScale)
	}
}

// DrawTextDocument renders a complete text document.
func (tsr *TextScreenRenderer) DrawTextDocument(doc *core.TextDocument) {
	if len(doc.Sections) > 0 && doc.Sections[0].Region == nil {
		doc.Layout()
	}

	for _, section := range doc.Sections {
		tsr.drawSection(section)
	}
}

func (tsr *TextScreenRenderer) calculateTransform(screen *core.TextScreen) rl.Matrix {
	model := rl.MatrixIdentity()
	model = rl.MatrixRotateX(core.DegToRad(screen.Rotation.X))
	model = rl.MatrixMultiply(model, rl.MatrixRotateY(core.DegToRad(screen.Rotation.Y+180.0)))
	model = rl.MatrixMultiply(model, rl.MatrixRotateZ(core.DegToRad(screen.Rotation.Z)))
	model = rl.MatrixMultiply(model, rl.MatrixTranslate(screen.Position.X, screen.Position.Y, screen.Position.Z))
	return model
}

func (tsr *TextScreenRenderer) drawScreenBackground(screen *core.TextScreen, transform rl.Matrix) {
	topLeft := rl.Vector3Transform(rl.Vector3{X: 0, Y: 0, Z: 0}, transform)
	topRight := rl.Vector3Transform(rl.Vector3{X: screen.Width, Y: 0, Z: 0}, transform)
	bottomRight := rl.Vector3Transform(rl.Vector3{X: screen.Width, Y: screen.Height, Z: 0}, transform)
	bottomLeft := rl.Vector3Transform(rl.Vector3{X: 0, Y: screen.Height, Z: 0}, transform)

	bgColor := coreToRlColor(screen.BackgroundColor)
	rl.DrawTriangle3D(topLeft, topRight, bottomRight, bgColor)
	rl.DrawTriangle3D(topLeft, bottomRight, bottomLeft, bgColor)
}

func (tsr *TextScreenRenderer) drawScreenBorder(screen *core.TextScreen, transform rl.Matrix) {
	topLeft := rl.Vector3Transform(rl.Vector3{X: 0, Y: 0, Z: 0}, transform)
	topRight := rl.Vector3Transform(rl.Vector3{X: screen.Width, Y: 0, Z: 0}, transform)
	bottomRight := rl.Vector3Transform(rl.Vector3{X: screen.Width, Y: screen.Height, Z: 0}, transform)
	bottomLeft := rl.Vector3Transform(rl.Vector3{X: 0, Y: screen.Height, Z: 0}, transform)

	borderColor := coreToRlColor(screen.BorderColor)
	if screen.Debug {
		borderColor = rl.Blue
	}

	rl.DrawLine3D(topLeft, topRight, borderColor)
	rl.DrawLine3D(topRight, bottomRight, borderColor)
	rl.DrawLine3D(bottomRight, bottomLeft, borderColor)
	rl.DrawLine3D(bottomLeft, topLeft, borderColor)
}

func (tsr *TextScreenRenderer) drawRegionBackground(region *core.TextRegion, transform rl.Matrix) {
	topLeft := rl.Vector3Transform(rl.Vector3{X: region.X, Y: region.Y, Z: 0}, transform)
	topRight := rl.Vector3Transform(rl.Vector3{X: region.X + region.Width, Y: region.Y, Z: 0}, transform)
	bottomRight := rl.Vector3Transform(rl.Vector3{X: region.X + region.Width, Y: region.Y + region.Height, Z: 0}, transform)
	bottomLeft := rl.Vector3Transform(rl.Vector3{X: region.X, Y: region.Y + region.Height, Z: 0}, transform)

	bgColor := coreToRlColor(region.BackgroundColor)
	rl.DrawTriangle3D(topLeft, topRight, bottomRight, bgColor)
	rl.DrawTriangle3D(topLeft, bottomRight, bottomLeft, bgColor)
}

func (tsr *TextScreenRenderer) drawRegionBorder(region *core.TextRegion, transform rl.Matrix) {
	topLeft := rl.Vector3Transform(rl.Vector3{X: region.X, Y: region.Y, Z: 0}, transform)
	topRight := rl.Vector3Transform(rl.Vector3{X: region.X + region.Width, Y: region.Y, Z: 0}, transform)
	bottomRight := rl.Vector3Transform(rl.Vector3{X: region.X + region.Width, Y: region.Y + region.Height, Z: 0}, transform)
	bottomLeft := rl.Vector3Transform(rl.Vector3{X: region.X, Y: region.Y + region.Height, Z: 0}, transform)

	borderColor := coreToRlColor(region.BorderColor)
	if region.Parent.Debug {
		borderColor = rl.Red
	}

	rl.DrawLine3D(topLeft, topRight, borderColor)
	rl.DrawLine3D(topRight, bottomRight, borderColor)
	rl.DrawLine3D(bottomRight, bottomLeft, borderColor)
	rl.DrawLine3D(bottomLeft, topLeft, borderColor)
}

func (tsr *TextScreenRenderer) drawLine(region *core.TextRegion, line string, position rl.Vector3, scale float32) {
	xOffset := float32(0)
	runes := []rune(line)

	// Iterate backwards through characters to compensate for 180° Y rotation mirror effect
	for i := len(runes) - 1; i >= 0; i-- {
		char := runes[i]

		if char < 32 || char > 126 {
			continue
		}

		glyph, exists := region.Font.Glyphs[int(char)-31]
		if !exists {
			xOffset += 8 * scale
			continue
		}

		var glyphWidth float32
		if glyph.RealWidth > 0 {
			spacing := float32(glyph.RealWidth)
			if spacing < 5 {
				spacing = 5
			}
			glyphWidth = spacing * scale
		} else {
			glyphWidth = float32(glyph.Width) * scale
		}

		glyphPos := rl.Vector3{
			X: position.X + xOffset + glyphWidth,
			Y: position.Y,
			Z: position.Z,
		}

		tsr.drawGlyph(region.Font, int(char), glyphPos, region.Color, scale)

		xOffset += glyphWidth
		xOffset += (1.0 + region.CharSpacing) * scale
	}
}

func (tsr *TextScreenRenderer) drawGlyph(font *core.HersheyFont, char int, position rl.Vector3, color core.Color, scale float32) {
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
		rl.DrawLine3D(start, end, rlColor)
	}
}

func (tsr *TextScreenRenderer) drawJustifiedLine(region *core.TextRegion, line string, x, y float32, scale float32, transform rl.Matrix) {
	words := strings.Split(line, " ")
	if len(words) <= 1 {
		pos := rl.Vector3Transform(rl.Vector3{X: x, Y: y, Z: 0}, transform)
		tsr.drawLine(region, line, pos, scale)
		return
	}

	totalWordsWidth := float32(0)
	for _, word := range words {
		totalWordsWidth += region.CalculateLineWidth(word, scale)
	}

	extraSpacePerGap := (region.Width - totalWordsWidth) / float32(len(words)-1)
	// Start from local right edge (x is already region.X + region.Width from caller)
	// and work leftward, placing words from last to first
	xPos := x - region.Width

	for i := len(words) - 1; i >= 0; i-- {
		word := words[i]
		wordWidth := region.CalculateLineWidth(word, scale)
		wordPos := xPos + wordWidth

		pos := rl.Vector3Transform(rl.Vector3{X: wordPos, Y: y, Z: 0}, transform)
		tsr.drawLine(region, word, pos, scale)

		xPos += wordWidth
		if i > 0 {
			xPos += extraSpacePerGap
		}
	}
}

func (tsr *TextScreenRenderer) drawSection(section *core.TextSection) {
	if section.Region == nil {
		return
	}

	contentFont := section.GetContentFont()
	titleFont := section.GetTitleFont()

	region := section.Region
	screenTransform := tsr.calculateTransform(region.Parent)

	if section.Title != "" && titleFont != nil {
		titleLines := float32(len(strings.Split(section.Title, "\n")))
		titleHeight := titleLines * float32(titleFont.Height) *
			section.TitleStyle.Scale * section.TitleStyle.LineSpacing
		titleGap := float32(titleFont.Height) * section.Style.Scale * 0.5

		// Title region at top of section (higher Y in 3D space)
		titleRegion := &core.TextRegion{
			X:           region.X,
			Y:           region.Y + region.Height - titleHeight,
			Width:       region.Width,
			Height:      titleHeight,
			Text:        section.Title,
			Font:        titleFont,
			Color:       section.TitleStyle.Color,
			Scale:       section.TitleStyle.Scale,
			LineSpacing: section.TitleStyle.LineSpacing,
			CharSpacing: section.TitleStyle.CharSpacing,
			HAlign:      section.TitleStyle.HAlign,
			VAlign:      section.TitleStyle.VAlign,
			WordWrap:    true,
			Parent:      region.Parent,
		}

		tsr.DrawTextRegion(titleRegion, screenTransform, region.Parent.Scale)

		// Content region below title (lower Y in 3D space)
		contentRegion := &core.TextRegion{
			X:           region.X,
			Y:           region.Y,
			Width:       region.Width,
			Height:      region.Height - titleHeight - titleGap,
			Text:        section.Content,
			Font:        contentFont,
			Color:       section.Style.Color,
			Scale:       section.Style.Scale,
			LineSpacing: section.Style.LineSpacing,
			CharSpacing: section.Style.CharSpacing,
			HAlign:      section.Style.HAlign,
			VAlign:      section.Style.VAlign,
			WordWrap:    true,
			Parent:      region.Parent,
		}

		tsr.DrawTextRegion(contentRegion, screenTransform, region.Parent.Scale)
	} else {
		region.Text = section.Content
		region.Font = contentFont
		region.Color = section.Style.Color
	}
}
