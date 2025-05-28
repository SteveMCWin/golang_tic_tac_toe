package game

import (
    // "log"
    "time"
)

// TODO: this needs to be tested in the game setting because it's acting really weird sometimes
type PlayerTimer struct {
    TimeLeft time.Duration
    isPaused bool
    Finished chan bool
    lastFrame time.Time
}

func MakePlayerTimer(d time.Duration) (pt *PlayerTimer){
    pt = &PlayerTimer{}
    pt.TimeLeft = d
    pt.isPaused = true
    pt.Finished = make(chan bool, 1)
    pt.lastFrame = time.Now()

    return
}

func (pt *PlayerTimer) Start() {
    // I think I have to update last frame here
    pt.isPaused = false
    go pt.run()
}

func (pt *PlayerTimer) Pause() {
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
