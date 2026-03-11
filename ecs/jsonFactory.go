package ecs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ComponentConstructor is a function that creates a new zero-value instance of a component.
type ComponentConstructor func() Component

// JSONFactory manages JSON-based entity blueprints and component construction.
type JSONFactory struct {
	blueprints map[string]map[string]map[string]interface{}
	registry   map[string]ComponentConstructor
}

// NewJSONFactory creates a new JSON-based entity factory.
func NewJSONFactory() *JSONFactory {
	return &JSONFactory{
		blueprints: make(map[string]map[string]map[string]interface{}),
		registry:   make(map[string]ComponentConstructor),
	}
}

// RegisterComponent registers a component constructor by name.
// The constructor should return a pointer to a new zero-value component.
func (f *JSONFactory) RegisterComponent(name string, constructor ComponentConstructor) {
	f.registry[name] = constructor
}

// LoadBlueprintsFromDir loads all JSON blueprint files from a directory.
// Each file should contain a map of blueprint names to component definitions.
func (f *JSONFactory) LoadBlueprintsFromDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read blueprint dir %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if err := f.LoadBlueprintsFromFile(path); err != nil {
			return err
		}
	}
	return nil
}

// LoadBlueprintsFromFile loads blueprints from a single JSON file.
// The file should contain a map where keys are blueprint names and values
// are maps of component names to component data.
//
// Example JSON:
//
//	{
//	    "player": {
//	        "HealthComponent": {"MaxHealth": 100, "Health": 100},
//	        "AppearanceComponent": {"SpriteName": "player_idle"}
//	    }
//	}
func (f *JSONFactory) LoadBlueprintsFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", path, err)
	}

	var fileBPs map[string]map[string]map[string]interface{}
	if err := json.Unmarshal(data, &fileBPs); err != nil {
		return fmt.Errorf("failed to parse %s: %w", path, err)
	}

	for name, bp := range fileBPs {
		f.blueprints[name] = bp
	}
	return nil
}

// BlueprintExists checks if a blueprint name is registered.
func (f *JSONFactory) BlueprintExists(name string) bool {
	_, ok := f.blueprints[name]
	return ok
}

// GetBlueprintNames returns all registered blueprint names.
func (f *JSONFactory) GetBlueprintNames() []string {
	names := make([]string, 0, len(f.blueprints))
	for name := range f.blueprints {
		names = append(names, name)
	}
	return names
}

// Create creates an entity from a named blueprint.
// Components are created by looking up their constructors in the registry,
// then populating them via JSON unmarshalling.
func (f *JSONFactory) Create(name string) (*Entity, error) {
	blueprint, ok := f.blueprints[name]
	if !ok {
		return nil, fmt.Errorf("no blueprint found: %s", name)
	}

	entity := &Entity{}
	entity.Blueprint = name

	for compName, params := range blueprint {
		comp, err := f.CreateComponent(compName, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create component %s for %s: %w", compName, name, err)
		}
		entity.AddComponent(comp)
	}

	return entity, nil
}

// CreateWithCallback creates an entity and calls the callback for each component
// before adding it to the entity. This allows custom initialization.
func (f *JSONFactory) CreateWithCallback(name string, callback func(comp Component) error) (*Entity, error) {
	blueprint, ok := f.blueprints[name]
	if !ok {
		return nil, fmt.Errorf("no blueprint found: %s", name)
	}

	entity := &Entity{}
	entity.Blueprint = name

	for compName, params := range blueprint {
		comp, err := f.CreateComponent(compName, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create component %s for %s: %w", compName, name, err)
		}
		if callback != nil {
			if err := callback(comp); err != nil {
				return nil, fmt.Errorf("callback failed for component %s: %w", compName, err)
			}
		}
		entity.AddComponent(comp)
	}

	return entity, nil
}

// CreateComponent creates a component instance from registry and populates it with JSON data.
func (f *JSONFactory) CreateComponent(name string, data map[string]interface{}) (Component, error) {
	constructor, ok := f.registry[name]
	if !ok {
		return nil, fmt.Errorf("component %s not registered", name)
	}

	comp := constructor()
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, comp); err != nil {
		return nil, err
	}

	return comp, nil
}
