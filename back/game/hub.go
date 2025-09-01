package game

import (
    "log"

	"tic_tac_toe.fun/models"
)

// Hub represents multiple queues in which players wait for an opponent.
type Hub struct {
    PlayerQueues map[GameMode]chan *Player  // There is one queue for each game mode so you can pair up users accordingly
}

// MakeHub initializes the PlayerQueues map and the channels coresponding to each game mode.
func MakeHub() *Hub {
    h := Hub{ PlayerQueues: make(map[GameMode](chan *Player)) }
    for i := GameMode(0); i < game_mode_size; i++ {
        h.PlayerQueues[i] = make(chan *Player, 20)
    }
    return &h
}

// AddPlayer places a player in a queue corresponding to the game mode the player selected.
// Gets called after a websocket connection is established with the user.
// Note that if the game mode that is passed in is note valid, the game mode will just be set to normal_180s.
func (h *Hub) AddPlayer(p *Player, game_mode GameMode) {     // puts player in one of the queues based on the game_mode
    if game_mode < 0 || game_mode > game_mode_size {   // this should never occur but in case it does, just set it to 0
        log.Println("Invalid game mode passed, setting it to 0")
        game_mode = 0
    }
    h.PlayerQueues[GameMode(game_mode)] <- p    // notifies the channel the player is waiting for a game
}

// HandleGames starts one queue (actually a channel) of players for each game mode. When there are two players inside the same queue,
// a new game is created and started.
func (h *Hub) HandleGames(db *models.DataBase) {   // starts a goroutine for every queue
    for i := GameMode(0); i < game_mode_size; i++ {
        go func () {
            for {
                g := NewGame(<-h.PlayerQueues[i], <-h.PlayerQueues[i], i)   // makes the game once there are 2 players in the same queue
                go g.Run(db)  // starts the game
            }
        }()
    }
}

