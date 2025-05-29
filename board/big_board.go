package board

import (
    "log"
)

type BigBoard struct {
    boardToPlayIn byte
    Boards [9]*Board
    Result byte
}

// TODO: handle when the boardToPlayIn is complete

func (bb *BigBoard) MakeMove(big_pos byte, pos byte, player byte) (b_state [][]byte, err error) {
    err = bb.checkIfMoveValid(big_pos, pos, player)
    if err != nil {
        return
    }

    bb.Boards[big_pos].MakeMove(pos, player)
    bb.boardToPlayIn = pos

    res := bb.Boards[big_pos].UpdateResult()
    if res == true {
        bb.UpdateResult()
    }

    b_state = make([][]byte, 9)

    for i := 0; i < 9; i++ {
        b_state[i] = make([]byte, 9)
        for j := 0; j < 9; j++ {
            b_state[i][j] = bb.Boards[i].Cells[j]
        }
    }

    return
}

// TODO: check for tie and stop the game accordingly
func (bb *BigBoard) UpdateResult() {
    for i := 0; i < 3; i++ {
        // check columns
        if bb.Boards[i].Result == 0 {
            continue
        }
        if bb.Boards[i].Result == bb.Boards[i+3].Result && bb.Boards[i].Result == bb.Boards[i+6].Result {
            b.Result = bb.Boards[i].Result
            return
        }
        if bb.Boards[3*i].Result == 0 {
            continue
        }
        if bb.Boards[3*i].Result == bb.Boards[3*i+1].Result && bb.Boards[3*i].Result == bb.Boards[3*i+2].Result {
            b.Result = bb.Boards[i].Result
            return
        }
    }

    if bb.Boards[4].Result == 0 {
        return
    }

    if bb.Boards[0].Result == bb.Boards[4].Result && bb.Boards[0].Result == bb.Boards[8].Result {
        b.Result = bb.Boards[i].Result
        return
    }

    if bb.Boards[2].Result == bb.Boards[4].Result && bb.Boards[2].Result == bb.Boards[6].Result {
        b.Result = bb.Boards[i].Result
        return
    }

    return
}

func (bb *BigBoard) checkIfMoveValid(big_pos byte, pos byte, player byte) (err error) {

    if big_pos >= '0' && big_pos <= '9' {
        big_pos = big_pos - '0'
    }

    if pos >= '0' && pos <= '9' {
        pos = pos - '0'
    }

    if big_pos != bb.boardToPlayIn {
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

