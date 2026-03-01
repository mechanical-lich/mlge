package particle

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/ecs"
)

// particle is an ephemeral value managed by ParticleSystem.
// It is never exposed to game code; game code interacts only via EmitterComponent.
type particle struct {
	x, y      float64
	vx, vy    float64
	life      float64 // remaining lifetime in seconds
	totalLife float64 // original lifetime used for t = 1 - life/totalLife lerp
}

// ParticleSystem manages per-entity particle emitters.
//
// It satisfies all three mlge system interfaces via structural typing:
//
//   - [ecs.SystemInterface]          (UpdateSystem / UpdateEntity)
//   - simulation.SimulationSystem    (UpdateSimulation / UpdateEntitySimulation)
//   - client.RenderSystem            (UpdateRender / UpdateEntityRender)
//
// Compile-time assertions in system_iface_test.go confirm compatibility.
// Pick the role that matches your game architecture and register accordingly.
//
// Updating particle state and rendering are separated by design.
// Call [ParticleSystem.Draw] from your state's Draw method to render particles
// after the rest of the world has been drawn.
//
// Create with [NewParticleSystem].
type ParticleSystem struct {
	// DT is the fixed timestep in seconds used to advance particle simulation.
	//   - Render-side: set to 1.0/60.0 or update each frame.
	//   - Simulation-side: set to 1.0/tickRate.
	DT float64

	// pools holds the live particle slice for each emitter entity.
	pools map[*ecs.Entity][]particle

	// accum holds the fractional emission accumulator per entity.
	// It carries the sub-integer particle debt between ticks so that low
	// emission rates (e.g. 0.5/s) still fire exactly on schedule.
	accum map[*ecs.Entity]float64
}

// NewParticleSystem returns a ready-to-use ParticleSystem with the given
// fixed timestep.
func NewParticleSystem(dt float64) *ParticleSystem {
	return &ParticleSystem{
		DT:    dt,
		pools: make(map[*ecs.Entity][]particle),
		accum: make(map[*ecs.Entity]float64),
	}
}

// Requires returns the component types an entity must have for the system to
// process it. Satisfies all three mlge system interfaces.
func (ps *ParticleSystem) Requires() []ecs.ComponentType {
	return []ecs.ComponentType{ComponentType}
}

// =============================================================================
// ecs.SystemInterface
// =============================================================================

// UpdateSystem is a no-op; all work happens per entity.
func (ps *ParticleSystem) UpdateSystem(_ any) error { return nil }

// UpdateEntity advances the particle emitter for a single entity.
func (ps *ParticleSystem) UpdateEntity(_ any, e *ecs.Entity) error { return ps.advance(e) }

// =============================================================================
// simulation.SimulationSystem
// =============================================================================

// UpdateSimulation is a no-op; all work happens per entity.
func (ps *ParticleSystem) UpdateSimulation(_ any) error { return nil }

// UpdateEntitySimulation advances the particle emitter for a single entity.
func (ps *ParticleSystem) UpdateEntitySimulation(_ any, e *ecs.Entity) error {
	return ps.advance(e)
}

// =============================================================================
// client.RenderSystem
// =============================================================================

// UpdateRender is a no-op; all work happens per entity.
func (ps *ParticleSystem) UpdateRender(_ any) error { return nil }

// UpdateEntityRender advances the particle emitter for a single entity.
func (ps *ParticleSystem) UpdateEntityRender(_ any, e *ecs.Entity) error { return ps.advance(e) }

// =============================================================================
// Core simulation
// =============================================================================

// advance is the shared simulation step called by all three interface variants.
func (ps *ParticleSystem) advance(e *ecs.Entity) error {
	comp, ok := e.Components[ComponentType]
	if !ok {
		return nil
	}
	cfg := comp.(EmitterComponent)
	dt := ps.DT

	// Age and compact existing particles.
	pool := ps.pools[e]
	live := pool[:0]
	for _, p := range pool {
		p.vy += cfg.Gravity * dt
		p.x += p.vx * dt
		p.y += p.vy * dt
		p.life -= dt
		if p.life > 0 {
			live = append(live, p)
		}
	}
	// Nil out trailing slots so GC can reclaim any backing memory.
	for i := len(live); i < len(pool); i++ {
		pool[i] = particle{}
	}
	pool = live

	max := cfg.MaxParticles
	if max <= 0 {
		max = DefaultMaxParticles
	}

	// Handle one-shot burst.
	if cfg.BurstCount > 0 {
		for i := 0; i < cfg.BurstCount && len(pool) < max; i++ {
			pool = append(pool, ps.spawnOne(cfg))
		}
		cfg.BurstCount = 0
		e.AddComponent(cfg) // write back the cleared burst count
	}

	// Handle continuous emission.
	if cfg.Active && cfg.EmitRate > 0 {
		acc := ps.accum[e] + cfg.EmitRate*dt
		toEmit := int(math.Floor(acc))
		ps.accum[e] = acc - float64(toEmit)
		for i := 0; i < toEmit && len(pool) < max; i++ {
			pool = append(pool, ps.spawnOne(cfg))
		}
	}

	ps.pools[e] = pool
	return nil
}

