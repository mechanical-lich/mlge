package particle

import (
	"image/color"
	"math"
	"testing"

	"github.com/mechanical-lich/mlge/ecs"
)

// newEntity is a test helper that creates an entity with the given emitter.
func newEntity(cfg EmitterComponent) *ecs.Entity {
	e := &ecs.Entity{Blueprint: "test"}
	e.AddComponent(cfg)
	return e
}

// stepN calls advance n times on entity e using system ps.
func stepN(ps *ParticleSystem, e *ecs.Entity, n int) {
	for i := 0; i < n; i++ {
		_ = ps.advance(e)
	}
}

// =============================================================================
// Emission
// =============================================================================

func TestBurstEmitsThenResets(t *testing.T) {
	ps := NewParticleSystem(1.0 / 60.0)
	e := newEntity(EmitterComponent{
		Active:       false,
		BurstCount:   10,
		MaxParticles: 20,
		LifeMin:      5,
		LifeMax:      5,
		SpeedMin:     0,
		SpeedMax:     0,
	})

	if err := ps.advance(e); err != nil {
		t.Fatalf("advance: %v", err)
	}
	if got := ps.ActiveCount(); got != 10 {
		t.Errorf("after burst: ActiveCount = %d, want 10", got)
	}

	// BurstCount must be cleared on the component after firing.
	comp := e.Components[ComponentType].(EmitterComponent)
	if comp.BurstCount != 0 {
		t.Errorf("BurstCount not reset: got %d, want 0", comp.BurstCount)
	}

	// Second advance must not emit more particles.
	_ = ps.advance(e)
	if got := ps.ActiveCount(); got != 10 {
		t.Errorf("second advance: ActiveCount = %d, want 10 (no extra emission)", got)
	}
}

func TestContinuousEmissionRate(t *testing.T) {
	// 60 particles/sec at dt=1/60 => exactly 1 particle per tick.
	ps := NewParticleSystem(1.0 / 60.0)
	e := newEntity(EmitterComponent{
		Active:       true,
		EmitRate:     60,
		MaxParticles: 1000,
		LifeMin:      100,
		LifeMax:      100,
		SpeedMin:     0,
		SpeedMax:     0,
	})

	stepN(ps, e, 10)

	if got := ps.ActiveCount(); got != 10 {
		t.Errorf("after 10 ticks at 60/s: ActiveCount = %d, want 10", got)
	}
}

func TestFractionalEmissionRate(t *testing.T) {
	// 0.5/s at dt=1s => nothing after tick 1, one particle after tick 2.
	ps := NewParticleSystem(1.0)
	e := newEntity(EmitterComponent{
		Active:       true,
		EmitRate:     0.5,
		MaxParticles: 100,
		LifeMin:      100,
		LifeMax:      100,
		SpeedMin:     0,
		SpeedMax:     0,
	})

	_ = ps.advance(e)
	if ps.ActiveCount() != 0 {
		t.Errorf("after 1 tick (0.5 acc): want 0, got %d", ps.ActiveCount())
	}
	_ = ps.advance(e)
	if ps.ActiveCount() != 1 {
		t.Errorf("after 2 ticks (1.0 acc): want 1, got %d", ps.ActiveCount())
	}
}

func TestInactiveEmitterDoesNotEmit(t *testing.T) {
	ps := NewParticleSystem(1.0 / 60.0)
	e := newEntity(EmitterComponent{
		Active:       false,
		EmitRate:     100,
		MaxParticles: 100,
		LifeMin:      100,
		LifeMax:      100,
	})

	stepN(ps, e, 60)

	if got := ps.ActiveCount(); got != 0 {
		t.Errorf("inactive emitter: ActiveCount = %d, want 0", got)
	}
}

func TestMaxParticlesCapIsRespected(t *testing.T) {
	const cap = 5
	ps := NewParticleSystem(1.0 / 60.0)
	e := newEntity(EmitterComponent{
		Active:       true,
		EmitRate:     1000,
		MaxParticles: cap,
		LifeMin:      100,
		LifeMax:      100,
		SpeedMin:     0,
		SpeedMax:     0,
	})

	stepN(ps, e, 3)

	if got := ps.ActiveCount(); got > cap {
		t.Errorf("ActiveCount %d exceeds MaxParticles %d", got, cap)
	}
}

