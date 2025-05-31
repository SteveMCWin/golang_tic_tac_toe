package unit_tests

import (
    // "fmt"
    "testing"
    "tic_tac_toe.fun/board"
)

// TODO: figure out how to force a tie, the current problem is that the boardsCompelte field needs to be 9 somehow
// func TestTie(t *testing.T) {
//     bb := &board.BigBoard{}
//     bb.Initialize()
//
//     var err error
//
//     bb.Boards[0].Result = 'x'
//     bb.Boards[1].Result = 'x'
//     bb.Boards[2].Result = 'o'
//     bb.Boards[3].Result = 'o'
//     bb.Boards[4].Result = 'o'
//     bb.Boards[5].Result = 'x'
//     bb.Boards[6].Result = 'x'
//     bb.Boards[7].Result = 'x'
//     bb.Boards[8].Cells = [9]byte{'o', 'o', 0, 0, 0, 0, 0, 0, 0}
//
//     err = bb.MakeMove(byte(8), byte(2), 'o')
//
//     for i := 0; i < 9; i++ {
//         if bb.Boards[i].Result != 0 {
//             fmt.Printf("Boards[%d].Result = %c\n", i, bb.Boards[i].Result)
//         }
//     }
//
//     if err != nil {
//         t.Error("error making move")
//     }
//
//     if bb.Result != board.BoardTie {
//         t.Errorf("should have been a tie but the winner is: %c", bb.Result)
//     }
//
// }

func TestGameXWins(t *testing.T) {

    bb := &board.BigBoard{}
    bb.Initialize()

    var err error

    err = bb.MakeMove(byte(5), byte(0), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(0), byte(5), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(2), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(7), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(5), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(8), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(8), byte(2), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(8), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(8), byte(5), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(5), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(1), 'x')
    if err != nil {
        t.Error("error making move")
    }//check board 5
    err = bb.MakeMove(byte(1), byte(2), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(6), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(6), byte(1), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(1), byte(1), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(1), byte(8), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(8), byte(8), 'x')
    if err != nil {
        t.Error("error making move")
    }//check board 8
    err = bb.MakeMove(byte(1), byte(5), 'o')
    if err != nil {
        t.Error("error making move")
    }//check board 1
    err = bb.MakeMove(byte(2), byte(4), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(4), byte(7), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(2), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(2), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(5), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(3), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(3), byte(7), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(7), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(8), 'x')
    if err != nil {
        t.Error("error making move")
    }//check board 7
    err = bb.MakeMove(byte(4), byte(2), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(0), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(0), byte(2), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(1), 'x')
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

    err = bb.MakeMove(byte(5), byte(0), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(0), byte(5), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(2), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(7), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(5), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(8), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(8), byte(2), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(8), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(8), byte(5), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(5), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(5), byte(1), 'o')
    if err != nil {
        t.Error("error making move")
    }//check board 5
    err = bb.MakeMove(byte(1), byte(2), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(6), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(6), byte(1), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(1), byte(1), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(1), byte(8), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(8), byte(8), 'o')
    if err != nil {
        t.Error("error making move")
    }//check board 8
    err = bb.MakeMove(byte(1), byte(5), 'x')
    if err != nil {
        t.Error("error making move")
    }//check board 1
    err = bb.MakeMove(byte(2), byte(4), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(4), byte(7), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(2), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(2), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(5), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(3), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(3), byte(7), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(7), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(7), byte(8), 'o')
    if err != nil {
        t.Error("error making move")
    }//check board 7
    err = bb.MakeMove(byte(4), byte(2), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(0), 'o')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(0), byte(2), 'x')
    if err != nil {
        t.Error("error making move")
    }
    err = bb.MakeMove(byte(2), byte(1), 'o')
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

    err := bb.MakeMove(byte(4), byte(3), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(byte(4), byte(2), 'o')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(byte(3), byte(1), 'o')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(byte(9), byte(2), 'x')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(byte(1), byte(9), 'x')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(byte(1), byte(1), 'r')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(byte(1), byte(4), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(byte(4), byte(3), 'o')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    bb = &board.BigBoard{}
    bb.Initialize()
    bb.Boards[2].Cells[0] = 'x'
    bb.Boards[2].Cells[1] = 'x'

    err = bb.MakeMove(byte(8), byte(2), 'o')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(byte(2), byte(2), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(byte(2), byte(7), 'o')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(byte(4), byte(2), 'o')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    err = bb.MakeMove(byte(2), byte(5), 'x')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    err = bb.MakeMove(byte(4), byte(3), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

}
