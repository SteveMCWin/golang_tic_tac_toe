package game

import (
    "log"
    "time"
    "math"
    "net/http"
    "encoding/json"

    "tic_tac_toe.fun/board"

	"github.com/gin-gonic/gin"
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
    p1move bool
    game_started bool
}

func ServePlay(c *gin.Context) {
    game_mode := c.Query("game_mode")
    if game_mode == "" {
        game_mode = "0"
    }
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

                    g.updateBoardVisuals()

                    if g.b.Result != 0 {
                        g.checkWinner()
                        return
                    }

                    g.p1move = false
                    g.player_timers[1].Start()
                    log.Printf("Resumed p2 timer at %s", g.player_timers[1].TimeLeft.String())
                }
            } /*else {
                log.Println("IT'S PLAYER 2'S MOVE")
            }*/
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

                    g.updateBoardVisuals()

                    if g.b.Result != 0 {
                        g.checkWinner()
                        return
                    }

                    g.p1move = true
                    g.player_timers[0].Start()
                    log.Printf("Resumed p1 timer at %s", g.player_timers[0].TimeLeft.String())
                }
            } /*else {
                log.Println("IT'S PLAYER 2'S MOVE")
            }*/
        case _ = <- g.players[0].exited:
            g.b.Result = 'o'
            g.checkWinner()
            return
        case _ = <- g.players[1].exited:
            g.b.Result = 'x'
            g.checkWinner()
            return
        case _ = <-g.player_timers[0].Finished:
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

func (g *Game) updateBoardVisuals() {
    new_b_state := parseBoardToJSON(g.b.BoardState)
    g.players[0].board_state <- new_b_state
    g.players[1].board_state <- new_b_state
}

func (g *Game) checkWinner() {
    switch g.b.Result {
    case 'x':
        g.players[0].u.GamesWon += 1
        log.Printf("Player x wins!!!")
    case 'o':
        g.players[1].u.GamesWon += 1
        log.Printf("Player o wins!!!")
    default:
        log.Printf("Tie!!!")
    }

    g.updatePlayerElo()

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