func TestDefaultMaxParticlesAppliedWhenZero(t *testing.T) {
	ps := NewParticleSystem(1.0 / 60.0)
	e := newEntity(EmitterComponent{
		Active:       true,
		EmitRate:     100000,
		MaxParticles: 0,
		LifeMin:      100,
		LifeMax:      100,
		SpeedMin:     0,
		SpeedMax:     0,
	})

	_ = ps.advance(e)

	if got := ps.ActiveCount(); got > DefaultMaxParticles {
		t.Errorf("ActiveCount %d exceeds DefaultMaxParticles %d", got, DefaultMaxParticles)
	}
}

// =============================================================================
// Aging and movement
// =============================================================================

func TestParticlesAgeAndDie(t *testing.T) {
	// dt=1s, lifetime=2s => particles survive 2 ticks then die on tick 3.
	ps := NewParticleSystem(1.0)
	e := newEntity(EmitterComponent{
		Active:       true,
		EmitRate:     5,
		MaxParticles: 100,
		LifeMin:      2,
		LifeMax:      2,
		SpeedMin:     0,
		SpeedMax:     0,
	})

	_ = ps.advance(e)
	afterTick1 := ps.ActiveCount()
	if afterTick1 == 0 {
		t.Fatal("expected particles after tick 1")
	}

	// Stop emitting so we only watch the existing particles die.
	cfg := e.Components[ComponentType].(EmitterComponent)
	cfg.Active = false
	e.AddComponent(cfg)

	_ = ps.advance(e)
	if ps.ActiveCount() != afterTick1 {
		t.Errorf("after tick 2: count changed from %d to %d before expiry",
			afterTick1, ps.ActiveCount())
	}

	_ = ps.advance(e)
	if ps.ActiveCount() != 0 {
		t.Errorf("after tick 3: expected 0, got %d", ps.ActiveCount())
	}
}

func TestParticleMovesWithVelocity(t *testing.T) {
	// One particle heading right at 100 px/s; after one dt=1s tick x must be 100.
	ps := NewParticleSystem(1.0)
	e := newEntity(EmitterComponent{
		X:              0,
		Y:              0,
		Active:         false,
		BurstCount:     1,
		MaxParticles:   1,
		LifeMin:        10,
		LifeMax:        10,
		DirectionAngle: 0,
		Spread:         0,
		SpeedMin:       100,
		SpeedMax:       100,
		Gravity:        0,
	})

	_ = ps.advance(e) // spawns particle
	_ = ps.advance(e) // physics tick

	pool := ps.pools[e]
	if len(pool) != 1 {
		t.Fatalf("expected 1 particle, got %d", len(pool))
	}
	p := pool[0]

	const want, tol = 100.0, 1e-6
	if math.Abs(p.x-want) > tol || math.Abs(p.y) > tol {
		t.Errorf("particle position = (%.4f, %.4f), want (100, 0)", p.x, p.y)
	}
}

func TestGravityAcceleratesParticle(t *testing.T) {
	ps := NewParticleSystem(1.0)
	e := newEntity(EmitterComponent{
		Active:       false,
		BurstCount:   1,
		MaxParticles: 1,
		LifeMin:      10,
		LifeMax:      10,
		SpeedMin:     0,
		SpeedMax:     0,
		Gravity:      20,
	})

	_ = ps.advance(e) // spawn: vy=0
	_ = ps.advance(e)
	p1 := ps.pools[e][0]
	_ = ps.advance(e)
	p2 := ps.pools[e][0]

	if p2.y <= p1.y {
		t.Errorf("gravity: y did not increase: p1.y=%.2f p2.y=%.2f", p1.y, p2.y)
	}
	if p2.vy <= p1.vy {
		t.Errorf("gravity: vy did not increase: p1.vy=%.2f p2.vy=%.2f", p1.vy, p2.vy)
	}
}

// =============================================================================
// No component
// =============================================================================

