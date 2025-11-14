package minui

import (
	"image/color"
)

// Style defines the visual styling for UI elements with inheritance support
type Style struct {
	// Layout
	Width     *int // nil means auto/fill
	Height    *int
	MinWidth  *int
	MinHeight *int
	MaxWidth  *int
	MaxHeight *int

	// Spacing
	Padding      *EdgeInsets
	Margin       *EdgeInsets
	BorderWidth  *int
	BorderRadius *int

	// Colors
	BackgroundColor *color.Color
	BorderColor     *color.Color
	ForegroundColor *color.Color // Text color

	// Background Image
	BackgroundImage *string // Resource ID

	// Font
	FontSize   *int
	FontBold   *bool
	FontItalic *bool

	// Alignment
	TextAlign *TextAlignment
	VertAlign *VerticalAlignment
	Align     *Alignment // For container child alignment

	// States (for interactive elements)
	HoverStyle    *Style
	ActiveStyle   *Style
	DisabledStyle *Style
	FocusStyle    *Style

	// Opacity
	Opacity *float64 // 0.0 to 1.0

	// Display
	Visible *bool
}

// EdgeInsets represents spacing on all four sides
type EdgeInsets struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// NewEdgeInsets creates edge insets with the same value on all sides
func NewEdgeInsets(all int) *EdgeInsets {
	return &EdgeInsets{Top: all, Right: all, Bottom: all, Left: all}
}

// NewEdgeInsetsLR creates edge insets with left-right and top-bottom values
func NewEdgeInsetsLR(vertical, horizontal int) *EdgeInsets {
	return &EdgeInsets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

// NewEdgeInsetsTRBL creates edge insets with individual values
func NewEdgeInsetsTRBL(top, right, bottom, left int) *EdgeInsets {
	return &EdgeInsets{Top: top, Right: right, Bottom: bottom, Left: left}
}

// TextAlignment represents horizontal text alignment
type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlignRight
)

// VerticalAlignment represents vertical alignment
type VerticalAlignment int

const (
	VertAlignTop VerticalAlignment = iota
	VertAlignMiddle
	VertAlignBottom
)

// Alignment represents container child alignment
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// DefaultStyle returns a basic default style
func DefaultStyle() *Style {
	visible := true
	opacity := 1.0
	fontSize := 14
	padding := NewEdgeInsets(0)
	margin := NewEdgeInsets(0)
	borderWidth := 0
	borderRadius := 0
	textAlign := TextAlignLeft
	vertAlign := VertAlignTop
	align := AlignStart
	bold := false
	italic := false

	fg := color.Color(color.RGBA{0, 0, 0, 255})
	bg := color.Color(color.RGBA{60, 60, 70, 255})
	border := color.Color(color.RGBA{80, 80, 90, 255})

	return &Style{
		Visible:         &visible,
		Opacity:         &opacity,
		FontSize:        &fontSize,
		FontBold:        &bold,
		FontItalic:      &italic,
		Padding:         padding,
		Margin:          margin,
		BorderWidth:     &borderWidth,
		BorderRadius:    &borderRadius,
		ForegroundColor: &fg,
		BackgroundColor: &bg,
		BorderColor:     &border,
		TextAlign:       &textAlign,
		VertAlign:       &vertAlign,
		Align:           &align,
	}
}

