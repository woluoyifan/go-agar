package game

import "go-agar/internal/util"

type Virus struct {
	Id     string
	X      float64
	Y      float64
	Radius float64
	mass   float64
}

func NewVirus() *Virus {
	mass := Config.VirusMinMass + util.RandomInRange(Config.VirusMinMass, Config.VirusMaxMass)
	radius := util.MassToRadius(mass)
	x, y := util.RandomPosition(radius, Config.GameWidth, Config.GameHeight)
	return &Virus{
		Id:     util.GenId(),
		X:      x,
		Y:      y,
		Radius: radius,
		mass:   mass,
	}
}

func (v *Virus) CircleStatus() (x, y, radius float64) {
	return v.X, v.Y, v.Radius
}

func (v *Virus) IsColliding(other CollidingCircle) bool {
	x, y, radius := other.CircleStatus()
	return util.IsCycleColliding(v.X, v.Y, v.Radius, x, y, radius)
}