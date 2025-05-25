package game

import (
    "log"
    "time"
    // "bytes"

    "tic_tac_toe.fun/users"
    // "tic_tac_toe.fun/board"

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
    u *users.User
    conn *websocket.Conn
    move chan byte  // used to send messages to the game handler
    board_state chan []byte // used to updated the state of the board
}

func MakePlayer(c *gin.Context, usr *users.User) (*Player, error) {
    connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        // log.Println(err)
        return nil, err
    }
    return &Player{u: usr, conn: connection, move: make(chan byte), board_state: make(chan []byte)}, nil
}

func (p *Player) ListenToSocket() {
    defer func() {
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
        p.move <- message[0]
    }
}

func (p *Player) ListenToServer() {
    ticker := time.NewTicker(pingPeriod)
	defer func() {
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
            w.Write(new_state)

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

