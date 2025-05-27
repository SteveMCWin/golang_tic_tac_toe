package game

import (
    "net/http"
    "html/template"
	"github.com/gin-gonic/gin"
)

type GameMode int

const (
    normal_5s GameMode = iota
    normal_10s
    normal_15s
    game_mode_size  // make sure to keep this the last on the list
)

type Hub struct {
    PlayerQueues map[GameMode]chan *Player
    // Games []*Game
}

func EnterHub(c *gin.Context) {
    tmpl, err := template.ParseFiles("templates/hub.html")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(c.Writer, gin.H{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func MakeHub() *Hub {
    h := Hub{/* Games: make([]*Game, 0), */ PlayerQueues: make(map[GameMode](chan *Player))}
    for i := GameMode(0); i < game_mode_size; i++ {
        h.PlayerQueues[i] = make(chan *Player, 20)
    }
    return &h
}

func (h *Hub) AddPlayer(p *Player, game_mode int) {
    if game_mode < 0 || game_mode > int(game_mode_size) {
        panic("Invalid game mode when adding player")
    }
    h.PlayerQueues[GameMode(game_mode)] <- p
}

func (h *Hub) HandleGames() {
    for i := GameMode(0); i < game_mode_size; i++ {
        go func () {
            for {
                g := NewGame(<-h.PlayerQueues[i], <-h.PlayerQueues[i], i)
                go g.Run()
                // h.Games = append(h.Games, g)
            }
        }()
    }
}

