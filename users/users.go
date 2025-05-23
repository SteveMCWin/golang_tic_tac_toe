package users

import (
    // "fmt"
    // "log"
    "errors"

	"github.com/gin-gonic/gin"

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

func LoadUserData(c *gin.Context) (usr User, err error) {

    usr = User{}

    user, err := c.Cookie("user_id")
    if err != nil {
        return
    }

    csrf, err := c.Cookie("csrf_token")
    if err != nil {
        return
    }

    sess, err := c.Cookie("session_token")
    if err != nil {
        return
    }

    err = Db.QueryRow("select id, username, email, session_token, csrftoken, provider from users where id = ?", user).Scan(
        &usr.Id,
        &usr.UserName,
        &usr.Email,
        &usr.SessionToken,
        &usr.CSRFToken,
        &usr.Provider,
    )

    if usr.SessionToken != sess || usr.CSRFToken != csrf {
        usr = User{}
        err = errors.New("Session token or csrf token missmatch")
    }

    return
}

func (usr *User) AddUser() (err error) {
    err = Db.QueryRow("select id where username like '?'", usr.UserName).Scan(&usr.Id)

    // the user hasn't logged in before so load him into the data base
    if err != nil {
        statement := "insert into users (username, email, session_token, csrf_token, provider) values (?, ?, ?, ?, ?) returning id"
        var stmt *sql.Stmt
        stmt, err = Db.Prepare(statement)
        if err != nil {
            // panic(err)
            return
        }
        defer stmt.Close()
        err = stmt.QueryRow(usr.UserName, usr.Email, usr.SessionToken, usr.CSRFToken, usr.Provider).Scan(&usr.Id)
        return
    }

    // the user has logged in before so just update the tokens
    _, err = Db.Exec("update users set session_token = ?, csrf_token = ? where id = ?", usr.SessionToken, usr.CSRFToken, usr.Id)
    return
}


