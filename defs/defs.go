package defs

const (
	NO_WINNER = -1
	NO_USER_ID = 0
	COOKIE_MAX_AGE = 86400 * 30
	STARTING_ELO = 800
	DEFAULT_GAMES_PLAYED = 0
	DEFAULT_GAMES_WON = 0
	WILD_BOARD = '?'   // wildBoard means the player can play in any of the mini-boards
	BOARD_TIE = '_'    // used in the BigBoard's Result field if it comes to a tie
	BOARD_HISTORY_DELIMITER = byte(';')
	EMPTY_CELL = 0
)
