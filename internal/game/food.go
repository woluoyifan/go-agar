package game

import (
	"go-agar/internal/util"
)

type Food struct {
	X      float64
	Y      float64
	Radius float64
	mass   float64
	Color  string
}

func NewFood() *Food {
	mass := Config.FoodMass
	radius := util.MassToRadius(Config.FoodMass)
	x, y := util.RandomPosition(radius, Config.GameWidth, Config.GameHeight)
	return &Food{
		X:      x,
		Y:      y,
		Radius: radius,
		mass:   mass,
		Color:  util.Color(),
	}
}

func (f *Food) CircleStatus() (x, y, radius float64) {
	return f.X, f.Y, f.Radius
}

func (f *Food) IsColliding(other CollidingCircle) bool {
	x, y, radius := other.CircleStatus()
	return util.IsCycleColliding(f.X, f.Y, f.Radius, x, y, radius)
}