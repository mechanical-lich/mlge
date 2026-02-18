---
layout: default
title: Text
nav_order: 13
---

# Text Rendering

`github.com/mechanical-lich/mlge/text`

Text rendering, measurement, and word wrapping using an embedded Roboto Regular font via Ebitengine's text/v2.

## Functions

### Draw

```go
func Draw(dst *ebiten.Image, txt string, size float64, x, y int, clr color.Color)
```

Draws text onto an image at the specified position with the given size and color.

```go
import "github.com/mechanical-lich/mlge/text"

text.Draw(screen, "Hello, World!", 16, 100, 50, color.White)
```

### Measure

```go
func Measure(txt string, size float64) (width float64, height float64)
```

Returns the width and height of the rendered text at the given size. Useful for centering or layout calculations.

```go
w, h := text.Measure("Score: 1000", 14)
```

### Wrap

```go
func Wrap(s string, maxChars, maxLines int) []string
```

Word-wraps a string to fit within character and line limits. Respects existing newlines in the input.

```go
lines := text.Wrap("This is a long message that needs to be wrapped", 20, 5)
for i, line := range lines {
    text.Draw(screen, line, 14, 10, 20 + i*18, color.White)
}
```
