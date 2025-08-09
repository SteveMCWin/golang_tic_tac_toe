package unit_tests

import (
	"testing"

	"tic_tac_toe.fun/board"
)

func TestGameXWins(t *testing.T) {

    bb := &board.BigBoard{
	}
    bb.Initialize()

    var err error

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(0),
		Player: 'x',
	})
    if err != nil {
        t.Error(err)
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(0),
		SmallPos: byte(5),
		Player: 'o',
	})
    if err != nil {
        t.Error(err)
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error(err)
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(7),
		Player: 'o',
	})
    if err != nil {
        t.Error(err)
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(5),
		Player: 'x',
	})
    if err != nil {
        t.Error(err)
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(8),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(8),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(8),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(8),
		SmallPos: byte(5),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(5),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(1),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 5
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(6),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(6),
		SmallPos: byte(1),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(1),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(8),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(8),
		SmallPos: byte(8),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 8
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(5),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 1
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(4),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(7),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(5),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(3),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(3),
		SmallPos: byte(7),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(7),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(8),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 7
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(0),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(0),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
	if bb.Result != 0 {
		t.Errorf("No player should have won yet")
	}
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(1),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }

    if bb.Boards[1].Result != 'o' {
        t.Errorf("Result should be 'o' but is %c", bb.Boards[1].Result)
    }
    if bb.Boards[2].Result != 'x' {
        t.Errorf("Result should be 'x' but is %c", bb.Boards[2].Result)
    }
    if bb.Boards[5].Result != 'x' {
        t.Errorf("Result should be 'x' but is %c", bb.Boards[5].Result)
    }
    if bb.Boards[7].Result != 'x' {
        t.Errorf("Result should be 'x' but is %c", bb.Boards[7].Result)
    }
    if bb.Boards[8].Result != 'x' {
        t.Errorf("Result should be 'x' but is %c", bb.Boards[8].Result)
    }

    if bb.Result != 'x' {
        t.Errorf("Winner should be 'x' but is %c", bb.Result)
    }

}

func TestGameOWins(t *testing.T) {

    bb := &board.BigBoard{}
    bb.Initialize()

    var err error

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(0),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(0),
		SmallPos: byte(5),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(7),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(5),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(8),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(8),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(8),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(8),
		SmallPos: byte(5),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(5),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(5),
		SmallPos: byte(1),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 5
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(6),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(6),
		SmallPos: byte(1),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(1),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(8),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(8),
		SmallPos: byte(8),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 8
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(5),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 1
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(4),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(7),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(5),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(3),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(3),
		SmallPos: byte(7),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(7),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(7),
		SmallPos: byte(8),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }//check board 7
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(0),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(0),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(1),
		Player: 'o',
	})
    if err != nil {
        t.Error("error making move")
    }

    if bb.Boards[1].Result != 'x' {
        t.Errorf("Result should be 'x' but is %c", bb.Boards[1].Result)
    }
    if bb.Boards[2].Result != 'o' {
        t.Errorf("Result should be 'o' but is %c", bb.Boards[2].Result)
    }
    if bb.Boards[5].Result != 'o' {
        t.Errorf("Result should be 'o' but is %c", bb.Boards[5].Result)
    }
    if bb.Boards[7].Result != 'o' {
        t.Errorf("Result should be 'o' but is %c", bb.Boards[7].Result)
    }
    if bb.Boards[8].Result != 'o' {
        t.Errorf("Result should be 'o' but is %c", bb.Boards[8].Result)
    }

    if bb.Result != 'o' {
        t.Errorf("Winner should be 'o' but is %c", bb.Result)
    }

}

func TestMoveValidity(t *testing.T) {

    bb := &board.BigBoard{}
    bb.Initialize()

    err := bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(3),
		Player: 'x',
	})
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(3),
		SmallPos: byte(1),
		Player: 'o',
	})
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(9),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(9),
		Player: 'x',
	})
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(1),
		Player: 'r',
	})
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(1),
		SmallPos: byte(4),
		Player: 'x',
	})
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(3),
		Player: 'o',
	})
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    bb = &board.BigBoard{}
    bb.Initialize()
    bb.Boards[2].Cells[0] = 'x'
    bb.Boards[2].Cells[1] = 'x'

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(8),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(2),
		Player: 'x',
	})
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(7),
		Player: 'o',
	})
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(2),
		Player: 'o',
	})
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(2),
		SmallPos: byte(5),
		Player: 'x',
	})
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(board.Move{ 
		BigPos: byte(4),
		SmallPos: byte(3),
		Player: 'x',
	})
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

}
