package users

import (
    "log"
    "time"
)

var LeaderBoard []*User

func updateLeaderboard() {
    rows, err := Db.Query("SELECT username, games_won FROM users ORDER BY games_won DESC LIMIT 10");
    if err != nil {
        log.Println("update leaderboard error:")
        log.Println(err)
    }
    defer rows.Close()

    err = rows.Err()
    if err != nil {
        log.Println("update leaderboard error:")
        log.Println(err)
    }

    LeaderBoard = make([]*User, 0)

    for rows.Next() {
        usr := User{}
        err = rows.Scan(&usr.UserName, &usr.GamesWon)
        if err != nil {
            log.Println("Error scanning leaderboard user")
            log.Println(err)
        }

        LeaderBoard = append(LeaderBoard, &usr)
    }
}

func RunLeaderBoard() {
    for {
        updateLeaderboard()
        log.Println("LeaderBoard:")
        for _, usr := range LeaderBoard {
            log.Printf("%s:\t%d", usr.UserName, usr.GamesWon)
        }
        time.Sleep(10 * time.Second)
    }
}
