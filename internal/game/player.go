package game

import (
	"go-agar/internal/util"
	"math"
	"sort"
	"time"
)

type Player struct {
	Id               string
	Name             string
	cells            []*Cell
	X                float64
	Y                float64
	targetX          float64
	targetY          float64
	Color            string
	TextColor        string
	lastSplit        time.Time
	fireFood         chan *MassFood
	split            chan *Cell
	MassTotal        float64
	VisibleFoods     []*Food
	VisibleMassFoods []*MassFood
	VisibleCells     []*Cell
	VisibleViruses   []*Virus
}

type personSlice []*Player

func (s personSlice) Len() int           { return len(s) }
func (s personSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s personSlice) Less(i, j int) bool { return s[i].MassTotal > s[j].MassTotal }

func SortPlayers(players []*Player) {
	sort.Sort(personSlice(players))
}

func NewPlayer(name string) *Player {
	id := util.GenId()
	if name == "" {
		name = Config.AnonymousUserNamePrefix + id[:6]
	}
	mass := Config.DefaultPlayerMass
	radius := util.MassToRadius(mass)
	//x, y := util.RandomPosition(radius, Config.GameWidth, Config.GameHeight)
	x, y := Config.GameWidth/2, Config.GameHeight/2

	p := &Player{
		Id:            id,
		Name:          name,
		X:             x,
		Y:             y,
		Color:         util.Color(),
		TextColor:     "#000000",
		fireFood:      make(chan *MassFood, Config.CellMaxNum*10),
		split:         make(chan *Cell, Config.CellMaxNum),
		MassTotal:     mass,
	}
	c := p.addCell()
	c.mass = mass
	c.Radius = radius
	return p
}

func (p *Player) IsDied() bool {
	return len(p.cells) == 0 || int(p.MassTotal) == 0
}

func (p *Player) SplitAll() {
	for _, v := range p.cells {
		p.Split(v)
	}
}

func (p *Player) Split(c *Cell) {
	if len(p.cells) >= Config.CellMaxNum || c.mass < Config.DefaultPlayerMass*2 {
		return
	}
	select {
	case p.split <- c:
	default:
	}
}

func (p *Player) MoveTo(x, y float64) {
	p.targetX = x
	p.targetY = y
}

func (p *Player) FireFood() {
	for _, c := range p.cells {
		fireMass := Config.FireFoodRate * c.mass
		if Config.FireFoodMass > fireMass {
			fireMass = Config.FireFoodMass
		}
		if c.mass-fireMass <= Config.DefaultPlayerMass {
			return
		}
		c.mass -= fireMass
		p.MassTotal -= fireMass
		mf := NewMassFood(p)
		mf.mass = fireMass
		mf.Radius = util.MassToRadius(fireMass)
		mf.X = c.X
		mf.Y = c.Y
		mf.targetX = p.X - c.X + p.targetX
		mf.targetY = p.Y - c.Y + p.targetY
		mf.speed = Config.FireFoodSpeed
		p.fireFood <- mf
	}
}

func (p *Player) FireFoods() <-chan *MassFood {
	return p.fireFood
}

func (p *Player) UpdateVisibleFoods(foods [] *Food) {
	visibleFoods := make([]*Food, len(foods))
	i := 0
	for _, v := range foods {
		if v.X > p.X-Config.ScreenWidth/2-20 &&
			v.X < p.X+Config.ScreenWidth/2+20 &&
			v.Y > p.Y-Config.ScreenHeight/2-20 &&
			v.Y < p.Y+Config.ScreenHeight/2+20 {
			visibleFoods[i] = v
			i++
		}
	}
	p.VisibleFoods = visibleFoods[:i]
}

func (p *Player) UpdateVisibleMassFoods(massFoods [] *MassFood) {
	visibleMassFoods := make([]*MassFood, len(massFoods))
	i := 0
	for _, v := range massFoods {
		if v.X > p.X-Config.ScreenWidth/2-v.Radius &&
			v.X < p.X+Config.ScreenWidth/2+v.Radius &&
			v.Y > p.Y-Config.ScreenHeight/2-v.Radius &&
			v.Y < p.Y+Config.ScreenHeight/2+v.Radius {
			visibleMassFoods[i] = v
			i++
		}
	}
	p.VisibleMassFoods = visibleMassFoods[:i]
}

func (p *Player) UpdateVisibleViruses(viruses []*Virus) {
	visibleViruses := make([]*Virus, len(viruses))
	i := 0
	for _, v := range viruses {
		if v.X > p.X-Config.ScreenWidth/2-20-v.Radius &&
			v.X < p.X+Config.ScreenWidth/2+20+v.Radius &&
			v.Y > p.Y-Config.ScreenHeight/2-20-v.Radius &&
			v.Y < p.Y+Config.ScreenHeight/2+20+v.Radius {
			visibleViruses[i] = v
			i++
		}
	}
	p.VisibleViruses = visibleViruses[:i]
}

