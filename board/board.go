package board

import (
    "errors"
)

// board layout
// 0 1 2
// 3 4 5
// 6 7 8

type Board [9]byte

func (b *Board) MakeMove(pos int, player byte) (err error) {
    if pos < 0 || pos > 8 || player != byte('x') || player != byte('o') {
        err = errors.New("invalid call to make_move:\nExpected pos 0-8, got", pos, "\nExpected player x or o, got", player)
        return
    }

    if b[pos] != 0 {
        return errors.New("invlaid call to make_move:\nCannov overwrite already played field")
    }

    b[pos] = player

    return
}

func (b *Board) CheckForWin() (res bool) {
    res = false

    for i := 0; i < 3; i++ {
        // check columns
        if b[i] == 0 {
            continue
        }
        if b[i] == b[i+3] && b[i] == b[i+6] {
            res = true
            return
        }
        if b[3*i] == b[3*i+1] && b[3*i] == b[3*i+2] {
            res = true
            return
        }
    }

    return
}
