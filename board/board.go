// Package board contains data and functions realted to the boards the players play on.
package board

import "tic_tac_toe.fun/defs"

// board layout
// 0 1 2
// 3 4 5
// 6 7 8

// Board stuct represents a normal tic tac toe board.
// A Board cannot result in a tie.
type Board struct {
    Cells [9]byte // Each cell conatins a value that represents the player 'x' or 'o', or the value defs.EMPTY_CELL
    Result byte // Represents the player who won the board, represented by 'x' or 'o', or the value defs.EMPTY_CELL if no player won yet.
    movesPlayed byte    // used to determine the winner of the board when it's full
}

// MakeMove places an 'x' or an 'o' inside the specified index inside the board.
// Note that after each move, there is a check to see if either of the playes won.
func (b *Board) MakeMove(m Move) { // there are no safeguards in this method because the BigBoard does all the safeguards for it
    b.Cells[m.SmallPos] = m.Player
    b.movesPlayed += 1
    b.UpdateResult()
}

// UndoMove sets a specified index with an empty cell and checks if that changed the result.
func (b *Board) UndoMove(m Move) {
	b.Cells[m.SmallPos] = defs.EMPTY_CELL
	b.movesPlayed -= 1
	b.UpdateResult()
}

// UpdateResult checks if either player won. Note that if the board is a tie, the player with most plays in it wins.
// In normal tic-tac-toe this would always result in x winning, but in 2D Tic Tac Toe, that is not the case.
func (b *Board) UpdateResult() {
    for i := range 3 {
        // check columns
        if b.Cells[i] == defs.EMPTY_CELL {
            continue
        }
        if b.Cells[i] == b.Cells[i+3] && b.Cells[i] == b.Cells[i+6] {
            b.Result = b.Cells[i]
            return 
        }
	}

	for i := range 3 {
        // check rows
        if b.Cells[3*i] == defs.EMPTY_CELL {
            continue
        }
        if b.Cells[3*i] == b.Cells[3*i+1] && b.Cells[3*i] == b.Cells[3*i+2] {
            b.Result = b.Cells[3*i]
            return
        }
    }

    // check diagonals
    if b.Cells[4] == defs.EMPTY_CELL {
		b.Result = defs.EMPTY_CELL
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
