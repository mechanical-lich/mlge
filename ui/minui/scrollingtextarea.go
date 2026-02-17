package minui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text"
)

// ScrollingTextArea is a multi-line text display with scrolling support
type ScrollingTextArea struct {
	*ElementBase
	Lines        []string
	ScrollOffset int
	LineHeight   int
	VisibleLines int
	WrapWidth    int

	draggingThumb bool
	dragOffsetY   int

	// Cached background for sprite-based rendering
	bgCache       *ebiten.Image
	bgCacheWidth  int
	bgCacheHeight int
}

// NewScrollingTextArea creates a new scrolling text area
func NewScrollingTextArea(id string, width, height int) *ScrollingTextArea {
	sta := &ScrollingTextArea{
		ElementBase:  NewElementBase(id),
		Lines:        make([]string, 0),
		ScrollOffset: 0,
		LineHeight:   18,
		WrapWidth:    width - 20, // Account for padding and scrollbar
	}

	sta.SetSize(width, height)

	// Calculate visible lines based on height
	sta.VisibleLines = int(math.Floor(float64(height-10) / float64(sta.LineHeight)))

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1
	borderRadius := 4
	padding := NewEdgeInsets(6)

	sta.style.BorderWidth = &borderWidth
	sta.style.BorderRadius = &borderRadius
	sta.style.Padding = padding

	return sta
}

// GetType returns the element type
func (sta *ScrollingTextArea) GetType() string {
	return "ScrollingTextArea"
}

// AddText appends text to the area, wrapping and scrolling to bottom
func (sta *ScrollingTextArea) AddText(txt string) {
	// Wrap the new text
	wrapped := text.Wrap(txt, sta.WrapWidth, 15)
	sta.Lines = append(sta.Lines, wrapped...)

	// Auto-scroll to bottom
	maxOffset := len(sta.Lines) - sta.VisibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}
	sta.ScrollOffset = maxOffset
}

// SetText replaces all text in the area
func (sta *ScrollingTextArea) SetText(txt string) {
	sta.Lines = text.Wrap(txt, sta.WrapWidth, 15)
	sta.ScrollOffset = 0
}

// Clear removes all text
func (sta *ScrollingTextArea) Clear() {
	sta.Lines = make([]string, 0)
	sta.ScrollOffset = 0
}

// Update handles scrolling input
func (sta *ScrollingTextArea) Update() {
	if !sta.visible || !sta.enabled {
		return
	}

	sta.UpdateHoverState()

	mx, my := ebiten.CursorPosition()

	// Get absolute position for calculations
	absX, absY := sta.GetAbsolutePosition()

	// Scrollbar calculations
	barX := absX + sta.bounds.Width - 16
	barY := absY + 4
	barW := 12
	barH := sta.bounds.Height - 8

	totalLines := len(sta.Lines)
	if totalLines <= sta.VisibleLines {
		sta.draggingThumb = false
	} else {
		thumbH := int(math.Max(float64(barH*sta.VisibleLines/totalLines), 16))
		maxThumbY := barY + barH - thumbH
		thumbY := barY
		if totalLines > sta.VisibleLines {
			thumbY = barY + (barH-thumbH)*sta.ScrollOffset/(totalLines-sta.VisibleLines)
		}

		// Handle thumb dragging
		mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
		mouseJustReleased := !mousePressed && sta.draggingThumb

		if sta.draggingThumb {
			if mouseJustReleased {
				sta.draggingThumb = false
			} else if mousePressed {
				newThumbY := my - sta.dragOffsetY
				if newThumbY < barY {
					newThumbY = barY
				}
				if newThumbY > maxThumbY {
					newThumbY = maxThumbY
				}
				scrollRange := barH - thumbH
				if scrollRange > 0 {
					sta.ScrollOffset = int(float64(newThumbY-barY) / float64(scrollRange) * float64(totalLines-sta.VisibleLines))
				}
				sta.clampScrollOffset()
			}
		} else if mousePressed && mx >= barX && mx < barX+barW && my >= thumbY && my < thumbY+thumbH {
			sta.draggingThumb = true
			sta.dragOffsetY = my - thumbY
		}
	}

	// Mouse wheel scroll (only if hovered)
	if sta.hovered && !sta.draggingThumb {
		_, yoff := ebiten.Wheel()
		if yoff != 0 {
			sta.ScrollOffset -= int(yoff * 3) // Scroll 3 lines per wheel tick
			sta.clampScrollOffset()
		}
	}
}

func (sta *ScrollingTextArea) clampScrollOffset() {
	if sta.ScrollOffset < 0 {
		sta.ScrollOffset = 0
	}
	maxOffset := len(sta.Lines) - sta.VisibleLines
	if maxOffset < 0 {
		maxOffset = 0
	}
	if sta.ScrollOffset > maxOffset {
		sta.ScrollOffset = maxOffset
	}
}

// Layout calculates dimensions
func (sta *ScrollingTextArea) Layout() {
	style := sta.GetComputedStyle()

	width := sta.bounds.Width
	height := sta.bounds.Height

	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min/max constraints
	width, height = ApplySizeConstraints(width, height, style)

	sta.bounds.Width = width
	sta.bounds.Height = height

	// Recalculate visible lines
	sta.VisibleLines = int(math.Floor(float64(height-10) / float64(sta.LineHeight)))
	sta.WrapWidth = width - 20
}

