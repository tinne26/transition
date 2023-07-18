package trigger

import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/u16"

var _ Trigger = (*TrigTemplate)(nil)

type TrigTemplate struct {
	area u16.Rect
	// ...
}

func NewTemplate(area u16.Rect) Trigger {
	return &TrigTemplate{
		area: area,
	}
}

func (self *TrigTemplate) Update(playerRect u16.Rect, _ *state.State, _ *audio.Soundscape) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	
	// ...

	return nil, nil
}

func (self *TrigTemplate) OnLevelEnter(_ *state.State) {}
func (self *TrigTemplate) OnLevelExit(_ *state.State) {}
func (self *TrigTemplate) OnDeath(_ *state.State) {}
