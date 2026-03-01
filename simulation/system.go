package simulation

import "github.com/mechanical-lich/mlge/ecs"

// SimulationSystem is the server-side counterpart to ecs.SystemInterface.
//
// Implement this for any logic that belongs to the authoritative simulation:
// physics, AI decision-making, chemistry, pathfinding, collision resolution.
// Never put rendering code in a SimulationSystem.
//
// The world parameter is the same game-specific context struct you'd pass to
// ecs.SystemManager (e.g., a *GameContext or *Level). Cast it inside the method.
type SimulationSystem interface {
	// UpdateSimulation runs once per server tick for global system logic
	// (e.g., advancing a global clock, updating spatial indices).
	UpdateSimulation(world any) error

	// UpdateEntitySimulation runs once per entity per tick, for entities whose
	// component set satisfies Requires().
	UpdateEntitySimulation(world any, entity *ecs.Entity) error

	// Requires returns the set of component types an entity must have for
	// UpdateEntitySimulation to be called on it. Return nil or an empty slice
	// to receive every entity.
	Requires() []ecs.ComponentType
}

// SimulationSystemManager holds an ordered list of SimulationSystems and
// drives them each server tick. It mirrors ecs.SystemManager in structure
// but is typed to SimulationSystem so the compiler prevents accidentally
// adding a render-side system here.
type SimulationSystemManager struct {
	systems            []SimulationSystem
	cachedRequirements [][]ecs.ComponentType
}

// AddSystem appends a system to the end of the execution order.
// Systems run in the order they are added each tick.
func (m *SimulationSystemManager) AddSystem(s SimulationSystem) {
	if m.systems == nil {
		m.systems = make([]SimulationSystem, 0)
		m.cachedRequirements = make([][]ecs.ComponentType, 0)
	}
	m.systems = append(m.systems, s)
	m.cachedRequirements = append(m.cachedRequirements, s.Requires())
}

// UpdateSystems calls UpdateSimulation on every registered system.
func (m *SimulationSystemManager) UpdateSystems(world any) error {
	for _, s := range m.systems {
		if err := s.UpdateSimulation(world); err != nil {
			return err
		}
	}
	return nil
}

// UpdateSystemsForEntities iterates every system, then every entity, calling
// UpdateEntitySimulation when the entity satisfies the system's Requires().
//
// Entities with the ecs.InanimateComponentType are skipped when that global
// is set, matching the behaviour of ecs.SystemManager.UpdateSystemsForEntities.
func (m *SimulationSystemManager) UpdateSystemsForEntities(world any, entities []*ecs.Entity) error {
	for i, s := range m.systems {
		required := m.cachedRequirements[i]
		for _, entity := range entities {
			if ecs.InanimateComponentType != "" && entity.HasComponent(ecs.InanimateComponentType) {
				continue
			}
			if entity.HasComponentsSlice(required) {
				if err := s.UpdateEntitySimulation(world, entity); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
