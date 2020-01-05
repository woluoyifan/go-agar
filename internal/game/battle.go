package game

import (
	"go-agar/internal/util"
	"math"
	"sync"
	"time"
)

type Battle struct {
	Id             string
	players        []*Player
	viruses        []*Virus
	foods          []*Food
	massFoods      []*MassFood
	LeaderBoard    []string
	startTime      time.Time
	endTime        time.Time
	stop           chan byte
	tick           chan byte
	Tick           <-chan byte
	joinExitLocker *sync.Mutex
}

func NewBattle() *Battle {
	tick := make(chan byte, 1)
	b := &Battle{
		Id:             util.GenId(),
		stop:           make(chan byte),
		tick:           tick,
		Tick:           tick,
		joinExitLocker: &sync.Mutex{},
	}
	go b.run()
	return b
}

func (b *Battle) Stop() {
	b.stop <- 1
}

func (b *Battle) IsAccess() bool {
	return len(b.players) < Config.BattlePlayerLimit && b.endTime.IsZero()
}

func (b *Battle) AddPlayer(p *Player) bool {
	if !b.IsAccess() {
		return false
	}
	b.joinExitLocker.Lock()
	defer b.joinExitLocker.Unlock()
	b.players = append(b.players, p)
	return true
}

func (b *Battle) RemovePlayer(p *Player) {
	b.joinExitLocker.Lock()
	defer b.joinExitLocker.Unlock()
	for i, p2 := range b.players {
		if p2 == p {
			b.players = append(b.players[:i], b.players[i+1:]...)
			return
		}
	}
}

func (b *Battle) clear() {
	b.joinExitLocker.Lock()
	defer b.joinExitLocker.Unlock()
	close(b.stop)
	close(b.tick)
	b.players = nil
	b.viruses = nil
	b.foods = nil
	b.massFoods = nil
}

func (b *Battle) run() {
	defer b.clear()
	b.startTime = time.Now()
	shortTicker := time.NewTicker(time.Duration(1000/Config.TickRate) * time.Millisecond)
	defer shortTicker.Stop()
	longTicker := time.NewTicker(1 * time.Second)
	defer longTicker.Stop()
	for {
		select {
		case <-b.stop:
			b.endTime = time.Now()
			return
		case <-shortTicker.C:
			b.shortTickLoop()
		case <-longTicker.C:
			b.longTickLoop()
		}
	}
}

func (b *Battle) shortTickLoop() {
	b.updatePlayer()
	for _, mf := range b.massFoods {
		mf.Move()
	}
	select {
	case b.tick <- 1:
	default:
	}
}

func (b *Battle) longTickLoop() {
	b.balance()
	b.addVirus()
	b.updateLeaderBoard()
}

func (b *Battle) updateLeaderBoard() {
	b.joinExitLocker.Lock()
	defer b.joinExitLocker.Unlock()
	SortPlayers(b.players)
	max := len(b.players)
	if max > 10 {
		max = 10
	}
	leaderBoard := make([]string, max)
	for i := 0; i < max; i++ {
		leaderBoard[i] = b.players[i].Name
	}
	b.LeaderBoard = leaderBoard
}

func (b *Battle) updatePlayer() {
	for _, p := range b.players {
		p.Update()
		b.handlePlayerFireFood(p)
		b.handlePlayerCollision(p)
		p.UpdateVisibleFoods(b.foods)
		p.UpdateVisibleMassFoods(b.massFoods)
		p.UpdateVisibleViruses(b.viruses)
		p.UpdateVisibleCells(b.players)
	}
}

func (b *Battle) handlePlayerFireFood(p *Player) {
	for {
		select {
		case mf := <-p.FireFoods():
			b.massFoods = append(b.massFoods, mf)
		default:
			return
		}
	}
}

func (b *Battle) handlePlayerCollision(p *Player) {
	for _, c := range p.cells {
		//avoid memory copy
		i := 0
		for _, f := range b.foods {
			if util.IsCycleColliding(f.X, f.Y, f.Radius, c.X, c.Y, c.Radius) {
				c.mass += f.mass
				p.MassTotal += f.mass
				continue
			}
			b.foods[i] = f
			i++
		}
		b.foods = b.foods[:i]

		//avoid memory copy
		i = 0
		for _, mf := range b.massFoods {
			if mf.speed == 0 && mf.player != p && mf.IsColliding(c) {
				c.mass += mf.mass
				p.MassTotal += mf.mass
				continue
			}
			b.massFoods[i] = mf
			i++
		}
		b.massFoods = b.massFoods[:i]

		//avoid memory copy
		i = 0
		for _, v := range b.viruses {
			if v.mass < c.mass && v.IsColliding(c) {
				p.Split(c)
				continue
			}
			b.viruses[i] = v
			i++
		}
		b.viruses = b.viruses[:i]

		c.Radius = util.MassToRadius(c.mass)

		for _, p2 := range b.players {
			if p2 == p {
				continue
			}
			i := 0
			for _, c2 := range p2.cells {
				if c2.mass > Config.DefaultPlayerMass &&
					c2.mass*Config.MassWinRate < c.mass &&
					c2.IsColliding(c) &&
					c.Radius > util.GetDistance(c.X, c.Y, c.Radius, c2.X, c2.Y, c2.Radius)*Config.CellMergeDistanceRate {
					p2.MassTotal -= c2.mass
					p.MassTotal += c2.mass
					c.mass += c2.mass
					continue
				}
				p2.cells[i] = c2
				i++
			}
			p2.cells = p2.cells[:i]
		}

		massRetentionPerTick := c.mass * (1 - Config.MassLoseRate/1000/float64(Config.TickRate))
		if massRetentionPerTick > Config.DefaultPlayerMass && p.MassTotal > Config.MinMassLose {
			p.MassTotal -= c.mass - massRetentionPerTick
			c.mass = massRetentionPerTick
		}
	}
}

func (b *Battle) balance() {
	totalMass := float64(len(b.foods)) * Config.FoodMass
	for _, p := range b.players {
		totalMass += p.MassTotal
	}

	massDiff := Config.GameMaxMass - totalMass
	maxDiff := Config.FoodMaxNum - len(b.foods)
	diff := int(massDiff/Config.FoodMass) - maxDiff
	add := int(math.Min(float64(diff), float64(maxDiff)))
	remove := int(-math.Max(float64(diff), float64(maxDiff)))
	if add > 0 {
		foods := make([]*Food, add)
		for i := 0; i < add; i++ {
			foods[i] = NewFood()
		}
		b.foods = append(b.foods, foods...)
	} else if remove > 0 {
		b.foods = b.foods[:len(b.foods)-remove]
	}
}

func (b *Battle) addVirus() {
	add := Config.VirusMaxNum - len(b.viruses)
	viruses := make([]*Virus, add)
	for i := 0; i < add; i++ {
		viruses[i] = NewVirus()
	}
	b.viruses = append(b.viruses, viruses...)
}