func TestAdvanceNoopOnMissingComponent(t *testing.T) {
	ps := NewParticleSystem(1.0 / 60.0)
	e := &ecs.Entity{Blueprint: "no-emitter"}

	if err := ps.advance(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.ActiveCount() != 0 {
		t.Errorf("expected 0 particles for entity without EmitterComponent")
	}
}

// =============================================================================
// Purge
// =============================================================================

func TestPurgeRemovesStaleEntities(t *testing.T) {
	ps := NewParticleSystem(1.0)

	alive := newEntity(EmitterComponent{
		Active: true, EmitRate: 10, MaxParticles: 100,
		LifeMin: 100, LifeMax: 100,
	})
	dead := newEntity(EmitterComponent{
		Active: true, EmitRate: 10, MaxParticles: 100,
		LifeMin: 100, LifeMax: 100,
	})

	_ = ps.advance(alive)
	_ = ps.advance(dead)

	if ps.ActiveCount() == 0 {
		t.Fatal("expected particles before purge")
	}

	ps.Purge([]*ecs.Entity{alive})

	if _, ok := ps.pools[dead]; ok {
		t.Errorf("pool for dead entity still present after Purge")
	}
	if _, ok := ps.accum[dead]; ok {
		t.Errorf("accum for dead entity still present after Purge")
	}
	if _, ok := ps.pools[alive]; !ok {
		t.Errorf("pool for alive entity removed by Purge")
	}
}

func TestPurgeEmptyLiveSet(t *testing.T) {
	ps := NewParticleSystem(1.0 / 60.0)
	e := newEntity(EmitterComponent{
		Active: true, EmitRate: 10, MaxParticles: 100,
		LifeMin: 100, LifeMax: 100,
	})
	_ = ps.advance(e)

	ps.Purge(nil)

	if ps.ActiveCount() != 0 {
		t.Errorf("after Purge(nil): expected 0, got %d", ps.ActiveCount())
	}
}

// =============================================================================
// ActiveCount
// =============================================================================

func TestActiveCountAcrossMultipleEntities(t *testing.T) {
	ps := NewParticleSystem(1.0)
	e1 := newEntity(EmitterComponent{
		Active: false, BurstCount: 3, MaxParticles: 10,
		LifeMin: 100, LifeMax: 100,
	})
	e2 := newEntity(EmitterComponent{
		Active: false, BurstCount: 5, MaxParticles: 10,
		LifeMin: 100, LifeMax: 100,
	})

	_ = ps.advance(e1)
	_ = ps.advance(e2)

	if got := ps.ActiveCount(); got != 8 {
		t.Errorf("ActiveCount = %d, want 8 (3+5)", got)
	}
}

// =============================================================================
// Requires
// =============================================================================

func TestRequiresReturnsEmitterComponentType(t *testing.T) {
	ps := NewParticleSystem(1.0 / 60.0)
	req := ps.Requires()
	if len(req) != 1 || req[0] != ComponentType {
		t.Errorf("Requires() = %v, want [%s]", req, ComponentType)
	}
}

// =============================================================================
// lerpColor
// =============================================================================

func TestLerpColorAtZero(t *testing.T) {
	a := color.RGBA{R: 100, G: 50, B: 200, A: 255}
	b := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	got := lerpColor(a, b, 0)
	if got != a {
		t.Errorf("lerpColor(a, b, 0) = %v, want %v", got, a)
	}
}

func TestLerpColorAtOne(t *testing.T) {
	a := color.RGBA{R: 100, G: 50, B: 200, A: 255}
	b := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	got := lerpColor(a, b, 1)
	if got != b {
		t.Errorf("lerpColor(a, b, 1) = %v, want %v", got, b)
	}
}

func TestLerpColorAtHalf(t *testing.T) {
	a := color.RGBA{R: 200, G: 100, B: 0, A: 255}
	b := color.RGBA{R: 0, G: 0, B: 200, A: 0}
	got := lerpColor(a, b, 0.5)
	// Each channel midpoint: R=100, G=50, B=100, A=127
	if got.R != 100 || got.G != 50 || got.B != 100 || got.A != 127 {
		t.Errorf("lerpColor at 0.5 = %v, want {100 50 100 127}", got)
	}
}

// =============================================================================
// Presets (smoke test)
// =============================================================================

func TestPresetsReturnNonZeroEmitters(t *testing.T) {
	tests := []struct {
		name string
		cfg  EmitterComponent
	}{
		{"FireEmitter", FireEmitter(0, 0)},
		{"SmokeEmitter", SmokeEmitter(0, 0)},
		{"SparkBurst", SparkBurst(0, 0, 20)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.LifeMin <= 0 {
				t.Errorf("%s: LifeMin <= 0", tt.name)
			}
			if tt.cfg.MaxParticles <= 0 {
				t.Errorf("%s: MaxParticles <= 0", tt.name)
			}
			ps := NewParticleSystem(1.0 / 60.0)
			e := newEntity(tt.cfg)
			if err := ps.advance(e); err != nil {
				t.Errorf("%s: advance error: %v", tt.name, err)
			}
		})
	}
}
