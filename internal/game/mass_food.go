package game

import (
	"go-agar/internal/util"
	"math"
)

type MassFood struct {
	Id      string
	X       float64
	Y       float64
	Radius  float64
	mass    float64
	speed   float64
	targetX float64
	targetY float64
	player  *Player
	Color   string
}

func NewMassFood(player *Player) *MassFood {
	return &MassFood{
		Id:     util.GenId(),
		player: player,
		Color:  player.Color,
	}
}

func (mf *MassFood) Move() {
	if mf.speed == 0 {
		return
	}
	deg := math.Atan2(mf.targetY, mf.targetX)
	deltaX := mf.speed * math.Cos(deg)
	deltaY := mf.speed * math.Sin(deg)
	mf.speed -= Config.FireFoodSpeedSlow
	if mf.speed < 0 {
		mf.speed = 0
	}
	mf.X += deltaX
	mf.Y += deltaY

	var borderCalc = mf.Radius + 5
	if mf.X > Config.GameWidth-borderCalc {
		mf.X = Config.GameWidth - borderCalc
	}
	if mf.Y > Config.GameHeight-borderCalc {
		mf.Y = Config.GameHeight - borderCalc
	}
	if mf.X < borderCalc {
		mf.X = borderCalc
	}
	if mf.Y < borderCalc {
		mf.Y = borderCalc
	}
}

func (mf *MassFood) CircleStatus() (x, y, radius float64) {
	return mf.X, mf.Y, mf.Radius
}

func (mf *MassFood) IsColliding(other CollidingCircle) bool {
	x, y, radius := other.CircleStatus()
	return util.IsCycleColliding(mf.X, mf.Y, mf.Radius, x, y, radius)
}