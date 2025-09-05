package board

import (
	"errors"
	"fmt"

	"tic_tac_toe.fun/defs"
	"tic_tac_toe.fun/stack"
)

// BigBoard is a collection of 9 normal (here usually referred to as mini) boards.
type BigBoard struct {
	boardToPlayIn  byte               // the next mini-board to make a move in
	Boards         [9]*Board          // Mini-Boards
	BoardState     [][]byte           // A simpler representation of the big board with all plays made so far. Used to update the front end ui.
	Result         byte               // Outcome of the game. The value at the end of the game is either 'x', 'o' or defs.BoardTie ('_')
	boardsComplete byte               // keeps track of how many mini-boards are complete to handle an over-all tie
	History        *stack.Stack[Move] // All of the plays made in order.
	RedoStack *stack.Stack[Move] // Used in replays. When you undo a move it is moved from the history stack to this stack. When you redo a move, it's the opposite
}

// Move represents a move on a big board by the player 'x' or the player 'o'.
type Move struct {
	BigPos   byte `json:"BigPos"`   // Index of a mini board inside of BigBoard
	SmallPos byte `json:"SmallPos"` // Index of a cell inside a mini board
	Player   byte `json:"Player"`   // 'x' or 'o'
}

// Initialize sets up a new board.
func (bb *BigBoard) Initialize() {
	for i := range 9 {
		bb.Boards[i] = &Board{}
	}

	bb.BoardState = make([][]byte, 9)

	for i := range 9 {
		bb.BoardState[i] = make([]byte, 9)
	}

	bb.boardToPlayIn = defs.WILD_BOARD // make sure the first play can be made in any of the mini-boards

	bb.History = stack.CreateStack[Move]()
	bb.RedoStack = stack.CreateStack[Move]()
}

func (bb *BigBoard) GetNextMiniBoard() byte {
	return bb.boardToPlayIn
}

// MakeMove updates the board corresponding to the Move passed in, if valid, and determines the next mini board to be played in.
// Note that on every MakeMove call, the history of the board gets updated.
// If the move made results in a mini board being won by a player, MakeMove checks if either player won the big board as well.
func (bb *BigBoard) MakeMove(m Move) error {

	if m.BigPos >= '0' && m.BigPos < '9' { // the positions passed in are (probably) ascii characters representing digits
		m.BigPos = m.BigPos - '0'
	}

	if m.SmallPos >= '0' && m.SmallPos < '9' {
		m.SmallPos = m.SmallPos - '0'
	}

	err := bb.isMakeMoveValid(m)
	if err != nil {
		return err
	}

	bb.Boards[m.BigPos].MakeMove(m) // updates the board through its mini-boards

	bb.History.Push(m)

	if bb.Boards[m.SmallPos].Result == 0 {
		bb.boardToPlayIn = m.SmallPos
	} else { // if the mini-board that should be played in next is already completed, allow the next player to play in any mini-board
		bb.boardToPlayIn = defs.WILD_BOARD
	}

	if bb.Boards[m.BigPos].Result != byte(0) { // if the board that was just played in gets completed, check for win and update counter
		bb.boardsComplete += 1
		bb.UpdateResult()
	}

	bb.updateBoardState(int(m.BigPos))
	// for i := range 9 { // update the board state partially, since only one mini-board can change per move
	// 	bb.BoardState[m.BigPos][i] = bb.Boards[m.BigPos].Cells[i]
	// }

	return nil
}

func (bb *BigBoard) updateBoardState(big_pos int) {
	for i := range 9 {
		bb.BoardState[big_pos][i] = bb.Boards[big_pos].Cells[i]
	}
}

// UndoLastMove removes the last move made on the board and from the history.
func (bb *BigBoard) UndoLastMove() error {

	m := bb.History.Top()
	if m == nil {
		return nil
	}

	if m.BigPos >= '0' && m.BigPos <= '9' { // the positions passed in are (probably) ascii characters representing digits
		m.BigPos = m.BigPos - '0'
	}

	if m.SmallPos >= '0' && m.SmallPos <= '9' {
		m.SmallPos = m.SmallPos - '0'
	}

	err := bb.isUndoMoveValid(*m)
	if err != nil {
		return err
	}

	bb.RedoStack.Push(*m)
	bb.History.Pop() // not checking if pop returns ok since it will always do so if it got to this point

	bb.Boards[m.BigPos].UndoMove(*m)

	bb.updateBoardState(int(m.BigPos))

	prev_m:= bb.History.Top()
	if prev_m == nil {
		bb.boardToPlayIn = defs.WILD_BOARD
	} else {
		bb.boardToPlayIn = prev_m.SmallPos
	}

	return nil
}

