package simulation

import (
	"context"
	"log"
	"time"

	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/transport"
)

// ServerConfig controls the simulation loop behaviour.
type ServerConfig struct {
	// TickRate is the number of simulation ticks per second (e.g., 20, 30, 60).
	// Defaults to 20 if zero.
	TickRate int

	// SnapshotEvery controls how often the server sends snapshots to clients.
	// A value of 1 sends a snapshot every tick. A value of 3 sends one every 3 ticks.
	// Defaults to 1 if zero.
	SnapshotEvery int
}

func (c *ServerConfig) tickRate() int {
	if c.TickRate <= 0 {
		return 20
	}
	return c.TickRate
}

func (c *ServerConfig) snapshotEvery() int {
	if c.SnapshotEvery <= 0 {
		return 1
	}
	return c.SnapshotEvery
}

// EntitySource is a function the Server calls each tick to obtain the current
// list of simulated entities. Typically this returns level.Entities or similar.
// Using a func instead of a field avoids holding a stale pointer if the entity
// slice is reallocated during world generation.
type EntitySource func() []*ecs.Entity

// Server runs the authoritative game simulation loop in a dedicated goroutine.
//
// It owns the game world, drives SimulationSystems at a fixed tick rate,
// processes commands received from the client, and sends world state snapshots
// back via the transport.
//
// Server has no Ebitengine dependencies and no rendering logic.
//
// There are two ways to drive the server:
//
// Independent (decoupled tick rate, own goroutine):
//
//	go server.Run()
//
// Linked (caller controls tick rate, e.g. from Ebitengine's Update):
//
//	func (g *Game) Update() error {
//	    server.Step()
//	    // ... render ...
//	}
type Server struct {
	config        ServerConfig
	world         any
	entitySource  EntitySource
	transport     transport.ServerTransport
	codec         transport.SnapshotCodec
	systems       SimulationSystemManager
	stateMachine  SimulationStateMachine
	tick          uint64
	snapshotEvery int
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewServer creates a Server but does not start the loop.
// Call AddSystem, SetState, then Run (typically in a goroutine).
func NewServer(
	config ServerConfig,
	world any,
	entitySource EntitySource,
	t transport.ServerTransport,
	codec transport.SnapshotCodec,
) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		config:        config,
		world:         world,
		entitySource:  entitySource,
		transport:     t,
		codec:         codec,
		snapshotEvery: config.snapshotEvery(),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// AddSystem registers a SimulationSystem. Call before Run.
func (s *Server) AddSystem(sys SimulationSystem) {
	s.systems.AddSystem(sys)
}

// SetState sets the initial SimulationState. Call before Run.
// If not set, the server runs systems without state machine logic.
func (s *Server) SetState(state SimulationState) {
	s.stateMachine.PushState(state)
}

// Run starts the simulation loop with its own ticker at the configured TickRate.
// It blocks until Stop() is called or the state machine empties.
// Use this for independent (decoupled) tick rates. Intended to run in a goroutine.
func (s *Server) Run() {
	tickRate := s.config.tickRate()
	interval := time.Duration(float64(time.Second) / float64(tickRate))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer s.transport.Close()

	log.Printf("[simulation] server started - tick rate: %d Hz, snapshot every: %d ticks", tickRate, s.snapshotEvery)

	for {
		select {
		case <-s.ctx.Done():
			log.Printf("[simulation] server stopped at tick %d", s.tick)
			return
		case <-ticker.C:
			s.step()
			if s.stateMachine.Current() == nil {
				log.Printf("[simulation] state machine empty - server stopping at tick %d", s.tick)
				return
			}
		}
	}
}

// Step advances the simulation by exactly one tick. The caller controls when
// and how often this is called, making it suitable for linked (lockstep) mode
// where the simulation ticks in sync with Ebitengine's Update().
//
// Returns true if the state machine is still active, false if it has emptied
// (meaning the simulation is done).
func (s *Server) Step() bool {
	s.step()
	return s.stateMachine.Current() != nil
}

// Stop signals the Run goroutine to exit cleanly.
func (s *Server) Stop() {
	s.cancel()
}

// Tick is the current server tick counter (read-only from outside the loop).
func (s *Server) Tick() uint64 {
	return s.tick
}

func (s *Server) step() {
	s.tick++

	// 1. Drain and process commands from clients.
	cmds := s.transport.ReceiveCommands()
	s.stateMachine.ProcessCommands(cmds)

	// 2. Run global system pass.
	if err := s.systems.UpdateSystems(s.world); err != nil {
		log.Printf("[simulation] tick %d UpdateSystems error: %v", s.tick, err)
	}

	// 3. Run per-entity system pass.
	entities := s.entitySource()
	if err := s.systems.UpdateSystemsForEntities(s.world, entities); err != nil {
		log.Printf("[simulation] tick %d UpdateSystemsForEntities error: %v", s.tick, err)
	}

	// 4. Advance state machine.
	s.stateMachine.Tick(s.world)

	// 5. Send snapshot if it's time.
	if s.tick%uint64(s.snapshotEvery) == 0 {
		snapshot := s.codec.Encode(s.tick, entities)
		s.transport.SendSnapshot(snapshot)
	}
}
