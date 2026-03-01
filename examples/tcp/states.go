package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/client"
	"github.com/mechanical-lich/mlge/simulation"
	"github.com/mechanical-lich/mlge/transport"
)

var colorBG = color.RGBA{20, 20, 30, 255}

// ---- Server-side state ------------------------------------------------------

// MainSimState is a trivial SimulationState that runs indefinitely.
type MainSimState struct{}

var _ simulation.SimulationState = (*MainSimState)(nil)

func (s *MainSimState) Tick(_ any) simulation.SimulationState { return nil }
func (s *MainSimState) ProcessCommand(_ *transport.Command)   {}
func (s *MainSimState) Done() bool                            { return false }

// ---- Client-side state ------------------------------------------------------

// MainClientState renders the balls received from the server via snapshots.
type MainClientState struct {
	world      *World
	latestTick uint64
}

var _ client.ClientState = (*MainClientState)(nil)

func (s *MainClientState) Done() bool { return false }

func (s *MainClientState) Update(snapshot *transport.Snapshot) client.ClientState {
	if snapshot != nil {
		s.latestTick = snapshot.Tick
	}
	return nil
}

func (s *MainClientState) Draw(screen *ebiten.Image) {
	screen.Fill(colorBG)

	for _, e := range s.world.Entities {
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
		"TCP MODE\nserver tick: %d\nballs: %d\nFPS: %.0f",
		s.latestTick, len(s.world.Entities), ebiten.ActualFPS(),
	))
}
