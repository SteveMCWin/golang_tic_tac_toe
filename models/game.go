package models

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
	"tic_tac_toe.fun/defs"
)

type GameRecord struct {
	Id int
	U1 *User
	U2 *User
	DateRecorded time.Time
	History string
	Winner byte
}

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

func (Db *DataBase) ReadGameRecordsForUser(user_id int) ([]*GameRecord, error) {

	if user_id == defs.NO_USER_ID {
		return nil, nil
	}

	res := make([]*GameRecord, 0)

	statement := "select id, p1, p2, date_recorded, moves, winner from game_records where ? in (p1, p2)"
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

func GetGameRecordWinner(record *GameRecord) int {
	if record.Winner == 'x' {
		return record.U1.Id
	}
	return record.U2.Id
}
