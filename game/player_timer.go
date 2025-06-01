package game

import (
    "time"
)

type PlayerTimer struct {
    TimeLeft time.Duration
    isPaused bool
    Finished chan bool
    lastFrame time.Time
    FischerTime time.Duration
}

func MakePlayerTimer(d time.Duration, ft time.Duration) (pt *PlayerTimer){
    pt = &PlayerTimer{}
    pt.TimeLeft = d
    pt.isPaused = true
    pt.Finished = make(chan bool, 1)
    pt.lastFrame = time.Now()
    pt.FischerTime = ft

    return
}

func (pt *PlayerTimer) Start() {
    pt.lastFrame = time.Now()
    pt.isPaused = false
    go pt.run()
}

func (pt *PlayerTimer) Pause() {
    pt.TimeLeft += pt.FischerTime
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

        time.Sleep(10 * time.Millisecond)
    }
}
