package gateway

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go-agar/internal/game"
	"strconv"
	"strings"
	"sync"
)

type Session struct {
	name      string
	conn      *websocket.Conn
	player    *game.Player
	battle    *game.Battle
	broadcast chan *Chat
	locker    *sync.Mutex
}

func NewSession(name string, conn *websocket.Conn) *Session {
	broadcast := make(chan *Chat, 10)
	return &Session{
		name:      name,
		conn:      conn,
		broadcast: broadcast,
		locker:    &sync.Mutex{},
	}
}

func (s *Session) close() {
	if e := s.conn.Close(); e == nil {
		if s.battle != nil && s.player != nil {
			s.battle.RemovePlayer(s.player)
			select {
			case s.broadcast <- NewSystemChat("player [ " + s.player.Name + " ] exit"):
			default:
			}
		}
	}
}

func (s *Session) ping() {
	s.send(ActionPing, "")
}

func (s *Session) notify(msg string) {
	s.send(ActionChat, msg)
}

func (s *Session) join(b *game.Battle) bool {
	if s.player == nil {
		s.player = game.NewPlayer(s.name)
	}
	if !b.AddPlayer(s.player) {
		return false
	}
	s.battle = b
	s.send(ActionGameSetup, fmt.Sprintf("%.0f|%.0f|%.0f|%.0f|%s", game.Config.GameWidth, game.Config.GameHeight, game.Config.ScreenWidth, game.Config.ScreenHeight, game.Config.VirusColor))
	select {
	case s.broadcast <- NewSystemChat("player [ " + s.player.Name + " ] join"):
	default:
	}
	return true
}

func (s *Session) pushPlayerStatus() {
	if s.player != nil {
		s.send(ActionPlayerStatus, s.fmtPlayerStatus())
		if s.player.IsDied() {
			s.close()
		}
	}
}

func (s *Session) pushLeaderBoard() {
	s.send(ActionLeaderBoard, s.fmtLeaderBoard())
}

func (s *Session) say(msg string) {
	if s.battle != nil && s.player != nil {
		select {
		case s.broadcast <- NewPlayerChat(s.player.Name + " : " + msg):
		default:
		}
	}
}

func (s *Session) move(c *websocket.Conn, payload string) {
	if p := s.player; p != nil {
		split := strings.Split(payload, ",")
		x, e := strconv.Atoi(split[0])
		if e != nil {
			return
		}
		y, e := strconv.Atoi(split[1])
		if e != nil {
			return
		}
		p.MoveTo(float64(x), float64(y))
	}
}

func (s *Session) fire() {
	if s.player != nil {
		s.player.FireFood()
	}
}

func (s *Session) split() {
	if s.player != nil {
		s.player.SplitAll()
	}
}

func (s *Session) send(msgType string, data string) {
	s.locker.Lock()
	defer s.locker.Unlock()
	retry := 0
	for retry < 3 {
		e := s.conn.WriteMessage(websocket.TextMessage, []byte(msgType+"|"+data))
		if e == nil {
			return
		}
		retry++
	}
	s.close()
}

func (s *Session) fmtLeaderBoard() string {
	b := s.battle
	if b == nil {
		return ""
	}
	if len(b.LeaderBoard) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, name := range b.LeaderBoard {
		sb.WriteString(name)
		sb.WriteString(",")
	}
	res := sb.String()
	return res[:len(res)-1]
}

func (s *Session) fmtPlayerStatus() string {
	p := s.player
	if p == nil {
		return ""
	}
	var result strings.Builder
	fmt.Fprintf(&result, "%s,%.2f,%.2f,%.2f", p.Name, p.X, p.Y, p.MassTotal)
	fmt.Fprintf(&result, "|%d", len(p.VisibleCells))
	for _, v := range p.VisibleCells {
		fmt.Fprintf(&result, "|%s,%s,%s,%.2f,%.2f,%.2f", v.Name, v.Color, v.TextColor, v.X, v.Y, v.Radius)
	}
	fmt.Fprintf(&result, "|%d", len(p.VisibleFoods))
	for _, v := range p.VisibleFoods {
		fmt.Fprintf(&result, "|%.2f,%.2f,%.2f,%s", v.X, v.Y, v.Radius, v.Color)
	}
	fmt.Fprintf(&result, "|%d", len(p.VisibleMassFoods))
	for _, v := range p.VisibleMassFoods {
		fmt.Fprintf(&result, "|%.2f,%.2f,%.2f,%s", v.X, v.Y, v.Radius, v.Color)
	}
	fmt.Fprintf(&result, "|%d", len(p.VisibleViruses))
	for _, v := range p.VisibleViruses {
		fmt.Fprintf(&result, "|%.2f,%.2f,%.2f", v.X, v.Y, v.Radius)
	}
	return result.String()
}
