package models

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"tic_tac_toe.fun/defs"
)

// GameRecord holds all data that is needed to recreate a game between two players.
type GameRecord struct {
	Id int
	U1 *User
	U2 *User
	DateRecorded time.Time
	History string // Represents the exact move sequence of the game
	Winner byte // The winner is represented with 'x', 'o', or '_' if the game resulted in a tie
}

// StoreGameRecord writes the GameRecord data into a database if at least one of the players played logged in.
func (Db *DataBase) StoreGameRecord(record GameRecord) error {
	// Don't store the game if neither user is logged in
	if record.U1.Id == defs.NO_USER_ID && record.U2.Id == defs.NO_USER_ID {
		return nil
	}

	statement := "insert into game_records (p1, p2, date_recorded, moves, winner) values (?, ?, ?, ?, ?)"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	winner_str := string([]byte{ record.Winner })
	_, err = stmt.Exec(record.U1.Id, record.U2.Id, record.DateRecorded, record.History, winner_str)
	return err
}

// ReadGameRecord returns a populated GameRecord with data based on the record id that is passed in,
// and also loads the user data of both players.
func (Db *DataBase) ReadGameRecord(rec_id int) (*GameRecord, error) {
	rec := &GameRecord{
		Id: rec_id,
		U1: new(User),
		U2: new(User),
	}

	var winner_str string

	statemet := "select p1, p2, date_recorded, moves, winner from game_records where id = ?"
	row := Db.Data.QueryRow(statemet, rec_id)
	err := row.Scan(
		&rec.U1.Id,
		&rec.U2.Id,
		&rec.DateRecorded,
		&rec.History,
		&winner_str,
	)
	if err != nil {
		return nil, err
	}

	rec.U1, err = Db.ReadUser(rec.U1.Id)
	if err != nil {
		return nil, err
	}

	rec.U2, err = Db.ReadUser(rec.U2.Id)
	if err != nil {
		return nil, err
	}

	rec.Winner = winner_str[0]

	return rec, nil
}

// ReadGameRecordsForUser returns a slice of all records in which one of the players is the player whose id is passed to the function.
// Note that if the user_id passed in is equal to the defs.NO_USER_ID, there will be no error, but the slice empty.
func (Db *DataBase) ReadGameRecordsForUser(user_id int) ([]*GameRecord, error) {

	res := make([]*GameRecord, 0)

	if user_id == defs.NO_USER_ID {
		return res, nil
	}

	statement := "select id, p1, p2, date_recorded, moves, winner from game_records where ? in (p1, p2) order by date_recorded desc"
	rows, err := Db.Data.Query(statement, user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rec := &GameRecord{
			U1: new(User),
			U2: new(User),
		}
		var winner_str string

		err := rows.Scan(
			&rec.Id,
			&rec.U1.Id,
			&rec.U2.Id,
			&rec.DateRecorded,
			&rec.History,
			&winner_str,
		)
		if err != nil {
			return nil, err
		}

		rec.U1, err = Db.ReadUser(rec.U1.Id)
		if err != nil {
			return nil, err
		}

		rec.U2, err = Db.ReadUser(rec.U2.Id)
		if err != nil {
			return nil, err
		}

		rec.Winner = winner_str[0]

		res = append(res, rec)
	}

	return res, nil

}

// GetGameRecordWinner is a convenience function that returns the id of a player who won the game based on the Winner field.
// Note that user 1 is always x and user 2 is always 0, if the game was tied, -1 is returned.
func GetGameRecordWinner(record *GameRecord) int {
	switch record.Winner {
	case 'x':
		return record.U1.Id
	case 'o':
		return record.U2.Id
	default:
		return -1
	}
}
