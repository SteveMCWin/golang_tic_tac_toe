package game

import (
    "log"

    "tic_tac_toe.fun/users"
    "tic_tac_toe.fun/board"
	"github.com/gin-gonic/gin"
)

type Player struct {
    u *User
    move chan byte
}

type Game struct {
    players [2]*Player
    b *Board
    p1move bool
}

func (g *Game) NewGame(u1 *User, u2 *User) {
    g.players[0] = &Player{u1, make(chan byte)}
    g.players[1] = &Player{u2, make(chan byte)}
    g.p1move = true // player 1 makes the first move
    // update the player stat: games played
}

func (g *Game) Run() {
    for {
        select {
        case pos := <- players[0].move:
            if p1move == true {
                b.MakeMove(pos, byte('x'))
                if res := b.CheckForWin(); res == true {
                    log.Println("PLAYER 1 WINS")
                    // do some other stuff ig
                    return
                }
                p1move = false
            } else {
                log.Println("IT'S PLAYER 2'S MOVE")
            }
        }
        case pos := <- players[1].move:
            if p1move == false {
                b.MakeMove(pos, byte('o'))
                if res := b.CheckForWin(); res == true {
                    log.Println("PLAYER 2 WINS")
                    // do some other stuff ig
                    return
                }
                p1move = true
            } else {
                log.Println("IT'S PLAYER 1'S MOVE")
            }
        }
    }
}

func HandleGame(c *gin.Context) {

}

