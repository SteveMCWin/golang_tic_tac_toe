package models

import (
    "log"
    "time"
)

type LeaderBoard struct {
    TopPlayers []*User
}

func (lb *LeaderBoard) InitLeaderBoard() {
    lb.TopPlayers = make([]*User, 0)
}

func (lb *LeaderBoard) RunLeaderBoard(Db *DataBase) {   // updates the leaderboard every 10 seconds, intended to be called in a goroutine
    for {
        lb.updateLeaderboard(Db)
        time.Sleep(10 * time.Second)
    }
}

func (lb *LeaderBoard) updateLeaderboard(Db *DataBase) {    // just pulls the top players (based on elo) from the db
    stmt := "SELECT username, avatar_url, elo FROM users ORDER BY elo DESC LIMIT 10"
    rows, err := Db.Data.Query(stmt);
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
