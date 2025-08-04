package unit_tests

import (
    "testing"
    "tic_tac_toe.fun/board"
)

// board layout
// 0 1 2
// 3 4 5
// 6 7 8

func TestMostPlaysWin(t *testing.T) {
    b := &board.Board{}

    b.MakeMove(board.Move{0, byte(0), 'x'})
    b.MakeMove(board.Move{0, byte(6), 'o'})
    b.MakeMove(board.Move{0, byte(3), 'x'})
    b.MakeMove(board.Move{0, byte(8), 'o'})
    b.MakeMove(board.Move{0, byte(7), 'x'})
    b.MakeMove(board.Move{0, byte(4), 'o'})
    b.MakeMove(board.Move{0, byte(5), 'x'})
    b.MakeMove(board.Move{0, byte(1), 'o'})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{0, byte(2), 'x'})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

    b = &board.Board{}

    b.MakeMove(board.Move{0, byte(0), 'o'})
    b.MakeMove(board.Move{0, byte(6), 'x'})
    b.MakeMove(board.Move{0, byte(3), 'o'})
    b.MakeMove(board.Move{0, byte(8), 'x'})
    b.MakeMove(board.Move{0, byte(7), 'o'})
    b.MakeMove(board.Move{0, byte(4), 'x'})
    b.MakeMove(board.Move{0, byte(5), 'o'})
    b.MakeMove(board.Move{0, byte(1), 'x'})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{0, byte(2), 'o'})
    if(b.Result != 'o') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

    b = &board.Board{}

    b.MakeMove(board.Move{0, byte(0), 'x'})
    b.MakeMove(board.Move{0, byte(3), 'x'})
    b.MakeMove(board.Move{0, byte(7), 'x'})
    b.MakeMove(board.Move{0, byte(5), 'x'})
    b.MakeMove(board.Move{0, byte(1), 'o'})
    b.MakeMove(board.Move{0, byte(4), 'o'})
    b.MakeMove(board.Move{0, byte(6), 'o'})
    b.MakeMove(board.Move{0, byte(8), 'o'})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{0, byte(2), 'x'})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

    b = &board.Board{}

    b.MakeMove(board.Move{0, byte(0), 'o'})
    b.MakeMove(board.Move{0, byte(3), 'o'})
    b.MakeMove(board.Move{0, byte(7), 'o'})
    b.MakeMove(board.Move{0, byte(5), 'o'})
    b.MakeMove(board.Move{0, byte(1), 'x'})
    b.MakeMove(board.Move{0, byte(4), 'x'})
    b.MakeMove(board.Move{0, byte(6), 'x'})
    b.MakeMove(board.Move{0, byte(8), 'x'})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{0, byte(2), 'o'})
    if(b.Result != 'o') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }
}

func TestSmallBoardOWin(t *testing.T) {
    b := &board.Board{}

    b.MakeMove(board.Move{0, byte(0), 'o'})
    b.MakeMove(board.Move{0, byte(1), 'x'})
    b.MakeMove(board.Move{0, byte(3), 'o'})
    b.MakeMove(board.Move{0, byte(6), 'x'})
    b.MakeMove(board.Move{0, byte(4), 'o'})
    b.MakeMove(board.Move{0, byte(5), 'x'})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{0, byte(8), 'o'})

    if(b.Result != 'o') {
        t.Errorf("Winner is supposed to be o but it's %c", b.Result)
    }
}

func TestSmallBoardXWin(t *testing.T) {
    b := &board.Board{}

    b.MakeMove(board.Move{0, byte(0), 'x'})
    b.MakeMove(board.Move{0, byte(1), 'o'})
    b.MakeMove(board.Move{0, byte(3), 'x'})
    b.MakeMove(board.Move{0, byte(6), 'o'})
    b.MakeMove(board.Move{0, byte(4), 'x'})
    b.MakeMove(board.Move{0, byte(5), 'o'})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{0, byte(8), 'x'})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }
}


