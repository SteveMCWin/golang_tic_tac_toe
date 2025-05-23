package users

import (
    "log"
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

func init() {
    var err error
    Db, err = sql.Open("sqlite3", "users/users.db")
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

    log.Println("user_id:", user)

    csrf, err := c.Cookie("csrf_token")
    if err != nil {
        return
    }

    sess, err := c.Cookie("session_token")
    if err != nil {
        return
    }

    err = Db.QueryRow("select id, username, email, session_token, csrf_token, provider from users where id = ?", user).Scan(
        &usr.Id,
        &usr.UserName,
        &usr.Email,
        &usr.SessionToken,
        &usr.CSRFToken,
        &usr.Provider,
    )

    if err != nil {
        log.Println("The user wasn't loaded for comparing\n")
        log.Println("error:", err)
    }

    if usr.SessionToken != sess || usr.CSRFToken != csrf {
        usr = User{}
        err = errors.New("Session token or csrf token missmatch")
        log.Println("EXPECTED")
        log.Println("sess:\t", sess)
        log.Println("csrf:\t", csrf)
        log.Println("RECIEVED")
        log.Println("sess:\t", usr.SessionToken)
        log.Println("csrf:\t", usr.CSRFToken)
    }

    return
}

func (usr *User) AddUser() (err error) {
    log.Println("Trying to get user named: ", usr.UserName)
    err = Db.QueryRow("select id from users where username like ?", usr.UserName).Scan(&usr.Id)

    // the user hasn't logged in before so load him into the data base
    if err != nil {
        log.Println("IT DID NOT RECOGNIZE THE USER")
        statement := "insert into users (username, email, session_token, csrf_token, provider) values (?, ?, ?, ?, ?) returning id"
        var stmt *sql.Stmt
        stmt, err = Db.Prepare(statement)
        if err != nil {
            return
        }
        defer stmt.Close()
        err = stmt.QueryRow(usr.UserName, usr.Email, usr.SessionToken, usr.CSRFToken, usr.Provider).Scan(&usr.Id)
        return
    }

    log.Println("IT RECOGNIZED THE USER FROM BEFORE")
    log.Println("THE NEW COOKIES ARE")
    log.Println("sess:\t", usr.SessionToken)
    log.Println("csrf:\t", usr.CSRFToken)

    // the user has logged in before so just update the tokens
    _, err = Db.Exec("update users set session_token = ?, csrf_token = ? where id = ?", usr.SessionToken, usr.CSRFToken, usr.Id)
    return
}


