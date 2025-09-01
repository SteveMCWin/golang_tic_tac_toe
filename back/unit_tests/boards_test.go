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

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(0),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(6),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(3),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(8),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(7),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(4),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(5),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(1),
		Player: 'o',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(2),
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

    b = &board.Board{}

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(0),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(6),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(3),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(8),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(7),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(4),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(5),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(1),
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(2),
		Player: 'o',
	})
    if(b.Result != 'o') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

    b = &board.Board{}

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(0),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(3),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(7),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(5),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(1),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(4),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(6),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(8),
		Player: 'o',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(2),
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

    b = &board.Board{}

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(0),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(3),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(7),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(5),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(1),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(4),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(6),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(8),
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(2),
		Player: 'o',
	})
    if(b.Result != 'o') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }
}

func TestSmallBoardOWin(t *testing.T) {
    b := &board.Board{}

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(0),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(1),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(3),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(6),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(4),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(5),
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(8),
		Player: 'o',
	})

    if(b.Result != 'o') {
        t.Errorf("Winner is supposed to be o but it's %c", b.Result)
    }
}

func TestSmallBoardXWin(t *testing.T) {
    b := &board.Board{}

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(0),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(1),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(3),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(6),
		Player: 'o',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(4),
		Player: 'x',
	})
    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(5),
		Player: 'o',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }

    b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: byte(8),
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

}

func TestSmallBoardCheckCols(t *testing.T) {
	b := &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 0,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 3,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 6,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

	b = &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 1,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 4,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 7,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

	b = &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 2,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 5,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 8,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }
}

func TestSmallBoardCheckRows(t *testing.T) {
	b := &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 0,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 1,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 2,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

	b = &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 3,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 4,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 5,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

	b = &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 6,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 7,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 8,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }
}

func TestSmallBoardCheckDiags(t *testing.T) {
	b := &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 0,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 4,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 8,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }

	b = &board.Board{}

	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 2,
		Player: 'x',
	})
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 4,
		Player: 'x',
	})
    if(b.Result != 0) {
        t.Errorf("Winner is supposed to be undecided but it's %c", b.Result)
    }
	b.MakeMove(board.Move{
		BigPos: 0,
		SmallPos: 6,
		Player: 'x',
	})
    if(b.Result != 'x') {
        t.Errorf("Winner is supposed to be x but it's %c", b.Result)
    }
}

