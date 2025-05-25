package game

import (
    "log"

    "tic_tac_toe.fun/users"
    "tic_tac_toe.fun/board"

	"github.com/gin-gonic/gin"
)

type Game struct {
    players [2]*Player
    b *board.Board
    p1move bool
}

func (g *Game) NewGame(c *gin.Context, u1 *users.User, u2 *users.User) error {
    var err error
    g.players[0], err = MakePlayer(c, u1)
    if err != nil {
        return err
    }
    g.players[1], err = MakePlayer(c, u2)
    if err != nil {
        return err
    }
    g.p1move = true // player 1 makes the first move
    // update the player stat: games played
    return nil
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
                    g.players[0].board_state <- b_state
                    g.players[1].board_state <- b_state
                    if res := g.b.CheckForWin(); res == true {
                        log.Println("PLAYER 1 WINS")
                        // do some other stuff ig
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
                    g.players[0].board_state <- b_state
                    g.players[1].board_state <- b_state
                    if res := g.b.CheckForWin(); res == true {
                        log.Println("PLAYER 2 WINS")
                        // do some other stuff ig
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

func HandleGame(c *gin.Context) {

}
