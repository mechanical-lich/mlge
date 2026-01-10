package ecs

// Entity - represents an entity which essentially is an array of components built from a blueprint.
type Entity struct {
	Components map[ComponentType]Component
	Blueprint  string
}

// AddComponent - Adds the provided component to the entity.
func (entity *Entity) AddComponent(c Component) {
	if entity.Components == nil {
		entity.Components = make(map[ComponentType]Component)
	}

	entity.Components[c.GetType()] = c
}

// HasComponent - Returns if the entity has the component
func (entity *Entity) HasComponent(name ComponentType) bool {
	return entity.Components[name] != nil
}

// HasComponents - takes component types and returns if entity has all of them.
func (entity *Entity) HasComponents(names ...ComponentType) bool {
	for _, name := range names {
		if entity.Components[name] == nil {
			return false
		}
	}

	return true
}

// HasComponentsSlice - same as HasComponents but takes a slice to avoid variadic allocation
func (entity *Entity) HasComponentsSlice(names []ComponentType) bool {
	for _, name := range names {
		if entity.Components[name] == nil {
			return false
		}
	}
	return true
}

// GetComponent - Gets the component as a component interface.
func (entity *Entity) GetComponent(name ComponentType) Component {
	return entity.Components[name]
}

// RemoveComponent - Removes the component from the entity.
func (entity *Entity) RemoveComponent(name ComponentType) {
	if entity.Components == nil {
		entity.Components = make(map[ComponentType]Component)
	}

	delete(entity.Components, name)
}
