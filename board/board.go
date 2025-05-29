package board

import (
    "fmt"
    "errors"
)

// board layout
// 0 1 2
// 3 4 5
// 6 7 8

type Board struc {
    Cells [9]byte
    Result byte
}

func (b *Board) MakeMove(pos byte, player byte) (b_state []byte, err error) {
    b.Cells[pos] = player

    b_state = make([]byte, 9)

    for i := 0; i < 9; i++ {
        b_state[i] = b[i]
    }

    return
}

func (b *Board) UpdateResult() bool {
    for i := 0; i < 3; i++ {
        // check columns
        if b.Cells[i] == 0 {
            continue
        }
        if b.Cells[i] == b.Cells[i+3] && b.Cells[i] == b.Cells[i+6] {
            b.Result = b.Cells[i]
            return true
        }
        if b.Cells[3*i] == 0 {
            continue
        }
        if b.Cells[3*i] == b.Cells[3*i+1] && b.Cells[3*i] == b.Cells[3*i+2] {
            b.Result = b.Cells[i]
            return true
        }
    }

    if b.Cells[4] == 0 {
        return false
    }

    if b.Cells[0] == b.Cells[4] && b.Cells[0] == b.Cells[8] {
        b.Result = b.Cells[i]
        return true
    }

    if b.Cells[2] == b.Cells[4] && b.Cells[2] == b.Cells[6] {
        b.Result = b.Cells[i]
        return true
    }

    // TODO: handle all cells full but no winner

    return false

}
