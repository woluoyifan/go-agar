package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-agar/internal/asset"
	"go-agar/internal/game"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	ActionError        = "00"
	ActionPing         = "01"
	ActionGameSetup    = "02"
	ActionChat         = "03"
	ActionPlayerStatus = "04"
	ActionMove         = "05"
	ActionFire         = "06"
	ActionSplit        = "07"
	ActionLeaderBoard  = "08"
)

type Gateway struct {
	wsCreator      *websocket.Upgrader
	sessionBattles map[*Session]*game.Battle
	battles        []*game.Battle
	battleLocker   *sync.Mutex
}

func NewGateway() (*Gateway, error) {
	return &Gateway{
		wsCreator: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		sessionBattles: make(map[*Session]*game.Battle),
		battleLocker:   &sync.Mutex{},
	}, nil
}

func (g *Gateway) Run() {
	defer g.Stop()
	if !game.Config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	engine.GET("/game", func(c *gin.Context) {
		go g.openSession(c)
		c.Done()
	})
	engine.GET("/", func(c *gin.Context) {
		bytes := asset.MustAsset("web/index.html")
		c.Data(http.StatusOK, "text/html", bytes)
	})
	engine.GET("/game.js", func(c *gin.Context) {
		bytes := asset.MustAsset("web/game.js")
		c.Data(http.StatusOK, "application/x-javascript", bytes)
	})
	engine.GET("/game.css", func(c *gin.Context) {
		bytes := asset.MustAsset("web/game.css")
		c.Data(http.StatusOK, "text/css", bytes)
	})
	println("server start on port [", game.Config.Port, "]")
	engine.Run(":" + strconv.Itoa(game.Config.Port))
}

func (g *Gateway) Stop() {
	for _, b := range g.battles {
		b.Stop()
	}
}

func (g *Gateway) openSession(context *gin.Context) {
	conn, err := g.wsCreator.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	session := NewSession(context.Query("name"), conn)
	defer g.closeSession(session)
	g.allocationBattle(session)
	for {
		if e := conn.SetReadDeadline(time.Now().Add(game.Config.MaxHeartbeatInterval)); e != nil {
			return
		}
		_, bytes, e := conn.ReadMessage()
		if e != nil {
			return
		}
		msg := string(bytes)
		msgType := msg[:2]
		switch msgType {
		case ActionError:
		case ActionPing:
			session.ping()
		case ActionChat:
			session.say(msg[3:])
		case ActionMove:
			session.move(conn, msg[3:])
		case ActionFire:
			session.fire()
		case ActionSplit:
			session.split()
		}
	}
}

func (g *Gateway) closeSession(s *Session) {
	s.close()
	g.battleLocker.Lock()
	b := g.sessionBattles[s]
	delete(g.sessionBattles, s)
	g.battleLocker.Unlock()
	for {
		select {
		case msg := <-s.broadcast:
			g.broadcast(b, msg)
		default:
			return
		}
	}
}

func (g *Gateway) allocationBattle(s *Session) {
	g.battleLocker.Lock()
	defer g.battleLocker.Unlock()
	if _, ok := g.sessionBattles[s]; ok {
		return
	}
	for _, b := range g.battles {
		if s.join(b) {
			g.sessionBattles[s] = b
			return
		}
	}
	b := game.NewBattle()
	g.battles = append(g.battles, b)
	go g.mountBattle(b)
	if s.join(b) {
		g.sessionBattles[s] = b
	} else {
		g.closeSession(s)
	}
}

func (g *Gateway) mountBattle(b *game.Battle) {
	leaderBoardTicker := time.NewTicker(time.Second)
	defer leaderBoardTicker.Stop()
	for {
		select {
		case _, ok := <-b.Tick:
			if !ok {
				for s, b2 := range g.sessionBattles {
					if b2 == b {
						g.closeSession(s)
					}
				}
				return
			}
			for s, b2 := range g.sessionBattles {
				if b2 == b {
					s.pushPlayerStatus()
					select {
					case c := <-s.broadcast:
						g.broadcast(b, c)
					default:
					}
				}
			}
		case <-leaderBoardTicker.C:
			for s, b2 := range g.sessionBattles {
				if b2 == b {
					s.pushLeaderBoard()
				}
			}
		}
	}
}

func (g *Gateway) broadcast(b *game.Battle, c *Chat) {
	for s, b2 := range g.sessionBattles {
		if b2 == b {
			s.send(ActionChat, c.Type+c.Data)
		}
	}
}
