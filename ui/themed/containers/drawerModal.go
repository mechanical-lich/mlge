package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	elements "github.com/mechanical-lich/mlge/ui/themed/elements"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"

	"github.com/mechanical-lich/mlge/utility"
)

type DrawSide string

const (
	DrawSideLeft   DrawSide = "left"
	DrawSideRight  DrawSide = "right"
	DrawSideTop    DrawSide = "top"
	DrawSideBottom DrawSide = "bottom"
)

// DrawerModal is a modal that slides in/out from the left, right, top, or bottom edge of the screen.
// Now supports positioning along the axis it is anchored to (Y for left/right, X for top/bottom).
type DrawerModal struct {
	elements.ElementBase

	Side       DrawSide // "left", "right", "top", or "bottom"
	Open       bool
	SlideSpeed int // pixels per frame

	offset       int // current slide offset (in px)
	targetOffset int // where the drawer should be (0=open, -Width/-Height or +Width/+Height=closed)
	bg           *ebiten.Image
	OnClose      func()

	CloseButton *elements.Button

	Children []elements.ElementInterface // Children elements
}

// NewDrawerModal creates a new DrawerModal.
// x and y are now respected as the position along the axis the drawer is anchored to.
func NewDrawerModal(name string, side DrawSide, x, y, width, height, slideSpeed int) *DrawerModal {
	var closeBtn *elements.Button
	switch side {
	case DrawSideLeft:
		closeBtn = elements.NewButton("close", width-24, 8, "X", "close")
	case DrawSideRight:
		closeBtn = elements.NewButton("close", 8, 8, "X", "close")
	case DrawSideTop:
		closeBtn = elements.NewButton("close", width-24, 8, "X", "close")
	case DrawSideBottom:
		closeBtn = elements.NewButton("close", width-24, 8, "X", "close")
	default:
		closeBtn = elements.NewButton("close", width-24, 8, "X", "close")
	}
	d := &DrawerModal{
		ElementBase: elements.ElementBase{
			Name:    name,
			X:       x,
			Y:       y,
			Width:   width,
			Height:  height,
			Visible: false,
			Op:      &ebiten.DrawImageOptions{},
		},
		Side:        side,
		Open:        false,
		SlideSpeed:  slideSpeed,
		CloseButton: closeBtn,
		Children:    []elements.ElementInterface{},
	}
	switch side {
	case DrawSideLeft:
		d.offset = -width
		d.targetOffset = -width
		// d.X = 0 // X is always 0 for left, Y is set by user
	case DrawSideRight:
		d.offset = getScreenWidth()
		d.targetOffset = getScreenWidth()
		d.X = getScreenWidth() - width // X is always right edge, Y is set by user
	case DrawSideTop:
		d.offset = -height
		d.targetOffset = -height
		// d.Y = 0 // Y is always 0 for top, X is set by user
	case DrawSideBottom:
		d.offset = getScreenHeight()
		d.targetOffset = getScreenHeight()
		d.Y = getScreenHeight() - height // Y is always bottom edge, X is set by user
	default:
		d.offset = -width
		d.targetOffset = -width
	}
	return d
}

// SetOpen opens or closes the drawer.
func (d *DrawerModal) SetOpen(open bool) {
	d.Open = open
	switch d.Side {
	case DrawSideLeft:
		if open {
			d.targetOffset = 0
		} else {
			d.targetOffset = -d.Width
		}
	case DrawSideRight:
		if open {
			d.targetOffset = getScreenWidth() - d.Width
		} else {
			d.targetOffset = getScreenWidth()
		}
	case DrawSideTop:
		if open {
			d.targetOffset = 0
		} else {
			d.targetOffset = -d.Height
		}
	case DrawSideBottom:
		if open {
			d.targetOffset = getScreenHeight() - d.Height
		} else {
			d.targetOffset = getScreenHeight()
		}
	}
	d.Visible = open || d.offset != d.targetOffset // keep drawing while animating
}

