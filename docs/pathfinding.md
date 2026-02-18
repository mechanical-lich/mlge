---
layout: default
title: Pathfinding
nav_order: 7
---

# A* Pathfinding

`github.com/mechanical-lich/mlge/path`

Allocation-efficient A* pathfinding with pre-allocated internal data structures and reusable pathfinder instances.

## The Pather Interface

Any node in your graph must implement the `Pather` interface:

```go
type Pather interface {
    PathID() int
    PathNeighborsAppend(neighbors []Pather) []Pather
    PathNeighborCost(to Pather) float64
    PathEstimatedCost(to Pather) float64
}
```

| Method | Description |
|--------|-------------|
| `PathID` | Unique identifier for this node (used internally for lookups) |
| `PathNeighborsAppend` | Appends reachable neighbors to the provided slice and returns it |
| `PathNeighborCost` | Actual cost to move from this node to a neighbor |
| `PathEstimatedCost` | Heuristic estimate of cost from this node to the target (e.g., Manhattan distance) |

## Quick Start

Use the convenience function with the default pathfinder:

```go
import "github.com/mechanical-lich/mlge/path"

result, distance, found := path.Path(startNode, endNode)
if found {
    for _, node := range result {
        // Walk the path
    }
}
```

## Custom Pathfinder

For better performance in hot paths, create a dedicated `AStar` instance with a pre-allocated node pool:

```go
pathfinder := path.NewAStar(2048) // Pre-allocate for ~2048 nodes

result, distance, found := pathfinder.Path(from, to)
```

The `AStar` struct reuses memory across calls, reducing GC pressure in per-frame pathfinding.

## Example: Grid-Based Pathfinding

```go
type Tile struct {
    X, Y     int
    Passable bool
    Grid     *Grid
}

func (t *Tile) PathID() int {
    return t.Y*t.Grid.Width + t.X
}

func (t *Tile) PathNeighborsAppend(neighbors []path.Pather) []path.Pather {
    for _, dir := range [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
        nx, ny := t.X+dir[0], t.Y+dir[1]
        if neighbor := t.Grid.Get(nx, ny); neighbor != nil && neighbor.Passable {
            neighbors = append(neighbors, neighbor)
        }
    }
    return neighbors
}

func (t *Tile) PathNeighborCost(to path.Pather) float64 {
    return 1.0
}

func (t *Tile) PathEstimatedCost(to path.Pather) float64 {
    other := to.(*Tile)
    dx := math.Abs(float64(t.X - other.X))
    dy := math.Abs(float64(t.Y - other.Y))
    return dx + dy // Manhattan distance
}
```

## Default Instance

A shared `AStar` instance is available as `path.Default`, pre-allocated with 1024 nodes. The `path.Path()` convenience function uses this instance.

```go
var Default *AStar = NewAStar(1024)
```
