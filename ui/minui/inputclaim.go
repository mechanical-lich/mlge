package minui

import "github.com/hajimehoshi/ebiten/v2"

// Input claim: while certain widgets are active (currently an expanded
// SelectBox), they hold exclusive mouse input — clicks outside the claimed
// rect must not register on widgets behind them. This is set once per frame
// by GUI.Update before any widget's Update runs, so widget update order
// doesn't matter.
//
// We also swallow the entire press→release cycle if the press began during a
// claim. Otherwise widgets that track pressed state across frames (Button,
// IconButton) would re-arm on the very next frame after the dropdown closes
// and fire OnClick on release, piggybacking on the dropdown click.

var (
	inputClaimRect *Rect
	inputClaimant  Element
	pressSwallowed bool
)

// IsInputClaimed reports whether some widget is currently claiming input
// (e.g. an expanded SelectBox dropdown), or the in-flight mouse press
// originated during a claim. All interactive widgets except the claimant
// should call this at the top of Update and skip click handling when it
// returns true.
func IsInputClaimed() bool { return inputClaimRect != nil || pressSwallowed }

// IsInputClaimedByOther reports whether some widget other than `e` is
// claiming input this frame. The claimant itself uses this so it can keep
// responding to clicks while every other interactive widget bails.
//
// Note: this does NOT consult the press-swallow flag. The flag exists to
// stop widgets that track pressed-state across frames from re-arming after
// a claimed press, but it must not block the claimant from handling its
// own click on the press frame itself.
func IsInputClaimedByOther(e Element) bool {
	return inputClaimant != nil && inputClaimant != e
}

// ClaimedRect returns the active claim rect, or false if none.
func ClaimedRect() (Rect, bool) {
	if inputClaimRect == nil {
		return Rect{}, false
	}
	return *inputClaimRect, true
}

func claimInput(rect Rect, by Element) {
	inputClaimRect = &rect
	inputClaimant = by
}

func resetInputClaim() {
	inputClaimRect = nil
	inputClaimant = nil
}

// PrepareInputClaims is called by code that hosts widgets without going
// through GUI.Update (e.g. a screen that calls Update on each widget
// directly). Call it once at the top of the host's Update with the root
// elements that may contain expanded SelectBoxes; it resets and recomputes
// the input claim so click-through to widgets behind dropdowns is blocked.
func PrepareInputClaims(roots ...Element) {
	resetInputClaim()
	scanInputClaims(roots)
}

// scanInputClaims walks the element tree looking for widgets that want to
// claim input this frame. Currently only an expanded SelectBox claims. It
// also maintains the press-swallow flag: any mouse-down that occurs while a
// claim is active locks out other widgets for the rest of that press cycle.
func scanInputClaims(roots []Element) {
	for _, root := range roots {
		if scanElementClaim(root) {
			break
		}
	}
	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if !mouseDown {
		pressSwallowed = false
	} else if inputClaimRect != nil {
		pressSwallowed = true
	}
}

func scanElementClaim(e Element) bool {
	if !e.IsVisible() {
		return false
	}
	if sb, ok := e.(*SelectBox); ok && sb.expanded {
		sx, sy := sb.GetAbsolutePosition()
		w := sb.bounds.Width
		h := sb.bounds.Height
		if sb.listBox != nil {
			h += sb.listBox.bounds.Height
		}
		claimInput(Rect{X: sx, Y: sy, Width: w, Height: h}, sb)
		return true
	}
	for _, child := range e.GetChildren() {
		if scanElementClaim(child) {
			return true
		}
	}
	return false
}
