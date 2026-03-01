package main

import (
	"log"
	"math/rand/v2"

	"github.com/mechanical-lich/mlge/client"
	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/simulation"
	"github.com/mechanical-lich/mlge/transport"
)

const (
	screenW   = 640
	screenH   = 480
	ballCount = 12
	tickRate  = 20 // server Hz
)

func main() {
	// Shared transport
	srvT, cliT := transport.NewLocalTransport()

	// Server-side world and entities
	srvWorld := &World{Width: screenW, Height: screenH}
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
		srvWorld.Entities = append(srvWorld.Entities, e)
	}

	codec := &BallCodec{}

	srv := simulation.NewServer(
		simulation.ServerConfig{TickRate: tickRate, SnapshotEvery: 1},
		srvWorld,
		func() []*ecs.Entity { return srvWorld.Entities },
		srvT,
		codec,
	)
	srv.AddSystem(&PhysicsSystem{DT: 1.0 / float64(tickRate)})
	srv.SetState(&MainSimState{})
	go srv.Run()

	// Client-side world (non-authoritative, updated from snapshots)
	cliWorld := &World{Width: screenW, Height: screenH}
	initialState := &MainClientState{world: cliWorld}

	c := client.NewClient(
		cliT,
		codec,
		initialState,
		cliWorld,
		func() []*ecs.Entity { return cliWorld.Entities },
		client.ClientConfig{
			ScreenWidth:  screenW,
			ScreenHeight: screenH,
			WindowTitle:  "mlge client/server example",
		},
	)

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
	srv.Stop()
}
