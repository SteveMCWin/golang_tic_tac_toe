package unit_tests

import (
    "testing"
    "tic_tac_toe.fun/board"
)

func TestMoveValidity(t *testing.T) {

    bb := &board.BigBoard{}
    bb.Initialize()

    _, err := bb.MakeMove(byte(4), byte(3), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    _, err = bb.MakeMove(byte(4), byte(2), 'o')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    _, err = bb.MakeMove(byte(3), byte(1), 'o')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    _, err = bb.MakeMove(byte(9), byte(2), 'x')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    _, err = bb.MakeMove(byte(1), byte(9), 'x')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    _, err = bb.MakeMove(byte(1), byte(1), 'r')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    _, err = bb.MakeMove(byte(1), byte(4), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    _, err = bb.MakeMove(byte(4), byte(3), 'o')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    bb = &board.BigBoard{}
    bb.Initialize()
    bb.Boards[2].Cells[0] = 'x'
    bb.Boards[2].Cells[1] = 'x'

    _, err = bb.MakeMove(byte(8), byte(2), 'o')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    _, err = bb.MakeMove(byte(2), byte(2), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    _, err = bb.MakeMove(byte(2), byte(7), 'o')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    _, err = bb.MakeMove(byte(4), byte(2), 'o')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

    _, err = bb.MakeMove(byte(2), byte(5), 'x')
    if err == nil {
        t.Error("Invalid move detected as valid")
    }

    _, err = bb.MakeMove(byte(4), byte(3), 'x')
    if err != nil {
        t.Error("Valid move detected as invalid")
    }

}
