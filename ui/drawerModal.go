package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
	"github.com/mechanical-lich/mlge/utility"
)

type DrawSide string

const (
	DrawSideLeft  DrawSide = "left"
	DrawSideRight DrawSide = "right"
)

// DrawerModal is a modal that slides in/out from the left or right edge of the screen.
type DrawerModal struct {
	ElementBase

	Views       map[string]GUIViewInterface
	CurrentView string

	Side       DrawSide // "left" or "right"
	Open       bool
	SlideSpeed int // pixels per frame

	offset       int // current slide offset (in px)
	targetOffset int // where the drawer should be (0=open, -Width or +Width=closed)
	bg           *ebiten.Image
	OnClose      func()
}

// NewDrawerModal creates a new DrawerModal.
func NewDrawerModal(name string, side DrawSide, width, height, slideSpeed int, initialView string, views map[string]GUIViewInterface) *DrawerModal {
	d := &DrawerModal{
		ElementBase: ElementBase{
			Name:    name,
			X:       0,
			Y:       0,
			Width:   width,
			Height:  height,
			Visible: false,
			op:      &ebiten.DrawImageOptions{},
		},
		Views:       views,
		CurrentView: initialView,
		Side:        side,
		Open:        false,
		SlideSpeed:  slideSpeed,
	}
	if side == DrawSideLeft {
		d.offset = -width
		d.targetOffset = -width
		d.X = 0
	} else {
		d.offset = getScreenWidth()
		d.targetOffset = getScreenWidth()
		d.X = getScreenWidth() - width
	}
	d.Y = 0
	return d
}

// SetView switches the drawer to a different view state.
func (d *DrawerModal) SetView(name string) {
	if _, ok := d.Views[name]; ok {
		d.CurrentView = name
	}
}

// SetOpen opens or closes the drawer.
func (d *DrawerModal) SetOpen(open bool) {
	d.Open = open
	if d.Side == DrawSideLeft {
		if open {
			d.targetOffset = 0
		} else {
			d.targetOffset = -d.Width
		}
	} else {
		if open {
			d.targetOffset = getScreenWidth() - d.Width
		} else {
			d.targetOffset = getScreenWidth()
		}
	}
	d.Visible = open || d.offset != d.targetOffset // keep drawing while animating
}

// Update handles sliding and delegates to the current view.
func (d *DrawerModal) Update(s state.StateInterface) {
	if !d.Visible && d.offset == d.targetOffset {
		return
	}

	// Slide toward target offset
	if d.offset < d.targetOffset {
		d.offset += d.SlideSpeed
		if d.offset > d.targetOffset {
			d.offset = d.targetOffset
		}
	} else if d.offset > d.targetOffset {
		d.offset -= d.SlideSpeed
		if d.offset < d.targetOffset {
			d.offset = d.targetOffset
		}
	}

	// Hide when fully closed
	if !d.Open && d.offset == d.targetOffset {
		d.Visible = false
	}

	// Only update view if visible
	if d.Visible {
		if v, ok := d.Views[d.CurrentView]; ok {
			if d.Side == DrawSideLeft {
				v.SetPosition(d.offset, d.Y)
			} else {
				v.SetPosition(d.offset, d.Y)
			}
			v.UpdateElements(s)
			v.Update(s)
		}
	}
}

// Draw renders the drawer background and current view.
func (d *DrawerModal) Draw(screen *ebiten.Image, s state.StateInterface, theme *Theme) {
	if !d.Visible && d.offset == d.targetOffset {
		return
	}

	if d.bg == nil || d.bg.Bounds().Dx() != d.Width || d.bg.Bounds().Dy() != d.Height {
		d.bg = ebiten.NewImage(d.Width, d.Height)
		utility.Draw9Slice(d.bg, 0, 0, d.Width, d.Height, theme.ModalNineSlice.SrcX, theme.ModalNineSlice.SrcY, theme.ModalNineSlice.TileSize, theme.ModalNineSlice.TileScale)
	}
	// Draw the drawer background
	d.op.GeoM.Reset()
	d.op.GeoM.Translate(float64(d.offset), float64(d.Y))
	screen.DrawImage(d.bg, d.op)

	// Draw current view
	if v, ok := d.Views[d.CurrentView]; ok {
		v.SetPosition(d.offset, d.Y)
		v.Draw(screen, s, theme)
		v.DrawElements(screen, s, theme)
	}
}

// GetInputFocused delegates to the current view.
func (d *DrawerModal) GetInputFocused() bool {
	if v, ok := d.Views[d.CurrentView]; ok {
		return v.GetInputFocused()
	}
	return false
}

func (d *DrawerModal) WithinBounds(mouseX, mouseY int) bool {
	return mouseX >= d.offset && mouseX <= d.offset+d.Width && mouseY >= d.Y && mouseY <= d.Y+d.Height
}

func getScreenWidth() int {
	w, _ := ebiten.WindowSize()
	return w
}

func (d *DrawerModal) GetName() string {
	return d.Name
}

func (d *DrawerModal) GetPosition() (int, int) {
	return d.X, d.Y
}

func (d *DrawerModal) SetPosition(x, y int) {
	d.X = x
	d.Y = y
}

func (d *DrawerModal) OpenModal() {
	d.SetOpen(true)
}

func (d *DrawerModal) CloseModal() {
	d.SetOpen(false)
	if d.OnClose != nil {
		d.OnClose()
	}
}

func (d *DrawerModal) IsOpen() bool {
	return d.Open
}

func (d *DrawerModal) IsVisible() bool {
	return d.Visible
}
