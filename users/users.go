package users

import (
    // "fmt"
    "log"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
    Id              int
    UserName        string
    Email           string
    SessionToken    string
    CSRFToken       string
    Provider        string
}

var Db *sql.DB

var Users map[string]User


// usr_name := gothic.GetFromSession("username", )

func init() {
    Users = make(map[string]User)
    var err error
    Db, err = sql.Open("sqlite3", "users.db")
    if err != nil {
        panic(err)
    }
}

func AddActiveUser(username, email, sessionToken, csrfToken, provider string) {
    log.Println("Added", username)
    Users[username] = User{-1, username, email, sessionToken, csrfToken, provider}
    log.Println("map: ", Users)
}


