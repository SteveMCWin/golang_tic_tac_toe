package board

import (
    "fmt"
    "errors"
)

// board layout
// 0 1 2
// 3 4 5
// 6 7 8

type Board [9]byte

func (b *Board) MakeMove(pos byte, player byte) (b_state []byte, err error) {
    if pos < 0 || pos > 8 || player != byte('x') || player != byte('o') {
        err = fmt.Errorf("invalid call to make_move:\nExpected pos 0-8, got %d\nExpected player x or o, got %c", pos, player)
        return
    }

    if b[pos] != 0 {
        err = errors.New("invlaid call to make_move:\nCannov overwrite already played field")
        return
    }

    b[pos] = player

    b_state = make([]byte, 9)

    for i := 0; i < 9; i++ {
        b_state[i] = b[i]
    }

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
