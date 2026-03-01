package main

import (
	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/transport"
)

// BallCodec implements transport.SnapshotCodec for the bouncing balls example.
// It serializes Position and Color, omitting server-only Velocity.
type BallCodec struct{}

var _ transport.SnapshotCodec = (*BallCodec)(nil)

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

func (c *BallCodec) Decode(snap *transport.Snapshot, world any) {
	w := world.(*World)
	for _, es := range snap.Entities {
		e := w.FindOrCreateEntity(es.ID)
		if pos, ok := es.Components[TypePosition]; ok {
			e.AddComponent(pos.(PositionComponent))
		}
		if col, ok := es.Components[TypeColor]; ok {
			e.AddComponent(col.(ColorComponent))
		}
	}
}
