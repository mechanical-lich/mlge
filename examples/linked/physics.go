package main

import (
	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/simulation"
)

// PhysicsSystem advances ball positions and bounces them off walls.
type PhysicsSystem struct {
	DT float64
}

var _ simulation.SimulationSystem = (*PhysicsSystem)(nil)

func (s *PhysicsSystem) Requires() []ecs.ComponentType {
	return []ecs.ComponentType{TypePosition, TypeVelocity}
}

func (s *PhysicsSystem) UpdateSimulation(_ any) error { return nil }

func (s *PhysicsSystem) UpdateEntitySimulation(world any, entity *ecs.Entity) error {
	w := world.(*World)

	pos := entity.Components[TypePosition].(PositionComponent)
	vel := entity.Components[TypeVelocity].(VelocityComponent)

	pos.X += vel.VX * s.DT
	pos.Y += vel.VY * s.DT

	if pos.X-ballRadius < 0 {
		pos.X = ballRadius
		vel.VX = -vel.VX
	}
	if pos.X+ballRadius > w.Width {
		pos.X = w.Width - ballRadius
		vel.VX = -vel.VX
	}
	if pos.Y-ballRadius < 0 {
		pos.Y = ballRadius
		vel.VY = -vel.VY
	}
	if pos.Y+ballRadius > w.Height {
		pos.Y = w.Height - ballRadius
		vel.VY = -vel.VY
	}

	entity.AddComponent(pos)
	entity.AddComponent(vel)
	return nil
}
