// Package models holds the golang representation of data models stored in the sqlite database
// and provides functions that operate on that data.
package models

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

// DataBase is used to make changes to the actual database.db file with sqlite querries.
type DataBase struct {
	Data    *sql.DB // the connection to the database through which all operations on the said database are preformed
	is_open bool
}

// Close handles the closing of a connection to the databse
func (dataBase *DataBase) Close() {
	dataBase.Data.Close()
	dataBase.is_open = false
}

// initializes the database
// if any parameters are passed, uses the test_database

// InitDatabase opens a connection to the database and loads the needed extensions.
func (Db *DataBase) InitDatabase(is_test ...bool) error {
	if Db.is_open {
		return errors.New("ERROR: Database already open")
	}

	dir, err := os.Getwd()
	log.Println("dir:", dir)
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

