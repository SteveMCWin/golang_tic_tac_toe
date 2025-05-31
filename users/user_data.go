package users

import (
    "io"
    "os"
    "log"
    "sync"
    "errors"
    "net/http"

	"github.com/gin-gonic/gin"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
    Id              int
    UserName        string
    Email           string
    AvatarURL       string
    SessionToken    string
    CSRFToken       string
    Provider        string
    // Game related
    GamesPlayed     int
    GamesWon        int
    Elo             int
}

// TODO: consider splitting the users table into two tables, users and players. they will share the keys
var Db *sql.DB
var once sync.Once

var MostWinsLb *LeaderBoard
var MostGamesLb *LeaderBoard
var MostEloLb *LeaderBoard

func InitDb() {
    once.Do(func() {
        var err error
        Db, err = sql.Open("sqlite3", "users/users.db")
        if err != nil {
            panic(err)
        }

        MostWinsLb, err = InitLeaderBoard(MostWins)
        if err != nil {
            panic(err)
        }

        MostGamesLb, err = InitLeaderBoard(MostGames)
        if err != nil {
            panic(err)
        }

        MostEloLb, err = InitLeaderBoard(MostElo)

        go MostWinsLb.RunLeaderBoard()
        go MostGamesLb.RunLeaderBoard()
        go MostEloLb.RunLeaderBoard()
    })
}

// change this to be (db *sql.Db) LoadUserData(c *gin.Context) (usr *User, err error)
func LoadUserData(c *gin.Context) (usr *User, err error) {

    usr = &User{}

    user, err := c.Cookie("user_id")
    if err != nil {
        usr.UserName = "Guest"
        usr.Elo = 800
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

    err = Db.QueryRow("select id, username, email, avatar_url, session_token, csrf_token, provider, games_played, games_won, elo from users where id = ?", user).Scan(
        &usr.Id,
        &usr.UserName,
        &usr.Email,
        &usr.AvatarURL,
        &usr.SessionToken,
        &usr.CSRFToken,
        &usr.Provider,
        &usr.GamesPlayed,
        &usr.GamesWon,
        &usr.Elo,
    )

    if err != nil {
        return
    }

    if usr.SessionToken != sess || usr.CSRFToken != csrf {
        usr = &User{}    // make sure to return an empty user if login fails
        err = errors.New("Session token or csrf token missmatch")
    }

    return
}

func (usr *User) StoreUser() (err error) {
    // log.Println("Trying to get user with email: ", usr.Email)
    if usr.Email == "" {
        err = errors.New("Could not store user data: email missing")
        return
    }

    var prov string
    err = Db.QueryRow("select id, provider from users where email like ?", usr.Email).Scan(&usr.Id, &prov)

    // log.Println(usr)
    // the user hasn't logged in before so load him into the data base
    if err != nil {
        log.Println("IT DID NOT RECOGNIZE THE USER")
        statement := "insert into users (username, email, avatar_url, session_token, csrf_token, provider, games_played, games_won, elo) values (?, ?, ?, ?, ?, ?, ?, ?, ?) returning id"
        var stmt *sql.Stmt
        stmt, err = Db.Prepare(statement)
        if err != nil {
            return
        }
        defer stmt.Close()
        err = stmt.QueryRow(usr.UserName, usr.Email, usr.AvatarURL, usr.SessionToken, usr.CSRFToken, usr.Provider, usr.GamesPlayed, usr.GamesWon, usr.Elo).Scan(&usr.Id)
        return
    }

    log.Println("IT RECOGNIZED THE USER FROM BEFORE")

    if prov != usr.Provider {
        _, err = Db.Exec("update users set session_token = ?, csrf_token = ? where id = ?",
                          usr.SessionToken, usr.CSRFToken, usr.Id)
    } else {
        _, err = Db.Exec("update users set username = ?, avatar_url = ?, session_token = ?, csrf_token = ?, provider = ? where id = ?",
                          usr.UserName, usr.AvatarURL, usr.SessionToken, usr.CSRFToken, usr.Provider, usr.Id)
    }
    return
}

// not used but still here in case I find it useful somehow
func (usr *User) storeUserPfp() (err error) {
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

func (usr *User) UpdateGameStats() (err error) {
    if usr.Id == 0 {
        return // cannot store data for a guest user
    }

    _, err = Db.Exec("update users set games_played = ?, games_won = ?, elo = ? where id = ?", usr.GamesPlayed, usr.GamesWon, usr.Elo, usr.Id)

    return

}






