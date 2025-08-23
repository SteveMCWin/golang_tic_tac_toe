// Package game contains types and functions that drive the logic of a realtime game of 2D Tic Tac Toe
package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	"time"

	"tic_tac_toe.fun/board"
	"tic_tac_toe.fun/defs"
	"tic_tac_toe.fun/models"
)

// GameMode is used to identify what type of game users want to play.
type GameMode int

const (
	normal_180s GameMode = iota
	normal_300s
	normal_600s

	fischer_60_3s
	fischer_60_5s
	fischer_180_3s
	fischer_180_5s

	game_mode_size // make sure to keep this the last on the list
)

// Game holds the data that is needed to handle the logic of the game.
type Game struct {
	players      [2]*Player
	b            *board.BigBoard
	p1_move      bool // used to prevent the other player playing when it's not their turn
	game_started bool
	gameRecord   models.GameRecord
}

// NewGame initializes a new game based on the GameMode passed in.
func NewGame(p1, p2 *Player, mode GameMode) *Game {
	g := &Game{}
	g.players[0] = p1
	g.players[1] = p2
	g.b = &board.BigBoard{}
	g.b.Initialize()
	g.p1_move = true // player 1 is always 'x' so they make the first move
	g.game_started = false
	// update user stats
	g.gameRecord = models.GameRecord{
		U1:           g.players[0].u,
		U2:           g.players[1].u,
		DateRecorded: time.Now(),
		History:       "", // start with empty recording because the game just started (duuh)
	}

	// NOTE: Could be moved to the player
	switch mode {
	case normal_180s:
		g.players[0].timer = MakePlayerTimer(180*time.Second, 0*time.Second)
		g.players[1].timer = MakePlayerTimer(180*time.Second, 0*time.Second)
		log.Println("CREATED GAME DURATION 180s")
	case normal_300s:
		g.players[0].timer = MakePlayerTimer(300*time.Second, 0*time.Second)
		g.players[1].timer = MakePlayerTimer(300*time.Second, 0*time.Second)
		log.Println("CREATED GAME DURATION 300s")
	case normal_600s:
		g.players[0].timer = MakePlayerTimer(600*time.Second, 0*time.Second)
		g.players[1].timer = MakePlayerTimer(600*time.Second, 0*time.Second)
		log.Println("CREATED GAME DURATION 600s")
	case fischer_60_3s:
		g.players[0].timer = MakePlayerTimer(60*time.Second, 3*time.Second)
		g.players[1].timer = MakePlayerTimer(60*time.Second, 3*time.Second)
		log.Println("CREATED GAME DURATION 60s+3s")
	case fischer_60_5s:
		g.players[0].timer = MakePlayerTimer(60*time.Second, 5*time.Second)
		g.players[1].timer = MakePlayerTimer(60*time.Second, 5*time.Second)
		log.Println("CREATED GAME DURATION 60s+5s")
	case fischer_180_3s:
		g.players[0].timer = MakePlayerTimer(180*time.Second, 3*time.Second)
		g.players[1].timer = MakePlayerTimer(180*time.Second, 3*time.Second)
		log.Println("CREATED GAME DURATION 180s+3s")
	case fischer_180_5s:
		g.players[0].timer = MakePlayerTimer(180*time.Second, 5*time.Second)
		g.players[1].timer = MakePlayerTimer(180*time.Second, 5*time.Second)
		log.Println("CREATED GAME DURATION 180s+5s")
	}

	return g
}

// Run listens for any activity from the players and updates the board accordingly. It is responsible for most of the game logic, including tracking whose move it is,
// managing player timers, checking if a player won, etc.
// When this function is exited, it calls makes sure the player game-related stats get updated, and that the game record gets saved to the database.
func (g *Game) Run(db *models.DataBase) {
	go g.players[0].ListenToSocket()
	go g.players[0].ListenToServer()

	go g.players[1].ListenToSocket()
	go g.players[1].ListenToServer()

	g.players[0].u.NumOfGamesPlayed += 1
	g.players[1].u.NumOfGamesPlayed += 1

	defer db.UpdateGameStats(g.players[0].u) // makes sure that the player stats stored only when the game finishes
	defer db.UpdateGameStats(g.players[1].u)

	// NOTE: doing this here instead for example the checkWinner function because here we have access to db
	defer func() {
		log.Println("Storing board history")
		g.boardHistoryToString()
		err := db.StoreGameRecord(g.gameRecord)
		if err != nil {
			log.Println("Error storing game record!")
			log.Println(err)
		}
	}()

	go func() {
		log.Println("GAME STARTS IN")
		log.Println("3")
		time.Sleep(time.Second)
		log.Println("2")
		time.Sleep(time.Second)
		log.Println("1")
		time.Sleep(time.Second)
		log.Println("GOOO")
		g.game_started = true
		g.players[0].timer.Start()
		g.updateBoardVisuals(defs.MSG_TYPE_START)
	}() // this being in a goroutine prevents the players from filling up the input channel before the game starts which may lead to the 'o' player going first

	for {
		select {
		case pos := <-g.players[0].move:
			if g.p1_move == true && g.game_started == true {
				// the position sent is a byte slice which is supposed to be of length 2, pos[0] being the mini board idx, pos[1] being the cell idx
				if len(pos) < 2 {
					log.Printf("Expected position to be of length >= 2, got legnth of %d", len(pos))
					continue
				}

				m := board.Move {
					BigPos: pos[0],
					SmallPos: pos[1],
					Player: 'x',
				}

				err := g.b.MakeMove(m) // this updates the board back-end logic/state
				if err != nil {
					log.Println(err)
				} else {
					// if the move was deemed valid and all, handle the timers, send the messages to update the front-end based on the back-end
					g.players[0].timer.Pause()

					g.updateBoardVisuals(defs.MSG_TYPE_STATE)
					if g.b.Result != 0 { // handles the game being finished and exits the function
						log.Println("The backend recognized the winner!!")
						g.checkWinner()
						return
					}

					g.p1_move = false
					g.players[1].timer.Start()

				}
			}
		case pos := <-g.players[1].move: // same as the case above, just for the other player
			if g.p1_move == false && g.game_started == true {
				if len(pos) < 2 {
					log.Printf("Expected position to be of length >= 2, got legnth of %d", len(pos))
					continue
				}

				m := board.Move {
					BigPos: pos[0],
					SmallPos: pos[1],
					Player: 'o',
				}

				err := g.b.MakeMove(m)
				if err != nil {
					log.Println(err)
				} else {
					g.players[1].timer.Pause()

					g.updateBoardVisuals(defs.MSG_TYPE_STATE)

					if g.b.Result != 0 {
						g.checkWinner()
						return
					}

					g.p1_move = true
					g.players[0].timer.Start()
				}
			}
		case _ = <-g.players[0].exited: // in case someone exits before the game is finished, finish the game and make the other player the winner
			g.b.Result = 'o'
			g.checkWinner()
			return
		case _ = <-g.players[1].exited:
			g.b.Result = 'x'
			g.checkWinner()
			return
		case _ = <-g.players[0].timer.Finished: // in case someone's timer runs out, the other player wins
			g.b.Result = 'o'
			g.checkWinner()
			return
		case _ = <-g.players[1].timer.Finished:
			g.b.Result = 'x'
			g.checkWinner()
			return
		}
	}
}

