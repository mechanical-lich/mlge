package main

import "github.com/mechanical-lich/mlge/ecs"

const ballRadius = 10.0

// World is the simple shared world type for this example.
// Both the server and the client hold a World; they are separate instances -
// the server's is authoritative, the client's is updated from snapshots.
type World struct {
	Width, Height float64
	Entities      []*ecs.Entity

	// entityByID caches client-side entities for O(1) snapshot decode lookups.
	entityByID map[string]*ecs.Entity
}

// FindOrCreateEntity returns the entity with the given ID, creating a new one
// if it doesn't yet exist. Used by the codec during snapshot decode.
func (w *World) FindOrCreateEntity(id string) *ecs.Entity {
	if w.entityByID == nil {
		w.entityByID = make(map[string]*ecs.Entity)
		for _, e := range w.Entities {
			if c, ok := e.Components[TypeID]; ok {
				w.entityByID[c.(IDComponent).ID] = e
			}
		}
	}
	if e, ok := w.entityByID[id]; ok {
		return e
	}
	e := &ecs.Entity{Blueprint: "ball"}
	e.AddComponent(IDComponent{ID: id})
	w.Entities = append(w.Entities, e)
	w.entityByID[id] = e
	return e
}
