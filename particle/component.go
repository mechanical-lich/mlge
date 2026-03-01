package particle

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/ecs"
)

// ComponentType is the ECS key for [EmitterComponent].
const ComponentType ecs.ComponentType = "particle.Emitter"

// EmitterComponent configures a particle emitter attached to an ECS entity.
//
// Attach it to an entity with entity.AddComponent(EmitterComponent{...}).
// Update X and Y each tick to follow the entity's position.
//
// The zero value is safe but produces no particles (Active is false, EmitRate
// is zero). Set at minimum Active, EmitRate, LifeMin/Max and a colour pair.
type EmitterComponent struct {
	// X and Y are the world-space origin of the emitter.
	// Update these each tick to match the parent entity's position.
	X, Y float64

	// StartColor and EndColor define the per-particle color gradient.
	// A particle at birth uses StartColor; at death it uses EndColor.
	// The alpha channel is also interpolated, so fading out is free.
	StartColor color.RGBA
	EndColor   color.RGBA

	// StartSize and EndSize define the particle radius in pixels over its
	// lifetime. StartSize is applied at birth, EndSize at death.
	StartSize, EndSize float64

	// EmitRate is the number of particles to emit per second during continuous
	// emission. Fractional rates (e.g. 0.5 = one every two seconds) are handled
	// with a per-entity accumulator so emission timing is exact.
	// Zero pauses continuous emission without clearing live particles.
	EmitRate float64

	// MaxParticles caps the total number of live particles for this emitter.
	// Zero uses [DefaultMaxParticles].
	MaxParticles int

	// DirectionAngle is the base emission angle in radians.
	// 0 = right (+X), -math.Pi/2 = up (-Y), math.Pi/2 = down (+Y).
	DirectionAngle float64

	// Spread is the half-angle deviation from DirectionAngle in radians.
	// Each particle is emitted at DirectionAngle +/- rand*Spread.
	// Use math.Pi for a full 360-degree omnidirectional emitter.
	Spread float64

	// SpeedMin and SpeedMax define the speed range in pixels per second.
	// A random value in [SpeedMin, SpeedMax] is chosen per particle.
	SpeedMin, SpeedMax float64

	// LifeMin and LifeMax define the lifetime range in seconds.
	// A random value in [LifeMin, LifeMax] is chosen per particle.
	// Both values must be positive; zero LifeMax is treated as LifeMin.
	LifeMin, LifeMax float64

	// Gravity is a constant acceleration added to the particle's Y velocity
	// each second, in pixels per second squared.
	// Positive values pull particles downward in screen space (Y increases down).
	// Negative values push particles upward.
	Gravity float64

	// Image is an optional sprite drawn centered on each particle instead of a
	// filled circle. The image is scaled so its width equals the current
	// particle diameter (StartSize to EndSize). Nil draws circles.
	Image *ebiten.Image

	// Active controls whether the emitter spawns new particles this tick.
	// Existing particles continue aging regardless of this flag, so setting
	// Active to false lets the current burst finish naturally.
	Active bool

	// BurstCount, when greater than zero, immediately emits exactly this many
	// particles on the next update tick, then resets to zero.
	// Use for one-shot effects: explosions, hit sparks, death pops.
	// Combine with Active=false to emit a burst without continuous emission.
	BurstCount int
}

// GetType satisfies [ecs.Component].
func (e EmitterComponent) GetType() ecs.ComponentType { return ComponentType }

// DefaultMaxParticles is used when EmitterComponent.MaxParticles is zero.
const DefaultMaxParticles = 256

// Preset constructors for common emitter configurations.

// FireEmitter returns an upward fire-like emitter at (x, y).
func FireEmitter(x, y float64) EmitterComponent {
	return EmitterComponent{
		X: x, Y: y,
		Active:         true,
		EmitRate:       40,
		MaxParticles:   200,
		StartColor:     color.RGBA{R: 255, G: 200, B: 50, A: 255},
		EndColor:       color.RGBA{R: 200, G: 20, B: 0, A: 0},
		StartSize:      6,
		EndSize:        1,
		DirectionAngle: -math.Pi / 2,
		Spread:         math.Pi / 6,
		SpeedMin:       40,
		SpeedMax:       90,
		LifeMin:        0.5,
		LifeMax:        1.2,
		Gravity:        -15,
	}
}

// SmokeEmitter returns a gentle upward smoke emitter at (x, y).
func SmokeEmitter(x, y float64) EmitterComponent {
	return EmitterComponent{
		X: x, Y: y,
		Active:         true,
		EmitRate:       8,
		MaxParticles:   100,
		StartColor:     color.RGBA{R: 180, G: 180, B: 180, A: 160},
		EndColor:       color.RGBA{R: 140, G: 140, B: 140, A: 0},
		StartSize:      4,
		EndSize:        14,
		DirectionAngle: -math.Pi / 2,
		Spread:         math.Pi / 8,
		SpeedMin:       15,
		SpeedMax:       30,
		LifeMin:        1.5,
		LifeMax:        3.0,
		Gravity:        -5,
	}
}

// SparkBurst returns a one-shot omnidirectional spark explosion at (x, y).
// Set it on an entity; the burst fires once and then goes quiet.
func SparkBurst(x, y float64, count int) EmitterComponent {
	return EmitterComponent{
		X: x, Y: y,
		Active:       false,
		BurstCount:   count,
		MaxParticles: count,
		StartColor:   color.RGBA{R: 255, G: 230, B: 80, A: 255},
		EndColor:     color.RGBA{R: 255, G: 80, B: 20, A: 0},
		StartSize:    3,
		EndSize:      1,
		Spread:       math.Pi,
		SpeedMin:     60,
		SpeedMax:     180,
		LifeMin:      0.3,
		LifeMax:      0.7,
		Gravity:      60,
	}
}
