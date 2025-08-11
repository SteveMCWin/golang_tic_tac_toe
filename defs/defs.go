// Package defs contains commonly used constants
package defs

import "time"

const (
	NO_USER_ID = 0

	COOKIE_MAX_AGE = 86400 * 30

	STARTING_ELO = 800
	DEFAULT_GAMES_PLAYED = 0
	DEFAULT_GAMES_WON = 0
	NO_WINNER = -1
	DELETED_USER_ID = -1

	WILD_BOARD = '?'   // wildBoard means the player can play in any of the mini-boards
	BOARD_TIE = '_'    // used in the BigBoard's Result field if it comes to a tie
	BOARD_HISTORY_DELIMITER = byte(';')
	EMPTY_CELL = 0
	PREV_MOVE = 60
	NEXT_MOVE = 62

	LEADERBOARD_UPDATE_DURATION = 10 * time.Second

	WRITE_WAIT = 10 * time.Second // Time allowed to write a message to the peer in a websocket connection.
	PONG_WAIT = 60 * time.Second // Time allowed to read the next pong message from the peer in a websocket connection.
	PING_PERIOD = (PONG_WAIT * 9) / 10 // Send pings to peer with this period in a websocket connection. Must be less than pongWait.
	MAX_WEBSOCKET_MESSAGE_SIZE = 512 // Maximum message size allowed from peer in a websocket connection.
)
