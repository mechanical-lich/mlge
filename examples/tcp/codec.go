package main

import (
	"encoding/json"
	"image/color"
	"log"

	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/transport"
)

// BallCodec implements transport.SnapshotCodec for the TCP bouncing-balls
// example.
//
// Encode is identical to the local example: build a Snapshot from live server
// entities, omitting VelocityComponent (server-private).
//
// Decode must handle the TCP case: component values arrive over JSON as
// map[string]interface{} rather than typed structs (because transport.ComponentData
// is any, and encoding/json decodes unknown objects to maps). The remarshal
// helper re-encodes such a value to JSON bytes and then unmarshals into the
// correct concrete type, recovering full type information.
type BallCodec struct{}

var _ transport.SnapshotCodec = (*BallCodec)(nil)

// Encode builds a Snapshot from the server's entity list.
func (c *BallCodec) Encode(tick uint64, entities []*ecs.Entity) *transport.Snapshot {
	snaps := make([]*transport.EntitySnapshot, 0, len(entities))
	for _, e := range entities {
		idC, ok := e.Components[TypeID]
		if !ok {
			continue
		}
		posC, ok := e.Components[TypePosition]
		if !ok {
			continue
		}
		colC, ok := e.Components[TypeColor]
		if !ok {
			continue
		}
		snaps = append(snaps, &transport.EntitySnapshot{
			ID:        idC.(IDComponent).ID,
			Blueprint: e.Blueprint,
			Components: map[ecs.ComponentType]transport.ComponentData{
				TypePosition: posC.(PositionComponent),
				TypeColor:    colC.(ColorComponent),
			},
		})
	}
	return transport.NewSnapshot(tick, snaps)
}

// Decode applies a received Snapshot to the client's World.
//
// It reconciles the full entity set: entities absent from the snapshot are
// removed from the world so the client does not accumulate stale entities
// when the server destroys them.
//
// When called with a snapshot that came over TCP, component values are
// map[string]interface{} (JSON default). The remarshal helper converts them
// back to typed structs.
func (c *BallCodec) Decode(snap *transport.Snapshot, world any) {
	w := world.(*World)

	// Build a set of IDs present in this snapshot.
	seen := make(map[string]bool, len(snap.Entities))

	for _, es := range snap.Entities {
		seen[es.ID] = true
		e := w.FindOrCreateEntity(es.ID)

		if raw, ok := es.Components[TypePosition]; ok {
			var pos PositionComponent
			if typed, ok := raw.(PositionComponent); ok {
				// Local transport: value is already the correct type.
				pos = typed
			} else if m, ok := raw.(map[string]any); ok {
				// TCP transport: JSON decoded to map - pull fields directly,
				// no re-marshal needed for simple numeric types.
				pos.X, _ = m["X"].(float64)
				pos.Y, _ = m["Y"].(float64)
			} else {
				log.Printf("codec: unexpected position type %T", raw)
				continue
			}
			e.AddComponent(pos)
		}

		if raw, ok := es.Components[TypeColor]; ok {
			var col ColorComponent
			if typed, ok := raw.(ColorComponent); ok {
				col = typed
			} else {
				// TCP path: JSON encodes color.RGBA as {"R":n,"G":n,"B":n,"A":n}.
				// Unmarshal into the intermediate struct first.
				var wire struct {
					RGBA struct {
						R, G, B, A uint8
					}
				}
				if err := remarshal(raw, &wire); err != nil {
					log.Printf("codec: decode color: %v", err)
					continue
				}
				col = ColorComponent{RGBA: color.RGBA{
					R: wire.RGBA.R,
					G: wire.RGBA.G,
					B: wire.RGBA.B,
					A: wire.RGBA.A,
				}}
			}
			e.AddComponent(col)
		}
	}

	// Remove entities that are no longer present in the snapshot.
	// Without this, destroyed server entities accumulate in the client world.
	w.RemoveAbsent(seen)
}

// remarshal is a helper that round-trips src through JSON into dst.
// It is used to convert map[string]interface{} values (produced by
// encoding/json when decoding into an any field) back into typed structs.
func remarshal(src any, dst any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