func (g *Game) updateBoardVisuals(msg_type string) { // sends the board's back-end data to the front-end to update the ui
	new_game_state := g.parseGameStateToJSON(msg_type) // the front-end expects the board state to be in json format
	g.players[0].game_state <- new_game_state
	g.players[1].game_state <- new_game_state
}

func (g *Game) checkWinner() { // the winner is determined from the boards Result field
	fmt.Println("CheckWinner called!!!!")
	switch g.b.Result {
	case 'x':
		g.players[0].u.NumOfGamesWon += 1
		log.Printf("Player x wins!!!")
		g.gameRecord.Winner = 'x'
	case 'o':
		g.players[1].u.NumOfGamesWon += 1
		log.Printf("Player o wins!!!")
		g.gameRecord.Winner = 'o'
	default:
		log.Printf("Tie!!!")
		g.gameRecord.Winner = defs.BOARD_TIE
	}

	g.calculateNewPlayerElo() // needs to be called after a winner is determined

}

func (g *Game) calculateNewPlayerElo() { // calculates each players new elo based on the outcome of the game, doesn't store it in db itself, only in the structs
	Qa := math.Pow(10.0, float64(g.players[0].u.Elo)/400.0)
	Qb := math.Pow(10.0, float64(g.players[1].u.Elo)/400.0)

	Ea := Qa / (Qa + Qb)
	Eb := 1.0 - Ea

	const K = 32

	var Sa float64
	var Sb float64

	switch g.b.Result {
	case 'x':
		Sa = 1.0
		Sb = 0.0
	case 'o':
		Sa = 0.0
		Sb = 1.0
	default:
		Sa = 0.5
		Sb = 0.5
	}

	g.players[0].u.Elo += int(math.Round(K * (Sa - Ea)))
	g.players[1].u.Elo += int(math.Round(K * (Sb - Eb)))
}

func (g *Game) parseGameStateToJSON(msg_type string) []byte { // note that the board state is turned into a normal byte array instead of a 2d byte array

	var board []string
	for _, b := range g.b.BoardState {
		for _, cell := range b {
			if cell == 'x' {
				board = append(board, "x")
			} else if cell == 'o' {
				board = append(board, "o")
			} else {
				board = append(board, "")
			}
		}
	}

	msg := struct {
		Type           string   `json:"type"`
		Board          []string `json:"board"`
		CompleteBoards []byte   `json:"complete_boards"`
		P1_time        int64    `json:"p1_time"`
		P2_time        int64    `json:"p2_time"`
		P1_move        bool     `json:"p1_move"`
	}{
		Type:    msg_type,
		Board:   board,
		CompleteBoards: make([]byte, 0),
		P1_time: g.players[0].timer.TimeLeft.Milliseconds(),
		P2_time: g.players[1].timer.TimeLeft.Milliseconds(),
		P1_move: g.p1_move,
	}

	for _, b := range g.b.Boards {
		msg.CompleteBoards = append(msg.CompleteBoards, b.Result)
	}

	res, err := json.Marshal(msg)
	if err != nil {
		log.Println("ERROR PARSING BOARD STATE TO JSON:")
		log.Println(err)
	}

	return res
}

func (g *Game) boardHistoryToString() {
	for g.b.History.Len() > 0 {
		m := g.b.History.Top()
		if m == nil {
			return
		}
		g.b.History.Pop()
		m_str := string([]byte{ m.Player, m.BigPos + '0', m.SmallPos + '0', defs.BOARD_HISTORY_DELIMITER })
		g.gameRecord.History = m_str + g.gameRecord.History // note that we must prepend the last-read move since we are working with a stack
	}
}