func (bb *BigBoard) RedoLastMove() error {
	m := bb.RedoStack.Top()
	if m == nil {
		return nil
	}

	err := bb.MakeMove(*m)
	if err != nil {
		return err
	}

	bb.RedoStack.Pop()

	return nil
}

// UpdateResult checks if eigher player won the big board.
func (bb *BigBoard) UpdateResult() { // checks if anyone won and if so, update the board's result

	// check columns
	for i := range 3 {
		if bb.Boards[i].Result == 0 {
			continue
		}
		if bb.Boards[i].Result == bb.Boards[i+3].Result && bb.Boards[i].Result == bb.Boards[i+6].Result {
			bb.Result = bb.Boards[i].Result
			return
		}
	}

	// check rows
	for i := range 3 {
		if bb.Boards[3*i].Result == 0 {
			continue
		}
		if bb.Boards[3*i].Result == bb.Boards[3*i+1].Result && bb.Boards[3*i].Result == bb.Boards[3*i+2].Result {
			bb.Result = bb.Boards[3*i].Result
			return
		}
	}

	// check diagonals
	if bb.Boards[4].Result == 0 {
		return
	}

	if bb.Boards[0].Result == bb.Boards[4].Result && bb.Boards[0].Result == bb.Boards[8].Result {
		bb.Result = bb.Boards[4].Result
		return
	}

	if bb.Boards[2].Result == bb.Boards[4].Result && bb.Boards[2].Result == bb.Boards[6].Result {
		bb.Result = bb.Boards[4].Result
		return
	}

	// if board full but neither player won, it's a tie
	if bb.boardsComplete >= 9 {
		bb.Result = defs.BOARD_TIE
		return
	}
}

func (bb *BigBoard) isMakeMoveValid(m Move) error { // checks if received input for a move is valid

	if m.BigPos != bb.boardToPlayIn && bb.boardToPlayIn != defs.WILD_BOARD {
		return fmt.Errorf("invalid call to make_move:\nm.BigPos and prev_pos missmatch: m.BigPos = %d, prev_pos = %d", m.BigPos, bb.boardToPlayIn)

	}

	if m.BigPos > 8 {
		return fmt.Errorf("invalid call to make_move:\nExpected m.BigPos 0-8, got %d", m.BigPos)

	}

	if m.SmallPos > 8 {
		return fmt.Errorf("invalid call to make_move:\nExpected m.SmallPos 0-8, got %d", m.SmallPos)

	}

	if m.Player != byte('x') && m.Player != byte('o') {
		return fmt.Errorf("invalid call to make_move:\nExpected player x or o, got %c", m.Player)

	}

	if bb.Boards[m.BigPos].Result != byte(0) {
		return fmt.Errorf("invalid call to make_move:\nBoard %d already complete", m.BigPos)

	}

	if bb.Boards[m.BigPos].Cells[m.SmallPos] != byte(0) {
		return errors.New("invlaid call to make_move:\nCannot overwrite already played field")

	}

	return nil
}

func (bb *BigBoard) isUndoMoveValid(m Move) error { // checks if received input for a move is valid

	// if m.BigPos != bb.boardToPlayIn && bb.boardToPlayIn != defs.WILD_BOARD {
	// 	return fmt.Errorf("invalid call to UndoMove:\nm.BigPos and prev_pos missmatch: m.BigPos = %d, prev_pos = %d", m.BigPos, bb.boardToPlayIn)
	// }

	if m.BigPos > 8 {
		return fmt.Errorf("invalid call to UndoMove:\nExpected m.BigPos 0-8, got %d", m.BigPos)
	}

	if m.SmallPos > 8 {
		return fmt.Errorf("invalid call to UndoMove:\nExpected m.SmallPos 0-8, got %d", m.SmallPos)
	}

	if m.Player != byte('x') && m.Player != byte('o') {
		return fmt.Errorf("invalid call to UndoMove:\nExpected m.Player x or o, got %c", m.Player)
	}

	if bb.Boards[m.BigPos].Cells[m.SmallPos] != m.Player {
		return errors.New("invlaid call to UndoMove:\nMissmatch between cell byte and passed in m.Player byte")
	}

	return nil
}
