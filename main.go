package main

import (
	"net/http"

	"tic_tac_toe.fun/game"
	"tic_tac_toe.fun/handlers"
	"tic_tac_toe.fun/models"
)

func main() {

	db := &models.DataBase{}
    db.InitDatabase()

	MostEloLb := &models.LeaderBoard{}
	MostEloLb.InitLeaderBoard()
	go MostEloLb.RunLeaderBoard(db)

    hub := game.MakeHub()
    go hub.HandleGames()

	handler := handlers.SetUpRouter(db, MostEloLb, hub)

	http.ListenAndServe(":5000", handler)
}

