package game

import (
    // "log"
)

type Hub struct {
    Games []*Game
    PlayerQueue chan *Player
}

func (h *Hub) HandleGames() {
    // make a list of all the players in queue
    // get 2 players from the queue and pass them to the NewGame
    // run the game i guess
    
    for {
        // log.Println("Waiting for playa oneeeee")
        p1 := <-h.PlayerQueue
        // log.Println("Waiting for playa twooooo")
        p2 := <-h.PlayerQueue

        // log.Println("MADE DA GAMEEEE")
        g := NewGame(p1, p2)
        go g.Run()
        h.Games = append(h.Games, g)

    }

}

