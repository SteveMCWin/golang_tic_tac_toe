package game

import (
    // "log"
    "time"
)

// TODO: perhaps rework this so it has a goroutine that just sleeps for a second and then decrements the time

type PlayerTimer struct {
    Tmr *time.Timer
    last_start_time time.Time
    TimeLeft time.Duration
    IsPaused bool
}

func MakePlayerTimer(d time.Duration) (pt *PlayerTimer){
    pt = &PlayerTimer{}
    pt.Tmr = time.NewTimer(10 * time.Second)
    pt.Tmr.Stop()
    pt.TimeLeft = d
    pt.IsPaused = true

    return
}

func (pt *PlayerTimer) Start() {

    if pt.IsPaused == false {
        return
    }
    
    pt.Tmr.Reset(pt.TimeLeft)
    pt.last_start_time = time.Now()
    pt.IsPaused = false
}

func (pt *PlayerTimer) Pause() {

    if pt.IsPaused == true {
        return
    }

    pt.IsPaused = true
    pt.TimeLeft = pt.TimeLeft - time.Since(pt.last_start_time)
    pt.Tmr.Stop()
}
