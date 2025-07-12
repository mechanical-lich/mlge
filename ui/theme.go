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
	// Add more as needed
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
}
