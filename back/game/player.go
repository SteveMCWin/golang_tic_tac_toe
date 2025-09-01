package game

import (
	"log"
	"strconv"
	"time"

	"tic_tac_toe.fun/defs"
	"tic_tac_toe.fun/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Player represents a user inside a game.
type Player struct {
    u *models.User
    conn *websocket.Conn
    move chan []byte    // used to send messages to the game handler regarding the move the player made
    game_state chan []byte // used to updated the state of the board on the front end
    exited chan bool    // used to signal the player exited the game
	timer *PlayerTimer
}

// ConnectPlayerToSocket  upgrades the users connection to a websocket connection and puts them in a queue for a game.
func ConnectPlayerToSocket(hub *Hub, usr *models.User, c *gin.Context) error {
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)   // set up websocket connection
	if err != nil {
		return err
	}

	new_player := &Player{u: usr, conn: connection, move: make(chan []byte), game_state: make(chan []byte), exited: make(chan bool)} // NOTE: timer is set up in at game start

	game_mode, err := strconv.Atoi(c.Query("game_mode"))    // the url is expected to contain a game_mode the player wants to play
	if err != nil { // if the value of game_mode is invalid, don't start the game
		log.Printf("Unexpected game mode: %s", c.Query("game_mode"))
		new_player.conn.Close()
		return err
	}

	hub.AddPlayer(new_player, GameMode(game_mode))

	return nil
}

// ListenToSocket is a layer between the player and the server. Once a player makes a move, this function passes it to the backend to get handled.
func (p *Player) ListenToSocket() { 
    defer func() {
        log.Printf("Player %s exited", p.u.UserName)
        p.exited <- true
        p.conn.Close()
    }()
    p.conn.SetReadLimit(defs.MAX_WEBSOCKET_MESSAGE_SIZE)
    p.conn.SetReadDeadline(time.Now().Add(defs.PONG_WAIT))
    p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(defs.PONG_WAIT)); return nil })

    for {
        _, message, err := p.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        p.move <- message[:2]   // the expected message is supposed to be just 2 characters corresponding to the position of the move the player made

    }
}

// ListenToServer is a layer between the server and the player. Once the server valides and realizes a board move, it lets the players know so their front end ui can get updated.
func (p *Player) ListenToServer() { // listens to the response of the server, which is a new board state that is used to update the front-end ui
    ticker := time.NewTicker(defs.PING_PERIOD)
	defer func() {
        p.exited <- true
		ticker.Stop()
		p.conn.Close()
    }()
    for {
        select {
        case new_state, ok := <- p.game_state:
            p.conn.SetWriteDeadline(time.Now().Add(defs.WRITE_WAIT))
            if !ok {
                p.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := p.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(new_state)  // sends the board state to the websocket

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            p.conn.SetWriteDeadline(time.Now().Add(defs.WRITE_WAIT))
            if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

