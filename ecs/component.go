package ecs

type ComponentType string

// Component base component interface
type Component interface {
	GetType() ComponentType
}
