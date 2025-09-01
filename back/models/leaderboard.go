package models

import (
    "log"
    "time"
	"tic_tac_toe.fun/defs"
)

// LeaderBoard is basically just a slice of users, which holds the best players based on elo.
type LeaderBoard struct {
    TopPlayers []*User
	UpdateDuration time.Duration
}

// InitLeaderBoard just creates an empty slice of users and sets the time duration it takes the leaderboard to update.
// If the duration is not passed, a default duration defined in defs.go is used.
func (lb *LeaderBoard) InitLeaderBoard(updateDuration ...time.Duration) {
    lb.TopPlayers = make([]*User, 0)
	if len(updateDuration) > 0 {
		lb.UpdateDuration = updateDuration[0]
	} else {
		lb.UpdateDuration = defs.LEADERBOARD_UPDATE_DURATION
	}
}

// RunLeaderBoard starts an infinite goroutine that updates the leaderboard after a duration determined when the leaderboard was initialized.
func (lb *LeaderBoard) RunLeaderBoard(Db *DataBase) {
	go func() {
		for {
			lb.updateLeaderboard(Db)
			time.Sleep(lb.UpdateDuration)
		}
	}()
}

func (lb *LeaderBoard) updateLeaderboard(Db *DataBase) { // just pulls the top players (based on elo) from the db
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
