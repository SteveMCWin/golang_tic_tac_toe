package models

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"tic_tac_toe.fun/defs"
)

// User is a subset of data provided by a provider when logging in, plus game related stats.
type User struct {
	Id        int    `json:"id"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Provider  string `json:"provider"` // Represents the last provider used to log in to the account with the stored email. Actually useless at the moment.
	// Game related
	NumOfGamesPlayed int `json:"numOfGamesPlayed"`
	NumOfGamesWon    int `json:"numOfGamesWon"`
	Elo              int `json:"elo"`
}

// ReadUser populates a User struct with data linked to the id passed to the function.
// If the id passed is defs.NO_USER_ID, the User struct returned will only have a UserName "Guest" and a starting elo represented by defs.STARTING_ELO.
func (Db *DataBase) ReadUser(user_id int) (*User, error) {

	usr := &User{Id: user_id}

	if user_id == defs.NO_USER_ID {
		usr = &User{UserName: "Guest", Elo: defs.STARTING_ELO}
		return usr, nil
	}

	err := Db.Data.QueryRow("select username, email, avatar_url, provider, games_played, games_won, elo from users where id = ?", user_id).Scan(
		&usr.UserName,
		&usr.Email,
		&usr.AvatarURL,
		&usr.Provider,
		&usr.NumOfGamesPlayed,
		&usr.NumOfGamesWon,
		&usr.Elo,
	)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

// StoreUser writes the data inside the passed in User struct to the database, unless the email from the User struct already
// exists in the database. In that case the username, avatar_url and the provider are updated in the databse.
func (Db *DataBase) StoreUser(usr *User) error { // writes the user to the db based on his email and gets the auto-generated id to set as a cookie in auth.go
	if usr.Email == "" {
		err := errors.New("Could not store user data: email missing")
		return err
	}

	var prov string
	// check if a user with this email has already logged in before and with which provider
	err := Db.Data.QueryRow("select id, provider from users where email like ?", usr.Email).Scan(&usr.Id, &prov)

	if err != nil { // the user hasn't logged in before so load him into the data base
		// store the user in the database
		statement := "insert into users (username, email, avatar_url, provider, games_played, games_won, elo) values (?, ?, ?, ?, ?, ?, ?) returning id"
		var stmt *sql.Stmt
		stmt, err = Db.Data.Prepare(statement)
		if err != nil {
			return err
		}
		defer stmt.Close()
		err = stmt.QueryRow(usr.UserName, usr.Email, usr.AvatarURL, usr.Provider, defs.DEFAULT_GAMES_PLAYED, defs.DEFAULT_GAMES_WON, defs.STARTING_ELO).Scan(&usr.Id)
		return err
	}

	// if the user has logged in before make sure the username and avatar get updated in case the user changed them
	_, err = Db.Data.Exec("update users set username = ?, avatar_url = ?, provider = ? where id = ?",
		usr.UserName, usr.AvatarURL, usr.Provider, usr.Id)

	return err
}

func (Db *DataBase) DeleteUser(user_id int) error {
	statement_delete_profile := "delete from users where id = ?"
	stmt_delete_profile, err := Db.Data.Prepare(statement_delete_profile)
	if err != nil {
		return err
	}

	defer stmt_delete_profile.Close()

	_, err = stmt_delete_profile.Exec(user_id)
	if err != nil {
		return err
	}

	log.Println("Going into UpdateRecordsDeletedUser")
	err = Db.UpdateRecordsDeletedUser(user_id)
	return err
}

// not used but still here in case I find it useful somehow
func (usr *User) storeUserPfp() (err error) { // used to locally store the user's avatar as an image file
	log.Println("AVATAR URL:", usr.AvatarURL)
	img_response, err := http.Get(usr.AvatarURL)
	if err != nil {
		return err
	}
	defer img_response.Body.Close()

	img_filepath := "users/pfps/pfp" + usr.UserName
	file, err := os.Create(img_filepath)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, img_response.Body)
	if err != nil {
		return err
	}

	return nil
}

// UpdateGameStats changes the values of the game related data inside the passed in User struct inside the databse.
// UpdateGameStats is intended to be called after a game concludes.
// It assumes the Elo, NumOfGamesPlayed and NumOfGamesWon inside the User struct have already been affected by the outcome of the game and shouldn't be altered.
func (Db *DataBase) UpdateGameStats(usr *User) error { // called at the end of the game to update only the gameplay related stats
	if usr.Id == defs.NO_USER_ID {
		return errors.New("Cannot update player stats when the player isn't logged in") // cannot store data for a guest user (all guest users have the id = 0)
	}

	_, err := Db.Data.Exec("update users set games_played = ?, games_won = ?, elo = ? where id = ?", usr.NumOfGamesPlayed, usr.NumOfGamesWon, usr.Elo, usr.Id)

	return err
}

func (Db *DataBase) SearchForUsers(username string, requesting_user_id int) ([]User, error) {
	rows, err := Db.Data.Query("select id, username, avatar_url from spellfix_users inner join users on word = username where word match ? and id != ?", username, requesting_user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	matches := make([]User, 0)

	for rows.Next() {
		var usr User
		err = rows.Scan(&usr.Id, &usr.UserName, &usr.AvatarURL)
		if err != nil {
			return nil, err
		}

		matches = append(matches, usr)
	}

	return matches, nil
}
