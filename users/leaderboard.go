package users

import (
    "log"
    "time"
    "errors"
)

type LbCriteria int

type LeaderBoard struct {
    TopPlayers []*User
    Criteria LbCriteria
}

const (
    MostWins LbCriteria = iota
    MostGames
    MostElo
    LbCriteriaSize
)

// var MostWinsLb *LeaderBoard
// var MostGamesLb *LeaderBoard
// var BestWinrateLb *LeaderBoard

func InitLeaderBoard(crit LbCriteria) (*LeaderBoard, error) {
    if crit < LbCriteria(0) || crit >= LbCriteriaSize {
        return nil, errors.New("Wrong leaderboard criteria")
    }

    lb := &LeaderBoard{ TopPlayers: make([]*User, 0), Criteria: crit}

    return lb, nil
}

func (lb *LeaderBoard) updateLeaderboard() {
    var stmt string
    switch lb.Criteria {
    case MostWins:
        stmt = "SELECT username, games_won FROM users ORDER BY games_won DESC LIMIT 10"
    case MostGames:
        stmt = "SELECT username, games_played FROM users ORDER BY games_played DESC LIMIT 10"
    case MostElo:
        stmt = "SELECT username, games_played FROM users ORDER BY elo DESC LIMIT 10"
    default:
        stmt = "SELECT username, games_won FROM users ORDER BY games_won DESC LIMIT 10"
    }
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
        err = rows.Scan(&usr.UserName, &usr.GamesWon)
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
        // log.Print("LeaderBoard: ")
        // if lb.Criteria == MostWins {
        //     log.Print("MostWins")
        // } else {
        //     log.Print("GamesPlayed")
        // }
        // for _, usr := range lb.TopPlayers {
        //     log.Printf("%s:\t%d", usr.UserName, usr.GamesWon)
        // }
        time.Sleep(10 * time.Second)
    }
}
