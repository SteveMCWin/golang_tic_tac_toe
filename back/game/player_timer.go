package game

import (
    "time"
)

// PlayerTimer is a countdown timer that can be paused and unpaused during a game.
// Note that the player timer can be a standard timer or a Fischer timer.
// To make a normal timer just set the fischer time to 0 when creating the player timer.
type PlayerTimer struct {
    TimeLeft time.Duration
    isPaused bool
    Finished chan bool // A value will be sent to this channel once the timer runs out
    lastFrame time.Time // last recorded time, used to decrease the TimeLeft based on how long ago it was recorded
    FischerTime time.Duration   // Used to add time to TimeLeft when a player ends their turn, look it up on wikipedia
}

// MakePlayerTimer creates a timer based on the time and fischer time passed. Note that the timer is paused on creation.
func MakePlayerTimer(d time.Duration, ft time.Duration) (pt *PlayerTimer) {
    pt = &PlayerTimer{}
    pt.TimeLeft = d
    pt.isPaused = true
    pt.Finished = make(chan bool, 1)
    pt.lastFrame = time.Now()
    pt.FischerTime = ft

    return
}

// Start unpauses the timer.
func (pt *PlayerTimer) Start() {
    pt.lastFrame = time.Now()
    pt.isPaused = false
    go pt.run()
}

// Pause does the opposite of Start
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

        time.Sleep(10 * time.Millisecond)   // this is here so the app has time to do something else, sacrificing some precision for some performance
    }
}
