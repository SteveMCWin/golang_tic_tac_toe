package models

import (
	"database/sql"
	"errors"
	"os"

	"github.com/mattn/go-sqlite3"
)

type DataBase struct {
	Data    *sql.DB
	is_open bool
}

func (dataBase *DataBase) Close() {
	dataBase.Data.Close()
	dataBase.is_open = false
}

// initializes the database
// if any parameters are passed, uses the test_database
func (Db *DataBase) InitDatabase(is_test ...bool) error {
	if Db.is_open {
		return errors.New("ERROR: Database already open")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	spellfix_relative_path := "/extensions/spellfix.so"

	sql.Register("sqlite3_with_extension",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				dir + spellfix_relative_path,
			},
		},
	)

	db_path := "data/"
	if len(is_test) == 0 {
		db_path = db_path + "database.db"
	} else {
		db_path = db_path + "test_database.db"
	}

	Db.Data, err = sql.Open("sqlite3_with_extension", db_path)
	if err != nil {
		return err
	}

	Db.is_open = true

	return nil
}

