package game

import (
    "log"
    "encoding/json"
    "time"

    "tic_tac_toe.fun/board"

	"github.com/gin-gonic/gin"
)

type Game struct {
    players [2]*Player
    player_timers [2]*PlayerTimer
    b *board.Board
    p1move bool
}

func Play(c *gin.Context) {
    c.File("templates/board.html")
}

// TODO: add time as a parameter
func NewGame(p1, p2 *Player) *Game {
    g := &Game{}
    g.players[0] = p1
    g.players[1] = p2
    g.b = &board.Board{}
    g.p1move = true // player 1 makes the first move
    // update user stats
    g.players[0].u.GamesPlayed += 1
    g.players[1].u.GamesPlayed += 1

    g.player_timers[0] = MakePlayerTimer(10 * time.Second)
    g.player_timers[1] = MakePlayerTimer(10 * time.Second)

    return g
}

func (g *Game) Run() {
    go g.players[0].ListenToSocket()
    go g.players[1].ListenToSocket()

    go g.players[0].ListenToServer()
    go g.players[1].ListenToServer()

    g.player_timers[0].Start()

    defer g.players[0].UpdatePlayerStats()
    defer g.players[1].UpdatePlayerStats()

    for {
        select {
        case pos := <- g.players[0].move:
            if g.p1move == true {
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
            if g.p1move == false {
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

