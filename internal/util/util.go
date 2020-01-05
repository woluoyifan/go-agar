package util

import (
	"github.com/satori/go.uuid"
	"math"
	"math/rand"
	"strings"
	"time"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

//质量半径转换
func MassToRadius(mass float64) float64 {
	return 4 + math.Sqrt(mass)*6
}

func Color() string {
	s:="0123456789abcdef"
	split :=strings.Split(s,"")
	var b strings.Builder
	b.WriteString("#")
	for i:=0;i<6;i++{
		b.WriteString(split[rand.Intn(16)])
	}
	return b.String()
}

//对角线线距减双方半径
func GetDistance(x1, y1, radius1, x2, y2, radius2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2)+math.Pow(y2-y1, 2)) - radius1 - radius2
}

func IsCycleColliding(x1,y1,radius1,x2,y2,radius2 float64) bool{
	radiusTotal := radius1 + radius2
	//fast fail
	if radiusTotal > math.Abs(x1-x2) && radiusTotal > math.Abs(y1-y2) {
		powDistance := math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2)
		if math.Pow(radiusTotal, 2) > powDistance {
			return true
		}
	}
	return false
}

func RandomInRange(from, to float64) float64 {
	return math.Floor(rand.Float64()*(to-from)) + from
}

func RandomPosition(radius, gameWidth, gameHeight float64) (float64, float64) {
	x := RandomInRange(radius, gameWidth-radius)
	y := RandomInRange(radius, gameHeight-radius)
	return x, y
}

func Log(n ,base float64) float64{
	logN := math.Log(n)
	if base == 0{
		return logN
	}
	return logN / math.Log(base)
}

func GenId() string {
	id := uuid.NewV4()
	return strings.ReplaceAll(id.String(), "-", "")
}