// spawnOne creates a single new particle using the current emitter config.
func (ps *ParticleSystem) spawnOne(cfg EmitterComponent) particle {
	angle := cfg.DirectionAngle + (rand.Float64()*2-1)*cfg.Spread
	speed := cfg.SpeedMin + rand.Float64()*(cfg.SpeedMax-cfg.SpeedMin)

	lifeMax := cfg.LifeMax
	if lifeMax <= 0 {
		lifeMax = cfg.LifeMin
	}
	life := cfg.LifeMin + rand.Float64()*(lifeMax-cfg.LifeMin)
	if life <= 0 {
		life = 1
	}

	return particle{
		x:         cfg.X,
		y:         cfg.Y,
		vx:        math.Cos(angle) * speed,
		vy:        math.Sin(angle) * speed,
		life:      life,
		totalLife: life,
	}
}

// =============================================================================
// Draw pass
// =============================================================================

// Draw renders all live particles to screen.
//
// Call this from your render state's Draw method after drawing the world:
//
//	func (s *MyState) Draw(screen *ebiten.Image) {
//	    drawWorld(screen)
//	    ps.Draw(screen)
//	}
func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for e, pool := range ps.pools {
		comp, ok := e.Components[ComponentType]
		if !ok {
			continue
		}
		cfg := comp.(EmitterComponent)
		for _, p := range pool {
			// t goes 0 (just born) to 1 (about to die).
			t := 1.0 - p.life/p.totalLife
			col := lerpColor(cfg.StartColor, cfg.EndColor, t)
			size := cfg.StartSize + (cfg.EndSize-cfg.StartSize)*t
			if size <= 0 {
				continue
			}
			if cfg.Image != nil {
				drawSprite(screen, cfg.Image, p.x, p.y, size, col)
			} else {
				vector.DrawFilledCircle(screen,
					float32(p.x), float32(p.y), float32(size), col, true)
			}
		}
	}
}

// drawSprite draws img centered at (x, y), scaled so its width equals size*2,
// and tinted by col.
func drawSprite(screen, img *ebiten.Image, x, y, size float64, col color.RGBA) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w == 0 || h == 0 {
		return
	}
	scale := (size * 2) / float64(w)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	op.ColorScale.SetR(float32(col.R) / 255)
	op.ColorScale.SetG(float32(col.G) / 255)
	op.ColorScale.SetB(float32(col.B) / 255)
	op.ColorScale.SetA(float32(col.A) / 255)
	screen.DrawImage(img, op)
}

// =============================================================================
// Lifecycle helpers
// =============================================================================

// Purge removes particle pools for entities that are no longer in the live set.
// Call this after reconciling the world entity list (e.g. after snapshot decode)
// to prevent stale pools accumulating for destroyed entities.
//
//	ps.Purge(world.Entities)
func (ps *ParticleSystem) Purge(live []*ecs.Entity) {
	liveSet := make(map[*ecs.Entity]bool, len(live))
	for _, e := range live {
		liveSet[e] = true
	}
	for e := range ps.pools {
		if !liveSet[e] {
			delete(ps.pools, e)
			delete(ps.accum, e)
		}
	}
}

// ActiveCount returns the total number of live particles across all emitters.
// Useful for debugging and performance monitoring.
func (ps *ParticleSystem) ActiveCount() int {
	n := 0
	for _, pool := range ps.pools {
		n += len(pool)
	}
	return n
}

// =============================================================================
// Helpers
// =============================================================================

// lerpColor linearly interpolates between a and b by t in [0, 1].
func lerpColor(a, b color.RGBA, t float64) color.RGBA {
	f := func(a, b uint8) uint8 {
		return uint8(float64(a) + (float64(b)-float64(a))*t)
	}
	return color.RGBA{
		R: f(a.R, b.R),
		G: f(a.G, b.G),
		B: f(a.B, b.B),
		A: f(a.A, b.A),
	}
}
