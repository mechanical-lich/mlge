package particle

// Package particle provides a generic particle emitter component and system
// for mlge games.
//
// # Overview
//
// Attach a [EmitterComponent] to any ECS entity to configure a particle
// emitter at that entity's position. Create a [ParticleSystem], add it to
// whatever system manager you prefer, and call [ParticleSystem.Draw] from
// your state's Draw method.
//
// # System interfaces
//
// [ParticleSystem] satisfies all three mlge system interfaces simultaneously.
// Pick the role that matches the architecture of the game:
//
//   - As an [ecs.SystemInterface] -- pass to [ecs.SystemManager.AddSystem].
//     Update ticks at the rate driven by the caller.
//
//   - As a [simulation.SimulationSystem] -- pass to
//     [simulation.SimulationSystemManager.AddSystem]. Particle state advances on
//     the authoritative server tick. Useful when particle state must be included//     in deterministic replays or sent over the network.
//
//   - As a [client.RenderSystem] -- pass to [client.Client.AddRenderSystem].
//     Particle state advances every Ebitengine frame. This is the most common
//     choice for pure visual effects that do not need server authority.
//
// # Draw pass
//
// Updating particle state and drawing are intentionally separated. After
// registering the system with a manager, call [ParticleSystem.Draw] from
// your render state's Draw method:
//
//	func (s *MyState) Draw(screen *ebiten.Image) {
//	    drawWorld(screen)
//	    ps.Draw(screen)
//	}
//
// # Example
//
//	ps := particle.NewParticleSystem(1.0 / 60.0)
//
//	fireEntity := &ecs.Entity{Blueprint: "fire"}
//	fireEntity.AddComponent(particle.FireEmitter(320, 240))
//
//	// Register as a render system (most common).
//	client.AddRenderSystem(ps)
//
//	// In state Draw:
//	ps.Draw(screen)
