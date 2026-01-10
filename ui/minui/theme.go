package minui

import "image/color"

// Colors holds the color scheme for a theme
type Colors struct {
	Primary       color.Color
	Secondary     color.Color
	Background    color.Color
	Surface       color.Color
	Text          color.Color
	TextSecondary color.Color
	Border        color.Color
	Focus         color.Color
	Error         color.Color
	Success       color.Color
	Warning       color.Color
	Disabled      color.Color
}

// SpriteCoords holds the source coordinates and dimensions for a sprite
type SpriteCoords struct {
	SrcX   int
	SrcY   int
	Width  int
	Height int
}

// NineSliceCoords holds the source coordinates for 9-slice rendering
type NineSliceCoords struct {
	SrcX      int
	SrcY      int
	TileSize  int
	TileScale int
}

// ScrollbarCoords holds sprite coordinates for scrollbar elements
type ScrollbarCoords struct {
	TrackX      int
	TrackY      int
	TrackWidth  int
	TrackHeight int
	ThumbX      int
	ThumbY      int
	ThumbWidth  int
	ThumbHeight int
}

// Theme defines the visual styling for the UI, including optional sprite coordinates
type Theme struct {
	Name   string
	Colors Colors

	// Sprite sheet resource ID (e.g., "ui")
	SpriteSheet string

	// Sprite coordinates for each element type
	// If these are set and SpriteSheet is non-empty, sprite rendering will be used
	ModalNineSlice    *NineSliceCoords
	Button            *SpriteCoords
	ButtonPressed     *SpriteCoords // Optional pressed state sprite
	RadioButton       *SpriteCoords
	RadioButtonActive *SpriteCoords // Optional selected state sprite
	Toggle            *SpriteCoords
	ToggleOn          *SpriteCoords // Optional on state sprite
	InputField        *SpriteCoords
	ScrollingTextArea *NineSliceCoords
	Scrollbar         *ScrollbarCoords
	Checkbox          *SpriteCoords
	CheckboxChecked   *SpriteCoords // Optional checked state sprite
	Slider            *struct {
		Track *SpriteCoords
		Thumb *SpriteCoords
	}
	Dropdown *SpriteCoords
}

// HasSprites returns true if the theme has sprite coordinates set
func (t *Theme) HasSprites() bool {
	return t != nil && t.SpriteSheet != ""
}

// HasModalSprites returns true if the theme has modal sprite coordinates
func (t *Theme) HasModalSprites() bool {
	return t.HasSprites() && t.ModalNineSlice != nil
}

// HasButtonSprites returns true if the theme has button sprite coordinates
func (t *Theme) HasButtonSprites() bool {
	return t.HasSprites() && t.Button != nil
}

// HasRadioButtonSprites returns true if the theme has radio button sprite coordinates
func (t *Theme) HasRadioButtonSprites() bool {
	return t.HasSprites() && t.RadioButton != nil
}

// HasToggleSprites returns true if the theme has toggle sprite coordinates
func (t *Theme) HasToggleSprites() bool {
	return t.HasSprites() && t.Toggle != nil
}

// HasScrollingTextAreaSprites returns true if the theme has scrolling text area sprite coordinates
func (t *Theme) HasScrollingTextAreaSprites() bool {
	return t.HasSprites() && t.ScrollingTextArea != nil
}

// HasScrollbarSprites returns true if the theme has scrollbar sprite coordinates
func (t *Theme) HasScrollbarSprites() bool {
	return t.HasSprites() && t.Scrollbar != nil
}

// Default color schemes
var DefaultColors = Colors{
	Primary:       color.RGBA{100, 150, 255, 255},
	Secondary:     color.RGBA{150, 150, 150, 255},
	Background:    color.RGBA{40, 40, 50, 255},
	Surface:       color.RGBA{60, 60, 70, 255},
	Text:          color.RGBA{255, 255, 255, 255},
	TextSecondary: color.RGBA{180, 180, 190, 255},
	Border:        color.RGBA{80, 80, 90, 255},
	Focus:         color.RGBA{120, 180, 255, 255},
	Error:         color.RGBA{255, 100, 100, 255},
	Success:       color.RGBA{100, 255, 150, 255},
	Warning:       color.RGBA{255, 200, 100, 255},
	Disabled:      color.RGBA{100, 100, 110, 255},
}

var DarkColors = Colors{
	Primary:       color.RGBA{80, 120, 200, 255},
	Secondary:     color.RGBA{120, 120, 120, 255},
	Background:    color.RGBA{25, 25, 30, 255},
	Surface:       color.RGBA{40, 40, 45, 255},
	Text:          color.RGBA{240, 240, 245, 255},
	TextSecondary: color.RGBA{160, 160, 170, 255},
	Border:        color.RGBA{60, 60, 70, 255},
	Focus:         color.RGBA{100, 150, 220, 255},
	Error:         color.RGBA{220, 80, 80, 255},
	Success:       color.RGBA{80, 200, 120, 255},
	Warning:       color.RGBA{220, 180, 80, 255},
	Disabled:      color.RGBA{80, 80, 90, 255},
}

