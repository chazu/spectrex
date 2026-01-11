// Package core provides text layout capabilities for the Spectrex framework.
// This file implements the text layout system for arranging and displaying text
// in virtual 2D screens within 3D space.
package core

import (
	"strings"
)

// TextAlign defines text alignment options within a text region.
type TextAlign int

const (
	AlignLeft TextAlign = iota
	AlignCenter
	AlignRight
	AlignJustified
)

// VerticalAlign defines vertical alignment options within a text region.
type VerticalAlign int

const (
	AlignTop VerticalAlign = iota
	AlignMiddle
	AlignBottom
)

// TextScreen represents a virtual 2D screen in 3D space for organizing text and regions.
type TextScreen struct {
	Position        Vec3
	Rotation        Vec3
	Width           float32
	Height          float32
	Scale           float32
	Regions         []*TextRegion
	Transparent     bool
	ShowBorder      bool
	BorderColor     Color
	BackgroundColor Color
	Debug           bool
}

// TextRegion represents a rectangular area within a TextScreen for text layout.
type TextRegion struct {
	X               float32
	Y               float32
	Width           float32
	Height          float32
	Text            string
	Font            *HersheyFont
	Color           Color
	Scale           float32
	LineSpacing     float32
	CharSpacing     float32
	HAlign          TextAlign
	VAlign          VerticalAlign
	WordWrap        bool
	MaxLines        int
	TruncateOverflow bool
	OverflowMarker  string
	Transparent     bool
	ShowBorder      bool
	BorderColor     Color
	BackgroundColor Color
	Parent          *TextScreen
}

// NewTextScreen creates a new virtual screen for text layout in 3D space.
func NewTextScreen(position Vec3, width, height, scale float32) *TextScreen {
	return &TextScreen{
		Position:        position,
		Width:           width,
		Height:          height,
		Scale:           scale,
		Regions:         make([]*TextRegion, 0),
		Debug:           false,
		Transparent:     true,
		ShowBorder:      false,
		BorderColor:     ColorWhite,
		BackgroundColor: ColorBlack,
	}
}

// AddRegion creates a new text region within the screen and returns it.
func (ts *TextScreen) AddRegion(x, y, width, height float32) *TextRegion {
	region := &TextRegion{
		X:                x,
		Y:                y,
		Width:            width,
		Height:           height,
		Scale:            1.0,
		LineSpacing:      1.2,
		CharSpacing:      0.0,
		HAlign:           AlignLeft,
		VAlign:           AlignTop,
		WordWrap:         true,
		TruncateOverflow: false,
		OverflowMarker:   "...",
		Transparent:      true,
		ShowBorder:       false,
		BorderColor:      ColorWhite,
		BackgroundColor:  ColorBlack,
		Parent:           ts,
	}
	ts.Regions = append(ts.Regions, region)
	return region
}

// SetTransparency sets whether the screen should be transparent.
func (ts *TextScreen) SetTransparency(transparent bool) {
	ts.Transparent = transparent
}

// SetBorder sets whether to show the screen border and its color.
func (ts *TextScreen) SetBorder(show bool, color Color) {
	ts.ShowBorder = show
	ts.BorderColor = color
}

// SetBackground sets the background color for the screen when not transparent.
func (ts *TextScreen) SetBackground(color Color) {
	ts.BackgroundColor = color
}

// SetDebug enables or disables debug visualization of regions.
func (ts *TextScreen) SetDebug(debug bool) {
	ts.Debug = debug
}

// GetTransformMatrix calculates the screen's transformation matrix.
func (ts *TextScreen) GetTransformMatrix() Matrix {
	model := MatrixIdentity()
	model = MatrixRotateX(DegToRad(ts.Rotation.X))
	model = model.Multiply(MatrixRotateY(DegToRad(ts.Rotation.Y + 180.0)))
	model = model.Multiply(MatrixRotateZ(DegToRad(ts.Rotation.Z)))
	model = model.Multiply(MatrixTranslate(ts.Position.X, ts.Position.Y, ts.Position.Z))
	return model
}

// SetContent sets the content and styling for a text region.
func (tr *TextRegion) SetContent(text string, font *HersheyFont, color Color) {
	tr.Text = text
	tr.Font = font
	tr.Color = color
}

// SetAlignment sets the horizontal and vertical alignment for a text region.
func (tr *TextRegion) SetAlignment(hAlign TextAlign, vAlign VerticalAlign) {
	tr.HAlign = hAlign
	tr.VAlign = vAlign
}

// SetSpacing sets the line and character spacing for a text region.
func (tr *TextRegion) SetSpacing(lineSpacing, charSpacing float32) {
	tr.LineSpacing = lineSpacing
	tr.CharSpacing = charSpacing
}

// SetTransparency sets whether the region should be transparent.
func (tr *TextRegion) SetTransparency(transparent bool) {
	tr.Transparent = transparent
}

// SetBorder sets whether to show the region border and its color.
func (tr *TextRegion) SetBorder(show bool, color Color) {
	tr.ShowBorder = show
	tr.BorderColor = color
}

// SetBackground sets the background color for the region when not transparent.
func (tr *TextRegion) SetBackground(color Color) {
	tr.BackgroundColor = color
}

