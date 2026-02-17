package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text"
	elements "github.com/mechanical-lich/mlge/ui/themed/elements"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
)

type TabbedContainer struct {
	elements.ElementBase
	Tabs        []string
	ActiveTab   int
	Containers  []ContainerInterface
	OnTabChange func(tabIndex int)
}

type ContainerInterface interface {
	Update()
	Draw(screen *ebiten.Image, theme *theming.Theme)
	SetPosition(x, y int)
	AddChild(child elements.ElementInterface)
	GetChild(name string) elements.ElementInterface
	RemoveChild(name string)
}

func NewTabbedContainer(name string, x, y, width, height int, tabs []string, containers []ContainerInterface) *TabbedContainer {
	return &TabbedContainer{
		ElementBase: elements.ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
			Op:     &ebiten.DrawImageOptions{},
		},
		Tabs:       tabs,
		ActiveTab:  0,
		Containers: containers,
	}
}

func (tc *TabbedContainer) Update() {
	cX, cY := ebiten.CursorPosition()
	absX, absY := tc.ElementBase.GetAbsolutePosition()

	// Tab click detection
	for i := range tc.Tabs {
		tabX := absX + i*100 // Tab width 100px
		tabY := absY
		if cX >= tabX && cX <= tabX+100 && cY >= tabY && cY <= tabY+30 {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				tc.ActiveTab = i
				if tc.OnTabChange != nil {
					tc.OnTabChange(i)
				}
			}
		}
	}

	// Update active container
	if tc.ActiveTab >= 0 && tc.ActiveTab < len(tc.Containers) {
		tc.Containers[tc.ActiveTab].SetPosition(tc.X, tc.Y+30)
		tc.Containers[tc.ActiveTab].Update()
	}
}

func (tc *TabbedContainer) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := tc.ElementBase.GetAbsolutePosition()

	// Draw tabs
	for i := range tc.Tabs {
		tabX := absX + i*100
		tabY := absY
		clr := color.RGBA{80, 80, 120, 255}
		if i == tc.ActiveTab {
			clr = color.RGBA{120, 120, 180, 255}
		}
		img := ebiten.NewImage(100, 30)
		img.Fill(clr)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tabX), float64(tabY))
		screen.DrawImage(img, op)
		text.Draw(screen, tc.Tabs[i], 15, tabX+10, tabY+8, color.White)
	}

	// Draw active container
	if tc.ActiveTab >= 0 && tc.ActiveTab < len(tc.Containers) {
		tc.Containers[tc.ActiveTab].Draw(screen, theme)
	}
}

// GetCurrentContainer returns the currently active container.
func (tc *TabbedContainer) GetCurrentContainer() ContainerInterface {
	if tc.ActiveTab >= 0 && tc.ActiveTab < len(tc.Containers) {
		return tc.Containers[tc.ActiveTab]
	}
	return nil
}
