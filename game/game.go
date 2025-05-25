package game

import (
    "log"
    "encoding/json"

    "tic_tac_toe.fun/board"

	"github.com/gin-gonic/gin"
)

type Game struct {
    players [2]*Player
    b *board.Board
    p1move bool
}

func NewGame(p1, p2 *Player) *Game {
    g := &Game{}
    g.players[0] = p1
    g.players[1] = p2
    g.b = &board.Board{}
    g.p1move = true // player 1 makes the first move
    // TODO: update the player stat: games played
    return g
}

func (g *Game) Run() {
    go g.players[0].ListenToSocket()
    go g.players[1].ListenToSocket()

    go g.players[0].ListenToServer()
    go g.players[1].ListenToServer()

    for {
        select {
        case pos := <- g.players[0].move:
            if g.p1move == true {
                b_state, err := g.b.MakeMove(pos, byte('x'))
                if err != nil {
                    log.Println(err)
                } else {
                    g.players[0].board_state <- parseBoardToJSON(b_state)
                    g.players[1].board_state <- parseBoardToJSON(b_state)
                    if res := g.b.CheckForWin(); res == true {
                        log.Println("PLAYER 1 WINS")
                        // TODO: do some other stuff ig
                        return
                    }
                    g.p1move = false
                }
            } else {
                log.Println("IT'S PLAYER 2'S MOVE")
            }
        case pos := <- g.players[1].move:
            if g.p1move == false {
                b_state, err := g.b.MakeMove(pos, byte('o'))
                if err != nil {
                    log.Println(err)
                } else {
                    g.players[0].board_state <- parseBoardToJSON(b_state)
                    g.players[1].board_state <- parseBoardToJSON(b_state)
                    if res := g.b.CheckForWin(); res == true {
                        log.Println("PLAYER 2 WINS")
                        // TODO: do some other stuff ig
                        return
                    }
                    g.p1move = true
                }
            } else {
                log.Println("IT'S PLAYER 1'S MOVE")
            }
        }
    }
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
    // log.Println("No encoding:", b_state)
    // log.Println("Json encoding:", string(res))
    if err != nil {
        log.Println("ERROR PARSING BOARD STATE TO JSON:")
        log.Println(err)
    }
    return res
}

func Play(c *gin.Context) {
    c.File("templates/board.html")
}
