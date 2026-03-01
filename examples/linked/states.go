package main

import (
	"github.com/mechanical-lich/mlge/simulation"
	"github.com/mechanical-lich/mlge/transport"
)

// MainSimState is a trivial SimulationState that runs forever.
type MainSimState struct{}

var _ simulation.SimulationState = (*MainSimState)(nil)

func (s *MainSimState) Tick(_ any) simulation.SimulationState { return nil }
func (s *MainSimState) ProcessCommand(_ *transport.Command)   {}
func (s *MainSimState) Done() bool                            { return false }
