package main

import (
	"fmt"
	"image/color"

	"github.com/mechanical-lich/mlge/ecs"
)

// ---- Component types --------------------------------------------------------

const (
	TypeID       ecs.ComponentType = "ID"
	TypePosition ecs.ComponentType = "Position"
	TypeVelocity ecs.ComponentType = "Velocity"
	TypeColor    ecs.ComponentType = "Color"
)

// IDComponent assigns a stable string identity to an entity so the snapshot
// codec can correlate server entities with client-side counterparts.
type IDComponent struct{ ID string }

func (c IDComponent) GetType() ecs.ComponentType { return TypeID }

// PositionComponent stores the entity's world-space position.
type PositionComponent struct{ X, Y float64 }

func (c PositionComponent) GetType() ecs.ComponentType { return TypePosition }

// VelocityComponent stores the entity's velocity in pixels per second.
type VelocityComponent struct{ VX, VY float64 }

func (c VelocityComponent) GetType() ecs.ComponentType { return TypeVelocity }

// ColorComponent stores the entity's render color.
type ColorComponent struct{ RGBA color.RGBA }

func (c ColorComponent) GetType() ecs.ComponentType { return TypeColor }

// ---- Helpers ----------------------------------------------------------------

func ballID(i int) string { return fmt.Sprintf("ball-%d", i) }

var palette = []color.RGBA{
	{220, 60, 60, 255},
	{60, 180, 60, 255},
	{60, 120, 220, 255},
	{220, 180, 40, 255},
	{180, 60, 220, 255},
	{40, 200, 200, 255},
	{220, 120, 40, 255},
	{200, 200, 200, 255},
	{220, 80, 140, 255},
	{80, 160, 80, 255},
	{80, 140, 220, 255},
	{220, 160, 80, 255},
}

func ballColor(i int) color.RGBA {
	return palette[i%len(palette)]
}
