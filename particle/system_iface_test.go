package particle_test

import (
	"github.com/mechanical-lich/mlge/client"
	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/particle"
	"github.com/mechanical-lich/mlge/simulation"
)

// Compile-time checks that ParticleSystem satisfies all three mlge system interfaces.
var (
	_ ecs.SystemInterface           = (*particle.ParticleSystem)(nil)
	_ simulation.SimulationSystem   = (*particle.ParticleSystem)(nil)
	_ client.RenderSystem           = (*particle.ParticleSystem)(nil)
)