var LightColors = Colors{
	Primary:       color.RGBA{70, 130, 220, 255},
	Secondary:     color.RGBA{140, 140, 140, 255},
	Background:    color.RGBA{245, 245, 250, 255},
	Surface:       color.RGBA{255, 255, 255, 255},
	Text:          color.RGBA{30, 30, 40, 255},
	TextSecondary: color.RGBA{100, 100, 110, 255},
	Border:        color.RGBA{200, 200, 210, 255},
	Focus:         color.RGBA{90, 150, 240, 255},
	Error:         color.RGBA{220, 60, 60, 255},
	Success:       color.RGBA{60, 180, 100, 255},
	Warning:       color.RGBA{220, 160, 60, 255},
	Disabled:      color.RGBA{180, 180, 190, 255},
}

// NewDefaultTheme creates a theme with default colors and no sprites (vector rendering)
func NewDefaultTheme() *Theme {
	return &Theme{
		Name:   "Default",
		Colors: DefaultColors,
	}
}

// NewDarkTheme creates a dark theme with no sprites (vector rendering)
func NewDarkTheme() *Theme {
	return &Theme{
		Name:   "Dark",
		Colors: DarkColors,
	}
}

// NewLightTheme creates a light theme with no sprites (vector rendering)
func NewLightTheme() *Theme {
	return &Theme{
		Name:   "Light",
		Colors: LightColors,
	}
}

// NewSpriteTheme creates a theme that uses sprites from a sprite sheet.
// This is a convenience function that sets up common sprite coordinates.
// You can customize the coordinates after creation.
func NewSpriteTheme(name, spriteSheet string) *Theme {
	return &Theme{
		Name:        name,
		Colors:      DefaultColors,
		SpriteSheet: spriteSheet,
	}
}

// WithModalNineSlice adds 9-slice modal sprite coordinates to the theme
func (t *Theme) WithModalNineSlice(srcX, srcY, tileSize, tileScale int) *Theme {
	t.ModalNineSlice = &NineSliceCoords{
		SrcX:      srcX,
		SrcY:      srcY,
		TileSize:  tileSize,
		TileScale: tileScale,
	}
	return t
}

// WithButton adds button sprite coordinates to the theme
func (t *Theme) WithButton(srcX, srcY, width, height int) *Theme {
	t.Button = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// WithButtonPressed adds pressed button sprite coordinates to the theme
func (t *Theme) WithButtonPressed(srcX, srcY, width, height int) *Theme {
	t.ButtonPressed = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// WithRadioButton adds radio button sprite coordinates to the theme
func (t *Theme) WithRadioButton(srcX, srcY, width, height int) *Theme {
	t.RadioButton = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// WithRadioButtonActive adds active/selected radio button sprite coordinates to the theme
func (t *Theme) WithRadioButtonActive(srcX, srcY, width, height int) *Theme {
	t.RadioButtonActive = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// WithToggle adds toggle sprite coordinates to the theme
func (t *Theme) WithToggle(srcX, srcY, width, height int) *Theme {
	t.Toggle = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// WithToggleOn adds on-state toggle sprite coordinates to the theme
func (t *Theme) WithToggleOn(srcX, srcY, width, height int) *Theme {
	t.ToggleOn = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// WithScrollingTextArea adds scrolling text area sprite coordinates to the theme
func (t *Theme) WithScrollingTextArea(srcX, srcY, tileSize, tileScale int) *Theme {
	t.ScrollingTextArea = &NineSliceCoords{
		SrcX:      srcX,
		SrcY:      srcY,
		TileSize:  tileSize,
		TileScale: tileScale,
	}
	return t
}

// WithScrollbar adds scrollbar sprite coordinates to the theme
func (t *Theme) WithScrollbar(trackX, trackY, trackW, trackH, thumbX, thumbY, thumbW, thumbH int) *Theme {
	t.Scrollbar = &ScrollbarCoords{
		TrackX:      trackX,
		TrackY:      trackY,
		TrackWidth:  trackW,
		TrackHeight: trackH,
		ThumbX:      thumbX,
		ThumbY:      thumbY,
		ThumbWidth:  thumbW,
		ThumbHeight: thumbH,
	}
	return t
}

// WithCheckbox adds checkbox sprite coordinates to the theme
func (t *Theme) WithCheckbox(srcX, srcY, width, height int) *Theme {
	t.Checkbox = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// WithCheckboxChecked adds checked checkbox sprite coordinates to the theme
func (t *Theme) WithCheckboxChecked(srcX, srcY, width, height int) *Theme {
	t.CheckboxChecked = &SpriteCoords{
		SrcX:   srcX,
		SrcY:   srcY,
		Width:  width,
		Height: height,
	}
	return t
}

// FantasySettlementsTheme creates a theme matching the original fantasy_settlements UI
func FantasySettlementsTheme() *Theme {
	return NewSpriteTheme("FantasySettlements", "ui").
		WithModalNineSlice(144, 0, 16, 2).
		WithButton(16, 64, 32, 16).
		WithButtonPressed(48, 64, 32, 16).
		WithRadioButton(16, 64, 32, 16).
		WithRadioButtonActive(48, 64, 32, 16).
		WithToggle(16, 64, 32, 16).
		WithToggleOn(48, 64, 32, 16).
		WithScrollingTextArea(192, 0, 16, 2).
		WithScrollbar(96, 48, 16, 48, 112, 48, 16, 16)
}
