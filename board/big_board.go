package board

import (
    "log"
    "fmt"
    "errors"
)

const wildBoard = '?'
const BoardTie = '_'

type BigBoard struct {
    boardToPlayIn byte
    Boards [9]*Board
    BoardState [][]byte
    Result byte
    boardsComplete byte
}

func (bb *BigBoard) Initialize() {
    for i := 0; i < 9; i++ {
        bb.Boards[i] = &Board{}
    }

    bb.BoardState = make([][]byte, 9)

    for i := 0; i < 9; i++ {
        bb.BoardState[i] = make([]byte, 9)
    }

    bb.boardToPlayIn = wildBoard
}

func (bb *BigBoard) MakeMove(big_pos byte, pos byte, player byte) (err error) {

    if big_pos >= '0' && big_pos <= '9' {
        big_pos = big_pos - '0'
    }

    if pos >= '0' && pos <= '9' {
        pos = pos - '0'
    }

    err = bb.checkIfMoveValid(big_pos, pos, player)
    if err != nil {
        return
    }

    log.Println("big_pos:", big_pos)
    log.Println("pos:", pos)
    bb.Boards[big_pos].MakeMove(pos, player)

    if bb.Boards[pos].Result == 0 {
        bb.boardToPlayIn = pos
    } else {
        bb.boardToPlayIn = wildBoard
    }

    if bb.Boards[big_pos].Result != byte(0) {
        bb.boardsComplete += 1
        bb.UpdateResult()
    }

    for i := 0; i < 9; i++ {
        bb.BoardState[big_pos][i] = bb.Boards[big_pos].Cells[i]
    }

    return
}

func (bb *BigBoard) UpdateResult() {
    for i := 0; i < 3; i++ {
        // check columns
        if bb.Boards[i].Result == 0 {
            continue
        }
        if bb.Boards[i].Result == bb.Boards[i+3].Result && bb.Boards[i].Result == bb.Boards[i+6].Result {
            bb.Result = bb.Boards[i].Result
            return
        }
        if bb.Boards[3*i].Result == 0 {
            continue
        }
        if bb.Boards[3*i].Result == bb.Boards[3*i+1].Result && bb.Boards[3*i].Result == bb.Boards[3*i+2].Result {
            bb.Result = bb.Boards[i].Result
            return
        }
    }

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

    if bb.boardsComplete >= 9 {
        bb.Result = BoardTie
        return
    }

    return
}

func (bb *BigBoard) checkIfMoveValid(big_pos byte, pos byte, player byte) (err error) {

    if big_pos != bb.boardToPlayIn && bb.boardToPlayIn != wildBoard {
        err = fmt.Errorf("invalid call to make_move:\nbig_pos and prev_pos missmatch: big_pos = %d, prev_pos = %d", big_pos, bb.boardToPlayIn)
        return
    }

    if big_pos < 0 || big_pos > 8 {
        err = fmt.Errorf("invalid call to make_move:\nExpected big_pos 0-8, got %d", big_pos)
        return
    }

    if pos < 0 || pos > 8 {
        err = fmt.Errorf("invalid call to make_move:\nExpected pos 0-8, got %d", pos)
        return
    }

    if player != byte('x') && player != byte('o') {
        err = fmt.Errorf("invalid call to make_move:\nExpected player x or o, got %c", player)
        return
    }

    if bb.Boards[big_pos].Result != byte(0) {
        err = fmt.Errorf("invalid call to make_move:\nBoard %d already complete", big_pos)
        return
    }

    if bb.Boards[big_pos].Cells[pos] != byte(0) {
        err = errors.New("invlaid call to make_move:\nCannov overwrite already played field")
        return
    }

    return
}

