package game

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"tic_tac_toe.fun/board"
	"tic_tac_toe.fun/models"
)

type GameReplay struct {
	b *board.BigBoard
	gameRecord *models.GameRecord
}

var cachedReplays map[int]GameReplay

func InitGameReplay(user_id int, rec *models.GameRecord) error {
	if cachedReplays == nil {
		cachedReplays = make(map[int]GameReplay)
	}

	new_board := board.BigBoard{}
	new_board.Initialize()

	moves, err := RecordHistoryToMoves(rec.History)
	if err != nil {
		log.Println("Error while converting record history to moves")
		return err
	}

	for i := len(moves)-1; i >= 0; i-- { // push moves in reverse order since the stack also reverses them
		new_board.RedoStack.Push(moves[i])
	}

	cachedReplays[user_id] = GameReplay{ b: &new_board, gameRecord: rec }

	return nil
}

func RecordHistoryToMoves(history string) ([]board.Move, error) {
	moves_str := strings.Split(history[:len(history)-1], ";")

	moves := make([]board.Move, len(moves_str))

	for i, move_str := range moves_str {
		if len(move_str) != 3 {
			log.Println("len(move_str):", len(move_str))
			log.Println("move_str:", move_str)
			return nil, errors.New("ERROR: expected len of move_str is 3")
		}

		moves[i] = board.Move {
			Player: move_str[0],
			BigPos: move_str[1] - '0',
			SmallPos: move_str[2] - '0',
		}
	}

	return moves, nil

}

func ReplayNextMove(user_id int) ([]byte, error) {

	replay := cachedReplays[user_id]
	err := replay.b.RedoLastMove()
	if err != nil {
		return nil, err
	}

	return parseReplayStateToJSON(&replay)
}

func ReplayPrevMove(user_id int) ([]byte, error) {

	replay := cachedReplays[user_id]
	err := replay.b.UndoLastMove()
	if err != nil {
		return nil, err
	}

	return parseReplayStateToJSON(&replay)
}

func parseReplayStateToJSON(replay *GameReplay) ([]byte, error) {
	var board []string
	for _, b := range replay.b.BoardState {
		for _, cell := range b {
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
		Type           string   `json:"type"`
		Board          []string `json:"board"`
		CompleteBoards []byte   `json:"complete_boards"`
	}{
		Type:    "state",
		Board:   board,
		CompleteBoards: make([]byte, 0),
	}

	for _, b := range replay.b.Boards {
		msg.CompleteBoards = append(msg.CompleteBoards, b.Result)
	}

	res, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return res, nil
}
