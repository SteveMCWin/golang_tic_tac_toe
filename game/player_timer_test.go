package game

import (
    "fmt"
    "testing"
    "time"
)

func TestTimerRunsOut(t *testing.T) {
	fmt.Println("Started TestTimerRunsOut")
	pt := MakePlayerTimer(1 * time.Second)
	pt.Start()

	select {
	case <-pt.Finished:
		fmt.Println("Timer finished as expected")
	case <-time.After(2 * time.Second):
		t.Error("Timer didn't finish on time!")
	}
}
//
// func TestTimerPause(t *testing.T) {
//     fmt.Println("Started TestTimerPause")
//     pt := MakePlayerTimer(2 * time.Second)
//     pt.Start()
//     pt_finished = false
//
//     go func() {
//         t := time.NewTimer(3 * time.Second)
//         select {
//         case <- t.C:
//             return
//         case <- pt.Finished:
//             t.Error("Timer finished even after being stopped")
//         }
//     }()
//
//     pt.Pause()
//     time.Sleep(3 * time.Second)
// }
//
// func TestTimerRace(t *testing.T) {
//     fmt.Println("Started TestTimerRace")
//     pt1 := MakePlayerTimer(3 * time.Second)
//     pt2 := MakePlayerTimer(3 * time.Second)
//
//     pt1.Start()
//
//     time.Sleep(1 * time.Second)
//     pt2.Start()
//
//     select {
//     case <-pt1.Finished:
//         return
//     case <-pt2.Finished:
//         t.Error("Timer 2 finished efore timer 1")
//     }
//
// }
//
// func TestTimerRace2(t *testing.T) {
//     fmt.Println("Started TestTimerRace2")
//     pt1 := MakePlayerTimer(3 * time.Second)
//     pt2 := MakePlayerTimer(3 * time.Second)
//
//     pt1.Start()
//
//     time.Sleep(1 * time.Second)
//     pt1.Stop()
//     pt2.Start()
//
//     select {
//     case <-pt1.Finished:
//         t.Error("Timer 1 finished efore timer 2")
//     case <-pt2.Finished:
//         return
//     }
//
// }
//
// func TestPauseAccuracy(t *testing.T) {
//     fmt.Println("Started TestPauseAccuracy")
//     delta := 20 * time.Millisecond
//
//     pt1 := MakePlayerTimer(3 * time.Second)
//     pt1.Start()
//     time.Sleep(1000 * time.Millisecond)
//     pt1.Stop()
//     if time.Abs(pt1.TimeLeft - 2000 * time.Millisecond) > delta {
//         t.Error("Timer unprecise. Expected about 2000ms, got", pt1.TimeLeft)
//     }
//     time.Sleep(200 * time.Millisecond)
//     pt1.Start()
//     time.Sleep(1000 * time.Millisecond)
//     pt1.Stop()
//     if time.Abs(pt1.TimeLeft - 1000 * time.Millisecond) > delta {
//         t.Error("Timer unprecise. Expected about 1000ms, got", pt1.TimeLeft)
//     }
// }
//
