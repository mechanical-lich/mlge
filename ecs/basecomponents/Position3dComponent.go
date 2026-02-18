package basecomponents

import "github.com/mechanical-lich/mlge/ecs"

const Position3d ecs.ComponentType = "Position3dComponent"

// Position3dComponent .
type Position3dComponent struct {
	X, Y, Z float64
}

func (pc Position3dComponent) GetType() ecs.ComponentType {
	return "Position3dComponent"
}

func (pc Position3dComponent) GetX() float64 {
	return pc.X
}
func (pc Position3dComponent) GetY() float64 {
	return pc.Y
}

func (pc Position3dComponent) GetZ() float64 {
	return pc.Z
}

func (pc *Position3dComponent) SetPosition(x float64, y float64, z float64) {
	pc.X = x
	pc.Y = y
	pc.Z = z
}
