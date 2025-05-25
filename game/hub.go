package game

import (

    // "tic_tac_toe.fun/users"

	// "github.com/gin-gonic/gin"
)

type Hub struct {
    games []*Game
    player_queue chan *Player
}

func (h *Hub) HandleGames() {
    // make a list of all the players in queue
    // get 2 players from the queue and pass them to the NewGame
    // run the game i guess
    
    for {
        p1 := <-h.player_queue
        p2 := <-h.player_queue

        g := NewGame(p1, p2)
        go g.Run()
        h.games = append(h.games, g)

    }

}

