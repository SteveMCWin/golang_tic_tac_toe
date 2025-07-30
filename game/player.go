package game

import (
	"log"
	"strconv"
	"time"

	"tic_tac_toe.fun/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Player struct {
    u *models.User
    conn *websocket.Conn
    move chan []byte    // used to send messages to the game handler regarding the move the player made
    board_state chan []byte // used to updated the state of the board on the front end
    exited chan bool    // used to signal the player exited the game
}

func ConnectPlayerToSocket(hub *Hub, usr *models.User, c *gin.Context) error {

	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)   // set up websocket connection
	if err != nil {
		return err
	}

	new_player := &Player{u: usr, conn: connection, move: make(chan []byte), board_state: make(chan []byte), exited: make(chan bool)}

	game_mode, err := strconv.Atoi(c.Query("game_mode"))    // the url is expected to contain a game_mode the player wants to play
	if err != nil { // if the value of game_mode is invalid, don't start the game
		log.Printf("Unexpected game mode: %s", c.Query("game_mode"))
		new_player.conn.Close()
		return err
	}

	hub.AddPlayer(new_player, game_mode)

	return nil
}

// func (p *Player) UpdatePlayerStats() {  // intended to be called at the end of the game
//     p.u.UpdateGameStats()
// }

func (p *Player) ListenToSocket() { // listens to the response from the websocket (the front end)
    defer func() {
        log.Printf("Player %s exited", p.u.UserName)
        p.exited <- true
        p.conn.Close()
    }()
    p.conn.SetReadLimit(maxMessageSize)
    p.conn.SetReadDeadline(time.Now().Add(pongWait))
    p.conn.SetPongHandler(func(string) error { p.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

    for {
        _, message, err := p.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        p.move <- message[:2]   // the expected message is supposed to be just 2 character corresponding to the position of the move the player made

    }
}

func (p *Player) ListenToServer() { // listens to the response of the server, which is a new board state that is used to update the front end ui
    ticker := time.NewTicker(pingPeriod)
	defer func() {
        p.exited <- true
		ticker.Stop()
		p.conn.Close()
    }()
    for {
        select {
        case new_state, ok := <- p.board_state:
            p.conn.SetWriteDeadline(time.Now().Add(writeWait))
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
            p.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

