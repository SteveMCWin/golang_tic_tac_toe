package game

import (
    "log"
    "time"
    "math"
    "net/http"
    // "encoding/json"

    "tic_tac_toe.fun/board"

	"github.com/gin-gonic/gin"
)

type GameMode int

const (
    normal_5s GameMode = iota
    normal_10s
    normal_15s

    game_mode_size  // make sure to keep this the last on the list
)

type Game struct {
    players [2]*Player
    player_timers [2]*PlayerTimer
    // game_start_timer *PlayerTimer
    b *board.BigBoard
    p1move bool
    game_started bool
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
    g.b = &board.BigBoard{}
    g.b.Initialize()
    g.p1move = true // player 1 makes the first move
    g.game_started = false
    // update user stats
    g.players[0].u.GamesPlayed += 1
    g.players[1].u.GamesPlayed += 1

    // g.game_start_timer = MakePlayerTimer(3 * time.Second)

    switch mode {
    case normal_5s:
        g.player_timers[0] = MakePlayerTimer(5 * time.Second, 0 * time.Second)
        g.player_timers[1] = MakePlayerTimer(5 * time.Second, 0 * time.Second)
        log.Println("CREATED GAME DURATION 5s")
    case normal_10s:
        g.player_timers[0] = MakePlayerTimer(10 * time.Second, 0 * time.Second)
        g.player_timers[1] = MakePlayerTimer(10 * time.Second, 0 * time.Second)
        log.Println("CREATED GAME DURATION 10s")
    case normal_15s:
        g.player_timers[0] = MakePlayerTimer(15 * time.Second, 0 * time.Second)
        g.player_timers[1] = MakePlayerTimer(15 * time.Second, 0 * time.Second)
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
        // remember to check for a draw!
        select {
        case pos := <- g.players[0].move:
            if g.p1move == true && g.game_started == true {
                if len(pos) < 2 {
                    log.Printf("Expected position to be of length >= 2, got legnth of %d", len(pos))
                    continue
                }
                err := g.b.MakeMove(pos[0], pos[1], 'x')
                if err != nil {
                    log.Println(err)
                } else {
                    log.Printf("Paused p1 timer at %s", g.player_timers[0].TimeLeft.String())
                    g.player_timers[0].Pause()

                    // move this to a function of it's own
                    new_b_state := parseBoardToJSON(g.b.BoardState)
                    g.players[0].board_state <- new_b_state
                    g.players[1].board_state <- new_b_state

                    if g.b.Result != 0 {
                        g.playerWon(0)
                        return
                    }
                    // move this to a function of it's own

                    g.p1move = false
                    g.player_timers[1].Start()
                    log.Printf("Resumed p2 timer at %s", g.player_timers[1].TimeLeft.String())
                }
            } else {
                log.Println("IT'S PLAYER 2'S MOVE")
            }
        case pos := <- g.players[1].move:
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

                    new_b_state := parseBoardToJSON(g.b.BoardState)
                    g.players[0].board_state <- new_b_state
                    g.players[1].board_state <- new_b_state

                    if g.b.Result != 0 {
                        g.playerWon(0)
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

func (g *Game) updatePlayerElo() {
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

func parseBoardToJSON(b_state [][]byte) []byte {

    // TODO: Update this to handle a 2d array

    // var board []string
    // for _, cell := range b_state {
    //     if cell == 'x' {
    //         board = append(board, "x")
    //     } else if cell == 'o' {
    //         board = append(board, "o")
    //     } else {
    //         board = append(board, "")
    //     }
    // }
    //
    // msg := struct {
    //     Type string     `json:"type"`
    //     Board []string  `json:"board"`
    // }{
    //     Type    : "state",
    //     Board   : board,
    // }
    //
    // res, err := json.Marshal(msg)
    // if err != nil {
    //     log.Println("ERROR PARSING BOARD STATE TO JSON:")
    //     log.Println(err)
    // }
    //
    // return res
    return make([]byte, 0)
}

