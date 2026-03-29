package minui

import "github.com/hajimehoshi/ebiten/v2"

// ImageWidget displays an *ebiten.Image as a managed UI element.
// The image is scaled to fit the widget's bounds.
type ImageWidget struct {
	*ElementBase
	Image *ebiten.Image
}

// NewImageWidget creates a new image widget with the given size.
func NewImageWidget(id string, width, height int) *ImageWidget {
	iw := &ImageWidget{
		ElementBase: NewElementBase(id),
	}
	iw.SetSize(width, height)
	return iw
}

// GetType returns the element type.
func (iw *ImageWidget) GetType() string { return "ImageWidget" }

// Update is a no-op for static image display.
func (iw *ImageWidget) Update() {}

// Layout is a no-op; size is set explicitly via SetSize.
func (iw *ImageWidget) Layout() {}

// Draw renders the image scaled to the widget's bounds.
func (iw *ImageWidget) Draw(screen *ebiten.Image) {
	if !iw.visible || iw.Image == nil {
		return
	}
	absX, absY := iw.GetAbsolutePosition()
	imgW := iw.Image.Bounds().Dx()
	imgH := iw.Image.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	if iw.bounds.Width > 0 && iw.bounds.Height > 0 && imgW > 0 && imgH > 0 {
		op.GeoM.Scale(
			float64(iw.bounds.Width)/float64(imgW),
			float64(iw.bounds.Height)/float64(imgH),
		)
	}
	op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(iw.Image, op)
}
