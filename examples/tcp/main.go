package main

// TCP transport example -- the same bouncing-balls simulation as
// examples/client_server, but the server and client connect over a real TCP
// socket instead of in-process channels.
//
// Run two terminals:
//
//	terminal 1 (server):  go run examples/tcp/*.go --mode server
//	terminal 2 (client):  go run examples/tcp/*.go --mode client
//
// Optional flags:
//
//	--addr   TCP address to listen on / dial (default ":7777")
//
// The server runs headless (no window). The client opens an Ebitengine window
// and renders the balls received from the server.
//
// Key difference from the local transport example: Command.Payload and
// EntitySnapshot.Components values travel over JSON, so on the receiving side
// they arrive as map[string]interface{}. The BallCodec.Decode method handles
// this with a re-marshal helper.

import (
	"flag"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"

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
	mode := flag.String("mode", "", "server or client (required)")
	addr := flag.String("addr", ":7777", "TCP address")
	flag.Parse()

	switch *mode {
	case "server":
		runServer(*addr)
	case "client":
		runClient(*addr)
	default:
		log.Println("usage: --mode server|client [--addr :7777]")
		os.Exit(1)
	}
}

// runServer starts a headless simulation that listens for TCP clients.
func runServer(addr string) {
	srvT, err := transport.NewTCPServerTransport(addr)
	if err != nil {
		log.Fatalf("server: listen %s: %v", addr, err)
	}
	defer srvT.Close()
	log.Printf("server: listening on %s", addr)

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
	srv := simulation.NewServer(
		simulation.ServerConfig{TickRate: tickRate, SnapshotEvery: 1},
		world,
		func() []*ecs.Entity { return world.Entities },
		srvT,
		codec,
	)
	srv.AddSystem(&PhysicsSystem{DT: 1.0 / float64(tickRate)})
	srv.SetState(&MainSimState{})

	// Run until SIGINT / SIGTERM.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go srv.Run()
	log.Println("server: running. Ctrl+C to quit.")
	<-stop
	log.Println("server: stopping.")
	srv.Stop()
}

// runClient connects to the server and opens the Ebitengine render window.
func runClient(addr string) {
	cliT, err := transport.NewTCPClientTransport(addr)
	if err != nil {
		log.Fatalf("client: dial %s: %v", addr, err)
	}
	defer cliT.Close()
	log.Printf("client: connected to %s", addr)

	cliWorld := &World{Width: screenW, Height: screenH}
	codec := &BallCodec{}

	c := client.NewClient(
		cliT,
		codec,
		&MainClientState{world: cliWorld},
		cliWorld,
		func() []*ecs.Entity { return cliWorld.Entities },
		client.ClientConfig{
			ScreenWidth:  screenW,
			ScreenHeight: screenH,
			WindowTitle:  "mlge tcp example -- client",
		},
	)

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