// SetOverflowHandling configures how text that doesn't fit in the region is handled.
func (tr *TextRegion) SetOverflowHandling(truncate bool, marker string) {
	tr.TruncateOverflow = truncate
	tr.OverflowMarker = marker
}

// CalculateLineWidth calculates the rendered width of a text line.
func (tr *TextRegion) CalculateLineWidth(line string, scale float32) float32 {
	if tr.Font == nil {
		return 0
	}

	totalWidth := float32(0)

	for _, char := range line {
		if char < 32 || char > 126 {
			continue
		}

		glyph, exists := tr.Font.Glyphs[int(char)-31]
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

		totalWidth += (1.0 + tr.CharSpacing) * scale
	}

	return totalWidth
}

// WrapText wraps the text to fit within the region width.
func (tr *TextRegion) WrapText() []string {
	effectiveScale := tr.Scale * tr.Parent.Scale

	rawLines := strings.Split(tr.Text, "\n")
	var wrappedLines []string

	for _, line := range rawLines {
		if line == "" {
			wrappedLines = append(wrappedLines, "")
			continue
		}

		words := strings.Split(line, " ")
		currentLine := ""
		currentWidth := float32(0)

		for _, word := range words {
			wordWidth := tr.CalculateLineWidth(word, effectiveScale)
			spaceWidth := tr.CalculateLineWidth(" ", effectiveScale)

			if currentWidth > 0 && currentWidth+wordWidth+spaceWidth > tr.Width {
				wrappedLines = append(wrappedLines, currentLine)
				currentLine = word
				currentWidth = wordWidth
			} else {
				if currentWidth > 0 {
					currentLine += " " + word
					currentWidth += spaceWidth + wordWidth
				} else {
					currentLine = word
					currentWidth = wordWidth
				}
			}
		}

		if currentLine != "" {
			wrappedLines = append(wrappedLines, currentLine)
		}
	}

	return wrappedLines
}

// TruncateLineToFit truncates a line of text to fit within a specified width.
func (tr *TextRegion) TruncateLineToFit(line string, maxWidth float32, scale float32) string {
	if tr.CalculateLineWidth(line, scale) <= maxWidth {
		return line
	}

	runes := []rune(line)
	left := 1
	right := len(runes)
	result := 0

	for left <= right {
		mid := (left + right) / 2
		substring := string(runes[:mid])
		width := tr.CalculateLineWidth(substring, scale)

		if width <= maxWidth {
			result = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return string(runes[:result])
}

// GetLines returns the processed lines ready for rendering.
func (tr *TextRegion) GetLines() []string {
	if tr.Font == nil || tr.Text == "" {
		return nil
	}

	effectiveScale := tr.Scale * tr.Parent.Scale

	var lines []string
	if tr.WordWrap {
		lines = tr.WrapText()
	} else {
		lines = strings.Split(tr.Text, "\n")
	}

	if tr.MaxLines > 0 && len(lines) > tr.MaxLines {
		if tr.TruncateOverflow && tr.OverflowMarker != "" {
			lastVisibleLine := lines[tr.MaxLines-1]
			maxLineWidth := tr.CalculateLineWidth(lastVisibleLine, effectiveScale)
			markerWidth := tr.CalculateLineWidth(tr.OverflowMarker, effectiveScale)

			if maxLineWidth+markerWidth > tr.Width {
				truncatedLine := tr.TruncateLineToFit(lastVisibleLine, tr.Width-markerWidth, effectiveScale)
				lines[tr.MaxLines-1] = truncatedLine + tr.OverflowMarker
			} else {
				lines[tr.MaxLines-1] = lastVisibleLine + tr.OverflowMarker
			}
		}
		lines = lines[:tr.MaxLines]
	}

	return lines
}

// CalculateTextHeight calculates the total height of the text block.
func (tr *TextRegion) CalculateTextHeight(lines []string) float32 {
	if tr.Font == nil || len(lines) == 0 {
		return 0
	}

	effectiveScale := tr.Scale * tr.Parent.Scale
	lineHeight := float32(tr.Font.Height) * effectiveScale
	totalHeight := lineHeight * float32(len(lines))

	if len(lines) > 1 {
		totalHeight += (float32(len(lines)-1) * lineHeight * (tr.LineSpacing - 1.0))
	}

	return totalHeight
}

// CalculateStartY calculates the starting Y position based on vertical alignment.
// Note: In 3D space Y increases upward, so "top" of region is at tr.Y + tr.Height.
// Text lines are rendered with decreasing Y (flowing downward on screen).
func (tr *TextRegion) CalculateStartY(totalTextHeight float32) float32 {
	// Offset to account for glyph ascent (characters extend above baseline)
	// Hershey fonts have ascent roughly 70% of the total height
	topPadding := float32(0)
	if tr.Font != nil && tr.Parent != nil {
		effectiveScale := tr.Scale * tr.Parent.Scale
		topPadding = float32(tr.Font.Height) * effectiveScale * 0.8
	}

	switch tr.VAlign {
	case AlignTop:
		// Start below top edge to account for glyph ascent
		return tr.Y + tr.Height - topPadding
	case AlignMiddle:
		// Center vertically
		return tr.Y + tr.Height - (tr.Height-totalTextHeight)/2
	case AlignBottom:
		// Start so last line ends at bottom of region
		return tr.Y + totalTextHeight
	default:
		return tr.Y + tr.Height - topPadding
	}
}
