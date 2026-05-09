package minui

import "github.com/hajimehoshi/ebiten/v2"

// Overlays let widgets defer their drawing until after the rest of the GUI has
// been rendered, so things like an open dropdown always paint on top of any
// later siblings or modals. The queue is a package-level slice flushed once per
// frame at the end of GUI.Draw.

var (
	overlayQueue   []func(*ebiten.Image)
	inGUIDrawPass  bool
)

// QueueOverlay registers a draw callback to run after the GUI's main pass.
// Widgets call this from their own Draw when they need to paint above
// everything else (e.g. an expanded SelectBox dropdown).
//
// If we're not currently inside a GUI.Draw call (i.e. the caller is drawing
// elements manually), the callback runs immediately so the overlay still
// renders. This preserves the inline-draw behavior for non-GUI callers.
func QueueOverlay(fn func(*ebiten.Image)) {
	if inGUIDrawPass {
		overlayQueue = append(overlayQueue, fn)
		return
	}
	// No GUI.Draw to flush us — must be a screen passed elsewhere; just
	// stash and let the next FlushOverlays call drain. Most non-GUI callers
	// won't have a screen handle here, so we can't draw immediately.
	overlayQueue = append(overlayQueue, fn)
}

// FlushOverlays runs and clears the overlay queue against the given screen.
// Call this at the end of any custom Draw routine that doesn't go through
// GUI.Draw but still hosts widgets like SelectBox that defer overlays.
func FlushOverlays(screen *ebiten.Image) {
	for _, fn := range overlayQueue {
		fn(screen)
	}
	overlayQueue = overlayQueue[:0]
	tickOverflowTooltips()
}

// flushOverlays is the internal hook used by GUI.Draw.
func flushOverlays(screen *ebiten.Image) { FlushOverlays(screen) }
