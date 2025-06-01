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

func InitDb() { // opens db connection and sets up the leaderboard and runs it because the leaderboard needs to wait for the db to be open
    once.Do(func() {
        var err error
        Db, err = sql.Open("sqlite3", "users/users.db")
        if err != nil {
            panic(err)
        }

        MostEloLb, err = InitLeaderBoard()
        if err != nil {
            panic(err)
        }

        go MostEloLb.RunLeaderBoard()
    })
}

func LoadUserData(c *gin.Context) (usr *User, err error) {  // gets a data from a user stored in the db based on the user_id stored in a cookie in the browser

    usr = &User{}

    user, err := c.Cookie("user_id")    // get the user_id from a cookie stored in the browser
    if err != nil {
        usr.UserName = "Guest"  // if there is not user_id cookie, the player is logged out/never logged in and is considered a guest
        usr.Elo = 800   // set the elo so that when a logged in user and a guest play, the logged in user loses/gains elo, which would not happen when the guest has 0 elo
        return
    }

    csrf, err := c.Cookie("csrf_token") // if there is no csrf_token the user cannot obtain the data in the db and must log in again
    if err != nil {
        return
    }

    sess, err := c.Cookie("session_token")  // if there is no session_token the user cannot obtain the data in the db and must log in again
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
    )   // just get everything haha

    if err != nil {
        return
    }

    // before returning the user check if the browser tokens are the same as the db tokens
    // if not, someone is trying to steal data
    if usr.SessionToken != sess || usr.CSRFToken != csrf {  
        usr = &User{}    // make sure to return an empty user if login fails
        err = errors.New("Session token or csrf token missmatch, please log in again")
    }

    return
}

func (usr *User) StoreUser() (err error) {  // writes the user to the db based on his email and gets the auto-generated id to set as a cookie in auth.go
    if usr.Email == "" {
        err = errors.New("Could not store user data: email missing")
        return
    }

    var prov string
    // check if a user with this email has already logged in before and with which provider
    err = Db.QueryRow("select id, provider from users where email like ?", usr.Email).Scan(&usr.Id, &prov)

    if err != nil { // the user hasn't logged in before so load him into the data base
        // store the user in the database
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

    // if the user has logged in before make sure the username and avatar get updated in case the user changed them
    _, err = Db.Exec("update users set username = ?, avatar_url = ?, session_token = ?, csrf_token = ?, provider = ? where id = ?",
                      usr.UserName, usr.AvatarURL, usr.SessionToken, usr.CSRFToken, usr.Provider, usr.Id)

    return
}

// not used but still here in case I find it useful somehow
func (usr *User) storeUserPfp() (err error) {   // used to locally store the user's avatar as an image file
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

func (usr *User) UpdateGameStats() (err error) {    // called at the end of the game to update only the gameplay related stats
    if usr.Id == 0 {
        return // cannot store data for a guest user (all guest users have the id = 0)
    }

    _, err = Db.Exec("update users set games_played = ?, games_won = ?, elo = ? where id = ?", usr.GamesPlayed, usr.GamesWon, usr.Elo, usr.Id)

    return

}






