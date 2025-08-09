package board

import (
	"errors"
	"fmt"
	// "log"

	"tic_tac_toe.fun/defs"
	"tic_tac_toe.fun/stack"
)

type BigBoard struct {
	boardToPlayIn  byte      // the next mini-board to make a move in
	Boards         [9]*Board // mini-boards
	BoardState     [][]byte  // this is used to update the front end ui
	Result         byte      // outcome-of the game, the value is either 0, 'x', 'o' or defs.BoardTie ('_')
	boardsComplete byte      // keeps track of how many mini-boards are complete to handle an over-all tie
	// History        string
	History *stack.Stack[Move]
}

type Move struct {
	BigPos   byte `json:"BigPos"`
	SmallPos byte `json:"SmallPos"`
	Player   byte `json:"Player"`
}

func (bb *BigBoard) Initialize() {
	for i := 0; i < 9; i++ {
		bb.Boards[i] = &Board{}
	}

	bb.BoardState = make([][]byte, 9)

	for i := 0; i < 9; i++ {
		bb.BoardState[i] = make([]byte, 9)
	}

	bb.boardToPlayIn = defs.WILD_BOARD // make sure the first play can be made in any of the mini-boards

	bb.History = stack.CreateStack[Move]()
}

// big pos represents the mini-board index, the pos represents the cell index
func (bb *BigBoard) MakeMove(m Move) error {

	if m.BigPos >= '0' && m.BigPos <= '9' { // the positions passed in are (probably) ascii characters representing digits
		m.BigPos = m.BigPos - '0'
	}

	if m.SmallPos >= '0' && m.SmallPos <= '9' {
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
		// if bb.Result != 0 {
		// 	bb.History += strings.ToUpper(string([]byte{ bb.Result }))
		// }
	}

	for i := 0; i < 9; i++ { // update the board state partially, since only one mini-board can change per move
		bb.BoardState[m.BigPos][i] = bb.Boards[m.BigPos].Cells[i]
	}

	return nil
}

func (bb *BigBoard) UndoLastMove() error {

	// if bb.History[len(bb.History)-1] != defs.BOARD_HISTORY_DELIMITER { // if the history is 'complete' and we want to remove the last move, chop off the winner character first
	// 	bb.History = bb.History[:len(bb.History)-1]
	// }
	m, ok := bb.History.Top()
	if !ok {
		return nil
	}

	bb.History.Pop() // not checking if pop returns ok since it will always do so if it got to this point

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

	bb.Boards[m.BigPos].UndoMove(*m)

	prev_m, ok := bb.History.Top()
	if !ok {
		bb.boardToPlayIn = defs.WILD_BOARD
	} else {
		bb.boardToPlayIn = prev_m.SmallPos
	}

	return nil
}

// func (bb *BigBoard) GenerateBoardResBitmap() (x_board []byte, o_board []byte) {
// 	x_board = make([]byte, 9)
// 	o_board = make([]byte, 9)
//
// 	for i := range 9 {
// 		if bb.Boards[i].Result == 'x' {
// 			x_board[i] = 'x'
// 		} else if bb.Boards[i].Result == 'o' {
// 			o_board[i] = 'o'
// 		}
// 	}
//
// 	return
// }
//
// func getXBitboards() [][]byte {
// 	return [][]byte{
// 		// Row wins
// 		{'x','x','x',
// 		  0 , 0 , 0 ,
// 		  0 , 0 , 0},
//
// 		{ 0 , 0 , 0 ,
// 		 'x','x','x',
// 		  0 , 0 , 0},
//
// 		{ 0 , 0 , 0 ,
// 		  0 , 0 , 0 ,
// 		 'x','x','x'},
//
// 		// Column wins
// 		{'x', 0 , 0 ,
// 		 'x', 0 , 0 ,
// 		 'x', 0 , 0},
//
// 		{ 0 ,'x', 0 ,
// 		  0 ,'x', 0 ,
// 		  0 ,'x', 0},
//
// 		{ 0 , 0 ,'x',
// 		  0 , 0 ,'x',
// 		  0 , 0 ,'x'},
//
// 		// Diagonal wins
// 		{'x', 0 , 0 ,
// 		  0 ,'x', 0 ,
// 		  0 , 0 ,'x'},
//
// 		{ 0 , 0 ,'x',
// 		  0 ,'x', 0 ,
// 		 'x', 0 , 0},
// 	}
// }
//
// func getOBitboards() [][]byte {
// 	return [][]byte {
// 		// Row wins
// 		{'o','o','o',
// 		  0 , 0 , 0 ,
// 		  0 , 0 , 0},
//
// 		{ 0 , 0 , 0 ,
// 		 'o','o','o',
// 		  0 , 0 , 0},
//
// 		{ 0 , 0 , 0 ,
// 		  0 , 0 , 0 ,
// 		 'o','o','o'},
//
// 		// Column wins
// 		{'o', 0 , 0 ,
// 		 'o', 0 , 0 ,
// 		 'o', 0 , 0},
//
// 		{ 0 ,'o', 0 ,
// 		  0 ,'o', 0 ,
// 		  0 ,'o', 0},
//
// 		{ 0 , 0 ,'o',
// 		  0 , 0 ,'o',
// 		  0 , 0 ,'o'},
//
// 		// Diagonal wins
// 		{'o', 0 , 0 ,
// 		  0 ,'o', 0 ,
// 		  0 , 0 ,'o'},
//
// 		{ 0 , 0 ,'o',
// 		  0 ,'o', 0 ,
// 		 'o', 0 , 0},
// 	}
// }

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

	return
}

func (bb *BigBoard) isMakeMoveValid(m Move) error { // checks if received input for a move is valid

	if m.BigPos != bb.boardToPlayIn && bb.boardToPlayIn != defs.WILD_BOARD {
		return fmt.Errorf("invalid call to make_move:\nm.BigPos and prev_pos missmatch: m.BigPos = %d, prev_pos = %d", m.BigPos, bb.boardToPlayIn)

	}

	if m.BigPos < 0 || m.BigPos > 8 {
		return fmt.Errorf("invalid call to make_move:\nExpected m.BigPos 0-8, got %d", m.BigPos)

	}

	if m.SmallPos < 0 || m.SmallPos > 8 {
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

	if m.BigPos != bb.boardToPlayIn && bb.boardToPlayIn != defs.WILD_BOARD {
		return fmt.Errorf("invalid call to UndoMove:\nm.BigPos and prev_pos missmatch: m.BigPos = %d, prev_pos = %d", m.BigPos, bb.boardToPlayIn)
	}

	if m.BigPos < 0 || m.BigPos > 8 {
		return fmt.Errorf("invalid call to UndoMove:\nExpected m.BigPos 0-8, got %d", m.BigPos)
	}

	if m.SmallPos < 0 || m.SmallPos > 8 {
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
