package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockSystem - a mock implementation of SystemInterface for testing purposes.
type MockSystem struct {
	UpdateSystemCalled bool
	UpdateEntityCalled bool
	RequiresCalled     bool
	RequiredComponents []ComponentType
	UpdateSystemData   interface{}
	UpdateEntityData   interface{}
	UpdateEntityEntity *Entity
}

func (m *MockSystem) UpdateSystem(data interface{}) error {
	m.UpdateSystemCalled = true
	m.UpdateSystemData = data
	return nil
}

func (m *MockSystem) UpdateEntity(data interface{}, entity *Entity) error {
	m.UpdateEntityCalled = true
	m.UpdateEntityData = data
	m.UpdateEntityEntity = entity
	return nil
}

func (m *MockSystem) Requires() []ComponentType {
	m.RequiresCalled = true
	return m.RequiredComponents
}

// TestAddSystem - Tests if a system is added to the SystemManager correctly.
func TestAddSystem(t *testing.T) {
	assert := assert.New(t)
	sm := &SystemManager{}
	mockSystem := &MockSystem{}

	sm.AddSystem(mockSystem)

	assert.Equal(1, len(sm.systems), "Expected exactly one system in the manager")
	assert.Equal(mockSystem, sm.systems[0], "Expected the added system to be the same as the mock system")
}

// TestUpdateSystems - Tests if UpdateSystems calls UpdateSystem on all systems.
func TestUpdateSystems(t *testing.T) {
	assert := assert.New(t)
	sm := &SystemManager{}
	mockSystem1 := &MockSystem{}
	mockSystem2 := &MockSystem{}

	sm.AddSystem(mockSystem1)
	sm.AddSystem(mockSystem2)

	err := sm.UpdateSystems("testData")

	assert.Nil(err, "Expected no error from UpdateSystems")
	assert.True(mockSystem1.UpdateSystemCalled, "Expected UpdateSystem to be called on mock system 1")
	assert.Equal("testData", mockSystem1.UpdateSystemData, "Expected correct data passed to UpdateSystem of mock system 1")
	assert.True(mockSystem2.UpdateSystemCalled, "Expected UpdateSystem to be called on mock system 2")
	assert.Equal("testData", mockSystem2.UpdateSystemData, "Expected correct data passed to UpdateSystem of mock system 2")
}

// TestUpdateSystemsForEntity - Tests if UpdateSystemsForEntity calls UpdateEntity on systems that require the entity's components.
func TestUpdateSystemsForEntity(t *testing.T) {
	assert := assert.New(t)
	sm := &SystemManager{}
	mockSystem1 := &MockSystem{RequiredComponents: []ComponentType{testComponentType}}
	mockSystem2 := &MockSystem{RequiredComponents: []ComponentType{"UnknownComponent"}}
	entity := &Entity{Blueprint: "testEntity"}
	testComponent := TestComponent{}
	entity.AddComponent(testComponent)

	sm.AddSystem(mockSystem1)
	sm.AddSystem(mockSystem2)

	err := sm.UpdateSystemsForEntity("testData", entity)

	assert.Nil(err, "Expected no error from UpdateSystemsForEntity")
	assert.True(mockSystem1.UpdateEntityCalled, "Expected UpdateEntity to be called on mock system 1")
	assert.Equal("testData", mockSystem1.UpdateEntityData, "Expected correct data passed to UpdateEntity of mock system 1")
	assert.Equal(entity, mockSystem1.UpdateEntityEntity, "Expected correct entity passed to UpdateEntity of mock system 1")
	assert.False(mockSystem2.UpdateEntityCalled, "Expected UpdateEntity not to be called on mock system 2")
}
