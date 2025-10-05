// Package models holds the golang representation of data models stored in the sqlite database
// and provides functions that operate on that data.
package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

// InitDatabase opens a connection to the database and loads the needed extensions.
// If a boolean parameter is passed (no matter it's value), a test database will be used.
// If a database doesn't exist in data/ it gets created and initalized with sql scripts from the same directory
func (Db *DataBase) InitDatabase(is_test ...bool) error {
	if Db.is_open {
		return errors.New("ERROR: Database already open")
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	spellfix_path := filepath.Join(dir, "extensions", "spellfix.so")

	sql.Register("sqlite3_with_extension",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				spellfix_path,
			},
		},
	)

	data_dir := "data"
	if err := os.MkdirAll(data_dir, 0755); err != nil {
		return err
	}

	db_name := "database.db"
	if len(is_test) != 0 {
		db_name = "test_database.db"
	}

	db_path := filepath.Join(data_dir, db_name)

	// Check if database already exists
	_, err = os.Stat(db_path)
	dbExists := !os.IsNotExist(err)

	Db.Data, err = sql.Open("sqlite3_with_extension", db_path)
	if err != nil {
		return err
	}

	if dbExists != true {
		sqlFiles := []string{
			"data/create_user_table.sql",
			"data/create_spellfix_user_table.sql",
			"data/create_sessions_table.sql",
			"data/create_game_table.sql",
		}

		for _, sqlFile := range sqlFiles {
			log.Printf("Running %s...\n", sqlFile)

			// Open the SQL file
			file, err := os.Open(sqlFile)
			if err != nil {
				log.Println("failed to open sql file:", sqlFile, err)
				return err
			}

			// Execute: sqlite3 database.db < sqlFile
			cmd := exec.Command("sqlite3", db_path)
			cmd.Stdin = file
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				file.Close()
				log.Println("failed to execute sql file:", sqlFile, err)
				return err
			}

			file.Close()
		}
	}

	Db.is_open = true

	return nil
}