// Update handles sliding and delegates to the current view.
func (d *DrawerModal) Update() {
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

	// Update close button position based on side and offset
	switch d.Side {
	case DrawSideLeft:
		d.CloseButton.SetPosition(d.offset+d.Width-24, d.Y+8)
	case DrawSideRight:
		d.CloseButton.SetPosition(d.offset+8, d.Y+8)
	case DrawSideTop:
		d.CloseButton.SetPosition(d.X+d.Width-24, d.offset+8)
	case DrawSideBottom:
		d.CloseButton.SetPosition(d.X+d.Width-24, d.offset+8)
	}

	// Handle close button click
	if d.CloseButton.IsJustClicked() {
		d.SetOpen(false)
		if d.OnClose != nil {
			d.OnClose()
		}
		return
	}

	// Only update children if visible
	if d.Visible {
		for _, child := range d.Children {
			child.Update()
		}
	}
}

// Draw renders the drawer background and current view.
func (d *DrawerModal) Draw(screen *ebiten.Image, theme *theming.Theme) {
	if !d.Visible && d.offset == d.targetOffset {
		return
	}

	if d.bg == nil || d.bg.Bounds().Dx() != d.Width || d.bg.Bounds().Dy() != d.Height {
		d.bg = ebiten.NewImage(d.Width, d.Height)
		utility.Draw9Slice(d.bg, 0, 0, d.Width, d.Height, theme.ModalNineSlice.SrcX, theme.ModalNineSlice.SrcY, theme.ModalNineSlice.TileSize, theme.ModalNineSlice.TileScale)
	}
	// Draw the drawer background
	d.Op.GeoM.Reset()
	switch d.Side {
	case DrawSideLeft, DrawSideRight:
		d.Op.GeoM.Translate(float64(d.offset), float64(d.Y))
	case DrawSideTop, DrawSideBottom:
		d.Op.GeoM.Translate(float64(d.X), float64(d.offset))
	}
	screen.DrawImage(d.bg, d.Op)

	// Draw close button
	d.CloseButton.Draw(screen, theme)

	// Draw children
	for _, child := range d.Children {
		child.Draw(screen, theme)
		fmt.Println(child.GetName())
		fmt.Println(child.GetAbsolutePosition())
	}
}

// GetInputFocused checks if any child is focused.
func (d *DrawerModal) GetInputFocused() bool {
	for _, child := range d.Children {
		if child.GetFocused() {
			return true
		}
	}
	return false
}

// AddChild adds a child element to the drawer modal.
func (d *DrawerModal) AddChild(child elements.ElementInterface) {
	child.SetParent(d)
	d.Children = append(d.Children, child)
}

// RemoveChild removes a child element from the drawer modal.
func (d *DrawerModal) RemoveChild(child elements.ElementInterface) {
	for i, c := range d.Children {
		if c == child {
			d.Children = append(d.Children[:i], d.Children[i+1:]...)
			child.SetParent(nil)
			break
		}
	}
}

// WithinBounds checks if a point is within the drawer, respecting axis anchoring.
func (d *DrawerModal) WithinBounds(mouseX, mouseY int) bool {
	switch d.Side {
	case DrawSideLeft, DrawSideRight:
		return mouseX >= d.offset && mouseX <= d.offset+d.Width && mouseY >= d.Y && mouseY <= d.Y+d.Height
	case DrawSideTop, DrawSideBottom:
		return mouseX >= d.X && mouseX <= d.X+d.Width && mouseY >= d.offset && mouseY <= d.offset+d.Height
	default:
		return false
	}
}

func getScreenWidth() int {
	w, _ := ebiten.WindowSize()
	return w
}

func getScreenHeight() int {
	_, h := ebiten.WindowSize()
	return h
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

func (m *DrawerModal) GetAbsolutePosition() (int, int) {
	return m.X, m.Y
}
