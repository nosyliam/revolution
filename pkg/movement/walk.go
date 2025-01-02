package movement

import (
	"github.com/nosyliam/revolution/pkg/common"
	"time"
)

type Direction int

const (
	Forward   Direction = Direction(common.Forward)
	Backwards Direction = Direction(common.Backward)
	Left      Direction = Direction(common.Left)
	Right     Direction = Direction(common.Right)
)

func walk(direction Direction, distance float64, macro common.Macro, async bool) {
	if len(macro.Stop) > 0 {
		return
	}
	finish := make(chan struct{})
	go func() {
		remaining := distance
		macro.EventBus.KeyDown(macro.Window.PID(), common.Key(direction))
		for remaining > 0 {
			change := macro.BuffDetect.Watch()
			speed := macro.BuffDetect.MoveSpeed()
			start := time.Now()
			select {
			case <-time.After(time.Duration(remaining/speed) * time.Second):
				macro.EventBus.KeyUp(macro.Window.PID(), common.Key(direction))
				if !async {
					finish <- struct{}{}
				}
				return
			case <-macro.Stop:
				macro.EventBus.KeyUp(macro.Window.PID(), common.Key(direction))
				return
			case resume := <-macro.Pause:
				macro.EventBus.KeyUp(macro.Window.PID(), common.Key(direction))
				<-resume
				macro.EventBus.KeyDown(macro.Window.PID(), common.Key(direction))
				remaining -= (time.Now().Sub(start)).Seconds() * speed
			case <-change:
				remaining -= (time.Now().Sub(start)).Seconds() * speed
			}
			close(change)
		}
	}()
	if !async {
		<-finish
	}
}

func Walk(direction Direction, distance float64, macro common.Macro) {
	walk(direction, distance, macro, false)
}

func WalkAsync(direction Direction, distance float64, macro common.Macro) {
	walk(direction, distance, macro, true)
}
