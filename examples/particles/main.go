package main

// Particle system example.
//
// Demonstrates [particle.ParticleSystem] with three emitter presets:
//   - Fire  -- continuous flame at the left post
//   - Smoke -- continuous smoke at the right post
//   - Spark -- one-shot burst on every mouse click
//
// No simulation server is needed: the particle system is advanced directly in
// Update() and rendered in Draw(), which is the typical client-only pattern.
//
// Run:
//
//	go run examples/particles/main.go

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/particle"
)

const (
	screenW = 640
	screenH = 480
)

var bgColor = color.RGBA{15, 15, 25, 255}

// Game holds the three emitter entities and the shared ParticleSystem.
type Game struct {
	ps    *particle.ParticleSystem
	fire  *ecs.Entity
	smoke *ecs.Entity
	spark *ecs.Entity // reused for every click; only BurstCount changes
}

func NewGame() *Game {
	ps := particle.NewParticleSystem(1.0 / 60.0)

	// Fire emitter – left side, midway down.
	fire := &ecs.Entity{Blueprint: "fire"}
	fire.AddComponent(particle.FireEmitter(160, 320))

	// Smoke emitter – right side, same height.
	smoke := &ecs.Entity{Blueprint: "smoke"}
	smoke.AddComponent(particle.SmokeEmitter(480, 320))

	// Spark entity – starts inactive; BurstCount is set on each click.
	spark := &ecs.Entity{Blueprint: "spark"}
	spark.AddComponent(particle.SparkBurst(0, 0, 0))

	return &Game{
		ps:    ps,
		fire:  fire,
		smoke: smoke,
		spark: spark,
	}
}

func (g *Game) Update() error {
	// Trigger a spark burst on left-click.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cfg := particle.SparkBurst(float64(x), float64(y), 40)
		g.spark.AddComponent(cfg)
	}

	_ = g.ps.UpdateEntity(nil, g.fire)
	_ = g.ps.UpdateEntity(nil, g.smoke)
	_ = g.ps.UpdateEntity(nil, g.spark)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bgColor)
	g.ps.Draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"particles: %d\nFPS: %.0f\nClick anywhere for sparks",
		g.ps.ActiveCount(), ebiten.ActualFPS(),
	))
}

func (g *Game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("mlge particle example")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
