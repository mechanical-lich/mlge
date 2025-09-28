package basecomponents

import "github.com/mechanical-lich/mlge/ecs/v2"

const Position2d ecs.ComponentType = "Position2dComponent"

// Position2dComponent .
type Position2dComponent struct {
	X, Y float64
}

func (pc Position2dComponent) GetType() ecs.ComponentType {
	return "Position2dComponent"
}

func (pc Position2dComponent) GetX() float64 {
	return pc.X
}
func (pc Position2dComponent) GetY() float64 {
	return pc.Y
}

func (pc *Position2dComponent) SetPosition(x float64, y float64) {
	pc.X = x
	pc.Y = y
}
