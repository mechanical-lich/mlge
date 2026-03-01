package client

import "github.com/mechanical-lich/mlge/ecs"

// RenderSystem is the client-side counterpart to simulation.SimulationSystem.
//
// Implement this for visual logic that must run every Ebitengine frame:
//   - Sprite animation frame advancement
//   - Position interpolation between snapshots
//   - Particle effects and screen-space VFX
//   - Audio cue triggers based on visible state changes
//
// RenderSystems operate on the client's local (non-authoritative) world copy
// and must never mutate simulation state directly.
type RenderSystem interface {
	// UpdateRender is called once per frame for global render logic
	// (e.g., advancing global animation timers, updating camera).
	UpdateRender(world any) error

	// UpdateEntityRender is called once per frame per entity whose component
	// set satisfies Requires().
	UpdateEntityRender(world any, entity *ecs.Entity) error

	// Requires returns the component types an entity must have for
	// UpdateEntityRender to be called on it.
	Requires() []ecs.ComponentType
}

// RenderSystemManager holds an ordered list of RenderSystems and
// drives them each Ebitengine frame.
type RenderSystemManager struct {
	systems            []RenderSystem
	cachedRequirements [][]ecs.ComponentType
}

// AddSystem appends a system to the execution order.
func (m *RenderSystemManager) AddSystem(s RenderSystem) {
	if m.systems == nil {
		m.systems = make([]RenderSystem, 0)
		m.cachedRequirements = make([][]ecs.ComponentType, 0)
	}
	m.systems = append(m.systems, s)
	m.cachedRequirements = append(m.cachedRequirements, s.Requires())
}

// UpdateSystems calls UpdateRender on every registered system.
func (m *RenderSystemManager) UpdateSystems(world any) error {
	for _, s := range m.systems {
		if err := s.UpdateRender(world); err != nil {
			return err
		}
	}
	return nil
}

// UpdateSystemsForEntities iterates every system, then every entity, calling
// UpdateEntityRender when the entity satisfies the system's Requires().
func (m *RenderSystemManager) UpdateSystemsForEntities(world any, entities []*ecs.Entity) error {
	for i, s := range m.systems {
		required := m.cachedRequirements[i]
		for _, entity := range entities {
			if ecs.InanimateComponentType != "" && entity.HasComponent(ecs.InanimateComponentType) {
				continue
			}
			if entity.HasComponentsSlice(required) {
				if err := s.UpdateEntityRender(world, entity); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