// Merge creates a new style by merging this style with a parent style
// Child style properties override parent style properties
func (s *Style) Merge(parent *Style) *Style {
	if parent == nil {
		return s
	}

	merged := &Style{}

	// Helper to pick child over parent
	merged.Width = s.Width
	if merged.Width == nil {
		merged.Width = parent.Width
	}

	merged.Height = s.Height
	if merged.Height == nil {
		merged.Height = parent.Height
	}

	merged.MinWidth = s.MinWidth
	if merged.MinWidth == nil {
		merged.MinWidth = parent.MinWidth
	}

	merged.MinHeight = s.MinHeight
	if merged.MinHeight == nil {
		merged.MinHeight = parent.MinHeight
	}

	merged.MaxWidth = s.MaxWidth
	if merged.MaxWidth == nil {
		merged.MaxWidth = parent.MaxWidth
	}

	merged.MaxHeight = s.MaxHeight
	if merged.MaxHeight == nil {
		merged.MaxHeight = parent.MaxHeight
	}

	merged.Padding = s.Padding
	if merged.Padding == nil {
		merged.Padding = parent.Padding
	}

	merged.Margin = s.Margin
	if merged.Margin == nil {
		merged.Margin = parent.Margin
	}

	merged.BorderWidth = s.BorderWidth
	if merged.BorderWidth == nil {
		merged.BorderWidth = parent.BorderWidth
	}

	merged.BorderRadius = s.BorderRadius
	if merged.BorderRadius == nil {
		merged.BorderRadius = parent.BorderRadius
	}

	merged.BackgroundColor = s.BackgroundColor
	if merged.BackgroundColor == nil {
		merged.BackgroundColor = parent.BackgroundColor
	}

	merged.BorderColor = s.BorderColor
	if merged.BorderColor == nil {
		merged.BorderColor = parent.BorderColor
	}

	merged.ForegroundColor = s.ForegroundColor
	if merged.ForegroundColor == nil {
		merged.ForegroundColor = parent.ForegroundColor
	}

	merged.BackgroundImage = s.BackgroundImage
	if merged.BackgroundImage == nil {
		merged.BackgroundImage = parent.BackgroundImage
	}

	merged.FontSize = s.FontSize
	if merged.FontSize == nil {
		merged.FontSize = parent.FontSize
	}

	merged.FontBold = s.FontBold
	if merged.FontBold == nil {
		merged.FontBold = parent.FontBold
	}

	merged.FontItalic = s.FontItalic
	if merged.FontItalic == nil {
		merged.FontItalic = parent.FontItalic
	}

	merged.TextAlign = s.TextAlign
	if merged.TextAlign == nil {
		merged.TextAlign = parent.TextAlign
	}

	merged.VertAlign = s.VertAlign
	if merged.VertAlign == nil {
		merged.VertAlign = parent.VertAlign
	}

	merged.Align = s.Align
	if merged.Align == nil {
		merged.Align = parent.Align
	}

	merged.Opacity = s.Opacity
	if merged.Opacity == nil {
		merged.Opacity = parent.Opacity
	}

	merged.Visible = s.Visible
	if merged.Visible == nil {
		merged.Visible = parent.Visible
	}

	// Merge state styles recursively
	if s.HoverStyle != nil {
		if parent.HoverStyle != nil {
			merged.HoverStyle = s.HoverStyle.Merge(parent.HoverStyle)
		} else {
			merged.HoverStyle = s.HoverStyle
		}
	} else {
		merged.HoverStyle = parent.HoverStyle
	}

	if s.ActiveStyle != nil {
		if parent.ActiveStyle != nil {
			merged.ActiveStyle = s.ActiveStyle.Merge(parent.ActiveStyle)
		} else {
			merged.ActiveStyle = s.ActiveStyle
		}
	} else {
		merged.ActiveStyle = parent.ActiveStyle
	}

	if s.DisabledStyle != nil {
		if parent.DisabledStyle != nil {
			merged.DisabledStyle = s.DisabledStyle.Merge(parent.DisabledStyle)
		} else {
			merged.DisabledStyle = s.DisabledStyle
		}
	} else {
		merged.DisabledStyle = parent.DisabledStyle
	}

	if s.FocusStyle != nil {
		if parent.FocusStyle != nil {
			merged.FocusStyle = s.FocusStyle.Merge(parent.FocusStyle)
		} else {
			merged.FocusStyle = s.FocusStyle
		}
	} else {
		merged.FocusStyle = parent.FocusStyle
	}

	return merged
}

// Copy creates a deep copy of the style
func (s *Style) Copy() *Style {
	if s == nil {
		return nil
	}

	copied := &Style{}

	if s.Width != nil {
		w := *s.Width
		copied.Width = &w
	}
	if s.Height != nil {
		h := *s.Height
		copied.Height = &h
	}
	if s.MinWidth != nil {
		w := *s.MinWidth
		copied.MinWidth = &w
	}
	if s.MinHeight != nil {
		h := *s.MinHeight
		copied.MinHeight = &h
	}
	if s.MaxWidth != nil {
		w := *s.MaxWidth
		copied.MaxWidth = &w
	}
	if s.MaxHeight != nil {
		h := *s.MaxHeight
		copied.MaxHeight = &h
	}

	if s.Padding != nil {
		p := *s.Padding
		copied.Padding = &p
	}
	if s.Margin != nil {
		m := *s.Margin
		copied.Margin = &m
	}
	if s.BorderWidth != nil {
		b := *s.BorderWidth
		copied.BorderWidth = &b
	}
	if s.BorderRadius != nil {
		r := *s.BorderRadius
		copied.BorderRadius = &r
	}

	if s.BackgroundColor != nil {
		c := *s.BackgroundColor
		copied.BackgroundColor = &c
	}
	if s.BorderColor != nil {
		c := *s.BorderColor
		copied.BorderColor = &c
	}
	if s.ForegroundColor != nil {
		c := *s.ForegroundColor
		copied.ForegroundColor = &c
	}

	if s.BackgroundImage != nil {
		img := *s.BackgroundImage
		copied.BackgroundImage = &img
	}

	if s.FontSize != nil {
		f := *s.FontSize
		copied.FontSize = &f
	}
	if s.FontBold != nil {
		b := *s.FontBold
		copied.FontBold = &b
	}
	if s.FontItalic != nil {
		i := *s.FontItalic
		copied.FontItalic = &i
	}

	if s.TextAlign != nil {
		a := *s.TextAlign
		copied.TextAlign = &a
	}
	if s.VertAlign != nil {
		a := *s.VertAlign
		copied.VertAlign = &a
	}
	if s.Align != nil {
		a := *s.Align
		copied.Align = &a
	}

	if s.Opacity != nil {
		o := *s.Opacity
		copied.Opacity = &o
	}
	if s.Visible != nil {
		v := *s.Visible
		copied.Visible = &v
	}

	copied.HoverStyle = s.HoverStyle.Copy()
	copied.ActiveStyle = s.ActiveStyle.Copy()
	copied.DisabledStyle = s.DisabledStyle.Copy()
	copied.FocusStyle = s.FocusStyle.Copy()

	return copied
}

// GetComputedStyle returns the appropriate style based on element state
func (s *Style) GetComputedStyle(hovered, active, disabled, focused bool) *Style {
	base := s

	if disabled && s.DisabledStyle != nil {
		return s.DisabledStyle.Merge(base)
	}

	if active && s.ActiveStyle != nil {
		return s.ActiveStyle.Merge(base)
	}

	if focused && s.FocusStyle != nil {
		return s.FocusStyle.Merge(base)
	}

	if hovered && s.HoverStyle != nil {
		return s.HoverStyle.Merge(base)
	}

	return base
}
