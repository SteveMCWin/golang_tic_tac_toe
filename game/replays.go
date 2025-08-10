package game

import (
	"encoding/json"

	"tic_tac_toe.fun/board"
	"tic_tac_toe.fun/models"
)

type GameReplay struct {
	b *board.BigBoard
	gameRecord *models.GameRecord
}

var cachedReplays map[int]GameReplay

func InitGameReplay(user_id int, rec *models.GameRecord) {
	if cachedReplays == nil {
		cachedReplays = make(map[int]GameReplay)
	}

	new_board := board.BigBoard{}
	new_board.Initialize()

	cachedReplays[user_id] = GameReplay{ b: &new_board, gameRecord: rec }
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
