package game

import (
    "log"
    "net/http"
	"github.com/gin-gonic/gin"
)

type GameDuration int

type Hub struct {
    PlayerQueues map[GameMode]chan *Player  // one queue for each game mode so you can pair up users accordingly
}

func ServeHub(c *gin.Context) {
    c.HTML(http.StatusOK, "hub.html", gin.H{})
}

func MakeHub() *Hub {
    h := Hub{ PlayerQueues: make(map[GameMode](chan *Player)) }
    for i := GameMode(0); i < game_mode_size; i++ {
        h.PlayerQueues[i] = make(chan *Player, 20)
    }
    return &h
}

func (h *Hub) AddPlayer(p *Player, game_mode int) {     // puts player in one of the queues based on the game_mode
    if game_mode < 0 || game_mode > int(game_mode_size) {   // this should never occur but in case it does, just set it to 0
        log.Println("Invalid game mode passed, setting it to 0")
        game_mode = 0
    }
    h.PlayerQueues[GameMode(game_mode)] <- p    // notifies the channel the player is waiting for a game
}

func (h *Hub) HandleGames() {   // starts a goroutine for every queue
    for i := GameMode(0); i < game_mode_size; i++ {
        go func () {
            for {
                g := NewGame(<-h.PlayerQueues[i], <-h.PlayerQueues[i], i)   // makes the game once there are 2 players in the same queue
                go g.Run()  // starts the game
            }
        }()
    }
}

