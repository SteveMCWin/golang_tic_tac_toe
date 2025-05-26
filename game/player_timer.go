package game

import (
    // "log"
    "time"
)

// TODO: perhaps rework this so it has a goroutine that just sleeps for a second and then decrements the time

type PlayerTimer struct {
    // Tmr *time.Timer
    // last_start_time time.Time
    // TimeLeft time.Duration
    // IsPaused bool
    TimeLeft time.Duration
    isPaused bool
    Finished chan bool
    lastFrame time.Time
}

func MakePlayerTimer(d time.Duration) (pt *PlayerTimer){
    // pt = &PlayerTimer{}
    // pt.Tmr = time.NewTimer(10 * time.Second)
    // pt.Tmr.Stop()
    // pt.TimeLeft = d
    // pt.IsPaused = true
    //
    // return

    pt = &PlayerTimer{}
    pt.TimeLeft = d
    pt.isPaused = true
    pt.Finished = make(chan bool)
    pt.lastFrame = time.Now()

    return
}

func (pt *PlayerTimer) Start() {

    // if pt.IsPaused == false {
    //     return
    // }
    // 
    // pt.Tmr.Reset(pt.TimeLeft)
    // pt.last_start_time = time.Now()
    // pt.IsPaused = false

    pt.isPaused = false
    go pt.run()
}

func (pt *PlayerTimer) Pause() {

    // if pt.IsPaused == true {
    //     return
    // }
    //
    // pt.IsPaused = true
    // pt.TimeLeft = pt.TimeLeft - time.Since(pt.last_start_time)
    // pt.Tmr.Stop()
    pt.isPaused = true
}

func (pt *PlayerTimer) run() {
    for {
        if pt.isPaused == true {
            return
        }
        dt := time.Since(pt.lastFrame)
        pt.lastFrame = time.Now()
        pt.TimeLeft -= dt
        if pt.TimeLeft <= time.Millisecond {
            pt.Finished <- true
            pt.isPaused = true
            return
        }
    }
}
