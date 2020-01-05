package game

type CollidingCircle interface {
	CircleStatus() (x,y,radius float64)
	IsColliding(other CollidingCircle) bool
}