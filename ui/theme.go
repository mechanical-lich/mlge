package ui

type Theme struct {
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
		ThumbX          int // Add for thumb sprite
		ThumbY          int
		ThumbWidth      int
		ThumbHeight     int
	}
}

var DefaultTheme = Theme{
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
		ThumbWidth:      8,
		ThumbHeight:     16,
	},
}