// Draw draws the scrolling text area
func (sta *ScrollingTextArea) Draw(screen *ebiten.Image) {
	if !sta.visible {
		return
	}

	style := sta.GetComputedStyle()
	theme := sta.GetTheme()
	absX, absY := sta.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  sta.bounds.Width,
		Height: sta.bounds.Height,
	}

	// Check if we should use sprite-based rendering
	if theme != nil && theme.HasScrollingTextAreaSprites() {
		// Use 9-slice sprite rendering for background
		if sta.bgCache == nil || sta.bgCacheWidth != sta.bounds.Width || sta.bgCacheHeight != sta.bounds.Height {
			sta.bgCache = ebiten.NewImage(sta.bounds.Width, sta.bounds.Height)
			Draw9SliceToImage(sta.bgCache, theme.SpriteSheet, theme.ScrollingTextArea)
			sta.bgCacheWidth = sta.bounds.Width
			sta.bgCacheHeight = sta.bounds.Height
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(absX), float64(absY))
		screen.DrawImage(sta.bgCache, op)
	} else {
		// Use vector-based rendering with theme support
		DrawBackgroundWithTheme(screen, absBounds, style, theme)
		DrawBorderWithTheme(screen, absBounds, style, theme)
	}

	// Draw text lines
	contentBounds := GetContentBounds(absBounds, style)
	start := sta.ScrollOffset
	end := start + sta.VisibleLines
	if end > len(sta.Lines) {
		end = len(sta.Lines)
	}

	// Get text color from style, then theme, then default
	textColor := color.RGBA{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		textColor = colorToRGBA(*style.ForegroundColor)
	} else if theme != nil {
		textColor = colorToRGBA(theme.Colors.Text)
	}

	for i := start; i < end; i++ {
		lineY := contentBounds.Y + (i-start)*sta.LineHeight
		text.Draw(screen, sta.Lines[i], 14.0, contentBounds.X, lineY, textColor)
	}

	// Draw scrollbar if needed
	if len(sta.Lines) > sta.VisibleLines {
		sta.drawScrollbar(screen, absBounds, theme)
	}
}

func (sta *ScrollingTextArea) drawScrollbar(screen *ebiten.Image, absBounds Rect, theme *Theme) {
	barX := absBounds.X + absBounds.Width - 14
	barY := absBounds.Y + 4
	barW := 10
	barH := absBounds.Height - 8

	totalLines := len(sta.Lines)
	if totalLines <= sta.VisibleLines {
		return
	}

	// Check if we should use sprite-based rendering
	if theme != nil && theme.HasScrollbarSprites() {
		// Draw scrollbar track using sprite
		trackBounds := Rect{X: barX, Y: barY, Width: barW, Height: barH}
		trackCoords := &SpriteCoords{
			SrcX:   theme.Scrollbar.TrackX,
			SrcY:   theme.Scrollbar.TrackY,
			Width:  theme.Scrollbar.TrackWidth,
			Height: theme.Scrollbar.TrackHeight,
		}
		DrawSprite(screen, theme.SpriteSheet, trackCoords, trackBounds)

		// Draw thumb using sprite
		thumbH := int(math.Max(float64(barH*sta.VisibleLines/totalLines), 16))
		scrollRange := barH - thumbH
		var thumbY int
		if totalLines > sta.VisibleLines && scrollRange > 0 {
			thumbY = barY + (scrollRange*sta.ScrollOffset)/(totalLines-sta.VisibleLines)
		} else {
			thumbY = barY
		}

		thumbBounds := Rect{X: barX, Y: thumbY, Width: barW, Height: thumbH}
		thumbCoords := &SpriteCoords{
			SrcX:   theme.Scrollbar.ThumbX,
			SrcY:   theme.Scrollbar.ThumbY,
			Width:  theme.Scrollbar.ThumbWidth,
			Height: theme.Scrollbar.ThumbHeight,
		}
		DrawSprite(screen, theme.SpriteSheet, thumbCoords, thumbBounds)
	} else {
		// Use vector-based rendering with theme colors
		// Draw scrollbar track
		trackColor := color.RGBA{60, 60, 70, 200}
		if theme != nil {
			trackColor = colorToRGBA(theme.Colors.Surface)
			trackColor.A = 200
		}
		trackBounds := Rect{X: barX, Y: barY, Width: barW, Height: barH}
		DrawRoundedRect(screen, trackBounds, 4, trackColor)

		// Draw thumb
		thumbH := int(math.Max(float64(barH*sta.VisibleLines/totalLines), 16))
		scrollRange := barH - thumbH
		var thumbY int
		if totalLines > sta.VisibleLines && scrollRange > 0 {
			thumbY = barY + (scrollRange*sta.ScrollOffset)/(totalLines-sta.VisibleLines)
		} else {
			thumbY = barY
		}

		thumbColor := color.RGBA{120, 120, 140, 255}
		if theme != nil {
			thumbColor = colorToRGBA(theme.Colors.Border)
		}
		if sta.hovered || sta.draggingThumb {
			thumbColor.R = min(thumbColor.R+20, 255)
			thumbColor.G = min(thumbColor.G+20, 255)
			thumbColor.B = min(thumbColor.B+20, 255)
		}

		thumbBounds := Rect{X: barX, Y: thumbY, Width: barW, Height: thumbH}
		DrawRoundedRect(screen, thumbBounds, 4, thumbColor)
	}
}
