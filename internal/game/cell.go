package game

import "go-agar/internal/util"

type Cell struct {
	Name      string
	X         float64
	Y         float64
	Radius    float64
	mass      float64
	speed     float64
	Color     string
	TextColor string
}

func NewCell(p *Player) *Cell {
	return &Cell{
		Name:      p.Name,
		X:         p.X,
		Y:         p.Y,
		Color:     p.Color,
		TextColor: p.TextColor,
	}
}

func (c *Cell) CircleStatus() (x, y, radius float64) {
	return c.X, c.Y, c.Radius
}

func (c *Cell) IsColliding(other CollidingCircle) bool {
	x, y, radius := other.CircleStatus()
	return util.IsCycleColliding(c.X, c.Y, c.Radius, x, y, radius)
}
