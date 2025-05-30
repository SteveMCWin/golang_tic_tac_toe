package board

// board layout
// 0 1 2
// 3 4 5
// 6 7 8

type Board struct {
    Cells [9]byte
    Result byte
    movesPlayed byte
}

// const boardTie = '_'

func (b *Board) MakeMove(pos byte, player byte) {
    b.Cells[pos] = player
    b.movesPlayed += 1
}

func (b *Board) UpdateResult() {
    for i := 0; i < 3; i++ {
        // check columns
        if b.Cells[i] == 0 {
            continue
        }
        if b.Cells[i] == b.Cells[i+3] && b.Cells[i] == b.Cells[i+6] {
            b.Result = b.Cells[i]
            return 
        }
        if b.Cells[3*i] == 0 {
            continue
        }
        if b.Cells[3*i] == b.Cells[3*i+1] && b.Cells[3*i] == b.Cells[3*i+2] {
            b.Result = b.Cells[i]
            return
        }
    }

    if b.Cells[4] == 0 {
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

    if b.movesPlayed >= 9 {
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

        // b.Result = boardTie
    }

}
