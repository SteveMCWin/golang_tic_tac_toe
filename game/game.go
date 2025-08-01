package game

import (
	"encoding/json"
	"log"
	"math"
	"strconv"
	"time"

	"tic_tac_toe.fun/board"
	"tic_tac_toe.fun/models"
)

type GameMode int

const (
    normal_180s GameMode = iota
    normal_300s
    normal_600s

    fischer_60_3s
    fischer_60_5s
    fischer_180_3s
    fischer_180_5s

    game_mode_size  // make sure to keep this the last on the list
)

type Game struct {
    players [2]*Player
    player_timers [2]*PlayerTimer
    b *board.BigBoard
    p1move bool // used to prevent the other player playing when it's not their turn
    game_started bool
	GameRecord string
}

func NewGame(p1, p2 *Player, mode GameMode) *Game {
    g := &Game{}
    g.players[0] = p1
    g.players[1] = p2
    g.b = &board.BigBoard{}
    g.b.Initialize()
    g.p1move = true // player 1 is always 'x' so they make the first move
    g.game_started = false
    // update user stats
    g.players[0].u.GamesPlayed += 1
    g.players[1].u.GamesPlayed += 1
	g.GameRecord = ""

    switch mode {
    case normal_180s:
        g.player_timers[0] = MakePlayerTimer(180 * time.Second, 0 * time.Second)
        g.player_timers[1] = MakePlayerTimer(180 * time.Second, 0 * time.Second)
        log.Println("CREATED GAME DURATION 180s")
    case normal_300s:
        g.player_timers[0] = MakePlayerTimer(300 * time.Second, 0 * time.Second)
        g.player_timers[1] = MakePlayerTimer(300 * time.Second, 0 * time.Second)
        log.Println("CREATED GAME DURATION 300s")
    case normal_600s:
        g.player_timers[0] = MakePlayerTimer(600 * time.Second, 0 * time.Second)
        g.player_timers[1] = MakePlayerTimer(600 * time.Second, 0 * time.Second)
        log.Println("CREATED GAME DURATION 600s")
    case fischer_60_3s:
        g.player_timers[0] = MakePlayerTimer(60 * time.Second, 3 * time.Second)
        g.player_timers[1] = MakePlayerTimer(60 * time.Second, 3 * time.Second)
        log.Println("CREATED GAME DURATION 60s+3s")
    case fischer_60_5s:
        g.player_timers[0] = MakePlayerTimer(60 * time.Second, 5 * time.Second)
        g.player_timers[1] = MakePlayerTimer(60 * time.Second, 5 * time.Second)
        log.Println("CREATED GAME DURATION 60s+3s")
    case fischer_180_3s:
        g.player_timers[0] = MakePlayerTimer(180 * time.Second, 3 * time.Second)
        g.player_timers[1] = MakePlayerTimer(180 * time.Second, 3 * time.Second)
        log.Println("CREATED GAME DURATION 180s+3s")
    case fischer_180_5s:
        g.player_timers[0] = MakePlayerTimer(180 * time.Second, 5 * time.Second)
        g.player_timers[1] = MakePlayerTimer(180 * time.Second, 5 * time.Second)
        log.Println("CREATED GAME DURATION 180s+5s")
    }

    return g
}

