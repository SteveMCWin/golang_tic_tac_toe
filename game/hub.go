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
    Games []*Game
    // PlayerQueue chan *Player
    PlayerQueues map[GameMode]chan *Player
    // PlayerExit chan *Player
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
    h := Hub{Games: make([]*Game, 0), PlayerQueues: make(map[GameMode](chan *Player))}
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
    // make a list of all the players in queue
    // get 2 players from the queue and pass them to the NewGame
    // run the game i guess

    for i := GameMode(0); i < game_mode_size; i++ {
        go func () {
            for {
                g := NewGame(<-h.PlayerQueues[i], <-h.PlayerQueues[i], i)
                go g.Run()
                h.Games = append(h.Games, g)
            }
        }()
    }
    
    // for {
    //     // log.Println("MADE DA GAMEEEE")
    //     g := NewGame(<-h.PlayerQueue, <-h.PlayerQueue)
    //     go g.Run()
    //     h.Games = append(h.Games, g)
    // }

}

