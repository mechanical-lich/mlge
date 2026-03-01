package transport

import (
	"time"

	"github.com/mechanical-lich/mlge/ecs"
)

// ComponentData holds the serialized state of a single component.
// In Phase 1 (local transport) this is a direct Go value - no encoding needed.
// In a future TCP transport, the SnapshotCodec is responsible for marshaling
// this to json.RawMessage or a proto bytes field before transmission.
type ComponentData = any

// EntitySnapshot captures the state of one entity at a given server tick.
type EntitySnapshot struct {
	// ID is a game-assigned identifier for the entity.
	// mlge entities have no built-in ID; games must assign one via a component
	// (e.g., an IDComponent) or another scheme. The SnapshotCodec is responsible
	// for populating this field.
	ID string

	// Blueprint is the entity's blueprint name, copied from ecs.Entity.Blueprint.
	Blueprint string

	// Components maps component type → component value.
	// The SnapshotCodec decides which components to include and how to encode them.
	Components map[ecs.ComponentType]ComponentData
}

// Snapshot is the authoritative world state produced by the server each tick.
// It is consumed by the client to update its local copy of the world for rendering.
type Snapshot struct {
	// Tick is the server tick number this snapshot was produced on.
	Tick uint64

	// Timestamp is when the snapshot was produced (UnixNano).
	Timestamp int64

	// Entities holds one snapshot per simulated entity.
	// Phase 1 always sends the full entity list (no delta compression).
	Entities []*EntitySnapshot
}

// NewSnapshot is a convenience constructor used by the Server.
func NewSnapshot(tick uint64, entities []*EntitySnapshot) *Snapshot {
	return &Snapshot{
		Tick:      tick,
		Timestamp: time.Now().UnixNano(),
		Entities:  entities,
	}
}

// SnapshotCodec is implemented by the game to control snapshot serialization.
//
// The server calls Encode every SnapshotEvery ticks to produce a Snapshot.
// The client calls Decode to apply the snapshot to its local world representation.
//
// This keeps mlge/transport free of any game-specific component knowledge.
type SnapshotCodec interface {
	// Encode serializes a slice of live entities into a Snapshot.
	// Only include components that the client needs to render the world.
	// Omit pure-simulation components (brain weights, AI state) to reduce payload.
	Encode(tick uint64, entities []*ecs.Entity) *Snapshot

	// Decode applies a received Snapshot to the client's local world.
	// Implementations typically find-or-create local entities by ID and overwrite
	// the relevant components with the snapshot values.
	Decode(snapshot *Snapshot, world any)
}
