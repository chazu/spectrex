// Package core provides text composition capabilities for the Spectrex framework.
// This file implements text composition features for creating complex text layouts.
package core

import (
	"strings"
)

// TextStyle defines a set of styling properties for text rendering.
type TextStyle struct {
	Font        *HersheyFont
	Color       Color
	Scale       float32
	LineSpacing float32
	CharSpacing float32
	HAlign      TextAlign
	VAlign      VerticalAlign
	WordWrap    bool
}

// TextDocument represents a complex text document with multiple regions
// and sections, providing high-level text layout capabilities.
type TextDocument struct {
	Screen    *TextScreen
	Sections  []*TextSection
	Columns   int
	Padding   float32
	PageStyle TextStyle
}

// TextSection represents a section of content within a document.
type TextSection struct {
	Title      string
	Content    string
	Style      TextStyle
	TitleStyle TextStyle
	Region     *TextRegion
	Document   *TextDocument
}

// NewTextDocument creates a new text document with the specified screen.
func NewTextDocument(screen *TextScreen, columns int, padding float32) *TextDocument {
	return &TextDocument{
		Screen:   screen,
		Sections: make([]*TextSection, 0),
		Columns:  columns,
		Padding:  padding,
		PageStyle: TextStyle{
			Scale:       1.0,
			LineSpacing: 1.2,
			CharSpacing: 0.0,
			HAlign:      AlignLeft,
			VAlign:      AlignTop,
			WordWrap:    true,
		},
	}
}

// AddSection adds a new section to the document and returns it.
func (doc *TextDocument) AddSection(title, content string) *TextSection {
	section := &TextSection{
		Title:      title,
		Content:    content,
		Style:      doc.PageStyle,
		TitleStyle: doc.PageStyle,
		Document:   doc,
	}

	section.TitleStyle.Scale = doc.PageStyle.Scale * 1.2

	doc.Sections = append(doc.Sections, section)
	return section
}

// Layout calculates the layout for all sections in the document.
// Uses Y-up coordinate system: higher Y values appear higher on screen.
func (doc *TextDocument) Layout() {
	if doc.Screen == nil || len(doc.Sections) == 0 {
		return
	}

	contentWidth := doc.Screen.Width - (2 * doc.Padding)

	columnCount := doc.Columns
	if columnCount < 1 {
		columnCount = 1
	}
	columnWidth := contentWidth / float32(columnCount)

	currentColumn := 0
	// Start from top of screen (high Y) and work down
	currentY := doc.Screen.Height - doc.Padding

	for _, section := range doc.Sections {
		contentLinesCount := len(strings.Split(section.Content, "\n"))
		contentLines := float32(contentLinesCount)

		if section.Style.WordWrap {
			avgCharsPerLine := columnWidth / (section.Style.Scale * 8)
			totalChars := float32(len(section.Content))
			estimatedLines := int(totalChars / avgCharsPerLine)
			if estimatedLines > contentLinesCount {
				contentLines = float32(estimatedLines)
			}
		}

		titleHeight := float32(0)
		if section.Title != "" && doc.PageStyle.Font != nil {
			titleLines := float32(len(strings.Split(section.Title, "\n")))
			titleHeight = titleLines * float32(doc.PageStyle.Font.Height) *
				section.TitleStyle.Scale * section.TitleStyle.LineSpacing
			titleHeight += float32(doc.PageStyle.Font.Height) * section.Style.Scale * 0.5
		}

		sectionHeight := titleHeight
		if doc.PageStyle.Font != nil {
			sectionHeight += contentLines * float32(doc.PageStyle.Font.Height) *
				section.Style.Scale * section.Style.LineSpacing
		}

		// Check if we need to move to next column (Y going below padding)
		if currentY-sectionHeight < doc.Padding {
			currentColumn++
			currentY = doc.Screen.Height - doc.Padding

			if currentColumn >= columnCount {
				currentColumn = columnCount - 1
			}
		}

		x := doc.Padding + float32(currentColumn)*columnWidth
		// Region Y is at bottom of section, height extends upward
		y := currentY - sectionHeight

		region := doc.Screen.AddRegion(x, y, columnWidth, sectionHeight)

		region.Scale = section.Style.Scale
		region.LineSpacing = section.Style.LineSpacing
		region.CharSpacing = section.Style.CharSpacing
		region.HAlign = section.Style.HAlign
		region.VAlign = section.Style.VAlign
		region.WordWrap = section.Style.WordWrap

		section.Region = region

		// Move down for next section (decrease Y)
		if doc.PageStyle.Font != nil {
			currentY -= sectionHeight + float32(doc.PageStyle.Font.Height)*section.Style.Scale
		} else {
			currentY -= sectionHeight + 20 // Default spacing
		}
	}
}

// SetStyle sets the style for a section.
func (section *TextSection) SetStyle(style TextStyle) {
	section.Style = style
}

// SetTitleStyle sets the style for a section's title.
func (section *TextSection) SetTitleStyle(style TextStyle) {
	section.TitleStyle = style
}

// GetContentFont returns the content font, falling back to document default.
func (section *TextSection) GetContentFont() *HersheyFont {
	if section.Style.Font != nil {
		return section.Style.Font
	}
	return section.Document.PageStyle.Font
}

// GetTitleFont returns the title font, falling back to document default.
func (section *TextSection) GetTitleFont() *HersheyFont {
	if section.TitleStyle.Font != nil {
		return section.TitleStyle.Font
	}
	return section.Document.PageStyle.Font
}
