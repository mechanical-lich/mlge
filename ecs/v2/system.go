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
	systems []SystemInterface
}

func (s *SystemManager) AddSystem(system SystemInterface) {
	if s.systems == nil {
		s.systems = make([]SystemInterface, 0)
	}

	s.systems = append(s.systems, system)
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
	for system := range s.systems {
		if entity.HasComponents(s.systems[system].Requires()...) {
			err := s.systems[system].UpdateEntity(world, entity)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
