// Package path provides an allocation-efficient A* pathfinding implementation.
// Unlike go-astar, this implementation pre-allocates data structures and reuses
// them between calls to minimize garbage collection pressure.
package path

// Pather is an interface for nodes that can be pathfound.
// It uses integer IDs for efficient map lookups.
type Pather interface {
	// PathID returns a unique integer ID for this node.
	// For grid-based worlds, this is typically: z*width*height + y*width + x
	PathID() int

	// PathNeighborsAppend appends neighbors to the provided slice and returns it.
	// This avoids allocating a new slice on each call.
	PathNeighborsAppend(neighbors []Pather) []Pather

	// PathNeighborCost returns the movement cost to a neighbor node.
	PathNeighborCost(to Pather) float64

	// PathEstimatedCost returns the heuristic estimate to the goal.
	PathEstimatedCost(to Pather) float64
}

// node wraps a Pather with A* data
type node struct {
	pather Pather
	cost   float64
	rank   float64
	parent int // index into nodes slice, -1 if no parent
	open   bool
	closed bool
	index  int // index in priority queue
}

// AStar is a reusable A* pathfinder that minimizes allocations.
type AStar struct {
	// nodes maps PathID -> node index
	nodeIndex map[int]int
	nodes     []node

	// Priority queue (min-heap by rank)
	openSet []int // indices into nodes

	// Neighbor buffer
	neighbors []Pather

	// Result buffer
	result []Pather
}

// NewAStar creates a new AStar pathfinder with pre-allocated capacity.
// estimatedNodes is the expected maximum number of nodes to explore.
func NewAStar(estimatedNodes int) *AStar {
	if estimatedNodes < 64 {
		estimatedNodes = 64
	}
	return &AStar{
		nodeIndex: make(map[int]int, estimatedNodes),
		nodes:     make([]node, 0, estimatedNodes),
		openSet:   make([]int, 0, estimatedNodes/4),
		neighbors: make([]Pather, 0, 8),
		result:    make([]Pather, 0, 64),
	}
}

// reset clears the pathfinder for a new search without deallocating
func (a *AStar) reset() {
	// Clear the map by deleting all keys (faster than creating new map)
	for k := range a.nodeIndex {
		delete(a.nodeIndex, k)
	}
	a.nodes = a.nodes[:0]
	a.openSet = a.openSet[:0]
	a.neighbors = a.neighbors[:0]
	a.result = a.result[:0]
}

// getOrCreateNode gets or creates a node for the given pather
func (a *AStar) getOrCreateNode(p Pather) int {
	id := p.PathID()
	if idx, ok := a.nodeIndex[id]; ok {
		return idx
	}

	idx := len(a.nodes)
	a.nodes = append(a.nodes, node{
		pather: p,
		cost:   1e18, // infinity
		rank:   1e18,
		parent: -1,
		open:   false,
		closed: false,
		index:  -1,
	})
	a.nodeIndex[id] = idx
	return idx
}

// Heap operations for openSet
func (a *AStar) heapPush(nodeIdx int) {
	n := &a.nodes[nodeIdx]
	n.index = len(a.openSet)
	a.openSet = append(a.openSet, nodeIdx)
	a.heapUp(n.index)
}

func (a *AStar) heapPop() int {
	if len(a.openSet) == 0 {
		return -1
	}
	top := a.openSet[0]
	last := len(a.openSet) - 1
	a.heapSwap(0, last)
	a.openSet = a.openSet[:last]
	if len(a.openSet) > 0 {
		a.heapDown(0)
	}
	a.nodes[top].index = -1
	return top
}

func (a *AStar) heapRemove(heapIdx int) {
	last := len(a.openSet) - 1
	if heapIdx != last {
		a.heapSwap(heapIdx, last)
		a.openSet = a.openSet[:last]
		if heapIdx < len(a.openSet) {
			a.heapDown(heapIdx)
			a.heapUp(heapIdx)
		}
	} else {
		a.openSet = a.openSet[:last]
	}
}

func (a *AStar) heapUp(idx int) {
	for idx > 0 {
		parent := (idx - 1) / 2
		if a.nodes[a.openSet[parent]].rank <= a.nodes[a.openSet[idx]].rank {
			break
		}
		a.heapSwap(parent, idx)
		idx = parent
	}
}

func (a *AStar) heapDown(idx int) {
	n := len(a.openSet)
	for {
		smallest := idx
		left := 2*idx + 1
		right := 2*idx + 2

		if left < n && a.nodes[a.openSet[left]].rank < a.nodes[a.openSet[smallest]].rank {
			smallest = left
		}
		if right < n && a.nodes[a.openSet[right]].rank < a.nodes[a.openSet[smallest]].rank {
			smallest = right
		}
		if smallest == idx {
			break
		}
		a.heapSwap(idx, smallest)
		idx = smallest
	}
}

func (a *AStar) heapSwap(i, j int) {
	a.openSet[i], a.openSet[j] = a.openSet[j], a.openSet[i]
	a.nodes[a.openSet[i]].index = i
	a.nodes[a.openSet[j]].index = j
}

// Path finds a path from 'from' to 'to'.
// Returns the path (from start to goal), total cost, and whether a path was found.
// The returned path slice is reused between calls - copy it if you need to keep it.
func (a *AStar) Path(from, to Pather) (path []Pather, distance float64, found bool) {
	a.reset()

	toID := to.PathID()

	// Initialize start node
	startIdx := a.getOrCreateNode(from)
	a.nodes[startIdx].cost = 0
	a.nodes[startIdx].rank = from.PathEstimatedCost(to)
	a.nodes[startIdx].open = true
	a.heapPush(startIdx)

	for len(a.openSet) > 0 {
		currentIdx := a.heapPop()
		current := &a.nodes[currentIdx]
		current.open = false
		current.closed = true

		// Check if we reached the goal
		if current.pather.PathID() == toID {
			// Reconstruct path
			a.result = a.result[:0]
			idx := currentIdx
			for idx != -1 {
				a.result = append(a.result, a.nodes[idx].pather)
				idx = a.nodes[idx].parent
			}
			// Reverse to get start->goal order
			for i, j := 0, len(a.result)-1; i < j; i, j = i+1, j-1 {
				a.result[i], a.result[j] = a.result[j], a.result[i]
			}
			return a.result, current.cost, true
		}

		// Explore neighbors
		a.neighbors = current.pather.PathNeighborsAppend(a.neighbors[:0])
		for _, neighbor := range a.neighbors {
			cost := current.cost + current.pather.PathNeighborCost(neighbor)
			neighborIdx := a.getOrCreateNode(neighbor)
			neighborNode := &a.nodes[neighborIdx]

			if cost < neighborNode.cost {
				if neighborNode.open {
					a.heapRemove(neighborNode.index)
					neighborNode.open = false
				}
				neighborNode.closed = false
			}

			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = currentIdx
				neighborNode.open = true
				a.heapPush(neighborIdx)
			}
		}
	}

	return nil, 0, false
}

// Default is a shared AStar instance for convenience.
// For concurrent use, create separate AStar instances.
var Default = NewAStar(1024)

// Path is a convenience function using the default AStar instance.
// Not safe for concurrent use - use NewAStar() for concurrent pathfinding.
func Path(from, to Pather) (path []Pather, distance float64, found bool) {
	return Default.Path(from, to)
}