func (p *Player) UpdateVisibleCells(players []*Player) {
	// The maximum number of cells does not exceed the maximum number of divisions performed by all users
	visibleCells := make([]*Cell, len(players)*Config.CellMaxNum)
	i := 0
	for _, p2 := range players {
		for _, v := range p2.cells {
			if v.X > p.X-Config.ScreenWidth/2-20-v.Radius &&
				v.X < p.X+Config.ScreenWidth/2+20+v.Radius &&
				v.Y > p.Y-Config.ScreenHeight/2-20-v.Radius &&
				v.Y < p.Y+Config.ScreenHeight/2+20+v.Radius {
				visibleCells[i] = v
				i++
			}
		}
	}
	p.VisibleCells = visibleCells[:i]
}

func (p *Player) Update() {
	p.updateSplit()
	x, y := float64(0), float64(0)
	for i := 0; i < len(p.cells); i++ {
		c := p.cells[i]
		if c.speed == 0 {
			c.speed = Config.CellDefaultSpeed
		}
		p.moveCell(c)
		p.mergeCell(i)
		if len(p.cells) > i {
			p.borderReboundCell(c)
			x += c.X
			y += c.Y
		}
	}
	clen := float64(len(p.cells))
	p.X = x / clen
	p.Y = y / clen
}

func (p *Player) updateSplit() {
	for {
		select {
		case c := <-p.split:
			p.splitCell(c)
		default:
			return
		}
	}
}

func (p *Player) moveCell(c *Cell) {
	targetX := p.X - c.X + p.targetX
	targetY := p.Y - c.Y + p.targetY
	dist := math.Sqrt(math.Pow(targetY, 2) + math.Pow(targetX, 2))
	deg := math.Atan2(targetY, targetX)
	//slow easy...
	slowDown := float64(1)
	if c.speed <= Config.CellDefaultSpeed {
		slowDown = util.Log(c.mass, Config.SlowBase) - Config.InitMassLog + 1
	}
	deltaY := c.speed * math.Sin(deg) / slowDown
	deltaX := c.speed * math.Cos(deg) / slowDown
	if c.speed > Config.CellDefaultSpeed {
		c.speed -= 0.5
	}
	//why 50 ?
	if dist < (50 + c.Radius) {
		deltaY *= dist / (50 + c.Radius)
		deltaX *= dist / (50 + c.Radius)
	}
	c.Y += deltaY
	c.X += deltaX
}

func (p *Player) mergeCell(i int) {
	c := p.cells[i]
	//merge or separate
	mergePermit := p.lastSplit.Add(Config.MergeInterval).Before(time.Now())
	for j := i + 1; j < len(p.cells); j++ {
		c2 := p.cells[j]
		distance := util.GetDistance(c.X, c.Y, 0, c2.X, c2.Y, 0)
		radiusTotal := c.Radius + c2.Radius
		if distance >= radiusTotal {
			continue
		}
		if mergePermit && radiusTotal > distance * Config.CellMergeDistanceRate {
			c.mass += c2.mass
			c.Radius = util.MassToRadius(c.mass)
			p.cells = append(p.cells[:j], p.cells[j+1:]...)
			continue
		}
		if c.X < c2.X {
			c.X --
		} else if c.X > c2.X {
			c.X ++
		}
		if c.Y < c2.Y {
			c.Y --
		} else if c.Y > c2.Y {
			c.Y ++
		}
	}
}

func (p *Player) borderReboundCell(c *Cell) {
	//why 3 ? it seems to overlap the border
	borderCalc := c.Radius / 3
	//border rebound
	if c.X > Config.GameWidth-borderCalc {
		c.X = Config.GameWidth - borderCalc
	}
	if c.X < borderCalc {
		c.X = borderCalc
	}
	if c.Y > Config.GameHeight-borderCalc {
		c.Y = Config.GameHeight - borderCalc
	}
	if c.Y < borderCalc {
		c.Y = borderCalc
	}
}

func (p *Player) splitCell(c *Cell) {
	if len(p.cells) >= Config.CellMaxNum || c.mass < Config.DefaultPlayerMass*2 {
		return
	}
	for _, v := range p.cells {
		if v == c {
			c.mass = c.mass / 2
			c.Radius = util.MassToRadius(c.mass)

			nc := p.addCell()
			nc.X = c.X
			nc.Y = c.Y
			nc.mass = c.mass
			nc.Radius = c.Radius
			nc.speed = Config.SplitSpeed

			p.lastSplit = time.Now()
			return
		}
	}
}

func (p *Player) addCell() *Cell {
	c := NewCell(p)
	p.cells = append(p.cells, c)
	return c
}
