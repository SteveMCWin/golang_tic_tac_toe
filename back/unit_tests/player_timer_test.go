package unit_tests

import (
    "testing"
    "time"
    "tic_tac_toe.fun/game"
)

func TestTimerRunsOut(t *testing.T) {
    t.Parallel()
	pt := game.MakePlayerTimer(1 * time.Second, 2 * time.Second)
	pt.Start()

	select {
	case <-pt.Finished:
	case <-time.After(2 * time.Second):
		t.Error("Timer didn't finish on time!")
	}
}

func TestTimerPause(t *testing.T) {
    t.Parallel()
	pt := game.MakePlayerTimer(1 * time.Second, 2 * time.Second)
	pt.Start()

    time.Sleep(50 * time.Millisecond)

    pt.Pause()

	select {
	case <-pt.Finished:
		t.Error("Timer finished even when stopped!")
	case <-time.After(2 * time.Second):
	}
}

func TestTimerRace(t *testing.T) {
    t.Parallel()
    pt1 := game.MakePlayerTimer(3 * time.Second, 2 * time.Second)
    pt2 := game.MakePlayerTimer(3 * time.Second, 2 * time.Second)

    pt1.Start()

    time.Sleep(1 * time.Second)
    pt2.Start()

    select {
    case <-pt1.Finished:
        return
    case <-pt2.Finished:
        t.Error("Timer 2 finished efore timer 1")
    }

}

func TestTimerRace2(t *testing.T) {
    t.Parallel()
    pt1 := game.MakePlayerTimer(3 * time.Second, 0 * time.Second)
    pt2 := game.MakePlayerTimer(3 * time.Second, 0 * time.Second)

    pt1.Start()

    time.Sleep(1 * time.Second)
    pt1.Pause()
    pt2.Start()

    select {
    case <-pt1.Finished:
        t.Error("Timer 1 finished efore timer 2")
    case <-pt2.Finished:
        return
    }

}

func TestPauseAccuracy(t *testing.T) {
    t.Parallel()
    delta := 20 * time.Millisecond

    pt1 := game.MakePlayerTimer(3 * time.Second, 0 * time.Second)
    pt1.Start()
    time.Sleep(1000 * time.Millisecond)
    pt1.Pause()
    if (pt1.TimeLeft - 2000 * time.Millisecond).Abs() > delta {
        t.Error("Timer unprecise. Expected about 2000ms, got", pt1.TimeLeft)
    }
    time.Sleep(200 * time.Millisecond)
    pt1.Start()
    time.Sleep(1000 * time.Millisecond)
    pt1.Pause()
    if (pt1.TimeLeft - 1000 * time.Millisecond).Abs() > delta {
        t.Error("Timer unprecise. Expected about 1000ms, got", pt1.TimeLeft)
    }
}

func TestFischerPauseAccuracy(t *testing.T) {
    t.Parallel()
    delta := 20 * time.Millisecond

    pt1 := game.MakePlayerTimer(3 * time.Second, 1 * time.Second)
    pt1.Start()
    time.Sleep(1000 * time.Millisecond)
    pt1.Pause()
    if (pt1.TimeLeft - 3000 * time.Millisecond).Abs() > delta {
        t.Error("Timer unprecise. Expected about 3000ms, got", pt1.TimeLeft)
    }
    time.Sleep(200 * time.Millisecond)
    pt1.Start()
    time.Sleep(2000 * time.Millisecond)
    pt1.Pause()
    if (pt1.TimeLeft - 2000 * time.Millisecond).Abs() > delta {
        t.Error("Timer unprecise. Expected about 2000ms, got", pt1.TimeLeft)
    }
}

