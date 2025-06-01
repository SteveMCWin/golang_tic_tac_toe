package users

import (
    "log"
    "time"
    // "errors"
    "net/http"
	"github.com/gin-gonic/gin"
)

type LeaderBoard struct {
    TopPlayers []*User
}

var MostEloLb *LeaderBoard

func InitLeaderBoard() (*LeaderBoard, error) {   // is called in the user_data's InitDb function
    lb := &LeaderBoard{ TopPlayers: make([]*User, 0) }

    return lb, nil
}

func (lb *LeaderBoard) updateLeaderboard() {
    stmt := "SELECT username, avatar_url, elo FROM users ORDER BY elo DESC LIMIT 10"
    rows, err := Db.Query(stmt);
    if err != nil {
        log.Println("update leaderboard error:")
        log.Println(err)
        return
    }
    defer rows.Close()

    err = rows.Err()
    if err != nil {
        log.Println("update leaderboard error:")
        log.Println(err)
        return
    }

    lb.TopPlayers = make([]*User, 0)

    for rows.Next() {
        usr := User{}
        err = rows.Scan(&usr.UserName, &usr.AvatarURL, &usr.Elo)
        if err != nil {
            log.Println("Error scanning leaderboard user")
            log.Println(err)
        }

        lb.TopPlayers = append(lb.TopPlayers, &usr)
    }
}

func (lb *LeaderBoard) RunLeaderBoard() {
    for {
        lb.updateLeaderboard()
        time.Sleep(10 * time.Second)
    }
}

func ServeLeaderboard(c *gin.Context) {
    c.HTML(http.StatusOK, "leaderboard.html", gin.H{
        "TopPlayers": MostEloLb.TopPlayers,
    })
}
