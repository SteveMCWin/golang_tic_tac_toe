package models

import (
	// "database/sql"
	// "errors"
	// "io"
	// "log"
	// "net/http"
	// "os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"tic_tac_toe.fun/defs"
)

func (Db *DataBase) StoreGameRecord(game_record string, p1_id, p2_id int) error {
	// Don't store the game if neither user is logged in
	if p1_id == defs.NO_USER_ID && p2_id == defs.NO_USER_ID {
		return nil
	}

	statement := "insert into game_records (p1, p2, date_recorded, moves) values (?, ?, ?, ?)"
	stmt, err := Db.Data.Prepare(statement)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(p1_id, p2_id, time.Now(), game_record)
	return err
}
