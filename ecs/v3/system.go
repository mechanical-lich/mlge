package ecs

// SystemInterface - interface that represents a system, world is an interface and should be cast to whatever data
// structure the game is currently using or that the system cares about.
type SystemInterface interface {
	UpdateSystem(data any) error
	UpdateEntity(data any, entity *Entity) error
	Requires() []ComponentType
}

// SystemManager - contains a list of systems and is responsible for calling their update functions on entities.
type SystemManager struct {
	systems            []SystemInterface
	cachedRequirements [][]ComponentType // Cache Requires() results
}

func (s *SystemManager) AddSystem(system SystemInterface) {
	if s.systems == nil {
		s.systems = make([]SystemInterface, 0)
		s.cachedRequirements = make([][]ComponentType, 0)
	}

	s.systems = append(s.systems, system)
	s.cachedRequirements = append(s.cachedRequirements, system.Requires())
}

func (s *SystemManager) UpdateSystems(world any) error {
	for system := range s.systems {
		err := s.systems[system].UpdateSystem(world)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateSystemsForEntity - Iterates through the systems for the specific entity
func (s *SystemManager) UpdateSystemsForEntity(world any, entity *Entity) error {
	for i, system := range s.systems {
		// Use cached requirements instead of calling Requires() each time
		if entity.HasComponentsSlice(s.cachedRequirements[i]) {
			err := system.UpdateEntity(world, entity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// InanimateComponentType can be set by the game to skip inanimate entities
var InanimateComponentType ComponentType = -1

func (s *SystemManager) UpdateSystemsForEntities(world any, entities []*Entity) error {
	for i, system := range s.systems {
		required := s.cachedRequirements[i]
		for _, entity := range entities {
			if InanimateComponentType >= 0 && entity.HasComponent(InanimateComponentType) {
				continue // Skip inanimate entities
			}
			if entity.HasComponentsSlice(required) {
				if err := system.UpdateEntity(world, entity); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
