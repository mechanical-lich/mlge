package minui

// Input claim: while certain widgets are active (currently an expanded
// SelectBox), they hold exclusive mouse input — clicks outside the claimed
// rect must not register on widgets behind them. This is set once per frame
// by GUI.Update before any widget's Update runs, so widget update order
// doesn't matter.

var inputClaimRect *Rect

// IsInputClaimed reports whether some widget is currently claiming input.
func IsInputClaimed() bool { return inputClaimRect != nil }

// IsInputClaimedOutside reports whether input is claimed AND the given
// screen-space point is outside the claim. Interactive widgets (Button,
// MenuItem, IconButton, Toggle, RadioButton, Checkbox, ...) should call this
// at the top of Update and skip click handling when it returns true.
func IsInputClaimedOutside(x, y int) bool {
	return inputClaimRect != nil && !inputClaimRect.Contains(x, y)
}

func claimInput(rect Rect) { inputClaimRect = &rect }
func resetInputClaim()     { inputClaimRect = nil }

// scanInputClaims walks the element tree looking for widgets that want to
// claim input this frame. Currently only an expanded SelectBox claims.
func scanInputClaims(roots []Element) {
	for _, root := range roots {
		if scanElementClaim(root) {
			return
		}
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
		claimInput(Rect{X: sx, Y: sy, Width: w, Height: h})
		return true
	}
	for _, child := range e.GetChildren() {
		if scanElementClaim(child) {
			return true
		}
	}
	return false
}
