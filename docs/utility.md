---
layout: default
title: Utilities
nav_order: 16
---

# Utility Functions

`github.com/mechanical-lich/mlge/utility`

A collection of math, geometry, array, and drawing helper functions.

## Integer Math

```go
import "github.com/mechanical-lich/mlge/utility"
```

| Function | Signature | Description |
|----------|-----------|-------------|
| `GetRandom` | `(low, high int) int` | Random integer in range [low, high] |
| `Distance` | `(x1, y1, x2, y2 int) float64` | Euclidean distance between two points |
| `Sgn` | `(a int) int` | Sign of an integer (-1, 0, or 1) |
| `Clamp` | `(value, min, max int) int` | Clamp value to range |
| `Wrap` | `(value, max int) int` | Wrap value around range [0, max) |
| `Abs` | `(a int) int` | Absolute value |
| `Max` | `(a, b int) int` | Maximum of two values |
| `Min` | `(a, b int) int` | Minimum of two values |

## Float Math

| Function | Signature | Description |
|----------|-----------|-------------|
| `ClampF` | `(value, min, max float64) float64` | Clamp float to range |
| `WrapF` | `(value, max float64) float64` | Wrap float around range |
| `AbsF` | `(a float64) float64` | Absolute value (float) |
| `MaxF` | `(a, b float64) float64` | Maximum of two floats |
| `MinF` | `(a, b float64) float64` | Minimum of two floats |
| `Lerp` | `(a, b, t float64) float64` | Linear interpolation |
| `LerpAngle` | `(a, b, t float64) float64` | Angle-aware linear interpolation |

## Geometry

| Function | Signature | Description |
|----------|-----------|-------------|
| `RectsOverlap` | `(x1, y1, w1, h1, x2, y2, w2, h2 int) bool` | Check if two rectangles overlap |
| `RectContains` | `(rx, ry, rw, rh, px, py int) bool` | Check if a point is inside a rectangle |

## Array Helpers

| Function | Signature | Description |
|----------|-----------|-------------|
| `Contains` | `(slice []string, item string) bool` | Check if a string slice contains an item |

## Drawing

### Draw9Slice

```go
func Draw9Slice(dst *ebiten.Image, x, y, w, h, srcX, srcY, tileSize, tileScale int)
```

Renders a 9-slice sprite from the `"ui"` texture in the resource manager. The source region at `(srcX, srcY)` is divided into a 3x3 grid of `tileSize` pixels, then drawn scaled to fill the target area `(x, y, w, h)`. Corners are drawn at fixed size, edges are stretched, and the center is tiled.

```go
// Draw a 200x100 panel using a 9-slice sprite starting at (0, 0) with 8px tiles at 2x scale
utility.Draw9Slice(screen, 10, 10, 200, 100, 0, 0, 8, 2)
```
