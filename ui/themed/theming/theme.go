package ui

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

type Theme struct {
	Name   string
	Colors Colors

	ModalNineSlice struct {
		SrcX      int
		SrcY      int
		TileSize  int
		TileScale int
	}
	Button struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}
	RadioButton struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}
	Toggle struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}
	InputField struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}
	ScrollingTextArea struct {
		SrcX            int
		SrcY            int
		TileSize        int
		TileScale       int
		ScrollBarX      int
		ScrollBarY      int
		ScrollBarWidth  int
		ScrollBarHeight int
		ThumbX          int
		ThumbY          int
		ThumbWidth      int
		ThumbHeight     int
	}
	Checkbox struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}
	Slider struct {
		TrackSrcX   int
		TrackSrcY   int
		TrackWidth  int
		TrackHeight int
		ThumbSrcX   int
		ThumbSrcY   int
		ThumbWidth  int
		ThumbHeight int
	}
	Dropdown struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}
}

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

var DefaultTheme = Theme{
	Name:   "Default",
	Colors: DefaultColors,
	ModalNineSlice: struct {
		SrcX      int
		SrcY      int
		TileSize  int
		TileScale int
	}{SrcX: 144, SrcY: 0, TileSize: 16, TileScale: 2},
	Button: struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}{SrcX: 16, SrcY: 64, Width: 32, Height: 16},
	RadioButton: struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}{SrcX: 16, SrcY: 64, Width: 32, Height: 16},
	Toggle: struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}{SrcX: 16, SrcY: 64, Width: 32, Height: 16},
	InputField: struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}{SrcX: 16, SrcY: 80, Width: 48, Height: 16},
	ScrollingTextArea: struct {
		SrcX            int
		SrcY            int
		TileSize        int
		TileScale       int
		ScrollBarX      int
		ScrollBarY      int
		ScrollBarWidth  int
		ScrollBarHeight int
		ThumbX          int
		ThumbY          int
		ThumbWidth      int
		ThumbHeight     int
	}{
		SrcX:            192,
		SrcY:            0,
		TileSize:        16,
		TileScale:       2,
		ScrollBarX:      96,
		ScrollBarY:      48,
		ScrollBarWidth:  16,
		ScrollBarHeight: 48,
		ThumbX:          112,
		ThumbY:          48,
		ThumbWidth:      16,
		ThumbHeight:     16,
	},
	Checkbox: struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}{SrcX: 80, SrcY: 64, Width: 16, Height: 16},
	Slider: struct {
		TrackSrcX   int
		TrackSrcY   int
		TrackWidth  int
		TrackHeight int
		ThumbSrcX   int
		ThumbSrcY   int
		ThumbWidth  int
		ThumbHeight int
	}{
		TrackSrcX:   64,
		TrackSrcY:   80,
		TrackWidth:  48,
		TrackHeight: 8,
		ThumbSrcX:   64,
		ThumbSrcY:   96,
		ThumbWidth:  12,
		ThumbHeight: 16,
	},
	Dropdown: struct {
		SrcX   int
		SrcY   int
		Width  int
		Height int
	}{SrcX: 16, SrcY: 96, Width: 48, Height: 16},
}
