package board

import "tic_tac_toe.fun/defs"

// board layout
// 0 1 2
// 3 4 5
// 6 7 8

type Board struct {
    Cells [9]byte
    Result byte
    movesPlayed byte    // used to determine the winner of the board when it's full
}

func (b *Board) MakeMove(m Move) {   // there are no safeguards in this method because the BigBoard does all the safeguards for it
    b.Cells[m.SmallPos] = m.Player
    b.movesPlayed += 1
    b.UpdateResult()    // check for winner after each move
}

func (b *Board) UndoMove(m Move) {
	b.Cells[m.SmallPos] = defs.EMPTY_CELL
	b.movesPlayed -= 1
	b.UpdateResult()
}

// checks for a winner the normal way first, if no winner but board full, make the one with most plays in the board win
// in normal tic-tac-toe this would always result in x winning, but not in 2d tic-tac-toe
func (b *Board) UpdateResult() {
    for i := 0; i < 3; i++ {
        // check columns
        if b.Cells[i] == defs.EMPTY_CELL {
            continue
        }
        if b.Cells[i] == b.Cells[i+3] && b.Cells[i] == b.Cells[i+6] {
            b.Result = b.Cells[i]
            return 
        }
        // check rows
        if b.Cells[3*i] == defs.EMPTY_CELL {
            continue
        }
        if b.Cells[3*i] == b.Cells[3*i+1] && b.Cells[3*i] == b.Cells[3*i+2] {
            b.Result = b.Cells[i]
            return
        }
    }

    // check diagonals
    if b.Cells[4] == defs.EMPTY_CELL {
        return 
    }

    if b.Cells[0] == b.Cells[4] && b.Cells[0] == b.Cells[8] {
        b.Result = b.Cells[4]
        return
    }

    if b.Cells[2] == b.Cells[4] && b.Cells[2] == b.Cells[6] {
        b.Result = b.Cells[4]
        return
    }

    if b.movesPlayed >= 9 { // if neither player won so far but the board is full, check who made most plays
        moves_made := make(map[byte]byte)

        for _, val := range b.Cells {
            moves_made[val] += 1
        }

        var winner byte
        var max_plays byte

        for key, val := range moves_made {
            if val > max_plays {
                max_plays = val
                winner = key
            }
        }

        b.Result = winner

    }

}
