package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestComponent struct{}

var testComponentType ComponentType = "TestComponent"

func (c TestComponent) GetType() ComponentType {
	return testComponentType
}

type TestComponent2 struct{}

var testComponent2Type ComponentType = "TestComponent2"

func (c TestComponent2) GetType() ComponentType {
	return testComponent2Type
}

func TestAddComponent(t *testing.T) {
	assert := assert.New(t)
	entity := &Entity{}
	entity.AddComponent(TestComponent{})
	assert.True(entity.HasComponent(testComponentType), "Expected entity to have the component %s", testComponentType)
}

func TestHasComponent(t *testing.T) {
	assert := assert.New(t)
	entity := &Entity{}
	entity.AddComponent(TestComponent{})

	assert.True(entity.HasComponent(testComponentType), "Expected entity to have the component %s", testComponentType)
}

func TestHasComponents(t *testing.T) {
	assert := assert.New(t)
	entity := &Entity{}
	testComponents := []Component{TestComponent{}, TestComponent2{}}
	componentTypes := []ComponentType{}
	for _, ct := range testComponents {
		entity.AddComponent(ct)
		componentTypes = append(componentTypes, ct.GetType())
	}

	assert.True(entity.HasComponents(componentTypes...), "Expected entity to have all components %v", componentTypes)
}

func TestGetComponent(t *testing.T) {
	assert := assert.New(t)
	entity := &Entity{}
	entity.AddComponent(TestComponent{})

	retrievedComponent := entity.GetComponent(testComponentType)

	assert.Equal(TestComponent{}, retrievedComponent, "Expected to retrieve the same component")
}

func TestRemoveComponent(t *testing.T) {
	assert := assert.New(t)
	entity := &Entity{}
	entity.AddComponent(TestComponent{})

	entity.RemoveComponent(testComponentType)

	assert.False(entity.HasComponent(testComponentType), "Expected entity to not have the component %s after removal", testComponentType)
}
