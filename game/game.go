package game

import (
    "log"
    "encoding/json"
    "time"
    "net/http"

    "tic_tac_toe.fun/board"

	"github.com/gin-gonic/gin"
)

type Game struct {
    players [2]*Player
    player_timers [2]*PlayerTimer
    // game_start_timer *PlayerTimer
    b *board.Board
    p1move bool
    game_started bool
    move_counter int
}

func ServePlay(c *gin.Context) {
    game_mode := c.Query("game_mode")
    log.Println("game_mode in Play:", game_mode)
    // c.File("templates/board.html")
    c.HTML(http.StatusOK, "board.html", gin.H{
        "game_mode": game_mode,
    })
}

func NewGame(p1, p2 *Player, mode GameMode) *Game {
    g := &Game{}
    g.players[0] = p1
    g.players[1] = p2
    g.b = &board.Board{}
    g.p1move = true // player 1 makes the first move
    g.game_started = false
    g.move_counter = 0
    // update user stats
    g.players[0].u.GamesPlayed += 1
    g.players[1].u.GamesPlayed += 1

    // g.game_start_timer = MakePlayerTimer(3 * time.Second)

    switch mode {
    case normal_5s:
        g.player_timers[0] = MakePlayerTimer(5 * time.Second)
        g.player_timers[1] = MakePlayerTimer(5 * time.Second)
        log.Println("CREATED GAME DURATION 5s")
    case normal_10s:
        g.player_timers[0] = MakePlayerTimer(10 * time.Second)
        g.player_timers[1] = MakePlayerTimer(10 * time.Second)
        log.Println("CREATED GAME DURATION 10s")
    case normal_15s:
        g.player_timers[0] = MakePlayerTimer(15 * time.Second)
        g.player_timers[1] = MakePlayerTimer(15 * time.Second)
        log.Println("CREATED GAME DURATION 15s")
    }

    return g
}

func (g *Game) Run() {
    go g.players[0].ListenToSocket()
    go g.players[0].ListenToServer()

    go g.players[1].ListenToSocket()
    go g.players[1].ListenToServer()

    defer g.players[0].UpdatePlayerStats()
    defer g.players[1].UpdatePlayerStats()

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
    }()

    for {
        if g.move_counter >= 9 {
            log.Printf("IT'S A DRAW")
            return
        }
        select {
        case pos := <- g.players[0].move:
            if g.p1move == true && g.game_started == true {
                g.move_counter += 1
                g.player_timers[0].Pause()
                log.Printf("Paused p1 timer at %s", g.player_timers[0].TimeLeft.String())
                b_state, err := g.b.MakeMove(pos, byte('x'))
                if err != nil {
                    log.Println(err)
                } else {
                    g.players[0].board_state <- parseBoardToJSON(b_state)
                    g.players[1].board_state <- parseBoardToJSON(b_state)
                    if res := g.b.CheckForWin(); res == true {
                        g.playerWon(0)
                        return
                    }
                    g.p1move = false
                    g.player_timers[1].Start()
                    log.Printf("Resumed p2 timer at %s", g.player_timers[1].TimeLeft.String())
                }
            } else {
                log.Println("IT'S PLAYER 2'S MOVE")
            }
        case pos := <- g.players[1].move:
            if g.p1move == false && g.game_started == true {
                g.move_counter += 1
                g.player_timers[1].Pause()
                log.Printf("Paused p2 timer at %s", g.player_timers[1].TimeLeft.String())
                b_state, err := g.b.MakeMove(pos, byte('o'))
                if err != nil {
                    log.Println(err)
                } else {
                    g.players[0].board_state <- parseBoardToJSON(b_state)
                    g.players[1].board_state <- parseBoardToJSON(b_state)
                    if res := g.b.CheckForWin(); res == true {
                        g.playerWon(1)
                        return
                    }
                    g.p1move = true
                    g.player_timers[0].Start()
                    log.Printf("Resumed p1 timer at %s", g.player_timers[0].TimeLeft.String())
                }
            } else {
                log.Println("IT'S PLAYER 1'S MOVE")
            }
        case _ = <- g.players[0].exited:
            g.playerWon(1)
            return
        case _ = <- g.players[1].exited:
            g.playerWon(0)
            return
        case _ = <-g.player_timers[0].Finished:
            g.playerWon(1)
            return
        case _ = <-g.player_timers[1].Finished:
            g.playerWon(0)
            return
        }
    }
}

func (g *Game) playerWon(player_idx int) {
    g.players[player_idx].u.GamesWon += 1
    log.Printf("PLAYER %d WINS\n", player_idx+1)
}

func parseBoardToJSON(b_state []byte) []byte {

    var board []string
    for _, cell := range b_state {
        if cell == 'x' {
            board = append(board, "x")
        } else if cell == 'o' {
            board = append(board, "o")
        } else {
            board = append(board, "")
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