func (g *Game) Run(db *models.DataBase) {  // listens for user and server actions and udpdates the game accordingly, also handles game logic such as switching player turns etc
    go g.players[0].ListenToSocket()
    go g.players[0].ListenToServer()

    go g.players[1].ListenToSocket()
    go g.players[1].ListenToServer()

    defer db.UpdateGameStats(g.players[0].u) // makes sure that the player stats stored only when the game finishes
    defer db.UpdateGameStats(g.players[1].u)

	defer func() {
		err := db.StoreGameRecord(g.GameRecord, g.players[0].u.Id, g.players[1].u.Id)
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
        g.player_timers[0].Start()
    }() // this being in a goroutine prevents the players from filling up the input channel before the game starts which may lead to the 'o' player going first

    for {
        select {
        case pos := <- g.players[0].move:
            if g.p1move == true && g.game_started == true {
                // the position sent is a byte slice which is supposed to be of length 2, pos[0] being the mini board idx, pos[1] being the cell idx
                if len(pos) < 2 {
                    log.Printf("Expected position to be of length >= 2, got legnth of %d", len(pos))
                    continue
                }

                err := g.b.MakeMove(pos[0], pos[1], 'x')    // this updates the board back-end logic/state
                if err != nil {
                    log.Println(err)
                } else {
                    log.Printf("Paused p1 timer at %s", g.player_timers[0].TimeLeft.String())
                    // if the move was deemed valid and all, handle the timers, send the messages to update the front-end based on the back-end
                    g.player_timers[0].Pause()

					g.RecordMove("x", pos[0], pos[1])

                    g.updateBoardVisuals()
                    if g.b.Result != 0 {    // handles the game being finished and exits the function
                        g.checkWinner()
                        return
                    }

                    g.p1move = false
                    g.player_timers[1].Start()

                    log.Printf("Resumed p2 timer at %s", g.player_timers[1].TimeLeft.String())
                }
            }
        case pos := <- g.players[1].move:   // same as the case above, just for the other player
            if g.p1move == false && g.game_started == true {
                if len(pos) < 2 {
                    log.Printf("Expected position to be of length >= 2, got legnth of %d", len(pos))
                    continue
                }
                err := g.b.MakeMove(pos[0], pos[1], 'o')
                if err != nil {
                    log.Println(err)
                } else {
                    log.Printf("Paused p2 timer at %s", g.player_timers[1].TimeLeft.String())
                    g.player_timers[1].Pause()

					g.RecordMove("o", pos[0], pos[1])

                    g.updateBoardVisuals()

                    if g.b.Result != 0 {
                        g.checkWinner()
                        return
                    }

                    g.p1move = true
                    g.player_timers[0].Start()
                    log.Printf("Resumed p1 timer at %s", g.player_timers[0].TimeLeft.String())
                }
            }
        case _ = <- g.players[0].exited:    // in case someone exits before the game is finished, finish the game and make the other player the winner
            g.b.Result = 'o'
            g.checkWinner()
            return
        case _ = <- g.players[1].exited:
            g.b.Result = 'x'
            g.checkWinner()
            return
        case _ = <-g.player_timers[0].Finished: // in case someone's timer runs out, the other player wins
            g.b.Result = 'o'
            g.checkWinner()
            return
        case _ = <-g.player_timers[1].Finished:
            g.b.Result = 'x'
            g.checkWinner()
            return
        }
    }
}

func (g *Game) RecordMove(player string, big_pos, small_pos byte) {
    if big_pos >= '0' && big_pos <= '9' {   // the positions passed in are probably ascii characters representing digis
        big_pos = big_pos - '0'
    }

    if small_pos >= '0' && small_pos <= '9' {
        small_pos = small_pos - '0'
    }

	g.GameRecord += player + strconv.Itoa(int(big_pos)) + strconv.Itoa(int(small_pos)) + ";"
}

func (g *Game) updateBoardVisuals() {   // sends the board's back-end data to the front-end to update the ui
    new_b_state := parseBoardToJSON(g.b.BoardState) // the front-end expects the board state to be in json format
    g.players[0].board_state <- new_b_state
    g.players[1].board_state <- new_b_state
}

func (g *Game) checkWinner() {  // the winner is determined from the boards Result field
    switch g.b.Result {
    case 'x':
        g.players[0].u.GamesWon += 1
        log.Printf("Player x wins!!!")
		g.GameRecord += "X;"
    case 'o':
        g.players[1].u.GamesWon += 1
        log.Printf("Player o wins!!!")
		g.GameRecord += "O;"
    default:
        log.Printf("Tie!!!")
    }

    g.calculateNewPlayerElo()   // needs to be called after a winner is determined

}

func (g *Game) calculateNewPlayerElo() {    // calculates each players new elo based on the outcome of the game, doesn't store it in db itself, only in the structs
    Qa := math.Pow(10.0, float64(g.players[0].u.Elo) / 400.0)
    Qb := math.Pow(10.0, float64(g.players[1].u.Elo) / 400.0)

    Ea := Qa/(Qa+Qb)
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

func parseBoardToJSON(b_state [][]byte) []byte {    // note that the board state is turned into a normal byte array instead of a 2d byte array

    var board []string
    for _, b := range b_state {
        for _, cell := range b{
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
        Type string     `json:"type"`
        Board []string  `json:"board"`
    }{
        Type    : "state",
        Board   : board,
    }

    res, err := json.Marshal(msg)
    if err != nil {
        log.Println("ERROR PARSING BOARD STATE TO JSON:")
        log.Println(err)
    }

    return res
}

