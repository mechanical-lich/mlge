package main

// Linked mode example -- the simulation ticks in lockstep with Ebitengine's
// Update() rather than running in a separate goroutine. This is the simplest
// architecture: one Game struct owns the server, the world, and the rendering.
//
// Compare with examples/client_server which uses `go srv.Run()` for a
// decoupled (independent) tick rate.

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/simulation"
	"github.com/mechanical-lich/mlge/transport"
)

const (
	screenW   = 640
	screenH   = 480
	ballCount = 12
)

var colorBG = color.RGBA{20, 20, 30, 255}

// Game implements ebiten.Game and drives both the simulation and the rendering.
// In linked mode the server never launches a goroutine -- the caller (Update)
// invokes srv.Step() directly each frame.
type Game struct {
	srv   *simulation.Server
	cliT  transport.ClientTransport
	codec transport.SnapshotCodec
	world *World
	tick  uint64
}

func (g *Game) Update() error {
	// Advance the simulation by one tick, synchronously.
	if !g.srv.Step() {
		// State machine emptied -- simulation is done.
		return ebiten.Termination
	}

	// Pull the latest snapshot produced by Step() and apply it.
	snap := g.cliT.ReceiveSnapshot()
	if snap != nil {
		g.codec.Decode(snap, g.world)
		g.tick = snap.Tick
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colorBG)

	for _, e := range g.world.Entities {
		posC, hasPos := e.Components[TypePosition]
		colC, hasCol := e.Components[TypeColor]
		if !hasPos || !hasCol {
			continue
		}
		pos := posC.(PositionComponent)
		col := colC.(ColorComponent)

		vector.DrawFilledCircle(
			screen,
			float32(pos.X), float32(pos.Y),
			ballRadius,
			col.RGBA,
			true,
		)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"LINKED MODE\nserver tick: %d\nballs: %d\nFPS: %.0f",
		g.tick, len(g.world.Entities), ebiten.ActualFPS(),
	))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func main() {
	// Transport -- still used so the architecture stays consistent.
	// In linked mode the channels just act as a single-item hand-off.
	srvT, cliT := transport.NewLocalTransport()

	// Server-side world.
	world := &World{Width: screenW, Height: screenH}
	for i := range ballCount {
		e := &ecs.Entity{Blueprint: "ball"}
		e.AddComponent(IDComponent{ID: ballID(i)})
		e.AddComponent(PositionComponent{
			X: rand.Float64() * screenW,
			Y: rand.Float64() * screenH,
		})
		e.AddComponent(VelocityComponent{
			VX: (rand.Float64()*2 - 1) * 120,
			VY: (rand.Float64()*2 - 1) * 120,
		})
		e.AddComponent(ColorComponent{ballColor(i)})
		world.Entities = append(world.Entities, e)
	}

	codec := &BallCodec{}

	// Ebitengine targets 60 TPS by default, so the effective simulation rate
	// is 60 Hz -- matching the frame rate exactly.
	srv := simulation.NewServer(
		simulation.ServerConfig{TickRate: 60, SnapshotEvery: 1},
		world,
		func() []*ecs.Entity { return world.Entities },
		srvT,
		codec,
	)
	srv.AddSystem(&PhysicsSystem{DT: 1.0 / 60.0})
	srv.SetState(&MainSimState{})

	// Client-side world (receives decoded snapshots).
	cliWorld := &World{Width: screenW, Height: screenH}

	game := &Game{
		srv:   srv,
		cliT:  cliT,
		codec: codec,
		world: cliWorld,
	}

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("mlge linked mode example")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
