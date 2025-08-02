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
	Record string
}

func (Db *DataBase) StoreGameRecord(record GameRecord) error {
	// Don't store the game if neither user is logged in
	if record.U1.Id == defs.NO_USER_ID && record.U2.Id == defs.NO_USER_ID {
		return nil
	}

	statement := "insert into game_records (p1, p2, date_recorded, moves) values (?, ?, ?, ?)"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(record.U1.Id, record.U2.Id, record.DateRecorded, record.Record)
	return err
}

func (Db *DataBase) ReadGameRecord(rec_id int) (*GameRecord, error) {
	rec := &GameRecord{
		Id: rec_id,
		U1: new(User),
		U2: new(User),
	}

	statemet := "select p1, p2, date_recorded, moves from game_records where id = ?"
	row := Db.Data.QueryRow(statemet, rec_id)
	err := row.Scan(
		&rec.U1.Id,
		&rec.U2.Id,
		&rec.DateRecorded,
		&rec.Record,
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

	return rec, nil
}

func (Db *DataBase) ReadGameRecordsForUser(user_id int) ([]*GameRecord, error) {

	if user_id == defs.NO_USER_ID {
		return nil, nil
	}

	res := make([]*GameRecord, 0)

	statement := "select id, p1, p2, date_recorded, moves from game_records where ? in (p1, p2)"
	rows, err := Db.Data.Query(statement, user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rec := &GameRecord{
			U1: new(User),
			U2: new(User),
		}
		err := rows.Scan(
			&rec.Id,
			&rec.U1.Id,
			&rec.U2.Id,
			&rec.DateRecorded,
			&rec.Record,
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

		res = append(res, rec)
	}

	return res, nil

}

func GetGameRecordWinner(record *GameRecord) int {
	record_len := len(record.Record)
	if record_len < 2 {
		return defs.NO_WINNER
	}

	if record.Record[record_len - 2] == 'X' {
		return record.U1.Id
	}

	return record.U2.Id
}
